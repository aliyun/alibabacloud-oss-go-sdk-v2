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

func TestMarshalInput_CreateTable(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *CreateTableRequest
	var input *oss.OperationInput
	var err error

	request = &CreateTableRequest{}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &CreateTableRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	request = &CreateTableRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Format")

	request = &CreateTableRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Format:    oss.Ptr("iceberg"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table")

	request = &CreateTableRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Format:    oss.Ptr("iceberg"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Metadata")

	request = &CreateTableRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Format:    oss.Ptr("iceberg"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Metadata")

	request = &CreateTableRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Format:    oss.Ptr("iceberg"),
		Table:     oss.Ptr("table"),
		Metadata: &TableMetadata{
			Iceberg: &MetadataIceberg{
				Schema: map[string]any{
					"fields": []map[string]any{
						{
							"name": "id", "type": "int", "required": true,
						},
						{
							"name": "name", "type": "string",
						},
					},
				},
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	jsonStr, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"format\":\"iceberg\",\"metadata\":{\"iceberg\":{\"schema\":{\"fields\":[{\"name\":\"id\",\"required\":true,\"type\":\"int\"},{\"name\":\"name\",\"type\":\"string\"}]}}},\"name\":\"table\"}")

	request = &CreateTableRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Namespace: oss.Ptr("space"),
		Format:    oss.Ptr("iceberg"),
		Table:     oss.Ptr("table"),
		Metadata: &TableMetadata{
			Iceberg: &MetadataIceberg{
				Schema: map[string]any{
					"fields": []map[string]any{
						{
							"name": "id", "type": "int", "required": true,
						},
						{
							"name": "name", "type": "string",
						},
					},
				},
			},
		},
		EncryptionConfiguration: &EncryptionConfiguration{
			KmsKeyArn:    oss.Ptr("arn"),
			SseAlgorithm: oss.Ptr("AES256"),
		},
		Tags: map[string]any{
			"k1": "v1", "k2": "v2",
		},
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	jsonStr, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"arn\",\"sseAlgorithm\":\"AES256\"},\"format\":\"iceberg\",\"metadata\":{\"iceberg\":{\"schema\":{\"fields\":[{\"name\":\"id\",\"required\":true,\"type\":\"int\"},{\"name\":\"name\",\"type\":\"string\"}]}}},\"name\":\"table\",\"tags\":{\"k1\":\"v1\",\"k2\":\"v2\"}}")
}

func TestUnmarshalOutput_CreateTable(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := `{
   "tableARN": "acs:osstable:cn-hangzhou:123:bucket/table_name/table/123456",
   "versionToken": "aaabbb"
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
	result := &CreateTableResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, oss.ToString(result.TableARN), "acs:osstable:cn-hangzhou:123:bucket/table_name/table/123456")
	assert.Equal(t, oss.ToString(result.VersionToken), "aaabbb")

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
	resultErr := &CreateTableResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetTable(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetTableRequest
	var input *oss.OperationInput
	var err error

	request = &GetTableRequest{}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"get-table": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetTableRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"get-table": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &GetTableRequest{
		Bucket: oss.Ptr("bucket"),
		Table:  oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"get-table": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &GetTableRequest{
		Bucket:    oss.Ptr("bucket"),
		Table:     oss.Ptr("table"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"get-table": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableArn.")

	request = &GetTableRequest{
		Bucket:    oss.Ptr("bucket"),
		Table:     oss.Ptr("table"),
		Namespace: oss.Ptr("space"),
		TableArn:  oss.Ptr("table-arn"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"get-table": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN.")

	request = &GetTableRequest{
		Bucket:         oss.Ptr("bucket"),
		Table:          oss.Ptr("table"),
		Namespace:      oss.Ptr("space"),
		TableArn:       oss.Ptr("table-arn"),
		TableBucketARN: oss.Ptr("table-bucket-arn"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"get-table": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["name"], "table")
	assert.Equal(t, input.Parameters["namespace"], "space")
	assert.Equal(t, input.Parameters["tableArn"], "table-arn")
	assert.Equal(t, input.Parameters["tableBucketARN"], "table-bucket-arn")
}

func TestUnmarshalOutput_GetTable(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "createdAt": "2026-02-31T10:56:21.000Z",
   "createdBy": "oss-create",
   "format": "demo-format",
   "metadataLocation": "location",
   "modifiedAt": "2026-03-01T10:56:21.000Z",
   "modifiedBy": "oss-modify",
   "name": "table",
   "namespace": [ "space" ],
   "namespaceId": "space-01",
   "ownerAccountId": "123",
   "tableARN": "acs:osstable:cn-hangzhou:123:bucket/table_bucket/table/table_123",
   "tableBucketId": "table_bucket_123",
   "type": "oss",
   "versionToken": "aaa",
   "warehouseLocation": "bbb"
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
	result := &GetTableResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.CreatedAt, "2026-02-31T10:56:21.000Z")
	assert.Equal(t, *result.CreatedBy, "oss-create")
	assert.Equal(t, *result.Format, "demo-format")
	assert.Equal(t, *result.MetadataLocation, "location")
	assert.Equal(t, *result.ModifiedAt, "2026-03-01T10:56:21.000Z")
	assert.Equal(t, *result.ModifiedBy, "oss-modify")
	assert.Equal(t, *result.Name, "table")
	assert.Equal(t, result.Namespace[0], "space")
	assert.Equal(t, *result.NamespaceId, "space-01")
	assert.Equal(t, *result.OwnerAccountId, "123")
	assert.Equal(t, *result.TableARN, "acs:osstable:cn-hangzhou:123:bucket/table_bucket/table/table_123")
	assert.Equal(t, *result.TableBucketId, "table_bucket_123")
	assert.Equal(t, *result.Type, "oss")
	assert.Equal(t, *result.VersionToken, "aaa")
	assert.Equal(t, *result.WarehouseLocation, "bbb")

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
	resultErr := &GetTableResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_ListTables(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *ListTablesRequest
	var input *oss.OperationInput
	var err error

	request = &ListTablesRequest{}
	input = &oss.OperationInput{
		OpName: "ListTables",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &ListTablesRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "ListTables",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &ListTablesRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "ListTables",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"get-table"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)

	request = &ListTablesRequest{
		Bucket:            oss.Ptr("bucket"),
		Namespace:         oss.Ptr("space"),
		ContinuationToken: oss.Ptr("token"),
		MaxTables:         int32(1000),
		Prefix:            oss.Ptr("prefix"),
	}
	input = &oss.OperationInput{
		OpName: "ListTables",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["namespace"], "space")
	assert.Equal(t, input.Parameters["continuationToken"], "token")
	assert.Equal(t, input.Parameters["maxTables"], "1000")
	assert.Equal(t, input.Parameters["prefix"], "prefix")
}

func TestUnmarshalOutput_ListTables(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
  "continuationToken": "AAMA-EFRSURBSGk2VFNsQXNjVHQ2QU05UU5YN2xkME53VWI3U1B5RTl6WEh1UTRVc",
  "tables": [
    {
      "createdAt": "2026-01-26T03:11:51.527997035Z",
      "modifiedAt": "2026-01-26T03:11:51.527997035Z",
      "name": "example_table",
      "namespace": [
        "my_namespace"
      ],
      "tableARN": "arn:aws:s3tables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/7568a090-50f8-4808-8c8d-930a2c264076",
      "type": "customer"
    },
    {
      "createdAt": "2026-01-26T03:16:46.622650810Z",
      "modifiedAt": "2026-01-26T03:16:46.622650810Z",
      "name": "example_table1",
      "namespace": [
        "my_namespace"
      ],
      "tableARN": "arn:aws:s3tables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/757c17c1-532e-4a45-b5b3-d8783374fc2a",
      "type": "customer"
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
	result := &ListTablesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	assert.Equal(t, *result.ContinuationToken, "AAMA-EFRSURBSGk2VFNsQXNjVHQ2QU05UU5YN2xkME53VWI3U1B5RTl6WEh1UTRVc")
	assert.Equal(t, len(result.Tables), 2)
	assert.Equal(t, *result.Tables[0].CreatedAt, "2026-01-26T03:11:51.527997035Z")
	assert.Equal(t, *result.Tables[0].ModifiedAt, "2026-01-26T03:11:51.527997035Z")
	assert.Equal(t, *result.Tables[0].Name, "example_table")
	assert.Equal(t, result.Tables[0].Namespace[0], "my_namespace")
	assert.Equal(t, *result.Tables[0].TableARN, "arn:aws:s3tables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/7568a090-50f8-4808-8c8d-930a2c264076")
	assert.Equal(t, *result.Tables[0].Type, "customer")

	assert.Equal(t, *result.Tables[1].CreatedAt, "2026-01-26T03:16:46.622650810Z")
	assert.Equal(t, *result.Tables[1].ModifiedAt, "2026-01-26T03:16:46.622650810Z")
	assert.Equal(t, *result.Tables[1].Name, "example_table1")
	assert.Equal(t, result.Tables[1].Namespace[0], "my_namespace")
	assert.Equal(t, *result.Tables[1].TableARN, "arn:aws:s3tables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/757c17c1-532e-4a45-b5b3-d8783374fc2a")
	assert.Equal(t, *result.Tables[1].Type, "customer")

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
	resultErr := &ListTablesResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteTable(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *DeleteTableRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteTableRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteTableRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &DeleteTableRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &DeleteTableRequest{
		Bucket:       oss.Ptr("bucket"),
		Namespace:    oss.Ptr("space"),
		Table:        oss.Ptr("table"),
		VersionToken: oss.Ptr("version_token"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["versionToken"], "version_token")
}

func TestUnmarshalOutput_DeleteTable(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	output = &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteTableResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	resultErr := &DeleteTableResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_RenameTable(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *RenameTableRequest
	var input *oss.OperationInput
	var err error

	request = &RenameTableRequest{}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"rename":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &RenameTableRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"rename":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &RenameTableRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"rename":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Table.")

	request = &RenameTableRequest{
		Bucket:    oss.Ptr("bucket"),
		Namespace: oss.Ptr("space"),
		Table:     oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"rename":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, NewNamespace.")

	request = &RenameTableRequest{
		Bucket:       oss.Ptr("bucket"),
		Namespace:    oss.Ptr("space"),
		Table:        oss.Ptr("table"),
		NewNamespace: oss.Ptr("new-space"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"rename":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, NewTable.")

	request = &RenameTableRequest{
		Bucket:       oss.Ptr("bucket"),
		Namespace:    oss.Ptr("space"),
		Table:        oss.Ptr("table"),
		NewNamespace: oss.Ptr("new-space"),
		NewTable:     oss.Ptr("new-table"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"rename":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, VersionToken.")

	request = &RenameTableRequest{
		Bucket:       oss.Ptr("bucket"),
		Namespace:    oss.Ptr("space"),
		Table:        oss.Ptr("table"),
		NewNamespace: oss.Ptr("new-space"),
		NewTable:     oss.Ptr("new-table"),
		VersionToken: oss.Ptr("version-token"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"rename":                        "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["tables"], "")
	assert.Equal(t, input.Parameters["space"], "")
	assert.Equal(t, input.Parameters["table"], "")
	assert.Equal(t, input.Parameters["rename"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"namespace\":\"new-space\",\"newName\":\"new-table\",\"versionToken\":\"version-token\"}")
}

func TestUnmarshalOutput_RenameTable(t *testing.T) {
	c := TablesClient{}
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
	result := &RenameTableResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	resultErr := &RenameTableResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}
