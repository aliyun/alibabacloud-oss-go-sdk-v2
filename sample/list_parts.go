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
	region     string
	endpoint   string
	bucketName string
	objectName string
	uploadId   string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
	flag.StringVar(&uploadId, "upload-id", "", "The upload id of the object.")
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

	if len(uploadId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, upload id required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithEndpoint(endpoint)

	client := oss.NewClient(cfg)

	request := &oss.ListPartsRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
		UploadId: oss.Ptr(uploadId),
	}

	p := client.NewListPartsPaginator(request)

	var i int
	fmt.Println("List Parts:")
	for p.HasNext() {
		i++
		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}
		for _, part := range page.Parts {
			fmt.Printf("Part Number:%v, ETag:%v, Last Modified:%v, Size:%v, HashCRC64:%v\n", part.PartNumber, oss.ToString(part.ETag), oss.ToTime(part.LastModified), part.Size, oss.ToString(part.HashCRC64))
		}
	}
}
