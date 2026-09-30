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

	request := &oss.CreateBucketMetadataTableConfigurationRequest{
		Bucket: oss.Ptr(bucketName),
		MetadataConfiguration: &oss.MetadataConfiguration{
			JournalTableConfiguration: &oss.JournalTableConfiguration{
				RecordExpiration: &oss.RecordExpiration{
					Expiration: oss.Ptr("ENABLED"),
					Days:       oss.Ptr(7),
				},
			},
			InventoryTableConfiguration: &oss.InventoryTableConfiguration{
				ConfigurationState: oss.Ptr("ENABLED"),
				EncryptionConfiguration: &oss.MetadataTableEncryptionConfiguration{
					SseAlgorithm: oss.Ptr("AES256"),
				},
			},
		},
	}
	result, err := client.CreateBucketMetadataTableConfiguration(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to create bucket metadata table configuration %v", err)
	}

	log.Printf("create bucket metadata table configuration result:%#v\n", result)
}
