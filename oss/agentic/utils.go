package agentic

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type agenticProvider struct {
	endpoint  *url.URL
	accountId string
	region    string
	suffix    string // "ab-apsr" or "bs-apsr"
	urlStyle  oss.UrlStyleType
}

func (p *agenticProvider) BuildBucketName(input *oss.OperationInput) (string, error) {
	if input.Bucket == nil {
		return "", nil
	}
	return fmt.Sprintf("%s-%s-%s-%s", *input.Bucket, p.accountId, p.region, p.suffix), nil
}

func (p *agenticProvider) BuildURL(input *oss.OperationInput) string {
	if input == nil || p.endpoint == nil {
		return ""
	}

	var host string
	var paths []string

	if input.Bucket == nil {
		host = p.endpoint.Host
	} else {
		fullName := fmt.Sprintf("%s-%s-%s-%s", *input.Bucket, p.accountId, p.region, p.suffix)
		switch p.urlStyle {
		default: // UrlStyleVirtualHosted
			host = fmt.Sprintf("%s.%s", fullName, p.endpoint.Host)
		case oss.UrlStylePath:
			host = p.endpoint.Host
			paths = append(paths, fullName)
			if input.Key == nil {
				paths = append(paths, "")
			}
		}
	}

	if input.Key != nil {
		paths = append(paths, oss.EscapePath(*input.Key, false))
	}

	path := "/" + strings.Join(paths, "/")
	return fmt.Sprintf("%s://%s%s", p.endpoint.Scheme, host, path)
}

// BucketSpaceHelper builds full bucket space names for use with a plain oss.Client.
type BucketSpaceHelper struct {
	accountId string
	region    string
}

// NewBucketSpaceHelper creates a BucketSpaceHelper from the account ID and region in cfg.
func NewBucketSpaceHelper(cfg *oss.Config) *BucketSpaceHelper {
	return &BucketSpaceHelper{
		accountId: oss.ToString(cfg.AccountId),
		region:    oss.ToString(cfg.Region),
	}
}

// ToBucketName builds the full bucket space name "{prefix}-{accountId}-{region}-bs-apsr" from a short prefix.
func (h *BucketSpaceHelper) ToBucketName(prefix string) string {
	return fmt.Sprintf("%s-%s-%s-bs-apsr", prefix, h.accountId, h.region)
}
