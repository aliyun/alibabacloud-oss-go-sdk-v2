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
	region           string
	tableBucketArn   string
	namespace        string
	name             string
	metadataLocation string
	versionToken     string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&tableBucketArn, "table-bucket-arn", "", "The arn of the table bucket.")
	flag.StringVar(&namespace, "namespace", "", "The name of the namespace.")
	flag.StringVar(&name, "name", "", "The name of the table.")
	flag.StringVar(&metadataLocation, "metadata-location", "", "The metadata location of the table.")
	flag.StringVar(&versionToken, "version-token", "", "The version token of the table.")
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

	if len(metadataLocation) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, table metadata location required")
	}

	if len(versionToken) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, table version token required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	result, err := client.UpdateTableMetadataLocation(context.TODO(), &tables.UpdateTableMetadataLocationRequest{
		TableBucketARN:   oss.Ptr(tableBucketArn),
		Namespace:        oss.Ptr(namespace),
		Name:             oss.Ptr(name),
		MetadataLocation: oss.Ptr(metadataLocation),
		VersionToken:     oss.Ptr(versionToken),
	})

	if err != nil {
		log.Fatalf("failed to update table metadata location %v", err)
	}

	log.Printf("update table metadata location result:%#v\n", result)
}
