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
	if len(bucketName) == 0 || len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name and region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(cfg)

	request := &oss.GetBucketStorageQuotaRequest{Bucket: oss.Ptr(bucketName)}
	result, err := client.GetBucketStorageQuota(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to get bucket storage quota %v", err)
	}
	log.Printf("bucket storage quota: storage quota:%d, mode:%s, current usage:%d\n",
		oss.ToInt64(result.QuotaConfiguration.StorageQuota),
		result.QuotaConfiguration.Mode,
		oss.ToInt64(result.QuotaConfiguration.CurrentUsage))
}
