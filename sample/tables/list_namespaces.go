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
	region         string
	tableBucketArn string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&tableBucketArn, "table-bucket-arn", "", "The arn of the table bucket.")
}

func main() {
	flag.Parse()
	if len(tableBucketArn) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, table bucket arn required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := tables.NewTablesClient(cfg)

	p := client.NewListNameSpacesPaginator(&tables.ListNamespacesRequest{
		TableBucketARN: oss.Ptr(tableBucketArn),
	})

	var i int
	log.Println("Namespaces:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		for _, b := range page.Namespaces {
			log.Printf("Namespace %v,%v,%v,%v,%v,%v\n", oss.ToString(b.CreatedAt), oss.ToString(b.CreatedBy), oss.ToString(b.OwnerAccountId), oss.ToString(b.TableBucketId), oss.ToString(b.NamespaceId), b.Namespace)
		}
	}
}
