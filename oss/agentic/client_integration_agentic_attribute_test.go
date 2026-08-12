//go:build integration

package agentic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

// TestAgenticBucketAttribute runs the five attribute CRUD groups on one shared
// bucket, then a final subtest for lifecycle connectivity. It must run last:
// disabling the bucket breaks the attribute calls.
func TestAgenticBucketAttribute(t *testing.T) {
	skipIfNotConfigured(t)
	client := getAgenticBucketClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer disableAndReap(bucket)

	t.Run("Acl", func(t *testing.T) {
		putResult, err := client.PutAgenticBucketAcl(context.TODO(), &PutAgenticBucketAclRequest{
			Bucket: oss.Ptr(bucket),
			Acl:    oss.BucketACLPrivate,
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, putResult.StatusCode)

		getResult, err := client.GetAgenticBucketAcl(context.TODO(), &GetAgenticBucketAclRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getResult.StatusCode)
		assert.Equal(t, string(oss.BucketACLPrivate), oss.ToString(getResult.ACL))
	})

	t.Run("Encryption", func(t *testing.T) {
		putResult, err := client.PutAgenticBucketEncryption(context.TODO(), &PutAgenticBucketEncryptionRequest{
			Bucket: oss.Ptr(bucket),
			ServerSideEncryptionRule: &ServerSideEncryptionRule{
				ApplyServerSideEncryptionByDefault: &ApplyServerSideEncryptionByDefault{
					SSEAlgorithm: oss.Ptr("AES256"),
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, putResult.StatusCode)

		getResult, err := client.GetAgenticBucketEncryption(context.TODO(), &GetAgenticBucketEncryptionRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getResult.StatusCode)
		assert.NotNil(t, getResult.ServerSideEncryptionRule)
		assert.NotNil(t, getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault)
		assert.Equal(t, "AES256", oss.ToString(getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm))

		deleteResult, err := client.DeleteAgenticBucketEncryption(context.TODO(), &DeleteAgenticBucketEncryptionRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.True(t, deleteResult.StatusCode == 200 || deleteResult.StatusCode == 204)
	})

	t.Run("Versioning", func(t *testing.T) {
		putResult, err := client.PutAgenticBucketVersioning(context.TODO(), &PutAgenticBucketVersioningRequest{
			Bucket: oss.Ptr(bucket),
			VersioningConfiguration: &VersioningConfiguration{
				Status: oss.VersionEnabled,
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, putResult.StatusCode)

		getResult, err := client.GetAgenticBucketVersioning(context.TODO(), &GetAgenticBucketVersioningRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getResult.StatusCode)
		assert.NotNil(t, getResult.VersioningConfiguration)
		assert.Equal(t, oss.VersionEnabled, getResult.VersioningConfiguration.Status)
	})

	t.Run("PublicAccessBlock", func(t *testing.T) {
		putResult, err := client.PutAgenticBucketPublicAccessBlock(context.TODO(), &PutAgenticBucketPublicAccessBlockRequest{
			Bucket: oss.Ptr(bucket),
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				BlockPublicAccess: oss.Ptr(true),
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, putResult.StatusCode)

		getResult, err := client.GetAgenticBucketPublicAccessBlock(context.TODO(), &GetAgenticBucketPublicAccessBlockRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getResult.StatusCode)
		assert.NotNil(t, getResult.PublicAccessBlockConfiguration)

		deleteResult, err := client.DeleteAgenticBucketPublicAccessBlock(context.TODO(), &DeleteAgenticBucketPublicAccessBlockRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.True(t, deleteResult.StatusCode == 200 || deleteResult.StatusCode == 204)
	})

	t.Run("Policy", func(t *testing.T) {
		_, _ = client.PutAgenticBucketPublicAccessBlock(context.TODO(), &PutAgenticBucketPublicAccessBlockRequest{
			Bucket: oss.Ptr(bucket),
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				BlockPublicAccess: oss.Ptr(false),
			},
		})

		policy := fmt.Sprintf(`{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["*"],"Resource":["acs:oss:*:%s:*"]}]}`, accountId_)

		putResult, err := client.PutAgenticBucketPolicy(context.TODO(), &PutAgenticBucketPolicyRequest{
			Bucket: oss.Ptr(bucket),
			Body:   strings.NewReader(policy),
		})

		assert.Nil(t, err)
		assert.Equal(t, 200, putResult.StatusCode)

		getResult, err := client.GetAgenticBucketPolicy(context.TODO(), &GetAgenticBucketPolicyRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, getResult.StatusCode)
		assert.Contains(t, getResult.Body, "oss:GetObject")

		deleteResult, err := client.DeleteAgenticBucketPolicy(context.TODO(), &DeleteAgenticBucketPolicyRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.Nil(t, err)
		assert.True(t, deleteResult.StatusCode == 200 || deleteResult.StatusCode == 204)
	})

	// Connectivity only: PutStatus(Disabled) succeeds, then Delete is not yet ready.
	t.Run("LifecycleDisableThenDelete", func(t *testing.T) {
		putResult, err := client.PutAgenticBucketStatus(context.TODO(), &PutAgenticBucketStatusRequest{
			Bucket: oss.Ptr(bucket),
			AgenticBucketStatus: &AgenticBucketStatus{
				Status: oss.Ptr("Disabled"),
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, 200, putResult.StatusCode)

		_, err = client.DeleteAgenticBucket(context.TODO(), &DeleteAgenticBucketRequest{
			Bucket: oss.Ptr(bucket),
		})
		assert.NotNil(t, err)
		var serr *oss.ServiceError
		if assert.True(t, errors.As(err, &serr)) {
			assert.True(t, serr.StatusCode == 409 || serr.Code == "AgenticBucketNotReady",
				"expected AgenticBucketNotReady/409, got code=%q status=%d", serr.Code, serr.StatusCode)
		}
	})
}
