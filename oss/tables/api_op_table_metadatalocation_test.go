package tables

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_GetTableMetadataLocation(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetTableMetadataLocationRequest
	var input *oss.OperationInput
	var err error

	request = &GetTableMetadataLocationRequest{}
	input = &oss.OperationInput{
		OpName: "GetTableMetadataLocation",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class":                 "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetTableMetadataLocationRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableMetadataLocation",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class":                 "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &GetTableMetadataLocationRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableMetadataLocation",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class":                 "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &GetTableMetadataLocationRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableMetadataLocation",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"metadata-location":             "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "GET")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["metadata-location"], "")
	assert.Equal(t, input.Parameters["space"], "")
	assert.Equal(t, input.Parameters["table"], "")
}

func TestUnmarshalOutput_GetTableMetadataLocation(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "metadataLocation": "location",
   "warehouseLocation": "bbb",
   "versionToken": "aaa"
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
	result := &GetTableMetadataLocationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.MetadataLocation, "location")
	assert.Equal(t, *result.VersionToken, "aaa")
	assert.Equal(t, *result.WarehouseLocation, "bbb")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetTableMetadataLocationResult{}
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
	result = &GetTableMetadataLocationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetTableMetadataLocationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_UpdateTableMetadataLocation(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *UpdateTableMetadataLocationRequest
	var input *oss.OperationInput
	var err error

	request = &UpdateTableMetadataLocationRequest{}
	input = &oss.OperationInput{
		OpName: "UpdateTableMetadataLocation",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class":                 "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &UpdateTableMetadataLocationRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "UpdateTableMetadataLocation",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class":                 "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &UpdateTableMetadataLocationRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "UpdateTableMetadataLocation",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class":                 "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &UpdateTableMetadataLocationRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "UpdateTableMetadataLocation",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class":                 "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, MetadataLocation.")

	request = &UpdateTableMetadataLocationRequest{
		Bucket:           oss.Ptr("oss-demo"),
		Namespace:        oss.Ptr("space"),
		Table:            oss.Ptr("table"),
		MetadataLocation: oss.Ptr("location"),
	}
	input = &oss.OperationInput{
		OpName: "UpdateTableMetadataLocation",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"storage-class":                 "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, VersionToken.")

	request = &UpdateTableMetadataLocationRequest{
		Bucket:           oss.Ptr("oss-demo"),
		Namespace:        oss.Ptr("space"),
		Table:            oss.Ptr("table"),
		MetadataLocation: oss.Ptr("location"),
		VersionToken:     oss.Ptr("version-token"),
	}
	input = &oss.OperationInput{
		OpName: "UpdateTableMetadataLocation",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"metadata-location":             "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["metadata-location"], "")
	assert.Equal(t, input.Parameters["space"], "")
	assert.Equal(t, input.Parameters["table"], "")
	body, err := io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(body), "{\"metadataLocation\":\"location\",\"versionToken\":\"version-token\"}")
}

func TestUnmarshalOutput_UpdateTableMetadataLocation(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "metadataLocation": "location",
   "name": "table",
   "namespace": [ "space" ],
   "tableARN": "acs:osstable:cn-hangzhou:123:bucket/table_bucket/table/table_123",
   "versionToken": "aaa"
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
	result := &UpdateTableMetadataLocationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.MetadataLocation, "location")
	assert.Equal(t, *result.Name, "table")
	assert.Equal(t, result.Namespace[0], "space")
	assert.Equal(t, *result.TableARN, "acs:osstable:cn-hangzhou:123:bucket/table_bucket/table/table_123")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &UpdateTableMetadataLocationResult{}
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
	result = &UpdateTableMetadataLocationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &UpdateTableMetadataLocationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
