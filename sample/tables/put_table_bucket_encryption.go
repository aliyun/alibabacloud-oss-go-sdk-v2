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
	region         string
	tableBucketArn string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&tableBucketArn, "table-bucket-arn", "", "The arn of the table bucket.")
}

func main() {
	flag.Parse()
	if len(tableBucketArn) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, table bucket arn required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	result, err := client.PutTableBucketEncryption(context.TODO(), &tables.PutTableBucketEncryptionRequest{
		TableBucketARN: oss.Ptr(tableBucketArn),
		EncryptionConfiguration: &tables.EncryptionConfiguration{
			SseAlgorithm: oss.Ptr("AES256"),
		},
	})

	if err != nil {
		log.Fatalf("failed to put table bucket encryption %v", err)
	}

	log.Printf("put table bucket encryption result:%#v\n", result)
}
