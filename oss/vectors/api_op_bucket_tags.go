package vectors

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutBucketTagsRequest = oss.PutBucketTagsRequest
type PutBucketTagsResult = oss.PutBucketTagsResult
type Tagging = oss.Tagging
type Tag = oss.Tag
type TagSet = oss.TagSet

// PutBucketTags Adds tags to or modifies the existing tags of a vector bucket.
func (c *VectorsClient) PutBucketTags(ctx context.Context, request *PutBucketTagsRequest, optFns ...func(*oss.Options)) (*PutBucketTagsResult, error) {
	var err error
	if request == nil {
		request = &PutBucketTagsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutBucketTags",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tagging": "",
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
	result := &PutBucketTagsResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetBucketTagsRequest = oss.GetBucketTagsRequest
type GetBucketTagsResult = oss.GetBucketTagsResult

// GetBucketTags Queries the tags of a vector bucket.
func (c *VectorsClient) GetBucketTags(ctx context.Context, request *GetBucketTagsRequest, optFns ...func(*oss.Options)) (*GetBucketTagsResult, error) {
	var err error
	if request == nil {
		request = &GetBucketTagsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetBucketTags",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tagging": "",
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
	result := &GetBucketTagsResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteBucketTagsRequest = oss.DeleteBucketTagsRequest
type DeleteBucketTagsResult = oss.DeleteBucketTagsResult

// DeleteBucketTags Deletes tags configured for a vector bucket.
func (c *VectorsClient) DeleteBucketTags(ctx context.Context, request *DeleteBucketTagsRequest, optFns ...func(*oss.Options)) (*DeleteBucketTagsResult, error) {
	var err error
	if request == nil {
		request = &DeleteBucketTagsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteBucketTags",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tagging": "",
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
	result := &DeleteBucketTagsResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
