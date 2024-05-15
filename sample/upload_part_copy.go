package main

import (
	"context"
	"flag"
	"log"
	"sync"

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

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	var wg sync.WaitGroup
	var parts []oss.UploadPart
	count := 3
	var mu sync.Mutex
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(partNumber int, i int) {
			defer wg.Done()
			partRequest := &oss.UploadPartCopyRequest{
				Bucket:     oss.Ptr(bucketName),
				Key:        oss.Ptr(objectName),
				PartNumber: int32(partNumber),
				UploadId:   oss.Ptr(uploadId),
			}
			partResult, err := client.UploadPartCopy(context.TODO(), partRequest)
			if err != nil {
				log.Fatalf("failed to upload part copy %d: %v", partNumber, err)
			}
			part := oss.UploadPart{
				PartNumber: partRequest.PartNumber,
				ETag:       partResult.ETag,
			}
			mu.Lock()
			parts = append(parts, part)
			mu.Unlock()
		}(i+1, i)
	}
	wg.Wait()
	log.Println("upload part copy success!")
}
