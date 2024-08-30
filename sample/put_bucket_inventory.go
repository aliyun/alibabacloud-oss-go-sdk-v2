package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
}

func main() {
	flag.Parse()
	var (
		accountId   = "account id of the bucket"
		inventoryId = "inventory id"
	)
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	putRequest := &oss.PutBucketInventoryRequest{
		Bucket:      oss.Ptr(bucketName),
		InventoryId: oss.Ptr(inventoryId),
		InventoryConfiguration: &oss.InventoryConfiguration{
			Id:        oss.Ptr(inventoryId),
			IsEnabled: oss.Ptr(true),
			Filter: &oss.InventoryFilter{
				Prefix:                   oss.Ptr("filterPrefix"),
				LastModifyBeginTimeStamp: oss.Ptr(int64(1637883649)),
				LastModifyEndTimeStamp:   oss.Ptr(int64(1638347592)),
				LowerSizeBound:           oss.Ptr(int64(1024)),
				UpperSizeBound:           oss.Ptr(int64(1048576)),
				StorageClass:             oss.Ptr("Standard,IA"),
			},
			Destination: &oss.InventoryDestination{
				&oss.InventoryOSSBucketDestination{
					Format:    oss.InventoryFormatCSV,
					AccountId: oss.Ptr(accountId),
					RoleArn:   oss.Ptr("acs:ram::" + accountId + ":role/AliyunOSSRole"),
					Bucket:    oss.Ptr("acs:oss:::" + bucketName),
				},
			},
			Schedule: &oss.InventorySchedule{
				oss.InventoryFrequencyDaily,
			},
			IncludedObjectVersions: oss.Ptr("All"),
		},
	}
	putResult, err := client.PutBucketInventory(context.TODO(), putRequest)
	if err != nil {
		log.Fatalf("failed to put bucket inventory %v", err)
	}

	log.Printf("put bucket inventory result:%#v\n", putResult)
}
