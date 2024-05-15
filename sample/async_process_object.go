package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
	objectName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
}

func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}
	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(cfg)
	style := "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1"
	process := fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", style, strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(bucketName)), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(objectName)), "="))
	request := &oss.AsyncProcessObjectRequest{
		Bucket:       oss.Ptr(bucketName),
		Key:          oss.Ptr(objectName),
		AsyncProcess: oss.Ptr(process),
	}
	result, err := client.AsyncProcessObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to async process object %v", err)
	}
	log.Printf("async process object result:%#v\n", result)
}
