package main

import (
	"context"
	"flag"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/dataprocess"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region      string
	bucketName  string
	datasetName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&datasetName, "dataset", "", "The name of the dataset.")
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

	if len(datasetName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, dataset name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := dataprocess.NewClient(cfg)

	request := &dataprocess.SemanticQueryRequest{
		Bucket:      oss.Ptr(bucketName),
		DatasetName: oss.Ptr(datasetName),
		Query:       oss.Ptr("Sunset by the sea"),
		MaxResults:  oss.Ptr(int32(10)),
		MediaTypes:  oss.Ptr(`["image"]`),
	}
	result, err := client.SemanticQuery(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to semantic query %v", err)
	}
	log.Printf("semantic query result:%#v\n", result)
}
