package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	cloudBoxEndpoint string
)

func init() {
	flag.StringVar(&cloudBoxEndpoint, "endpoint", "", "The endpoint of cloud box.")
}

func main() {
	flag.Parse()
	if len(cloudBoxEndpoint) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, the endpoint of cloud box required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithEndpoint(cloudBoxEndpoint)

	client := oss.NewClient(cfg)

	request := &oss.ListCloudBoxesRequest{}

	result, err := client.ListCloudBoxes(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to list cloud boxes %v", err)
	}
	log.Printf(" list cloud boxes result:%#v\n", result)
}
