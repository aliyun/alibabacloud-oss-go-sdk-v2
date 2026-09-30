package vectors

import (
	"context"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutVectorIndexRequest struct {
	// The name of the vector bucket.
	Bucket         *string        `input:"host,bucket,required"`
	IndexName      *string        `input:"body,indexName,json,required"`
	DataType       *string        `input:"body,dataType,json,required"`
	Dimension      *int           `input:"body,dimension,json,required"`
	DistanceMetric *string        `input:"body,distanceMetric,json,required"`
	Metadata       map[string]any `input:"body,metadata,json"`

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

	oss.ResultCommon
}

type VectorIndex struct {
	CreateTime     *time.Time     `json:"createTime"`
	DataType       *string        `json:"dataType"`
	Dimension      *int           `json:"dimension"`
	DistanceMetric *string        `json:"distanceMetric"`
	IndexName      *string        `json:"indexName"`
	Metadata       map[string]any `json:"metadata"`
	Status         *string        `json:"status"`
	BucketArn      *string        `json:"bucketArn"`

	// deprecated
	VectorBucketName *string `json:"vectorBucketName"`

	Mode                *string              `json:"mode"`
	SchemaConfiguration *SchemaConfiguration `json:"schemaConfiguration"`
}

// SchemaConfiguration defines the schema configuration of a vector index.
type SchemaConfiguration struct {
	// The container that stores the field configurations.
	Fields []SchemaField `json:"fields,omitempty"`
}

// SchemaField defines the configuration of a single field in the vector index schema.
type SchemaField struct {
	// The name of the field.
	Name *string `json:"name,omitempty"`

	// The type of the field. Valid values: vector, string, long, double, ip, geoPoint.
	Type FieldType   `json:"type,omitempty"`

	// The data type of the vector field. Valid values: float32.
	// This parameter is required only when Type is set to vector.
	DataType VectorDataType `json:"dataType,omitempty"`

	// The dimension of the vector field.
	// This parameter is required only when Type is set to vector.
	Dimension *int `json:"dimension,omitempty"`

	// The distance metric of the vector field. Valid values: euclidean, cosine, inner_product.
	// This parameter is required only when Type is set to vector.
	DistanceMetric DistanceMetricType `json:"distanceMetric,omitempty"`

	// Specifies whether the field is an array.
	IsArray *bool `json:"isArray,omitempty"`

	// Specifies whether the field is a partition key.
	IsPartitionKey *bool `json:"isPartitionKey,omitempty"`

	// Specifies whether exact match is enabled for the string field.
	ExactMatch *bool `json:"exactMatch,omitempty"`

	// The text search configuration of the string field.
	Text *TextConfiguration `json:"text,omitempty"`
}

// TextConfiguration defines the text search configuration for a string field.
type TextConfiguration struct {
	// Specifies whether text search is enabled for the field.
	Enabled *bool `json:"enabled,omitempty"`

	// The analyzer used for text search. Valid values: standard, split.
	Analyzer *string `json:"analyzer,omitempty"`

	// The parameters of the analyzer.
	AnalyzerParameters *AnalyzerParameters `json:"analyzerParameters,omitempty"`
}

// AnalyzerParameters defines the parameters of the analyzer used for text search.
type AnalyzerParameters struct {
	// Specifies whether the analyzer is case-sensitive.
	CaseSensitive *bool `json:"caseSensitive,omitempty"`

	// Specifies whether words are delimited. This parameter is valid only when
	// Analyzer is set to standard.
	DelimitWord *bool `json:"delimitWord,omitempty"`

	// The delimiter used to split words. This parameter is valid only when
	// Analyzer is set to split.
	Delimiter *string `json:"delimiter,omitempty"`
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
	MaxResults int `input:"body,maxResults,json"`

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

type PutVectorIndexFusionRequest struct {
	// The name of the vector bucket.
	Bucket              *string              `input:"host,bucket,required"`
	IndexName           *string              `input:"body,indexName,json,required"`
	Mode                *string              `input:"body,mode,json,required"`
	SchemaConfiguration *SchemaConfiguration `input:"body,schemaConfiguration,json,required"`

	oss.RequestCommon
}

type PutVectorIndexFusionResult struct {
	oss.ResultCommon
}

// PutVectorIndexFusion Creates a vector Index by fusion mode.
func (c *VectorsClient) PutVectorIndexFusion(ctx context.Context, request *PutVectorIndexFusionRequest, optFns ...func(*oss.Options)) (*PutVectorIndexFusionResult, error) {
	var err error
	if request == nil {
		request = &PutVectorIndexFusionRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutVectorIndexFusion",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndexFusion": "",
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

	result := &PutVectorIndexFusionResult{}

	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
