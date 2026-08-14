package agentic

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutAgenticBucketPublicAccessBlock(t *testing.T) {
	c := AgenticBucketClient{}
	var request *PutAgenticBucketPublicAccessBlockRequest
	var input *oss.OperationInput
	var err error

	request = &PutAgenticBucketPublicAccessBlockRequest{}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketPublicAccessBlock",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr("my-agentic"),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			BlockPublicAccess: oss.Ptr(true),
		},
	}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketPublicAccessBlock",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.NotNil(t, input.Body)
}

func TestUnmarshalOutput_PutAgenticBucketPublicAccessBlock(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutAgenticBucketPublicAccessBlockResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
}

func TestMarshalInput_GetAgenticBucketPublicAccessBlock(t *testing.T) {
	c := AgenticBucketClient{}

	request := &GetAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketPublicAccessBlock",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetAgenticBucketPublicAccessBlock(t *testing.T) {
	c := AgenticBucketClient{}
	body := `<?xml version="1.0" encoding="UTF-8"?>
<PublicAccessBlockConfiguration>
  <BlockPublicAccess>true</BlockPublicAccess>
</PublicAccessBlockConfiguration>`

	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetAgenticBucketPublicAccessBlockResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotNil(t, result.PublicAccessBlockConfiguration)
	assert.True(t, *result.PublicAccessBlockConfiguration.BlockPublicAccess)
}

func TestMarshalInput_DeleteAgenticBucketPublicAccessBlock(t *testing.T) {
	c := AgenticBucketClient{}

	request := &DeleteAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input := &oss.OperationInput{
		OpName: "DeleteAgenticBucketPublicAccessBlock",
		Method: "DELETE",
		Bucket: request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteAgenticBucketPublicAccessBlock(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteAgenticBucketPublicAccessBlockResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
}
