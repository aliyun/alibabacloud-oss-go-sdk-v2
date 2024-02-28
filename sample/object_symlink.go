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
	region      string
	endpoint    string
	bucketName  string
	objectName  string
	symlinkName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
	flag.StringVar(&symlinkName, "symlink", "", "The destination object to which the symbolic link points. ")
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
	if len(symlinkName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, symlink name required")
	}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithEndpoint(endpoint)
	client := oss.NewClient(cfg)
	putRequest := &oss.PutSymlinkRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
		Target: oss.Ptr(symlinkName),
	}
	putResult, err := client.PutSymlink(context.TODO(), putRequest)
	if err != nil {
		log.Fatalf("failed to put symlink %v", err)
	}
	log.Printf("put symlink result:%#v\n", putResult)
	getRequest := &oss.GetSymlinkRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	getResult, err := client.GetSymlink(context.TODO(), getRequest)
	if err != nil {
		log.Fatalf("failed to get symlink %v", err)
	}
	log.Printf("get symlink result:%#v\n", getResult)
}
