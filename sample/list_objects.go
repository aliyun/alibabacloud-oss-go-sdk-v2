package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region          string
	endpoint        string
	bucketName      string
	objectPrefix    string
	objectDelimiter string
	maxKeys         int
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The `name` of the bucket.")
	flag.StringVar(&objectPrefix, "prefix", "", "[optional]`object prefix` of the keys to list.")
	flag.StringVar(&objectDelimiter, "delimiter", "",
		"[optional]`object key delimiter` used by List objects to group object keys.")
	flag.IntVar(&maxKeys, "max-keys", 0,
		"[optional]The maximum number of `keys per page` to retrieve at once.")
}

// Lists all objects in a bucket using paginator
func main() {
	flag.Parse()
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

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
	request := &oss.ListObjectsRequest{
		Bucket: oss.Ptr(bucketName),
	}

	if len(objectPrefix) > 0 {
		request.Prefix = oss.Ptr(objectPrefix)
	}

	if len(objectDelimiter) > 0 {
		request.Delimiter = oss.Ptr(objectDelimiter)
	}

	if maxKeys > 0 {
		request.MaxKeys = int32(maxKeys)
	}

	// Create the Paginator for the ListObjects operation.
	p := client.NewListObjectsPaginator(request)

	// Iterate through the object pages
	var i int
	log.Println("Objects:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, obj := range page.Contents {
			log.Printf("Object:%v, %v, %v\n", oss.ToString(obj.Key), obj.Size, oss.ToTime(obj.LastModified))
		}
	}
}
