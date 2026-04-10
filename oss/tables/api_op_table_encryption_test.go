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

func TestMarshalInput_GetTableEncryption(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetTableEncryptionRequest
	var input *oss.OperationInput
	var err error

	request = &GetTableEncryptionRequest{}
	input = &oss.OperationInput{
		OpName: "GetTableEncryption",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/encryption", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn")

	request = &GetTableEncryptionRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableEncryption",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/encryption", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	request = &GetTableEncryptionRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableEncryption",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/encryption", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name.")

	request = &GetTableEncryptionRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("space"),
		Name:      oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableEncryption",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/encryption", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/encryption")
}

func TestUnmarshalOutput_GetTableEncryption(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "encryptionConfiguration": { 
      "kmsKeyArn": "",
      "sseAlgorithm": "AES256"
   }
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
	result := &GetTableEncryptionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.EncryptionConfiguration.KmsKeyArn, "")
	assert.Equal(t, *result.EncryptionConfiguration.SseAlgorithm, "AES256")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetTableEncryptionResult{}
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
	result = &GetTableEncryptionResult{}
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
	result = &GetTableEncryptionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
