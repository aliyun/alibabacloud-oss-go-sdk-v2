package main

import (
	"bufio"
	"context"
	"flag"
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
	bucketName string
	objectName string
	letters    = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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
	body := randBody(400000)
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
				UploadId:   oss.Ptr(uploadId),
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
	log.Println("upload part success!")
}

func randBody(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}
