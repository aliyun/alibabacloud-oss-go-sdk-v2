package tables

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableBucketStorageClassRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

type GetTableBucketStorageClassResult struct {
	// The container that stores the storage class of the table bucket.
	StorageClassConfiguration *StorageClassConfiguration `output:"body,storageClassConfiguration,json"`

	oss.ResultCommon
}

// GetTableBucketStorageClass Queries the storage class of a bucket.
func (c *TablesClient) GetTableBucketStorageClass(ctx context.Context, request *GetTableBucketStorageClassRequest, optFns ...func(*oss.Options)) (*GetTableBucketStorageClassResult, error) {
	var err error
	if request == nil {
		request = &GetTableBucketStorageClassRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableBucketStorageClass",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &GetTableBucketStorageClassResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type PutTableBucketStorageClassRequest struct {
	// This name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	// The storage class of the table bucket.
	StorageClassConfiguration *StorageClassConfiguration `input:"body,storageClassConfiguration,json,required"`

	oss.RequestCommon
}

type PutTableBucketStorageClassResult struct {
	oss.ResultCommon
}

// PutTableBucketStorageClass set storage class for the table bucket.
func (c *TablesClient) PutTableBucketStorageClass(ctx context.Context, request *PutTableBucketStorageClassRequest, optFns ...func(*oss.Options)) (*PutTableBucketStorageClassResult, error) {
	var err error
	if request == nil {
		request = &PutTableBucketStorageClassRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutTableBucketStorageClass",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutTableBucketStorageClassResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
