package vectors

import (
	"context"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutVectorBucketRequest struct {
	// The name of the bucket to create.
	Bucket *string `input:"host,bucket,required"`

	// The ID of the resource group.
	ResourceGroupId *string `input:"header,x-oss-resource-group-id"`

	// The tagging of the bucket.
	Tagging *string `input:"header,x-oss-bucket-tagging"`

	oss.RequestCommon
}

type PutVectorBucketResult struct {
	oss.ResultCommon
}

// PutVectorBucket Creates a vector bucket.
func (c *VectorsClient) PutVectorBucket(ctx context.Context, request *PutVectorBucketRequest, optFns ...func(*oss.Options)) (*PutVectorBucketResult, error) {
	var err error
	if request == nil {
		request = &PutVectorBucketRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutVectorBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: "application/json",
		},
		Bucket: request.Bucket,
	}

	if err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutVectorBucketResult{}

	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetVectorBucketRequest struct {
	// The name of the bucket containing the objects
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

type GetVectorBucketResult struct {
	// The container that stores the bucket information.
	BucketInfo *VectorBucketInfo `output:"body,BucketInfo>Bucket,json"`

	oss.ResultCommon
}

// VectorBucketInfo defines Bucket information
type VectorBucketInfo struct {
	// The name of the bucket.
	Name *string `json:"Name"`

	// The region in which the bucket is located.
	Location *string `json:"Location"`

	// The time when the bucket is created. The time is in UTc.client.
	CreationDate *time.Time `json:"CreationDate"`

	// The public endpoint that is used to access the bucket over the Internet.
	ExtranetEndpoint *string `json:"ExtranetEndpoint"`

	// The internal endpoint that is used to access the bucket from Elastic
	IntranetEndpoint *string `json:"IntranetEndpoint"`

	// The ID of the resource group to which the bucket belongs.
	ResourceGroupId *string `json:"ResourceGroupId"`
}

// GetVectorBucket Queries information about a bucket.
func (c *VectorsClient) GetVectorBucket(ctx context.Context, request *GetVectorBucketRequest, optFns ...func(*oss.Options)) (*GetVectorBucketResult, error) {
	var err error
	if request == nil {
		request = &GetVectorBucketRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetVectorBucket",
		Method: "GET",
		Parameters: map[string]string{
			"bucketInfo": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetVectorBucketResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type ListVectorBucketsRequest struct {
	// The name of the bucket from which the list operation begins.
	Marker *string `input:"query,marker"`

	// The maximum number of buckets that can be returned in the single query.
	// Valid values: 1 to 1000.
	MaxKeys int32 `input:"query,max-keys"`

	// The prefix that the names of returned buckets must contain.
	Prefix *string `input:"query,prefix"` // Limits the response to keys that begin with the specified prefix

	// The ID of the resource group.
	ResourceGroupId *string `input:"header,x-oss-resource-group-id"`

	oss.RequestCommon
}

type ListVectorBucketsResult struct {
	// The prefix contained in the names of the returned bucket.
	Prefix *string `json:"Prefix"`

	// The name of the bucket after which the ListVectorBuckets  operation starts.
	Marker *string `json:"Marker"` // The marker filter.

	// The maximum number of buckets that can be returned for the request.
	MaxKeys int32 `json:"MaxKeys"`

	// Indicates whether all results are returned.
	// true: Only part of the results are returned for the request.
	// false: All results are returned for the request.
	IsTruncated bool `json:"IsTruncated"`

	// The marker for the next ListVectorBuckets request, which can be used to return the remaining results.
	NextMarker *string `json:"NextMarker"`

	// The container that stores information about buckets.
	Buckets Buckets `json:"Buckets"`

	oss.ResultCommon
}

// Buckets The container that stores information about buckets.
type Buckets struct {
	Bucket []VectorBucketProperties `json:"Bucket"`
}

type VectorBucketProperties struct {
	// The name of the bucket.
	Name *string `json:"Name"`

	// The data center in which the bucket is located.
	Location *string `json:"Location"`

	// The time when the bucket was created. Format: yyyy-mm-ddThh:mm:ss.timezone.
	CreationDate *time.Time `json:"CreationDate"`

	// The public endpoint used to access the bucket over the Internet.
	ExtranetEndpoint *string `json:"ExtranetEndpoint"`

	// The internal endpoint that is used to access the bucket from ECS instances
	// that reside in the same region as the bucket.
	IntranetEndpoint *string `json:"IntranetEndpoint"`

	// The region in which the bucket is located.
	Region *string `json:"Region"`

	// The ID of the resource group to which the bucket belongs.
	ResourceGroupId *string `json:"ResourceGroupId"`
}

// ListVectorBuckets Lists vector buckets that belong to the current account.
func (c *VectorsClient) ListVectorBuckets(ctx context.Context, request *ListVectorBucketsRequest, optFns ...func(*oss.Options)) (*ListVectorBucketsResult, error) {
	var err error
	if request == nil {
		request = &ListVectorBucketsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "ListVectorBuckets",
		Method: "GET",
	}
	if err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &ListVectorBucketsResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson2); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteVectorBucketRequest struct {
	// The name of the bucket to delete.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

type DeleteVectorBucketResult struct {
	oss.ResultCommon
}

// DeleteVectorBucket Deletes a vector bucket.
func (c *VectorsClient) DeleteVectorBucket(ctx context.Context, request *DeleteVectorBucketRequest, optFns ...func(*oss.Options)) (*DeleteVectorBucketResult, error) {
	var err error
	if request == nil {
		request = &DeleteVectorBucketRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteVectorBucket",
		Method: "DELETE",
		Bucket: request.Bucket,
	}
	if err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteVectorBucketResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
