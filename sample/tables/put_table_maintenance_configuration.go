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
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketArn, "bucket-arn", "", "The arn of the table bucket.")
	flag.StringVar(&nameSpace, "name-space", "", "The name of the name space.")
	flag.StringVar(&table, "table", "", "The name of the table.")
}

func main() {
	flag.Parse()

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

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

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	// icebergCompaction type
	result, err := client.PutTableMaintenanceConfiguration(context.TODO(), &tables.PutTableMaintenanceConfigurationRequest{
		TableBucketARN: oss.Ptr(bucketArn),
		Namespace: oss.Ptr(nameSpace),
		Name:     oss.Ptr(table),
		Type:      oss.Ptr("icebergCompaction"),
		Value: &tables.TableMaintenanceValue{
			Status: oss.Ptr("enabled"),
			Settings: &tables.TableMaintenanceSettings{
				IcebergCompaction: &tables.IcebergCompactionSettingsDetail{
					TargetFileSizeMB: oss.Ptr(400),
					Strategy:         oss.Ptr("auto"),
				},
			},
		},
	})

	// icebergSnapshotManagement type
	//result, err := client.PutTableMaintenanceConfiguration(context.TODO(), &tables.PutTableMaintenanceConfigurationRequest{
	//	TableBucketARN: oss.Ptr(bucketArn),
	//	Namespace: oss.Ptr(nameSpace),
	//	Name:     oss.Ptr(table),
	//	Type:      oss.Ptr("icebergSnapshotManagement"),
	//	Value: &tables.TableMaintenanceValue{
	//		Status: oss.Ptr("enabled"),
	//		Settings: &tables.TableMaintenanceSettings{
	//			IcebergSnapshotManagement: &tables.IcebergSnapshotManagementSettingsDetail{
	//				MaxSnapshotAgeHours: oss.Ptr(350),
	//				MinSnapshotsToKeep:  oss.Ptr(1),
	//			},
	//		},
	//	},
	//})

	if err != nil {
		log.Fatalf("failed to put table maintenance configuration %v", err)
	}

	log.Printf("put table maintenance configuration result:%#v\n", result)
}
