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
	var id = "your bucket data redundancy transition task id"
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

	request := &oss.DeleteBucketDataRedundancyTransitionRequest{
		Bucket:                     oss.Ptr(bucketName),
		RedundancyTransitionTaskid: oss.Ptr(id),
	}
	result, err := client.DeleteBucketDataRedundancyTransition(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to delete bucket data redundancy transition %v", err)
	}

	log.Printf("delete bucket data redundancy transition result:%#v\n", result)
}
