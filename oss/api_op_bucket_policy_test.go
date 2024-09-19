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

func TestMarshalInput_PutBucketPolicy(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketPolicyRequest
	var input *OperationInput
	var err error

	request = &PutBucketPolicyRequest{}
	input = &OperationInput{
		OpName: "PutBucketPolicy",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutBucketPolicyRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucketPolicy",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, Body.")
	putBody := `{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:oss:*:1234567890:*/*"]}]}`
	request = &PutBucketPolicyRequest{
		Bucket: Ptr("oss-demo"),
		Body:   strings.NewReader(putBody),
	}
	input = &OperationInput{
		OpName: "PutBucketPolicy",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:PutObject\",\"oss:GetObject\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:oss:*:1234567890:*/*\"]}]}")
}

func TestUnmarshalOutput_PutBucketPolicy(t *testing.T) {
	c := Client{}
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
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketPolicyResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketPolicyResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketPolicyResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketPolicy(t *testing.T) {
	c := Client{}
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
	err = c.marshalInput(request, input, updateContentMd5)
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketPolicy(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	putBody := `{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:oss:*:1234567890:*/*"]}]}`
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
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, result.Body, putBody)

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>66C2FF09FDF07830343C72EC</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
		<EC>0015-00000101</EC>
	</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	body, err = ioutil.ReadAll(output.Body)
	assert.Nil(t, err)
	result = &GetBucketPolicyResult{
		Body: string(body),
	}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	body, err = ioutil.ReadAll(output.Body)
	assert.Nil(t, err)
	result = &GetBucketPolicyResult{
		Body: string(body),
	}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteBucketPolicy(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteBucketPolicyRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketPolicyRequest{}
	input = &OperationInput{
		OpName: "DeleteBucketPolicy",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteBucketPolicyRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketPolicy",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteBucketPolicy(t *testing.T) {
	c := Client{}
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
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	body := `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>66C2FF09FDF07830343C72EC</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
		<EC>0015-00000101</EC>
	</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketPolicyResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketPolicyResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketPolicyStatus(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketPolicyStatusRequest
	var input *OperationInput
	var err error

	request = &GetBucketPolicyStatusRequest{}
	input = &OperationInput{
		OpName: "GetBucketPolicyStatus",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"policyStatus": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policyStatus"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketPolicyStatusRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketPolicyStatus",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"policyStatus": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policyStatus"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketPolicyStatus(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<PolicyStatus>
   <IsPublic>true</IsPublic>
</PolicyStatus>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(strings.NewReader(body)),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketPolicyStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.True(t, *result.PolicyStatus.IsPublic)
	body = `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>66C2FF09FDF07830343C72EC</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
		<EC>0015-00000101</EC>
	</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketPolicyStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketPolicyStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
