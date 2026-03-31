package tables

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateNamespace(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *CreateNamespaceRequest
	var input *oss.OperationInput
	var err error

	request = &CreateNamespaceRequest{}
	input = &oss.OperationInput{
		OpName: "CreateNamespace",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &CreateNamespaceRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "CreateNamespace",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	request = &CreateNamespaceRequest{
		Bucket: oss.Ptr("bucket"),
		Namespace: []string{
			"test-namespace",
		},
	}
	input = &oss.OperationInput{
		OpName: "CreateNamespace",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, input.Parameters["namespaces"], "")
	jsonStr, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"namespace\":[\"test-namespace\"]}")
}

func TestUnmarshalOutput_CreateNamespace(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := `{
   "namespace": [ "test-namespace" ],
   "tableBucketARN": "bucket-arn"
}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &CreateNamespaceResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, len(result.Namespace), 1)
	assert.Equal(t, result.Namespace[0], "test-namespace")
	assert.Equal(t, *result.TableBucketARN, "bucket-arn")

	output = &oss.OperationOutput{
		StatusCode: 409,
		Status:     "BucketAlreadyExist",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &CreateNamespaceResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "BucketAlreadyExist")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetNamespace(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetNamespaceRequest
	var input *oss.OperationInput
	var err error

	request = &GetNamespaceRequest{}
	input = &oss.OperationInput{
		OpName: "GetNamespace",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetNamespaceRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetNamespace",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	request = &GetNamespaceRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("test-namespace"),
	}
	input = &oss.OperationInput{
		OpName: "GetNamespace",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	input.Parameters[*request.Namespace] = ""
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, input.Parameters["namespaces"], "")
}

func TestUnmarshalOutput_GetNamespace(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "createdAt": "2013-07-31T10:56:21.000Z",
   "createdBy": "aliyun",
   "namespace": ["123"],
   "namespaceId": "123",
   "ownerAccountId": "123456",
   "tableBucketId": "1"
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
	result := &GetNamespaceResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.CreatedBy, "aliyun")
	assert.Equal(t, *result.CreatedAt, "2013-07-31T10:56:21.000Z")
	assert.Equal(t, *result.OwnerAccountId, "123456")
	assert.Equal(t, result.Namespace[0], "123")
	assert.Equal(t, *result.NamespaceId, "123")
	assert.Equal(t, *result.TableBucketId, "1")
}

func TestMarshalInput_ListNamespaces(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *ListNamespacesRequest
	var input *oss.OperationInput
	var err error

	request = &ListNamespacesRequest{}
	input = &oss.OperationInput{
		OpName: "ListNamespaces",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &ListNamespacesRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "ListNamespaces",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)

	request = &ListNamespacesRequest{
		Bucket:        oss.Ptr("bucket"),
		MaxNamespaces: 10,
		Prefix:        oss.Ptr("/"),
	}
	input = &oss.OperationInput{
		OpName: "ListNamespaces",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["maxNamespaces"], "10")
	assert.Equal(t, input.Parameters["prefix"], "/")

	request = &ListNamespacesRequest{
		Bucket:            oss.Ptr("bucket"),
		MaxNamespaces:     10,
		Prefix:            oss.Ptr("/"),
		ContinuationToken: oss.Ptr("123"),
	}
	input = &oss.OperationInput{
		OpName: "ListNamespaces",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["maxNamespaces"], "10")
	assert.Equal(t, input.Parameters["prefix"], "/")
	assert.Equal(t, input.Parameters["continuationToken"], "123")
}

func TestUnmarshalOutput_ListNamespaces(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := `{
  "continuationToken": "token-123",
  "Namespaces": [{
    "createdAt": "2026-01-31T10:56:21.000Z",
    "createdBy": "aliyun",
    "namespace": ["demo-space"],
    "namespaceId": "123",
    "ownerAccountId": "123456",
    "tableBucketId": "1"
  },
  {
     "createdAt": "2026-02-31T10:56:21.000Z",
    "createdBy": "aliyun",
    "namespace": ["oss-space"],
    "namespaceId": "123457",
    "ownerAccountId": "123456",
    "tableBucketId": "2"
  }]
}`

	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"5374A2880232A65C2300****"},
		},
	}
	result := &ListNamespacesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5374A2880232A65C2300****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.ContinuationToken, "token-123")
	assert.Equal(t, len(result.Namespaces), 2)
	assert.Equal(t, *result.Namespaces[0].CreatedAt, "2026-01-31T10:56:21.000Z")
	assert.Equal(t, *result.Namespaces[0].CreatedBy, "aliyun")
	assert.Equal(t, result.Namespaces[0].Namespace[0], "demo-space")
	assert.Equal(t, *result.Namespaces[0].NamespaceId, "123")
	assert.Equal(t, *result.Namespaces[0].OwnerAccountId, "123456")
	assert.Equal(t, *result.Namespaces[0].TableBucketId, "1")

	assert.Equal(t, *result.Namespaces[1].CreatedAt, "2026-02-31T10:56:21.000Z")
	assert.Equal(t, *result.Namespaces[1].CreatedBy, "aliyun")
	assert.Equal(t, result.Namespaces[1].Namespace[0], "oss-space")
	assert.Equal(t, *result.Namespaces[1].NamespaceId, "123457")
	assert.Equal(t, *result.Namespaces[1].OwnerAccountId, "123456")
	assert.Equal(t, *result.Namespaces[1].TableBucketId, "2")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"Content-Type":     {"application/json"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson2)
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
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	resultErr := &ListNamespacesResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteNamespace(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *DeleteNamespaceRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteNamespaceRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteNamespace",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteNamespaceRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteNamespace",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &DeleteNamespaceRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("oss-space"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteNamespace",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"namespaces": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"namespaces"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
}

func TestUnmarshalOutput_DeleteNamespace(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	output = &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5C3D9778CC1C2AEDF85B****"},
		},
	}
	result := &DeleteNamespaceResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C3D9778CC1C2AEDF85B****")

	output = &oss.OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
}
