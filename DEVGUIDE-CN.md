# 开发者指南
## [English](DEVGUIDE.md)

阿里云对象存储（Object Storage Service，简称OSS），是阿里云对外提供的海量、安全、低成本、高可靠的云存储服务。用户可以通过调用API，在任何应用、任何时间、任何地点上传和下载数据，也可以通过用户Web控制台对数据进行简单的管理。OSS适合存放任意文件类型，适合各种网站、开发企业及开发者使用。

该开发套件隐藏了许多较低级别的实现，例如身份验证、请求重试和错误处理, 通过其提供的接口，让您不用复杂编程即可访问阿里云OSS服务。

该开发套件同时提供实用的模块，例如上传和下载管理器，自动将大对象分成多块并行传输。

您可以参阅该指南，来帮助您安装、配置和使用该开发套件。

跳转到:

* [安装](#安装)
* [配置](#配置)
* [接口说明](#接口说明)
* [场景示例](#场景示例)
* [迁移指南](#迁移指南)

# 安装

## 环境准备

使用Go 1.18及以上版本。
请参考[Go安装](https://golang.org/doc/install)下载和安装Go编译运行环境。
您可以执行以下命令查看Go语言版本。
```
go version
```

## 安装SDK

### Go Mod 方式
在 go.mod 文件中添加以下依赖。
```
require (
    github.com/aliyun/alibabacloud-oss-go-sdk-v2 latest
)
```

### 源码方式
```
go get github.com/aliyun/alibabacloud-oss-go-sdk-v2
```
 
## 验证SDK
运行以下代码查看SDK版本：
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

# 配置
您可以配置服务客户端的常用设置，例如超时、日志级别和重试配置，大多数设置都是可选的。
但是，对于每个客户端，您必须指定区域和凭证。 SDK使用这些信息签署请求并将其发送到正确的区域。

此部分的其它主题
* [区域](#区域)
* [凭证](#凭证)
* [访问域名](#访问域名)
* [HTTP客户端](#http客户端)
* [重试](#重试)
* [日志](#日志)
* [配置参数汇总](#配置参数汇总)

## 加载配置
配置客户端的设置有多种方法，以下是推荐的模式。

```
package main

import (
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func main() {
  var (
    // 以华东1（杭州）为例
    region = "cn-hangzhou"

    // 以从环境变量中获取访问凭证为例
    provider credentials.CredentialsProvider = credentials.NewEnvironmentVariableCredentialsProvider()
  )

  cfg := oss.LoadDefaultConfig().
    WithCredentialsProvider(provider).
    WithRegion(region)
  }
```

## 区域
指定区域时，您可以指定向何处发送请求，例如 cn-hangzhou 或 cn-shanghai。有关所支持的区域列表，请参阅 [OSS访问域名和数据中心](https://www.alibabacloud.com/help/zh/oss/user-guide/regions-and-endpoints)。
SDK 没有默认区域，您需要加载配置时使用`config.WithRegion`作为参数显式设置区域。例如
```
cfg := oss.LoadDefaultConfig().WithRegion("cn-hangzhou")
```

>**说明**：该SDK默认使用v4签名，所以必须指定该参数。

## 凭证

SDK需要凭证（访问密钥）来签署对 OSS 的请求, 所以您需要显式指定这些信息。当前支持凭证配置如下：
* [环境变量](#环境变量)
* [ECS实例角色](#ecs实例角色)
* [静态凭证](#静态凭证)
* [外部进程](#外部进程)
* [RAM角色](#ram角色)
* [OIDC角色SSO](#oidc角色sso)
* [自定义凭证提供者](#自定义凭证提供者)

### 环境变量

SDK 支持从环境变量获取凭证，支持的环境变量名如下：
* OSS_ACCESS_KEY_ID
* OSS_ACCESS_KEY_SECRET
* OSS_SESSION_TOKEN（可选）

以下展示了如何配置环境变量。

1. Linux、OS X 或 Unix
```
$ export OSS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
$ export OSS_ACCESS_KEY_SECRET=YOUR_ACCESS_KEY_SECRET
$ export OSS_SESSION_TOKEN=TOKEN
```

2. Windows
```
$ set OSS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
$ set OSS_ACCESS_KEY_SECRET=YOUR_ACCESS_KEY_SECRET
$ set OSS_SESSION_TOKEN=TOKEN
```

使用环境变量凭证

```
provider := credentials.NewEnvironmentVariableCredentialsProvider()
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

### ECS实例角色

如果你需要在阿里云的云服务器ECS中访问您的OSS，您可以通过ECS实例RAM角色的方式访问OSS。实例RAM角色允许您将一个角色关联到云服务器实例，在实例内部基于STS临时凭证通过指定方法访问OSS。

使用ECS实例角色凭证

1. 指定实例角色，例如角色名为 EcsRoleExample
```
provider := credentials.NewEcsRoleCredentialsProvider(func(ercpo *credentials.EcsRoleCredentialsProviderOptions) {
	ercpo.RamRole = "EcsRoleExample"
})
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```
   
2. 不指定实例角色
```
provider := credentials.NewEcsRoleCredentialsProvider()
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```
当不指定实例角色名时，会自动查询角色名。

### 静态凭证

您可以在应用程序中对凭据进行硬编码，显式设置要使用的访问密钥。

> **注意:** 请勿将凭据嵌入应用程序中，此方法仅用于测试目的。

1. 长期凭证
```
provider := credentials.NewStaticCredentialsProvider("AKId", "AKSecrect")
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

2. 临时凭证
```
provider := credentials.NewStaticCredentialsProvider("AKId", "AKSecrect", "Token")
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

### 外部进程

您可以在应用程序中，通过外部进程获取凭证。
> **注意:**
> </br>生成凭证的命令不可由未经批准的进程或用户访问，则可能存在安全风险。
> </br>生成凭证的命令不会把任何秘密信息写入 stderr 或 stdout，因为该信息可能会被捕获或记录，可能会将其向未经授权的用户公开。

外部命令返回的凭证，支持长期凭证和临时凭证，其格式如下：
1. 长期凭证
```
{
  "AccessKeyId" : "AKId",
  "AccessKeySecret" : "AKSecrect",
}
```

2. 临时凭证
```
{
  "AccessKeyId" : "AKId",
  "AccessKeySecret" : "AKSecrect",
  "Expiration" : "2023-12-29T07:45:02Z",
  "SecurityToken" : "token",
}
```

以 test-command 命令为例，该命令返回长期凭证

```
process := "test-command"
provider := credentials.NewProcessCredentialsProvider(process)
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

以 test-command-sts 命令为例，该命令返回临时凭证，每次请求凭证都不一样

```
process := "test-command-sts"
cprovider := credentials.NewProcessCredentialsProvider(process)
// NewCredentialsFetcherProvider 根据 'Expiration' 时间，自动刷新凭证
provider := credentials.NewCredentialsFetcherProvider(credentials.CredentialsFetcherFunc(func(ctx context.Context) (credentials.Credentials, error) {
  return cprovider.GetCredentials(ctx)
}))
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

### RAM角色

如果您需要授权访问或跨账号访问OSS，您可以通过RAM用户扮演对应RAM角色的方式授权访问或跨账号访问OSS。

SDK 不直接提供该访问凭证实现，需要结合阿里云凭证库[credentials-go](https://github.com/aliyun/credentials-go)，具体配置如下:

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

### OIDC角色SSO

您也可以在应用或服务中使用OIDC认证访问OSS服务，关于OIDC角色SSO的更多信息，请参见[OIDC角色SSO概览](https://www.alibabacloud.com/help/zh/ram/user-guide/overview-of-oidc-based-sso)。

SDK 不直接提供该访问凭证实现，需要结合阿里云凭证库[credentials-go](https://github.com/aliyun/credentials-go)，具体配置如下:

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

### 自定义凭证提供者

当以上凭证配置方式不满足要求时，您可以自定义获取凭证的方式。SDK 支持多种实现方式。

1. 实现 credentials.CredentialsProvider 接口
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
  // 返回长期凭证
  return credentials.Credentials{AccessKeyID: "id", AccessKeySecret: "secret"}, nil
  // 返回临时凭证
  //return credentials.Credentials{AccessKeyID: "id", AccessKeySecret: "secret",	SecurityToken: "token"}, nil
}

provider := NewCustomerCredentialsProvider()
cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)

```

2. 通过 credentials.CredentialsProviderFunc

credentials.CredentialsProviderFunc 是 credentials.CredentialsProvider 的 易用性封装。

```
provider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
  // 返回长期凭证
  return credentials.Credentials{AccessKeyID: "id", AccessKeySecret: "secret"}, nil
  // 返回临时凭证
  //return credentials.Credentials{AccessKeyID: "id", AccessKeySecret: "secret",	SecurityToken: "token"}, nil
})

cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)
```

3. 通过 credentials.CredentialsFetcherFunc

credentials.CredentialsFetcherFunc 是 credentials.CredentialsFetcher 易用性接口。

credentials.CredentialsFetcher 具备 根据 'Expiration' 时间，自动刷新凭证的能力，当您需要定期更新凭证时，请使用该方式。

```
customerProvider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
  var (
    akId     string
    akSecret string
    token    string
    expires  *time.Time
  )

  // 获取临时凭证 和 凭证的过期时间
  ...

  // 遇到错误
  if err != nil {
    return credentials.Credentials{}, err
  }

  // 成功
  return credentials.Credentials{AccessKeyID: akId, AccessKeySecret: akSecret,	SecurityToken: token, Expires: expires}, nil
})

provider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
  return customerProvider.GetCredentials()
})

cfg := oss.LoadDefaultConfig().WithCredentialsProvider(provider)

```

## 访问域名

您可以通过Endpoint参数，自定义服务请求的访问域名。

当不指定时，SDK根据Region信息，构造公网访问域名。例如当Region为'cn-hangzhou'时，构造出来的访问域名为'oss-cn-hangzhou.aliyuncs.com'。

您可以通过修改配置参数，构造出其它访问域名，例如 内网访问域名，传输加速访问域名 和 双栈(IPV6,IPV4)访问域名。有关OSS访问域名规则，请参考[OSS访问域名使用规则](https://www.alibabacloud.com/help/zh/oss/user-guide/oss-domain-names)。

当通过自定义域名访问OSS服务时，您需要指定该配置参数。在使用自定义域名发送请求时，请先绑定自定域名至Bucket默认域名，具体操作详见 [绑定自定义域名](https://www.alibabacloud.com/help/zh/oss/user-guide/map-custom-domain-names-5)。


### 使用标准域名访问

以 访问 Region 'cn-hangzhou' 为例

1. 使用公网域名

```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou")

或者

cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("oss-cn-hanghzou.aliyuncs.com")
```

2. 使用内网域名

```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithUseInternalEndpoint(true)

或者

cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("oss-cn-hanghzou-internal.aliyuncs.com")
```
   
3. 使用传输加速域名
```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithUseAccelerateEndpoint(true)

或者

cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("oss-accelerate.aliyuncs.com")
```   
   
4. 使用双栈域名
```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithUseDualStackEndpoint(true)

或者

cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("cn-hangzhou.oss.aliyuncs.com")
```   

### 使用自定义域名访问

以 'www.example-***.com' 域名 绑定到 'cn-hangzhou' 区域 的 bucket-example 存储空间为例

```
cfg := oss.LoadDefaultConfig().
  WithRegion("cn-hangzhou").
  WithEndpoint("www.example-***.com").
  WithUseCName(true)
```

### 访问专有云或专有域

```
var (
  region = "YOUR Region"
  endpoint = "YOUR Endpoint"
)

cfg := oss.LoadDefaultConfig().
  WithRegion(region).
  WithEndpoint(endpoint)
```

## HTTP客户端

在大多数情况下，使用具有默认值的默认HTTP客户端 能够满足业务需求。您也可以更改HTTP 客户端，或者更改 HTTP 客户端的默认配置，以满足特定环境下的使用需求。

本部分将介绍如何设置 和 创建 HTTP 客户端。

### 设置HTTP客户端常用配置

通过config修改常用的配置，支持参数如下：

|参数名字 | 说明 | 示例 
|:-------|:-------|:-------
|ConnectTimeout|建立连接的超时时间, 默认值为 5 秒|WithConnectTimeout(10 * time.Second)
|ReadWriteTimeout|应用读写数据的超时时间, 默认值为 10 秒|WithReadWriteTimeout(30 * time.Second)
|InsecureSkipVerify|是否跳过SSL证书校验，默认检查SSL证书|WithInsecureSkipVerify(true)
|EnabledRedirect|是否开启HTTP重定向, 默认不开启|WithEnabledRedirect(true)
|ProxyHost|设置代理服务器|WithProxyHost("http://user:passswd@proxy.example-***.com")
|ProxyFromEnvironment|通过环境变量设置代理服务器|WithProxyFromEnvironment(true)
|UploadBandwidthlimit|整体的上传带宽限制，单位为 KiB/s|WithUploadBandwidthlimit(10*1024)
|DownloadBandwidthlimit|整体的下载带宽限制，单位为 KiB/s|WithDownloadBandwidthlimit(10*1024)

示例

```
cfg := oss.LoadDefaultConfig().
  WithConnectTimeout(10 * time.Second).
  WithUploadBandwidthlimit(10*1024)
```

### 自定义HTTP客户端

当常用配置参数无法满足场景需求时，可以使用 WithHTTPClient 替换默认的 HTTP 客户端。

在以下示例未提到的设置参数，请参考 [Transport](https://pkg.go.dev/net/http#Transport) 文档。

```
import (
  "crypto/tls"
  "net/http"
  "time"

  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
  "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/transport"
)

...

// 常用超时或其它设置
transConfig := transport.Config{
  // 连接超时, 默认值 5秒
  //ConnectTimeout: oss.Ptr(10 * time.Second),

  // 应用读写数据的超时时间, 默认值 10秒
  //ReadWriteTimeout: oss.Ptr(20 * time.Second),

  // 连接空闲超时时间, 默认值 50 秒
  //IdleConnectionTimeout: oss.Ptr(40 * time.Second),

  // 网络连接的保持期限, 默认值 30 秒
  //KeepAliveTimeout: oss.Ptr(40 * time.Second),

  // 是否打开启HTTP重定向，默认不启用
  //EnabledRedirect: oss.Ptr(true),
}

// http.Transport 设置
var transports []func(*http.Transport)

// 最大连接数，默认值 100
//transports = append(transports, transport.MaxConnections(200))

// 如果请求有“Expect: 100-Continue”标头，则此设置表示完全写入请求标头后等待服务器第一个响应标头的最长时间，默认 1秒
//transports = append(transports, transport.ExpectContinueTimeout(2*time.Second))

// TLS的最低版本，默认值 TLS 1.2
//transports = append(transports, transport.TLSMinVersion(tls.VersionTLS13))

// 是否跳过证书检查，默认不跳过
//transports = append(transports, transport.InsecureSkipVerify(true))

// 其它 http.Transport 参数设置
//transports = append(transports, func(t *http.Transport) {
//  t.DisableCompression
//})

customClient := transport.NewHttpClient(&transConfig, transports...)

cfg := oss.LoadDefaultConfig().WithHttpClient(customClient)
```

## 重试

您可以配置对HTTP请求的重试行为。

### 默认重试策略

当没有配置重试策略时，SDK 使用 retry.Standard 作为客户端的默认实现，其默认配置如下：

|参数名称 | 说明 | 默认值 
|:-------|:-------|:-------
|MaxAttempts|最大尝试次数| 3
|MaxBackoff|最大退避时间| 20秒, 20 * time.Second
|BaseDelay|基础延迟| 200毫秒, 200 * time.Millisecond
|Backoff|退避算法| FullJitter 退避,  [0.0, 1.0) * min(2 ^ attempts * baseDealy, maxBackoff)
|ErrorRetryables|可重试的错误| 具体的错误信息，请参见[重试错误](oss/retry/retryable_error.go)

当发生可重试错误时，将使用其提供的配置来延迟并随后重试该请求。请求的总体延迟会随着重试次数而增加，如果默认配置不满足您的场景需求时，需要配置重试参数 或者修改重试实现。

### 调整最大尝试次数

您可以通过以下两种方式修改最大尝试次数。例如 最多尝试 5  次 

```
cfg := oss.LoadDefaultConfig().WithRetryMaxAttempts(5)

或者

cfg := oss.LoadDefaultConfig().WithRetryer(retry.NewStandard(func(ro *retry.RetryOptions) {
  ro.MaxAttempts = 5
}))
```

### 调整退避延迟

例如 调整 BaseDelay 为 500毫秒，最大退避时间为 25秒

```
cfg := oss.LoadDefaultConfig().WithRetryer(retry.NewStandard(func(ro *retry.RetryOptions) {
  ro.MaxBackoff = 25 * time.Second
  ro.BaseDelay = 500 * time.Millisecond
}))
```

### 调整退避算法

例如 使用固定时间退避算法，每次延迟2秒 

```
cfg := oss.LoadDefaultConfig().WithRetryer(retry.NewStandard(func(ro *retry.RetryOptions) {
  ro.Backoff = retry.NewFixedDelayBackoff(2 * time.Second)
}))
```

### 调整重试错误

例如 在原有基础上，新增自定义可重试错误

```
type CustomErrorCodeRetryable struct {
}

func (*CustomErrorCodeRetryable) IsErrorRetryable(err error) bool {
  // 判断错误
  // return true
  return false
}

errorRetryables := retry.DefaultErrorRetryables
errorRetryables = append(errorRetryables, &CustomErrorCodeRetryable{})

cfg := oss.LoadDefaultConfig().WithRetryer(retry.NewStandard(func(ro *retry.RetryOptions) {
  ro.ErrorRetryables = errorRetryables
}))
```

### 禁用重试

当您希望禁用所有重试尝试时，可以使用 retry.NopRetryer 实现
```
cfg := oss.LoadDefaultConfig().WithRetryer(&retry.NopRetryer{})
```


## 日志

为了方便追查问题，SDK提供了日志记录功能，您可以在应用程序中启用调试信息以调试和诊断请求问题。

当需要启用日志记录功能时，您需要配置日志级别。当不设置日志接口时，默认将日志信息发送到进程的标准输出(stdout).

日志级别：oss.LogError, oss.LogWarn, oss.LogInfo, oss.LogDebug

日志接口: oss.LogPrinter, oss.LogPrinterFunc

例如，开启日志功能，设置日志级别为 Info，输出到标准错误输出(stderr)

```
cfg := oss.LoadDefaultConfig().
  WithLogLevel(oss.LogInfo).
  WithLogPrinter(oss.LogPrinterFunc(func(a ...any) {
    fmt.Fprint(os.Stderr, a...)
  }))
```

## 配置参数汇总

支持的配置参数：

|参数名字 | 说明 | 示例 
|:-------|:-------|:-------
|Region|(必选)请求发送的区域, 必选|WithRegion("cn-hangzhou")
|CredentialsProvider|(必选)设置访问凭证|WithCredentialsProvider(provider)
|Endpoint|访问域名|WithEndpoint("oss-cn-hanghzou.aliyuncs.com")
|HttpClient|HTTP客户都端|WithHttpClient(customClient)
|RetryMaxAttempts|HTTP请求时的最大尝试次数, 默认值为 3|WithRetryMaxAttempts(5)
|Retryer|HTTP请求时的重试实现|WithRetryer(customRetryer)
|ConnectTimeout|建立连接的超时时间, 默认值为 5 秒|WithConnectTimeout(10 * time.Second)
|ReadWriteTimeout|应用读写数据的超时时间, 默认值为 10 秒|WithReadWriteTimeout(30 * time.Second)
|InsecureSkipVerify|是否跳过SSL证书校验，默认检查SSL证书|WithInsecureSkipVerify(true)
|EnabledRedirect|是否开启HTTP重定向, 默认不开启|WithEnabledRedirect(true)
|ProxyHost|设置代理服务器|WithProxyHost("http://user:passswd@proxy.example-***.com")
|ProxyFromEnvironment|通过环境变量设置代理服务器|WithProxyFromEnvironment(true)
|UploadBandwidthlimit|整体的上传带宽限制，单位为 KiB/s|WithUploadBandwidthlimit(10*1024)
|DownloadBandwidthlimit|整体的下载带宽限制，单位为 KiB/s|WithDownloadBandwidthlimit(10*1024)
|SignatureVersion|签名版本，默认值为v4|WithSignatureVersion(oss.SignatureVersionV1)
|LogLevel|设置日志级别|WithLogLevel(oss.LogInfo)
|LogPrinter|设置日志打印接口|WithLogPrinter(customPrinter)
|DisableSSL|不使用https请求，默认使用https|WithDisableSSL(true)
|UsePathStyle|使用路径请求风格，即二级域名请求风格，默认为bucket托管域名|WithUsePathStyle(true)
|UseCName|是否使用自定义域名访问，默认不使用|WithUseCName(true)
|UseDualStackEndpoint|是否使用双栈域名访问，默认不使用|WithUseDualStackEndpoint(true)
|UseAccelerateEndpoint|是否使用传输加速域名访问，默认不使用|WithUseAccelerateEndpoint(true)
|UseInternalEndpoint|是否使用内网域名访问，默认不使用|WithUseInternalEndpoint(true)
|DisableUploadCRC64Check|上传时关闭CRC64校验，默认开启CRC64校验|WithDisableUploadCRC64Check(true)
|DisableDownloadCRC64Check|下载时关闭CRC64校验，默认开启CRC64校验|WithDisableDownloadCRC64Check(true)
|AdditionalHeaders|指定额外的签名请求头，V4签名下有效|WithAdditionalHeaders([]string{"content-length"})
|UserAgent|指定额外的User-Agent信息|WithUserAgent("user identifier")

# 接口说明

本部分介绍SDK提供的接口, 以及如何使用这些接口。

此部分的其它主题
* [基础接口](#基础接口)
* [预签名接口](#预签名接口)
* [分页器](#分页器)
* [传输管理器](#传输管理器)
* [类文件(File-Like)](#类文件file-like)
* [客户端加密](#客户端加密)
* [其它接口](#其它接口)
* [上传下载接口对比](#上传下载接口对比)

## 基础接口

SDK 提供了 与 REST API 对应的接口，把这类接口叫做 基础接口 或者 低级别API。您可以通过这些接口访问OSS的服务，例如创建存储空间，更新和删除存储空间的配置等。

这些接口采用了相同的命名规则，其接口定义如下：

```
func (c *Client) <OperationName>(ctx context.Context, request *<OperationName>Request, optFns ...func(*Options)) (result *<OperationName>Result, err error)
```

**参数列表**：
|参数名|类型|说明
|:-------|:-------|:-------
|ctx|context.Context|请求的上下文，可以用来设置请求的总时限
|request|*\<OperationName\>Request|设置具体接口的请求参数，例如bucket，key
|optFns|...func(*Options)|(可选)接口级的配置参数, 例如修改此次调用接口时读写超时

**返回值列表**：
|返回值名|类型|说明
|:-------|:-------|:-------
|result|*\<OperationName\>Result|接口返回值，当 err 为nil 时有效
|err|error|请求的状态，当请求失败时，err 不为 nil

**指针参数**:
<br/>'\<OperationName\>Request' 类型里，输入参数需要通过指针方式传递给接口，'\<OperationName\>Result' 类型里，返回参数需要通过指针方式返回给调用者。例如 ListObjctsRequest 和 ListObjectsResult：

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
为了方便接口调用，SDK 提供了辅助函数'oss.Ptr'完成非指针类型到指针类型的转换，同时也提供辅助函数'oss.To\<Type\>' 安全地从指针类型转换成非指针类型。
例如 oss.Ptr函数把 string 转成 *string, 相反, oss.ToString 从 *string 转成 string。

示例：
1. 创建存储空间

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

2. 拷贝对象, 同时设置接口级的读写超时

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

更多的示例，请参考 sample 目录

## 预签名接口

您可以使用预签名接口生成预签名URL，授予对存储空间中对象的限时访问权限，或者允许他人将特定对象的上传到存储空间。在过期时间之前，您可以多次使用预签名URL。

预签名接口定义如下：
```
func (c *Client) Presign(ctx context.Context, request any, optFns ...func(*PresignOptions)) (result *PresignResult, err error)
```

**参数列表**：
|参数名|类型|说明
|:-------|:-------|:-------
|ctx|context.Context|请求的上下文
|request|any|设置需要生成签名URL的接口名，和 '\<OperationName\>Request' 一致
|optFns|...func(*PresignOptions)|(可选)，设置过期时间，如果不指定，默认有效期为15分钟

**返回值列表**：
|返回值名|类型|说明
|:-------|:-------|:-------
|result|*PresignResult|返回结果，包含 预签名URL，HTTP 方法，过期时间 和 参与签名的请求头
|err|error|请求的状态，当请求失败时，err 不为 nil

**request参数支持的类型**：
|类型|对应的接口
|:-------|:-------
|*GetObjectRequest|GetObject
|*PutObjectRequest|PutObject
|*HeadObjectRequest|HeadObject
|*InitiateMultipartUploadRequest|InitiateMultipartUpload
|*UploadPartRequest|UploadPart
|*CompleteMultipartUploadRequest|CompleteMultipartUpload
|*AbortMultipartUploadRequest|AbortMultipartUpload

**PresignOptions选项**
|选项值|类型|说明
|:-------|:-------|:-------
|Expires|time.Duration|从当前时间开始，多长时间过期。例如 设置一个有效期为30分钟，30 * time.Minute
|Expiration|time.Time|绝对过期时间

> **注意:** 在签名版本4下，有效期最长为7天。同时设置 Expiration 和 Expires时，优先取 Expiration。

**PresignResult返回值**：
|参数名|类型|说明
|:-------|:-------|:-------
|Method|string|HTTP 方法，和 接口对应，例如GetObject接口，返回 GET
|URL|string|预签名 URL
|Expiration|time.Time| 签名URL的过期时间
|SignedHeaders|map[string]string|被签名的请求头，例如PutObject接口，设置了Content-Type 时，会返回 Content-Type 的信息。


示例
1. 为对象生成预签名 URL，然后下载对象（GET 请求）
```
client := oss.NewClient(cfg)

result, err := client.Presign(context.TODO(), &oss.GetObjectRequest{
  Bucket: oss.Ptr("bucket"),
  Key:    oss.Ptr("key"),
})

resp, err := http.Get(result.URL)
```

2. 为上传生成预签名 URL, 设置自定义元数据，有效期为10分钟，然后上传文件（PUT 请求）
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

更多的示例，请参考 sample 目录

## 分页器

对于列举类接口，当响应结果太大而无法在单个响应中返回时，都会返回分页结果，该结果同时包含一个用于检索下一页结果的标记。当需要获取下一页结果时，您需要在发送请求时设置该标记。

对常用的列举接口，V2 SDK 提供了分页器（Paginator），支持自动分页，当进行多次调用时，自动为您获取下一页结果。使用分页器时，您只需要编写处理结果的代码。

分页器 包含了 分页器对象 '\<OperationName\>Paginator' 和 分页器创建方法 'New\<OperationName\>Paginator'。分页器创建方法返回一个分页器对象，该对象实现了 'HasNext' 和 'NextPage' 方法，分别用于判断是否还有更多页, 并调用操作来获取下一页。

分页器创建方法 'New\<OperationName\>Paginator' 里的 request 参数类型 与 '\<OperationName\>' 接口中的 reqeust 参数类型一致。

'\<OperationName\>Paginator.NextPage' 返回的结果类型 和 '\<OperationName\>' 接口 返回的结果类型 一致。

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

支持的分页器对象如下：
|分页器对象|创建方法|对应的列举接口
|:-------|:-------|:-------
|ListObjectsPaginator|NewListObjectsPaginator|ListObjects, 列举存储空间中的对象信息
|ListObjectsV2Paginator|NewListObjectsV2Paginator|ListObjectsV2, 列举存储空间中的对象信息
|ListObjectVersionsPaginator|NewListObjectVersionsPaginator|ListObjectVersions, 列举存储空间中的对象版本信息
|ListBucketsPaginator|NewListBucketsPaginator|ListBuckets, 列举存储空间
|ListPartsPaginator|NewListPartsPaginator|ListParts, 列举指定Upload ID所属的所有已经上传成功分片
|ListMultipartUploadsPaginator|NewListMultipartUploadsPaginator|ListMultipartUploads, 列举存储空间中的执行中的分片上传事件

PaginatorOptions 选项说明：
|参数|说明
|:-------|:-------
|Limit|指定返回结果的最大数


以 ListObjects 为例，分页器遍历所有对象 和 手动分页遍历所有对象 对比

```
// 分页器遍历所有对象
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
// 手动分页遍历所有对象
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

## 传输管理器

针对大文件的传输场景，新增了 'Uploader'，'Downloader' 和 'Copier' 模块，分别管理对象的 上传，下载 和 拷贝。

### 上传管理器(Uploader)

上传管理器 利用分片上传接口，把大文件或者流分成多个较小的分片并发上传，提升上传的性能。
</br>针对文件的上传场景，还提供了断点续传的能力，即在上传过程中，记录已完成的分片状态，如果出现网络中断、程序异常退出等问题导致文件上传失败，甚至重试多次仍无法完成上传，再次上传时，可以通过断点记录文件恢复上传。

```
type Uploader struct {
  ...
}

func (c *Client) NewUploader(optFns ...func(*UploaderOptions)) *Uploader 

func (u *Uploader) UploadFrom(ctx context.Context, request *PutObjectRequest, body io.Reader, optFns ...func(*UploaderOptions)) (*UploadResult, error)

func (u *Uploader) UploadFile(ctx context.Context, request *PutObjectRequest, filePath string, optFns ...func(*UploaderOptions)) (*UploadResult, error)
```

**参数列表**：
|参数名|类型|说明
|:-------|:-------|:-------
|ctx|context.Context|请求的上下文
|request|*PutObjectRequest|上传对象的请求参数，和 PutObject 接口的 请求参数一致
|body|io.Reader|需要上传的流。当 body 只支持io.Reader类型，必须先把数据缓冲在内存中，然后才能上传该部分。当 body 同时支持 io.Reader, io.Seeker 和 io.ReaderAt 类型时，不需要把数据缓存在内存里。
|filePath|string|本地文件路径
|optFns|...func(*UploaderOptions)|(可选)，配置选项


**UploaderOptions选项说明：**
|参数|类型|说明
|:-------|:-------|:-------
|PartSize|int64|指定分片大小，默认值为 6MiB
|ParallelNum|int|指定上传任务的并发数，默认值为 3。针对的是单次调用的并发限制，而不是全局的并发限制
|LeavePartsOnError|bool|当上传失败时，是否保留已上传的分片，默认不保留 
|EnableCheckpoint|bool|是否记录断点上传信息，默认不记录
|CheckpointDir|string|指定记录文件的保存路径，例如 /local/dir/, 当EnableCheckpoint 为 true时有效


当使用NewUploader实例化实例时，您可以指定多个配置选项来自定义对象的上传行为。也可以在每次调用上传接口时，指定多个配置选项来自定义每次上传对象的行为。

设置Uploader的配置参数
```
u := client.NewUploader(func(uo *oss.UploaderOptions) {
  uo.PartSize = 10 * 1024 * 1024
})
```

设置每次上传请求的配置参数
```
request := &oss.PutObjectRequest{Bucket: oss.Ptr("bucket"), Key: oss.Ptr("key")}
result, err := u.UploadFile(context.TODO(), request, "/local/dir/example", func(uo *oss.UploaderOptions) {
  uo.PartSize = 10 * 1024 * 1024
})
```

示例

1. 使用 Uploader上传流

```
...
client := oss.NewClient(cfg)

u := client.NewUploader()

var r io.Reader
// TODO 绑定io.Reader 实例到 r

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

2. 使用 Uploader上传文件

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

3. 上传文件，并开启断点续传功能
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


### 下载管理器(Downloader)

下载管理器 利用范围下载，把大文件分成多个较小的分片并发下载，提升下载的性能。
</br>该接口提供了断点续传的能力，即在下载过程中，记录已完成的分片状态，如果出现网络中断、程序异常退出等问题导致文件下载失败，甚至重试多次仍无法完成下载，再次下载时，可以通过断点记录文件恢复下载。

```
type Downloader struct {
  ...
}

func (c *Client) NewDownloader(optFns ...func(*DownloaderOptions)) *Downloader

func (d *Downloader) DownloadFile(ctx context.Context, request *GetObjectRequest, filePath string, optFns ...func(*DownloaderOptions)) (result *DownloadResult, err error)
```

**参数列表**：
|参数名|类型|说明
|:-------|:-------|:-------
|ctx|context.Context|请求的上下文
|request|*GetObjectRequest|下载对象的请求参数，和 GetObject 接口的 请求参数一致
|filePath|string|本地文件路径
|optFns|...func(*DownloaderOptions)|(可选)，配置选项


**DownloaderOptions选项说明：**
|参数|类型|说明
|:-------|:-------|:-------
|PartSize|int64|指定分片大小，默认值为 6MiB
|ParallelNum|int|指定上传任务的并发数，默认值为 3。针对的是单次调用的并发限制，而不是全局的并发限制
|EnableCheckpoint|bool|是否记录断点下载信息，默认不记录
|CheckpointDir|string|指定记录文件的保存路径，例如 /local/dir/, 当EnableCheckpoint 为 true时有效
|VerifyData|bool|恢复下载时，是否要校验已下载数据的CRC64值，默认不校验, 当EnableCheckpoint 为 true时有效
|UseTempFile |bool|下载文件时，是否使用临时文件，默认使用。先下载到 临时文件上，当成功后，再重命名为目标文件


当使用NewDownloader实例化实例时，您可以指定多个配置选项来自定义对象的下载行为。也可以在每次调用下载接口时，指定多个配置选项来自定义每次下载对象的行为。

设置Downloader的配置参数
```
d := client.NewDownloader(func(do *oss.DownloaderOptions) {
  do.PartSize = 10 * 1024 * 1024
})
```

设置每次下载请求的配置参数
```
request := &oss.GetObjectRequest{Bucket: oss.Ptr("bucket"), Key: oss.Ptr("key")}
d.DownloadFile(context.TODO(), request, "/local/dir/example", func(do *oss.DownloaderOptions) {
  do.PartSize = 10 * 1024 * 1024
})
```

示例

1. 使用 Downloader 下载到本地文件

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

### 拷贝管理器(Copier)
当需要将对象从存储空间复制到另外一个存储空间，或者修改对象的属性时，您可以通过拷贝接口 或者分片拷贝接口来完成这个操作。
</br>这两个接口有其适用的场景，例如：
* 拷贝接口(CopyObject) 只适合拷贝 5GiB 以下的对象；
* 分片拷贝接口(UploadPartCopy) 不支持 元数据指令(x-oss-metadata-directive) 和 标签指令(x-oss-tagging-directive) 参数, 
拷贝时，您需要主动去设置需要复制的元数据和标签。
* 服务端优化了拷贝(CopyObject)接口，使其具备浅拷贝的能力，在特定的场景下也支持拷贝大文件。

拷贝管理器提供了通用的拷贝接口，隐藏了接口的差异和实现细节，根据拷贝的请求参数，自动选择合适的接口复制对象。

```
type Copier struct {
  ...
}

func (c *Client) NewCopier(optFns ...func(*CopierOptions)) *Copier

func (c *Copier) Copy(ctx context.Context, request *CopyObjectRequest, optFns ...func(*CopierOptions)) (*CopyResult, error)
```

**参数列表**：
|参数名|类型|说明
|:-------|:-------|:-------
|ctx|context.Context|请求的上下文
|request|*CopyObjectRequest|拷贝对象的请求参数，和 CopyObject 接口的 请求参数一致
|optFns|...func(*CopierOptions)|(可选)，配置选项


**CopierOptions选项说明：**
|参数|类型|说明
|:-------|:-------|:-------
|PartSize|int64|指定分片大小，默认值为 64MiB
|ParallelNum|int|指定上传任务的并发数，默认值为 3。针对的是单次调用的并发限制，而不是全局的并发限制
|MultipartCopyThreshold|int64|使用分片拷贝的阈值，默认值为 200MiB
|LeavePartsOnError|bool|当拷贝失败时，是否保留已拷贝的分片，默认不保留 
|DisableShallowCopy|bool|不使用浅拷贝行为，默认使用


当使用NewCopier实例化实例时，您可以指定多个配置选项来自定义对象的下载行为。也可以在每次调用下载接口时，指定多个配置选项来自定义每次下载对象的行为。

设置Copier的配置参数
```
d := client.NewCopier(func(co *oss.CopierOptions) {
  co.PartSize = 100 * 1024 * 1024
})
```

设置每次拷贝请求的配置参数
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

> **注意:**
> </br>拷贝对象时，CopyObjectRequest.MetadataDirective 决定了对象元数据的拷贝行为，默认 复制 源对象标签
> </br>拷贝对象时，CopyObjectRequest.TaggingDirective 决定了对象标签的拷贝行为，默认 复制 源对象标签 


示例

1. 拷贝文件，默认会复制 元数据 和 标签
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

2. 拷贝文件，只拷贝数据，不拷贝元数据和标签
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

3. 修改 对象的存储类型 为标准类型

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

## 类文件(File-Like)

新增了File-Like接口，提供了模仿文件的读写行为来操作存储空间里的对象。

支持以下两种方式：
* 只读(ReadOnlyFile)
* 追加写(AppendOnlyFile)

### 只读文件(ReadOnlyFile)

以只读方式访问存储空间的对象。在只读方式上，提供了 单流 和 并发+预取 两种模式，您可以根据场景需要，调整并发数，以提升读的速度。同时，内部实现了连接断掉重连的机制，在一些比较复杂的网络环境下，具备更好的鲁棒性。

```
type ReadOnlyFile struct {
...
}

func (c *Client) OpenFile(ctx context.Context, bucket string, key string, optFns ...func(*OpenOptions)) (file *ReadOnlyFile, err error)
```

**参数列表**：
|参数名|类型|说明
|:-------|:-------|:-------
|ctx|context.Context|请求的上下文
|bucket|string|设置存储空间名字
|key|string|设置对象名
|optFns|...func(*OpenOptions)|(可选)，打开文件时的配置选项

**返回值列表**：
|返回值名|类型|说明
|:-------|:-------|:-------
|file|*ReadOnlyFile|只读文件的实例，当 err 为nil 时有效
|err|error|打开只读文件的状态，当失败时，err 不为 nil

**OpenOptions选项说明：**
|参数|类型|说明
|:-------|:-------|:-------
|Offset|int64|打开文件时的初始偏移量，默认值是0
|VersionId|*string|指定对象的版本号，多版本下有效
|RequestPayer|*string|启用了请求者付费模式时，需要设置为'requester'
|EnablePrefetch|bool|是否启用预取模式，默认不启用
|PrefetchNum|int|预取块的数量，默认值为3。启用预取模式时有效
|ChunkSize|int64|每个预取块的大小，默认值为6MiB。启用预取模式时有效
|PrefetchThreshold|int64|持续顺序读取多少字节后进入到预取模式，默认值为20MiB。启用预取模式时有效

**ReadOnlyFile接口：**
|接口名|说明
|:-------|:-------
|Close() error|关闭文件句柄，释放资源，例如内存，活动的socket 等
|Read(p []byte) (int, error)|从数据源中读取长度为len(p)的字节，存储到p中，返回读取的字节数和遇到的错误
|Seek(offset int64, whence int) (int64, error)|用于设置下一次读或写的偏移量。其中whence的取值：0：相对于头部，1：相对于当前偏移量，2：相对于尾部
|Stat() (os.FileInfo, error)|获取对象的信息，包括 对象大小，最后修改时间 以及元信息

> **注意:** 当预取模式打开时，如果出现多次乱序读时，则会自动退回单流模式。

示例 

1. 以单流模式，读取整个对象
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

2. 启用预取模式，读取整个对象
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

3. 通过Seek方法，从指定位置开始读取剩余的数据

```
...
client := oss.NewClient(cfg)

f, err := client.OpenFile(context.TODO(), "bucket", "key")

if err != nil {
  log.Fatalf("failed to open file %v", err)
}

defer f.Close()

// 获取对象信息
info, _ := f.Stat()

// 基本属性
fmt.Printf("size:%v, mtime:%v\n", info.Size(), info.ModTime())

// 对象元数据
if header, ok := info.Sys().(http.Header); ok {
  fmt.Printf("content-type:%v\n", header.Get(oss.HTTPHeaderContentType))
}

// 设置文件的偏移值，例如 从123开始
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

### 追加写文件(AppendOnlyFile)

调用AppendObject接口以追加写的方式上传数据。如果对象不存在，则创建追加类型的对象。如果对象存在，并且不为追加类型的对象时，则返回错误。

```
type AppendOnlyFile struct {
...
}

func (c *Client) AppendFile(ctx context.Context, bucket string, key string, optFns ...func(*AppendOptions)) (*AppendOnlyFile, error)
```

**参数列表**：
|参数名|类型|说明
|:-------|:-------|:-------
|ctx|context.Context|请求的上下文
|bucket|string|设置存储空间名字
|key|string|设置对象名
|optFns|...func(*AppendOptions)|(可选)，追加文件时的配置选项

**返回值列表**：
|返回值名|类型|说明
|:-------|:-------|:-------
|file|*AppendOnlyFile|追加文件的实例，当 err 为nil 时有效
|err|error|打开追加文件时的状态，当失败时，err 不为 nil

**AppendOptions选项说明：**
|参数|类型|说明
|:-------|:-------|:-------
|RequestPayer|*string|启用了请求者付费模式时，需要设置为'requester'
|CreateParameter|*AppendObjectRequest|用于首次上传时，设置对象的元信息，包括ContentType，Metadata，权限，存储类型 等

**AppendOnlyFile接口：**
|接口名|说明
|:-------|:-------
|Close() error|关闭文件句柄，释放资源
|Write(b []byte) (int, error)|将b中的数据写入到数据流中，返回写入的字节数和遇到的错误
|WriteFrom(r io.Reader) (int64, error)|将r中的数据写入到数据流中，返回写入的字节数和遇到的错误
|Stat() (os.FileInfo, error)|获取对象的信息，包括 对象大小，最后修改时间 以及元信息


示例 

1. 把多个本地文件合并成一个文件
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

2. 合并数据时，同时设置对象的权限和存储类型
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


## 客户端加密

客户端加密是在数据上传至OSS之前，由用户在本地对数据进行加密处理，确保只有密钥持有者才能解密数据，增强数据在传输和存储过程中的安全性。

> **注意:** 
> </br>使用客户端加密功能时，您需要对主密钥的完整性和正确性负责。
> </br>在对加密数据进行复制或者迁移时，您需要对加密元数据的完整性和正确性负责。

如果您需要了解OSS客户端加密实现的原理，请参考OSS用户指南中的[客户端加密](https://www.alibabacloud.com/help/zh/oss/user-guide/client-side-encryption)。

使用客户端加密，首先您需要实例化加密客户端，然后调用其提供的接口进行操作。您的对象将作为请求的一部分自动加密和解密。

```
type EncryptionClient struct {
  ...
}

func NewEncryptionClient(c *Client, masterCipher crypto.MasterCipher, optFns ...func(*EncryptionClientOptions)) (eclient *EncryptionClient, err error)
```

**参数列表**：
|参数名|类型|说明
|:-------|:-------|:-------
|c|*Client| 非加密客户端实例
|masterCipher|crypto.MasterCipher|主密钥实例，用于加密和解密数据密钥
|optFns|...func(*EncryptionClientOptions)|(可选)，加密客户端配置选项

**返回值列表**：
|返回值名|类型|说明
|:-------|:-------|:-------
|eclient|*EncryptionClient|加密客户端实例, 当 err 为 nil 时有效
|err|error|创建加密客户端的状态，当失败时，err 不为 nil

**EncryptionClientOptions选项说明：**
|参数|类型|说明
|:-------|:-------|:-------
|MasterCiphers|[]crypto.MasterCipher|主密钥实例组, 用于解密数据密钥。

**EncryptionClient接口：**
|基础接口名|说明
|:-------|:-------
|GetObjectMeta|获取对象的部分元信息
|HeadObject|获取对象的部元信息
|GetObject|下载对象，并自动解密
|PutObject|上传对象，并自动加密
|InitiateMultipartUpload|初始化一个分片上传事件 和 分片加密上下文（EncryptionMultiPartContext）
|UploadPart|初始化一个分片上传事件, 调用该接口上传分片数据，并自动加密。调用该接口时，需要设置 分片加密上下文
|CompleteMultipartUpload|在将所有分片数据上传完成后，调用该接口合并成一个文件
|AbortMultipartUpload|取消分片上传事件,并删除对应的分片数据
|ListParts|列举指定上传事件所属的所有已经上传成功分片
|**高级接口名**|**说明**
|NewDownloader|创建下载管理器实例
|NewUploader|创建上传管理器实例
|OpenFile|创建ReadOnlyFile实例
|**辅助接口名**|**说明**
|Unwrap|获取非加密客户端实例，可以通过该实例访问其它基础接口

> **说明:** EncryptionClient 采用了 和 Client 一样的接口命名规则 和 调用方式，有关接口的详细用法，请参考指南的其它章节说明。

### 使用RSA主密钥

**创建RAS加密客户端**

```
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/crypto"

cfg := oss.LoadDefaultConfig().
  WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
  WithRegion("your region")

client := oss.NewClient(cfg)

// 创建一个主密钥的描述信息，创建后不允许修改。主密钥描述信息和主密钥一一对应。
// 如果所有的对象都使用相同的主密钥，主密钥描述信息可以为空，但后续不支持更换主密钥。
// 如果主密钥描述信息为空，解密时无法判断使用的是哪个主密钥。
// 强烈建议为每个主密钥都配置主密钥描述信息，由客户端保存主密钥和描述信息之间的对应关系。
materialDesc := make(map[string]string)
materialDesc["desc"] = "your master encrypt key material describe information"

// 创建只包含 主密钥 的 加密客户端
mc, err := crypto.CreateMasterRsa(materialDesc, "yourRsaPublicKey", "yourRsaPrivateKey")
eclient, err := NewEncryptionClient(client, mc)

// 创建包含主密钥 和 多个解密密钥的 加密客户端
// 当解密时，先匹配解密密钥的描述信息，如果不匹配，则使用主密钥解密
//decryptMC := []crypto.MasterCipher{
//	// TODO
//}
//eclient, err := oss.NewEncryptionClient(client, mc, func(eco *oss.EncryptionClientOptions) {
//	eco.MasterCiphers = decryptMC
//})
```

**使用加密客户端上传或者下载**
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

**使用加密客户端以分片方式上传数据**
</br>以上传500K内存数据为例 
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

// 加密客户端 需要 设置分片大小和总文件大小
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

  // 加密客户端 需要 设置分片加密上下文
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

### 使用自定义主密钥
当RSA主密钥方式无法满足需求时，您可自定主密钥的加密实现。主密钥的接口定义如下：
```
type MasterCipher interface {
  Encrypt([]byte) ([]byte, error)
  Decrypt([]byte) ([]byte, error)
  GetWrapAlgorithm() string
  GetMatDesc() string
}
```
**MasterCipher接口说明**
|接口名|说明
|:-------|:-------
|Encrypt|加密 数据加密密钥 和 加密数据的初始值(IV)
|Decrypt|解密 数据加密密钥  和 加密数据的初始值(IV)
|GetWrapAlgorithm|返回 数据密钥的加密算法信息，建议采用 算法/模式/填充 格式，例如RSA/NONE/PKCS1Padding
|GetMatDesc|返回 主密钥的描述信息，JSON格式

例如

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

## 其它接口

为了方便用户使用，封装了一些易用性接口。当前扩展的接口如下：

|接口名 | 说明
|:-------|:-------
|IsObjectExist|判断对象(object)是否存在
|IsBucketExist|判断存储空间(bucket)是否存在
|PutObjectFromFile|上传本地文件到存储空间
|GetObjectToFile|下载对象到本地文件

### IsObjectExist/IsBucketExist

这两个接口的返回值为 (bool, error), 当 error 为 nil时，如果bool 为 true，表示存在，如果 bool值为 false，表示不存在。当 error 不为 nil时，表示无法从该错误信息判断 是否存在。

```
func (c *Client) IsObjectExist(ctx context.Context, bucket string, key string, optFns ...func(*IsObjectExistOptions)) (bool, error)
func (c *Client) IsBucketExist(ctx context.Context, bucket string, optFns ...func(*Options)) (bool, error)
```

例如 判断对象是否存在

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

使用简单上传(PutObject)接口 把本地文件上传到存储空间，该接口不支持并发。

```
func (c *Client) PutObjectFromFile(ctx context.Context, request *PutObjectRequest, filePath string, optFns ...func(*Options)) (*PutObjectResult, error) 
```

示例

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

使用GetObject接口，把存储空间的对象下载到本地文件，该接口不支持并发。

```
func (c *Client) GetObjectToFile(ctx context.Context, request *GetObjectRequest, filePath string, optFns ...func(*Options)) (*GetObjectResult, error) 
```

示例

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

## 上传下载接口对比

提供了各种上传下载接口，您可以根据使用场景，选择适合的接口。

**上传接口**
|接口名 | 说明
|:-------|:-------
|Client.PutObject|简单上传, 最大支持5GiB</br>支持CRC64数据校验(默认启用)</br>支持进度条</br>请求body类型为io.Reader, 当支持io.Seeker类型时，具备失败重传
|Client.PutObjectFromFile|与Client.PutObject接口能力一致</br>请求body数据来源于文件路径
|分片上传接口</br>Client.InitiateMultipartUpload</br>Client.UploadPart</br>Client.CompleteMultipartUpload|分片上传，单个分片最大5GiB，文件最大48.8TiB</br>UploadPart接口支持CRC64校验(默认启用)</br>UploadPart接口支持进度条</br>UploadPart请求body类型为io.Reader, 当支持io.Seeker类型时，具备失败重传
|Uploader.UploadFrom|封装了简单上传 和 分片上传接口，最大支持48.8TiB</br>支持CRC64数据校验(默认启用)</br>支持进度条</br>请求body参数类型为io.Reader，当同时支持 io.Reader, io.Seeker 和 io.ReaderAt 类型时，不需要把数据缓存在内存里，否则 必须先把数据缓冲在内存中，然后才能上传该部分
|Uploader.UploadFile|与Uploader.UploadFrom接口能力一致</br>请求body数据来源于文件路径</br>支持断点续传
|Client.AppendObject|追加上传, 最终文件最大支持5GiB</br>支持CRC64数据校验(默认启用)</br>支持进度条</br>请求body类型为io.Reader, 当支持io.Seeker类型时，具备失败重传(该接口为非幂等接口，重传时可能出现失败)
|AppendOnlyFile接口</br>AppendOnlyFile.Write</br>AppendOnlyFile.WriteFrom|与Client.AppendObject接口能力一致</br>优化了重传时失败后容错处理

**下载接口**
|接口名| 说明
|:-------|:-------
|Client.GetObject|流式下载, 响应体为io.ReadCloser类型</br>不直接支持CRC64校验</br>不直接支持进度条</br>流式读数据阶段，不支持失败重连
|Client.GetObjectToFile|下载到本地文件</br>单连接下载</br>支持CRC64数据校验(默认启用)</br>支持进度条</br>支持失败重连
|Downloader.DownloadFile&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|采用分片方式下载到本地文件</br>支持自定义分片大小和并发数</br>支持CRC64数据校验(默认启用)</br>支持进度条</br>支持失败重连</br>支持断点续传</br>先写临时文件，再重命名(可配置，默认启用)
|ReadOnlyFile接口</br>ReadOnlyFile.Read</br>ReadOnlyFile.Seek</br>ReadOnlyFile.Close|File-Like形式接口, 提供io.Reader, io.Seeker 和 io.Closer接口</br>具备Seek能力</br>支持单流模式(默认)</br>支持异步预取模式，提升读的速度</br>支持自定义预取块和预取数</br>不直接支持CRC64校验</br>不直接支持进度条</br>支持失败重连


# 场景示例

本部分将从使用场景出发, 介绍如何使用SDK。

包含的主题
* [设置进度条](#设置进度条)
* [数据校验](#数据校验)

## 设置进度条

在对象的上传，下载 和 拷贝 场景下，您可以设置进度条，用于查看对象的传输状态。

**支持设置进度条的请求参数**
|支持的请求参数| 用法
|:-------|:-------
|PutObjectRequest|PutObjectRequest.ProgressFunc
|GetObjectRequest|GetObjectRequest.ProgressFunc
|CopyObjectRequest|CopyObjectRequest.ProgressFunc
|AppendObjectRequest|AppendObjectRequest.ProgressFunc
|UploadPartRequest|UploadPartRequest.ProgressFunc

**ProgressFunc定义和参数说明**
```
type ProgressFunc func(increment, transferred, total int64)
```
|参数名|类型|说明
|:-------|:-------|:-------
|increment|int64|本次回调传输的数据大小,单位字节
|transferred|int64|已传输的数据大小，单位为字节
|total|int64|本次请求的数据大小，单位为字节，如果为 -1，表示获无法获取总大小


示例

1. 上传时，设置进度条，以PutObject 为例

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

2. 下载时，设置进度条，以GetObjectToFile为例
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

3. 流式下载时，设置进度条，以GetObject 为例
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

## 数据校验

OSS提供基于MD5和CRC64的数据校验，确保请求的过程中的数据完整性。

## MD5校验

当向OSS发送请求时，如果设置了Content-MD5，OSS会根据接收的内容计算MD5。当OSS计算的MD5值和上传提供的MD5值不一致时，则返回InvalidDigest异常，从而保证数据的完整性。

基础接口里，除了 PutObject, AppendObject, UploadPart 接口外，会自动计算MD5, 并设置Content-MD5, 保证请求的完整性。

如果您需要在 PutObject, AppendObject, UploadPart 接口里使用MD5校验，可以参考以下写法

```
...
client := oss.NewClient(cfg)

var body io.Reader

// 计算Content-Md5, 如果 body 不是 io.ReadSeeker 类型, 则先读到缓存里，再计算MD5
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

## CRC64校验

上传对象时，默认开启CRC64数据校验，以确保数据的完整性，例如 PutObject, AppendObject, UploadPart 等接口。

下载对象时，
* 如果是下载到本地文件，默认开启CRC64数据校验，以确保数据的完整性，例如 Downloader.DownloadFile 和 GetObjectToFile 接口。
* 如果是流式读类型的接口，不会做CRC64校验，例如 GetObject 和 ReadOnlyFile.Read 接口。

如果您需要在流式读接口里使用CRC64校验，可以参考以下写法

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

// 响应头返回的是整个文件的CRC64值，如果是范围下载，不支持CRC64校验
// 206 Partial Content 表示是范围下载
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

如果您需要关闭CRC64校验，通过Config.WithDisableDownloadCRC64Check 和  Config.WithDisableUploadCRC64Check 配置，例如
```
cfg := oss.LoadDefaultConfig().
  WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
  WithRegion(region).
  WithDisableDownloadCRC64Check(true).
  WithDisableUploadCRC64Check(true)

client := oss.NewClient(cfg)
```


# 迁移指南

本部分介绍如何从V1 版本([aliyun-oss-go-sdk](https://github.com/aliyun/aliyun-oss-go-sdk)) 迁移到 V2 版本。

## 最低 GO 版本

V2 版本 要求 Go 版本最低为 1.18。

## 导入路径

V2 版本使用新的代码仓库，同时也对代码结构进行了调整，按照功能模块组织，以下是这些模块路径和说明：

|模块路径 | 说明 
|:-------|:-------
|github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss|SDK核心，接口 和 高级接口实现
|github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials|访问凭证相关
|github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/retry|重试相关
|github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer|签名相关
|github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/transport|HTTP客户端相关
|github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/crypto|客户端加密相关

示例 

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
  // 根据需要，导入 retry，transport 或者 signer
  //"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/xxxx"
)
```

## 配置加载

V2 版本简化了配置设置方式，全部迁移到 [config](oss/config.go) 下，并提供了以With为前缀的辅助函数，方便以编程方式覆盖缺省配置。

V2 默认使用 V4签名，所以必须配置区域（Region）。

V2 支持从区域（Region）信息构造 访问域名(Endpoint), 当访问的是公有云时，可以不设置Endpoint。

示例

```
// v1
import (
  "github.com/aliyun/aliyun-oss-go-sdk/oss"
)
...

// 环境变量中获取访问凭证
provider, err := oss.NewEnvironmentVariableCredentialsProvider()

// 设置HTTP连接超时时间为20秒，HTTP读取或写入超时时间为60秒。
time := oss.Timeout(20,60)

// 不校验SSL证书校验
verifySsl := oss.InsecureSkipVerify(true)

// 设置日志
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

// 环境变量中获取访问凭证
provider := credentials.NewEnvironmentVariableCredentialsProvider()

cfg := oss.LoadDefaultConfig().
  WithCredentialsProvider(provider).
  // 设置HTTP连接超时时间为20秒
  WithConnectTimeout(20 * time.Second).
  // HTTP读取或写入超时时间为60秒
  WithReadWriteTimeout(60 * time.Second).
  // 不校验SSL证书校验
  WithInsecureSkipVerify(true).
  // 设置日志
  WithLogLevel(oss.LogInfo).
  // 设置区域
  WithRegion("cn-hangzhou")

client := oss.NewClient(cfg)
```

## 创建Client

V2 版本 把 Client 的创建 函数 从 New 修改 为 NewClient， 同时 创建函数 不在支持传入Endpoint 以及 access key id 和 access key secrect 参数。

示例

```
// v1
client, err := oss.New(endpoint, "ak", "sk")
```

```
// v2
client := oss.NewClient(cfg)
```

## 调用API操作

基础 API 接口 都 合并为 单一操作方法 '\<OperationName\>'，操作的请求参数为 '\<OperationName\>Request'，操作的返回值为 '\<OperationName\>Result'。这些操作方法都 迁移到 Client下，同时需要设置 context.Context。如下格式：

```
func (c *Client) <OperationName>(ctx context.Context, request *<OperationName>Request, optFns ...func(*Options)) (*<OperationName>Result，, error) 
```

关于API接口的详细使用说明，请参考[基础接口](#基础接口)。

示例

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

## 预签名

V2 版本 把 预签名接口 名字从 SignURL 修改为 Presign，同时把 接口 迁移到 Client 下。接口形式如下：

```
func (c *Client) Presign(ctx context.Context, request any, optFns ...func(*PresignOptions)) (*PresignResult, error)
```

对于 request 参数，其类型 与 API 接口中的 '\<OperationName\>Request' 一致。

对于返回结果，除了返回 预签名 URL 外，还返回 HTTP 方法，过期时间 和 被签名的请求头，如下：
```
type PresignResult struct {
  Method        string
  URL           string
  Expiration    time.Time
  SignedHeaders map[string]string
}
```

关于预签名的详细使用说明，请参考[预签名接口](#预签名接口)。

以 生成下载对象的预签名URL 为例，如何从 V1 迁移到 V2

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

## 断点续传接口

V2 版本使用 传输管理器 'Uploader'，'Downloader' 和 'Copier' 分别 管理 对象的 上传，下载 和 拷贝。 同时移除了原有的 断点续传接口 Bucket.UploadFile，Bucket.DownloadFile 和 Bucket.CopyFile。

接口对比如下：

|场景|v2|v1
|:-------|:-------|:-------
|上传文件|Uploader.UploadFile|Bucket.UploadFile
|上传流(io.Reader)|Uploader.UploadFrom|不支持
|下载到文件|Downloader.DownloadFile|Bucket.DownloadFile
|拷贝对象|Copier.Copy|Bucket.CopyFile

默认参数的变化

|场景|v2|v1
|:-------|:-------|:-------
|上传-分片默认值|6 MiB|通过参数设置
|上传-并发默认值|3|1
|上传-阈值|分片大小|无
|上传-记录checkpoint|支持|支持
|下载-分片默认值|6 MiB|通过参数设置
|下载-并发默认值|3|1
|下载-阈值|分片大小|无
|下载-记录checkpoint|支持|支持
|拷贝-分片默认值|64 MiB|无
|拷贝-并发默认值|3|1
|拷贝-阈值|200 MiB|无
|拷贝-记录checkpoint|不支持|支持

阈值(上传/下载拷贝) 表示 对象/文件 大小 大于该值时，使用分片方式(上传/下载/拷贝)。

关于传输管理器的详细使用说明，请参考[传输管理器](#传输管理器)。

## 客户端加密

V2 版本 使用 EncryptionClient 来提供 客户端加密功能，同时也对API 接口做了精简，采用了 和 Client 一样的接口命名规则 和 调用方式。

另外，该版本 仅保留 基于 RSA 自主管理的主密钥 的参考实现。

对于 KSM 的实现，可以参考[sample/crypto/kms.go](sample/crypto/kms.go)。

关于客户端加密的详细使用说明，请参考[客户端加密](#客户端加密)。

以 使用主密钥RSA 上传对象为例，如何从 V1 迁移到 V2

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

## 重试

V2 版本 默认开启对HTTP请求的重试行为。从 V1 版本迁移到 V2 时，您需要移除原有的重试代码，避免放大重试次数。
