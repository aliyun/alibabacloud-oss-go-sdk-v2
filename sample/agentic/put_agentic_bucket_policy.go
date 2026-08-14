package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

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

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithAccountId(accountId)
	if len(endpoint) > 0 {
		cfg.WithEndpoint(endpoint)
	}

	client := agentic.NewAgenticBucketClient(cfg)

	policy := fmt.Sprintf(`{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["*"],"Resource":["acs:oss:*:%s:*"]}]}`, accountId)

	request := &agentic.PutAgenticBucketPolicyRequest{
		Bucket: oss.Ptr(bucket),
		Body:   strings.NewReader(policy),
	}
	result, err := client.PutAgenticBucketPolicy(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put agentic bucket policy %v", err)
	}
	log.Printf("put agentic bucket policy result:%#v\n", result)
}
