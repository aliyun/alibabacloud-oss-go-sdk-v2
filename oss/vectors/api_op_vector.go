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
	Bucket    *string `input:"host,bucket,required"`
	IndexName *string `input:"body,indexName,json,required"`
	// The knn vector queries, exactly as they are sent to the service. A single vector query is
	// also represented as a one-element list.
	//
	// Each element is the raw JSON object of a query, so an attribute that the service adds later
	// is passed through without an SDK change. Build an element with Knn.ToMap, or write the
	// map directly to set an attribute that the SDK does not model yet.
	Knn   []map[string]any `input:"body,knn,json"`
	Query map[string]any   `input:"body,query,json"`
	// The multi-way hybrid retriever, exactly as it is sent to the service. It carries the raw
	// JSON object, so an attribute that the service adds later is passed through without an SDK
	// change.
	//
	// The object has a single type key: simple, knn, rrf or weight. Build the value with the
	// corresponding ToMap, e.g. map[string]any{"rrf": RrfRetriever{...}.ToMap()}, or write the map
	// directly to set an attribute that the SDK does not model yet.
	Retriever            map[string]any `input:"body,retriever,json"`
	ReturnMetadata       *bool          `input:"body,returnMetadata,json"`
	ReturnMetadataFields []string       `input:"body,returnMetadataFields,json"`
	PartitionKeys        []string       `input:"body,partitionKeys,json"`
	Limit                *int           `input:"body,limit,json"`
	NextToken            *string        `input:"body,nextToken,json"`
	Sort                 []Sort         `input:"body,sort,json"`
	oss.RequestCommon
}

// Knn defines a single knn vector query of the QueryVectorsFusion operation.
type Knn struct {
	Field         *string  `json:"field,omitempty"`
	QueryVector   any      `json:"queryVector,omitempty"`
	TopK          *int     `json:"topK,omitempty"`
	Filter        any      `json:"filter,omitempty"`
	NumCandidates *int     `json:"numCandidates,omitempty"`
	Boost         *float32 `json:"boost,omitempty"`
}

// ToMap returns the knn query as the raw JSON object that the service expects, so that it can be
// used as an element of QueryVectorsFusionRequest.Knn. Attributes that the SDK does not model can
// be added to the returned map directly.
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
	// Each element is the raw JSON object of a retriever component, so an attribute that the SDK
	// does not model is passed through without an SDK change. Build an element with
	// RetrieverComponent.ToMap, or write the map directly.
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

// WeightRetriever defines the weight compound retriever that merges the results of the sub
// retrievers by weight.
type WeightRetriever struct {
	// The number of the top results taken from each sub retriever. Default value: 100.
	WindowSize *int32 `json:"windowSize,omitempty"`

	// The sub retrievers. It contains 1 to 3 elements.
	//
	// Each element is the raw JSON object of a retriever component, so an attribute that the SDK
	// does not model is passed through without an SDK change. Build an element with
	// RetrieverComponent.ToMap, or write the map directly.
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

// RetrieverComponent defines the component of a compound retriever (RrfRetriever or
// WeightRetriever). It wraps a sub retriever with the fusion weight, and optionally the score
// normalizer.
type RetrieverComponent struct {
	// The nested retriever as the raw JSON object. It can be a leaf retriever (knn/simple) or a
	// nested compound retriever (rrf/weight).
	//
	// Build it with the leaf ToMap wrapped by its type key, e.g.
	// map[string]any{"knn": Knn{...}.ToMap()}, or write the map directly.
	Retriever map[string]any `json:"retriever,omitempty"`

	// The weight of this component. It is a non-negative float32 number. Default value: 1.0.
	Weight *float32 `json:"weight,omitempty"`

	// The normalizer of the score. Valid values: none, minMax, l2. The values are also declared as
	// the NormalizerType constants, e.g. oss.Ptr(string(NormalizerTypeMinMax)).
	// It applies to the weight compound retriever only.
	Normalizer *string `json:"normalizer,omitempty"`
}

// ToMap returns the component as the raw JSON object of an element of RrfRetriever.Retrievers or
// WeightRetriever.Retrievers. Attributes that the SDK does not model can be added to the returned
// map directly.
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

type SortOptions struct {
	// The sort order of the field. Valid values: asc, desc.
	// The values are also declared as the SortOrderType constants, e.g. oss.Ptr(string(SortOrderTypeAsc)).
	Order *string `json:"order,omitempty"`
}

type Sort map[string]SortOptions

type QueryVectorsFusionResult struct {
	Vectors   []QueryVectorsFusionSummary `json:"vectors"`
	NextToken *string                     `json:"nextToken"`

	oss.ResultCommon
}

type QueryVectorsFusionSummary struct {
	Key      *string        `json:"key"`
	Metadata map[string]any `json:"metadata"`
	Score    *float32       `json:"score"`
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
