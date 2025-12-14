# Release Notes - v1.0.0

**Release Date**: December 14, 2025

## 🎉 Initial Release

This is the first stable release of the **EC2 Kubernetes Operator** - a production-ready Kubernetes operator for managing AWS EC2 instances as native Kubernetes resources.

## 🚀 What's New

### Core Functionality

- **Declarative EC2 Management**: Manage EC2 instances using Kubernetes Custom Resources
- **Full Lifecycle Support**: Automated creation, monitoring, and deletion of EC2 instances
- **AWS SDK v2**: Built with the latest AWS SDK for Go v2
- **Multi-Region**: Deploy instances across any AWS region

### Key Features

#### 🔧 Instance Configuration
- Instance type selection (t2.micro, t3.micro, etc.)
- AMI ID specification
- VPC subnet assignment
- Security group attachment
- SSH key pair configuration
- User data scripts for initialization
- EBS storage configuration (size, type, IOPS)
- Public IP association
- Custom resource tags

#### 📊 Status Tracking
The operator tracks and reports:
- Instance ID
- Instance state (pending, running, stopping, terminated)
- Public IP address
- Private IP address
- Public DNS name
- Private DNS name

#### 🔒 Security
- **Secure Credentials**: AWS credentials stored in Kubernetes Secrets
- **No Hardcoded Secrets**: Zero credentials in source code
- **RBAC Enabled**: Full role-based access control
- **Validation**: Pre-flight credential validation
- **Security Documentation**: Best practices guide included

#### 📦 Easy Deployment
- **Helm Chart**: Install with a single command
- **Metrics**: Prometheus-compatible metrics endpoint
- **Health Checks**: Kubernetes-native readiness and liveness probes
- **Leader Election**: High-availability support

## 📋 Installation

### Prerequisites
- Kubernetes cluster (1.24+)
- kubectl configured
- Helm 3.x
- AWS account with EC2 permissions

### Quick Start

1. **Create AWS credentials secret**:
```bash
kubectl create secret generic aws-credentials \
  --from-literal=AWS_ACCESS_KEY_ID='your-key' \
  --from-literal=AWS_SECRET_ACCESS_KEY='your-secret' \
  --namespace=default
```

2. **Install the operator**:
```bash
helm install ec2-operator ./dist/chart --namespace=default
```

3. **Create an EC2 instance**:
```bash
kubectl apply -f config/samples/compute_v1_ec2instance_minimal.yaml
```

4. **Check the instance**:
```bash
kubectl get ec2instances
kubectl describe ec2instance ec2instance-sample
```

## 📖 Usage Example

Create an EC2 instance with this simple manifest:

```yaml
apiVersion: compute.cloud.com/v1
kind: Ec2instance
metadata:
  name: my-web-server
spec:
  instanceType: t3.micro
  amiId: ami-0c2af51e265bd5e0e
  region: ap-south-1
  keyPair: my-key-pair
  subnet: subnet-xxxxxxxxx
  securityGroups:
    - sg-xxxxxxxxx
  associatePublicIP: true
  tags:
    Environment: production
    Application: web-server
```

Apply it:
```bash
kubectl apply -f my-instance.yaml
```

The operator will:
1. Validate your AWS credentials
2. Create the EC2 instance in AWS
3. Wait for it to reach "running" state
4. Update the Kubernetes resource status with instance details

Delete it:
```bash
kubectl delete ec2instance my-web-server
```

The operator will:
1. Terminate the EC2 instance in AWS
2. Wait for termination to complete
3. Remove the Kubernetes resource

## 🛠️ Configuration Options

### Spec Fields
- `instanceType`: EC2 instance type (required)
- `amiId`: Amazon Machine Image ID (required)
- `region`: AWS region (required)
- `keyPair`: SSH key pair name (optional)
- `subnet`: VPC subnet ID (optional)
- `securityGroups`: List of security group IDs (optional)
- `tags`: Custom tags map (optional)
- `userData`: Base64-encoded user data script (optional)
- `storage`: EBS configuration (optional)
- `associatePublicIP`: Assign public IP (optional, default: false)

### Storage Configuration
```yaml
storage:
  volumeSize: 30
  volumeType: gp3
  iops: 3000
  throughput: 125
  deleteOnTermination: true
```

## 📚 Documentation

- **README.md**: Project overview and setup
- **CHANGELOG.md**: Complete version history
- **AWS Credentials Setup**: Step-by-step credential configuration
- **API Reference**: Generated CRD documentation

## 🔗 Sample Manifests

Included samples for different regions:
- `compute_v1_ec2instance.yaml`: Full-featured example
- `compute_v1_ec2instance_minimal.yaml`: Minimal configuration
- Examples optimized for ap-south-1, us-east-1, us-west-2

## 🧪 Testing

Run the test suite:
```bash
make test
```

Build and deploy locally:
```bash
# Build the operator
make docker-build IMG=my-registry/operator-repo:v1.0.0

# Push to registry
make docker-push IMG=my-registry/operator-repo:v1.0.0

# Deploy
make deploy IMG=my-registry/operator-repo:v1.0.0
```

## 📊 Metrics

The operator exposes Prometheus metrics on port 8443 (HTTPS with authentication):

```bash
# Create a token
TOKEN=$(kubectl create token operator-repo-controller-manager -n default)

# Access metrics
kubectl port-forward deployment/operator-repo-controller-manager 8443:8443
curl -k -H "Authorization: Bearer $TOKEN" https://localhost:8443/metrics
```

Metrics include:
- Reconciliation counts and duration
- Work queue depth
- Go runtime metrics
- API client metrics

## 🐛 Known Issues

None at this time.

## 🔮 Future Roadmap

- Instance update/modification support
- AWS IAM roles for service accounts (IRSA)
- CloudWatch integration
- SNS event notifications
- Cost tracking and reporting
- Multi-instance group management
- Auto-scaling integration

## 📦 Artifacts

- **Docker Image**: `docker.io/farhaanshamsee/kubernetes-operator:v1.0.0`
- **Helm Chart**: `dist/chart/` (version 0.1.0)
- **CRD**: `ec2instances.compute.cloud.com`

## 🤝 Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📄 License

Apache License 2.0

## 👤 Author

**Farhaan Shamsee** ([@farhaan-shamsee](https://github.com/farhaan-shamsee))

## 🙏 Acknowledgments

Built with:
- [Kubebuilder](https://kubebuilder.io/)
- [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)
- [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2)

---

**Full Changelog**: https://github.com/farhaan-shamsee/operator-repo/blob/main/CHANGELOG.md
