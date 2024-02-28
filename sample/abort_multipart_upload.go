package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

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

	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithEndpoint(endpoint)

	client := oss.NewClient(cfg)

	initRequest := &oss.InitiateMultipartUploadRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}

	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	if err != nil {
		log.Fatalf("failed to initiate multi part upload %v", err)
	}

	partRequest := &oss.UploadPartRequest{
		Bucket:     oss.Ptr(bucketName),
		Key:        oss.Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   oss.Ptr(*initResult.UploadId),
		Body:       strings.NewReader("hi upload part request"),
	}
	var parts []oss.UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	if err != nil {
		log.Fatalf("failed to upload part %v", err)
	}
	part := oss.UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	request := &oss.AbortMultipartUploadRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
		UploadId: oss.Ptr(*initResult.UploadId),
	}
	result, err := client.AbortMultipartUpload(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to abort multipart upload %v", err)
	}
	log.Printf("abort multipart upload result:%#v\n", result)
}
