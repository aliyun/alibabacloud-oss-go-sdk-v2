package vectors

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketResourceGroup(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutBucketResourceGroupRequest
	var input *oss.OperationInput
	var err error

	request = &PutBucketResourceGroupRequest{}
	input = &oss.OperationInput{
		OpName: "PutBucketResourceGroup",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutBucketResourceGroupRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "PutBucketResourceGroup",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, BucketResourceGroupConfiguration.")

	request = &PutBucketResourceGroupRequest{
		Bucket: oss.Ptr("oss-demo"),
		BucketResourceGroupConfiguration: &BucketResourceGroupConfiguration{
			ResourceGroupId: oss.Ptr("rg-aekz****"),
		},
	}
	input = &oss.OperationInput{
		OpName: "PutBucketResourceGroup",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"BucketResourceGroupConfiguration\":{\"ResourceGroupId\":\"rg-aekz****\"}}")
}

func TestUnmarshalOutput_PutBucketResourceGroup(t *testing.T) {
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
	result := &PutBucketResourceGroupResult{}
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
	result = &PutBucketResourceGroupResult{}
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
	result = &PutBucketResourceGroupResult{}
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
	result = &PutBucketResourceGroupResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetBucketResourceGroup(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *GetBucketResourceGroupRequest
	var input *oss.OperationInput
	var err error

	request = &GetBucketResourceGroupRequest{}
	input = &oss.OperationInput{
		OpName: "GetBucketResourceGroup",
		Method: "GET",
		Parameters: map[string]string{
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketResourceGroupRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "GetBucketResourceGroup",
		Method: "GET",
		Parameters: map[string]string{
			"resourceGroup": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketResourceGroup(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	putBody := `{
  "BucketResourceGroupConfiguration": {
    "ResourceGroupId": "rg-aekz****"
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &GetBucketResourceGroupResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.BucketResourceGroupConfiguration.ResourceGroupId, "rg-aekz****")

	// without BucketResourceGroupConfiguration root name
	putBody = `{
  "InvalidRoot": {
    "ResourceGroupId": "rg-aekz****"
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetBucketResourceGroupResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "deserialization failed, Found key InvalidRoot, but expect BucketResourceGroupConfiguration")

	// invalid json
	putBody = `
  "BucketResourceGroupConfiguration": {
    "ResourceGroupId": "rg-aekz****"
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetBucketResourceGroupResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "deserialization failed, invalid character")
}
