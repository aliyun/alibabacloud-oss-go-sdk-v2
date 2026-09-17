package oss

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
)

// MetadataTableEncryptionConfiguration specifies the encryption configuration for a metadata table.
type MetadataTableEncryptionConfiguration struct {
	// The server-side encryption algorithm. Valid values: AES256.
	SseAlgorithm *string `xml:"SseAlgorithm"`

	// The ARN of the KMS key. This parameter is required only when SseAlgorithm is set to oss:kms.
	KmsKeyArn *string `xml:"KmsKeyArn"`
}

// RecordExpiration specifies the record expiration configuration for the journal table.
type RecordExpiration struct {
	// Specifies whether to enable record expiration. Valid values: ENABLED, DISABLED.
	Expiration *string `xml:"Expiration"`

	// The retention period of records. This parameter is required only when Expiration is set to ENABLED. Valid values: >= 7.
	Days *int `xml:"Days"`
}

// JournalTableConfiguration specifies the configuration of the journal table.
type JournalTableConfiguration struct {
	// The record expiration configuration.
	RecordExpiration *RecordExpiration `xml:"RecordExpiration"`

	// The encryption configuration.
	EncryptionConfiguration *MetadataTableEncryptionConfiguration `xml:"EncryptionConfiguration"`
}

// InventoryTableConfiguration specifies the configuration of the inventory table.
type InventoryTableConfiguration struct {
	// Specifies whether to enable the inventory table. Valid values: ENABLED, DISABLED.
	ConfigurationState *string `xml:"ConfigurationState"`

	// The encryption configuration.
	EncryptionConfiguration *MetadataTableEncryptionConfiguration `xml:"EncryptionConfiguration"`
}

// MetadataConfiguration specifies the metadata table configuration of a bucket.
type MetadataConfiguration struct {
	// The journal table configuration.
	JournalTableConfiguration *JournalTableConfiguration `xml:"JournalTableConfiguration"`

	// The inventory table configuration.
	InventoryTableConfiguration *InventoryTableConfiguration `xml:"InventoryTableConfiguration"`
}

// DestinationResult specifies the destination information of the metadata table.
type DestinationResult struct {
	// The type of the table bucket. The value is fixed to oss.
	TableBucketType *string `xml:"TableBucketType"`

	// The ARN of the table bucket.
	TableBucketArn *string `xml:"TableBucketArn"`

	// The namespace of the table bucket. The value is in the format of b_<bucket-name>.
	TableNamespace *string `xml:"TableNamespace"`
}

// MetadataTableConfigurationError specifies the error information of a metadata table.
type MetadataTableConfigurationError struct {
	// The error code.
	ErrorCode *string `xml:"ErrorCode"`

	// The error message.
	ErrorMessage *string `xml:"ErrorMessage"`
}

// JournalTableConfigurationResult specifies the configuration result of the journal table.
type JournalTableConfigurationResult struct {
	// The status of the journal table. Valid values: CREATING, ACTIVE, FAILED.
	TableStatus *string `xml:"TableStatus"`

	// The name of the journal table. The value is fixed to journal.
	TableName *string `xml:"TableName"`

	// The ARN of the journal table.
	TableArn *string `xml:"TableArn"`

	// The record expiration configuration.
	RecordExpiration *RecordExpiration `xml:"RecordExpiration"`

	// The encryption configuration.
	EncryptionConfiguration *MetadataTableEncryptionConfiguration `xml:"EncryptionConfiguration"`

	// The error information. This parameter is returned only when TableStatus is FAILED.
	Error *MetadataTableConfigurationError `xml:"Error"`
}

// InventoryTableConfigurationResult specifies the configuration result of the inventory table.
type InventoryTableConfigurationResult struct {
	// Specifies whether the inventory table is enabled. Valid values: ENABLED, DISABLED.
	ConfigurationState *string `xml:"ConfigurationState"`

	// The status of the inventory table. Valid values: CREATING, BACKFILLING, ACTIVE, FAILED.
	TableStatus *string `xml:"TableStatus"`

	// The name of the inventory table. The value is fixed to inventory.
	TableName *string `xml:"TableName"`

	// The ARN of the inventory table.
	TableArn *string `xml:"TableArn"`

	// The encryption configuration.
	EncryptionConfiguration *MetadataTableEncryptionConfiguration `xml:"EncryptionConfiguration"`

	// The error information. This parameter is returned only when TableStatus is FAILED.
	Error *MetadataTableConfigurationError `xml:"Error"`
}

// MetadataConfigurationResult specifies the metadata table configuration result of a bucket.
type MetadataConfigurationResult struct {
	// The destination information.
	DestinationResult *DestinationResult `xml:"DestinationResult"`

	// The journal table configuration result.
	JournalTableConfigurationResult *JournalTableConfigurationResult `xml:"JournalTableConfigurationResult"`

	// The inventory table configuration result.
	InventoryTableConfigurationResult *InventoryTableConfigurationResult `xml:"InventoryTableConfigurationResult"`
}

type CreateBucketMetadataTableConfigurationRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The request body schema.
	MetadataConfiguration *MetadataConfiguration `input:"body,MetadataConfiguration,xml,required"`

	RequestCommon
}

type CreateBucketMetadataTableConfigurationResult struct {
	ResultCommon
}

// CreateBucketMetadataTableConfiguration Creates the metadata table configuration for a bucket.
func (c *Client) CreateBucketMetadataTableConfiguration(ctx context.Context, request *CreateBucketMetadataTableConfigurationRequest, optFns ...func(*Options)) (*CreateBucketMetadataTableConfigurationResult, error) {
	var err error
	if request == nil {
		request = &CreateBucketMetadataTableConfigurationRequest{}
	}
	input := &OperationInput{
		OpName: "CreateBucketMetadataTableConfiguration",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metadataConfiguration"})

	if err = c.marshalInput(request, input, MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &CreateBucketMetadataTableConfigurationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type GetBucketMetadataTableConfigurationRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	RequestCommon
}

type GetBucketMetadataTableConfigurationResult struct {
	// The container that stores the metadata table configuration result.
	GetBucketMetadataConfigurationResult *GetBucketMetadataConfigurationResult `output:"body,GetBucketMetadataConfigurationResult,xml"`

	ResultCommon
}

type GetBucketMetadataConfigurationResult struct {
	MetadataConfigurationResult *MetadataConfigurationResult `xml:"MetadataConfigurationResult"`
}

// GetBucketMetadataTableConfiguration Queries the metadata table configuration of a bucket.
func (c *Client) GetBucketMetadataTableConfiguration(ctx context.Context, request *GetBucketMetadataTableConfigurationRequest, optFns ...func(*Options)) (*GetBucketMetadataTableConfigurationResult, error) {
	var err error
	if request == nil {
		request = &GetBucketMetadataTableConfigurationRequest{}
	}
	input := &OperationInput{
		OpName: "GetBucketMetadataTableConfiguration",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metadataConfiguration"})

	if err = c.marshalInput(request, input, MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetBucketMetadataTableConfigurationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteBucketMetadataTableConfigurationRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	RequestCommon
}

type DeleteBucketMetadataTableConfigurationResult struct {
	ResultCommon
}

// DeleteBucketMetadataTableConfiguration Deletes the metadata table configuration of a bucket.
func (c *Client) DeleteBucketMetadataTableConfiguration(ctx context.Context, request *DeleteBucketMetadataTableConfigurationRequest, optFns ...func(*Options)) (*DeleteBucketMetadataTableConfigurationResult, error) {
	var err error
	if request == nil {
		request = &DeleteBucketMetadataTableConfigurationRequest{}
	}
	input := &OperationInput{
		OpName: "DeleteBucketMetadataTableConfiguration",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metadataConfiguration"})

	if err = c.marshalInput(request, input, MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteBucketMetadataTableConfigurationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type UpdateBucketMetadataInventoryTableConfigurationRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The request body schema.
	InventoryTableConfiguration *InventoryTableConfiguration `input:"body,InventoryTableConfiguration,xml,required"`

	RequestCommon
}

type UpdateBucketMetadataInventoryTableConfigurationResult struct {
	ResultCommon
}

// UpdateBucketMetadataInventoryTableConfiguration Enables or disables the inventory table of a bucket.
func (c *Client) UpdateBucketMetadataInventoryTableConfiguration(ctx context.Context, request *UpdateBucketMetadataInventoryTableConfigurationRequest, optFns ...func(*Options)) (*UpdateBucketMetadataInventoryTableConfigurationResult, error) {
	var err error
	if request == nil {
		request = &UpdateBucketMetadataInventoryTableConfigurationRequest{}
	}
	input := &OperationInput{
		OpName: "UpdateBucketMetadataInventoryTableConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataInventoryTable": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metadataInventoryTable"})

	if err = c.marshalInput(request, input, MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &UpdateBucketMetadataInventoryTableConfigurationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type UpdateBucketMetadataJournalTableConfigurationRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The request body schema.
	JournalTableConfiguration *JournalTableConfiguration `input:"body,JournalTableConfiguration,xml,required"`

	RequestCommon
}

type UpdateBucketMetadataJournalTableConfigurationResult struct {
	ResultCommon
}

// UpdateBucketMetadataJournalTableConfiguration Updates the record expiration configuration of the journal table for a bucket.
func (c *Client) UpdateBucketMetadataJournalTableConfiguration(ctx context.Context, request *UpdateBucketMetadataJournalTableConfigurationRequest, optFns ...func(*Options)) (*UpdateBucketMetadataJournalTableConfigurationResult, error) {
	var err error
	if request == nil {
		request = &UpdateBucketMetadataJournalTableConfigurationRequest{}
	}
	input := &OperationInput{
		OpName: "UpdateBucketMetadataJournalTableConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataJournalTable": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metadataJournalTable"})

	if err = c.marshalInput(request, input, MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &UpdateBucketMetadataJournalTableConfigurationResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
