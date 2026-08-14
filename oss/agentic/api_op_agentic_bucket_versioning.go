package agentic

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// --- PutAgenticBucketVersioning ---

// PutAgenticBucketVersioningRequest is the request for the PutAgenticBucketVersioning operation.
type PutAgenticBucketVersioningRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	// The versioning configuration of the agentic bucket.
	VersioningConfiguration *VersioningConfiguration `input:"body,VersioningConfiguration,xml,required"`

	oss.RequestCommon
}

// PutAgenticBucketVersioningResult is the result for the PutAgenticBucketVersioning operation.
type PutAgenticBucketVersioningResult struct {
	oss.ResultCommon
}

// PutAgenticBucketVersioning Configures the versioning state of an agentic bucket.
func (c *AgenticBucketClient) PutAgenticBucketVersioning(ctx context.Context, request *PutAgenticBucketVersioningRequest, optFns ...func(*oss.Options)) (*PutAgenticBucketVersioningResult, error) {
	var err error
	if request == nil {
		request = &PutAgenticBucketVersioningRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutAgenticBucketVersioning",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"versioning":    "",
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

	result := &PutAgenticBucketVersioningResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- GetAgenticBucketVersioning ---

// GetAgenticBucketVersioningRequest is the request for the GetAgenticBucketVersioning operation.
type GetAgenticBucketVersioningRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// GetAgenticBucketVersioningResult is the result for the GetAgenticBucketVersioning operation.
type GetAgenticBucketVersioningResult struct {
	// The versioning configuration of the agentic bucket.
	VersioningConfiguration *VersioningConfiguration `output:"body,VersioningConfiguration,xml"`

	oss.ResultCommon
}

// GetAgenticBucketVersioning Queries the versioning state of an agentic bucket.
func (c *AgenticBucketClient) GetAgenticBucketVersioning(ctx context.Context, request *GetAgenticBucketVersioningRequest, optFns ...func(*oss.Options)) (*GetAgenticBucketVersioningResult, error) {
	var err error
	if request == nil {
		request = &GetAgenticBucketVersioningRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketVersioning",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"versioning":    "",
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

	result := &GetAgenticBucketVersioningResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
