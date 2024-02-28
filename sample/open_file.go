package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	endpoint   string
	bucketName string
	objectName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
}

func main() {
	flag.Parse()
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(endpoint) == 0 {
		endpoint = fmt.Sprintf("oss-%v.aliyuncs.com", region)
	}

	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithEndpoint(endpoint)
	client := oss.NewClient(cfg)

	f, err := client.OpenFile(context.TODO(), bucketName, objectName)
	defer f.Close()
	if err != nil {
		log.Fatalf("failed to open file %v", err)
	}
	stat, err := f.Stat()
	if err != nil {
		log.Fatalf("failed to stat file %v", err)
	}
	var contentBuilder strings.Builder
	buf := make([]byte, stat.Size())
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("failed to read object: %v", err)
		}
		if n == 0 {
			break
		}
		contentBuilder.Write(buf[:n])
	}
	fmt.Printf("content:%#v\n", contentBuilder.String())
}
