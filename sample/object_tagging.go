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
	region     string
	endpoint   string
	bucketName string
	objectName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
}

func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}
	if len(endpoint) == 0 {
		endpoint = fmt.Sprintf("oss-%v.aliyuncs.com", region)
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
		WithRegion(region).
		WithEndpoint(endpoint)
	client := oss.NewClient(cfg)
	putRequest := &oss.PutObjectTaggingRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
		Tagging: &oss.Tagging{
			TagSet: oss.TagSet{
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
	getRequest := &oss.GetObjectTaggingRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	getResult, err := client.GetObjectTagging(context.TODO(), getRequest)
	if err != nil {
		log.Fatalf("failed to get object tagging %v", err)
	}
	log.Printf("get object tagging result:%#v\n", getResult)
}
