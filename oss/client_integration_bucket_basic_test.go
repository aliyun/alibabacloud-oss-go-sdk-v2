//go:build integration

package oss

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

	bucketTagName := bucketPrefix + strconv.Itoa(1)
	putTagsRequest := &PutBucketTagsRequest{
		Bucket: Ptr(bucketTagName),
		Tagging: &Tagging{
			&TagSet{
				[]Tag{
					{
						Ptr("k"),
						Ptr("v"),
					},
				},
			},
		},
	}
	_, err = client.PutBucketTags(context.TODO(), putTagsRequest)
	listRequest = &ListBucketsRequest{
		Prefix:  Ptr(bucketPrefix),
		Tagging: Ptr("\"k\":\"v\""),
	}

	result, err = client.ListBuckets(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, len(result.Buckets), 1)

	listRequest = &ListBucketsRequest{
		Prefix:   Ptr(bucketPrefix),
		TagKey:   Ptr("k"),
		TagValue: Ptr("v"),
	}

	result, err = client.ListBuckets(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, len(result.Buckets), 1)

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
	assert.True(t, *info.BucketInfo.BlockPublicAccess)
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
	assert.Greater(t, stat.LastModifiedTime, int64(0))
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
