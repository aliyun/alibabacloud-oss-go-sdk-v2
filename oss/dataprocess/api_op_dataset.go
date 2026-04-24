package dataprocess

import (
	"context"
	"encoding/xml"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// Dataset represents a dataset in OSS Data Process
type Dataset struct {
	XMLName                  xml.Name            `xml:"Dataset"`
	BindCount                *int64              `xml:"BindCount,omitempty"`
	CreateTime               *string             `xml:"CreateTime,omitempty"`
	DatasetMaxBindCount      *int64              `xml:"DatasetMaxBindCount,omitempty"`
	DatasetMaxEntityCount    *int64              `xml:"DatasetMaxEntityCount,omitempty"`
	DatasetMaxFileCount      *int64              `xml:"DatasetMaxFileCount,omitempty"`
	DatasetMaxRelationCount  *int64              `xml:"DatasetMaxRelationCount,omitempty"`
	DatasetMaxTotalFileSize  *int64              `xml:"DatasetMaxTotalFileSize,omitempty"`
	DatasetName              *string             `xml:"DatasetName,omitempty"`
	Description              *string             `xml:"Description,omitempty"`
	FileCount                *int64              `xml:"FileCount,omitempty"`
	TemplateId               *string             `xml:"TemplateId,omitempty"`
	TotalFileSize            *int64              `xml:"TotalFileSize,omitempty"`
	UpdateTime               *string             `xml:"UpdateTime,omitempty"`
	WorkflowParameters       *WorkflowParameters `xml:"WorkflowParameters,omitempty"`
	WorkflowParametersString *string             `xml:"WorkflowParametersString,omitempty"`
	DatasetConfig            *DatasetConfig      `xml:"DatasetConfig,omitempty"`
}

// WorkflowParameters represents the workflow parameters configuration
type WorkflowParameters struct {
	XMLName           xml.Name            `xml:"WorkflowParameters"`
	WorkflowParameter []WorkflowParameter `xml:"WorkflowParameter"`
}

// WorkflowParameter represents a single workflow parameter
type WorkflowParameter struct {
	XMLName     xml.Name `xml:"WorkflowParameter"`
	Name        *string  `xml:"Name,omitempty"`
	Value       *string  `xml:"Value,omitempty"`
	Description *string  `xml:"Description,omitempty"`
}

// DatasetConfig represents the dataset configuration
type DatasetConfig struct {
	XMLName  xml.Name        `xml:"DatasetConfig"`
	Insights *InsightsConfig `xml:"Insights,omitempty"`
}

// InsightsConfig represents the insights configuration
type InsightsConfig struct {
	XMLName     xml.Name `xml:"Insights"`
	EnableLabel *bool    `xml:"EnableLabel,omitempty"`
	EnableOCR   *bool    `xml:"EnableOCR,omitempty"`
	EnableFace  *bool    `xml:"EnableFace,omitempty"`
	EnableImage *bool    `xml:"EnableImage,omitempty"`
	EnableVideo *bool    `xml:"EnableVideo,omitempty"`
	EnableAudio *bool    `xml:"EnableAudio,omitempty"`
	Language    *string  `xml:"Language,omitempty"`
}

// CreateDatasetRequest defines the request for creating a dataset
type CreateDatasetRequest struct {
	Bucket             *string `input:"host,bucket,required"`
	DatasetName        *string `input:"query,datasetName,required"`
	Description        *string `input:"query,description"`
	TemplateId         *string `input:"query,templateId"`
	ClusterType        *string `input:"query,clusterType"`
	WorkflowParameters *string `input:"query,workflowParameters"`
	DatasetConfig      *string `input:"query,datasetConfig"`
	oss.RequestCommon
}

// CreateDatasetResult defines the result for CreateDataset operation
type CreateDatasetResult struct {
	XMLName xml.Name `xml:"CreateDatasetResponse"`
	Dataset *Dataset `xml:"Dataset"`
	oss.ResultCommon
}

// CreateDataset creates a dataset.
func (c *Client) CreateDataset(ctx context.Context, request *CreateDatasetRequest, optFns ...func(*oss.Options)) (*CreateDatasetResult, error) {
	var err error
	if request == nil {
		request = &CreateDatasetRequest{}
	}

	input := &oss.OperationInput{
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

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &CreateDatasetResult{}

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

// GetDatasetRequest defines the request for getting dataset information
type GetDatasetRequest struct {
	Bucket         *string `input:"host,bucket,required"`
	DatasetName    *string `input:"query,datasetName,required"`
	WithStatistics *bool   `input:"query,withStatistics"`
	oss.RequestCommon
}

// GetDatasetResult defines the result for GetDataset operation
type GetDatasetResult struct {
	XMLName xml.Name `xml:"GetDatasetResponse"`
	Dataset *Dataset `xml:"Dataset"`
	oss.ResultCommon
}

// GetDataset gets the information of a dataset.
func (c *Client) GetDataset(ctx context.Context, request *GetDatasetRequest, optFns ...func(*oss.Options)) (*GetDatasetResult, error) {
	var err error
	if request == nil {
		request = &GetDatasetRequest{}
	}

	input := &oss.OperationInput{
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

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetDatasetResult{}

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

// UpdateDatasetRequest defines the request for updating a dataset
type UpdateDatasetRequest struct {
	Bucket             *string `input:"host,bucket,required"`
	DatasetName        *string `input:"query,datasetName,required"`
	Description        *string `input:"query,description"`
	TemplateId         *string `input:"query,templateId"`
	WorkflowParameters *string `input:"query,workflowParameters"`
	DatasetConfig      *string `input:"query,datasetConfig"`
	oss.RequestCommon
}

// UpdateDatasetResult defines the result for UpdateDataset operation
type UpdateDatasetResult struct {
	XMLName xml.Name `xml:"UpdateDatasetResponse"`
	Dataset *Dataset `xml:"Dataset"`
	oss.ResultCommon
}

// UpdateDataset updates a dataset.
func (c *Client) UpdateDataset(ctx context.Context, request *UpdateDatasetRequest, optFns ...func(*oss.Options)) (*UpdateDatasetResult, error) {
	var err error
	if request == nil {
		request = &UpdateDatasetRequest{}
	}

	input := &oss.OperationInput{
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

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &UpdateDatasetResult{}

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

// DeleteDatasetRequest defines the request for deleting a dataset
type DeleteDatasetRequest struct {
	Bucket      *string `input:"host,bucket,required"`
	DatasetName *string `input:"query,datasetName,required"`
	oss.RequestCommon
}

// DeleteDatasetResult defines the result for DeleteDataset operation
type DeleteDatasetResult struct {
	oss.ResultCommon
}

// DeleteDataset deletes a dataset.
func (c *Client) DeleteDataset(ctx context.Context, request *DeleteDatasetRequest, optFns ...func(*oss.Options)) (*DeleteDatasetResult, error) {
	var err error
	if request == nil {
		request = &DeleteDatasetRequest{}
	}

	input := &oss.OperationInput{
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

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteDatasetResult{}

	if err = c.client.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

// ListDatasetsRequest defines the request for listing datasets
type ListDatasetsRequest struct {
	Bucket     *string `input:"host,bucket,required"`
	MaxResults *int64  `input:"query,maxResults"`
	NextToken  *string `input:"query,nextToken"`
	Prefix     *string `input:"query,prefix"`
	oss.RequestCommon
}

// ListDatasetsResult defines the result for ListDatasets operation
type ListDatasetsResult struct {
	XMLName    xml.Name  `xml:"ListDatasetsResponse"`
	Datasets   []Dataset `xml:"Datasets>Dataset,omitempty"`
	NextToken  *string   `xml:"NextToken,omitempty"`
	MaxResults *int64    `xml:"MaxResults,omitempty"`
	oss.ResultCommon
}

// ListDatasets lists datasets.
func (c *Client) ListDatasets(ctx context.Context, request *ListDatasetsRequest, optFns ...func(*oss.Options)) (*ListDatasetsResult, error) {
	var err error
	if request == nil {
		request = &ListDatasetsRequest{}
	}

	input := &oss.OperationInput{
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

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &ListDatasetsResult{}

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
