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
	uid        string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the vector bucket.")
	flag.StringVar(&uid, "uid", "", "The id of vector account.")
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

	if len(uid) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, uid required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithUserId(uid)

	client := vectors.NewVectorsClient(cfg)

	request := &vectors.QueryVectorsRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr("index"),
		Filter:    oss.Ptr(`{"$and":[{"type":{"$in":["comedy","documentary"]}},{"year":{"$gte":2020}}]}`),
		QueryVector: map[string]interface{}{
			"float32": []float32{float32(32)},
		},
		ReturnMetadata: oss.Ptr(true),
		ReturnDistance: oss.Ptr(true),
		TopK:           oss.Ptr(10),
	}
	result, err := client.QueryVectors(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to query vectors %v", err)
	}
	log.Printf("query vectors result:%#v\n", result)
}
