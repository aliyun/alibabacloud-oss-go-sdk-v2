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
	var (
		wormId = "worm id"
	)
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

	request := &oss.ExtendBucketWormRequest{
		Bucket: oss.Ptr(bucketName),
		WormId: oss.Ptr(wormId),
		ExtendWormConfiguration: &oss.ExtendWormConfiguration{
			2,
		},
	}
	result, err := client.ExtendBucketWorm(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to extend bucket worm %v", err)
	}

	log.Printf("extend bucket worm result:%#v\n", result)
}
