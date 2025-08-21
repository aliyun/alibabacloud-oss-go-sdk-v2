package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketTags_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutBucketTagsRequest
	var input *OperationInput
	var err error

	request = &PutBucketTagsRequest{}
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	input = &OperationInput{
		OpName: "PutBucketTags",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: func() string {
				if request.Headers != nil && request.Headers[HTTPHeaderContentType] != "" {
					return request.Headers[HTTPHeaderContentType]
				}
				return contentTypeXML
			}(),
		},
		Parameters: map[string]string{
			"tagging": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"tagging"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutBucketTagsRequest{
		Bucket: Ptr("oss-demo"),
	}
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	input = &OperationInput{
		OpName: "PutBucketTags",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: func() string {
				if request.Headers != nil && request.Headers[HTTPHeaderContentType] != "" {
					return request.Headers[HTTPHeaderContentType]
				}
				return contentTypeXML
			}(),
		},
		Parameters: map[string]string{
			"tagging": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"tagging"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, Tagging.")

	request = &PutBucketTagsRequest{
		Bucket: Ptr("oss-demo"),
		Tagging: &Tagging{
			&TagSet{
				[]Tag{
					{
						Ptr("key1"),
						Ptr("value1"),
					},
					{
						Ptr("key2"),
						Ptr("value2"),
					},
				},
			},
		},
	}
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	input = &OperationInput{
		OpName: "PutBucketTags",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: func() string {
				if request.Headers != nil && request.Headers[HTTPHeaderContentType] != "" {
					return request.Headers[HTTPHeaderContentType]
				}
				return contentTypeXML
			}(),
		},
		Parameters: map[string]string{
			"tagging": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"tagging"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"Tagging\":{\"TagSet\":{\"Tag\":[{\"Key\":\"key1\",\"Value\":\"value1\"},{\"Key\":\"key2\",\"Value\":\"value2\"}]}}}")

	request = &PutBucketTagsRequest{
		Bucket: Ptr("oss-demo"),
		Tagging: &Tagging{
			&TagSet{
				[]Tag{
					{
						Ptr("key1"),
						Ptr("value1"),
					},
				},
			},
		},
	}
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	input = &OperationInput{
		OpName: "PutBucketTags",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: func() string {
				if request.Headers != nil && request.Headers[HTTPHeaderContentType] != "" {
					return request.Headers[HTTPHeaderContentType]
				}
				return contentTypeXML
			}(),
		},
		Parameters: map[string]string{
			"tagging": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"tagging"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"Tagging\":{\"TagSet\":{\"Tag\":[{\"Key\":\"key1\",\"Value\":\"value1\"}]}}}")
}

func TestUnmarshalOutput_PutBucketTags_ForVectorBucket(t *testing.T) {
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
	result := &PutBucketTagsResult{}
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
	result = &PutBucketTagsResult{}
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
	result = &PutBucketTagsResult{}
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
	result = &PutBucketTagsResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetBucketTags_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *GetBucketTagsRequest
	var input *OperationInput
	var err error

	request = &GetBucketTagsRequest{}
	input = &OperationInput{
		OpName: "GetBucketTags",
		Method: "GET",
		Parameters: map[string]string{
			"tagging": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"tagging"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketTagsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketTags",
		Method: "GET",
		Parameters: map[string]string{
			"tagging": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"tagging"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketTags_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `{
  "Tagging": {
    "TagSet": {
      "Tag": [
        {
          "Key": "testa",
          "Value": "testv1"
        },
        {
          "Key": "testb",
          "Value": "testv2"
        }
      ]
    }
  }
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
	result := &GetBucketTagsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyXmlOrJson)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	assert.Equal(t, len(result.Tagging.TagSet.Tags), 2)
	assert.Equal(t, *result.Tagging.TagSet.Tags[0].Key, "testa")
	assert.Equal(t, *result.Tagging.TagSet.Tags[1].Value, "testv2")

	body = `{
  "Tagging": {
    "TagSet": {
      "Tag": [
        {
          "Key": "testa",
          "Value": "testv1"
        }
      ]
    }
  }
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
	result = &GetBucketTagsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyXmlOrJson)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	assert.Equal(t, len(result.Tagging.TagSet.Tags), 1)
	assert.Equal(t, *result.Tagging.TagSet.Tags[0].Key, "testa")
	assert.Equal(t, *result.Tagging.TagSet.Tags[0].Value, "testv1")


	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetBucketTagsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyXmlOrJson)
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
	result = &GetBucketTagsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyXmlOrJson)
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
	result = &GetBucketTagsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyXmlOrJson)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteBucketTags_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *DeleteBucketTagsRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketTagsRequest{}
	input = &OperationInput{
		OpName: "DeleteBucketTags",
		Method: "DELETE",
		Parameters: map[string]string{
			"tagging": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"tagging"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteBucketTagsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketTags",
		Method: "DELETE",
		Parameters: map[string]string{
			"tagging": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"tagging"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &DeleteBucketTagsRequest{
		Bucket:  Ptr("oss-demo"),
		Tagging: Ptr("k1,k2"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketTags",
		Method: "DELETE",
		Parameters: map[string]string{
			"tagging": "",
		},
		Bucket: request.Bucket,
	}
	if request.Tagging != nil {
		input.Parameters["tagging"] = *request.Tagging
	}
	input.OpMetadata.Set(signer.SubResource, []string{"tagging"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["tagging"], "k1,k2")
}

func TestUnmarshalOutput_DeleteBucketTags_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &DeleteBucketTagsResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &DeleteBucketTagsResult{}
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
	result = &DeleteBucketTagsResult{}
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
	result = &DeleteBucketTagsResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
