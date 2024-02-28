package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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

	localFile := "upload.file"
	file, err := os.Create(localFile)
	defer file.Close()
	if err != nil {
		log.Fatalf("failed to create file %v", err)
	}
	_, err = io.WriteString(file, "hi oss")
	if err != nil {
		log.Fatalf("failed to write content %v", err)
	}
	result, err := client.PutObjectFromFile(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}, localFile)
	if err != nil {
		log.Fatalf("failed to put object from file %v", err)
	}
	log.Printf("put object from file result:%#v\n", result)
}
