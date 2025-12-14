# operator-repo

[![Release](https://img.shields.io/github/v/release/farhaan-shamsee/operator-repo)](https://github.com/farhaan-shamsee/operator-repo/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/farhaan-shamsee/operator-repo)](https://goreportcard.com/report/github.com/farhaan-shamsee/operator-repo)

EC2 Instance Kubernetes Operator - Manage AWS EC2 instances declaratively using Kubernetes Custom Resources.

## Description

This operator allows you to manage AWS EC2 instances as Kubernetes resources. It provides a declarative way to create, update, and delete EC2 instances using familiar Kubernetes tools and patterns.

**Features:**
- Create and manage EC2 instances via Kubernetes Custom Resources
- Automatic cleanup with finalizers
- Support for VPC, security groups, SSH keys, and user data
- Tag-based naming and organization
- Parallel instance creation and deletion

## Getting Started

### Prerequisites
- go version v1.24.0+
- docker version 17.03+
- kubectl version v1.11.3+
- Access to a Kubernetes v1.11.3+ cluster
- **AWS Account with EC2 permissions**
- **AWS credentials (Access Key ID and Secret Access Key)**

### AWS Credentials Setup

**IMPORTANT:** Never commit AWS credentials to git!

1. Copy the example environment file:
   ```sh
   cp .env.example .env
   ```

2. Edit `.env` and add your AWS credentials:
   ```sh
   AWS_ACCESS_KEY_ID=your-actual-access-key
   AWS_SECRET_ACCESS_KEY=your-actual-secret-key
   AWS_DEFAULT_REGION=ap-south-1
   ```

3. Load the credentials when running locally:
   ```sh
   source .env
   make run
   ```

**For production deployments**, use Kubernetes Secrets:
```sh
kubectl create secret generic aws-credentials \
  --from-literal=AWS_ACCESS_KEY_ID=your-key \
  --from-literal=AWS_SECRET_ACCESS_KEY=your-secret
```

## Installation

### Quick Install with Helm (Recommended)

1. Create AWS credentials secret:
```sh
kubectl create secret generic aws-credentials \
  --from-literal=AWS_ACCESS_KEY_ID=your-key \
  --from-literal=AWS_SECRET_ACCESS_KEY=your-secret \
  --namespace=default
```

2. Install the operator:
```sh
helm install ec2-operator ./dist/chart --namespace=default
```

3. Create your first EC2 instance:
```sh
kubectl apply -f config/samples/compute_v1_ec2instance_minimal.yaml
```

4. Check the status:
```sh
kubectl get ec2instances
kubectl describe ec2instance ec2instance-sample
```

### Install from GitHub Release

```sh
# Download the latest release
wget https://github.com/farhaan-shamsee/operator-repo/releases/download/v1.0.0/operator-repo-v1.0.0.tar.gz

# Extract
tar -xzf operator-repo-v1.0.0.tar.gz

# Install
cd operator-repo-v1.0.0
helm install ec2-operator ./dist/chart --namespace=default
```

### To Deploy on the cluster (Development)
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/operator-repo:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/operator-repo:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/operator-repo:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/operator-repo/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

