# Sample Configurations

This directory contains sample YAML files demonstrating different configurations for the operator's Custom Resources.

## EC2 Instance Samples

### `compute_v1_ec2instance.yaml`
Full-featured EC2 instance with all configuration options including:
- Custom instance type, AMI, region
- VPC subnet and security groups
- SSH key pair
- User data script
- EBS storage configuration
- Public IP assignment
- Custom tags

### `compute_v1_ec2instance_minimal.yaml`
Minimal EC2 instance configuration with only required fields.

## S3 Bucket Samples

### `compute_v1_s3bucket.yaml`
Standard S3 bucket configuration with:
- Private ACL
- STANDARD storage class
- Region specification

### `compute_v1_s3bucket_minimal.yaml`
Minimal S3 bucket with only required fields:
- Bucket name
- Region

### `compute_v1_s3bucket_versioned.yaml`
S3 bucket with versioning enabled:
- Versioning: Enabled
- Private ACL
- STANDARD storage class

### `compute_v1_s3bucket_glacier.yaml`
S3 bucket configured for archival storage:
- Storage class: GLACIER
- Versioning: Suspended
- Private ACL
- Different region (us-west-2)

### `compute_v1_s3bucket_public_read.yaml`
S3 bucket with public read access:
- ACL: public-read
- ⚠️ **Warning**: Requires disabling AWS Block Public Access settings

### `compute_v1_s3bucket_multiregion.yaml`
S3 bucket in EU region demonstrating multi-region support:
- Region: eu-central-1
- Versioning enabled
- Private ACL

## Usage

### Apply a single sample:
```bash
kubectl apply -f config/samples/compute_v1_s3bucket_minimal.yaml
```

### Apply all samples:
```bash
kubectl apply -k config/samples/
```

### Delete a sample:
```bash
kubectl delete -f config/samples/compute_v1_s3bucket_minimal.yaml
```

### Delete all samples:
```bash
kubectl delete -k config/samples/
```

## Common ACL Values

- `private` - Only bucket owner has access (default, recommended)
- `public-read` - Anyone can read objects (requires disabling Block Public Access)
- `public-read-write` - Anyone can read/write (not recommended)
- `authenticated-read` - Any authenticated AWS user can read

## Storage Classes

- `STANDARD` - General purpose (default)
- `STANDARD_IA` - Infrequent access
- `ONEZONE_IA` - Single AZ, infrequent access
- `GLACIER` - Archive storage, slower retrieval
- `GLACIER_IR` - Glacier Instant Retrieval
- `DEEP_ARCHIVE` - Lowest cost archive
- `INTELLIGENT_TIERING` - Automatic tier optimization

## Versioning Options

- `Enabled` - Keep all versions of objects
- `Suspended` - Stop creating new versions (existing versions remain)
- Empty/omitted - Versioning not configured

## Notes

⚠️ **Bucket Names**: Must be globally unique across all AWS accounts

⚠️ **Public ACLs**: AWS blocks public access by default for security

⚠️ **Region Constraints**: Non us-east-1 regions require LocationConstraint (handled automatically)

⚠️ **Deletion**: S3 buckets must be empty before deletion (handled automatically by the operator)
