package vectors

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutVectorsRequest struct {
	// The name of the vector bucket.
	Bucket *string `input:"host,bucket,required"`

	IndexName *string `input:"body,indexName,json,required"`

	Vectors []map[string]any `input:"body,vectors,json,required"`

	oss.RequestCommon
}

type PutVectorsResult struct {
	oss.ResultCommon
}

// PutVectors Creates a vector.
func (c *VectorsClient) PutVectors(ctx context.Context, request *PutVectorsRequest, optFns ...func(*oss.Options)) (*PutVectorsResult, error) {
	var err error
	if request == nil {
		request = &PutVectorsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectors": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutVectorsResult{}

	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetVectorsRequest struct {
	// The name of the vector bucket.
	Bucket *string `input:"host,bucket,required"`

	IndexName      *string  `input:"body,indexName,json,required"`
	Keys           []string `input:"body,keys,json,required"`
	ReturnData     *bool    `input:"body,returnData,json"`
	ReturnMetadata *bool    `input:"body,returnMetadata,json"`

	oss.RequestCommon
}

type GetVectorsResult struct {
	Vectors []map[string]any `json:"vectors"`

	oss.ResultCommon
}

// GetVectors Get a vector.
func (c *VectorsClient) GetVectors(ctx context.Context, request *GetVectorsRequest, optFns ...func(*oss.Options)) (*GetVectorsResult, error) {
	var err error
	if request == nil {
		request = &GetVectorsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"getVectors": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetVectorsResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type ListVectorsRequest struct {
	// The name of the vector bucket.
	Bucket *string `input:"host,bucket,required"`

	IndexName *string `input:"body,indexName,json,required"`

	// The maximum number of indexes that can be returned.
	MaxResults *int `input:"body,maxResults,json"`

	NextToken *string `input:"body,nextToken,json"`

	ReturnData *bool `input:"body,returnData,json"`

	ReturnMetadata *bool `input:"body,ReturnMetadata,json"`

	SegmentCount *int `input:"body,SegmentCount,json"`

	SegmentIndex *int `input:"body,SegmentIndex,json"`

	oss.RequestCommon
}

type ListVectorsResult struct {
	// The marker for the next ListVectors request, which can be used to return the remaining results.
	NextToken *string `json:"NextToken"`

	// The container that stores information about vector.
	Vectors []map[string]any `json:"Vectors"`

	oss.ResultCommon
}

// ListVectors Lists vectors that belong to the current account.
func (c *VectorsClient) ListVectors(ctx context.Context, request *ListVectorsRequest, optFns ...func(*oss.Options)) (*ListVectorsResult, error) {
	var err error
	if request == nil {
		request = &ListVectorsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectors": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &ListVectorsResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteVectorsRequest struct {
	// The name of the vector bucket.
	Bucket    *string  `input:"host,bucket,required"`
	IndexName *string  `input:"body,indexName,json,required"`
	Keys      []string `input:"body,keys,json,required"`

	oss.RequestCommon
}

type DeleteVectorsResult struct {
	oss.ResultCommon
}

// DeleteVectors Deletes a vector.
func (c *VectorsClient) DeleteVectors(ctx context.Context, request *DeleteVectorsRequest, optFns ...func(*oss.Options)) (*DeleteVectorsResult, error) {
	var err error
	if request == nil {
		request = &DeleteVectorsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectors": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteVectorsResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type QueryVectorsRequest struct {
	// The name of the vector bucket.
	Bucket         *string                `input:"host,bucket,required"`
	IndexName      *string                `input:"body,indexName,json,required"`
	QueryVector    map[string]interface{} `input:"body,queryVector,json,required"`
	TopK           *int                   `input:"body,topK,json,required"`
	Filter         *string                `input:"body,filter,json"`
	ReturnDistance *bool                  `input:"body,returnDistance,json"`
	ReturnMetadata *bool                  `input:"body,returnMetadata,json"`

	oss.RequestCommon
}

type QueryVectorsResult struct {
	Vectors []map[string]any `json:"vectors"`

	oss.ResultCommon
}

// QueryVectors Query a vector.
func (c *VectorsClient) QueryVectors(ctx context.Context, request *QueryVectorsRequest, optFns ...func(*oss.Options)) (*QueryVectorsResult, error) {
	var err error
	if request == nil {
		request = &QueryVectorsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"queryVectors": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &QueryVectorsResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
