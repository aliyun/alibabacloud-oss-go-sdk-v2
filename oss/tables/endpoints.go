package tables

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type endpointProvider struct {
	endpoint     *url.URL
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
			// acs:osstables:cn-hangzhou:account:bucket/bucketName
			vals1 := strings.Split(*input.Bucket, ":")
			vals2 := strings.Split(vals1[4], "/")
			host = fmt.Sprintf("%s-%s.%s", vals2[1], vals1[3], p.endpoint.Host)
		case oss.UrlStylePath:
			host = p.endpoint.Host
			paths = append(paths, *input.Bucket)
			if input.Key == nil {
				paths = append(paths, "")
			}
		}
	}

	if input.Key != nil {
		paths = append(paths, *input.Key)
	}

	path = "/" + strings.Join(paths, "/")

	return fmt.Sprintf("%s://%s%s", p.endpoint.Scheme, host, path)
}
