package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	endpoint   string
	bucketName string
	objectName string
	letters    = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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
	var wg sync.WaitGroup
	var parts []oss.UploadPart
	count := 3
	body := randStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, _ := io.ReadAll(bufReader)
	partSize := len(body) / count
	var mu sync.Mutex
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(partNumber int, partSize int, i int) {
			defer wg.Done()
			partRequest := &oss.UploadPartRequest{
				Bucket:     oss.Ptr(bucketName),
				Key:        oss.Ptr(objectName),
				PartNumber: int32(partNumber),
				UploadId:   oss.Ptr(*initResult.UploadId),
				Body:       strings.NewReader(string(content[i*partSize : (i+1)*partSize])),
			}
			partResult, err := client.UploadPart(context.TODO(), partRequest)
			if err != nil {
				log.Fatalf("failed to upload part %d: %v", partNumber, err)
			}
			part := oss.UploadPart{
				PartNumber: partRequest.PartNumber,
				ETag:       partResult.ETag,
			}
			mu.Lock()
			parts = append(parts, part)
			mu.Unlock()
		}(i+1, partSize, i)
	}
	wg.Wait()
	request := &oss.CompleteMultipartUploadRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
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

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}
