package oss

import (
	"bufio"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/crypto"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

var (
	// Endpoint/ID/Key
	region_           = os.Getenv("OSS_TEST_REGION")
	endpoint_         = os.Getenv("OSS_TEST_ENDPOINT")
	accessID_         = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	accessKey_        = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")
	ramRoleArn_       = os.Getenv("OSS_TEST_RAM_ROLE_ARN")
	accountID_        = os.Getenv("OSS_TEST_RAM_UID")
	signatureVersion_ = os.Getenv("OSS_TEST_SIGNATURE_VERSION")

	// payer
	payerAccessID_  = os.Getenv("OSS_TEST_PAYER_ACCESS_KEY_ID")
	payerAccessKey_ = os.Getenv("OSS_TEST_PAYER_ACCESS_KEY_SECRET")
	payerUID_       = os.Getenv("OSS_TEST_PAYER_UID")

	// path style
	pathStyleBucket_ = os.Getenv("OSS_TEST_PATHSTYLE_BUCKET")
	pathStyleRegion_ = os.Getenv("OSS_TEST_PATHSTYLE_REGION")

	apEnable = os.Getenv("OSS_TEST_AP_ENABLE")

	instance_ *Client
	testOnce_ sync.Once

	kmdIdMap_ = map[string]string{}
)

var (
	bucketNamePrefix = "go-sdk-test-bucket-"
	objectNamePrefix = "go-sdk-test-object-"
	letters          = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func getDefaultClient() *Client {
	testOnce_.Do(func() {
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
			WithRegion(region_).
			WithEndpoint(endpoint_).
			WithSignatureVersion(getSignatrueVersion())

		instance_ = NewClient(cfg)
	})
	return instance_
}

func getClient(region, endpoint string) *Client {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithSignatureVersion(getSignatrueVersion())

	return NewClient(cfg)
}

func getClientUseStsToken(region, endpoint string) *Client {
	resp, err := stsAssumeRole(accessID_, accessKey_, ramRoleArn_)
	if err != nil {
		return nil
	}
	accessId := resp.Credentials.AccessKeyId
	accessKey := resp.Credentials.AccessKeySecret
	token := resp.Credentials.SecurityToken
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessId, accessKey, token)).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithSignatureVersion(getSignatrueVersion())

	return NewClient(cfg)
}

func getClientUseStsTokenV2(cfg *Config) *Client {
	resp, err := stsAssumeRole(accessID_, accessKey_, ramRoleArn_)
	if err != nil {
		return nil
	}
	accessId := resp.Credentials.AccessKeyId
	accessKey := resp.Credentials.AccessKeySecret
	token := resp.Credentials.SecurityToken
	cfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessId, accessKey, token))
	/*
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessId, accessKey, token)).
			WithRegion(region).
			WithEndpoint(endpoint).
			WithSignatureVersion(getSignatrueVersion())
	*/
	return NewClient(cfg)
}

func getClientWithCredentialsProvider(region, endpoint string, cred credentials.CredentialsProvider) *Client {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(cred).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithSignatureVersion(getSignatrueVersion())

	return NewClient(cfg)
}

func getKmsID(region string) string {
	if id, ok := kmdIdMap_[region]; ok {
		return id
	}

	client := getClient(region, fmt.Sprintf("oss-%s.aliyuncs.com", region))
	bucketName := bucketNamePrefix + randLowStr(6)

	if _, err := client.PutBucket(context.TODO(), &PutBucketRequest{Bucket: Ptr(bucketName)}); err != nil {
		return ""
	}

	kmdId := ""
	if _, err := client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket:               Ptr(bucketName),
		Key:                  Ptr("kms-id"),
		ServerSideEncryption: Ptr("KMS")}); err == nil {

		if result, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr("kms-id")}); err == nil {
			kmdId = ToString(result.ServerSideEncryptionKeyId)
			kmdIdMap_[region] = kmdId
		}
	}
	client.DeleteObject(context.TODO(), &DeleteObjectRequest{Bucket: Ptr(bucketName), Key: Ptr("kms-id")})
	client.DeleteBucket(context.TODO(), &DeleteBucketRequest{Bucket: Ptr(bucketName)})
	return kmdId
}

func getSignatrueVersion() SignatureVersionType {
	switch signatureVersion_ {
	case "v1":
		return SignatureVersionV1
	default:
		return SignatureVersionV4
	}
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

func cleanBucket(bucketInfo BucketProperties, t *testing.T) {
	assert.NotEmpty(t, *bucketInfo.Name)
	var c *Client
	if strings.Contains(endpoint_, *bucketInfo.ExtranetEndpoint) ||
		strings.Contains(endpoint_, *bucketInfo.IntranetEndpoint) {
		c = getDefaultClient()
	} else {
		c = getClient(*bucketInfo.Region, *bucketInfo.ExtranetEndpoint)
	}
	assert.NotNil(t, c)
	cleanObjects(c, *bucketInfo.Name, t)
}

func deleteBucket(bucketName string, t *testing.T) {
	assert.NotEmpty(t, bucketName)
	var c *Client
	c = getDefaultClient()
	assert.NotNil(t, c)
	cleanObjects(c, bucketName, t)
}

func cleanBuckets(prefix string, t *testing.T) {
	c := getDefaultClient()
	for {
		request := &ListBucketsRequest{
			Prefix: Ptr(prefix),
		}
		result, err := c.ListBuckets(context.TODO(), request)
		assert.Nil(t, err)
		if len(result.Buckets) == 0 {
			return
		}
		for _, b := range result.Buckets {
			cleanBucket(b, t)
		}
	}
}

func cleanObjects(c *Client, bucketName string, t *testing.T) {
	var err error
	var listRequest *ListObjectsRequest
	var delObjRequest *DeleteObjectRequest
	var lor *ListObjectsResult
	marker := ""
	for {
		listRequest = &ListObjectsRequest{
			Bucket: Ptr(bucketName),
			Marker: Ptr(marker),
		}
		lor, err = c.ListObjects(context.TODO(), listRequest)
		assert.Nil(t, err)
		var deleteObjects []DeleteObject
		for _, object := range lor.Contents {
			deleteObjects = append(deleteObjects, DeleteObject{Key: object.Key})
		}
		if len(deleteObjects) > 0 {
			_, err = c.DeleteMultipleObjects(context.TODO(), &DeleteMultipleObjectsRequest{
				Bucket:  Ptr(bucketName),
				Objects: deleteObjects,
			})
			assert.Nil(t, err)
		}

		if !lor.IsTruncated {
			break
		}
		if lor.NextMarker != nil {
			marker = *lor.NextMarker
		}
	}
	var listUploadRequest *ListMultipartUploadsRequest
	var abortRequest *AbortMultipartUploadRequest
	var lsRes *ListMultipartUploadsResult
	keyMarker := ""
	uploadIdMarker := ""
	for {
		listUploadRequest = &ListMultipartUploadsRequest{
			Bucket:         Ptr(bucketName),
			KeyMarker:      Ptr(keyMarker),
			UploadIdMarker: Ptr(uploadIdMarker),
		}
		lsRes, err = c.ListMultipartUploads(context.TODO(), listUploadRequest)
		assert.Nil(t, err)
		for _, upload := range lsRes.Uploads {
			abortRequest = &AbortMultipartUploadRequest{
				Bucket:   Ptr(bucketName),
				Key:      Ptr(*upload.Key),
				UploadId: Ptr(*upload.UploadId),
			}
			_, err = c.AbortMultipartUpload(context.TODO(), abortRequest)
			assert.Nil(t, err)
		}
		if !lsRes.IsTruncated {
			break
		}
		keyMarker = *lsRes.NextKeyMarker
		uploadIdMarker = *lsRes.NextUploadIdMarker
	}
	var lsVersionRq *ListObjectVersionsRequest
	var lsVersionRs *ListObjectVersionsResult
	versionKeyMarker := ""
	VersionIdMarker := ""
	for {
		lsVersionRq = &ListObjectVersionsRequest{
			Bucket:          Ptr(bucketName),
			KeyMarker:       Ptr(versionKeyMarker),
			VersionIdMarker: Ptr(VersionIdMarker),
		}
		lsVersionRs, err = c.ListObjectVersions(context.TODO(), lsVersionRq)
		assert.Nil(t, err)
		for _, object := range lsVersionRs.ObjectDeleteMarkers {
			delObjRequest = &DeleteObjectRequest{
				Bucket:    Ptr(bucketName),
				Key:       Ptr(*object.Key),
				VersionId: Ptr(*object.VersionId),
			}
			_, err = c.DeleteObject(context.TODO(), delObjRequest)
			assert.Nil(t, err)
		}
		for _, object := range lsVersionRs.ObjectVersions {
			delObjRequest = &DeleteObjectRequest{
				Bucket:    Ptr(bucketName),
				Key:       Ptr(*object.Key),
				VersionId: Ptr(*object.VersionId),
			}
			_, err = c.DeleteObject(context.TODO(), delObjRequest)
			assert.Nil(t, err)
		}
		if !lsVersionRs.IsTruncated {
			break
		}
		versionKeyMarker = *lsVersionRs.NextKeyMarker
		VersionIdMarker = *lsVersionRs.NextVersionIdMarker
	}
	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = c.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
}

func deleteAccessPoint(c *Client, bucketName string, t *testing.T) {
	var err error
	var listRequest *ListAccessPointsRequest
	var lap *ListAccessPointsResult
	var delRequest *DeleteAccessPointRequest
	var token = ""
	var getResult *GetAccessPointResult
	var serr = &ServiceError{}
	for {
		listRequest = &ListAccessPointsRequest{
			Bucket:            Ptr(bucketName),
			ContinuationToken: Ptr(token),
		}
		lap, err = c.ListAccessPoints(context.TODO(), listRequest)
		assert.Nil(t, err)
		if len(lap.AccessPoints) > 0 {
			for _, accessPoint := range lap.AccessPoints {
				switch *accessPoint.Status {
				case "creating":
					time.Sleep(3 * time.Second)
					for {
						getResult, err = c.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
							Bucket:          Ptr(bucketName),
							AccessPointName: accessPoint.AccessPointName,
						})
						assert.Nil(t, err)
						if *getResult.AccessPointStatus != "creating" {
							break
						} else {
							time.Sleep(3 * time.Second)
						}
					}
					delRequest = &DeleteAccessPointRequest{
						Bucket:          Ptr(bucketName),
						AccessPointName: accessPoint.AccessPointName,
					}
					_, err = c.DeleteAccessPoint(context.TODO(), delRequest)
					assert.Nil(t, err)
				case "deleting":
					time.Sleep(5 * time.Second)
				default:
					delRequest = &DeleteAccessPointRequest{
						Bucket:          Ptr(bucketName),
						AccessPointName: accessPoint.AccessPointName,
					}
					_, err = c.DeleteAccessPoint(context.TODO(), delRequest)
					assert.Nil(t, err)
					time.Sleep(3 * time.Second)

				}
				for {
					getResult, err = c.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
						Bucket:          Ptr(bucketName),
						AccessPointName: accessPoint.AccessPointName,
					})
					if err != nil {
						errors.As(err, &serr)
						if serr.StatusCode == 404 && serr.Code == "NoSuchAccessPoint" {
							break
						}
					}
					time.Sleep(3 * time.Second)
				}
			}

			if !*lap.IsTruncated {
				break
			}
			if lap.NextContinuationToken != nil {
				token = *lap.NextContinuationToken
			}
		} else {
			break
		}

	}
}

type credentialsForSts struct {
	AccessKeyId     string
	AccessKeySecret string
	Expiration      time.Time
	SecurityToken   string
}

type assumedRoleUserForSts struct {
	Arn           string
	AssumedRoleId string
}

type responseForSts struct {
	Credentials     credentialsForSts
	AssumedRoleUser assumedRoleUserForSts
	RequestId       string
}

func stsAssumeRole(accessKeyId string, accessKeySecret string, roleArn string) (*responseForSts, error) {
	// StsSignVersion sts sign version
	StsSignVersion := "1.0"
	// StsAPIVersion sts api version
	StsAPIVersion := "2015-04-01"
	// StsHost sts host
	StsHost := "https://sts.aliyuncs.com/"
	// TimeFormat time fomrat
	TimeFormat := "2006-01-02T15:04:05Z"
	// RespBodyFormat  respone body format
	RespBodyFormat := "JSON"
	// PercentEncode '/'
	PercentEncode := "%2F"
	// HTTPGet http get method
	HTTPGet := "GET"
	rand.Seed(time.Now().UnixNano())
	uuid := fmt.Sprintf("Nonce-%d", rand.Intn(10000))
	queryStr := "SignatureVersion=" + StsSignVersion
	queryStr += "&Format=" + RespBodyFormat
	queryStr += "&Timestamp=" + url.QueryEscape(time.Now().UTC().Format(TimeFormat))
	queryStr += "&RoleArn=" + url.QueryEscape(roleArn)
	queryStr += "&RoleSessionName=" + "oss_test_sess"
	queryStr += "&AccessKeyId=" + accessKeyId
	queryStr += "&SignatureMethod=HMAC-SHA1"
	queryStr += "&Version=" + StsAPIVersion
	queryStr += "&Action=AssumeRole"
	queryStr += "&SignatureNonce=" + uuid
	queryStr += "&DurationSeconds=" + strconv.FormatInt(3600, 10)

	// Sort query string
	queryParams, err := url.ParseQuery(queryStr)
	if err != nil {
		return nil, err
	}

	strToSign := HTTPGet + "&" + PercentEncode + "&" + url.QueryEscape(queryParams.Encode())

	// Generate signature
	hashSign := hmac.New(sha1.New, []byte(accessKeySecret+"&"))
	hashSign.Write([]byte(strToSign))
	signature := base64.StdEncoding.EncodeToString(hashSign.Sum(nil))

	// Build url
	assumeURL := StsHost + "?" + queryStr + "&Signature=" + url.QueryEscape(signature)

	// Send Request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(assumeURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	// Handle Response
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	result := responseForSts{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func before(_ *testing.T) func(t *testing.T) {

	//fmt.Println("setup test case")
	return after
}

func after(t *testing.T) {
	cleanBuckets(bucketNamePrefix, t)
	//fmt.Println("teardown  test case")
}

func clearAp(t *testing.T) {
	c := getDefaultClient()
	for {
		request := &ListBucketsRequest{
			Prefix: Ptr(bucketNamePrefix),
		}
		result, err := c.ListBuckets(context.TODO(), request)
		assert.Nil(t, err)
		if len(result.Buckets) == 0 {
			return
		}
		for _, b := range result.Buckets {
			deleteAccessPoint(c, *b.Name, t)
		}
		if !result.IsTruncated {
			break
		}
	}
}

func dumpErrIfNotNil(err error) {
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}
}

func TestInvokeOperation(t *testing.T) {
	after := before(t)
	defer after(t)
	BucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	input := &OperationInput{
		OpName: "PutBucket",
		Bucket: Ptr(BucketName),
		Method: "PUT",
	}
	client := getDefaultClient()
	_, err := client.InvokeOperation(context.TODO(), input)
	assert.Nil(t, err)

	_, err = client.InvokeOperation(context.TODO(), nil)
	assert.NotNil(t, err)
}

func TestInvokeOperation_BucketPolicy(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	calcMd5 := func(input string) string {
		if len(input) == 0 {
			return "1B2M2Y8AsgTpgAmY7PhCfg=="
		}
		h := md5.New()
		h.Write([]byte(input))
		return base64.StdEncoding.EncodeToString(h.Sum(nil))
	}

	// PutBucketPolicy
	policy := `{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Allow","Resource":["acs:oss:*:*:*"]}]}`
	input := &OperationInput{
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
		Bucket: Ptr(bucketName),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	output, err := client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)

	// GetBucketPolicy
	input = &OperationInput{
		OpName: "GetBucketPolicy",
		Method: "GET",
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: Ptr(bucketName),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	output, err = client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)
	policy1, err := io.ReadAll(output.Body)
	assert.NoError(t, err)
	if output.Body != nil {
		output.Body.Close()
	}
	assert.NotEmpty(t, policy1)

	// DeleteBucketPolicy
	input = &OperationInput{
		OpName: "DeleteBucketPolicy",
		Method: "DELETE",
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: Ptr(bucketName),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	output, err = client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)
	// discard body
	_, err = io.ReadAll(output.Body)
	assert.NoError(t, err)
	if output.Body != nil {
		output.Body.Close()
	}
}

func TestListBuckets(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketPrefix := bucketNamePrefix + randLowStr(6)
	client := getDefaultClient()
	//TODO
	var bucketName string
	count := 10
	for i := 0; i < count; i++ {
		bucketName = bucketPrefix + strconv.Itoa(i)
		putRequest := &PutBucketRequest{
			Bucket: Ptr(bucketName),
		}
		_, err := client.PutBucket(context.TODO(), putRequest)
		assert.NoError(t, err)
		assert.Nil(t, err)
	}

	listRequest := &ListBucketsRequest{
		Prefix: Ptr(bucketPrefix),
	}

	result, err := client.ListBuckets(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, len(result.Buckets), count)

	_, err = client.ListBuckets(context.TODO(), nil)
	assert.Nil(t, err)
}

func TestPutBucket(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	result, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.Status, "200 OK")
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id") != "", true)
	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	_, err = client.PutBucket(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	bucketName = bucketNamePrefix + randLowStr(6)
	putRequest = &PutBucketRequest{
		Bucket: Ptr(bucketName),
		CreateBucketConfiguration: &CreateBucketConfiguration{
			StorageClass:       StorageClassStandard,
			DataRedundancyType: DataRedundancyLRS,
		},
	}
	result, err = client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.Status, "200 OK")
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id") != "", true)
	delRequest = &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
}

func TestDeleteBucket(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.Status, "204 No Content")
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id") != "", true)

	_, err = client.DeleteBucket(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	result, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, strings.Contains(serr.Message, "not exist"), true)
	assert.Equal(t, serr.RequestID != "", true)
}

func TestListObjects(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	request := &ListObjectsRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.ListObjects(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.Name, bucketName)
	assert.Equal(t, len(result.Contents), 0)
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Empty(t, result.Prefix)
	assert.Empty(t, result.Marker)
	assert.Empty(t, result.Delimiter)
	assert.Equal(t, result.IsTruncated, false)
	bucketNotExist := bucketNamePrefix + "not-exist" + randLowStr(5)
	request = &ListObjectsRequest{
		Bucket: Ptr(bucketNotExist),
	}
	_, err = client.ListObjects(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	_, err = client.ListObjects(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
}

func TestListObjectsV2(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	request := &ListObjectsV2Request{
		Bucket: Ptr(bucketName),
	}
	result, err := client.ListObjectsV2(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.Name, bucketName)
	assert.Equal(t, len(result.Contents), 0)
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Empty(t, result.Prefix)
	assert.Empty(t, result.StartAfter)
	assert.Empty(t, result.Delimiter)
	assert.Equal(t, result.IsTruncated, false)
	bucketNotExist := bucketNamePrefix + "not-exist" + randLowStr(5)
	request = &ListObjectsV2Request{
		Bucket: Ptr(bucketNotExist),
	}
	_, err = client.ListObjectsV2(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	_, err = client.ListObjectsV2(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
}

func TestGetBucketInfo(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	getRequest := &GetBucketInfoRequest{
		Bucket: Ptr(bucketName),
	}
	info, err := client.GetBucketInfo(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, *info.BucketInfo.Name, bucketName)
	assert.Equal(t, *info.BucketInfo.AccessMonitor, "Disabled")
	assert.NotEmpty(t, *info.BucketInfo.CreationDate)
	assert.True(t, strings.Contains(*info.BucketInfo.ExtranetEndpoint, ".aliyuncs.com"))
	assert.True(t, strings.Contains(*info.BucketInfo.IntranetEndpoint, "internal.aliyuncs.com"))
	assert.True(t, strings.Contains(*info.BucketInfo.Location, "oss-"))
	assert.True(t, strings.Contains(*info.BucketInfo.StorageClass, "Standard"))
	assert.Equal(t, *info.BucketInfo.TransferAcceleration, "Disabled")
	assert.Equal(t, *info.BucketInfo.CrossRegionReplication, "Disabled")
	assert.NotEmpty(t, *info.BucketInfo.ResourceGroupId)
	assert.NotEmpty(t, *info.BucketInfo.Owner.DisplayName)
	assert.NotEmpty(t, *info.BucketInfo.Owner.DisplayName)
	assert.Equal(t, *info.BucketInfo.ACL, "private")
	assert.Empty(t, info.BucketInfo.BucketPolicy.LogBucket)
	assert.Empty(t, info.BucketInfo.BucketPolicy.LogPrefix)

	assert.Equal(t, *info.BucketInfo.SseRule.SSEAlgorithm, "")
	assert.False(t, *info.BucketInfo.BlockPublicAccess)
	assert.Nil(t, info.BucketInfo.SseRule.KMSDataEncryption)
	assert.Nil(t, info.BucketInfo.SseRule.KMSMasterKeyID)
	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
	_, err = client.GetBucketInfo(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	_, err = client.GetBucketInfo(context.TODO(), getRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetBucketLocation(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	getRequest := &GetBucketLocationRequest{
		Bucket: Ptr(bucketName),
	}
	info, err := client.GetBucketLocation(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Contains(t, *info.LocationConstraint, region_)
	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	_, err = client.GetBucketLocation(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	_, err = client.GetBucketLocation(context.TODO(), getRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetBucketStat(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	getRequest := &GetBucketStatRequest{
		Bucket: Ptr(bucketName),
	}
	stat, err := client.GetBucketStat(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), stat.Storage)
	assert.Equal(t, int64(0), stat.ObjectCount)
	assert.Equal(t, int64(0), stat.MultipartUploadCount)
	assert.Equal(t, int64(0), stat.LiveChannelCount)
	assert.Equal(t, int64(0), stat.LastModifiedTime)
	assert.Equal(t, int64(0), stat.StandardStorage)
	assert.Equal(t, int64(0), stat.StandardObjectCount)
	assert.Equal(t, int64(0), stat.InfrequentAccessStorage)
	assert.Equal(t, int64(0), stat.InfrequentAccessRealStorage)
	assert.Equal(t, int64(0), stat.InfrequentAccessObjectCount)
	assert.Equal(t, int64(0), stat.ArchiveStorage)
	assert.Equal(t, int64(0), stat.ArchiveRealStorage)
	assert.Equal(t, int64(0), stat.ArchiveObjectCount)
	assert.Equal(t, int64(0), stat.ColdArchiveStorage)
	assert.Equal(t, int64(0), stat.ColdArchiveRealStorage)
	assert.Equal(t, int64(0), stat.ColdArchiveObjectCount)
	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	_, err = client.GetBucketStat(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	_, err = client.GetBucketStat(context.TODO(), getRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPutBucketAcl(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	request := &PutBucketAclRequest{
		Bucket: Ptr(bucketName),
		Acl:    BucketACLPublicRead,
	}
	result, err := client.PutBucketAcl(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(2 * time.Second)
	infoRequest := &GetBucketInfoRequest{
		Bucket: Ptr(bucketName),
	}

	info, err := client.GetBucketInfo(context.TODO(), infoRequest)
	assert.Nil(t, err)
	assert.Equal(t, string(BucketACLPublicRead), *info.BucketInfo.ACL)
	request = &PutBucketAclRequest{
		Bucket: Ptr(bucketName),
		Acl:    BucketACLPrivate,
	}
	result, err = client.PutBucketAcl(context.TODO(), request)
	assert.Nil(t, err)

	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(2 * time.Second)
	info, err = client.GetBucketInfo(context.TODO(), infoRequest)
	assert.Nil(t, err)
	assert.Equal(t, string(BucketACLPrivate), *info.BucketInfo.ACL)

	_, err = client.PutBucketAcl(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutBucketAclRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.PutBucketAcl(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
}

func TestGetBucketAcl(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	request := &GetBucketAclRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.GetBucketAcl(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, BucketACLType(*result.ACL), BucketACLPrivate)
	assert.NotEmpty(t, *result.Owner.ID)
	assert.NotEmpty(t, *result.Owner.DisplayName)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	result, err = client.GetBucketAcl(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketAclRequest{
		Bucket: Ptr(bucketName),
	}
	result, err = client.GetBucketAcl(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPutObject(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	result, err := client.PutObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)

	request = &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(content),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)

	var serr *ServiceError
	request = &PutObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		Body:     strings.NewReader(content),
		Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 203, serr.StatusCode)
	assert.Equal(t, "CallbackFailed", serr.Code)
	assert.Equal(t, "Error status : 301.", serr.Message)
	assert.Equal(t, "0007-00000203", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	result, err = client.PutObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	request = &PutObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	//Body is bigger than Content-Length
	request = &PutObjectRequest{
		Bucket:        Ptr(bucketName),
		Key:           Ptr(objectName),
		ContentLength: Ptr(int64(len(content) - 10)),
		Body:          strings.NewReader(content),
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " transport connection broken")

}

func TestGetObject(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)
	getRequest := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObject(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.ContentLength, int64(len(content)))

	getRequest = &GetObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	result, err = client.GetObject(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.ContentLength, int64(len(content)))
	_, err = client.GetObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	getRequest = &GetObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	_, err = client.GetObject(context.TODO(), getRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestCopyObject(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	objectCopyName := objectNamePrefix + randLowStr(6) + "copy"
	copyRequest := &CopyObjectRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		SourceKey: Ptr(objectCopyName),
	}
	result, err := client.CopyObject(context.TODO(), copyRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchKey", serr.Code)
	assert.Equal(t, "The specified key does not exist.", serr.Message)

	copyRequest = &CopyObjectRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectCopyName),
		SourceKey: Ptr(objectName),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.Nil(t, result.VersionId)

	copyRequest = &CopyObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectCopyName),
		SourceKey:    Ptr(objectName),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.Nil(t, result.VersionId)

	_, err = client.CopyObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	copyRequest = &CopyObjectRequest{
		Bucket:    Ptr(bucketNameNotExist),
		Key:       Ptr(objectCopyName),
		SourceKey: Ptr(objectName),
	}
	_, err = client.CopyObject(context.TODO(), copyRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	metaRequest := &GetObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	metaResult, err := client.GetObjectMeta(context.TODO(), metaRequest)
	assert.Nil(t, err)
	sourceVersionId := *metaResult.VersionId

	copyRequest = &CopyObjectRequest{
		Bucket:          Ptr(bucketName),
		Key:             Ptr(objectCopyName),
		SourceKey:       Ptr(objectName),
		SourceVersionId: Ptr(sourceVersionId),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, result.VersionId)
	assert.Equal(t, *result.SourceVersionId, sourceVersionId)

	bucketCopyName := bucketNamePrefix + randLowStr(6) + "copy"
	putRequest = &PutBucketRequest{
		Bucket: Ptr(bucketCopyName),
	}
	client = getDefaultClient()
	_, err = client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	copyRequest = &CopyObjectRequest{
		Bucket:       Ptr(bucketCopyName),
		Key:          Ptr(objectCopyName),
		SourceKey:    Ptr(objectName),
		SourceBucket: Ptr(bucketName),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.Nil(t, result.VersionId)
}

func TestAppendObject(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	var result *AppendObjectResult
	content := randLowStr(100)
	request := &AppendObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.AppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AppendObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		Body:     strings.NewReader(content),
		Position: Ptr(int64(0)),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)))
	assert.NotEmpty(t, result.HashCRC64)

	nextPosition := result.NextPosition
	request = &AppendObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(content),
		Position:     Ptr(nextPosition),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)*2))
	assert.NotEmpty(t, result.HashCRC64)

	nextPosition = result.NextPosition
	request = &AppendObjectRequest{
		Bucket:                   Ptr(bucketName),
		Key:                      Ptr(objectName),
		Body:                     strings.NewReader(content),
		Position:                 Ptr(nextPosition),
		ServerSideDataEncryption: Ptr("SM4"),
		ServerSideEncryption:     Ptr("KMS"),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)*3))
	assert.NotEmpty(t, result.HashCRC64)

	objectName2 := objectName + "-kms-sm4"
	request = &AppendObjectRequest{
		Bucket:                   Ptr(bucketName),
		Key:                      Ptr(objectName2),
		Body:                     strings.NewReader(content),
		Position:                 Ptr(int64(0)),
		ServerSideDataEncryption: Ptr("SM4"),
		ServerSideEncryption:     Ptr("KMS"),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.NotEmpty(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)))
	assert.NotEmpty(t, result.HashCRC64)

	nextPosition = result.NextPosition
	request = &AppendObjectRequest{
		Bucket:                   Ptr(bucketName),
		Key:                      Ptr(objectName2),
		Body:                     strings.NewReader(content),
		Position:                 Ptr(nextPosition),
		ServerSideDataEncryption: Ptr("SM4"),
		ServerSideEncryption:     Ptr("KMS"),
		TrafficLimit:             int64(100 * 1024 * 8),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.NotEmpty(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)*2))
	assert.NotEmpty(t, result.HashCRC64)

	_, err = client.AppendObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	var serr *ServiceError
	request = &AppendObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		Body:     strings.NewReader(content),
		Position: Ptr(int64(0)),
	}
	_, err = client.AppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "PositionNotEqualToLength", serr.Code)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketName + "-not-exist"
	request = &AppendObjectRequest{
		Bucket:   Ptr(bucketNameNotExist),
		Key:      Ptr(objectName),
		Body:     strings.NewReader(content),
		Position: Ptr(int64(0)),
	}
	_, err = client.AppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDeleteObject(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	exist, err := client.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	assert.True(t, exist)

	delRequest := &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "204 No Content", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	assert.False(t, exist)

	objectNameNotExist := objectNamePrefix + randLowStr(6) + "-not-exist"
	delRequest = &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameNotExist),
	}
	result, err = client.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "204 No Content", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectNameNotExist)
	assert.Nil(t, err)
	assert.False(t, exist)

	delRequest = &DeleteObjectRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		VersionId: Ptr("null"),
	}
	result, err = client.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "204 No Content", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	_, err = client.DeleteObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	delRequest = &DeleteObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectNamePrefix),
	}
	_, err = client.DeleteObject(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDeleteMultipleObjects(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	delRequest := &DeleteMultipleObjectsRequest{
		Bucket:  Ptr(bucketName),
		Objects: []DeleteObject{{Key: Ptr(objectName)}},
	}
	result, err := client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "200 OK", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Len(t, result.DeletedObjects, 1)
	assert.Equal(t, *result.DeletedObjects[0].Key, objectName)

	str := "\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"
	objectNameSpecial := objectNamePrefix + randLowStr(6) + str
	content = randLowStr(10)
	request = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	delRequest = &DeleteMultipleObjectsRequest{
		Bucket:       Ptr(bucketName),
		Objects:      []DeleteObject{{Key: Ptr(objectNameSpecial)}},
		EncodingType: Ptr("url"),
	}
	result, err = client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "200 OK", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Len(t, result.DeletedObjects, 1)
	assert.Equal(t, *result.DeletedObjects[0].Key, objectNameSpecial)

	_, err = client.DeleteMultipleObjects(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	delRequest = &DeleteMultipleObjectsRequest{
		Bucket:  Ptr(bucketNameNotExist),
		Objects: []DeleteObject{{Key: Ptr(objectNameSpecial)}},
	}
	_, err = client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestHeadObject(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	headRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.HeadObject(context.TODO(), headRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.ContentLength, int64(len(content)))
	assert.NotEmpty(t, *result.ContentMD5)
	assert.NotEmpty(t, *result.ObjectType)
	assert.NotEmpty(t, *result.StorageClass)
	assert.NotEmpty(t, *result.ETag)
	_, err = client.HeadObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	headRequest = &HeadObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.HeadObject(context.TODO(), headRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectMeta(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	headRequest := &GetObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObjectMeta(context.TODO(), headRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.ContentLength, int64(len(content)))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)

	_, err = client.GetObjectMeta(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	headRequest = &GetObjectMetaRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.GetObjectMeta(context.TODO(), headRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestRestoreObject(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(content),
		StorageClass: StorageClassColdArchive,
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	restoreRequest := &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.RestoreObject(context.TODO(), restoreRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 202)
	assert.Equal(t, result.Status, "202 Accepted")
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))

	_, err = client.RestoreObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	restoreRequest = &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err = client.RestoreObject(context.TODO(), restoreRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "RestoreAlreadyInProgress", serr.Code)
	assert.Equal(t, "The restore operation is in progress.", serr.Message)
	assert.NotEmpty(t, serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	restoreRequest = &RestoreObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	_, err = client.RestoreObject(context.TODO(), restoreRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPutObjectAcl(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	objectRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	_, err = client.PutObject(context.TODO(), objectRequest)
	assert.Nil(t, err)
	request := &PutObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Acl:    ObjectACLPublicRead,
	}
	result, err := client.PutObjectAcl(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	infoRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	_, err = client.HeadObject(context.TODO(), infoRequest)
	assert.Nil(t, err)
	_, err = client.PutObjectAcl(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "-not-exist"
	request = &PutObjectAclRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		Acl:    ObjectACLPublicRead,
	}
	_, err = client.PutObjectAcl(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectAcl(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	objectRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Acl:    ObjectACLPublicReadWrite,
	}
	_, err = client.PutObject(context.TODO(), objectRequest)
	assert.Nil(t, err)
	request := &GetObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObjectAcl(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, ObjectACLType(*result.ACL), ObjectACLPublicReadWrite)
	assert.NotEmpty(t, *result.Owner.ID)
	assert.NotEmpty(t, *result.Owner.DisplayName)

	_, err = client.GetObjectAcl(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	objectNameNotExist := objectName + "-not-exist"
	request = &GetObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameNotExist),
	}
	result, err = client.GetObjectAcl(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchKey", serr.Code)
	assert.Equal(t, "The specified key does not exist.", serr.Message)
	assert.Equal(t, "0026-00000001", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestInitiateMultipartUpload(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, initResult.StatusCode)
	assert.NotEmpty(t, initResult.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, *initResult.Bucket, bucketName)
	assert.Equal(t, *initResult.Key, objectName)
	assert.NotEmpty(t, *initResult.UploadId)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.InitiateMultipartUpload(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "-not-exist"
	initRequest = &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	_, err = client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestUploadPart(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	partRequest := &UploadPartRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		PartNumber:   int32(1),
		UploadId:     Ptr(*initResult.UploadId),
		Body:         strings.NewReader("upload part 1"),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, initResult.StatusCode)
	assert.NotEmpty(t, partResult.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *partResult.ETag)
	assert.NotEmpty(t, *partResult.ContentMD5)
	assert.NotEmpty(t, *partResult.HashCRC64)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.UploadPart(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	partRequest = &UploadPartRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		PartNumber:   int32(2),
		UploadId:     Ptr(*initResult.UploadId),
		Body:         strings.NewReader("upload part 2"),
		TrafficLimit: int64(100 * 1024 * 8),
	}

	_, err = client.UploadPart(context.TODO(), partRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000104", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestUploadPartCopy(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	body := randLowStr(100000)
	objectSrcName := objectNamePrefix + randLowStr(6) + "src"
	objRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectSrcName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), objRequest)
	assert.Nil(t, err)
	objectDestName := objectNamePrefix + randLowStr(6) + "dest"
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectDestName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	copyRequest := &UploadPartCopyRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectDestName),
		PartNumber:   int32(1),
		UploadId:     Ptr(*initResult.UploadId),
		SourceKey:    Ptr(objectSrcName),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	copyResult, err := client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, copyResult.StatusCode)
	assert.NotEmpty(t, copyResult.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *copyResult.ETag)
	assert.NotEmpty(t, *copyResult.LastModified)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	metaRequest := &GetObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectSrcName),
	}
	metaResult, err := client.GetObjectMeta(context.TODO(), metaRequest)
	assert.Nil(t, err)
	sourceVersionId := *metaResult.VersionId

	copyRequest = &UploadPartCopyRequest{
		Bucket:          Ptr(bucketName),
		Key:             Ptr(objectDestName),
		PartNumber:      int32(1),
		UploadId:        Ptr(*initResult.UploadId),
		SourceKey:       Ptr(objectSrcName),
		SourceVersionId: Ptr(sourceVersionId),
	}
	copyResult, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, copyResult.StatusCode)
	assert.NotEmpty(t, copyResult.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *copyResult.ETag)
	assert.NotEmpty(t, *copyResult.LastModified)
	assert.NotEmpty(t, *copyResult.VersionId)
	assert.Equal(t, *copyResult.VersionId, sourceVersionId)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectDestName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.UploadPartCopy(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	copyRequest = &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		SourceKey:  Ptr(objectSrcName),
	}
	_, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000311", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestCompleteMultipartUpload(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	body := randLowStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := io.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part1)),
	}
	var parts []UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(2),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part2)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(3),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part3)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	request := &CompleteMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
		CompleteMultipartUpload: &CompleteMultipartUpload{
			Parts: parts,
		},
	}
	result, err := client.CompleteMultipartUpload(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.Location)
	assert.Equal(t, *result.Bucket, bucketName)
	assert.Equal(t, *result.Key, objectName)
	getObj := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	getObjresult, err := client.GetObject(context.TODO(), getObj)
	assert.Nil(t, err)
	data, _ := io.ReadAll(getObjresult.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(data), body)

	objectDestName := objectNamePrefix + randLowStr(6) + "dest" + "\f\v"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectDestName),
	}
	initCopyResult, err := client.InitiateMultipartUpload(context.TODO(), initCopyRequest)
	assert.Nil(t, err)
	copyRequest := &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initCopyResult.UploadId),
		SourceKey:  Ptr(objectName),
	}
	_, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	request = &CompleteMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectDestName),
		UploadId:    Ptr(*initCopyResult.UploadId),
		CompleteAll: Ptr("yes"),
	}
	result, err = client.CompleteMultipartUpload(context.TODO(), request)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.Location)
	assert.Equal(t, *result.Bucket, bucketName)
	assert.Equal(t, *result.Key, objectDestName)

	initCopyResult, err = client.InitiateMultipartUpload(context.TODO(), initCopyRequest)
	assert.Nil(t, err)

	copyRequest = &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initCopyResult.UploadId),
		SourceKey:  Ptr(objectName),
	}
	copyResult, err := client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)

	copyPart := UploadPart{
		PartNumber: copyRequest.PartNumber,
		ETag:       copyResult.ETag,
	}
	var serr *ServiceError
	request = &CompleteMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectDestName),
		UploadId:    Ptr(*initCopyResult.UploadId),
		CompleteAll: Ptr("yes"),
		CompleteMultipartUpload: &CompleteMultipartUpload{
			Parts: []UploadPart{
				copyPart,
			},
		},
	}
	result, err = client.CompleteMultipartUpload(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 400, serr.StatusCode)
	assert.Equal(t, "InvalidArgument", serr.Code)
	assert.Equal(t, "Should not speficy both complete all header and http body.", serr.Message)
	assert.Equal(t, "0042-00000216", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.CompleteMultipartUpload(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CompleteMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectDestName),
		UploadId:    Ptr(*initCopyResult.UploadId),
		CompleteAll: Ptr("yes"),
		Callback:    Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
	}
	result, err = client.CompleteMultipartUpload(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 203, serr.StatusCode)
	assert.Equal(t, "CallbackFailed", serr.Code)
	assert.Equal(t, "Error status : 301.", serr.Message)
	assert.Equal(t, "0007-00000203", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestAbortMultipartUpload(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	result, err := client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))

	_, err = client.AbortMultipartUpload(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	abortRequest = &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000002", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestListMultipartUploads(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	objectName := objectNamePrefix + randLowStr(6) + "\v\n\f"
	body := randLowStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := io.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]

	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part1)),
	}
	var parts []UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(2),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part2)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(3),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part3)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)

	putObj := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(randLowStr(1000)),
	}

	_, err = client.PutObject(context.TODO(), putObj)
	assert.Nil(t, err)
	objectDestName := objectNamePrefix + randLowStr(6) + "dest" + "\f\v\n"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectDestName),
	}

	initCopyResult, err := client.InitiateMultipartUpload(context.TODO(), initCopyRequest)
	assert.Nil(t, err)
	copyRequest := &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initCopyResult.UploadId),
		SourceKey:  Ptr(objectName),
	}
	_, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)

	listRequest := &ListMultipartUploadsRequest{
		Bucket: Ptr(bucketName),
	}
	listResult, err := client.ListMultipartUploads(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, *listResult.Bucket, bucketName)
	assert.Empty(t, *listResult.KeyMarker, bucketName)
	assert.Len(t, listResult.Uploads, 2)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.ListMultipartUploads(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	bucketNameNotExist := bucketName + "-not-exist"
	listRequest = &ListMultipartUploadsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	listResult, err = client.ListMultipartUploads(context.TODO(), listRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestListParts(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6) + "-\v\n\f"
	body := randLowStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := io.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]

	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)

	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part1)),
	}
	var parts []UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)

	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(2),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part2)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)

	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(3),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part3)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)

	listRequest := &ListPartsRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	listResult, err := client.ListParts(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, *listResult.Bucket, bucketName)
	assert.Equal(t, *listResult.Key, objectName)
	assert.Equal(t, *listResult.UploadId, *initResult.UploadId)
	assert.Equal(t, *listResult.StorageClass, "Standard")
	assert.Equal(t, listResult.IsTruncated, false)
	assert.Equal(t, listResult.PartNumberMarker, int32(0))
	assert.Equal(t, listResult.NextPartNumberMarker, int32(3))
	assert.Equal(t, listResult.MaxParts, int32(1000))
	assert.Len(t, listResult.Parts, count)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.ListParts(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	bucketNameNotExist := bucketName + "-not-exist"
	listRequest = &ListPartsRequest{
		Bucket:   Ptr(bucketNameNotExist),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	listResult, err = client.ListParts(context.TODO(), listRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPutBucketVersioning(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	request := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	result, err := client.PutBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	request = &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	result, err = client.PutBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	_, err = client.PutBucketVersioning(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutBucketVersioningRequest{
		Bucket: Ptr(bucketNameNotExist),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	result, err = client.PutBucketVersioning(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetBucketVersioning(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	request := &GetBucketVersioningRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.GetBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Nil(t, result.VersionStatus)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request = &GetBucketVersioningRequest{
		Bucket: Ptr(bucketName),
	}
	result, err = client.GetBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.VersionStatus, "Enabled")

	versionRequest = &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionSuspended,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request = &GetBucketVersioningRequest{
		Bucket: Ptr(bucketName),
	}
	result, err = client.GetBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.VersionStatus, "Suspended")

	_, err = client.GetBucketVersioning(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &GetBucketVersioningRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	result, err = client.GetBucketVersioning(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestListObjectVersions(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request := &GetBucketVersioningRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.GetBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.VersionStatus, "Enabled")

	// put object v1
	content1 := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6) + "\v\f\n"
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content1),
	}
	putObjResult, err := client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	versionIdV1 := putObjResult.Headers.Get("x-oss-version-id")
	assert.True(t, len(versionIdV1) > 0)

	// put object v2
	content2 := randLowStr(200)
	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content2),
	}
	putObjResult, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	versionIdV2 := putObjResult.Headers.Get("x-oss-version-id")
	assert.True(t, len(versionIdV2) > 0)
	assert.NotEqual(t, versionIdV1, versionIdV2)

	delObjRequest := &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	delObjResult, err := client.DeleteObject(context.TODO(), delObjRequest)
	assert.Nil(t, err)
	assert.True(t, delObjResult.DeleteMarker)
	markVersionId := delObjResult.Headers.Get("x-oss-version-id")
	assert.True(t, len(markVersionId) > 0)

	delObjRequest = &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	delObjResult, err = client.DeleteObject(context.TODO(), delObjRequest)
	assert.Nil(t, err)
	assert.True(t, delObjResult.DeleteMarker)
	markVersionIdAgain := delObjResult.Headers.Get("x-oss-version-id")
	assert.True(t, len(markVersionIdAgain) > 0)
	assert.NotEqual(t, markVersionId, markVersionIdAgain)

	versions := &ListObjectVersionsRequest{
		Bucket: Ptr(bucketName),
	}
	versionsResult, err := client.ListObjectVersions(context.TODO(), versions)
	assert.Nil(t, err)
	assert.Len(t, versionsResult.ObjectDeleteMarkers, 2)
	assert.Len(t, versionsResult.ObjectVersions, 2)

	versions = &ListObjectVersionsRequest{
		Bucket: Ptr(bucketName),
		IsMix:  true,
	}
	versionsResult, err = client.ListObjectVersions(context.TODO(), versions)
	assert.Nil(t, err)
	assert.Len(t, versionsResult.ObjectVersionsDeleteMarkers, 4)

	_, err = client.ListObjectVersions(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	versions = &ListObjectVersionsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	versionsResult, err = client.ListObjectVersions(context.TODO(), versions)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPutSymlink(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	symlinkName := objectName + "-symlink"
	request := &PutSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
		Target: Ptr(objectName),
	}
	result, err := client.PutSymlink(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request = &PutSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
		Target: Ptr(objectName),
	}
	result, err = client.PutSymlink(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.VersionId)

	_, err = client.PutSymlink(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutSymlinkRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(symlinkName),
		Target: Ptr(objectName),
	}
	result, err = client.PutSymlink(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetSymlink(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	symlinkName := objectName + "-symlink"
	putSymRequest := &PutSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
		Target: Ptr(objectName),
	}
	_, err = client.PutSymlink(context.TODO(), putSymRequest)
	assert.Nil(t, err)

	request := &GetSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
	}
	result, err := client.GetSymlink(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, result.ETag)
	assert.Equal(t, *result.Target, objectName)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request = &GetSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
	}
	result, err = client.GetSymlink(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, result.ETag)
	assert.Equal(t, *result.Target, objectName)
	assert.NotEmpty(t, *result.VersionId)

	_, err = client.GetSymlink(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &GetSymlinkRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(symlinkName),
	}
	result, err = client.GetSymlink(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPutObjectTagging(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	request := &PutObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
	}
	result, err := client.PutObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	putObjResult, err := client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	versionId := *putObjResult.VersionId
	request = &PutObjectTaggingRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		VersionId: Ptr(versionId),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
				},
			},
		},
	}
	result, err = client.PutObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.VersionId, versionId)

	_, err = client.PutObjectTagging(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutObjectTaggingRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
	}
	result, err = client.PutObjectTagging(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectTagging(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	request := &GetObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Len(t, result.Tags, 0)

	putTagRequest := &PutObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
	}
	_, err = client.PutObjectTagging(context.TODO(), putTagRequest)
	assert.Nil(t, err)

	request = &GetObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err = client.GetObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Len(t, result.Tags, 2)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	putObjResult, err := client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	versionId := *putObjResult.VersionId

	request = &GetObjectTaggingRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		VersionId: Ptr(versionId),
	}
	result, err = client.GetObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Len(t, result.Tags, 0)

	_, err = client.GetObjectTagging(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &GetObjectTaggingRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.GetObjectTagging(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDeleteObjectTagging(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	putTagRequest := &PutObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
	}
	_, err = client.PutObjectTagging(context.TODO(), putTagRequest)
	assert.Nil(t, err)

	request := &DeleteObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.DeleteObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	putObjResult, err := client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	versionId := *putObjResult.VersionId

	request = &DeleteObjectTaggingRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		VersionId: Ptr(versionId),
	}
	result, err = client.DeleteObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	_, err = client.DeleteObjectTagging(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &DeleteObjectTaggingRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.DeleteObjectTagging(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestCreateSelectObjectMeta(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := "name,school,company,age\r\nLora Francis,School A,Staples Inc,27\r\n" + "Eleanor Little,School B,\"Conectiv, Inc\",43\r\n" + "Rosie Hughes,School C,Western Gas Resources Inc,44\r\n" + "Lawrence Ross,School D,MetLife Inc.,24"
	objectNameCsv := objectNamePrefix + randLowStr(6) + ".csv"
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	csvMeta := &CreateSelectObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		MetaRequest: &CsvMetaRequest{
			OverwriteIfExists: Ptr(true),
		},
	}
	result, err := client.CreateSelectObjectMeta(context.TODO(), csvMeta)
	assert.Nil(t, err)
	assert.Equal(t, result.RowsCount, int64(5))

	body = "{\n" +
		"\t\"name\": \"Lora Francis\",\n" +
		"\t\"age\": 27,\n" +
		"\t\"company\": \"Staples Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Eleanor Little\",\n" +
		"\t\"age\": 43,\n" +
		"\t\"company\": \"Conectiv, Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Rosie Hughes\",\n" +
		"\t\"age\": 44,\n" +
		"\t\"company\": \"Western Gas Resources Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Lawrence Ross\",\n" +
		"\t\"age\": 24,\n" +
		"\t\"company\": \"MetLife Inc.\"\n" +
		"}"
	objectNameJson := objectNamePrefix + randLowStr(6) + ".json"
	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameJson),
		Body:   strings.NewReader(string(body)),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	csvMeta = &CreateSelectObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameJson),
		MetaRequest: &JsonMetaRequest{
			InputSerialization: &InputSerialization{
				JSON: &InputSerializationJSON{
					JSONType: Ptr("LINES"),
				},
			},
		},
	}
	result, err = client.CreateSelectObjectMeta(context.TODO(), csvMeta)
	assert.Nil(t, err)
	assert.Equal(t, result.RowsCount, int64(4))

	_, err = client.CreateSelectObjectMeta(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	csvMeta = &CreateSelectObjectMetaRequest{
		Bucket:      Ptr(bucketNameNotExist),
		Key:         Ptr(objectNameCsv),
		MetaRequest: &CsvMetaRequest{},
	}
	result, err = client.CreateSelectObjectMeta(context.TODO(), csvMeta)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestSelectObject(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := "name,school,company,age\r\nLora Francis,School A,Staples Inc,27\r\n" + "Eleanor Little,School B,\"Conectiv, Inc\",43\r\n" + "Rosie Hughes,School C,Western Gas Resources Inc,44\r\n" + "Lawrence Ross,School D,MetLife Inc.,24"
	objectNameCsv := objectNamePrefix + randLowStr(6) + ".csv"
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	request := &SelectObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select name from ossobject"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				OutputHeader: Ptr(true),
			},
		},
	}
	result, err := client.SelectObject(context.TODO(), request)
	assert.Nil(t, err)
	dataByte, err := io.ReadAll(result.Body)
	assert.Equal(t, string(dataByte), "name\nLora Francis\nEleanor Little\nRosie Hughes\nLawrence Ross\n")

	body = "{\n" +
		"\t\"name\": \"Lora Francis\",\n" +
		"\t\"age\": 27,\n" +
		"\t\"company\": \"Staples Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Eleanor Little\",\n" +
		"\t\"age\": 43,\n" +
		"\t\"company\": \"Conectiv, Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Rosie Hughes\",\n" +
		"\t\"age\": 44,\n" +
		"\t\"company\": \"Western Gas Resources Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Lawrence Ross\",\n" +
		"\t\"age\": 24,\n" +
		"\t\"company\": \"MetLife Inc.\"\n" +
		"}"
	objectNameJson := objectNamePrefix + randLowStr(6) + ".json"
	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameJson),
		Body:   strings.NewReader(string(body)),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	request = &SelectObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select name from ossobject"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				OutputHeader: Ptr(true),
			},
		},
	}
	result, err = client.SelectObject(context.TODO(), request)
	assert.Nil(t, err)
	dataByte, err = io.ReadAll(result.Body)
	assert.Equal(t, string(dataByte), "name\nLora Francis\nEleanor Little\nRosie Hughes\nLawrence Ross\n")

	_, err = client.SelectObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &SelectObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectNameCsv),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select name from ossobject"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				OutputHeader: Ptr(true),
			},
		},
	}
	result, err = client.SelectObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPresign(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	// PutObjRequest
	body := randLowStr(1000)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	req, err := http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	c := &http.Client{}
	resp, err := c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	// GetObjRequest
	getObjRequest := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration := time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	data, _ := io.ReadAll(resp.Body)
	assert.Equal(t, string(data), body)

	// HeadObjRequest
	headObjRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), headObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header.Get(HTTPHeaderContentLength), strconv.Itoa(len(body)))

	// MultiPart
	objectNameMultipart := objectNamePrefix + randLowStr(6) + "-multi-part"
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipart),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), initRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	initResult := &InitiateMultipartUploadResult{}
	err = xml.Unmarshal(data, initResult)
	assert.Equal(t, *initResult.Key, objectNameMultipart)
	uploadId := initResult.UploadId

	//UploadPart
	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectNameMultipart),
		PartNumber: int32(1),
		UploadId:   uploadId,
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), partRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	var parts []UploadPart
	uploadResult := &UploadPartResult{}
	err = xml.Unmarshal(data, uploadResult)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       Ptr(resp.Header.Get("ETag")),
	}
	parts = append(parts, part)
	completeRequest := &CompleteMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipart),
		UploadId: uploadId,
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), completeRequest, PresignExpiration(expiration))
	assert.Nil(t, err)

	//Complete
	upload := CompleteMultipartUpload{
		Parts: parts,
	}
	xmlData, err := xml.Marshal(upload)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(string(xmlData)))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	headObjRequest = &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipart),
	}
	headResult, err := client.HeadObject(context.TODO(), headObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, headResult.Headers.Get(HTTPHeaderContentLength), strconv.FormatInt(int64(len(body)), 10))
	assert.Equal(t, *headResult.ObjectType, "Multipart")

	// Test Abort
	objectNameMultipartCopy := objectNamePrefix + randLowStr(6) + "-multi-part-copy"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipartCopy),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), initCopyRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	initCopyResult := &InitiateMultipartUploadResult{}
	err = xml.Unmarshal(data, initCopyResult)
	assert.Equal(t, *initCopyResult.Key, objectNameMultipartCopy)
	copyUploadId := *initCopyResult.UploadId

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipartCopy),
		UploadId: Ptr(copyUploadId),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), abortRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 204)

	listPartsRequest := &ListPartsRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipartCopy),
		UploadId: Ptr(copyUploadId),
	}
	_, err = client.ListParts(context.TODO(), listPartsRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000002", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	deleteBucket(bucketName, t)
}

func TestPresignExtra(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	// PutObjRequest
	body := randLowStr(1000)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	req, err := http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	c := &http.Client{}
	resp, err := c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	cfgV1 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithSignatureVersion(SignatureVersionV1).WithAdditionalHeaders([]string{"content-length"})

	clientV1 := NewClient(cfgV1)
	getObjRequest := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration := time.Now().Add(1 * time.Second)
	result, err = clientV1.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Equal(t, map[string]string(nil), result.SignedHeaders)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	data, _ := io.ReadAll(resp.Body)
	assert.Equal(t, string(data), body)

	getObjRequest = &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(-1 * time.Second)
	result, err = clientV1.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Equal(t, map[string]string(nil), result.SignedHeaders)

	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 403)

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"content-length": "1000",
			},
		},
	}
	result, err = clientV1.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string(nil), result.SignedHeaders)
	assert.NotEmpty(t, result.Expiration)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	cfgV4 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithSignatureVersion(SignatureVersionV4).WithAdditionalHeaders([]string{"content-length"})

	clientV4 := NewClient(cfgV4)
	getObjRequest = &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(1 * time.Second)
	result, err = clientV4.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Equal(t, map[string]string(nil), result.SignedHeaders)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	data, _ = io.ReadAll(resp.Body)
	assert.Equal(t, string(data), body)

	getObjRequest = &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(8 * 24 * time.Hour)
	result, err = clientV4.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "expires should be not greater than 604800(seven days)")

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"content-length": "1000",
			},
		},
	}
	result, err = clientV4.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"Content-Length": "1000"}, result.SignedHeaders)
	assert.NotEmpty(t, result.Expiration)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader("hi oss"))
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 403)

	cfgV4 = LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithSignatureVersion(SignatureVersionV4).WithAdditionalHeaders([]string{"email", "name"})

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"email": "demo@aliyun.com",
				"name":  "aliyun",
			},
		},
	}
	clientV4 = NewClient(cfgV4)
	result, err = clientV4.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"Email": "demo@aliyun.com",
		"Name": "aliyun"}, result.SignedHeaders)
	assert.NotEmpty(t, result.Expiration)

	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 400)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)

	header := make(http.Header)
	for key, value := range result.SignedHeaders {
		header[key] = []string{value}
	}
	req.Header = header
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	deleteBucket(bucketName, t)
}

func TestPresignWithStsToken(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	client := getClientUseStsToken(region_, endpoint_)
	assert.NotNil(t, client)

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := randLowStr(1000)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	req, err := http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	c := &http.Client{}
	resp, err := c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	getObjRequest := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration := time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	data, _ := io.ReadAll(resp.Body)
	assert.Equal(t, string(data), body)

	headObjRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), headObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header.Get(HTTPHeaderContentLength), fmt.Sprint(len(body)))

	objectNameMultipart := objectNamePrefix + randLowStr(6) + "-multi-part"
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipart),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), initRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	initResult := &InitiateMultipartUploadResult{}
	err = xml.Unmarshal(data, initResult)
	assert.Equal(t, *initResult.Key, objectNameMultipart)
	uploadId := initResult.UploadId

	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectNameMultipart),
		PartNumber: int32(1),
		UploadId:   uploadId,
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), partRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	parts := []UploadPart{}
	uploadResult := &UploadPartResult{}
	err = xml.Unmarshal(data, uploadResult)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       Ptr(resp.Header.Get("ETag")),
	}
	parts = append(parts, part)
	completeRequest := &CompleteMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipart),
		UploadId: uploadId,
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), completeRequest, PresignExpiration(expiration))
	assert.Nil(t, err)

	upload := CompleteMultipartUpload{
		Parts: parts,
	}
	xmlData, err := xml.Marshal(upload)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(string(xmlData)))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	headObjRequest = &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipart),
	}
	headResult, err := client.HeadObject(context.TODO(), headObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, headResult.Headers.Get(HTTPHeaderContentLength), strconv.FormatInt(int64(len(body)), 10))
	assert.Equal(t, *headResult.ObjectType, "Multipart")

	objectNameMultipartCopy := objectNamePrefix + randLowStr(6) + "-multi-part-copy"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipartCopy),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), initCopyRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	initCopyResult := &InitiateMultipartUploadResult{}
	err = xml.Unmarshal(data, initCopyResult)
	assert.Equal(t, *initCopyResult.Key, objectNameMultipartCopy)
	copyUploadId := *initCopyResult.UploadId

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipartCopy),
		UploadId: Ptr(copyUploadId),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), abortRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 204)

	listPartsRequest := &ListPartsRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipartCopy),
		UploadId: Ptr(copyUploadId),
	}
	_, err = client.ListParts(context.TODO(), listPartsRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000002", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	cleanObjects(client, bucketName, t)
}

func TestPaginator(t *testing.T) {
	after := before(t)
	defer after(t)
	var bucketName string
	client := getDefaultClient()
	count := 10
	bucketNameTestPrefix := bucketNamePrefix + randLowStr(6) + "-"
	for i := 0; i < count; i++ {
		bucketName = bucketNameTestPrefix + strconv.Itoa(i)
		putRequest := &PutBucketRequest{
			Bucket: Ptr(bucketName),
		}
		_, err := client.PutBucket(context.TODO(), putRequest)
		assert.Nil(t, err)
	}

	lbRequest := &ListBucketsRequest{
		MaxKeys: int32(4),
		Prefix:  Ptr(bucketNameTestPrefix),
	}
	lbPaginator := client.NewListBucketsPaginator(lbRequest)
	countBucket := 0
	for lbPaginator.HasNext() {
		result, err := lbPaginator.NextPage(context.TODO())
		assert.Nil(t, err)
		countBucket += len(result.Buckets)
	}
	assert.Equal(t, countBucket, count)

	lbPaginator = client.NewListBucketsPaginator(nil)
	countBucket = 0
	for lbPaginator.HasNext() {
		result, err := lbPaginator.NextPage(context.TODO())
		assert.Nil(t, err)
		countBucket += len(result.Buckets)
	}
	assert.True(t, countBucket >= count)

	listBucket, err := client.ListBuckets(context.TODO(), &ListBucketsRequest{
		Prefix: Ptr(bucketNameTestPrefix),
	})
	assert.Nil(t, err)
	bucketNameTest := *listBucket.Buckets[0].Name

	var objName string
	countObj := 10
	objectNameTestPrefix := objectNamePrefix + randLowStr(6) + "-"
	for i := 0; i < countObj; i++ {
		objName = objectNameTestPrefix + strconv.Itoa(i)
		putRequest := &PutObjectRequest{
			Bucket: Ptr(bucketNameTest),
			Key:    Ptr(objName),
		}
		_, err = client.PutObject(context.TODO(), putRequest)
		assert.Nil(t, err)
	}

	var listObjCount int
	listObjRequest := &ListObjectsRequest{
		Bucket:  Ptr(bucketNameTest),
		MaxKeys: int32(4),
	}
	listObjPaginator := client.NewListObjectsPaginator(listObjRequest)
	for listObjPaginator.HasNext() {
		result, err := listObjPaginator.NextPage(context.TODO())
		assert.Nil(t, err)
		listObjCount += len(result.Contents)
	}
	assert.Equal(t, countObj, listObjCount)

	listObjPaginator = client.NewListObjectsPaginator(nil)
	listObjCount = 0
	for listObjPaginator.HasNext() {
		_, err = listObjPaginator.NextPage(context.TODO())
		assert.NotNil(t, err)
		break
	}

	var listObjCountV2 int
	listObjV2Request := &ListObjectsV2Request{
		Bucket:  Ptr(bucketNameTest),
		MaxKeys: int32(4),
	}
	listObjV2Paginator := client.NewListObjectsV2Paginator(listObjV2Request)
	for listObjV2Paginator.HasNext() {
		result, err := listObjV2Paginator.NextPage(context.TODO())
		assert.Nil(t, err)
		listObjCountV2 += len(result.Contents)
	}
	assert.Equal(t, countObj, listObjCountV2)

	listObjV2Paginator = client.NewListObjectsV2Paginator(nil)
	listObjCountV2 = 0
	for listObjPaginator.HasNext() {
		_, err = listObjPaginator.NextPage(context.TODO())
		assert.NotNil(t, err)
		break
	}

	var listObjVersionCount, listObjDeleted int
	lovRequest := &ListObjectVersionsRequest{
		Bucket:  Ptr(bucketNameTest),
		MaxKeys: int32(4),
	}
	lovPaginator := client.NewListObjectVersionsPaginator(lovRequest)
	for lovPaginator.HasNext() {
		result, err := lovPaginator.NextPage(context.TODO())
		assert.Nil(t, err)
		listObjVersionCount += len(result.ObjectVersions)
		listObjDeleted += len(result.ObjectDeleteMarkers)
	}
	assert.Equal(t, countObj, listObjVersionCount)
	assert.Equal(t, 0, listObjDeleted)

	lovPaginator = client.NewListObjectVersionsPaginator(nil)
	for lovPaginator.HasNext() {
		_, err = lovPaginator.NextPage(context.TODO())
		assert.NotNil(t, err)
		break
	}

	var objMultiName string
	countObjMulti := 20
	for i := 0; i < countObjMulti; i++ {
		objMultiName = objectNameTestPrefix + "multi-part-" + strconv.Itoa(i)
		_, err = client.InitiateMultipartUpload(context.TODO(), &InitiateMultipartUploadRequest{
			Bucket: Ptr(bucketNameTest),
			Key:    Ptr(objMultiName),
		})
		assert.Nil(t, err)
	}

	var countUploads int
	lmuRequest := &ListMultipartUploadsRequest{
		Bucket:     Ptr(bucketNameTest),
		MaxUploads: int32(8),
	}
	lmuPaginator := client.NewListMultipartUploadsPaginator(lmuRequest)
	for lmuPaginator.HasNext() {
		result, err := lmuPaginator.NextPage(context.TODO())
		assert.Nil(t, err)
		countUploads += len(result.Uploads)
	}
	assert.Equal(t, countObjMulti, countUploads)

	lmuPaginator = client.NewListMultipartUploadsPaginator(nil)
	for lmuPaginator.HasNext() {
		_, err = lmuPaginator.NextPage(context.TODO())
		assert.NotNil(t, err)
		break
	}

	uploadsResult, err := client.ListMultipartUploads(context.TODO(), &ListMultipartUploadsRequest{
		Bucket: Ptr(bucketNameTest),
	})
	assert.Nil(t, err)

	objectName := *uploadsResult.Uploads[0].Key
	uploadId := *uploadsResult.Uploads[0].UploadId
	data := randLowStr(1024 * 1024 * 20)
	countPart := 20
	lenStr := len(data)
	avgLen := lenStr / countPart
	result := make([]string, 0)
	for i := 0; i < lenStr; i += avgLen {
		end := i + avgLen
		if end > lenStr {
			end = lenStr
		}
		result = append(result, data[i:end])
	}

	for k, content := range result {
		_, err = client.UploadPart(context.TODO(), &UploadPartRequest{
			Bucket:     Ptr(bucketNameTest),
			Key:        Ptr(objectName),
			UploadId:   Ptr(uploadId),
			PartNumber: int32(k + 1),
			Body:       strings.NewReader(content),
		})
		assert.Nil(t, err)
	}

	var countPartResult int
	lpRequest := &ListPartsRequest{
		Bucket:   Ptr(bucketNameTest),
		Key:      Ptr(objectName),
		UploadId: Ptr(uploadId),
		MaxParts: int32(6),
	}
	lpPaginator := client.NewListPartsPaginator(lpRequest)
	for lpPaginator.HasNext() {
		result, err := lpPaginator.NextPage(context.TODO())
		assert.Nil(t, err)
		countPartResult += len(result.Parts)
	}
	assert.Equal(t, countPart, countPartResult)

	lpPaginator = client.NewListPartsPaginator(nil)
	for lmuPaginator.HasNext() {
		_, err = lpPaginator.NextPage(context.TODO())
		assert.NotNil(t, err)
		break
	}
}

func TestEncryptionClient(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	objectName := objectNamePrefix + randLowStr(6)

	length := 3*100*1024 + 123
	partSize := int64(200 * 1024)
	partsNum := length/int(partSize) + 1
	data := []byte(randStr(length))
	hashData := NewCRC64(0)
	hashData.Write(data)

	client := getDefaultClient()
	assert.NotNil(t, client)

	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	mc, err := crypto.CreateMasterRsa(map[string]string{"tag": "value"}, rsaPublicKey, rsaPrivateKey)
	assert.Nil(t, err)
	eclient, err := NewEncryptionClient(client, mc)
	assert.Nil(t, err)

	initResult, err := eclient.InitiateMultipartUpload(context.TODO(), &InitiateMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectName),
		CSEPartSize: Ptr(partSize),
		CSEDataSize: Ptr(int64(length)),
	})
	assert.Nil(t, err)
	assert.NotNil(t, initResult)
	assert.NotNil(t, initResult.CSEMultiPartContext)
	assert.NotNil(t, initResult.CSEMultiPartContext.ContentCipher)
	assert.Equal(t, partSize, initResult.CSEMultiPartContext.PartSize)
	assert.Equal(t, int64(length), initResult.CSEMultiPartContext.DataSize)

	var parts UploadParts
	for i := 0; i < partsNum; i++ {
		start := i * int(partSize)
		end := start + int(partSize)
		end = minInt(end, length)
		var contentLength *int64 = nil
		if i%2 == 0 {
			contentLength = Ptr(int64(end - start))
		}
		upResult, err := eclient.UploadPart(context.TODO(), &UploadPartRequest{
			Bucket:              Ptr(bucketName),
			Key:                 Ptr(objectName),
			UploadId:            initResult.UploadId,
			PartNumber:          int32(i + 1),
			CSEMultiPartContext: initResult.CSEMultiPartContext,
			ContentLength:       contentLength,
			Body:                bytes.NewReader(data[start:end]),
		})
		assert.Nil(t, err)
		assert.NotNil(t, upResult)
		parts = append(parts, UploadPart{PartNumber: int32(i + 1), ETag: upResult.ETag})
	}

	lsResult, err := eclient.ListParts(context.TODO(), &ListPartsRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: initResult.UploadId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, lsResult)

	sort.Sort(parts)
	cmResult, err := eclient.CompleteMultipartUpload(context.TODO(), &CompleteMultipartUploadRequest{
		Bucket:                  Ptr(bucketName),
		Key:                     Ptr(objectName),
		UploadId:                initResult.UploadId,
		CompleteMultipartUpload: &CompleteMultipartUpload{Parts: parts},
	})
	assert.Nil(t, err)
	assert.NotNil(t, cmResult)

	// GetObject
	gResult, err := eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err := io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.Len(t, gData, length)
	assert.EqualValues(t, data, gData)

	assert.NotEmpty(t, gResult.Headers.Get(OssClientSideEncryptionKey))
	assert.NotEmpty(t, gResult.Headers.Get(OssClientSideEncryptionStart))
	assert.Equal(t, crypto.AesCtrAlgorithm, gResult.Headers.Get(OssClientSideEncryptionCekAlg))
	assert.Equal(t, crypto.RsaCryptoWrap, gResult.Headers.Get(OssClientSideEncryptionWrapAlg))
	assert.Equal(t, "{\"tag\":\"value\"}", gResult.Headers.Get(OssClientSideEncryptionMatDesc))
	assert.Equal(t, fmt.Sprint(partSize), gResult.Headers.Get(OssClientSideEncryptionPartSize))
	assert.Equal(t, fmt.Sprint(length), gResult.Headers.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, gResult.Headers.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, gResult.Headers.Get(OssClientSideEncryptionUnencryptedContentMD5))

	// HeadObject
	hResult, err := eclient.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.NotEmpty(t, hResult.Headers.Get(OssClientSideEncryptionKey))
	assert.NotEmpty(t, hResult.Headers.Get(OssClientSideEncryptionStart))
	assert.Equal(t, crypto.AesCtrAlgorithm, hResult.Headers.Get(OssClientSideEncryptionCekAlg))
	assert.Equal(t, crypto.RsaCryptoWrap, hResult.Headers.Get(OssClientSideEncryptionWrapAlg))
	assert.Equal(t, "{\"tag\":\"value\"}", hResult.Headers.Get(OssClientSideEncryptionMatDesc))
	assert.Equal(t, fmt.Sprint(partSize), hResult.Headers.Get(OssClientSideEncryptionPartSize))
	assert.Equal(t, fmt.Sprint(length), hResult.Headers.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, hResult.Headers.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, hResult.Headers.Get(OssClientSideEncryptionUnencryptedContentMD5))
	assert.Equal(t, int64(length), hResult.ContentLength)

	// HeadObject
	gmResult, err := eclient.GetObjectMeta(context.TODO(), &GetObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionKey))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionStart))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionCekAlg))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionWrapAlg))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionMatDesc))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionPartSize))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionUnencryptedContentMD5))
	assert.Equal(t, int64(length), gmResult.ContentLength)

	// Downloader with not 16 align partSize
	d := eclient.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 123 * 1024
	})
	assert.NotNil(t, d)
	assert.Equal(t, int64(123*1024), d.options.PartSize)
	assert.Equal(t, 3, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"
	dResult, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName)},
		localFile)
	defer os.Remove(localFile)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(gData)), dResult.Written)
	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)
	io.Copy(hash, rfile)
	rfile.Close()
	assert.Equal(t, hash.Sum64(), hashData.Sum64())

	//Use ReadOnlyFile
	f, err := eclient.OpenFile(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	for i := 13; i < 42; i++ {
		for len := 100*1024 + 123; len < 100*1024+123+17; len++ {
			_, err := f.Seek(int64(i), io.SeekStart)
			assert.Nil(t, err)
			gData, err := io.ReadAll(io.LimitReader(f, int64(len)))
			assert.Nil(t, err)
			assert.EqualValues(t, data[i:i+len], gData)
		}
	}
	f.Close()
	time.Sleep(2 * time.Second)

	// Use Uploader
	lastEtag := hResult.Headers.Get(HTTPHeaderETag)
	assert.NotEmpty(t, lastEtag)
	u := eclient.NewUploader()
	assert.NotNil(t, u)
	urResult, err := u.UploadFrom(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
		},
		bytes.NewReader(data),
		func(uo *UploaderOptions) {
			uo.ParallelNum = 2
			uo.PartSize = 100 * 1024
		},
	)
	if !assert.Nil(t, err) {
		fmt.Printf("%s", err.Error())
	}
	assert.NotNil(t, urResult)

	// GetObject again
	gResult, err = eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err = io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.Len(t, gData, length)
	assert.EqualValues(t, data, gData)

	assert.NotEmpty(t, gResult.Headers.Get(OssClientSideEncryptionKey))
	assert.NotEmpty(t, gResult.Headers.Get(OssClientSideEncryptionStart))
	assert.Equal(t, crypto.AesCtrAlgorithm, gResult.Headers.Get(OssClientSideEncryptionCekAlg))
	assert.Equal(t, crypto.RsaCryptoWrap, gResult.Headers.Get(OssClientSideEncryptionWrapAlg))
	assert.Equal(t, "{\"tag\":\"value\"}", gResult.Headers.Get(OssClientSideEncryptionMatDesc))
	assert.Equal(t, fmt.Sprint(100*1024), gResult.Headers.Get(OssClientSideEncryptionPartSize))
	assert.Equal(t, fmt.Sprint(length), gResult.Headers.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, gResult.Headers.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, gResult.Headers.Get(OssClientSideEncryptionUnencryptedContentMD5))

	assert.NotEqual(t, lastEtag, ToString(gResult.ETag))
}

func TestClientExtension(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)
	bucketNameNoExist := bucketName + "-no-exist"
	objectNameNoExist := objectName + "-no-exist"

	client := getDefaultClient()
	assert.NotNil(t, client)

	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	assert.NotNil(t, noPermClient)

	errorClient := getClientWithCredentialsProvider("", "",
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	assert.NotNil(t, errorClient)

	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)

	// IsBucketExist
	exist, err := client.IsBucketExist(context.TODO(), bucketName)
	assert.Nil(t, err)
	assert.True(t, exist)

	exist, err = client.IsBucketExist(context.TODO(), bucketNameNoExist)
	assert.Nil(t, err)
	assert.False(t, exist)

	exist, err = noPermClient.IsBucketExist(context.TODO(), bucketName)
	assert.Nil(t, err)
	assert.True(t, exist)

	exist, err = noPermClient.IsBucketExist(context.TODO(), bucketNameNoExist)
	assert.Nil(t, err)
	assert.False(t, exist)

	exist, err = errorClient.IsBucketExist(context.TODO(), bucketName)
	assert.NotNil(t, err)
	var serr *ServiceError
	assert.False(t, errors.As(err, &serr))

	// IsObjectExist
	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	assert.True(t, exist)

	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectNameNoExist)
	assert.Nil(t, err)
	assert.False(t, exist)

	exist, err = client.IsObjectExist(context.TODO(), bucketNameNoExist, objectName)
	assert.NotNil(t, err)
	assert.False(t, exist)
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "NoSuchBucket", serr.Code)

	exist, err = client.IsObjectExist(context.TODO(), bucketNameNoExist, objectNameNoExist)
	assert.NotNil(t, err)
	assert.False(t, exist)
	assert.NotNil(t, serr)
	assert.Equal(t, "NoSuchBucket", serr.Code)

	exist, err = noPermClient.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.NotNil(t, err)
	assert.False(t, exist)
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)

	exist, err = noPermClient.IsObjectExist(context.TODO(), bucketNameNoExist, objectName)
	assert.NotNil(t, err)
	assert.False(t, exist)
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "NoSuchBucket", serr.Code)

	exist, err = errorClient.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.NotNil(t, err)
	assert.False(t, exist)
	assert.False(t, errors.As(err, &serr))

	//PutObjectFromFile
	objectNameFromFile := objectName + "-from-file"
	var localFile = randStr(8) + ".txt"
	length := 1234
	content := randStr(length)
	hashContent := NewCRC64(0)
	hashContent.Write([]byte(content))
	createFile(t, localFile, content)
	defer func() { os.Remove(localFile) }()

	result, err := client.PutObjectFromFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameFromFile),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	gResult, err := client.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameFromFile),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, fmt.Sprint(hashContent.Sum64()), ToString(gResult.HashCRC64))
	_, err = io.ReadAll(gResult.Body)
	gResult.Body.Close()
	dumpErrIfNotNil(err)

	// Use Uploader, set meta and acl
	objectNameBig := objectName + "-big"
	bigLength := 5*100*1024 + 1234
	bigContent := randStr(bigLength)
	bigHash := NewCRC64(0)
	bigHash.Write([]byte(bigContent))
	u := client.NewUploader()
	assert.NotNil(t, u)
	urResult, err := u.UploadFrom(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectNameBig),
			Metadata: map[string]string{
				"author": "test",
				"magic":  "123",
			},
			Acl: ObjectACLPublicRead,
		},
		bytes.NewReader([]byte(bigContent)),
		func(uo *UploaderOptions) {
			uo.ParallelNum = 3
			uo.PartSize = 100 * 1024
		},
	)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.NotNil(t, urResult)

	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectNameBig)
	assert.Nil(t, err)
	assert.True(t, exist)

	hResult, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameBig),
	})
	assert.Nil(t, err)
	assert.NotNil(t, hResult)
	assert.Contains(t, hResult.Headers.Get(HTTPHeaderETag), "-6")
	assert.Equal(t, "Multipart", hResult.Headers.Get(HeaderOssObjectType))
	assert.Equal(t, "test", hResult.Headers.Get("x-oss-meta-author"))
	assert.Equal(t, "123", hResult.Headers.Get("x-oss-meta-magic"))

	aclResult, err := client.GetObjectAcl(context.TODO(), &GetObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameBig),
	})
	assert.Nil(t, err)
	assert.NotNil(t, hResult)
	assert.Equal(t, "public-read", ToString(aclResult.ACL))

	// Downloader with not align partSize
	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 100*1024 + 123
	})
	assert.NotNil(t, d)
	assert.Equal(t, int64(100*1024+123), d.options.PartSize)
	assert.Equal(t, 3, d.options.ParallelNum)
	localFileBig := randStr(8) + "-downloader"
	dResult, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectNameBig)},
		localFileBig)

	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, int64(bigLength), dResult.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFileBig)
	assert.Nil(t, err)
	defer func() {
		rfile.Close()
		os.Remove(localFileBig)
	}()
	io.Copy(hash, rfile)
	assert.Equal(t, bigHash.Sum64(), hash.Sum64())

	//Use ReadOnlyFile
	f, err := client.OpenFile(context.TODO(), bucketName, objectNameBig)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	for i := 13; i < 42; i++ {
		for len := 100*1024 + 123; len < 100*1024+123+17; len++ {
			_, err := f.Seek(int64(i), io.SeekStart)
			assert.Nil(t, err)
			gData, err := io.ReadAll(io.LimitReader(f, int64(len)))
			assert.Nil(t, err)
			assert.EqualValues(t, []byte(bigContent)[i:i+len], gData)
		}
	}
	f.Close()

	// AppenableFile
	objectNameAppend := objectName + "-append"
	dataa1 := []byte("helle world")
	dataa2 := []byte(randStr(12345))
	dataa3 := []byte(randStr(100*1024*5 + 13))
	var localFileData3 = randStr(8) + ".txt"
	createFile(t, localFileData3, string(dataa3))
	defer func() {
		os.Remove(localFileData3)
	}()

	af, err := client.AppendFile(context.TODO(), bucketName, objectNameAppend)
	n, err := af.Write(dataa1)
	assert.Nil(t, err)
	assert.Equal(t, len(dataa1), n)

	hResult, err = client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend),
	})
	assert.Nil(t, err)
	assert.NotNil(t, hResult)
	assert.Equal(t, int64(n), hResult.ContentLength)
	nl, err := af.WriteFrom(bytes.NewReader(dataa2))
	assert.Equal(t, int64(len(dataa2)), nl)

	filedataa3, err := os.Open(localFileData3)
	assert.Nil(t, err)
	nl, err = io.Copy(af, filedataa3)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(dataa3)), nl)
	defer func() {
		filedataa3.Close()
	}()

	af.Close()
	hashA := NewCRC64(0)
	hashA.Write(dataa1)
	hashA.Write(dataa2)
	hashA.Write(dataa3)
	hResult, err = client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend),
	})
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(hashA.Sum64()), ToString(hResult.HashCRC64))

	//GetObjectToFile
	var localFileToFile = randStr(8) + "-to-file"
	defer func() {
		os.Remove(localFileToFile)
	}()
	gResult, err = client.GetObjectToFile(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend)},
		localFileToFile,
	)
	assert.Nil(t, err)
	gResult, err = client.GetObjectToFile(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend),
		ProgressFn: func(increment, transferred, total int64) {
			//fmt.Printf("increment:%v, transferred:%v, total:%v\n", increment, transferred, total)
		}},
		localFileToFile,
	)
	assert.Nil(t, err)
	hash = NewCRC64(0)
	rfiletoFile, err := os.Open(localFileToFile)
	assert.Nil(t, err)
	defer func() {
		rfiletoFile.Close()
	}()
	io.Copy(hash, rfiletoFile)
	assert.Equal(t, hashA.Sum64(), hash.Sum64())
}

func TestProcessObject(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6) + ".jpg"
	objectDestName := objectNamePrefix + randLowStr(6) + "dest.jpg"

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}

	filePath := "../sample/example.jpg"
	_, err = client.PutObjectFromFile(context.TODO(), putObjRequest, filePath)
	assert.Nil(t, err)

	request := &ProcessObjectRequest{
		Bucket:  Ptr(bucketName),
		Key:     Ptr(objectName),
		Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte(objectDestName)))),
	}
	result, err := client.ProcessObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, result.Bucket, "")
	assert.NotEmpty(t, result.FileSize)
	assert.Equal(t, result.Object, objectDestName)
	assert.Equal(t, result.ProcessStatus, "OK")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &ProcessObjectRequest{
		Bucket:  Ptr(bucketNameNotExist),
		Key:     Ptr(objectName),
		Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte(objectDestName)))),
	}
	result, err = client.ProcessObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestAsyncProcessObject(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6) + ".mp4"
	objectDestName := objectNamePrefix + randLowStr(6) + "dest.mp4"

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	putObjrequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	videoUrl := "https://oss-console-img-demo-cn-hangzhou.oss-cn-hangzhou.aliyuncs.com/video.mp4?spm=a2c4g.64555.0.0.515675979u4B8w&file=video.mp4"
	fileName := "video.mp4"
	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = http.Get(videoUrl)
		if err != nil {
			continue
		}
	}
	assert.Nil(t, err)
	defer resp.Body.Close()
	defer os.Remove(fileName)

	file, err := os.Create(fileName)
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	assert.Nil(t, err)
	_, err = client.PutObjectFromFile(context.TODO(), putObjrequest, fileName)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	style := "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1"
	process := fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", style, strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(bucketName)), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(objectDestName)), "="))
	request := &AsyncProcessObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		AsyncProcess: Ptr(process),
	}
	var serr *ServiceError
	_, err = client.AsyncProcessObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, "Imm Client", serr.Code)
	assert.Contains(t, serr.Message, "ResourceNotFound, The specified resource Attachment is not found")
	assert.NotEmpty(t, serr.RequestID)

	time.Sleep(1 * time.Second)
	bucketNameNotExist := bucketName + "-not-exist"
	request = &AsyncProcessObjectRequest{
		Bucket:       Ptr(bucketNameNotExist),
		Key:          Ptr(objectName),
		AsyncProcess: Ptr(process),
	}
	_, err = client.AsyncProcessObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectWithProcess(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6) + ".jpg"

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}

	filePath := "../sample/example.jpg"
	_, err = client.PutObjectFromFile(context.TODO(), putObjRequest, filePath)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	style := "image/resize,m_fixed,w_100,h_100/rotate,90"
	getObjRequest := &GetObjectRequest{
		Bucket:  Ptr(bucketName),
		Key:     Ptr(objectName),
		Process: Ptr(style),
	}

	downloadFile := "example-download.jpg"
	_, err = client.GetObjectToFile(context.TODO(), getObjRequest, downloadFile)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	content, err := os.ReadFile(downloadFile)
	assert.Nil(t, err)

	result, err := client.GetObject(context.TODO(), getObjRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	content2, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, content2, content)

	sign, err := client.Presign(context.TODO(), getObjRequest)
	req, err := http.NewRequest(sign.Method, sign.URL, nil)
	assert.Nil(t, err)
	c := &http.Client{}
	resp, err := c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	time.Sleep(1 * time.Second)

	content3, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, content3, content)

	os.Remove(downloadFile)
}

func TestPutBucketRequestPayment(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	request := &PutBucketRequestPaymentRequest{
		Bucket: Ptr(bucketName),
		PaymentConfiguration: &RequestPaymentConfiguration{
			Payer: Requester,
		},
	}

	result, err := client.PutBucketRequestPayment(context.TODO(), request)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutBucketRequestPaymentRequest{
		Bucket: Ptr(bucketNameNotExist),
		PaymentConfiguration: &RequestPaymentConfiguration{
			Payer: Requester,
		},
	}
	result, err = client.PutBucketRequestPayment(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetBucketRequestPayment(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	request := &GetBucketRequestPaymentRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.GetBucketRequestPayment(context.TODO(), request)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.Payer, "BucketOwner")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &GetBucketRequestPaymentRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	result, err = client.GetBucketRequestPayment(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPaymentWithRequester(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	policyInfo := `
	{
		"Version":"1",
		"Statement":[
			{
				"Action":[
					"oss:*"
				],
				"Effect":"Allow",
				"Principal":["` + payerUID_ + `"],
				"Resource":["acs:oss:*:*:` + bucketName + `", "acs:oss:*:*:` + bucketName + `/*"]
			}
		]
	}`
	input := &OperationInput{
		OpName: "PutBucketPolicy",
		Bucket: Ptr(bucketName),
		Method: "PUT",
		Parameters: map[string]string{
			"policy": "",
		},
		Body: strings.NewReader(policyInfo),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	_, err = client.InvokeOperation(context.TODO(), input)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	request := &PutBucketRequestPaymentRequest{
		Bucket: Ptr(bucketName),
		PaymentConfiguration: &RequestPaymentConfiguration{
			Payer: Requester,
		},
	}
	_, err = client.PutBucketRequestPayment(context.TODO(), request)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	body := randStr(100)
	creClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider(payerAccessID_, payerAccessKey_))

	objectName := objectNamePrefix + randStr(6)

	putObjReq := &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(body),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObject(context.TODO(), putObjReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getObjReq := &GetObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	getObjResult, err := creClient.GetObject(context.TODO(), getObjReq)
	assert.Nil(t, err)
	getObjData, _ := io.ReadAll(getObjResult.Body)
	assert.Equal(t, string(getObjData), body)
	time.Sleep(1 * time.Second)

	objectCopyName := objectName + "-copy"
	copyRequest := &CopyObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectCopyName),
		SourceKey:    Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	objectAppendName := objectName + "-append"
	appendRequest := &AppendObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectAppendName),
		Body:         strings.NewReader(body),
		Position:     Ptr(int64(0)),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.AppendObject(context.TODO(), appendRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	delRequest := &DeleteObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	delObjsRequest := &DeleteMultipleObjectsRequest{
		Bucket:       Ptr(bucketName),
		Objects:      []DeleteObject{{Key: Ptr(objectAppendName)}},
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.DeleteMultipleObjects(context.TODO(), delObjsRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	headRequest := &HeadObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectCopyName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.HeadObject(context.TODO(), headRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	metaRequest := &GetObjectMetaRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectCopyName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.GetObjectMeta(context.TODO(), metaRequest)
	assert.Nil(t, err)

	objectRestoreName := objectName + "-restore"
	putObjReq = &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectRestoreName),
		Body:         strings.NewReader(body),
		StorageClass: StorageClassColdArchive,
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObject(context.TODO(), putObjReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	restoreRequest := &RestoreObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectRestoreName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.RestoreObject(context.TODO(), restoreRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	putObjReq = &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(body),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObject(context.TODO(), putObjReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	putAclRequest := &PutObjectAclRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Acl:          ObjectACLPublicRead,
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObjectAcl(context.TODO(), putAclRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getAclRequest := &GetObjectAclRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.GetObjectAcl(context.TODO(), getAclRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	objectMultiName := objectName + "-multi"
	body = randLowStr(360000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := io.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]
	initRequest := &InitiateMultipartUploadRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		RequestPayer: Ptr("requester"),
	}
	initResult, err := creClient.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	contents := []string{string(part1), string(part2), string(part3)}
	var parts []UploadPart
	var wg sync.WaitGroup
	wg.Add(len(contents))
	for i, content1 := range contents {
		partRequest := &UploadPartRequest{
			Bucket:       Ptr(bucketName),
			Key:          Ptr(objectMultiName),
			PartNumber:   int32(i + 1),
			UploadId:     Ptr(*initResult.UploadId),
			Body:         strings.NewReader(content1),
			RequestPayer: Ptr("requester"),
		}
		partResult, err := creClient.UploadPart(context.TODO(), partRequest)
		assert.Nil(t, err)

		part := UploadPart{
			PartNumber: partRequest.PartNumber,
			ETag:       partResult.ETag,
		}
		parts = append(parts, part)
		wg.Done()
	}

	comRequest := &CompleteMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectMultiName),
		UploadId: Ptr(*initResult.UploadId),
		CompleteMultipartUpload: &CompleteMultipartUpload{
			Parts: parts,
		},
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.CompleteMultipartUpload(context.TODO(), comRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	initRequest = &InitiateMultipartUploadRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		RequestPayer: Ptr("requester"),
	}
	initResult, err = creClient.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	copyMultiRequest := &UploadPartCopyRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		PartNumber:   int32(1),
		UploadId:     Ptr(*initResult.UploadId),
		SourceKey:    Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.UploadPartCopy(context.TODO(), copyMultiRequest)
	assert.Nil(t, err)

	listMultiRequest := &ListMultipartUploadsRequest{
		Bucket:       Ptr(bucketName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListMultipartUploads(context.TODO(), listMultiRequest)

	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	listRequest := &ListPartsRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		UploadId:     Ptr(*initResult.UploadId),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListParts(context.TODO(), listRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		UploadId:     Ptr(*initResult.UploadId),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	symlinkName := objectName + "-symlink"
	putSymRequest := &PutSymlinkRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(symlinkName),
		Target:       Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutSymlink(context.TODO(), putSymRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getSymRequest := &GetSymlinkRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(symlinkName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.GetSymlink(context.TODO(), getSymRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	putTagRequest := &PutObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObjectTagging(context.TODO(), putTagRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getTagRequest := &GetObjectTaggingRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.GetObjectTagging(context.TODO(), getTagRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	delTagRequest := &DeleteObjectTaggingRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.DeleteObjectTagging(context.TODO(), delTagRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	listObjReq := &ListObjectsRequest{
		Bucket:       Ptr(bucketName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListObjects(context.TODO(), listObjReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	listObjReqV2 := &ListObjectsV2Request{
		Bucket:       Ptr(bucketName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListObjectsV2(context.TODO(), listObjReqV2)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	listObjVersionReq := &ListObjectVersionsRequest{
		Bucket:       Ptr(bucketName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListObjectVersions(context.TODO(), listObjVersionReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	putObjReq = &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(body),
		RequestPayer: Ptr("bucketOwner"),
	}
	_, err = creClient.PutObject(context.TODO(), putObjReq)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Access denied for requester pay bucket")
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "AccessDenied", serr.Code)
	assert.Equal(t, "Access denied for requester pay bucket", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}

func TestClientExtensionWithPayer(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)
	client := getDefaultClient()
	assert.NotNil(t, client)

	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	policyInfo := `
	{
		"Version":"1",
		"Statement":[
			{
				"Action":[
					"oss:*"
				],
				"Effect":"Allow",
				"Principal":["` + payerUID_ + `"],
				"Resource":["acs:oss:*:*:` + bucketName + `", "acs:oss:*:*:` + bucketName + `/*"]
			}
		]
	}`
	input := &OperationInput{
		OpName: "PutBucketPolicy",
		Bucket: Ptr(bucketName),
		Method: "PUT",
		Parameters: map[string]string{
			"policy": "",
		},
		Body: strings.NewReader(policyInfo),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	_, err = client.InvokeOperation(context.TODO(), input)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	request := &PutBucketRequestPaymentRequest{
		Bucket: Ptr(bucketName),
		PaymentConfiguration: &RequestPaymentConfiguration{
			Payer: Requester,
		},
	}
	_, err = client.PutBucketRequestPayment(context.TODO(), request)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	creClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider(payerAccessID_, payerAccessKey_))

	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)

	// IsObjectExist
	exist, err := creClient.IsObjectExist(context.TODO(), bucketName, objectName, func(op *IsObjectExistOptions) {
		op.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	assert.True(t, exist)

	//PutObjectFromFile
	objectNameFromFile := objectName + "-from-file"
	var localFile = randStr(8) + ".txt"
	length := 1234
	content := randStr(length)
	hashContent := NewCRC64(0)
	hashContent.Write([]byte(content))
	createFile(t, localFile, content)
	defer func() { os.Remove(localFile) }()

	result, err := creClient.PutObjectFromFile(context.TODO(), &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectNameFromFile),
		RequestPayer: Ptr("requester"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Use Uploader, set meta and acl
	objectNameBig := objectName + "-big"
	bigLength := 5*100*1024 + 1234
	bigContent := randStr(bigLength)
	bigHash := NewCRC64(0)
	bigHash.Write([]byte(bigContent))
	u := creClient.NewUploader()
	assert.NotNil(t, u)
	urResult, err := u.UploadFrom(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectNameBig),
			Metadata: map[string]string{
				"author": "test",
				"magic":  "123",
			},
			Acl:          ObjectACLPublicRead,
			RequestPayer: Ptr("requester"),
		},
		bytes.NewReader([]byte(bigContent)),
		func(uo *UploaderOptions) {
			uo.ParallelNum = 3
			uo.PartSize = 100 * 1024
		},
	)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.NotNil(t, urResult)

	exist, err = creClient.IsObjectExist(context.TODO(), bucketName, objectNameBig, func(op *IsObjectExistOptions) {
		op.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	assert.True(t, exist)

	// Downloader with not align partSize
	d := creClient.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 100*1024 + 123
	})
	assert.NotNil(t, d)
	assert.Equal(t, int64(100*1024+123), d.options.PartSize)
	assert.Equal(t, 3, d.options.ParallelNum)
	localFileBig := randStr(8) + "-downloader"
	dResult, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket:       Ptr(bucketName),
			Key:          Ptr(objectNameBig),
			RequestPayer: Ptr("requester"),
		},
		localFileBig)

	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, int64(bigLength), dResult.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFileBig)
	assert.Nil(t, err)
	defer func() {
		rfile.Close()
		os.Remove(localFileBig)
	}()
	io.Copy(hash, rfile)
	assert.Equal(t, bigHash.Sum64(), hash.Sum64())

	//Use ReadOnlyFile
	f, err := creClient.OpenFile(context.TODO(), bucketName, objectNameBig, func(op *OpenOptions) {
		op.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	for i := 13; i < 42; i++ {
		for len := 100*1024 + 123; len < 100*1024+123+17; len++ {
			_, err := f.Seek(int64(i), io.SeekStart)
			assert.Nil(t, err)
			gData, err := io.ReadAll(io.LimitReader(f, int64(len)))
			assert.Nil(t, err)
			assert.EqualValues(t, []byte(bigContent)[i:i+len], gData)
		}
	}
	f.Close()

	// AppenableFile
	objectNameAppend := objectName + "-append"
	dataa3 := []byte(randStr(100*1024*5 + 13))
	var localFileData3 = randStr(8) + ".txt"
	createFile(t, localFileData3, string(dataa3))
	defer func() {
		os.Remove(localFileData3)
	}()

	af, err := creClient.AppendFile(context.TODO(), bucketName, objectNameAppend, func(op *AppendOptions) {
		op.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	_, err = af.Write([]byte(content))
	assert.Nil(t, err)
	_, err = af.WriteFrom(strings.NewReader(content))
	assert.Nil(t, err)
	_, err = af.Stat()
	assert.Nil(t, err)
	//GetObjectToFile
	var localFileToFile = randStr(8) + "-to-file"
	defer func() {
		os.Remove(localFileToFile)
	}()
	_, err = creClient.GetObjectToFile(context.TODO(), &GetObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectNameAppend),
		RequestPayer: Ptr("requester"),
	},
		localFileToFile,
	)
	assert.Nil(t, err)
	_, err = creClient.GetObjectToFile(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend),
	},
		localFileToFile,
	)
	assert.NotNil(t, err)
}

func TestClientAppendFile(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6) + ".append"

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	f, err := client.AppendFile(context.TODO(), bucketName, objectName, func(ao *AppendOptions) {
		ao.CreateParameter = &AppendObjectRequest{
			Acl:          ObjectACLPublicRead,
			CacheControl: Ptr("no-cache"),
			Metadata: map[string]string{
				"user": "jack",
			},
			Tagging: Ptr("key=value"),
		}
	})
	assert.Nil(t, err)

	n, err := f.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, int(5), n)

	n, err = f.Write([]byte(" world"))
	assert.Nil(t, err)
	assert.Equal(t, int(6), n)

	nn, err := f.WriteFrom(strings.NewReader(" 123"))
	assert.Nil(t, err)
	assert.Equal(t, int64(4), nn)

	info, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(15), info.Size())
	header, ok := info.Sys().(http.Header)
	assert.True(t, ok)

	assert.Equal(t, "no-cache", header.Get("Cache-Control"))
	assert.Equal(t, "jack", header.Get("x-oss-meta-user"))
	assert.Equal(t, "1", header.Get("x-oss-tagging-count"))

	f.Close()

	//Open Again
	f, err = client.AppendFile(context.TODO(), bucketName, objectName, func(ao *AppendOptions) {
		ao.CreateParameter = &AppendObjectRequest{
			Acl:          ObjectACLPublicRead,
			CacheControl: Ptr("no-cache"),
			Metadata: map[string]string{
				"user": "jack",
			},
			Tagging: Ptr("key=value"),
		}
	})
	n, err = f.Write([]byte("abc"))
	assert.Nil(t, err)
	info, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(18), info.Size())
	f.Close()

	objectName1 := objectNamePrefix + randLowStr(6) + "-1-.append"
	f, err = client.AppendFile(context.TODO(), bucketName, objectName1, func(ao *AppendOptions) {
		ao.CreateParameter = &AppendObjectRequest{
			Acl:          ObjectACLPublicRead,
			CacheControl: Ptr("no-cache"),
			Metadata: map[string]string{
				"user": "jack-1",
			},
			Tagging: Ptr("key1=value1"),
		}
	})
	assert.Nil(t, err)

	nn, err = f.WriteFrom(strings.NewReader("123"))
	assert.Nil(t, err)
	assert.Equal(t, int64(3), nn)

	nn, err = f.WriteFrom(strings.NewReader("-abc-321"))
	assert.Nil(t, err)
	assert.Equal(t, int64(8), nn)
	info, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(11), info.Size())
	header, ok = info.Sys().(http.Header)
	assert.True(t, ok)

	assert.Equal(t, "no-cache", header.Get("Cache-Control"))
	assert.Equal(t, "jack-1", header.Get("x-oss-meta-user"))
	assert.Equal(t, "1", header.Get("x-oss-tagging-count"))
	f.Close()
}

func TestServiceError(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.DeleteBucket(context.TODO(), putRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The specified bucket does not exist.")
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
	assert.NotNil(t, serr.Headers)
	assert.NotEmpty(t, serr.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, serr.Headers.Get("server"))
}

func TestBucketLogging(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	request := &PutBucketLoggingRequest{
		Bucket: Ptr(bucketName),
		BucketLoggingStatus: &BucketLoggingStatus{
			&LoggingEnabled{
				TargetBucket: Ptr(bucketName),
				TargetPrefix: Ptr("TargetPrefix"),
			},
		},
	}
	result, err := client.PutBucketLogging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketLoggingRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketLogging(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *getResult.BucketLoggingStatus.LoggingEnabled.TargetBucket, bucketName)
	assert.Equal(t, *getResult.BucketLoggingStatus.LoggingEnabled.TargetPrefix, "TargetPrefix")
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketLoggingRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketLogging(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.Equal(t, "204 No Content", delResult.Status)
	assert.NotEmpty(t, delResult.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, delResult.Headers.Get("Date"))
	time.Sleep(1 * time.Second)

	putUserRequest := &PutUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketName),
		UserDefinedLogFieldsConfiguration: &UserDefinedLogFieldsConfiguration{
			HeaderSet: &LoggingHeaderSet{
				[]string{"header1", "header2", "header3"},
			},
			ParamSet: &LoggingParamSet{
				[]string{"param1", "param2"},
			},
		},
	}
	putUserResult, err := client.PutUserDefinedLogFieldsConfig(context.TODO(), putUserRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putUserResult.StatusCode)
	assert.NotEmpty(t, putUserResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getUserRequest := &GetUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketName),
	}
	getUserResult, err := client.GetUserDefinedLogFieldsConfig(context.TODO(), getUserRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getUserResult.StatusCode)
	assert.NotEmpty(t, getUserResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, 3, len(getUserResult.UserDefinedLogFieldsConfiguration.HeaderSet.Headers))
	assert.Equal(t, 2, len(getUserResult.UserDefinedLogFieldsConfiguration.ParamSet.Parameters))
	time.Sleep(1 * time.Second)

	delUserRequest := &DeleteUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketName),
	}
	delUserResult, err := client.DeleteUserDefinedLogFieldsConfig(context.TODO(), delUserRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delUserResult.StatusCode)
	assert.Equal(t, "204 No Content", delUserResult.Status)
	assert.NotEmpty(t, delUserResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutBucketLoggingRequest{
		Bucket: Ptr(bucketNameNotExist),
		BucketLoggingStatus: &BucketLoggingStatus{
			&LoggingEnabled{
				TargetBucket: Ptr("TargetBucket"),
				TargetPrefix: Ptr("TargetPrefix"),
			},
		},
	}
	result, err = client.PutBucketLogging(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getRequest = &GetBucketLoggingRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	getResult, err = client.GetBucketLogging(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketLoggingRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	delResult, err = client.DeleteBucketLogging(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	putUserRequest = &PutUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
		UserDefinedLogFieldsConfiguration: &UserDefinedLogFieldsConfiguration{
			HeaderSet: &LoggingHeaderSet{
				[]string{"header1", "header2", "header3"},
			},
			ParamSet: &LoggingParamSet{
				[]string{"param1", "param2"},
			},
		},
	}
	serr = &ServiceError{}
	putUserResult, err = client.PutUserDefinedLogFieldsConfig(context.TODO(), putUserRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getUserRequest = &GetUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	getUserResult, err = client.GetUserDefinedLogFieldsConfig(context.TODO(), getUserRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delUserRequest = &DeleteUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	delUserResult, err = client.DeleteUserDefinedLogFieldsConfig(context.TODO(), delUserRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketWorm(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	initWorm := &InitiateBucketWormRequest{
		Bucket: Ptr(bucketName),
		InitiateWormConfiguration: &InitiateWormConfiguration{
			Ptr(int32(1)),
		},
	}

	initResult, err := client.InitiateBucketWorm(context.TODO(), initWorm)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketWormRequest{
		Bucket: Ptr(bucketName),
	}

	getResult, err := client.GetBucketWorm(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *getResult.WormConfiguration.WormId, *initResult.WormId)
	assert.NotEmpty(t, *getResult.WormConfiguration.CreationDate)
	assert.NotEmpty(t, getResult.WormConfiguration.RetentionPeriodInDays)
	assert.NotEmpty(t, getResult.WormConfiguration.State)
	time.Sleep(1 * time.Second)

	abortRequest := &AbortBucketWormRequest{
		Bucket: Ptr(bucketName),
	}
	abortResult, err := client.AbortBucketWorm(context.TODO(), abortRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, abortResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	completeRequest := &CompleteBucketWormRequest{
		Bucket: Ptr(bucketName),
		WormId: initResult.WormId,
	}
	_, err = client.CompleteBucketWorm(context.TODO(), completeRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	extendRequest := &ExtendBucketWormRequest{
		Bucket: Ptr(bucketName),
		WormId: initResult.WormId,
	}
	_, err = client.ExtendBucketWorm(context.TODO(), extendRequest)
	assert.NotNil(t, err)
	assert.Equal(t, "missing required field, ExtendWormConfiguration.", err.Error())
	time.Sleep(1 * time.Second)

	extendRequest = &ExtendBucketWormRequest{
		Bucket: Ptr(bucketName),
		WormId: initResult.WormId,
		ExtendWormConfiguration: &ExtendWormConfiguration{
			Ptr(int32(2)),
		},
	}
	serr = &ServiceError{}
	_, err = client.ExtendBucketWorm(context.TODO(), extendRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketPolicy(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketPolicyRequest{
		Bucket: Ptr(bucketName),
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

	putResult, err := client.PutBucketPolicy(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketPolicyRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketPolicy(context.TODO(), getRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, getResult.Body)

	statusRequest := &GetBucketPolicyStatusRequest{
		Bucket: Ptr(bucketName),
	}
	statusResult, err := client.GetBucketPolicyStatus(context.TODO(), statusRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, statusResult.StatusCode)
	assert.NotEmpty(t, statusResult.Headers.Get("X-Oss-Request-Id"))
	assert.False(t, *statusResult.PolicyStatus.IsPublic)

	delRequest := &DeleteBucketPolicyRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketPolicy(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketPolicyRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketPolicy(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	statusRequest = &GetBucketPolicyStatusRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	statusResult, err = client.GetBucketPolicyStatus(context.TODO(), statusRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketPolicyRequest{
		Bucket: Ptr(bucketNameNotExist),
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
	serr = &ServiceError{}
	putResult, err = client.PutBucketPolicy(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest = &DeleteBucketPolicyRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	delResult, err = client.DeleteBucketPolicy(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketTransferAcceleration(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketTransferAccelerationRequest{
		Bucket: Ptr(bucketName),
		TransferAccelerationConfiguration: &TransferAccelerationConfiguration{
			Ptr(true),
		},
	}
	putResult, err := client.PutBucketTransferAcceleration(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketTransferAccelerationRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketTransferAcceleration(context.TODO(), getRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.True(t, *getResult.TransferAccelerationConfiguration.Enabled)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketTransferAccelerationRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketTransferAcceleration(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketTransferAccelerationRequest{
		Bucket: Ptr(bucketNameNotExist),
		TransferAccelerationConfiguration: &TransferAccelerationConfiguration{
			Ptr(true),
		},
	}
	serr = &ServiceError{}
	putResult, err = client.PutBucketTransferAcceleration(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketArchiveDirectRead(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketArchiveDirectReadRequest{
		Bucket: Ptr(bucketName),
		ArchiveDirectReadConfiguration: &ArchiveDirectReadConfiguration{
			Ptr(true),
		},
	}
	putResult, err := client.PutBucketArchiveDirectRead(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketArchiveDirectReadRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketArchiveDirectRead(context.TODO(), getRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.True(t, *getResult.ArchiveDirectReadConfiguration.Enabled)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketArchiveDirectReadRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketArchiveDirectRead(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketArchiveDirectReadRequest{
		Bucket: Ptr(bucketNameNotExist),
		ArchiveDirectReadConfiguration: &ArchiveDirectReadConfiguration{
			Ptr(true),
		},
	}
	serr = &ServiceError{}
	putResult, err = client.PutBucketArchiveDirectRead(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketWebsite(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketWebsiteRequest{
		Bucket: Ptr(bucketName),
		WebsiteConfiguration: &WebsiteConfiguration{
			IndexDocument: &IndexDocument{
				Suffix:        Ptr("index.html"),
				SupportSubDir: Ptr(true),
				Type:          Ptr(int64(0)),
			},
			ErrorDocument: &ErrorDocument{
				Key:        Ptr("error.html"),
				HttpStatus: Ptr(int64(404)),
			},
		},
	}
	putResult, err := client.PutBucketWebsite(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketWebsiteRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketWebsite(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, getResult.WebsiteConfiguration, putRequest.WebsiteConfiguration)
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketWebsiteRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketWebsite(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketWebsiteRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketWebsite(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketWebsiteRequest{
		Bucket: Ptr(bucketNameNotExist),
		WebsiteConfiguration: &WebsiteConfiguration{
			IndexDocument: &IndexDocument{
				Suffix:        Ptr("index.html"),
				SupportSubDir: Ptr(true),
				Type:          Ptr(int64(0)),
			},
			ErrorDocument: &ErrorDocument{
				Key:        Ptr("error.html"),
				HttpStatus: Ptr(int64(404)),
			},
		},
	}
	serr = &ServiceError{}
	putResult, err = client.PutBucketWebsite(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketWebsiteRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	delResult, err = client.DeleteBucketWebsite(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketHttpsConfig(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketHttpsConfigRequest{
		Bucket: Ptr(bucketName),
		HttpsConfiguration: &HttpsConfiguration{
			TLS: &TLS{
				Enable:      Ptr(true),
				TLSVersions: []string{"TLSv1.2", "TLSv1.3"},
			},
		},
	}
	putResult, err := client.PutBucketHttpsConfig(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketHttpsConfigRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketHttpsConfig(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.True(t, *getResult.HttpsConfiguration.TLS.Enable)
	assert.Equal(t, len(getResult.HttpsConfiguration.TLS.TLSVersions), 2)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketHttpsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketHttpsConfig(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketHttpsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
		HttpsConfiguration: &HttpsConfiguration{
			TLS: &TLS{
				Enable:      Ptr(true),
				TLSVersions: []string{"TLSv1.2", "TLSv1.3"},
			},
		},
	}
	putResult, err = client.PutBucketHttpsConfig(context.TODO(), putRequest)
	serr = &ServiceError{}
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketResourceGroup(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	getRequest := &GetBucketResourceGroupRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketResourceGroup(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *getResult.BucketResourceGroupConfiguration.ResourceGroupId)
	time.Sleep(1 * time.Second)

	putRequest := &PutBucketResourceGroupRequest{
		Bucket: Ptr(bucketName),
		BucketResourceGroupConfiguration: &BucketResourceGroupConfiguration{
			getResult.BucketResourceGroupConfiguration.ResourceGroupId,
		},
	}
	putResult, err := client.PutBucketResourceGroup(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketResourceGroupRequest{
		Bucket: Ptr(bucketName),
		BucketResourceGroupConfiguration: &BucketResourceGroupConfiguration{
			Ptr("rg-not-exist"),
		},
	}
	putResult, err = client.PutBucketResourceGroup(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "ResourceGroupIdPreCheckError", serr.Code)
	assert.Equal(t, "The resource group id precheck error", serr.Message)
	assert.Equal(t, "0039-00000003", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getRequest = &GetBucketResourceGroupRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketResourceGroup(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
}

func TestBucketTags(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketTagsRequest{
		Bucket: Ptr(bucketName),
		Tagging: &Tagging{
			&TagSet{
				[]Tag{
					{
						Ptr("key1"),
						Ptr("value1"),
					},
					{
						Ptr("key2"),
						Ptr("value2"),
					},
					{
						Ptr("key3"),
						Ptr("value3"),
					},
				},
			},
		},
	}
	putResult, err := client.PutBucketTags(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketTagsRequest{
		Bucket:  Ptr(bucketName),
		Tagging: Ptr("key1,key3"),
	}
	delResult, err := client.DeleteBucketTags(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketTagsRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketTags(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, len(getResult.Tagging.TagSet.Tags), 1)
	assert.Equal(t, *getResult.Tagging.TagSet.Tags[0].Key, "key2")
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketTagsRequest{
		Bucket: Ptr(bucketNameNotExist),
		Tagging: &Tagging{
			&TagSet{
				[]Tag{
					{
						Ptr("key1"),
						Ptr("value1"),
					},
					{
						Ptr("key2"),
						Ptr("value2"),
					},
					{
						Ptr("key3"),
						Ptr("value3"),
					},
				},
			},
		},
	}
	putResult, err = client.PutBucketTags(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketTagsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketTags(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest = &DeleteBucketTagsRequest{
		Bucket:  Ptr(bucketNameNotExist),
		Tagging: Ptr("key1,key3"),
	}
	delResult, err = client.DeleteBucketTags(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketEncryption(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketEncryptionRequest{
		Bucket: Ptr(bucketName),
		ServerSideEncryptionRule: &ServerSideEncryptionRule{
			&ApplyServerSideEncryptionByDefault{
				SSEAlgorithm:      Ptr("KMS"),
				KMSDataEncryption: Ptr("SM4"),
			},
		},
	}
	putResult, err := client.PutBucketEncryption(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketEncryptionRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketEncryption(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm, "KMS")
	assert.Equal(t, *getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.KMSDataEncryption, "SM4")
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketEncryptionRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketEncryption(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketEncryptionRequest{
		Bucket: Ptr(bucketNameNotExist),
		ServerSideEncryptionRule: &ServerSideEncryptionRule{
			&ApplyServerSideEncryptionByDefault{
				SSEAlgorithm:      Ptr("KMS"),
				KMSDataEncryption: Ptr("SM4"),
			},
		},
	}
	putResult, err = client.PutBucketEncryption(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketEncryptionRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketEncryption(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest = &DeleteBucketEncryptionRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	delResult, err = client.DeleteBucketEncryption(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketReferer(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketRefererRequest{
		Bucket: Ptr(bucketName),
		RefererConfiguration: &RefererConfiguration{
			AllowEmptyReferer: Ptr(false),
			RefererList: &RefererList{
				Referers: []string{""},
			},
		},
	}
	putResult, err := client.PutBucketReferer(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketRefererRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketReferer(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.False(t, *getResult.RefererConfiguration.AllowEmptyReferer)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketRefererRequest{
		Bucket: Ptr(bucketNameNotExist),
		RefererConfiguration: &RefererConfiguration{
			AllowEmptyReferer: Ptr(false),
		},
	}
	putResult, err = client.PutBucketReferer(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketRefererRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketReferer(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketInventory(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	id := "report1" + randStr(6)
	putRequest := &PutBucketInventoryRequest{
		Bucket:      Ptr(bucketName),
		InventoryId: Ptr(id),
		InventoryConfiguration: &InventoryConfiguration{
			Id:        Ptr(id),
			IsEnabled: Ptr(true),
			Filter: &InventoryFilter{
				Prefix:                   Ptr("filterPrefix"),
				LastModifyBeginTimeStamp: Ptr(int64(1637883649)),
				LastModifyEndTimeStamp:   Ptr(int64(1638347592)),
				LowerSizeBound:           Ptr(int64(1024)),
				UpperSizeBound:           Ptr(int64(1048576)),
				StorageClass:             Ptr("Standard,IA"),
			},
			Destination: &InventoryDestination{
				&InventoryOSSBucketDestination{
					Format:    InventoryFormatCSV,
					AccountId: Ptr(accountID_),
					RoleArn:   Ptr("acs:ram::" + accountID_ + ":role/AliyunOSSRole"),
					Bucket:    Ptr("acs:oss:::" + bucketName),
					Prefix:    Ptr("prefix1"),
				},
			},
			Schedule: &InventorySchedule{
				InventoryFrequencyDaily,
			},
			IncludedObjectVersions: Ptr("All"),
			OptionalFields: &OptionalFields{
				Fields: []InventoryOptionalFieldType{
					InventoryOptionalFieldSize,
					InventoryOptionalFieldLastModifiedDate,
					InventoryOptionalFieldETag,
					InventoryOptionalFieldStorageClass,
					InventoryOptionalFieldIsMultipartUploaded,
					InventoryOptionalFieldEncryptionStatus,
				},
			},
		},
	}
	putResult, err := client.PutBucketInventory(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketInventoryRequest{
		Bucket:      Ptr(bucketName),
		InventoryId: Ptr(id),
	}
	getResult, err := client.GetBucketInventory(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listRequest := &ListBucketInventoryRequest{
		Bucket: Ptr(bucketName),
	}
	listResult, err := client.ListBucketInventory(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, len(listResult.ListInventoryConfigurationsResult.InventoryConfigurations), 1)
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketInventoryRequest{
		Bucket:      Ptr(bucketName),
		InventoryId: Ptr(id),
	}
	delResult, err := client.DeleteBucketInventory(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketInventoryRequest{
		Bucket:      Ptr(bucketNameNotExist),
		InventoryId: Ptr(id),
		InventoryConfiguration: &InventoryConfiguration{
			Id:        Ptr(id),
			IsEnabled: Ptr(true),
			Filter: &InventoryFilter{
				Prefix:                   Ptr("filterPrefix"),
				LastModifyBeginTimeStamp: Ptr(int64(1637883649)),
				LastModifyEndTimeStamp:   Ptr(int64(1638347592)),
				LowerSizeBound:           Ptr(int64(1024)),
				UpperSizeBound:           Ptr(int64(1048576)),
				StorageClass:             Ptr("Standard,IA"),
			},
			Destination: &InventoryDestination{
				&InventoryOSSBucketDestination{
					Format:    InventoryFormatCSV,
					AccountId: Ptr(accountID_),
					RoleArn:   Ptr("acs:ram::" + accountID_ + ":role/AliyunOSSRole"),
					Bucket:    Ptr("acs:oss:::" + bucketName),
					Prefix:    Ptr("prefix1"),
				},
			},
			Schedule: &InventorySchedule{
				InventoryFrequencyDaily,
			},
			IncludedObjectVersions: Ptr("All"),
			OptionalFields: &OptionalFields{
				Fields: []InventoryOptionalFieldType{
					InventoryOptionalFieldSize,
					InventoryOptionalFieldLastModifiedDate,
					InventoryOptionalFieldETag,
					InventoryOptionalFieldStorageClass,
					InventoryOptionalFieldIsMultipartUploaded,
					InventoryOptionalFieldEncryptionStatus,
				},
			},
		},
	}
	putResult, err = client.PutBucketInventory(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketInventoryRequest{
		Bucket:      Ptr(bucketNameNotExist),
		InventoryId: Ptr(id),
	}
	getResult, err = client.GetBucketInventory(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	listRequest = &ListBucketInventoryRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	listResult, err = client.ListBucketInventory(context.TODO(), listRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketInventoryRequest{
		Bucket:      Ptr(bucketNameNotExist),
		InventoryId: Ptr(id),
	}
	delResult, err = client.DeleteBucketInventory(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketAccessMonitor(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketAccessMonitorRequest{
		Bucket: Ptr(bucketName),
		AccessMonitorConfiguration: &AccessMonitorConfiguration{
			Status: AccessMonitorStatusEnabled,
		},
	}
	putResult, err := client.PutBucketAccessMonitor(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketAccessMonitorRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketAccessMonitor(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, getResult.AccessMonitorConfiguration.Status, AccessMonitorStatusEnabled)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketAccessMonitorRequest{
		Bucket: Ptr(bucketNameNotExist),
		AccessMonitorConfiguration: &AccessMonitorConfiguration{
			Status: AccessMonitorStatusEnabled,
		},
	}
	putResult, err = client.PutBucketAccessMonitor(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketAccessMonitorRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketAccessMonitor(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketStyle(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutStyleRequest{
		Bucket:    Ptr(bucketName),
		StyleName: Ptr("demo"),
		Style: &StyleContent{
			Content: Ptr("image/resize,p_50"),
		},
	}
	putResult, err := client.PutStyle(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetStyleRequest{
		Bucket:    Ptr(bucketName),
		StyleName: Ptr("demo"),
	}
	getResult, err := client.GetStyle(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listRequest := &ListStyleRequest{
		Bucket: Ptr(bucketName),
	}
	listResult, err := client.ListStyle(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, len(listResult.StyleList.Styles), 1)
	time.Sleep(1 * time.Second)

	delRequest := &DeleteStyleRequest{
		Bucket:    Ptr(bucketName),
		StyleName: Ptr("demo"),
	}
	delResult, err := client.DeleteStyle(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutStyleRequest{
		Bucket:    Ptr(bucketNameNotExist),
		StyleName: Ptr("demo"),
		Style: &StyleContent{
			Content: Ptr("image/resize,p_50"),
		},
	}
	putResult, err = client.PutStyle(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetStyleRequest{
		Bucket:    Ptr(bucketNameNotExist),
		StyleName: Ptr("demo"),
	}
	getResult, err = client.GetStyle(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	listRequest = &ListStyleRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	listResult, err = client.ListStyle(context.TODO(), listRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteStyleRequest{
		Bucket:    Ptr(bucketNameNotExist),
		StyleName: Ptr("demo"),
	}
	delResult, err = client.DeleteStyle(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketReplication(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	targetBucketName := bucketNamePrefix + "-target-" + randLowStr(6)
	request = &PutBucketRequest{
		Bucket: Ptr(targetBucketName),
	}
	client1 := getClient("cn-beijing", "https://oss-cn-beijing.aliyuncs.com")
	_, err = client1.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketReplicationRequest{
		Bucket: Ptr(bucketName),
		ReplicationConfiguration: &ReplicationConfiguration{
			[]ReplicationRule{
				{
					RTC: &ReplicationTimeControl{
						Status: Ptr("enabled"),
					},
					Destination: &ReplicationDestination{
						Bucket:       Ptr(targetBucketName),
						Location:     Ptr("oss-cn-beijing"),
						TransferType: TransferTypeInternal,
					},
					HistoricalObjectReplication: HistoricalObjectReplicationEnabled,
					SourceSelectionCriteria: &ReplicationSourceSelectionCriteria{
						&SseKmsEncryptedObjects{
							Status: StatusEnabled,
						},
					},
				},
			},
		},
	}
	putResult, err := client.PutBucketReplication(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketReplicationRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketReplication(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.Equal(t, 1, len(getResult.ReplicationConfiguration.Rules))
	time.Sleep(1 * time.Second)

	getLocationRequest := &GetBucketReplicationLocationRequest{
		Bucket: Ptr(bucketName),
	}
	getLocationResult, err := client.GetBucketReplicationLocation(context.TODO(), getLocationRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getLocationResult.StatusCode)
	time.Sleep(1 * time.Second)

	getProgressRequest := &GetBucketReplicationProgressRequest{
		Bucket: Ptr(bucketName),
		RuleId: getResult.ReplicationConfiguration.Rules[0].ID,
	}
	getProgressResult, err := client.GetBucketReplicationProgress(context.TODO(), getProgressRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getProgressResult.StatusCode)
	assert.Equal(t, 1, len(getProgressResult.ReplicationProgress.Rules))
	time.Sleep(1 * time.Second)

	rtcRequest := &PutBucketRtcRequest{
		Bucket: Ptr(bucketName),
		RtcConfiguration: &RtcConfiguration{
			RTC: &ReplicationTimeControl{
				Status: Ptr("disabled"),
			},
			ID: getResult.ReplicationConfiguration.Rules[0].ID,
		},
	}
	rtcResult, err := client.PutBucketRtc(context.TODO(), rtcRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, rtcResult.StatusCode)
	assert.NotEmpty(t, rtcResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketReplicationRequest{
		Bucket: Ptr(bucketName),
		ReplicationRules: &ReplicationRules{
			[]string{*getResult.ReplicationConfiguration.Rules[0].ID},
		},
	}
	delResult, err := client.DeleteBucketReplication(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketReplicationRequest{
		Bucket: Ptr(bucketNameNotExist),
		ReplicationConfiguration: &ReplicationConfiguration{
			[]ReplicationRule{
				{
					RTC: &ReplicationTimeControl{
						Status: Ptr("enabled"),
					},
					Destination: &ReplicationDestination{
						Bucket:       Ptr(targetBucketName),
						Location:     Ptr("oss-cn-cn-hangzhou"),
						TransferType: TransferTypeOssAcc,
					},
					HistoricalObjectReplication: HistoricalObjectReplicationEnabled,
					SourceSelectionCriteria: &ReplicationSourceSelectionCriteria{
						&SseKmsEncryptedObjects{
							Status: StatusEnabled,
						},
					},
				},
			},
		},
	}
	putResult, err = client.PutBucketReplication(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketReplicationRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	_, err = client.GetBucketReplication(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getLocationRequest = &GetBucketReplicationLocationRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getLocationResult, err = client.GetBucketReplicationLocation(context.TODO(), getLocationRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getProgressRequest = &GetBucketReplicationProgressRequest{
		Bucket: Ptr(bucketNameNotExist),
		RuleId: getResult.ReplicationConfiguration.Rules[0].ID,
	}
	getProgressResult, err = client.GetBucketReplicationProgress(context.TODO(), getProgressRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	rtcRequest = &PutBucketRtcRequest{
		Bucket: Ptr(bucketNameNotExist),
		RtcConfiguration: &RtcConfiguration{
			RTC: &ReplicationTimeControl{
				Status: Ptr("disabled"),
			},
			ID: getResult.ReplicationConfiguration.Rules[0].ID,
		},
	}
	rtcResult, err = client.PutBucketRtc(context.TODO(), rtcRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketReplicationRequest{
		Bucket: Ptr(bucketNameNotExist),
		ReplicationRules: &ReplicationRules{
			[]string{*getResult.ReplicationConfiguration.Rules[0].ID},
		},
	}
	delResult, err = client.DeleteBucketReplication(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketCors(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketCorsRequest{
		Bucket: Ptr(bucketName),
		CORSConfiguration: &CORSConfiguration{
			CORSRules: []CORSRule{
				{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"PUT", "GET"},
				},
			},
		},
	}
	putResult, err := client.PutBucketCors(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	key := objectNamePrefix + randStr(6)
	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(key),
		Body:   strings.NewReader("hi oss"),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	optionRequest := &OptionObjectRequest{
		Bucket:                     Ptr(bucketName),
		Key:                        Ptr(key),
		Origin:                     Ptr("http://www.example.com"),
		AccessControlRequestMethod: Ptr("PUT"),
	}
	optionResult, err := client.OptionObject(context.TODO(), optionRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, optionResult.StatusCode)
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketCorsRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketCors(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.Equal(t, 1, len(getResult.CORSConfiguration.CORSRules))
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketCorsRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketCors(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketCorsRequest{
		Bucket: Ptr(bucketNameNotExist),
		CORSConfiguration: &CORSConfiguration{
			CORSRules: []CORSRule{
				{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"PUT", "GET"},
				},
			},
		},
	}
	putResult, err = client.PutBucketCors(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketCorsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketCors(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketCorsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	delResult, err = client.DeleteBucketCors(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	optionRequest = &OptionObjectRequest{
		Bucket:                     Ptr(bucketNameNotExist),
		Key:                        Ptr(key),
		Origin:                     Ptr("http://www.example.com"),
		AccessControlRequestMethod: Ptr("PUT"),
	}
	optionResult, err = client.OptionObject(context.TODO(), optionRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketLifecycle(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	createRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), createRequest)
	assert.Nil(t, err)

	putRequest := &PutBucketLifecycleRequest{
		Bucket: Ptr(bucketName),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					Status: Ptr("Enabled"),
					ID:     Ptr("rule"),
					Prefix: Ptr("log/"),
					Transitions: []LifecycleRuleTransition{
						{
							Days:         Ptr(int32(30)),
							StorageClass: StorageClassIA,
						},
					},
				},
			},
		},
	}

	putResult, err := client.PutBucketLifecycle(context.TODO(), putRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))

	getRequest := &GetBucketLifecycleRequest{
		Bucket: Ptr(bucketName),
	}

	result, err := client.GetBucketLifecycle(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.LifecycleConfiguration)
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketLifecycleRequest{
		Bucket: Ptr(bucketName),
	}

	delResult, err := client.DeleteBucketLifecycle(context.TODO(), delRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 204, delResult.StatusCode)
	assert.Equal(t, "204 No Content", delResult.Status)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest.Bucket = Ptr(bucketNameNotExist)
	_, err = client.PutBucketLifecycle(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))

	getRequest.Bucket = Ptr(bucketNameNotExist)
	serr = &ServiceError{}
	_, err = client.GetBucketLifecycle(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest.Bucket = Ptr(bucketNameNotExist)
	serr = &ServiceError{}
	delResult, err = client.DeleteBucketLifecycle(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestAddressingModePathStyle(t *testing.T) {
	after := before(t)
	defer after(t)

	if pathStyleBucket_ == "" {
		assert.Fail(t, "Please specify the bucket that supports the path style request.")
	}

	region := pathStyleRegion_
	if pathStyleRegion_ == "" {
		region = region_
	}

	cfg := LoadDefaultConfig().
		WithRegion(region).
		WithSignatureVersion(getSignatrueVersion()).
		WithUsePathStyle(true)

	client := getClientUseStsTokenV2(cfg)

	bucketName := pathStyleBucket_
	objectName := "key-path-style"

	// bucket
	request := &ListObjectsRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.ListObjects(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.Name, bucketName)
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Empty(t, result.Prefix)
	assert.Empty(t, result.Marker)
	assert.Empty(t, result.Delimiter)

	// bucket + subresource
	aclResult, err := client.PutBucketAcl(context.TODO(), &PutBucketAclRequest{
		Bucket: Ptr(bucketName),
		Acl:    BucketACLPrivate,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, aclResult.StatusCode)
	assert.NotEmpty(t, aclResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(2 * time.Second)

	infoResult, err := client.GetBucketInfo(context.TODO(), &GetBucketInfoRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, string(BucketACLPrivate), *infoResult.BucketInfo.ACL)

	// bucket + key
	_, err = client.DeleteObject(context.TODO(), &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)

	exist, err := client.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	assert.False(t, exist)

	content := randLowStr(10)
	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	})
	assert.Nil(t, err)

	hoResult, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.Equal(t, hoResult.ContentLength, int64(len(content)))

	// bucket + key subresource
	goaResult, err := client.GetObjectAcl(context.TODO(), &GetObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, goaResult.StatusCode)
	assert.NotEmpty(t, goaResult.Headers.Get(HeaderOssRequestID))

	// presign HeadObjRequest
	expiration := time.Now().Add(100 * time.Second)
	preResult, err := client.Presign(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err := http.NewRequest(preResult.Method, preResult.URL, nil)
	c := &http.Client{}
	assert.Nil(t, err)
	resp, _ := c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header.Get(HTTPHeaderContentLength), strconv.Itoa(len(content)))

	assert.Contains(t, preResult.URL, fmt.Sprintf("oss-%s.aliyuncs.com/%s/%s?", region, bucketName, objectName))
}

func TestPublicAccessBlock(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	client := getDefaultClient()

	putResult, err := client.PutPublicAccessBlock(context.TODO(), &PutPublicAccessBlockRequest{
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetPublicAccessBlock(context.TODO(), &GetPublicAccessBlockRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeletePublicAccessBlock(context.TODO(), &DeletePublicAccessBlockRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	putResult, err = noPermClient.PutPublicAccessBlock(context.TODO(), &PutPublicAccessBlockRequest{
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getResult, err = noPermClient.GetPublicAccessBlock(context.TODO(), &GetPublicAccessBlockRequest{})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delResult, err = noPermClient.DeletePublicAccessBlock(context.TODO(), &DeletePublicAccessBlockRequest{})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketPublicAccessBlock(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)
	putResult, err := client.PutBucketPublicAccessBlock(context.TODO(), &PutBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetBucketPublicAccessBlock(context.TODO(), &GetBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteBucketPublicAccessBlock(context.TODO(), &DeleteBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	putResult, err = noPermClient.PutBucketPublicAccessBlock(context.TODO(), &PutBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getResult, err = noPermClient.GetBucketPublicAccessBlock(context.TODO(), &GetBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delResult, err = noPermClient.DeleteBucketPublicAccessBlock(context.TODO(), &DeleteBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestCname(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	createResult, err := client.CreateCnameToken(context.TODO(), &CreateCnameTokenRequest{
		Bucket: Ptr(bucketName),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotNil(t, createResult.CnameToken)
	time.Sleep(1 * time.Second)

	getResult, err := client.GetCnameToken(context.TODO(), &GetCnameTokenRequest{
		Bucket: Ptr(bucketName),
		Cname:  Ptr("example.com"),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotNil(t, getResult.CnameToken)
	time.Sleep(1 * time.Second)

	listResult, err := client.ListCname(context.TODO(), &ListCnameRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, len(listResult.Cnames), 0)
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteCname(context.TODO(), &DeleteCnameRequest{
		Bucket: Ptr(bucketName),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	_, err = client.PutCname(context.TODO(), &PutCnameRequest{
		Bucket: Ptr(bucketName),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "NeedVerifyDomainOwnership", serr.Code)
	assert.Equal(t, "0018-00000115", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	createResult, err = client.CreateCnameToken(context.TODO(), &CreateCnameTokenRequest{
		Bucket: Ptr(bucketNameNotExist),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.GetCnameToken(context.TODO(), &GetCnameTokenRequest{
		Bucket: Ptr(bucketNameNotExist),
		Cname:  Ptr("example.com"),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.ListCname(context.TODO(), &ListCnameRequest{
		Bucket: Ptr(bucketNameNotExist),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DeleteCname(context.TODO(), &DeleteCnameRequest{
		Bucket: Ptr(bucketNameNotExist),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestMetaQuery(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	bucketNameAiSearch := bucketNamePrefix + randLowStr(6)
	_, err = client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketNameAiSearch),
	})
	assert.Nil(t, err)

	openRequest := &OpenMetaQueryRequest{
		Bucket: Ptr(bucketName),
	}
	openResult, err := client.OpenMetaQuery(context.TODO(), openRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, openResult.StatusCode)
	assert.NotEmpty(t, openResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetMetaQueryStatusRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetMetaQueryStatus(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	time.Sleep(1 * time.Second)

	doRequest := &DoMetaQueryRequest{
		Bucket: Ptr(bucketName),
		MetaQuery: &MetaQuery{
			Query: Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
			Sort:  Ptr("Size"),
			Order: Ptr(MetaQueryOrderAsc),
			Aggregations: &MetaQueryAggregations{
				[]MetaQueryAggregation{
					{
						Field:     Ptr("Size"),
						Operation: Ptr("sum"),
					},
					{
						Field:     Ptr("Size"),
						Operation: Ptr("max"),
					},
				},
			},
		},
	}
	doResult, err := client.DoMetaQuery(context.TODO(), doRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, doResult.StatusCode)
	assert.Empty(t, *doResult.NextToken)
	assert.Equal(t, len(doResult.Files), 0)
	assert.Equal(t, len(doResult.Aggregations), 2)
	time.Sleep(1 * time.Second)

	closeRequest := &CloseMetaQueryRequest{
		Bucket: Ptr(bucketName),
	}
	closeResult, err := client.CloseMetaQuery(context.TODO(), closeRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, closeResult.StatusCode)
	time.Sleep(1 * time.Second)

	openRequest = &OpenMetaQueryRequest{
		Bucket: Ptr(bucketNameAiSearch),
		Mode:   Ptr("semantic"),
	}
	openResult, err = client.OpenMetaQuery(context.TODO(), openRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, openResult.StatusCode)
	assert.NotEmpty(t, openResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	doRequest = &DoMetaQueryRequest{
		Bucket: Ptr(bucketNameAiSearch),
		Mode:   Ptr("semantic"),
		MetaQuery: &MetaQuery{
			MaxResults:  Ptr(int64(99)),
			Query:       Ptr("Overlook the snow-covered forest"),
			MediaType:   Ptr("image"),
			SimpleQuery: Ptr(`{"Operation":"gt", "Field": "Size", "Value": "30"}`),
		},
	}
	doResult, err = client.DoMetaQuery(context.TODO(), doRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, doResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	openRequest = &OpenMetaQueryRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	openResult, err = client.OpenMetaQuery(context.TODO(), openRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	openRequest = &OpenMetaQueryRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	openResult, err = client.OpenMetaQuery(context.TODO(), openRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getRequest = &GetMetaQueryStatusRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetMetaQueryStatus(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	doRequest = &DoMetaQueryRequest{
		Bucket: Ptr(bucketNameNotExist),
		MetaQuery: &MetaQuery{
			Query: Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
			Sort:  Ptr("Size"),
			Order: Ptr(MetaQueryOrderAsc),
			Aggregations: &MetaQueryAggregations{
				[]MetaQueryAggregation{
					{
						Field:     Ptr("Size"),
						Operation: Ptr("sum"),
					},
					{
						Field:     Ptr("Size"),
						Operation: Ptr("max"),
					},
				},
			},
		},
	}
	doResult, err = client.DoMetaQuery(context.TODO(), doRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	closeRequest = &CloseMetaQueryRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	closeResult, err = client.CloseMetaQuery(context.TODO(), closeRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestAccessPoint(t *testing.T) {
	if apEnable == "" {
		return
	}
	after := before(t)
	defer after(t)
	defer clearAp(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	accessPointName := "ap-01-" + randLowStr(5)
	createResult, err := client.CreateAccessPoint(context.TODO(), &CreateAccessPointRequest{
		Bucket: Ptr(bucketName),
		CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
			AccessPointName: Ptr(accessPointName),
			NetworkOrigin:   Ptr("internet"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listResult, err := client.ListAccessPoints(context.TODO(), &ListAccessPointsRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	policy := `{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["` + accountID_ + `"],"Resource":["acs:oss:` + region_ + `:` + accountID_ + `:accesspoint/` + accessPointName + `","acs:oss:` + region_ + `:` + accountID_ + `:accesspoint/` + accessPointName + `/object/*"]}]}`
	putPolicyResult, err := client.PutAccessPointPolicy(context.TODO(), &PutAccessPointPolicyRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
		Body:            strings.NewReader(policy),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putPolicyResult.StatusCode)
	assert.NotEmpty(t, putPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getPolicyResult, err := client.GetAccessPointPolicy(context.TODO(), &GetAccessPointPolicyRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putPolicyResult.StatusCode)
	assert.NotEmpty(t, getPolicyResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, policy, getPolicyResult.Body)
	time.Sleep(1 * time.Second)

	delPolicyResult, err := client.DeleteAccessPointPolicy(context.TODO(), &DeleteAccessPointPolicyRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delPolicyResult.StatusCode)
	assert.NotEmpty(t, delPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	for {
		getResult, err = client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
			Bucket:          Ptr(bucketName),
			AccessPointName: Ptr(accessPointName),
		})
		if *getResult.AccessPointStatus != "creating" {
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	delResult, err := client.DeleteAccessPoint(context.TODO(), &DeleteAccessPointRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	createResult, err = client.CreateAccessPoint(context.TODO(), &CreateAccessPointRequest{
		Bucket: Ptr(bucketNameNotExist),
		CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
			AccessPointName: Ptr(accessPointName),
			NetworkOrigin:   Ptr("internet"),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getResult, err = client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	listResult, err = client.ListAccessPoints(context.TODO(), &ListAccessPointsRequest{
		Bucket: Ptr(bucketNameNotExist),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	putPolicyResult, err = client.PutAccessPointPolicy(context.TODO(), &PutAccessPointPolicyRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
		Body:            strings.NewReader(policy),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getPolicyResult, err = client.GetAccessPointPolicy(context.TODO(), &GetAccessPointPolicyRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delPolicyResult, err = client.DeleteAccessPointPolicy(context.TODO(), &DeleteAccessPointPolicyRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delResult, err = client.DeleteAccessPoint(context.TODO(), &DeleteAccessPointRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestAccessPointPublicAccessBlock(t *testing.T) {
	if apEnable == "" {
		return
	}
	after := before(t)
	defer after(t)
	defer clearAp(t)
	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	accessPointName := "ap-01-" + randLowStr(5)
	createResult, err := client.CreateAccessPoint(context.TODO(), &CreateAccessPointRequest{
		Bucket: Ptr(bucketName),
		CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
			AccessPointName: Ptr(accessPointName),
			NetworkOrigin:   Ptr("internet"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	putResult, err := client.PutAccessPointPublicAccessBlock(context.TODO(), &PutAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetAccessPointPublicAccessBlock(context.TODO(), &GetAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.True(t, *getResult.PublicAccessBlockConfiguration.BlockPublicAccess)
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteAccessPointPublicAccessBlock(context.TODO(), &DeleteAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putResult, err = client.PutAccessPointPublicAccessBlock(context.TODO(), &PutAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getResult, err = client.GetAccessPointPublicAccessBlock(context.TODO(), &GetAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delResult, err = client.DeleteAccessPointPublicAccessBlock(context.TODO(), &DeleteAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	for {
		getPointResult, err := client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
			Bucket:          Ptr(bucketName),
			AccessPointName: Ptr(accessPointName),
		})
		assert.Nil(t, err)
		if *getPointResult.AccessPointStatus != "creating" {
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	delPointResult, err := client.DeleteAccessPoint(context.TODO(), &DeleteAccessPointRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delPointResult.StatusCode)
	time.Sleep(1 * time.Second)
}

func TestCleanRestoredObject(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	objectRequest := &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		StorageClass: StorageClassColdArchive,
	}
	_, err = client.PutObject(context.TODO(), objectRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	_, err = client.RestoreObject(context.TODO(), &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	_, err = client.CleanRestoredObject(context.TODO(), &CleanRestoredObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "ArchiveRestoreNotFinished", serr.Code)
	assert.Equal(t, "The archive file's restore is not finished.", serr.Message)
	assert.Equal(t, "0016-00000719", serr.EC)
}

func TestAccessPointForObjectProcess(t *testing.T) {
	if apEnable == "" {
		return
	}
	after := before(t)
	defer after(t)
	defer clearAp(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)
	accessPointName := "ap-01-" + randLowStr(5)
	createResult, err := client.CreateAccessPoint(context.TODO(), &CreateAccessPointRequest{
		Bucket: Ptr(bucketName),
		CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
			AccessPointName: Ptr(accessPointName),
			NetworkOrigin:   Ptr("internet"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	objectProcessName := "fc-ap-01-" + randLowStr(5)
	arn := "acs:fc:" + region_ + ":" + accountID_ + ":services/test-oss-fc.LATEST/functions/" + objectProcessName
	roleArn := "acs:ram::" + accountID_ + ":role/aliyunfcdefaultrole"
	createObjectResult, err := client.CreateAccessPointForObjectProcess(context.TODO(), &CreateAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		CreateAccessPointForObjectProcessConfiguration: &CreateAccessPointForObjectProcessConfiguration{
			AccessPointName: Ptr(accessPointName),
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: []string{"GetObject-Range"},
				TransformationConfigurations: []TransformationConfiguration{
					{
						Actions: &AccessPointActions{
							[]string{"GetObject"},
						},
						ContentTransformation: &ContentTransformation{
							FunctionArn:           Ptr(arn),
							FunctionAssumeRoleArn: Ptr(roleArn),
						},
					},
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createObjectResult.StatusCode)
	assert.NotEmpty(t, createObjectResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	getResult, err := client.GetAccessPointForObjectProcess(context.TODO(), &GetAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	listResult, err := client.ListAccessPointsForObjectProcess(context.TODO(), &ListAccessPointsForObjectProcessRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	policy := `{"Version":"1","Statement":[{"Action":["oss:GetObject"],"Effect":"Allow","Principal":["` + accountID_ + `"],"Resource":["acs:oss:` + region_ + `:` + accountID_ + `:accesspointforobjectprocess/` + objectProcessName + `/object/*"]}]}`
	putPolicyResult, err := client.PutAccessPointPolicyForObjectProcess(context.TODO(), &PutAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		Body:                            strings.NewReader(policy),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putPolicyResult.StatusCode)
	assert.NotEmpty(t, putPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	getPolicyResult, err := client.GetAccessPointPolicyForObjectProcess(context.TODO(), &GetAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putPolicyResult.StatusCode)
	assert.NotEmpty(t, getPolicyResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, getPolicyResult.Body, policy)
	time.Sleep(1 * time.Second)
	delPolicyResult, err := client.DeleteAccessPointPolicyForObjectProcess(context.TODO(), &DeleteAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delPolicyResult.StatusCode)
	assert.NotEmpty(t, delPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	putConfigResult, err := client.PutAccessPointConfigForObjectProcess(context.TODO(), &PutAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		PutAccessPointConfigForObjectProcessConfiguration: &PutAccessPointConfigForObjectProcessConfiguration{
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: []string{"GetObject-Range"},
				TransformationConfigurations: []TransformationConfiguration{
					{
						Actions: &AccessPointActions{
							[]string{"GetObject"},
						},
						ContentTransformation: &ContentTransformation{
							FunctionArn:           Ptr(arn),
							FunctionAssumeRoleArn: Ptr(roleArn),
						},
					},
				},
			},
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				Ptr(true),
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putConfigResult.StatusCode)
	assert.NotEmpty(t, putConfigResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	getConfigResult, err := client.GetAccessPointConfigForObjectProcess(context.TODO(), &GetAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getConfigResult.StatusCode)
	assert.NotEmpty(t, getConfigResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	for {
		getResult, err := client.GetAccessPointForObjectProcess(context.TODO(), &GetAccessPointForObjectProcessRequest{
			Bucket:                          Ptr(bucketName),
			AccessPointForObjectProcessName: Ptr(objectProcessName),
		})
		assert.Nil(t, err)
		if *getResult.AccessPointForObjectProcessStatus != "creating" {
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	delResult, err := client.DeleteAccessPointForObjectProcess(context.TODO(), &DeleteAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	var serr *ServiceError
	for {
		_, err := client.GetAccessPointForObjectProcess(context.TODO(), &GetAccessPointForObjectProcessRequest{
			Bucket:                          Ptr(bucketName),
			AccessPointForObjectProcessName: Ptr(objectProcessName),
		})
		if err != nil {
			errors.As(err, &serr)
			if serr.StatusCode == 404 && serr.Code == "NoSuchAccessPointForObjectProcess" {
				break
			}
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	for {
		gResult, err := client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
			Bucket:          Ptr(bucketName),
			AccessPointName: Ptr(accessPointName),
		})
		assert.Nil(t, err)
		if *gResult.AccessPointStatus != "creating" {
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	dResult, err := client.DeleteAccessPoint(context.TODO(), &DeleteAccessPointRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, dResult.StatusCode)
	time.Sleep(1 * time.Second)
	bucketNameNotExist := bucketName + "-not-exist"
	_, err = client.CreateAccessPointForObjectProcess(context.TODO(), &CreateAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		CreateAccessPointForObjectProcessConfiguration: &CreateAccessPointForObjectProcessConfiguration{
			AccessPointName: Ptr(accessPointName),
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: []string{"GetObject-Range"},
				TransformationConfigurations: []TransformationConfiguration{
					{
						Actions: &AccessPointActions{
							[]string{"GetObject"},
						},
						ContentTransformation: &ContentTransformation{
							FunctionArn:           Ptr(arn),
							FunctionAssumeRoleArn: Ptr(roleArn),
						},
					},
				},
			},
		},
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	getResult, err = client.GetAccessPointForObjectProcess(context.TODO(), &GetAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	listResult, err = noPermClient.ListAccessPointsForObjectProcess(context.TODO(), &ListAccessPointsForObjectProcessRequest{})
	serr = &ServiceError{}
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	putPolicyResult, err = client.PutAccessPointPolicyForObjectProcess(context.TODO(), &PutAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		Body:                            strings.NewReader(policy),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	getPolicyResult, err = client.GetAccessPointPolicyForObjectProcess(context.TODO(), &GetAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	delPolicyResult, err = client.DeleteAccessPointPolicyForObjectProcess(context.TODO(), &DeleteAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	putConfigResult, err = client.PutAccessPointConfigForObjectProcess(context.TODO(), &PutAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		PutAccessPointConfigForObjectProcessConfiguration: &PutAccessPointConfigForObjectProcessConfiguration{
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: []string{"GetObject-Range"},
				TransformationConfigurations: []TransformationConfiguration{
					{
						Actions: &AccessPointActions{
							[]string{"GetObject"},
						},
						ContentTransformation: &ContentTransformation{
							FunctionArn:           Ptr(arn),
							FunctionAssumeRoleArn: Ptr(roleArn),
						},
					},
				},
			},
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				Ptr(true),
			},
		},
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	_, err = client.GetAccessPointConfigForObjectProcess(context.TODO(), &GetAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	_, err = client.DeleteAccessPointForObjectProcess(context.TODO(), &DeleteAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestRedundancyTransition(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
		CreateBucketConfiguration: &CreateBucketConfiguration{
			StorageClass:       StorageClassStandard,
			DataRedundancyType: DataRedundancyLRS,
		},
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)
	createResult, err := client.CreateBucketDataRedundancyTransition(context.TODO(), &CreateBucketDataRedundancyTransitionRequest{
		Bucket:               Ptr(bucketName),
		TargetRedundancyType: Ptr("ZRS"),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetBucketDataRedundancyTransition(context.TODO(), &GetBucketDataRedundancyTransitionRequest{
		Bucket:                     Ptr(bucketName),
		RedundancyTransitionTaskid: createResult.BucketDataRedundancyTransition.TaskId,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listResult, err := client.ListBucketDataRedundancyTransition(context.TODO(), &ListBucketDataRedundancyTransitionRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listUserResult, err := client.ListUserDataRedundancyTransition(context.TODO(), &ListUserDataRedundancyTransitionRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 200, listUserResult.StatusCode)
	assert.NotEmpty(t, listUserResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteBucketDataRedundancyTransition(context.TODO(), &DeleteBucketDataRedundancyTransitionRequest{
		Bucket:                     Ptr(bucketName),
		RedundancyTransitionTaskid: createResult.BucketDataRedundancyTransition.TaskId,
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	bucketNameNotExist := bucketName + "-not-exist"
	_, err = client.CreateBucketDataRedundancyTransition(context.TODO(), &CreateBucketDataRedundancyTransitionRequest{
		Bucket:               Ptr(bucketNameNotExist),
		TargetRedundancyType: Ptr("ZRS"),
	})
	serr := &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	_, err = client.GetBucketDataRedundancyTransition(context.TODO(), &GetBucketDataRedundancyTransitionRequest{
		Bucket:                     Ptr(bucketNameNotExist),
		RedundancyTransitionTaskid: createResult.BucketDataRedundancyTransition.TaskId,
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	_, err = client.DeleteBucketDataRedundancyTransition(context.TODO(), &DeleteBucketDataRedundancyTransitionRequest{
		Bucket:                     Ptr(bucketNameNotExist),
		RedundancyTransitionTaskid: createResult.BucketDataRedundancyTransition.TaskId,
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	_, err = client.ListBucketDataRedundancyTransition(context.TODO(), &ListBucketDataRedundancyTransitionRequest{
		Bucket: Ptr(bucketNameNotExist),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	_, err = noPermClient.ListUserDataRedundancyTransition(context.TODO(), &ListUserDataRedundancyTransitionRequest{})
	serr = &ServiceError{}
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestRegions(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	client := getDefaultClient()
	result, err := client.DescribeRegions(context.TODO(), &DescribeRegionsRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	result, err = client.DescribeRegions(context.TODO(), &DescribeRegionsRequest{
		Regions: Ptr("oss-cn-hangzhou"),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, len(result.RegionInfoList.RegionInfos), 1)
	time.Sleep(1 * time.Second)

	serr := &ServiceError{}
	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	_, err = noPermClient.DescribeRegions(context.TODO(), &DescribeRegionsRequest{})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = noPermClient.DescribeRegions(context.TODO(), &DescribeRegionsRequest{
		Regions: Ptr("oss-cn-hangzhou"),
	})
	serr = &ServiceError{}
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDownloaderTruncate(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	objectLen := 100*1024*4 + 123
	content := randLowStr(objectLen)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	result, err := client.PutObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 100 * 1024
		do.UseTempFile = false
	})
	assert.NotNil(t, d)
	assert.Equal(t, 3, d.options.ParallelNum)
	localFile := randStr(8) + "-downloader"
	defer os.Remove(localFile)

	// file size is 0, file size < get size
	dResult, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
		},
		localFile)

	assert.Nil(t, err)

	assert.Equal(t, objectLen, (int)(dResult.Written))
	gotCotent, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, []byte(content), gotCotent)

	//range get
	// file size > get size
	dResult, err = d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
			Range:  Ptr("bytes=123-123456"),
		},
		localFile)

	gotCotent, err = os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, 123456-123+1, (int)(dResult.Written))
	assert.Equal(t, []byte(content[123:123456+1]), []byte(gotCotent))

	// file size == get size
	content1 := randLowStr(123456 - 123 + 1)
	os.WriteFile(localFile, []byte(content1), 0644)
	gotCotent, err = os.ReadFile(localFile)
	assert.Equal(t, []byte(content1), []byte(gotCotent))
	dResult, err = d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
			Range:  Ptr("bytes=123-123456"),
		},
		localFile)

	gotCotent, err = os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, 123456-123+1, (int)(dResult.Written))
	assert.Equal(t, []byte(content[123:123456+1]), []byte(gotCotent))
}

func TestUploaderWithSequential(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	objectLen := 100*1024*4 + 123
	content := randLowStr(objectLen)
	partSize := int64(100 * 1024)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.PartSize = partSize
		},
	)

	result, err := u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}, strings.NewReader(content))
	assert.Nil(t, err)
	assert.NotNil(t, result)

	request := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	headResult, err := client.HeadObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, headResult.ContentMD5)

	result, err = u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"sequential": "",
			},
		},
	}, strings.NewReader(content))
	assert.Nil(t, err)
	assert.NotNil(t, result)

	request = &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	headResult, err = client.HeadObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.NotEmpty(t, headResult.ContentMD5)
}

func TestUploaderUploadFromTeeReader(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	partSize := int64(100 * 1024)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.PartSize = partSize
		},
	)

	// empty reader
	content := strings.Repeat("a", 0)
	reader := io.TeeReader(strings.NewReader(content), io.Discard)

	result, err := u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}, reader)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	request := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	headResult, err := client.HeadObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), headResult.ContentLength)

	// 2 * part-size
	content = strings.Repeat("a", 200*1024)
	reader = io.TeeReader(strings.NewReader(content), io.Discard)
	result, err = u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}, strings.NewReader(content))
	assert.Nil(t, err)
	assert.NotNil(t, result)

	request = &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	headResult, err = client.HeadObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, int64(200*1024), headResult.ContentLength)
}
