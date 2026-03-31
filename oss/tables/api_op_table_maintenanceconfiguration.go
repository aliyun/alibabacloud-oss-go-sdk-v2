package tables

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableMaintenanceConfigurationRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	oss.RequestCommon
}

type GetTableMaintenanceConfigurationResult struct {
	// The container that stores the maintenance configuration of the table bucket.
	Configuration *MaintenanceConfiguration `output:"body,configuration,json"`

	TableARN *string `output:"body,TableARN,json"`

	oss.ResultCommon
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
		Parameters: map[string]string{
			"maintenance-job-status":        "",
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
	result := &GetTableMaintenanceConfigurationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type PutTableMaintenanceConfigurationRequest struct {
	// This name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	// The container that stores the maintenance configuration of the table bucket.
	IcebergUnreferencedFileRemoval *IcebergUnreferencedFileRemoval `input:"body,icebergUnreferencedFileRemoval,json,required"`

	oss.RequestCommon
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
		Parameters: map[string]string{
			"maintenance-job-status":        "",
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

	result := &PutTableMaintenanceConfigurationResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
