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
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the vector bucket.")
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

	client := oss.NewVectorsClient(cfg)

	request := &oss.PutVectorIndexRequest{
		Bucket:         oss.Ptr(bucketName),
		DataType:       oss.Ptr("string"),
		Dimension:      oss.Ptr(128),
		DistanceMetric: oss.Ptr("cosine"),
		IndexName:      oss.Ptr("exampleIndex"),
		Metadata: map[string]interface{}{
			"nonFilterableMetadataKeys": []string{"foo", "bar"},
		},
	}
	result, err := client.PutVectorIndex(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put vectors %v", err)
	}
	log.Printf("put vectors result:%#v\n", result)
}
