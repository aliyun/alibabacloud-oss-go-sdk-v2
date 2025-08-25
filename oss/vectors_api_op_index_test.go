package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutVectorIndexRequest
	var input *OperationInput
	var err error

	request = &PutVectorIndexRequest{}
	input = &OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"PutVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutVectorIndexRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"PutVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &PutVectorIndexRequest{
		Bucket:         Ptr("oss-demo"),
		DataType:       Ptr("string"),
		Dimension:      Ptr(128),
		DistanceMetric: Ptr("cosine"),
		IndexName:      Ptr("exampleIndex"),
		Metadata: map[string]interface{}{
			"nonFilterableMetadataKeys": []string{"foo", "bar"},
		},
	}
	input = &OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"PutVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["PutVectorIndex"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\",\"metadata\":{\"nonFilterableMetadataKeys\":[\"foo\",\"bar\"]}}")
}

func TestUnmarshalOutput_PutVectorIndex(t *testing.T) {
	c := VectorsClient{}
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
	result := &PutVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	body := `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *GetVectorIndexRequest
	var input *OperationInput
	var err error

	request = &GetVectorIndexRequest{}
	input = &OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"GetVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &GetVectorIndexRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"GetVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &GetVectorIndexRequest{
		Bucket:    Ptr("oss-demo"),
		IndexName: Ptr("demo"),
	}
	input = &OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"GetVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["GetVectorIndex"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"demo\"}")
}

func TestUnmarshalOutput_GetVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `{
   "index": { 
      "createTime": "2025-08-02T10:49:17.289372919+08:00",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "string",
      "metadata": { 
         "nonFilterableMetadataKeys": ["foo", "bar"]
      },
      "status": "running"
   },
   "vectorBucketName": "bucket"
}`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &GetVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Index.CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.Local))
	assert.Equal(t, *result.Index.DataType, "string")
	assert.Equal(t, *result.Index.Dimension, 128)
	assert.Equal(t, *result.Index.DistanceMetric, "string")
	assert.Equal(t, *result.Index.IndexName, "string")
	assert.Len(t, result.Index.Metadata["nonFilterableMetadataKeys"], 2)
	if metadataValue, ok := result.Index.Metadata["nonFilterableMetadataKeys"]; ok {
		if keys, ok := metadataValue.([]interface{}); ok {
			assert.Equal(t, keys[0].(string), "foo")
			assert.Equal(t, keys[1].(string), "bar")
		}
	}
	assert.Equal(t, *result.Index.Status, "running")
	assert.Equal(t, *result.VectorBucketName, "bucket")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_ListVectorIndexes(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *ListVectorIndexesRequest
	var input *OperationInput
	var err error

	request = &ListVectorIndexesRequest{}
	input = &OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &ListVectorIndexesRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectorIndexes"], "")
	assert.Equal(t, input.Headers[HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")

	request = &ListVectorIndexesRequest{
		Bucket: Ptr("oss-demo"),
		Prefix: Ptr("prefix"),
	}
	input = &OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectorIndexes"], "")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[HTTPHeaderContentType], contentTypeJSON)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"prefix\":\"prefix\"}")

	request = &ListVectorIndexesRequest{
		Bucket:     Ptr("oss-demo"),
		MaxResults: Ptr(100),
		NextToken:  Ptr("123"),
		Prefix:     Ptr("prefix"),
	}
	input = &OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectorIndexes"], "")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[HTTPHeaderContentType], contentTypeJSON)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"maxResults\":100,\"nextToken\":\"123\",\"prefix\":\"prefix\"}")
}

func TestUnmarshalOutput_ListVectorIndexes(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `{
  "indexes": [
    { 
      "createTime": "2025-08-02T10:49:17.289372919+08:00",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "demo1",
      "metadata": { 
        "nonFilterableMetadataKeys": ["foo", "bar"]
      },
      "status": "running",
      "vectorBucketName": "bucket"
    },
    { 
      "createTime": "2025-08-20T10:49:17.289372919+08:00",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "demo2",
      "metadata": { 
        "nonFilterableMetadataKeys": ["foo2", "bar2"]
      },
      "status": "deleting",
      "vectorBucketName": "bucket2"
    }
  ],
  "nextToken": "123"
}`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &ListVectorIndexesResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Indexes), 2)
	assert.Equal(t, *result.Indexes[0].CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.Local))
	assert.Equal(t, *result.Indexes[0].DataType, "string")
	assert.Equal(t, *result.Indexes[0].Dimension, 128)
	assert.Equal(t, *result.Indexes[0].DistanceMetric, "string")
	assert.Equal(t, *result.Indexes[0].IndexName, "demo1")
	assert.Len(t, result.Indexes[0].Metadata["nonFilterableMetadataKeys"], 2)
	if metadataValue, ok := result.Indexes[0].Metadata["nonFilterableMetadataKeys"]; ok {
		if keys, ok := metadataValue.([]interface{}); ok {
			assert.Equal(t, keys[0].(string), "foo")
			assert.Equal(t, keys[1].(string), "bar")
		}
	}
	assert.Equal(t, *result.Indexes[0].Status, "running")
	assert.Equal(t, *result.Indexes[0].VectorBucketName, "bucket")

	assert.Equal(t, *result.Indexes[1].CreateTime, time.Date(2025, time.August, 20, 10, 49, 17, 289372919, time.Local))
	assert.Equal(t, *result.Indexes[1].DataType, "string")
	assert.Equal(t, *result.Indexes[1].Dimension, 128)
	assert.Equal(t, *result.Indexes[1].DistanceMetric, "string")
	assert.Equal(t, *result.Indexes[1].IndexName, "demo2")
	if metadataValue, ok := result.Indexes[1].Metadata["nonFilterableMetadataKeys"]; ok {
		if keys, ok := metadataValue.([]interface{}); ok {
			assert.Equal(t, keys[0].(string), "foo2")
			assert.Equal(t, keys[1].(string), "bar2")
		}
	}
	assert.Equal(t, *result.Indexes[1].VectorBucketName, "bucket2")
	assert.Equal(t, *result.Indexes[1].Status, "deleting")

	assert.Equal(t, *result.NextToken, "123")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorIndexesResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorIndexesResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorIndexesResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *DeleteVectorIndexRequest
	var input *OperationInput
	var err error

	request = &DeleteVectorIndexRequest{}
	input = &OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"DeleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &DeleteVectorIndexRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"DeleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &DeleteVectorIndexRequest{
		Bucket:    Ptr("oss-demo"),
		IndexName: Ptr("demo"),
	}
	input = &OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"DeleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["DeleteVectorIndex"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"demo\"}")
}

func TestUnmarshalOutput_DeleteVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Body:       io.NopCloser(bytes.NewReader([]byte(nil))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &DeleteVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &DeleteVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	body := `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &DeleteVectorIndexResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
