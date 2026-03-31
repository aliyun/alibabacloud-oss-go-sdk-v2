package tables

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_SetTableMaintenanceJobStatus(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *SetTableMaintenanceJobStatusRequest
	var input *oss.OperationInput
	var err error

	request = &SetTableMaintenanceJobStatusRequest{}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &SetTableMaintenanceJobStatusRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &SetTableMaintenanceJobStatusRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &SetTableMaintenanceJobStatusRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Status.")

	request = &SetTableMaintenanceJobStatusRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
		Status: &MaintenanceJobStatus{
			Job: &MaintenanceJob{
				FailureMessage:   oss.Ptr("no message"),
				LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
				Status:           oss.Ptr("success"),
			},
		},
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, VersionToken.")

	request = &SetTableMaintenanceJobStatusRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
		Status: &MaintenanceJobStatus{
			Job: &MaintenanceJob{
				FailureMessage:   oss.Ptr("no message"),
				LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
				Status:           oss.Ptr("success"),
			},
		},
		VersionToken: oss.Ptr("token"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, input.Headers["x-oss-tables-operation"], "")
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["maintenance-job-status"], "")
	assert.Equal(t, input.Parameters["space"], "")
	assert.Equal(t, input.Parameters["table"], "")
	body, err := io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(body), "{\"status\":{\"job\":{\"failureMessage\":\"no message\",\"lastRunTimestamp\":\"2026-02-31T10:56:21.000Z\",\"status\":\"success\"}},\"versionToken\":\"token\"}")
}

func TestMarshalInput_SetTableMaintenanceJobStatusByTableArn(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *SetTableMaintenanceJobStatusByTableArnRequest
	var input *oss.OperationInput
	var err error

	request = &SetTableMaintenanceJobStatusByTableArnRequest{}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableArn.")

	request = &SetTableMaintenanceJobStatusByTableArnRequest{
		TableArn: oss.Ptr("acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Status.")

	request = &SetTableMaintenanceJobStatusByTableArnRequest{
		TableArn: oss.Ptr("acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123"),
		Status: &MaintenanceJobStatus{
			Job: &MaintenanceJob{
				FailureMessage:   oss.Ptr("no message"),
				LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
				Status:           oss.Ptr("success"),
			},
		},
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, VersionToken.")

	request = &SetTableMaintenanceJobStatusByTableArnRequest{
		TableArn: oss.Ptr("acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123"),
		Status: &MaintenanceJobStatus{
			Job: &MaintenanceJob{
				FailureMessage:   oss.Ptr("no message"),
				LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
				Status:           oss.Ptr("success"),
			},
		},
		VersionToken: oss.Ptr("token"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, input.Headers["x-oss-table-arn"], "acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123")
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["maintenance-job-status"], "")
	assert.Equal(t, input.Parameters["space"], "")
	assert.Equal(t, input.Parameters["table"], "")
	body, err := io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(body), "{\"status\":{\"job\":{\"failureMessage\":\"no message\",\"lastRunTimestamp\":\"2026-02-31T10:56:21.000Z\",\"status\":\"success\"}},\"versionToken\":\"token\"}")
}

func TestUnmarshalOutput_SetTableMaintenanceJobStatus(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "versionToken": "aaa"
}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &SetTableMaintenanceJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.VersionToken, "aaa")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &SetTableMaintenanceJobStatusResult{}
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
	result = &SetTableMaintenanceJobStatusResult{}
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
	result = &SetTableMaintenanceJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

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
		Parameters: map[string]string{
			"maintenance-job-status":        "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetTableMaintenanceJobStatusRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &GetTableMaintenanceJobStatusRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &GetTableMaintenanceJobStatusRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "GET")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["maintenance-job-status"], "")
	assert.Equal(t, input.Parameters["space"], "")
	assert.Equal(t, input.Parameters["table"], "")
	assert.Equal(t, input.Headers["x-oss-table-arn"], "table-arn")
}

func TestMarshalInput_GetTableMaintenanceJobStatusByTableArn(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetTableMaintenanceJobStatusByTableArnRequest
	var input *oss.OperationInput
	var err error

	request = &GetTableMaintenanceJobStatusByTableArnRequest{}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableArn.")

	request = &GetTableMaintenanceJobStatusByTableArnRequest{
		TableArn: oss.Ptr("table-arn"),
	}
	input = &oss.OperationInput{
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
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "GET")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["maintenance-job-status"], "")
	assert.Equal(t, input.Headers["x-oss-table-arn"], "table-arn")
}

func TestUnmarshalOutput_GetTableMaintenanceJobStatus(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "status": { 
      "job" : { 
         "failureMessage": "no message",
         "lastRunTimestamp": "2026-02-31T10:56:21.000Z",
         "status": "success"
      }
   },
   "versionToken": "aaa",
   "tableARN": "test-arn"
}`
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
	assert.Equal(t, *result.MaintenanceJobStatus.Job.FailureMessage, "no message")
	assert.Equal(t, *result.MaintenanceJobStatus.Job.Status, "success")
	assert.Equal(t, *result.MaintenanceJobStatus.Job.LastRunTimestamp, "2026-02-31T10:56:21.000Z")
	assert.Equal(t, *result.VersionToken, "aaa")
	assert.Equal(t, *result.TableARN, "test-arn")

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
