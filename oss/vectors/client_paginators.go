package vectors

import (
	"context"
	"fmt"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// ListVectorBucketsPaginator is a paginator for ListVectorBuckets
type ListVectorBucketsPaginator struct {
	options     oss.PaginatorOptions
	client      *VectorsClient
	request     *ListVectorBucketsRequest
	marker      *string
	firstPage   bool
	isTruncated bool
}

func (c *VectorsClient) NewListVectorBucketsPaginator(request *ListVectorBucketsRequest, optFns ...func(*oss.PaginatorOptions)) *ListVectorBucketsPaginator {
	if request == nil {
		request = &ListVectorBucketsRequest{}
	}

	options := oss.PaginatorOptions{}
	options.Limit = request.MaxKeys

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListVectorBucketsPaginator{
		options:     options,
		client:      c,
		request:     request,
		marker:      request.Marker,
		firstPage:   true,
		isTruncated: false,
	}
}

// HasNext Returns true if there’s a next page.
func (p *ListVectorBucketsPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListVectorBuckets page.
func (p *ListVectorBucketsPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListVectorBucketsResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.Marker = p.marker

	var limit int32
	if p.options.Limit > 0 {
		limit = p.options.Limit
	}
	request.MaxKeys = limit

	result, err := p.client.ListVectorBuckets(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.isTruncated = result.IsTruncated
	p.marker = result.NextMarker

	return result, nil
}

// ListVectorIndexesPaginator is a paginator for ListVectorIndexes
type ListVectorIndexesPaginator struct {
	options     oss.PaginatorOptions
	client      *VectorsClient
	request     *ListVectorIndexesRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *VectorsClient) NewListVectorIndexesPaginator(request *ListVectorIndexesRequest, optFns ...func(*oss.PaginatorOptions)) *ListVectorIndexesPaginator {
	if request == nil {
		request = &ListVectorIndexesRequest{}
	}

	options := oss.PaginatorOptions{}

	options.Limit = int32(request.MaxResults)

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListVectorIndexesPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.NextToken,
		firstPage:   true,
		isTruncated: false,
	}
}

// HasNext Returns true if there’s a next page.
func (p *ListVectorIndexesPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListVectorIndexes page.
func (p *ListVectorIndexesPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListVectorIndexesResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.NextToken = p.nextToken

	var limit int32
	if p.options.Limit > 0 {
		limit = p.options.Limit
	}
	request.MaxResults = int(limit)

	result, err := p.client.ListVectorIndexes(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.nextToken = result.NextToken
	p.isTruncated = oss.ToString(p.nextToken) != ""

	return result, nil
}

// ListVectorsPaginator is a paginator for ListVectors
type ListVectorsPaginator struct {
	options     oss.PaginatorOptions
	client      *VectorsClient
	request     *ListVectorsRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *VectorsClient) NewListVectorsPaginator(request *ListVectorsRequest, optFns ...func(*oss.PaginatorOptions)) *ListVectorsPaginator {
	if request == nil {
		request = &ListVectorsRequest{}
	}

	options := oss.PaginatorOptions{}
	options.Limit = int32(request.MaxResults)

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListVectorsPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.NextToken,
		firstPage:   true,
		isTruncated: true,
	}
}

// HasNext Returns true if there’s a next page.
func (p *ListVectorsPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListVectors page.
func (p *ListVectorsPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListVectorsResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.NextToken = p.nextToken

	var limit int32
	if p.options.Limit > 0 {
		limit = p.options.Limit
	}
	request.MaxResults = int(limit)

	result, err := p.client.ListVectors(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.nextToken = result.NextToken
	p.isTruncated = oss.ToString(p.nextToken) != ""

	return result, nil
}
