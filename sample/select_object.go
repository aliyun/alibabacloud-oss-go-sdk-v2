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

	request := &oss.SelectObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
		SelectRequest: &oss.SelectRequest{
			Expression: oss.Ptr("select name from ossobject"),
			InputSerializationSelect: oss.InputSerializationSelect{
				CsvBodyInput: &oss.CSVSelectInput{
					FileHeaderInfo: oss.Ptr("Use"),
				},
			},
			OutputSerializationSelect: oss.OutputSerializationSelect{
				OutputHeader: oss.Ptr(true),
			},
		},
	}

	result, err := client.SelectObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to select object %v", err)
	}
	log.Printf("select object result:%#v\n", result)
}
