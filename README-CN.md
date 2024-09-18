# Alibaba Cloud OSS SDK for Go v2

[![GitHub version](https://badge.fury.io/gh/aliyun%2Falibabacloud-oss-go-sdk-v2.svg)](https://badge.fury.io/gh/aliyun%2Falibabacloud-oss-go-sdk-v2)

alibabacloud-oss-go-sdk-v2 是OSS在Go编译语言下的第二版SDK

## [README in English](README.md)

## 关于
> - 此Go SDK基于[阿里云对象存储服务](http://www.aliyun.com/product/oss/)官方API构建。
> - 阿里云对象存储（Object Storage Service，简称OSS），是阿里云对外提供的海量，安全，低成本，高可靠的云存储服务。
> - OSS适合存放任意文件类型，适合各种网站、开发企业及开发者使用。
> - 使用此SDK，用户可以方便地在任何应用、任何时间、任何地点上传，下载和管理数据。

## 运行环境
> - Go 1.18及以上。

## 安装方法
### GitHub安装
> - 执行命令`go get github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss`获取远程代码包。
> - 在您的代码中使用`import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"`引入OSS Go SDK的包。

## 快速使用
#### 获取存储空间列表（List Bucket）
```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func main() {
	var (
		region = "cn-hangzhou"
	)

	// Using the SDK's default configuration
	// loading credentials values from the environment variables
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	// Create the Paginator for the ListBuckets operation.
	p := client.NewListBucketsPaginator(&oss.ListBucketsRequest{})

	// Iterate through the bucket pages
	var i int
	fmt.Println("Buckets:")
	for p.HasNext() {
		i++
		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}
		// Print the bucket found
		for _, b := range page.Buckets {
			fmt.Printf("Bucket:%v, %v, %v\n", oss.ToString(b.Name), oss.ToString(b.StorageClass), oss.ToString(b.Location))
		}
	}
}
```

#### 获取文件列表（List Objects）
```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func main() {
	var (
		region   = "cn-hangzhou"
		bucketName = "your bucket name"
	)

	// Using the SDK's default configuration
	// loading credentials values from the environment variables
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	// Create the Paginator for the ListObjectsV2 operation.
	p := client.NewListObjectsV2Paginator(&oss.ListObjectsV2Request{
		Bucket: oss.Ptr(bucketName),
	})

	// Iterate through the object pages
	var i int
	fmt.Println("Objects:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Print the objects found
		for _, obj := range page.Contents {
			fmt.Printf("Object:%v, %v, %v\n", oss.ToString(obj.Key), obj.Size, oss.ToTime(obj.LastModified))
		}
	}
}
```

#### 上传文件（Put Object）
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func main() {
	var (
		region   = "cn-hangzhou"
		bucketName = "your bucket name"
		objectName = "your object name"
		localFile  = "your local file path"
	)

	// Using the SDK's default configuration
	// loading credentials values from the environment variables
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	file, err := os.Open(localFile)
	if err != nil {
		log.Fatalf("failed to open file %v", err)
	}
	defer file.Close()

	result, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
		Body:   file,
	})

	if err != nil {
		log.Fatalf("failed to put object %v", err)
	}

	fmt.Printf("put object sucessfully, ETag :%v\n", result.ETag)
}
```

## 更多示例
请参看`sample`目录

### 运行示例
> - 进入示例程序目录 `sample`。
> - 通过环境变量，配置访问凭证, `export OSS_ACCESS_KEY_ID="your access key id"`, `export OSS_ACCESS_KEY_SECRET="your access key secrect"`
> - 以 list_buckets.go 为例，执行 `go run list_buckets.go -region cn-hangzhou`。

## 资源
[开发者指南](DEVGUIDE-CN.md) - 参阅该指南，来帮助您安装、配置和使用该开发套件。

## 许可协议
> - Apache-2.0, 请参阅 [许可文件](LICENSE)

