package tables

import (
	"context"
	"io"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutTablePolicyRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	// The request parameters.
	Body io.Reader `input:"body,nop,required"`

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
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
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
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	oss.RequestCommon
}

type GetTablePolicyResult struct {
	Body string

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
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(output.Body)
	defer output.Body.Close()
	if err != nil {
		return nil, err
	}
	result := &GetTablePolicyResult{
		Body: string(body),
	}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type DeleteTablePolicyRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

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
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
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
