package tables

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/arn"
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

	// The schema of the table in the v2 format. Compared with Schema, it also describes
	// the struct, list and map field types and their nested combinations.
	//
	// Build an element with IcebergSchemaV2.ToMap, or write the map directly.
	//
	// Only one of Schema and SchemaV2 can be set in the same request.
	SchemaV2 map[string]any `json:"schemaV2,omitempty"`
}

// IcebergSchemaV2 defines the schema of an iceberg table in the v2 format. Compared with the
// schema in the v1 format, it also describes the struct, list and map field types and their
// nested combinations.
//
// The request accepts the raw JSON object, see IcebergSchemaV2.ToMap.
type IcebergSchemaV2 struct {
	// The type of the schema. The value is struct.
	Type *string `json:"type,omitempty"`

	// The id of the schema.
	SchemaID *int64 `json:"schema-id,omitempty"`

	// The ids of the fields that identify a row.
	IdentifierFieldIDs []int64 `json:"identifier-field-ids,omitempty"`

	// The fields of the schema.
	Fields IcebergFields `json:"fields,omitempty"`
}

// ToMap returns the schema as the raw JSON object that the service expects, so that the result
// can be assigned to IcebergMetadata.SchemaV2.
func (s IcebergSchemaV2) ToMap() map[string]any {
	m := make(map[string]any)
	if s.Type != nil {
		m["type"] = *s.Type
	}
	if s.SchemaID != nil {
		m["schema-id"] = *s.SchemaID
	}
	if len(s.IdentifierFieldIDs) > 0 {
		m["identifier-field-ids"] = s.IdentifierFieldIDs
	}
	if len(s.Fields) > 0 {
		m["fields"] = s.Fields.ToMaps()
	}
	return m
}

// IcebergField defines a field of the iceberg table schema.
//
// Type is either the name of a primitive type, such as "long" and "string", or one of
// IcebergStructType, IcebergListType and IcebergMapType for the nested field types.
type IcebergField struct {
	// The id of the field.
	ID *int64 `json:"id,omitempty"`

	// The name of the field.
	Name *string `json:"name,omitempty"`

	// The type of the field.
	Type any `json:"type,omitempty"`

	// Specifies whether the field is required.
	Required *bool `json:"required,omitempty"`

	// The description of the field.
	Doc *string `json:"doc,omitempty"`
}

// ToMap returns the field definition as the raw JSON object that the service expects, so that the
// result can be used as an element of IcebergSchemaV2.Fields.
func (s IcebergField) ToMap() map[string]any {
	m := make(map[string]any)
	if s.ID != nil {
		m["id"] = *s.ID
	}
	if s.Name != nil {
		m["name"] = *s.Name
	}
	if s.Type != nil {
		m["type"] = icebergFieldTypeToAny(s.Type)
	}
	if s.Required != nil {
		m["required"] = *s.Required
	}
	if s.Doc != nil {
		m["doc"] = *s.Doc
	}
	return m
}

// IcebergFields is a list of fields of the iceberg table schema. It converts the strongly-typed
// IcebergField model to the raw JSON objects that the service expects.
type IcebergFields []IcebergField

// ToMaps returns the fields as the raw JSON objects that the service expects, so that the result
// can be assigned to IcebergSchemaV2.Fields. It is the batch counterpart of IcebergField.ToMap.
func (s IcebergFields) ToMaps() []map[string]any {
	maps := make([]map[string]any, 0, len(s))
	for _, field := range s {
		maps = append(maps, field.ToMap())
	}
	return maps
}

// IcebergStructType defines the struct field type of the iceberg table schema.
type IcebergStructType struct {
	// The type of the field. The value is struct.
	Type *string `json:"type,omitempty"`

	// The fields of the struct.
	Fields IcebergFields `json:"fields,omitempty"`
}

// ToMap returns the struct type as the raw JSON object that the service expects, so that the result
// can be used as the Type of an IcebergField.
func (s IcebergStructType) ToMap() map[string]any {
	m := make(map[string]any)
	if s.Type != nil {
		m["type"] = *s.Type
	}
	if len(s.Fields) > 0 {
		m["fields"] = s.Fields.ToMaps()
	}
	return m
}

// IcebergListType defines the list field type of the iceberg table schema.
type IcebergListType struct {
	// The type of the field. The value is list.
	Type *string `json:"type,omitempty"`

	// The id of the element.
	ElementID *int64 `json:"element-id,omitempty"`

	// The type of the element, either the name of a primitive type or one of the nested field types.
	Element any `json:"element,omitempty"`

	// Specifies whether the element is required.
	ElementRequired *bool `json:"element-required,omitempty"`
}

// ToMap returns the list type as the raw JSON object that the service expects, so that the result
// can be used as the Type of an IcebergField.
func (s IcebergListType) ToMap() map[string]any {
	m := make(map[string]any)
	if s.Type != nil {
		m["type"] = *s.Type
	}
	if s.ElementID != nil {
		m["element-id"] = *s.ElementID
	}
	if s.Element != nil {
		m["element"] = icebergFieldTypeToAny(s.Element)
	}
	if s.ElementRequired != nil {
		m["element-required"] = *s.ElementRequired
	}
	return m
}

// IcebergMapType defines the map field type of the iceberg table schema.
type IcebergMapType struct {
	// The type of the field. The value is map.
	Type *string `json:"type,omitempty"`

	// The id of the key.
	KeyID *int64 `json:"key-id,omitempty"`

	// The type of the key, either the name of a primitive type or one of the nested field types.
	Key any `json:"key,omitempty"`

	// The id of the value.
	ValueID *int64 `json:"value-id,omitempty"`

	// The type of the value, either the name of a primitive type or one of the nested field types.
	Value any `json:"value,omitempty"`

	// Specifies whether the value is required.
	ValueRequired *bool `json:"value-required,omitempty"`
}

// ToMap returns the map type as the raw JSON object that the service expects, so that the result
// can be used as the Type of an IcebergField.
func (s IcebergMapType) ToMap() map[string]any {
	m := make(map[string]any)
	if s.Type != nil {
		m["type"] = *s.Type
	}
	if s.KeyID != nil {
		m["key-id"] = *s.KeyID
	}
	if s.Key != nil {
		m["key"] = icebergFieldTypeToAny(s.Key)
	}
	if s.ValueID != nil {
		m["value-id"] = *s.ValueID
	}
	if s.Value != nil {
		m["value"] = icebergFieldTypeToAny(s.Value)
	}
	if s.ValueRequired != nil {
		m["value-required"] = *s.ValueRequired
	}
	return m
}

// icebergFieldTypeToAny converts the nested field types to the raw JSON objects, the primitive
// types are returned as they are.
func icebergFieldTypeToAny(t any) any {
	switch v := t.(type) {
	case IcebergStructType:
		return v.ToMap()
	case *IcebergStructType:
		if v != nil {
			return v.ToMap()
		}
	case IcebergListType:
		return v.ToMap()
	case *IcebergListType:
		if v != nil {
			return v.ToMap()
		}
	case IcebergMapType:
		return v.ToMap()
	case *IcebergMapType:
		if v != nil {
			return v.ToMap()
		}
	}
	return t
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
	if err = checkCreateTableRequest(request); err != nil {
		return nil, err
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
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
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
	input.Bucket, err = parseBucketArn(request)
	if err != nil {
		return nil, err
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
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
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
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
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
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
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
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

func checkCreateTableRequest(request *CreateTableRequest) error {
	if request.Metadata != nil && request.Metadata.Iceberg != nil {
		iceberg := request.Metadata.Iceberg
		if iceberg.Schema != nil && iceberg.SchemaV2 != nil {
			return fmt.Errorf("Schema and SchemaV2 are mutually exclusive, only one of them can be provided")
		}
	}
	return nil
}

func checkGetTableRequest(request *GetTableRequest) error {
	if request.TableBucketARN != nil {
		if request.TableARN != nil {
			return fmt.Errorf("must provide either table arn alone OR all of (table bucket arn, namespace, table name) together")
		}
		if request.Namespace == nil {
			return oss.NewErrParamRequired("Namespace")
		}
		if request.Name == nil {
			return oss.NewErrParamRequired("Name")
		}
	} else if request.TableARN != nil {
		// PASS
	} else {
		return oss.NewErrParamRequired("TableBucketARN")
	}

	return nil
}

func checkRenameTableRequest(request *RenameTableRequest) error {
	if request.NewName == nil && request.NewNamespace == nil {
		return fmt.Errorf("either NewTable or NewNamespace must be provided")
	}
	return nil
}

func parseBucketArn(request *GetTableRequest) (*string, error) {
	if request.TableARN != nil {
		parsedArn, err := arn.ParseArn(oss.ToString(request.TableARN))
		if err != nil {
			return nil, err
		}
		if oss.ToString(parsedArn.Resource().ResourceType()) == "bucket" && strings.HasPrefix(oss.ToString(parsedArn.Resource().Qualifier()), "table/") {
			if vals := strings.Split(oss.ToString(request.TableARN), "/table"); len(vals) > 0 {
				return oss.Ptr(vals[0]), err
			}
		} else {
			return nil, fmt.Errorf("malformed table arn")
		}
	}

	return request.TableBucketARN, nil
}
