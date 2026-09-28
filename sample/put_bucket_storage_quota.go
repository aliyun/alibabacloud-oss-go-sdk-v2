package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region       string
	bucketName   string
	storageQuota int64
	quotaMode    string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.Int64Var(&storageQuota, "storage-quota", 0, "The maximum storage capacity in bytes.")
	flag.StringVar(&quotaMode, "mode", "Strict", "The storage quota enforcement mode: Strict or Warning.")
}

func main() {
	flag.Parse()
	if len(bucketName) == 0 || len(region) == 0 || storageQuota <= 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name, region, and positive storage quota required")
	}
	if quotaMode != string(oss.StorageQuotaModeStrict) && quotaMode != string(oss.StorageQuotaModeWarning) {
		flag.PrintDefaults()
		log.Fatalf("invalid quota mode, must be Strict or Warning")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(cfg)

	request := &oss.PutBucketStorageQuotaRequest{
		Bucket: oss.Ptr(bucketName),
		QuotaConfiguration: &oss.QuotaConfiguration{
			StorageQuota: oss.Ptr(storageQuota),
			Mode:         oss.StorageQuotaModeType(quotaMode),
		},
	}
	result, err := client.PutBucketStorageQuota(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket storage quota %v", err)
	}
	log.Printf("put bucket storage quota result:%#v\n", result)
}
