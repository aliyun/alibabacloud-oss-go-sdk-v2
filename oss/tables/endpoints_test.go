package tables

import (
	"net/url"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestEndpointProviderBuildURL(t *testing.T) {
	endpoint, _ := url.Parse("https://cn-hangzhou.oss-tables.aliyuncs.com")

	t.Run("nil input", func(t *testing.T) {
		p := &endpointProvider{endpoint: endpoint, endpointType: oss.UrlStyleVirtualHosted}
		assert.Equal(t, "", p.BuildURL(nil))
	})

	t.Run("nil endpoint", func(t *testing.T) {
		p := &endpointProvider{endpoint: nil, endpointType: oss.UrlStyleVirtualHosted}
		assert.Equal(t, "", p.BuildURL(&oss.OperationInput{}))
	})

	t.Run("no bucket no key", func(t *testing.T) {
		p := &endpointProvider{endpoint: endpoint, endpointType: oss.UrlStyleVirtualHosted}
		input := &oss.OperationInput{}
		result := p.BuildURL(input)
		assert.Equal(t, "https://cn-hangzhou.oss-tables.aliyuncs.com/", result)
	})

	t.Run("virtual hosted with bucket", func(t *testing.T) {
		p := &endpointProvider{endpoint: endpoint, endpointType: oss.UrlStyleVirtualHosted}
		input := &oss.OperationInput{
			Bucket: oss.Ptr("acs:osstables:cn-hangzhou:123456:bucket/myBucket"),
		}
		result := p.BuildURL(input)
		assert.Equal(t, "https://myBucket-123456.cn-hangzhou.oss-tables.aliyuncs.com/", result)
	})

	t.Run("virtual hosted with bucket and key", func(t *testing.T) {
		p := &endpointProvider{endpoint: endpoint, endpointType: oss.UrlStyleVirtualHosted}
		input := &oss.OperationInput{
			Bucket: oss.Ptr("acs:osstables:cn-hangzhou:123456:bucket/myBucket"),
			Key:    oss.Ptr("myKey"),
		}
		result := p.BuildURL(input)
		assert.Equal(t, "https://myBucket-123456.cn-hangzhou.oss-tables.aliyuncs.com/myKey", result)
	})

	t.Run("path style with bucket no key", func(t *testing.T) {
		p := &endpointProvider{endpoint: endpoint, endpointType: oss.UrlStylePath}
		input := &oss.OperationInput{
			Bucket: oss.Ptr("acs:osstables:cn-hangzhou:123456:bucket/myBucket"),
		}
		result := p.BuildURL(input)
		assert.Equal(t, "https://cn-hangzhou.oss-tables.aliyuncs.com/", result)
	})

	t.Run("path style with bucket and key", func(t *testing.T) {
		p := &endpointProvider{endpoint: endpoint, endpointType: oss.UrlStylePath}
		input := &oss.OperationInput{
			Bucket: oss.Ptr("acs:osstables:cn-hangzhou:123456:bucket/myBucket"),
			Key:    oss.Ptr("myKey"),
		}
		result := p.BuildURL(input)
		assert.Equal(t, "https://cn-hangzhou.oss-tables.aliyuncs.com/myKey", result)
	})

	t.Run("path style no bucket with key", func(t *testing.T) {
		p := &endpointProvider{endpoint: endpoint, endpointType: oss.UrlStylePath}
		input := &oss.OperationInput{
			Key: oss.Ptr("someKey"),
		}
		result := p.BuildURL(input)
		assert.Equal(t, "https://cn-hangzhou.oss-tables.aliyuncs.com/someKey", result)
	})

	t.Run("http scheme", func(t *testing.T) {
		httpEndpoint, _ := url.Parse("http://cn-hangzhou.oss-tables.aliyuncs.com")
		p := &endpointProvider{endpoint: httpEndpoint, endpointType: oss.UrlStyleVirtualHosted}
		input := &oss.OperationInput{
			Bucket: oss.Ptr("acs:osstables:cn-hangzhou:999999:bucket/testBucket"),
		}
		result := p.BuildURL(input)
		assert.Equal(t, "http://testBucket-999999.cn-hangzhou.oss-tables.aliyuncs.com/", result)
	})

	t.Run("key with slash", func(t *testing.T) {
		p := &endpointProvider{endpoint: endpoint, endpointType: oss.UrlStyleVirtualHosted}
		input := &oss.OperationInput{
			Bucket: oss.Ptr("acs:osstables:cn-hangzhou:123456:bucket/myBucket"),
			Key:    oss.Ptr("path/to/key"),
		}
		result := p.BuildURL(input)
		assert.Equal(t, "https://myBucket-123456.cn-hangzhou.oss-tables.aliyuncs.com/path/to/key", result)
	})
}
