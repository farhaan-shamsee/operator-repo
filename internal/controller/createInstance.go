package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	computev1 "github.com/farhaan-shamsee/operator-repo/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// stringOrNil returns nil if s is empty, otherwise returns a pointer to the string
func stringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return aws.String(s)
}

func createEc2Instance(ctx context.Context, ec2Instance *computev1.Ec2instance) (createdInstanceInfo *computev1.CreatedInstanceInfo, err error) {
	l := log.FromContext(ctx) // Use context-aware logger instead of global logger

	l.Info("=== STARTING EC2 INSTANCE CREATION ===",
		"ami", ec2Instance.Spec.AMIId,
		"instanceType", ec2Instance.Spec.InstanceType,
		"region", ec2Instance.Spec.Region)

	// create the client for ec2 instance
	ec2Client, err := awsClient(ec2Instance.Spec.Region)
	if err != nil {
		l.Error(err, "Failed to create AWS client")
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	// create the input for the run instances
	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(ec2Instance.Spec.AMIId),
		InstanceType: ec2types.InstanceType(ec2Instance.Spec.InstanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		KeyName:      stringOrNil(ec2Instance.Spec.KeyPair),
		SubnetId:     stringOrNil(ec2Instance.Spec.Subnet),
		UserData:     stringOrNil(ec2Instance.Spec.UserData),
	}

	// Add security groups if provided
	if len(ec2Instance.Spec.SecurityGroups) > 0 {
		runInput.SecurityGroupIds = ec2Instance.Spec.SecurityGroups
	}

	// Prepare tags - always include a Name tag based on the Kubernetes resource name
	tags := []ec2types.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(ec2Instance.Name), // Use Kubernetes resource name
		},
		{
			Key:   aws.String("ManagedBy"),
			Value: aws.String("ec2instance-operator"),
		},
		{
			Key:   aws.String("Namespace"),
			Value: aws.String(ec2Instance.Namespace),
		},
	}

	// Add user-defined tags from spec
	for key, value := range ec2Instance.Spec.Tags {
		tags = append(tags, ec2types.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	// Add tags to the instance creation request
	if len(tags) > 0 {
		runInput.TagSpecifications = []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeInstance,
				Tags:         tags,
			},
		}
	}

	l.Info("=== CALLING AWS RunInstances API ===")

	// run the instances
	result, err := ec2Client.RunInstances(ctx, runInput)
	if err != nil {
		l.Error(err, "Failed to create EC2 instance")
		return nil, fmt.Errorf("failed to create EC2 instance: %w", err)
	}

	if len(result.Instances) == 0 {
		l.Error(nil, "No instances returned in RunInstanceOutput")
		fmt.Println("No instances returned in RunInstanceOutput")
		return nil, nil
	}

	// Till here, the instance is created and we have
	// Instance ID, private dns and IP, instance type and image id.
	inst := result.Instances[0]
	l.Info("=== EC2 INSTANCE CREATED SUCCESSFULLY ===",
		"InstanceID", *inst.InstanceId,
		"State", inst.State.Name)

	// Now we need to wait for the instance to be running and then get the public ip and dns
	l.Info("=== WAITING FOR INSTANCE TO BE IN 'running' STATE ===")

	runWaiter := ec2.NewInstanceRunningWaiter(ec2Client)
	maxWaitTime := 3 * time.Minute

	err = runWaiter.Wait(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{*inst.InstanceId},
	}, maxWaitTime) //we wait for instance to be in running state

	if err != nil {
		l.Error(err, "Failed to wait for the instance to come to running state")
		return nil, fmt.Errorf("failed to wait for the instance to be in running: %w", err)
	}

	// After creating the instance, we waited and now we describe to
	// 1. Get the public IP and dns as it takes some time for it
	// 2. Getting the state of the instance.
	// We do this so we can send the instance's state to the status of the custom resource. for user to see with kubectl get ec2instances
	l.Info("=== CALLING AWS DescribeInstances API TO GET INSTANCE DETAILS ===")
	describeInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{*inst.InstanceId},
	}

	describeResult, err := ec2Client.DescribeInstances(ctx, describeInput)
	if err != nil {
		l.Error(err, "Failed to describe EC2 instance")
		return nil, fmt.Errorf("failed to describe EC2 instance: %w", err)
	}

	// You get "invalid memory address or nil pointer dereference" here if any of the following are true:
	// - result.Instances is nil or has length 0
	// - Any of the pointer fields (e.g., PublicIpAddress, PrivateIpAddress, etc.) are nil

	// To avoid this, always check for nil and length before dereferencing:

	// Wait for a bit to allow instance fields to be populated

	fmt.Printf("Private IP of the instance: %v", derefString(inst.PrivateIpAddress))
	fmt.Printf("State of the instance: %v", describeResult.Reservations[0].Instances[0].State.Name)
	fmt.Printf("Private DNS of the instance: %v", derefString(inst.PrivateDnsName))
	fmt.Printf("Instance ID of the instance: %v", derefString(inst.InstanceId))
	fmt.Println("Instance Type of the instance: ", inst.InstanceType)
	fmt.Printf("Image ID of the instance: %v", derefString(inst.ImageId))
	fmt.Printf("Key Name of the instance: %v", derefString(inst.KeyName))

	// block until the instance is running
	// blockUntilInstanceRunning(ctx, ec2Instance.Status.InstanceID, ec2Instance)

	// Get the instance details safely (public IP/DNS might be nil for private subnets)
	instance := describeResult.Reservations[0].Instances[0]

	createdInstanceInfo = &computev1.CreatedInstanceInfo{
		InstanceId: *instance.InstanceId,
		State:      string(instance.State.Name),
		PublicIP:   derefString(instance.PublicIpAddress),
		PrivateIP:  derefString(instance.PrivateIpAddress),
		PublicDNS:  derefString(instance.PublicDnsName),
		PrivateDNS: derefString(instance.PrivateDnsName),
	}

	l.Info("=== EC2 INSTANCE CREATION COMPLETED ===",
		"InstanceID", createdInstanceInfo.InstanceId,
		"State", createdInstanceInfo.State,
		"PublicIP", createdInstanceInfo.PublicIP,
	)

	return createdInstanceInfo, nil
}

// derefString is a helper function to safely dereference *string
func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return "<nil>"
}
