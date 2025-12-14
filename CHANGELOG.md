# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-14

### Added

#### Core Features

- **EC2 Instance Management**: Full lifecycle management of AWS EC2 instances via Kubernetes Custom Resources
- **Custom Resource Definition (CRD)**: `Ec2instance` resource with comprehensive spec and status fields
- **Automated Instance Creation**: Deploy EC2 instances by creating Kubernetes resources
- **Automated Instance Deletion**: Cleanup EC2 instances when Kubernetes resources are deleted
- **Finalizer Pattern**: Proper resource cleanup using Kubernetes finalizers to prevent orphaned AWS resources

#### AWS Integration

- **AWS SDK v2 Integration**: Modern AWS SDK for Go v2 implementation
- **Multi-Region Support**: Deploy instances in any AWS region
- **Instance Configuration**: Support for:
  - Instance type selection
  - AMI ID specification
  - VPC subnet assignment
  - Security group attachment
  - SSH key pair configuration
  - User data scripts
  - EBS storage configuration
  - Public IP association
  - Custom tags
- **Tag-Based Naming**: Automatic Name tag creation from Kubernetes resource metadata
- **Instance State Tracking**: Monitor instance state and network details in Kubernetes status

#### Security

- **Secure Credential Management**: AWS credentials from Kubernetes Secrets
- **Helm Secret Integration**: Automatic secret mounting in operator deployment
- **Credential Validation**: Pre-flight validation of AWS credentials
- **Security Hardening**: 
  - Comprehensive `.gitignore` patterns
  - Example `.env.example` template
  - Security best practices documentation
  - No credentials in source code

#### Deployment & Operations

- **Helm Chart**: Production-ready Helm chart for easy deployment
- **RBAC Configuration**: Complete role-based access control setup
- **Health & Readiness Probes**: Kubernetes health checks on ports 8081
- **Metrics Endpoint**: Prometheus metrics on port 8443 (HTTPS with authentication)
- **Leader Election**: Support for high-availability deployments
- **Resource Limits**: Sensible CPU and memory limits configured

#### Error Handling

- **Context-Aware Logging**: Structured logging with request context
- **AWS Error Handling**: Proper error propagation from AWS API
- **Reconciliation Loop**: Robust error handling and retry logic
- **Status Updates**: Instance state reflected in Kubernetes status
- **Deletion Safety**: Checks for instance existence before termination

#### Documentation

- **README**: Comprehensive project documentation with:
  - Quick start guide
  - Architecture overview
  - Security section
  - Development instructions
  - Testing guidelines
- **AWS Credentials Setup Guide**: Step-by-step instructions for secret configuration
- **API Documentation**: Generated API reference for CRD
- **Sample Manifests**: Example EC2 instance definitions for multiple regions

#### Developer Experience

- **Makefile**: Common development tasks automated
- **Code Generation**: Controller-gen integration for deepcopy and CRD generation
- **Testing Framework**: Test suite structure with envtest
- **Regional Examples**: Sample YAML files optimized for ap-south-1, us-east-1, and us-west-2

### Technical Details

#### Dependencies
- Go 1.24.5
- controller-runtime v0.21.0
- AWS SDK for Go v2
- Kubernetes API 0.32.0

#### Supported Configurations
- **Regions**: All AWS regions supported
- **Instance Types**: All EC2 instance types (examples use t3.micro/t2.micro based on region)
- **AMIs**: Any public or private AMI
- **Networking**: VPC, subnet, security group, public IP customization
- **Storage**: EBS volume configuration with size, type, and IOPS settings

### Fixed
- Import path corrections for Go module structure
- Deepcopy generation issues with CRD types
- Infinite reconciliation loop from duplicate finalizers
- Empty string parameters causing AWS API errors
- Region-specific instance type availability (t3.micro for ap-south-1)

### Changed
- Migrated to AWS SDK v2 from v1
- Updated to controller-runtime v0.21.0
- Enhanced error messages with context
- Improved status field structure for better observability

### Security
- Added comprehensive `.gitignore` for credentials and secrets
- Created `.env.example` template for local development
- Documented security best practices in README
- Implemented secret-based credential injection via Helm

## [Unreleased]

### Planned
- E2E test implementation
- Multiple instance management optimization
- AWS IAM role support (IRSA for EKS)
- Instance update/modification support
- CloudWatch integration for metrics
- SNS notifications for instance events
- Cost tracking and reporting

---

## Release Links

- [1.0.0] - Initial Release - 2025-12-14

## Upgrade Notes

### Upgrading to 1.0.0
This is the initial release. No upgrade path needed.

## Breaking Changes

None - this is the initial release.

## Contributors

- Farhaan Shamsee (@farhaan-shamsee)

---

**Note**: For detailed installation and usage instructions, see [README.md](README.md).
