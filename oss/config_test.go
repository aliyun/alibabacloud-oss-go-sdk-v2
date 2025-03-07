package oss

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/retry"
	"github.com/stretchr/testify/assert"
)

func TestConfigDefault(t *testing.T) {
	config := LoadDefaultConfig()
	assert.Nil(t, config.RetryMaxAttempts)
	assert.Nil(t, config.SignatureVersion)
	assert.Nil(t, config.Region)
	assert.Nil(t, config.Endpoint)
	assert.Nil(t, config.Retryer)
	assert.Nil(t, config.HttpClient)
	assert.Nil(t, config.CredentialsProvider)
	assert.Nil(t, config.UsePathStyle)
	assert.Nil(t, config.UseCName)
	assert.Nil(t, config.ConnectTimeout)
	assert.Nil(t, config.ReadWriteTimeout)
	assert.Nil(t, config.InsecureSkipVerify)
	assert.Nil(t, config.EnabledRedirect)
	assert.Nil(t, config.ProxyHost)
	assert.Nil(t, config.ProxyFromEnvironment)
	assert.Nil(t, config.UploadBandwidthlimit)
	assert.Nil(t, config.DownloadBandwidthlimit)
	assert.Nil(t, config.LogLevel)
	assert.Nil(t, config.LogPrinter)

	assert.Nil(t, config.DisableSSL)
	assert.Nil(t, config.UseDualStackEndpoint)
	assert.Nil(t, config.UseAccelerateEndpoint)
	assert.Nil(t, config.UseInternalEndpoint)

	assert.Nil(t, config.DisableUploadCRC64Check)
	assert.Nil(t, config.DisableDownloadCRC64Check)

	assert.Nil(t, config.AdditionalHeaders)
	assert.Nil(t, config.UserAgent)

	assert.Nil(t, config.CloudBoxId)
	assert.Nil(t, config.EnableAutoDetectCloudBoxId)

	config.WithSignatureVersion(SignatureVersionV1)
	assert.Equal(t, SignatureVersionV1, *config.SignatureVersion)

	config.WithRegion("region")
	assert.Equal(t, "region", *config.Region)

	config.WithEndpoint("Endpoint")
	assert.Equal(t, "Endpoint", *config.Endpoint)

	config.WithRetryMaxAttempts(5)
	assert.Equal(t, 5, *config.RetryMaxAttempts)

	config.WithRetryer(retry.NopRetryer{})
	assert.EqualValues(t, retry.NopRetryer{}, config.Retryer)

	config.WithHttpClient(http.DefaultClient)
	assert.EqualValues(t, http.DefaultClient, config.HttpClient)

	cred := credentials.NewAnonymousCredentialsProvider()
	config.WithCredentialsProvider(cred)
	assert.EqualValues(t, cred, config.CredentialsProvider)

	config.WithUsePathStyle(true)
	assert.Equal(t, true, *config.UsePathStyle)

	config.WithUseCName(true)
	assert.Equal(t, true, *config.UseCName)

	config.WithConnectTimeout(100 * time.Second)
	assert.Equal(t, 100*time.Second, *config.ConnectTimeout)

	config.WithReadWriteTimeout(50 * time.Second)
	assert.Equal(t, 50*time.Second, *config.ReadWriteTimeout)

	config.WithInsecureSkipVerify(true)
	assert.Equal(t, true, *config.InsecureSkipVerify)

	config.WithEnabledRedirect(true)
	assert.Equal(t, true, *config.EnabledRedirect)

	config.WithProxyHost("proxy")
	assert.Equal(t, "proxy", *config.ProxyHost)

	config.WithProxyFromEnvironment(true)
	assert.Equal(t, true, *config.ProxyFromEnvironment)

	config.WithUploadBandwidthlimit(int64(30))
	assert.Equal(t, int64(30), *config.UploadBandwidthlimit)

	config.WithDownloadBandwidthlimit(int64(60))
	assert.Equal(t, int64(60), *config.DownloadBandwidthlimit)

	config.WithLogLevel(LogError)
	assert.Equal(t, LogError, *config.LogLevel)

	config.WithLogPrinter(LogPrinterFunc(func(_ ...any) {}))
	assert.NotNil(t, config.LogPrinter)

	config.WithDisableSSL(true)
	assert.Equal(t, true, *config.DisableSSL)

	config.WithUseDualStackEndpoint(true)
	assert.Equal(t, true, *config.UseDualStackEndpoint)

	config.WithUseAccelerateEndpoint(true)
	assert.Equal(t, true, *config.UseAccelerateEndpoint)

	config.WithUseInternalEndpoint(true)
	assert.Equal(t, true, *config.UseInternalEndpoint)

	config.WithDisableUploadCRC64Check(true)
	assert.Equal(t, true, *config.DisableUploadCRC64Check)

	config.WithDisableDownloadCRC64Check(true)
	assert.Equal(t, true, *config.DisableDownloadCRC64Check)

	config.WithAdditionalHeaders([]string{"content-length"})
	assert.NotNil(t, config.AdditionalHeaders)
	assert.Len(t, config.AdditionalHeaders, 1)
	assert.Equal(t, "content-length", config.AdditionalHeaders[0])

	config.WithUserAgent("custom-ua")
	assert.NotNil(t, config.UserAgent)
	assert.Equal(t, "custom-ua", *config.UserAgent)

	config.WithCloudBoxId("cb-1234")
	assert.Equal(t, "cb-1234", *config.CloudBoxId)

	config.WithEnableAutoDetectCloudBoxId(true)
	assert.Equal(t, true, *config.EnableAutoDetectCloudBoxId)

	config.WithEnableAutoDetectCloudBoxId(false)
	assert.Equal(t, false, *config.EnableAutoDetectCloudBoxId)
}

func TestLogLevelEnvironmentVariable(t *testing.T) {
	oriloglevel := os.Getenv("OSS_SDK_LOG_LEVEL")
	defer func() {
		if oriloglevel == "" {
			os.Unsetenv("OSS_SDK_LOG_LEVEL")
		} else {
			os.Setenv("OSS_SDK_LOG_LEVEL", oriloglevel)
		}
	}()

	os.Setenv("OSS_SDK_LOG_LEVEL", "debug")
	config := LoadDefaultConfig()
	assert.Equal(t, LogDebug, *config.LogLevel)

	os.Setenv("OSS_SDK_LOG_LEVEL", "info")
	config = LoadDefaultConfig()
	assert.Equal(t, LogInfo, *config.LogLevel)

	os.Setenv("OSS_SDK_LOG_LEVEL", "error")
	config = LoadDefaultConfig()
	assert.Equal(t, LogError, *config.LogLevel)

	os.Setenv("OSS_SDK_LOG_LEVEL", "warn")
	config = LoadDefaultConfig()
	assert.Equal(t, LogWarn, *config.LogLevel)

	os.Setenv("OSS_SDK_LOG_LEVEL", "")
	config = LoadDefaultConfig()
	assert.Nil(t, config.LogLevel)

	os.Setenv("OSS_SDK_LOG_LEVEL", "off")
	config = LoadDefaultConfig()
	assert.Nil(t, config.LogLevel)
}
