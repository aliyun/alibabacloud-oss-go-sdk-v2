package main

import (
	"context"
	"flag"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/dataprocess"
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

	client := dataprocess.NewClient(cfg)

	request := &dataprocess.ListDatasetsRequest{
		Bucket: oss.Ptr(bucketName),
	}
	p := client.NewListDatasetsPaginator(request)

	var i int
	log.Println("Datasets:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, dataset := range page.Datasets {
			log.Printf("DatasetName:%v, Description:%v, Create Time:%v, Update Time:%v, Template Id:%v, File Count:%v, Dataset Max Bind Count:%v,Dataset Max File Count:%v, Dataset Max Entity Count:%v, Dataset Max Relation Count:%v, Dataset Max Total File Size:%v, Bind Count:%v, Total File Size:%v\n", oss.ToString(dataset.DatasetName), oss.ToString(dataset.Description), oss.ToString(dataset.CreateTime), oss.ToString(dataset.UpdateTime), oss.ToString(dataset.TemplateId), oss.ToInt64(dataset.FileCount), oss.ToInt64(dataset.DatasetMaxBindCount), oss.ToInt64(dataset.DatasetMaxFileCount), oss.ToInt64(dataset.DatasetMaxEntityCount), oss.ToInt64(dataset.DatasetMaxRelationCount), oss.ToInt64(dataset.DatasetMaxTotalFileSize), oss.ToInt64(dataset.BindCount), oss.ToInt64(dataset.TotalFileSize))
		}
	}
}
