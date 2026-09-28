package agentic

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
)

// PutAgenticBucketStorageQuotaRequest is the request for the PutAgenticBucketStorageQuota operation.
type PutAgenticBucketStorageQuotaRequest struct {
	// Bucket is the name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	// QuotaConfiguration is the default storage quota configuration for new bucket spaces.
	QuotaConfiguration *oss.QuotaConfiguration `input:"body,QuotaConfiguration,xml,required"`

	oss.RequestCommon
}

// PutAgenticBucketStorageQuotaResult is the result for the PutAgenticBucketStorageQuota operation.
type PutAgenticBucketStorageQuotaResult struct {
	oss.ResultCommon
}

// PutAgenticBucketStorageQuota configures the default storage quota for new bucket spaces.
func (c *AgenticBucketClient) PutAgenticBucketStorageQuota(ctx context.Context, request *PutAgenticBucketStorageQuotaRequest, optFns ...func(*oss.Options)) (*PutAgenticBucketStorageQuotaResult, error) {
	if request == nil {
		request = &PutAgenticBucketStorageQuotaRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutAgenticBucketStorageQuota",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"quota":         "",
		},
		Bucket: request.Bucket,
	}
	if err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &PutAgenticBucketStorageQuotaResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, nil
}

// GetAgenticBucketStorageQuotaRequest is the request for the GetAgenticBucketStorageQuota operation.
type GetAgenticBucketStorageQuotaRequest struct {
	// Bucket is the name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// GetAgenticBucketStorageQuotaResult is the result for the GetAgenticBucketStorageQuota operation.
type GetAgenticBucketStorageQuotaResult struct {
	// QuotaConfiguration is the default storage quota configuration for new bucket spaces.
	QuotaConfiguration *oss.QuotaConfiguration `output:"body,QuotaConfiguration,xml"`

	oss.ResultCommon
}

// GetAgenticBucketStorageQuota queries the default storage quota for new bucket spaces.
func (c *AgenticBucketClient) GetAgenticBucketStorageQuota(ctx context.Context, request *GetAgenticBucketStorageQuotaRequest, optFns ...func(*oss.Options)) (*GetAgenticBucketStorageQuotaResult, error) {
	if request == nil {
		request = &GetAgenticBucketStorageQuotaRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketStorageQuota",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"quota":         "",
		},
		Bucket: request.Bucket,
	}
	if err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &GetAgenticBucketStorageQuotaResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, nil
}

// DeleteAgenticBucketStorageQuotaRequest is the request for the DeleteAgenticBucketStorageQuota operation.
type DeleteAgenticBucketStorageQuotaRequest struct {
	// Bucket is the name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// DeleteAgenticBucketStorageQuotaResult is the result for the DeleteAgenticBucketStorageQuota operation.
type DeleteAgenticBucketStorageQuotaResult struct {
	oss.ResultCommon
}

// DeleteAgenticBucketStorageQuota deletes the default storage quota configuration for new bucket spaces.
func (c *AgenticBucketClient) DeleteAgenticBucketStorageQuota(ctx context.Context, request *DeleteAgenticBucketStorageQuotaRequest, optFns ...func(*oss.Options)) (*DeleteAgenticBucketStorageQuotaResult, error) {
	if request == nil {
		request = &DeleteAgenticBucketStorageQuotaRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteAgenticBucketStorageQuota",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"quota":         "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"agenticBucket", "quota"})
	if err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &DeleteAgenticBucketStorageQuotaResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, nil
}
