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

	result, err := client.DoDataPipeLineAction(context.TODO(), &oss.DoDataPipeLineActionRequest{
		Action: oss.Ptr("listDataPipelineConfigurations"),
	})
	if err != nil {
		log.Fatalf("failed to do data pipe line action %v", err)
	}
	log.Printf("do data pipe line action result:%#v\n", result)
}
