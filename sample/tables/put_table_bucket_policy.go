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

	result, err := client.PutTableBucketPolicy(context.TODO(), &tables.PutTableBucketPolicyRequest{
		TableBucketARN: oss.Ptr(tableBucketArn),
		ResourcePolicy: oss.Ptr(`
			{
			   "Version":"1",
			   "Statement":[
			   {
				 "Action":[
				   "oss:GetTable",
				],
				"Effect":"Deny",
				"Principal":["1234567890"],
				"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket"]
			   }
			  ]
			 }
		`),
	})

	if err != nil {
		log.Fatalf("failed to put table bucket policy %v", err)
	}

	log.Printf("put table bucket policy result:%#v\n", result)
}
