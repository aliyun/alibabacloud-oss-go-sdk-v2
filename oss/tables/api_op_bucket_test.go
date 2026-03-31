package tables

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateTableBucket(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *CreateTableBucketRequest
	var input *oss.OperationInput
	var err error

	request = &CreateTableBucketRequest{}
	input = &oss.OperationInput{
		OpName: "CreateTableBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &CreateTableBucketRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTableBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	jsonStr, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"name\":\"oss-demo\"}")

	request = &CreateTableBucketRequest{
		Bucket: oss.Ptr("oss-demo"),
		EncryptionConfiguration: &EncryptionConfiguration{
			KmsKeyArn:    oss.Ptr("arn"),
			SseAlgorithm: oss.Ptr("AES256"),
		},
		StorageClassConfiguration: &StorageClassConfiguration{
			StorageClass: oss.StorageClassStandard,
		},
		Tags: map[string]any{
			"k1": "v1", "k2": "v2",
		},
	}
	input = &oss.OperationInput{
		OpName: "CreateTableBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	jsonStr, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"arn\",\"sseAlgorithm\":\"AES256\"},\"name\":\"oss-demo\",\"storageClassConfiguration\":{\"storageClass\":\"Standard\"},\"tags\":{\"k1\":\"v1\",\"k2\":\"v2\"}}")
}

func TestUnmarshalOutput_CreateTableBucket(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := `{"arn": "test-arn"}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &CreateTableBucketResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, oss.ToString(result.Arn), "test-arn")

	output = &oss.OperationOutput{
		StatusCode: 409,
		Status:     "BucketAlreadyExist",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &CreateTableBucketResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "BucketAlreadyExist")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetTableBucket(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetTableBucketRequest
	var input *oss.OperationInput
	var err error

	request = &GetTableBucketRequest{}
	input = &oss.OperationInput{
		OpName: "GetTableBucket",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetTableBucketRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableBucket",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
}

func TestUnmarshalOutput_GetTableBucket(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "arn": "test-arn",
   "createdAt": "2013-07-31T10:56:21.000Z",
   "name": "oss-bucket",
   "ownerAccountId": "123456",
   "tableBucketId": "123",
   "type": "oss"
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
	result := &GetTableBucketResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.Arn, "test-arn")
	assert.Equal(t, *result.Name, "oss-bucket")
	assert.Equal(t, *result.CreatedAt, "2013-07-31T10:56:21.000Z")
	assert.Equal(t, *result.OwnerAccountId, "123456")
	assert.Equal(t, *result.TableBucketId, "123")
	assert.Equal(t, *result.Type, "oss")
}

func TestMarshalInput_ListTableBuckets(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *ListTableBucketsRequest
	var input *oss.OperationInput
	var err error

	request = &ListTableBucketsRequest{}
	input = &oss.OperationInput{
		OpName: "ListTableBuckets",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
	}
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)

	request = &ListTableBucketsRequest{
		MaxBuckets: 10,
		Prefix:     oss.Ptr("/"),
	}
	input = &oss.OperationInput{
		OpName: "ListTableBuckets",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["maxBuckets"], "10")
	assert.Equal(t, input.Parameters["prefix"], "/")

	request = &ListTableBucketsRequest{
		MaxBuckets:        10,
		Prefix:            oss.Ptr("/"),
		ContinuationToken: oss.Ptr("123"),
	}
	input = &oss.OperationInput{
		OpName: "ListTableBuckets",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["maxBuckets"], "10")
	assert.Equal(t, input.Parameters["prefix"], "/")
	assert.Equal(t, input.Parameters["continuationToken"], "123")
}

func TestUnmarshalOutput_ListTableBuckets(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := `{
  "continuationToken": "token-123",
  "tableBuckets": [{
    "tableBucketArn": "test-arn",
    "createdAt": "2014-02-17T18:12:43.000Z",
    "name": "app-base-oss",
    "ownerAccountId": "123456",
    "tableBucketId": "123"
  },
  {
    "tableBucketArn": "test-arn",
    "createdAt": "2014-02-18T18:12:43.000Z",
    "name": "app-base-oss2",
    "ownerAccountId": "123456",
    "tableBucketId": "124"
  }]
}`

	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"5374A2880232A65C2300****"},
		},
	}
	result := &ListTableBucketsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5374A2880232A65C2300****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.ContinuationToken, "token-123")
	assert.Equal(t, len(result.Buckets), 2)
	assert.Equal(t, *result.Buckets[0].CreatedAt, "2014-02-17T18:12:43.000Z")
	assert.Equal(t, *result.Buckets[0].TableBucketArn, "test-arn")
	assert.Equal(t, *result.Buckets[0].Name, "app-base-oss")
	assert.Equal(t, *result.Buckets[0].TableBucketId, "123")
	assert.Equal(t, *result.Buckets[0].OwnerAccountId, "123456")

	assert.Equal(t, *result.Buckets[1].CreatedAt, "2014-02-18T18:12:43.000Z")
	assert.Equal(t, *result.Buckets[1].TableBucketArn, "test-arn")
	assert.Equal(t, *result.Buckets[1].Name, "app-base-oss2")
	assert.Equal(t, *result.Buckets[1].TableBucketId, "124")
	assert.Equal(t, *result.Buckets[1].OwnerAccountId, "123456")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
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
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	resultErr := &ListTableBucketsResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteTableBucket(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *DeleteTableBucketRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteTableBucketRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteTableBucket",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &DeleteTableBucketRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTableBucket",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
}

func TestUnmarshalOutput_DeleteTableBucket(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	output = &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5C3D9778CC1C2AEDF85B****"},
		},
	}
	result := &DeleteTableBucketResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C3D9778CC1C2AEDF85B****")

	output = &oss.OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
}
