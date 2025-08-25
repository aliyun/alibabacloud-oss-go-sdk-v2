package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutVectorsRequest
	var input *OperationInput
	var err error

	request = &PutVectorsRequest{}
	input = &OperationInput{
		OpName: "PutVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"PutVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutVectorsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"PutVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &PutVectorsRequest{
		Bucket:    Ptr("oss-demo"),
		IndexName: Ptr("exampleIndex"),
		Vectors: []map[string]interface{}{
			{
				"key": "vector1",
				"data": map[string]interface{}{
					"float32": []float32{1.2, 2.5, 3},
				},
				"metadata": map[string]interface{}{
					"Key1": 32,
					"Key2": "value2",
					"Key3": []string{"1", "2", "3"},
					"Key4": false,
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"PutVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["PutVectors"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), `{"indexName":"exampleIndex","vectors":[{"data":{"float32":[1.2,2.5,3]},"key":"vector1","metadata":{"Key1":32,"Key2":"value2","Key3":["1","2","3"],"Key4":false}}]}`)
}

func TestUnmarshalOutput_PutVectors(t *testing.T) {
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
	result := &PutVectorsResult{}
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
	result = &PutVectorsResult{}
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
	result = &PutVectorsResult{}
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
	result = &PutVectorsResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
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
	var input *OperationInput
	var err error

	request = &GetVectorsRequest{}
	input = &OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"GetVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &GetVectorsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"GetVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &GetVectorsRequest{
		Bucket:    Ptr("oss-demo"),
		IndexName: Ptr("index"),
	}
	input = &OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"GetVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Keys")

	request = &GetVectorsRequest{
		Bucket:         Ptr("oss-demo"),
		IndexName:      Ptr("index"),
		Keys:           []string{"key1", "key2", "key3"},
		ReturnData:     Ptr(true),
		ReturnMetadata: Ptr(false),
	}
	input = &OperationInput{
		OpName: "GetVectors",
		Method: "POST",
		Parameters: map[string]string{
			"GetVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
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
	var output *OperationOutput
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
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &GetVectorsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
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
			dataMap, ok := dataVal.(map[string]interface{})
			assert.True(t, ok)
			if float32Val, exists := dataMap["float32"]; exists {
				float32Data, ok := float32Val.([]interface{})
				assert.True(t, ok)
				assert.Equal(t, float32Data[0], float64(2.2))
			}
		}

		if metadataVal, exists := vector["metadata"]; exists {
			metadataMap, ok := metadataVal.(map[string]interface{})
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

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetVectorsResult{}
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
	result = &GetVectorsResult{}
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
	result = &GetVectorsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
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
	var input *OperationInput
	var err error

	request = &ListVectorsRequest{}
	input = &OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &ListVectorsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &ListVectorsRequest{
		Bucket:         Ptr("oss-demo"),
		IndexName:      Ptr("index"),
		MaxResults:     Ptr(100),
		NextToken:      Ptr("123"),
		ReturnMetadata: Ptr(true),
		ReturnData:     Ptr(false),
		SegmentCount:   Ptr(int(10)),
		SegmentIndex:   Ptr(3),
	}
	input = &OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectors"], "")
	assert.Equal(t, input.Headers[HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[HTTPHeaderContentType], contentTypeJSON)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"ReturnMetadata\":true,\"SegmentCount\":10,\"SegmentIndex\":3,\"indexName\":\"index\",\"maxResults\":100,\"nextToken\":\"123\",\"returnData\":false}")

	request = &ListVectorsRequest{
		Bucket:         Ptr("oss-demo"),
		IndexName:      Ptr("index"),
		ReturnMetadata: Ptr(true),
		ReturnData:     Ptr(false),
		SegmentCount:   Ptr(int(10)),
		SegmentIndex:   Ptr(3),
	}
	input = &OperationInput{
		OpName: "ListVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"ListVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectors"], "")
	assert.Equal(t, input.Headers[HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[HTTPHeaderContentType], contentTypeJSON)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"ReturnMetadata\":true,\"SegmentCount\":10,\"SegmentIndex\":3,\"indexName\":\"index\",\"returnData\":false}")

}

func TestUnmarshalOutput_ListVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
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
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &ListVectorsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
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

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorsResult{}
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
	result = &ListVectorsResult{}
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
	result = &ListVectorsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
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
	var input *OperationInput
	var err error

	request = &DeleteVectorsRequest{}
	input = &OperationInput{
		OpName: "DeleteVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"DeleteVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &DeleteVectorsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"DeleteVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName.")

	request = &DeleteVectorsRequest{
		Bucket:    Ptr("oss-demo"),
		IndexName: Ptr("index"),
		Keys: []string{
			"key1", "key2",
		},
	}
	input = &OperationInput{
		OpName: "DeleteVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"DeleteVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
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
	result := &DeleteVectorsResult{}
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
	result = &DeleteVectorsResult{}
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
	result = &DeleteVectorsResult{}
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
	result = &DeleteVectorsResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
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
	var input *OperationInput
	var err error

	request = &QueryVectorsRequest{}
	input = &OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"QueryVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &QueryVectorsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"QueryVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &QueryVectorsRequest{
		Bucket:    Ptr("oss-demo"),
		IndexName: Ptr("index"),
		Filter:    Ptr(`{"$and":[{"type":{"$in":["comedy","documentary"]}},{"year":{"$gte":2020}}]}`),
		QueryVector: map[string]interface{}{
			"float32": []float32{float32(32)},
		},
		ReturnMetadata: Ptr(true),
		ReturnDistance: Ptr(true),
		TopK:           Ptr(10),
	}
	input = &OperationInput{
		OpName: "QueryVectors",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"QueryVectors": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"QueryVectors"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["QueryVectors"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"filter\":\"{\\\"$and\\\":[{\\\"type\\\":{\\\"$in\\\":[\\\"comedy\\\",\\\"documentary\\\"]}},{\\\"year\\\":{\\\"$gte\\\":2020}}]}\",\"indexName\":\"index\",\"queryVector\":{\"float32\":[32]},\"returnDistance\":true,\"returnMetadata\":true,\"topK\":10}")
}

func TestUnmarshalOutput_QueryVectors(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
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
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &QueryVectorsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
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

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &QueryVectorsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
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
	result = &QueryVectorsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
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
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &QueryVectorsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefault)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
