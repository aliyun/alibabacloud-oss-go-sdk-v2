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

	result, err := client.CreateTable(context.TODO(), &tables.CreateTableRequest{
		TableBucketARN: oss.Ptr(tableBucketArn),
		Namespace:      oss.Ptr(namespace),
		Name:           oss.Ptr(name),
		Format:         oss.Ptr("ICEBERG"),
		Metadata: &tables.TableMetadata{
			Iceberg: &tables.IcebergMetadata{
				SchemaV2: tables.IcebergSchemaV2{
					Type:               oss.Ptr("struct"),
					IdentifierFieldIDs: []int64{10},
					Fields: tables.IcebergFields{
						{
							ID: oss.Ptr(int64(10)), Name: oss.Ptr("order_id"), Type: "long",
							Required: oss.Ptr(true), Doc: oss.Ptr("The identifier of the order."),
						},
						{
							ID: oss.Ptr(int64(20)), Name: oss.Ptr("customer"), Required: oss.Ptr(false),
							Type: tables.IcebergStructType{
								Type: oss.Ptr("struct"),
								Fields: tables.IcebergFields{
									{ID: oss.Ptr(int64(21)), Name: oss.Ptr("name"), Type: "string"},
									{ID: oss.Ptr(int64(22)), Name: oss.Ptr("level"), Type: "int"},
								},
							},
						},
						{
							ID: oss.Ptr(int64(30)), Name: oss.Ptr("tags"), Required: oss.Ptr(false),
							Type: tables.IcebergListType{
								Type: oss.Ptr("list"), ElementID: oss.Ptr(int64(31)),
								Element: "string", ElementRequired: oss.Ptr(false),
							},
						},
						{
							ID: oss.Ptr(int64(40)), Name: oss.Ptr("attributes"), Required: oss.Ptr(false),
							Type: tables.IcebergMapType{
								Type: oss.Ptr("map"), KeyID: oss.Ptr(int64(41)), Key: "string",
								ValueID: oss.Ptr(int64(42)), Value: "string", ValueRequired: oss.Ptr(false),
							},
						},
					},
				}.ToMap(),
			},
		},
	})

	if err != nil {
		log.Fatalf("failed to create table with schema v2 %v", err)
	}

	log.Printf("create table with schema v2 result:%#v\n", result)
}
