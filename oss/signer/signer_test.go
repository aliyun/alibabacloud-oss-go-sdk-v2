package signer

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

func TestSigningContext(t *testing.T) {
	r := SigningContext{}
	assert.Empty(t, r.Product)
	assert.Empty(t, r.Region)
	assert.Empty(t, r.Bucket)
	assert.Empty(t, r.Key)
	assert.Nil(t, r.Request)
	assert.Empty(t, r.SubResource)

	assert.Empty(t, r.Credentials)
	assert.Empty(t, r.StringToSign)
	assert.Empty(t, r.SignedHeaders)
	assert.Empty(t, r.Time)
}

func TestNopSigner(t *testing.T) {
	r := NopSigner{}
	assert.Nil(t, r.Sign(context.TODO(), nil))
}

func TestV1AuthHeader(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())

	//case 1
	request, _ := http.NewRequest("PUT", "http://examplebucket.oss-cn-hangzhou.aliyuncs.com", nil)
	request.Header = http.Header{}
	request.Header.Add("Content-MD5", "eB5eJF1ptWaXm4bijSPyxw==")
	request.Header.Add("Content-Type", "text/html")
	request.Header.Add("x-oss-meta-author", "alice")
	request.Header.Add("x-oss-meta-magic", "abracadabra")
	request.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Key:         ptr("nelson"),
		Request:     request,
		Credentials: &cred,
		Time:        signTime,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)

	signToString := "PUT\neB5eJF1ptWaXm4bijSPyxw==\ntext/html\nWed, 28 Dec 2022 10:27:41 GMT\nx-oss-date:Wed, 28 Dec 2022 10:27:41 GMT\nx-oss-meta-author:alice\nx-oss-meta-magic:abracadabra\n/examplebucket/nelson"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS ak:kSHKmLxlyEAKtZPkJhG9bZb5k7M=", request.Header.Get("Authorization"))

	//case 2
	request, _ = http.NewRequest("PUT", "http://examplebucket.oss-cn-hangzhou.aliyuncs.com/?acl", nil)
	request.Header = http.Header{}
	request.Header.Add("Content-MD5", "eB5eJF1ptWaXm4bijSPyxw==")
	request.Header.Add("Content-Type", "text/html")
	request.Header.Add("x-oss-meta-author", "alice")
	request.Header.Add("x-oss-meta-magic", "abracadabra")
	request.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Key:         ptr("nelson"),
		Request:     request,
		Credentials: &cred,
		Time:        signTime,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "PUT\neB5eJF1ptWaXm4bijSPyxw==\ntext/html\nWed, 28 Dec 2022 10:27:41 GMT\nx-oss-date:Wed, 28 Dec 2022 10:27:41 GMT\nx-oss-meta-author:alice\nx-oss-meta-magic:abracadabra\n/examplebucket/nelson?acl"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS ak:/afkugFbmWDQ967j1vr6zygBLQk=", request.Header.Get("Authorization"))

	//case 3
	request, _ = http.NewRequest("GET", "http://examplebucket.oss-cn-hangzhou.aliyuncs.com/?resourceGroup&non-resousce=null", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Request:     request,
		Credentials: &cred,
		SubResource: []string{"resourceGroup"},
		Time:        signTime,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "GET\n\n\nWed, 28 Dec 2022 10:27:41 GMT\nx-oss-date:Wed, 28 Dec 2022 10:27:41 GMT\n/examplebucket/?resourceGroup"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS ak:vkQmfuUDyi1uDi3bKt67oemssIs=", request.Header.Get("Authorization"))

	//case 4
	request, _ = http.NewRequest("GET", "http://examplebucket.oss-cn-hangzhou.aliyuncs.com/?resourceGroup&acl", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Request:     request,
		Credentials: &cred,
		SubResource: []string{"resourceGroup"},
		Time:        signTime,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "GET\n\n\nWed, 28 Dec 2022 10:27:41 GMT\nx-oss-date:Wed, 28 Dec 2022 10:27:41 GMT\n/examplebucket/?acl&resourceGroup"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS ak:x3E5TgOvl/i7PN618s5mEvpJDYk=", request.Header.Get("Authorization"))
}

func TestV1InvalidArgument(t *testing.T) {
	signer := &SignerV1{}
	signCtx := &SigningContext{}
	err := signer.Sign(context.TODO(), signCtx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext.Credentials is null or empty")

	provider := credentials.NewStaticCredentialsProvider("", "sk")
	cred, _ := provider.GetCredentials(context.TODO())
	signCtx = &SigningContext{
		Credentials: &cred,
	}
	err = signer.Sign(context.TODO(), signCtx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext.Credentials is null or empty")

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())
	signCtx = &SigningContext{
		Credentials: &cred,
	}
	err = signer.Sign(context.TODO(), signCtx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext.Request is null")

	err = signer.Sign(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext is null")
}

func TestV1AuthQuery(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	//case 1
	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())
	request, _ := http.NewRequest("GET", "http://bucket.oss-cn-hangzhou.aliyuncs.com/key?versionId=versionId", nil)
	request.Header = http.Header{}
	signTime, _ = http.ParseTime("Sun, 12 Nov 2023 16:43:40 GMT")

	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("key"),
		Request:         request,
		Credentials:     &cred,
		Time:            signTime,
		AuthMethodQuery: true,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)
	signUrl := "http://bucket.oss-cn-hangzhou.aliyuncs.com/key?Expires=1699807420&OSSAccessKeyId=ak&Signature=dcLTea%2BYh9ApirQ8o8dOPqtvJXQ%3D&versionId=versionId"
	assert.Equal(t, signUrl, request.URL.String())

	//case 2
	provider = credentials.NewStaticCredentialsProvider("ak", "sk", "token")
	cred, _ = provider.GetCredentials(context.TODO())
	request, _ = http.NewRequest("GET", "http://bucket.oss-cn-hangzhou.aliyuncs.com/key%2B123?versionId=versionId", nil)
	request.Header = http.Header{}
	signTime, _ = http.ParseTime("Sun, 12 Nov 2023 16:56:44 GMT")
	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("key+123"),
		Request:         request,
		Credentials:     &cred,
		Time:            signTime,
		AuthMethodQuery: true,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)
	signUrl = "http://bucket.oss-cn-hangzhou.aliyuncs.com/key%2B123?Expires=1699808204&OSSAccessKeyId=ak&Signature=jzKYRrM5y6Br0dRFPaTGOsbrDhY%3D&security-token=token&versionId=versionId"
	assert.Equal(t, signUrl, request.URL.String())
}

func TestV1AuthQuery1(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	//case 1
	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())
	request, _ := http.NewRequest("GET", "http://bucket.oss-cn-hangzhou.aliyuncs.com/key?versionId=versionId", nil)
	request.Header = http.Header{}
	signTime, _ = http.ParseTime("Sun, 12 Nov 2023 16:43:40 GMT")

	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("key"),
		Request:         request,
		Credentials:     &cred,
		Time:            signTime,
		AuthMethodQuery: true,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)
	signUrl := "http://bucket.oss-cn-hangzhou.aliyuncs.com/key?Expires=1699807420&OSSAccessKeyId=ak&Signature=dcLTea%2BYh9ApirQ8o8dOPqtvJXQ%3D&versionId=versionId"
	assert.Equal(t, signUrl, request.URL.String())

	//case 2
	provider = credentials.NewStaticCredentialsProvider("ak", "sk", "token")
	cred, _ = provider.GetCredentials(context.TODO())
	request, _ = http.NewRequest("GET", "http://bucket.oss-cn-hangzhou.aliyuncs.com/key%2B123?versionId=versionId", nil)
	request.Header = http.Header{}
	signTime, _ = http.ParseTime("Sun, 12 Nov 2023 16:56:44 GMT")
	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("key+123"),
		Request:         request,
		Credentials:     &cred,
		Time:            signTime,
		AuthMethodQuery: true,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)
	signUrl = "http://bucket.oss-cn-hangzhou.aliyuncs.com/key%2B123?Expires=1699808204&OSSAccessKeyId=ak&Signature=jzKYRrM5y6Br0dRFPaTGOsbrDhY%3D&security-token=token&versionId=versionId"
	assert.Equal(t, signUrl, request.URL.String())
}

func TestV4InvalidArgument(t *testing.T) {
	signer := &SignerV4{}
	signCtx := &SigningContext{}
	err := signer.Sign(context.TODO(), signCtx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext.Credentials is null or empty")

	provider := credentials.NewStaticCredentialsProvider("", "sk")
	cred, _ := provider.GetCredentials(context.TODO())
	signCtx = &SigningContext{
		Credentials: &cred,
	}
	err = signer.Sign(context.TODO(), signCtx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext.Credentials is null or empty")

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())
	signCtx = &SigningContext{
		Credentials: &cred,
	}
	err = signer.Sign(context.TODO(), signCtx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext.Request is null")

	err = signer.Sign(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext is null")
}

func TestV4AuthHeader(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())

	//case 1
	request, _ := http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "text/plain")
	request.Header.Add("x-oss-content-sha256", "UNSIGNED-PAYLOAD")
	signTime = time.Unix(1702743657, 0).UTC()
	signCtx = &SigningContext{
		Bucket:      ptr("bucket"),
		Key:         ptr("1234+-/123/1.txt"),
		Request:     request,
		Credentials: &cred,
		Product:     ptr("oss"),
		Region:      ptr("cn-hangzhou"),
		Time:        signTime,
	}

	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")

	request.URL.RawQuery = values.Encode()

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	authPat := "OSS4-HMAC-SHA256 Credential=ak/20231216/cn-hangzhou/oss/aliyun_v4_request,Signature=e21d18daa82167720f9b1047ae7e7f1ce7cb77a31e8203a7d5f4624fa0284afe"
	//assert.Equal(t, signToString, signCtx.StringToSign)
	//assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, authPat, request.Header.Get("Authorization"))
}

func TestV4AuthHeaderWithCloudBox(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())

	request, _ := http.NewRequest("PUT", "http://bucket.cb-123.cn-hangzhou.oss-cloudbox.aliyuncs.com/1234+-/123/1.txt", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "text/plain")
	request.Header.Add("x-oss-content-sha256", "UNSIGNED-PAYLOAD")
	signTime = time.Unix(1702743657, 0).UTC()

	signCtx = &SigningContext{
		Bucket:      ptr("bucket"),
		Key:         ptr("1234+-/123/1.txt"),
		Request:     request,
		Credentials: &cred,
		Product:     ptr("oss-cloudbox"),
		Region:      ptr("cb-123"),
		Time:        signTime,
	}

	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")
	request.URL.RawQuery = values.Encode()
	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	authPat := "OSS4-HMAC-SHA256 Credential=ak/20231216/cb-123/oss-cloudbox/aliyun_v4_request,Signature=94ce1f12c17d148ea681030275a94449d3357f5b5b21133996eec80af3e08a43"

	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, authPat, request.Header.Get("Authorization"))
}

func TestV4AuthHeaderToken(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk", "token")
	cred, _ = provider.GetCredentials(context.TODO())

	//case 1
	request, _ := http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "text/plain")
	request.Header.Add("x-oss-content-sha256", "UNSIGNED-PAYLOAD")
	signTime = time.Unix(1702784856, 0).UTC()
	signCtx = &SigningContext{
		Bucket:      ptr("bucket"),
		Key:         ptr("1234+-/123/1.txt"),
		Request:     request,
		Credentials: &cred,
		Product:     ptr("oss"),
		Region:      ptr("cn-hangzhou"),
		Time:        signTime,
	}

	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")

	request.URL.RawQuery = values.Encode()

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	authPat := "OSS4-HMAC-SHA256 Credential=ak/20231217/cn-hangzhou/oss/aliyun_v4_request,Signature=b94a3f999cf85bcdc00d332fbd3734ba03e48382c36fa4d5af5df817395bd9ea"
	assert.Equal(t, authPat, request.Header.Get("Authorization"))
}

func TestV4AuthHeaderWithAdditionalHeaders(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())

	//case 1
	request, _ := http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "text/plain")
	request.Header.Add("x-oss-content-sha256", "UNSIGNED-PAYLOAD")
	signTime = time.Unix(1702747512, 0).UTC()
	signCtx = &SigningContext{
		Bucket:            ptr("bucket"),
		Key:               ptr("1234+-/123/1.txt"),
		Request:           request,
		Credentials:       &cred,
		Product:           ptr("oss"),
		Region:            ptr("cn-hangzhou"),
		Time:              signTime,
		AdditionalHeaders: []string{"ZAbc", "abc"},
	}

	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")

	request.URL.RawQuery = values.Encode()

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	authPat := "OSS4-HMAC-SHA256 Credential=ak/20231216/cn-hangzhou/oss/aliyun_v4_request,AdditionalHeaders=abc;zabc,Signature=4a4183c187c07c8947db7620deb0a6b38d9fbdd34187b6dbaccb316fa251212f"
	assert.Equal(t, authPat, request.Header.Get("Authorization"))

	// with default signed header
	request, _ = http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "text/plain")
	request.Header.Add("x-oss-content-sha256", "UNSIGNED-PAYLOAD")
	signTime = time.Unix(1702747512, 0).UTC()
	signCtx = &SigningContext{
		Bucket:            ptr("bucket"),
		Key:               ptr("1234+-/123/1.txt"),
		Request:           request,
		Credentials:       &cred,
		Product:           ptr("oss"),
		Region:            ptr("cn-hangzhou"),
		Time:              signTime,
		AdditionalHeaders: []string{"x-oss-no-exist", "ZAbc", "x-oss-head1", "abc"},
	}

	values = url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")

	request.URL.RawQuery = values.Encode()

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	authPat = "OSS4-HMAC-SHA256 Credential=ak/20231216/cn-hangzhou/oss/aliyun_v4_request,AdditionalHeaders=abc;zabc,Signature=4a4183c187c07c8947db7620deb0a6b38d9fbdd34187b6dbaccb316fa251212f"
	assert.Equal(t, authPat, request.Header.Get("Authorization"))
}

func TestV4AuthQuery(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())

	//case 1
	request, _ := http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "application/octet-stream")

	signTime = time.Unix(1702781677, 0)
	time := time.Unix(1702782276, 0)
	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("1234+-/123/1.txt"),
		Request:         request,
		Credentials:     &cred,
		Product:         ptr("oss"),
		Region:          ptr("cn-hangzhou"),
		AuthMethodQuery: true,
		Time:            time,
		signTime:        &signTime,
	}

	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")

	request.URL.RawQuery = values.Encode()

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	querys := signCtx.Request.URL.Query()

	assert.Equal(t, "OSS4-HMAC-SHA256", querys.Get("x-oss-signature-version"))
	assert.Equal(t, "599", querys.Get("x-oss-expires"))
	assert.Equal(t, "ak/20231217/cn-hangzhou/oss/aliyun_v4_request", querys.Get("x-oss-credential"))
	assert.Equal(t, "a39966c61718be0d5b14e668088b3fa07601033f6518ac7b523100014269c0fe", querys.Get("x-oss-signature"))
	assert.Equal(t, "", querys.Get("x-oss-additional-headers"))
}

func TestV4AuthQueryToken(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk", "token")
	cred, _ = provider.GetCredentials(context.TODO())

	request, _ := http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "application/octet-stream")

	signTime = time.Unix(1702785388, 0)
	time := time.Unix(1702785987, 0)
	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("1234+-/123/1.txt"),
		Request:         request,
		Credentials:     &cred,
		Product:         ptr("oss"),
		Region:          ptr("cn-hangzhou"),
		AuthMethodQuery: true,
		Time:            time,
		signTime:        &signTime,
	}

	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")

	request.URL.RawQuery = values.Encode()

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	querys := signCtx.Request.URL.Query()

	assert.Equal(t, "OSS4-HMAC-SHA256", querys.Get("x-oss-signature-version"))
	assert.Equal(t, "20231217T035628Z", querys.Get("x-oss-date"))
	assert.Equal(t, "599", querys.Get("x-oss-expires"))
	assert.Equal(t, "ak/20231217/cn-hangzhou/oss/aliyun_v4_request", querys.Get("x-oss-credential"))
	assert.Equal(t, "3817ac9d206cd6dfc90f1c09c00be45005602e55898f26f5ddb06d7892e1f8b5", querys.Get("x-oss-signature"))
	assert.Equal(t, "", querys.Get("x-oss-additional-headers"))
}

func TestV4AuthQueryWithAdditionalHeaders(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())

	//case 1
	request, _ := http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "application/octet-stream")

	signTime = time.Unix(1702783809, 0)
	//time := time.Unix(1702784408, 0)
	signCtx = &SigningContext{
		Bucket:            ptr("bucket"),
		Key:               ptr("1234+-/123/1.txt"),
		Request:           request,
		Credentials:       &cred,
		Product:           ptr("oss"),
		Region:            ptr("cn-hangzhou"),
		AuthMethodQuery:   true,
		Time:              time.Unix(1702784408, 0),
		signTime:          &signTime,
		AdditionalHeaders: []string{"ZAbc", "abc"},
	}

	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")

	request.URL.RawQuery = values.Encode()

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	querys := signCtx.Request.URL.Query()

	assert.Equal(t, "OSS4-HMAC-SHA256", querys.Get("x-oss-signature-version"))
	assert.Equal(t, "20231217T033009Z", querys.Get("x-oss-date"))
	assert.Equal(t, "599", querys.Get("x-oss-expires"))
	assert.Equal(t, "ak/20231217/cn-hangzhou/oss/aliyun_v4_request", querys.Get("x-oss-credential"))
	assert.Equal(t, "6bd984bfe531afb6db1f7550983a741b103a8c58e5e14f83ea474c2322dfa2b7", querys.Get("x-oss-signature"))
	assert.Equal(t, "abc;zabc", querys.Get("x-oss-additional-headers"))

	// with default signed header
	request, _ = http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "application/octet-stream")

	signTime = time.Unix(1702783809, 0)
	//time = time.Unix(1702784408, 0)
	signCtx = &SigningContext{
		Bucket:            ptr("bucket"),
		Key:               ptr("1234+-/123/1.txt"),
		Request:           request,
		Credentials:       &cred,
		Product:           ptr("oss"),
		Region:            ptr("cn-hangzhou"),
		AuthMethodQuery:   true,
		Time:              time.Unix(1702784408, 0),
		signTime:          &signTime,
		AdditionalHeaders: []string{"x-oss-no-exist", "abc", "x-oss-head1", "ZAbc"},
	}

	values = url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")

	request.URL.RawQuery = values.Encode()

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	querys = signCtx.Request.URL.Query()

	assert.Equal(t, "OSS4-HMAC-SHA256", querys.Get("x-oss-signature-version"))
	assert.Equal(t, "20231217T033009Z", querys.Get("x-oss-date"))
	assert.Equal(t, "599", querys.Get("x-oss-expires"))
	assert.Equal(t, "ak/20231217/cn-hangzhou/oss/aliyun_v4_request", querys.Get("x-oss-credential"))
	assert.Equal(t, "6bd984bfe531afb6db1f7550983a741b103a8c58e5e14f83ea474c2322dfa2b7", querys.Get("x-oss-signature"))
	assert.Equal(t, "abc;zabc", querys.Get("x-oss-additional-headers"))
}

func TestV4AuthQueryWithCloudBox(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())

	//case 1
	request, _ := http.NewRequest("PUT", "http://bucket.cb-123.cn-hangzhou.oss-cloudbox.aliyuncs.com/1234+-/123/1.txt", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "application/octet-stream")

	signTime = time.Unix(1702781677, 0)
	//time := time.Unix(1702784408, 0)
	signCtx = &SigningContext{
		Bucket:            ptr("bucket"),
		Key:               ptr("1234+-/123/1.txt"),
		Request:           request,
		Credentials:       &cred,
		Product:           ptr("oss-cloudbox"),
		Region:            ptr("cb-123"),
		AuthMethodQuery:   true,
		Time:              time.Unix(1702782276, 0),
		signTime:          &signTime,
		AdditionalHeaders: []string{"ZAbc", "abc"},
	}

	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")

	request.URL.RawQuery = values.Encode()

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	querys := signCtx.Request.URL.Query()
	assert.Equal(t, "OSS4-HMAC-SHA256", querys.Get("x-oss-signature-version"))
	assert.Equal(t, "20231217T025437Z", querys.Get("x-oss-date"))
	assert.Equal(t, "599", querys.Get("x-oss-expires"))
	assert.Equal(t, "ak/20231217/cb-123/oss-cloudbox/aliyun_v4_request", querys.Get("x-oss-credential"))
	assert.Equal(t, "07284191b9b4978ac3520cd39ee2dea2747eda454089359371ff463a6c7ba20f", querys.Get("x-oss-signature"))
	assert.Equal(t, "abc;zabc", querys.Get("x-oss-additional-headers"))

	// with default signed header
	request, _ = http.NewRequest("PUT", "http://bucket.cb-123.cn-hangzhou.oss-cloudbox.aliyuncs.com/1234+-/123/1.txt", nil)
	request.Header = http.Header{}
	request.Header.Add("x-oss-head1", "value")
	request.Header.Add("abc", "value")
	request.Header.Add("ZAbc", "value")
	request.Header.Add("XYZ", "value")
	request.Header.Add("content-type", "application/octet-stream")

	signTime = time.Unix(1702783809, 0)
	//time = time.Unix(1702784408, 0)
	signCtx = &SigningContext{
		Bucket:            ptr("bucket"),
		Key:               ptr("1234+-/123/1.txt"),
		Request:           request,
		Credentials:       &cred,
		Product:           ptr("oss-cloudbox"),
		Region:            ptr("cb-123"),
		AuthMethodQuery:   true,
		Time:              time.Unix(1702784408, 0),
		signTime:          &signTime,
		AdditionalHeaders: []string{"x-oss-no-exist", "abc", "x-oss-head1", "ZAbc"},
	}

	values = url.Values{}
	values.Add("param1", "value1")
	values.Add("+param1", "value3")
	values.Add("|param1", "value4")
	values.Add("+param2", "")
	values.Add("|param2", "")
	values.Add("param2", "")
	request.URL.RawQuery = values.Encode()
	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	querys = signCtx.Request.URL.Query()

	assert.Equal(t, "OSS4-HMAC-SHA256", querys.Get("x-oss-signature-version"))
	assert.Equal(t, "20231217T033009Z", querys.Get("x-oss-date"))
	assert.Equal(t, "599", querys.Get("x-oss-expires"))
	assert.Equal(t, "ak/20231217/cb-123/oss-cloudbox/aliyun_v4_request", querys.Get("x-oss-credential"))
	assert.Equal(t, "16782cc8a7a554523db055eb804b508522e7e370073108ad88ee2f47496701dd", querys.Get("x-oss-signature"))
	assert.Equal(t, "abc;zabc", querys.Get("x-oss-additional-headers"))
}
