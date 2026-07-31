package main

import (
	"context"
	"errors"
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

	partSize := int64(100 * 1024)
	full := []byte(strings.Repeat("y", int(partSize)*4)) // 4 parts of source data

	// session 1: write 2 parts, snapshot the checkpoint, then simulate a crash
	f1, err := client.NewWriteOnlyFile(context.TODO(), bucketName, objectName, func(o *oss.WriteOnlyOptions) {
		o.PartSize = partSize
		o.ParallelNum = 1
	})
	if err != nil {
		log.Fatalf("open failed: %v", err)
	}
	if _, err = f1.Write(full[:partSize*2]); err != nil {
		log.Fatalf("write failed: %v", err)
	}
	if err = f1.Flush(); err != nil {
		log.Fatalf("flush failed: %v", err)
	}
	cp := f1.StatCheckpoint() // persist cp + bucket/key somewhere durable in real usage
	log.Printf("checkpoint before crash: offset=%d uploadId=%s", cp.Offset, cp.UploadId)
	// (crash: f1 is abandoned without Close)

	// session 2: resume from the checkpoint and finish
	f2, err := client.OpenWriteOnlyFile(context.TODO(), bucketName, objectName, cp, func(o *oss.WriteOnlyOptions) {
		o.ParallelNum = 1
	})
	if errors.Is(err, oss.ErrCheckpointInvalid) {
		// the checkpoint no longer matches server state (upload gone, parts
		// changed, offset/crc drift): discard it and start a fresh upload.
		log.Printf("checkpoint rejected (%v); restarting from scratch", err)
		f2, err = client.NewWriteOnlyFile(context.TODO(), bucketName, objectName, func(o *oss.WriteOnlyOptions) {
			o.PartSize = partSize
			o.ParallelNum = 1
		})
		if err != nil {
			log.Fatalf("fresh open failed: %v", err)
		}
		if _, err = f2.Write(full); err != nil {
			log.Fatalf("fresh write failed: %v", err)
		}
		if err = f2.Close(); err != nil {
			log.Fatalf("fresh close failed: %v", err)
		}
		log.Printf("restarted and committed: oss://%s/%s", bucketName, objectName)
		return
	}
	if err != nil {
		log.Fatalf("resume open failed: %v", err)
	}
	resumed := f2.StatCheckpoint()
	log.Printf("resumed at offset=%d", resumed.Offset)
	if _, err = f2.Write(full[resumed.Offset:]); err != nil {
		log.Fatalf("replay write failed: %v", err)
	}
	if err = f2.Close(); err != nil {
		log.Fatalf("close(commit) failed: %v", err)
	}
	log.Printf("resume committed: oss://%s/%s", bucketName, objectName)
}
