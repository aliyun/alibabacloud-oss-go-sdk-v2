package oss

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockCreateJobSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateJobRequest
	CheckOutputFn  func(t *testing.T, o *CreateJobResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateJobRequest><Operation><DeleteObjectTagging></DeleteObjectTagging></Operation><Report><Bucket>bucket</Bucket><Enabled>true</Enabled><Prefix>batch-reports/</Prefix><ReportScope>AllTasks</ReportScope></Report><KeyPrefixManifestGenerator><SourceBucket>bucket</SourceBucket><Prefix>batch-manifests/</Prefix></KeyPrefixManifestGenerator><Priority>10</Priority><RoleArn>acs:ram::1234567890:role/oss-sdk-batch-test</RoleArn></CreateJobRequest>")
		},
		&CreateJobRequest{
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
		},
		func(t *testing.T, o *CreateJobResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateJobRequest><Operation><PutObjectTagging><TagSet><Tag><Key>key1</Key><Value>value1</Value></Tag></TagSet></PutObjectTagging></Operation><ClientRequestToken>unique-token-123</ClientRequestToken><Manifest><Location><ETag>d41d8cd98f00b204e9800998ecf8427e</ETag><Bucket>manifest-bucket</Bucket><Object>manifest.csv</Object></Location><Spec><Fields>Bucket,Key</Fields><Format>OSS_BatchOperations_CSV_20250611</Format></Spec></Manifest><Description>批量设置对象标签任务</Description><Priority>10</Priority><RoleArn>acs:ram::1234567890:role/oss-sdk-batch-test</RoleArn></CreateJobRequest>")
		},
		&CreateJobRequest{
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
		},
		func(t *testing.T, o *CreateJobResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockCreateJob_Success(t *testing.T) {
	for _, c := range testMockCreateJobSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateJob(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateJobErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateJobRequest
	CheckOutputFn  func(t *testing.T, o *CreateJobResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateJobRequest><Operation><DeleteObjectTagging></DeleteObjectTagging></Operation><Report><Bucket>bucket</Bucket><Enabled>true</Enabled><Prefix>batch-reports/</Prefix><ReportScope>AllTasks</ReportScope></Report><KeyPrefixManifestGenerator><SourceBucket>bucket</SourceBucket><Prefix>batch-manifests/</Prefix></KeyPrefixManifestGenerator><Priority>10</Priority><RoleArn>acs:ram::1234567890:role/oss-sdk-batch-test</RoleArn></CreateJobRequest>")

		},
		&CreateJobRequest{
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
		},
		func(t *testing.T, o *CreateJobResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateJobRequest><Operation><DeleteObjectTagging></DeleteObjectTagging></Operation><Report><Bucket>bucket</Bucket><Enabled>true</Enabled><Prefix>batch-reports/</Prefix><ReportScope>AllTasks</ReportScope></Report><KeyPrefixManifestGenerator><SourceBucket>bucket</SourceBucket><Prefix>batch-manifests/</Prefix></KeyPrefixManifestGenerator><Priority>10</Priority><RoleArn>acs:ram::1234567890:role/oss-sdk-batch-test</RoleArn></CreateJobRequest>")
		},
		&CreateJobRequest{
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
		},
		func(t *testing.T, o *CreateJobResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockCreateJob_Error(t *testing.T) {
	for _, c := range testMockCreateJobErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateJob(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDescribeJobSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DescribeJobRequest
	CheckOutputFn  func(t *testing.T, o *DescribeJobResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob&batchJobId=test-batch-job-id", strUrl)
		},
		&DescribeJobRequest{
			BatchJobId: Ptr("test-batch-job-id"),
		},
		func(t *testing.T, o *DescribeJobResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDescribeJob_Success(t *testing.T) {
	for _, c := range testMockDescribeJobSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DescribeJob(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDescribeJobErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DescribeJobRequest
	CheckOutputFn  func(t *testing.T, o *DescribeJobResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob&batchJobId=test-batch-job-id", strUrl)

		},
		&DescribeJobRequest{
			BatchJobId: Ptr("test-batch-job-id"),
		},
		func(t *testing.T, o *DescribeJobResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob&batchJobId=test-batch-job-id", strUrl)
		},
		&DescribeJobRequest{
			BatchJobId: Ptr("test-batch-job-id"),
		},
		func(t *testing.T, o *DescribeJobResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDescribeJob_Error(t *testing.T) {
	for _, c := range testMockDescribeJobErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DescribeJob(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListJobsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListJobsRequest
	CheckOutputFn  func(t *testing.T, o *ListJobsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob", strUrl)
		},
		&ListJobsRequest{},
		func(t *testing.T, o *ListJobsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob&batchJobStatuses=Active+%7C+Cancelled+%7C+Cancelling&continuation-token=token&max-keys=100", strUrl)
		},
		&ListJobsRequest{
			BatchJobStatuses:  Ptr("Active | Cancelled | Cancelling"),
			MaxKeys:           int32(100),
			ContinuationToken: Ptr("token"),
		},
		func(t *testing.T, o *ListJobsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockListJobs_Success(t *testing.T) {
	for _, c := range testMockListJobsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListJobs(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListJobsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListJobsRequest
	CheckOutputFn  func(t *testing.T, o *ListJobsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob", strUrl)

		},
		&ListJobsRequest{},
		func(t *testing.T, o *ListJobsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJob", strUrl)
		},
		&ListJobsRequest{},
		func(t *testing.T, o *ListJobsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockListJobs_Error(t *testing.T) {
	for _, c := range testMockListJobsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListJobs(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateJobPrioritySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateJobPriorityRequest
	CheckOutputFn  func(t *testing.T, o *UpdateJobPriorityResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<UpdateJobPriorityResult>
   <JobId>MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=</JobId>
   <Priority>20</Priority>
</UpdateJobPriorityResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJobId=test-batch-job-id&batchJobPriority&targetPriority=20", strUrl)
		},
		&UpdateJobPriorityRequest{
			BatchJobId:     Ptr("test-batch-job-id"),
			TargetPriority: Ptr(int32(20)),
		},
		func(t *testing.T, o *UpdateJobPriorityResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=", *o.JobId)
			assert.Equal(t, int32(20), *o.Priority)
		},
	},
}

func TestMockUpdateJobPriority_Success(t *testing.T) {
	for _, c := range testMockUpdateJobPrioritySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateJobPriority(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateJobPriorityErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateJobPriorityRequest
	CheckOutputFn  func(t *testing.T, o *UpdateJobPriorityResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJobId=test-batch-job-id&batchJobPriority&targetPriority=20", strUrl)

		},
		&UpdateJobPriorityRequest{
			BatchJobId:     Ptr("test-batch-job-id"),
			TargetPriority: Ptr(int32(20)),
		},
		func(t *testing.T, o *UpdateJobPriorityResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJobId=test-batch-job-id&batchJobPriority&targetPriority=20", strUrl)
		},
		&UpdateJobPriorityRequest{
			BatchJobId:     Ptr("test-batch-job-id"),
			TargetPriority: Ptr(int32(20)),
		},
		func(t *testing.T, o *UpdateJobPriorityResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockUpdateJobPriority_Error(t *testing.T) {
	for _, c := range testMockUpdateJobPriorityErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateJobPriority(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateJobStatusSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateJobStatusRequest
	CheckOutputFn  func(t *testing.T, o *UpdateJobStatusResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<UpdateJobStatusResult>
   <JobId>MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=</JobId>
   <Status>Cancelling</Status>
   <StatusUpdateReason>User requested cancellation</StatusUpdateReason>
</UpdateJobStatusResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJobId=test-batch-job-id&batchJobStatus&requestedJobStatus=Cancelled&statusUpdateReason=User+requested+cancellation", strUrl)
		},
		&UpdateJobStatusRequest{
			BatchJobId:         Ptr("test-batch-job-id"),
			RequestedJobStatus: Ptr("Cancelled"),
			StatusUpdateReason: Ptr("User requested cancellation"),
		},
		func(t *testing.T, o *UpdateJobStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "MzRjZGU2NGQ3YTY5NGRhMTkxZmZhYzY5OTM5YTcxYWU=", *o.JobId)
			assert.Equal(t, "Cancelling", *o.JobStatus)
			assert.Equal(t, "User requested cancellation", *o.StatusUpdateReason)
		},
	},
}

func TestMockUpdateJobStatus_Success(t *testing.T) {
	for _, c := range testMockUpdateJobStatusSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateJobStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateJobStatusErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateJobStatusRequest
	CheckOutputFn  func(t *testing.T, o *UpdateJobStatusResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJobId=test-batch-job-id&batchJobStatus&requestedJobStatus=Cancelled", strUrl)

		},
		&UpdateJobStatusRequest{
			BatchJobId:         Ptr("test-batch-job-id"),
			RequestedJobStatus: Ptr("Cancelled"),
		},
		func(t *testing.T, o *UpdateJobStatusResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?batchJobId=test-batch-job-id&batchJobStatus&requestedJobStatus=Cancelled", strUrl)
		},
		&UpdateJobStatusRequest{
			BatchJobId:         Ptr("test-batch-job-id"),
			RequestedJobStatus: Ptr("Cancelled"),
		},
		func(t *testing.T, o *UpdateJobStatusResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockUpdateJobStatus_Error(t *testing.T) {
	for _, c := range testMockUpdateJobStatusErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateJobStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
