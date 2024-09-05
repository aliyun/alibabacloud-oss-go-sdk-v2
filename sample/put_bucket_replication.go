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
		targetBucket   = "target bucket name"
		targetLocation = "oss-cn-beijing"
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

	request := &oss.PutBucketReplicationRequest{
		Bucket: oss.Ptr(bucketName),
		ReplicationConfiguration: &oss.ReplicationConfiguration{
			[]oss.ReplicationRule{
				{
					RTC: &oss.ReplicationTimeControl{
						Status: oss.Ptr("enabled"),
					},
					Destination: &oss.ReplicationDestination{
						Bucket:       oss.Ptr(targetBucket),
						Location:     oss.Ptr(targetLocation),
						TransferType: oss.TransferTypeOssAcc,
					},
					HistoricalObjectReplication: oss.HistoricalObjectReplicationEnabled,
				},
			},
		},
	}
	result, err := client.PutBucketReplication(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put bucket replication %v", err)
	}

	log.Printf("put bucket replication result:%#v\n", result)
}
