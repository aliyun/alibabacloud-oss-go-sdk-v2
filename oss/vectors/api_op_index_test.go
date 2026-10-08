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
	fields := result.Index.SchemaConfiguration.Fields
	assert.Len(t, fields, 11)

	assert.Equal(t, fields[0]["name"], "vector_1")
	assert.Equal(t, fields[0]["type"], "vector")
	assert.Equal(t, fields[0]["dataType"], "float32")
	assert.Equal(t, fields[0]["dimension"], float64(1024))
	assert.Equal(t, fields[0]["distanceMetric"], "euclidean")

	assert.Equal(t, fields[1]["name"], "vector_2")
	assert.Equal(t, fields[1]["type"], "vector")
	assert.Equal(t, fields[1]["dataType"], "float32")
	assert.Equal(t, fields[1]["dimension"], float64(512))
	assert.Equal(t, fields[1]["distanceMetric"], "cosine")

	assert.Equal(t, fields[2]["name"], "timestamps")
	assert.Equal(t, fields[2]["type"], "long")
	assert.Equal(t, fields[2]["isArray"], true)

	assert.Equal(t, fields[3]["name"], "price")
	assert.Equal(t, fields[3]["type"], "double")

	assert.Equal(t, fields[4]["name"], "ip")
	assert.Equal(t, fields[4]["type"], "ip")

	assert.Equal(t, fields[5]["name"], "location")
	assert.Equal(t, fields[5]["type"], "geoPoint")

	assert.Equal(t, fields[6]["name"], "tag")
	assert.Equal(t, fields[6]["type"], "string")

	assert.Equal(t, fields[7]["name"], "user_id")
	assert.Equal(t, fields[7]["type"], "string")
	assert.Equal(t, fields[7]["isPartitionKey"], true)

	assert.Equal(t, fields[8]["name"], "tags")
	assert.Equal(t, fields[8]["type"], "string")
	assert.Equal(t, fields[8]["isArray"], true)

	assert.Equal(t, fields[9]["name"], "title_1")
	assert.Equal(t, fields[9]["type"], "string")
	assert.Equal(t, fields[9]["exactMatch"], true)
	text1, ok := fields[9]["text"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, text1["enabled"], true)
	assert.Equal(t, text1["analyzer"], "standard")
	parameters1, ok := text1["analyzerParameters"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, parameters1["caseSensitive"], true)
	assert.Equal(t, parameters1["delimitWord"], false)

	assert.Equal(t, fields[10]["name"], "title_2")
	assert.Equal(t, fields[10]["type"], "string")
	assert.Equal(t, fields[10]["exactMatch"], false)
	text2, ok := fields[10]["text"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, text2["enabled"], true)
	assert.Equal(t, text2["analyzer"], "split")
	parameters2, ok := text2["analyzerParameters"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, parameters2["caseSensitive"], true)
	assert.Equal(t, parameters2["delimiter"], " ")

	schemas, err := result.Index.SchemaConfiguration.FieldSchemas()
	assert.Nil(t, err)
	assert.Len(t, schemas, 11)
	assert.Equal(t, *schemas[0].Name, "vector_1")
	assert.Equal(t, *schemas[0].Type, "vector")
	assert.Equal(t, *schemas[0].DataType, "float32")
	assert.Equal(t, *schemas[0].Dimension, 1024)
	assert.Equal(t, *schemas[0].DistanceMetric, "euclidean")
	assert.Equal(t, *schemas[9].Text.Enabled, true)
	assert.Equal(t, *schemas[9].Text.Analyzer, "standard")
	assert.Equal(t, *schemas[9].Text.AnalyzerParameters.DelimitWord, false)
	assert.Equal(t, *schemas[10].Text.AnalyzerParameters.Delimiter, " ")

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

func TestSchemaConfiguration_FieldSchemas(t *testing.T) {
	configuration := &SchemaConfiguration{
		Fields: []map[string]any{
			FieldSchema{
				Name:           oss.Ptr("vector_1"),
				Type:           oss.Ptr("vector"),
				DataType:       oss.Ptr("float32"),
				Dimension:      oss.Ptr(1024),
				DistanceMetric: oss.Ptr("euclidean"),
				Text: &TextSchema{
					Analyzer: oss.Ptr("standard"),
				},
			}.ToMap(),
			{
				"name": "tag",
				"type": "string",
			},
		},
	}
	schemas, err := configuration.FieldSchemas()
	assert.Nil(t, err)
	assert.Len(t, schemas, 2)
	assert.Equal(t, *schemas[0].Name, "vector_1")
	assert.Equal(t, *schemas[0].Type, "vector")
	assert.Equal(t, *schemas[0].DataType, "float32")
	assert.Equal(t, *schemas[0].Dimension, 1024)
	assert.Equal(t, *schemas[0].DistanceMetric, "euclidean")
	assert.Equal(t, *schemas[0].Text.Analyzer, "standard")
	assert.Equal(t, *schemas[1].Name, "tag")
	assert.Equal(t, *schemas[1].Type, "string")

	empty, err := (&SchemaConfiguration{}).FieldSchemas()
	assert.Nil(t, err)
	assert.Nil(t, empty)

	invalid, err := (&SchemaConfiguration{
		Fields: []map[string]any{{"dimension": "not-a-number"}},
	}).FieldSchemas()
	assert.NotNil(t, err)
	assert.Nil(t, invalid)
}

func TestFieldSchemas_ToMaps(t *testing.T) {
	schemas := FieldSchemas{
		{
			Name:           oss.Ptr("vector_1"),
			Type:           oss.Ptr("vector"),
			DataType:       oss.Ptr("float32"),
			Dimension:      oss.Ptr(1024),
			DistanceMetric: oss.Ptr("euclidean"),
			Text: &TextSchema{
				Analyzer: oss.Ptr("standard"),
			},
		},
		{
			Name: oss.Ptr("tag"),
			Type: oss.Ptr("string"),
		},
	}

	assert.Equal(t, []map[string]any{
		{
			"name":           "vector_1",
			"type":           "vector",
			"dataType":       "float32",
			"dimension":      1024,
			"distanceMetric": "euclidean",
			"text":           map[string]any{"analyzer": "standard"},
		},
		{
			"name": "tag",
			"type": "string",
		},
	}, schemas.ToMaps())

	assert.Empty(t, FieldSchemas(nil).ToMaps())
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
			Fields: []map[string]any{
				FieldSchema{
					Name:           oss.Ptr("vector_1"),
					Type:           oss.Ptr("vector"),
					DataType:       oss.Ptr("float32"),
					Dimension:      oss.Ptr(1024),
					DistanceMetric: oss.Ptr("euclidean"),
				}.ToMap(),
				FieldSchema{
					Name:           oss.Ptr("vector_2"),
					Type:           oss.Ptr("vector"),
					DataType:       oss.Ptr("float32"),
					Dimension:      oss.Ptr(512),
					DistanceMetric: oss.Ptr("cosine"),
				}.ToMap(),
				FieldSchema{
					Name:    oss.Ptr("timestamps"),
					Type:    oss.Ptr("long"),
					IsArray: oss.Ptr(true),
				}.ToMap(),
				FieldSchema{
					Name: oss.Ptr("price"),
					Type: oss.Ptr("double"),
				}.ToMap(),
				FieldSchema{
					Name: oss.Ptr("ip"),
					Type: oss.Ptr("ip"),
				}.ToMap(),
				FieldSchema{
					Name: oss.Ptr("location"),
					Type: oss.Ptr("geoPoint"),
				}.ToMap(),
				FieldSchema{
					Name: oss.Ptr("tag"),
					Type: oss.Ptr("string"),
				}.ToMap(),
				FieldSchema{
					Name:           oss.Ptr("user_id"),
					Type:           oss.Ptr("string"),
					IsPartitionKey: oss.Ptr(true),
				}.ToMap(),
				FieldSchema{
					Name:    oss.Ptr("tags"),
					Type:    oss.Ptr("string"),
					IsArray: oss.Ptr(true),
				}.ToMap(),
				FieldSchema{
					Name:       oss.Ptr("title_1"),
					Type:       oss.Ptr("string"),
					ExactMatch: oss.Ptr(true),
					Text: &TextSchema{
						Enabled:  oss.Ptr(true),
						Analyzer: oss.Ptr("standard"),
						AnalyzerParameters: &AnalyzerParameters{
							CaseSensitive: oss.Ptr(true),
							DelimitWord:   oss.Ptr(false),
						},
					},
				}.ToMap(),
				FieldSchema{
					Name:       oss.Ptr("title_2"),
					Type:       oss.Ptr("string"),
					ExactMatch: oss.Ptr(false),
					Text: &TextSchema{
						Enabled:  oss.Ptr(true),
						Analyzer: oss.Ptr("split"),
						AnalyzerParameters: &AnalyzerParameters{
							CaseSensitive: oss.Ptr(true),
							Delimiter:     oss.Ptr(" "),
						},
					},
				}.ToMap(),
			},
		},
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["putVectorIndex"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"exampleIndex\",\"mode\":\"fusion\",\"schemaConfiguration\":{\"fields\":[{\"dataType\":\"float32\",\"dimension\":1024,\"distanceMetric\":\"euclidean\",\"name\":\"vector_1\",\"type\":\"vector\"},{\"dataType\":\"float32\",\"dimension\":512,\"distanceMetric\":\"cosine\",\"name\":\"vector_2\",\"type\":\"vector\"},{\"isArray\":true,\"name\":\"timestamps\",\"type\":\"long\"},{\"name\":\"price\",\"type\":\"double\"},{\"name\":\"ip\",\"type\":\"ip\"},{\"name\":\"location\",\"type\":\"geoPoint\"},{\"name\":\"tag\",\"type\":\"string\"},{\"isPartitionKey\":true,\"name\":\"user_id\",\"type\":\"string\"},{\"isArray\":true,\"name\":\"tags\",\"type\":\"string\"},{\"exactMatch\":true,\"name\":\"title_1\",\"text\":{\"analyzer\":\"standard\",\"analyzerParameters\":{\"caseSensitive\":true,\"delimitWord\":false},\"enabled\":true},\"type\":\"string\"},{\"exactMatch\":false,\"name\":\"title_2\",\"text\":{\"analyzer\":\"split\",\"analyzerParameters\":{\"caseSensitive\":true,\"delimiter\":\" \"},\"enabled\":true},\"type\":\"string\"}]}}")
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
