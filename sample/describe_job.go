package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region string
	jobId  string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&jobId, "job-id", "", "The id of the job.")
}

func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	request := &oss.DescribeJobRequest{
		BatchJobId: oss.Ptr(jobId),
	}

	result, err := client.DescribeJob(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to describe job %v", err)
	}
	log.Printf("describe job result: %v\n", result.Job)
}
