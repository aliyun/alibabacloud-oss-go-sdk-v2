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
		WithRegion(region).WithUid(uid)

	client := oss.NewVectorsClient(cfg)

	request := &oss.DeleteBucketPolicyRequest{
		Bucket: oss.Ptr(bucketName),
	}
	result, err := client.DeleteBucketPolicy(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to delete vector bucket policy %v", err)
	}

	log.Printf("delete vector bucket policy result:%#v\n", result)
}
