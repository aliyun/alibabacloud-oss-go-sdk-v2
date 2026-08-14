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

	request := &agentic.GetAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr(bucket),
	}
	result, err := client.GetAgenticBucketEncryption(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to get agentic bucket encryption %v", err)
	}

	if rule := result.ServerSideEncryptionRule; rule != nil && rule.ApplyServerSideEncryptionByDefault != nil {
		apply := rule.ApplyServerSideEncryptionByDefault
		log.Printf("get agentic bucket encryption result: sse algorithm:%v, kms master key id:%v, kms data encryption:%v\n",
			oss.ToString(apply.SSEAlgorithm), oss.ToString(apply.KMSMasterKeyID), oss.ToString(apply.KMSDataEncryption))
	}
}
