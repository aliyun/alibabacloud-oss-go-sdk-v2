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

// buildFullName joins the parts of a full bucket name "{bucket}-{accountId}-{region}-{suffix}".
func buildFullName(bucket, accountId, region, suffix string) string {
	return fmt.Sprintf("%s-%s-%s-%s", bucket, accountId, region, suffix)
}

// resolveBucketName expands a short bucket name into its full name, validating
// that the required accountId and region are present.
func (p *agenticProvider) resolveBucketName(bucket string) (string, error) {
	if p.accountId == "" {
		return "", oss.NewErrParamRequired("AccountId")
	}
	if p.region == "" {
		return "", oss.NewErrParamRequired("Region")
	}
	return buildFullName(bucket, p.accountId, p.region, p.suffix), nil
}

func (p *agenticProvider) BuildBucketName(input *oss.OperationInput) (string, error) {
	if input.Bucket == nil {
		return "", nil
	}
	return p.resolveBucketName(*input.Bucket)
}

func (p *agenticProvider) BuildURL(input *oss.OperationInput) (string, error) {
	if input == nil || p.endpoint == nil {
		return "", nil
	}

	var host string
	var paths []string

	if input.Bucket == nil {
		host = p.endpoint.Host
	} else {
		fullName, err := p.resolveBucketName(*input.Bucket)
		if err != nil {
			return "", err
		}
		switch p.urlStyle {
		default: // UrlStyleVirtualHosted
			// In virtual-hosted style the full name becomes the leftmost DNS label,
			// whose length must not exceed 63 characters.
			if len(fullName) > 63 {
				return "", fmt.Errorf("the host label %q exceeds the maximum length of 63 characters", fullName)
			}
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
	return fmt.Sprintf("%s://%s%s", p.endpoint.Scheme, host, path), nil
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
	return buildFullName(prefix, h.accountId, h.region, "bs-apsr")
}
