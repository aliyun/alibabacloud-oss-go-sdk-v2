package oss

import (
	"context"
)

func (c *VectorsClient) InvokeOperation(ctx context.Context, input *OperationInput, optFns ...func(*Options)) (*OperationOutput, error) {
	return c.client.InvokeOperation(ctx, input, optFns...)
}
