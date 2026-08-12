//go:build integration

package agentic

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

// TestAgenticBucketBasic: one shared bucket, exercising Create/Get/List/PutStatus.
// Lifecycle Disable+Delete is asserted in TestAgenticBucketAttribute instead.
func TestAgenticBucketBasic(t *testing.T) {
	skipIfNotConfigured(t)
	client := getAgenticBucketClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
		CreateAgenticBucketConfiguration: &CreateAgenticBucketConfiguration{
			StorageClass:       oss.StorageClassStandard,
			DataRedundancyType: oss.DataRedundancyLRS,
		},
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer disableAndReap(bucket)

	t.Run("Get", func(t *testing.T) {
		getResult, err := client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getResult.StatusCode)
		assert.NotNil(t, getResult.AgenticBucketInfo)
		assert.Contains(t, oss.ToString(getResult.AgenticBucketInfo.Name), bucket)
	})

	t.Run("List", func(t *testing.T) {
		// A newly created bucket may take a moment to appear in the list, so poll.
		found := false
		for attempt := 0; attempt < 5 && !found; attempt++ {
			if attempt > 0 {
				time.Sleep(10 * time.Second)
			}
			paginator := client.NewListAgenticBucketsPaginator(&ListAgenticBucketsRequest{})
			for paginator.HasNext() {
				page, err := paginator.NextPage(context.TODO())
				assert.Nil(t, err)
				for _, b := range page.AgenticBuckets {
					if b.Name != nil && strings.Contains(*b.Name, bucket) {
						found = true
					}
				}
			}
		}
		assert.True(t, found, "created agentic bucket should appear in list")
	})

	t.Run("PutStatusEnabled", func(t *testing.T) {
		putResult, err := client.PutAgenticBucketStatus(context.TODO(), &PutAgenticBucketStatusRequest{
			Bucket: oss.Ptr(bucket),
			AgenticBucketStatus: &AgenticBucketStatus{
				Status: oss.Ptr("Enabled"),
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, putResult.StatusCode)
	})
}

// TestAgenticBucketServerErrors checks error propagation with invalid credentials.
func TestAgenticBucketServerErrors(t *testing.T) {
	skipIfNotConfigured(t)
	client := getInvalidAkClient()
	bucket := genBucketName()

	var serr *oss.ServiceError

	// Create with invalid AK
	_, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, 403, serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)

	// Get with invalid AK
	serr = nil
	_, err = client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, 404, serr.StatusCode)

	// List with invalid AK
	serr = nil
	_, err = client.ListAgenticBuckets(context.TODO(), &ListAgenticBucketsRequest{})
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, 403, serr.StatusCode)
}
