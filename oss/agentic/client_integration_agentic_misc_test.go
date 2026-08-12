//go:build integration

package agentic

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

// Miscellaneous agentic integration scenarios that do not belong to the Basic,
// Attribute, or Space suites. Add future one-offs here.

func getAgenticBucketClientPathStyle() *AgenticBucketClient {
	return NewAgenticBucketClient(getTestConfig().WithUsePathStyle(true))
}

func getBucketSpaceClientPathStyle() *oss.Client {
	return NewBucketSpaceClient(getTestConfig().WithUsePathStyle(true))
}

// isSecondLevelDomainForbidden reports whether err is the server signalling that
// path-style (second-level domain) addressing is not allowed on this endpoint.
func isSecondLevelDomainForbidden(err error) bool {
	var serr *oss.ServiceError
	return errors.As(err, &serr) && serr.Code == "SecondLevelDomainForbidden"
}

// TestAgenticPathStyle validates path-style addressing across both the agentic
// bucket client and the bucket space client. Path-style may be disabled on the
// endpoint; the probe detects SecondLevelDomainForbidden and skips rather than
// fails, since that is an endpoint capability and not an SDK defect.
func TestAgenticPathStyle(t *testing.T) {
	skipIfNotConfigured(t)

	// Create the bucket with the default (virtual-hosted) client so the fixture
	// stands regardless of whether path-style turns out to be allowed.
	client := getAgenticBucketClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer disableAndReap(bucket)

	psClient := getAgenticBucketClientPathStyle()

	// Probe: a path-style GET on the bucket. ListAgenticBuckets is service-level
	// (no bucket label) so its URL is identical in both styles and cannot probe
	// path-style; GetAgenticBucket carries the bucket and does.
	getResult, err := psClient.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	if isSecondLevelDomainForbidden(err) {
		t.Skip("path-style addressing is not allowed on this endpoint (SecondLevelDomainForbidden)")
	}
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.AgenticBucketInfo)

	// Create one bucket space (via the default client) shared by the path-style
	// subtests below.
	bsClient := getBucketSpaceClient()
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

	// Agentic bucket client over path-style.
	t.Run("AgenticClientListBucketSpaces", func(t *testing.T) {
		listResult, err := psClient.ListBucketSpaces(context.TODO(), &ListBucketSpacesRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, listResult.StatusCode)
	})

	// Bucket space client over path-style.
	psBsClient := getBucketSpaceClientPathStyle()

	t.Run("BucketSpaceClientObjectLifecycle", func(t *testing.T) {
		key := "go-sdk-test-object-" + randStr(6)
		putObjectResult, err := psBsClient.PutObject(context.TODO(), &oss.PutObjectRequest{
			Bucket: oss.Ptr(bucket),
			Key:    oss.Ptr(key),
			Body:   strings.NewReader("hello path-style"),
		})
		if isSecondLevelDomainForbidden(err) {
			t.Skip("path-style addressing is not allowed on this endpoint (SecondLevelDomainForbidden)")
		}
		dumpErrIfNotNil(err)
		assert.Nil(t, err)
		assert.Equal(t, 200, putObjectResult.StatusCode)

		defer func() {
			_, _ = psBsClient.DeleteObject(context.TODO(), &oss.DeleteObjectRequest{
				Bucket: oss.Ptr(bucket),
				Key:    oss.Ptr(key),
			})
		}()

		getObjectResult, err := psBsClient.GetObject(context.TODO(), &oss.GetObjectRequest{
			Bucket: oss.Ptr(bucket),
			Key:    oss.Ptr(key),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getObjectResult.StatusCode)
		getObjectResult.Body.Close()
	})

	t.Run("BucketSpaceClientGetBucketAcl", func(t *testing.T) {
		getAclResult, err := psBsClient.GetBucketAcl(context.TODO(), &oss.GetBucketAclRequest{
			Bucket: oss.Ptr(bucket),
		})
		if isSecondLevelDomainForbidden(err) {
			t.Skip("path-style addressing is not allowed on this endpoint (SecondLevelDomainForbidden)")
		}
		assert.Nil(t, err)
		assert.Equal(t, 200, getAclResult.StatusCode)
	})
}
