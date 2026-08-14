package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/agentic"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region        string
	bucket        string
	endpoint      string
	accountId     string
	agenticBucket string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucket, "bucket", "", "The name of the bucket space.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&accountId, "account-id", "", "The account id.")
	flag.StringVar(&agenticBucket, "agentic-bucket", "", "The name of the agentic bucket that the bucket space belongs to.")
}

func main() {
	flag.Parse()
	if len(bucket) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}
	if len(accountId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, account id required")
	}
	if len(agenticBucket) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, agentic bucket name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithAccountId(accountId)
	if len(endpoint) > 0 {
		cfg.WithEndpoint(endpoint)
	}

	client := agentic.NewBucketSpaceClient(cfg)

	// The bucket space must be created under an agentic bucket, identified by its
	// full name "{bucket}-{accountId}-{region}-ab-apsr".
	request := &oss.PutBucketRequest{
		Bucket:        oss.Ptr(bucket),
		AgenticBucket: oss.Ptr(fmt.Sprintf("%s-%s-%s-ab-apsr", agenticBucket, accountId, region)),
	}
	result, err := client.PutBucket(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket space %v", err)
	}
	log.Printf("put bucket space result:%#v\n", result)
}
