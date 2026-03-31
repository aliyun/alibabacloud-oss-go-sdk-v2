package tables

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableBucketMaintenanceConfigurationRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

type GetTableBucketMaintenanceConfigurationResult struct {
	// The container that stores the maintenance configuration of the table bucket.
	Configuration *MaintenanceConfiguration `output:"body,configuration,json"`

	TableBucketARN *string `output:"body,tableBucketARN,json"`

	oss.ResultCommon
}

type IcebergUnreferencedFileRemoval struct {
	Settings *MaintenanceSettings `json:"settings"`

	Status *string `json:"status"`
}

type MaintenanceSettings struct {
	UnreferencedDays *int64 `json:"unreferencedDays"`
	NonCurrentDays   *int64 `json:"nonCurrentDays"`
}

// GetTableBucketMaintenanceConfiguration Queries the maintenance config of a bucket.
func (c *TablesClient) GetTableBucketMaintenanceConfiguration(ctx context.Context, request *GetTableBucketMaintenanceConfigurationRequest, optFns ...func(*oss.Options)) (*GetTableBucketMaintenanceConfigurationResult, error) {
	var err error
	if request == nil {
		request = &GetTableBucketMaintenanceConfigurationRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableBucketMaintenanceConfiguration",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"maintenance": "",
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
	result := &GetTableBucketMaintenanceConfigurationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type PutTableBucketMaintenanceConfigurationRequest struct {
	// This name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	// The container that stores the maintenance configuration of the table bucket.
	IcebergUnreferencedFileRemoval *IcebergUnreferencedFileRemoval `input:"body,icebergUnreferencedFileRemoval,json,required"`

	oss.RequestCommon
}

type MaintenanceConfiguration struct {
	IcebergUnreferencedFileRemoval *IcebergUnreferencedFileRemoval `json:"icebergUnreferencedFileRemoval"`
}

type PutTableBucketMaintenanceConfigurationResult struct {
	oss.ResultCommon
}

// PutTableBucketMaintenanceConfiguration set maintenance config for the table bucket.
func (c *TablesClient) PutTableBucketMaintenanceConfiguration(ctx context.Context, request *PutTableBucketMaintenanceConfigurationRequest, optFns ...func(*oss.Options)) (*PutTableBucketMaintenanceConfigurationResult, error) {
	var err error
	if request == nil {
		request = &PutTableBucketMaintenanceConfigurationRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutTableBucketMaintenanceConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"maintenance": "",
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

	result := &PutTableBucketMaintenanceConfigurationResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
