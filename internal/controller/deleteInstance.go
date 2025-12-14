package controller

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	computev1 "github.com/farhaan-shamsee/operator-repo/api/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func deleteEc2Instance(ctx context.Context, ec2Instance *computev1.Ec2instance) (bool, error) {
	l := logf.FromContext(ctx)

	l.Info("Deleting EC2 instance", "instanceID", ec2Instance.Status.InstanceID)

	ec2Client, err := awsClient(ec2Instance.Spec.Region)
	if err != nil {
		l.Error(err, "Failed to create AWS client")
		return false, err
	}

	terminateResult, err := ec2Client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{ec2Instance.Status.InstanceID},
	})

	if err != nil {
		l.Error(err, "Failed to terminate instance")
		return false, err
	}

	l.Info("Instance termination instantiated",
		"instanceId", ec2Instance.Status.InstanceID,
		"currentState", terminateResult.TerminatingInstances[0].CurrentState.Name)

	waiter := ec2.NewInstanceTerminatedWaiter(ec2Client)
	maxWaitTime := 5 * time.Minute
	waitParams := &ec2.DescribeInstancesInput{
		InstanceIds: []string{ec2Instance.Status.InstanceID},
	}

	l.Info("Waiting for the instance to be terminated",
		"InstanceId", ec2Instance.Status.InstanceID,
		"maxWaitTime", maxWaitTime)

	err = waiter.Wait(ctx, waitParams, maxWaitTime)
	if err != nil {
		l.Error(err, "Failed to wait for instance termination",
			"instanceId", ec2Instance.Status.InstanceID)
		return false, err
	}

	l.Info("EC2 Instance successfully terminated", "instanceID", ec2Instance.Status.InstanceID)
	return true, nil
}
