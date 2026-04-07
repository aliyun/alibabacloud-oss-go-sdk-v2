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
	newTable  string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketArn, "bucket-arn", "", "The arn of the table bucket.")
	flag.StringVar(&nameSpace, "name-space", "", "The name of the name space.")
	flag.StringVar(&table, "table", "", "The name of the table.")
	flag.StringVar(&newTable, "new-table", "", "The new name of the table.")
}

func main() {
	flag.Parse()

	if len(bucketArn) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket arn required")
	}

	if len(nameSpace) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, name space required")
	}

	if len(table) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, table name required")
	}

	if len(newTable) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, table new name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	result, err := client.RenameTable(context.TODO(), &tables.RenameTableRequest{
		BucketArn: oss.Ptr(bucketArn),
		Namespace: oss.Ptr(nameSpace),
		Table:     oss.Ptr(table),
		NewTable:  oss.Ptr(newTable),
	})

	if err != nil {
		log.Fatalf("failed to rename table %v", err)
	}

	log.Printf("rename table result:%#v\n", result)
}
