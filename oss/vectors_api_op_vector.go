package oss

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
)

type PutVectorsRequest struct {
	// The name of the vector bucket.
	Bucket *string `input:"host,bucket,required"`

	IndexName *string `input:"body,indexName,json,required"`

	Vectors []map[string]interface{} `input:"body,vectors,json,required"`

	RequestCommon
}

type PutVectorsResult struct {
	ResultCommon
}

// PutVectors Creates a vector.
func (c *VectorsClient) PutVectors(ctx context.Context, request *PutVectorsRequest, optFns ...func(*Options)) (*PutVectorsResult, error) {
	var err error
	if request == nil {
		request = &PutVectorsRequest{}
	}
	input := &OperationInput{
		OpName: "PutVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"PutVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectors"})
	if err = c.client.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &PutVectorsResult{}

	if err = c.client.unmarshalOutput(result, output, discardBody); err != nil {
		return nil, c.client.toClientError(err, "UnmarshalOutputFail", output)
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

	RequestCommon
}

type GetVectorsResult struct {
	Vectors []map[string]interface{} `json:"vectors"`

	ResultCommon
}

// GetVectors Get a vector.
func (c *VectorsClient) GetVectors(ctx context.Context, request *GetVectorsRequest, optFns ...func(*Options)) (*GetVectorsResult, error) {
	var err error
	if request == nil {
		request = &GetVectorsRequest{}
	}
	input := &OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"GetVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectors"})
	if err = c.client.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &GetVectorsResult{}

	if err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault); err != nil {
		return nil, c.client.toClientError(err, "UnmarshalOutputFail", output)
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

	RequestCommon
}

type ListVectorsResult struct {
	// The marker for the next ListVectors request, which can be used to return the remaining results.
	NextToken *string `json:"NextToken"`

	// The container that stores information about vector.
	Vectors []map[string]interface{} `json:"Vectors"`

	ResultCommon
}

// ListVectors Lists vectors that belong to the current account.
func (c *VectorsClient) ListVectors(ctx context.Context, request *ListVectorsRequest, optFns ...func(*Options)) (*ListVectorsResult, error) {
	var err error
	if request == nil {
		request = &ListVectorsRequest{}
	}
	input := &OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectors"})
	if err = c.client.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.client.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &ListVectorsResult{}
	if err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault); err != nil {
		return nil, c.client.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteVectorsRequest struct {
	// The name of the vector bucket.
	Bucket    *string  `input:"host,bucket,required"`
	IndexName *string  `input:"body,indexName,json,required"`
	Keys      []string `input:"body,keys,json,required"`

	RequestCommon
}

type DeleteVectorsResult struct {
	ResultCommon
}

// DeleteVectors Deletes a vector.
func (c *VectorsClient) DeleteVectors(ctx context.Context, request *DeleteVectorsRequest, optFns ...func(*Options)) (*DeleteVectorsResult, error) {
	var err error
	if request == nil {
		request = &DeleteVectorsRequest{}
	}
	input := &OperationInput{
		OpName: "DeleteVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"DeleteVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectors"})
	if err = c.client.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &DeleteVectorsResult{}
	if err = c.client.unmarshalOutput(result, output, discardBody); err != nil {
		return nil, c.client.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type QueryVectorsRequest struct {
	// The name of the vector bucket.
	Bucket         *string                `input:"host,bucket,required"`
	IndexName      *string                `input:"body,indexName,json,required"`
	QueryVector    map[string]interface{} `input:"body,queryVector,json,required"`
	TopK           *int                   `input:"body,topK,json,required"`
	Filter         *string                `input:"body,filter,json,required"`
	ReturnDistance *bool                  `input:"body,returnDistance,json,required"`
	ReturnMetadata *bool                  `input:"body,returnMetadata,json,required"`

	RequestCommon
}

type QueryVectorsResult struct {
	Vectors []map[string]interface{} `json:"vectors"`

	ResultCommon
}

// QueryVectors Query a vector.
func (c *VectorsClient) QueryVectors(ctx context.Context, request *QueryVectorsRequest, optFns ...func(*Options)) (*QueryVectorsResult, error) {
	var err error
	if request == nil {
		request = &QueryVectorsRequest{}
	}
	input := &OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"QueryVectors"})
	if err = c.client.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &QueryVectorsResult{}
	if err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault); err != nil {
		return nil, c.client.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
