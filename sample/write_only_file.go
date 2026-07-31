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
	objectName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
}

func main() {
	flag.Parse()
	if len(bucketName) == 0 || len(region) == 0 || len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region/bucket/object required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(cfg)

	f, err := client.NewWriteOnlyFile(context.TODO(), bucketName, objectName, func(o *oss.WriteOnlyOptions) {
		o.PartSize = 100 * 1024
		o.ParallelNum = 3
	})
	if err != nil {
		log.Fatalf("failed to open write-only file %v", err)
	}

	// write several blocks; poll the checkpoint to see durable progress
	for i := 0; i < 8; i++ {
		block := []byte(strings.Repeat("x", 50*1024))
		if _, err = f.Write(block); err != nil {
			log.Fatalf("write failed: %v", err)
		}
		cp := f.StatCheckpoint()
		if cp.LastError != nil {
			log.Fatalf("background upload error: %v", cp.LastError)
		}
		pos, _ := f.Seek(0, 1) // io.SeekCurrent
		log.Printf("written cursor=%d durable offset=%d uploadId=%s", pos, cp.Offset, cp.UploadId)
	}

	if err = f.Close(); err != nil {
		// Close does not auto-abort on failure; parts are kept for resume.
		log.Fatalf("close(commit) failed: %v", err)
	}
	log.Printf("write-only file committed: oss://%s/%s", bucketName, objectName)
}
