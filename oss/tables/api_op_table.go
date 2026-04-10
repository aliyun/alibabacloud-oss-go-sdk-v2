package tables

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type CreateTableRequest struct {
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	Namespace *string `input:"nop,namespace,required"`

	Name *string `input:"body,name,json,required"`

	Format *string `input:"body,format,json,required"`

	Metadata *TableMetadata `input:"body,metadata,json"`

	// The encryption of the table .
	EncryptionConfiguration *EncryptionConfiguration `input:"body,encryptionConfiguration,json"`

	oss.RequestCommon
}

type TableMetadata struct {
	Iceberg *IcebergMetadata `json:"iceberg,omitempty"`
}

type IcebergMetadata struct {
	Schema map[string]any `json:"schema,omitempty"`
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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
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
	TableBucketARN *string `input:"query,tableBucketARN"`

	// The name of the table.
	Name *string `input:"query,name"`

	Namespace *string `input:"query,namespace"`

	TableARN *string `input:"query,tableArn"`

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
		Key: oss.Ptr("get-table"),
	}
	if err = checkGetTableRequest(request); err != nil {
		return nil, err
	}
	input.Bucket = parseBucketArn(request)
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
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
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

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
	Tables []TableSummary `json:"tables"`

	oss.ResultCommon
}

type TableSummary struct {
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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
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
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	Namespace *string `input:"nop,namespace,required"`

	Name *string `input:"nop,name,required"`

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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
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
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	Namespace *string `input:"nop,namespace,required"`

	Name *string `input:"nop,name,required"`

	NewNamespace *string `input:"body,newNamespaceName,json"`

	NewName *string `input:"body,newName,json"`

	VersionToken *string `input:"body,versionToken,json"`

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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/rename", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	if err = checkRenameTableRequest(request); err != nil {
		return nil, err
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
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

func checkGetTableRequest(request *GetTableRequest) error {
	if request.TableARN == nil && (request.TableBucketARN == nil || request.Namespace == nil || request.Name == nil) {
		return fmt.Errorf("must provide either table arn alone OR all of (table bucket arn, namespace, table name) together")
	}
	if request.TableARN != nil && (request.TableBucketARN != nil || request.Namespace != nil || request.Name != nil) {
		return fmt.Errorf("must provide either table arn alone OR all of (table bucket arn, namespace, table name) together")
	}
	return nil
}

func checkRenameTableRequest(request *RenameTableRequest) error {
	if request.NewName == nil && request.NewNamespace == nil {
		return fmt.Errorf("either NewTable or NewNamespace must be provided")
	}
	return nil
}

func parseBucketArn(request *GetTableRequest) *string {
	switch {
	case request.TableBucketARN != nil:
		return request.TableBucketARN
	case request.TableARN != nil:
		if vals := strings.Split(oss.ToString(request.TableARN), "/table"); len(vals) > 0 {
			return oss.Ptr(vals[0])
		}
	}
	return nil
}
