# Release Notes - v1.2.0

**Release Date**: January 3, 2026

## S3 Bucket Management Support

This release adds comprehensive S3 bucket management capabilities to the operator, expanding beyond EC2 instances to support AWS storage resources.

## What's New

### S3 Bucket Controller

- **Declarative S3 Management**: Manage S3 buckets using Kubernetes Custom Resources
- **Full Lifecycle Support**: Automated creation, monitoring, and deletion of S3 buckets
- **Automatic Cleanup**: Buckets are emptied before deletion to comply with AWS requirements

### Key Features

#### S3 Bucket Configuration

- Bucket name specification
- Region selection
- ACL configuration (private, public-read, etc.)
- Versioning control (Enabled/Suspended)
- Storage class selection (STANDARD, GLACIER, etc.)
- LocationConstraint handling for non-us-east-1 regions

#### S3 Status Tracking

The operator tracks and reports:
- Bucket ARN (Amazon Resource Name)
- Location (AWS response path)
- Creation status
- Last sync time (RFC3339 format)

#### Enhanced AWS Client

- **Generic AWS Config**: Refactored to support multiple AWS services
- **Reusable Configuration**: Single AWS config function for all services
- **Service-Specific Clients**: EC2, S3, and extensible for future services

## Improvements

### Code Quality

- Improved AWS client architecture with generic configuration
- Better separation of concerns between services
- Consistent error handling across controllers
- Proper timestamp formatting using RFC3339

### Bug Fixes

- Fixed deletion flow to prevent status updates on deleted resources
- Corrected S3 bucket ARN construction (works for general-purpose buckets)
- Added proper return statement after finalizer removal

## API Changes

### New Custom Resource Definition

```yaml
apiVersion: compute.cloud.com/v1
kind: S3Bucket
```

### New Status Fields

- `bucketARN`: Full ARN of the S3 bucket
- `location`: AWS location response
- `created`: Boolean indicating creation status
- `lastSyncTime`: Last reconciliation timestamp

## Security

Same security posture as v1.0.0:
- AWS credentials via environment variables
- No hardcoded secrets
- RBAC enabled
- Kubernetes finalizers for safe deletion

## Sample Resources

New sample configuration available:
- `config/samples/compute_v1_s3bucket.yaml`
- `config/samples/compute_v1_s3bucket_minimal.yaml`
- `config/samples/compute_v1_s3bucket_versioned.yaml`
- `config/samples/compute_v1_s3bucket_glacier.yaml`
- `config/samples/compute_v1_s3bucket_public_read.yaml`
- `config/samples/compute_v1_s3bucket_multiregion.yaml`

## Architecture

### Controllers

- **EC2 Instance Controller**: Unchanged, stable from v1.0.0
- **S3 Bucket Controller**: New in v1.2.0

### Shared Components

- Unified AWS configuration (`getAWSConfig`)
- Service-specific client initialization
- Consistent reconciliation patterns

## Upgrade Notes

### Breaking Changes

None - fully backward compatible with v1.0.0

### Migration Guide

1. Update CRDs: `make manifests && make install`
2. Update operator image to v1.2.0
3. Apply S3Bucket resources as needed

## Known Limitations

- S3 buckets must be emptied before deletion (handled automatically)
- Bucket names must be globally unique (AWS constraint)
- Directory buckets (S3 Express One Zone) not yet supported
- Versioning and lifecycle policies require manual configuration post-creation

## Future Enhancements

Planned for future releases:
- S3 bucket policy management
- Lifecycle rule configuration
- CORS configuration
- Website hosting settings
- Replication configuration

## Supported Resources

- EC2 Instances (v1.0.0)
- S3 Buckets (v1.2.0)
- Additional AWS resources coming soon

## Credits

Thanks to the Kubernetes operator-sdk and AWS SDK teams for their excellent tools and documentation.

---

**Full Changelog**: https://github.com/farhaan-shamsee/operator-repo/compare/v1.0.0...v1.2.0
