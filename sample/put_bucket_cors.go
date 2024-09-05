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

	request := &oss.PutBucketCorsRequest{
		Bucket: oss.Ptr(bucketName),
		CORSConfiguration: &oss.CORSConfiguration{
			CORSRules: []oss.CORSRule{
				{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"PUT", "GET"},
					AllowedHeaders: []string{"Authorization"},
				},
				{
					AllowedOrigins: []string{"http://example.com", "http://example.net"},
					AllowedMethods: []string{"GET"},
					AllowedHeaders: []string{"Authorization"},
					ExposeHeaders:  []string{"x-oss-test", "x-oss-test1"},
					MaxAgeSeconds:  oss.Ptr(int64(100)),
				},
			},
			ResponseVary: oss.Ptr(false),
		},
	}
	result, err := client.PutBucketCors(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket cors %v", err)
	}

	log.Printf("put bucket cors result:%#v\n", result)
}
