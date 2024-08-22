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
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
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

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	request := &oss.PutBucketTransferAccelerationRequest{
		Bucket: oss.Ptr(bucketName),
		TransferAccelerationConfiguration: &oss.TransferAccelerationConfiguration{
			oss.Ptr(true),
		},
	}
	result, err := client.PutBucketTransferAcceleration(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket transfer acceleration %v", err)
	}

	log.Printf("put bucket transfer acceleration result:%#v\n", result)
}
