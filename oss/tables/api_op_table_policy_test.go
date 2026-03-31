package tables

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutTablePolicy(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *PutTablePolicyRequest
	var input *oss.OperationInput
	var err error

	request = &PutTablePolicyRequest{}
	input = &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutTablePolicyRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	request = &PutTablePolicyRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table")

	request = &PutTablePolicyRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Body")

	request = &PutTablePolicyRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("xfz-table-bucket"),
		Body:      strings.NewReader(`{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/xfz-table-bucket/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`),
	}
	input = &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["policy"], "")
	assert.Equal(t, input.Parameters["space"], "")
	assert.Equal(t, input.Parameters["table"], "")
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	jsonStr, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), `{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/xfz-table-bucket/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`)
}

func TestUnmarshalOutput_PutTablePolicy(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutTablePolicyResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	body := `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	resultErr := &PutTablePolicyResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetTablePolicy(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetTablePolicyRequest
	var input *oss.OperationInput
	var err error

	request = &GetTablePolicyRequest{}
	input = &oss.OperationInput{
		OpName: "GetTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetTablePolicyRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &GetTablePolicyRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "GetTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &GetTablePolicyRequest{
		Bucket:    oss.Ptr("bucket"),
		Table:     oss.Ptr("table"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "GetTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["policy"], "")
	assert.Equal(t, input.Parameters["space"], "")
	assert.Equal(t, input.Parameters["table"], "")
}

func TestUnmarshalOutput_GetTablePolicy(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/xfz-table-bucket/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	policy, err := io.ReadAll(output.Body)
	assert.Nil(t, err)
	defer output.Body.Close()
	result := &GetTablePolicyResult{
		Body: string(policy),
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Body, `{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/xfz-table-bucket/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`)

	body = `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	policy, err = io.ReadAll(output.Body)
	assert.Nil(t, err)
	result = &GetTablePolicyResult{
		Body: string(policy),
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteTablePolicy(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *DeleteTablePolicyRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteTablePolicyRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteTablePolicyRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &DeleteTablePolicyRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &DeleteTablePolicyRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"policy":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["policy"], "")
	assert.Equal(t, input.Parameters["space"], "")
	assert.Equal(t, input.Parameters["table"], "")
}

func TestUnmarshalOutput_DeleteTablePolicy(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	output = &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteTablePolicyResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	resultErr := &DeleteTablePolicyResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}
