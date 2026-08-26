package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateJob(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CreateJobRequest
	var input *OperationInput
	var err error

	request = &CreateJobRequest{}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, CreateJob.")

	request = &CreateJobRequest{
		CreateJob: &CreateJobConfig{
			Operation: &Operation{
				DeleteObjectTagging: Ptr(""),
			},
			KeyPrefixManifestGenerator: &KeyPrefixManifestGenerator{
				SourceBucket: Ptr("bucket"),
				Prefix:       Ptr("batch-manifests/"),
			},
			Report: &Report{
				Bucket:      Ptr("bucket"),
				Enabled:     Ptr(true),
				Prefix:      Ptr("batch-reports/"),
				ReportScope: Ptr("AllTasks"),
			},
			Priority: Ptr(int32(10)),
			RoleArn:  Ptr("acs:ram::1234567890:role/oss-sdk-batch-test"),
		},
	}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<CreateJobRequest><Operation><DeleteObjectTagging></DeleteObjectTagging></Operation><Report><Bucket>bucket</Bucket><Enabled>true</Enabled><Prefix>batch-reports/</Prefix><ReportScope>AllTasks</ReportScope></Report><KeyPrefixManifestGenerator><SourceBucket>bucket</SourceBucket><Prefix>batch-manifests/</Prefix></KeyPrefixManifestGenerator><Priority>10</Priority><RoleArn>acs:ram::1234567890:role/oss-sdk-batch-test</RoleArn></CreateJobRequest>")

	request = &CreateJobRequest{
		CreateJob: &CreateJobConfig{
			Operation: &Operation{
				PutObjectTagging: &Tagging{
					TagSet: &TagSet{
						Tags: []Tag{
							{
								Key:   Ptr("key1"),
								Value: Ptr("value1"),
							},
						},
					},
				},
			},
			ClientRequestToken: Ptr("unique-token-123"),
			Manifest: &Manifest{
				Location: &Location{
					ETag:   Ptr("d41d8cd98f00b204e9800998ecf8427e"),
					Bucket: Ptr("manifest-bucket"),
					Object: Ptr("manifest.csv"),
				},
				Spec: &Spec{
					Fields: Ptr("Bucket,Key"),
					Format: Ptr("OSS_BatchOperations_CSV_20250611"),
				},
			},
			Description: Ptr("批量设置对象标签任务"),
			Priority:    Ptr(int32(10)),
			RoleArn:     Ptr("acs:ram::1234567890:role/oss-sdk-batch-test"),
		},
	}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<CreateJobRequest><Operation><PutObjectTagging><TagSet><Tag><Key>key1</Key><Value>value1</Value></Tag></TagSet></PutObjectTagging></Operation><ClientRequestToken>unique-token-123</ClientRequestToken><Manifest><Location><ETag>d41d8cd98f00b204e9800998ecf8427e</ETag><Bucket>manifest-bucket</Bucket><Object>manifest.csv</Object></Location><Spec><Fields>Bucket,Key</Fields><Format>OSS_BatchOperations_CSV_20250611</Format></Spec></Manifest><Description>批量设置对象标签任务</Description><Priority>10</Priority><RoleArn>acs:ram::1234567890:role/oss-sdk-batch-test</RoleArn></CreateJobRequest>")
}

func TestUnmarshalOutput_CreateJob(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<CreateJobResult>
  <JobId>MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU</JobId>
</CreateJobResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &CreateJobResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.JobId, "MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateJobResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateJobResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateJobResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DescribeJob(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DescribeJobRequest
	var input *OperationInput
	var err error

	request = &DescribeJobRequest{}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BatchJobId.")

	request = &DescribeJobRequest{
		BatchJobId: Ptr("MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU"),
	}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["batchJobId"], "MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU")
}

func TestUnmarshalOutput_DescribeJob(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<DescribeJobResult>
  <Job>
    <ConfirmationRequired>false</ConfirmationRequired>
    <CreationTime>1749983400</CreationTime>
    <JobId>MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=</JobId>
    <Operation>
      <RestoreObject>
        <Days>7</Days>
        <Tier>Standard</Tier>
      </RestoreObject>
    </Operation>
    <Report>
      <Bucket>report-bucket</Bucket>
      <Enabled>true</Enabled>
      <Prefix>reports/</Prefix>
      <ReportScope>AllTasks</ReportScope>
    </Report>
    <Manifest>
      <Location>
        <ETag>d41d8cd98f00b204e9800998ecf8427e</ETag>
        <Bucket>manifest-bucket</Bucket>
        <Object>manifest.csv</Object>
        <VersionId>3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY+MTRCxf3vjVBH40Nr8X8gdRQBpUMLUo</VersionId>
      </Location>
      <Spec>
        <Fields>Bucket,Key</Fields>
        <Format>OSS_BatchOperations_CSV_20250611</Format>
      </Spec>
    </Manifest>
    <Description>批量恢复归档对象任务</Description>
    <Priority>10</Priority>
    <RoleArn>arn:acs:ram::uid:role/BatchOperationRole</RoleArn>
    <StatusUpdateReason>Task completed successfully</StatusUpdateReason>
    <ProgressSummary>
      <NumberOfTasksFailed>0</NumberOfTasksFailed>
      <NumberOfTasksSucceeded>1000</NumberOfTasksSucceeded>
      <Timers>
        <ElapsedTimeInActiveSeconds>3600</ElapsedTimeInActiveSeconds>
      </Timers>
      <TotalNumberOfTasks>1000</TotalNumberOfTasks>
    </ProgressSummary>
    <Status>Complete</Status>
    <TerminationDate>1749987000</TerminationDate>
  </Job>
</DescribeJobResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &DescribeJobResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Job.ConfirmationRequired, false)
	assert.Equal(t, *result.Job.CreationTime, int64(1749983400))
	assert.Equal(t, *result.Job.Status, "Complete")
	assert.Equal(t, *result.Job.TerminationDate, int64(1749987000))
	assert.Equal(t, *result.Job.Description, "批量恢复归档对象任务")
	assert.Equal(t, *result.Job.Priority, int32(10))
	assert.Equal(t, *result.Job.RoleArn, "arn:acs:ram::uid:role/BatchOperationRole")
	assert.Equal(t, *result.Job.TerminationDate, int64(1749987000))
	assert.Equal(t, *result.Job.JobId, "MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=")
	assert.Equal(t, *result.Job.Operation.RestoreObject.Days, int32(7))
	assert.Equal(t, *result.Job.Operation.RestoreObject.Tier, "Standard")
	assert.Equal(t, *result.Job.Report.Bucket, "report-bucket")
	assert.Equal(t, *result.Job.Report.Enabled, true)
	assert.Equal(t, *result.Job.Report.Prefix, "reports/")
	assert.Equal(t, *result.Job.Report.ReportScope, "AllTasks")
	assert.Equal(t, *result.Job.Manifest.Location.ETag, "d41d8cd98f00b204e9800998ecf8427e")
	assert.Equal(t, *result.Job.Manifest.Location.Bucket, "manifest-bucket")
	assert.Equal(t, *result.Job.Manifest.Location.Object, "manifest.csv")
	assert.Equal(t, *result.Job.Manifest.Location.VersionId, "3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY+MTRCxf3vjVBH40Nr8X8gdRQBpUMLUo")
	assert.Equal(t, *result.Job.Manifest.Spec.Fields, "Bucket,Key")
	assert.Equal(t, *result.Job.Manifest.Spec.Format, "OSS_BatchOperations_CSV_20250611")

	body = `<DescribeJobResult>
			  <Job>
				<JobId>testJob123</JobId>
				<Status>Failed</Status>
				<KeyPrefixManifestGenerator>
				  <SourceBucket>source-bucket</SourceBucket>
				  <Prefix>data/</Prefix>
				</KeyPrefixManifestGenerator>
				<FailureReasons>
					<JobFailure>
						<FailureCode>InternalError</FailureCode>
						<FailureReason>Internal service error</FailureReason>
					</JobFailure>
				</FailureReasons>
			  </Job>
			</DescribeJobResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &DescribeJobResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Job.JobId, "testJob123")
	assert.Equal(t, *result.Job.Status, "Failed")
	assert.Equal(t, *result.Job.KeyPrefixManifestGenerator.SourceBucket, "source-bucket")
	assert.Equal(t, *result.Job.KeyPrefixManifestGenerator.Prefix, "data/")
	assert.Equal(t, *result.Job.FailureReasons.JobFailure.FailureCode, "InternalError")
	assert.Equal(t, *result.Job.FailureReasons.JobFailure.FailureReason, "Internal service error")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DescribeJobResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DescribeJobResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DescribeJobResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListJobs(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListJobsRequest
	var input *OperationInput
	var err error

	request = &ListJobsRequest{}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &ListJobsRequest{
		BatchJobStatuses:  Ptr("Complete"),
		MaxKeys:           int32(100),
		ContinuationToken: Ptr("token"),
	}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["batchJobStatuses"], "Complete")
	assert.Equal(t, input.Parameters["max-keys"], "100")
	assert.Equal(t, input.Parameters["continuation-token"], "token")

	request = &ListJobsRequest{
		BatchJobStatuses:  Ptr("Active | Cancelled | Cancelling"),
		MaxKeys:           int32(100),
		ContinuationToken: Ptr("token"),
	}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["batchJobStatuses"], "Active | Cancelled | Cancelling")
	assert.Equal(t, input.Parameters["max-keys"], "100")
	assert.Equal(t, input.Parameters["continuation-token"], "token")
}

func TestUnmarshalOutput_ListJobs(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListJobsResult>
  <NextToken>OTY2NzU4</NextToken>
  <Jobs>
    <JobListDescriptor>
      <CreationTime>1785825069</CreationTime>
      <Description>Batch delete object tagging job</Description>
      <JobId>OTliNzRjYzE3OGM0NDlmZmI5MzZlZmQ0OGUxZDA5Zjk=</JobId>
      <Operation>DeleteObjectTagging</Operation>
      <Priority>20</Priority>
      <ProgressSummary>
        <NumberOfTasksFailed>-1</NumberOfTasksFailed>
        <NumberOfTasksSucceeded>-1</NumberOfTasksSucceeded>
        <TotalNumberOfTasks>-1</TotalNumberOfTasks>
        <Timers>
          <ElapsedTimeInActiveSeconds>0</ElapsedTimeInActiveSeconds>
        </Timers>
      </ProgressSummary>
      <TerminationDate>1785825076</TerminationDate>
      <Status>Failed</Status>
    </JobListDescriptor>
    <JobListDescriptor>
      <CreationTime>1785831719</CreationTime>
      <Description>Batch put object acl job</Description>
      <JobId>ZTY0ZGQ0YjY5ZjVhNDI5MDhkYjE5ZGEzYmU5M2YxODA=</JobId>
      <Operation>PutObjectAcl</Operation>
      <Priority>20</Priority>
      <ProgressSummary>
        <NumberOfTasksFailed>0</NumberOfTasksFailed>
        <NumberOfTasksSucceeded>0</NumberOfTasksSucceeded>
        <TotalNumberOfTasks>0</TotalNumberOfTasks>
        <Timers>
          <ElapsedTimeInActiveSeconds>8</ElapsedTimeInActiveSeconds>
        </Timers>
      </ProgressSummary>
      <TerminationDate>1785831749</TerminationDate>
      <Status>Complete</Status>
    </JobListDescriptor>
    <JobListDescriptor>
      <CreationTime>1785831976</CreationTime>
      <Description>Batch put object acl job</Description>
      <JobId>YTBkNDMxZWFlZmUxNDE4MTkwZTUwNzY3YzE5NDUxMjc=</JobId>
      <Operation>PutObjectAcl</Operation>
      <Priority>10</Priority>
      <ProgressSummary>
        <NumberOfTasksFailed>3</NumberOfTasksFailed>
        <NumberOfTasksSucceeded>0</NumberOfTasksSucceeded>
        <TotalNumberOfTasks>3</TotalNumberOfTasks>
        <Timers>
          <ElapsedTimeInActiveSeconds>9</ElapsedTimeInActiveSeconds>
        </Timers>
      </ProgressSummary>
      <TerminationDate>1785832020</TerminationDate>
      <Status>Failed</Status>
    </JobListDescriptor>
  </Jobs>
</ListJobsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &ListJobsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Len(t, result.Jobs.JobListDescriptor, 3)
	assert.Equal(t, *result.NextToken, "OTY2NzU4")
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].CreationTime, int64(1785825069))
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].Description, "Batch delete object tagging job")
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].JobId, "OTliNzRjYzE3OGM0NDlmZmI5MzZlZmQ0OGUxZDA5Zjk=")
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].Operation, "DeleteObjectTagging")
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].Priority, int32(20))
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].Status, "Failed")
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].ProgressSummary.NumberOfTasksFailed, int32(-1))
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].TerminationDate, int64(1785825076))
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].ProgressSummary.NumberOfTasksSucceeded, int32(-1))
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].ProgressSummary.TotalNumberOfTasks, int32(-1))
	assert.Equal(t, *result.Jobs.JobListDescriptor[0].ProgressSummary.Timers.ElapsedTimeInActiveSeconds, int32(0))

	assert.Equal(t, *result.Jobs.JobListDescriptor[1].CreationTime, int64(1785831719))
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].Description, "Batch put object acl job")
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].JobId, "ZTY0ZGQ0YjY5ZjVhNDI5MDhkYjE5ZGEzYmU5M2YxODA=")
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].Operation, "PutObjectAcl")
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].Priority, int32(20))
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].Status, "Complete")
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].ProgressSummary.NumberOfTasksFailed, int32(0))
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].TerminationDate, int64(1785831749))
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].ProgressSummary.NumberOfTasksSucceeded, int32(0))
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].ProgressSummary.TotalNumberOfTasks, int32(0))
	assert.Equal(t, *result.Jobs.JobListDescriptor[1].ProgressSummary.Timers.ElapsedTimeInActiveSeconds, int32(8))
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].CreationTime, int64(1785831976))
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].Description, "Batch put object acl job")
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].JobId, "YTBkNDMxZWFlZmUxNDE4MTkwZTUwNzY3YzE5NDUxMjc=")
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].Operation, "PutObjectAcl")
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].Priority, int32(10))
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].Status, "Failed")
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].ProgressSummary.NumberOfTasksFailed, int32(3))
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].TerminationDate, int64(1785832020))
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].ProgressSummary.NumberOfTasksSucceeded, int32(0))
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].ProgressSummary.TotalNumberOfTasks, int32(3))
	assert.Equal(t, *result.Jobs.JobListDescriptor[2].ProgressSummary.Timers.ElapsedTimeInActiveSeconds, int32(9))

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListJobsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListJobsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListJobsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_UpdateJobPriority(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UpdateJobPriorityRequest
	var input *OperationInput
	var err error

	request = &UpdateJobPriorityRequest{}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BatchJobId.")

	request = &UpdateJobPriorityRequest{
		BatchJobId: Ptr("testBatchJobId"),
	}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, TargetPriority.")

	request = &UpdateJobPriorityRequest{
		BatchJobId:     Ptr("testBatchJobId"),
		TargetPriority: Ptr(int32(10)),
	}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["batchJobId"], "testBatchJobId")
	assert.Equal(t, input.Parameters["targetPriority"], "10")
}

func TestUnmarshalOutput_UpdateJobPriority(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<UpdateJobPriorityResult>
   <JobId>MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=</JobId>
   <Priority>20</Priority>
</UpdateJobPriorityResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &UpdateJobPriorityResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.JobId, "MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=")
	assert.Equal(t, *result.Priority, int32(20))

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &UpdateJobPriorityResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &UpdateJobPriorityResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &UpdateJobPriorityResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_UpdateJobStatus(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UpdateJobStatusRequest
	var input *OperationInput
	var err error

	request = &UpdateJobStatusRequest{}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, BatchJobId.")

	request = &UpdateJobStatusRequest{
		BatchJobId: Ptr("testBatchJobId"),
	}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, RequestedJobStatus.")

	request = &UpdateJobStatusRequest{
		BatchJobId:         Ptr("testBatchJobId"),
		RequestedJobStatus: Ptr("Cancelled"),
		StatusUpdateReason: Ptr("User requested cancellation"),
	}
	input = &OperationInput{
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
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["batchJobId"], "testBatchJobId")
	assert.Equal(t, input.Parameters["requestedJobStatus"], "Cancelled")
	assert.Equal(t, input.Parameters["statusUpdateReason"], "User requested cancellation")
}

func TestUnmarshalOutput_UpdateJobStatus(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<UpdateJobStatusResult>
   <JobId>MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=</JobId>
   <Status>Cancelling</Status>
   <StatusUpdateReason>User requested cancellation</StatusUpdateReason>
</UpdateJobStatusResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &UpdateJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.JobId, "MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=")
	assert.Equal(t, *result.JobStatus, "Cancelling")
	assert.Equal(t, *result.StatusUpdateReason, "User requested cancellation")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &UpdateJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &UpdateJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &UpdateJobStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
