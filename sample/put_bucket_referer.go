package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
}

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

	client := oss.NewClient(cfg)

	request := &oss.PutBucketRefererRequest{
		Bucket: oss.Ptr(bucketName),
		RefererConfiguration: &oss.RefererConfiguration{
			AllowEmptyReferer: oss.Ptr(false),
			RefererList: &oss.RefererList{
				Referers: []string{""},
			},
		},
	}
	result, err := client.PutBucketReferer(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket referer %v", err)
	}

	log.Printf("put bucket referer result:%#v\n", result)
}
