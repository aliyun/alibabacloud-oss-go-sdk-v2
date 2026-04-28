package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_DoDataPipeLineAction(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DoDataPipeLineActionRequest
	var input *OperationInput
	var err error

	request = &DoDataPipeLineActionRequest{}
	input = &OperationInput{
		OpName: "DoDataPipeLineAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
		},
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Action.")

	request = &DoDataPipeLineActionRequest{
		Action: Ptr("listDataPipelineConfigurations"),
	}
	input = &OperationInput{
		OpName: "DoDataPipeLineAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
		},
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["action"], "listDataPipelineConfigurations")

	request = &DoDataPipeLineActionRequest{
		Action: Ptr("getDataPipelineConfiguration"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"dataPipelineName": "data-pipeline",
			},
		},
	}
	input = &OperationInput{
		OpName: "DoDataPipeLineAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
		},
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["action"], "getDataPipelineConfiguration")
	assert.Equal(t, input.Parameters["dataPipelineName"], "data-pipeline")

	request = &DoDataPipeLineActionRequest{
		Action: Ptr("deleteDataPipelineConfiguration"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"dataPipelineName": "data-pipeline",
			},
		},
	}
	input = &OperationInput{
		OpName: "DoDataPipeLineAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
		},
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["action"], "deleteDataPipelineConfiguration")
	assert.Equal(t, input.Parameters["dataPipelineName"], "data-pipeline")

	request = &DoDataPipeLineActionRequest{
		Action: Ptr("restartDataPipeline"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"dataPipelineName": "data-pipeline",
			},
		},
	}
	input = &OperationInput{
		OpName: "DoDataPipeLineAction",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
		},
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["action"], "restartDataPipeline")
	assert.Equal(t, input.Parameters["dataPipelineName"], "data-pipeline")
}

func TestUnmarshalOutput_DoDataPipeLineAction(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	var getBody []byte
	body := `<?xml version="1.0" encoding="UTF-8" ?>
<DataPipelineConfiguration>
  <DataPipelineName>my-data-pipeline</DataPipelineName>
  <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>
  <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>
  <Status>Running</Status>
  <Sources>
      <InputBucket>my-bucket</InputBucket>
      <InputDataScope>All</InputDataScope>
      <IgnoreDelete>true</IgnoreDelete>
      <FilterConfiguration>
          <PrefixSet>prefix1/</PrefixSet>
          <PrefixSet>prefix2/prefix3/</PrefixSet>
          <ObjectMediaTypes>text</ObjectMediaTypes>
          <ObjectMediaTypes>image</ObjectMediaTypes>
          <ObjectMediaTypes>video</ObjectMediaTypes>
      </FilterConfiguration>
  </Sources>
  <DataPipelineEmbeddingConfiguration>
      <EmbeddingProvider>bailian</EmbeddingProvider>
      <ApiKey>xxxx</ApiKey>
      <Model>qwen2.5-vl-embedding</Model>
      <FPS>1</FPS>
  </DataPipelineEmbeddingConfiguration>
  <Destination>
      <VectorBucketName>my-vector-bucket</VectorBucketName>
      <VectorIndexNames>my-index</VectorIndexNames>
      <VectorKeyPrefix></VectorKeyPrefix>
      <ObjectTagToMetadata>key1</ObjectTagToMetadata>
      <ObjectTagToMetadata>key2</ObjectTagToMetadata>
      <UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata>
  </Destination>
  <DataPipelineError>
      <ErrorMode>ignoreAndRecord</ErrorMode>
      <ErrorBucket>my-error-bucket</ErrorBucket>
      <ErrorPrefix>error-output/</ErrorPrefix>
  </DataPipelineError>
  <CreateTime>2021-06-29T14:50:13.011643661+08:00</CreateTime>
</DataPipelineConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &DoDataPipeLineActionResult{
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

	body = `<?xml version="1.0" encoding="UTF-8" ?>
<ListDataPipelineConfigurationsResult>
  <DataPipelineConfigurations>
    <DataPipelineConfiguration>
      <DataPipelineName>my-data-pipeline</DataPipelineName>
      <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>
      <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>
      <Status>Running</Status>
      <Sources>
          <InputBucket>my-bucket</InputBucket>
          <InputDataScope>All</InputDataScope>
          <IgnoreDelete>true</IgnoreDelete>
          <FilterConfiguration>
              <PrefixSet>prefix1/</PrefixSet>
              <PrefixSet>prefix2/prefix3/</PrefixSet>
              <ObjectMediaTypes>text</ObjectMediaTypes>
              <ObjectMediaTypes>image</ObjectMediaTypes>
              <ObjectMediaTypes>video</ObjectMediaTypes>
          </FilterConfiguration>
      </Sources>
      <DataPipelineEmbeddingConfiguration>
          <EmbeddingProvider>bailian</EmbeddingProvider>
          <ApiKey>xxxx</ApiKey>
          <Model>qwen2.5-vl-embedding</Model>
          <FPS>1</FPS>
      </DataPipelineEmbeddingConfiguration>
      <Destination>
          <VectorBucketName>my-vector-bucket</VectorBucketName>
          <VectorIndexNames>my-index</VectorIndexNames>
          <VectorKeyPrefix></VectorKeyPrefix>
          <ObjectTagToMetadata>key1</ObjectTagToMetadata>
          <ObjectTagToMetadata>key2</ObjectTagToMetadata>
          <UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata>
      </Destination>
      <DataPipelineError>
          <ErrorMode>ignoreAndRecord</ErrorMode>
          <ErrorBucket>my-error-bucket</ErrorBucket>
          <ErrorPrefix>error-output/</ErrorPrefix>
      </DataPipelineError>
      <CreateTime>2021-06-29T14:50:13.011643661+08:00</CreateTime>
    </DataPipelineConfiguration>
  </DataPipelineConfigurations>
  <NextToken>xxx</NextToken>
</ListDataPipelineConfigurationsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DoDataPipeLineActionResult{
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
	result = &DoDataPipeLineActionResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
