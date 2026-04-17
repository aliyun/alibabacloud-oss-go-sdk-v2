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

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

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

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	result, err := client.GetTableMaintenanceConfiguration(context.TODO(), &tables.GetTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr(tableBucketArn),
		Namespace:      oss.Ptr(namespace),
		Name:           oss.Ptr(name),
	})

	if err != nil {
		log.Fatalf("failed to get table maintenance configuration %v", err)
	}

	log.Printf("get table maintenance configuration result:%#v\n", result)
}
