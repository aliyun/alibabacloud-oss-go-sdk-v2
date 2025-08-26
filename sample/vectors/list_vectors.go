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
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	if len(uid) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, uid required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithUserId(uid)

	client := vectors.NewVectorsClient(cfg)

	request := &vectors.ListVectorsRequest{
		Bucket:         oss.Ptr(bucketName),
		IndexName:      oss.Ptr("index"),
		MaxResults:     oss.Ptr(100),
		NextToken:      oss.Ptr("123"),
		ReturnMetadata: oss.Ptr(true),
		ReturnData:     oss.Ptr(false),
		SegmentCount:   oss.Ptr(int(10)),
		SegmentIndex:   oss.Ptr(3),
	}

	p := client.NewListVectorsPaginator(request)

	var i int
	log.Println("Vectors:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		for _, v := range page.Vectors {
			log.Printf("vector:%v\n", v)
		}
	}
}
