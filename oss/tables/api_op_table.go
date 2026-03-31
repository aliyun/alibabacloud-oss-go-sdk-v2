package tables

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type CreateTableRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,Namespace,required"`

	Format *string `input:"body,format,json,required"`

	Table *string `input:"body,name,json,required"`

	Metadata *TableMetadata `input:"body,metadata,json,required"`

	// The encryption of the table .
	EncryptionConfiguration *EncryptionConfiguration `input:"body,encryptionConfiguration,json"`

	// The storage class of the table .
	StorageClassConfiguration *StorageClassConfiguration `input:"body,storageClassConfiguration,json"`

	// The tagging of the table .
	Tags map[string]any `input:"body,tags,json"`

	oss.RequestCommon
}

type TableMetadata struct {
	Iceberg *MetadataIceberg `json:"iceberg"`
}

type MetadataIceberg struct {
	Schema map[string]any `json:"schema"`
}

type CreateTableResult struct {
	TableARN *string `json:"tableARN"`

	VersionToken *string `json:"versionToken"`

	oss.ResultCommon
}

// CreateTable Creates a table.
func (c *TablesClient) CreateTable(ctx context.Context, request *CreateTableRequest, optFns ...func(*oss.Options)) (*CreateTableResult, error) {
	var err error
	if request == nil {
		request = &CreateTableRequest{}
	}
	input := &oss.OperationInput{
		OpName: "CreateTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &CreateTableResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetTableRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the table.
	Table *string `input:"query,name,required"`

	Namespace *string `input:"query,namespace,required"`

	TableArn *string `input:"query,tableArn,required"`

	TableBucketARN *string `input:"query,tableBucketARN,required"`

	oss.RequestCommon
}

type GetTableResult struct {
	CreatedAt         *string  `json:"createdAt"`
	CreatedBy         *string  `json:"createdBy"`
	Format            *string  `json:"format"`
	MetadataLocation  *string  `json:"metadataLocation"`
	ModifiedAt        *string  `json:"modifiedAt"`
	ModifiedBy        *string  `json:"modifiedBy"`
	Name              *string  `json:"name"`
	Namespace         []string `json:"namespace"`
	NamespaceId       *string  `json:"namespaceId"`
	OwnerAccountId    *string  `json:"ownerAccountId"`
	TableARN          *string  `json:"tableARN"`
	TableBucketId     *string  `json:"tableBucketId"`
	Type              *string  `json:"type"`
	VersionToken      *string  `json:"versionToken"`
	WarehouseLocation *string  `json:"warehouseLocation"`

	oss.ResultCommon
}

// GetTable Queries information about a table.
func (c *TablesClient) GetTable(ctx context.Context, request *GetTableRequest, optFns ...func(*oss.Options)) (*GetTableResult, error) {
	var err error
	if request == nil {
		request = &GetTableRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTable",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"get-table": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetTableResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type ListTablesRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"query,namespace,required"`

	// The token from which the ListTables operation must start.
	ContinuationToken *string `input:"query,continuationToken"`

	// The maximum number of s that can be returned in the single query.
	// Valid values: 1 to 1000.
	MaxTables int32 `input:"query,maxTables"`

	// The prefix that the names of returned s must contain.
	Prefix *string `input:"query,prefix"` // Limits the response to keys that begin with the specified prefix

	oss.RequestCommon
}

type ListTablesResult struct {
	// The token from which the ListTables operation must start.
	ContinuationToken *string `json:"continuationToken"`

	// The container that stores information about s.
	Tables []TableProperties `json:"tables"`

	oss.ResultCommon
}

type TableProperties struct {
	CreatedAt  *string  `json:"createdAt"`
	ModifiedAt *string  `json:"modifiedAt"`
	Name       *string  `json:"name"`
	Namespace  []string `json:"namespace"`
	TableARN   *string  `json:"tableARN"`
	Type       *string  `json:"type"`
}

// ListTables Lists table s that belong to the current account.
func (c *TablesClient) ListTables(ctx context.Context, request *ListTablesRequest, optFns ...func(*oss.Options)) (*ListTablesResult, error) {
	var err error
	if request == nil {
		request = &ListTablesRequest{}
	}
	input := &oss.OperationInput{
		OpName: "ListTables",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &ListTablesResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteTableRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"query,namespace,required"`

	Table *string `input:"query,name,required"`

	VersionToken *string `input:"query,versionToken"`

	oss.RequestCommon
}

type DeleteTableResult struct {
	oss.ResultCommon
}

// DeleteTable Deletes a table.
func (c *TablesClient) DeleteTable(ctx context.Context, request *DeleteTableRequest, optFns ...func(*oss.Options)) (*DeleteTableResult, error) {
	var err error
	if request == nil {
		request = &DeleteTableRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteTable",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteTableResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type RenameTableRequest struct {
	// The name of the table bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	NewNamespace *string `input:"body,namespace,required,json"`

	NewTable *string `input:"body,newName,required,json"`

	VersionToken *string `input:"body,versionToken,required,json"`

	oss.RequestCommon
}

type RenameTableResult struct {
	oss.ResultCommon
}

// RenameTable Rename a table .
func (c *TablesClient) RenameTable(ctx context.Context, request *RenameTableRequest, optFns ...func(*oss.Options)) (*RenameTableResult, error) {
	var err error
	if request == nil {
		request = &RenameTableRequest{}
	}
	input := &oss.OperationInput{
		OpName: "RenameTable",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
			"rename":                        "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &RenameTableResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
