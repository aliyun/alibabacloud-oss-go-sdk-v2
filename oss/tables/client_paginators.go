package tables

import (
	"context"
	"fmt"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// ListTableBucketsPaginator is a paginator for ListTableBuckets
type ListTableBucketsPaginator struct {
	options     oss.PaginatorOptions
	client      *TablesClient
	request     *ListTableBucketsRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *TablesClient) NewListTableBucketsPaginator(request *ListTableBucketsRequest, optFns ...func(*oss.PaginatorOptions)) *ListTableBucketsPaginator {
	if request == nil {
		request = &ListTableBucketsRequest{}
	}

	options := oss.PaginatorOptions{}
	options.Limit = request.MaxBuckets

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListTableBucketsPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.ContinuationToken,
		firstPage:   true,
		isTruncated: false,
	}
}

// HasNext Returns true if there’s a next page.
func (p *ListTableBucketsPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListTableBuckets page.
func (p *ListTableBucketsPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListTableBucketsResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.ContinuationToken = p.nextToken

	var limit int32
	if p.options.Limit > 0 {
		limit = p.options.Limit
	}
	request.MaxBuckets = limit

	result, err := p.client.ListTableBuckets(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.nextToken = result.ContinuationToken
	p.isTruncated = oss.ToString(p.nextToken) != ""

	return result, nil
}

// ListNamespacesPaginator is a paginator for ListNamespaces.
type ListNamespacesPaginator struct {
	options     oss.PaginatorOptions
	client      *TablesClient
	request     *ListNamespacesRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *TablesClient) NewListNameSpacesPaginator(request *ListNamespacesRequest, optFns ...func(*oss.PaginatorOptions)) *ListNamespacesPaginator {
	if request == nil {
		request = &ListNamespacesRequest{}
	}

	options := oss.PaginatorOptions{}
	options.Limit = request.MaxNamespaces

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListNamespacesPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.ContinuationToken,
		firstPage:   true,
		isTruncated: false,
	}
}

// HasNext Returns true if there’s a next page.
func (p *ListNamespacesPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListNamespaces page.
func (p *ListNamespacesPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListNamespacesResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.ContinuationToken = p.nextToken

	var limit int32
	if p.options.Limit > 0 {
		limit = p.options.Limit
	}
	request.MaxNamespaces = limit

	result, err := p.client.ListNamespaces(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.nextToken = result.ContinuationToken
	p.isTruncated = oss.ToString(p.nextToken) != ""

	return result, nil
}

// ListTablesPaginator is a paginator for ListTables.
type ListTablesPaginator struct {
	options     oss.PaginatorOptions
	client      *TablesClient
	request     *ListTablesRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *TablesClient) NewListTablesPaginator(request *ListTablesRequest, optFns ...func(*oss.PaginatorOptions)) *ListTablesPaginator {
	if request == nil {
		request = &ListTablesRequest{}
	}

	options := oss.PaginatorOptions{}
	options.Limit = request.MaxTables

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListTablesPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.ContinuationToken,
		firstPage:   true,
		isTruncated: false,
	}
}

// HasNext Returns true if there’s a next page.
func (p *ListTablesPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListTables page.
func (p *ListTablesPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListTablesResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.ContinuationToken = p.nextToken

	var limit int32
	if p.options.Limit > 0 {
		limit = p.options.Limit
	}
	request.MaxTables = limit

	result, err := p.client.ListTables(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.nextToken = result.ContinuationToken
	p.isTruncated = oss.ToString(p.nextToken) != ""

	return result, nil
}
