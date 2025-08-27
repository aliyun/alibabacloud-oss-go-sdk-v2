package main

import (
	"context"
	"flag"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/vectors"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
	uid        string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the vector bucket.")
	flag.StringVar(&uid, "uid", "", "The id of vector account.")
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

	if len(uid) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, uid required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithUserId(uid)

	client := vectors.NewVectorsClient(cfg)

	request := &vectors.PutBucketLoggingRequest{
		Bucket: oss.Ptr(bucketName),
		BucketLoggingStatus: &vectors.BucketLoggingStatus{
			&vectors.LoggingEnabled{
				TargetBucket: oss.Ptr("TargetBucket"),
				TargetPrefix: oss.Ptr("TargetPrefix"),
				LoggingRole:  oss.Ptr("AliyunOSSLoggingDefaultRole"),
			},
		},
	}
	result, err := client.PutBucketLogging(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put vector bucket logging %v", err)
	}

	log.Printf("put vector bucket logging result:%#v\n", result)
}
