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
	nameSpace string
	table     string
	tableArn  string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketArn, "bucket-arn", "", "The arn of the table bucket.")
	flag.StringVar(&nameSpace, "name-space", "", "The name of the name space.")
	flag.StringVar(&table, "table", "", "The name of the table.")
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
		BucketArn: oss.Ptr(bucketArn),
		Namespace: oss.Ptr(nameSpace),
		Name:     oss.Ptr(table),
	})

	// function two
	//result, err := client.GetTable(context.TODO(), &tables.GetTableRequest{
	//	TableArn: oss.Ptr(tableArn),
	//})

	if err != nil {
		log.Fatalf("failed to get table %v", err)
	}

	log.Printf("get table result:%#v\n", result)
}
