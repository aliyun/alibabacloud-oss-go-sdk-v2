//go:build integration

package oss

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketStorageQuota(t *testing.T) {
	after := before(t)
	defer after(t)

	bucket := bucketNamePrefix + randLowStr(6)
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{Bucket: Ptr(bucket)})
	assert.NoError(t, err)

	quota := int64(10737418240)
	putResult, err := client.PutBucketStorageQuota(context.TODO(), &PutBucketStorageQuotaRequest{
		Bucket: Ptr(bucket),
		QuotaConfiguration: &QuotaConfiguration{
			StorageQuota: &quota,
			Mode:         StorageQuotaModeStrict,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(time.Second)

	getResult, err := client.GetBucketStorageQuota(context.TODO(), &GetBucketStorageQuotaRequest{Bucket: Ptr(bucket)})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	if assert.NotNil(t, getResult.QuotaConfiguration) {
		assert.Equal(t, quota, ToInt64(getResult.QuotaConfiguration.StorageQuota))
		assert.Equal(t, StorageQuotaModeStrict, getResult.QuotaConfiguration.Mode)
	}

	deleteResult, err := client.DeleteBucketStorageQuota(context.TODO(), &DeleteBucketStorageQuotaRequest{Bucket: Ptr(bucket)})
	assert.NoError(t, err)
	assert.True(t, deleteResult.StatusCode == http.StatusOK || deleteResult.StatusCode == http.StatusNoContent)

	assertNoSuchBucket := func(err error) {
		t.Helper()
		assert.Error(t, err)
		var serviceErr *ServiceError
		assert.True(t, errors.As(err, &serviceErr))
		assert.Equal(t, http.StatusNotFound, serviceErr.StatusCode)
		assert.Equal(t, "NoSuchBucket", serviceErr.Code)
		assert.NotEmpty(t, serviceErr.RequestID)
	}
	bucketNotExist := bucketNamePrefix + randLowStr(6) + "-not-exist"

	_, err = client.PutBucketStorageQuota(context.TODO(), &PutBucketStorageQuotaRequest{
		Bucket: Ptr(bucketNotExist),
		QuotaConfiguration: &QuotaConfiguration{
			StorageQuota: &quota,
			Mode:         StorageQuotaModeStrict,
		},
	})
	assertNoSuchBucket(err)

	_, err = client.GetBucketStorageQuota(context.TODO(), &GetBucketStorageQuotaRequest{Bucket: Ptr(bucketNotExist)})
	assertNoSuchBucket(err)

	_, err = client.DeleteBucketStorageQuota(context.TODO(), &DeleteBucketStorageQuotaRequest{Bucket: Ptr(bucketNotExist)})
	assertNoSuchBucket(err)
}
