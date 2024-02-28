package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	endpoint   string
	bucketName string
	maxUploads int
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&endpoint, "endpoint", "", "The domain names that other services can use to access OSS.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.IntVar(&maxUploads, "max-uploads", 0, "[optional]The maximum number of `keys per page` to retrieve at once.")
}

func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(endpoint) == 0 {
		endpoint = fmt.Sprintf("oss-%v.aliyuncs.com", region)
	}

	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).
		WithEndpoint(endpoint)

	client := oss.NewClient(cfg)

	request := &oss.ListMultipartUploadsRequest{
		Bucket: oss.Ptr(bucketName),
	}

	if maxUploads > 0 {
		request.MaxUploads = int32(maxUploads)
	}

	p := client.NewListMultipartUploadsPaginator(request)

	var i int
	fmt.Println("List Multipart Uploads:")
	for p.HasNext() {
		i++
		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}
		// Log the objects found
		for _, u := range page.Uploads {
			fmt.Printf("Upload key:%v,upload id:%v, initiated:%v\n", oss.ToString(u.Key), oss.ToString(u.UploadId), oss.ToTime(u.Initiated))
		}
	}
}
