package tables

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
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
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn.")

	request = &CreateNamespaceRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "CreateNamespace",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	request = &CreateNamespaceRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: []string{
			"test_namespace",
		},
	}
	input = &oss.OperationInput{
		OpName: "CreateNamespace",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, input.Parameters["namespaces"], "")
	jsonStr, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"namespace\":[\"test_namespace\"]}")
}

func TestUnmarshalOutput_CreateNamespace(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := `{
   "namespace": [ "test-namespace" ],
   "tableBucketARN": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"
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
	assert.Equal(t, result.Namespace, []string{"test-namespace"})
	assert.Equal(t, *result.TableBucketARN, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")

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
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s/%s", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn.")

	request = &GetNamespaceRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetNamespace",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s/%s", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	request = &GetNamespaceRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("test_namespace"),
	}
	input = &oss.OperationInput{
		OpName: "GetNamespace",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s/%s", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/test_namespace")
}

func TestUnmarshalOutput_GetNamespace(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "createdAt": "2026-04-03T09:00:44.014637+00:00",
   "createdBy": "1234567890",
   "namespace": ["my_space"],
   "namespaceId": "0a8fcd4d-a22a-42a4-a3f6-d4a88027018f",
   "ownerAccountId": "1234567890",
   "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5"
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
	assert.Equal(t, *result.CreatedBy, "1234567890")
	assert.Equal(t, *result.CreatedAt, "2026-04-03T09:00:44.014637+00:00")
	assert.Equal(t, *result.OwnerAccountId, "1234567890")
	assert.Equal(t, result.Namespace[0], "my_space")
	assert.Equal(t, *result.NamespaceId, "0a8fcd4d-a22a-42a4-a3f6-d4a88027018f")
	assert.Equal(t, *result.TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")
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
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn.")

	request = &ListNamespacesRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "ListNamespaces",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)

	request = &ListNamespacesRequest{
		BucketArn:     oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		MaxNamespaces: 10,
		Prefix:        oss.Ptr("/"),
	}
	input = &oss.OperationInput{
		OpName: "ListNamespaces",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["maxNamespaces"], "10")
	assert.Equal(t, input.Parameters["prefix"], "/")

	request = &ListNamespacesRequest{
		BucketArn:         oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s", url.QueryEscape(oss.ToString(request.BucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["maxNamespaces"], "10")
	assert.Equal(t, input.Parameters["prefix"], "/")
	assert.Equal(t, input.Parameters["continuationToken"], "123")
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket")
}

func TestUnmarshalOutput_ListNamespaces(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := `{
  "continuationToken": "CgxteV9uYW1lc3BhY2U-",
  "Namespaces": [{
    "createdAt": "2026-04-03T08:54:25.205905+00:00",
    "createdBy": "1760225545089999",
    "namespace": ["my_namespace"],
    "namespaceId": "22af7160-82b5-4d6a-b9fb-4d14c6e01199",
    "ownerAccountId": "1760225545089999",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f4"
  },
  {
     "createdAt": "2026-04-03T08:59:25.205905+00:00",
    "createdBy": "1760225545089999",
    "namespace": ["demo_namespace"],
    "namespaceId": "22af7160-82b5-4d6a-b9fb-4d14c6e01198",
    "ownerAccountId": "1760225545089999",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5"
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
	assert.Equal(t, *result.ContinuationToken, "CgxteV9uYW1lc3BhY2U-")
	assert.Equal(t, len(result.Namespaces), 2)
	assert.Equal(t, *result.Namespaces[0].CreatedAt, "2026-04-03T08:54:25.205905+00:00")
	assert.Equal(t, *result.Namespaces[0].CreatedBy, "1760225545089999")
	assert.Equal(t, result.Namespaces[0].Namespace[0], "my_namespace")
	assert.Equal(t, *result.Namespaces[0].NamespaceId, "22af7160-82b5-4d6a-b9fb-4d14c6e01199")
	assert.Equal(t, *result.Namespaces[0].OwnerAccountId, "1760225545089999")
	assert.Equal(t, *result.Namespaces[0].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f4")

	assert.Equal(t, *result.Namespaces[1].CreatedAt, "2026-04-03T08:59:25.205905+00:00")
	assert.Equal(t, *result.Namespaces[1].CreatedBy, "1760225545089999")
	assert.Equal(t, result.Namespaces[1].Namespace[0], "demo_namespace")
	assert.Equal(t, *result.Namespaces[1].NamespaceId, "22af7160-82b5-4d6a-b9fb-4d14c6e01198")
	assert.Equal(t, *result.Namespaces[1].OwnerAccountId, "1760225545089999")
	assert.Equal(t, *result.Namespaces[1].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")

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
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s/%s", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn.")

	request = &DeleteNamespaceRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteNamespace",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s/%s", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &DeleteNamespaceRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("oss_space"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteNamespace",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s/%s", url.QueryEscape(oss.ToString(request.BucketArn)), oss.ToString(request.Namespace))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/oss_space")
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
