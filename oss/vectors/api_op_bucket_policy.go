package vectors

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutBucketPolicyRequest = oss.PutBucketPolicyRequest
type PutBucketPolicyResult = oss.PutBucketPolicyResult

// PutBucketPolicy Configures a policy for a vector bucket.
func (c *VectorsClient) PutBucketPolicy(ctx context.Context, request *PutBucketPolicyRequest, optFns ...func(*oss.Options)) (*PutBucketPolicyResult, error) {
	return c.clientImpl.PutBucketPolicy(ctx, request, optFns...)
}

type GetBucketPolicyRequest = oss.GetBucketPolicyRequest
type GetBucketPolicyResult = oss.GetBucketPolicyResult

// GetBucketPolicy Queries the policies configured for a vector bucket.
func (c *VectorsClient) GetBucketPolicy(ctx context.Context, request *GetBucketPolicyRequest, optFns ...func(*oss.Options)) (*GetBucketPolicyResult, error) {
	return c.clientImpl.GetBucketPolicy(ctx, request, optFns...)
}

type DeleteBucketPolicyRequest = oss.DeleteBucketPolicyRequest
type DeleteBucketPolicyResult = oss.DeleteBucketPolicyResult

// DeleteBucketPolicy Deletes a policy for a vector bucket.
func (c *VectorsClient) DeleteBucketPolicy(ctx context.Context, request *DeleteBucketPolicyRequest, optFns ...func(*oss.Options)) (*DeleteBucketPolicyResult, error) {
	return c.clientImpl.DeleteBucketPolicy(ctx, request, optFns...)
}
