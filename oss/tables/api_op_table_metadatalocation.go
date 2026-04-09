package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableMetadataLocationRequest struct {
	BucketArn *string `input:"nop,bucketArn,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	oss.RequestCommon
}

type GetTableMetadataLocationResult struct {
	MetadataLocation  *string `json:"metadataLocation"`
	WarehouseLocation *string `json:"warehouseLocation"`
	VersionToken      *string `json:"versionToken"`

	oss.ResultCommon
}

// GetTableMetadataLocation Queries the metadata location of a table.
func (c *TablesClient) GetTableMetadataLocation(ctx context.Context, request *GetTableMetadataLocationRequest, optFns ...func(*oss.Options)) (*GetTableMetadataLocationResult, error) {
	var err error
	if request == nil {
		request = &GetTableMetadataLocationRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableMetadataLocation",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/metadata-location", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Table)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &GetTableMetadataLocationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type UpdateTableMetadataLocationRequest struct {
	BucketArn *string `input:"nop,bucketArn,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	MetadataLocation *string `input:"body,metadataLocation,required,json"`

	VersionToken *string `input:"body,versionToken,required,json"`

	oss.RequestCommon
}

type UpdateTableMetadataLocationResult struct {
	MetadataLocation *string  `json:"metadataLocation"`
	Name             *string  `json:"name"`
	Namespace        []string `json:"namespace"`
	TableArn         *string  `json:"tableARN"`
	VersionToken     *string  `json:"versionToken"`

	oss.ResultCommon
}

// UpdateTableMetadataLocation Update the metadata location of a table.
func (c *TablesClient) UpdateTableMetadataLocation(ctx context.Context, request *UpdateTableMetadataLocationRequest, optFns ...func(*oss.Options)) (*UpdateTableMetadataLocationResult, error) {
	var err error
	if request == nil {
		request = &UpdateTableMetadataLocationRequest{}
	}
	input := &oss.OperationInput{
		OpName: "UpdateTableMetadataLocation",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.BucketArn,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/metadata-location", url.QueryEscape(oss.ToString(request.BucketArn)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Table)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}
	result := &UpdateTableMetadataLocationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
