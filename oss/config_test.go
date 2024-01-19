package oss

import (
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/retry"
	"github.com/stretchr/testify/assert"
)

func TestConfigDefault(t *testing.T) {
	config := LoadDefaultConfig()
	assert.Equal(t, 3, config.RetryMaxAttempts)
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

	config.WithSignatureVersion(SignatureVersionV1)
	assert.Equal(t, SignatureVersionV1, *config.SignatureVersion)

	config.WithRegion("region")
	assert.Equal(t, "region", *config.Region)

	config.WithEndpoint("Endpoint")
	assert.Equal(t, "Endpoint", *config.Endpoint)

	config.WithRetryMaxAttempts(5)
	assert.Equal(t, 5, config.RetryMaxAttempts)

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
}
