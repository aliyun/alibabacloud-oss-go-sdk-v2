package agentic

import (
	"net/url"
	"strings"
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
	got, err := p.BuildURL(input)
	assert.Nil(t, err)
	assert.Equal(t, "https://my-agentic-1234567890123456-cn-hangzhou-ab-apsr.oss-cn-hangzhou.aliyuncs.com/", got)

	// Without bucket (endpoint host as-is)
	input = &oss.OperationInput{}
	got, err = p.BuildURL(input)
	assert.Nil(t, err)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/", got)

	// Nil input
	got, err = p.BuildURL(nil)
	assert.Nil(t, err)
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
	url2, err := p2.BuildURL(input)
	assert.Nil(t, err)
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
	got, err := p.BuildURL(input)
	assert.Nil(t, err)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/my-agentic-1234567890123456-cn-hangzhou-ab-apsr/", got)

	// With bucket and key
	input = &oss.OperationInput{
		Bucket: oss.Ptr("my-agentic"),
		Key:    oss.Ptr("test.txt"),
	}
	got, err = p.BuildURL(input)
	assert.Nil(t, err)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/my-agentic-1234567890123456-cn-hangzhou-ab-apsr/test.txt", got)

	// Without bucket (endpoint host as-is)
	input = &oss.OperationInput{}
	got, err = p.BuildURL(input)
	assert.Nil(t, err)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/", got)
}

func TestAgenticProviderImplementsEndpointProviderE(t *testing.T) {
	var _ oss.EndpointProviderE = (*agenticProvider)(nil)
	var _ oss.BucketNameResolver = (*agenticProvider)(nil)
}

func TestAgenticProviderMissingRequiredFields(t *testing.T) {
	endpoint, _ := url.Parse("https://oss-cn-hangzhou.aliyuncs.com")
	input := &oss.OperationInput{Bucket: oss.Ptr("my-agentic")}

	// Missing accountId
	p := &agenticProvider{endpoint: endpoint, region: "cn-hangzhou", suffix: "ab-apsr"}
	_, err := p.BuildURL(input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "AccountId")
	_, err = p.BuildBucketName(input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "AccountId")

	// Missing region
	p = &agenticProvider{endpoint: endpoint, accountId: "1234567890123456", suffix: "ab-apsr"}
	_, err = p.BuildURL(input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Region")
	_, err = p.BuildBucketName(input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Region")

	// No bucket: validation is skipped, no error
	name, err := p.BuildBucketName(&oss.OperationInput{})
	assert.Nil(t, err)
	assert.Equal(t, "", name)
}

func TestAgenticEndpointProviderHostLabelTooLong(t *testing.T) {
	endpoint, _ := url.Parse("https://oss-cn-hangzhou.aliyuncs.com")
	// fullName = "{bucket}-1234567890123456-cn-hangzhou-ab-apsr" -> len(bucket) + 37
	suffixPart := "-1234567890123456-cn-hangzhou-ab-apsr"
	p := &agenticProvider{
		endpoint:  endpoint,
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "ab-apsr",
	}

	// Boundary: fullName == 63 (bucket 26) is allowed in virtual-hosted style
	okName := strings.Repeat("a", 26)
	assert.Equal(t, 63, len(okName+suffixPart))
	got, err := p.BuildURL(&oss.OperationInput{Bucket: oss.Ptr(okName)})
	assert.Nil(t, err)
	assert.Equal(t, "https://"+okName+suffixPart+".oss-cn-hangzhou.aliyuncs.com/", got)

	// Over limit: fullName == 64 (bucket 27) is rejected in virtual-hosted style
	longName := strings.Repeat("a", 27)
	assert.Equal(t, 64, len(longName+suffixPart))
	_, err = p.BuildURL(&oss.OperationInput{Bucket: oss.Ptr(longName)})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "exceeds the maximum length of 63 characters")

	// Path style has no DNS label limit, so the same long name is fine
	pathStyle := &agenticProvider{
		endpoint:  endpoint,
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "ab-apsr",
		urlStyle:  oss.UrlStylePath,
	}
	got, err = pathStyle.BuildURL(&oss.OperationInput{Bucket: oss.Ptr(longName)})
	assert.Nil(t, err)
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com/"+longName+suffixPart+"/", got)
}

func TestUrlStyleVirtualHostedAliasString(t *testing.T) {
	assert.Equal(t, "virtual-hosted-alias-style", oss.UrlStyleVirtualHostedAlias.String())
}

func TestAgenticAliasStyleSigningNameIsFull(t *testing.T) {
	// In alias style the signing name (BuildBucketName) is still the full name,
	// so accountId and region remain required.
	p := &agenticProvider{
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "ab-apsr",
		urlStyle:  oss.UrlStyleVirtualHostedAlias,
	}
	name, err := p.BuildBucketName(&oss.OperationInput{Bucket: oss.Ptr("my-agentic")})
	assert.Nil(t, err)
	assert.Equal(t, "my-agentic-1234567890123456-cn-hangzhou-ab-apsr", name)

	// Missing accountId still fails, even in alias style.
	bad := &agenticProvider{region: "cn-hangzhou", suffix: "ab-apsr", urlStyle: oss.UrlStyleVirtualHostedAlias}
	_, err = bad.BuildBucketName(&oss.OperationInput{Bucket: oss.Ptr("my-agentic")})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "AccountId")
}

func TestAgenticEndpointProviderAliasStyle(t *testing.T) {
	endpoint, _ := url.Parse("https://abc.com")
	p := &agenticProvider{
		endpoint:  endpoint,
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "bs-apsr",
		urlStyle:  oss.UrlStyleVirtualHostedAlias,
	}

	// Host uses the short alias label; the full name (used for signing) is not on the host.
	input := &oss.OperationInput{Bucket: oss.Ptr("my-sandbox"), Key: oss.Ptr("test.txt")}
	got, err := p.BuildURL(input)
	assert.Nil(t, err)
	assert.Equal(t, "https://my-sandbox-alias-bs-apsr.abc.com/test.txt", got)

	// Without bucket, endpoint host as-is.
	got, err = p.BuildURL(&oss.OperationInput{})
	assert.Nil(t, err)
	assert.Equal(t, "https://abc.com/", got)

	// Missing accountId fails: the signing name cannot be resolved.
	bad := &agenticProvider{endpoint: endpoint, region: "cn-hangzhou", suffix: "bs-apsr", urlStyle: oss.UrlStyleVirtualHostedAlias}
	_, err = bad.BuildURL(input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "AccountId")
}

func TestAgenticEndpointProviderAliasHostLabelTooLong(t *testing.T) {
	endpoint, _ := url.Parse("https://abc.com")
	// alias label = "{bucket}-alias-ab-apsr" -> len(bucket) + 14
	suffixPart := "-alias-ab-apsr"
	assert.Equal(t, 14, len(suffixPart))
	p := &agenticProvider{
		endpoint:  endpoint,
		accountId: "1234567890123456",
		region:    "cn-hangzhou",
		suffix:    "ab-apsr",
		urlStyle:  oss.UrlStyleVirtualHostedAlias,
	}

	// Boundary: alias label == 63 (bucket 49) is allowed, even though the full
	// signing name for the same bucket far exceeds 63.
	okName := strings.Repeat("a", 49)
	assert.Equal(t, 63, len(okName+suffixPart))
	got, err := p.BuildURL(&oss.OperationInput{Bucket: oss.Ptr(okName)})
	assert.Nil(t, err)
	assert.Equal(t, "https://"+okName+suffixPart+".abc.com/", got)

	// Over limit: alias label == 64 (bucket 50) is rejected.
	longName := strings.Repeat("a", 50)
	assert.Equal(t, 64, len(longName+suffixPart))
	_, err = p.BuildURL(&oss.OperationInput{Bucket: oss.Ptr(longName)})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "exceeds the maximum length of 63 characters")
}

func TestBucketSpaceHelper(t *testing.T) {
	cfg := &oss.Config{}
	cfg.WithAccountId("1234567890123456")
	cfg.WithRegion("cn-hangzhou")

	helper := NewBucketSpaceHelper(cfg)
	name := helper.ToBucketName("my-sandbox")
	assert.Equal(t, "my-sandbox-1234567890123456-cn-hangzhou-bs-apsr", name)
}
