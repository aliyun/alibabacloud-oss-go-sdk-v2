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
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The `name` of the bucket.")
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

	request := &oss.ListObjectVersionsRequest{
		Bucket: oss.Ptr(bucketName),
	}
	p := client.NewListObjectVersionsPaginator(request)

	var i int
	log.Println("Object Versions:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, obj := range page.ObjectVersions {
			log.Printf("Object:%v, VersionId:%v, IsLatest:%v, Size:%v, ETag:%v, Storage Class:%v,  Last Modified:%v\n", oss.ToString(obj.Key), oss.ToString(obj.VersionId), obj.IsLatest, obj.Size, oss.ToString(obj.ETag), oss.ToString(obj.StorageClass), oss.ToTime(obj.LastModified))
		}
	}
}
