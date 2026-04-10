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

func TestMarshalInput_PutTableMaintenanceConfiguration(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *PutTableMaintenanceConfigurationRequest
	var input *oss.OperationInput
	var err error

	request = &PutTableMaintenanceConfigurationRequest{}
	input = &oss.OperationInput{
		OpName: "PutTableMaintenanceConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance/%s", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)), oss.ToString(request.Type))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn.")

	request = &PutTableMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "PutTableMaintenanceConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance/%s", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)), oss.ToString(request.Type))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &PutTableMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("oss-space"),
	}
	input = &oss.OperationInput{
		OpName: "PutTableMaintenanceConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance/%s", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)), oss.ToString(request.Type))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name.")

	request = &PutTableMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("oss-space"),
		Name:      oss.Ptr("oss-table"),
	}
	input = &oss.OperationInput{
		OpName: "PutTableMaintenanceConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance/%s", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)), oss.ToString(request.Type))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Type.")

	request = &PutTableMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("oss-space"),
		Name:      oss.Ptr("oss-table"),
		Type:      oss.Ptr("icebergCompaction"),
	}
	input = &oss.OperationInput{
		OpName: "PutTableMaintenanceConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance/%s", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)), oss.ToString(request.Type))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Value.")

	request = &PutTableMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("oss-space"),
		Name:      oss.Ptr("oss-table"),
		Type:      oss.Ptr("icebergSnapshotManagement"),
		Value: &TableMaintenanceValue{
			Status: oss.Ptr("enabled"),
			Settings: &TableMaintenanceSettings{
				IcebergSnapshotManagement: &IcebergSnapshotManagementSettingsDetail{
					MaxSnapshotAgeHours: oss.Ptr(int(350)),
					MinSnapshotsToKeep:  oss.Ptr(1),
				},
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "PutTableMaintenanceConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance/%s", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)), oss.ToString(request.Type))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/oss-space/oss-table/maintenance/icebergSnapshotManagement")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"value\":{\"settings\":{\"icebergSnapshotManagement\":{\"maxSnapshotAgeHours\":350,\"minSnapshotsToKeep\":1}},\"status\":\"enabled\"}}")

	request = &PutTableMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("oss-space"),
		Name:      oss.Ptr("oss-table"),
		Type:      oss.Ptr("icebergCompaction"),
		Value: &TableMaintenanceValue{
			Status: oss.Ptr("enabled"),
			Settings: &TableMaintenanceSettings{
				IcebergCompaction: &IcebergCompactionSettingsDetail{
					TargetFileSizeMB: oss.Ptr(400),
					Strategy:         oss.Ptr("auto"),
				},
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "PutTableMaintenanceConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance/%s", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)), oss.ToString(request.Type))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/oss-space/oss-table/maintenance/icebergCompaction")
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"value\":{\"settings\":{\"icebergCompaction\":{\"strategy\":\"auto\",\"targetFileSizeMB\":400}},\"status\":\"enabled\"}}")
}

func TestUnmarshalOutput_PutTableMaintenanceConfiguration(t *testing.T) {
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
	result := &PutTableMaintenanceConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutTableMaintenanceConfigurationResult{}
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
	result = &PutTableMaintenanceConfigurationResult{}
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
	result = &PutTableMaintenanceConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetTableMaintenanceConfiguration(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var request *GetTableMaintenanceConfigurationRequest
	var input *oss.OperationInput
	var err error

	request = &GetTableMaintenanceConfigurationRequest{}
	input = &oss.OperationInput{
		OpName: "GetTableMaintenanceConfiguration",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BucketArn.")

	request = &GetTableMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableMaintenanceConfiguration",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace.")

	request = &GetTableMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("oss-space"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableMaintenanceConfiguration",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name.")

	request = &GetTableMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		Namespace: oss.Ptr("oss-space"),
		Name:      oss.Ptr("oss-table"),
	}
	input = &oss.OperationInput{
		OpName: "GetTableMaintenanceConfiguration",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/maintenance", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
	assert.Equal(t, *input.Key, "tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/oss-space/oss-table/maintenance")
}

func TestUnmarshalOutput_GetTableMaintenanceConfiguration(t *testing.T) {
	c := TablesClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "tableARN": "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/table_id",
   "configuration": {
      "icebergCompaction": {
         "status": "enabled",
         "settings": {
            "icebergCompaction": {
               "targetFileSizeMB": 512,
               "strategy": "binpack"
            }
         }
      },
      "icebergSnapshotManagement": {
         "status": "enabled",
         "settings": {
            "icebergSnapshotManagement": {
               "minSnapshotsToKeep": 1,
               "maxSnapshotAgeHours": 720
            }
         }
      }
   }
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
	result := &GetTableMaintenanceConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)

	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, *result.Configuration.IcebergCompaction.Settings.IcebergCompaction.TargetFileSizeMB, 512)
	assert.Equal(t, *result.Configuration.IcebergCompaction.Settings.IcebergCompaction.Strategy, "binpack")
	assert.Equal(t, *result.Configuration.IcebergCompaction.Status, "enabled")
	assert.Equal(t, *result.Configuration.IcebergSnapshotManagement.Settings.IcebergSnapshotManagement.MaxSnapshotAgeHours, 720)
	assert.Equal(t, *result.Configuration.IcebergSnapshotManagement.Settings.IcebergSnapshotManagement.MinSnapshotsToKeep, 1)
	assert.Equal(t, *result.Configuration.IcebergSnapshotManagement.Status, "enabled")
	assert.Equal(t, *result.TableArn, "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/table_id")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetTableMaintenanceConfigurationResult{}
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
	result = &GetTableMaintenanceConfigurationResult{}
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
	result = &GetTableMaintenanceConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
