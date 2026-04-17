package tables

import (
	"context"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func newTestTablesClient() *TablesClient {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou")
	return NewTablesClient(cfg)
}

// ==================== Table Bucket API ====================

func TestClient_CreateTableBucket_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Create minimal request - Name is in body, validated server-side
	req := &CreateTableBucketRequest{}
	_, err := client.CreateTableBucket(context.Background(), req)
	assert.Error(t, err)
	// Name validation happens server-side for body fields
}

func TestClient_GetTableBucket_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetTableBucketRequest{}
	_, err := client.GetTableBucket(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// not table bucket arn
	req = &GetTableBucketRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
	}
	_, err = client.GetTableBucket(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name (contains ?)
	req = &GetTableBucketRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
	}
	_, err = client.GetTableBucket(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_DeleteTableBucket_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &DeleteTableBucketRequest{}
	_, err := client.DeleteTableBucket(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// not table bucket arn
	req = &DeleteTableBucketRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
	}
	_, err = client.DeleteTableBucket(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name (contains ?)
	req = &DeleteTableBucketRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
	}
	_, err = client.DeleteTableBucket(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_ListTableBuckets_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// ListTableBuckets has no required fields
	req := &ListTableBucketsRequest{}
	_, err := client.ListTableBuckets(context.Background(), req)
	// May fail with network error, but not field validation error
	_ = err
}

// ==================== Namespace API ====================

func TestClient_CreateNamespace_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &CreateNamespaceRequest{
		Namespace: []string{"test-ns"},
	}
	_, err := client.CreateNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// not table bucket arn
	req = &CreateNamespaceRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      []string{"test-ns"},
	}
	_, err = client.CreateNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &CreateNamespaceRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      []string{"test-ns"},
	}
	_, err = client.CreateNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_GetNamespace_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetNamespaceRequest{
		Namespace: oss.Ptr("test-ns"),
	}
	_, err := client.GetNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &GetNamespaceRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
	}
	_, err = client.GetNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// not table bucket arn
	req = &GetNamespaceRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.GetNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &GetNamespaceRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.GetNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_DeleteNamespace_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &DeleteNamespaceRequest{
		Namespace: oss.Ptr("test-ns"),
	}
	_, err := client.DeleteNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &DeleteNamespaceRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
	}
	_, err = client.DeleteNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// not table bucket arn
	req = &DeleteNamespaceRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.DeleteNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &DeleteNamespaceRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.DeleteNamespace(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_ListNamespaces_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &ListNamespacesRequest{}
	_, err := client.ListNamespaces(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// not table bucket arn
	req = &ListNamespacesRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
	}
	_, err = client.ListNamespaces(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &ListNamespacesRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
	}
	_, err = client.ListNamespaces(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

// ==================== Table API ====================

func TestClient_CreateTable_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &CreateTableRequest{
		Namespace: oss.Ptr("test-ns"),
	}
	_, err := client.CreateTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &CreateTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
	}
	_, err = client.CreateTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// not table bucket arn - must provide all required body fields first
	req = &CreateTableRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
		Format:         oss.Ptr("iceberg"),
	}
	_, err = client.CreateTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &CreateTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
		Format:         oss.Ptr("iceberg"),
	}
	_, err = client.CreateTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_GetTable_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// GetTable can use either TableBucketARN+Namespace+Name OR TableARN
	// Test with invalid TableARN
	req := &GetTableRequest{
		TableARN: oss.Ptr("invalid-arn"),
	}
	_, err := client.GetTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// Test with TableARN containing invalid bucket name
	req = &GetTableRequest{
		TableARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/invalid?bucket/table/test"),
	}
	_, err = client.GetTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")

	// Test without any required field (neither TableARN nor TableBucketARN+Namespace+Name)
	req = &GetTableRequest{}
	_, err = client.GetTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")
}

func TestClient_DeleteTable_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &DeleteTableRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.DeleteTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &DeleteTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.DeleteTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &DeleteTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.DeleteTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// not table bucket arn
	req = &DeleteTableRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.DeleteTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &DeleteTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.DeleteTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_ListTables_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &ListTablesRequest{}
	_, err := client.ListTables(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &ListTablesRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
	}
	_, err = client.ListTables(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// not table bucket arn
	req = &ListTablesRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.ListTables(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &ListTablesRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.ListTables(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_RenameTable_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &RenameTableRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("old-table"),
		NewName:   oss.Ptr("new-table"),
	}
	_, err := client.RenameTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &RenameTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("old-table"),
		NewName:        oss.Ptr("new-table"),
	}
	_, err = client.RenameTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &RenameTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
		NewName:        oss.Ptr("new-table"),
	}
	_, err = client.RenameTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// not table bucket arn
	req = &RenameTableRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("old-table"),
		NewName:        oss.Ptr("new-table"),
	}
	_, err = client.RenameTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &RenameTableRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("old-table"),
		NewName:        oss.Ptr("new-table"),
	}
	_, err = client.RenameTable(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

// ==================== Table Bucket Config API ====================

func TestClient_PutTableBucketEncryption_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &PutTableBucketEncryptionRequest{}
	_, err := client.PutTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - EncryptionConfiguration
	req = &PutTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
	}
	_, err = client.PutTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, EncryptionConfiguration")

	// not table bucket arn - must provide all required body fields first
	req = &PutTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		EncryptionConfiguration: &EncryptionConfiguration{
			SseAlgorithm: oss.Ptr("AES256"),
		},
	}
	_, err = client.PutTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &PutTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		EncryptionConfiguration: &EncryptionConfiguration{
			SseAlgorithm: oss.Ptr("AES256"),
		},
	}
	_, err = client.PutTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_GetTableBucketEncryption_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetTableBucketEncryptionRequest{}
	_, err := client.GetTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// not table bucket arn
	req = &GetTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
	}
	_, err = client.GetTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &GetTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
	}
	_, err = client.GetTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_DeleteTableBucketEncryption_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &DeleteTableBucketEncryptionRequest{}
	_, err := client.DeleteTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// not table bucket arn
	req = &DeleteTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
	}
	_, err = client.DeleteTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &DeleteTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
	}
	_, err = client.DeleteTableBucketEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_PutTableBucketPolicy_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &PutTableBucketPolicyRequest{}
	_, err := client.PutTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - ResourcePolicy
	req = &PutTableBucketPolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
	}
	_, err = client.PutTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, ResourcePolicy")

	// not table bucket arn - must provide all required body fields first
	req = &PutTableBucketPolicyRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		ResourcePolicy: oss.Ptr(`{"Version":"1"}`),
	}
	_, err = client.PutTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &PutTableBucketPolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		ResourcePolicy: oss.Ptr(`{"Version":"1"}`),
	}
	_, err = client.PutTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_GetTableBucketPolicy_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetTableBucketPolicyRequest{}
	_, err := client.GetTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// not table bucket arn
	req = &GetTableBucketPolicyRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
	}
	_, err = client.GetTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &GetTableBucketPolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
	}
	_, err = client.GetTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_DeleteTableBucketPolicy_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &DeleteTableBucketPolicyRequest{}
	_, err := client.DeleteTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// not table bucket arn
	req = &DeleteTableBucketPolicyRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
	}
	_, err = client.DeleteTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &DeleteTableBucketPolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
	}
	_, err = client.DeleteTableBucketPolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_PutTableBucketMaintenanceConfiguration_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &PutTableBucketMaintenanceConfigurationRequest{}
	_, err := client.PutTableBucketMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Type
	req = &PutTableBucketMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
	}
	_, err = client.PutTableBucketMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Type")

	// Required field - Value
	req = &PutTableBucketMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Type:           oss.Ptr("full"),
	}
	_, err = client.PutTableBucketMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Value")

	// not table bucket arn - must provide all required body fields first
	req = &PutTableBucketMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Type:           oss.Ptr("full"),
		Value:          &MaintenanceValue{Status: oss.Ptr("enabled")},
	}
	_, err = client.PutTableBucketMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &PutTableBucketMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Type:           oss.Ptr("full"),
		Value:          &MaintenanceValue{Status: oss.Ptr("enabled")},
	}
	_, err = client.PutTableBucketMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_GetTableBucketMaintenanceConfiguration_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetTableBucketMaintenanceConfigurationRequest{}
	_, err := client.GetTableBucketMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// not table bucket arn
	req = &GetTableBucketMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
	}
	_, err = client.GetTableBucketMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &GetTableBucketMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
	}
	_, err = client.GetTableBucketMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

// ==================== Table Config API ====================

func TestClient_GetTableEncryption_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetTableEncryptionRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.GetTableEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &GetTableEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &GetTableEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.GetTableEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// not table bucket arn
	req = &GetTableEncryptionRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &GetTableEncryptionRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableEncryption(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_PutTablePolicy_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &PutTablePolicyRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.PutTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &PutTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.PutTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &PutTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.PutTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// Required field - ResourcePolicy
	req = &PutTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.PutTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, ResourcePolicy")

	// not table bucket arn - must provide all required body fields first
	req = &PutTablePolicyRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
		ResourcePolicy: oss.Ptr(`{"Version":"1"}`),
	}
	_, err = client.PutTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &PutTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
		ResourcePolicy: oss.Ptr(`{"Version":"1"}`),
	}
	_, err = client.PutTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_GetTablePolicy_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetTablePolicyRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.GetTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &GetTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &GetTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.GetTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// not table bucket arn
	req = &GetTablePolicyRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &GetTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_DeleteTablePolicy_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &DeleteTablePolicyRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.DeleteTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &DeleteTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.DeleteTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &DeleteTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.DeleteTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// not table bucket arn
	req = &DeleteTablePolicyRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.DeleteTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &DeleteTablePolicyRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.DeleteTablePolicy(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_PutTableMaintenanceConfiguration_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &PutTableMaintenanceConfigurationRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.PutTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &PutTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.PutTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &PutTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.PutTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// Required field - Type
	req = &PutTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.PutTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Type")

	// Required field - Value
	req = &PutTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
		Type:           oss.Ptr("full"),
	}
	_, err = client.PutTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Value")

	// not table bucket arn - must provide all required body fields first
	req = &PutTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
		Type:           oss.Ptr("full"),
		Value:          &TableMaintenanceValue{Status: oss.Ptr("enabled")},
	}
	_, err = client.PutTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &PutTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
		Type:           oss.Ptr("full"),
		Value:          &TableMaintenanceValue{Status: oss.Ptr("enabled")},
	}
	_, err = client.PutTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_GetTableMaintenanceConfiguration_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetTableMaintenanceConfigurationRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.GetTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &GetTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &GetTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.GetTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// not table bucket arn
	req = &GetTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &GetTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableMaintenanceConfiguration(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_GetTableMaintenanceJobStatus_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetTableMaintenanceJobStatusRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.GetTableMaintenanceJobStatus(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &GetTableMaintenanceJobStatusRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableMaintenanceJobStatus(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &GetTableMaintenanceJobStatusRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.GetTableMaintenanceJobStatus(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// not table bucket arn
	req = &GetTableMaintenanceJobStatusRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableMaintenanceJobStatus(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &GetTableMaintenanceJobStatusRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableMaintenanceJobStatus(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_GetTableMetadataLocation_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &GetTableMetadataLocationRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.GetTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &GetTableMetadataLocationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &GetTableMetadataLocationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.GetTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// not table bucket arn
	req = &GetTableMetadataLocationRequest{
		TableBucketARN: oss.Ptr("bucket-name"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &GetTableMetadataLocationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.GetTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}

func TestClient_UpdateTableMetadataLocation_FieldValidation(t *testing.T) {
	client := newTestTablesClient()
	assert.NotNil(t, client)

	// Required field - TableBucketARN
	req := &UpdateTableMetadataLocationRequest{
		Namespace: oss.Ptr("test-ns"),
		Name:      oss.Ptr("test-table"),
	}
	_, err := client.UpdateTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, TableBucketARN")

	// Required field - Namespace
	req = &UpdateTableMetadataLocationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.UpdateTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Namespace")

	// Required field - Name
	req = &UpdateTableMetadataLocationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
	}
	_, err = client.UpdateTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name")

	// Required field - MetadataLocation
	req = &UpdateTableMetadataLocationRequest{
		TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:      oss.Ptr("test-ns"),
		Name:           oss.Ptr("test-table"),
	}
	_, err = client.UpdateTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, MetadataLocation")

	// Required field - VersionToken
	req = &UpdateTableMetadataLocationRequest{
		TableBucketARN:   oss.Ptr("acs:osstables:cn-beijing:123456:bucket/valid-bucket"),
		Namespace:        oss.Ptr("test-ns"),
		Name:             oss.Ptr("test-table"),
		MetadataLocation: oss.Ptr("oss://bucket/path"),
	}
	_, err = client.UpdateTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field, VersionToken")

	// not table bucket arn - must provide all required body fields first
	req = &UpdateTableMetadataLocationRequest{
		TableBucketARN:   oss.Ptr("bucket-name"),
		Namespace:        oss.Ptr("test-ns"),
		Name:             oss.Ptr("test-table"),
		MetadataLocation: oss.Ptr("oss://bucket/path"),
		VersionToken:     oss.Ptr("v1"),
	}
	_, err = client.UpdateTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed ARN")

	// table bucket arn with invalid bucket name
	req = &UpdateTableMetadataLocationRequest{
		TableBucketARN:   oss.Ptr("acs:osstables:cn-beijing:123456:bucket/test-table?1234"),
		Namespace:        oss.Ptr("test-ns"),
		Name:             oss.Ptr("test-table"),
		MetadataLocation: oss.Ptr("oss://bucket/path"),
		VersionToken:     oss.Ptr("v1"),
	}
	_, err = client.UpdateTableMetadataLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket resource is invalid")
}
