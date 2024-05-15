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
		uploadId = "your upload id"
	)
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}
	if len(uploadId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, upload id required")
	}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)
	putObjRequest := &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	putObjResult, err := client.Presign(context.TODO(), putObjRequest)
	if err != nil {
		log.Fatalf("failed to put object presign %v", err)
	}
	log.Printf("put object presign result:%#v\n", putObjResult)
	log.Printf("put object url:%#v\n", putObjResult.URL)

	getObjRequest := &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	getObjResult, err := client.Presign(context.TODO(), getObjRequest)
	if err != nil {
		log.Fatalf("failed to get object presign %v", err)
	}
	log.Printf("get object presign result:%#v\n", getObjResult)
	log.Printf("get object url:%#v\n", getObjResult.URL)

	headObjRequest := &oss.HeadObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	headObjResult, err := client.Presign(context.TODO(), headObjRequest)
	if err != nil {
		log.Fatalf("failed to head object presign %v", err)
	}
	log.Printf("head object presign result:%#v\n", headObjResult)
	log.Printf("head object url:%#v\n", headObjResult.URL)

	initObjRequest := &oss.InitiateMultipartUploadRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	initObjResult, err := client.Presign(context.TODO(), initObjRequest)
	if err != nil {
		log.Fatalf("failed to initiate multipart upload object presign %v", err)
	}
	log.Printf("initiate multipart upload result:%#v\n", initObjResult)
	log.Printf("initiate multipart upload url:%#v\n", headObjResult.URL)

	partObjRequest := &oss.UploadPartRequest{
		Bucket:     oss.Ptr(bucketName),
		Key:        oss.Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   oss.Ptr(uploadId),
	}
	partObjResult, err := client.Presign(context.TODO(), partObjRequest)
	if err != nil {
		log.Fatalf("failed to upload part presign %v", err)
	}
	log.Printf("upload part result:%#v\n", partObjResult)
	log.Printf("upload part url:%#v\n", headObjResult.URL)

	completeObjRequest := &oss.CompleteMultipartUploadRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
		UploadId: oss.Ptr(uploadId),
	}
	completeObjResult, err := client.Presign(context.TODO(), completeObjRequest)
	if err != nil {
		log.Fatalf("failed to complete multipart upload presign %v", err)
	}
	log.Printf("complete multipart upload result:%#v\n", completeObjResult)
	log.Printf("complete multipart upload result url:%#v\n", completeObjResult.URL)

	abortObjRequest := &oss.AbortMultipartUploadRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
		UploadId: oss.Ptr(uploadId),
	}
	abortObjResult, err := client.Presign(context.TODO(), abortObjRequest)
	if err != nil {
		log.Fatalf("failed to abort multipart upload presign %v", err)
	}
	log.Printf("abort multipart upload result:%#v\n", abortObjResult)
	log.Printf("abort multipart upload result url:%#v\n", abortObjResult.URL)

}
