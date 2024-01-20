package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	endpoint   string
	bucketName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The `name` of the bucket.")
}

// a example of showing how to get the bucket info.
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
	request := &oss.GetBucketInfoRequest{
		Bucket: oss.Ptr(bucketName),
	}

	// Send request
	result, err := client.GetBucketInfo(context.TODO(), request)

	if err != nil {
		log.Fatalf("failed to get bucket info %v", err)
	}

	// Print the result
	out, _ := json.MarshalIndent(result.BucketInfo, "", "  ")
	log.Printf("Result:\n%v", string(out))
}
