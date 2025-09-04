package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/vectors"
)

var (
	region    string
	accountId string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&accountId, "account-id", "", "The id of vector account.")
}

func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(accountId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, accounId required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithAccountId(accountId)

	client := vectors.NewVectorsClient(cfg)

	request := &vectors.ListVectorBucketsRequest{}

	p := client.NewListVectorBucketsPaginator(request)

	var i int
	log.Println("Vector Buckets:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, b := range page.Buckets.Bucket {
			log.Printf("Bucket:%v, %v, %v\n", oss.ToString(b.Name), oss.ToString(b.ResourceGroupId), oss.ToString(b.Location))
		}
	}
}
