package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutTableBucketEncryptionRequest struct {
	// The name of the bucket.
	BucketArn *string `input:"nop,bucketArn,required"`

	// The encryption of the table bucket.
	EncryptionConfiguration *EncryptionConfiguration `input:"body,encryptionConfiguration,json,required"`

	oss.RequestCommon
}

type PutTableBucketEncryptionResult struct {
	oss.ResultCommon
}

// PutTableBucketEncryption Configures encryption rules for a bucket.
func (c *TablesClient) PutTableBucketEncryption(ctx context.Context, request *PutTableBucketEncryptionRequest, optFns ...func(*oss.Options)) (*PutTableBucketEncryptionResult, error) {
	var err error
	if request == nil {
		request = &PutTableBucketEncryptionRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutTableBucketEncryption",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutTableBucketEncryptionResult{}

	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetTableBucketEncryptionRequest struct {
	BucketArn *string `input:"nop,bucketArn,required"`

	oss.RequestCommon
}

type GetTableBucketEncryptionResult struct {
	EncryptionConfiguration *EncryptionConfiguration `output:"body,encryptionConfiguration,json"`

	oss.ResultCommon
}

// GetTableBucketEncryption Queries the encryption rules configured for a bucket.
func (c *TablesClient) GetTableBucketEncryption(ctx context.Context, request *GetTableBucketEncryptionRequest, optFns ...func(*oss.Options)) (*GetTableBucketEncryptionResult, error) {
	var err error
	if request == nil {
		request = &GetTableBucketEncryptionRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableBucketEncryption",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetTableBucketEncryptionResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type DeleteTableBucketEncryptionRequest struct {
	BucketArn *string `input:"nop,bucketArn,required"`

	oss.RequestCommon
}

type DeleteTableBucketEncryptionResult struct {
	oss.ResultCommon
}

// DeleteTableBucketEncryption Deletes encryption rules for a bucket.
func (c *TablesClient) DeleteTableBucketEncryption(ctx context.Context, request *DeleteTableBucketEncryptionRequest, optFns ...func(*oss.Options)) (*DeleteTableBucketEncryptionResult, error) {
	var err error
	if request == nil {
		request = &DeleteTableBucketEncryptionRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteTableBucketEncryption",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteTableBucketEncryptionResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
