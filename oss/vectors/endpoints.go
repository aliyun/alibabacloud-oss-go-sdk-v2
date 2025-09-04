package vectors

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type endpointProvider struct {
	endpoint     *url.URL
	accountId    string
	endpointType oss.UrlStyleType
}

func (p *endpointProvider) BuildURL(input *oss.OperationInput) string {
	if input == nil || p.endpoint == nil {
		return ""
	}
	var host string
	var path string
	var paths []string

	if input.Bucket == nil {
		host = p.endpoint.Host
	} else {
		switch p.endpointType {
		default: // UrlStyleVirtualHosted
			host = fmt.Sprintf("%s-%s.%s", *input.Bucket, p.accountId, p.endpoint.Host)
		case oss.UrlStylePath:
			host = p.endpoint.Host
			paths = append(paths, *input.Bucket)
			if input.Key == nil {
				paths = append(paths, "")
			}
		}
	}

	if input.Key != nil {
		paths = append(paths, oss.EscapePath(*input.Key, false))
	}

	path = "/" + strings.Join(paths, "/")

	return fmt.Sprintf("%s://%s%s", p.endpoint.Scheme, host, path)
}
