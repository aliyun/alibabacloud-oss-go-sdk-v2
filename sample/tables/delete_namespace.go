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
	region    string
	bucketArn string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketArn, "bucket-arn", "", "The bucket arn of the bucket.")
}

func main() {
	flag.Parse()
	if len(bucketArn) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket arn required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	result, err := client.DeleteNamespace(context.TODO(), &tables.DeleteNamespaceRequest{
		BucketArn: oss.Ptr(bucketArn),
		Namespace: oss.Ptr("my_space"),
	})

	if err != nil {
		log.Fatalf("failed to delete namespace %v", err)
	}

	log.Printf("delete namespace result:%#v\n", result)
}
