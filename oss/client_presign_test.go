package oss

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func TestPresignPresignOptions(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV1)

	client := NewClient(cfg)

	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.Equal(t, expiration, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")

	expires := 50 * time.Minute
	expiration = time.Now().Add(expires)
	result, err = client.Presign(context.TODO(), request, PresignExpires(expires))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.True(t, result.Expiration.Unix()-expiration.Unix() < 2)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", result.Expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")

	cfgV4 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV4)

	client = NewClient(cfgV4)

	request = &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}
	expiration = time.Now().Add(1 * time.Hour)
	result, err = client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.Equal(t, expiration, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-date=%v", expiration.Add(-1*time.Hour).UTC().Format("20060102T150405Z")))
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-expires=%v", (1*time.Hour).Seconds()))
	assert.Contains(t, result.URL, "x-oss-signature=")
	credential := fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, result.URL, "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, result.URL, "x-oss-signature-version=OSS4-HMAC-SHA256")

	expires = 50 * time.Minute
	expiration = time.Now().Add(expires)
	result, err = client.Presign(context.TODO(), request, PresignExpires(expires))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.True(t, result.Expiration.Unix()-expiration.Unix() < 2)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-date=%v", expiration.Add(-50*time.Minute).UTC().Format("20060102T150405Z")))
	credential = fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, result.URL, "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, result.URL, "x-oss-signature=")
	assert.Contains(t, result.URL, "x-oss-credential="+url.QueryEscape("ak/"))
	assert.Contains(t, result.URL, "x-oss-signature-version=OSS4-HMAC-SHA256")
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-expires=%v", (50*time.Minute).Seconds()))
}

func TestPresignWithToken(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk", "token")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV1)

	client := NewClient(cfg)

	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")
	assert.Contains(t, result.URL, "security-token=token")

	cfgV4 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk", "token")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV4)

	client = NewClient(cfgV4)

	request = &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}

	expiration = time.Now().Add(1 * time.Hour)
	result, err = client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "x-oss-security-token=token")
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-date=%v", expiration.Add(-1*time.Hour).UTC().Format("20060102T150405Z")))
	credential := fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, result.URL, "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, result.URL, "x-oss-signature=")
	assert.Contains(t, result.URL, "x-oss-signature-version=OSS4-HMAC-SHA256")
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-expires=%v", (1*time.Hour).Seconds()))
}

func TestPresignWithHeader(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV1)

	client := NewClient(cfg)

	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"Content-Type": "application/octet-stream",
			},
		},
	}

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Len(t, result.SignedHeaders, 1)
	assert.Equal(t, "application/octet-stream", result.SignedHeaders["Content-Type"])
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")

	cfgV4 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV4)

	client = NewClient(cfgV4)
	request = &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"Content-Type": "application/octet-stream",
			},
		},
	}
	expiration = time.Now().Add(1 * time.Hour)
	result, err = client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	//fmt.Printf("result.SignedHeaders:%#v\n", result.SignedHeaders)
	assert.Len(t, result.SignedHeaders, 1)
	assert.Equal(t, "application/octet-stream", result.SignedHeaders["Content-Type"])
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-date=%v", expiration.Add(-1*time.Hour).UTC().Format("20060102T150405Z")))
	credential := fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, result.URL, "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, result.URL, "x-oss-signature=")
	assert.Contains(t, result.URL, "x-oss-signature-version=OSS4-HMAC-SHA256")
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-expires=%v", (1*time.Hour).Seconds()))
}

func TestPresignWithAdditionalHeaders(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV1).
		WithAdditionalHeaders([]string{"email", "name"})

	client := NewClient(cfg)

	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"Content-Type": "application/octet-stream",
				"email":        "demo@aliyun.com",
				"name":         "aliyun",
			},
		},
	}

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Len(t, result.SignedHeaders, 1)
	assert.Equal(t, "application/octet-stream", result.SignedHeaders["Content-Type"])
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")

	cfgV4 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV4).
		WithAdditionalHeaders([]string{"email", "name"})

	client = NewClient(cfgV4)
	request = &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"email": "demo@aliyun.com",
				"name":  "aliyun",
			},
		},
	}
	expiration = time.Now().Add(1 * time.Hour)
	result, err = client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Len(t, result.SignedHeaders, 2)
	assert.Equal(t, "demo@aliyun.com", result.SignedHeaders["Email"])
	assert.Equal(t, "aliyun", result.SignedHeaders["Name"])
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-date=%v", expiration.Add(-1*time.Hour).UTC().Format("20060102T150405Z")))
	credential := fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, result.URL, "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, result.URL, "x-oss-signature=")
	assert.Contains(t, result.URL, "x-oss-signature-version=OSS4-HMAC-SHA256")
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-expires=%v", (1*time.Hour).Seconds()))
	assert.Contains(t, result.URL, "x-oss-additional-headers=email%3Bname")

}

func TestPresignWithQuery(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV1)

	client := NewClient(cfg)

	reqeust := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"x-oss-process": "abc",
			},
		},
	}

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), reqeust, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")
	assert.Contains(t, result.URL, "x-oss-process=abc")

	cfgV4 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV4)

	client = NewClient(cfgV4)

	reqeust = &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"x-oss-process": "abc",
			},
		},
	}
	expiration = time.Now().Add(1 * time.Hour)
	result, err = client.Presign(context.TODO(), reqeust, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-date=%v", expiration.Add(-1*time.Hour).UTC().Format("20060102T150405Z")))
	assert.Contains(t, result.URL, "x-oss-signature=")
	credential := fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, result.URL, "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, result.URL, "x-oss-signature-version=OSS4-HMAC-SHA256")
	assert.Contains(t, result.URL, "x-oss-process=abc")
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-expires=%v", (1*time.Hour).Seconds()))
}

func TestPresignOperationInput(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV1)

	client := NewClient(cfg)

	request := &GetObjectRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("key"),
		VersionId: Ptr("versionId"),
	}

	expiration, _ := http.ParseTime("Sun, 12 Nov 2023 16:43:40 GMT")
	result, err := client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, "Expires=1699807420")
	assert.Contains(t, result.URL, "Signature=dcLTea%2BYh9ApirQ8o8dOPqtvJXQ%3D")

	//token
	cfg = LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk", "token")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV1)

	client = NewClient(cfg)

	request = &GetObjectRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("key+123"),
		VersionId: Ptr("versionId"),
	}

	expiration, _ = http.ParseTime("Sun, 12 Nov 2023 16:56:44 GMT")
	result, err = client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key%2B123?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, "Expires=1699808204")
	assert.Contains(t, result.URL, "Signature=jzKYRrM5y6Br0dRFPaTGOsbrDhY%3D")
	assert.Contains(t, result.URL, "security-token=token")
	assert.Contains(t, result.URL, "versionId=versionId")

	cfgV4 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV4)

	client = NewClient(cfgV4)

	request = &GetObjectRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("key"),
		VersionId: Ptr("versionId"),
	}

	expiration, _ = http.ParseTime("Sun, 12 Nov 2023 16:43:40 GMT")
	currentTime := time.Now()
	result, err = client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-date=%v", currentTime.UTC().Format("20060102T150405Z")))
	assert.Contains(t, result.URL, "x-oss-signature=")
	credential := fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", currentTime.UTC().Format("20060102"))
	assert.Contains(t, result.URL, "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, result.URL, "x-oss-signature-version=OSS4-HMAC-SHA256")
	diff := expiration.Unix() - currentTime.Unix()
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-expires=%v", diff))
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "versionId=versionId")

	cfgV4 = LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk", "token")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV4)
	client = NewClient(cfgV4)

	request = &GetObjectRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("key+123"),
		VersionId: Ptr("versionId"),
	}

	expiration, _ = http.ParseTime("Sun, 12 Nov 2023 16:56:44 GMT")
	currentTime = time.Now()
	result, err = client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key%2B123?")
	assert.Contains(t, result.URL, "x-oss-security-token=token")
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-date=%v", currentTime.UTC().Format("20060102T150405Z")))
	assert.Contains(t, result.URL, "x-oss-signature=")
	credential = fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", currentTime.UTC().Format("20060102"))
	assert.Contains(t, result.URL, "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, result.URL, "x-oss-signature-version=OSS4-HMAC-SHA256")
	diff = expiration.Unix() - currentTime.Unix()
	assert.Contains(t, result.URL, fmt.Sprintf("x-oss-expires=%v", diff))
	assert.Contains(t, result.URL, "versionId=versionId")
}

func TestPresignWithError(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV1)

	client := NewClient(cfg)

	// unsupport request
	request := &ListObjectsRequest{
		Bucket: Ptr("bucket"),
	}
	_, err := client.Presign(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request *oss.ListObjectsRequest")

	// request is nil
	_, err = client.Presign(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field, request")

	getRequest := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key+123"),
	}
	_, err = client.Presign(context.TODO(), getRequest, PresignExpires(8*24*time.Hour))
	assert.Nil(t, err)

	cfg = LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com").
		WithSignatureVersion(SignatureVersionV4)

	client = NewClient(cfg)

	// unsupport request
	request = &ListObjectsRequest{
		Bucket: Ptr("bucket"),
	}
	_, err = client.Presign(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request *oss.ListObjectsRequest")

	// request is nil
	_, err = client.Presign(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field, request")

	getRequest = &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key+123"),
	}
	_, err = client.Presign(context.TODO(), getRequest, PresignExpires(8*24*time.Hour))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "expires should be not greater than 604800(seven days)")
}
