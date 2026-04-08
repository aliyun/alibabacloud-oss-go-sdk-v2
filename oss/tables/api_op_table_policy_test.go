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
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn")

	request = &PutTablePolicyRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	request = &PutTablePolicyRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table")

	request = &PutTablePolicyRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, ResourcePolicy.")

	request = &PutTablePolicyRequest{
		BucketArn:      oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
		Table:          oss.Ptr("xfz-table-bucket"),
		ResourcePolicy: oss.Ptr("{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"),
	}
	input = &oss.OperationInput{
		OpName: "PutTablePolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)

	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/xfz-table-bucket/policy")
	jsonStr, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), `{"resourcePolicy":"{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"}`)
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

	body := `{"message": "AccessDenied"}`
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
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn.")

	request = &GetTablePolicyRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &GetTablePolicyRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "GetTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &GetTablePolicyRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Table:     oss.Ptr("table"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "GetTablePolicy",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy")
}

func TestUnmarshalOutput_GetTablePolicy(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{"resourcePolicy":"{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &GetTablePolicyResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.ResourcePolicy, `{"Version":"1","Statement":[{"Action":["oss:GetTable"],"Effect":"Allow","Principal":["9876543210"],"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*"]}]}`)

	body = `{
    "message": "AccessDenied"
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
	result = &GetTablePolicyResult{}
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
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn.")

	request = &DeleteTablePolicyRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTablePolicy",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &DeleteTablePolicyRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTablePolicy",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &DeleteTablePolicyRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTablePolicy",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/policy", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace), oss.ToString(request.Table))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method,"DELETE")
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy")
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
