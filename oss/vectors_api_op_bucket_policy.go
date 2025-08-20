package oss

import (
	"context"
)

// PutBucketPolicy Configures a policy for a vector bucket.
func (c *VectorsClient) PutBucketPolicy(ctx context.Context, request *PutBucketPolicyRequest, optFns ...func(*Options)) (*PutBucketPolicyResult, error) {
	return c.client.PutBucketPolicy(ctx, request, optFns...)
}

// GetBucketPolicy Queries the policies configured for a vector bucket.
func (c *VectorsClient) GetBucketPolicy(ctx context.Context, request *GetBucketPolicyRequest, optFns ...func(*Options)) (*GetBucketPolicyResult, error) {
	return c.client.GetBucketPolicy(ctx, request, optFns...)
}

// DeleteBucketPolicy Deletes a policy for a vector bucket.
func (c *VectorsClient) DeleteBucketPolicy(ctx context.Context, request *DeleteBucketPolicyRequest, optFns ...func(*Options)) (*DeleteBucketPolicyResult, error) {
	return c.client.DeleteBucketPolicy(ctx, request, optFns...)
}