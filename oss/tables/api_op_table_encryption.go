package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type GetTableEncryptionRequest struct {
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	Namespace *string `input:"nop,namespace,required"`

	Name *string `input:"nop,name,required"`

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
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("tables/%s/%s/%s/encryption", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)), url.QueryEscape(oss.ToString(request.Name)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
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
