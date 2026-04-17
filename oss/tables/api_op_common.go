package tables

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

func (c *TablesClient) InvokeOperation(ctx context.Context, input *oss.OperationInput, optFns ...func(*oss.Options)) (*oss.OperationOutput, error) {
	return c.clientImpl.InvokeOperation(ctx, input, optFns...)
}
