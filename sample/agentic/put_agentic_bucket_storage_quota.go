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
	region       string
	bucket       string
	endpoint     string
	accountId    string
	storageQuota int64
	quotaMode    string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucket, "bucket", "", "The name of the agentic bucket.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&accountId, "account-id", "", "The account id.")
	flag.Int64Var(&storageQuota, "storage-quota", 0, "The default maximum storage capacity in bytes for new bucket spaces.")
	flag.StringVar(&quotaMode, "mode", "Strict", "The storage quota enforcement mode: Strict or Warning.")
}

func main() {
	flag.Parse()
	if len(bucket) == 0 || len(region) == 0 || len(accountId) == 0 || storageQuota <= 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name, region, account id, and positive storage quota required")
	}
	if quotaMode != string(oss.StorageQuotaModeStrict) && quotaMode != string(oss.StorageQuotaModeWarning) {
		flag.PrintDefaults()
		log.Fatalf("invalid quota mode, must be Strict or Warning")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithAccountId(accountId)
	if len(endpoint) > 0 {
		cfg.WithEndpoint(endpoint)
	}
	client := agentic.NewAgenticBucketClient(cfg)

	request := &agentic.PutAgenticBucketStorageQuotaRequest{
		Bucket: oss.Ptr(bucket),
		QuotaConfiguration: &oss.QuotaConfiguration{
			StorageQuota: oss.Ptr(storageQuota),
			Mode:         oss.StorageQuotaModeType(quotaMode),
		},
	}
	result, err := client.PutAgenticBucketStorageQuota(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put agentic bucket storage quota %v", err)
	}
	log.Printf("put agentic bucket storage quota result:%#v\n", result)
}
