//go:build integration

package oss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

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

func TestDefaultRequestHeaders(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)

	newClient := func(headers map[string]string, ak, sk string) *Client {
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(ak, sk)).
			WithRegion(region_).
			WithEndpoint(endpoint_).
			WithSignatureVersion(getSignatrueVersion()).
			WithDefaultRequestHeaders(headers)
		return NewClient(cfg)
	}

	// normal case, a signable x-oss-* default header must take part in the
	// signature and reach the server
	client := newClient(map[string]string{
		"x-oss-meta-app": "sdk-default-headers",
		"x-my-trace-id":  "trace-" + randLowStr(6),
	}, accessID_, accessKey_)

	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{Bucket: Ptr(bucketName)})
	assert.Nil(t, err)

	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader("hello default headers"),
	})
	assert.Nil(t, err)

	hResult, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.Equal(t, "sdk-default-headers", hResult.Metadata["app"])

	gResult, err := client.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	data, err := io.ReadAll(gResult.Body)
	gResult.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, "hello default headers", string(data))
	assert.Equal(t, "sdk-default-headers", gResult.Metadata["app"])

	// the request wins over the default on conflict
	objectNameOverride := objectName + "-override"
	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameOverride),
		Metadata: map[string]string{"app": "from-request"},
		Body:     strings.NewReader("hi"),
	})
	assert.Nil(t, err)

	hResult, err = client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameOverride),
	})
	assert.Nil(t, err)
	assert.Equal(t, "from-request", hResult.Metadata["app"])

	// a default header rejected by the server surfaces as a service error
	forbidClient := newClient(map[string]string{
		"x-oss-forbid-overwrite": "true",
	}, accessID_, accessKey_)

	_, err = forbidClient.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader("overwrite me"),
	})
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "FileAlreadyExists", serr.Code)
	assert.NotEmpty(t, serr.RequestID)

	// default headers do not disturb the not-found path
	_, err = client.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName + "-not-exist"),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchKey", serr.Code)
	assert.Equal(t, "The specified key does not exist.", serr.Message)

	// nor the invalid credentials path
	noPermClient := newClient(map[string]string{
		"x-oss-meta-app": "sdk-default-headers",
	}, "invalid-ak", "invalid-sk")

	_, err = noPermClient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.NotEmpty(t, serr.RequestID)
}
