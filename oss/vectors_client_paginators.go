package oss

import (
	"context"
	"fmt"
)

// ListVectorBucketsPaginator is a paginator for ListVectorBuckets
type ListVectorBucketsPaginator struct {
	options     PaginatorOptions
	client      *VectorsClient
	request     *ListVectorBucketsRequest
	marker      *string
	firstPage   bool
	isTruncated bool
}

func (c *VectorsClient) NewListVectorBucketsPaginator(request *ListVectorBucketsRequest, optFns ...func(*PaginatorOptions)) *ListVectorBucketsPaginator {
	if request == nil {
		request = &ListVectorBucketsRequest{}
	}

	options := PaginatorOptions{}
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

// NextPage retrieves the next ListBuckets page.
func (p *ListVectorBucketsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListVectorBucketsResult, error) {
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
