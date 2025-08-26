package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/vectors"
)

var (
	region     string
	bucketName string
	uid        string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the vector bucket.")
	flag.StringVar(&uid, "uid", "", "The id of vector account.")
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

	if len(uid) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, uid required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithUserId(uid)

	client := vectors.NewVectorsClient(cfg)

	request := &vectors.PutBucketPolicyRequest{
		Bucket: oss.Ptr(bucketName),
		Body: strings.NewReader(`{
			   "Version":"1",
			   "Statement":[
			   {
				 "Action":[
				   "ossvector:PutVectors",
				   "ossvector:GetVectors"
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
