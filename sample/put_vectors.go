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

	request := &oss.PutVectorsRequest{
		Bucket: oss.Ptr(bucketName),
		IndexName: oss.Ptr("exampleIndex"),
		Vectors: []map[string]interface{}{
			{
				"key": "vector1",
				"data": map[string]interface{}{
					"float32": []float32{1.2, 2.5, 3},
				},
				"metadata": map[string]interface{}{
					"Key1": 32,
					"Key2": "value2",
					"Key3": []string{"1", "2", "3"},
					"Key4": false,
				},
			},
		},
	}
	result, err := client.PutVectors(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put vectors %v", err)
	}
	log.Printf("put vectors result:%#v\n", result)
}
