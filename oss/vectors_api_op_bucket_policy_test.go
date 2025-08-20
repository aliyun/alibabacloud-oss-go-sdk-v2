package oss

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketPolicy_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutBucketPolicyRequest
	var input *OperationInput
	var err error

	request = &PutBucketPolicyRequest{}
	input = &OperationInput{
		OpName: "PutBucketPolicy",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutBucketPolicyRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucketPolicy",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, Body.")
	putBody := `{"Version":"1","Statement":[{"Action":["ossvector:PutVectors","ossvector:GetVectors"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:ossvector:cn-hangzhou:1234567890:*"]}]}`
	request = &PutBucketPolicyRequest{
		Bucket: Ptr("oss-demo"),
		Body:   strings.NewReader(putBody),
	}
	input = &OperationInput{
		OpName: "PutBucketPolicy",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"ossvector:PutVectors\",\"ossvector:GetVectors\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:ossvector:cn-hangzhou:1234567890:*\"]}]}")
}

func TestUnmarshalOutput_PutBucketPolicy_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutBucketPolicyResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutBucketPolicyResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutBucketPolicyResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutBucketPolicyResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetBucketPolicy_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *GetBucketPolicyRequest
	var input *OperationInput
	var err error

	request = &GetBucketPolicyRequest{}
	input = &OperationInput{
		OpName: "GetBucketPolicy",
		Method: "GET",
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketPolicyRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketPolicy",
		Method: "GET",
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketPolicy_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	putBody := `{"Version":"1","Statement":[{"Action":["ossvector:PutVectors","ossvector:GetVectors"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:ossvector:cn-hangzhou:1234567890:*"]}]}`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	body, err := ioutil.ReadAll(output.Body)
	assert.Nil(t, err)
	result := &GetBucketPolicyResult{
		Body: string(body),
	}
	err = c.client.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, result.Body, putBody)

	putBody = `{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "66C2FF09FDF07830343C72EC",
    "HostId": "bucket.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "bucket",
    "EC": "0015-00000101"
  }
}`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	body, err = ioutil.ReadAll(output.Body)
	assert.Nil(t, err)
	result = &GetBucketPolicyResult{
		Body: string(body),
	}
	err = c.client.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	putBody = `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	body, err = ioutil.ReadAll(output.Body)
	assert.Nil(t, err)
	result = &GetBucketPolicyResult{
		Body: string(body),
	}
	err = c.client.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteBucketPolicy_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *DeleteBucketPolicyRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketPolicyRequest{}
	input = &OperationInput{
		OpName: "DeleteBucketPolicy",
		Method: "DELETE",
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteBucketPolicyRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketPolicy",
		Method: "DELETE",
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteBucketPolicy_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteBucketPolicyResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	body := `{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "66C2FF09FDF07830343C72EC",
    "HostId": "bucket.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "bucket",
    "EC": "0015-00000101"
  }
}`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &DeleteBucketPolicyResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	body = `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &DeleteBucketPolicyResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
