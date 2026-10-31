//go:build integration

package oss

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestBucketWorm(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	initWorm := &InitiateBucketWormRequest{
		Bucket: Ptr(bucketName),
		InitiateWormConfiguration: &InitiateWormConfiguration{
			Ptr(int32(1)),
		},
	}

	initResult, err := client.InitiateBucketWorm(context.TODO(), initWorm)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketWormRequest{
		Bucket: Ptr(bucketName),
	}

	getResult, err := client.GetBucketWorm(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *getResult.WormConfiguration.WormId, *initResult.WormId)
	assert.NotEmpty(t, *getResult.WormConfiguration.CreationDate)
	assert.NotEmpty(t, getResult.WormConfiguration.RetentionPeriodInDays)
	assert.NotEmpty(t, getResult.WormConfiguration.State)
	time.Sleep(1 * time.Second)

	abortRequest := &AbortBucketWormRequest{
		Bucket: Ptr(bucketName),
	}
	abortResult, err := client.AbortBucketWorm(context.TODO(), abortRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, abortResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	completeRequest := &CompleteBucketWormRequest{
		Bucket: Ptr(bucketName),
		WormId: initResult.WormId,
	}
	_, err = client.CompleteBucketWorm(context.TODO(), completeRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	extendRequest := &ExtendBucketWormRequest{
		Bucket: Ptr(bucketName),
		WormId: initResult.WormId,
	}
	_, err = client.ExtendBucketWorm(context.TODO(), extendRequest)
	assert.NotNil(t, err)
	assert.Equal(t, "missing required field, ExtendWormConfiguration.", err.Error())
	time.Sleep(1 * time.Second)

	extendRequest = &ExtendBucketWormRequest{
		Bucket: Ptr(bucketName),
		WormId: initResult.WormId,
		ExtendWormConfiguration: &ExtendWormConfiguration{
			Ptr(int32(2)),
		},
	}
	serr = &ServiceError{}
	_, err = client.ExtendBucketWorm(context.TODO(), extendRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketObjectWormConfiguration(t *testing.T) {
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	_, err = client.PutBucketVersioning(context.TODO(), &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	})
	assert.Nil(t, err)

	putRequest := &PutBucketObjectWormConfigurationRequest{
		Bucket: Ptr(bucketName),
		ObjectWormConfiguration: &ObjectWormConfiguration{
			ObjectWormEnabled: Ptr("Enabled"),
			Rule: &ObjectWormRule{
				DefaultRetention: &ObjectWormDefaultRetention{
					Mode:  Ptr("GOVERNANCE"),
					Years: Ptr(int32(1)),
				},
			},
		},
	}
	putResult, err := client.PutBucketObjectWormConfiguration(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketObjectWormConfigurationRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketObjectWormConfiguration(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketObjectWormConfigurationRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketObjectWormConfiguration(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketObjectWormConfigurationRequest{
		Bucket: Ptr(bucketNameNotExist),
		ObjectWormConfiguration: &ObjectWormConfiguration{
			ObjectWormEnabled: Ptr("Enabled"),
			Rule: &ObjectWormRule{
				DefaultRetention: &ObjectWormDefaultRetention{
					Mode:  Ptr("GOVERNANCE"),
					Years: Ptr(int32(1)),
				},
			},
		},
	}
	serr = &ServiceError{}
	putResult, err = client.PutBucketObjectWormConfiguration(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "MalformedXML", serr.Code)
	assert.Equal(t, "The XML you provided was not well-formed or did not validate against our published schema.", serr.Message)
	assert.Equal(t, "0015-00000231", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DeleteBucket(context.TODO(), &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
}

func TestObjectWorm(t *testing.T) {
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader("hi oss"),
	})
	assert.Nil(t, err)

	_, err = client.PutBucketVersioning(context.TODO(), &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	})
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	input := &OperationInput{
		Bucket: Ptr(bucketName),
		OpName: "PutBucketObjectWormConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: "application/xml",
		},
		Parameters: map[string]string{
			"objectWorm": "",
		},
		Body: io.NopCloser(bytes.NewReader([]byte(`<ObjectWormConfiguration>
			<ObjectWormEnabled>Enabled</ObjectWormEnabled>
				<Rule>
					<DefaultRetention>
						<Mode>COMPLIANCE</Mode>
						<Days>1</Days>
					</DefaultRetention>
				</Rule>
			</ObjectWormConfiguration>`))),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"objectWorm"})
	_, err = client.InvokeOperation(context.TODO(), input)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	date := time.Now().UTC().Add(3 * time.Second).Format("2006-01-02T15:04:05.000Z")
	putResult, err := client.PutObjectRetention(context.TODO(), &PutObjectRetentionRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Retention: &ObjectWormRetention{
			Mode:            Ptr("COMPLIANCE"),
			RetainUntilDate: Ptr(date),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetObjectRetention(context.TODO(), &GetObjectRetentionRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	putResult2, err := client.PutObjectLegalHold(context.TODO(), &PutObjectLegalHoldRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		LegalHold: &ObjectWormLegalHold{
			Status: Ptr("OFF"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult2.StatusCode)
	assert.NotEmpty(t, putResult2.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult2, err := client.GetObjectLegalHold(context.TODO(), &GetObjectLegalHoldRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult2.StatusCode)
	assert.NotEmpty(t, getResult2.LegalHold.Status, "OFF")
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	_, err = client.PutObjectRetention(context.TODO(), &PutObjectRetentionRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		Retention: &ObjectWormRetention{
			Mode:            Ptr("COMPLIANCE"),
			RetainUntilDate: Ptr(date),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	serr = &ServiceError{}
	_, err = client.GetObjectRetention(context.TODO(), &GetObjectRetentionRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	serr = &ServiceError{}
	_, err = client.PutObjectLegalHold(context.TODO(), &PutObjectLegalHoldRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		LegalHold: &ObjectWormLegalHold{
			Status: Ptr("ON"),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	serr = &ServiceError{}
	_, err = client.GetObjectLegalHold(context.TODO(), &GetObjectLegalHoldRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	time.Sleep(4 * time.Second)
	cleanObjects(client, bucketName, t)
}
