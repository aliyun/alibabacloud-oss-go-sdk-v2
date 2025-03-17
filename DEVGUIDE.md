# Developer Guide
[简体中文](DEVGUIDE-CN.md)

Alibaba Cloud Object Storage Service (OSS) is a secure, cost-effective, and highly reliable cloud storage service that allows you to store large amounts of data. You can upload and download data from any application at any time and anywhere by calling API operations. You can also simply manage data by using the web console. OSS can store all types of files and is suitable for various websites, enterprises, and developers.

This development kit hides many lower-level operations, such as identity authentication, request retry, and error handling. You can access OSS by calling API operations without complex programming.

The development kit also provides practical modules, such as Uploader and Downloader, to automatically split large objects into multiple parts and transfer the parts in parallel.

You can refer to this developer guide to install, configure, and use the development kit.

Go to:

* [Installation](#installation)
* [Configuration](#configuration)
* [API operations](#api-operations)
* [Sample scenarios](#sample-scenarios)
* [Migration guide](#migration-guide)

# Installation

## Prerequisites

Go 1.18 or later is installed.
For more information about how to download and install Go, visit [Download and install](https://golang.org/doc/install).
You can run the following command to check the version of Go:
```
go version
```

## Install OSS SDK for Go

### Use Go Mod
Add the following dependencies to the go.mod file:
```
require (
    github.com/aliyun/alibabacloud-oss-go-sdk-v2 latest
)
```

### Use the source code
```
go get github.com/aliyun/alibabacloud-oss-go-sdk-v2
```

## Verify OSS SDK for Go
Run the following code to check the version of OSS SDK for Go:
```
package main

import (
  "fmt"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

func main() {
  fmt.Println("OSS Go SDK Version: ", oss.Version())
}
```

# Configuration
You can configure common settings for a client, such as the timeout period, log level, and retry policies. Most settings are optional.
However, you must specify the region and credentials for each client.  OSS SDK for Go uses the information to sign requests and send them to the correct region.

Subtopics in this section
* [Region](#Region)
* [Credentials](#credentials)
* [Endpoints](#endpoint)
* [HTTP client](#http-client)
* [Retry](#retry)
* [Logs](#logs)
* [Configuration parameters](#configuration-parameters)

## Load configurations
You can use several methods to configure a client. We recommend that you run the following sample code to configure a client:

```
package main

import (
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func main() {
  var (
    // In this example, the China (Hangzhou) region is used.
    region = "cn-hangzhou"

    // In this example, the credential is obtained from environment variables.
    provider credentials.CredentialsProvider = credentials.NewEnvironmentVariableCredentialsProvider()
  )

  cfg := oss.LoadDefaultConfig().
    WithCredentialsProvider(provider).
    WithRegion(region)
  }
```

## Region
You can specify a region to which you want the request to be sent, such as cn-hangzhou or cn-shanghai. For more information about the supported regions, see [Regions and endpoints](https://www.alibabacloud.com/help/en/oss/user-guide/regions-and-endpoints).
OSS SDK for Go does not have a default region. You must specify the `config.WithRegion` parameter to explicitly specify a region when you load the configurations. Example:
```
cfg := oss.LoadDefaultConfig().WithRegion("cn-hangzhou")
```

>**Note**: OSS SDK for Go uses a V4 signature by default. In this case, you must specify this parameter.

## Credentials

OSS SDK for Go requires credentials (AccessKey pair) to sign requests sent to OSS. In this case, you must explicitly specify the credentials. The following credential configurations are supported:
* [Environment variables](#environment-variables)
* [ECS instance role](#ecs-instance-role)
* [Static credentials](#static-credentials)
* [External processes](#external-processes)
* [RAM role](#ram-role)
* [OIDC-based SSO](#oidc-based-sso)
* [Custom credential provider](#custom-credential-provider)

### Environment variables

OSS SDK for Go supports obtaining credentials from environment variables. The following environment variables are supported:
* OSS_ACCESS_KEY_ID
* OSS_ACCESS_KEY_SECRET
* OSS_SESSION_TOKEN (Optional) 

The following sample code provides examples on how to configure environment variables.

1. Use Linux, OS X, or Unix
```
$ export OSS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
$ export OSS_ACCESS_KEY_SECRET=YOUR_ACCESS_KEY_SECRET
$ export OSS_SESSION_TOKEN=TOKEN
```

2. Use Windows
```
$ set OSS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
$ set OSS_ACCESS_KEY_SECRET=YOUR_ACCESS_KEY_SECRET
$ set OSS_SESSION_TOKEN=TOKEN
```

Use the credentials obtained from environment variables

```
provider := credentials.NewEnvironmentVariableCredentialsProvider()
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

### ECS instance role

If you want to access OSS from an Elastic Compute Service (ECS) instance, you can use a RAM role that is attached to the ECS instance to access OSS. You can attach a RAM role to an ECS instance to access OSS resources from the instance by using temporary access credentials that are obtained from Security Token Service (STS).

Use ECS instance role credentials

1. Specify the ECS instance role. Example: EcsRoleExample.
```
provider := credentials.NewEcsRoleCredentialsProvider(func(ercpo *credentials.EcsRoleCredentialsProviderOptions) {
	ercpo.RamRole = "EcsRoleExample"
})
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

2. Do not specify the ECS instance role
```
provider := credentials.NewEcsRoleCredentialsProvider()
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```
If you do not specify the ECS instance role name, the role name is automatically queried.

### Static credentials

You can hardcode the static credentials in your application to explicitly specify the AccessKey pair that you want to use to access OSS.

> **Note**: Do not embed the static credentials in the application. This method is used only for testing.

1. Long-term credentials
```
provider := credentials.NewStaticCredentialsProvider("AKId", "AKSecrect")
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

2. Temporary credentials
```
provider := credentials.NewStaticCredentialsProvider("AKId", "AKSecrect", "Token")
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

### External processes

You can obtain credentials from an external process in your application.
> **Note**:
> </br>Security risks may arise if unapproved processes or users run the commands that generate credentials.
> </br>The commands that generate credentials do not write any secret information to stderr or stdout because the information may be captured or recorded and may be exposed to unauthorized users.

You can run an external command to obtain long-term and temporary credentials. Format:
1. Long-term credentials
```
{
  "AccessKeyId" : "AKId",
  "AccessKeySecret" : "AKSecrect",
}
```

2. Temporary credentials
```
{
  "AccessKeyId" : "AKId",
  "AccessKeySecret" : "AKSecrect",
  "Expiration" : "2023-12-29T07:45:02Z",
  "SecurityToken" : "token",
}
```

Run the test-command command to obtain long-term credentials:

```
process := "test-command"
provider := credentials.NewProcessCredentialsProvider(process)
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

Run the test-command-sts command to obtain temporary credentials. The temporary credentials are different for each request.

```
process := "test-command-sts"
cprovider := credentials.NewProcessCredentialsProvider(process)
// NewCredentialsFetcherProvider automatically refreshes the credentials based on the Expiration parameter.
provider := credentials.NewCredentialsFetcherProvider(credentials.CredentialsFetcherFunc(func(ctx context.Context) (credentials.Credentials, error) {
  return cprovider.GetCredentials(ctx)
}))
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

### RAM role

If you want to authorize a RAM user to access OSS or access OSS across accounts, you can authorize the RAM user to assume a RAM role.

OSS SDK for Go does not directly provide access credentials. You need to use the [credentials-go](https://github.com/aliyun/credentials-go) Alibaba Cloud credential library. Example:

```
import (
  "context"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
  openapicred "github.com/aliyun/credentials-go/credentials"
)

// ...

config := new(openapicred.Config).
  // Which type of credential you want
  SetType("ram_role_arn").
  // AccessKeyId of your account
  SetAccessKeyId("AccessKeyId").
  // AccessKeySecret of your account
  SetAccessKeySecret("AccessKeySecret").
  // Format: acs:ram::USER_Id:role/ROLE_NAME
  SetRoleArn("RoleArn").
  // Role Session Name
  SetRoleSessionName("RoleSessionName").
  // Not required, limit the permissions of STS Token
  SetPolicy("Policy").
  // Not required, limit the Valid time of STS Token
  SetRoleSessionExpiration(3600)
	
arnCredential, gerr := openapicred.NewCredential(config)
provider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
  if gerr != nil {
    return credentials.Credentials{}, gerr
  }
  cred, err := arnCredential.GetCredential()
  if err != nil {
    return credentials.Credentials{}, err
  }
  return credentials.Credentials{
    AccessKeyID:     *cred.AccessKeyId,
    AccessKeySecret: *cred.AccessKeySecret,
    SecurityToken:   *cred.SecurityToken,
  }, nil
})

cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)

```

### OIDC-based SSO

You can also use the OpenID Connect (OIDC) authentication protocol in applications or services to access OSS. For more information about OIDC-based single sign-on (SSO), see [Overview of OIDC-based SSO](https://www.alibabacloud.com/help/en/ram/user-guide/overview-of-oidc-based-sso).

OSS SDK for Go does not directly provide the access credentials. You need to use the [credentials-go](https://github.com/aliyun/credentials-go) Alibaba Cloud credential library. Example:

```
import (
  "context"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
  openapicred "github.com/aliyun/credentials-go/credentials"
)

// ...

config := new(openapicred.Config).
  SetType("oidc_role_arn").
  SetOIDCProviderArn("OIDCProviderArn").
  SetOIDCTokenFilePath("OIDCTokenFilePath").
  SetRoleSessionName("RoleSessionName").
  SetPolicy("Policy").
  SetRoleArn("RoleArn").
  SetSessionExpiration(3600)
	
arnCredential, gerr := openapicred.NewCredential(config)
provider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
  if gerr != nil {
    return credentials.Credentials{}, gerr
  }
  cred, err := arnCredential.GetCredential()
  if err != nil {
    return credentials.Credentials{}, err
  }
  return credentials.Credentials{
    AccessKeyID:     *cred.AccessKeyId,
    AccessKeySecret: *cred.AccessKeySecret,
    SecurityToken:   *cred.SecurityToken,
  }, nil
})

cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)

```

### Custom credential provider

If the preceding credential configuration methods do not meet your requirements, you can specify the method that you want to use to obtain credentials. The following methods are supported:

1. Use credentials.CredentialsProvider
```
import (
  "context"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

type CustomerCredentialsProvider struct {
  // TODO
}

func NewCustomerCredentialsProvider() CredentialsProvider {
  return &CustomerCredentialsProvider{}
}

func (s CustomerCredentialsProvider) GetCredentials(_ context.Context) (credentials.Credentials, error) {
  // Return long-term credentials.
  return credentials.Credentials{AccessKeyID: "id", AccessKeySecret: "secret"}, nil
  // Return temporary credentials.
  //return credentials.Credentials{AccessKeyID: "id", AccessKeySecret: "secret",	SecurityToken: "token"}, nil
}

provider := NewCustomerCredentialsProvider()
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)

```

2. Use credentials.CredentialsProviderFunc

credentials.CredentialsProviderFunc is an easy-to-use encapsulation for the credentials.CredentialsProvider.

```
provider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
  // Return long-term credentials.
  return credentials.Credentials{AccessKeyID: "id", AccessKeySecret: "secret"}, nil
  // Return temporary credentials.
  //return credentials.Credentials{AccessKeyID: "id", AccessKeySecret: "secret",	SecurityToken: "token"}, nil
})

cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

3. Use credentials.CredentialsFetcherFunc

credentials.CredentialsFetcherFunc is an easy-to-use API operation for credentials.CredentialsFetcher.

credentials.CredentialsFetcher automatically refreshes credentials based on the Expiration parameter. You can use this method when you need to periodically update the credentials.

```
customerProvider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
  var (
    akId     string
    akSecret string
    token    string
    expires  *time.Time
  )

  // Obtain the temporary credentials and the expiration time of the credentials.
  ...

  // An error occurred.
  if err != nil {
    return credentials.Credentials{}, err
  }

  // The operation is successful.
  return credentials.Credentials{AccessKeyID: akId, AccessKeySecret: akSecret,	SecurityToken: token, Expires: expires}, nil
})

provider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
  return customerProvider.GetCredentials()
})

cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)

```

## Endpoint

You can use the Endpoint parameter to specify the endpoint of a request.

If the Endpoint parameter is not specified, OSS SDK for Go creates a public endpoint based on the region. For example, if the value of the Region parameter is cn-hangzhou, oss-cn-hangzhou.aliyuncs.com is created as a public endpoint.

You can modify parameters to create other endpoints, such as internal endpoints, transfer acceleration endpoints, and dual-stack endpoints that support IPv6 and IPv4. For more information about OSS domain name rules, see [OSS domain names](https://www.alibabacloud.com/help/en/oss/user-guide/oss-domain-names).

If you use a custom domain name to access OSS, you must specify this parameter. When you use a custom domain name to access a bucket, you must map the custom domain name to the default domain name of the bucket. For more information, see [Map a custom domain name to the default domain name of a bucket](https://www.alibabacloud.com/help/en/oss/user-guide/map-custom-domain-names-5).


### Access OSS by using standard domain names

In the following examples, the Region parameter is set to cn-hangzhou.

1. Use a public endpoint

```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou")

Or

cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("oss-cn-hanghzou.aliyuncs.com")
```

2. Use an internal endpoint

```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithUseInternalEndpoint(true)

Or

cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("oss-cn-hanghzou-internal.aliyuncs.com")
```

3. Use an OSS-accelerated endpoint
```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithUseAccelerateEndpoint(true)

Or

cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("oss-accelerate.aliyuncs.com")
```

4. Use a dual-stack endpoint
```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithUseDualStackEndpoint(true)

Or

cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("cn-hangzhou.oss.aliyuncs.com")
```

### Access OSS by using a custom domain name

In this example, the www.example-***.com domain name is mapped to the bucket-example bucket in the cn-hangzhou region.

```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("www.example-***.com").
  WithUseCName(true)
```

### Access private cloud or private domain

```
var (
  region = "YOUR Region"
  endpoint = "YOUR Endpoint"
)

cfg := oss.LoadDefaultConfig().
  WithRegion(region).
  WithEndpoint(endpoint)
```

## HTTP client

In most cases, the default HTTP client that uses the default configurations can meet the business requirements. You can also change the HTTP client or change the default configurations of the HTTP client to meet the requirements of specific environments.

The following section describes how to configure and create an HTTP client.

### Common configurations for an HTTP client

Modify common configurations by using config. The following table describes the parameters that you can configure.

| Parameter | Description | Example |
|:-------|:-------|:-------
| ConnectTimeout | The timeout period for establishing a connection. Default value: 5. Unit: seconds. | WithConnectTimeout(10 * time.Second) |
| ReadWriteTimeout | The timeout period for the application to read and write data. Default value: 10. Unit: seconds. | WithReadWriteTimeout(30 * time.Second) |
| InsecureSkipVerify | Specifies whether to skip SSL certificate verification. By default, the SSL certificates are verified. | WithInsecureSkipVerify(true) |
| EnabledRedirect | Specifies whether to enable HTTP redirection. By default, HTTP redirection is disabled. | WithEnabledRedirect(true) |
| ProxyHost | Specifies a proxy server. | WithProxyHost("http://user:passswd@proxy.example-***.com") |
| ProxyFromEnvironment | Specifies a proxy server by using environment variables. | WithProxyFromEnvironment(true) |
| UploadBandwidthlimit | The upper limit for the total upload bandwidth. Unit: KiB/s. | WithUploadBandwidthlimit(10*1024) |
| DownloadBandwidthlimit | The upper limit for the total download bandwidth. Unit: KiB/s. | WithDownloadBandwidthlimit(10*1024) |

Example

```
cfg := oss.LoadDefaultConfig().
  WithConnectTimeout(10 * time.Second).
  WithUploadBandwidthlimit(10*1024)
```

### Specify a custom HTTP client

If common parameters cannot meet the requirements of your business scenario, you can use WithHTTPClient to replace the default HTTP client.

For more information about the parameters that are not mentioned in the following example, visit [Transport](https://pkg.go.dev/net/http#Transport).

```
import (
  "crypto/tls"
  "net/http"
  "time"

  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/transport"
)

...

// Specify common timeout period or other parameters.
transConfig := transport.Config{
  // Specify the timeout period for a connection. Default value: 5. Unit: seconds.
  //ConnectTimeout: oss.Ptr(10 * time.Second),

  // Specify the timeout period for the application to read and write data. Default value: 10. Unit: seconds.
  //ReadWriteTimeout: oss.Ptr(20 * time.Second),

  // Specify the timeout period for an idle connection. Default value: 50. Unit: seconds.
  //IdleConnectionTimeout: oss.Ptr(40 * time.Second),

  // Specify the retention period of a network connection. Default value: 30. Unit: seconds.
  //KeepAliveTimeout: oss.Ptr(40 * time.Second),

  // Specify whether to enable HTTP redirection. By default, HTTP redirection is disabled.
  //EnabledRedirect: oss.Ptr(true),
}

// Specify http.Transport parameters.
var transports []func(*http.Transport)

// Specify the maximum number of connections. Default value: 100.
//transports = append(transports, transport.MaxConnections(200))

// If a request contains the Expect: 100-Continue header, it indicates the maximum period of time to wait for the first response header returned from the server after request headers are completely written. Default value: 1. Unit: seconds.
//transports = append(transports, transport.ExpectContinueTimeout(2*time.Second))

// Specify the earliest version of Transport Layer Security (TLS). Default value: TLS 1.2.
//transports = append(transports, transport.TLSMinVersion(tls.VersionTLS13))

// Specify whether to skip the SSL certificate verification. By default, the SSL certificate is verified.
//transports = append(transports, transport.InsecureSkipVerify(true))

// Specify other http.Transport parameters.
//transports = append(transports, func(t *http.Transport) {
//  t.DisableCompression
//})

customClient := transport.NewHttpClient(&transConfig, transports...)

cfg := oss.LoadDefaultConfig().WithHttpClient(customClient)
```

## Retry

You can specify the retry behaviors for HTTP requests.

### Default retry policy

If you do not specify a retry policy, OSS SDK for Go uses retry.Standard as the default retry policy of the client. Default configurations:

| Parameter | Description | Default value |
|:-------|:-------|:-------
| MaxAttempts | The maximum number of retries. | 3 |
| MaxBackoff | The maximum backoff time. | 20 seconds, 20 * time.Second |
| BaseDelay | The base delay. | 200 milliseconds, 200 * time.Millisecond |
| Backoff | The backoff algorithm. | FullJitter backoff, [0.0, 1.0) * min(2 ^ attempts * baseDealy, maxBackoff) |
| ErrorRetryables | The retryable errors. | For more information, visit [retryable_error.go](oss/retry/retryable_error.go). |

When a retryable error occurs, the system uses the provided configurations to delay and retry the request. The overall latency of a request increases as the number of retries increases. If the default configurations do not meet your business requirements, you must configure or modify the retry parameters.

### Modify the maximum number of retries

You can use one of the following methods to modify the maximum number of retries. For example, set the maximum number of retries to 5.

```
cfg := oss.LoadDefaultConfig().WithRetryMaxAttempts(5)

Or

cfg := oss.LoadDefaultConfig().WithRetryer(retry.NewStandard(func(ro *retry.RetryOptions) {
  ro.MaxAttempts = 5
}))
```

### Modify the backoff delay

For example, you can set BaseDelay to 500 milliseconds and the maximum backoff time to 25 seconds.

```
cfg := oss.LoadDefaultConfig().WithRetryer(retry.NewStandard(func(ro *retry.RetryOptions) {
  ro.MaxBackoff = 25 * time.Second
  ro.BaseDelay = 500 * time.Millisecond
}))
```

### Modify the backoff algorithm

For example, you can use a fixed-time backoff algorithm that has a delay of 2 seconds each time.

```
cfg := oss.LoadDefaultConfig().WithRetryer(retry.NewStandard(func(ro *retry.RetryOptions) {
  ro.Backoff = &retry.NewFixedDelayBackoff(2 * time.Second)
}))
```

### Change retryable errors

For example, you can add custom retryable errors.

```
type CustomErrorCodeRetryable struct {
}

func (*CustomErrorCodeRetryable) IsErrorRetryable(err error) bool {
  // Determine the error.
  // return true
  return false
}

errorRetryables := retry.DefaultErrorRetryables
errorRetryables = append(errorRetryables, &CustomErrorCodeRetryable{})

cfg := oss.LoadDefaultConfig().WithRetryer(retry.NewStandard(func(ro *retry.RetryOptions) {
  ro.ErrorRetryables = errorRetryables
}))
```

### Disable retry

If you want to disable all retry parameters, use retry.NopRetry.
```
cfg := oss.LoadDefaultConfig().WithRetryer(&retry.NopRetryer{})
```


## Logs

To facilitate troubleshooting, OSS SDK for Go provides the logging feature that uses debugging information in your application to debug and diagnose request issues.

If you want to use the logging feature, you must configure the log level. If the logging operation is not specified, logs are sent to the standard output (stdout) of the process by default.

Log level: oss.LogError, oss.LogWarn, oss.LogInfo, oss.LogDebug

Log operation: oss.LogPrinter, oss.LogPrinterFunc

For example, to enable the logging feature, set the log level to Info and the output to the standard error output (stderr)

```
cfg := oss.LoadDefaultConfig().
  WithLogLevel(oss.LogInfo).
  WithLogPrinter(oss.LogPrinterFunc(func(a ...any) {
    fmt.Fprint(os.Stderr, a...)
  }))
```

## Configuration parameters

Supported configuration parameters

| Parameter | Description | Example |
|:-------|:-------|:-------
| Region | (Required) The region to which the request is sent. | WithRegion("cn-hangzhou") |
| CredentialsProvider | (Required) The access credentials. | WithCredentialsProvider(provider) |
| Endpoint | The endpoint used to access OSS. | WithEndpoint("oss-cn-hanghzou.aliyuncs.com") |
| HttpClient | The HTTP client. | WithHttpClient(customClient) |
| RetryMaxAttempts | The maximum number of HTTP retries. Default value: 3. | WithRetryMaxAttempts(5) |
| Retryer | The retry configurations for HTTP requests. | WithRetryer(customRetryer) |
| ConnectTimeout | The timeout period for establishing a connection. Default value: 5. Unit: seconds. | WithConnectTimeout(10 * time.Second) |
| ReadWriteTimeout | The timeout period for the application to read and write data. Default value: 10. Unit: seconds. | WithReadWriteTimeout(30 * time.Second) |
| InsecureSkipVerify | Specifies whether to skip SSL certificate verification. By default, the SSL certificates are verified. | WithInsecureSkipVerify(true) |
| EnabledRedirect | Specifies whether to enable HTTP redirection. By default, HTTP redirection is disabled. | WithEnabledRedirect(true) |
| ProxyHost | Specifies a proxy server. | WithProxyHost("http://user:passswd@proxy.example-***.com") |
| ProxyFromEnvironment | Specifies a proxy server by using environment variables. | WithProxyFromEnvironment(true) |
| UploadBandwidthlimit | The upper limit for the total upload bandwidth. Unit: KiB/s. | WithUploadBandwidthlimit(10*1024) |
| DownloadBandwidthlimit | The upper limit for the total download bandwidth. Unit: KiB/s. | WithDownloadBandwidthlimit(10*1024) |
| SignatureVersion | The signature version. Default value: v4 | WithSignatureVersion(oss.SignatureVersionV1) |
| LogLevel | The log level. | WithLogLevel(oss.LogInfo) |
| LogPrinter | The log printing operation. | WithLogPrinter(customPrinter) |
| DisableSSL | Specifies that HTTPS is not used for requests. HTTPS is used by default. | WithDisableSSL(true) |
| UsePathStyle | The path request style, which is also known as the root domain name request style. By default, the default domain name of the bucket is used. | WithUsePathStyle(true) |
| UseCName | Specifies whether to use a custom domain name to access OSS. By default, a custom domain name is not used. | WithUseCName(true) |
| UseDualStackEndpoint | Specifies whether to use a dual-stack endpoint to access OSS. By default, a dual-stack endpoint is not used. | WithUseDualStackEndpoint(true) |
| UseAccelerateEndpoint | Specifies whether to use an OSS-accelerated endpoint to access OSS. By default, an OSS-accelerated endpoint is not used. | WithUseAccelerateEndpoint(true) |
| UseInternalEndpoint | Specifies whether to use an internal endpoint to access OSS. By default, an internal endpoint is not used. | WithUseInternalEndpoint(true) |
| DisableUploadCRC64Check | Specifies that CRC-64 is disabled during object upload. By default, CRC-64 is enabled. | WithDisableUploadCRC64Check(true) |
| DisableDownloadCRC64Check | Specifies that CRC-64 is disabled during object download. By default, CRC-64 is enabled. | WithDisableDownloadCRC64Check(true) |
|AdditionalHeaders| Specifies that additional headers to be signed. It's valid in V4 signature.|WithAdditionalHeaders([]string{"content-length"})
|UserAgent|Specifies user identifier appended to the User-Agent header.|WithUserAgent("user identifier")


# API operations

This section describes the API operations provided by OSS SDK for Go and how to call these API operations.

Subtopics in this section
* [Basic operations](#basic-operations)
* [Pre-signed URL](#pre-signed-url)
* [Paginator](#paginator)
* [Transfer Managers](#transfer-managers)
* [File-Like](#file-like)
* [Client-side encryption](#client-side-encryption)
* [Other API operations](#other-api-operations)
* [Comparison between upload and download operations](#comparison-between-upload-and-download-operations)

## Basic operations

OSS SDK for Go provides operations corresponding to RESTful APIs, which are called basic operations or low-level API operations. You can call the basic operations to manage OSS, such as creating a bucket and updating and deleting the configurations of a bucket.

The basic operations use the same naming conventions and use the following syntax:

```
func (c *Client) <OperationName>(ctx context.Context, request *<OperationName>Request, optFns ...func(*Options)) (result *<OperationName>Result, err error)
```

**Request parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|ctx|context.Context|The context of the request, which can be used to specify the total time limit of the request.
|request|*\<OperationName\>Request|Specifies the request parameters of a specific operation, such as bucket and key.
|optFns|...func(*Options)|Optional. Operation-level configuration parameters, such as the parameter used to modify the read and write timeout period when you call the operation this time.

**Response parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|result|*\<OperationName\>Result|The response to the operation. This parameter is valid when the value of err is nil.
|err|error|The status of the request. If the request fails, the value of err cannot be nil.

**Pointer parameters**:
<br/>In the parameters of the \<OperationName\>Request type, request parameters need to be passed to the operation in pointer mode. In the parameters of the \<OperationName\>Result type, response parameters need to be returned to the caller in pointer mode. The following sample code provides examples for ListObjectsRequest and ListObjectsResult:

```
type ListObjectsRequest struct {
  Bucket *string `input:"host,bucket,required"`
  Delimiter *string `input:"query,delimiter"`
  ...
  RequestPayer *string `input:"header,x-oss-request-payer"`
  RequestCommon
}

type ListObjectsResult struct {
  Name *string `xml:"Name"`
  Prefix *string `xml:"Prefix"`
  ...
  ResultCommon
}
```
To facilitate API calls, OSS SDK for Go provides the oss.Ptr function to convert parameters from non-pointer mode to pointer mode and provides the oss.To\<Type\> function to safely convert parameters from a pointer mode to a non-pointer mode.
For example, the oss.Ptr function converts string to * string. Conversely, oss.ToString converts *string to string.

Examples:
1. Create a bucket

```
client := oss.NewClient(cfg)
result, err := client.PutBucket(context.TODO(), &oss.PutBucketRequest{
  Bucket:          oss.Ptr("bucket"),
  Acl:             oss.BucketACLPrivate,
  ResourceGroupId: oss.Ptr("resource-group-id"),
  CreateBucketConfiguration: &oss.CreateBucketConfiguration{
    StorageClass: oss.StorageClassIA,
  },
})

if err != nil {
  log.Fatalf("failed to PutBucket %v", err)
}

fmt.Printf("PutBucket result:%v", result)
```

2. Copy an object and specify the operation-level read and write timeout period

```
client := oss.NewClient(cfg)
result, err := client.CopyObject(context.TODO(),
  &oss.CopyObjectRequest{
    Bucket:       oss.Ptr("bucket"),
    Key:          oss.Ptr("key"),
    SourceBucket: oss.Ptr("source-bucket"),
    SourceKey:    oss.Ptr("source-key"),
  },
  func(o *oss.Options) {
    o.OpReadWriteTimeout = oss.Ptr(30 * time.Second)
  },
)

if err != nil {
  log.Fatalf("failed to PutBucket %v", err)
}

fmt.Printf("CopyObject result, etg:%v", oss.ToString(result.ETag))
```

For more examples, refer to the sample directory.

## Pre-signed URL

You can call a specific operation to generate a pre-signed URL and use the pre-signed URL to grant temporary access to objects in a bucket or allow other users to upload specific objects to a bucket. You can use a pre-signed URL multiple times before the URL expires.

Syntax
```
func (c *Client) Presign(ctx context.Context, request any, optFns ...func(*PresignOptions)) (result *PresignResult, err error)
```

**Request parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|ctx|context.Context|The context of the request.
|request|any|Specifies the name of the API operation that is used to generate a signed URL. The value must be the same as the value of parameters of the <OperationName>Request type.
|optFns|...func(*PressignOptions)|Optional. Specifies the validity period of the pre-signed URL. If you do not specify this parameter, the pre-signed URL uses the default value, which is 15 minutes.

**Response parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|result|*PresignResult|The returned results, including the pre-signed URL, HTTP method, validity period, and request headers specified in the request.
|err|error|The status of the request. If the request fails, the value of err cannot be nil.

**Supported types of request parameters**
|Type|Operation
|:-------|:-------
|*GetObjectRequest|GetObject
|*PutObjectRequest|PutObject
|*HeadObjectRequest|HeadObject
|*InitiateMultipartUploadRequest|InitiateMultipartUpload
|*UploadPartRequest|UploadPart
|*CompleteMultipartUploadRequest|CompleteMultipartUpload
|*AbortMultipartUploadRequest|AbortMultipartUpload

**PressignOptions**
|Option|Type|Description
|:-------|:-------|:-------
|Expires|time.Duration|The validity period of the pre-signed URL. For example, if you want to set the validity period to 30 minutes, set Expires to 30 * time.Minute.
|Expiration|time.Time|The absolute expiration time of the pre-signed URL.

> **Note**: If you use the V4 signature algorithm, the validity period can be up to seven days. If you specify both Expiration and Expires, Expiration takes precedence.

**Response parameters of PresignResult**
|Parameter|Type|Description
|:-------|:-------|:-------
|Method|string|The HTTP method, which corresponds to the operation. For example, the HTTP method of the GetObject operation is GET.
|URL|string|The pre-signed URL.
|Expiration|time.Time|The time when the pre-signed URL expires.
|SignedHeaders|map[string]string|The request headers specified in the request. For example, if Content-Type is specified for PutObject, information about Content-Type is returned.


Examples
1. Generate a pre-signed URL for an object and download the object (GET request)
```
client := oss.NewClient(cfg)

result, err := client.Presign(context.TODO(), &oss.GetObjectRequest{
  Bucket: oss.Ptr("bucket"),
  Key:    oss.Ptr("key"),
})

resp, err := http.Get(result.URL)
```

2. Generate a pre-signed URL whose validity period is 10 minutes to upload an object, specify user metadata, and then upload the object (PUT request)
```
client := oss.NewClient(cfg)

result, err := client.Presign(context.TODO(), &oss.PutObjectRequest{
  Bucket:   oss.Ptr("bucket"),
  Key:      oss.Ptr("key"),
  Metadata: map[string]string{"user": "jack"}},
  oss.PresignExpires(10*time.Minute),
)

req, _ := http.NewRequest(result.Method, result.URL, nil)

for k, v := range result.SignedHeaders {
  req.Header.Add(k, v)
}

resp, err := http.DefaultClient.Do(req)
```

For more examples, refer to the sample directory.

## Paginator

For the list operations, a paged result, which contains a tag for retrieving the next page of results, is returned if the response results are too large to be returned in a single response. If you want to obtain the next page of results, you must specify the tag when you send the request.

OSS SDK for Go V2 provides a paginator that supports automatic pagination. If you call an API operation multiple times, OSS SDK for Go V2 automatically obtains the results of the next page. When you use the paginator, you need to only compile the code that is used to process the results.

The paginator contains an object in the \<OperationName\>Paginator format and the paginator creation method in the New\<OperationName\>Paginator format. The paginator creation method returns a paginator object that implements the HasNext and NextPage methods. The HasNext method is used to determine whether more pages exist and the NextPage method is used to call an API operation to obtain the next page.

The request parameter type of New\<OperationName\>Paginator is the same as that of \<OperationName\>.

The returned result type of \<OperationName\>Paginator.NextPage is the same as that of \<OperationName\>.

```
type <OperationName>Paginator struct {
...
}

func (p *<OperationName>Paginator) HasNext() bool {
	...
}

func (p *<OperationName>Paginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*<OperationName>Result, error) {
 ...
}

func (c *Client) New<OperationName>Paginator(request *<OperationName>Request, optFns ...func(*PaginatorOptions)) *<OperationName>Paginator
```

The following paginator objects are supported:
|Paginator object|Creation method|Corresponding list operation
|:-------|:-------|:-------
|ListObjectsPaginator|NewListObjectsPaginator|ListObjects: lists objects in a bucket.
|ListObjectsV2Paginator|NewListObjectsV 2Paginator|ListObjectsV2: lists objects in a bucket.
|ListObjectVersionsPaginator|NewListObjectVersionsPaginator|ListObjectVersions: lists object versions in a bucket.
|ListBucketsPaginator|NewListBucketsPaginator|ListBuckets: lists buckets.
|ListPartsPaginator|NewListPartsPaginator|ListParts: lists all uploaded parts of an upload task that has a specific upload ID.
|ListMultipartUploadsPaginator|NewListMultipartUploadsPaginator|ListMultipartUploads: lists the running multipart upload tasks in a bucket.

PaginatorOptions
|Parameter|Description
|:-------|:-------
|Limit|The maximum number of returned results.


In this example, ListObjects is used to describe how the paginator traverses all objects and how all objects are manually traversed.

```
// The paginator traverses all objects.
...
client := oss.NewClient(cfg)

p := client.NewListObjectsPaginator(&oss.ListObjectsRequest{
  Bucket: oss.Ptr("examplebucket"),
})

for p.HasNext() {
  page, err := p.NextPage(context.TODO())
  if err != nil {
    log.Fatalf("failed to get page %v", err)
  }

  for _, b := range page.Contents {
    fmt.Printf("Object:%v, %v, %v\n", oss.ToString(b.Key), oss.ToString(b.StorageClass), oss.ToTime(b.LastModified))
  }
}
```

```
// All objects are manually traversed.
...
client := oss.NewClient(cfg)

var marker *string
for {
  result, err := client.ListObjects(context.TODO(), &oss.ListObjectsRequest{
    Bucket: oss.Ptr("examplebucket"),
    Marker: marker,
  })
  if err != nil {
    log.Fatalf("failed to ListObjects %v", err)
  }

  for _, b := range result.Contents {
    fmt.Printf("Object:%v, %v, %v\n", oss.ToString(b.Key), oss.ToString(b.StorageClass), oss.ToTime(b.LastModified))
  }

  if result.IsTruncated {
    marker = result.NextMarker
  } else {
    break
  }
}
```

## Transfer Managers

For large object transfer scenarios, the Uploader, Downloader, and Copier modules are added to manage the upload, download, and copy of objects, respectively.

### Uploader

The uploader calls the multipart upload operation to split a large local file or stream into multiple smaller parts and upload the parts in parallel to improve upload performance.
</br>The multipart upload operation provides the resumable upload feature to record the upload progress. If a local file fails to be uploaded due to network interruptions or program crashes, and the local file cannot be uploaded after multiple retries, you can resume the upload from the position that is recorded in the checkpoint file.

```
type Uploader struct {
  ...
}

func (c *Client) NewUploader(optFns ...func(*UploaderOptions)) *Uploader

func (u *Uploader) UploadFrom(ctx context.Context, request *PutObjectRequest, body io.Reader, optFns ...func(*UploaderOptions)) (*UploadResult, error)

func (u *Uploader) UploadFile(ctx context.Context, request *PutObjectRequest, filePath string, optFns ...func(*UploaderOptions)) (*UploadResult, error)
```

**Request parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|ctx|context.Context|The context of the request.
|request|*PutObjectRequest|The request parameters of the upload request, which must be the same as the request parameters of the PutObject operation.
|body|io.Reader|The stream that you want to upload. If the body parameter only supports the io.Reader type, you must buffer the stream in the memory before you can upload it. If the body parameter supports the io.Reader, io.Seeker, and io.ReaderAt types, you do not need to cache the stream in the memory.
|filePath|string|The path of the local file that you want to upload.
|optFns|...func(*UploaderOptions)|Optional. The configuration options.


**UploaderOptions**
|Option|Type|Description
|:-------|:-------|:-------
|PartSize|int64|The part size. Default value: 6 MiB.
|ParallelNum|int|The number of the upload tasks in parallel. Default value: 3. ParallelNum takes effect only for a single operation.
|LeavePartsOnError|bool|Specifies whether to retain the uploaded parts when an upload task fails. By default, the uploaded parts are not retained.
|EnableCheckpoint|bool|Specifies whether to record the resumable upload progress in the checkpoint file. By default, no resumable upload progress is recorded.
|CheckpointDir|string|The path in which the checkpoint file is stored. Example: /local/dir/. This parameter is valid only if EnableCheckpoint is set to true.


If you use NewUploader to create an instance, you can specify several configuration parameters to specify custom object upload behaviors. You can also specify multiple configuration parameters to specify custom object upload behaviors each time you call an upload operation.

Specify the configuration parameters of Uploader
```
u := client.NewUploader(func(uo *oss.UploaderOptions) {
  uo.PartSize = 10 * 1024 * 1024
})
```

Specify configuration parameters for each download request
```
request := &oss.PutObjectRequest{Bucket: oss.Ptr("bucket"), Key: oss.Ptr("key")}
result, err := u.UploadFile(context.TODO(), request, "/local/dir/example", func(uo *oss.UploaderOptions) {
  uo.PartSize = 10 * 1024 * 1024
})
```

Examples

1. Use Uploader to upload a stream

```
...
client := oss.NewClient(cfg)

u := client.NewUploader()

var r io.Reader
// Use TODO to attach the io.Reader instance to r.

result, err := u.UploadFrom(context.TODO(),
  &oss.PutObjectRequest{
    Bucket: oss.Ptr("bucket"),
    Key:    oss.Ptr("key"),
  },
  r,
)

if err != nil {
  log.Fatalf("failed to UploadFile %v", err)
}

fmt.Printf("upload done, etag %v\n", oss.ToString(result.ETag))
```

2. Use Uploader to upload an object

```
...
client := oss.NewClient(cfg)

u := client.NewUploader()

result, err := u.UploadFile(context.TODO(),
  &oss.PutObjectRequest{
    Bucket: oss.Ptr("bucket"),
    Key:    oss.Ptr("key"),
  },
  "/local/dir/example",
)

if err != nil {
  log.Fatalf("failed to UploadFile %v", err)
}

fmt.Printf("upload done, etag %v\n", oss.ToString(result.ETag))
```

3. Use resumable upload to upload a local file
```
...
client := oss.NewClient(cfg)
u := client.NewUploader(func(uo *oss.UploaderOptions) {
  uo.CheckpointDir = "/local/dir/"
  uo.EnableCheckpoint = true
})

result, err := u.UploadFile(context.TODO(),
  &oss.PutObjectRequest{
    Bucket: oss.Ptr("bucket"),
    Key:    oss.Ptr("key"),
  },
  "/local/dir/example"
)

if err != nil {
  log.Fatalf("failed to UploadFile %v", err)
}

fmt.Printf("upload done, etag %v\n", oss.ToString(result.ETag))
```


### Downloader

Downloader uses range download to split a large object into multiple smaller parts and download the parts in parallel to improve download performance.
</br>The range download operation provides the resumable upload feature to record the download progress. If an object fails to be downloaded due to network interruptions or program crashes, and the object cannot be downloaded after multiple retries, you can resume the download from the position that is recorded in the checkpoint file.

```
type Downloader struct {
  ...
}

func (c *Client) NewDownloader(optFns ...func(*DownloaderOptions)) *Downloader

func (d *Downloader) DownloadFile(ctx context.Context, request *GetObjectRequest, filePath string, optFns ...func(*DownloaderOptions)) (result *DownloadResult, err error)
```

**Request parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|ctx|context.Context|The context of the request.
|request|*GetObjectRequest|The request parameters of the download request, which must be the same as the request parameters of the GetObject operation.
|filePath|string|The local path in which you want to store the downloaded object.
|optFns|...func(*DownloaderOptions)|Optional. The configuration options.


**DownloaderOptions**
|Option|Type|Description
|:-------|:-------|:-------
|PartSize|int64|The part size. Default value: 6 MiB.
|ParallelNum|int|The number of download tasks in parallel. Default value: 3. ParallelNum takes effect only for a single operation.
|EnableCheckpoint|bool|Specifies whether to record the download progress in the checkpoint file. By default, no download progress is recorded.
|CheckpointDir|string|The path in which the checkpoint file is stored. Example: /local/dir/. This parameter is valid only if EnableCheckpoint is set to true.
|VerifyData|bool|Specifies whether to verify the CRC-64 of the downloaded object when the download is resumed. By default, the CRC-64 is not verified. This parameter is valid only if EnableCheckpoint is set to true.
|UseTempFile|bool|Specifies whether to use a temporary file when you download an object. A temporary file is used by default. The object is downloaded to the temporary file. Then, the temporary file is renamed and uses the same name as the object that you want to download.


When you use NewDownloader to create an instance, you can specify several configuration parameters to specify custom object download behaviors. You can also specify multiple configuration parameters to specify custom object download behaviors each time you call a download operation.

Specify configuration parameters for Downloader
```
d := client.NewDownloader(func(do *oss.DownloaderOptions) {
  do.PartSize = 10 * 1024 * 1024
})
```

Specify configuration parameters for each download request
```
request := &oss.GetObjectRequest{Bucket: oss.Ptr("bucket"), Key: oss.Ptr("key")}
d.DownloadFile(context.TODO(), request, "/local/dir/example", func(do *oss.DownloaderOptions) {
  do.PartSize = 10 * 1024 * 1024
})
```

Example

1. Use Downloader to download an object as a local file

```
...
client := oss.NewClient(cfg)

d := client.NewDownloader()

d.DownloadFile(context.TODO(),
  &oss.GetObjectRequest{
    Bucket: oss.Ptr("bucket"),
    Key:    oss.Ptr("key"),
  },
  "/local/dir/example",
)
```

### Copier
If you want to copy an object from a bucket to another bucket or modify the attributes of an object, you can call the CopyObject operation or the UploadPartCopy operation.
</br>These two API operations are suitable for scenarios, such as:
* You can call the CopyObject operation to copy objects smaller than 5 GiB.
* The UploadPartCopy operation does not support the metadata command (x-oss-metadata-directive) and tag command (x-oss-tagging-directive) parameters.
   When you copy data, you must configure the metadata and tags that you want to copy.
* The server optimizes the CopyObject operation and the operation has the shallow copy capability to copy large objects in specific scenarios.

Copier provides common copy operations, hides the differences and implementation details of the operations, and automatically selects the appropriate operation to copy objects according to the request parameters of the copy task.

```
type Copier struct {
  ...
}

func (c *Client) NewCopier(optFns ...func(*CopierOptions)) *Copier

func (c *Copier) Copy(ctx context.Context, request *CopyObjectRequest, optFns ...func(*CopierOptions)) (*CopyResult, error)
```

**Request parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|ctx|context.Context|The context of the request.
|request|*CopyObjectRequest|The request parameters of the copy request, which must be the same as the request parameters of the CopyObject operation.
|optFns|...func(*CopierOptions)|Optional. The configuration options.


**CopierOptions:**
|Option|Type|Description
|:-------|:-------|:-------
|PartSize|int64|The part size. Default value: 64 MiB.
|ParallelNum|int|The number of the copy tasks in parallel. Default value: 3. ParallelNum takes effect only for a single operation.
|MultipartCopyThreshold|int64|The minimum object size for calling the multipart copy operation. Default value: 200 MiB.
|LeavePartsOnError|bool|Specifies whether to retain the copied parts when the copy task fails. By default, the copied parts are not retained.
|DisableShallowCopy|bool|Specifies that the shallow copy capability is not used. By default, the shallow copy capability is used.


When you use NewCopier to create an instance, you can specify several configuration parameters to specify custom object copy behaviors. You can also specify multiple configuration parameters to specify custom object copy behaviors each time you call a copy operation.

Specify configuration parameters for Copier
```
d := client.NewCopier(func(co *oss.CopierOptions) {
  co.PartSize = 100 * 1024 * 1024
})
```

Specify configuration parameters for each copy request
```
request := &oss.CopyObjectRequest{
  Bucket:       oss.Ptr("bucket"),
  Key:          oss.Ptr("key"),
  SourceBucket: oss.Ptr("src-bucket"),
  SourceKey:    oss.Ptr("src-key"),
}
copier.Copy(context.TODO(), request, func(co *oss.CopierOptions) {
  co.PartSize = 100 * 1024 * 1024
})
```

> **Note**:
> </br>When you copy an object, the CopyObjectRequest.MetadataDirective parameter specifies whether to copy the object metadata. By default, the metadata of a source object is copied.
> </br>When you copy an object, the CopyObjectRequest.TaggingDirective parameter specifies whether to copy the object tags. By default, the tags of a source object tag are copied.


Examples

1. Copy an object, including the object metadata and tags
```
...
client := oss.NewClient(cfg)
copier := client.NewCopier()

result, err := copier.Copy(context.TODO(), &oss.CopyObjectRequest{
  Bucket:       oss.Ptr("bucket"),
  Key:          oss.Ptr("key"),
  SourceBucket: oss.Ptr("src-bucket"),
  SourceKey:    oss.Ptr("src-key"),
})

if err != nil {
  log.Fatalf("failed to UploadFile %v", err)
}

fmt.Printf("copy done, etag %v\n", oss.ToString(result.ETag))
```

2. Copy an object without the object metadata and tags
```
...
client := oss.NewClient(cfg)
copier := client.NewCopier()

result, err := copier.Copy(context.TODO(), &oss.CopyObjectRequest{
  Bucket:            oss.Ptr("bucket"),
  Key:               oss.Ptr("key"),
  SourceBucket:      oss.Ptr("src-bucket"),
  SourceKey:         oss.Ptr("src-key"),
  MetadataDirective: oss.Ptr("Replace"),
  TaggingDirective:  oss.Ptr("Replace"),
})

if err != nil {
  log.Fatalf("failed to UploadFile %v", err)
}

fmt.Printf("copy done, etag %v\n", oss.ToString(result.ETag))
```

3. Change the storage class of an copied object to Standard

```
...
client := oss.NewClient(cfg)
copier := client.NewCopier()

result, err := copier.Copy(context.TODO(), &oss.CopyObjectRequest{
  Bucket:       oss.Ptr("bucket"),
  Key:          oss.Ptr("key"),
  SourceBucket: oss.Ptr("src-bucket"),
  SourceKey:    oss.Ptr("src-key"),
  StorageClass: oss.StorageClassStandard,
})

if err != nil {
  log.Fatalf("failed to UploadFile %v", err)
}

fmt.Printf("copy done, etag %v\n", oss.ToString(result.ETag))
```

## File-Like

The File-Like operation is added to simulate the read and write behaviors on objects in a bucket.

The following methods are supported:
* ReadOnlyFile
* AppendOnlyFile

### ReadOnlyFile

You can only read objects in a bucket. The ReadOnlyFile method provides the Single Stream and Prefetch modes. You can change the number of tasks in parallel to improve the read speed. At the same time, the ReadOnlyFile method provides a reconnection mechanism, which has strong robustness in a more complex network environment.

```
type ReadOnlyFile struct {
...
}

func (c *Client) OpenFile(ctx context.Context, bucket string, key string, optFns ...func(*OpenOptions)) (file *ReadOnlyFile, err error)
```

**Request parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|ctx|context.Context|The context of the request.
|bucket|string|The name of the bucket.
|key|string|The name of the object.
|optFns|...func(*OpenOptions)|Optional. The configuration options when you open an object.

**Response parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|file|*ReadOnlyFile|The instance of ReadOnlyFile. This parameter is valid when the value of err is nil.
|err|error|The status of ReadOnlyFile. If an error occurs, the value of err cannot be nil.

**OpenOptions:**
|Option|Type|Description
|:-------|:-------|:-------
|Offset|int64|The initial offset when the object is opened. Default value: 0.
|VersionId|*string|The version number of the specified object. This parameter is valid only if multiple versions of the object exist.
|RequestPayer|*string|Specifies that if pay-by-requester is enabled, RequestPayer must be set to requester.
|EnablePrefetch|bool|Specifies whether to enable the Prefetch mode. By default, the prefetch mode is disabled.
|PrefetchNum|int|The number of prefetched chunks. Default value: 3. This parameter is valid when the Prefetch mode is enabled.
|ChunkSize|int64|The size of each prefetched chunk. Default value: 6 MiB. This parameter is valid when the Prefetch mode is enabled.
|PrefetchThreshold|int64|The number of bytes to be read in sequence before the prefetch mode is used. Default value: 20 MiB. This parameter is valid when the Prefetch mode is enabled.

**ReadOnlyFile**:
|Operation name|Description
|:-------|:-------
|Close() error|Closes the file handles to release resources, such as memory and active sockets.
|Read(p []byte) (int, error)|Reads a byte whose length is len(p) from the data source, stores the byte in p, and returns the number of bytes that are read and the encountered errors.
|Seek(offset int64, whence int) (int64, error)|Specifies the offset for the next read or write. Valid values of whence: 0: the head. 1: the current offset. 2: the tail.
|Stat() (os.FileInfo, error)|Queries the object information, including the object size, last modified time, and metadata.

> **Note**: If the Prefetch mode is enabled and multiple out-of-order reads occur, the Single Stream mode is automatically used.

Examples

1. Read the entire object by using the Single Stream mode
```
...
client := oss.NewClient(cfg)

f, err := client.OpenFile(context.TODO(), "bucket", "key")

if err != nil {
  log.Fatalf("failed to open file %v", err)
}
defer f.Close()

written, err := io.Copy(io.Discard, f)

if err != nil {
  log.Fatalf("failed to read file %v", err)
}

fmt.Print("read data count:%v", written)
```

2. Read the entire object by using the Prefetch mode
```
...
client := oss.NewClient(cfg)

f, err := client.OpenFile(context.TODO(),
  "bucket",
  "key",
  func(oo *oss.OpenOptions) {
    oo.EnablePrefetch = true
  })

if err != nil {
  log.Fatalf("failed to open file %v", err)
}

defer f.Close()

written, err := io.Copy(io.Discard, f)

if err != nil {
  log.Fatalf("failed to read file %v", err)
}

fmt.Print("read data count:%v", written)
```

3. Read remaining data from a specific position by using the Seek method

```
...
client := oss.NewClient(cfg)

f, err := client.OpenFile(context.TODO(), "bucket", "key")

if err != nil {
  log.Fatalf("failed to open file %v", err)
}

defer f.Close()

// Query the object information.
info, _ := f.Stat()

// Query basic attributes.
fmt.Printf("size:%v, mtime:%v\n", info.Size(), info.ModTime())

// Query object metadata.
if header, ok := info.Sys().(http.Header); ok {
  fmt.Printf("content-type:%v\n", header.Get(oss.HTTPHeaderContentType))
}

// Specify the offset of the object, such as starting from 123.
_, err = f.Seek(123, io.SeekStart)
if err != nil {
  log.Fatalf("failed to seek file %v", err)
}

written, err := io.Copy(io.Discard, f)

if err != nil {
  log.Fatalf("failed to read file %v", err)
}

fmt.Print("read data count:%v", written)
```

### AppendOnlyFile

You can call the AppendObject operation to upload data by appending data to an existing object. If the object does not exist, the AppendObject operation creates an appendable object. If the object exists but is not the appendable type, an error is returned.

```
type AppendOnlyFile struct {
...
}

func (c *Client) AppendFile(ctx context.Context, bucket string, key string, optFns ...func(*AppendOptions)) (*AppendOnlyFile, error)
```

**Request parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|ctx|context.Context|The context of the request.
|bucket|string|The name of the bucket.
|key|string|The name of the object.
|optFns|...func(*AppendOptions)|Optional. The configuration options when you append data to an existing object.

**Response parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|file|*AppendOnlyFile|The instance of AppendOnlyFile. This parameter is valid when the value of err is nil.
|err|error|The status of AppendOnlyFile. If an error occurs, the value of err cannot be nil.

**AppendOptions**:
|Option|Type|Description
|:-------|:-------|:-------
|RequestPayer|*string|Specifies that if pay-by-requester is enabled, RequestPayer must be set to requester.
|CreateParameter|*AppendObjectRequest|Specifies the object metadata, including the content type, metadata, ACL, and storage class of the object, the first time you upload an object.

**AppendOnlyFile**:
|Operation name|Description
|:-------|:-------
|Close() error|Closes the file handles to release resources.
|Write (b []byte) (int, error)|Writes the data in b to the data stream and returns the number of written bytes and the encountered errors.
|WriteFrom(r io.Reader) (int64, error)|Writes the data in r to the data stream and returns the number of written bytes and the encountered errors.
|Stat() (os.FileInfo, error)|Queries the object information, including the object size, last modified time, and metadata.


Examples

1. Combine multiple local files into a file
```
...
client := oss.NewClient(cfg)

f, err := client.AppendFile(context.TODO(), "bucket", "key")
if err != nil {
  log.Fatalf("failed to append file %v", err)
}

defer f.Close()

// example1.txt
lf, err := os.Open("/local/dir/example1.txt")
if err != nil {
  log.Fatalf("failed to local file %v", err)
}

_, err = f.WriteFrom(lf)

if err != nil {
  log.Fatalf("failed to append file %v", err)
}
lf.Close()

// example2.txt
lf, err = os.Open("/local/dir/example2.txt")
if err != nil {
  log.Fatalf("failed to local file %v", err)
}

_, err = f.WriteFrom(lf)

if err != nil {
  log.Fatalf("failed to append file %v", err)
}
lf.Close()

// example3.txt
lb, err := os.ReadFile("/local/dir/example3.txt")

_, err = f.Write(lb)

if err != nil {
  log.Fatalf("failed to append file %v", err)
}
```

2. Specify the ACL and storage class of an object when you combine data
```
...
client := oss.NewClient(cfg)

f, err := client.AppendFile(context.TODO(),
  "bucket",
  "key",
  func(ao *oss.AppendOptions) {
    ao.CreateParameter = &oss.AppendObjectRequest{
      Acl: oss.ObjectACLPrivate,
      Metadata: map[string]string{
        "user": "jack",
      },
      Tagging: oss.Ptr("key=value"),
    }
  },
)

if err != nil {
  log.Fatalf("failed to append file %v", err)
}

defer f.Close()

_, err = f.Write([]byte("hello"))

if err != nil {
  log.Fatalf("failed to append file %v", err)
}

_, err = f.Write([]byte("world"))

if err != nil {
  log.Fatalf("failed to append file %v", err)
}

info, err := f.Stat()
if err != nil {
  log.Fatalf("failed to stat file %v", err)
}

fmt.Printf("size:%v, mtime:%v\n", info.Size(), info.ModTime())

if header, ok := info.Sys().(http.Header); ok {
  fmt.Printf("user:%v\n", header.Get("x-oss-meta-user"))
}
```


## Client-side encryption

If client-side encryption is enabled, objects are locally encrypted before they are uploaded to OSS. Only the owner of the customer master key (CMK) can decrypt the objects. This way, data security during data transmission and storage is enhanced.

> **Note**:
> </br>If you enable client-side encryption, you must ensure the integrity and validity of the CMK.
> </br>When you copy or migrate encrypted data, you are responsible for the integrity and validity of the object metadata related to client-side encryption.

For more information, see [Client-side encryption](https://www.alibabacloud.com/help/en/oss/user-guide/client-side-encryption).

To use client-side encryption, you must create an instance for the client for client-side encryption and call related API operations. Your objects are automatically encrypted and decrypted as part of the request.

```
type EncryptionClient struct {
  ...
}

func NewEncryptionClient(c *Client, masterCipher crypto.MasterCipher, optFns ...func(*EncryptionClientOptions)) (eclient *EncryptionClient, err error)
```

**Request parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|c|*Client|The instance of the client for non-encryption.
|masterCipher|crypto.MasterCipher|The CMK instance for encrypting and decrypting data keys.
|optFns|...func(*EncryptionClientOptions)|Optional. The configuration options of the client for client-side encryption.

**Response parameters**
|Parameter|Type|Description
|:-------|:-------|:-------
|eclient|*EncryptionClient|The instance of the client for client-side encryption. This parameter is valid when the value of err is nil.
|err|error|The status of the created client for client-side encryption. If an error occurs, the value of err cannot be nil.

**EncryptionClientOptions:**
|Option|Type|Description
|:-------|:-------|:-------
|MasterCiphers|[]crypto.MasterCipher|The instance group of CMKs, which is used to decrypt data keys.

**The API operations of EncryptionClient**
|Basic operation|Description
|:-------|:-------
|GetObjectMeta|Queries part of the object metadata.
|HeadObject|Queries part of the object metadata.
|GetObject|Downloads an object and automatically decrypts it.
|PutObject|Uploads an object and automatically encrypts it.
|InitiateMultipartUpload|Initiates a multipart upload task and the context for encryption in multipart upload (EncryptionMultiPartContext).
|UploadPart|Initiates a multipart upload task, uploads parts, and automatically encrypts the parts. When you call the API operation, you must specify context for encryption in multipart upload.
|CompleteMultipartUpload|Combines all parts into a complete object after all parts are uploaded.
|AbortMultipartUpload|Cancels a multipart upload task and deletes the uploaded parts.
|ListParts|Lists all parts that are uploaded by using a specific upload ID.
|**Advanced operation**|**Description**
|NewDownloader|Creates a downloader instance.
|NewUploader|Creates an uploader instance.
|OpenFile|Creates a ReadOnlyFile instance.
|**Auxiliary opeation**|**Description**
|Unwrap|Queries an instance of the client for non-encryption and uses the instance to perform other basic operations.

> **Note**: EncryptionClient uses the same operation naming rules and calling methods as Client. For more information, refer to other chapters of the guide.

### Use an RSA-based CMK

**Create a client for RSA-based client-side encryption**

```
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/crypto"

cfg := oss.LoadDefaultConfig().
  WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
  WithRegion("your region")

client := oss.NewClient(cfg)

// Specify a description for the CMK. The description cannot be changed after it is specified. You can specify only one description for a CMK.
// If all objects use the same CMK, the description of the CMK can be empty. However, you cannot change the CMK.
// If you do not specify a description for the CMK, the client cannot determine which CMK to use for decryption.
// We recommend that you specify a description for each CMK. The client saves the mappings between the CMK and the description.
materialDesc := make(map[string]string)
materialDesc["desc"] = "your master encrypt key material describe information"

// Create a client that contains only a CMK for client-side encryption.
mc, err := crypto.CreateMasterRsa(materialDesc, "yourRsaPublicKey", "yourRsaPrivateKey")
eclient, err := NewEncryptionClient(client, mc)

// Create a client that contains a CMK and multiple decryption keys for client-side encryption.
// When you decrypt objects, the description of the decryption keys is decrypted first. If you cannot use the decryption keys to decrypt objects, the CMK is used.
//decryptMC := []crypto.MasterCipher{
//	// TODO
//}
//eclient, err := oss.NewEncryptionClient(client, mc, func(eco *oss.EncryptionClientOptions) {
//	eco.MasterCiphers = decryptMC
//})
```

**Use a client for client-side encryption to upload or download objects**
```
...
eclient, err := NewEncryptionClient(client, mc)

// Use PutObject
_, err = eclient.PutObject(context.TODO(), &oss.PutObjectRequest{
  Bucket: oss.Ptr("bucket"),
  Key:    oss.Ptr("key"),
  Body:   bytes.NewReader([]byte("hello world")),
})

if err != nil {
  log.Fatalf("failed to PutObject %v", err)
}

// Use GetObject
gresult, err := eclient.GetObject(context.TODO(), &oss.GetObjectRequest{
  Bucket: oss.Ptr("bucket"),
  Key:    oss.Ptr("key"),
})

if err != nil {
  log.Fatalf("failed to GetObject %v", err)
}

io.Copy(io.Discard, gresult.Body)
gresult.Body.Close()

// Use Downloader
d := eclient.NewDownloader()
_, err = d.DownloadFile(context.TODO(),
  &oss.GetObjectRequest{
    Bucket: oss.Ptr("bucket"),
    Key:    oss.Ptr("key"),
  },
  "/local/dir/example",
)
if err != nil {
  log.Fatalf("failed to DownloadFile %v", err)
}

// Use Uploader
u := eclient.NewUploader()
_, err = u.UploadFile(context.TODO(),
  &oss.PutObjectRequest{
    Bucket: oss.Ptr("bucket"),
    Key:    oss.Ptr("key"),
  },
  "/local/dir/example",
)
if err != nil {
  log.Fatalf("failed to UploadFile %v", err)
}

// Use ReadOnlyFile
f, err := eclient.OpenFile(context.TODO(), "bucket", "key")
if err != nil {
  log.Fatalf("failed to OpenFile %v", err)
}
defer f.Close()

_, err = io.Copy(io.Discard, f)
if err != nil {
  log.Fatalf("failed to Copy %v", err)
}
```

**Use a client for client-side encryption to upload data by running a multipart upload task**
</br>Example: upload memory data that is smaller than 500 KB
```
...
eclient, err := NewEncryptionClient(client, mc)

var (
  bucketName string = "bucket"
  objectName string = "key"
  length            = int64(500 * 1024)
  partSize          = int64(200 * 1024)
  partsNum          = int(length/partSize + 1)
  data              = make([]byte, length, length)
)

// Specify the part size and the total object size by using the client for client-side encryption.
initResult, err := eclient.InitiateMultipartUpload(context.TODO(), &oss.InitiateMultipartUploadRequest{
  Bucket: oss.Ptr(bucketName),
  Key:    oss.Ptr(objectName),
  CSEPartSize: oss.Ptr(partSize),
  CSEDataSize: oss.Ptr(length),
})

var parts oss.UploadParts
for i := 0; i < partsNum; i++ {
  start := int64(i) * partSize
  end := start + partSize
  if end > length {
    end = length
  }

  // Specify the context for encryption in the multipart upload task by using the client for client-side encryption.
  upResult, err := eclient.UploadPart(context.TODO(), &oss.UploadPartRequest{
    Bucket:     oss.Ptr(bucketName),
    Key:        oss.Ptr(objectName),
    UploadId:   initResult.UploadId,
    PartNumber: int32(i + 1),
    CSEMultiPartContext: initResult.CSEMultiPartContext,
    Body:                bytes.NewReader(data[start:end]),
  })

  if err != nil {
    log.Fatalf("failed to UploadPart %v", err)
  }
  parts = append(parts, oss.UploadPart{PartNumber: int32(i + 1), ETag: upResult.ETag})
}

sort.Sort(parts)
_, err = eclient.CompleteMultipartUpload(context.TODO(), &oss.CompleteMultipartUploadRequest{
  Bucket:                  oss.Ptr(bucketName),
  Key:                     oss.Ptr(objectName),
  UploadId:                initResult.UploadId,
  CompleteMultipartUpload: &oss.CompleteMultipartUpload{Parts: parts},
})

if err != nil {
  log.Fatalf("failed to CompleteMultipartUpload %v", err)
}
```

### Use a custom CMK
If the RSA-based CMK cannot meet your requirements, you can use a custom CMK. Syntax of a custom CMK:
```
type MasterCipher interface {
  Encrypt([]byte) ([]byte, error)
  Decrypt([]byte) ([]byte, error)
  GetWrapAlgorithm() string
  GetMatDesc() string
}
```
**The API operations of MasterCipher**
|Operation name|Description
|:-------|:-------
|Encrypt|Encrypts the data key and the initial values (IV) of the encrypted data.
|Decrypt|Decrypts the data key and the initial values (IV) of the encrypted data.
|GetWrapAlgorithm|Returns the encryption algorithm of the data key. Recommended format: algorithm/mode/padding. Example: RSA/NONE/PKCS1Padding.
|GetMatDesc|Returns the description of the CMK in JSON format.

Example:

```
...
type MasterCustomCipher struct {
  MatDesc    string
  SecrectKey string
}

func (mrc MasterCustomCipher) GetWrapAlgorithm() string {
  return "Custom/None/NoPadding"
}

func (mrc MasterCustomCipher) GetMatDesc() string {
  return mrc.MatDesc
}

func (mrc MasterCustomCipher) Encrypt(plainData []byte) ([]byte, error) {
  // TODO
}

func (mrc MasterCustomCipher) Decrypt(cryptoData []byte) ([]byte, error) {
  // TODO
}

func MasterCustomCipher(matDesc map[string]string, secrectKey string) (crypto.MasterCipher, error) {
  var jsonDesc string
  if len(matDesc) > 0 {
    b, err := json.Marshal(matDesc)
    if err != nil {
      return nil, err
    }
    jsonDesc = string(b)
  }
  return MasterCustomCipher{MatDesc: jsonDesc, SecrectKey: secrectKey}, nil
}

client := oss.NewClient(cfg)
materialDesc := make(map[string]string)
materialDesc["desc"] = "your master encrypt key material describe information"
mc, err := MasterCustomCipher(materialDesc, "yourSecrectKey")
eclient, err := NewEncryptionClient(client, mc)
```

## Other API operations

The following easy-to-use operations are encapsulated to improve user experience.  

| Operation | Description |
|:-------|:-------
| IsObjectExist | Determines whether an object exists. |
| IsBucketExist | Determines whether a bucket exists. |
| PutObjectFromFile | Uploads a local file to a bucket. |
| GetObjectToFile | Downloads an object to the local computer. |

### IsObjectExist/IsBucketExist

The return values of IsObjectExist and IsBucketExist are (bool, error). If the value of the error parameter is nil and the value of bool is true, the object or the bucket exists. If the value of the error parameter is nil and the value of bool is false, the object or the bucket does not exist. If the value of error is not nil, the error message cannot be used to determine whether the object or the bucket exists.

```
func (c *Client) IsObjectExist(ctx context.Context, bucket string, key string, optFns ...func(*IsObjectExistOptions)) (bool, error)
func (c *Client) IsBucketExist(ctx context.Context, bucket string, optFns ...func(*Options)) (bool, error)
```

Example: determine whether an object exists

```
client := oss.NewClient(cfg)

existed, err := client.IsObjectExist(context.TODO(), "examplebucket", "exampleobject")
//existed, err := client.IsObjectExist(context.TODO(), "examplebucket", "exampleobject", func(ioeo *oss.IsObjectExistOptions) {
//	//ioeo.VersionId = oss.Ptr("versionId")
//	//ioeo.RequestPayer = oss.Ptr("requester")
//})

if err != nil {
  // Error
} else {
  fmt.Printf("object existed :%v", existed)
}
```

### PutObjectFromFile

Call the PutObject operation to upload a local file to a bucket. The operation does not support concurrent upload.

```
func (c *Client) PutObjectFromFile(ctx context.Context, request *PutObjectRequest, filePath string, optFns ...func(*Options)) (*PutObjectResult, error)
```

Example

```
client := oss.NewClient(cfg)

result, err := client.PutObjectFromFile(context.TODO(),
  &oss.PutObjectRequest{
    Bucket: oss.Ptr("examplebucket"),
    Key:    oss.Ptr("exampleobject"),
  },
  "/local/dir/example",
)
```

### GetObjectToFile

Call the GetObject operation to download an object in a bucket to the local computer. The operation does not support concurrent downloads.

```
func (c *Client) GetObjectToFile(ctx context.Context, request *GetObjectRequest, filePath string, optFns ...func(*Options)) (*GetObjectResult, error)
```

Example

```
client := oss.NewClient(cfg)

result, err := client.GetObjectToFile(context.TODO(),
  &oss.GetObjectRequest{
    Bucket: oss.Ptr("examplebucket"),
    Key:    oss.Ptr("exampleobject"),
  },
  "/local/dir/example",
)
```

## Comparison between upload and download operations

Various upload and download operations are provided and you can select appropriate operations based on your business scenarios.

**Upload operations**
|Operation name|Description
|:-------|:-------
|Client.PutObject|Performs simple upload to upload a local file of up to 5 GiB.</br>Supports CRC-64 (enabled by default).</br>Supports the progress bar.</br>Supports the request body whose type is io.Reader. If the type of the request body is io.Seeker, when the upload task fails, the local file is reuploaded.
|Client.PutObjectFromFile|Provides the same capability as Client.PutObject.</br>Obtains the request body from the path of the local file.
|Multipart upload operations</br>Client.InitiateMultipartUpload</br>Client.UploadPart</br>Client.CompleteMultipartUpload|Performs multipart upload to upload a local file whose size is up to 48.8 TiB and whose part size is up to 5 GiB.</br>UploadPart supports CRC-64 (enabled by default).</br>UploadPart supports the progress bar.</br>UploadPart supports the request body whose type is io.Reader. If the type of the request body is io.Seeker, when the upload task fails, the local file is reuploaded.
|Uploader.UploadFrom|Encapsulates the simple upload and multipart upload operations and uploads a local file of up to 48.8 TiB.</br>Supports CRC-64 (enabled by default).</br>Supports the progress bar.</br>Supports the request body whose type is io.Reader. If the io.Reader, io.Seeker, and io.ReaderAt types of the request body are supported, data do not need to be cached in memory. Otherwise, data must be cached in memory and then be uploaded.
|Uploader.UploadFile|Provides the same capability as Uploader.UploadFrom.</br>Obtains the request body from the path of the local file.</br>Supports resumable upload.
|Client.AppendObject|Performs append upload to upload a local file of up to 5 GiB.</br>Supports CRC-64 (enabled by default).</br>Supports the progress bar.</br>Supports the request body whose type is io.Reader. If the type of the request body is io.Seeker, when the upload task fails, the local file is reuploaded. The operation is idempotent and data reupload may fail.
|AppendOnlyFile operations</br>AppendOnlyFile.Write</br>AppendOnlyFile.WriteFrom|Provides the same capability as Client.AppendObject.</br>Optimizes fault tolerance after the data reupload fails.

**Download operations**
|Operation name|Description
|:-------|:-------
|Client.GetObject|Performs streaming download. The type of the response body is io.ReadCloser.</br>Does not directly support CRC-64.</br>Does not directly support the progress bar.</br>Does not support reconnection for failed connections during streaming read.
|Client.GetObjectToFile|Downloads an object to the local computer.</br>Supports single connection download.</br>Supports CRC-64 (enabled by default).</br>Supports the progress bar.</br>Supports reconnection for failed connections.
|Downloader.DownloadFile|Performs multipart download to download an object to the local computer.</br>Supports custom part size and concurrent downloads.</br>Supports CRC-64 (enabled by default).</br>Supports the progress bar.</br>Supports reconnection for failed connections.</br>Supports resumable download.</br>Writes a temporary object to OSS and renames the temporary object (enabled by default). You can modify the configurations.
|ReadOnlyFile operations</br>ReadOnlyFile.Read</br>ReadOnlyFile.Seek</br>ReadOnlyFile.Close|File-Like operations that support the io.Reader, io.Seeker, and io.Closer types.</br>Provides the Seek capability.</br>Supports the single-stream mode (default).</br>Supports the asynchronous prefetch mode to improve the read speed.</br>Supports the custom prefetch of blocks and the number of prefetch chunks.</br>Does not directly support CRC-64.</br>Does not directly support the progress bar.</br>Supports reconnection for failed connections.


# Sample scenarios

This section describes how to use OSS SDK for Go in different scenarios.

Scenarios
* [Specify the progress bar](# Specify the progress bar)
* [Data verification](# Data verification)

## Specify the progress bar

In object upload, download, and copy scenarios, you can specify the progress bar to view the transmission progress of an object.

**Supported request parameters for the progress bar**
|Supported request parameter|Method
|:-------|:-------
|PutObjectRequest|PutObjectRequest.ProgressFunc
|GetObjectRequest|GetObjectRequest.ProgressFunc
|CopyObjectRequest|CopyObjectRequest.ProgressFunc
|AppendObjectRequest|AppendObjectRequest.ProgressFunc
|UploadPartRequest|UploadPartRequest.ProgressFunc

**Syntax and parameters for ProgressFunc **
```
type ProgressFunc func(increment, transferred, total int64)
```
| Parameter | Type | Description |
|:-------|:-------|:-------
| increment | int64 | The size of the data transmitted by this callback. Unit: bytes. |
| transferred | int64 | The size of transmitted data. Unit: bytes. |
| total | int64 | The size of the requested data. Unit: bytes. If the value of this parameter is -1, it specifies that the size cannot be obtained. |


Examples

1. Specify the progress bar when you upload a local file by calling PutObject

```
...
client := oss.NewClient(cfg)
client.PutObject(context.TODO(), &oss.PutObjectRequest{
  Bucket: oss.Ptr("bucket"),
  Key:    oss.Ptr("key"),
  ProgressFn: func(increment, transferred, total int64) {
    fmt.Printf("increment:%v, transferred:%v, total:%v\n", increment, transferred, total)
  },
})


```

2. Specify the progress bar when you download an object by calling GetObjectToFile
```
...
client := oss.NewClient(cfg)
client.GetObjectToFile(context.TODO(),
  &oss.GetObjectRequest{
    Bucket: oss.Ptr("bucket"),
    Key:    oss.Ptr("key"),
    ProgressFn: func(increment, transferred, total int64) {
      fmt.Printf("increment:%v, transferred:%v, total:%v\n", increment, transferred, total)
    },
  },
  "/local/dir/example",
)
```

3. Specify the progress bar when you download an object by performing streaming download
```
...
client := oss.NewClient(cfg)

result, err := client.GetObject(context.TODO(), &oss.GetObjectRequest{
  Bucket: oss.Ptr("bucket"),
  Key:    oss.Ptr("key"),
})

if err != nil {
  log.Fatalf("fail to GetObject %v", err)
}

prop := oss.NewProgress(
  func(increment, transferred, total int64) {
    fmt.Printf("increment:%v, transferred:%v, total:%v\n", increment, transferred, total)
  },
  result.ContentLength,
)

io.ReadAll(io.TeeReader(result.Body, prop))
```

## Data verification

OSS provides MD5 verification and CRC-64 to ensure data integrity during requests.

## MD5 verification

When a request is sent to OSS, if the Content-MD5 header is specified, OSS calculates the MD5 hash based on the received content. If the MD5 hash calculated by OSS is different from the MD5 hash configured in the upload request, the InvalidDigest error code is returned. This allows OSS to ensure data integrity for object upload.

Except for PutObject, AppendObject, and UploadPart, the basic API operations automatically calculate the MD5 hash and specify the Content-MD5 header to ensure the integrity of the request.

If you want to use MD5 verification in PutObject, AppendObject, or UploadPart, use the following syntax:

```
...
client := oss.NewClient(cfg)

var body io.Reader

// Calculate the Content-Md5 header. If the request body is not of the io.ReadSeeker type, the data is cached and an MD5 hash is calculated.
calcMd5 := func(input io.Reader) (io.Reader, string, error) {
  if input == nil {
    return input, "1B2M2Y8AsgTpgAmY7PhCfg==", nil
  }
  var (
    r  io.ReadSeeker
    ok bool
  )
  if r, ok = input.(io.ReadSeeker); !ok {
    buf, err := io.ReadAll(input)
    if err != nil {
      return input, "", err
    }
    r = bytes.NewReader(buf)
  }

  curPos, err := r.Seek(0, io.SeekCurrent)
  if err != nil {
    return input, "", err
  }
  h := md5.New()
  _, err = io.Copy(h, r)
  if err != nil {
    return input, "", err
  }
  _, err = r.Seek(curPos, io.SeekStart)
  if err != nil {
    return input, "", err
  }

  return r, base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

body, md5, err := calcMd5(body)

if err != nil {
  log.Fatalf("fail to calcMd5, %v", err)
}

result, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
  Bucket:     oss.Ptr("bucket"),
  Key:        oss.Ptr("key"),
  ContentMD5: oss.Ptr(md5),
  Body:       body,
})

if err != nil {
  log.Fatalf("fail to PutObject, %v", err)
}

fmt.Printf("PutObject result, etg:%v", oss.ToString(result.ETag))
```

## CRC-64

When you upload an object by calling an API operation, such as PutObject, AppendObject, and UploadPart, CRC-64 is enabled by default to ensure data integrity.

When you download an object, take note of the following items:
* If you download an object to a local computer, CRC-64 is enabled to ensure data integrity by default. For example, CRC-64 is enabled for the Downloader.DownloadFile and GetObjectToFile operations.
* If you call streaming read operations such as the GetObject and ReadOnlyFile.Read operations, CRC-64 is not enabled.

If you want to enable CRC-64 in the streaming read operations, you can use the following syntax:

```
...
client := oss.NewClient(cfg)

result, err := client.GetObject(context.TODO(), &oss.GetObjectRequest{
  Bucket: oss.Ptr("bucket"),
  Key:    oss.Ptr("key"),
})

if err != nil {
  log.Fatalf("fail to GetObject, %v", err)
}
defer func() {
  if result.Body != nil {
    result.Body.Close()
  }
}()

var h hash.Hash64
var r io.Reader = result.Body

// The CRC-64 value of the entire object is returned in the response. If you perform range download, CRC-64 is not supported.
// 206 Partial Content indicates that range download is performed.
if result.StatusCode == 200 {
  h = oss.NewCRC64(0)
  r = io.TeeReader(result.Body, h)
}
_, err = io.Copy(io.Discard, r)

if err != nil {
  log.Fatalf("fail to Copy, %v", err)
}

if h != nil && result.HashCRC64 != nil {
  ccrc := fmt.Sprint(h.Sum64())
  scrc := oss.ToString(result.HashCRC64)
  if ccrc != scrc {
    log.Fatalf("crc is inconsistent, client %s, server %s", ccrc, scrc)
  }
}
```

To disable CRC-64, set Config.WithDisableDownloadCRC64Check and Config.WithDisableUploadCRC64Check to true. Example:
```
cfg := oss.LoadDefaultConfig().
  WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
  WithRegion(region).
  WithDisableDownloadCRC64Check(true).
  WithDisableUploadCRC64Check(true)

client := oss.NewClient(cfg)
```


# Migration guide

This section describes how to upgrade OSS SDK for Go from V1 ([aliyun-oss-go-sdk](https://github.com/aliyun/aliyun-oss-go-sdk)) to V2.

## Earliest version for Go

OSS SDK for Go V2 requires that the version for Go must be Go 1.18 or later.

## Import path

OSS SDK for Go V2 uses a new code repository. The code structure is adjusted and organized by functional module. The following table describes the paths and descriptions of these modules.

| Module path | Description |
|:-------|:-------
| github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss | The core of OSS SDK for Go, which is used to call basic and advance API operations. |
| github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials | The access credentials. |
| github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/retry | The retry policies. |
| github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer | The signatures. |
| github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/transport | The HTTP clients. |
| github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/crypto | The client-side encryption configurations. |

Examples

```
// v1
import (
  "github.com/aliyun/aliyun-oss-go-sdk/oss"
)
```

```
// v2
import (
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
  // Import retry, transport, or signer based on your business requirements.
  //"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/xxxx"
)
```

## Configuration loading

OSS SDK for Go V2 simplifies configurations and imports the configurations to [config](oss/config.go). OSS SDK for Go V2 provides auxiliary functions prefixed with With to facilitate the programmatically overwriting of the default configurations.

OSS SDK for Go V2 uses V4 signatures by default. In this case, you must specify the region.

OSS SDK for Go V2 allows you to create an endpoint based on the region information. If you access resources in the public cloud, you do not need to create an endpoint.

Examples

```
// v1
import (
  "github.com/aliyun/aliyun-oss-go-sdk/oss"
)
...

// Obtain access credentials from environment variables.
provider, err := oss.NewEnvironmentVariableCredentialsProvider()

// Set the timeout period of an HTTP connection to 20 and the read or write timeout period of an HTTP connection to 60. Unit: seconds.
time := oss.Timeout(20,60)

// Do not verify SSL certificates.
verifySsl := oss.InsecureSkipVerify(true)

// Specify logs.
logLevel := oss.SetLogLevel(oss.LogInfo)

// Endpoint
endpoint := "oss-cn-hangzhou.aliyuncs.com"

client, err := oss.New(endpoint, "", "", oss.SetCredentialsProvider(&provider), time, verifySsl, logLevel)
```

```
// v2
import (
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

...

// Obtain access credentials from environment variables.
provider := credentials.NewEnvironmentVariableCredentialsProvider()

cfg := oss.LoadDefaultConfig().
  WithCredentialsProvider(provider).
  // Set the timeout period of an HTTP connection to 20. Unit: seconds.
  WithConnectTimeout(20 * time.Second).
  // Set the read or write timeout period of an HTTP connection to 60. Unit: seconds.
  WithReadWriteTimeout(60 * time.Second).
  // Do not verify SSL certificates.
  WithInsecureSkipVerify(true).
  // Specify logs.
  WithLogLevel(oss.LogInfo).
  // Specify the region.
  WithRegion("cn-hangzhou")

client := oss.NewClient(cfg)
```

## Create a client

In OSS SDK for Go V2, the client creation function is changed from New to NewClient. In addition, the client creation function no longer supports the endpoint, ak, and sk parameters.

Examples

```
// v1
client, err := oss.New(endpoint, "ak", "sk")
```

```
// v2
client := oss.NewClient(cfg)
```

## Call API operations

Basic API operations are merged into a single operation method in the \<OperationName\> format, the request parameters of an operation are merged into \<OperationName\>Request, and the response parameters of an operation are merged into \<OperationName\>Result. The operation methods are imported to Client, and context.Context needs to be specified at the same time. Syntax:

```
func (c *Client) <OperationName>(ctx context.Context, request *<OperationName>Request, optFns ...func(*Options)) (*<OperationName>Result,, error)
```

For more information, see [Basic API operations](# Basic operations).

Examples

```
// v1
import "github.com/aliyun/aliyun-oss-go-sdk/oss"

provider, err := oss.NewEnvironmentVariableCredentialsProvider()

client, err := oss.New("yourEndpoint", "", "", oss.SetCredentialsProvider(&provider))  

bucket, err := client.Bucket("examplebucket")

err = bucket.PutObject("exampleobject.txt", bytes.NewReader([]byte("example data")))
```

```
// v2
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

cfg := oss.LoadDefaultConfig().
  WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
  WithRegion("your region")

client := oss.NewClient(cfg)

result, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
  Bucket: oss.Ptr("examplebucket"),
  Key:    oss.Ptr("exampleobject.txt"),
  Body:   bytes.NewReader([]byte("example data")),
})
```

## Generate a pre-signed URL

In OSS SDK for Go V2, the name of the operation used to generate a pre-signed URL is changed from SignURL to Pressign, and the operation is imported to Client. Syntax:

```
func (c *Client) Presign(ctx context.Context, request any, optFns ...func(*PresignOptions)) (*PresignResult, error)
```

The type of request parameters is the same as \<OperationName\>Request in the API operation.

The response contains a pre-signed URL, the HTTP method, the expiration time of the URL, and the signed request headers. Example:
```
type PresignResult struct {
  Method        string
  URL           string
  Expiration    time.Time
  SignedHeaders map[string]string
}
```

For more information, see [Operation used to generate a pre-signed URL](# Operation used to generate a pre-signed URL).

The following sample code provides an example on how to migrate an object from OSS SDK for Go V1 to OSS SDK for Go V2 by generating a pre-signed URL that is used to download the object:

```
// v1
import "github.com/aliyun/aliyun-oss-go-sdk/oss"

provider, err := oss.NewEnvironmentVariableCredentialsProvider()

client, err := oss.New("yourEndpoint", "", "", oss.SetCredentialsProvider(&provider))  

bucket, err := client.Bucket("examplebucket")

signedURL, err := bucket.SignURL("exampleobject.txt", oss.HTTPGet, 60)

fmt.Printf("Sign Url:%s\n", signedURL)
```

```
// v2
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

cfg := oss.LoadDefaultConfig().
  WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
  WithRegion("your region")

client := oss.NewClient(cfg)

result, err := client.Presign(
  context.TODO(),
  &oss.GetObjectRequest{
    Bucket: oss.Ptr("examplebucket"),
    Key:    oss.Ptr("exampleobject.txt"),
  },
  oss.PresignExpires(60*time.Second),
)

fmt.Printf("Sign Method:%v\n", result.Method)
fmt.Printf("Sign Url:%v\n", result.URL)
fmt.Printf("Sign Expiration:%v\n", result.Expiration)
for k, v := range result.SignedHeaders {
  fmt.Printf("SignedHeader %v:%v\n", k, v)
}
```

## Resumable transfer operations

OSS SDK for Go V2 uses data transmission managers, including Uploader, Downloader, and Copier, to manage the upload, download, and copy of objects respectively.  The original resumable transfer operations (Bucket.UploadFile, Bucket.DownloadFile, and Bucket.CopyFile) are removed.

The following table describes the resumable transfer operations.

| Scenario | v2 | v1 |
|:-------|:-------|:-------
| Upload an object | Uploader.UploadFile | Bucket.UploadFile |
| Upload a stream (io.Reader) | Uploader.UploadFrom | Not supported |
| Download an object to a local computer | Downloader.DownloadFile | Bucket.DownloadFile |
| Copy an object | Copier.Copy | Bucket.CopyFile |

Changes to default values

| Scenario | v2 | v1 |
|:-------|:-------|:-------
| Object upload-part size | 6 MiB | Configure part size by specifying parameters. |
| Object upload-default value for concurrency | 3 | 1 |
| Object upload-size threshold | The size of the part. | None |
| Object upload-record upload progress in the checkpoint file | Supported | Supported |
| Object download-part size | 6 MiB | Configure part size by specifying parameters. |
| Object download-default value for concurrency | 3 | 1 |
| Object download-size threshold | The size of the part. | None |
| Object download-record download progress in the checkpoint file | Supported | Supported |
| Object copy-part size | 64 MiB | None |
| Object copy-default value for concurrency | 3 | 1 |
| Object copy-size threshold | 200 MiB | None |
| Object copy-record copy progress in the checkpoint file | Not supported | Supported |

The object upload-size threshold, object download-size threshold, or object copy-size threshold parameters indicate that when the object size is greater than the value of the parameters, multipart upload, multipart download, or multipart copy, respectively, is performed.

For more information about how to use data transmission managers, see [Data transmission managers](# Data transmission managers).

## Client-side encryption

OSS SDK for Go V2 uses EncryptionClient to provide client encryption. OSS SDK for Go V2 also simplifies the API operations used to perform client-side encryption and adopts the same operation naming rules and calling methods as Client.

In addition, OSS SDK for Go V2 only provides reference for how to perform client-side encryption by using an RSA-based, self-managed CMK.

For more information about how to perform client-side encryption by using Key Management Service (KMS), visit [sample/crypto/kms.go](sample/crypto/kms.go).

For more information about client-side encryption, see [Client-side encryption](# Client-side encryption).

The following sample code provides examples on how to use an RSA-based CMK to perform client-side encryption when you upload an object by using OSS SDK for Go V1 and OSS SDK for Go V2:

```
// v1
import "github.com/aliyun/aliyun-oss-go-sdk/oss"
import "github.com/aliyun/aliyun-oss-go-sdk/oss/crypto"

provider, err := oss.NewEnvironmentVariableCredentialsProvider()

client, err := oss.New("yourEndpoint", "", "", oss.SetCredentialsProvider(&provider))  

materialDesc := make(map[string]string)
materialDesc["desc"] = "your master encrypt key material describe information"

masterRsaCipher, err := osscrypto.CreateMasterRsa(materialDesc, "yourRsaPublicKey", "yourRsaPrivateKey")

contentProvider := osscrypto.CreateAesCtrCipher(masterRsaCipher)

cryptoBucket, err := osscrypto.GetCryptoBucket(client, "examplebucket", contentProvider)

err = cryptoBucket.PutObject("exampleobject.txt", bytes.NewReader([]byte("example data")))
```

```
// v2
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/crypto"

cfg := oss.LoadDefaultConfig().
  WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
  WithRegion("your region")

client := oss.NewClient(cfg)

materialDesc := make(map[string]string)
materialDesc["desc"] = "your master encrypt key material describe information"

mc, err := crypto.CreateMasterRsa(materialDesc, "yourRsaPublicKey", "yourRsaPrivateKey")
eclient, err := NewEncryptionClient(client, mc)

result, err := eclient.PutObject(context.TODO(), &PutObjectRequest{
  Bucket: Ptr("examplebucket"),
  Key:    Ptr("exampleobject.txt"),
  Body:   bytes.NewReader([]byte("example data")),
})
```

## Retry

By default, OSS SDK for Go V2 allows retries for HTTP requests. When you upgrade OSS SDK for Go V1 to OSS SDK for Go V2, you must remove the original retry code to avoid increasing the number of retries.
