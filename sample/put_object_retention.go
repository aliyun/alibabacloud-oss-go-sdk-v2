package main

import (
	"context"
	"flag"
	"log"
	"time"

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

	date := time.Now().UTC().Add(5 * time.Hour).Format("2006-01-02T15:04:05.000Z")
	putRequest := &oss.PutObjectRetentionRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
		Retention: &oss.Retention{
			Mode:            oss.Ptr("COMPLIANCE"),
			RetainUntilDate: oss.Ptr(date),
		},
	}
	putResult, err := client.PutObjectRetention(context.TODO(), putRequest)
	if err != nil {
		log.Fatalf("failed to put object retention %v", err)
	}
	log.Printf("put object retention result:%#v\n", putResult)
}
