package tables

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableEncryptionRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	Namespace *string `input:"nop,namespace,required"`

	Table *string `input:"nop,name,required"`

	oss.RequestCommon
}

type GetTableEncryptionResult struct {
	EncryptionConfiguration *EncryptionConfiguration `output:"body,encryptionConfiguration,json"`

	oss.ResultCommon
}

// GetTableEncryption Queries the encryption rules configured for a table.
func (c *TablesClient) GetTableEncryption(ctx context.Context, request *GetTableEncryptionRequest, optFns ...func(*oss.Options)) (*GetTableEncryptionResult, error) {
	var err error
	if request == nil {
		request = &GetTableEncryptionRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetTableEncryption",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"encryption":                    "",
			"tables":                        "",
			oss.ToString(request.Namespace): "",
			oss.ToString(request.Table):     "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetTableEncryptionResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
