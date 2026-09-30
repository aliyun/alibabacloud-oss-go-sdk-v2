//go:build integration

package oss

import (
	"context"
	"errors"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketMetadataTableConfiguration(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(5)
	client := getDefaultClient()

	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	// Create
	createResult, err := client.CreateBucketMetadataTableConfiguration(context.TODO(), &CreateBucketMetadataTableConfigurationRequest{
		Bucket: Ptr(bucketName),
		MetadataConfiguration: &MetadataConfiguration{
			JournalTableConfiguration: &JournalTableConfiguration{
				RecordExpiration: &RecordExpiration{
					Expiration: Ptr("ENABLED"),
					Days:       Ptr(7),
				},
			},
			InventoryTableConfiguration: &InventoryTableConfiguration{
				ConfigurationState: Ptr("ENABLED"),
				EncryptionConfiguration: &MetadataTableEncryptionConfiguration{
					SseAlgorithm: Ptr("AES256"),
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// Get
	getResult, err := client.GetBucketMetadataTableConfiguration(context.TODO(), &GetBucketMetadataTableConfigurationRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotNil(t, getResult.GetBucketMetadataConfigurationResult)
	assert.NotNil(t, getResult.GetBucketMetadataConfigurationResult.MetadataConfigurationResult)
	metaResult := getResult.GetBucketMetadataConfigurationResult.MetadataConfigurationResult
	assert.NotNil(t, metaResult.DestinationResult)
	assert.Equal(t, "oss", *metaResult.DestinationResult.TableBucketType)
	assert.NotEmpty(t, *metaResult.DestinationResult.TableBucketArn)
	assert.True(t, strings.HasPrefix(*metaResult.DestinationResult.TableNamespace, "b_"))
	assert.NotNil(t, metaResult.JournalTableConfigurationResult)
	assert.Equal(t, "journal", *metaResult.JournalTableConfigurationResult.TableName)
	assert.Equal(t, "ENABLED", *metaResult.JournalTableConfigurationResult.RecordExpiration.Expiration)
	assert.Equal(t, int32(7), *metaResult.JournalTableConfigurationResult.RecordExpiration.Days)
	assert.NotNil(t, metaResult.InventoryTableConfigurationResult)
	assert.Equal(t, "ENABLED", *metaResult.InventoryTableConfigurationResult.ConfigurationState)
	assert.Equal(t, "inventory", *metaResult.InventoryTableConfigurationResult.TableName)
	time.Sleep(1 * time.Second)

	// Update inventory
	updateInvResult, err := client.UpdateBucketMetadataInventoryTableConfiguration(context.TODO(), &UpdateBucketMetadataInventoryTableConfigurationRequest{
		Bucket: Ptr(bucketName),
		InventoryTableConfiguration: &InventoryTableConfiguration{
			ConfigurationState: Ptr("DISABLED"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, updateInvResult.StatusCode)
	assert.NotEmpty(t, updateInvResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// Update journal
	updateJournalResult, err := client.UpdateBucketMetadataJournalTableConfiguration(context.TODO(), &UpdateBucketMetadataJournalTableConfigurationRequest{
		Bucket: Ptr(bucketName),
		JournalTableConfiguration: &JournalTableConfiguration{
			RecordExpiration: &RecordExpiration{
				Expiration: Ptr("DISABLED"),
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, updateJournalResult.StatusCode)
	assert.NotEmpty(t, updateJournalResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// Get after updates
	getResult2, err := client.GetBucketMetadataTableConfiguration(context.TODO(), &GetBucketMetadataTableConfigurationRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult2.StatusCode)
	assert.Equal(t, "DISABLED", *getResult2.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.InventoryTableConfigurationResult.ConfigurationState)
	assert.Equal(t, "DISABLED", *getResult2.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.RecordExpiration.Expiration)
	time.Sleep(1 * time.Second)

	// Delete
	delResult, err := client.DeleteBucketMetadataTableConfiguration(context.TODO(), &DeleteBucketMetadataTableConfigurationRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	// Get after delete
	_, err = client.GetBucketMetadataTableConfiguration(context.TODO(), &GetBucketMetadataTableConfigurationRequest{
		Bucket: Ptr(bucketName),
	})
	assert.NotNil(t, err)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "MetadataConfigurationNotFound", serr.Code)
	assert.NotEmpty(t, serr.RequestID)

	// test server error with invalid AK
	invalidAkClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))

	_, err = invalidAkClient.CreateBucketMetadataTableConfiguration(context.TODO(), &CreateBucketMetadataTableConfigurationRequest{
		Bucket: Ptr(bucketName),
		MetadataConfiguration: &MetadataConfiguration{
			JournalTableConfiguration: &JournalTableConfiguration{
				RecordExpiration: &RecordExpiration{
					Expiration: Ptr("ENABLED"),
					Days:       Ptr(7),
				},
			},
		},
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.GetBucketMetadataTableConfiguration(context.TODO(), &GetBucketMetadataTableConfigurationRequest{
		Bucket: Ptr(bucketName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.UpdateBucketMetadataInventoryTableConfiguration(context.TODO(), &UpdateBucketMetadataInventoryTableConfigurationRequest{
		Bucket: Ptr(bucketName),
		InventoryTableConfiguration: &InventoryTableConfiguration{
			ConfigurationState: Ptr("DISABLED"),
		},
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.UpdateBucketMetadataJournalTableConfiguration(context.TODO(), &UpdateBucketMetadataJournalTableConfigurationRequest{
		Bucket: Ptr(bucketName),
		JournalTableConfiguration: &JournalTableConfiguration{
			RecordExpiration: &RecordExpiration{
				Expiration: Ptr("DISABLED"),
			},
		},
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = invalidAkClient.DeleteBucketMetadataTableConfiguration(context.TODO(), &DeleteBucketMetadataTableConfigurationRequest{
		Bucket: Ptr(bucketName),
	})
	assert.NotNil(t, err)
	serr = nil
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "Forbidden", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}
