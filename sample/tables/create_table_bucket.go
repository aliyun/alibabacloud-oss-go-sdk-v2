package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/tables"
)

var (
	region string
	bucket string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucket, "bucket", "", "The bucket name of the bucket.")
}

func main() {
	flag.Parse()
	if len(bucket) == 0 {
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

	client := tables.NewTablesClient(cfg)

	result, err := client.CreateTableBucket(context.TODO(), &tables.CreateTableBucketRequest{
		Bucket: oss.Ptr(bucket),
	})

	if err != nil {
		log.Fatalf("failed to create table bucket %v", err)
	}

	log.Printf("create table bucket result:%#v\n", result)
}
