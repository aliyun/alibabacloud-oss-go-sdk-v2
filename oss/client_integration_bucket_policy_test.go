//go:build integration

package oss

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestInvokeOperation_BucketPolicy(t *testing.T) {
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

	_, err = client.PutBucketPublicAccessBlock(context.TODO(), &PutBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			BlockPublicAccess: Ptr(false),
		},
	})
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)

	calcMd5 := func(input string) string {
		if len(input) == 0 {
			return "1B2M2Y8AsgTpgAmY7PhCfg=="
		}
		h := md5.New()
		h.Write([]byte(input))
		return base64.StdEncoding.EncodeToString(h.Sum(nil))
	}

	// PutBucketPolicy
	policy := `{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Allow","Resource":["acs:oss:*:*:*"]}]}`
	input := &OperationInput{
		OpName: "PutBucketPolicy",
		Method: "PUT",
		Parameters: map[string]string{
			"policy": "",
		},
		// Add Content-md5
		Headers: map[string]string{
			"Content-MD5": calcMd5(policy),
		},
		Body:   strings.NewReader(policy),
		Bucket: Ptr(bucketName),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	output, err := client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)

	// GetBucketPolicy
	input = &OperationInput{
		OpName: "GetBucketPolicy",
		Method: "GET",
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: Ptr(bucketName),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	output, err = client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)
	policy1, err := io.ReadAll(output.Body)
	assert.NoError(t, err)
	if output.Body != nil {
		output.Body.Close()
	}
	assert.NotEmpty(t, policy1)

	// DeleteBucketPolicy
	input = &OperationInput{
		OpName: "DeleteBucketPolicy",
		Method: "DELETE",
		Parameters: map[string]string{
			"policy": "",
		},
		Bucket: Ptr(bucketName),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	output, err = client.InvokeOperation(context.TODO(), input)
	assert.NoError(t, err)
	// discard body
	_, err = io.ReadAll(output.Body)
	assert.NoError(t, err)
	if output.Body != nil {
		output.Body.Close()
	}
}

func TestBucketPolicy(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketPolicyRequest{
		Bucket: Ptr(bucketName),
		Body: strings.NewReader(`{
   "Version":"1",
   "Statement":[
   {
     "Action":[
       "oss:PutObject",
       "oss:GetObject"
    ],
    "Effect":"Deny",
    "Principal":["1234567890"],
    "Resource":["acs:oss:*:1234567890:*/*"]
   }
  ]
 }`),
	}

	putResult, err := client.PutBucketPolicy(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketPolicyRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketPolicy(context.TODO(), getRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, getResult.Body)

	statusRequest := &GetBucketPolicyStatusRequest{
		Bucket: Ptr(bucketName),
	}
	statusResult, err := client.GetBucketPolicyStatus(context.TODO(), statusRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, statusResult.StatusCode)
	assert.NotEmpty(t, statusResult.Headers.Get("X-Oss-Request-Id"))
	assert.False(t, *statusResult.PolicyStatus.IsPublic)

	delRequest := &DeleteBucketPolicyRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketPolicy(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketPolicyRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketPolicy(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	statusRequest = &GetBucketPolicyStatusRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	statusResult, err = client.GetBucketPolicyStatus(context.TODO(), statusRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketPolicyRequest{
		Bucket: Ptr(bucketNameNotExist),
		Body: strings.NewReader(`{
   "Version":"1",
   "Statement":[
   {
     "Action":[
       "oss:PutObject",
       "oss:GetObject"
    ],
    "Effect":"Deny",
    "Principal":["1234567890"],
    "Resource":["acs:oss:*:1234567890:*/*"]
   }
  ]
 }`),
	}
	serr = &ServiceError{}
	putResult, err = client.PutBucketPolicy(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest = &DeleteBucketPolicyRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	delResult, err = client.DeleteBucketPolicy(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
