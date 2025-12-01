package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
	objectName string
	position   int64
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
	flag.Int64Var(&position, "position", -1, "The position of the append object.")
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

	if position < 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, position required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	request := &oss.SealAppendObjectRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
		Position: oss.Ptr(position),
	}
	result, err := client.SealAppendObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to seal append object %v", err)
	}
	log.Printf("seal append object result:%#v\n", result)
}
