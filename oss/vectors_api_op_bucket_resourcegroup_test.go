package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketResourceGroup_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutBucketResourceGroupRequest
	var input *OperationInput
	var err error

	request = &PutBucketResourceGroupRequest{}
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	input = &OperationInput{
		OpName: "PutBucketResourceGroup",
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
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"resourceGroup"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutBucketResourceGroupRequest{
		Bucket: Ptr("oss-demo"),
	}
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	input = &OperationInput{
		OpName: "PutBucketResourceGroup",
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
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, BucketResourceGroupConfiguration.")
	request = &PutBucketResourceGroupRequest{
		Bucket: Ptr("oss-demo"),
		BucketResourceGroupConfiguration: &BucketResourceGroupConfiguration{
			Ptr("rg-aekz****"),
		},
	}
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	input = &OperationInput{
		OpName: "PutBucketResourceGroup",
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
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"resourceGroup"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"BucketResourceGroupConfiguration\":{\"ResourceGroupId\":\"rg-aekz****\"}}")
}

func TestUnmarshalOutput_PutBucketResourceGroup_ForVectorBucket(t *testing.T) {
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
	result := &PutBucketResourceGroupResult{}
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
	result = &PutBucketResourceGroupResult{}
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
	result = &PutBucketResourceGroupResult{}
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
	result = &PutBucketResourceGroupResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetBucketResourceGroup_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *GetBucketResourceGroupRequest
	var input *OperationInput
	var err error

	request = &GetBucketResourceGroupRequest{}
	input = &OperationInput{
		OpName: "GetBucketResourceGroup",
		Method: "GET",
		Parameters: map[string]string{
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"resourceGroup"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketResourceGroupRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketResourceGroup",
		Method: "GET",
		Parameters: map[string]string{
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"resourceGroup"})
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketResourceGroup_ForVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	putBody := `{
  "BucketResourceGroupConfiguration": {
    "ResourceGroupId": "rg-aekz****"
  }
}`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	assert.Nil(t, err)
	result := &GetBucketResourceGroupResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyXmlOrJson)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.BucketResourceGroupConfiguration.ResourceGroupId, "rg-aekz****")

	putBody = `{
 "Error": {
   "Code": "NoSuchBucket",
   "Message": "The specified bucket does not exist.",
   "RequestId": "66C2FF09FDF07830343C72EC",
   "HostId": "bucket.oss-cn-hangzhou.aliyuncs.com",
   "BucketName": "bucket",
   "EC": "0015-00000101"
 }
}`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketResourceGroupResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyXmlOrJson)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	putBody = `{
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
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketResourceGroupResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyXmlOrJson)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
