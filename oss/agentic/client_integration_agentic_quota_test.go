//go:build integration

package agentic

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestAgenticBucketStorageQuota(t *testing.T) {
	skipIfNotConfigured(t)
	client := getAgenticBucketClient()
	bucket := genBucketName()
	quota := int64(10737418240)

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, createResult.StatusCode)
	defer disableAndReap(bucket)

	putResult, err := client.PutAgenticBucketStorageQuota(context.TODO(), &PutAgenticBucketStorageQuotaRequest{
		Bucket: oss.Ptr(bucket),
		QuotaConfiguration: &oss.QuotaConfiguration{
			StorageQuota: &quota,
			Mode:         oss.StorageQuotaModeStrict,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, putResult.StatusCode)

	getResult, err := client.GetAgenticBucketStorageQuota(context.TODO(), &GetAgenticBucketStorageQuotaRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResult.StatusCode)
	if assert.NotNil(t, getResult.QuotaConfiguration) {
		assert.Equal(t, quota, oss.ToInt64(getResult.QuotaConfiguration.StorageQuota))
		assert.Equal(t, oss.StorageQuotaModeStrict, getResult.QuotaConfiguration.Mode)
	}

	deleteResult, err := client.DeleteAgenticBucketStorageQuota(context.TODO(), &DeleteAgenticBucketStorageQuotaRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.NoError(t, err)
	assert.True(t, deleteResult.StatusCode == http.StatusOK || deleteResult.StatusCode == http.StatusNoContent)
}

func TestAgenticBucketStorageQuotaErrors(t *testing.T) {
	skipIfNotConfigured(t)
	client := getInvalidAkClient()
	bucket := genBucketName()
	quota := int64(10737418240)

	var serr *oss.ServiceError

	_, err := client.PutAgenticBucketStorageQuota(context.TODO(), &PutAgenticBucketStorageQuotaRequest{
		Bucket: oss.Ptr(bucket),
		QuotaConfiguration: &oss.QuotaConfiguration{
			StorageQuota: &quota,
			Mode:         oss.StorageQuotaModeStrict,
		},
	})
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, http.StatusNotFound, serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)

	serr = nil
	_, err = client.GetAgenticBucketStorageQuota(context.TODO(), &GetAgenticBucketStorageQuotaRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, http.StatusNotFound, serr.StatusCode)

	serr = nil
	_, err = client.DeleteAgenticBucketStorageQuota(context.TODO(), &DeleteAgenticBucketStorageQuotaRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, http.StatusNotFound, serr.StatusCode)
}
