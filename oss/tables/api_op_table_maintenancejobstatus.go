package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableMaintenanceJobStatusRequest struct {
	BucketArn *string `input:"nop,bucketArn,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	TableArn *string `input:"header,x-oss-table-arn"`

	oss.RequestCommon
}

type GetTableMaintenanceJobStatusResult struct {
	JobStatus *MaintenanceJobStatus `json:"status"`
	TableArn  *string               `json:"tableARN"`

	oss.ResultCommon
}

type MaintenanceJobStatus struct {
	IcebergCompaction              *StatusDetail `json:"icebergCompaction,omitempty"`
	IcebergSnapshotManagement      *StatusDetail `json:"icebergSnapshotManagement,omitempty"`
	IcebergUnreferencedFileRemoval *StatusDetail `json:"icebergUnreferencedFileRemoval,omitempty"`
}

type StatusDetail struct {
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
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance-job-status", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
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
