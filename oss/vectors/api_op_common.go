package vectors

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

func (c *VectorsClient) InvokeOperation(ctx context.Context, input *oss.OperationInput, optFns ...func(*oss.Options)) (*oss.OperationOutput, error) {
	return c.clientImpl.InvokeOperation(ctx, input, optFns...)
}
