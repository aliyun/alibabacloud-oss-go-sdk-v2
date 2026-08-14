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

	request := &agentic.GetAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	}
	result, err := client.GetAgenticBucket(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to get agentic bucket %v", err)
	}

	if info := result.AgenticBucketInfo; info != nil {
		log.Printf("agentic bucket info: name:%v, owner:%v, region:%v, storage class:%v, status:%v, acl:%v, versioning:%v, create time:%v\n",
			oss.ToString(info.Name), oss.ToString(info.Owner), oss.ToString(info.Region),
			oss.ToString(info.StorageClass), oss.ToString(info.Status), oss.ToString(info.ACL),
			oss.ToString(info.Versioning), oss.ToString(info.CreateTime))
	}
}
