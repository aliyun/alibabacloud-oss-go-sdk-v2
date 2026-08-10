package main

import (
	"context"
	"flag"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"log"
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

	if len(jobId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, job-id required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(cfg)

	result, err := client.UpdateJobStatus(context.TODO(), &oss.UpdateJobStatusRequest{
		BatchJobId:         oss.Ptr(jobId),
		RequestedJobStatus: oss.Ptr("Cancelled"),
	})
	if err != nil {
		log.Fatalf("failed to update job status, %v", err)
	}
	log.Printf("update job status result:%#v\n", result)
}
