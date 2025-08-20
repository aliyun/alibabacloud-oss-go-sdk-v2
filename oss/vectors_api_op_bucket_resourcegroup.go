package oss

import (
	"context"
)

// PutBucketResourceGroup Modifies the ID of the resource group to which a vectors bucket belongs.
func (c *VectorsClient) PutBucketResourceGroup(ctx context.Context, request *PutBucketResourceGroupRequest, optFns ...func(*Options)) (*PutBucketResourceGroupResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	return c.client.PutBucketResourceGroup(ctx, request, optFns...)
}

// GetBucketResourceGroup Queries the ID of the resource group to which a vectors bucket belongs.
func (c *VectorsClient) GetBucketResourceGroup(ctx context.Context, request *GetBucketResourceGroupRequest, optFns ...func(*Options)) (*GetBucketResourceGroupResult, error) {
	return c.client.GetBucketResourceGroup(ctx, request, optFns...)
}