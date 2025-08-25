package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region string
	uid    string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&uid, "uid", "", "The id of vector account.")
}

func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(uid) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, uid required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithUid(uid)

	client := oss.NewVectorsClient(cfg)

	request := &oss.ListVectorBucketsRequest{}

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
