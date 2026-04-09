package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type CreateTableBucketRequest struct {
	// The name of the table bucket to create.
	Bucket *string `input:"body,name,json,required"`

	// The encryption of the table bucket.
	EncryptionConfiguration *EncryptionConfiguration `input:"body,encryptionConfiguration,json"`

	oss.RequestCommon
}

type EncryptionConfiguration struct {
	KmsKeyArn *string `json:"kmsKeyArn,omitempty"`

	SseAlgorithm *string `json:"sseAlgorithm,omitempty"`
}

type CreateTableBucketResult struct {
	BucketArn *string `json:"arn"`

	oss.ResultCommon
}

// CreateTableBucket Creates a table bucket.
func (c *TablesClient) CreateTableBucket(ctx context.Context, request *CreateTableBucketRequest, optFns ...func(*oss.Options)) (*CreateTableBucketResult, error) {
	var err error
	if request == nil {
		request = &CreateTableBucketRequest{}
	}
	input := &oss.OperationInput{
		OpName: "CreateTableBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("buckets"),
	}

	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &CreateTableBucketResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetTableBucketRequest struct {
	BucketArn *string `input:"nop,bucketArn,required"`

	oss.RequestCommon
}

type GetTableBucketResult struct {
	BucketArn      *string `json:"arn"`
	CreatedAt      *string `json:"createdAt"`
	Name           *string `json:"name"`
	OwnerAccountId *string `json:"ownerAccountId"`
	TableBucketId  *string `json:"tableBucketId"`
	Type           *string `json:"type"`

	oss.ResultCommon
}

// GetTableBucket Queries information about a table bucket.
func (c *TablesClient) GetTableBucket(ctx context.Context, request *GetTableBucketRequest, optFns ...func(*oss.Options)) (*GetTableBucketResult, error) {
	var err error
	if request == nil {
		request = &GetTableBucketRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableBucket",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetTableBucketResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type ListTableBucketsRequest struct {
	// The token from which the ListTableBuckets operation must start.
	ContinuationToken *string `input:"query,continuationToken"`

	// The maximum number of buckets that can be returned in the single query.
	// Valid values: 1 to 1000.
	MaxBuckets int32 `input:"query,maxBuckets"`

	// The prefix that the names of returned buckets must contain.
	Prefix *string `input:"query,prefix"` // Limits the response to keys that begin with the specified prefix

	oss.RequestCommon
}

type ListTableBucketsResult struct {
	// The token from which the ListTableBuckets operation must start.
	ContinuationToken *string `json:"continuationToken"`

	// The container that stores information about buckets.
	Buckets []TableBucketProperties `json:"tableBuckets"`

	oss.ResultCommon
}

type TableBucketProperties struct {
	BucketArn      *string `json:"arn"`
	CreatedAt      *string `json:"createdAt"`
	Name           *string `json:"name"`
	OwnerAccountId *string `json:"ownerAccountId"`
	TableBucketId  *string `json:"tableBucketId"`
	Type           *string `json:"type"`
}

// ListTableBuckets Lists table buckets that belong to the current account.
func (c *TablesClient) ListTableBuckets(ctx context.Context, request *ListTableBucketsRequest, optFns ...func(*oss.Options)) (*ListTableBucketsResult, error) {
	var err error
	if request == nil {
		request = &ListTableBucketsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "ListTableBuckets",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("buckets"),
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &ListTableBucketsResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteTableBucketRequest struct {
	// The bucket arn of the table bucket to delete.
	BucketArn *string `input:"nop,bucketArn,required"`

	oss.RequestCommon
}

type DeleteTableBucketResult struct {
	oss.ResultCommon
}

// DeleteTableBucket Deletes a table bucket.
func (c *TablesClient) DeleteTableBucket(ctx context.Context, request *DeleteTableBucketRequest, optFns ...func(*oss.Options)) (*DeleteTableBucketResult, error) {
	var err error
	if request == nil {
		request = &DeleteTableBucketRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteTableBucket",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteTableBucketResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
