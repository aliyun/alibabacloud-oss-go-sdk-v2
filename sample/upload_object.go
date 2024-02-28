package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	endpoint   string
	bucketName string
	objectName string
	uploadType string
	filePath   string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
	flag.StringVar(&uploadType, "type", "", "The upload type of the object.")
	flag.StringVar(&filePath, "file", "", "The name of the file.")
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

	if len(uploadType) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, upload type required")
	}

	if uploadType != "put" && uploadType != "append" && uploadType != "multi" && uploadType != "uploader" && uploadType != "from" {
		log.Fatalf("invalid parameters, upload type value in the optional value:put|uploader|multi|append|from")
	}

	var content string
	if len(filePath) > 0 {
		if filePath == "-" {
			content, _ = readFromStdin()
		} else {
			content, _ = readFromFile(filePath)
		}
	}

	if len(content) == 0 {
		content = "hi oss"
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithEndpoint(endpoint)

	client := oss.NewClient(cfg)

	switch uploadType {
	case "put":
		request := &oss.PutObjectRequest{
			Bucket: oss.Ptr(bucketName),
			Key:    oss.Ptr(objectName),
			Body:   strings.NewReader(content),
		}

		result, err := client.PutObject(context.TODO(), request)
		if err != nil {
			log.Fatalf("failed to put object %v", err)
		}
		log.Printf("put object result:%#v\n", result)
	case "append":
		af, err := client.AppendFile(context.TODO(), bucketName, objectName)
		if err != nil {
			log.Fatalf("failed to append file %v", err)
		}
		log.Printf("append file af:%#v\n", af)
		n, err := af.Write([]byte(content))
		if err != nil {
			log.Fatalf("failed to af write %v", err)
		}
		log.Printf("af write n:%#v\n", n)
	case "multi":
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
		partSize := len(content) / count
		if partSize < 10500 {
			log.Fatalf("the file content is too small!")
		}
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
	case "uploader":
		u := oss.NewUploader(client)
		localFile := "upload.file.tmp"
		file, err := os.Create(localFile)
		if err != nil {
			log.Fatalf("failed to create file %v", err)
		}
		defer file.Close()
		_, err = io.WriteString(file, content)
		if err != nil {
			log.Fatalf("failed to write content %v", err)
		}
		result, err := u.UploadFile(context.TODO(), &oss.PutObjectRequest{
			Bucket: oss.Ptr(bucketName),
			Key:    oss.Ptr(objectName),
		}, localFile)
		os.Remove(localFile)
		if err != nil {
			log.Fatalf("failed to upload file %v", err)
		}
		log.Printf("upload file result:%#v\n", result)
	case "from":
		u := oss.NewUploader(client)
		result, err := u.UploadFrom(context.TODO(), &oss.PutObjectRequest{
			Bucket: oss.Ptr(bucketName),
			Key:    oss.Ptr(objectName),
		}, strings.NewReader(content))
		if err != nil {
			log.Fatalf("failed to upload form %v", err)
		}
		log.Printf("upload form result:%#v\n", result)
	}
}

func readFromStdin() (string, error) {
	var builder strings.Builder
	reader := bufio.NewReader(os.Stdin)
	for {
		content, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("failed to read from stdin: %v", err)
		}
		builder.WriteString(content)
		if err == io.EOF {
			break
		}
	}
	return builder.String(), nil
}

func readFromFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read from file: %v", err)
	}
	return string(content), nil
}
