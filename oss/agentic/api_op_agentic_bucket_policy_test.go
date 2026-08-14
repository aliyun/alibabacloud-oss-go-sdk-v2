package agentic

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutAgenticBucketPolicy(t *testing.T) {
	c := AgenticBucketClient{}
	var request *PutAgenticBucketPolicyRequest
	var input *oss.OperationInput
	var err error

	request = &PutAgenticBucketPolicyRequest{}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketPolicy",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	policyBody := `{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["*"],"Resource":["acs:oss:*:*:my-agentic/*"]}]}`
	request = &PutAgenticBucketPolicyRequest{
		Bucket: oss.Ptr("my-agentic"),
		Body:   strings.NewReader(policyBody),
	}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketPolicy",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.NotNil(t, input.Body)
}

func TestUnmarshalOutput_PutAgenticBucketPolicy(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutAgenticBucketPolicyResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
}

func TestMarshalInput_GetAgenticBucketPolicy(t *testing.T) {
	c := AgenticBucketClient{}

	request := &GetAgenticBucketPolicyRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketPolicy",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetAgenticBucketPolicy(t *testing.T) {
	c := AgenticBucketClient{}
	policyBody := `{"Version":"1","Statement":[]}`

	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(policyBody))),
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetAgenticBucketPolicyResult{}

	body, _ := io.ReadAll(output.Body)
	result.Body = string(body)
	err := c.clientImpl.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, policyBody, result.Body)
}

func TestMarshalInput_DeleteAgenticBucketPolicy(t *testing.T) {
	c := AgenticBucketClient{}

	request := &DeleteAgenticBucketPolicyRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input := &oss.OperationInput{
		OpName: "DeleteAgenticBucketPolicy",
		Method: "DELETE",
		Bucket: request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteAgenticBucketPolicy(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteAgenticBucketPolicyResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
}
