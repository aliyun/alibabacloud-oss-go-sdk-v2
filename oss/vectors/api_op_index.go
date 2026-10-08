package vectors

import (
	"context"
	"encoding/json"
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
	//
	// Each element is the raw JSON object of a field, so an attribute the service adds later is
	// passed through without an SDK change. Build an element with FieldSchema.ToMap, convert a whole
	// list with FieldSchemas(...).ToMaps, or write the map directly to set an attribute that the SDK
	// does not model yet.
	Fields []map[string]any `json:"fields,omitempty"`
}

// FieldSchemas returns the field configurations as the strongly-typed FieldSchema model.
//
// This is a convenience view over Fields: an attribute that the SDK does not model is not visible
// here, and a value that does not fit the typed model makes the conversion fail. Use Fields to
// access the complete definition.
//
// It returns nil when no field is set.
func (s SchemaConfiguration) FieldSchemas() ([]FieldSchema, error) {
	data, err := json.Marshal(s.Fields)
	if err != nil {
		return nil, err
	}
	var schemas []FieldSchema
	if err = json.Unmarshal(data, &schemas); err != nil {
		return nil, err
	}
	return schemas, nil
}

// FieldSchema defines the configuration of a single field in the vector index schema.
type FieldSchema struct {
	// The name of the field.
	Name *string `json:"name,omitempty"`

	// The type of the field. Valid values: vector, string, long, double, bool, ip, geoPoint.
	// The values are also declared as the FieldType constants, e.g. oss.Ptr(string(FieldTypeVector)).
	Type *string `json:"type,omitempty"`

	// The data type of the vector field. Valid values: float32.
	// This parameter is required only when Type is set to vector.
	// The value is also declared as the VectorDataTypeFloat32 constant.
	DataType *string `json:"dataType,omitempty"`

	// The dimension of the vector field.
	// This parameter is required only when Type is set to vector.
	Dimension *int `json:"dimension,omitempty"`

	// The distance metric of the vector field. Valid values: euclidean, cosine, ip.
	// This parameter is required only when Type is set to vector.
	// The values are also declared as the DistanceMetricType constants, e.g. oss.Ptr(string(DistanceMetricTypeCosine)).
	DistanceMetric *string `json:"distanceMetric,omitempty"`

	// Specifies whether the field is an array.
	IsArray *bool `json:"isArray,omitempty"`

	// Specifies whether the field is a partition key.
	IsPartitionKey *bool `json:"isPartitionKey,omitempty"`

	// Specifies whether exact match is enabled for the string field.
	ExactMatch *bool `json:"exactMatch,omitempty"`

	// The text search configuration of the string field.
	Text *TextSchema `json:"text,omitempty"`
}

// ToMap returns the field definition as the raw JSON object that the service expects, so that it
// can be used as an element of SchemaConfiguration.Fields. Attributes that the SDK does not model
// can be added to the returned map directly.
func (s FieldSchema) ToMap() map[string]any {
	m := make(map[string]any)
	if s.Name != nil {
		m["name"] = *s.Name
	}
	if s.Type != nil {
		m["type"] = *s.Type
	}
	if s.DataType != nil {
		m["dataType"] = *s.DataType
	}
	if s.Dimension != nil {
		m["dimension"] = *s.Dimension
	}
	if s.DistanceMetric != nil {
		m["distanceMetric"] = *s.DistanceMetric
	}
	if s.IsArray != nil {
		m["isArray"] = *s.IsArray
	}
	if s.IsPartitionKey != nil {
		m["isPartitionKey"] = *s.IsPartitionKey
	}
	if s.ExactMatch != nil {
		m["exactMatch"] = *s.ExactMatch
	}
	if s.Text != nil {
		m["text"] = s.Text.toMap()
	}
	return m
}

// FieldSchemas is a list of field configurations. It converts the strongly-typed FieldSchema model
// to the raw JSON objects that the service expects.
type FieldSchemas []FieldSchema

// ToMaps returns the field configurations as the raw JSON objects that the service expects, so that
// the result can be assigned to SchemaConfiguration.Fields. It is the batch counterpart of
// FieldSchema.ToMap: attributes that the SDK does not model can be added to each returned map
// directly.
func (s FieldSchemas) ToMaps() []map[string]any {
	maps := make([]map[string]any, 0, len(s))
	for _, schema := range s {
		maps = append(maps, schema.ToMap())
	}
	return maps
}

// TextSchema defines the text search configuration for a string field.
type TextSchema struct {
	// Specifies whether text search is enabled for the field.
	Enabled *bool `json:"enabled,omitempty"`

	// The analyzer used for text search. Valid values: standard, split.
	Analyzer *string `json:"analyzer,omitempty"`

	// The parameters of the analyzer.
	AnalyzerParameters *AnalyzerParameters `json:"analyzerParameters,omitempty"`
}

// toMap returns the text configuration as the raw JSON object of the "text" attribute.
func (s TextSchema) toMap() map[string]any {
	m := make(map[string]any)
	if s.Enabled != nil {
		m["enabled"] = *s.Enabled
	}
	if s.Analyzer != nil {
		m["analyzer"] = *s.Analyzer
	}
	if s.AnalyzerParameters != nil {
		p := make(map[string]any)
		if s.AnalyzerParameters.CaseSensitive != nil {
			p["caseSensitive"] = *s.AnalyzerParameters.CaseSensitive
		}
		if s.AnalyzerParameters.DelimitWord != nil {
			p["delimitWord"] = *s.AnalyzerParameters.DelimitWord
		}
		if s.AnalyzerParameters.Delimiter != nil {
			p["delimiter"] = *s.AnalyzerParameters.Delimiter
		}
		m["analyzerParameters"] = p
	}
	return m
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
