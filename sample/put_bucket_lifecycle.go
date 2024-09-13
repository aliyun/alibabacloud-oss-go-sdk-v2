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

	request := &oss.PutBucketLifecycleRequest{
		Bucket: oss.Ptr(bucketName),
		LifecycleConfiguration: &oss.LifecycleConfiguration{
			Rules: []oss.LifecycleRule{
				{
					Status: oss.Ptr("Enabled"),
					ID:     oss.Ptr("rule"),
					Prefix: oss.Ptr("log/"),
					Transitions: []oss.LifecycleRuleTransition{
						{
							Days:         oss.Ptr(int32(30)),
							StorageClass: oss.StorageClassIA,
						},
					},
				},
			},
		},
	}
	result, err := client.PutBucketLifecycle(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket lifecycle %v", err)
	}

	log.Printf("put bucket lifecycle result:%#v\n", result)
}
