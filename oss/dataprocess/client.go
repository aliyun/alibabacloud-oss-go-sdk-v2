package dataprocess

import (
	"fmt"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// Client is the client for accessing OSS Data Process API
type Client struct {
	client *oss.Client
}

// NewClient creates a new DataProcess client with the given configuration
func NewClient(cfg *oss.Config, optFns ...func(*oss.Options)) *Client {
	return &Client{
		client: oss.NewClient(cfg, optFns...),
	}
}

// Unwrap returns the underlying OSS client
func (c *Client) Unwrap() *oss.Client { return c.client }

// toClientError converts an error to a client error
func (c *Client) toClientError(err error, code string, output *oss.OperationOutput) error {
	if err == nil {
		return nil
	}

	return &oss.ClientError{
		Code: code,
		Message: fmt.Sprintf("execute %s fail, error code is %s, request id:%s",
			output.Input.OpName,
			code,
			output.Headers.Get(oss.HeaderOssRequestID),
		),
		Err: err,
	}
}
