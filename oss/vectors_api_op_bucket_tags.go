package oss

import (
	"context"
)

// PutBucketTags Adds tags to or modifies the existing tags of a vector bucket.
func (c *VectorsClient) PutBucketTags(ctx context.Context, request *PutBucketTagsRequest, optFns ...func(*Options)) (*PutBucketTagsResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	return c.client.PutBucketTags(ctx, request, optFns...)
}

// GetBucketTags Queries the tags of a vector bucket.
func (c *VectorsClient) GetBucketTags(ctx context.Context, request *GetBucketTagsRequest, optFns ...func(*Options)) (*GetBucketTagsResult, error) {
	return c.client.GetBucketTags(ctx, request, optFns...)
}

// DeleteBucketTags Deletes tags configured for a vector bucket.
func (c *VectorsClient) DeleteBucketTags(ctx context.Context, request *DeleteBucketTagsRequest, optFns ...func(*Options)) (*DeleteBucketTagsResult, error) {
	return c.client.DeleteBucketTags(ctx, request, optFns...)
}
