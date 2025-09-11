package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/vectors"
)

var (
	region     string
	bucketName string
	accountId  string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the vector bucket.")
	flag.StringVar(&accountId, "account-id", "", "The id of vector account.")
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

	if len(accountId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, accounId required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithAccountId(accountId)

	client := vectors.NewVectorsClient(cfg)

	request := &vectors.PutVectorsRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr("exampleIndex"),
		Vectors: []map[string]any{
			{
				"key": "vector1",
				"data": map[string]any{
					"float32": []float32{1.2, 2.5, 3},
				},
				"metadata": map[string]any{
					"Key1": "value2",
					"Key2": []string{"1", "2", "3"},
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
