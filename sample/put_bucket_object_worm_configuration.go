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

	request := &oss.PutBucketObjectWormConfigurationRequest{
		Bucket: oss.Ptr(bucketName),
		ObjectWormConfiguration: &oss.ObjectWormConfiguration{
			ObjectWormEnabled: oss.Ptr("Enabled"),
			Rule: &oss.ObjectWormConfigurationRule{
				DefaultRetention: &oss.ObjectWormConfigurationRuleDefaultRetention{
					Mode: oss.Ptr("COMPLIANCE"),
					Days: oss.Ptr(int32(1)),
				},
			},
		},
	}

	result, err := client.PutBucketObjectWormConfiguration(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket object worm configuration %v", err)
	}
	log.Printf("put bucket object worm configuration result:%#v\n", result)
}
