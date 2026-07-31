package oss

import (
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestCopierClientCopierOptions(t *testing.T) {

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou")

	client := NewClient(cfg)

	// Default
	c := NewCopier(client)
	assert.Equal(t, DefaultCopyParallel, c.options.ParallelNum)
	assert.Equal(t, DefaultCopyPartSize, c.options.PartSize)
	assert.Equal(t, DefaultCopyThreshold, c.options.MultipartCopyThreshold)
	assert.Equal(t, false, c.options.LeavePartsOnError)
	assert.Equal(t, false, c.options.DisableShallowCopy)
	assert.Equal(t, 0, len(c.options.ClientOptions))
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

	// Set From Client
	c = NewCopier(client, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = true
		co.ClientOptions = []func(do *Options){func(do *Options) {}}
		co.MetadataProperties = &HeadObjectResult{}
		co.TagProperties = &GetObjectTaggingResult{}
	})
	assert.Equal(t, int(2), c.options.ParallelNum)
	assert.Equal(t, int64(1024*1024), c.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), c.options.MultipartCopyThreshold)
	assert.Equal(t, true, c.options.LeavePartsOnError)
	assert.Equal(t, true, c.options.DisableShallowCopy)
	assert.Equal(t, 1, len(c.options.ClientOptions))
	// only supports setting from c.Copy
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

	// Use WithXXX
	c = NewCopier(client, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = true
		co.ClientOptions = []func(do *Options){func(do *Options) {}}
		co.MetadataProperties = &HeadObjectResult{}
		co.TagProperties = &GetObjectTaggingResult{}
	},
		WithCopierParallelNum(5), WithCopierPartSize(2*1024*1024))
	assert.Equal(t, int(5), c.options.ParallelNum)
	assert.Equal(t, int64(2*1024*1024), c.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), c.options.MultipartCopyThreshold)
	assert.Equal(t, true, c.options.LeavePartsOnError)
	assert.Equal(t, true, c.options.DisableShallowCopy)
	assert.Equal(t, 1, len(c.options.ClientOptions))
	// only supports setting from c.Copy
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

}

func TestCopierApiCopierOptions(t *testing.T) {

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou")

	client := NewClient(cfg)

	// Default
	c := NewCopier(client)
	assert.Equal(t, DefaultCopyParallel, c.options.ParallelNum)
	assert.Equal(t, DefaultCopyPartSize, c.options.PartSize)
	assert.Equal(t, DefaultCopyThreshold, c.options.MultipartCopyThreshold)
	assert.Equal(t, false, c.options.LeavePartsOnError)
	assert.Equal(t, false, c.options.DisableShallowCopy)
	assert.Equal(t, 0, len(c.options.ClientOptions))
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

	// Set From Client
	deleget, err := c.newDelegate(context.TODO(), &CopyObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		SourceBucket: Ptr("src-bucket"),
		SourceKey:    Ptr("src-key"),
	}, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = true
		co.ClientOptions = []func(do *Options){func(do *Options) {}}
		co.MetadataProperties = &HeadObjectResult{}
		co.TagProperties = &GetObjectTaggingResult{}
	})
	assert.NoError(t, err)

	assert.Equal(t, int(2), deleget.options.ParallelNum)
	assert.Equal(t, int64(1024*1024), deleget.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), deleget.options.MultipartCopyThreshold)
	assert.Equal(t, true, deleget.options.LeavePartsOnError)
	assert.Equal(t, true, deleget.options.DisableShallowCopy)
	assert.Equal(t, 1, len(deleget.options.ClientOptions))
	// only supports setting from c.Copy
	assert.NotNil(t, deleget.options.MetadataProperties)
	assert.NotNil(t, deleget.options.TagProperties)

	// Use WithXXX
	deleget, err = c.newDelegate(context.TODO(), &CopyObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		SourceBucket: Ptr("src-bucket"),
		SourceKey:    Ptr("src-key"),
	}, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = false
		co.ClientOptions = []func(do *Options){func(do *Options) {}, func(do *Options) {}}
		co.MetadataProperties = nil
		co.TagProperties = &GetObjectTaggingResult{}
	},
		WithCopierParallelNum(5),
		WithCopierPartSize(2*1024*1024),
	)

	assert.NoError(t, err)
	assert.Equal(t, int(5), deleget.options.ParallelNum)
	assert.Equal(t, int64(2*1024*1024), deleget.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), deleget.options.MultipartCopyThreshold)
	assert.Equal(t, true, deleget.options.LeavePartsOnError)
	assert.Equal(t, false, deleget.options.DisableShallowCopy)
	assert.Equal(t, 2, len(deleget.options.ClientOptions))
	assert.Nil(t, deleget.options.MetadataProperties)
	assert.NotNil(t, deleget.options.TagProperties)

}

func TestCopierShallowCopyFlags(t *testing.T) {

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou")

	client := NewClient(cfg)

	// Default
	c := NewCopier(client)
	assert.Equal(t, DefaultCopyParallel, c.options.ParallelNum)
	assert.Equal(t, DefaultCopyPartSize, c.options.PartSize)
	assert.Equal(t, DefaultCopyThreshold, c.options.MultipartCopyThreshold)
	assert.Equal(t, false, c.options.LeavePartsOnError)
	assert.Equal(t, false, c.options.DisableShallowCopy)
	assert.Equal(t, false, c.options.NoCheckSSE)
	assert.Equal(t, false, c.options.NoCheckCrossBucket)
	assert.Equal(t, 0, len(c.options.ClientOptions))
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

	// Set From Client
	deleget, err := c.newDelegate(context.TODO(), &CopyObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		SourceBucket: Ptr("src-bucket"),
		SourceKey:    Ptr("src-key"),
	}, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = true
		co.ClientOptions = []func(do *Options){func(do *Options) {}}
		co.MetadataProperties = &HeadObjectResult{}
		co.TagProperties = &GetObjectTaggingResult{}
		co.NoCheckSSE = true
		co.NoCheckCrossBucket = true
	})
	assert.NoError(t, err)

	assert.Equal(t, int(2), deleget.options.ParallelNum)
	assert.Equal(t, int64(1024*1024), deleget.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), deleget.options.MultipartCopyThreshold)
	assert.Equal(t, true, deleget.options.LeavePartsOnError)
	assert.Equal(t, true, deleget.options.DisableShallowCopy)
	assert.Equal(t, 1, len(deleget.options.ClientOptions))
	// only supports setting from c.Copy
	assert.NotNil(t, deleget.options.MetadataProperties)
	assert.NotNil(t, deleget.options.TagProperties)

	assert.Equal(t, true, deleget.options.NoCheckSSE)
	assert.Equal(t, true, deleget.options.NoCheckCrossBucket)

	// Use WithXXX
	deleget, err = c.newDelegate(context.TODO(), &CopyObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		SourceBucket: Ptr("src-bucket"),
		SourceKey:    Ptr("src-key"),
	}, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = false
		co.ClientOptions = []func(do *Options){func(do *Options) {}, func(do *Options) {}}
		co.MetadataProperties = nil
		co.TagProperties = &GetObjectTaggingResult{}
	},
		WithCopierParallelNum(5),
		WithCopierPartSize(2*1024*1024),
		WithCopierNoCheckCrossBucket(true),
		WithCopierNoCheckSSE(true),
	)

	assert.NoError(t, err)
	assert.Equal(t, int(5), deleget.options.ParallelNum)
	assert.Equal(t, int64(2*1024*1024), deleget.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), deleget.options.MultipartCopyThreshold)
	assert.Equal(t, true, deleget.options.LeavePartsOnError)
	assert.Equal(t, false, deleget.options.DisableShallowCopy)
	assert.Equal(t, 2, len(deleget.options.ClientOptions))
	assert.Nil(t, deleget.options.MetadataProperties)
	assert.NotNil(t, deleget.options.TagProperties)

	assert.Equal(t, true, deleget.options.NoCheckSSE)
	assert.Equal(t, true, deleget.options.NoCheckCrossBucket)
}

func TestCopierShallowCopyFallbackOnEntityTooLarge(t *testing.T) {
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		method := r.Method
		path := r.URL.Path
		query := r.URL.RawQuery

		switch {
		// Head source object
		case method == "HEAD" && path == "/bucket/src-key":
			w.Header().Set("Content-Length", fmt.Sprintf("%d", 2*1024*1024))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("ETag", "\"src-etag\"")
			w.WriteHeader(200)

		// CopyObject returns EntityTooLarge
		case method == "PUT" && path == "/bucket/dst-key" && query == "":
			w.Header().Set("Content-Type", "application/xml")
			w.Header().Set("x-oss-request-id", "req-copy-fail")
			w.WriteHeader(400)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>EntityTooLarge</Code>
  <Message>Entity Too Large</Message>
  <RequestId>req-copy-fail</RequestId>
</Error>`))

		// InitiateMultipartUpload
		case method == "POST" && path == "/bucket/dst-key" && strings.Contains(query, "uploads"):
			w.Header().Set("Content-Type", "application/xml")
			w.Header().Set("x-oss-request-id", "req-init")
			w.WriteHeader(200)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
  <Bucket>bucket</Bucket>
  <Key>dst-key</Key>
  <UploadId>upload-id-123</UploadId>
</InitiateMultipartUploadResult>`))

		// UploadPartCopy
		case method == "PUT" && path == "/bucket/dst-key" && strings.Contains(query, "partNumber="):
			w.Header().Set("Content-Type", "application/xml")
			w.Header().Set("x-oss-request-id", "req-part-copy")
			w.WriteHeader(200)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
  <LastModified>2024-01-01T00:00:00.000Z</LastModified>
  <ETag>"part-etag"</ETag>
</CopyPartResult>`))

		// CompleteMultipartUpload
		case method == "POST" && path == "/bucket/dst-key":
			w.Header().Set("Content-Type", "application/xml")
			w.Header().Set("x-oss-request-id", "req-complete")
			w.WriteHeader(200)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult>
  <Location>http://bucket.oss-cn-hangzhou.aliyuncs.com/dst-key</Location>
  <Bucket>bucket</Bucket>
  <Key>dst-key</Key>
  <ETag>"final-etag"</ETag>
</CompleteMultipartUploadResult>`))
		default:
			t.Fatalf("unexpected request #%d: %s %s?%s", count, method, path, query)
		}
		t.Logf("request #%d: %s %s?%s", count, method, path, query)
	}))
	defer server.Close()

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL)

	client := NewClient(cfg)
	copier := NewCopier(client, func(co *CopierOptions) {
		co.MultipartCopyThreshold = 1024 * 1024
		co.PartSize = 1024 * 1024
		co.ParallelNum = 1
	})

	result, err := copier.Copy(context.TODO(), &CopyObjectRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("dst-key"),
		SourceKey: Ptr("src-key"),
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "upload-id-123", *result.UploadId)
	assert.Equal(t, "\"final-etag\"", *result.ETag)
	assert.Equal(t, int32(6), atomic.LoadInt32(&requestCount))
}

func TestCopierShallowCopyFallbackMultiCopyFails(t *testing.T) {
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		method := r.Method
		path := r.URL.Path
		query := r.URL.RawQuery

		switch {
		// Head source object
		case method == "HEAD" && path == "/bucket/src-key":
			w.Header().Set("Content-Length", fmt.Sprintf("%d", 2*1024*1024))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("ETag", "\"src-etag\"")
			w.WriteHeader(200)

		// CopyObject returns EntityTooLarge -> triggers fallback
		case method == "PUT" && path == "/bucket/dst-key" && query == "":
			w.Header().Set("Content-Type", "application/xml")
			w.Header().Set("x-oss-request-id", "req-copy-fail")
			w.WriteHeader(400)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>EntityTooLarge</Code>
  <Message>Entity Too Large</Message>
  <RequestId>req-copy-fail</RequestId>
</Error>`))

		// InitiateMultipartUpload fails (non-retriable), so multiCopy fails
		case method == "POST" && path == "/bucket/dst-key" && strings.Contains(query, "uploads"):
			w.Header().Set("Content-Type", "application/xml")
			w.Header().Set("x-oss-request-id", "req-init-fail")
			w.WriteHeader(403)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>Access Denied</Message>
  <RequestId>req-init-fail</RequestId>
</Error>`))
		default:
			t.Fatalf("unexpected request #%d: %s %s?%s", count, method, path, query)
		}
		t.Logf("request #%d: %s %s?%s", count, method, path, query)
	}))
	defer server.Close()

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithRetryMaxAttempts(1)

	client := NewClient(cfg)
	copier := NewCopier(client, func(co *CopierOptions) {
		co.MultipartCopyThreshold = 1024 * 1024
		co.PartSize = 1024 * 1024
		co.ParallelNum = 1
	})

	result, err := copier.Copy(context.TODO(), &CopyObjectRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("dst-key"),
		SourceKey: Ptr("src-key"),
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "AccessDenied")
	// HEAD + CopyObject(EntityTooLarge) + InitiateMultipartUpload(AccessDenied)
	assert.Equal(t, int32(3), atomic.LoadInt32(&requestCount))
}

func TestCopierShallowCopyNoFallbackOnOtherError(t *testing.T) {
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		method := r.Method
		path := r.URL.Path
		query := r.URL.RawQuery

		switch {
		// Head source object
		case method == "HEAD" && path == "/bucket/src-key":
			w.Header().Set("Content-Length", fmt.Sprintf("%d", 2*1024*1024))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("ETag", "\"src-etag\"")
			w.WriteHeader(200)

		// CopyObject returns a non-EntityTooLarge error -> must NOT fall back
		case method == "PUT" && path == "/bucket/dst-key" && query == "":
			w.Header().Set("Content-Type", "application/xml")
			w.Header().Set("x-oss-request-id", "req-copy-denied")
			w.WriteHeader(403)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>Access Denied</Message>
  <RequestId>req-copy-denied</RequestId>
</Error>`))
		default:
			t.Fatalf("unexpected request #%d: %s %s?%s (fallback should not happen)", count, method, path, query)
		}
		t.Logf("request #%d: %s %s?%s", count, method, path, query)
	}))
	defer server.Close()

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithRetryMaxAttempts(1)

	client := NewClient(cfg)
	copier := NewCopier(client, func(co *CopierOptions) {
		co.MultipartCopyThreshold = 1024 * 1024
		co.PartSize = 1024 * 1024
		co.ParallelNum = 1
	})

	result, err := copier.Copy(context.TODO(), &CopyObjectRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("dst-key"),
		SourceKey: Ptr("src-key"),
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "AccessDenied")
	// Only HEAD + CopyObject; no InitiateMultipartUpload means no fallback happened
	assert.Equal(t, int32(2), atomic.LoadInt32(&requestCount))
}
