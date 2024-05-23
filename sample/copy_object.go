package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region         string
	bucketName     string
	objectName     string
	destBucketName string
	destObjectName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "src-object", "", "The name of the source object.")
	flag.StringVar(&destBucketName, "dest-bucket", "", "The name of the destination bucket.")
	flag.StringVar(&destObjectName, "dest-object", "", "The name of the destination object.")
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

	if len(destBucketName) == 0 {
		destBucketName = bucketName
	}

	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, src object name required")
	}

	if len(destObjectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, destination object name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	request := &oss.CopyObjectRequest{
		Bucket:       oss.Ptr(destBucketName),
		Key:          oss.Ptr(destObjectName),
		SourceKey:    oss.Ptr(objectName),
		SourceBucket: oss.Ptr(bucketName),
	}
	result, err := client.CopyObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to copy object %v", err)
	}
	log.Printf("copy object result:%#v\n", result)
}
