//go:build integration

package tables

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var (
	// Endpoint/ID/Key
	region_     = os.Getenv("OSS_TEST_REGION")
	endpoint_   = os.Getenv("OSS_TEST_TABLES_ENDPOINT")
	accessID_   = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	accessKey_  = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")
	accountUid_ = os.Getenv("OSS_TEST_ACCOUNT_ID")

	instance_ *TablesClient
	testOnce_ sync.Once
)

var (
	bucketNamePrefix = "go-sdk-test-bucket-"
	spaceNamePrefix  = "go_sdk_space"
	tableNamePrefix  = "go_sdk_table"
	letters          = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func getDefaultClient() *TablesClient {
	testOnce_.Do(func() {
		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
			WithRegion(region_).
			WithEndpoint(endpoint_)
		instance_ = NewTablesClient(cfg)
	})
	return instance_
}

func getInvalidAkClient() *TablesClient {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("invalid-ak", "invalid-sk")).
		WithRegion(region_).
		WithEndpoint(endpoint_)
	return NewTablesClient(cfg)
}

func getClient(region, endpoint string) *TablesClient {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region).
		WithEndpoint(endpoint)
	return NewTablesClient(cfg)
}

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func randLowStr(n int) string {
	return strings.ToLower(randStr(n))
}

func cleanBucket(bucketInfo TableBucketProperties, t *testing.T) {
	assert.NotEmpty(t, *bucketInfo.Arn)
	var c *TablesClient
	c = getDefaultClient()
	assert.NotNil(t, c)
	cleanTablesAndNamespaces(c, *bucketInfo.Arn, t)
}

func deleteBucket(bucketArn string, t *testing.T) {
	assert.NotEmpty(t, bucketArn)
	var c *TablesClient
	c = getDefaultClient()
	assert.NotNil(t, c)
	cleanTablesAndNamespaces(c, bucketArn, t)
}

func cleanBuckets(prefix string, t *testing.T) {
	c := getDefaultClient()
	for {
		request := &ListTableBucketsRequest{
			Prefix: oss.Ptr(prefix),
		}
		result, err := c.ListTableBuckets(context.TODO(), request)
		assert.Nil(t, err)
		if len(result.Buckets) == 0 {
			return
		}
		for _, b := range result.Buckets {
			cleanBucket(b, t)
		}
	}
}

func cleanTablesAndNamespaces(c *TablesClient, bucketArn string, t *testing.T) {
	var err error
	var listNamespaceRequest *ListNamespacesRequest
	var listTablesRequest *ListTablesRequest
	listNamespaceRequest = &ListNamespacesRequest{
		BucketArn: oss.Ptr(bucketArn),
	}
	pagNamespaces := c.NewListNameSpacesPaginator(listNamespaceRequest)
	for pagNamespaces.HasNext() {
		page, err := pagNamespaces.NextPage(context.TODO())
		dumpErrIfNotNil(err)
		assert.Nil(t, err)
		for _, namespace := range page.Namespaces {
			listTablesRequest = &ListTablesRequest{
				BucketArn: oss.Ptr(bucketArn),
				Namespace: oss.Ptr(namespace.Namespace[0]),
			}
			pagTables := c.NewListTablesPaginator(listTablesRequest)
			for pagTables.HasNext() {
				page2, err := pagTables.NextPage(context.TODO())
				dumpErrIfNotNil(err)
				assert.Nil(t, err)
				for _, table := range page2.Tables {
					_, err = c.DeleteTable(context.TODO(), &DeleteTableRequest{
						BucketArn: oss.Ptr(bucketArn),
						Namespace: oss.Ptr(namespace.Namespace[0]),
						Name:      table.Name,
					})
					assert.Nil(t, err)
				}
			}
			_, err = c.DeleteNamespace(context.TODO(), &DeleteNamespaceRequest{
				BucketArn: oss.Ptr(bucketArn),
				Namespace: oss.Ptr(namespace.Namespace[0]),
			})
			dumpErrIfNotNil(err)
			assert.Nil(t, err)
		}
	}
	_, err = c.DeleteTableBucket(context.TODO(), &DeleteTableBucketRequest{
		BucketArn: oss.Ptr(bucketArn),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)

}

func before(_ *testing.T) func(t *testing.T) {

	//fmt.Println("setup test case")
	return after
}

func after(t *testing.T) {
	cleanBuckets(bucketNamePrefix, t)
	//fmt.Println("teardown  test case")
}

func dumpErrIfNotNil(err error) {
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}
}

func calcMd5(input string) string {
	if len(input) == 0 {
		return "1B2M2Y8AsgTpgAmY7PhCfg=="
	}
	h := md5.New()
	h.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func TestInvokeOperation(t *testing.T) {
	after := before(t)
	defer after(t)
	BucketName := bucketNamePrefix + randLowStr(5)

	body := `{"name":"` + BucketName + `"}`
	contentMd5 := calcMd5(body)
	//TODO
	input := &oss.OperationInput{
		OpName: "CreateTableBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
			oss.HTTPHeaderContentMD5:  contentMd5,
		},
		Key:  oss.Ptr("buckets"),
		Body: strings.NewReader(body),
	}

	client := getDefaultClient()
	_, err := client.InvokeOperation(context.TODO(), input)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)

	_, err = client.InvokeOperation(context.TODO(), nil)
	assert.NotNil(t, err)
}

func TestInvokeOperation_TableBucketPolicy(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	// PutBucketPolicy
	policy := `{"resourcePolicy":"{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:osstable:cn-beijing:1234567890:bucket/demo-bucket\"]}]}"}`
	input := &oss.OperationInput{
		OpName: "PutTableBucketPolicy",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
			oss.HTTPHeaderContentMD5:  calcMd5(policy),
		},
		Bucket: bucketArn,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/policy", url.QueryEscape(oss.ToString(bucketArn)))),
		Body:   strings.NewReader(policy),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	output, err := client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)

	// GetBucketPolicy
	input = &oss.OperationInput{
		OpName: "GetBucketPolicy",
		Method: "GET",
		Bucket: bucketArn,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/policy", url.QueryEscape(oss.ToString(bucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	output, err = client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)
	policy1, err := io.ReadAll(output.Body)
	assert.NoError(t, err)
	if output.Body != nil {
		output.Body.Close()
	}
	assert.NotEmpty(t, policy1)

	// DeleteBucketPolicy
	input = &oss.OperationInput{
		OpName: "DeleteBucketPolicy",
		Method: "DELETE",
		Bucket: bucketArn,
		Key:    oss.Ptr(fmt.Sprintf("buckets/%s/policy", url.QueryEscape(oss.ToString(bucketArn)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	output, err = client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)
	// discard body
	_, err = io.ReadAll(output.Body)
	assert.NoError(t, err)
	if output.Body != nil {
		output.Body.Close()
	}
}

func TestTableBucket(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	_, err = client.GetTableBucket(context.TODO(), &GetTableBucketRequest{
		BucketArn: bucketArn,
	})
	assert.Nil(t, err)

	list, err := client.ListTableBuckets(context.TODO(), &ListTableBucketsRequest{
		Prefix: oss.Ptr(bucketNamePrefix),
	})
	assert.Nil(t, err)
	assert.True(t, len(list.Buckets) > 0)

	_, err = client.DeleteTableBucket(context.TODO(), &DeleteTableBucketRequest{
		BucketArn: bucketArn,
	})
	assert.Nil(t, err)

	// test server error
	invalidAkClient := getInvalidAkClient()
	bucketNameNotExist := bucketNamePrefix + "not-exist"

	_, err = invalidAkClient.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketNameNotExist),
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.GetTableBucket(context.TODO(), &GetTableBucketRequest{
		BucketArn: bucketArn,
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "Not Found", serr.Code)
	assert.Equal(t, "The specified table bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.ListTableBuckets(context.TODO(), &ListTableBucketsRequest{
		Prefix: oss.Ptr(bucketNamePrefix),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DeleteTableBucket(context.TODO(), &DeleteTableBucketRequest{
		BucketArn: bucketArn,
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "Not Found", serr.Code)
	assert.Equal(t, "The specified table bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestTableBucketEncryption(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	putResult, err := client.PutTableBucketEncryption(context.TODO(), &PutTableBucketEncryptionRequest{
		BucketArn: bucketArn,
		EncryptionConfiguration: &EncryptionConfiguration{
			SseAlgorithm: oss.Ptr("AES256"),
		},
	})
	assert.Nil(t, err)

	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetTableBucketEncryption(context.TODO(), &GetTableBucketEncryptionRequest{
		BucketArn: bucketArn,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteTableBucketEncryption(context.TODO(), &DeleteTableBucketEncryptionRequest{
		BucketArn: bucketArn,
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()

	_, err = invalidAkClient.PutTableBucketEncryption(context.TODO(), &PutTableBucketEncryptionRequest{
		BucketArn: bucketArn,
		EncryptionConfiguration: &EncryptionConfiguration{
			SseAlgorithm: oss.Ptr("AES256"),
		},
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.GetTableBucketEncryption(context.TODO(), &GetTableBucketEncryptionRequest{
		BucketArn: bucketArn,
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.DeleteTableBucketEncryption(context.TODO(), &DeleteTableBucketEncryptionRequest{
		BucketArn: bucketArn,
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestTableBucketPolicy(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	putResult, err := client.PutTableBucketPolicy(context.TODO(), &PutTableBucketPolicyRequest{
		BucketArn:      bucketArn,
		ResourcePolicy: oss.Ptr(`{"Version":"1","Statement":[{"Action":["oss:GetTable"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:osstable:cn-beijing:1234567890:bucket/demo-bucket"]}]}`),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetTableBucketPolicy(context.TODO(), &GetTableBucketPolicyRequest{
		BucketArn: bucketArn,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteTableBucketPolicy(context.TODO(), &DeleteTableBucketPolicyRequest{
		BucketArn: bucketArn,
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()

	_, err = invalidAkClient.PutTableBucketPolicy(context.TODO(), &PutTableBucketPolicyRequest{
		BucketArn: bucketArn,
		ResourcePolicy: oss.Ptr(`(
		{
			"Version":"1",
			"Statement":[{
				"Action":["oss:GetTable"],
				"Effect":"Deny",
				"Principal":["1234567890"],
				"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket"]
			}]
		}`),
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.GetTableBucketPolicy(context.TODO(), &GetTableBucketPolicyRequest{
		BucketArn: bucketArn,
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "Not Found", serr.Code)
	assert.Equal(t, "The specified bucket policy does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.DeleteTableBucketPolicy(context.TODO(), &DeleteTableBucketPolicyRequest{
		BucketArn: bucketArn,
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestTableBucketMaintenance(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	putResult, err := client.PutTableBucketMaintenanceConfiguration(context.TODO(), &PutTableBucketMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
		Type:      oss.Ptr("icebergUnreferencedFileRemoval"),
		Value: &MaintenanceValue{
			Settings: &MaintenanceSettings{
				IcebergUnreferencedFileRemoval: &SettingsDetail{
					UnreferencedDays: oss.Ptr(1),
					NonCurrentDays:   oss.Ptr(10),
				},
			},
			Status: oss.Ptr("enabled"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetTableBucketMaintenanceConfiguration(context.TODO(), &GetTableBucketMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()
	_, err = invalidAkClient.PutTableBucketMaintenanceConfiguration(context.TODO(), &PutTableBucketMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
		Type:      oss.Ptr("icebergUnreferencedFileRemoval"),
		Value: &MaintenanceValue{
			Settings: &MaintenanceSettings{
				IcebergUnreferencedFileRemoval: &SettingsDetail{
					UnreferencedDays: oss.Ptr(1),
					NonCurrentDays:   oss.Ptr(10),
				},
			},
			Status: oss.Ptr("enabled"),
		},
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.GetTableBucketMaintenanceConfiguration(context.TODO(), &GetTableBucketMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestNamespace(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	spaceName := spaceNamePrefix + "_" + randLowStr(5)
	putResult, err := client.CreateNamespace(context.TODO(), &CreateNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: []string{spaceName},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetNamespace(context.TODO(), &GetNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listResult, err := client.ListNamespaces(context.TODO(), &ListNamespacesRequest{
		BucketArn: bucketArn,
		Prefix:    oss.Ptr(spaceNamePrefix),
	})
	assert.Nil(t, err)
	assert.True(t, len(listResult.Namespaces) > 0)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteNamespace(context.TODO(), &DeleteNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()
	_, err = invalidAkClient.CreateNamespace(context.TODO(), &CreateNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: []string{spaceName},
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.GetNamespace(context.TODO(), &GetNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "Not Found", serr.Code)
	assert.Equal(t, "The specified namespace does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.ListNamespaces(context.TODO(), &ListNamespacesRequest{
		BucketArn: bucketArn,
		Prefix:    oss.Ptr(spaceNamePrefix),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DeleteNamespace(context.TODO(), &DeleteNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "Not Found", serr.Code)
	assert.Equal(t, "The specified namespace does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestTable(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	spaceName := spaceNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateNamespace(context.TODO(), &CreateNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: []string{spaceName},
	})
	assert.Nil(t, err)

	tableName := tableNamePrefix + "_" + randLowStr(5)
	putResult, err := client.CreateTable(context.TODO(), &CreateTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Format:    oss.Ptr("ICEBERG"),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetTable(context.TODO(), &GetTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listResult, err := client.ListTables(context.TODO(), &ListTablesRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Prefix:    oss.Ptr(tableNamePrefix),
	})
	assert.Nil(t, err)
	assert.True(t, len(listResult.Tables) > 0)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	newTableName := tableNamePrefix + "_" + randLowStr(6)
	reResult, err := client.RenameTable(context.TODO(), &RenameTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		NewName:   oss.Ptr(newTableName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, reResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteTable(context.TODO(), &DeleteTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(newTableName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()
	_, err = invalidAkClient.CreateTable(context.TODO(), &CreateTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Format:    oss.Ptr("ICEBERG"),
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.GetTable(context.TODO(), &GetTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "Not Found", serr.Code)
	assert.Equal(t, "The specified table does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.ListTables(context.TODO(), &ListTablesRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Prefix:    oss.Ptr(tableNamePrefix),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.RenameTable(context.TODO(), &RenameTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		NewName:   oss.Ptr(newTableName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DeleteTable(context.TODO(), &DeleteTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "Not Found", serr.Code)
	assert.Equal(t, "The specified table does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestTableMetadataLocation(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	spaceName := spaceNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateNamespace(context.TODO(), &CreateNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: []string{spaceName},
	})
	assert.Nil(t, err)

	tableName := tableNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateTable(context.TODO(), &CreateTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Format:    oss.Ptr("ICEBERG"),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getResult, err := client.GetTableMetadataLocation(context.TODO(), &GetTableMetadataLocationRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()
	_, err = invalidAkClient.GetTableMetadataLocation(context.TODO(), &GetTableMetadataLocationRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.UpdateTableMetadataLocation(context.TODO(), &UpdateTableMetadataLocationRequest{
		BucketArn:        bucketArn,
		Namespace:        oss.Ptr(spaceName),
		Name:             oss.Ptr(tableName),
		VersionToken:     getResult.VersionToken,
		MetadataLocation: getResult.MetadataLocation,
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "Bad Request", serr.Code)
	assert.Equal(t, "The specified metadata location is invalid.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestTableEncryption(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	spaceName := spaceNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateNamespace(context.TODO(), &CreateNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: []string{spaceName},
	})
	assert.Nil(t, err)

	tableName := tableNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateTable(context.TODO(), &CreateTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Format:    oss.Ptr("ICEBERG"),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getResult, err := client.GetTableEncryption(context.TODO(), &GetTableEncryptionRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()
	_, err = invalidAkClient.GetTableEncryption(context.TODO(), &GetTableEncryptionRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestTablePolicy(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	spaceName := spaceNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateNamespace(context.TODO(), &CreateNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: []string{spaceName},
	})
	assert.Nil(t, err)

	tableName := tableNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateTable(context.TODO(), &CreateTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Format:    oss.Ptr("ICEBERG"),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	policy := `{
			   "Version":"1",
			   "Statement":[
			   {
				 "Action":[
				   "oss:GetTable"
				],
				"Effect":"Deny",
				"Principal":["1234567890"],
				"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/*"]
			   }
			  ]
			 }`
	putResult, err := client.PutTablePolicy(context.TODO(), &PutTablePolicyRequest{
		BucketArn:      bucketArn,
		Namespace:      oss.Ptr(spaceName),
		Name:           oss.Ptr(tableName),
		ResourcePolicy: oss.Ptr(policy),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetTablePolicy(context.TODO(), &GetTablePolicyRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteTablePolicy(context.TODO(), &DeleteTablePolicyRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()
	_, err = invalidAkClient.PutTablePolicy(context.TODO(), &PutTablePolicyRequest{
		BucketArn:      bucketArn,
		Namespace:      oss.Ptr(spaceName),
		Name:           oss.Ptr(tableName),
		ResourcePolicy: oss.Ptr(policy),
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.GetTablePolicy(context.TODO(), &GetTablePolicyRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "Not Found", serr.Code)
	assert.Equal(t, "The specified table policy does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.DeleteTablePolicy(context.TODO(), &DeleteTablePolicyRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestTableMaintenance(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	spaceName := spaceNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateNamespace(context.TODO(), &CreateNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: []string{spaceName},
	})
	assert.Nil(t, err)

	tableName := tableNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateTable(context.TODO(), &CreateTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Format:    oss.Ptr("ICEBERG"),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	putResult, err := client.PutTableMaintenanceConfiguration(context.TODO(), &PutTableMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Type:      oss.Ptr("icebergCompaction"),
		Value: &TableMaintenanceValue{
			Status: oss.Ptr("enabled"),
			Settings: &TableMaintenanceSettings{
				IcebergCompaction: &IcebergCompactionSettingsDetail{
					TargetFileSizeMB: oss.Ptr(400),
					Strategy:         oss.Ptr("auto"),
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	putResult, err = client.PutTableMaintenanceConfiguration(context.TODO(), &PutTableMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Type:      oss.Ptr("icebergSnapshotManagement"),
		Value: &TableMaintenanceValue{
			Status: oss.Ptr("enabled"),
			Settings: &TableMaintenanceSettings{
				IcebergSnapshotManagement: &IcebergSnapshotManagementSettingsDetail{
					MaxSnapshotAgeHours: oss.Ptr(350),
					MinSnapshotsToKeep:  oss.Ptr(1),
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetTableMaintenanceConfiguration(context.TODO(), &GetTableMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()
	_, err = invalidAkClient.PutTableMaintenanceConfiguration(context.TODO(), &PutTableMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Type:      oss.Ptr("icebergSnapshotManagement"),
		Value: &TableMaintenanceValue{
			Status: oss.Ptr("enabled"),
			Settings: &TableMaintenanceSettings{
				IcebergSnapshotManagement: &IcebergSnapshotManagementSettingsDetail{
					MaxSnapshotAgeHours: oss.Ptr(350),
					MinSnapshotsToKeep:  oss.Ptr(1),
				},
			},
		},
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.GetTableMaintenanceConfiguration(context.TODO(), &GetTableMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestTableMaintenanceJobStatus(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	result, err := client.CreateTableBucket(context.TODO(), &CreateTableBucketRequest{
		Bucket: oss.Ptr(bucketName),
	})
	assert.Nil(t, err)
	bucketArn := result.Arn

	spaceName := spaceNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateNamespace(context.TODO(), &CreateNamespaceRequest{
		BucketArn: bucketArn,
		Namespace: []string{spaceName},
	})
	assert.Nil(t, err)

	tableName := tableNamePrefix + "_" + randLowStr(5)
	_, err = client.CreateTable(context.TODO(), &CreateTableRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Format:    oss.Ptr("ICEBERG"),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	_, err = client.PutTableMaintenanceConfiguration(context.TODO(), &PutTableMaintenanceConfigurationRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
		Type:      oss.Ptr("icebergCompaction"),
		Value: &TableMaintenanceValue{
			Status: oss.Ptr("enabled"),
			Settings: &TableMaintenanceSettings{
				IcebergCompaction: &IcebergCompactionSettingsDetail{
					TargetFileSizeMB: oss.Ptr(400),
					Strategy:         oss.Ptr("auto"),
				},
			},
		},
	})
	assert.Nil(t, err)

	getResult, err := client.GetTableMaintenanceJobStatus(context.TODO(), &GetTableMaintenanceJobStatusRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// test server error
	invalidAkClient := getInvalidAkClient()
	_, err = invalidAkClient.GetTableMaintenanceJobStatus(context.TODO(), &GetTableMaintenanceJobStatusRequest{
		BucketArn: bucketArn,
		Namespace: oss.Ptr(spaceName),
		Name:      oss.Ptr(tableName),
	})
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}
