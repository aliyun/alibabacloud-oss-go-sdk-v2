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
	endpoint  string
	accountId string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&accountId, "account-id", "", "The account id.")
}

func main() {
	flag.Parse()
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

	request := &agentic.ListAgenticBucketsRequest{}
	p := client.NewListAgenticBucketsPaginator(request)

	var i int
	log.Println("Agentic Buckets:")
	for p.HasNext() {
		i++
		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}
		for _, b := range page.AgenticBuckets {
			log.Printf("Agentic Bucket: name:%v, storage class:%v, data redundancy type:%v, create time:%v\n",
				oss.ToString(b.Name), oss.ToString(b.StorageClass), oss.ToString(b.DataRedundancyType), oss.ToString(b.CreateTime))
		}
	}
}
