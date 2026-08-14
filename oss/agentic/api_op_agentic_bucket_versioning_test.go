package agentic

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutAgenticBucketVersioning(t *testing.T) {
	c := AgenticBucketClient{}
	var request *PutAgenticBucketVersioningRequest
	var input *oss.OperationInput
	var err error

	request = &PutAgenticBucketVersioningRequest{}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketVersioning",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutAgenticBucketVersioningRequest{
		Bucket: oss.Ptr("my-agentic"),
		VersioningConfiguration: &VersioningConfiguration{
			Status: oss.VersionEnabled,
		},
	}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketVersioning",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.NotNil(t, input.Body)
}

func TestUnmarshalOutput_PutAgenticBucketVersioning(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutAgenticBucketVersioningResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
}

func TestMarshalInput_GetAgenticBucketVersioning(t *testing.T) {
	c := AgenticBucketClient{}

	request := &GetAgenticBucketVersioningRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketVersioning",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetAgenticBucketVersioning(t *testing.T) {
	c := AgenticBucketClient{}
	body := `<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
  <Status>Enabled</Status>
</VersioningConfiguration>`

	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetAgenticBucketVersioningResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotNil(t, result.VersioningConfiguration)
	assert.Equal(t, oss.VersionEnabled, result.VersioningConfiguration.Status)
}
