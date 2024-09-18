# Alibaba Cloud OSS SDK for Go v2

[![GitHub version](https://badge.fury.io/gh/aliyun%2Falibabacloud-oss-go-sdk-v2.svg)](https://badge.fury.io/gh/aliyun%2Falibabacloud-oss-go-sdk-v2)

alibabacloud-oss-go-sdk-v2 is the v2 of the OSS SDK for the Go programming language

## [README in Chinese](README-CN.md)

## About
> - This Go SDK is based on the official APIs of [Alibaba Cloud OSS](http://www.aliyun.com/product/oss/).
> - Alibaba Cloud Object Storage Service (OSS) is a cloud storage service provided by Alibaba Cloud, featuring massive capacity, security, a low cost, and high reliability. 
> - The OSS can store any type of files and therefore applies to various websites, development enterprises and developers.
> - With this SDK, you can upload, download and manage data on any app anytime and anywhere conveniently. 

## Running Environment
> - Go 1.18 or above. 

## Installing
### Install the SDK through GitHub
> - Run the `go get github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss` command to get the remote code package.
> - Use `import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"` in your code to introduce OSS Go SDK package.

## Getting Started
#### List Bucket
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

#### List Objects
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
		region     = "cn-hangzhou"
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

#### Put Object
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
		region     = "cn-hangzhou"
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

##  Complete Example
More example projects can be found in the `sample` folder 

### Running Example
> - Go to the sample code folder `sample`。
> - Configure credentials values from the environment variables, like `export OSS_ACCESS_KEY_ID="your access key id"`, `export OSS_ACCESS_KEY_SECRET="your access key secrect"`
> - Take list_buckets.go as an example，run `go run list_buckets.go -region cn-hangzhou` command。

## Resources
[Developer Guide](DEVGUIDE.md) - Use this document to learn how to get started and use this sdk.

## License
> - Apache-2.0, see [license file](LICENSE)
