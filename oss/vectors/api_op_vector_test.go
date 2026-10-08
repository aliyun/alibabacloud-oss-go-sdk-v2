package vectors

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutVectorsRequest
	var input *oss.OperationInput
	var err error

	request = &PutVectorsRequest{}
	input = &oss.OperationInput{
		OpName: "PutVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutVectorsRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &PutVectorsRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
		Vectors: []map[string]any{
			{
				"key": "vector1",
				"data": map[string]any{
					"float32": []float32{1.2, 2.5, 3},
				},
				"metadata": map[string]any{
					"Key1": "value2",
					"Key2": []string{"1", "2", "3"},
				},
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "PutVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["PutVectors"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), `{"indexName":"exampleIndex","vectors":[{"data":{"float32":[1.2,2.5,3]},"key":"vector1","metadata":{"Key1":"value2","Key2":["1","2","3"]}}]}`)

	request = &PutVectorsRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
		Vectors: []map[string]any{
			{
				"key": "doc-001",
				"data": map[string]any{
					"vector_field_name_1": []float32{0.1, 0.2, 0.3, 0.4, 0.5},
				},
				"metadata": map[string]any{
					"title":    "Introduction to Vector Search",
					"category": []string{"technology", "ai"},
				},
			},
			{
				"key": "doc-002",
				"data": map[string]any{
					"vector_field_name_1": []float32{0.1, 0.2, 0.3, 0.4, 0.5},
					"vector_field_name_2": []float32{0.1, 0.2, 0.3, 0.4, 0.5},
				},
				"metadata": map[string]any{
					"title":    "Introduction to Vector Search",
					"category": []string{"technology", "ai"},
				},
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "PutVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["PutVectors"], "")
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), `{"indexName":"exampleIndex","vectors":[{"data":{"vector_field_name_1":[0.1,0.2,0.3,0.4,0.5]},"key":"doc-001","metadata":{"category":["technology","ai"],"title":"Introduction to Vector Search"}},{"data":{"vector_field_name_1":[0.1,0.2,0.3,0.4,0.5],"vector_field_name_2":[0.1,0.2,0.3,0.4,0.5]},"key":"doc-002","metadata":{"category":["technology","ai"],"title":"Introduction to Vector Search"}}]}`)
}

func TestUnmarshalOutput_PutVectors(t *testing.T) {
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
	result := &PutVectorsResult{}
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
	result = &PutVectorsResult{}
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
	result = &PutVectorsResult{}
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
	result = &PutVectorsResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *GetVectorsRequest
	var input *oss.OperationInput
	var err error

	request = &GetVectorsRequest{}
	input = &oss.OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"getVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &GetVectorsRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"getVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &GetVectorsRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("index"),
	}
	input = &oss.OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"getVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Keys")

	request = &GetVectorsRequest{
		Bucket:         oss.Ptr("oss-demo"),
		IndexName:      oss.Ptr("index"),
		Keys:           []string{"key1", "key2", "key3"},
		ReturnData:     oss.Ptr(true),
		ReturnMetadata: oss.Ptr(false),
	}
	input = &oss.OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"getVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["GetVectors"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"index\",\"keys\":[\"key1\",\"key2\",\"key3\"],\"returnData\":true,\"returnMetadata\":false}")
}

func TestUnmarshalOutput_GetVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "indexName": "index",
   "vectors": [ 
      { 
         "data": {
            "float32":[2.2]
         },
         "key": "key",
         "metadata": {
             "Key1": "value1",
             "Key2": "value2"
         }
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
	result := &GetVectorsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Vectors), 1)
	for _, vector := range result.Vectors {
		if keyVal, exists := vector["key"]; exists {
			keyStr, ok := keyVal.(string)
			assert.True(t, ok)
			assert.Equal(t, keyStr, "key")
		}

		// 访问 data 字段
		if dataVal, exists := vector["data"]; exists {
			dataMap, ok := dataVal.(map[string]any)
			assert.True(t, ok)
			if float32Val, exists := dataMap["float32"]; exists {
				float32Data, ok := float32Val.([]any)
				assert.True(t, ok)
				assert.Equal(t, float32Data[0], float64(2.2))
			}
		}

		if metadataVal, exists := vector["metadata"]; exists {
			metadataMap, ok := metadataVal.(map[string]any)
			assert.True(t, ok)
			if key1Val, exists := metadataMap["Key1"]; exists {
				key1Data, ok := key1Val.(string)
				assert.True(t, ok)
				assert.Equal(t, key1Data, "value1")
			}
			if key2Val, exists := metadataMap["Key2"]; exists {
				key2Data, ok := key2Val.(string)
				assert.True(t, ok)
				assert.Equal(t, key2Data, "value2")
			}
		}
	}

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetVectorsResult{}
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
	result = &GetVectorsResult{}
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
	result = &GetVectorsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_ListVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *ListVectorsRequest
	var input *oss.OperationInput
	var err error

	request = &ListVectorsRequest{}
	input = &oss.OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &ListVectorsRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &ListVectorsRequest{
		Bucket:         oss.Ptr("oss-demo"),
		IndexName:      oss.Ptr("index"),
		MaxResults:     100,
		NextToken:      oss.Ptr("123"),
		ReturnMetadata: oss.Ptr(true),
		ReturnData:     oss.Ptr(false),
		SegmentCount:   oss.Ptr(int(10)),
		SegmentIndex:   oss.Ptr(3),
	}
	input = &oss.OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectors"], "")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"index\",\"maxResults\":100,\"nextToken\":\"123\",\"returnData\":false,\"returnMetadata\":true,\"segmentCount\":10,\"segmentIndex\":3}")

	request = &ListVectorsRequest{
		Bucket:         oss.Ptr("oss-demo"),
		IndexName:      oss.Ptr("index"),
		ReturnMetadata: oss.Ptr(true),
		ReturnData:     oss.Ptr(false),
		SegmentCount:   oss.Ptr(int(10)),
		SegmentIndex:   oss.Ptr(3),
	}
	input = &oss.OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectors"], "")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"index\",\"returnData\":false,\"returnMetadata\":true,\"segmentCount\":10,\"segmentIndex\":3}")

}

func TestUnmarshalOutput_ListVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "nextToken": "123",
   "vectors": [ 
      { 
         "data": {
            "float32":[32]
         },
         "key": "key",
         "metadata": {
             "Key1": "value1",
             "Key2": "value2"
         }
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
	result := &ListVectorsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Vectors), 1)
	//assert.Equal(t, result.Vectors[0].Data.Float32[0], float32(32))
	//assert.Equal(t, *result.Vectors[0].Key, "key")
	//assert.Equal(t, (*result.Vectors[0].Metadata)["Key1"], "value1")
	//assert.Equal(t, (*result.Vectors[0].Metadata)["Key2"], "value2")
	assert.Equal(t, *result.NextToken, "123")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorsResult{}
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
	result = &ListVectorsResult{}
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
	result = &ListVectorsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *DeleteVectorsRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteVectorsRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &DeleteVectorsRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName.")

	request = &DeleteVectorsRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("index"),
		Keys: []string{
			"key1", "key2",
		},
	}
	input = &oss.OperationInput{
		OpName: "DeleteVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["DeleteVectors"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"index\",\"keys\":[\"key1\",\"key2\"]}")
}

func TestUnmarshalOutput_DeleteVectors(t *testing.T) {
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
	result := &DeleteVectorsResult{}
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
	result = &DeleteVectorsResult{}
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
	result = &DeleteVectorsResult{}
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
	result = &DeleteVectorsResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_QueryVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *QueryVectorsRequest
	var input *oss.OperationInput
	var err error

	request = &QueryVectorsRequest{}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &QueryVectorsRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &QueryVectorsRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("index"),
	}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, QueryVector")

	request = &QueryVectorsRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("index"),
		QueryVector: map[string]any{
			"float32": []float32{float32(32)},
		},
	}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TopK")

	request = &QueryVectorsRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("index"),
		QueryVector: map[string]any{
			"float32": []float32{float32(32)},
		},
		TopK: oss.Ptr(10),
	}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["QueryVectors"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"index\",\"queryVector\":{\"float32\":[32]},\"topK\":10}")

	request = &QueryVectorsRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("index"),
		Filter: map[string]any{
			"$and": []map[string]any{
				{
					"type": map[string]any{
						"$in": []string{"comedy", "documentary"},
					},
				},
			},
		},
		QueryVector: map[string]any{
			"float32": []float32{float32(32)},
		},
		ReturnMetadata: oss.Ptr(true),
		ReturnDistance: oss.Ptr(true),
		TopK:           oss.Ptr(10),
	}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["QueryVectors"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"filter\":{\"$and\":[{\"type\":{\"$in\":[\"comedy\",\"documentary\"]}}]},\"indexName\":\"index\",\"queryVector\":{\"float32\":[32]},\"returnDistance\":true,\"returnMetadata\":true,\"topK\":10}")
}

func TestUnmarshalOutput_QueryVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "vectors": [ 
      { 
         "data": {
            "float32":[32]
         },
         "key": "key",
         "metadata": {
             "key1": "value1",
             "key2": "value2"
         }
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
	result := &QueryVectorsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Vectors), 1)
	//assert.Equal(t, result.Vectors[0].Data.Float32[0], float32(32))
	//assert.Equal(t, *result.Vectors[0].Key, "key")
	//assert.Equal(t, (*result.Vectors[0].Metadata)["key1"], "value1")
	//assert.Equal(t, (*result.Vectors[0].Metadata)["key2"], "value2")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &QueryVectorsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
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
	result = &QueryVectorsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

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
	result = &QueryVectorsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_QueryVectorsFusion(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *QueryVectorsFusionRequest
	var input *oss.OperationInput
	var err error

	request = &QueryVectorsFusionRequest{}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &QueryVectorsFusionRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("index"),
	}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)

	request = &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("index"),
		Knn: []map[string]any{
			Knn{
				Field:         oss.Ptr("demo"),
				QueryVector:   map[string]any{"float32": []float32{float32(32)}},
				TopK:          oss.Ptr(10),
				NumCandidates: oss.Ptr(9),
				Filter: map[string]any{
					"meta_field_1": map[string]any{
						"$eq": "abc",
					},
				},
				Boost: oss.Ptr(float32(1)),
			}.ToMap(),
		},
	}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)

	request = &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("index"),
		Knn: []map[string]any{
			Knn{
				Field:         oss.Ptr("demo"),
				QueryVector:   map[string]any{"float32": []float32{float32(32)}},
				TopK:          oss.Ptr(10),
				NumCandidates: oss.Ptr(9),
				Filter: map[string]any{
					"meta_field_1": map[string]any{
						"$eq": "abc",
					},
				},
				Boost: oss.Ptr(float32(1)),
			}.ToMap(),
		},
		Retriever: map[string]any{
			"simple": SimpleRetriever{
				Query: map[string]any{
					"$and": []map[string]any{
						{"type": map[string]any{"$in": []string{"a", "b"}}},
						{"year": map[string]any{"$gte": 2020}},
					},
				},
			}.ToMap(),
		},
		ReturnMetadata:       oss.Ptr(true),
		ReturnMetadataFields: []string{"key1", "key2"},
		PartitionKeys:        []string{"key1", "key2"},
		Limit:                oss.Ptr(10),
		NextToken:            oss.Ptr("nextToken"),
		Sort: []Sort{
			{"field_a": SortOptions{Order: oss.Ptr(SortOrderTypeAsc)}},
			{"_score": SortOptions{Order: oss.Ptr(SortOrderTypeDesc)}},
			{"_primaryKey": SortOptions{Order: oss.Ptr(SortOrderTypeAsc)}},
		},
	}
	input = &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["queryVectorsFusion"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"index\",\"knn\":[{\"boost\":1,\"field\":\"demo\",\"filter\":{\"meta_field_1\":{\"$eq\":\"abc\"}},\"numCandidates\":9,\"queryVector\":{\"float32\":[32]},\"topK\":10}],\"limit\":10,\"nextToken\":\"nextToken\",\"partitionKeys\":[\"key1\",\"key2\"],\"retriever\":{\"simple\":{\"query\":{\"$and\":[{\"type\":{\"$in\":[\"a\",\"b\"]}},{\"year\":{\"$gte\":2020}}]}}},\"returnMetadata\":true,\"returnMetadataFields\":[\"key1\",\"key2\"],\"sort\":[{\"field_a\":{\"order\":\"asc\"}},{\"_score\":{\"order\":\"desc\"}},{\"_primaryKey\":{\"order\":\"asc\"}}]}")
}

func TestUnmarshalOutput_QueryVectorsFusion(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "vectors": [ 
      { 
         "data": {
            "float32":[32]
         },
         "key": "key",
         "metadata": {
             "key1": "value1",
             "key2": "value2"
         }
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
	result := &QueryVectorsFusionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Vectors), 1)
	v := result.Vectors[0]
	assert.Equal(t, *v.Key, "key")
	assert.Equal(t, v.Metadata["key1"], "value1")
	assert.Equal(t, v.Metadata["key2"], "value2")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &QueryVectorsFusionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
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
	result = &QueryVectorsFusionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

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
	result = &QueryVectorsFusionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func marshalQueryVectorsFusionBody(t *testing.T, request *QueryVectorsFusionRequest) string {
	t.Helper()
	c := VectorsClient{}
	input := &oss.OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	err := c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	return string(body)
}

// textMatchQuery builds the {"<field>":{"$textMatch":{...}}} query used by the retriever cases.
func textMatchQuery(field string, value string, boost *float32) map[string]any {
	textMatch := map[string]any{"value": value}
	if boost != nil {
		textMatch["boost"] = *boost
	}
	return map[string]any{field: map[string]any{"$textMatch": textMatch}}
}

func TestMarshalInput_QueryVectorsFusion_RetrieverSimple(t *testing.T) {
	request := &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("fusion-index"),
		Retriever: map[string]any{
			"simple": SimpleRetriever{
				Query: textMatchQuery("title_field", "hello world", oss.Ptr(float32(2))),
			}.ToMap(),
		},
	}
	assert.Equal(t, `{"indexName":"fusion-index","retriever":{"simple":{"query":{"title_field":{"$textMatch":{"boost":2,"value":"hello world"}}}}}}`, marshalQueryVectorsFusionBody(t, request))
}

func TestMarshalInput_QueryVectorsFusion_RetrieverKnnFromMap(t *testing.T) {
	request := &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("fusion-index"),
		Retriever: map[string]any{
			"knn": Knn{
				Field:       oss.Ptr("vector_field"),
				QueryVector: []int{10, 22, 77},
			}.ToMap(),
		},
	}
	assert.Equal(t, `{"indexName":"fusion-index","retriever":{"knn":{"field":"vector_field","queryVector":[10,22,77]}}}`, marshalQueryVectorsFusionBody(t, request))
}

func TestMarshalInput_QueryVectorsFusion_RetrieverRrf(t *testing.T) {
	request := &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("fusion-index"),
		Retriever: map[string]any{
			"rrf": RrfRetriever{
				K:          oss.Ptr(int32(50)),
				WindowSize: oss.Ptr(int32(100)),
				Retrievers: []map[string]any{
					RetrieverComponent{
						Retriever: map[string]any{
							"knn": Knn{Field: oss.Ptr("vector"), QueryVector: []int{10, 22, 77}}.ToMap(),
						},
						Weight: oss.Ptr(float32(1)),
					}.ToMap(),
					RetrieverComponent{
						Retriever: map[string]any{
							"simple": SimpleRetriever{
								Query: textMatchQuery("title", "hello world", oss.Ptr(float32(2))),
							}.ToMap(),
						},
						Weight: oss.Ptr(float32(2)),
					}.ToMap(),
				},
			}.ToMap(),
		},
	}
	assert.Equal(t, `{"indexName":"fusion-index","retriever":{"rrf":{"k":50,"retrievers":[{"retriever":{"knn":{"field":"vector","queryVector":[10,22,77]}},"weight":1},{"retriever":{"simple":{"query":{"title":{"$textMatch":{"boost":2,"value":"hello world"}}}}},"weight":2}],"windowSize":100}}}`, marshalQueryVectorsFusionBody(t, request))
}

func TestMarshalInput_QueryVectorsFusion_RetrieverWeight(t *testing.T) {
	request := &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("fusion-index"),
		Retriever: map[string]any{
			"weight": WeightRetriever{
				WindowSize: oss.Ptr(int32(100)),
				Retrievers: []map[string]any{
					RetrieverComponent{
						Retriever: map[string]any{
							"knn": Knn{Field: oss.Ptr("vector"), QueryVector: []int{10, 22, 77}}.ToMap(),
						},
						Weight:     oss.Ptr(float32(0.7)),
						Normalizer: oss.Ptr("minMax"),
					}.ToMap(),
					RetrieverComponent{
						Retriever: map[string]any{
							"simple": SimpleRetriever{
								Query: textMatchQuery("title", "hello world", oss.Ptr(float32(2))),
							}.ToMap(),
						},
						Weight:     oss.Ptr(float32(0.3)),
						Normalizer: oss.Ptr("minMax"),
					}.ToMap(),
				},
			}.ToMap(),
		},
	}
	assert.Equal(t, `{"indexName":"fusion-index","retriever":{"weight":{"retrievers":[{"normalizer":"minMax","retriever":{"knn":{"field":"vector","queryVector":[10,22,77]}},"weight":0.7},{"normalizer":"minMax","retriever":{"simple":{"query":{"title":{"$textMatch":{"boost":2,"value":"hello world"}}}}},"weight":0.3}],"windowSize":100}}}`, marshalQueryVectorsFusionBody(t, request))
}

func TestMarshalInput_QueryVectorsFusion_RetrieverNested(t *testing.T) {
	request := &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("fusion-index"),
		Retriever: map[string]any{
			"rrf": RrfRetriever{
				K:          oss.Ptr(int32(50)),
				WindowSize: oss.Ptr(int32(200)),
				Retrievers: []map[string]any{
					RetrieverComponent{
						Retriever: map[string]any{
							"simple": SimpleRetriever{
								Query: textMatchQuery("title", "hello world", nil),
							}.ToMap(),
						},
						Weight: oss.Ptr(float32(1)),
					}.ToMap(),
					RetrieverComponent{
						Retriever: map[string]any{
							"weight": WeightRetriever{
								WindowSize: oss.Ptr(int32(100)),
								Retrievers: []map[string]any{
									RetrieverComponent{
										Retriever: map[string]any{
											"knn": Knn{Field: oss.Ptr("text_vector"), QueryVector: []int{10, 22, 77}}.ToMap(),
										},
										Weight:     oss.Ptr(float32(0.7)),
										Normalizer: oss.Ptr("minMax"),
									}.ToMap(),
									RetrieverComponent{
										Retriever: map[string]any{
											"knn": Knn{Field: oss.Ptr("image_vector"), QueryVector: []int{21, 35, 66}}.ToMap(),
										},
										Weight:     oss.Ptr(float32(0.3)),
										Normalizer: oss.Ptr("minMax"),
									}.ToMap(),
								},
							}.ToMap(),
						},
						Weight: oss.Ptr(float32(1.2)),
					}.ToMap(),
				},
			}.ToMap(),
		},
	}
	assert.Equal(t, `{"indexName":"fusion-index","retriever":{"rrf":{"k":50,"retrievers":[{"retriever":{"simple":{"query":{"title":{"$textMatch":{"value":"hello world"}}}}},"weight":1},{"retriever":{"weight":{"retrievers":[{"normalizer":"minMax","retriever":{"knn":{"field":"text_vector","queryVector":[10,22,77]}},"weight":0.7},{"normalizer":"minMax","retriever":{"knn":{"field":"image_vector","queryVector":[21,35,66]}},"weight":0.3}],"windowSize":100}},"weight":1.2}],"windowSize":200}}}`, marshalQueryVectorsFusionBody(t, request))
}

func TestMarshalInput_QueryVectorsFusion_RetrieverNormalizerEnum(t *testing.T) {
	request := &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("fusion-index"),
		Retriever: map[string]any{
			"weight": WeightRetriever{
				WindowSize: oss.Ptr(int32(100)),
				Retrievers: []map[string]any{
					RetrieverComponent{
						Retriever: map[string]any{
							"knn": Knn{Field: oss.Ptr("text_vector"), QueryVector: []int{10, 22, 77}}.ToMap(),
						},
						Weight:     oss.Ptr(float32(0.7)),
						Normalizer: oss.Ptr(string(NormalizerTypeMinMax)),
					}.ToMap(),
					RetrieverComponent{
						Retriever: map[string]any{
							"knn": Knn{Field: oss.Ptr("image_vector"), QueryVector: []int{21, 35, 66}}.ToMap(),
						},
						Weight:     oss.Ptr(float32(0.3)),
						Normalizer: oss.Ptr(string(NormalizerTypeL2)),
					}.ToMap(),
				},
			}.ToMap(),
		},
	}
	assert.Equal(t, `{"indexName":"fusion-index","retriever":{"weight":{"retrievers":[{"normalizer":"minMax","retriever":{"knn":{"field":"text_vector","queryVector":[10,22,77]}},"weight":0.7},{"normalizer":"l2","retriever":{"knn":{"field":"image_vector","queryVector":[21,35,66]}},"weight":0.3}],"windowSize":100}}}`, marshalQueryVectorsFusionBody(t, request))
}

func TestMarshalInput_QueryVectorsFusion_RetrieverRawMapPassthrough(t *testing.T) {
	request := &QueryVectorsFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("fusion-index"),
		Retriever: map[string]any{
			"rrf": map[string]any{
				"k":           50,
				"futureParam": "x",
				"retrievers": []map[string]any{
					{
						"retriever":         map[string]any{"knn": map[string]any{"field": "vector", "queryVector": []int{10, 22, 77}}},
						"weight":            1.0,
						"futureWeightParam": 42,
					},
				},
			},
		},
	}
	assert.Equal(t, `{"indexName":"fusion-index","retriever":{"rrf":{"futureParam":"x","k":50,"retrievers":[{"futureWeightParam":42,"retriever":{"knn":{"field":"vector","queryVector":[10,22,77]}},"weight":1}]}}}`, marshalQueryVectorsFusionBody(t, request))
}
