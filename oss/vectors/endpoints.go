package vectors

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type endpointProvider struct {
	endpoint *url.URL
	acountId string
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
		host = fmt.Sprintf("%s-%s.%s", p.acountId, *input.Bucket, p.endpoint.Host)
	}

	if input.Key != nil {
		paths = append(paths, oss.EscapePath(*input.Key, false))
	}

	path = "/" + strings.Join(paths, "/")

	return fmt.Sprintf("%s://%s%s", p.endpoint.Scheme, host, path)
}
