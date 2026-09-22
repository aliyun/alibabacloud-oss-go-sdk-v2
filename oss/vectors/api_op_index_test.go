package vectors

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutVectorIndexRequest
	var input *oss.OperationInput
	var err error

	request = &PutVectorIndexRequest{}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutVectorIndexRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &PutVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DataType")

	request = &PutVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
		DataType:  oss.Ptr("string"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Dimension")

	request = &PutVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		DataType:  oss.Ptr("string"),
		IndexName: oss.Ptr("exampleIndex"),
		Dimension: oss.Ptr(128),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DistanceMetric")

	request = &PutVectorIndexRequest{
		Bucket:         oss.Ptr("oss-demo"),
		DataType:       oss.Ptr("string"),
		Dimension:      oss.Ptr(128),
		DistanceMetric: oss.Ptr("cosine"),
		IndexName:      oss.Ptr("exampleIndex"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["putVectorIndex"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\"}")

	request = &PutVectorIndexRequest{
		Bucket:         oss.Ptr("oss-demo"),
		DataType:       oss.Ptr("string"),
		Dimension:      oss.Ptr(128),
		DistanceMetric: oss.Ptr("cosine"),
		IndexName:      oss.Ptr("exampleIndex"),
		Metadata: map[string]any{
			"nonFilterableMetadataKeys": []string{"foo", "bar"},
		},
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["putVectorIndex"], "")
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\",\"metadata\":{\"nonFilterableMetadataKeys\":[\"foo\",\"bar\"]}}")
}

func TestUnmarshalOutput_PutVectorIndex(t *testing.T) {
	c := VectorsClient{}
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
	result := &PutVectorIndexResult{}
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
	result = &PutVectorIndexResult{}
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
	result = &PutVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
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
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
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
	var input *oss.OperationInput
	var err error

	request = &GetVectorIndexRequest{}
	input = &oss.OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"getVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &GetVectorIndexRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"getVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &GetVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("demo"),
	}
	input = &oss.OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"getVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
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
	var output *oss.OperationOutput
	var err error
	body := `{
   "index": { 
      "createTime": "2025-08-02T10:49:17.289372919Z",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "string",
      "metadata": { 
         "nonFilterableMetadataKeys": ["foo", "bar"]
      },
      "status": "running",
      "bucketArn": "acs:oss:::test-bucket"
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
	result := &GetVectorIndexResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	//dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Index.CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.UTC))
	assert.Equal(t, *result.Index.DataType, "string")
	assert.Equal(t, *result.Index.Dimension, 128)
	assert.Equal(t, *result.Index.DistanceMetric, "string")
	assert.Equal(t, *result.Index.IndexName, "string")
	assert.Len(t, result.Index.Metadata["nonFilterableMetadataKeys"], 2)
	if metadataValue, ok := result.Index.Metadata["nonFilterableMetadataKeys"]; ok {
		if keys, ok := metadataValue.([]any); ok {
			assert.Equal(t, keys[0].(string), "foo")
			assert.Equal(t, keys[1].(string), "bar")
		}
	}
	assert.Equal(t, *result.Index.Status, "running")
	assert.Equal(t, *result.Index.BucketArn, "acs:oss:::test-bucket")

	body = `{
  "index": {
    "indexName": "exampleIndex",
    "mode": "fusion",
    "dataType": "string",
    "dimension": 128,
    "status": "running",
    "bucketArn": "acs:oss:::test-bucket",
    "schemaConfiguration": {
      "fields": [
        {
          "dataType": "float32",
          "dimension": 1024,
          "distanceMetric": "euclidean",
          "name": "vector_1",
          "type": "vector"
        },
        {
          "dataType": "float32",
          "dimension": 512,
          "distanceMetric": "cosine",
          "name": "vector_2",
          "type": "vector"
        },
        {
          "isArray": true,
          "name": "timestamps",
          "type": "long"
        },
        {
          "name": "price",
          "type": "double"
        },
        {
          "name": "ip",
          "type": "ip"
        },
        {
          "name": "location",
          "type": "geoPoint"
        },
        {
          "name": "tag",
          "type": "string"
        },
        {
          "isPartitionKey": true,
          "name": "user_id",
          "type": "string"
        },
        {
          "isArray": true,
          "name": "tags",
          "type": "string"
        },
        {
          "exactMatch": true,
          "name": "title_1",
          "text": {
            "analyzer": "standard",
            "analyzerParameters": {
              "caseSensitive": true,
              "delimitWord": false
            },
            "enabled": true
          },
          "type": "string"
        },
        {
          "exactMatch": false,
          "name": "title_2",
          "text": {
            "analyzer": "split",
            "analyzerParameters": {
              "caseSensitive": true,
              "delimiter": " "
            },
            "enabled": true
          },
          "type": "string"
        }
      ]
    }
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
	result = &GetVectorIndexResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	//dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Index.DataType, string(FieldTypeString))
	assert.Equal(t, *result.Index.IndexName, "exampleIndex")
	assert.Equal(t, *result.Index.Mode, "fusion")
	assert.Equal(t, *result.Index.Dimension, 128)
	assert.Equal(t, *result.Index.Status, "running")
	assert.Equal(t, *result.Index.BucketArn, "acs:oss:::test-bucket")
	assert.Len(t, result.Index.SchemaConfiguration.Fields, 11)
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[0].DataType, VectorDataTypeFloat32)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[0].Dimension, int(1024))
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[0].DistanceMetric, DistanceMetricTypeEuclidean)

	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[1].Name, "vector_2")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[1].Type, FieldTypeVector)
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[1].DataType, VectorDataTypeFloat32)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[1].Dimension, int(512))
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[1].DistanceMetric, DistanceMetricTypeCosine)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[2].Name, "timestamps")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[2].Type, FieldTypeLong)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[2].IsArray, true)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[2].Name, "timestamps")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[2].Type, FieldTypeLong)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[3].Name, "price")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[3].Type, FieldTypeDouble)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[4].Name, "ip")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[4].Type, FieldTypeIp)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[5].Name, "location")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[5].Type, FieldTypeGeoPoint)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[6].Name, "tag")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[6].Type, FieldTypeString)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[7].IsPartitionKey, true)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[7].Name, "user_id")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[7].Type, FieldTypeString)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[8].IsArray, true)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[8].Name, "tags")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[8].Type, FieldTypeString)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[9].Name, "title_1")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[9].Type, FieldTypeString)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[9].ExactMatch, true)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[9].Text.Enabled, true)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[9].Text.Analyzer, "standard")
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[9].Text.AnalyzerParameters.CaseSensitive, true)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[9].Text.AnalyzerParameters.DelimitWord, false)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[10].Name, "title_2")
	assert.Equal(t, result.Index.SchemaConfiguration.Fields[10].Type, FieldTypeString)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[10].ExactMatch, false)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[10].Text.Enabled, true)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[10].Text.Analyzer, "split")
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[10].Text.AnalyzerParameters.CaseSensitive, true)
	assert.Equal(t, *result.Index.SchemaConfiguration.Fields[10].Text.AnalyzerParameters.Delimiter, " ")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetVectorIndexResult{}
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
	result = &GetVectorIndexResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
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
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetVectorIndexResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
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
	var input *oss.OperationInput
	var err error

	request = &ListVectorIndexesRequest{}
	input = &oss.OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &ListVectorIndexesRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectorIndexes"], "")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")

	request = &ListVectorIndexesRequest{
		Bucket: oss.Ptr("oss-demo"),
		Prefix: oss.Ptr("prefix"),
	}
	input = &oss.OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectorIndexes"], "")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"prefix\":\"prefix\"}")

	request = &ListVectorIndexesRequest{
		Bucket:     oss.Ptr("oss-demo"),
		MaxResults: 100,
		NextToken:  oss.Ptr("123"),
		Prefix:     oss.Ptr("prefix"),
	}
	input = &oss.OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectorIndexes"], "")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"maxResults\":100,\"nextToken\":\"123\",\"prefix\":\"prefix\"}")
}

func TestUnmarshalOutput_ListVectorIndexes(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
  "indexes": [
    { 
      "createTime": "2025-08-02T10:49:17.289372919Z",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "demo1",
      "metadata": { 
        "nonFilterableMetadataKeys": ["foo", "bar"]
      },
      "status": "running",
	  "bucketArn": "acs:oss:::test-bucket"
    },
    { 
      "createTime": "2025-08-20T10:49:17.289372919Z",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "demo2",
      "metadata": { 
        "nonFilterableMetadataKeys": ["foo2", "bar2"]
      },
      "status": "deleting",
	  "bucketArn": "acs:oss:::test-bucket"
    }
  ],
  "nextToken": "123"
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
	result := &ListVectorIndexesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Indexes), 2)
	assert.Equal(t, *result.Indexes[0].CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.UTC))
	assert.Equal(t, *result.Indexes[0].DataType, "string")
	assert.Equal(t, *result.Indexes[0].Dimension, 128)
	assert.Equal(t, *result.Indexes[0].DistanceMetric, "string")
	assert.Equal(t, *result.Indexes[0].IndexName, "demo1")
	assert.Len(t, result.Indexes[0].Metadata["nonFilterableMetadataKeys"], 2)
	if metadataValue, ok := result.Indexes[0].Metadata["nonFilterableMetadataKeys"]; ok {
		if keys, ok := metadataValue.([]any); ok {
			assert.Equal(t, keys[0].(string), "foo")
			assert.Equal(t, keys[1].(string), "bar")
		}
	}
	assert.Equal(t, *result.Indexes[0].Status, "running")
	assert.Equal(t, *result.Indexes[0].BucketArn, "acs:oss:::test-bucket")

	assert.Equal(t, *result.Indexes[1].CreateTime, time.Date(2025, time.August, 20, 10, 49, 17, 289372919, time.UTC))
	assert.Equal(t, *result.Indexes[1].DataType, "string")
	assert.Equal(t, *result.Indexes[1].Dimension, 128)
	assert.Equal(t, *result.Indexes[1].DistanceMetric, "string")
	assert.Equal(t, *result.Indexes[1].IndexName, "demo2")
	if metadataValue, ok := result.Indexes[1].Metadata["nonFilterableMetadataKeys"]; ok {
		if keys, ok := metadataValue.([]any); ok {
			assert.Equal(t, keys[0].(string), "foo2")
			assert.Equal(t, keys[1].(string), "bar2")
		}
	}
	assert.Equal(t, *result.Indexes[1].BucketArn, "acs:oss:::test-bucket")
	assert.Equal(t, *result.Indexes[1].Status, "deleting")
	assert.Equal(t, *result.NextToken, "123")

	body = `{
  "nextToken": "CAESCG15aC1mESgEeKaGA",
  "indexes": [
    {
      "vectorBucketName": "oss-vector-bucket1",
      "indexName": "index1",
      "createTime": "2025-08-20T10:49:17.289372919Z",
      "dataType": "float32",
      "dimension": 512,
      "distanceMetric": "euclidean",
      "status": "enable",
      "mode": "fusion"
    },
    {
      "vectorBucketName": "oss-vector-bucket1",
      "indexName": "index2",
      "createTime": "2025-08-20T10:49:17.289372919Z",
      "status": "enable",
      "mode": "standard"
    }
  ]
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
	result = &ListVectorIndexesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Indexes), 2)
	assert.Equal(t, *result.NextToken, "CAESCG15aC1mESgEeKaGA")
	assert.Equal(t, *result.Indexes[0].VectorBucketName, "oss-vector-bucket1")
	assert.Equal(t, *result.Indexes[0].IndexName, "index1")
	assert.Equal(t, *result.Indexes[0].CreateTime, time.Date(2025, time.August, 20, 10, 49, 17, 289372919, time.UTC))
	assert.Equal(t, *result.Indexes[0].DataType, "float32")
	assert.Equal(t, *result.Indexes[0].Dimension, 512)
	assert.Equal(t, *result.Indexes[0].DistanceMetric, "euclidean")
	assert.Equal(t, *result.Indexes[0].Status, "enable")
	assert.Equal(t, *result.Indexes[0].Mode, "fusion")
	assert.Equal(t, *result.Indexes[0].VectorBucketName, "oss-vector-bucket1")
	assert.Equal(t, *result.Indexes[1].Mode, "standard")
	assert.Equal(t, *result.Indexes[1].IndexName, "index2")
	assert.Equal(t, *result.Indexes[1].VectorBucketName, "oss-vector-bucket1")
	assert.Equal(t, *result.Indexes[1].Status, "enable")
	assert.Equal(t, *result.Indexes[1].CreateTime, time.Date(2025, time.August, 20, 10, 49, 17, 289372919, time.UTC))

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorIndexesResult{}
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
	result = &ListVectorIndexesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
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
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorIndexesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
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
	var input *oss.OperationInput
	var err error

	request = &DeleteVectorIndexRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &DeleteVectorIndexRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &DeleteVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("demo"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
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
	var output *oss.OperationOutput
	var err error
	output = &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Body:       io.NopCloser(bytes.NewReader([]byte(nil))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &DeleteVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &DeleteVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
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
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &DeleteVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_PutVectorIndexFusion(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutVectorIndexFusionRequest
	var input *oss.OperationInput
	var err error

	request = &PutVectorIndexFusionRequest{}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutVectorIndexFusionRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &PutVectorIndexFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Mode")

	request = &PutVectorIndexFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
		Mode:      oss.Ptr("fusion"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, SchemaConfiguration")

	request = &PutVectorIndexFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Mode:      oss.Ptr("fusion"),
		IndexName: oss.Ptr("exampleIndex"),
		SchemaConfiguration: &SchemaConfiguration{
			Fields: []SchemaField{
				{
					Name:           oss.Ptr("vector_1"),
					Type:           FieldTypeVector,
					DataType:       VectorDataTypeFloat32,
					Dimension:      oss.Ptr(1024),
					DistanceMetric: DistanceMetricTypeEuclidean,
				},
				{
					Name:           oss.Ptr("vector_2"),
					Type:           FieldTypeVector,
					DataType:       VectorDataTypeFloat32,
					Dimension:      oss.Ptr(512),
					DistanceMetric: DistanceMetricTypeCosine,
				},
				{
					Name:    oss.Ptr("timestamps"),
					Type:    FieldTypeLong,
					IsArray: oss.Ptr(true),
				},
				{
					Name: oss.Ptr("price"),
					Type: FieldTypeDouble,
				},
				{
					Name: oss.Ptr("ip"),
					Type: FieldTypeIp,
				},
				{
					Name: oss.Ptr("location"),
					Type: FieldTypeGeoPoint,
				},
				{
					Name: oss.Ptr("tag"),
					Type: FieldTypeString,
				},
				{
					Name:           oss.Ptr("user_id"),
					Type:           FieldTypeString,
					IsPartitionKey: oss.Ptr(true),
				},
				{
					Name:    oss.Ptr("tags"),
					Type:    FieldTypeString,
					IsArray: oss.Ptr(true),
				},
				{
					Name:       oss.Ptr("title_1"),
					Type:       FieldTypeString,
					ExactMatch: oss.Ptr(true),
					Text: &TextConfiguration{
						Enabled:  oss.Ptr(true),
						Analyzer: oss.Ptr("standard"),
						AnalyzerParameters: &AnalyzerParameters{
							CaseSensitive: oss.Ptr(true),
							DelimitWord:   oss.Ptr(false),
						},
					},
				},
				{
					Name:       oss.Ptr("title_2"),
					Type:       FieldTypeString,
					ExactMatch: oss.Ptr(false),
					Text: &TextConfiguration{
						Enabled:  oss.Ptr(true),
						Analyzer: oss.Ptr("split"),
						AnalyzerParameters: &AnalyzerParameters{
							CaseSensitive: oss.Ptr(true),
							Delimiter:     oss.Ptr(" "),
						},
					},
				},
			},
		},
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["putVectorIndex"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"exampleIndex\",\"mode\":\"fusion\",\"schemaConfiguration\":{\"fields\":[{\"name\":\"vector_1\",\"type\":\"vector\",\"dataType\":\"float32\",\"dimension\":1024,\"distanceMetric\":\"euclidean\"},{\"name\":\"vector_2\",\"type\":\"vector\",\"dataType\":\"float32\",\"dimension\":512,\"distanceMetric\":\"cosine\"},{\"name\":\"timestamps\",\"type\":\"long\",\"isArray\":true},{\"name\":\"price\",\"type\":\"double\"},{\"name\":\"ip\",\"type\":\"ip\"},{\"name\":\"location\",\"type\":\"geoPoint\"},{\"name\":\"tag\",\"type\":\"string\"},{\"name\":\"user_id\",\"type\":\"string\",\"isPartitionKey\":true},{\"name\":\"tags\",\"type\":\"string\",\"isArray\":true},{\"name\":\"title_1\",\"type\":\"string\",\"exactMatch\":true,\"text\":{\"enabled\":true,\"analyzer\":\"standard\",\"analyzerParameters\":{\"caseSensitive\":true,\"delimitWord\":false}}},{\"name\":\"title_2\",\"type\":\"string\",\"exactMatch\":false,\"text\":{\"enabled\":true,\"analyzer\":\"split\",\"analyzerParameters\":{\"caseSensitive\":true,\"delimiter\":\" \"}}}]}}")
}

func TestUnmarshalOutput_PutVectorIndexFusion(t *testing.T) {
	c := VectorsClient{}
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
	result := &PutVectorIndexFusionResult{}
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
	result = &PutVectorIndexFusionResult{}
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
	result = &PutVectorIndexFusionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
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
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexFusionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
