package dataprocess

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateDataset(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CreateDatasetRequest
	var input *oss.OperationInput
	var err error

	request = &CreateDatasetRequest{}
	input = &oss.OperationInput{
		OpName: "CreateDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "createDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &CreateDatasetRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "CreateDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "createDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &CreateDatasetRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
	}
	input = &oss.OperationInput{
		OpName: "CreateDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "createDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")

	request = &CreateDatasetRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		Description: oss.Ptr("this is a demo"),
		DatasetConfig: oss.Ptr(`{
  "Insights": {
    "Language": "zh"
  }
}`),
	}
	input = &oss.OperationInput{
		OpName: "CreateDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "createDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["description"], "this is a demo")
	assert.Equal(t, input.Parameters["datasetConfig"], "{\n  \"Insights\": {\n    \"Language\": \"zh\"\n  }\n}")

	request = &CreateDatasetRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		Description: oss.Ptr("this is a demo"),
		WorkflowParameters: oss.Ptr(WorkflowParameters{
			WorkflowParameter: []WorkflowParameter{
				{
					Name:  oss.Ptr("VideoInsightEnable"),
					Value: oss.Ptr("True"),
				},
				{
					Name:  oss.Ptr("ImageInsightEnable"),
					Value: oss.Ptr("True"),
				},
			},
		}.ToParameterValue()),
		DatasetConfig: oss.Ptr((&DatasetConfig{
			Insights: &InsightsConfig{
				Language: oss.Ptr("zh"),
				Image: &InsightsImage{
					Label: &InsightsImageLabel{
						System: &InsightsImageSystem{
							Enable: oss.Ptr(true),
						},
						UserDefined: &InsightsImageUserDefined{
							Enable: oss.Ptr(true),
							Mode:   oss.Ptr("Enhanced"),
							Labels: []LabelItem{
								{
									Name:        oss.Ptr("未佩戴安全帽"),
									Description: oss.Ptr("画面中存在人员，且人员头部未佩戴安全帽"),
								},
							},
						},
					},
				},
			},
		}).ToParameterValue()),
	}
	input = &oss.OperationInput{
		OpName: "CreateDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "createDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["description"], "this is a demo")
	assert.Equal(t, input.Parameters["datasetConfig"], "{\"Insights\":{\"Language\":\"zh\",\"Image\":{\"Label\":{\"System\":{\"Enable\":true},\"UserDefined\":{\"Enable\":true,\"Mode\":\"Enhanced\",\"Labels\":[{\"Name\":\"未佩戴安全帽\",\"Description\":\"画面中存在人员，且人员头部未佩戴安全帽\"}]}}}}}")
	assert.Equal(t, input.Parameters["workflowParameters"], "[{\"Name\":\"VideoInsightEnable\",\"Value\":\"True\"},{\"Name\":\"ImageInsightEnable\",\"Value\":\"True\"}]")
}

func TestUnmarshalOutput_CreateDataset(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<CreateDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters></WorkflowParameters>
<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>
<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig>
<Insights>
<Language>zh</Language>
<Image>
  <Label>
    <System>
      <Enable>true</Enable>
    </System>
    <UserDefined>
      <Enable>true</Enable>
      <Mode>Enhanced</Mode>
      <Labels>
        <Label>
          <Name>未佩戴安全帽</Name>
          <Description>画面中存在人员，且人员头部未佩戴安全帽</Description>
        </Label>
      </Labels>
    </UserDefined>
  </Label>
</Image>
</Insights></DatasetConfig>
</Dataset>
</CreateDatasetResponse>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &CreateDatasetResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Dataset.DatasetName, "test-dataset")
	assert.Equal(t, *result.Dataset.CreateTime, "2026-04-22T11:39:28.148283473+08:00")
	assert.Equal(t, *result.Dataset.UpdateTime, "2026-04-22T11:39:28.148283473+08:00")
	assert.Equal(t, *result.Dataset.Description, "this is a demo")
	assert.Equal(t, *result.Dataset.DatasetMaxFileCount, int64(100000000))
	assert.Equal(t, *result.Dataset.DatasetMaxEntityCount, int64(10000000000))
	assert.Equal(t, *result.Dataset.DatasetMaxRelationCount, int64(100000000000))
	assert.Equal(t, *result.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
	assert.Equal(t, *result.Dataset.DatasetConfig.Insights.Language, "zh")
	assert.Equal(t, *result.Dataset.DatasetConfig.Insights.Image.Label.System.Enable, true)
	assert.Equal(t, *result.Dataset.DatasetConfig.Insights.Image.Label.UserDefined.Enable, true)
	assert.Equal(t, *result.Dataset.DatasetConfig.Insights.Image.Label.UserDefined.Mode, "Enhanced")
	assert.Equal(t, *result.Dataset.DatasetConfig.Insights.Image.Label.UserDefined.Labels[0].Name, "未佩戴安全帽")
	assert.Equal(t, *result.Dataset.DatasetConfig.Insights.Image.Label.UserDefined.Labels[0].Description, "画面中存在人员，且人员头部未佩戴安全帽")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateDatasetResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetDataset(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetDatasetRequest
	var input *oss.OperationInput
	var err error

	request = &GetDatasetRequest{}
	input = &oss.OperationInput{
		OpName: "GetDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "getDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetDatasetRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "getDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &GetDatasetRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
	}
	input = &oss.OperationInput{
		OpName: "GetDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "getDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")

	request = &GetDatasetRequest{
		Bucket:         oss.Ptr("bucket"),
		DatasetName:    oss.Ptr("your_dataset"),
		WithStatistics: oss.Ptr(true),
	}
	input = &oss.OperationInput{
		OpName: "GetDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "getDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["withStatistics"], "true")
}

func TestUnmarshalOutput_GetDataset(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<GetDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters>
      <WorkflowParameter><Name>ImageInsightEnable</Name><Value>True</Value></WorkflowParameter>
</WorkflowParameters>
<CreateTime>2026-04-21T18:17:58.727923181+08:00</CreateTime>
<UpdateTime>2026-04-21T18:17:58.727923181+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh-Hans</Language></Insights></DatasetConfig>
<FileCount>3456</FileCount>
<TotalFileSize>10737418240</TotalFileSize>
</Dataset>
</GetDatasetResponse>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetDatasetResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Dataset.DatasetName, "test-dataset")
	assert.Equal(t, *result.Dataset.WorkflowParameters.WorkflowParameter[0].Name, "ImageInsightEnable")
	assert.Equal(t, *result.Dataset.WorkflowParameters.WorkflowParameter[0].Value, "True")
	assert.Equal(t, *result.Dataset.CreateTime, "2026-04-21T18:17:58.727923181+08:00")
	assert.Equal(t, *result.Dataset.UpdateTime, "2026-04-21T18:17:58.727923181+08:00")
	assert.Equal(t, *result.Dataset.Description, "this is a demo")
	assert.Equal(t, *result.Dataset.DatasetMaxFileCount, int64(100000000))
	assert.Equal(t, *result.Dataset.DatasetMaxEntityCount, int64(10000000000))
	assert.Equal(t, *result.Dataset.DatasetMaxRelationCount, int64(100000000000))
	assert.Equal(t, *result.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
	assert.Equal(t, *result.Dataset.DatasetConfig.Insights.Language, "zh-Hans")
	assert.Equal(t, *result.Dataset.FileCount, int64(3456))
	assert.Equal(t, *result.Dataset.TotalFileSize, int64(10737418240))

	body = `<GetDatasetResponse>
<Dataset>
<DatasetName>photos-2026</DatasetName>
<WorkflowParameters></WorkflowParameters>
<WorkflowParametersString></WorkflowParametersString>
<TemplateId>Official:OSSBasicMeta</TemplateId>
<CreateTime>2026-04-21T18:17:58.727923181+08:00</CreateTime>
<UpdateTime>2026-04-21T18:17:58.727923181+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxBindCount>10</DatasetMaxBindCount>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh-Hans</Language></Insights></DatasetConfig>
</Dataset>
</GetDatasetResponse>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetDatasetResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Dataset.DatasetName, "photos-2026")
	assert.Equal(t, *result.Dataset.CreateTime, "2026-04-21T18:17:58.727923181+08:00")
	assert.Equal(t, *result.Dataset.UpdateTime, "2026-04-21T18:17:58.727923181+08:00")
	assert.Equal(t, *result.Dataset.Description, "this is a demo")
	assert.Equal(t, *result.Dataset.DatasetMaxFileCount, int64(100000000))
	assert.Equal(t, *result.Dataset.DatasetMaxEntityCount, int64(10000000000))
	assert.Equal(t, *result.Dataset.DatasetMaxRelationCount, int64(100000000000))
	assert.Equal(t, *result.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
	assert.Equal(t, *result.Dataset.DatasetConfig.Insights.Language, "zh-Hans")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "Not Found",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetDatasetResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "Not Found")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_UpdateDataset(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UpdateDatasetRequest
	var input *oss.OperationInput
	var err error

	request = &UpdateDatasetRequest{}
	input = &oss.OperationInput{
		OpName: "UpdateDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "updateDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &UpdateDatasetRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "UpdateDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "updateDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &UpdateDatasetRequest{
		Bucket:             oss.Ptr("bucket"),
		DatasetName:        oss.Ptr("your_dataset"),
		Description:        oss.Ptr("this is a demo"),
		WorkflowParameters: oss.Ptr(`[{"Name": "ImageInsightEnable","Value": "True","Description": "The source bucket for data processing"}]`),
		DatasetConfig:      oss.Ptr(`{"Insights":{"Language":"zh","Image":{"Caption":{"Enable":"true"}},"Video":{"Caption":{"Enable":"true"},"Label":{"System":{"Enable":"true"},"UserDefined":{"Enable":"true"},"Highlight":{"Enable":"true"}},"MultiStream":{"Enable":"true"}}}}}`),
	}
	input = &oss.OperationInput{
		OpName: "UpdateDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "updateDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["description"], "this is a demo")
	assert.Equal(t, input.Parameters["workflowParameters"], "[{\"Name\": \"ImageInsightEnable\",\"Value\": \"True\",\"Description\": \"The source bucket for data processing\"}]")
	assert.Equal(t, input.Parameters["datasetConfig"], "{\"Insights\":{\"Language\":\"zh\",\"Image\":{\"Caption\":{\"Enable\":\"true\"}},\"Video\":{\"Caption\":{\"Enable\":\"true\"},\"Label\":{\"System\":{\"Enable\":\"true\"},\"UserDefined\":{\"Enable\":\"true\"},\"Highlight\":{\"Enable\":\"true\"}},\"MultiStream\":{\"Enable\":\"true\"}}}}}")

	request = &UpdateDatasetRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		Description: oss.Ptr("this is a demo"),
		WorkflowParameters: oss.Ptr(WorkflowParameters{
			WorkflowParameter: []WorkflowParameter{
				{
					Name:        oss.Ptr("ImageInsightEnable"),
					Value:       oss.Ptr("True"),
					Description: oss.Ptr("The source bucket for data processing"),
				},
			},
		}.ToParameterValue()),
		DatasetConfig: oss.Ptr((&DatasetConfig{
			Insights: &InsightsConfig{
				Language: oss.Ptr("zh"),
				Image: &InsightsImage{
					Caption: &InsightsImageCaption{
						Enable: oss.Ptr(true),
					},
				},
				Video: &InsightsVideo{
					Caption: &InsightsVideoCaption{
						Enable: oss.Ptr(true),
					},
					Label: &InsightsVideoLabel{
						System: &InsightsVideoSystem{
							Enable: oss.Ptr(true),
						},
						UserDefined: &InsightsVideoUserDefined{
							Enable: oss.Ptr(true),
							Labels: []LabelItem{
								{
									Name:        oss.Ptr("有人摔倒"),
									Description: oss.Ptr("画面中有人由站立或行走变为倒在地面的动作"),
								},
							},
						},
					},
					MultiStream: &InsightsVideoMultiStream{
						Enable: oss.Ptr(true),
					},
				},
			},
		}).ToParameterValue()),
	}
	input = &oss.OperationInput{
		OpName: "UpdateDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "updateDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["description"], "this is a demo")
	assert.Equal(t, input.Parameters["workflowParameters"], "[{\"Name\":\"ImageInsightEnable\",\"Value\":\"True\",\"Description\":\"The source bucket for data processing\"}]")
	assert.Equal(t, input.Parameters["datasetConfig"], "{\"Insights\":{\"Language\":\"zh\",\"Image\":{\"Caption\":{\"Enable\":true}},\"Video\":{\"Caption\":{\"Enable\":true},\"Label\":{\"System\":{\"Enable\":true},\"UserDefined\":{\"Enable\":true,\"Labels\":[{\"Name\":\"有人摔倒\",\"Description\":\"画面中有人由站立或行走变为倒在地面的动作\"}]}},\"MultiStream\":{\"Enable\":true}}}}")
}

func TestUnmarshalOutput_UpdateDataset(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<UpdateDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters></WorkflowParameters>
<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>
<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh</Language></Insights></DatasetConfig>
</Dataset>
</UpdateDatasetResponse>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &UpdateDatasetResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Dataset.DatasetName, "test-dataset")
	assert.Equal(t, *result.Dataset.CreateTime, "2026-04-22T11:39:28.148283473+08:00")
	assert.Equal(t, *result.Dataset.UpdateTime, "2026-04-22T11:39:28.148283473+08:00")
	assert.Equal(t, *result.Dataset.Description, "this is a demo")
	assert.Equal(t, *result.Dataset.DatasetMaxFileCount, int64(100000000))
	assert.Equal(t, *result.Dataset.DatasetMaxEntityCount, int64(10000000000))
	assert.Equal(t, *result.Dataset.DatasetMaxRelationCount, int64(100000000000))
	assert.Equal(t, *result.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
	assert.Equal(t, *result.Dataset.DatasetConfig.Insights.Language, "zh")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &UpdateDatasetResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListDatasets(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListDatasetsRequest
	var input *oss.OperationInput
	var err error

	request = &ListDatasetsRequest{}
	input = &oss.OperationInput{
		OpName: "ListDatasets",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "listDatasets",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &ListDatasetsRequest{
		Bucket:     oss.Ptr("bucket"),
		MaxResults: oss.Ptr(int64(10)),
		NextToken:  oss.Ptr("1986505809429276:oss_1234567890_demo-bucket:test-dataset"),
		Prefix:     oss.Ptr("prefix"),
	}
	input = &oss.OperationInput{
		OpName: "ListDatasets",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "listDatasets",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["maxResults"], "10")
	assert.Equal(t, input.Parameters["nextToken"], "1986505809429276:oss_1234567890_demo-bucket:test-dataset")
	assert.Equal(t, input.Parameters["prefix"], "prefix")
}

func TestUnmarshalOutput_ListDatasets(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<ListDatasetsResponse>
<NextToken>1986505809429276:oss_1234567890_demo-bucket:test-dataset</NextToken>
<Datasets>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>
<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>
</Dataset>
</Datasets>
</ListDatasetsResponse>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &ListDatasetsResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Datasets), 1)
	assert.Equal(t, *result.NextToken, "1986505809429276:oss_1234567890_demo-bucket:test-dataset")
	assert.Equal(t, *result.Datasets[0].DatasetName, "test-dataset")
	assert.Equal(t, *result.Datasets[0].CreateTime, "2026-04-22T11:39:28.148283473+08:00")
	assert.Equal(t, *result.Datasets[0].UpdateTime, "2026-04-22T11:39:28.148283473+08:00")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListDatasetsResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteDataset(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteDatasetRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteDatasetRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "deleteDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteDatasetRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "deleteDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &DeleteDatasetRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteDataset",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "deleteDataset",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
}

func TestUnmarshalOutput_DeleteDataset(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &DeleteDatasetResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "Not Found",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteDatasetResult{}
	err = c.client.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "Not Found")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
