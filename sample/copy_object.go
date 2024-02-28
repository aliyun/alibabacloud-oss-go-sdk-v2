package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"log"
)

var (
	region         string
	endpoint       string
	bucketName     string
	objectName     string
	destBucketName string
	destObjectName string
	uploadType     string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "src-object", "", "The name of the source object.")
	flag.StringVar(&destBucketName, "dest-bucket", "", "The name of the destination bucket.")
	flag.StringVar(&destObjectName, "dest-object", "", "The name of the destination object.")
	flag.StringVar(&uploadType, "type", "", "The upload type of the object.")
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

	if len(destBucketName) == 0 {
		destBucketName = bucketName
	}

	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, src object name required")
	}

	if len(destObjectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, destination object name required")
	}

	if uploadType != "copy" && uploadType != "multi" {
		log.Fatalf("invalid parameters, upload type value in the optional value:copy|multi")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithEndpoint(endpoint)

	client := oss.NewClient(cfg)

	switch uploadType {
	case "copy":
		request := &oss.CopyObjectRequest{
			Bucket:       oss.Ptr(destBucketName),
			Key:          oss.Ptr(destObjectName),
			SourceKey:    oss.Ptr(objectName),
			SourceBucket: oss.Ptr(bucketName),
		}
		result, err := client.CopyObject(context.TODO(), request)
		if err != nil {
			log.Fatalf("failed to copy object %v", err)
		}
		log.Printf("copy object result:%#v\n", result)
	case "multi":
		initRequest := &oss.InitiateMultipartUploadRequest{
			Bucket: oss.Ptr(bucketName),
			Key:    oss.Ptr(destObjectName),
		}

		initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
		if err != nil {
			log.Fatalf("failed to initiate multi part upload %v", err)
		}

		partRequest := &oss.UploadPartCopyRequest{
			Bucket:     oss.Ptr(bucketName),
			Key:        oss.Ptr(destObjectName),
			PartNumber: int32(1),
			UploadId:   oss.Ptr(*initResult.UploadId),
			SourceKey:  oss.Ptr(objectName),
		}
		var parts []oss.UploadPart
		partResult, err := client.UploadPartCopy(context.TODO(), partRequest)
		if err != nil {
			log.Fatalf("failed to upload part copy 1 %v", err)
		}
		part := oss.UploadPart{
			PartNumber: partRequest.PartNumber,
			ETag:       partResult.ETag,
		}
		parts = append(parts, part)
		request := &oss.CompleteMultipartUploadRequest{
			Bucket:   oss.Ptr(bucketName),
			Key:      oss.Ptr(destObjectName),
			UploadId: oss.Ptr(*initResult.UploadId),
			CompleteMultipartUpload: &oss.CompleteMultipartUpload{
				Parts: parts,
			},
		}
		result, err := client.CompleteMultipartUpload(context.TODO(), request)
		if err != nil {
			log.Fatalf("failed to complete multipart upload %v", err)
		}
		log.Printf("complete multipart upload result:%#v\n", result)
	}
}
