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

	request := &vectors.QueryVectorsFusionRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
		Knn: []map[string]any{
			vectors.Knn{
				Field:         oss.Ptr("demo"),
				QueryVector:   map[string]any{"float32": []float32{float32(32)}},
				TopK:          oss.Ptr(10),
				NumCandidates: oss.Ptr(9),
				Filter: map[string]any{
					"meta_field_1": map[string]any{
						"$eq": "abc",
					},
				},
				Boost: oss.Ptr(float32(1)),
			}.ToMap(),
		},
		Retriever: map[string]any{
			"simple": vectors.SimpleRetriever{
				Query: map[string]any{
					"$and": []map[string]any{
						{"type": map[string]any{"$in": []string{"a", "b"}}},
						{"year": map[string]any{"$gte": 2020}},
					},
				},
			}.ToMap(),
		},
		ReturnMetadata:       oss.Ptr(true),
		ReturnMetadataFields: []string{"key1", "key2"},
		PartitionKeys:        []string{"key1", "key2"},
		Limit:                oss.Ptr(10),
		NextToken:            oss.Ptr("nextToken"),
		Sort: []vectors.Sort{
			{"field_a": vectors.SortOptions{Order: oss.Ptr("asc")}},
			{"_score": vectors.SortOptions{Order: oss.Ptr("desc")}},
			{"_primaryKey": vectors.SortOptions{Order: oss.Ptr("asc")}},
		},
	}
	result, err := client.QueryVectorsFusion(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to query vectors fusion: %v", err)
	}
	log.Printf("query vectors fusion result:%#v\n", result)
}
