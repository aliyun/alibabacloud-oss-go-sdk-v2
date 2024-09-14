package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
}

func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	request := &oss.PutPublicAccessBlockRequest{
		PublicAccessBlockConfiguration: &oss.PublicAccessBlockConfiguration{
			oss.Ptr(true),
		},
	}
	putResult, err := client.PutPublicAccessBlock(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put public access block %v", err)
	}

	log.Printf("put public access block result:%#v\n", putResult)
}
