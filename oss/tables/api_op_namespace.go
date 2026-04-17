package tables

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type CreateNamespaceRequest struct {
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	// The namespace.
	Namespace []string `input:"body,namespace,json,required"`

	oss.RequestCommon
}

type CreateNamespaceResult struct {
	Namespace []string `json:"namespace"`

	TableBucketARN *string `json:"tableBucketARN"`

	oss.ResultCommon
}

// CreateNamespace Creates a namespace.
func (c *TablesClient) CreateNamespace(ctx context.Context, request *CreateNamespaceRequest, optFns ...func(*oss.Options)) (*CreateNamespaceResult, error) {
	var err error
	if request == nil {
		request = &CreateNamespaceRequest{}
	}
	input := &oss.OperationInput{
		OpName: "CreateNamespace",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &CreateNamespaceResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetNamespaceRequest struct {
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	Namespace *string `input:"nop,namespace,required"`

	oss.RequestCommon
}

type GetNamespaceResult struct {
	CreatedAt      *string  `json:"createdAt"`
	CreatedBy      *string  `json:"createdBy"`
	Namespace      []string `json:"namespace"`
	NamespaceId    *string  `json:"namespaceId"`
	OwnerAccountId *string  `json:"ownerAccountId"`
	TableBucketId  *string  `json:"tableBucketId"`

	oss.ResultCommon
}

// GetNamespace Queries information about a table bucket.
func (c *TablesClient) GetNamespace(ctx context.Context, request *GetNamespaceRequest, optFns ...func(*oss.Options)) (*GetNamespaceResult, error) {
	var err error
	if request == nil {
		request = &GetNamespaceRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetNamespace",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetNamespaceResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type ListNamespacesRequest struct {
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	// The token from which the ListNamespaces operation must start.
	ContinuationToken *string `input:"query,continuationToken"`

	// The maximum number of namespaces that can be returned in the single query.
	// Valid values: 1 to 1000.
	MaxNamespaces int32 `input:"query,maxNamespaces"`

	// The prefix that the names of returned buckets must contain.
	Prefix *string `input:"query,prefix"` // Limits the response to keys that begin with the specified prefix

	oss.RequestCommon
}

type ListNamespacesResult struct {
	// The token from which the ListNamespaces operation must start.
	ContinuationToken *string `json:"continuationToken"`

	// The container that stores information about namespaces.
	Namespaces []NamespaceSummary `json:"namespaces"`

	oss.ResultCommon
}

type NamespaceSummary struct {
	CreatedAt      *string  `json:"createdAt"`
	CreatedBy      *string  `json:"createdBy"`
	Namespace      []string `json:"namespace"`
	NamespaceId    *string  `json:"namespaceId"`
	OwnerAccountId *string  `json:"ownerAccountId"`
	TableBucketId  *string  `json:"tableBucketId"`
}

// ListNamespaces Lists vector buckets that belong to the current account.
func (c *TablesClient) ListNamespaces(ctx context.Context, request *ListNamespacesRequest, optFns ...func(*oss.Options)) (*ListNamespacesResult, error) {
	var err error
	if request == nil {
		request = &ListNamespacesRequest{}
	}
	input := &oss.OperationInput{
		OpName: "ListNamespaces",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	if err = c.marshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &ListNamespacesResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DeleteNamespaceRequest struct {
	TableBucketARN *string `input:"nop,tableBucketARN,required"`

	// The namespace to delete.
	Namespace *string `input:"nop,namespace,required"`

	oss.RequestCommon
}

type DeleteNamespaceResult struct {
	oss.ResultCommon
}

// DeleteNamespace Deletes a namespace.
func (c *TablesClient) DeleteNamespace(ctx context.Context, request *DeleteNamespaceRequest, optFns ...func(*oss.Options)) (*DeleteNamespaceResult, error) {
	var err error
	if request == nil {
		request = &DeleteNamespaceRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteNamespace",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Bucket: request.TableBucketARN,
		Key:    oss.Ptr(fmt.Sprintf("namespaces/%s/%s", url.QueryEscape(oss.ToString(request.TableBucketARN)), url.QueryEscape(oss.ToString(request.Namespace)))),
	}
	input.OpMetadata.Add(oss.OpMetaKeyIsBucketArn, true)
	if err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteNamespaceResult{}
	if err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
