package tables

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableStorageClassRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	oss.RequestCommon
}

type GetTableStorageClassResult struct {
	// The container that stores the storage class of the table bucket.
	StorageClassConfiguration *StorageClassConfiguration `output:"body,storageClassConfiguration,json"`

	oss.ResultCommon
}

// GetTableStorageClass Queries the storage class of a table.
func (c *TablesClient) GetTableStorageClass(ctx context.Context, request *GetTableStorageClassRequest, optFns ...func(*oss.Options)) (*GetTableStorageClassResult, error) {
	var err error
	if request == nil {
		request = &GetTableStorageClassRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableStorageClass",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class":                 "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
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
	result := &GetTableStorageClassResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
