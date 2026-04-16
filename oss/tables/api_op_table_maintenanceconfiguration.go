package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableMaintenanceConfigurationRequest struct {
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	Namespace *string `input:"nop,namespace,required"`

	Name *string `input:"nop,name,required"`

	TableARN *string `input:"header,x-oss-table-arn"`

	oss.RequestCommon
}

type GetTableMaintenanceConfigurationResult struct {
	// The container that stores the maintenance configuration of the table bucket.
	Configuration *TableMaintenanceConfiguration `output:"body,configuration,json"`

	TableARN *string `output:"body,tableARN,json"`

	oss.ResultCommon
}

type TableMaintenanceConfiguration struct {
	IcebergCompaction *IcebergCompaction `json:"icebergCompaction"`

	IcebergSnapshotManagement *IcebergSnapshotManagement `json:"icebergSnapshotManagement"`
}

type IcebergCompaction struct {
	Settings *IcebergCompactionSettings `json:"settings"`

	Status *string `json:"status"`
}

type IcebergCompactionSettings struct {
	IcebergCompaction *IcebergCompactionSettingsDetail `json:"icebergCompaction"`
}

type IcebergCompactionSettingsDetail struct {
	Strategy *string `json:"strategy"`

	TargetFileSizeMB *int `json:"targetFileSizeMB"`
}

type IcebergSnapshotManagement struct {
	Settings *IcebergSnapshotManagementSettings `json:"settings"`

	Status *string `json:"status"`
}

type IcebergSnapshotManagementSettings struct {
	IcebergSnapshotManagement *IcebergSnapshotManagementSettingsDetail `json:"icebergSnapshotManagement"`
}

type IcebergSnapshotManagementSettingsDetail struct {
	MaxSnapshotAgeHours *int `json:"maxSnapshotAgeHours"`

	MinSnapshotsToKeep *int `json:"minSnapshotsToKeep"`
}

// GetTableMaintenanceConfiguration Queries the maintenance config of a table.
func (c *TablesClient) GetTableMaintenanceConfiguration(ctx context.Context, request *GetTableMaintenanceConfigurationRequest, optFns ...func(*oss.Options)) (*GetTableMaintenanceConfigurationResult, error) {
	var err error
	if request == nil {
		request = &GetTableMaintenanceConfigurationRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableMaintenanceConfiguration",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &GetTableMaintenanceConfigurationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type PutTableMaintenanceConfigurationRequest struct {
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	Namespace *string `input:"nop,namespace,required"`

	Name *string `input:"nop,name,required"`

	Type *string `input:"nop,type,required"`

	// The container that stores the maintenance configuration of the table.
	Value *TableMaintenanceValue `input:"body,value,json,required"`

	oss.RequestCommon
}

type TableMaintenanceValue struct {
	Settings *TableMaintenanceSettings `json:"settings"`

	Status *string `json:"status"`
}

type TableMaintenanceSettings struct {
	IcebergCompaction *IcebergCompactionSettingsDetail `json:"icebergCompaction,omitempty"`

	IcebergSnapshotManagement *IcebergSnapshotManagementSettingsDetail `json:"icebergSnapshotManagement,omitempty"`
}

type PutTableMaintenanceConfigurationResult struct {
	oss.ResultCommon
}

// PutTableMaintenanceConfiguration set maintenance config for the table.
func (c *TablesClient) PutTableMaintenanceConfiguration(ctx context.Context, request *PutTableMaintenanceConfigurationRequest, optFns ...func(*oss.Options)) (*PutTableMaintenanceConfigurationResult, error) {
	var err error
	if request == nil {
		request = &PutTableMaintenanceConfigurationRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutTableMaintenanceConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)), url.QueryEscape(oss.ToString(request.Type)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutTableMaintenanceConfigurationResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
