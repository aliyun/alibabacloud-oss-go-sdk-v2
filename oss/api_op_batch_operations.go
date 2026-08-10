package oss

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
)

type CreateJobRequest struct {
	CreateJob *CreateJobConfig `input:"body,CreateJobRequest,xml,required"`

	RequestCommon
}

type CreateJobConfig struct {
	ConfirmationRequired *bool `xml:"ConfirmationRequired"`

	Operation *Operation `xml:"Operation"`

	Report *Report `xml:"Report"`

	ClientRequestToken *string `xml:"ClientRequestToken"`

	Manifest *Manifest `xml:"Manifest"`

	KeyPrefixManifestGenerator *KeyPrefixManifestGenerator `xml:"KeyPrefixManifestGenerator"`

	Description *string `xml:"Description"`

	Priority *int32 `xml:"Priority"`

	RoleArn *string `xml:"RoleArn"`
}

// KeyPrefixManifestGenerator Generates a manifest based on a key prefix in a source bucket.
// Used as an alternative to Manifest in CreateJob.
type KeyPrefixManifestGenerator struct {
	SourceBucket *string `xml:"SourceBucket"`

	Prefix *string `xml:"Prefix"`
}

type Operation struct {
	PutObjectTagging *Tagging `xml:"PutObjectTagging"`

	DeleteObjectTagging *string `xml:"DeleteObjectTagging"`

	AddObjectTagging *Tagging `xml:"AddObjectTagging"`

	PutObjectAcl *PutObjectAcl `xml:"PutObjectAcl"`

	RestoreObject *RestoreObject `xml:"RestoreObject"`
}

type PutObjectAcl struct {
	ObjectAcl ObjectACLType `xml:"ObjectAcl"`
}

type RestoreObject struct {
	Days *int32 `xml:"Days"`

	Tier *string `xml:"Tier"`
}

type Report struct {
	Bucket *string `xml:"Bucket"`

	Enabled *bool `xml:"Enabled"`

	Prefix *string `xml:"Prefix"`

	ReportScope *string `xml:"ReportScope"`
}

type Manifest struct {
	Location *Location `xml:"Location"`

	Spec *Spec `xml:"Spec"`
}

type Location struct {
	ETag *string `xml:"ETag"`

	Bucket *string `xml:"Bucket"`

	Object *string `xml:"Object"`

	VersionId *string `xml:"VersionId"`
}

type Spec struct {
	Fields *string `xml:"Fields"`
	Format *string `xml:"Format"`
}

type CreateJobResult struct {
	JobId *string `xml:"JobId"`

	ResultCommon
}

// CreateJob Creates a batch job that performs a specified operation on multiple objects.
func (c *Client) CreateJob(ctx context.Context, request *CreateJobRequest, optFns ...func(*Options)) (*CreateJobResult, error) {
	var err error
	if request == nil {
		request = &CreateJobRequest{}
	}
	input := &OperationInput{
		OpName: "CreateJob",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"batchJob": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"batchJob"})
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &CreateJobResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type DescribeJobRequest struct {
	BatchJobId *string `input:"query,batchJobId,required"`

	RequestCommon
}

type DescribeJobResult struct {
	Job *Job `xml:"Job"`

	ResultCommon
}

type Job struct {
	ConfirmationRequired       *bool                       `xml:"ConfirmationRequired"`
	CreationTime               *int64                      `xml:"CreationTime"`
	JobId                      *string                     `xml:"JobId"`
	Operation                  *Operation                  `xml:"Operation"`
	Report                     *Report                     `xml:"Report"`
	Manifest                   *Manifest                   `xml:"Manifest"`
	Description                *string                     `xml:"Description"`
	Priority                   *int32                      `xml:"Priority"`
	RoleArn                    *string                     `xml:"RoleArn"`
	StatusUpdateReason         *string                     `xml:"StatusUpdateReason"`
	ProgressSummary            *ProgressSummary            `xml:"ProgressSummary"`
	KeyPrefixManifestGenerator *KeyPrefixManifestGenerator `xml:"KeyPrefixManifestGenerator"`
	FailureReasons             *JobFailureReasons          `xml:"FailureReasons"`
	Status                     *string                     `xml:"Status"`
	TerminationDate            *int64                      `xml:"TerminationDate"`
}

type ProgressSummary struct {
	NumberOfTasksFailed    *int32  `xml:"NumberOfTasksFailed"`
	NumberOfTasksSucceeded *int32  `xml:"NumberOfTasksSucceeded"`
	Timers                 *Timers `xml:"Timers"`
	TotalNumberOfTasks     *int32  `xml:"TotalNumberOfTasks"`
}

type Timers struct {
	ElapsedTimeInActiveSeconds *int32 `xml:"ElapsedTimeInActiveSeconds"`
}

type JobFailureReasons struct {
	JobFailure *JobFailure `xml:"JobFailure"`
}

// JobFailure represents the failure detail of a batch operation job.
type JobFailure struct {
	FailureCode   *string `xml:"FailureCode"`
	FailureReason *string `xml:"FailureReason"`
}

func (c *Client) DescribeJob(ctx context.Context, request *DescribeJobRequest, optFns ...func(*Options)) (*DescribeJobResult, error) {
	var err error
	if request == nil {
		request = &DescribeJobRequest{}
	}
	input := &OperationInput{
		OpName: "DescribeJob",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"batchJob": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"batchJob"})
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &DescribeJobResult{}

	if err = c.unmarshalOutput(result, output, discardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type ListJobsRequest struct {
	BatchJobStatuses *string `input:"query,batchJobStatuses"`

	MaxKeys int32 `input:"query,max-keys"`

	ContinuationToken *string `input:"query,continuation-token"`

	RequestCommon
}

type ListJobsResult struct {
	Jobs *Jobs `xml:"Jobs"`

	NextToken *string `xml:"NextToken"`

	ResultCommon
}

type Jobs struct {
	JobListDescriptor []JobListDescriptor `xml:"JobListDescriptor"`
}

type JobListDescriptor struct {
	CreationTime    *int64           `xml:"CreationTime"`
	Description     *string          `xml:"Description"`
	JobId           *string          `xml:"JobId"`
	Operation       *string          `xml:"Operation"`
	Priority        *int32           `xml:"Priority"`
	Status          *string          `xml:"Status"`
	ProgressSummary *ProgressSummary `xml:"ProgressSummary"`
	TerminationDate *int64           `xml:"TerminationDate"`
}

func (c *Client) ListJobs(ctx context.Context, request *ListJobsRequest, optFns ...func(*Options)) (*ListJobsResult, error) {
	var err error
	if request == nil {
		request = &ListJobsRequest{}
	}
	input := &OperationInput{
		OpName: "ListJobs",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"batchJob": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"batchJob"})
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &ListJobsResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type UpdateJobPriorityRequest struct {
	BatchJobId *string `input:"query,batchJobId,required"`

	TargetPriority *int32 `input:"query,targetPriority,required"`

	RequestCommon
}

type UpdateJobPriorityResult struct {
	JobId *string `xml:"JobId"`

	Priority *int32 `xml:"Priority"`

	ResultCommon
}

func (c *Client) UpdateJobPriority(ctx context.Context, request *UpdateJobPriorityRequest, optFns ...func(*Options)) (*UpdateJobPriorityResult, error) {
	var err error
	if request == nil {
		request = &UpdateJobPriorityRequest{}
	}
	input := &OperationInput{
		OpName: "UpdateJobPriority",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"batchJobPriority": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"batchJobPriority"})
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &UpdateJobPriorityResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type UpdateJobStatusRequest struct {
	BatchJobId *string `input:"query,batchJobId,required"`

	RequestedJobStatus *string `input:"query,requestedJobStatus,required"`

	StatusUpdateReason *string `input:"query,statusUpdateReason"`

	RequestCommon
}

type UpdateJobStatusResult struct {
	JobId *string `xml:"JobId"`

	JobStatus *string `xml:"Status"`

	StatusUpdateReason *string `xml:"StatusUpdateReason"`

	ResultCommon
}

func (c *Client) UpdateJobStatus(ctx context.Context, request *UpdateJobStatusRequest, optFns ...func(*Options)) (*UpdateJobStatusResult, error) {
	var err error
	if request == nil {
		request = &UpdateJobStatusRequest{}
	}
	input := &OperationInput{
		OpName: "UpdateJobStatus",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"batchJobStatus": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"batchJobStatus"})
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &UpdateJobStatusResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
