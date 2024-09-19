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

	request := &oss.UpdateUserAntiDDosInfoRequest{
		DefenderInstance: oss.Ptr(instanceId),
		DefenderStatus:   oss.Ptr("HaltDefending"),
	}
	result, err := client.UpdateUserAntiDDosInfo(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to update user anti ddos info %v", err)
	}

	log.Printf("update user anti ddos info result:%#v\n", result)
}
