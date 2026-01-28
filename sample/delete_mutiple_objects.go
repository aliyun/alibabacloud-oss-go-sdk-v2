package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
	objects    string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objects, "objects", "", "The name of the objects.")
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

	if len(objects) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, objects name required")
	}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(cfg)
	var DeleteObjects []oss.DeleteObject
	objectSlice := strings.Split(objects, ",")
	for _, name := range objectSlice {
		DeleteObjects = append(DeleteObjects, oss.DeleteObject{Key: oss.Ptr(name)})
	}
	request := &oss.DeleteMultipleObjectsRequest{
		Bucket: oss.Ptr(bucketName),
		Delete: &oss.Delete{
			Objects: DeleteObjects,
		},
	}
	result, err := client.DeleteMultipleObjects(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to delete multiple objects %v", err)
	}
	log.Printf("delete multiple objects result:%#v\n", result)
}
