package agentic

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// --- PutAgenticBucketPublicAccessBlock ---

// PutAgenticBucketPublicAccessBlockRequest is the request for the PutAgenticBucketPublicAccessBlock operation.
type PutAgenticBucketPublicAccessBlockRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	// The Block Public Access configuration of the agentic bucket.
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `input:"body,PublicAccessBlockConfiguration,xml,required"`

	oss.RequestCommon
}

// PutAgenticBucketPublicAccessBlockResult is the result for the PutAgenticBucketPublicAccessBlock operation.
type PutAgenticBucketPublicAccessBlockResult struct {
	oss.ResultCommon
}

// PutAgenticBucketPublicAccessBlock Configures the Block Public Access of an agentic bucket.
func (c *AgenticBucketClient) PutAgenticBucketPublicAccessBlock(ctx context.Context, request *PutAgenticBucketPublicAccessBlockRequest, optFns ...func(*oss.Options)) (*PutAgenticBucketPublicAccessBlockResult, error) {
	var err error
	if request == nil {
		request = &PutAgenticBucketPublicAccessBlockRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutAgenticBucketPublicAccessBlock",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket":     "",
			"publicAccessBlock": "",
		},
		Bucket: request.Bucket,
	}

	if err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutAgenticBucketPublicAccessBlockResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- GetAgenticBucketPublicAccessBlock ---

// GetAgenticBucketPublicAccessBlockRequest is the request for the GetAgenticBucketPublicAccessBlock operation.
type GetAgenticBucketPublicAccessBlockRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// GetAgenticBucketPublicAccessBlockResult is the result for the GetAgenticBucketPublicAccessBlock operation.
type GetAgenticBucketPublicAccessBlockResult struct {
	// The Block Public Access configuration of the agentic bucket.
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `output:"body,PublicAccessBlockConfiguration,xml"`

	oss.ResultCommon
}

// GetAgenticBucketPublicAccessBlock Queries the Block Public Access configuration of an agentic bucket.
func (c *AgenticBucketClient) GetAgenticBucketPublicAccessBlock(ctx context.Context, request *GetAgenticBucketPublicAccessBlockRequest, optFns ...func(*oss.Options)) (*GetAgenticBucketPublicAccessBlockResult, error) {
	var err error
	if request == nil {
		request = &GetAgenticBucketPublicAccessBlockRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketPublicAccessBlock",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket":     "",
			"publicAccessBlock": "",
		},
		Bucket: request.Bucket,
	}

	if err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetAgenticBucketPublicAccessBlockResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- DeleteAgenticBucketPublicAccessBlock ---

// DeleteAgenticBucketPublicAccessBlockRequest is the request for the DeleteAgenticBucketPublicAccessBlock operation.
type DeleteAgenticBucketPublicAccessBlockRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// DeleteAgenticBucketPublicAccessBlockResult is the result for the DeleteAgenticBucketPublicAccessBlock operation.
type DeleteAgenticBucketPublicAccessBlockResult struct {
	oss.ResultCommon
}

// DeleteAgenticBucketPublicAccessBlock Deletes the Block Public Access configuration of an agentic bucket.
func (c *AgenticBucketClient) DeleteAgenticBucketPublicAccessBlock(ctx context.Context, request *DeleteAgenticBucketPublicAccessBlockRequest, optFns ...func(*oss.Options)) (*DeleteAgenticBucketPublicAccessBlockResult, error) {
	var err error
	if request == nil {
		request = &DeleteAgenticBucketPublicAccessBlockRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteAgenticBucketPublicAccessBlock",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket":     "",
			"publicAccessBlock": "",
		},
		Bucket: request.Bucket,
	}

	if err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteAgenticBucketPublicAccessBlockResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
