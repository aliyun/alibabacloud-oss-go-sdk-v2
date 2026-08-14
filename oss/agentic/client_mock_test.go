package agentic

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func assertURL(t *testing.T, rawURL string, expectBase string, expectParams []string) {
	t.Helper()
	u, err := url.Parse(rawURL)
	assert.Nil(t, err)
	assert.Equal(t, expectBase, u.Scheme+"://"+u.Host+u.Path)
	for _, p := range expectParams {
		assert.True(t, u.Query().Has(p), "missing query param: %s in %s", p, rawURL)
	}
	assert.Equal(t, len(expectParams), len(u.Query()), "unexpected query params in %s", rawURL)
}

type urlCaptureTransport struct {
	RequestURL    string
	RequestMethod string
	RequestHost   string
}

func (t *urlCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.RequestURL = req.URL.String()
	t.RequestMethod = req.Method
	t.RequestHost = req.Host
	return &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Header: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"mock-request-id"},
		},
		Body: io.NopCloser(strings.NewReader("")),
	}, nil
}

func newMockAgenticBucketClient(region, accountId string) (*AgenticBucketClient, *urlCaptureTransport) {
	transport := &urlCaptureTransport{}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion(region).
		WithAccountId(accountId).
		WithHttpClient(&http.Client{Transport: transport})

	client := NewAgenticBucketClient(cfg)
	return client, transport
}

func newMockBucketSpaceClient(region, accountId, endpoint string) (*oss.Client, *urlCaptureTransport) {
	transport := &urlCaptureTransport{}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion(region).
		WithAccountId(accountId).
		WithEndpoint(endpoint).
		WithHttpClient(&http.Client{Transport: transport})

	client := NewBucketSpaceClient(cfg)
	return client, transport
}

func TestMockAgenticBucketClient_CreateAgenticBucket(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assert.Equal(t, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/?agenticBucket", transport.RequestURL)
}

func TestMockAgenticBucketClient_DeleteAgenticBucket(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.DeleteAgenticBucket(context.TODO(), &DeleteAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "DELETE", transport.RequestMethod)
	assert.Equal(t, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/?agenticBucket", transport.RequestURL)
}

func TestMockAgenticBucketClient_GetAgenticBucket(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assert.Equal(t, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/?agenticBucket", transport.RequestURL)
}

func TestMockAgenticBucketClient_HostLabelTooLong(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	// fullName = "{bucket}-123456-cn-hangzhou-ab-apsr" -> len(bucket)+27; a 37-char bucket
	// yields a 64-char label, exceeding the 63-character limit in virtual-hosted style.
	longBucket := strings.Repeat("a", 37)
	_, err := client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr(longBucket),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "exceeds the maximum length of 63 characters")
	// The request must not have been sent.
	assert.Equal(t, "", transport.RequestMethod)
	assert.Equal(t, "", transport.RequestURL)
}

func TestMockAgenticBucketClient_ListAgenticBuckets(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.ListAgenticBuckets(context.TODO(), &ListAgenticBucketsRequest{})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/?agenticBucket", transport.RequestURL)
}

func TestMockAgenticBucketClient_PutAgenticBucketAcl(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.PutAgenticBucketAcl(context.TODO(), &PutAgenticBucketAclRequest{
		Bucket: oss.Ptr("my-agentic"),
		Acl:    oss.BucketACLPrivate,
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "acl"})
}

func TestMockAgenticBucketClient_GetAgenticBucketAcl(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.GetAgenticBucketAcl(context.TODO(), &GetAgenticBucketAclRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "acl"})
}

func TestMockAgenticBucketClient_PutAgenticBucketEncryption(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.PutAgenticBucketEncryption(context.TODO(), &PutAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr("my-agentic"),
		ServerSideEncryptionRule: &ServerSideEncryptionRule{
			ApplyServerSideEncryptionByDefault: &ApplyServerSideEncryptionByDefault{
				SSEAlgorithm: oss.Ptr("AES256"),
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "encryption"})
}

func TestMockAgenticBucketClient_PutAgenticBucketVersioning(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.PutAgenticBucketVersioning(context.TODO(), &PutAgenticBucketVersioningRequest{
		Bucket: oss.Ptr("my-agentic"),
		VersioningConfiguration: &VersioningConfiguration{
			Status: oss.VersionEnabled,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "versioning"})
}

func TestMockAgenticBucketClient_PutAgenticBucketPolicy(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	policy := `{"Version":"1","Statement":[]}`
	_, err := client.PutAgenticBucketPolicy(context.TODO(), &PutAgenticBucketPolicyRequest{
		Bucket: oss.Ptr("my-agentic"),
		Body:   strings.NewReader(policy),
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "policy"})
}

func TestMockAgenticBucketClient_PutAgenticBucketPublicAccessBlock(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.PutAgenticBucketPublicAccessBlock(context.TODO(), &PutAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr("my-agentic"),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			BlockPublicAccess: oss.Ptr(true),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "publicAccessBlock"})
}

func TestMockAgenticBucketClient_GetAgenticBucketEncryption(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.GetAgenticBucketEncryption(context.TODO(), &GetAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "encryption"})
}

func TestMockAgenticBucketClient_DeleteAgenticBucketEncryption(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.DeleteAgenticBucketEncryption(context.TODO(), &DeleteAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "DELETE", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "encryption"})
}

func TestMockAgenticBucketClient_GetAgenticBucketVersioning(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.GetAgenticBucketVersioning(context.TODO(), &GetAgenticBucketVersioningRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "versioning"})
}

func TestMockAgenticBucketClient_GetAgenticBucketPolicy(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.GetAgenticBucketPolicy(context.TODO(), &GetAgenticBucketPolicyRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "policy"})
}

func TestMockAgenticBucketClient_DeleteAgenticBucketPolicy(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.DeleteAgenticBucketPolicy(context.TODO(), &DeleteAgenticBucketPolicyRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "DELETE", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "policy"})
}

func TestMockAgenticBucketClient_GetAgenticBucketPublicAccessBlock(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.GetAgenticBucketPublicAccessBlock(context.TODO(), &GetAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "publicAccessBlock"})
}

func TestMockAgenticBucketClient_DeleteAgenticBucketPublicAccessBlock(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.DeleteAgenticBucketPublicAccessBlock(context.TODO(), &DeleteAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "DELETE", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "publicAccessBlock"})
}

func TestMockAgenticBucketClient_PutAgenticBucketStatus(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.PutAgenticBucketStatus(context.TODO(), &PutAgenticBucketStatusRequest{
		Bucket: oss.Ptr("my-agentic"),
		AgenticBucketStatus: &AgenticBucketStatus{
			Status: oss.Ptr("enabled"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "status"})
}

func TestMockAgenticBucketClient_ListBucketSpaces(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-hangzhou", "123456")

	_, err := client.ListBucketSpaces(context.TODO(), &ListBucketSpacesRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-agentic-123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket", "bucketSpace"})
}

func newMockAgenticBucketClientPathStyle(region, accountId string) (*AgenticBucketClient, *urlCaptureTransport) {
	transport := &urlCaptureTransport{}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion(region).
		WithAccountId(accountId).
		WithUsePathStyle(true).
		WithHttpClient(&http.Client{Transport: transport})

	client := NewAgenticBucketClient(cfg)
	return client, transport
}

func TestMockAgenticBucketClient_GetAgenticBucket_PathStyle(t *testing.T) {
	client, transport := newMockAgenticBucketClientPathStyle("cn-hangzhou", "123456")

	_, err := client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://oss-cn-hangzhou.aliyuncs.com/my-agentic-123456-cn-hangzhou-ab-apsr/", []string{"agenticBucket"})
}

func TestMockAgenticBucketClient_ListAgenticBuckets_PathStyle(t *testing.T) {
	client, transport := newMockAgenticBucketClientPathStyle("cn-hangzhou", "123456")

	_, err := client.ListAgenticBuckets(context.TODO(), &ListAgenticBucketsRequest{})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://oss-cn-hangzhou.aliyuncs.com/", []string{"agenticBucket"})
}

func TestMockBucketSpaceClient_PutObject_PathStyle(t *testing.T) {
	transport := &urlCaptureTransport{}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithAccountId("123456").
		WithEndpoint("user-cname.test.com").
		WithUsePathStyle(true).
		WithHttpClient(&http.Client{Transport: transport})

	client := NewBucketSpaceClient(cfg)

	_, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr("my-space"),
		Key:    oss.Ptr("test.txt"),
		Body:   strings.NewReader("hello"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assert.Equal(t, "https://user-cname.test.com/my-space-123456-cn-hangzhou-bs-apsr/test.txt", transport.RequestURL)
}

func TestMockAgenticBucketClient_RegionInURL(t *testing.T) {
	client, transport := newMockAgenticBucketClient("cn-shanghai", "999888")

	_, err := client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr("test-bucket"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "https://test-bucket-999888-cn-shanghai-ab-apsr.oss-cn-shanghai.aliyuncs.com/?agenticBucket", transport.RequestURL)
}

func TestMockBucketSpaceClient_PutObject(t *testing.T) {
	client, transport := newMockBucketSpaceClient("cn-hangzhou", "123456", "user-cname.test.com")

	_, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr("my-space"),
		Key:    oss.Ptr("test.txt"),
		Body:   strings.NewReader("hello"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assert.Equal(t, "https://my-space-123456-cn-hangzhou-bs-apsr.user-cname.test.com/test.txt", transport.RequestURL)
}

func TestMockBucketSpaceClient_PutObject_Internal(t *testing.T) {
	transport := &urlCaptureTransport{}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithAccountId("123456").
		WithUseInternalEndpoint(true).
		WithHttpClient(&http.Client{Transport: transport})

	client := NewBucketSpaceClient(cfg)

	_, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr("my-space"),
		Key:    oss.Ptr("test.txt"),
		Body:   strings.NewReader("hello"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assert.Equal(t, "https://my-space-123456-cn-hangzhou-bs-apsr.oss-cn-hangzhou-internal.aliyuncs.com/test.txt", transport.RequestURL)
}

func TestMockBucketSpaceClient_GetBucketInfo(t *testing.T) {
	client, transport := newMockBucketSpaceClient("cn-hangzhou", "123456", "user-cname.test.com")

	_, err := client.GetBucketInfo(context.TODO(), &oss.GetBucketInfoRequest{
		Bucket: oss.Ptr("my-space"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assert.Equal(t, "https://my-space-123456-cn-hangzhou-bs-apsr.user-cname.test.com/?bucketInfo", transport.RequestURL)
}

func newMockClientWithHelper(region, accountId, endpoint string) (*oss.Client, *BucketSpaceHelper, *urlCaptureTransport) {
	transport := &urlCaptureTransport{}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion(region).
		WithAccountId(accountId).
		WithEndpoint(endpoint).
		WithHttpClient(&http.Client{Transport: transport})

	client := oss.NewClient(cfg)
	helper := NewBucketSpaceHelper(cfg)
	return client, helper, transport
}

func TestMockClientWithHelper_PutObject(t *testing.T) {
	client, helper, transport := newMockClientWithHelper("cn-hangzhou", "123456", "user-cname.test.com")

	bucket := helper.ToBucketName("my-space")
	_, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucket),
		Key:    oss.Ptr("dir/test.txt"),
		Body:   strings.NewReader("hello"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assert.Equal(t, "https://my-space-123456-cn-hangzhou-bs-apsr.user-cname.test.com/dir/test.txt", transport.RequestURL)
}

func TestMockClientWithHelper_GetObject(t *testing.T) {
	client, helper, transport := newMockClientWithHelper("cn-hangzhou", "123456", "user-cname.test.com")

	bucket := helper.ToBucketName("my-space")
	_, err := client.GetObject(context.TODO(), &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucket),
		Key:    oss.Ptr("test.txt"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assert.Equal(t, "https://my-space-123456-cn-hangzhou-bs-apsr.user-cname.test.com/test.txt", transport.RequestURL)
}

func TestMockClientWithHelper_DeleteObject(t *testing.T) {
	client, helper, transport := newMockClientWithHelper("cn-hangzhou", "123456", "user-cname.test.com")

	bucket := helper.ToBucketName("my-space")
	_, err := client.DeleteObject(context.TODO(), &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(bucket),
		Key:    oss.Ptr("test.txt"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "DELETE", transport.RequestMethod)
	assert.Equal(t, "https://my-space-123456-cn-hangzhou-bs-apsr.user-cname.test.com/test.txt", transport.RequestURL)
}

func TestMockClientWithHelper_ListObjects(t *testing.T) {
	client, helper, transport := newMockClientWithHelper("cn-hangzhou", "123456", "user-cname.test.com")

	bucket := helper.ToBucketName("my-space")
	_, err := client.ListObjectsV2(context.TODO(), &oss.ListObjectsV2Request{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assertURL(t, transport.RequestURL, "https://my-space-123456-cn-hangzhou-bs-apsr.user-cname.test.com/", []string{"list-type", "encoding-type"})
}

// withAliasStyle enables the short-alias host form by setting the addressing
// style before the agentic provider optFn reads it.
func withAliasStyle(o *oss.Options) {
	o.UrlStyle = oss.UrlStyleVirtualHostedAlias
}

func newMockAgenticBucketClientAliasStyle(region, accountId, endpoint string) (*AgenticBucketClient, *urlCaptureTransport) {
	transport := &urlCaptureTransport{}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion(region).
		WithAccountId(accountId).
		WithEndpoint(endpoint).
		WithHttpClient(&http.Client{Transport: transport})

	client := NewAgenticBucketClient(cfg, withAliasStyle)
	return client, transport
}

func TestMockAgenticBucketClient_AliasStyle(t *testing.T) {
	client, transport := newMockAgenticBucketClientAliasStyle("cn-hangzhou", "123456", "abc.com")

	_, err := client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	// Host uses the short literal "alias" label; accountId/region stay off the host.
	assert.Equal(t, "https://my-agentic-alias-ab-apsr.abc.com/?agenticBucket", transport.RequestURL)

	transport = &urlCaptureTransport{}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithAccountId("123456").
		WithEndpoint("abc.com").
		WithHttpClient(&http.Client{Transport: transport}).WithUseVirtualHostedAlias(true)

	client = NewAgenticBucketClient(cfg)
	_, err = client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	// Host uses the short literal "alias" label; accountId/region stay off the host.
	assert.Equal(t, "https://my-agentic-alias-ab-apsr.abc.com/?agenticBucket", transport.RequestURL)
}

func TestMockAgenticBucketClient_AliasStyle_MissingAccountId(t *testing.T) {
	// accountId empty: the signing name (full name) cannot be resolved, so the
	// request fails even though the alias host label needs no accountId.
	client, transport := newMockAgenticBucketClientAliasStyle("cn-hangzhou", "", "abc.com")

	_, err := client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "AccountId")
	// The request must not have been sent.
	assert.Equal(t, "", transport.RequestMethod)
	assert.Equal(t, "", transport.RequestURL)
}

func TestMockAgenticBucketClient_AliasStyle_HostLabelWithinLongName(t *testing.T) {
	// The full signing name far exceeds 63 chars, but the short alias label fits,
	// which is exactly the reason the alias form exists.
	client, transport := newMockAgenticBucketClientAliasStyle("cn-hangzhou", "123456", "abc.com")

	// aliasLabel = "{bucket}-alias-ab-apsr" -> len(bucket)+14; a 49-char bucket
	// yields a 63-char label (the boundary), while the full name would be far longer.
	bucket := strings.Repeat("a", 49)
	_, err := client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.Equal(t, "GET", transport.RequestMethod)
	assert.Equal(t, "https://"+bucket+"-alias-ab-apsr.abc.com/?agenticBucket", transport.RequestURL)
}

func TestMockBucketSpaceClient_AliasStyle(t *testing.T) {
	transport := &urlCaptureTransport{}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithAccountId("123456").
		WithEndpoint("abc.com").
		WithHttpClient(&http.Client{Transport: transport})

	client := NewBucketSpaceClient(cfg, withAliasStyle)

	_, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr("my-space"),
		Key:    oss.Ptr("test.txt"),
		Body:   strings.NewReader("hello"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assert.Equal(t, "https://my-space-alias-bs-apsr.abc.com/test.txt", transport.RequestURL)

	cfg = oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithAccountId("123456").
		WithEndpoint("abc.com").
		WithHttpClient(&http.Client{Transport: transport}).WithUseVirtualHostedAlias(true)

	client = NewBucketSpaceClient(cfg)

	_, err = client.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr("my-space"),
		Key:    oss.Ptr("test.txt"),
		Body:   strings.NewReader("hello"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "PUT", transport.RequestMethod)
	assert.Equal(t, "https://my-space-alias-bs-apsr.abc.com/test.txt", transport.RequestURL)
}

func TestMockBucketSpaceClient_Presign(t *testing.T) {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithAccountId("123456").
		WithEndpoint("user-cname.test.com").
		WithSignatureVersion(oss.SignatureVersionV4)

	client := NewBucketSpaceClient(cfg)

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), &oss.GetObjectRequest{
		Bucket: oss.Ptr("my-space"),
		Key:    oss.Ptr("test.txt"),
	}, oss.PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	u, perr := url.Parse(result.URL)
	assert.Nil(t, perr)
	assert.Equal(t, "https://my-space-123456-cn-hangzhou-bs-apsr.user-cname.test.com/test.txt", u.Scheme+"://"+u.Host+u.Path)
	assert.Contains(t, result.URL, "x-oss-signature=")
	assert.Contains(t, result.URL, "x-oss-credential=")
}

func TestMockBucketSpaceClient_Presign_AliasStyle(t *testing.T) {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithAccountId("123456").
		WithEndpoint("abc.com").
		WithSignatureVersion(oss.SignatureVersionV4)

	client := NewBucketSpaceClient(cfg, withAliasStyle)

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), &oss.GetObjectRequest{
		Bucket: oss.Ptr("my-space"),
		Key:    oss.Ptr("test.txt"),
	}, oss.PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	u, perr := url.Parse(result.URL)
	assert.Nil(t, perr)
	assert.Equal(t, "https://my-space-alias-bs-apsr.abc.com/test.txt", u.Scheme+"://"+u.Host+u.Path)
	assert.Contains(t, result.URL, "x-oss-signature=")
	assert.Contains(t, result.URL, "x-oss-credential=")
}

func TestMockBucketSpaceClient_Presign_MissingAccountId(t *testing.T) {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("abc.com").
		WithSignatureVersion(oss.SignatureVersionV4)

	client := NewBucketSpaceClient(cfg, withAliasStyle)

	_, err := client.Presign(context.TODO(), &oss.GetObjectRequest{
		Bucket: oss.Ptr("my-space"),
		Key:    oss.Ptr("test.txt"),
	}, oss.PresignExpiration(time.Now().Add(1*time.Hour)))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "AccountId")
}
