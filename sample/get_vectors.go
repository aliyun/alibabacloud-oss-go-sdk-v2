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

	request := &oss.GetVectorsRequest{
		Bucket:         oss.Ptr(bucketName),
		IndexName:      oss.Ptr("index"),
		Keys:           []string{"key1", "key2", "key3"},
		ReturnData:     oss.Ptr(true),
		ReturnMetadata: oss.Ptr(false),
	}
	result, err := client.GetVectors(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to get vectors %v", err)
	}
	log.Printf("get vectors result:%#v\n", result)
}
