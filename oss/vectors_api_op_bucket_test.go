package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutVectorBucketRequest
	var input *OperationInput
	var err error

	request = &PutVectorBucketRequest{}
	input = &OperationInput{
		OpName: "PutVectorBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutVectorBucketRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutVectorBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)

	request = &PutVectorBucketRequest{
		Bucket:          Ptr("oss-demo"),
		ResourceGroupId: Ptr("rg-aek27tc****"),
		Tagging:         Ptr("k1=v1&k2=v2"),
	}
	input = &OperationInput{
		OpName: "PutVectorBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-resource-group-id"], "rg-aek27tc****")
	assert.Equal(t, input.Headers["x-oss-bucket-tagging"], "k1=v1&k2=v2")
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
}

func TestUnmarshalOutput_PutVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutVectorBucketResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	output = &OperationOutput{
		StatusCode: 409,
		Status:     "BucketAlreadyExist",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &PutVectorBucketResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "BucketAlreadyExist")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *GetVectorBucketRequest
	var input *OperationInput
	var err error

	request = &GetVectorBucketRequest{}
	input = &OperationInput{
		OpName: "GetVectorBucket",
		Method: "GET",
		Parameters: map[string]string{
			"bucketInfo": "",
		},
		Bucket: request.Bucket,
	}
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetVectorBucketRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetVectorBucket",
		Method: "GET",
		Parameters: map[string]string{
			"bucketInfo": "",
		},
		Bucket: request.Bucket,
	}
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `{
  "BucketInfo": {
    "Bucket": {
      "CreationDate": "2013-07-31T10:56:21.000Z",
      "ExtranetEndpoint": "oss-cn-hangzhou.aliyuncs.com",
      "IntranetEndpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
      "Location": "oss-cn-hangzhou",
      "Name": "oss-example",
      "ResourceGroupId": "rg-aek27tc********"
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
	result := &GetVectorBucketResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefaultV2)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.BucketInfo.Name, "oss-example")
	assert.Equal(t, *result.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.Location, "oss-cn-hangzhou")
	assert.Equal(t, *result.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
	assert.Equal(t, *result.BucketInfo.ResourceGroupId, "rg-aek27tc********")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefaultV2)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefaultV2)
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
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/json"},
		},
	}
	resultErr := &GetVectorBucketResult{}
	err = c.client.unmarshalOutput(resultErr, output, unmarshalBodyDefaultV2)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_ListVectorBuckets(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *ListVectorBucketsRequest
	var input *OperationInput
	var err error

	request = &ListVectorBucketsRequest{}
	input = &OperationInput{
		OpName: "ListVectorBuckets",
		Method: "GET",
	}
	err = c.client.marshalInput(request, input)
	assert.Nil(t, err)

	request = &ListVectorBucketsRequest{
		Marker:          Ptr(""),
		MaxKeys:         10,
		Prefix:          Ptr("/"),
		ResourceGroupId: Ptr("rg-aek27tc********"),
	}
	input = &OperationInput{
		OpName: "ListVectorBuckets",
		Method: "GET",
	}
	err = c.client.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["max-keys"], "10")
	assert.Equal(t, input.Parameters["prefix"], "/")
	assert.Equal(t, input.Parameters["marker"], "")
	assert.Equal(t, input.Headers["x-oss-resource-group-id"], "rg-aek27tc********")
}

func TestUnmarshalOutput_ListVectorBuckets(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	body := `{
  "ListAllMyBucketsResult": {
    "Buckets": {
      "Bucket": [
        {
          "CreationDate": "2014-02-17T18:12:43.000Z",
          "ExtranetEndpoint": "oss-cn-shanghai.aliyuncs.com",
          "IntranetEndpoint": "oss-cn-shanghai-internal.aliyuncs.com",
          "Location": "oss-cn-shanghai",
          "Name": "app-base-oss",
          "Region": "cn-shanghai",
          "ResourceGroupId": "rg-aek27ta********"
        },
        {
          "CreationDate": "2014-02-25T11:21:04.000Z",
          "ExtranetEndpoint": "oss-cn-hangzhou.aliyuncs.com",
          "IntranetEndpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
          "Location": "oss-cn-hangzhou",
          "Name": "mybucket",
          "Region": "cn-hangzhou",
          "ResourceGroupId": "rg-aek27tc********"
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
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"5374A2880232A65C2300****"},
		},
	}
	result := &ListVectorBucketsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefaultV2)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5374A2880232A65C2300****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, len(result.Buckets.Bucket), 2)
	assert.Equal(t, *result.Buckets.Bucket[0].CreationDate, time.Date(2014, time.February, 17, 18, 12, 43, 0, time.UTC))
	assert.Equal(t, *result.Buckets.Bucket[0].ExtranetEndpoint, "oss-cn-shanghai.aliyuncs.com")
	assert.Equal(t, *result.Buckets.Bucket[0].IntranetEndpoint, "oss-cn-shanghai-internal.aliyuncs.com")
	assert.Equal(t, *result.Buckets.Bucket[0].Name, "app-base-oss")
	assert.Equal(t, *result.Buckets.Bucket[0].Region, "cn-shanghai")
	assert.Equal(t, *result.Buckets.Bucket[0].Location, "oss-cn-shanghai")
	assert.Equal(t, *result.Buckets.Bucket[0].ResourceGroupId, "rg-aek27ta********")

	assert.Equal(t, *result.Buckets.Bucket[1].CreationDate, time.Date(2014, time.February, 25, 11, 21, 04, 0, time.UTC))
	assert.Equal(t, *result.Buckets.Bucket[1].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.Buckets.Bucket[1].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.Buckets.Bucket[1].Name, "mybucket")
	assert.Equal(t, *result.Buckets.Bucket[1].Region, "cn-hangzhou")
	assert.Equal(t, *result.Buckets.Bucket[1].Location, "oss-cn-hangzhou")
	assert.Equal(t, *result.Buckets.Bucket[1].ResourceGroupId, "rg-aek27tc********")

	body = `{
  "ListAllMyBucketsResult": {
    "Prefix": "my",
    "Marker": "mybucket",
    "MaxKeys": 10,
    "IsTruncated": true,
    "NextMarker": "mybucket10",
    "Buckets": {
      "Bucket": [{
        "CreationDate": "2014-05-14T11:18:32.000Z",
        "ExtranetEndpoint": "oss-cn-hangzhou.aliyuncs.com",
        "IntranetEndpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
        "Location": "oss-cn-hangzhou",
        "Name": "mybucket01",
        "Region": "cn-hangzhou",
        "ResourceGroupId": "rg-aek27tc********"
      }]
    }
  }
}`

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"5374A2880232A65C2300****"},
		},
	}
	result = &ListVectorBucketsResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefaultV2)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5374A2880232A65C2300****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.Prefix, "my")
	assert.Equal(t, *result.Marker, "mybucket")
	assert.Equal(t, result.MaxKeys, int32(10))
	assert.Equal(t, result.IsTruncated, true)
	assert.Equal(t, *result.NextMarker, "mybucket10")
	assert.Equal(t, len(result.Buckets.Bucket), 1)

	assert.Equal(t, *result.Buckets.Bucket[0].CreationDate, time.Date(2014, time.May, 14, 11, 18, 32, 0, time.UTC))
	assert.Equal(t, *result.Buckets.Bucket[0].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.Buckets.Bucket[0].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.Buckets.Bucket[0].Name, "mybucket01")
	assert.Equal(t, *result.Buckets.Bucket[0].Region, "cn-hangzhou")
	assert.Equal(t, *result.Buckets.Bucket[0].Location, "oss-cn-hangzhou")
	assert.Equal(t, *result.Buckets.Bucket[0].ResourceGroupId, "rg-aek27tc********")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefaultV2)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

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
	resultErr := &ListVectorBucketsResult{}
	err = c.client.unmarshalOutput(resultErr, output, unmarshalBodyDefaultV2)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *DeleteVectorBucketRequest
	var input *OperationInput
	var err error

	request = &DeleteVectorBucketRequest{}
	input = &OperationInput{
		OpName: "DeleteVectorBucket",
		Method: "DELETE",
		Bucket: request.Bucket,
	}
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteVectorBucketRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteVectorBucket",
		Method: "DELETE",
		Bucket: request.Bucket,
	}
	err = c.client.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteVectorBucket(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5C3D9778CC1C2AEDF85B****"},
		},
	}
	result := &DeleteVectorBucketResult{}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C3D9778CC1C2AEDF85B****")
	
	output = &OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	err = c.client.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
}