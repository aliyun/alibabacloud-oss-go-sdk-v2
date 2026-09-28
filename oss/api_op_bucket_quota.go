package oss

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
)

// StorageQuotaModeType is the enforcement mode for a bucket storage quota.
type StorageQuotaModeType string

const (
	// StorageQuotaModeStrict rejects writes after the storage quota is exceeded.
	StorageQuotaModeStrict StorageQuotaModeType = "Strict"

	// StorageQuotaModeWarning allows writes after the storage quota is exceeded and returns a warning.
	StorageQuotaModeWarning StorageQuotaModeType = "Warning"
)

// QuotaConfiguration is the storage quota configuration of a bucket.
type QuotaConfiguration struct {
	// StorageQuota is the maximum storage capacity in bytes.
	StorageQuota *int64 `xml:"StorageQuota"`

	// Mode is the enforcement mode for the storage quota.
	Mode StorageQuotaModeType `xml:"Mode"`

	// CurrentUsage is the most recently calculated storage usage in bytes.
	CurrentUsage *int64 `xml:"CurrentUsage,omitempty"`
}

// PutBucketStorageQuotaRequest is the request for the PutBucketStorageQuota operation.
type PutBucketStorageQuotaRequest struct {
	// Bucket is the name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// QuotaConfiguration is the storage quota configuration to apply.
	QuotaConfiguration *QuotaConfiguration `input:"body,QuotaConfiguration,xml,required"`

	RequestCommon
}

// PutBucketStorageQuotaResult is the result for the PutBucketStorageQuota operation.
type PutBucketStorageQuotaResult struct {
	ResultCommon
}

// PutBucketStorageQuota configures the storage quota of a bucket.
func (c *Client) PutBucketStorageQuota(ctx context.Context, request *PutBucketStorageQuotaRequest, optFns ...func(*Options)) (*PutBucketStorageQuotaResult, error) {
	if request == nil {
		request = &PutBucketStorageQuotaRequest{}
	}
	input := &OperationInput{
		OpName: "PutBucketStorageQuota",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"quota": ""},
		Bucket:     request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"quota"})
	if err := c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}
	result := &PutBucketStorageQuotaResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, nil
}

// GetBucketStorageQuotaRequest is the request for the GetBucketStorageQuota operation.
type GetBucketStorageQuotaRequest struct {
	// Bucket is the name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	RequestCommon
}

// GetBucketStorageQuotaResult is the result for the GetBucketStorageQuota operation.
type GetBucketStorageQuotaResult struct {
	// QuotaConfiguration is the storage quota configuration of the bucket.
	QuotaConfiguration *QuotaConfiguration `output:"body,QuotaConfiguration,xml"`

	ResultCommon
}

// GetBucketStorageQuota queries the storage quota of a bucket.
func (c *Client) GetBucketStorageQuota(ctx context.Context, request *GetBucketStorageQuotaRequest, optFns ...func(*Options)) (*GetBucketStorageQuotaResult, error) {
	if request == nil {
		request = &GetBucketStorageQuotaRequest{}
	}
	input := &OperationInput{
		OpName:     "GetBucketStorageQuota",
		Method:     "GET",
		Parameters: map[string]string{"quota": ""},
		Bucket:     request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"quota"})
	if err := c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}
	result := &GetBucketStorageQuotaResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, nil
}

// DeleteBucketStorageQuotaRequest is the request for the DeleteBucketStorageQuota operation.
type DeleteBucketStorageQuotaRequest struct {
	// Bucket is the name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	RequestCommon
}

// DeleteBucketStorageQuotaResult is the result for the DeleteBucketStorageQuota operation.
type DeleteBucketStorageQuotaResult struct {
	ResultCommon
}

// DeleteBucketStorageQuota deletes the storage quota configuration of a bucket.
func (c *Client) DeleteBucketStorageQuota(ctx context.Context, request *DeleteBucketStorageQuotaRequest, optFns ...func(*Options)) (*DeleteBucketStorageQuotaResult, error) {
	if request == nil {
		request = &DeleteBucketStorageQuotaRequest{}
	}
	input := &OperationInput{
		OpName:     "DeleteBucketStorageQuota",
		Method:     "DELETE",
		Parameters: map[string]string{"quota": ""},
		Bucket:     request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"quota"})
	if err := c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}
	result := &DeleteBucketStorageQuotaResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, nil
}
