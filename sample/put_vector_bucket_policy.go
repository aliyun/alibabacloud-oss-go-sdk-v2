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
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the vector bucket.")
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

	client := oss.NewVectorsClient(cfg)

	request := &oss.PutBucketPolicyRequest{
		Bucket: oss.Ptr(bucketName),
		Body: strings.NewReader(`{
			   "Version":"1",
			   "Statement":[
			   {
				 "Action":[
				   "oss:PutVectors",
				   "oss:GetVectors"
				],
				"Effect":"Deny",
				"Principal":["1234567890"],
				"Resource":["acs:ossvector:cn-hangzhou:1234567890:*"]
			   }
			  ]
			 }`),
	}
	result, err := client.PutBucketPolicy(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put vector bucket policy %v", err)
	}

	log.Printf("put vector bucket policy result:%#v\n", result)
}
