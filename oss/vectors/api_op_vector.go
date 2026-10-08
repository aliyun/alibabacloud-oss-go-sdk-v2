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
	MaxResults int `input:"body,maxResults,json"`

	NextToken *string `input:"body,nextToken,json"`

	ReturnData *bool `input:"body,returnData,json"`

	ReturnMetadata *bool `input:"body,returnMetadata,json"`

	SegmentCount *int `input:"body,segmentCount,json"`

	SegmentIndex *int `input:"body,segmentIndex,json"`

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
	Bucket         *string        `input:"host,bucket,required"`
	IndexName      *string        `input:"body,indexName,json,required"`
	QueryVector    map[string]any `input:"body,queryVector,json,required"`
	TopK           *int           `input:"body,topK,json,required"`
	Filter         any            `input:"body,filter,json"`
	ReturnDistance *bool          `input:"body,returnDistance,json"`
	ReturnMetadata *bool          `input:"body,returnMetadata,json"`

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

type QueryVectorsFusionRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the index.
	IndexName *string `input:"body,indexName,json,required"`

	// The knn vector queries, exactly as they are sent to the service. A single vector query is
	// also represented as a one-element list.
	//
	// Build an element with Knn.ToMap, or write the map directly.
	Knn []map[string]any `input:"body,knn,json"`

	// The conditions of the scalar query and the full text query.
	Query map[string]any `input:"body,query,json"`

	// The multi-way hybrid retriever, exactly as it is sent to the service. It carries the raw
	// JSON object.
	//
	// The object has a single type key: simple, knn, rrf or weight. Build the value with the
	// ToRetriever method of the corresponding retriever, e.g. RrfRetriever{...}.ToRetriever(), or
	// write the map directly.
	Retriever map[string]any `input:"body,retriever,json"`

	// Whether to return the metadata. Default value: false.
	ReturnMetadata *bool `input:"body,returnMetadata,json"`

	// The metadata fields to return. It takes effect only when returnMetadata is true.
	ReturnMetadataFields []string `input:"body,returnMetadataFields,json"`

	// The partition keys to access. The server routes the query to the specified partitions.
	PartitionKeys []string `input:"body,partitionKeys,json"`

	// The number of the rows returned by the request. Default value: 10.
	Limit *int `input:"body,limit,json"`

	// The token for the next page. It is supported only when the request contains the query
	// parameter.
	NextToken *string `input:"body,nextToken,json"`

	// The sort fields. A maximum of 3 sort fields are supported.
	Sort []Sort `input:"body,sort,json"`

	oss.RequestCommon
}

// Knn defines a single knn vector query of the QueryVectorsFusion operation.
type Knn struct {
	// The name of the vector field.
	Field *string `json:"field,omitempty"`

	// The query vector.
	QueryVector any `json:"queryVector,omitempty"`

	// The number of the top K results returned by the knn query. Default value: 10.
	TopK *int `json:"topK,omitempty"`

	// The pre-filter of the knn query. The syntax is the same as the query parameter.
	Filter any `json:"filter,omitempty"`

	// The number of the candidates. It enlarges the search scope to improve the recall rate.
	// The value must be greater than or equal to topK.
	NumCandidates *int `json:"numCandidates,omitempty"`

	// The weight of the query. It is a non-negative float32 number. Default value: 1.0.
	Boost *float32 `json:"boost,omitempty"`
}

// ToMap returns the knn query as the raw JSON object that the service expects, so that it can be
// used as an element of QueryVectorsFusionRequest.Knn.
func (s Knn) ToMap() map[string]any {
	m := make(map[string]any)
	if s.Field != nil {
		m["field"] = *s.Field
	}
	if s.QueryVector != nil {
		m["queryVector"] = s.QueryVector
	}
	if s.TopK != nil {
		m["topK"] = *s.TopK
	}
	if s.Filter != nil {
		m["filter"] = s.Filter
	}
	if s.NumCandidates != nil {
		m["numCandidates"] = *s.NumCandidates
	}
	if s.Boost != nil {
		m["boost"] = *s.Boost
	}
	return m
}

// ToRetriever returns the knn query as the raw JSON object of the "knn" attribute. Use it as the
// value of the "knn" key of QueryVectorsFusionRequest.Retriever or of RetrieverComponent.Retriever.
func (s Knn) ToRetriever() map[string]any {
	return map[string]any{"knn": s.ToMap()}
}

// SimpleRetriever defines the simple retriever that queries documents by the specified conditions.
type SimpleRetriever struct {
	// The query conditions. The syntax is the same as the query parameter of the
	// QueryVectorsFusion operation.
	Query map[string]any `json:"query,omitempty"`
}

// ToMap returns the simple retriever as the raw JSON object of the "simple" attribute. Use it as
// the value of the "simple" key of QueryVectorsFusionRequest.Retriever.
func (s SimpleRetriever) ToMap() map[string]any {
	m := make(map[string]any)
	if s.Query != nil {
		m["query"] = s.Query
	}
	return m
}

// ToRetriever returns the simple retriever as the raw JSON object of the "simple" attribute. Use it
// as the value of the "simple" key of QueryVectorsFusionRequest.Retriever or of
// RetrieverComponent.Retriever.
func (s SimpleRetriever) ToRetriever() map[string]any {
	return map[string]any{"simple": s.ToMap()}
}

// RrfRetriever defines the rrf compound retriever that merges the results of the sub retrievers by
// the Reciprocal Rank Fusion algorithm.
type RrfRetriever struct {
	// The constant k in the RRF formula: score = sum of 1/(k + rank). Valid values: 1 to 65536.
	// Default value: 50.
	K *int32 `json:"k,omitempty"`

	// The number of the top results taken from each sub retriever. Default value: 100.
	WindowSize *int32 `json:"windowSize,omitempty"`

	// The sub retrievers. It contains 1 to 3 elements.
	//
	// Build an element with RetrieverComponent.ToMap, or write the map directly.
	Retrievers []map[string]any `json:"retrievers,omitempty"`
}

// ToMap returns the rrf retriever as the raw JSON object of the "rrf" attribute. Use it as the
// value of the "rrf" key of QueryVectorsFusionRequest.Retriever.
func (s RrfRetriever) ToMap() map[string]any {
	m := make(map[string]any)
	if s.K != nil {
		m["k"] = *s.K
	}
	if s.WindowSize != nil {
		m["windowSize"] = *s.WindowSize
	}
	if s.Retrievers != nil {
		m["retrievers"] = s.Retrievers
	}
	return m
}

// ToRetriever returns the rrf retriever as the raw JSON object of the "rrf" attribute. Use it as
// the value of the "rrf" key of QueryVectorsFusionRequest.Retriever or of
// RetrieverComponent.Retriever.
func (s RrfRetriever) ToRetriever() map[string]any {
	return map[string]any{"rrf": s.ToMap()}
}

// WeightRetriever defines the weight compound retriever that merges the results of the sub
// retrievers by weight.
type WeightRetriever struct {
	// The number of the top results taken from each sub retriever. Default value: 100.
	WindowSize *int32 `json:"windowSize,omitempty"`

	// The sub retrievers. It contains 1 to 3 elements.
	//
	// Build an element with RetrieverComponent.ToMap, or write the map directly.
	Retrievers []map[string]any `json:"retrievers,omitempty"`
}

// ToMap returns the weight retriever as the raw JSON object of the "weight" attribute. Use it as
// the value of the "weight" key of QueryVectorsFusionRequest.Retriever.
func (s WeightRetriever) ToMap() map[string]any {
	m := make(map[string]any)
	if s.WindowSize != nil {
		m["windowSize"] = *s.WindowSize
	}
	if s.Retrievers != nil {
		m["retrievers"] = s.Retrievers
	}
	return m
}

// ToRetriever returns the weight retriever as the raw JSON object of the "weight" attribute. Use it
// as the value of the "weight" key of QueryVectorsFusionRequest.Retriever or of
// RetrieverComponent.Retriever.
func (s WeightRetriever) ToRetriever() map[string]any {
	return map[string]any{"weight": s.ToMap()}
}

// RetrieverComponent defines the component of a compound retriever (RrfRetriever or
// WeightRetriever). It wraps a sub retriever with the fusion weight, and optionally the score
// normalizer.
type RetrieverComponent struct {
	// The nested retriever as the raw JSON object. It can be a leaf retriever (knn/simple) or a
	// nested compound retriever (rrf/weight).
	//
	// Build it with the ToRetriever method of the nested retriever, e.g. Knn{...}.ToRetriever(), or
	// write the map directly.
	Retriever map[string]any `json:"retriever,omitempty"`

	// The weight of this component. It is a non-negative float32 number. Default value: 1.0.
	Weight *float32 `json:"weight,omitempty"`

	// The normalizer of the score. Valid values: none, minMax, l2. The values are also declared as
	// the NormalizerType constants, e.g. oss.Ptr(string(NormalizerTypeMinMax)).
	// It applies to the weight compound retriever only.
	Normalizer *string `json:"normalizer,omitempty"`
}

// ToMap returns the component as the raw JSON object of an element of RrfRetriever.Retrievers or
// WeightRetriever.Retrievers.
func (s RetrieverComponent) ToMap() map[string]any {
	m := make(map[string]any)
	if s.Retriever != nil {
		m["retriever"] = s.Retriever
	}
	if s.Weight != nil {
		m["weight"] = *s.Weight
	}
	if s.Normalizer != nil {
		m["normalizer"] = *s.Normalizer
	}
	return m
}

// SortOptions defines the sorting of a single sort field of the QueryVectorsFusion operation.
type SortOptions struct {
	// The sort order of the field. Valid values: asc, desc.
	// The values are also declared as the SortOrderType constants, e.g. oss.Ptr(string(SortOrderTypeAsc)).
	Order *string `json:"order,omitempty"`
}

// Sort maps a field name to its sort options.
type Sort map[string]SortOptions

type QueryVectorsFusionResult struct {
	// The list of the query result vectors.
	Vectors []QueryVectorsFusionSummary `json:"vectors"`

	// The token for the next page of vectors.
	NextToken *string `json:"nextToken"`

	oss.ResultCommon
}

// QueryVectorsFusionSummary is the summary of the vectors returned by the QueryVectorsFusion
// operation.
type QueryVectorsFusionSummary struct {
	// The key of the vector.
	Key *string `json:"key"`

	// The metadata of the vector.
	Metadata map[string]any `json:"metadata"`

	// The relevance score of the vector.
	Score *float32 `json:"score"`
}

// QueryVectorsFusion Query a vector by fusion mode.
func (c *VectorsClient) QueryVectorsFusion(ctx context.Context, request *QueryVectorsFusionRequest, optFns ...func(*oss.Options)) (*QueryVectorsFusionResult, error) {
	var err error
	if request == nil {
		request = &QueryVectorsFusionRequest{}
	}
	input := &oss.OperationInput{
		OpName: "QueryVectorsFusion",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"queryVectorsFusion": "",
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

	result := &QueryVectorsFusionResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
