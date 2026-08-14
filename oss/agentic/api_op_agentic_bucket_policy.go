package agentic

import (
	"context"
	"io"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// --- PutAgenticBucketPolicy ---

// PutAgenticBucketPolicyRequest is the request for the PutAgenticBucketPolicy operation.
type PutAgenticBucketPolicyRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	// The policy of the agentic bucket, in JSON format.
	Body io.Reader `input:"body,nop,required"`

	oss.RequestCommon
}

// PutAgenticBucketPolicyResult is the result for the PutAgenticBucketPolicy operation.
type PutAgenticBucketPolicyResult struct {
	oss.ResultCommon
}

// PutAgenticBucketPolicy Configures the policy of an agentic bucket.
func (c *AgenticBucketClient) PutAgenticBucketPolicy(ctx context.Context, request *PutAgenticBucketPolicyRequest, optFns ...func(*oss.Options)) (*PutAgenticBucketPolicyResult, error) {
	var err error
	if request == nil {
		request = &PutAgenticBucketPolicyRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutAgenticBucketPolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: "application/json",
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"policy":        "",
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

	result := &PutAgenticBucketPolicyResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- GetAgenticBucketPolicy ---

// GetAgenticBucketPolicyRequest is the request for the GetAgenticBucketPolicy operation.
type GetAgenticBucketPolicyRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// GetAgenticBucketPolicyResult is the result for the GetAgenticBucketPolicy operation.
type GetAgenticBucketPolicyResult struct {
	// The policy of the agentic bucket, in JSON format.
	Body string

	oss.ResultCommon
}

// GetAgenticBucketPolicy Queries the policy of an agentic bucket.
func (c *AgenticBucketClient) GetAgenticBucketPolicy(ctx context.Context, request *GetAgenticBucketPolicyRequest, optFns ...func(*oss.Options)) (*GetAgenticBucketPolicyResult, error) {
	var err error
	if request == nil {
		request = &GetAgenticBucketPolicyRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketPolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"policy":        "",
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

	body, err := io.ReadAll(output.Body)
	defer output.Body.Close()
	if err != nil {
		return nil, err
	}

	result := &GetAgenticBucketPolicyResult{
		Body: string(body),
	}
	if err = c.clientImpl.UnmarshalOutput(result, output); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- DeleteAgenticBucketPolicy ---

// DeleteAgenticBucketPolicyRequest is the request for the DeleteAgenticBucketPolicy operation.
type DeleteAgenticBucketPolicyRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// DeleteAgenticBucketPolicyResult is the result for the DeleteAgenticBucketPolicy operation.
type DeleteAgenticBucketPolicyResult struct {
	oss.ResultCommon
}

// DeleteAgenticBucketPolicy Deletes the policy of an agentic bucket.
func (c *AgenticBucketClient) DeleteAgenticBucketPolicy(ctx context.Context, request *DeleteAgenticBucketPolicyRequest, optFns ...func(*oss.Options)) (*DeleteAgenticBucketPolicyResult, error) {
	var err error
	if request == nil {
		request = &DeleteAgenticBucketPolicyRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteAgenticBucketPolicy",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"policy":        "",
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

	result := &DeleteAgenticBucketPolicyResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
