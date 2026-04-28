package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_DoMetaQueryAction(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DoMetaQueryActionRequest
	var input *OperationInput
	var err error

	request = &DoMetaQueryActionRequest{}
	input = &OperationInput{
		OpName: "DoMetaQueryAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DoMetaQueryActionRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "DoMetaQueryAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery"})
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Action.")

	request = &DoMetaQueryActionRequest{
		Bucket: Ptr("bucket"),
		Action: Ptr("listDatasets"),
	}
	input = &OperationInput{
		OpName: "DoMetaQueryAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "DoMetaQueryAction",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["action"], "listDatasets")

	request = &DoMetaQueryActionRequest{
		Bucket: Ptr("bucket"),
		Action: Ptr("createDataset"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"datasetName":   "your_dataset",
				"description":   "this is a demo",
				"templateId":    "Official:OSSBasicMeta",
				"clusterType":   "auto",
				"datasetConfig": "{\n  \"Insights\": {\n    \"EnableLabel\": true,\n    \"EnableOCR\": true,\n    \"EnableFace\": true,\n    \"EnableImage\": true,\n    \"EnableVideo\": true,\n    \"EnableAudio\": true,\n    \"Language\": \"zh\"\n  }\n}",
			},
		},
	}
	input = &OperationInput{
		OpName: "DoMetaQueryAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["action"], "createDataset")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["description"], "this is a demo")
	assert.Equal(t, input.Parameters["templateId"], "Official:OSSBasicMeta")
	assert.Equal(t, input.Parameters["clusterType"], "auto")
	assert.Equal(t, input.Parameters["datasetConfig"], "{\n  \"Insights\": {\n    \"EnableLabel\": true,\n    \"EnableOCR\": true,\n    \"EnableFace\": true,\n    \"EnableImage\": true,\n    \"EnableVideo\": true,\n    \"EnableAudio\": true,\n    \"Language\": \"zh\"\n  }\n}")

	request = &DoMetaQueryActionRequest{
		Bucket: Ptr("bucket"),
		Action: Ptr("getDataset"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"datasetName": "your_dataset",
			},
		},
	}
	input = &OperationInput{
		OpName: "DoMetaQueryAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["action"], "getDataset")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")

	request = &DoMetaQueryActionRequest{
		Bucket: Ptr("bucket"),
		Action: Ptr("deleteDataset"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"datasetName": "your_dataset",
			},
		},
	}
	input = &OperationInput{
		OpName: "DoMetaQueryAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["action"], "deleteDataset")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
}

func TestUnmarshalOutput_DoMetaQueryAction(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	var getBody []byte
	body := `<CreateDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters></WorkflowParameters>
<WorkflowParametersString></WorkflowParametersString>
<TemplateId>Official:OSSBasicMeta</TemplateId>
<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>
<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxBindCount>10</DatasetMaxBindCount>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh</Language></Insights></DatasetConfig>
</Dataset>
</CreateDatasetResponse>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &DoMetaQueryActionResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	getBody, err = io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(getBody), body)

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DoMetaQueryActionResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
