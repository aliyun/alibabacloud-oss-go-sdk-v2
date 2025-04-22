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
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	// case 1: meta search
	request := &oss.OpenMetaQueryRequest{
		Bucket: oss.Ptr(bucketName),
	}
	result, err := client.OpenMetaQuery(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to open meta query %v", err)
	}

	log.Printf("open meta query result:%#v\n", result)

	// case 2: ai search
	//request = &oss.OpenMetaQueryRequest{
	//	Bucket: oss.Ptr(bucketName),
	//	Mode:   oss.Ptr("semantic"),
	//}
	//result, err = client.OpenMetaQuery(context.TODO(), request)
	//if err != nil {
	//	log.Fatalf("failed to open meta query %v", err)
	//}
	//
	//log.Printf("open meta query result:%#v\n", result)
}
