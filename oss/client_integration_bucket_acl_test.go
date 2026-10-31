//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	_, err = client.PutBucketPublicAccessBlock(context.TODO(), &PutBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			BlockPublicAccess: Ptr(false),
		},
	})
	time.Sleep(10 * time.Second)

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
