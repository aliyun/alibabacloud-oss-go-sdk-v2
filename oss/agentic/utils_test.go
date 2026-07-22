package agentic

import (
	"net/url"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestAgenticBucketNameResolver(t *testing.T) {
	r := &agenticProvider{
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "ab-apsr",
	}

	input := &oss.OperationInput{
		Bucket: oss.Ptr("my-agentic"),
	}
	name, err := r.BuildBucketName(input)
	assert.Nil(t, err)
	assert.Equal(t, "my-agentic-1234567890123456-cn-hangzhou-ab-apsr", name)

	input = &oss.OperationInput{}
	name, err = r.BuildBucketName(input)
	assert.Nil(t, err)
	assert.Equal(t, "", name)
}

func TestBucketSpaceNameResolver(t *testing.T) {
	r := &agenticProvider{
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "bs-apsr",
	}

	input := &oss.OperationInput{
		Bucket: oss.Ptr("my-sandbox"),
	}
	name, err := r.BuildBucketName(input)
	assert.Nil(t, err)
	assert.Equal(t, "my-sandbox-1234567890123456-cn-hangzhou-bs-apsr", name)
}

func TestAgenticEndpointProvider(t *testing.T) {
	endpoint, _ := url.Parse("https://oss-cn-hangzhou.aliyuncs.com")
	p := &agenticProvider{
		endpoint:  endpoint,
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "ab-apsr",
	}

	// With bucket (prefix expanded to full name, prepended to endpoint host)
	input := &oss.OperationInput{
		Bucket: oss.Ptr("my-agentic"),
	}
	got := p.BuildURL(input)
	assert.Equal(t, "https://my-agentic-1234567890123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", got)

	// Without bucket (endpoint host as-is)
	input = &oss.OperationInput{}
	got = p.BuildURL(input)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/", got)

	// Nil input
	got = p.BuildURL(nil)
	assert.Equal(t, "", got)

	// BucketSpace suffix with internal endpoint
	endpoint2, _ := url.Parse("https://oss-cn-hangzhou-internal.aliyuncs.com")
	p2 := &agenticProvider{
		endpoint:  endpoint2,
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "bs-apsr",
	}

	input = &oss.OperationInput{
		Bucket: oss.Ptr("my-sandbox"),
		Key:    oss.Ptr("test.txt"),
	}
	url2 := p2.BuildURL(input)
	assert.Equal(t, "https://my-sandbox-1234567890123456-cn-hangzhou-bs-apsr.oss-cn-hangzhou-internal.aliyuncs.com/test.txt", url2)
}

func TestAgenticEndpointProviderPathStyle(t *testing.T) {
	endpoint, _ := url.Parse("https://oss-cn-hangzhou.aliyuncs.com")
	p := &agenticProvider{
		endpoint:  endpoint,
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "ab-apsr",
		urlStyle:  oss.UrlStylePath,
	}

	// With bucket (full name in path, endpoint host as-is)
	input := &oss.OperationInput{
		Bucket: oss.Ptr("my-agentic"),
	}
	got := p.BuildURL(input)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/my-agentic-1234567890123456-cn-hangzhou-ab-apsr/", got)

	// With bucket and key
	input = &oss.OperationInput{
		Bucket: oss.Ptr("my-agentic"),
		Key:    oss.Ptr("test.txt"),
	}
	got = p.BuildURL(input)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/my-agentic-1234567890123456-cn-hangzhou-ab-apsr/test.txt", got)

	// Without bucket (endpoint host as-is)
	input = &oss.OperationInput{}
	got = p.BuildURL(input)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/", got)
}

func TestBucketSpaceHelper(t *testing.T) {
	cfg := &oss.Config{}
	cfg.WithAccountId("1234567890123456")
	cfg.WithRegion("cn-hangzhou")

	helper := NewBucketSpaceHelper(cfg)
	name := helper.ToBucketName("my-sandbox")
	assert.Equal(t, "my-sandbox-1234567890123456-cn-hangzhou-bs-apsr", name)
}
