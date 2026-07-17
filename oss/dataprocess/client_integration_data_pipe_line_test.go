//go:build integration

package dataprocess

import (
	"context"
	"errors"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataPipeLine(t *testing.T) {
	var err error
	endpoint := "https://oss-" + region_ + ".aliyuncs.com"
	client := getClient(region_, endpoint)
	var serr *oss.ServiceError

	_, err = client.PutDataPipelineConfiguration(context.TODO(), &PutDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr("data-pipeline"),
		Role:             oss.Ptr("not-exist-role"),
		DataPipelineConfiguration: &DataPipelineConfiguration{
			DataPipelineDescription: oss.Ptr("使用百炼多模态模型为业务数据向量化"),
			Sources: []DataPipelineSource{
				{
					InputBucket:    oss.Ptr("bucket"),
					InputDataScope: oss.Ptr("All"),
					FilterConfiguration: &DataPipelineSourceFilterConfiguration{
						PrefixSet:        []string{"prefix1"},
						ObjectMediaTypes: []string{"text"},
					},
				},
			},
			DataPipelineEmbeddingConfiguration: &DataPipelineEmbeddingConfiguration{
				ApiKey:            oss.Ptr("sk-123323423423423423424242425436457657567"),
				EmbeddingProvider: oss.Ptr("bailian"),
				FPS:               oss.Ptr(float64(1)),
				Model:             oss.Ptr("qwen2.5-vl-embedding"),
			},
			Destination: &DataPipelineDestination{
				VectorBucketName:    oss.Ptr("my-vector-bucket"),
				VectorIndexNames:    []string{"index"},
				VectorKeyPrefix:     oss.Ptr("prefix"),
				ObjectTagToMetadata: []string{"key1"},
				UsermetaToMetadata:  []string{"x-oss-meta-key1"},
			},
			DataPipelineError: &DataPipelineError{
				ErrorMode:   oss.Ptr("ignoreAndRecord"),
				ErrorBucket: oss.Ptr("my-error-bucket"),
				ErrorPrefix: oss.Ptr("error-output/"),
			},
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "InvalidArgument", serr.Code)
	assert.Equal(t, "The role parameter is not a valid RAM role ARN (expect acs:ram::<account>:role/<name>).", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.ListDataPipelineConfigurations(context.TODO(), &ListDataPipelineConfigurationsRequest{})
	assert.Nil(t, err)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.GetDataPipelineConfiguration(context.TODO(), &GetDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr("not-exist-data-pipeline"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchDataPipeline", serr.Code)
	assert.Equal(t, "The specified resource DataPipeline is not found.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.PauseDataPipeline(context.TODO(), &PauseDataPipelineRequest{
		DataPipelineName: oss.Ptr("not-exist-data-pipeline"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchDataPipeline", serr.Code)
	assert.Equal(t, "The specified resource DataPipeline is not found.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.RestartDataPipeline(context.TODO(), &RestartDataPipelineRequest{
		DataPipelineName: oss.Ptr("not-exist-data-pipeline"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchDataPipeline", serr.Code)
	assert.Equal(t, "The specified resource DataPipeline is not found.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DeleteDataPipelineConfiguration(context.TODO(), &DeleteDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr("not-exist-data-pipeline"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchDataPipeline", serr.Code)
	assert.Equal(t, "The specified resource DataPipeline is not found.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}
