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

func TestMarshalInput_PutTableBucketEncryption(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *PutTableBucketEncryptionRequest
	var input *oss.OperationInput
	var err error

	request = &PutTableBucketEncryptionRequest{}
	input = &oss.OperationInput{
		OpName: "PutTableBucketEncryption",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN.")

	request = &PutTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "PutTableBucketEncryption",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, EncryptionConfiguration.")

	request = &PutTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		EncryptionConfiguration: &EncryptionConfiguration{
			KmsKeyArn:    oss.Ptr(""),
			SseAlgorithm: oss.Ptr("AES256"),
		},
	}
	input = &oss.OperationInput{
		OpName: "PutTableBucketEncryption",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"\",\"sseAlgorithm\":\"AES256\"}}")
}

func TestUnmarshalOutput_PutTableBucketEncryption(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutTableBucketEncryptionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutTableBucketEncryptionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
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
	result = &PutTableBucketEncryptionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
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
	result = &PutTableBucketEncryptionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetTableBucketEncryption(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetTableBucketEncryptionRequest
	var input *oss.OperationInput
	var err error

	request = &GetTableBucketEncryptionRequest{}
	input = &oss.OperationInput{
		OpName: "GetTableBucketEncryption",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	request = &GetTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableBucketEncryption",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
}

func TestUnmarshalOutput_GetTableBucketEncryption(t *testing.T) {
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
	result := &GetTableBucketEncryptionResult{}
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
	result = &GetTableBucketEncryptionResult{}
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
	result = &GetTableBucketEncryptionResult{}
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
	result = &GetTableBucketEncryptionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteTableBucketEncryption(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *DeleteTableBucketEncryptionRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteTableBucketEncryptionRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteTableBucketEncryption",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN.")

	request = &DeleteTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTableBucketEncryption",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/encryption", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
}

func TestUnmarshalOutput_DeleteTableBucketEncryption(t *testing.T) {
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
	result := &DeleteTableBucketEncryptionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &DeleteTableBucketEncryptionResult{}
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
	result = &DeleteTableBucketEncryptionResult{}
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
	result = &DeleteTableBucketEncryptionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
