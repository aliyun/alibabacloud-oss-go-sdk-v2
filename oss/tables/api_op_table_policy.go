package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutTablePolicyRequest struct {
	BucketArn *string `input:"nop,bucketArn,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	ResourcePolicy *string `input:"body,resourcePolicy,required,json"`

	oss.RequestCommon
}

type PutTablePolicyResult struct {
	oss.ResultCommon
}

// PutTablePolicy create a table policy.
func (c *TablesClient) PutTablePolicy(ctx context.Context, request *PutTablePolicyRequest, optFns ...func(*oss.Options)) (*PutTablePolicyResult, error) {
	var err error
	if request == nil {
		request = &PutTablePolicyRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Table)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutTablePolicyResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetTablePolicyRequest struct {
	BucketArn *string `input:"nop,bucketArn,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	oss.RequestCommon
}

type GetTablePolicyResult struct {
	ResourcePolicy *string `output:"body,resourcePolicy,json"`

	oss.ResultCommon
}

// GetTablePolicy Queries a table policy.
func (c *TablesClient) GetTablePolicy(ctx context.Context, request *GetTablePolicyRequest, optFns ...func(*oss.Options)) (*GetTablePolicyResult, error) {
	var err error
	if request == nil {
		request = &GetTablePolicyRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Table)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetTablePolicyResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type DeleteTablePolicyRequest struct {
	BucketArn *string `input:"nop,bucketArn,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	oss.RequestCommon
}

type DeleteTablePolicyResult struct {
	oss.ResultCommon
}

// DeleteTablePolicy delete a table policy.
func (c *TablesClient) DeleteTablePolicy(ctx context.Context, request *DeleteTablePolicyRequest, optFns ...func(*oss.Options)) (*DeleteTablePolicyResult, error) {
	var err error
	if request == nil {
		request = &DeleteTablePolicyRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteTablePolicy",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Table)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteTablePolicyResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
