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
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
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

	if len(endpoint) == 0 {
		endpoint = fmt.Sprintf("oss-%v.aliyuncs.com", region)
	}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithEndpoint(endpoint)

	client := oss.NewClient(cfg)

	// Set the request
	putRequest := &oss.PutBucketAclRequest{
		Bucket: oss.Ptr(bucketName),
		Acl:    oss.BucketACLPrivate,
	}
	putResult, err := client.PutBucketAcl(context.TODO(), putRequest)
	if err != nil {
		log.Fatalf("failed to put bucket acl %v", err)
	}
	log.Printf("put bucket acl result:%#v\n", putResult)

	getRequest := &oss.GetBucketAclRequest{
		Bucket: oss.Ptr(bucketName),
	}
	getResult, err := client.GetBucketAcl(context.TODO(), getRequest)
	if err != nil {
		log.Fatalf("failed to put bucket acl %v", err)
	}
	log.Printf("get bucket acl result:%#v\n", getResult)

	log.Printf("get bucket acl:%#v\n", oss.ToString(getResult.ACL))
}
