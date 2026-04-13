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
	tableArn       string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&tableBucketArn, "table-bucket-arn", "", "The arn of the table bucket.")
	flag.StringVar(&namespace, "namespace", "", "The name of the namespace.")
	flag.StringVar(&name, "name", "", "The name of the table.")
	flag.StringVar(&tableArn, "table-arn", "", "The arn of the table.")
}

func main() {
	flag.Parse()

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	// function one
	result, err := client.GetTable(context.TODO(), &tables.GetTableRequest{
		TableBucketARN: oss.Ptr(tableBucketArn),
		Namespace:      oss.Ptr(namespace),
		Name:           oss.Ptr(name),
	})

	// function two
	//result, err := client.GetTable(context.TODO(), &tables.GetTableRequest{
	//	TableARN: oss.Ptr(tableArn),
	//})

	if err != nil {
		log.Fatalf("failed to get table %v", err)
	}

	log.Printf("get table result:%#v\n", result)
}
