package tables

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_GetTableMaintenanceJobStatus(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetTableMaintenanceJobStatusRequest
	var input *oss.OperationInput
	var err error

	request = &GetTableMaintenanceJobStatusRequest{}
	input = &oss.OperationInput{
		OpName: "GetTableMaintenanceJobStatus",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance-job-status", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn.")

	request = &GetTableMaintenanceJobStatusRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableMaintenanceJobStatus",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance-job-status", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &GetTableMaintenanceJobStatusRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableMaintenanceJobStatus",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance-job-status", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &GetTableMaintenanceJobStatusRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
		TableArn:  oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/demo-table/table-123456"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableMaintenanceJobStatus",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance-job-status", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "GET")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance-job-status")
	assert.Equal(t, input.Headers["x-oss-table-arn"], "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/demo-table/table-123456")
}

func TestUnmarshalOutput_GetTableMaintenanceJobStatus(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
    "status": {
        "icebergCompaction": {
            "failureMessage": "internal error",
            "lastRunTimestamp": "2026-04-08T08:37:07.988892Z",
            "status": "Failed"},
        "icebergSnapshotManagement": {
            "failureMessage": "internal error",
            "lastRunTimestamp": "2026-04-08T08:36:07.426846Z",
            "status": "Failed"},
        "icebergUnreferencedFileRemoval": {"status": "Disabled"}},
    "tableARN": "acs:osstables:cn-beijing:123456:bucket/demo-bucket/table/eb998f10-d20c-4f22-9a76-ed64e9668f56"}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &GetTableMaintenanceJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.JobStatus.IcebergCompaction.FailureMessage, "internal error")
	assert.Equal(t, *result.JobStatus.IcebergCompaction.Status, "Failed")
	assert.Equal(t, *result.JobStatus.IcebergCompaction.LastRunTimestamp, "2026-04-08T08:37:07.988892Z")
	assert.Equal(t, *result.JobStatus.IcebergSnapshotManagement.FailureMessage, "internal error")
	assert.Equal(t, *result.JobStatus.IcebergSnapshotManagement.Status, "Failed")
	assert.Equal(t, *result.JobStatus.IcebergSnapshotManagement.LastRunTimestamp, "2026-04-08T08:36:07.426846Z")
	assert.Equal(t, *result.JobStatus.IcebergUnreferencedFileRemoval.Status, "Disabled")
	assert.Equal(t, *result.TableArn, "acs:osstables:cn-beijing:123456:bucket/demo-bucket/table/eb998f10-d20c-4f22-9a76-ed64e9668f56")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetTableMaintenanceJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetTableMaintenanceJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetTableMaintenanceJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
