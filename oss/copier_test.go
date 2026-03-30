package oss

import (
	"context"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func TestCopierClientCopierOptions(t *testing.T) {

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou")

	client := NewClient(cfg)

	// Default
	c := NewCopier(client)
	assert.Equal(t, DefaultCopyParallel, c.options.ParallelNum)
	assert.Equal(t, DefaultCopyPartSize, c.options.PartSize)
	assert.Equal(t, DefaultCopyThreshold, c.options.MultipartCopyThreshold)
	assert.Equal(t, false, c.options.LeavePartsOnError)
	assert.Equal(t, false, c.options.DisableShallowCopy)
	assert.Equal(t, 0, len(c.options.ClientOptions))
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

	// Set From Client
	c = NewCopier(client, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = true
		co.ClientOptions = []func(do *Options){func(do *Options) {}}
		co.MetadataProperties = &HeadObjectResult{}
		co.TagProperties = &GetObjectTaggingResult{}
	})
	assert.Equal(t, int(2), c.options.ParallelNum)
	assert.Equal(t, int64(1024*1024), c.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), c.options.MultipartCopyThreshold)
	assert.Equal(t, true, c.options.LeavePartsOnError)
	assert.Equal(t, true, c.options.DisableShallowCopy)
	assert.Equal(t, 1, len(c.options.ClientOptions))
	// only supports setting from c.Copy
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

	// Use WithXXX
	c = NewCopier(client, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = true
		co.ClientOptions = []func(do *Options){func(do *Options) {}}
		co.MetadataProperties = &HeadObjectResult{}
		co.TagProperties = &GetObjectTaggingResult{}
	},
		WithCopierParallelNum(5), WithCopierPartSize(2*1024*1024))
	assert.Equal(t, int(5), c.options.ParallelNum)
	assert.Equal(t, int64(2*1024*1024), c.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), c.options.MultipartCopyThreshold)
	assert.Equal(t, true, c.options.LeavePartsOnError)
	assert.Equal(t, true, c.options.DisableShallowCopy)
	assert.Equal(t, 1, len(c.options.ClientOptions))
	// only supports setting from c.Copy
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

}

func TestCopierApiCopierOptions(t *testing.T) {

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou")

	client := NewClient(cfg)

	// Default
	c := NewCopier(client)
	assert.Equal(t, DefaultCopyParallel, c.options.ParallelNum)
	assert.Equal(t, DefaultCopyPartSize, c.options.PartSize)
	assert.Equal(t, DefaultCopyThreshold, c.options.MultipartCopyThreshold)
	assert.Equal(t, false, c.options.LeavePartsOnError)
	assert.Equal(t, false, c.options.DisableShallowCopy)
	assert.Equal(t, 0, len(c.options.ClientOptions))
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

	// Set From Client
	deleget, err := c.newDelegate(context.TODO(), &CopyObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		SourceBucket: Ptr("src-bucket"),
		SourceKey:    Ptr("src-key"),
	}, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = true
		co.ClientOptions = []func(do *Options){func(do *Options) {}}
		co.MetadataProperties = &HeadObjectResult{}
		co.TagProperties = &GetObjectTaggingResult{}
	})
	assert.NoError(t, err)

	assert.Equal(t, int(2), deleget.options.ParallelNum)
	assert.Equal(t, int64(1024*1024), deleget.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), deleget.options.MultipartCopyThreshold)
	assert.Equal(t, true, deleget.options.LeavePartsOnError)
	assert.Equal(t, true, deleget.options.DisableShallowCopy)
	assert.Equal(t, 1, len(deleget.options.ClientOptions))
	// only supports setting from c.Copy
	assert.NotNil(t, deleget.options.MetadataProperties)
	assert.NotNil(t, deleget.options.TagProperties)

	// Use WithXXX
	deleget, err = c.newDelegate(context.TODO(), &CopyObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		SourceBucket: Ptr("src-bucket"),
		SourceKey:    Ptr("src-key"),
	}, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = false
		co.ClientOptions = []func(do *Options){func(do *Options) {}, func(do *Options) {}}
		co.MetadataProperties = nil
		co.TagProperties = &GetObjectTaggingResult{}
	},
		WithCopierParallelNum(5),
		WithCopierPartSize(2*1024*1024),
	)

	assert.NoError(t, err)
	assert.Equal(t, int(5), deleget.options.ParallelNum)
	assert.Equal(t, int64(2*1024*1024), deleget.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), deleget.options.MultipartCopyThreshold)
	assert.Equal(t, true, deleget.options.LeavePartsOnError)
	assert.Equal(t, false, deleget.options.DisableShallowCopy)
	assert.Equal(t, 2, len(deleget.options.ClientOptions))
	assert.Nil(t, deleget.options.MetadataProperties)
	assert.NotNil(t, deleget.options.TagProperties)

}

func TestCopierShallowCopyFlags(t *testing.T) {

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou")

	client := NewClient(cfg)

	// Default
	c := NewCopier(client)
	assert.Equal(t, DefaultCopyParallel, c.options.ParallelNum)
	assert.Equal(t, DefaultCopyPartSize, c.options.PartSize)
	assert.Equal(t, DefaultCopyThreshold, c.options.MultipartCopyThreshold)
	assert.Equal(t, false, c.options.LeavePartsOnError)
	assert.Equal(t, false, c.options.DisableShallowCopy)
	assert.Equal(t, false, c.options.NoCheckSSE)
	assert.Equal(t, false, c.options.NoCheckCrossBucket)
	assert.Equal(t, 0, len(c.options.ClientOptions))
	assert.Nil(t, c.options.MetadataProperties)
	assert.Nil(t, c.options.TagProperties)

	// Set From Client
	deleget, err := c.newDelegate(context.TODO(), &CopyObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		SourceBucket: Ptr("src-bucket"),
		SourceKey:    Ptr("src-key"),
	}, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = true
		co.ClientOptions = []func(do *Options){func(do *Options) {}}
		co.MetadataProperties = &HeadObjectResult{}
		co.TagProperties = &GetObjectTaggingResult{}
		co.NoCheckSSE = true
		co.NoCheckCrossBucket = true
	})
	assert.NoError(t, err)

	assert.Equal(t, int(2), deleget.options.ParallelNum)
	assert.Equal(t, int64(1024*1024), deleget.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), deleget.options.MultipartCopyThreshold)
	assert.Equal(t, true, deleget.options.LeavePartsOnError)
	assert.Equal(t, true, deleget.options.DisableShallowCopy)
	assert.Equal(t, 1, len(deleget.options.ClientOptions))
	// only supports setting from c.Copy
	assert.NotNil(t, deleget.options.MetadataProperties)
	assert.NotNil(t, deleget.options.TagProperties)

	assert.Equal(t, true, deleget.options.NoCheckSSE)
	assert.Equal(t, true, deleget.options.NoCheckCrossBucket)

	// Use WithXXX
	deleget, err = c.newDelegate(context.TODO(), &CopyObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		SourceBucket: Ptr("src-bucket"),
		SourceKey:    Ptr("src-key"),
	}, func(co *CopierOptions) {
		co.ParallelNum = 2
		co.PartSize = 1024 * 1024
		co.MultipartCopyThreshold = 5 * 1024 * 1024
		co.LeavePartsOnError = true
		co.DisableShallowCopy = false
		co.ClientOptions = []func(do *Options){func(do *Options) {}, func(do *Options) {}}
		co.MetadataProperties = nil
		co.TagProperties = &GetObjectTaggingResult{}
	},
		WithCopierParallelNum(5),
		WithCopierPartSize(2*1024*1024),
		WithCopierNoCheckCrossBucket(true),
		WithCopierNoCheckSSE(true),
	)

	assert.NoError(t, err)
	assert.Equal(t, int(5), deleget.options.ParallelNum)
	assert.Equal(t, int64(2*1024*1024), deleget.options.PartSize)
	assert.Equal(t, int64(5*1024*1024), deleget.options.MultipartCopyThreshold)
	assert.Equal(t, true, deleget.options.LeavePartsOnError)
	assert.Equal(t, false, deleget.options.DisableShallowCopy)
	assert.Equal(t, 2, len(deleget.options.ClientOptions))
	assert.Nil(t, deleget.options.MetadataProperties)
	assert.NotNil(t, deleget.options.TagProperties)

	assert.Equal(t, true, deleget.options.NoCheckSSE)
	assert.Equal(t, true, deleget.options.NoCheckCrossBucket)
}
