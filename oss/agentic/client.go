package agentic

import (
	"fmt"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// AgenticBucketClient provides access to the agentic bucket management APIs.
type AgenticBucketClient struct {
	clientImpl *oss.Client
}

// NewAgenticBucketClient creates a client for the agentic bucket management APIs.
//
// The client automatically resolves the actual bucket name from the short name
// passed in each request. The final bucket name is constructed as
// "{bucket}-{accountId}-{region}-ab-apsr", where accountId and region come from
// cfg.AccountId and cfg.Region. For example, with bucket "my-agentic",
// account "123456" and region "cn-hangzhou", the resolved name is
// "my-agentic-123456-cn-hangzhou-ab-apsr". Therefore cfg.AccountId and
// cfg.Region are required.
func NewAgenticBucketClient(cfg *oss.Config, optFns ...func(*oss.Options)) *AgenticBucketClient {
	newCfg := cfg.Copy()
	updateUserAgent(&newCfg)

	region := oss.ToString(newCfg.Region)
	accountId := oss.ToString(newCfg.AccountId)

	agenticOptFn := func(options *oss.Options) {
		p := &agenticProvider{
			endpoint:  options.Endpoint,
			accountId: accountId,
			region:    region,
			suffix:    "ab-apsr",
			urlStyle:  options.UrlStyle,
		}
		options.BucketNameResolver = p
		options.EndpointProviderE = p
	}

	allOptFns := append(optFns, agenticOptFn)
	return &AgenticBucketClient{
		clientImpl: oss.NewClient(&newCfg, allOptFns...),
	}
}

// NewBucketSpaceClient creates an oss.Client that operates on the bucket spaces of an agentic bucket.
//
// Pass the short bucket name in each request; the client automatically resolves
// it to the full name "{bucket}-{accountId}-{region}-bs-apsr", where accountId
// and region come from cfg.AccountId and cfg.Region. For example, with bucket
// "my-space", account "123456" and region "cn-hangzhou", the resolved name is
// "my-space-123456-cn-hangzhou-bs-apsr". Therefore cfg.AccountId and cfg.Region
// are required.
//
// If you prefer to use a plain oss.Client instead, build the full bucket name
// yourself with BucketSpaceHelper.ToBucketName and pass that as the request Bucket.
func NewBucketSpaceClient(cfg *oss.Config, optFns ...func(*oss.Options)) *oss.Client {
	newCfg := cfg.Copy()
	updateUserAgent(&newCfg)

	region := oss.ToString(newCfg.Region)
	accountId := oss.ToString(newCfg.AccountId)

	bsOptFn := func(options *oss.Options) {
		p := &agenticProvider{
			endpoint:  options.Endpoint,
			accountId: accountId,
			region:    region,
			suffix:    "bs-apsr",
			urlStyle:  options.UrlStyle,
		}
		options.BucketNameResolver = p
		options.EndpointProviderE = p
	}

	allOptFns := append(optFns, bsOptFn)
	return oss.NewClient(&newCfg, allOptFns...)
}

func updateUserAgent(cfg *oss.Config) {
	userAgent := "agentic-client"
	if cfg.UserAgent != nil {
		userAgent = fmt.Sprintf("%s/%s", userAgent, oss.ToString(cfg.UserAgent))
	}
	cfg.UserAgent = oss.Ptr(userAgent)
}

const (
	contentTypeDefault = "application/octet-stream"
	contentTypeXML     = "application/xml"
)
