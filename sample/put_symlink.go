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
	var (
		symlinkName = "your symlink name"
	)
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
	if len(symlinkName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, symlink name required")
	}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
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
}
