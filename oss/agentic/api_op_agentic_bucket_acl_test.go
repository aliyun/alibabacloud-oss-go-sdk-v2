package agentic

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutAgenticBucketAcl(t *testing.T) {
	c := AgenticBucketClient{}
	var request *PutAgenticBucketAclRequest
	var input *oss.OperationInput
	var err error

	request = &PutAgenticBucketAclRequest{}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketAcl",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutAgenticBucketAclRequest{
		Bucket: oss.Ptr("my-agentic"),
		Acl:    oss.BucketACLPrivate,
	}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketAcl",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, string(oss.BucketACLPrivate), input.Headers["x-oss-acl"])
}

func TestUnmarshalOutput_PutAgenticBucketAcl(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutAgenticBucketAclResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
}

func TestMarshalInput_GetAgenticBucketAcl(t *testing.T) {
	c := AgenticBucketClient{}
	var request *GetAgenticBucketAclRequest
	var input *oss.OperationInput
	var err error

	request = &GetAgenticBucketAclRequest{}
	input = &oss.OperationInput{
		OpName: "GetAgenticBucketAcl",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)

	request = &GetAgenticBucketAclRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input = &oss.OperationInput{
		OpName: "GetAgenticBucketAcl",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetAgenticBucketAcl(t *testing.T) {
	c := AgenticBucketClient{}
	body := `<?xml version="1.0" encoding="UTF-8"?>
<AccessControlPolicy>
  <Owner>
    <ID>1234567890123456</ID>
    <DisplayName>owner-name</DisplayName>
  </Owner>
  <AccessControlList>
    <Grant>private</Grant>
  </AccessControlList>
</AccessControlPolicy>`

	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetAgenticBucketAclResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "private", *result.ACL)
	assert.NotNil(t, result.Owner)
	assert.Equal(t, "1234567890123456", *result.Owner.ID)
}
