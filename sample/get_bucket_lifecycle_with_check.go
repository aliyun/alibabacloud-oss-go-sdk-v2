package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
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

	// --- Mode 1: default (UnknownFieldsCheckFail) ---
	// When the response contains XML fields unknown to the SDK, the call
	// returns *oss.UnknownFieldsError and the result is still populated so
	// the caller can act on Body/UnknownFields.
	log.Println("=== Mode 1: default (fail on unknown fields) ===")
	result, err := client.GetBucketLifecycle(context.TODO(), &oss.GetBucketLifecycleRequest{
		Bucket: oss.Ptr(bucketName),
	})
	var ufe *oss.UnknownFieldsError
	switch {
	case errors.As(err, &ufe):
		log.Printf("unknown fields detected: %v", ufe.MissingPaths)
		log.Printf("raw body length: %d bytes", len(result.Body))
		putLifecycleRaw(client, bucketName, result.Body)
	case err != nil:
		log.Fatalf("failed to get bucket lifecycle: %v", err)
	default:
		log.Printf("no unknown fields; config: %#v", result.LifecycleConfiguration)
	}

	// --- Mode 2: UnknownFieldsCheckReport ---
	// Diagnose only — never fail the call. result.UnknownFields lists any
	// paths the SDK would drop on write-back.
	log.Println("=== Mode 2: report-only ===")
	result, err = client.GetBucketLifecycle(context.TODO(), &oss.GetBucketLifecycleRequest{
		Bucket:             oss.Ptr(bucketName),
		UnknownFieldsCheck: oss.Ptr(oss.UnknownFieldsCheckReport),
	})
	if err != nil {
		log.Fatalf("failed to get bucket lifecycle: %v", err)
	}
	log.Printf("unknown fields: %v", result.UnknownFields)

	// --- Mode 3: UnknownFieldsCheckSkip ---
	// Legacy behavior — minimal overhead, no Body/UnknownFields populated.
	log.Println("=== Mode 3: skip check ===")
	result, err = client.GetBucketLifecycle(context.TODO(), &oss.GetBucketLifecycleRequest{
		Bucket:             oss.Ptr(bucketName),
		UnknownFieldsCheck: oss.Ptr(oss.UnknownFieldsCheckSkip),
	})
	if err != nil {
		log.Fatalf("failed to get bucket lifecycle: %v", err)
	}
	log.Printf("config: %#v (Body=%d, UnknownFields=%v)",
		result.LifecycleConfiguration, len(result.Body), result.UnknownFields)
}

// putLifecycleRaw writes the raw XML back via the generic InvokeOperation
// entry point, preserving any server-side fields the typed struct would drop.
// If you need to patch the body, modify rawBody before calling this.
func putLifecycleRaw(client *oss.Client, bucket string, rawBody []byte) {
	input := &oss.OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: "application/xml",
		},
		Parameters: map[string]string{"lifecycle": ""},
		Bucket:     oss.Ptr(bucket),
		Body:       bytes.NewReader(rawBody),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})

	out, err := client.InvokeOperation(context.TODO(), input)
	if err != nil {
		log.Fatalf("raw PutBucketLifecycle failed: %v", err)
	}
	log.Printf("raw PUT ok, status=%d", out.StatusCode)
}
