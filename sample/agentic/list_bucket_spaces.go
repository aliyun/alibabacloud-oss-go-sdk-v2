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
	prefix    string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucket, "bucket", "", "The name of the agentic bucket.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&accountId, "account-id", "", "The account id.")
	flag.StringVar(&prefix, "prefix", "", "The prefix that returned bucket space names must contain.")
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

	request := &agentic.ListBucketSpacesRequest{
		Bucket: oss.Ptr(bucket),
	}
	if len(prefix) > 0 {
		request.Prefix = oss.Ptr(prefix)
	}
	p := client.NewListBucketSpacesPaginator(request)

	var i int
	log.Println("Bucket Spaces:")
	for p.HasNext() {
		i++
		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}
		for _, s := range page.BucketSpaces {
			log.Printf("Bucket Space: name:%v, location:%v, storage class:%v, creation date:%v\n",
				oss.ToString(s.Name), oss.ToString(s.Location), oss.ToString(s.StorageClass), oss.ToString(s.CreationDate))
		}
	}
}
