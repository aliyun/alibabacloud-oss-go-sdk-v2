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
		accessPointName = "access point name"
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

	request := &oss.CreateAccessPointRequest{
		Bucket: oss.Ptr(bucketName),
		CreateAccessPointConfiguration: &oss.CreateAccessPointConfiguration{
			AccessPointName: oss.Ptr(accessPointName),
			NetworkOrigin:   oss.Ptr("internet"),
		},
	}
	result, err := client.CreateAccessPoint(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to create access point %v", err)
	}

	log.Printf("create access point result:%#v\n", result)
}
