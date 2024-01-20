package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region       string
	endpoint     string
	bucketName   string
	bucketPrefix string
	maxKeys      int
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketPrefix, "prefix", "", "[optional]`bucket prefix` of the bucket name to list.")
	flag.IntVar(&maxKeys, "max-keys", 0, "[optional]The maximum number of `keys per page` to retrieve at once.")
}

// Lists all objects in a bucket using paginator
func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	if len(endpoint) > 0 {
		cfg.WithEndpoint(endpoint)
	}

	client := oss.NewClient(cfg)

	// Set the request
	request := &oss.ListBucketsRequest{}

	if len(bucketPrefix) > 0 {
		request.Prefix = oss.Ptr(bucketPrefix)
	}

	if maxKeys > 0 {
		request.MaxKeys = int32(maxKeys)
	}

	// Create the Paginator for the ListBuckets operation.
	p := client.NewListBucketsPaginator(request)

	// Iterate through the object pages
	var i int
	fmt.Println("Buckets:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, b := range page.Buckets {
			fmt.Printf("Bucket:%v, %v, %v\n", oss.ToString(b.Name), oss.ToString(b.StorageClass), oss.ToString(b.Location))
		}
	}
}
