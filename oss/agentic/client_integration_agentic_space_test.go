//go:build integrationignore

package agentic

import (
	"context"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

// TestAgenticBucketSpace: ListBucketSpaces, object read/write via the bucket space
// client, and the BucketSpaceHelper name builder, over one shared bucket.
func TestAgenticBucketSpace(t *testing.T) {
	skipIfNotConfigured(t)
	client := getAgenticBucketClient()
	bsClient := getBucketSpaceClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer disableAndReap(bucket)

	t.Run("ListBucketSpaces", func(t *testing.T) {
		listResult, err := client.ListBucketSpaces(context.TODO(), &ListBucketSpacesRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, listResult.StatusCode)
	})

	t.Run("ObjectLifecycle", func(t *testing.T) {
		// Create a bucket space using the short name.
		putBucketResult, err := bsClient.PutBucket(context.TODO(), &oss.PutBucketRequest{
			Bucket: oss.Ptr(bucket),
		})
		dumpErrIfNotNil(err)
		assert.Nil(t, err)
		assert.Equal(t, 200, putBucketResult.StatusCode)

		defer func() {
			_, _ = bsClient.DeleteBucket(context.TODO(), &oss.DeleteBucketRequest{
				Bucket: oss.Ptr(bucket),
			})
		}()

		key := "go-sdk-test-object-" + randStr(6)
		putObjectResult, err := bsClient.PutObject(context.TODO(), &oss.PutObjectRequest{
			Bucket: oss.Ptr(bucket),
			Key:    oss.Ptr(key),
			Body:   strings.NewReader("hello world"),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, putObjectResult.StatusCode)

		defer func() {
			_, _ = bsClient.DeleteObject(context.TODO(), &oss.DeleteObjectRequest{
				Bucket: oss.Ptr(bucket),
				Key:    oss.Ptr(key),
			})
		}()

		getObjectResult, err := bsClient.GetObject(context.TODO(), &oss.GetObjectRequest{
			Bucket: oss.Ptr(bucket),
			Key:    oss.Ptr(key),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getObjectResult.StatusCode)
		getObjectResult.Body.Close()
	})

	t.Run("SpaceHelper", func(t *testing.T) {
		helper := NewBucketSpaceHelper(getTestConfig())
		name := helper.ToBucketName(bucket)
		expected := buildFullName(bucket, accountId_, region_, "bs-apsr")
		assert.Equal(t, expected, name)
	})
}
