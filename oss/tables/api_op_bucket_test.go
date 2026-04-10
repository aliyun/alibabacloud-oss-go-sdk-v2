package tables

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
		Bucket: request.Name,
		Key:    oss.Ptr("buckets"),
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	request = &CreateTableBucketRequest{
		Name: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTableBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Name,
		Key:    oss.Ptr("buckets"),
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, *input.Key, "buckets")
	jsonStr, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"name\":\"oss-demo\"}")

	request = &CreateTableBucketRequest{
		Name: oss.Ptr("oss-demo"),
		EncryptionConfiguration: &EncryptionConfiguration{
			KmsKeyArn:    oss.Ptr("arn"),
			SseAlgorithm: oss.Ptr("AES256"),
		},
	}
	input = &oss.OperationInput{
		OpName: "CreateTableBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Name,
		Key:    oss.Ptr("buckets"),
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, *input.Key, "buckets")
	jsonStr, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"arn\",\"sseAlgorithm\":\"AES256\"},\"name\":\"oss-demo\"}")
}

func TestUnmarshalOutput_CreateTableBucket(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := `{"arn": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"}`
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
	assert.Equal(t, oss.ToString(result.Arn), "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")

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

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "BadErrorResponse",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
	result = &CreateTableBucketResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "BadErrorResponse")
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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN.")

	request = &GetTableBucketRequest{
		TableBucketARN: oss.Ptr("invlid-arn"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableBucket",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed ARN - doesn't start with 'acs:'")

	request = &GetTableBucketRequest{
		TableBucketARN: oss.Ptr("acs:osstables"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableBucket",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed ARN - no service specified")

	request = &GetTableBucketRequest{
		TableBucketARN: oss.Ptr("acs:osstables:oss-cn-beijing"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableBucket",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed ARN - no region specified")

	request = &GetTableBucketRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableBucket",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Bucket, request.TableBucketARN)
	assert.Equal(t, *input.Key, "buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket")
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, *input.Key, "buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket")
}

func TestUnmarshalOutput_GetTableBucket(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "arn": "acs:osstables:cn-beijing:12345657890:bucket/demo-bucket",
   "createdAt": "2026-04-01T09:42:50.000000+00:00",
   "name": "demo-bucket",
   "ownerAccountId": "12345657890",
   "tableBucketId": "50859410-3482-401c-b500-605c22848ef4",
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
	assert.Equal(t, *result.Arn, "acs:osstables:cn-beijing:12345657890:bucket/demo-bucket")
	assert.Equal(t, *result.Name, "demo-bucket")
	assert.Equal(t, *result.CreatedAt, "2026-04-01T09:42:50.000000+00:00")
	assert.Equal(t, *result.OwnerAccountId, "12345657890")
	assert.Equal(t, *result.TableBucketId, "50859410-3482-401c-b500-605c22848ef4")
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
		Key: oss.Ptr("buckets"),
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
		Key: oss.Ptr("buckets"),
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Key, "buckets")
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
		Key: oss.Ptr("buckets"),
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Key, "buckets")
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
  "continuationToken": "Cj5hY3M6b3NzdGFibGVzOmNuLWJlaWppbmc6MTc2MDIyNTU0NTA4NDMzMTpidWNrZXQvZGVtby13YWxrZXItMQ--",
  "tableBuckets": [{
    "arn": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket",
    "createdAt": "2026-04-02T05:27:31.000000+00:00",
    "name": "demo-bucket",
    "ownerAccountId": "1234567890",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5",
    "type": "customer"
  },
  {
    "arn": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket-1",
    "createdAt": "2026-04-02T05:27:32.000000+00:00",
    "name": "demo-bucket-1",
    "ownerAccountId": "1234567890",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f3",
    "type": "customer"
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
	assert.Equal(t, *result.ContinuationToken, "Cj5hY3M6b3NzdGFibGVzOmNuLWJlaWppbmc6MTc2MDIyNTU0NTA4NDMzMTpidWNrZXQvZGVtby13YWxrZXItMQ--")
	assert.Equal(t, len(result.TableBuckets), 2)
	assert.Equal(t, *result.TableBuckets[0].CreatedAt, "2026-04-02T05:27:31.000000+00:00")
	assert.Equal(t, *result.TableBuckets[0].Arn, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *result.TableBuckets[0].Name, "demo-bucket")
	assert.Equal(t, *result.TableBuckets[0].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")
	assert.Equal(t, *result.TableBuckets[0].OwnerAccountId, "1234567890")
	assert.Equal(t, *result.TableBuckets[0].Type, "customer")

	assert.Equal(t, *result.TableBuckets[1].CreatedAt, "2026-04-02T05:27:32.000000+00:00")
	assert.Equal(t, *result.TableBuckets[1].Arn, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket-1")
	assert.Equal(t, *result.TableBuckets[1].Name, "demo-bucket-1")
	assert.Equal(t, *result.TableBuckets[1].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f3")
	assert.Equal(t, *result.TableBuckets[1].OwnerAccountId, "1234567890")
	assert.Equal(t, *result.TableBuckets[1].Type, "customer")

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
    "message": "AccessDenied"
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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	request = &DeleteTableBucketRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTableBucket",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Bucket, request.TableBucketARN)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, *input.Key, "buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket")
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
