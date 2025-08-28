package vectors

import (
	"context"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
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

	oss.RequestCommon
}

type PutVectorIndexResult struct {
	oss.ResultCommon
}

// PutVectorIndex Creates a vector Index.
func (c *VectorsClient) PutVectorIndex(ctx context.Context, request *PutVectorIndexRequest, optFns ...func(*oss.Options)) (*PutVectorIndexResult, error) {
	var err error
	if request == nil {
		request = &PutVectorIndexRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
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

	result := &PutVectorIndexResult{}

	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetVectorIndexRequest struct {
	// The name of the vector bucket.
	Bucket *string `input:"host,bucket,required"`

	IndexName *string `input:"body,indexName,json,required"`

	oss.RequestCommon
}

type GetVectorIndexResult struct {
	Index *VectorIndex `json:"index"`

	VectorBucketName *string `json:"vectorBucketName"`

	oss.ResultCommon
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
func (c *VectorsClient) GetVectorIndex(ctx context.Context, request *GetVectorIndexRequest, optFns ...func(*oss.Options)) (*GetVectorIndexResult, error) {
	var err error
	if request == nil {
		request = &GetVectorIndexRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"getVectorIndex": "",
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

	result := &GetVectorIndexResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
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

	oss.RequestCommon
}

type ListVectorIndexesResult struct {
	// The marker for the next ListVectorIndexes request, which can be used to return the remaining results.
	NextToken *string `json:"NextToken"`

	// The container that stores information about indexes.
	Indexes []VectorIndex `json:"Indexes"`

	oss.ResultCommon
}

// ListVectorIndexes Lists vector indexes that belong to the current account.
func (c *VectorsClient) ListVectorIndexes(ctx context.Context, request *ListVectorIndexesRequest, optFns ...func(*oss.Options)) (*ListVectorIndexesResult, error) {
	var err error
	if request == nil {
		request = &ListVectorIndexesRequest{}
	}
	input := &oss.OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectorIndexes": "",
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

	result := &ListVectorIndexesResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteVectorIndexRequest struct {
	// The name of the vector bucket.
	Bucket *string `input:"host,bucket,required"`

	IndexName *string `input:"body,indexName,json,required"`

	oss.RequestCommon
}

type DeleteVectorIndexResult struct {
	oss.ResultCommon
}

// DeleteVectorIndex Deletes a vector index.
func (c *VectorsClient) DeleteVectorIndex(ctx context.Context, request *DeleteVectorIndexRequest, optFns ...func(*oss.Options)) (*DeleteVectorIndexResult, error) {
	var err error
	if request == nil {
		request = &DeleteVectorIndexRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectorIndex": "",
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

	result := &DeleteVectorIndexResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
