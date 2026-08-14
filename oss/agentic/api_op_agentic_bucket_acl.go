package agentic

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// --- PutAgenticBucketAcl ---

// PutAgenticBucketAclRequest is the request for the PutAgenticBucketAcl operation.
type PutAgenticBucketAclRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	// The access control list (ACL) of the agentic bucket.
	Acl BucketACLType `input:"header,x-oss-acl,required"`

	oss.RequestCommon
}

// PutAgenticBucketAclResult is the result for the PutAgenticBucketAcl operation.
type PutAgenticBucketAclResult struct {
	oss.ResultCommon
}

// PutAgenticBucketAcl Configures the access control list (ACL) of an agentic bucket.
func (c *AgenticBucketClient) PutAgenticBucketAcl(ctx context.Context, request *PutAgenticBucketAclRequest, optFns ...func(*oss.Options)) (*PutAgenticBucketAclResult, error) {
	var err error
	if request == nil {
		request = &PutAgenticBucketAclRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutAgenticBucketAcl",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"acl":           "",
		},
		Bucket: request.Bucket,
	}

	if err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &PutAgenticBucketAclResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- GetAgenticBucketAcl ---

// GetAgenticBucketAclRequest is the request for the GetAgenticBucketAcl operation.
type GetAgenticBucketAclRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// GetAgenticBucketAclResult is the result for the GetAgenticBucketAcl operation.
type GetAgenticBucketAclResult struct {
	// The access control list (ACL) of the agentic bucket.
	ACL *string `xml:"AccessControlList>Grant"`

	// The owner of the agentic bucket.
	Owner *Owner `xml:"Owner"`

	oss.ResultCommon
}

// GetAgenticBucketAcl Queries the access control list (ACL) of an agentic bucket.
func (c *AgenticBucketClient) GetAgenticBucketAcl(ctx context.Context, request *GetAgenticBucketAclRequest, optFns ...func(*oss.Options)) (*GetAgenticBucketAclResult, error) {
	var err error
	if request == nil {
		request = &GetAgenticBucketAclRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucketAcl",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"acl":           "",
		},
		Bucket: request.Bucket,
	}

	if err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetAgenticBucketAclResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXml); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
