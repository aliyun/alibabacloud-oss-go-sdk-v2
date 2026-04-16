package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableBucketMaintenanceConfigurationRequest struct {
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	oss.RequestCommon
}

type GetTableBucketMaintenanceConfigurationResult struct {
	// The container that stores the maintenance configuration of the table bucket.
	Configuration *MaintenanceConfiguration `output:"body,configuration,json"`

	TableBucketARN *string `output:"body,tableBucketARN,json"`

	oss.ResultCommon
}

type MaintenanceConfiguration struct {
	IcebergUnreferencedFileRemoval *IcebergUnreferencedFileRemoval `json:"icebergUnreferencedFileRemoval"`
}

type IcebergUnreferencedFileRemoval struct {
	Settings *IcebergSettings `json:"settings"`

	Status *string `json:"status"`
}

type IcebergSettings struct {
	IcebergUnreferencedFileRemoval *SettingsDetail `json:"icebergUnreferencedFileRemoval"`
}

type SettingsDetail struct {
	NonCurrentDays   *int `json:"nonCurrentDays,omitempty"`
	UnreferencedDays *int `json:"unreferencedDays,omitempty"`
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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/maintenance", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
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
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	Type *string `input:"nop,type,required"`

	// The container that stores the maintenance configuration of the table bucket.
	Value *MaintenanceValue `input:"body,value,json,required"`

	oss.RequestCommon
}

type MaintenanceValue struct {
	Settings *MaintenanceSettings `json:"settings,omitempty"`

	Status *string `json:"status,omitempty"`
}

type MaintenanceSettings struct {
	IcebergUnreferencedFileRemoval *SettingsDetail `json:"icebergUnreferencedFileRemoval,omitempty"`
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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/maintenance/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Type)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
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
