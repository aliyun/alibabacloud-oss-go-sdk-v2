package vectors

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type BucketResourceGroupConfiguration = oss.BucketResourceGroupConfiguration
type PutBucketResourceGroupRequest = oss.PutBucketResourceGroupRequest
type PutBucketResourceGroupResult = oss.PutBucketResourceGroupResult

// PutBucketResourceGroup Modifies the ID of the resource group to which a vectors bucket belongs.
func (c *VectorsClient) PutBucketResourceGroup(ctx context.Context, request *PutBucketResourceGroupRequest, optFns ...func(*oss.Options)) (*PutBucketResourceGroupResult, error) {
	var err error
	if request == nil {
		request = &PutBucketResourceGroupRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutBucketResourceGroup",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &PutBucketResourceGroupResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type GetBucketResourceGroupRequest = oss.GetBucketResourceGroupRequest
type GetBucketResourceGroupResult = oss.GetBucketResourceGroupResult

// GetBucketResourceGroup Queries the ID of the resource group to which a vectors bucket belongs.
func (c *VectorsClient) GetBucketResourceGroup(ctx context.Context, request *GetBucketResourceGroupRequest, optFns ...func(*oss.Options)) (*GetBucketResourceGroupResult, error) {
	var err error
	if request == nil {
		request = &GetBucketResourceGroupRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetBucketResourceGroup",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}

	if err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &GetBucketResourceGroupResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
