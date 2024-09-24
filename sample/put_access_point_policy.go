package main

import (
	"context"
	"flag"
	"log"
	"strings"

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
		accessPointName = "access point name"
		policy          = `{
		   "Version":"1",
		   "Statement":[
		   {
			 "Action":[
			   "oss:PutObject",
			   "oss:GetObject"
			],
			"Effect":"Deny",
			"Principal":["27737962156157xxxx"],
			"Resource":[
			   "acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01",
			   "acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01/object/*"
			 ]
		   }
		  ]
		 }`
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

	request := &oss.PutAccessPointPolicyRequest{
		Bucket:          oss.Ptr(bucketName),
		AccessPointName: oss.Ptr(accessPointName),
		Body:            strings.NewReader(policy),
	}
	result, err := client.PutAccessPointPolicy(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put access point policy %v", result)
	}

	log.Printf("put access point policy result:%#v\n", result)
}
