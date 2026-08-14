package agentic

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// --- PutAgenticBucketEncryption ---

// PutAgenticBucketEncryptionRequest is the request for the PutAgenticBucketEncryption operation.
type PutAgenticBucketEncryptionRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	// The server-side encryption rule of the agentic bucket.
	ServerSideEncryptionRule *ServerSideEncryptionRule `input:"body,ServerSideEncryptionRule,xml,required"`

	oss.RequestCommon
}

// PutAgenticBucketEncryptionResult is the result for the PutAgenticBucketEncryption operation.
type PutAgenticBucketEncryptionResult struct {
	oss.ResultCommon
}

// PutAgenticBucketEncryption Configures the server-side encryption rule of an agentic bucket.
func (c *AgenticBucketClient) PutAgenticBucketEncryption(ctx context.Context, request *PutAgenticBucketEncryptionRequest, optFns ...func(*oss.Options)) (*PutAgenticBucketEncryptionResult, error) {
	var err error
	if request == nil {
		request = &PutAgenticBucketEncryptionRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutAgenticBucketEncryption",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"encryption":    "",
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

	result := &PutAgenticBucketEncryptionResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- GetAgenticBucketEncryption ---

// GetAgenticBucketEncryptionRequest is the request for the GetAgenticBucketEncryption operation.
type GetAgenticBucketEncryptionRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// GetAgenticBucketEncryptionResult is the result for the GetAgenticBucketEncryption operation.
type GetAgenticBucketEncryptionResult struct {
	// The server-side encryption rule of the agentic bucket.
	ServerSideEncryptionRule *ServerSideEncryptionRule `output:"body,ServerSideEncryptionRule,xml"`

	oss.ResultCommon
}

// GetAgenticBucketEncryption Queries the server-side encryption rule of an agentic bucket.
func (c *AgenticBucketClient) GetAgenticBucketEncryption(ctx context.Context, request *GetAgenticBucketEncryptionRequest, optFns ...func(*oss.Options)) (*GetAgenticBucketEncryptionResult, error) {
	var err error
	if request == nil {
		request = &GetAgenticBucketEncryptionRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketEncryption",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"encryption":    "",
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

	result := &GetAgenticBucketEncryptionResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- DeleteAgenticBucketEncryption ---

// DeleteAgenticBucketEncryptionRequest is the request for the DeleteAgenticBucketEncryption operation.
type DeleteAgenticBucketEncryptionRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// DeleteAgenticBucketEncryptionResult is the result for the DeleteAgenticBucketEncryption operation.
type DeleteAgenticBucketEncryptionResult struct {
	oss.ResultCommon
}

// DeleteAgenticBucketEncryption Deletes the server-side encryption rule of an agentic bucket.
func (c *AgenticBucketClient) DeleteAgenticBucketEncryption(ctx context.Context, request *DeleteAgenticBucketEncryptionRequest, optFns ...func(*oss.Options)) (*DeleteAgenticBucketEncryptionResult, error) {
	var err error
	if request == nil {
		request = &DeleteAgenticBucketEncryptionRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteAgenticBucketEncryption",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"encryption":    "",
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

	result := &DeleteAgenticBucketEncryptionResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
