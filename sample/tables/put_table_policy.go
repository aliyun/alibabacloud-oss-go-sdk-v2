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
	namespace      string
	name           string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&tableBucketArn, "table-bucket-arn", "", "The arn of the table bucket.")
	flag.StringVar(&namespace, "namespace", "", "The name of the namespace.")
	flag.StringVar(&name, "name", "", "The name of the table.")
}

func main() {
	flag.Parse()
	if len(tableBucketArn) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, table bucket arn required")
	}

	if len(namespace) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, namespace name required")
	}

	if len(name) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, table name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	result, err := client.PutTablePolicy(context.TODO(), &tables.PutTablePolicyRequest{
		TableBucketARN: oss.Ptr(tableBucketArn),
		Namespace:      oss.Ptr(namespace),
		Name:           oss.Ptr(name),
		ResourcePolicy: oss.Ptr(`
			{
			   "Version":"1",
			   "Statement":[
			   {
				 "Action":[
				   "oss:GetTable"
				],
				"Effect":"Deny",
				"Principal":["1234567890"],
				"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/*"]
			   }
			  ]
			 }
		`),
	})

	if err != nil {
		log.Fatalf("failed to put table policy %v", err)
	}

	log.Printf("put table policy result:%#v\n", result)
}
