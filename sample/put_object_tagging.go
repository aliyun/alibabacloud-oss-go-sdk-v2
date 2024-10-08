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
	objectName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
}

func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}
	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(cfg)
	putRequest := &oss.PutObjectTaggingRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
		Tagging: &oss.Tagging{
			TagSet: &oss.TagSet{
				Tags: []oss.Tag{
					{
						Key:   oss.Ptr("k1"),
						Value: oss.Ptr("v1"),
					},
					{
						Key:   oss.Ptr("k2"),
						Value: oss.Ptr("v2"),
					},
				},
			},
		},
	}
	putResult, err := client.PutObjectTagging(context.TODO(), putRequest)
	if err != nil {
		log.Fatalf("failed to put object tagging %v", err)
	}
	log.Printf("put object tagging result:%#v\n", putResult)
}
