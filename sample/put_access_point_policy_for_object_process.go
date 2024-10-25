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

	policy := `{"Version":"1","Statement":[{"Action":["oss:GetObject"],"Effect":"Allow","Principal":["` + accountId + `"],"Resource":["acs:oss:` + region + `:` + accountId + `:accesspointforobjectprocess/` + objectProcessName + `/object/*"]}]}`
	request := &oss.PutAccessPointPolicyForObjectProcessRequest{
		Bucket:                          oss.Ptr(bucketName),
		AccessPointForObjectProcessName: oss.Ptr(objectProcessName),
		Body:                            strings.NewReader(policy),
	}
	result, err := client.PutAccessPointPolicyForObjectProcess(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put access point policy for object process %v", err)
	}

	log.Printf("put access point policy for object process result:%#v\n", result)
}
