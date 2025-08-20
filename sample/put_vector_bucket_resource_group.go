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
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the vector bucket.")
}

func main() {
	flag.Parse()
	var (
		groupId = "resource group id"
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

	client := oss.NewVectorsClient(cfg)

	request := &oss.PutBucketResourceGroupRequest{
		Bucket: oss.Ptr(bucketName),
		BucketResourceGroupConfiguration: &oss.BucketResourceGroupConfiguration{
			oss.Ptr(groupId),
		},
	}
	result, err := client.PutBucketResourceGroup(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put vector bucket resource group %v", err)
	}

	log.Printf("put vector bucket resource group result:%#v\n", result)
}
