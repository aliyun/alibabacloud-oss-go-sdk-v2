package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutTableBucketPolicyRequest struct {
	// The bucket arn of the bucket.
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	ResourcePolicy *string `input:"body,resourcePolicy,required,json"`

	oss.RequestCommon
}

type PutTableBucketPolicyResult struct {
	oss.ResultCommon
}

// PutTableBucketPolicy Configures a policy for a table bucket.
func (c *TablesClient) PutTableBucketPolicy(ctx context.Context, request *PutTableBucketPolicyRequest, optFns ...func(*oss.Options)) (*PutTableBucketPolicyResult, error) {
	var err error
	if request == nil {
		request = &PutTableBucketPolicyRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutTableBucketPolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/policy", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutTableBucketPolicyResult{}

	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetTableBucketPolicyRequest struct {
	// The bucket arn of the bucket.
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	oss.RequestCommon
}
type GetTableBucketPolicyResult struct {
	ResourcePolicy *string `json:"resourcePolicy"`

	oss.ResultCommon
}

// GetTableBucketPolicy Queries the policies configured for a table bucket.
func (c *TablesClient) GetTableBucketPolicy(ctx context.Context, request *GetTableBucketPolicyRequest, optFns ...func(*oss.Options)) (*GetTableBucketPolicyResult, error) {
	var err error
	if request == nil {
		request = &GetTableBucketPolicyRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableBucketPolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/policy", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetTableBucketPolicyResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type DeleteTableBucketPolicyRequest struct {
	// The bucket arn of the bucket.
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	oss.RequestCommon
}

type DeleteTableBucketPolicyResult struct {
	oss.ResultCommon
}

// DeleteTableBucketPolicy Deletes a policy for a table bucket.
func (c *TablesClient) DeleteTableBucketPolicy(ctx context.Context, request *DeleteTableBucketPolicyRequest, optFns ...func(*oss.Options)) (*DeleteTableBucketPolicyResult, error) {
	var err error
	if request == nil {
		request = &DeleteTableBucketPolicyRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteTableBucketPolicy",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/policy", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &DeleteTableBucketPolicyResult{}

	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
