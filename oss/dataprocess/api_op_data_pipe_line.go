package dataprocess

import (
	"context"
	"encoding/xml"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PutDataPipelineConfigurationRequest struct {
	DataPipelineName          *string                    `input:"query,dataPipelineName,required"`
	Role                      *string                    `input:"query,role,required"`
	DataPipelineConfiguration *DataPipelineConfiguration `input:"body,DataPipelineConfiguration,xml,required"`
	oss.RequestCommon
}

type DataPipelineConfiguration struct {
	DataPipelineDescription            *string                             `xml:"DataPipelineDescription"`
	Sources                            []DataPipelineSource                `xml:"Sources"`
	DataPipelineEmbeddingConfiguration *DataPipelineEmbeddingConfiguration `xml:"DataPipelineEmbeddingConfiguration"`
	Destination                        *DataPipelineDestination            `xml:"Destination"`
	DataPipelineError                  *DataPipelineError                  `xml:"DataPipelineError"`
	DataPipelineName                   *string                             `xml:"DataPipelineName,omitempty"`
	DataPipelineRole                   *string                             `xml:"DataPipelineRole,omitempty"`
	Status                             *string                             `xml:"Status,omitempty"`
	Phase                              *string                             `xml:"Phase,omitempty"`
	CreateTime                         *string                             `xml:"CreateTime,omitempty"`

	ModelTier                            *string                               `xml:"ModelTier"`
	DataPipelineDataProcessConfiguration *DataPipelineDataProcessConfiguration `xml:"DataPipelineDataProcessConfiguration"`
}

type DataPipelineSource struct {
	InputBucket         *string                                `xml:"InputBucket"`
	InputDataScope      *string                                `xml:"InputDataScope"`
	IgnoreDelete        *bool                                  `xml:"IgnoreDelete"`
	FilterConfiguration *DataPipelineSourceFilterConfiguration `xml:"FilterConfiguration"`
}

type DataPipelineEmbeddingConfiguration struct {
	EmbeddingProvider *string  `xml:"EmbeddingProvider"`
	ApiKey            *string  `xml:"ApiKey"`
	Model             *string  `xml:"Model"`
	FPS               *float64 `xml:"FPS"`
}

type DataPipelineDataProcessConfiguration struct {
	SearchMode *string               `xml:"SearchMode"`
	Insights   *DataPipelineInsights `xml:"Insights"`
}

type DataPipelineInsights struct {
	Image *DataPipelineInsightsImage `xml:"Image"`
	Video *DataPipelineInsightsVideo `xml:"Video"`
}

type DataPipelineInsightsImage struct {
	Caption *DataPipelineInsightsCaption `xml:"Caption"`
}

type DataPipelineInsightsVideo struct {
	Caption        *DataPipelineInsightsCaption        `xml:"Caption"`
	FrameEmbedding *DataPipelineInsightsFrameEmbedding `xml:"FrameEmbedding"`
}

type DataPipelineInsightsCaption struct {
	Prompt *string `xml:"Prompt"`
}

type DataPipelineInsightsFrameEmbedding struct {
	Snapshot *DataPipelineInsightsSnapshot `xml:"Snapshot"`
}

type DataPipelineInsightsSnapshot struct {
	Mode     *string  `xml:"Mode"`
	Interval *float64 `xml:"Interval"`
	Number   *int64   `xml:"Number"`
}

type DataPipelineDestination struct {
	VectorBucketName    *string  `xml:"VectorBucketName"`
	VectorKeyPrefix     *string  `xml:"VectorKeyPrefix"`
	VectorIndexNames    []string `xml:"VectorIndexNames"`
	ObjectTagToMetadata []string `xml:"ObjectTagToMetadata"`
	UsermetaToMetadata  []string `xml:"UsermetaToMetadata"`

	ImageEmbedding         *VectorDestination `xml:"ImageEmbedding"`
	ImageTextEmbedding     *VectorDestination `xml:"ImageTextEmbedding"`
	VideoFrameEmbedding    *VectorDestination `xml:"VideoFrameEmbedding"`
	VideoTextEmbedding     *VectorDestination `xml:"VideoTextEmbedding"`
	DocumentChunkEmbedding *VectorDestination `xml:"DocumentChunkEmbedding"`
}

type VectorDestination struct {
	Bucket    *string `xml:"Bucket"`
	IndexName *string `xml:"IndexName"`
	Prefix    *string `xml:"Prefix"`
}

type DataPipelineError struct {
	ErrorMode   *string `xml:"ErrorMode"`
	ErrorBucket *string `xml:"ErrorBucket"`
	ErrorPrefix *string `xml:"ErrorPrefix"`
}

type DataPipelineSourceFilterConfiguration struct {
	PrefixSet        []string `xml:"PrefixSet"`
	ObjectMediaTypes []string `xml:"ObjectMediaTypes"`
}

type PutDataPipelineConfigurationResult struct {
	oss.ResultCommon
}

func (c *Client) PutDataPipelineConfiguration(ctx context.Context, request *PutDataPipelineConfigurationRequest, optFns ...func(*oss.Options)) (*PutDataPipelineConfigurationResult, error) {
	var err error
	if request == nil {
		request = &PutDataPipelineConfigurationRequest{}
	}

	input := &oss.OperationInput{
		OpName: "PutDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"action":       "putDataPipelineConfiguration",
			"dataPipeline": "",
		},
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutDataPipelineConfigurationResult{}

	if err = c.client.UnmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

type GetDataPipelineConfigurationRequest struct {
	DataPipelineName *string `input:"query,dataPipelineName,required"`
	oss.RequestCommon
}

type GetDataPipelineConfigurationResult struct {
	DataPipelineConfiguration *DataPipelineConfiguration `output:"body,DataPipelineConfiguration,xml"`

	oss.ResultCommon
}

func (c *Client) GetDataPipelineConfiguration(ctx context.Context, request *GetDataPipelineConfigurationRequest, optFns ...func(*oss.Options)) (*GetDataPipelineConfigurationResult, error) {
	var err error
	if request == nil {
		request = &GetDataPipelineConfigurationRequest{}
	}

	input := &oss.OperationInput{
		OpName: "GetDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"action":       "getDataPipelineConfiguration",
			"dataPipeline": "",
		},
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetDataPipelineConfigurationResult{}

	if err = c.client.UnmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

type DeleteDataPipelineConfigurationRequest struct {
	DataPipelineName *string `input:"query,dataPipelineName,required"`
	oss.RequestCommon
}

type DeleteDataPipelineConfigurationResult struct {
	oss.ResultCommon
}

func (c *Client) DeleteDataPipelineConfiguration(ctx context.Context, request *DeleteDataPipelineConfigurationRequest, optFns ...func(*oss.Options)) (*DeleteDataPipelineConfigurationResult, error) {
	var err error
	if request == nil {
		request = &DeleteDataPipelineConfigurationRequest{}
	}

	input := &oss.OperationInput{
		OpName: "DeleteDataPipelineConfiguration",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"action":       "deleteDataPipelineConfiguration",
			"dataPipeline": "",
		},
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteDataPipelineConfigurationResult{}

	if err = c.client.UnmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

type ListDataPipelineConfigurationsRequest struct {
	MaxResults *int64  `input:"query,maxResults"`
	Prefix     *string `input:"query,prefix"`
	NextToken  *string `input:"query,nextToken"`

	InputBucket *string `input:"query,inputBucket"`
	oss.RequestCommon
}

type ListDataPipelineConfigurationsResult struct {
	XMLName                    xml.Name                    `xml:"ListDataPipelineConfigurationsResult"`
	DataPipelineConfigurations []DataPipelineConfiguration `xml:"DataPipelineConfigurations>DataPipelineConfiguration"`
	NextToken                  *string                     `xml:"NextToken"`
	oss.ResultCommon
}

func (c *Client) ListDataPipelineConfigurations(ctx context.Context, request *ListDataPipelineConfigurationsRequest, optFns ...func(*oss.Options)) (*ListDataPipelineConfigurationsResult, error) {
	var err error
	if request == nil {
		request = &ListDataPipelineConfigurationsRequest{}
	}

	input := &oss.OperationInput{
		OpName: "ListDataPipelineConfigurations",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"action":       "listDataPipelineConfigurations",
			"dataPipeline": "",
		},
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &ListDataPipelineConfigurationsResult{}

	if err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	}); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

type PauseDataPipelineRequest struct {
	Bucket           *string `input:"host,bucket,required"`
	DataPipelineName *string `input:"query,dataPipelineName,required"`
	oss.RequestCommon
}

type PauseDataPipelineResult struct {
	oss.ResultCommon
}

func (c *Client) PauseDataPipeline(ctx context.Context, request *PauseDataPipelineRequest, optFns ...func(*oss.Options)) (*PauseDataPipelineResult, error) {
	var err error
	if request == nil {
		request = &PauseDataPipelineRequest{}
	}

	input := &oss.OperationInput{
		OpName: "PauseDataPipeline",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"action":       "pauseDataPipeline",
			"dataPipeline": "",
		},
		Bucket: request.Bucket,
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PauseDataPipelineResult{}

	if err = c.client.UnmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

type RestartDataPipelineRequest struct {
	DataPipelineName *string `input:"query,dataPipelineName,required"`
	oss.RequestCommon
}

type RestartDataPipelineResult struct {
	oss.ResultCommon
}

func (c *Client) RestartDataPipeline(ctx context.Context, request *RestartDataPipelineRequest, optFns ...func(*oss.Options)) (*RestartDataPipelineResult, error) {
	var err error
	if request == nil {
		request = &RestartDataPipelineRequest{}
	}

	input := &oss.OperationInput{
		OpName: "RestartDataPipeline",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"action":       "restartDataPipeline",
			"dataPipeline": "",
		},
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &RestartDataPipelineResult{}

	if err = c.client.UnmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}
