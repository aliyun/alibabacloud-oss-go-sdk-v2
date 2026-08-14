package agentic

import (
	"context"
	"fmt"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// ListAgenticBucketsPaginator is a paginator for ListAgenticBuckets
type ListAgenticBucketsPaginator struct {
	options     oss.PaginatorOptions
	client      *AgenticBucketClient
	request     *ListAgenticBucketsRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *AgenticBucketClient) NewListAgenticBucketsPaginator(request *ListAgenticBucketsRequest, optFns ...func(*oss.PaginatorOptions)) *ListAgenticBucketsPaginator {
	if request == nil {
		request = &ListAgenticBucketsRequest{}
	}

	options := oss.PaginatorOptions{}
	if request.MaxKeys != nil {
		options.Limit = int32(*request.MaxKeys)
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListAgenticBucketsPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.ContinuationToken,
		firstPage:   true,
		isTruncated: false,
	}
}

func (p *ListAgenticBucketsPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

func (p *ListAgenticBucketsPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListAgenticBucketsResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.ContinuationToken = p.nextToken

	if p.options.Limit > 0 {
		limit := int(p.options.Limit)
		request.MaxKeys = &limit
	}

	result, err := p.client.ListAgenticBuckets(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.isTruncated = oss.ToBool(result.IsTruncated)
	p.nextToken = result.NextContinuationToken

	return result, nil
}

// ListBucketSpacesPaginator is a paginator for ListBucketSpaces
type ListBucketSpacesPaginator struct {
	options     oss.PaginatorOptions
	client      *AgenticBucketClient
	request     *ListBucketSpacesRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *AgenticBucketClient) NewListBucketSpacesPaginator(request *ListBucketSpacesRequest, optFns ...func(*oss.PaginatorOptions)) *ListBucketSpacesPaginator {
	if request == nil {
		request = &ListBucketSpacesRequest{}
	}

	options := oss.PaginatorOptions{}
	if request.MaxKeys != nil {
		options.Limit = int32(*request.MaxKeys)
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListBucketSpacesPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.ContinuationToken,
		firstPage:   true,
		isTruncated: false,
	}
}

func (p *ListBucketSpacesPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

func (p *ListBucketSpacesPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListBucketSpacesResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.ContinuationToken = p.nextToken

	if p.options.Limit > 0 {
		limit := int(p.options.Limit)
		request.MaxKeys = &limit
	}

	result, err := p.client.ListBucketSpaces(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.isTruncated = oss.ToBool(result.IsTruncated)
	p.nextToken = result.NextContinuationToken

	return result, nil
}
