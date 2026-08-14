package agentic

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type ServerSideEncryptionRule = oss.ServerSideEncryptionRule
type ApplyServerSideEncryptionByDefault = oss.ApplyServerSideEncryptionByDefault
type VersioningConfiguration = oss.VersioningConfiguration
type PublicAccessBlockConfiguration = oss.PublicAccessBlockConfiguration
type Owner = oss.Owner
type BucketACLType = oss.BucketACLType
type StorageClassType = oss.StorageClassType
type DataRedundancyType = oss.DataRedundancyType

func (c *AgenticBucketClient) InvokeOperation(ctx context.Context, input *oss.OperationInput, optFns ...func(*oss.Options)) (*oss.OperationOutput, error) {
	return c.clientImpl.InvokeOperation(ctx, input, optFns...)
}
