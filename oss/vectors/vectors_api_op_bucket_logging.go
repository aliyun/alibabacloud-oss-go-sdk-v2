package vectors

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutBucketLoggingRequest = oss.PutBucketLoggingRequest
type BucketLoggingStatus = oss.BucketLoggingStatus
type LoggingEnabled = oss.LoggingEnabled
type PutBucketLoggingResult = oss.PutBucketLoggingResult

// PutBucketLogging Enables logging for a vector bucket.
func (c *VectorsClient) PutBucketLogging(ctx context.Context, request *PutBucketLoggingRequest, optFns ...func(*oss.Options)) (*PutBucketLoggingResult, error) {
	var err error
	if request == nil {
		request = &PutBucketLoggingRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutBucketLogging",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"logging": "",
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

	result := &PutBucketLoggingResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type GetBucketLoggingRequest = oss.GetBucketLoggingRequest
type GetBucketLoggingResult = oss.GetBucketLoggingResult

// GetBucketLogging Queries the configurations of access log collection of a vector bucket.
func (c *VectorsClient) GetBucketLogging(ctx context.Context, request *GetBucketLoggingRequest, optFns ...func(*oss.Options)) (*GetBucketLoggingResult, error) {
	var err error
	if request == nil {
		request = &GetBucketLoggingRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetBucketLogging",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"logging": "",
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

	result := &GetBucketLoggingResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteBucketLoggingRequest = oss.DeleteBucketLoggingRequest
type DeleteBucketLoggingResult = oss.DeleteBucketLoggingResult

// DeleteBucketLogging Disables the logging feature for a vector bucket.
func (c *VectorsClient) DeleteBucketLogging(ctx context.Context, request *DeleteBucketLoggingRequest, optFns ...func(*oss.Options)) (*DeleteBucketLoggingResult, error) {
	var err error
	if request == nil {
		request = &DeleteBucketLoggingRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteBucketLogging",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"logging": "",
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

	result := &DeleteBucketLoggingResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
