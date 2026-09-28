package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/agentic"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region    string
	bucket    string
	endpoint  string
	accountId string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucket, "bucket", "", "The name of the agentic bucket.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&accountId, "account-id", "", "The account id.")
}

func main() {
	flag.Parse()
	if len(bucket) == 0 || len(region) == 0 || len(accountId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name, region, and account id required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithAccountId(accountId)
	if len(endpoint) > 0 {
		cfg.WithEndpoint(endpoint)
	}
	client := agentic.NewAgenticBucketClient(cfg)

	request := &agentic.GetAgenticBucketStorageQuotaRequest{Bucket: oss.Ptr(bucket)}
	result, err := client.GetAgenticBucketStorageQuota(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to get agentic bucket storage quota %v", err)
	}
	log.Printf("agentic bucket storage quota: storage quota:%d, mode:%s, current usage:%d\n",
		oss.ToInt64(result.QuotaConfiguration.StorageQuota),
		result.QuotaConfiguration.Mode,
		oss.ToInt64(result.QuotaConfiguration.CurrentUsage))
}
