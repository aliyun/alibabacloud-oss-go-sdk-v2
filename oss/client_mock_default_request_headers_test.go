package oss

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func testDefaultRequestHeadersConfig(serverURL string, headers map[string]string) *Config {
	return LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint(serverURL).
		WithDefaultRequestHeaders(headers)
}

func TestMockDefaultRequestHeaders_AppliedToEveryRequest(t *testing.T) {
	var got []http.Header
	server := testSetupMockServer(t, 200,
		map[string]string{"x-oss-request-id": "534B371674E88A4D8906****"},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?><ListBucketResult></ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			got = append(got, r.Header.Clone())
		})
	defer server.Close()

	client := NewClient(testDefaultRequestHeadersConfig(server.URL, map[string]string{
		"x-my-trace-id":       "abc123",
		"x-oss-request-payer": "requester",
	}))

	_, err := client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   strings.NewReader("hi"),
	})
	assert.Nil(t, err)

	_, err = client.ListObjectsV2(context.TODO(), &ListObjectsV2Request{Bucket: Ptr("bucket")})
	assert.Nil(t, err)

	assert.Len(t, got, 2)
	for _, h := range got {
		assert.Equal(t, "abc123", h.Get("x-my-trace-id"))
		assert.Equal(t, "requester", h.Get("x-oss-request-payer"))
	}
}

func TestMockDefaultRequestHeaders_RequestWinsOnConflict(t *testing.T) {
	var got http.Header
	server := testSetupMockServer(t, 200, nil, nil,
		func(t *testing.T, r *http.Request) {
			got = r.Header.Clone()
		})
	defer server.Close()

	// "content-type" differs in case from the canonical "Content-Type" the
	// operation sets, to cover case-insensitive conflict detection.
	client := NewClient(testDefaultRequestHeadersConfig(server.URL, map[string]string{
		"content-type":        "application/from-default",
		"x-oss-storage-class": "IA",
		"x-my-trace-id":       "abc123",
	}))

	_, err := client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		ContentType:  Ptr("application/from-request"),
		StorageClass: StorageClassStandard,
		Body:         strings.NewReader("hi"),
	})
	assert.Nil(t, err)

	assert.Equal(t, "application/from-request", got.Get("Content-Type"))
	assert.Len(t, got.Values("Content-Type"), 1)
	assert.Equal(t, string(StorageClassStandard), got.Get("x-oss-storage-class"))
	assert.Len(t, got.Values("x-oss-storage-class"), 1)
	// the non-conflicting one still gets through
	assert.Equal(t, "abc123", got.Get("x-my-trace-id"))
}

func TestMockDefaultRequestHeaders_CanNotOverrideUserAgent(t *testing.T) {
	var got http.Header
	server := testSetupMockServer(t, 200, nil, nil,
		func(t *testing.T, r *http.Request) {
			got = r.Header.Clone()
		})
	defer server.Close()

	client := NewClient(testDefaultRequestHeadersConfig(server.URL, map[string]string{
		"User-Agent": "evil-ua",
	}))

	_, err := client.ListObjectsV2(context.TODO(), &ListObjectsV2Request{Bucket: Ptr("bucket")})
	assert.Nil(t, err)

	assert.NotEqual(t, "evil-ua", got.Get("User-Agent"))
	assert.True(t, strings.HasPrefix(got.Get("User-Agent"), SdkName+"/"))
}

func TestMockDefaultRequestHeaders_ContentLengthNotCorrupted(t *testing.T) {
	var gotLength int64
	var gotBody string
	server := testSetupMockServer(t, 200, nil, nil,
		func(t *testing.T, r *http.Request) {
			gotLength = r.ContentLength
			data := make([]byte, 64)
			n, _ := r.Body.Read(data)
			gotBody = string(data[:n])
		})
	defer server.Close()

	client := NewClient(testDefaultRequestHeadersConfig(server.URL, map[string]string{
		"Content-Length": "999",
	}))

	_, err := client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   strings.NewReader("hello"),
	})
	assert.Nil(t, err)

	assert.EqualValues(t, 5, gotLength)
	assert.Equal(t, "hello", gotBody)
}

func TestMockDefaultRequestHeaders_EmptyKeyOrValueIgnored(t *testing.T) {
	var got http.Header
	server := testSetupMockServer(t, 200, nil, nil,
		func(t *testing.T, r *http.Request) {
			got = r.Header.Clone()
		})
	defer server.Close()

	client := NewClient(testDefaultRequestHeadersConfig(server.URL, map[string]string{
		"":              "no-name",
		"x-my-empty":    "",
		"x-my-trace-id": "abc123",
	}))
	assert.Len(t, client.options.DefaultRequestHeaders, 1)

	_, err := client.ListObjectsV2(context.TODO(), &ListObjectsV2Request{Bucket: Ptr("bucket")})
	assert.Nil(t, err)

	assert.Empty(t, got.Get("x-my-empty"))
	assert.Equal(t, "abc123", got.Get("x-my-trace-id"))
}

func TestMockDefaultRequestHeaders_ConfigMapMutationAfterNewClient(t *testing.T) {
	var got http.Header
	server := testSetupMockServer(t, 200, nil, nil,
		func(t *testing.T, r *http.Request) {
			got = r.Header.Clone()
		})
	defer server.Close()

	headers := map[string]string{"x-my-trace-id": "abc123"}
	client := NewClient(testDefaultRequestHeadersConfig(server.URL, headers))

	headers["x-my-trace-id"] = "mutated"
	headers["x-my-added"] = "added"

	_, err := client.ListObjectsV2(context.TODO(), &ListObjectsV2Request{Bucket: Ptr("bucket")})
	assert.Nil(t, err)

	assert.Equal(t, "abc123", got.Get("x-my-trace-id"))
	assert.Empty(t, got.Get("x-my-added"))
}

func TestMockDefaultRequestHeaders_NotSetByDefault(t *testing.T) {
	var got http.Header
	server := testSetupMockServer(t, 200, nil, nil,
		func(t *testing.T, r *http.Request) {
			got = r.Header.Clone()
		})
	defer server.Close()

	client := NewClient(LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL))
	assert.Nil(t, client.options.DefaultRequestHeaders)

	_, err := client.ListObjectsV2(context.TODO(), &ListObjectsV2Request{Bucket: Ptr("bucket")})
	assert.Nil(t, err)
	assert.Empty(t, got.Get("x-my-trace-id"))
}

func TestMockDefaultRequestHeaders_TakePartInSignature(t *testing.T) {
	server := testSetupMockServer(t, 200, nil, nil, func(t *testing.T, r *http.Request) {})
	defer server.Close()

	// a fixed expiration makes the signing time, and therefore the signature,
	// deterministic across the calls being compared
	expiration := time.Now().Add(time.Hour)
	presign := func(headers map[string]string) *PresignResult {
		client := NewClient(testDefaultRequestHeadersConfig(server.URL, headers))
		result, err := client.Presign(context.TODO(),
			&GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")},
			PresignExpiration(expiration))
		assert.Nil(t, err)
		return result
	}
	signatureOf := func(result *PresignResult) string {
		u, err := url.Parse(result.URL)
		assert.Nil(t, err)
		sig := u.Query().Get("x-oss-signature")
		assert.NotEmpty(t, sig)
		return sig
	}

	// a signable default header changes the signature
	payer := presign(map[string]string{"x-oss-request-payer": "requester"})
	other := presign(map[string]string{"x-oss-request-payer": "somebody-else"})
	assert.NotEqual(t, signatureOf(payer), signatureOf(other))

	// and is reported back by Presign, so the caller can replay it
	assert.Equal(t, "requester", payer.SignedHeaders["X-Oss-Request-Payer"])

	// an unsignable default header does not
	traceA := presign(map[string]string{"x-my-trace-id": "a"})
	traceB := presign(map[string]string{"x-my-trace-id": "b"})
	assert.Equal(t, signatureOf(traceA), signatureOf(traceB))
	assert.Empty(t, traceA.SignedHeaders["X-My-Trace-Id"])
}
