//go:build integration

package agentic

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region_    = os.Getenv("OSS_TEST_REGION")
	endpoint_  = os.Getenv("OSS_TEST_ENDPOINT")
	accessID_  = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	accessKey_ = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")
	accountId_ = os.Getenv("OSS_TEST_ACCOUNT_ID")

	instance_ *AgenticBucketClient
	testOnce_ sync.Once
)

var (
	bucketNamePrefix = getBucketNamePrefix()
	letters          = []rune("abcdefghijklmnopqrstuvwxyz")
)

// getBucketNamePrefix returns the test bucket prefix; the "ab" marker is what the reaper filters on.
// Prefix plus the random part must stay within 23 characters: the resolved name
// {bucket}-{accountId}-{region}-ab-apsr becomes a DNS host label capped at 63, and the account id
// (16) plus the longest region (14) plus the separators and the tail take the other 40.
func getBucketNamePrefix() string {
	if val := os.Getenv("OSS_TEST_BUCKET_PREFIX"); val != "" {
		return val + "go-ab"
	}
	return "sdk-oss-test-go-ab"
}

func getTestConfig() *oss.Config {
	return oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithAccountId(accountId_)
}

func getAgenticBucketClient() *AgenticBucketClient {
	testOnce_.Do(func() {
		instance_ = NewAgenticBucketClient(getTestConfig())
	})
	return instance_
}

func getInvalidAkClient() *AgenticBucketClient {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("invalid-ak", "invalid-sk")).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithAccountId(accountId_)

	return NewAgenticBucketClient(cfg)
}

func getBucketSpaceClient() *oss.Client {
	return NewBucketSpaceClient(getTestConfig())
}

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func genBucketName() string {
	return bucketNamePrefix + randStr(5)
}

// disableAndReap is the shared scenario teardown: disable this run's bucket, then
// reap buckets left disabled by previous runs.
func disableAndReap(bucket string) {
	c := getAgenticBucketClient()
	_, _ = c.PutAgenticBucketStatus(context.TODO(), &PutAgenticBucketStatusRequest{
		Bucket:              oss.Ptr(bucket),
		AgenticBucketStatus: &AgenticBucketStatus{Status: oss.Ptr("Disabled")},
	})
	reapDisabledAgenticBuckets()
}

// toShortName strips the resolved tail so a listed name can be passed back to a
// client that re-expands short names. suffix is "ab-apsr" or "bs-apsr".
func toShortName(name, suffix string) string {
	return strings.TrimSuffix(name, fmt.Sprintf("-%s-%s-%s", accountId_, region_, suffix))
}

// reapDisabledAgenticBuckets deletes leftover buckets from previous runs that carry
// our prefix and are already Disabled (Enabled ones may belong to a concurrent run),
// emptying their bucket spaces first. Best-effort: all errors are swallowed.
func reapDisabledAgenticBuckets() {
	c := getAgenticBucketClient()

	paginator := c.NewListAgenticBucketsPaginator(&ListAgenticBucketsRequest{})
	for paginator.HasNext() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return
		}
		for _, summary := range page.AgenticBuckets {
			if !strings.HasPrefix(oss.ToString(summary.Name), bucketNamePrefix) {
				continue
			}
			bucket := toShortName(oss.ToString(summary.Name), "ab-apsr")
			// The list summary has no status, so fetch it; only reclaim Disabled.
			info, err := c.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
				Bucket: oss.Ptr(bucket),
			})
			if err != nil || info.AgenticBucketInfo == nil ||
				oss.ToString(info.AgenticBucketInfo.Status) != "Disabled" {
				continue
			}
			reapBucketSpaces(bucket)
			_, _ = c.DeleteAgenticBucket(context.TODO(), &DeleteAgenticBucketRequest{
				Bucket: oss.Ptr(bucket),
			})
		}
	}
}

// reapBucketSpaces empties and deletes every bucket space of a Disabled agentic
// bucket. Best-effort: errors are swallowed.
func reapBucketSpaces(bucket string) {
	c := getAgenticBucketClient()
	bsClient := getBucketSpaceClient()

	spacePaginator := c.NewListBucketSpacesPaginator(&ListBucketSpacesRequest{
		Bucket: oss.Ptr(bucket),
	})
	for spacePaginator.HasNext() {
		spacePage, err := spacePaginator.NextPage(context.TODO())
		if err != nil {
			return
		}
		for _, space := range spacePage.BucketSpaces {
			spaceName := toShortName(oss.ToString(space.Name), "bs-apsr")
			objPaginator := bsClient.NewListObjectsV2Paginator(&oss.ListObjectsV2Request{
				Bucket: oss.Ptr(spaceName),
			})
			for objPaginator.HasNext() {
				objPage, err := objPaginator.NextPage(context.TODO())
				if err != nil {
					break
				}
				for _, obj := range objPage.Contents {
					_, _ = bsClient.DeleteObject(context.TODO(), &oss.DeleteObjectRequest{
						Bucket: oss.Ptr(spaceName),
						Key:    obj.Key,
					})
				}
			}
			_, _ = bsClient.DeleteBucket(context.TODO(), &oss.DeleteBucketRequest{
				Bucket: oss.Ptr(spaceName),
			})
		}
	}
}

func dumpErrIfNotNil(err error) {
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}
}

func skipIfNotConfigured(t *testing.T) {
	if accountId_ == "" || region_ == "" {
		t.Skip("agentic integration test requires OSS_TEST_ACCOUNT_ID and OSS_TEST_REGION")
	}
}
