package agentic

import (
	"context"
	"encoding/xml"
	"io"
	"reflect"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// --- CreateAgenticBucket ---

// CreateAgenticBucketRequest is the request for the CreateAgenticBucket operation.
type CreateAgenticBucketRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	// The configuration for the CreateAgenticBucket operation.
	CreateAgenticBucketConfiguration *CreateAgenticBucketConfiguration `input:"body,CreateAgenticBucketConfiguration,xml"`

	oss.RequestCommon
}

// CreateAgenticBucketConfiguration is the configuration for the CreateAgenticBucket operation.
type CreateAgenticBucketConfiguration struct {
	XMLName xml.Name `xml:"CreateAgenticBucketConfiguration"`

	// The storage class of the agentic bucket.
	StorageClass StorageClassType `xml:"StorageClass,omitempty"`

	// The data redundancy type of the agentic bucket.
	DataRedundancyType DataRedundancyType `xml:"DataRedundancyType,omitempty"`
}

// CreateAgenticBucketResult is the result for the CreateAgenticBucket operation.
type CreateAgenticBucketResult struct {
	oss.ResultCommon
}

// CreateAgenticBucket Creates an agentic bucket.
func (c *AgenticBucketClient) CreateAgenticBucket(ctx context.Context, request *CreateAgenticBucketRequest, optFns ...func(*oss.Options)) (*CreateAgenticBucketResult, error) {
	var err error
	if request == nil {
		request = &CreateAgenticBucketRequest{}
	}
	input := &oss.OperationInput{
		OpName: "CreateAgenticBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
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

	result := &CreateAgenticBucketResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- DeleteAgenticBucket ---

// DeleteAgenticBucketRequest is the request for the DeleteAgenticBucket operation.
type DeleteAgenticBucketRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// DeleteAgenticBucketResult is the result for the DeleteAgenticBucket operation.
type DeleteAgenticBucketResult struct {
	oss.ResultCommon
}

// DeleteAgenticBucket Deletes an agentic bucket.
func (c *AgenticBucketClient) DeleteAgenticBucket(ctx context.Context, request *DeleteAgenticBucketRequest, optFns ...func(*oss.Options)) (*DeleteAgenticBucketResult, error) {
	var err error
	if request == nil {
		request = &DeleteAgenticBucketRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteAgenticBucket",
		Method: "DELETE",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
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

	result := &DeleteAgenticBucketResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- GetAgenticBucket ---

// GetAgenticBucketRequest is the request for the GetAgenticBucket operation.
type GetAgenticBucketRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

// GetAgenticBucketResult is the result for the GetAgenticBucket operation.
type GetAgenticBucketResult struct {
	// The information about the agentic bucket.
	AgenticBucketInfo *AgenticBucketInfo `output:"body,AgenticBucketInfo,xml"`

	oss.ResultCommon
}

// AgenticBucketInfo is the information about an agentic bucket.
type AgenticBucketInfo struct {
	// The name of the agentic bucket.
	Name *string `xml:"Name"`

	// The owner of the agentic bucket.
	Owner *string `xml:"Owner"`

	// The region in which the agentic bucket is located.
	Region *string `xml:"Region"`

	// The storage class of the agentic bucket.
	StorageClass *string `xml:"StorageClass"`

	// The data redundancy type of the agentic bucket.
	DataRedundancyType *string `xml:"DataRedundancyType"`

	// The status of the agentic bucket.
	Status *string `xml:"Status"`

	// The resource type of the agentic bucket.
	BucketResourceType *string `xml:"BucketResourceType"`

	// The time when the agentic bucket was created.
	CreateTime *string `xml:"CreateTime"`

	// The access control list (ACL) of the agentic bucket.
	ACL *string `xml:"ACL"`

	// The Block Public Access configuration of the agentic bucket.
	PublicAccessBlock *string `xml:"PublicAccessBlock"`

	// The server-side encryption rule of the agentic bucket.
	ServerSideEncryptionRule *ServerSideEncryptionRule `xml:"ServerSideEncryptionRule"`

	// The versioning state of the agentic bucket.
	Versioning *string `xml:"Versioning"`

	// The policy of the agentic bucket.
	BucketPolicy *string `xml:"BucketPolicy"`
}

// GetAgenticBucket Queries the information about an agentic bucket.
func (c *AgenticBucketClient) GetAgenticBucket(ctx context.Context, request *GetAgenticBucketRequest, optFns ...func(*oss.Options)) (*GetAgenticBucketResult, error) {
	var err error
	if request == nil {
		request = &GetAgenticBucketRequest{}
	}
	input := &oss.OperationInput{
		OpName: "GetAgenticBucket",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
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

	result := &GetAgenticBucketResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- ListAgenticBuckets ---

// ListAgenticBucketsRequest is the request for the ListAgenticBuckets operation.
type ListAgenticBucketsRequest struct {
	// The token from which the list operation starts. You must specify the value of NextContinuationToken that is returned in the previous response as the value of ContinuationToken.
	ContinuationToken *string `input:"query,continuation-token"`

	// The maximum number of agentic buckets that can be returned.
	MaxKeys *int `input:"query,max-keys"`

	oss.RequestCommon
}

// ListAgenticBucketsResult is the result for the ListAgenticBuckets operation.
type ListAgenticBucketsResult struct {
	// The region in which the agentic buckets are located.
	Region *string `xml:"Region"`

	// The owner of the agentic buckets.
	Owner *string `xml:"Owner"`

	// The token from which the list operation started.
	ContinuationToken *string `xml:"ContinuationToken"`

	// The token from which the next list operation starts.
	NextContinuationToken *string `xml:"NextContinuationToken"`

	// Indicates whether the returned results are truncated.
	IsTruncated *bool `xml:"IsTruncated"`

	// The list of agentic buckets.
	AgenticBuckets []AgenticBucketSummary `xml:"AgenticBuckets>AgenticBucket"`

	oss.ResultCommon
}

// AgenticBucketSummary is the summary of an agentic bucket in ListAgenticBuckets.
type AgenticBucketSummary struct {
	// The name of the agentic bucket.
	Name *string `xml:"Name"`

	// The storage class of the agentic bucket.
	StorageClass *string `xml:"StorageClass"`

	// The data redundancy type of the agentic bucket.
	DataRedundancyType *string `xml:"DataRedundancyType"`

	// The time when the agentic bucket was created.
	CreateTime *string `xml:"CreateTime"`
}

// ListAgenticBuckets Lists the agentic buckets that belong to the current account.
func (c *AgenticBucketClient) ListAgenticBuckets(ctx context.Context, request *ListAgenticBucketsRequest, optFns ...func(*oss.Options)) (*ListAgenticBucketsResult, error) {
	var err error
	if request == nil {
		request = &ListAgenticBucketsRequest{}
	}
	input := &oss.OperationInput{
		OpName: "ListAgenticBuckets",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
		},
	}
	if err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.clientImpl.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &ListAgenticBucketsResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXml); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- PutAgenticBucketStatus ---

// PutAgenticBucketStatusRequest is the request for the PutAgenticBucketStatus operation.
type PutAgenticBucketStatusRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	// The status configuration of the agentic bucket.
	AgenticBucketStatus *AgenticBucketStatus `input:"body,AgenticBucketStatus,xml,required"`

	oss.RequestCommon
}

// AgenticBucketStatus is the status configuration of an agentic bucket.
type AgenticBucketStatus struct {
	XMLName xml.Name `xml:"AgenticBucketStatus"`

	// The status of the agentic bucket.
	Status *string `xml:"Status"`
}

// PutAgenticBucketStatusResult is the result for the PutAgenticBucketStatus operation.
type PutAgenticBucketStatusResult struct {
	oss.ResultCommon
}

// PutAgenticBucketStatus Configures the status of an agentic bucket.
func (c *AgenticBucketClient) PutAgenticBucketStatus(ctx context.Context, request *PutAgenticBucketStatusRequest, optFns ...func(*oss.Options)) (*PutAgenticBucketStatusResult, error) {
	var err error
	if request == nil {
		request = &PutAgenticBucketStatusRequest{}
	}
	input := &oss.OperationInput{
		OpName: "PutAgenticBucketStatus",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"status":        "",
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

	result := &PutAgenticBucketStatusResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// --- ListBucketSpaces ---

// ListBucketSpacesRequest is the request for the ListBucketSpaces operation.
type ListBucketSpacesRequest struct {
	// The name of the agentic bucket.
	Bucket *string `input:"host,bucket,required"`

	// The prefix that the names of the returned bucket spaces must contain.
	Prefix *string `input:"query,prefix"`

	// The token from which the list operation starts. You must specify the value of NextContinuationToken that is returned in the previous response as the value of ContinuationToken.
	ContinuationToken *string `input:"query,continuation-token"`

	// The name of the bucket space after which the list operation begins, sorted alphabetically.
	StartAfter *string `input:"query,start-after"`

	// The maximum number of bucket spaces that can be returned.
	MaxKeys *int `input:"query,max-keys"`

	oss.RequestCommon
}

// ListBucketSpacesResult is the result for the ListBucketSpaces operation.
type ListBucketSpacesResult struct {
	XMLName xml.Name `xml:"ListBucketSpacesResult"`

	// The owner of the bucket spaces.
	Owner *Owner `xml:"Owner"`

	// The list of bucket spaces.
	BucketSpaces []BucketSpaceSummary `xml:"BucketSpaces>BucketSpace"`

	// The prefix that the names of the returned bucket spaces contain.
	Prefix *string `xml:"Prefix"`

	// The maximum number of bucket spaces that can be returned.
	MaxKeys *int `xml:"MaxKeys"`

	// The token from which the list operation started.
	ContinuationToken *string `xml:"ContinuationToken"`

	// The token from which the next list operation starts.
	NextContinuationToken *string `xml:"NextContinuationToken"`

	// The name of the bucket space after which the list operation began.
	StartAfter *string `xml:"StartAfter"`

	// Indicates whether the returned results are truncated.
	IsTruncated *bool `xml:"IsTruncated"`

	oss.ResultCommon
}

// BucketSpaceSummary is the summary of a bucket space in ListBucketSpaces.
type BucketSpaceSummary struct {
	// The name of the bucket space.
	Name *string `xml:"Name"`

	// The region in which the bucket space is located.
	Location *string `xml:"Location"`

	// The time when the bucket space was created.
	CreationDate *string `xml:"CreationDate"`

	// The storage class of the bucket space.
	StorageClass *string `xml:"StorageClass"`
}

// ListBucketSpaces Lists the bucket spaces in an agentic bucket.
func (c *AgenticBucketClient) ListBucketSpaces(ctx context.Context, request *ListBucketSpacesRequest, optFns ...func(*oss.Options)) (*ListBucketSpacesResult, error) {
	var err error
	if request == nil {
		request = &ListBucketSpacesRequest{}
	}
	input := &oss.OperationInput{
		OpName: "ListBucketSpaces",
		Method: "GET",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"agenticBucket": "",
			"bucketSpace":   "",
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

	result := &ListBucketSpacesResult{}
	if err = c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXml); err != nil {
		return nil, c.clientImpl.ToClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

// XML unmarshal helpers

func unmarshalBodyXml(result any, output *oss.OperationOutput) error {
	var err error
	var body []byte
	if output.Body != nil {
		defer output.Body.Close()
		if body, err = io.ReadAll(output.Body); err != nil {
			return err
		}
	}
	if len(body) > 0 {
		if err = xml.Unmarshal(body, result); err != nil {
			err = &oss.DeserializationError{
				Err:      err,
				Snapshot: body,
			}
		}
	}
	return err
}

func unmarshalBodyXmlMix(result any, output *oss.OperationOutput) error {
	var err error
	var body []byte
	if output.Body != nil {
		defer output.Body.Close()
		if body, err = io.ReadAll(output.Body); err != nil {
			return err
		}
	}

	if len(body) == 0 {
		return nil
	}

	val := reflect.ValueOf(result)
	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || output == nil {
		return nil
	}

	t := val.Type()
	idx := -1
	for k := 0; k < t.NumField(); k++ {
		if tag, ok := t.Field(k).Tag.Lookup("output"); ok {
			tokens := strings.Split(tag, ",")
			if len(tokens) < 2 {
				continue
			}
			switch tokens[0] {
			case "body":
				idx = k
			}
		}
	}

	if idx >= 0 {
		dst := val.Field(idx)
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		err = xml.Unmarshal(body, dst.Interface())
	} else {
		err = xml.Unmarshal(body, result)
	}

	if err != nil {
		err = &oss.DeserializationError{
			Err:      err,
			Snapshot: body,
		}
	}

	return err
}
