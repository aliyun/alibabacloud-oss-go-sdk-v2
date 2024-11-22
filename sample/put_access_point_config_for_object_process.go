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
	var (
		accountId         = "your account id"
		objectProcessName = "access point for object process name"
	)
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

	arn := "acs:fc:" + region + ":" + accountId + ":services/test-oss-fc.LATEST/functions/" + objectProcessName
	roleArn := "acs:ram::" + accountId + ":role/aliyunfcdefaultrole"
	request := &oss.PutAccessPointConfigForObjectProcessRequest{
		Bucket:                          oss.Ptr(bucketName),
		AccessPointForObjectProcessName: oss.Ptr(objectProcessName),
		PutAccessPointConfigForObjectProcessConfiguration: &oss.PutAccessPointConfigForObjectProcessConfiguration{
			ObjectProcessConfiguration: &oss.ObjectProcessConfiguration{
				AllowedFeatures: []string{"GetObject-Range"},
				TransformationConfigurations: []oss.TransformationConfiguration{
					{
						Actions: &oss.Actions{
							[]string{"GetObject"},
						},
						ContentTransformation: &oss.ContentTransformation{
							FunctionArn:           oss.Ptr(arn),
							FunctionAssumeRoleArn: oss.Ptr(roleArn),
						},
					},
				},
			},
			PublicAccessBlockConfiguration: &oss.PublicAccessBlockConfiguration{
				oss.Ptr(true),
			},
		},
	}
	result, err := client.PutAccessPointConfigForObjectProcess(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put access point config for object process %v", err)
	}

	log.Printf("put access point config for object process result:%#v\n", result)
}