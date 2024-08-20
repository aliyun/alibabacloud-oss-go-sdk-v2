package main

import (
	"context"
	"flag"
	"log"
	"strings"

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

	request := &oss.PutBucketPolicyRequest{
		Bucket: oss.Ptr(bucketName),
		Body: strings.NewReader(`{
			   "Version":"1",
			   "Statement":[
			   {
				 "Action":[
				   "oss:PutObject",
				   "oss:GetObject"
				],
				"Effect":"Deny",
				"Principal":["1234567890"],
				"Resource":["acs:oss:*:1234567890:*/*"]
			   }
			  ]
			 }`),
	}
	result, err := client.PutBucketPolicy(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket policy %v", err)
	}

	log.Printf("put bucket policy result:%#v\n", result)
}
