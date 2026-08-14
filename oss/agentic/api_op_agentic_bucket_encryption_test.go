package agentic

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutAgenticBucketEncryption(t *testing.T) {
	c := AgenticBucketClient{}
	var request *PutAgenticBucketEncryptionRequest
	var input *oss.OperationInput
	var err error

	request = &PutAgenticBucketEncryptionRequest{}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketEncryption",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr("my-agentic"),
		ServerSideEncryptionRule: &ServerSideEncryptionRule{
			ApplyServerSideEncryptionByDefault: &ApplyServerSideEncryptionByDefault{
				SSEAlgorithm: oss.Ptr("AES256"),
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketEncryption",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.NotNil(t, input.Body)
}

func TestUnmarshalOutput_PutAgenticBucketEncryption(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutAgenticBucketEncryptionResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
}

func TestMarshalInput_GetAgenticBucketEncryption(t *testing.T) {
	c := AgenticBucketClient{}

	request := &GetAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketEncryption",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetAgenticBucketEncryption(t *testing.T) {
	c := AgenticBucketClient{}
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ServerSideEncryptionRule>
  <ApplyServerSideEncryptionByDefault>
    <SSEAlgorithm>KMS</SSEAlgorithm>
    <KMSMasterKeyID>9468da86-3509-4f8d-a61e-6eab****</KMSMasterKeyID>
    <KMSDataEncryption>SM4</KMSDataEncryption>
  </ApplyServerSideEncryptionByDefault>
</ServerSideEncryptionRule>`

	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetAgenticBucketEncryptionResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotNil(t, result.ServerSideEncryptionRule)
	assert.Equal(t, "KMS", *result.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
	assert.Equal(t, "9468da86-3509-4f8d-a61e-6eab****", *result.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID)
	assert.Equal(t, "SM4", *result.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.KMSDataEncryption)
}

func TestMarshalInput_DeleteAgenticBucketEncryption(t *testing.T) {
	c := AgenticBucketClient{}

	request := &DeleteAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input := &oss.OperationInput{
		OpName: "DeleteAgenticBucketEncryption",
		Method: "DELETE",
		Bucket: request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteAgenticBucketEncryption(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteAgenticBucketEncryptionResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
}
