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

func TestMarshalInput_PutDataPipelineConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutDataPipelineConfigurationRequest
	var input *oss.OperationInput
	var err error

	request = &PutDataPipelineConfigurationRequest{}
	input = &oss.OperationInput{
		OpName: "PutDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "putDataPipelineConfiguration",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DataPipelineName.")

	request = &PutDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr("data-pipeline"),
	}
	input = &oss.OperationInput{
		OpName: "PutDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "putDataPipelineConfiguration",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Role.")

	request = &PutDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr("data-pipeline"),
		Role:             oss.Ptr("AliyunOSSDataPipelineRole"),
	}
	input = &oss.OperationInput{
		OpName: "PutDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "putDataPipelineConfiguration",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DataPipelineConfiguration.")

	request = &PutDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr("data-pipeline"),
		Role:             oss.Ptr("AliyunOSSDataPipelineRole"),
		DataPipelineConfiguration: &DataPipelineConfiguration{
			DataPipelineDescription: oss.Ptr("使用百炼多模态模型为业务数据向量化"),
			Sources: []DataPipelineSource{
				{
					InputBucket:    oss.Ptr("bucket"),
					InputDataScope: oss.Ptr("All"),
					FilterConfiguration: &DataPipelineSourceFilterConfiguration{
						PrefixSet:        []string{"prefix1", "prefix2/prefix3"},
						ObjectMediaTypes: []string{"text", "image", "video"},
					},
				},
			},
			DataPipelineEmbeddingConfiguration: &DataPipelineEmbeddingConfiguration{
				ApiKey:            oss.Ptr("your_api_key"),
				EmbeddingProvider: oss.Ptr("bailian"),
				FPS:               oss.Ptr(float64(1)),
				Model:             oss.Ptr("qwen2.5-vl-embedding"),
			},
			Destination: &DataPipelineDestination{
				VectorBucketName:    oss.Ptr("my-vector-bucket"),
				VectorIndexNames:    []string{"my-index"},
				VectorKeyPrefix:     oss.Ptr("prefix"),
				ObjectTagToMetadata: []string{"key1", "key2"},
				UsermetaToMetadata:  []string{"x-oss-meta-key1"},
			},
			DataPipelineError: &DataPipelineError{
				ErrorMode:   oss.Ptr("ignoreAndRecord"),
				ErrorBucket: oss.Ptr("my-error-bucket"),
				ErrorPrefix: oss.Ptr("error-output/"),
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "PutDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "putDataPipelineConfiguration",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["action"], "putDataPipelineConfiguration")
	assert.Equal(t, input.Parameters["dataPipelineName"], "data-pipeline")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<DataPipelineConfiguration><DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription><Sources><InputBucket>bucket</InputBucket><InputDataScope>All</InputDataScope><FilterConfiguration><PrefixSet>prefix1</PrefixSet><PrefixSet>prefix2/prefix3</PrefixSet><ObjectMediaTypes>text</ObjectMediaTypes><ObjectMediaTypes>image</ObjectMediaTypes><ObjectMediaTypes>video</ObjectMediaTypes></FilterConfiguration></Sources><DataPipelineEmbeddingConfiguration><EmbeddingProvider>bailian</EmbeddingProvider><ApiKey>your_api_key</ApiKey><Model>qwen2.5-vl-embedding</Model><FPS>1</FPS></DataPipelineEmbeddingConfiguration><Destination><VectorBucketName>my-vector-bucket</VectorBucketName><VectorKeyPrefix>prefix</VectorKeyPrefix><VectorIndexNames>my-index</VectorIndexNames><ObjectTagToMetadata>key1</ObjectTagToMetadata><ObjectTagToMetadata>key2</ObjectTagToMetadata><UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata></Destination><DataPipelineError><ErrorMode>ignoreAndRecord</ErrorMode><ErrorBucket>my-error-bucket</ErrorBucket><ErrorPrefix>error-output/</ErrorPrefix></DataPipelineError></DataPipelineConfiguration>")

	request = &PutDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr("data-pipeline-v2"),
		Role:             oss.Ptr("AliyunOSSDataPipelineRole"),
		DataPipelineConfiguration: &DataPipelineConfiguration{
			Sources: []DataPipelineSource{{
				InputBucket:    oss.Ptr("source-bucket"),
				InputDataScope: oss.Ptr("All"),
				IgnoreDelete:   oss.Ptr(false),
			}},
			ModelTier: oss.Ptr("standard"),
			DataPipelineDataProcessConfiguration: &DataPipelineDataProcessConfiguration{
				SearchMode: oss.Ptr("balanced"),
				Insights: &DataPipelineInsights{
					Image: &DataPipelineInsightsImage{Caption: &DataPipelineInsightsCaption{Prompt: oss.Ptr("Describe the image.")}},
					Video: &DataPipelineInsightsVideo{
						Caption: &DataPipelineInsightsCaption{Prompt: oss.Ptr("Describe each video scene.")},
						FrameEmbedding: &DataPipelineInsightsFrameEmbedding{Snapshot: &DataPipelineInsightsSnapshot{
							Mode:     oss.Ptr("interval"),
							Interval: oss.Ptr(1.0),
						}},
					},
				},
			},
			Destination: &DataPipelineDestination{
				ImageEmbedding:         &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("image"), Prefix: oss.Ptr("v2")},
				ImageTextEmbedding:     &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("image-text"), Prefix: oss.Ptr("v2")},
				VideoFrameEmbedding:    &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("video-frame"), Prefix: oss.Ptr("v2")},
				VideoTextEmbedding:     &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("video-text"), Prefix: oss.Ptr("v2")},
				DocumentChunkEmbedding: &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("document"), Prefix: oss.Ptr("v2")},
				ObjectTagToMetadata:    []string{"category"},
				UsermetaToMetadata:     []string{"x-oss-meta-source"},
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "PutDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "putDataPipelineConfiguration",
		},
	}

	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NoError(t, err)
	assert.Equal(t, "data-pipeline-v2", input.Parameters["dataPipelineName"])
	body, err = io.ReadAll(input.Body)
	assert.NoError(t, err)
	assert.Equal(t, `<DataPipelineConfiguration><Sources><InputBucket>source-bucket</InputBucket><InputDataScope>All</InputDataScope><IgnoreDelete>false</IgnoreDelete></Sources><Destination><ObjectTagToMetadata>category</ObjectTagToMetadata><UsermetaToMetadata>x-oss-meta-source</UsermetaToMetadata><ImageEmbedding><Bucket>vector-bucket</Bucket><IndexName>image</IndexName><Prefix>v2</Prefix></ImageEmbedding><ImageTextEmbedding><Bucket>vector-bucket</Bucket><IndexName>image-text</IndexName><Prefix>v2</Prefix></ImageTextEmbedding><VideoFrameEmbedding><Bucket>vector-bucket</Bucket><IndexName>video-frame</IndexName><Prefix>v2</Prefix></VideoFrameEmbedding><VideoTextEmbedding><Bucket>vector-bucket</Bucket><IndexName>video-text</IndexName><Prefix>v2</Prefix></VideoTextEmbedding><DocumentChunkEmbedding><Bucket>vector-bucket</Bucket><IndexName>document</IndexName><Prefix>v2</Prefix></DocumentChunkEmbedding></Destination><ModelTier>standard</ModelTier><DataPipelineDataProcessConfiguration><SearchMode>balanced</SearchMode><Insights><Image><Caption><Prompt>Describe the image.</Prompt></Caption></Image><Video><Caption><Prompt>Describe each video scene.</Prompt></Caption><FrameEmbedding><Snapshot><Mode>interval</Mode><Interval>1</Interval></Snapshot></FrameEmbedding></Video></Insights></DataPipelineDataProcessConfiguration></DataPipelineConfiguration>`, string(body))
}

func TestUnmarshalOutput_PutDataPipelineConfiguration(t *testing.T) {
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
	result := &PutDataPipelineConfigurationResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutDataPipelineConfigurationResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetDataPipelineConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetDataPipelineConfigurationRequest
	var input *oss.OperationInput
	var err error

	request = &GetDataPipelineConfigurationRequest{}
	input = &oss.OperationInput{
		OpName: "GetDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "getDataPipelineConfiguration",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DataPipelineName.")

	request = &GetDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr("data-pipeline"),
	}
	input = &oss.OperationInput{
		OpName: "GetDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "getDataPipelineConfiguration",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["action"], "getDataPipelineConfiguration")
	assert.Equal(t, input.Parameters["dataPipelineName"], "data-pipeline")
}

func TestUnmarshalOutput_GetDataPipelineConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8" ?>
<DataPipelineConfiguration>
    <DataPipelineName>my-data-pipeline</DataPipelineName>
    <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>
    <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>
    <Status>Running</Status>
    <Phase>IncrementalScanning</Phase>
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
        <ApiKey>sk-12345678901234556</ApiKey>
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
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &GetDataPipelineConfigurationResult{}
	err = c.client.UnmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineDescription, "使用百炼多模态模型为业务数据向量化")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineName, "my-data-pipeline")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineRole, "my-data-pipeline-role")
	assert.Equal(t, *result.DataPipelineConfiguration.Status, "Running")
	assert.Equal(t, *result.DataPipelineConfiguration.Phase, "IncrementalScanning")
	assert.Equal(t, *result.DataPipelineConfiguration.Sources[0].InputBucket, "my-bucket")
	assert.Equal(t, *result.DataPipelineConfiguration.Sources[0].InputDataScope, "All")
	assert.Equal(t, *result.DataPipelineConfiguration.Sources[0].IgnoreDelete, true)
	assert.Equal(t, result.DataPipelineConfiguration.Sources[0].FilterConfiguration.PrefixSet[0], "prefix1/")
	assert.Equal(t, result.DataPipelineConfiguration.Sources[0].FilterConfiguration.PrefixSet[1], "prefix2/prefix3/")
	assert.Equal(t, result.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[0], "text")
	assert.Equal(t, result.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[1], "image")
	assert.Equal(t, result.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[2], "video")
	assert.Equal(t, result.DataPipelineConfiguration.Sources[0].FilterConfiguration.PrefixSet[1], "prefix2/prefix3/")
	assert.Equal(t, result.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[0], "text")
	assert.Equal(t, result.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[1], "image")
	assert.Equal(t, result.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[2], "video")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineEmbeddingConfiguration.EmbeddingProvider, "bailian")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineEmbeddingConfiguration.ApiKey, "sk-12345678901234556")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineEmbeddingConfiguration.Model, "qwen2.5-vl-embedding")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineEmbeddingConfiguration.FPS, float64(1))
	assert.Equal(t, *result.DataPipelineConfiguration.Destination.VectorBucketName, "my-vector-bucket")
	assert.Equal(t, result.DataPipelineConfiguration.Destination.VectorIndexNames[0], "my-index")
	assert.Equal(t, *result.DataPipelineConfiguration.Destination.VectorKeyPrefix, "")
	assert.Equal(t, result.DataPipelineConfiguration.Destination.ObjectTagToMetadata[0], "key1")
	assert.Equal(t, result.DataPipelineConfiguration.Destination.ObjectTagToMetadata[1], "key2")
	assert.Equal(t, result.DataPipelineConfiguration.Destination.UsermetaToMetadata[0], "x-oss-meta-key1")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineError.ErrorBucket, "my-error-bucket")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineError.ErrorPrefix, "error-output/")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineError.ErrorMode, "ignoreAndRecord")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineError.ErrorPrefix, "error-output/")
	assert.Equal(t, *result.DataPipelineConfiguration.DataPipelineError.ErrorMode, "ignoreAndRecord")
	assert.Equal(t, *result.DataPipelineConfiguration.CreateTime, "2021-06-29T14:50:13.011643661+08:00")

	body = `<DataPipelineConfiguration>
	<ModelTier>standard</ModelTier>
	<DataPipelineDataProcessConfiguration>
		<SearchMode>balanced</SearchMode>
		<Insights>
			<Image><Caption><Prompt>Describe the image.</Prompt></Caption></Image>
			<Video><Caption><Prompt>Describe each video scene.</Prompt></Caption><FrameEmbedding><Snapshot><Mode>dhash</Mode><Number>10</Number></Snapshot></FrameEmbedding></Video>
		</Insights>
	</DataPipelineDataProcessConfiguration>
	<Destination>
		<ImageEmbedding><Bucket>vector-bucket</Bucket><IndexName>image</IndexName><Prefix>v2</Prefix></ImageEmbedding>
		<ImageTextEmbedding><Bucket>vector-bucket</Bucket><IndexName>image-text</IndexName><Prefix>v2</Prefix></ImageTextEmbedding>
		<VideoFrameEmbedding><Bucket>vector-bucket</Bucket><IndexName>video-frame</IndexName><Prefix>v2</Prefix></VideoFrameEmbedding>
		<VideoTextEmbedding><Bucket>vector-bucket</Bucket><IndexName>video-text</IndexName><Prefix>v2</Prefix></VideoTextEmbedding>
		<DocumentChunkEmbedding><Bucket>vector-bucket</Bucket><IndexName>document</IndexName><Prefix>v2</Prefix></DocumentChunkEmbedding>
		<ObjectTagToMetadata>category</ObjectTagToMetadata>
		<UsermetaToMetadata>x-oss-meta-source</UsermetaToMetadata>
	</Destination>
</DataPipelineConfiguration>`
	output = &oss.OperationOutput{Body: io.NopCloser(bytes.NewReader([]byte(body)))}
	result = &GetDataPipelineConfigurationResult{}

	err = (&Client{}).client.UnmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.NoError(t, err)
	configuration := result.DataPipelineConfiguration
	assert.Equal(t, "standard", oss.ToString(configuration.ModelTier))
	assert.Equal(t, "balanced", oss.ToString(configuration.DataPipelineDataProcessConfiguration.SearchMode))
	assert.Equal(t, "Describe the image.", oss.ToString(configuration.DataPipelineDataProcessConfiguration.Insights.Image.Caption.Prompt))
	assert.Equal(t, "Describe each video scene.", oss.ToString(configuration.DataPipelineDataProcessConfiguration.Insights.Video.Caption.Prompt))
	assert.Equal(t, "dhash", oss.ToString(configuration.DataPipelineDataProcessConfiguration.Insights.Video.FrameEmbedding.Snapshot.Mode))
	assert.Equal(t, int64(10), oss.ToInt64(configuration.DataPipelineDataProcessConfiguration.Insights.Video.FrameEmbedding.Snapshot.Number))
	assert.Equal(t, "image", oss.ToString(configuration.Destination.ImageEmbedding.IndexName))
	assert.Equal(t, "image-text", oss.ToString(configuration.Destination.ImageTextEmbedding.IndexName))
	assert.Equal(t, "video-frame", oss.ToString(configuration.Destination.VideoFrameEmbedding.IndexName))
	assert.Equal(t, "video-text", oss.ToString(configuration.Destination.VideoTextEmbedding.IndexName))
	assert.Equal(t, "document", oss.ToString(configuration.Destination.DocumentChunkEmbedding.IndexName))
	assert.Equal(t, []string{"category"}, configuration.Destination.ObjectTagToMetadata)
	assert.Equal(t, []string{"x-oss-meta-source"}, configuration.Destination.UsermetaToMetadata)

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetDataPipelineConfigurationResult{}
	err = c.client.UnmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteDataPipelineConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteDataPipelineConfigurationRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteDataPipelineConfigurationRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "deleteDataPipelineConfiguration",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DataPipelineName.")

	request = &DeleteDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr("data-pipeline"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "deleteDataPipelineConfiguration",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["action"], "deleteDataPipelineConfiguration")
	assert.Equal(t, input.Parameters["dataPipelineName"], "data-pipeline")
}

func TestUnmarshalOutput_DeleteDataPipelineConfiguration(t *testing.T) {
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
	result := &DeleteDataPipelineConfigurationResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteDataPipelineConfigurationResult{}
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

func TestMarshalInput_ListDataPipelineConfigurations(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListDataPipelineConfigurationsRequest
	var input *oss.OperationInput
	var err error

	request = &ListDataPipelineConfigurationsRequest{}
	input = &oss.OperationInput{
		OpName: "ListDataPipelineConfigurations",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "listDataPipelineConfigurations",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)

	request = &ListDataPipelineConfigurationsRequest{
		MaxResults:  oss.Ptr(int64(100)),
		NextToken:   oss.Ptr("next-token"),
		Prefix:      oss.Ptr("prefix"),
		InputBucket: oss.Ptr("source-bucket"),
	}
	input = &oss.OperationInput{
		OpName: "ListDataPipelineConfigurations",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "listDataPipelineConfigurations",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["action"], "listDataPipelineConfigurations")
	assert.Equal(t, input.Parameters["maxResults"], "100")
	assert.Equal(t, input.Parameters["nextToken"], "next-token")
	assert.Equal(t, input.Parameters["prefix"], "prefix")
	assert.Equal(t, input.Parameters["inputBucket"], "source-bucket")
}

func TestUnmarshalOutput_ListDataPipelineConfigurations(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8" ?>
<ListDataPipelineConfigurationsResult>
  <DataPipelineConfigurations>
    <DataPipelineConfiguration>
      <DataPipelineName>my-data-pipeline</DataPipelineName>
      <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>
      <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>
      <Status>Running</Status>
      <Phase>IncrementalScanning</Phase>
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
          <ApiKey>sk-12345678901234556</ApiKey>
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
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &ListDataPipelineConfigurationsResult{}
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
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, len(result.DataPipelineConfigurations), 1)
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineDescription, "使用百炼多模态模型为业务数据向量化")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineName, "my-data-pipeline")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineRole, "my-data-pipeline-role")
	assert.Equal(t, *result.DataPipelineConfigurations[0].Status, "Running")
	assert.Equal(t, *result.DataPipelineConfigurations[0].Phase, "IncrementalScanning")
	assert.Equal(t, *result.DataPipelineConfigurations[0].Sources[0].InputBucket, "my-bucket")
	assert.Equal(t, *result.DataPipelineConfigurations[0].Sources[0].InputDataScope, "All")
	assert.Equal(t, *result.DataPipelineConfigurations[0].Sources[0].IgnoreDelete, true)
	assert.Equal(t, result.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.PrefixSet[0], "prefix1/")
	assert.Equal(t, result.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.PrefixSet[1], "prefix2/prefix3/")
	assert.Equal(t, result.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[0], "text")
	assert.Equal(t, result.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[1], "image")
	assert.Equal(t, result.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[2], "video")
	assert.Equal(t, result.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.PrefixSet[1], "prefix2/prefix3/")
	assert.Equal(t, result.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[0], "text")
	assert.Equal(t, result.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[1], "image")
	assert.Equal(t, result.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[2], "video")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineEmbeddingConfiguration.EmbeddingProvider, "bailian")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineEmbeddingConfiguration.ApiKey, "sk-12345678901234556")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineEmbeddingConfiguration.Model, "qwen2.5-vl-embedding")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineEmbeddingConfiguration.FPS, float64(1))
	assert.Equal(t, *result.DataPipelineConfigurations[0].Destination.VectorBucketName, "my-vector-bucket")
	assert.Equal(t, result.DataPipelineConfigurations[0].Destination.VectorIndexNames[0], "my-index")
	assert.Equal(t, *result.DataPipelineConfigurations[0].Destination.VectorKeyPrefix, "")
	assert.Equal(t, result.DataPipelineConfigurations[0].Destination.ObjectTagToMetadata[0], "key1")
	assert.Equal(t, result.DataPipelineConfigurations[0].Destination.ObjectTagToMetadata[1], "key2")
	assert.Equal(t, result.DataPipelineConfigurations[0].Destination.UsermetaToMetadata[0], "x-oss-meta-key1")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineError.ErrorBucket, "my-error-bucket")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineError.ErrorPrefix, "error-output/")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineError.ErrorMode, "ignoreAndRecord")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineError.ErrorPrefix, "error-output/")
	assert.Equal(t, *result.DataPipelineConfigurations[0].DataPipelineError.ErrorMode, "ignoreAndRecord")
	assert.Equal(t, *result.DataPipelineConfigurations[0].CreateTime, "2021-06-29T14:50:13.011643661+08:00")
	assert.Equal(t, *result.NextToken, "xxx")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListDataPipelineConfigurationsResult{}
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

func TestMarshalInput_PauseDataPipeline(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PauseDataPipelineRequest
	var input *oss.OperationInput
	var err error

	request = &PauseDataPipelineRequest{}
	input = &oss.OperationInput{
		OpName: "PauseDataPipeline",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "pauseDataPipeline",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PauseDataPipelineRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "PauseDataPipeline",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "pauseDataPipeline",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DataPipelineName.")

	request = &PauseDataPipelineRequest{
		Bucket:           oss.Ptr("bucket"),
		DataPipelineName: oss.Ptr("data-pipeline"),
	}
	input = &oss.OperationInput{
		OpName: "PauseDataPipeline",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "pauseDataPipeline",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["action"], "pauseDataPipeline")
	assert.Equal(t, input.Parameters["dataPipelineName"], "data-pipeline")
}

func TestUnmarshalOutput_PauseDataPipeline(t *testing.T) {
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
	result := &PauseDataPipelineResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PauseDataPipelineResult{}
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

func TestMarshalInput_RestartDataPipeline(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *RestartDataPipelineRequest
	var input *oss.OperationInput
	var err error

	request = &RestartDataPipelineRequest{}
	input = &oss.OperationInput{
		OpName: "RestartDataPipeline",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "restartDataPipeline",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DataPipelineName.")

	request = &RestartDataPipelineRequest{
		DataPipelineName: oss.Ptr("data-pipeline"),
	}
	input = &oss.OperationInput{
		OpName: "RestartDataPipeline",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"dataPipeline": "",
			"action":       "restartDataPipeline",
		},
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["action"], "restartDataPipeline")
	assert.Equal(t, input.Parameters["dataPipelineName"], "data-pipeline")
}

func TestUnmarshalOutput_RestartDataPipeline(t *testing.T) {
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
	result := &RestartDataPipelineResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &RestartDataPipelineResult{}
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
