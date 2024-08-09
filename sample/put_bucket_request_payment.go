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
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
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

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	putRequest := &oss.PutBucketRequestPaymentRequest{
		Bucket: oss.Ptr(bucketName),
		PaymentConfiguration: &oss.RequestPaymentConfiguration{
			Payer: oss.Requester,
		},
	}
	putResult, err := client.PutBucketRequestPayment(context.TODO(), putRequest)
	if err != nil {
		log.Fatalf("failed to put bucket request payment %v", err)
	}

	log.Printf("put bucket request payment result:%#v\n", putResult)
}
