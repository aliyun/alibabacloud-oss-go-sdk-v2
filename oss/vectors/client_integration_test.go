package vectors

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
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
	endpoint_   = os.Getenv("OSS_TEST_ENDPOINT")
	accessID_   = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	accessKey_  = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")
	accountUid_ = os.Getenv("OSS_TEST_ACCOUNT_ID")

	instance_ *VectorsClient
	testOnce_ sync.Once
)

var (
	bucketNamePrefix = "go-sdk-test-bucket-"
	indexNamePrefix  = "goSdkIndex"
	letters          = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func getDefaultClient() *VectorsClient {
	testOnce_.Do(func() {
		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
			WithRegion(region_).
			WithEndpoint(endpoint_).
			WithAccountId(accountUid_)

		instance_ = NewVectorsClient(cfg)
	})
	return instance_
}

func getInvalidAkClient() *VectorsClient {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("invalid-ak", "invalid-sk")).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithAccountId(accountUid_)

	return NewVectorsClient(cfg)
}

func getClient(region, endpoint string) *VectorsClient {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithAccountId(accountUid_)

	return NewVectorsClient(cfg)
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

func cleanBucket(bucketInfo VectorBucketProperties, t *testing.T) {
	assert.NotEmpty(t, *bucketInfo.Name)
	var c *VectorsClient
	if (strings.Contains(endpoint_, *bucketInfo.ExtranetEndpoint) ||
		strings.Contains(endpoint_, *bucketInfo.IntranetEndpoint)) || strings.Contains(endpoint_, "drill-") {
		c = getDefaultClient()
	} else {
		c = getClient(*bucketInfo.Region, *bucketInfo.ExtranetEndpoint)
	}
	assert.NotNil(t, c)
	cleanIndexes(c, *bucketInfo.Name, t)
}

func deleteBucket(bucketName string, t *testing.T) {
	assert.NotEmpty(t, bucketName)
	var c *VectorsClient
	c = getDefaultClient()
	assert.NotNil(t, c)
	cleanIndexes(c, bucketName, t)
}

func cleanBuckets(prefix string, t *testing.T) {
	c := getDefaultClient()
	for {
		request := &ListVectorBucketsRequest{
			Prefix: oss.Ptr(prefix),
		}
		result, err := c.ListVectorBuckets(context.TODO(), request)
		assert.Nil(t, err)
		if len(result.Buckets) == 0 {
			return
		}
		for _, b := range result.Buckets {
			cleanBucket(b, t)
		}
	}
}

func cleanIndexes(c *VectorsClient, name string, t *testing.T) {
	var err error
	var bucketName string

	if strings.HasPrefix(name, "acs:ossvector") {
		lastIndex := strings.LastIndex(name, ":")
		if lastIndex != -1 {
			bucketName = name[lastIndex+1:]
		}
	} else {
		bucketName = name
	}

	var listIndexesRequest *ListVectorIndexesRequest
	listIndexesRequest = &ListVectorIndexesRequest{
		Bucket: oss.Ptr(bucketName),
	}
	pagIndexes := c.NewListVectorIndexesPaginator(listIndexesRequest)
	var i int
	for pagIndexes.HasNext() {
		i++
		page, err := pagIndexes.NextPage(context.TODO())
		dumpErrIfNotNil(err)
		assert.Nil(t, err)
		for _, index := range page.Indexes {
			_, err = c.DeleteVectorIndex(context.TODO(), &DeleteVectorIndexRequest{
				Bucket:    oss.Ptr(bucketName),
				IndexName: index.IndexName,
			})
			assert.Nil(t, err)
		}
	}

	delRequest := &DeleteVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}
	_, err = c.DeleteVectorBucket(context.TODO(), delRequest)
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

func TestInvokeOperation(t *testing.T) {
	after := before(t)
	defer after(t)
	BucketName := bucketNamePrefix + randLowStr(5)
	//TODO
	input := &oss.OperationInput{
		OpName: "PutVectorBucket",
		Bucket: oss.Ptr(BucketName),
		Headers: map[string]string{
			oss.HTTPHeaderContentType: "application/json",
		},
		Method: "PUT",
	}

	client := getDefaultClient()
	_, err := client.InvokeOperation(context.TODO(), input)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)

	_, err = client.InvokeOperation(context.TODO(), nil)
	assert.NotNil(t, err)
}

func TestInvokeOperation_BucketPolicy(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	putRequest := &PutVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutVectorBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	//TODO
	calcMd5 := func(input string) string {
		if len(input) == 0 {
			return "1B2M2Y8AsgTpgAmY7PhCfg=="
		}
		h := md5.New()
		h.Write([]byte(input))
		return base64.StdEncoding.EncodeToString(h.Sum(nil))
	}

	// PutBucketPolicy
	policy := `{"Version":"1","Statement":[{"Action":["ossvector:PutVectors","ossvector:GetVectors"],"Effect":"Deny","Principal":["` + accountUid_ + `"],"Resource":["acs:ossvector:` + region_ + `:` + accountUid_ + `:*"]}]}`
	input := &oss.OperationInput{
		OpName: "PutBucketPolicy",
		Method: "PUT",
		Parameters: map[string]string{
			"policy": "",
		},
		// Add Content-md5
		Headers: map[string]string{
			"Content-MD5": calcMd5(policy),
		},
		Body:   strings.NewReader(policy),
		Bucket: oss.Ptr(bucketName),
	}
	output, err := client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)

	// GetBucketPolicy
	input = &oss.OperationInput{
		OpName: "GetBucketPolicy",
		Method: "GET",
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: oss.Ptr(bucketName),
	}
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
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: oss.Ptr(bucketName),
	}
	output, err = client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)
	// discard body
	_, err = io.ReadAll(output.Body)
	assert.NoError(t, err)
	if output.Body != nil {
		output.Body.Close()
	}
}

func TestVectorsBucket(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()
	invalidAkClient := getInvalidAkClient()
	// PutVectorBucket
	putRequest := &PutVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}
	_, err := client.PutVectorBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	// GetVectorBucket
	getRequest := &GetVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}
	_, err = client.GetVectorBucket(context.TODO(), getRequest)
	assert.Nil(t, err)

	// ListVectorBuckets
	listRequest := &ListVectorBucketsRequest{
		Prefix: oss.Ptr(bucketNamePrefix),
	}
	listResult, err := client.ListVectorBuckets(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.True(t, len(listResult.Buckets) > 0)

	// DeleteVectorBucket
	delRequest := &DeleteVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}
	_, err = client.DeleteVectorBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.True(t, len(listResult.Buckets) > 0)

	// test server error
	bucketNameNotExist := bucketNamePrefix + "not-exist"

	putRequest = &PutVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}
	_, err = invalidAkClient.PutVectorBucket(context.TODO(), putRequest)
	assert.NotNil(t, err)
	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	getRequest = &GetVectorBucketRequest{
		Bucket: oss.Ptr(bucketNameNotExist),
	}
	_, err = client.GetVectorBucket(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	listRequest = &ListVectorBucketsRequest{
		Prefix: oss.Ptr(bucketNamePrefix),
	}
	_, err = invalidAkClient.ListVectorBuckets(context.TODO(), listRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteVectorBucketRequest{
		Bucket: oss.Ptr(bucketNameNotExist),
	}
	_, err = client.DeleteVectorBucket(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

//func TestBucketLogging(t *testing.T) {
//	after := before(t)
//	defer after(t)
//	//TODO
//	bucketName := bucketNamePrefix + randLowStr(6)
//	putRequest := &PutVectorBucketRequest{
//		Bucket: oss.Ptr(bucketName),
//	}
//	client := getDefaultClient()
//	_, err := client.PutVectorBucket(context.TODO(), putRequest)
//	assert.Nil(t, err)
//
//	targetBucketName := bucketNamePrefix + randLowStr(6)
//	putOssRequest := &oss.PutBucketRequest{
//		Bucket: oss.Ptr(targetBucketName),
//	}
//	_, err = client.clientImpl.PutBucket(context.TODO(), putOssRequest)
//	assert.Nil(t, err)
//
//	request := &PutBucketLoggingRequest{
//		Bucket: oss.Ptr(bucketName),
//		BucketLoggingStatus: &BucketLoggingStatus{
//			&LoggingEnabled{
//				TargetBucket: oss.Ptr(bucketName),
//				TargetPrefix: oss.Ptr("TargetPrefix"),
//			},
//		},
//	}
//	result, err := client.PutBucketLogging(context.TODO(), request)
//	dumpErrIfNotNil(err)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, result.StatusCode)
//	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
//	time.Sleep(1 * time.Second)
//
//	getRequest := &GetBucketLoggingRequest{
//		Bucket: oss.Ptr(bucketName),
//	}
//	getResult, err := client.GetBucketLogging(context.TODO(), getRequest)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, getResult.StatusCode)
//	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
//	assert.Equal(t, *getResult.BucketLoggingStatus.LoggingEnabled.TargetBucket, bucketName)
//	assert.Equal(t, *getResult.BucketLoggingStatus.LoggingEnabled.TargetPrefix, "TargetPrefix")
//	time.Sleep(1 * time.Second)
//
//	delRequest := &DeleteBucketLoggingRequest{
//		Bucket: oss.Ptr(bucketName),
//	}
//	delResult, err := client.DeleteBucketLogging(context.TODO(), delRequest)
//	assert.Nil(t, err)
//	assert.Equal(t, 204, delResult.StatusCode)
//	assert.Equal(t, "204 No Content", delResult.Status)
//	assert.NotEmpty(t, delResult.Headers.Get("x-oss-request-id"))
//	assert.NotEmpty(t, delResult.Headers.Get("Date"))
//	time.Sleep(1 * time.Second)
//
//	var serr *oss.ServiceError
//	bucketNameNotExist := bucketName + "-not-exist"
//	request = &PutBucketLoggingRequest{
//		Bucket: oss.Ptr(bucketNameNotExist),
//		BucketLoggingStatus: &BucketLoggingStatus{
//			&LoggingEnabled{
//				TargetBucket: oss.Ptr("TargetBucket"),
//				TargetPrefix: oss.Ptr("TargetPrefix"),
//			},
//		},
//	}
//	result, err = client.PutBucketLogging(context.TODO(), request)
//	assert.NotNil(t, err)
//	errors.As(err, &serr)
//	assert.Equal(t, int(404), serr.StatusCode)
//	assert.Equal(t, "NoSuchBucket", serr.Code)
//	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
//	assert.NotEmpty(t, serr.RequestID)
//
//	getRequest = &GetBucketLoggingRequest{
//		Bucket: oss.Ptr(bucketNameNotExist),
//	}
//	serr = &oss.ServiceError{}
//	getResult, err = client.GetBucketLogging(context.TODO(), getRequest)
//	assert.NotNil(t, err)
//	errors.As(err, &serr)
//	assert.Equal(t, int(404), serr.StatusCode)
//	assert.Equal(t, "NoSuchBucket", serr.Code)
//	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
//	assert.NotEmpty(t, serr.RequestID)
//
//	delRequest = &DeleteBucketLoggingRequest{
//		Bucket: oss.Ptr(bucketNameNotExist),
//	}
//	serr = &oss.ServiceError{}
//	delResult, err = client.DeleteBucketLogging(context.TODO(), delRequest)
//	assert.NotNil(t, err)
//	errors.As(err, &serr)
//	assert.Equal(t, int(404), serr.StatusCode)
//	assert.Equal(t, "NoSuchBucket", serr.Code)
//	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
//	assert.NotEmpty(t, serr.RequestID)
//}

func TestBucketPolicy(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutVectorBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketPolicyRequest{
		Bucket: oss.Ptr(bucketName),
		Body: strings.NewReader(`{
  "Version":"1",
  "Statement":[
  {
    "Action":[
      "ossvector:PutVectors",
      "ossvector:GetVectors"
   ],
   "Effect":"Deny",
   "Principal":["` + accountUid_ + `"],
   "Resource":["acs:ossvector:` + region_ + `:` + accountUid_ + `:*"]
  }
 ]
}`),
	}

	putResult, err := client.PutBucketPolicy(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketPolicyRequest{
		Bucket: oss.Ptr(bucketName),
	}
	getResult, err := client.GetBucketPolicy(context.TODO(), getRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, getResult.Body)

	delRequest := &DeleteBucketPolicyRequest{
		Bucket: oss.Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketPolicy(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *oss.ServiceError
	bucketNameNotExist := bucketNamePrefix + "-not-exist"
	getRequest = &GetBucketPolicyRequest{
		Bucket: oss.Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketPolicy(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketPolicyRequest{
		Bucket: oss.Ptr(bucketNameNotExist),
		Body: strings.NewReader(`{
  "Version":"1",
  "Statement":[
  {
    "Action":[
      "oss:PutObject",
      "oss:GetObject"
   ],
   "Effect":"Deny",
   "Principal":["1234567890"],
   "Resource":["acs:oss:*:1234567890:*/*"]
  }
 ]
}`),
	}
	serr = &oss.ServiceError{}
	putResult, err = client.PutBucketPolicy(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest = &DeleteBucketPolicyRequest{
		Bucket: oss.Ptr(bucketNameNotExist),
	}
	serr = &oss.ServiceError{}
	delResult, err = client.DeleteBucketPolicy(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestIndex(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutVectorBucket(context.TODO(), request)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)

	putRequest := &PutVectorIndexRequest{
		Bucket:         oss.Ptr(bucketName),
		IndexName:      oss.Ptr("index"),
		DataType:       oss.Ptr("float32"),
		DistanceMetric: oss.Ptr("cosine"),
		Dimension:      oss.Ptr(int(10)),
	}
	putResult, err := client.PutVectorIndex(context.TODO(), putRequest)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetVectorIndexRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr("index"),
	}
	getResult, err := client.GetVectorIndex(context.TODO(), getRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, getResult.Index)

	listRequest := &ListVectorIndexesRequest{
		Bucket: oss.Ptr(bucketName),
	}
	listResult, err := client.ListVectorIndexes(context.TODO(), listRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, listResult.Indexes)

	delRequest := &DeleteVectorIndexRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr("index"),
	}
	delResult, err := client.DeleteVectorIndex(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *oss.ServiceError
	bucketNameNotExist := bucketNamePrefix + "-not-exist"
	getRequest = &GetVectorIndexRequest{
		Bucket:    oss.Ptr(bucketNameNotExist),
		IndexName: oss.Ptr("index"),
	}
	getResult, err = client.GetVectorIndex(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutVectorIndexRequest{
		Bucket:         oss.Ptr(bucketNameNotExist),
		IndexName:      oss.Ptr("index"),
		DataType:       oss.Ptr("demo"),
		DistanceMetric: oss.Ptr("oss"),
		Dimension:      oss.Ptr(int(10)),
	}
	serr = &oss.ServiceError{}
	putResult, err = client.PutVectorIndex(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	listRequest = &ListVectorIndexesRequest{
		Bucket: oss.Ptr(bucketNameNotExist),
	}
	serr = &oss.ServiceError{}
	listResult, err = client.ListVectorIndexes(context.TODO(), listRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest = &DeleteVectorIndexRequest{
		Bucket:    oss.Ptr(bucketNameNotExist),
		IndexName: oss.Ptr("index"),
	}
	serr = &oss.ServiceError{}
	delResult, err = client.DeleteVectorIndex(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestVectors(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	indexName := indexNamePrefix + randLowStr(5)
	request := &PutVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutVectorBucket(context.TODO(), request)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)

	putRequest := &PutVectorIndexRequest{
		Bucket:         oss.Ptr(bucketName),
		IndexName:      oss.Ptr(indexName),
		DataType:       oss.Ptr("float32"),
		DistanceMetric: oss.Ptr("cosine"),
		Dimension:      oss.Ptr(int(3)),
	}
	_, err = client.PutVectorIndex(context.TODO(), putRequest)
	assert.Nil(t, err)

	putVectorsRequest := &PutVectorsRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
		Vectors: []map[string]any{
			{
				"key": "vector1",
				"data": map[string]any{
					"float32": []float32{1.2, 2.5, 3},
				},
				"metadata": map[string]any{
					"Key2": "value2",
					"Key3": []string{"1", "2", "3"},
				},
			},
		},
	}
	putResult, err := client.PutVectors(context.TODO(), putVectorsRequest)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetVectorsRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
		Keys:      []string{"key1", "key2", "key3"},
	}
	getResult, err := client.GetVectors(context.TODO(), getRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))

	listRequest := &ListVectorsRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
	}
	listResult, err := client.ListVectors(context.TODO(), listRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, listResult.Vectors)

	queryRequest := &QueryVectorsRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
		Filter: map[string]any{
			"$and": []map[string]any{
				{
					"type": map[string]any{
						"$in": []string{"comedy", "documentary"},
					},
				},
			},
		},
		QueryVector: map[string]any{
			"float32": []float32{1, 2, 3},
		},
		ReturnMetadata: oss.Ptr(true),
		ReturnDistance: oss.Ptr(true),
		TopK:           oss.Ptr(10),
	}
	queryResult, err := client.QueryVectors(context.TODO(), queryRequest)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, queryResult.StatusCode)
	assert.NotEmpty(t, queryResult.Headers.Get("X-Oss-Request-Id"))

	delRequest := &DeleteVectorsRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
		Keys:      []string{"key1", "key2", "key3"},
	}
	delResult, err := client.DeleteVectors(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *oss.ServiceError
	bucketNameNotExist := bucketNamePrefix + "-not-exist"
	getRequest = &GetVectorsRequest{
		Bucket:    oss.Ptr(bucketNameNotExist),
		IndexName: oss.Ptr(indexName),
		Keys:      []string{"key1", "key2", "key3"},
	}
	getResult, err = client.GetVectors(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putVectorsRequest = &PutVectorsRequest{
		Bucket:    oss.Ptr(bucketNameNotExist),
		IndexName: oss.Ptr(indexName),
		Vectors: []map[string]any{
			{
				"key": "vector1",
				"data": map[string]any{
					"float32": []float32{1.2, 2.5, 3},
				},
				"metadata": map[string]any{
					"Key1": 32,
					"Key2": "value2",
					"Key3": []string{"1", "2", "3"},
					"Key4": false,
				},
			},
		},
	}
	serr = &oss.ServiceError{}
	putResult, err = client.PutVectors(context.TODO(), putVectorsRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	listRequest = &ListVectorsRequest{
		Bucket:    oss.Ptr(bucketNameNotExist),
		IndexName: oss.Ptr(indexName),
	}
	serr = &oss.ServiceError{}
	listResult, err = client.ListVectors(context.TODO(), listRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	queryRequest = &QueryVectorsRequest{
		Bucket:    oss.Ptr(bucketNameNotExist),
		IndexName: oss.Ptr(indexName),
		Filter: map[string]any{
			"$and": []map[string]any{
				{
					"type": map[string]any{
						"$in": []string{"comedy", "documentary"},
					},
				},
			},
		},
		QueryVector: map[string]any{
			"float32": []float32{1, 2, 3},
		},
		ReturnMetadata: oss.Ptr(true),
		ReturnDistance: oss.Ptr(true),
		TopK:           oss.Ptr(10),
	}
	queryResult, err = client.QueryVectors(context.TODO(), queryRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteVectorsRequest{
		Bucket:    oss.Ptr(bucketNameNotExist),
		IndexName: oss.Ptr(indexName),
		Keys:      []string{"key1", "key2", "key3"},
	}
	serr = &oss.ServiceError{}
	delResult, err = client.DeleteVectors(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	del := &DeleteVectorIndexRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
	}
	_, err = client.DeleteVectorIndex(context.TODO(), del)
	assert.Nil(t, err)
}

func TestPaginator(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutVectorBucketRequest{
		Bucket: oss.Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutVectorBucket(context.TODO(), request)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)

	indexName := indexNamePrefix + randStr(5)

	putRequest := &PutVectorIndexRequest{
		Bucket:         oss.Ptr(bucketName),
		IndexName:      oss.Ptr(indexName),
		DataType:       oss.Ptr("float32"),
		DistanceMetric: oss.Ptr("cosine"),
		Dimension:      oss.Ptr(int(3)),
	}
	_, err = client.PutVectorIndex(context.TODO(), putRequest)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)

	for i := 1; i <= 10; i++ {
		putRequest := &PutVectorIndexRequest{
			Bucket:         oss.Ptr(bucketName),
			IndexName:      oss.Ptr(indexName + strconv.Itoa(i)),
			DataType:       oss.Ptr("float32"),
			DistanceMetric: oss.Ptr("cosine"),
			Dimension:      oss.Ptr(int(3)),
		}
		_, err = client.PutVectorIndex(context.TODO(), putRequest)
		assert.Nil(t, err)
	}

	for i := 1; i <= 10; i++ {
		putVectorsRequest := &PutVectorsRequest{
			Bucket:    oss.Ptr(bucketName),
			IndexName: oss.Ptr(indexName),
			Vectors: []map[string]any{
				{
					"key": "vector" + strconv.Itoa(i),
					"data": map[string]any{
						"float32": []float32{1.2, 2.5, 3},
					},
					"metadata": map[string]any{
						"Key2": "value2",
						"Key3": []string{"1", "2", "3"},
					},
				},
			},
		}
		putResult, err := client.PutVectors(context.TODO(), putVectorsRequest)
		assert.Nil(t, err)
		assert.Equal(t, 200, putResult.StatusCode)
		assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	}

	bucketRequest := &ListVectorBucketsRequest{}
	p := client.NewListVectorBucketsPaginator(bucketRequest)
	for p.HasNext() {
		page, err := p.NextPage(context.TODO())
		assert.Nil(t, err)
		assert.True(t, len(page.Buckets) > 0)
	}

	indexRequest := &ListVectorIndexesRequest{
		Bucket:     oss.Ptr(bucketName),
		MaxResults: 5,
	}
	pIndex := client.NewListVectorIndexesPaginator(indexRequest)
	var m int
	for pIndex.HasNext() {
		m++
		pageIndex, err := pIndex.NextPage(context.TODO())
		assert.Nil(t, err)
		assert.True(t, len(pageIndex.Indexes) > 0)
	}
	assert.Equal(t, m, 3)
	vectorsRequest := &ListVectorsRequest{
		Bucket:     oss.Ptr(bucketName),
		IndexName:  oss.Ptr(indexName),
		MaxResults: 5,
	}
	pVectors := client.NewListVectorsPaginator(vectorsRequest)
	var n int
	for pVectors.HasNext() {
		n++
		pageVectors, err := pVectors.NextPage(context.TODO())
		assert.Nil(t, err)
		assert.True(t, len(pageVectors.Vectors) > 0)
	}
	assert.Equal(t, n, 2)
	for i := 1; i <= 10; i++ {
		delRequest := &DeleteVectorsRequest{
			Bucket:    oss.Ptr(bucketName),
			IndexName: oss.Ptr(indexName + strconv.Itoa(i)),
			Keys:      []string{"key2", "key3"},
		}
		_, err = client.DeleteVectors(context.TODO(), delRequest)
		assert.Nil(t, err)
	}

	for i := 1; i <= 10; i++ {
		del := &DeleteVectorIndexRequest{
			Bucket:    oss.Ptr(bucketName),
			IndexName: oss.Ptr(indexName + strconv.Itoa(i)),
		}
		_, err = client.DeleteVectorIndex(context.TODO(), del)
		assert.Nil(t, err)
	}
	del := &DeleteVectorIndexRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
	}
	_, err = client.DeleteVectorIndex(context.TODO(), del)
	assert.Nil(t, err)
}
