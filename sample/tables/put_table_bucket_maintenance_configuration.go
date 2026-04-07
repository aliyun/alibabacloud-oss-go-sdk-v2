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
	flag.StringVar(&bucketArn, "bucket-arn", "", "The arn of the table bucket.")
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

	result, err := client.PutTableBucketMaintenanceConfiguration(context.TODO(), &tables.PutTableBucketMaintenanceConfigurationRequest{
		BucketArn: oss.Ptr(bucketArn),
		Type:      oss.Ptr("icebergUnreferencedFileRemoval"),
		Value: &tables.MaintenanceValue{
			Settings: &tables.MaintenanceSettings{
				IcebergUnreferencedFileRemoval: &tables.SettingsDetail{
					UnreferencedDays: oss.Ptr(int(4)),
					NonCurrentDays:   oss.Ptr(10),
				},
			},
			Status: oss.Ptr("disabled"),
		},
	})

	if err != nil {
		log.Fatalf("failed to put table bucket maintenance configuration %v", err)
	}

	log.Printf("put table bucket maintenance configuration result:%#v\n", result)
}
