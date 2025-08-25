package oss

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"time"
)

type PutVectorIndexRequest struct {
	// The name of the vector bucket.
	Bucket           *string                `input:"host,bucket,required"`
	IndexName        *string                `input:"body,indexName,json,required"`
	CreateTime       *time.Time             `input:"body,createTime,json"`
	DataType         *string                `input:"body,dataType,json"`
	Dimension        *int                   `input:"body,dimension,json"`
	DistanceMetric   *string                `input:"body,distanceMetric,json"`
	Metadata         map[string]interface{} `input:"body,metadata,json"`
	Status           *string                `input:"body,status,json"`
	VectorBucketName *string                `input:"body,vectorBucketName,json"`

	RequestCommon
}

type PutVectorIndexResult struct {
	ResultCommon
}

// PutVectorIndex Creates a vector Index.
func (c *VectorsClient) PutVectorIndex(ctx context.Context, request *PutVectorIndexRequest, optFns ...func(*Options)) (*PutVectorIndexResult, error) {
	var err error
	if request == nil {
		request = &PutVectorIndexRequest{}
	}
	input := &OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"PutVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	if err = c.client.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &PutVectorIndexResult{}

	if err = c.client.unmarshalOutput(result, output, discardBody); err != nil {
		return nil, c.client.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetVectorIndexRequest struct {
	// The name of the vector bucket.
	Bucket *string `input:"host,bucket,required"`

	IndexName *string `input:"body,indexName,json,required"`

	RequestCommon
}

type GetVectorIndexResult struct {
	Index *VectorIndex `json:"index"`

	VectorBucketName *string `json:"vectorBucketName"`

	ResultCommon
}

type VectorIndex struct {
	CreateTime       *time.Time             `json:"createTime"`
	DataType         *string                `json:"dataType"`
	Dimension        *int                   `json:"dimension"`
	DistanceMetric   *string                `json:"distanceMetric"`
	IndexName        *string                `json:"indexName"`
	Metadata         map[string]interface{} `json:"metadata"`
	Status           *string                `json:"status"`
	VectorBucketName *string                `json:"vectorBucketName"`
}

// GetVectorIndex Get a vector Index.
func (c *VectorsClient) GetVectorIndex(ctx context.Context, request *GetVectorIndexRequest, optFns ...func(*Options)) (*GetVectorIndexResult, error) {
	var err error
	if request == nil {
		request = &GetVectorIndexRequest{}
	}
	input := &OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"GetVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	if err = c.client.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &GetVectorIndexResult{}

	if err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault); err != nil {
		return nil, c.client.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type ListVectorIndexesRequest struct {
	// The name of the vector bucket.
	Bucket *string `input:"host,bucket,required"`

	NextToken *string `input:"body,nextToken,json"`

	// The maximum number of indexes that can be returned.
	MaxResults *int `input:"body,maxResults,json"`

	// The prefix that the names of returned indexes must contain.
	Prefix *string `input:"body,prefix,json"`

	RequestCommon
}

type ListVectorIndexesResult struct {
	// The marker for the next ListVectorIndexes request, which can be used to return the remaining results.
	NextToken *string `json:"NextToken"`

	// The container that stores information about indexes.
	Indexes []VectorIndex `json:"Indexes"`

	ResultCommon
}

// ListVectorIndexes Lists vector indexes that belong to the current account.
func (c *VectorsClient) ListVectorIndexes(ctx context.Context, request *ListVectorIndexesRequest, optFns ...func(*Options)) (*ListVectorIndexesResult, error) {
	var err error
	if request == nil {
		request = &ListVectorIndexesRequest{}
	}
	input := &OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	if err = c.client.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.client.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &ListVectorIndexesResult{}
	if err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault); err != nil {
		return nil, c.client.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteVectorIndexRequest struct {
	// The name of the vector bucket.
	Bucket *string `input:"host,bucket,required"`

	IndexName *string `input:"body,indexName,json,required"`

	RequestCommon
}

type DeleteVectorIndexResult struct {
	ResultCommon
}

// DeleteVectorIndex Deletes a vector index.
func (c *VectorsClient) DeleteVectorIndex(ctx context.Context, request *DeleteVectorIndexRequest, optFns ...func(*Options)) (*DeleteVectorIndexResult, error) {
	var err error
	if request == nil {
		request = &DeleteVectorIndexRequest{}
	}
	input := &OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"DeleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	if err = c.client.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &DeleteVectorIndexResult{}
	if err = c.client.unmarshalOutput(result, output, discardBody); err != nil {
		return nil, c.client.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
