package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/tables"
)

var (
	region    string
	bucketArn string
	nameSpace string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketArn, "bucket-arn", "", "The arn of the table bucket.")
	flag.StringVar(&nameSpace, "name-space", "", "The name of the name space.")
}

func main() {
	flag.Parse()
	if len(bucketArn) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket arn required")
	}

	if len(nameSpace) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, name space required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	p := client.NewListTablesPaginator(&tables.ListTablesRequest{
		TableBucketARN: oss.Ptr(bucketArn),
		Namespace: oss.Ptr(nameSpace),
	})

	var i int
	log.Println("Tables:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		for _, b := range page.Tables {
			log.Printf("Table %v,%v,%v,%v,%v,%v\n", oss.ToString(b.CreatedAt), oss.ToString(b.TableARN), oss.ToString(b.ModifiedAt), oss.ToString(b.Type), oss.ToString(b.Name), b.Namespace)
		}
	}
}
