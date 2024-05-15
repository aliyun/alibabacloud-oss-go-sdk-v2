package main

import (
	"context"
	"flag"
	"io"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
	objectName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
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

	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	request := &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	result, err := client.GetObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to get object %v", err)
	}
	defer result.Body.Close()
	log.Printf("get object result:%#v\n", result)
	data, _ := io.ReadAll(result.Body)
	log.Printf("body:%s\n", data)
}
