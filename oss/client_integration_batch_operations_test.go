//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func TestBatchOperations(t *testing.T) {
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

	createResult, err := client.CreateJob(context.TODO(), &CreateJobRequest{
		CreateJob: &CreateJobConfig{
			Operation: &Operation{
				PutObjectAcl: &PutObjectAcl{
					ObjectAcl: ObjectACLPrivate,
				},
			},
			ClientRequestToken: Ptr(randLowStr(10)),
			Report: &Report{
				Bucket:  Ptr(bucketName),
				Enabled: Ptr(true),
			},
			KeyPrefixManifestGenerator: &KeyPrefixManifestGenerator{
				SourceBucket: Ptr(bucketName),
				Prefix:       Ptr("batch-manifests/"),
			},
			Priority: Ptr(int32(10)),
			RoleArn:  Ptr("acs:ram::" + accountID_ + ":role/oss-sdk-batch-test"),
		},
	})
	assert.Nil(t, err)

	describeResult, err := client.DescribeJob(context.TODO(), &DescribeJobRequest{
		BatchJobId: createResult.JobId,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, describeResult.StatusCode)
	assert.NotEmpty(t, describeResult.Headers.Get("X-Oss-Request-Id"))

	listResult, err := client.ListJobs(context.TODO(), &ListJobsRequest{
		MaxKeys: int32(1),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	assert.True(t, len(listResult.Jobs.JobListDescriptor) > 0)

	upResult, err := client.UpdateJobPriority(context.TODO(), &UpdateJobPriorityRequest{
		BatchJobId:     createResult.JobId,
		TargetPriority: Ptr(int32(1)),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, upResult.StatusCode)
	assert.NotEmpty(t, upResult.Headers.Get("X-Oss-Request-Id"))

	upResult2, err := client.UpdateJobStatus(context.TODO(), &UpdateJobStatusRequest{
		BatchJobId:         createResult.JobId,
		RequestedJobStatus: Ptr("Cancelled"),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, upResult2.StatusCode)
	assert.NotEmpty(t, upResult2.Headers.Get("X-Oss-Request-Id"))

	var serr *ServiceError
	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	_, err = noPermClient.CreateJob(context.TODO(), &CreateJobRequest{
		CreateJob: &CreateJobConfig{
			Operation: &Operation{
				PutObjectAcl: &PutObjectAcl{
					ObjectAcl: ObjectACLPrivate,
				},
			},
			ClientRequestToken: Ptr(randLowStr(10)),
			Report: &Report{
				Bucket:  Ptr(bucketName),
				Enabled: Ptr(true),
			},
			KeyPrefixManifestGenerator: &KeyPrefixManifestGenerator{
				SourceBucket: Ptr(bucketName),
				Prefix:       Ptr("batch-manifests/"),
			},
			Priority: Ptr(int32(10)),
			RoleArn:  Ptr("acs:ram::" + accountID_ + ":role/oss-sdk-batch-test"),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = noPermClient.DescribeJob(context.TODO(), &DescribeJobRequest{
		BatchJobId: createResult.JobId,
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = noPermClient.ListJobs(context.TODO(), &ListJobsRequest{
		MaxKeys: int32(1),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	_, err = noPermClient.UpdateJobPriority(context.TODO(), &UpdateJobPriorityRequest{
		BatchJobId:     createResult.JobId,
		TargetPriority: Ptr(int32(1)),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = noPermClient.UpdateJobStatus(context.TODO(), &UpdateJobStatusRequest{
		BatchJobId:         createResult.JobId,
		RequestedJobStatus: Ptr("Cancelled"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
}
