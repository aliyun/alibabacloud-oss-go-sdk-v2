package oss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEndpointScheme(t *testing.T) {

	assert.Equal(t, "", addEndpointScheme("", true))
	assert.Equal(t, "", addEndpointScheme("", false))

	assert.Equal(t, "https://123", addEndpointScheme("123", false))
	assert.Equal(t, "http://123", addEndpointScheme("123", true))

	assert.Equal(t, "http://123", addEndpointScheme("http://123", false))
	assert.Equal(t, "ftp://123", addEndpointScheme("ftp://123", false))
}

func TestEndpointFromRegion(t *testing.T) {
	// EndpointPublic
	assert.Equal(t, "https://oss-.aliyuncs.com", endpointFromRegion("", false, EndpointPublic))
	assert.Equal(t, "http://oss-.aliyuncs.com", endpointFromRegion("", true, EndpointPublic))
	assert.Equal(t, "https://oss-cn-hangzhou.aliyuncs.com", endpointFromRegion("cn-hangzhou", false, EndpointPublic))
	assert.Equal(t, "http://oss-cn-hangzhou.aliyuncs.com", endpointFromRegion("cn-hangzhou", true, EndpointPublic))

	// EndpointInternal
	assert.Equal(t, "https://oss--internal.aliyuncs.com", endpointFromRegion("", false, EndpointInternal))
	assert.Equal(t, "http://oss--internal.aliyuncs.com", endpointFromRegion("", true, EndpointInternal))
	assert.Equal(t, "https://oss-cn-hangzhou-internal.aliyuncs.com", endpointFromRegion("cn-hangzhou", false, EndpointInternal))
	assert.Equal(t, "http://oss-cn-hangzhou-internal.aliyuncs.com", endpointFromRegion("cn-hangzhou", true, EndpointInternal))

	// EndpointAccelerate
	assert.Equal(t, "https://oss-accelerate.aliyuncs.com", endpointFromRegion("", false, EndpointAccelerate))
	assert.Equal(t, "http://oss-accelerate.aliyuncs.com", endpointFromRegion("", true, EndpointAccelerate))
	assert.Equal(t, "https://oss-accelerate.aliyuncs.com", endpointFromRegion("cn-hangzhou", false, EndpointAccelerate))
	assert.Equal(t, "http://oss-accelerate.aliyuncs.com", endpointFromRegion("cn-hangzhou", true, EndpointAccelerate))

	// EndpointAccelerateOverseas
	assert.Equal(t, "https://oss-accelerate-overseas.aliyuncs.com", endpointFromRegion("", false, EndpointAccelerateOverseas))
	assert.Equal(t, "http://oss-accelerate-overseas.aliyuncs.com", endpointFromRegion("", true, EndpointAccelerateOverseas))
	assert.Equal(t, "https://oss-accelerate-overseas.aliyuncs.com", endpointFromRegion("cn-hangzhou", false, EndpointAccelerateOverseas))
	assert.Equal(t, "http://oss-accelerate-overseas.aliyuncs.com", endpointFromRegion("cn-hangzhou", true, EndpointAccelerateOverseas))

	// EndpointDualStack
	assert.Equal(t, "https://.oss.aliyuncs.com", endpointFromRegion("", false, EndpointDualStack))
	assert.Equal(t, "http://.oss.aliyuncs.com", endpointFromRegion("", true, EndpointDualStack))
	assert.Equal(t, "https://cn-hangzhou.oss.aliyuncs.com", endpointFromRegion("cn-hangzhou", false, EndpointDualStack))
	assert.Equal(t, "http://cn-hangzhou.oss.aliyuncs.com", endpointFromRegion("cn-hangzhou", true, EndpointDualStack))
}

func TestIsValidRegion(t *testing.T) {
	assert.True(t, isValidRegion("123-345"))
	assert.True(t, isValidRegion("abc"))
	assert.True(t, isValidRegion("abc-1234"))

	assert.False(t, isValidRegion("ABC"))
	assert.False(t, isValidRegion("#?23"))
	assert.False(t, isValidRegion(""))
	assert.False(t, isValidRegion("_"))
}
