package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
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
	request := &oss.ProcessObjectRequest{
		Bucket:  oss.Ptr(bucketName),
		Key:     oss.Ptr(objectName),
		Process: oss.Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte(objectName)))),
	}
	result, err := client.ProcessObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to process object %v", err)
	}
	log.Printf("process object result:%#v\n", result)
}
