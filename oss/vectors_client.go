package oss

import (
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
)

const VectorsUserAgentPrefix = "VectorsBucket"

type VectorsClient struct {
	client Client
}

func NewVectorsClient(cfg *Config, optFns ...func(*Options)) *VectorsClient {
	newCfg := cfg.Copy()
	resolveVectorsUserAgent(&newCfg)
	resolveVectorsEndpoint(&newCfg)
	vectorsOptFn := func(options *Options) {
		options.Signer = &signer.SignerVectorsV4{}
	}
	allOptFns := append(optFns, vectorsOptFn)
	client := NewClient(&newCfg, allOptFns...)
	return &VectorsClient{
		client: *client,
	}
}

func resolveVectorsEndpoint(cfg *Config) {
	disableSSL := ToBool(cfg.DisableSSL)
	endpoint := ToString(cfg.Endpoint)
	region := ToString(cfg.Region)
	if len(endpoint) > 0 {
		endpoint = addEndpointScheme(endpoint, disableSSL)
	} else if isValidRegion(region) {
		endpoint = vectorsEndpointFromRegion(
			region,
			disableSSL,
			func() EndpointType {
				if ToBool(cfg.UseInternalEndpoint) {
					return EndpointInternal
				}
				return EndpointPublic
			}(),
		)
	}

	if endpoint == "" {
		return
	}

	cfg.Endpoint = &endpoint
}

func resolveVectorsUserAgent(cfg *Config) {
	if cfg.UserAgent == nil {
		cfg.UserAgent = Ptr(VectorsUserAgentPrefix)
		return
	}
	cfg.UserAgent = Ptr(fmt.Sprintf("%s/%s", VectorsUserAgentPrefix, ToString(cfg.UserAgent)))
}
