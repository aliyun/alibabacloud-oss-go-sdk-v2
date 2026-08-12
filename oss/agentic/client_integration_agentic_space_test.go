//go:build integration

package agentic

import (
	"context"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

// TestAgenticBucketSpace: ListBucketSpaces, bucket/object interfaces via the bucket
// space client, and the BucketSpaceHelper name builder driving a plain oss.Client,
// over one shared bucket space.
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

	// Create one bucket space (short name) shared by the subtests below.
	putBucketResult, err := bsClient.PutBucket(context.TODO(), &oss.PutBucketRequest{
		Bucket:        oss.Ptr(bucket),
		AgenticBucket: oss.Ptr(buildFullName(bucket, accountId_, region_, "ab-apsr")),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, putBucketResult.StatusCode)

	defer func() {
		_, _ = bsClient.DeleteBucket(context.TODO(), &oss.DeleteBucketRequest{
			Bucket: oss.Ptr(bucket),
		})
	}()

	t.Run("ListBucketSpaces", func(t *testing.T) {
		listResult, err := client.ListBucketSpaces(context.TODO(), &ListBucketSpacesRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, listResult.StatusCode)
	})

	t.Run("BucketLifecycle", func(t *testing.T) {
		putAclResult, err := bsClient.PutBucketAcl(context.TODO(), &oss.PutBucketAclRequest{
			Bucket: oss.Ptr(bucket),
			Acl:    oss.BucketACLPrivate,
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, putAclResult.StatusCode)

		getAclResult, err := bsClient.GetBucketAcl(context.TODO(), &oss.GetBucketAclRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getAclResult.StatusCode)
		assert.Equal(t, string(oss.BucketACLPrivate), oss.ToString(getAclResult.ACL))
	})

	t.Run("ObjectLifecycle", func(t *testing.T) {
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

	// Drive the same space through a plain oss.Client using a helper-built full name.
	t.Run("SpaceHelper", func(t *testing.T) {
		helper := NewBucketSpaceHelper(getTestConfig())
		fullName := helper.ToBucketName(bucket)
		assert.Equal(t, buildFullName(bucket, accountId_, region_, "bs-apsr"), fullName)

		plainClient := oss.NewClient(getTestConfig())

		getInfoResult, err := plainClient.GetBucketInfo(context.TODO(), &oss.GetBucketInfoRequest{
			Bucket: oss.Ptr(fullName),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getInfoResult.StatusCode)

		key := "go-sdk-test-object-" + randStr(6)
		putObjectResult, err := plainClient.PutObject(context.TODO(), &oss.PutObjectRequest{
			Bucket: oss.Ptr(fullName),
			Key:    oss.Ptr(key),
			Body:   strings.NewReader("hello helper"),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, putObjectResult.StatusCode)

		defer func() {
			_, _ = plainClient.DeleteObject(context.TODO(), &oss.DeleteObjectRequest{
				Bucket: oss.Ptr(fullName),
				Key:    oss.Ptr(key),
			})
		}()

		getObjectResult, err := plainClient.GetObject(context.TODO(), &oss.GetObjectRequest{
			Bucket: oss.Ptr(fullName),
			Key:    oss.Ptr(key),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getObjectResult.StatusCode)
		getObjectResult.Body.Close()
	})
}
