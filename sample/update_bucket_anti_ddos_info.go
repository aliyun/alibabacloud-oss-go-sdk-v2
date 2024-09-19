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
	var (
		instanceId = "defender instance id"
	)
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

	request := &oss.UpdateBucketAntiDDosInfoRequest{
		Bucket:           oss.Ptr(bucketName),
		DefenderInstance: oss.Ptr(instanceId),
		DefenderStatus:   oss.Ptr("Init"),
	}
	result, err := client.UpdateBucketAntiDDosInfo(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to update bucket anti ddos info %v", err)
	}

	log.Printf("update bucket anti ddos info result:%#v\n", result)
}
