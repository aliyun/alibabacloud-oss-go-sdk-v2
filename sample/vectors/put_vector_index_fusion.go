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
	accountId  string
	indexName  string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the vector bucket.")
	flag.StringVar(&accountId, "account-id", "", "The id of vector account.")
	flag.StringVar(&indexName, "index", "", "The name of vector index.")
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

	if len(accountId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, accounId required")
	}

	if len(indexName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, index required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithAccountId(accountId)

	client := vectors.NewVectorsClient(cfg)

	request := &vectors.PutVectorIndexFusionRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
		Mode:      oss.Ptr("fusion"),
		SchemaConfiguration: &vectors.SchemaConfiguration{
			Fields: []map[string]any{
				vectors.FieldSchema{
					Name:           oss.Ptr("vector_1"),
					Type:           oss.Ptr("vector"),
					DataType:       oss.Ptr("float32"),
					Dimension:      oss.Ptr(1024),
					DistanceMetric: oss.Ptr("cosine"),
				}.ToMap(),
				vectors.FieldSchema{
					Name:           oss.Ptr("vector_2"),
					Type:           oss.Ptr("vector"),
					DataType:       oss.Ptr("float32"),
					Dimension:      oss.Ptr(512),
					DistanceMetric: oss.Ptr("cosine"),
				}.ToMap(),
				vectors.FieldSchema{
					Name:    oss.Ptr("timestamps"),
					Type:    oss.Ptr("long"),
					IsArray: oss.Ptr(true),
				}.ToMap(),
				vectors.FieldSchema{
					Name: oss.Ptr("price"),
					Type: oss.Ptr("double"),
				}.ToMap(),
				vectors.FieldSchema{
					Name: oss.Ptr("ip"),
					Type: oss.Ptr("ip"),
				}.ToMap(),
				vectors.FieldSchema{
					Name: oss.Ptr("location"),
					Type: oss.Ptr("geoPoint"),
				}.ToMap(),
				vectors.FieldSchema{
					Name: oss.Ptr("tag"),
					Type: oss.Ptr("string"),
				}.ToMap(),
				vectors.FieldSchema{
					Name:           oss.Ptr("user_id"),
					Type:           oss.Ptr("string"),
					IsPartitionKey: oss.Ptr(true),
				}.ToMap(),
				vectors.FieldSchema{
					Name:    oss.Ptr("tags"),
					Type:    oss.Ptr("string"),
					IsArray: oss.Ptr(true),
				}.ToMap(),
				vectors.FieldSchema{
					Name:       oss.Ptr("title_1"),
					Type:       oss.Ptr("string"),
					ExactMatch: oss.Ptr(true),
					Text: &vectors.TextSchema{
						Enabled:  oss.Ptr(true),
						Analyzer: oss.Ptr("standard"),
						AnalyzerParameters: &vectors.AnalyzerParameters{
							CaseSensitive: oss.Ptr(true),
							DelimitWord:   oss.Ptr(false),
						},
					},
				}.ToMap(),
				vectors.FieldSchema{
					Name:       oss.Ptr("title_2"),
					Type:       oss.Ptr("string"),
					ExactMatch: oss.Ptr(false),
					Text: &vectors.TextSchema{
						Enabled:  oss.Ptr(true),
						Analyzer: oss.Ptr("split"),
						AnalyzerParameters: &vectors.AnalyzerParameters{
							CaseSensitive: oss.Ptr(true),
							Delimiter:     oss.Ptr(" "),
						},
					},
				}.ToMap(),
			},
		},
	}
	result, err := client.PutVectorIndexFusion(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put vector index fusion%v", err)
	}
	log.Printf("put vector index fusion result:%#v\n", result)
}
