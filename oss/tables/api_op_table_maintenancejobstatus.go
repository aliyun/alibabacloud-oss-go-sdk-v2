package tables

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableMaintenanceJobStatusRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	oss.RequestCommon
}

type GetTableMaintenanceJobStatusResult struct {
	MaintenanceJobStatus *MaintenanceJobStatus `json:"status"`
	TableARN             *string               `json:"tableARN"`
	VersionToken         *string               `json:"versionToken"`

	oss.ResultCommon
}

type MaintenanceJobStatus struct {
	Job *MaintenanceJob `json:"job"`
}

type MaintenanceJob struct {
	FailureMessage   *string `json:"failureMessage"`
	LastRunTimestamp *string `json:"lastRunTimestamp"`
	Status           *string `json:"status"`
}

// GetTableMaintenanceJobStatus Queries the table maintenance job status of a table.
func (c *TablesClient) GetTableMaintenanceJobStatus(ctx context.Context, request *GetTableMaintenanceJobStatusRequest, optFns ...func(*oss.Options)) (*GetTableMaintenanceJobStatusResult, error) {
	var err error
	if request == nil {
		request = &GetTableMaintenanceJobStatusRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableMaintenanceJobStatus",
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
	result := &GetTableMaintenanceJobStatusResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type GetTableMaintenanceJobStatusByTableArnRequest struct {
	TableArn *string `input:"header,x-oss-table-arn,required"`

	oss.RequestCommon
}

// GetTableMaintenanceJobStatusByTableArn Queries the table maintenance job status of a table by table arn.
func (c *TablesClient) GetTableMaintenanceJobStatusByTableArn(ctx context.Context, request *GetTableMaintenanceJobStatusByTableArnRequest, optFns ...func(*oss.Options)) (*GetTableMaintenanceJobStatusResult, error) {
	var err error
	if request == nil {
		request = &GetTableMaintenanceJobStatusByTableArnRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableMaintenanceJobStatus",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"maintenance-job-status": "",
			"tables":                 "",
		},
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &GetTableMaintenanceJobStatusResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type SetTableMaintenanceJobStatusRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	//TablesOperation *string `input:"header,x-oss-tables-operation,required"`

	Status *MaintenanceJobStatus `input:"body,status,json,required"`

	VersionToken *string `input:"body,versionToken,json,required"`

	oss.RequestCommon
}

type SetTableMaintenanceJobStatusResult struct {
	VersionToken *string `json:"versionToken"`

	oss.ResultCommon
}

// SetTableMaintenanceJobStatus Set the table maintenance job status of a table by name.
func (c *TablesClient) SetTableMaintenanceJobStatus(ctx context.Context, request *SetTableMaintenanceJobStatusRequest, optFns ...func(*oss.Options)) (*SetTableMaintenanceJobStatusResult, error) {
	var err error
	if request == nil {
		request = &SetTableMaintenanceJobStatusRequest{}
	}
	input := &oss.OperationInput{
		OpName: "SetTableMaintenanceJobStatus",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
			"x-oss-tables-operation":  "",
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
	result := &SetTableMaintenanceJobStatusResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type SetTableMaintenanceJobStatusByTableArnRequest struct {
	TableArn *string `input:"header,x-oss-table-arn,required"`

	Status *MaintenanceJobStatus `input:"body,status,json,required"`

	VersionToken *string `input:"body,versionToken,json,required"`

	oss.RequestCommon
}

// SetTableMaintenanceJobStatusByTableArn Set the table maintenance job status of a table by table arn.
func (c *TablesClient) SetTableMaintenanceJobStatusByTableArn(ctx context.Context, request *SetTableMaintenanceJobStatusByTableArnRequest, optFns ...func(*oss.Options)) (*SetTableMaintenanceJobStatusResult, error) {
	var err error
	if request == nil {
		request = &SetTableMaintenanceJobStatusByTableArnRequest{}
	}
	input := &oss.OperationInput{
		OpName: "SetTableMaintenanceJobStatus",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
			"x-oss-tables-operation":  "",
		},
		Parameters: map[string]string{
			"maintenance-job-status": "",
			"tables":                 "",
		},
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &SetTableMaintenanceJobStatusResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
