package tables

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutTableBucketPolicyRequest = oss.PutBucketPolicyRequest
type PutTableBucketPolicyResult = oss.PutBucketPolicyResult

// PutTableBucketPolicy Configures a policy for a table bucket.
func (c *TablesClient) PutTableBucketPolicy(ctx context.Context, request *PutTableBucketPolicyRequest, optFns ...func(*oss.Options)) (*PutTableBucketPolicyResult, error) {
	return c.clientImpl.PutBucketPolicy(ctx, request, optFns...)
}

type GetTableBucketPolicyRequest = oss.GetBucketPolicyRequest
type GetTableBucketPolicyResult = oss.GetBucketPolicyResult

// GetTableBucketPolicy Queries the policies configured for a table bucket.
func (c *TablesClient) GetTableBucketPolicy(ctx context.Context, request *GetTableBucketPolicyRequest, optFns ...func(*oss.Options)) (*GetTableBucketPolicyResult, error) {
	return c.clientImpl.GetBucketPolicy(ctx, request, optFns...)
}

type DeleteTableBucketPolicyRequest = oss.DeleteBucketPolicyRequest
type DeleteTableBucketPolicyResult = oss.DeleteBucketPolicyResult

// DeleteTableBucketPolicy Deletes a policy for a table bucket.
func (c *TablesClient) DeleteTableBucketPolicy(ctx context.Context, request *DeleteTableBucketPolicyRequest, optFns ...func(*oss.Options)) (*DeleteTableBucketPolicyResult, error) {
	return c.clientImpl.DeleteBucketPolicy(ctx, request, optFns...)
}
