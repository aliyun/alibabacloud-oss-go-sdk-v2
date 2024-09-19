package oss

import (
	"context"
)

type BucketAntiDDOSConfiguration struct {
	// The custom domain names that you want to protect.
	Domains []string `xml:"Cnames>Domain"`
}

type BucketAntiDDOSInfo struct {
	// The ID of the bucket owner.
	Owner *string `xml:"Owner"`

	// The time when the Anti-DDoS instance was created. The value is a timestamp.
	Ctime *int64 `xml:"Ctime"`

	// The status of the Anti-DDoS instance. Valid values:*   Init*   Defending*   HaltDefending
	Status *string `xml:"Status"`

	// The custom domain names.
	Domains []string `xml:"Cnames>Domain"`

	// The ID of the Anti-DDoS instance.
	InstanceId *string `xml:"InstanceId"`

	// The name of the bucket for which the Anti-DDoS instance is created.
	Bucket *string `xml:"Bucket"`

	// The time when the Anti-DDoS instance was last updated. The value is a timestamp.
	Mtime *int64 `xml:"Mtime"`

	// The time when the Anti-DDoS instance was activated. The value is a timestamp.
	ActiveTime *int64 `xml:"ActiveTime"`

	// The type of the Anti-DDoS instance. Valid value: AntiDDos Premimum.
	Type *string `xml:"Type"`
}

type UserAntiDDOSInfo struct {
	// The time when the Anti-DDoS instance was created. The value is a timestamp.
	Ctime *int64 `xml:"Ctime"`

	// The time when the Anti-DDoS instance was last updated. The value is a timestamp.
	Mtime *int64 `xml:"Mtime"`

	// The time when the Anti-DDoS instance was activated. The value is a timestamp.
	ActiveTime *int64 `xml:"ActiveTime"`

	// The status of the Anti-DDoS instance. Valid values:*   Init*   Defending*   HaltDefending
	Status *string `xml:"Status"`

	// The ID of the Anti-DDoS instance.
	InstanceId *string `xml:"InstanceId"`

	// The ID of the owner of the Anti-DDoS instance.
	Owner *string `xml:"Owner"`
}

type AntiDDOSListConfiguration struct {
	// The Anti-DDoS instances whose names are alphabetically after the specified marker.
	Marker *string `xml:"Marker"`

	// Indicates whether all Anti-DDoS instances are returned.- true: All Anti-DDoS instances are returned.- false: Not all Anti-DDoS instances are returned.
	IsTruncated *bool `xml:"IsTruncated"`

	// The container that stores information about the Anti-DDoS instance.
	AntiDDOSConfigurations []BucketAntiDDOSInfo `xml:"AntiDDOSConfiguration"`
}

type UpdateUserAntiDDosInfoRequest struct {
	// The Anti-DDoS instance ID.
	DefenderInstance *string `input:"header,x-oss-defender-instance,required"`

	// The new status of the Anti-DDoS instance. Set the value to HaltDefending, which indicates that the Anti-DDos protection is disabled for a bucket.
	DefenderStatus *string `input:"header,x-oss-defender-status,required"`

	RequestCommon
}

type UpdateUserAntiDDosInfoResult struct {
	ResultCommon
}

// UpdateUserAntiDDosInfo Modifies the status of an Anti-DDoS instance.
func (c *Client) UpdateUserAntiDDosInfo(ctx context.Context, request *UpdateUserAntiDDosInfoRequest, optFns ...func(*Options)) (*UpdateUserAntiDDosInfoResult, error) {
	var err error
	if request == nil {
		request = &UpdateUserAntiDDosInfoRequest{}
	}
	input := &OperationInput{
		OpName: "UpdateUserAntiDDosInfo",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &UpdateUserAntiDDosInfoResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type UpdateBucketAntiDDosInfoRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The Anti-DDoS instance ID.
	DefenderInstance *string `input:"header,x-oss-defender-instance,required"`

	// The new status of the Anti-DDoS instance. Valid values:*   Init: You must specify the custom domain name that you want to protect.*   Defending: You can select whether to specify the custom domain name that you want to protect.*   HaltDefending: You do not need to specify the custom domain name that you want to protect.
	DefenderStatus *string `input:"header,x-oss-defender-status,required"`

	// The request body schema.
	BucketAntiDDOSConfiguration *BucketAntiDDOSConfiguration `input:"body,AntiDDOSConfiguration,xml"`

	RequestCommon
}

type UpdateBucketAntiDDosInfoResult struct {
	ResultCommon
}

// UpdateBucketAntiDDosInfo Updates the status of an Anti-DDoS instance of a bucket.
func (c *Client) UpdateBucketAntiDDosInfo(ctx context.Context, request *UpdateBucketAntiDDosInfoRequest, optFns ...func(*Options)) (*UpdateBucketAntiDDosInfoResult, error) {
	var err error
	if request == nil {
		request = &UpdateBucketAntiDDosInfoRequest{}
	}
	input := &OperationInput{
		OpName: "UpdateBucketAntiDDosInfo",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &UpdateBucketAntiDDosInfoResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type ListBucketAntiDDosInfoRequest struct {
	// The name of the Anti-DDoS instance from which the list starts. The Anti-DDoS instances whose names are alphabetically after the value of marker are returned.  You can set marker to an empty string in the first request. If IsTruncated is returned in the response and the value of IsTruncated is true, you must use the value of Marker in the response as the value of marker in the next request.
	Marker *string `input:"query,marker"`

	// The maximum number of Anti-DDoS instances that can be returned.Valid values: 1 to 100.Default value: 100.
	MaxKeys *string `input:"query,max-keys"`

	RequestCommon
}

type ListBucketAntiDDosInfoResult struct {
	// The container that stores the protection list of an Anti-DDoS instance of a bucket.
	AntiDDOSListConfiguration *AntiDDOSListConfiguration `output:"body,AntiDDOSListConfiguration,xml"`

	ResultCommon
}

// ListBucketAntiDDosInfo Queries the protection list of an Anti-DDoS instance of a bucket.
func (c *Client) ListBucketAntiDDosInfo(ctx context.Context, request *ListBucketAntiDDosInfoRequest, optFns ...func(*Options)) (*ListBucketAntiDDosInfoResult, error) {
	var err error
	if request == nil {
		request = &ListBucketAntiDDosInfoRequest{}
	}
	input := &OperationInput{
		OpName: "ListBucketAntiDDosInfo",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"bucketAntiDDos": "",
		},
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &ListBucketAntiDDosInfoResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type InitUserAntiDDosInfoRequest struct {
	RequestCommon
}

type InitUserAntiDDosInfoResult struct {
	// The ID of the Anti-DDoS instance.
	DefenderInstance *string `output:"header,x-oss-defender-instance"`

	ResultCommon
}

// InitUserAntiDDosInfo Creates an Anti-DDoS instance.
func (c *Client) InitUserAntiDDosInfo(ctx context.Context, request *InitUserAntiDDosInfoRequest, optFns ...func(*Options)) (*InitUserAntiDDosInfoResult, error) {
	var err error
	if request == nil {
		request = &InitUserAntiDDosInfoRequest{}
	}
	input := &OperationInput{
		OpName: "InitUserAntiDDosInfo",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &InitUserAntiDDosInfoResult{}

	if err = c.unmarshalOutput(result, output, unmarshalHeader, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type InitBucketAntiDDosInfoRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The ID of the Anti-DDoS instance.
	DefenderInstance *string `input:"header,x-oss-defender-instance,required"`

	// The type of the Anti-DDoS instance. Set the value to AntiDDos Premimum.
	DefenderType *string `input:"header,x-oss-defender-type,required"`

	// The request body schema.
	BucketAntiDDOSConfiguration *BucketAntiDDOSConfiguration `input:"body,AntiDDOSConfiguration,xml"`

	RequestCommon
}

type InitBucketAntiDDosInfoResult struct {
	// The ID of the Anti-DDoS instance.
	DefenderInstance *string `output:"header,x-oss-defender-instance"`

	ResultCommon
}

// InitBucketAntiDDosInfo Initializes an Anti-DDoS instance for a bucket.
func (c *Client) InitBucketAntiDDosInfo(ctx context.Context, request *InitBucketAntiDDosInfoRequest, optFns ...func(*Options)) (*InitBucketAntiDDosInfoResult, error) {
	var err error
	if request == nil {
		request = &InitBucketAntiDDosInfoRequest{}
	}
	input := &OperationInput{
		OpName: "InitBucketAntiDDosInfo",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &InitBucketAntiDDosInfoResult{}

	if err = c.unmarshalOutput(result, output, unmarshalHeader, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetUserAntiDDosInfoRequest struct {
	RequestCommon
}

type GetUserAntiDDosInfoResult struct {
	// The container that stores information about the Anti-DDoS instance.
	AntiDDOSConfigurations []UserAntiDDOSInfo `xml:"AntiDDOSConfiguration"`

	ResultCommon
}

// GetUserAntiDDosInfo Queries the information about an Anti-DDoS instance of an Alibaba Cloud account.
func (c *Client) GetUserAntiDDosInfo(ctx context.Context, request *GetUserAntiDDosInfoRequest, optFns ...func(*Options)) (*GetUserAntiDDosInfoResult, error) {
	var err error
	if request == nil {
		request = &GetUserAntiDDosInfoRequest{}
	}
	input := &OperationInput{
		OpName: "GetUserAntiDDosInfo",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &GetUserAntiDDosInfoResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
