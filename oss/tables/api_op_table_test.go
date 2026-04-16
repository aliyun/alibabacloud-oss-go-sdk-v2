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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN.")

	request = &CreateTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &CreateTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name.")

	request = &CreateTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
		Name:           oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Format.")

	request = &CreateTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
		Name:           oss.Ptr("table"),
		Format:         oss.Ptr("ICEBERG"),
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space")
	jsonStr, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"format\":\"ICEBERG\",\"name\":\"table\"}")

	request = &CreateTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
		Name:           oss.Ptr("table"),
		Format:         oss.Ptr("ICEBERG"),
		Metadata: &TableMetadata{
			Iceberg: &IcebergMetadata{
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
			KmsKeyArn:    oss.Ptr(""),
			SseAlgorithm: oss.Ptr("AES256"),
		},
	}
	input = &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["Content-Type"], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space")
	jsonStr, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(jsonStr), "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"\",\"sseAlgorithm\":\"AES256\"},\"format\":\"ICEBERG\",\"metadata\":{\"iceberg\":{\"schema\":{\"fields\":[{\"name\":\"id\",\"required\":true,\"type\":\"int\"},{\"name\":\"name\",\"type\":\"string\"}]}}},\"name\":\"table\"}")
}

func TestUnmarshalOutput_CreateTable(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := `{
   "tableARN": "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/16dc6c23-7a64-4f55-af2f-ee243524a5cc",
   "versionToken": "8c651fb37897499092bd95e1bc2816a9"
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
	assert.Equal(t, oss.ToString(result.TableARN), "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/16dc6c23-7a64-4f55-af2f-ee243524a5cc")
	assert.Equal(t, oss.ToString(result.VersionToken), "8c651fb37897499092bd95e1bc2816a9")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
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

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	resultErr = &CreateTableResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 400)
	assert.Equal(t, resultErr.Status, "Bad Request")
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
		Key: oss.Ptr("get-table"),
	}
	err = checkGetTableRequest(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must provide either table arn alone OR all of (table bucket arn, namespace, table name) together")

	request = &GetTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("get-table"),
	}
	err = checkGetTableRequest(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must provide either table arn alone OR all of (table bucket arn, namespace, table name) together")

	request = &GetTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Name:           oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("get-table"),
	}
	err = checkGetTableRequest(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must provide either table arn alone OR all of (table bucket arn, namespace, table name) together")

	request = &GetTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Name:           oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("get-table"),
	}
	err = checkGetTableRequest(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must provide either table arn alone OR all of (table bucket arn, namespace, table name) together")

	request = &GetTableRequest{
		TableARN:       oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9"),
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("get-table"),
	}
	err = checkGetTableRequest(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must provide either table arn alone OR all of (table bucket arn, namespace, table name) together")

	request = &GetTableRequest{
		TableARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket1/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("get-table"),
	}
	err = checkGetTableRequest(request)
	assert.Nil(t, err)
	input.Bucket, err = parseBucketArn(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed table arn")

	request = &GetTableRequest{
		TableARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table-1313/f13de3a6-de93-4801-bd7f-a09c124177d9"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("get-table"),
	}
	err = checkGetTableRequest(request)
	assert.Nil(t, err)
	input.Bucket, err = parseBucketArn(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed table arn")

	request = &GetTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
		Name:           oss.Ptr("table"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("get-table"),
	}
	err = checkGetTableRequest(request)
	assert.Nil(t, err)
	input.Bucket, err = parseBucketArn(request)
	assert.Nil(t, err)
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "get-table")

	request = &GetTableRequest{
		TableARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9"),
	}
	input = &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Key: oss.Ptr("get-table"),
	}
	err = checkGetTableRequest(request)
	assert.Nil(t, err)
	input.Bucket, err = parseBucketArn(request)
	assert.Nil(t, err)
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "get-table")
}

func TestUnmarshalOutput_GetTable(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
    "createdAt": "2026-04-07T05:27:18.397920+00:00",
    "createdBy": "1234567890",
    "format": "ICEBERG",
    "metadataLocation": "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss/metadata/00000-edb683a9-ce46-492a-a495-35e5b2f7a649.metadata.json",
    "modifiedAt": "2026-04-07T05:27:18.397920+00:00",
    "modifiedBy": "1234567890",
    "name": "my_table",
    "namespace": ["my_namespace"],
    "namespaceId": "22af7160-82b5-4d6a-b9fb-4d14c6e01198",
    "ownerAccountId": "1234567890",
    "tableARN": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5",
    "type": "customer",
    "versionToken": "365f934c6e234f35ace5ae48f0a0d871",
    "warehouseLocation": "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss"}`
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
	assert.Equal(t, *result.CreatedAt, "2026-04-07T05:27:18.397920+00:00")
	assert.Equal(t, *result.CreatedBy, "1234567890")
	assert.Equal(t, *result.Format, "ICEBERG")
	assert.Equal(t, *result.MetadataLocation, "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss/metadata/00000-edb683a9-ce46-492a-a495-35e5b2f7a649.metadata.json")
	assert.Equal(t, *result.ModifiedAt, "2026-04-07T05:27:18.397920+00:00")
	assert.Equal(t, *result.ModifiedBy, "1234567890")
	assert.Equal(t, *result.Name, "my_table")
	assert.Equal(t, result.Namespace[0], "my_namespace")
	assert.Equal(t, *result.NamespaceId, "22af7160-82b5-4d6a-b9fb-4d14c6e01198")
	assert.Equal(t, *result.OwnerAccountId, "1234567890")
	assert.Equal(t, *result.TableARN, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9")
	assert.Equal(t, *result.TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")
	assert.Equal(t, *result.Type, "customer")
	assert.Equal(t, *result.VersionToken, "365f934c6e234f35ace5ae48f0a0d871")
	assert.Equal(t, *result.WarehouseLocation, "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss")

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN.")

	request = &ListTablesRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "ListTables",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &ListTablesRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "ListTables",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket")

	request = &ListTablesRequest{
		TableBucketARN:    oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket")
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
    "continuationToken": "Ci9teV90YWJsZV8xQGViOTk4ZjEwLWQyMGMtNGYyMi05YTc2LWVkNjRlOTY2OGY1Ng--",
    "tables": [{
            "createdAt": "2026-04-07T02:15:12.186626+00:00",
            "modifiedAt": "2026-04-07T02:15:12.186626+00:00",
            "name": "my_table",
            "namespace": ["my_namespace"],
            "tableARN": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/0e5f5125-ec94-4a82-a630-d6feade09217",
            "type": "customer"},
        {
            "createdAt": "2026-04-07T02:39:11.988947+00:00",
            "modifiedAt": "2026-04-07T02:39:11.988947+00:00",
            "name": "my_table_1",
            "namespace": ["my_namespace"],
            "tableARN": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/eb998f10-d20c-4f22-9a76-ed64e9668f56",
            "type": "customer"}]}`
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

	assert.Equal(t, *result.ContinuationToken, "Ci9teV90YWJsZV8xQGViOTk4ZjEwLWQyMGMtNGYyMi05YTc2LWVkNjRlOTY2OGY1Ng--")
	assert.Equal(t, len(result.Tables), 2)
	assert.Equal(t, *result.Tables[0].CreatedAt, "2026-04-07T02:15:12.186626+00:00")
	assert.Equal(t, *result.Tables[0].ModifiedAt, "2026-04-07T02:15:12.186626+00:00")
	assert.Equal(t, *result.Tables[0].Name, "my_table")
	assert.Equal(t, result.Tables[0].Namespace[0], "my_namespace")
	assert.Equal(t, *result.Tables[0].TableARN, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/0e5f5125-ec94-4a82-a630-d6feade09217")
	assert.Equal(t, *result.Tables[0].Type, "customer")

	assert.Equal(t, *result.Tables[1].CreatedAt, "2026-04-07T02:39:11.988947+00:00")
	assert.Equal(t, *result.Tables[1].ModifiedAt, "2026-04-07T02:39:11.988947+00:00")
	assert.Equal(t, *result.Tables[1].Name, "my_table_1")
	assert.Equal(t, result.Tables[1].Namespace[0], "my_namespace")
	assert.Equal(t, *result.Tables[1].TableARN, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/eb998f10-d20c-4f22-9a76-ed64e9668f56")
	assert.Equal(t, *result.Tables[1].Type, "customer")

	body = `{"message": "AccessDenied"}`
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
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN.")

	request = &DeleteTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTable",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &DeleteTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTable",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name.")

	request = &DeleteTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
		Name:           oss.Ptr("table"),
		VersionToken:   oss.Ptr("version_token"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteTable",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table")
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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/rename", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	err = checkRenameTableRequest(request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "either NewTable or NewNamespace must be provided")

	request = &RenameTableRequest{
		NewName: oss.Ptr("new_table"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/rename", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	err = checkRenameTableRequest(request)
	assert.Nil(t, err)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN.")

	request = &RenameTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		NewName:        oss.Ptr("new_table"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/rename", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	err = checkRenameTableRequest(request)
	assert.Nil(t, err)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &RenameTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
		NewName:        oss.Ptr("new_table"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/rename", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	err = checkRenameTableRequest(request)
	assert.Nil(t, err)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name.")

	request = &RenameTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
		Name:           oss.Ptr("table"),
		NewNamespace:   oss.Ptr("new-space"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/rename", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	err = checkRenameTableRequest(request)
	assert.Nil(t, err)
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/rename")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"newNamespaceName\":\"new-space\"}")

	request = &RenameTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace:      oss.Ptr("space"),
		Name:           oss.Ptr("table"),
		NewNamespace:   oss.Ptr("new-space"),
		NewName:        oss.Ptr("new-table"),
		VersionToken:   oss.Ptr("365f934c6e234f35ace5ae48f0a0d871"),
	}
	input = &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/rename", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	err = checkRenameTableRequest(request)
	assert.Nil(t, err)
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/rename")
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"newName\":\"new-table\",\"newNamespaceName\":\"new-space\",\"versionToken\":\"365f934c6e234f35ace5ae48f0a0d871\"}")
}

func TestUnmarshalOutput_RenameTable(t *testing.T) {
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
	result := &RenameTableResult{}
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
	resultErr := &RenameTableResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/json")
}
