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
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
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

	// case 1:meta search
	request := &oss.DoMetaQueryRequest{
		Bucket: oss.Ptr(bucketName),
		MetaQuery: &oss.MetaQuery{
			Query: oss.Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
			Sort:  oss.Ptr("Size"),
			Order: oss.Ptr(oss.MetaQueryOrderAsc),
			Aggregations: &oss.MetaQueryAggregations{
				[]oss.MetaQueryAggregation{
					{
						Field:     oss.Ptr("Size"),
						Operation: oss.Ptr("sum"),
					},
					{
						Field:     oss.Ptr("Size"),
						Operation: oss.Ptr("max"),
					},
				},
			},
		},
	}
	result, err := client.DoMetaQuery(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to do meta query %v", err)
	}

	log.Printf("do meta query result:%#v\n", result)

	// case 2: ai search
	//request = &oss.DoMetaQueryRequest{
	//	Bucket: oss.Ptr(bucketName),
	//	Mode:   oss.Ptr("semantic"),
	//	MetaQuery: &oss.MetaQuery{
	//		MaxResults:  oss.Ptr(int64(99)),
	//		Query:       oss.Ptr("Overlook the snow-covered forest"),
	//		MediaType:   oss.Ptr("image"),
	//		SimpleQuery: oss.Ptr(`{"Operation":"gt", "Field": "Size", "Value": "30"}`),
	//	},
	//}
	//result, err = client.DoMetaQuery(context.TODO(), request)
	//if err != nil {
	//	log.Fatalf("failed to do meta query %v", err)
	//}
	//
	//log.Printf("do meta query result:%#v\n", result)
}
