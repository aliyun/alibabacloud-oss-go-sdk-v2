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

	request := &oss.PutBucketTagsRequest{
		Bucket: oss.Ptr(bucketName),
		Tagging: &oss.Tagging{
			&oss.TagSet{
				[]oss.Tag{
					{
						oss.Ptr("k1"),
						oss.Ptr("v1"),
					},
				},
			},
		},
	}
	result, err := client.PutBucketTags(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put vector bucket tags %v", err)
	}

	log.Printf("put vector bucket tags result:%#v\n", result)
}
