package main

import (
	"context"
	"flag"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"log"
)

var (
	region     string
	bucketName string
	accountId  string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&accountId, "account-id", "", "The ID of the account that creates the job.")
}

func main() {
	flag.Parse()
	if len(bucketName) == 0 || len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region/bucket required")
	}

	if len(accountId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, account-id required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(cfg)

	result, err := client.CreateJob(context.TODO(), &oss.CreateJobRequest{
		CreateJob: &oss.CreateJobConfig{
			Operation: &oss.Operation{
				PutObjectAcl: &oss.PutObjectAcl{
					ObjectAcl: oss.ObjectACLPublicRead,
				},
			},
			ClientRequestToken: oss.Ptr("1234567890"),
			KeyPrefixManifestGenerator: &oss.KeyPrefixManifestGenerator{
				SourceBucket: oss.Ptr(bucketName),
				Prefix:       oss.Ptr("batch-manifests/"),
			},
			Report: &oss.Report{
				Bucket:      oss.Ptr(bucketName),
				Enabled:     oss.Ptr(true),
				Prefix:      oss.Ptr("batch-reports/"),
				ReportScope: oss.Ptr("AllTasks"),
			},
			Priority: oss.Ptr(int32(10)),
			RoleArn:  oss.Ptr("acs:ram::" + accountId + ":role/oss-sdk-batch-test"),
		},
	})
	if err != nil {
		log.Fatalf("failed to create job, %v", err)
	}
	log.Printf("create job result:%#v\n", result)
}
