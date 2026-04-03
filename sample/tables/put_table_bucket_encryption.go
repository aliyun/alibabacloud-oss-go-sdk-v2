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

	result, err := client.PutTableBucketEncryption(context.TODO(), &tables.PutTableBucketEncryptionRequest{
		BucketArn: oss.Ptr(bucketArn),
		EncryptionConfiguration: &tables.EncryptionConfiguration{
			SseAlgorithm: oss.Ptr("AES256"),
		},
	})

	if err != nil {
		log.Fatalf("failed to put table bucket encryption %v", err)
	}

	log.Printf("put table bucket encryption result:%#v\n", result)
}
