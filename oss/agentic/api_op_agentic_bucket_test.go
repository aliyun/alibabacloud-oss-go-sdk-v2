package agentic

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateAgenticBucket(t *testing.T) {
	c := AgenticBucketClient{}
	var request *CreateAgenticBucketRequest
	var input *oss.OperationInput
	var err error

	request = &CreateAgenticBucketRequest{}
	input = &oss.OperationInput{
		OpName: "CreateAgenticBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CreateAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input = &oss.OperationInput{
		OpName: "CreateAgenticBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)

	request = &CreateAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
		CreateAgenticBucketConfiguration: &CreateAgenticBucketConfiguration{
			StorageClass:       oss.StorageClassStandard,
			DataRedundancyType: oss.DataRedundancyLRS,
		},
	}
	input = &oss.OperationInput{
		OpName: "CreateAgenticBucket",
		Method: "PUT",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.NotNil(t, input.Body)
}

func TestUnmarshalOutput_CreateAgenticBucket(t *testing.T) {
	c := AgenticBucketClient{}
	var output *oss.OperationOutput
	var err error

	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &CreateAgenticBucketResult{}
	err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, "534B371674E88A4D8906****", result.Headers.Get("X-Oss-Request-Id"))
}

func TestMarshalInput_DeleteAgenticBucket(t *testing.T) {
	c := AgenticBucketClient{}
	var request *DeleteAgenticBucketRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteAgenticBucketRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteAgenticBucket",
		Method: "DELETE",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteAgenticBucket",
		Method: "DELETE",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteAgenticBucket(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteAgenticBucketResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
}

func TestMarshalInput_GetAgenticBucket(t *testing.T) {
	c := AgenticBucketClient{}
	var request *GetAgenticBucketRequest
	var input *oss.OperationInput
	var err error

	request = &GetAgenticBucketRequest{}
	input = &oss.OperationInput{
		OpName: "GetAgenticBucket",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetAgenticBucketRequest{
		Bucket: oss.Ptr("my-agentic"),
	}
	input = &oss.OperationInput{
		OpName: "GetAgenticBucket",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetAgenticBucket(t *testing.T) {
	c := AgenticBucketClient{}
	body := `<?xml version="1.0" encoding="UTF-8"?>
<AgenticBucketInfo>
  <Name>my-agentic-1234567890123456-cn-hangzhou-ab-apsr</Name>
  <Owner>1234567890123456</Owner>
  <Region>cn-hangzhou</Region>
  <StorageClass>Standard</StorageClass>
  <DataRedundancyType>LRS</DataRedundancyType>
  <Status>enabled</Status>
  <BucketResourceType>AgenticBucket</BucketResourceType>
  <CreateTime>2024-01-01T00:00:00.000Z</CreateTime>
  <ACL>private</ACL>
  <PublicAccessBlock>true</PublicAccessBlock>
  <Versioning>Enabled</Versioning>
  <BucketPolicy>{"Version":"1"}</BucketPolicy>
  <ServerSideEncryptionRule>
    <ApplyServerSideEncryptionByDefault>
      <SSEAlgorithm>AES256</SSEAlgorithm>
    </ApplyServerSideEncryptionByDefault>
  </ServerSideEncryptionRule>
</AgenticBucketInfo>`

	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetAgenticBucketResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotNil(t, result.AgenticBucketInfo)
	assert.Equal(t, "my-agentic-1234567890123456-cn-hangzhou-ab-apsr", *result.AgenticBucketInfo.Name)
	assert.Equal(t, "1234567890123456", *result.AgenticBucketInfo.Owner)
	assert.Equal(t, "cn-hangzhou", *result.AgenticBucketInfo.Region)
	assert.Equal(t, "Standard", *result.AgenticBucketInfo.StorageClass)
	assert.Equal(t, "LRS", *result.AgenticBucketInfo.DataRedundancyType)
	assert.Equal(t, "enabled", *result.AgenticBucketInfo.Status)
	assert.Equal(t, "AgenticBucket", *result.AgenticBucketInfo.BucketResourceType)
	assert.Equal(t, "private", *result.AgenticBucketInfo.ACL)
	assert.Equal(t, "Enabled", *result.AgenticBucketInfo.Versioning)
	assert.Equal(t, "AES256", *result.AgenticBucketInfo.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
}

func TestMarshalInput_ListAgenticBuckets(t *testing.T) {
	c := AgenticBucketClient{}
	var request *ListAgenticBucketsRequest
	var input *oss.OperationInput
	var err error

	request = &ListAgenticBucketsRequest{}
	input = &oss.OperationInput{
		OpName: "ListAgenticBuckets",
		Method: "GET",
	}
	err = c.clientImpl.MarshalInput(request, input)
	assert.Nil(t, err)

	maxKeys := 10
	request = &ListAgenticBucketsRequest{
		ContinuationToken: oss.Ptr("token123"),
		MaxKeys:           &maxKeys,
	}
	input = &oss.OperationInput{
		OpName: "ListAgenticBuckets",
		Method: "GET",
	}
	err = c.clientImpl.MarshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, "token123", input.Parameters["continuation-token"])
	assert.Equal(t, "10", input.Parameters["max-keys"])
}

func TestUnmarshalOutput_ListAgenticBuckets(t *testing.T) {
	c := AgenticBucketClient{}
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListAgenticBucketsResult>
  <Region>cn-hangzhou</Region>
  <Owner>1234567890123456</Owner>
  <IsTruncated>true</IsTruncated>
  <ContinuationToken>token1</ContinuationToken>
  <NextContinuationToken>token2</NextContinuationToken>
  <AgenticBuckets>
    <AgenticBucket>
      <Name>agentic-1</Name>
      <StorageClass>Standard</StorageClass>
      <DataRedundancyType>LRS</DataRedundancyType>
      <CreateTime>2024-01-01T00:00:00.000Z</CreateTime>
    </AgenticBucket>
    <AgenticBucket>
      <Name>agentic-2</Name>
      <StorageClass>IA</StorageClass>
      <DataRedundancyType>ZRS</DataRedundancyType>
      <CreateTime>2024-02-01T00:00:00.000Z</CreateTime>
    </AgenticBucket>
  </AgenticBuckets>
</ListAgenticBucketsResult>`

	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &ListAgenticBucketsResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "cn-hangzhou", *result.Region)
	assert.Equal(t, "1234567890123456", *result.Owner)
	assert.True(t, *result.IsTruncated)
	assert.Equal(t, "token1", *result.ContinuationToken)
	assert.Equal(t, "token2", *result.NextContinuationToken)
	assert.Equal(t, 2, len(result.AgenticBuckets))
	assert.Equal(t, "agentic-1", *result.AgenticBuckets[0].Name)
	assert.Equal(t, "Standard", *result.AgenticBuckets[0].StorageClass)
	assert.Equal(t, "agentic-2", *result.AgenticBuckets[1].Name)
	assert.Equal(t, "IA", *result.AgenticBuckets[1].StorageClass)
}

func TestMarshalInput_PutAgenticBucketStatus(t *testing.T) {
	c := AgenticBucketClient{}
	var request *PutAgenticBucketStatusRequest
	var input *oss.OperationInput
	var err error

	request = &PutAgenticBucketStatusRequest{}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketStatus",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutAgenticBucketStatusRequest{
		Bucket: oss.Ptr("my-agentic"),
		AgenticBucketStatus: &AgenticBucketStatus{
			Status: oss.Ptr("enabled"),
		},
	}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketStatus",
		Method: "PUT",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.NotNil(t, input.Body)
}

func TestUnmarshalOutput_PutAgenticBucketStatus(t *testing.T) {
	c := AgenticBucketClient{}
	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutAgenticBucketStatusResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
}

func TestMarshalInput_ListBucketSpaces(t *testing.T) {
	c := AgenticBucketClient{}
	var request *ListBucketSpacesRequest
	var input *oss.OperationInput
	var err error

	request = &ListBucketSpacesRequest{}
	input = &oss.OperationInput{
		OpName: "ListBucketSpaces",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	maxKeys := 20
	request = &ListBucketSpacesRequest{
		Bucket:            oss.Ptr("my-agentic"),
		Prefix:            oss.Ptr("sandbox-"),
		ContinuationToken: oss.Ptr("token1"),
		StartAfter:        oss.Ptr("sandbox-000"),
		MaxKeys:           &maxKeys,
	}
	input = &oss.OperationInput{
		OpName: "ListBucketSpaces",
		Method: "GET",
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, "sandbox-", input.Parameters["prefix"])
	assert.Equal(t, "token1", input.Parameters["continuation-token"])
	assert.Equal(t, "sandbox-000", input.Parameters["start-after"])
	assert.Equal(t, "20", input.Parameters["max-keys"])
}

func TestUnmarshalOutput_ListBucketSpaces(t *testing.T) {
	c := AgenticBucketClient{}
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketSpacesResult>
  <Owner>
    <ID>1234567890123456</ID>
    <DisplayName>owner-name</DisplayName>
  </Owner>
  <Prefix>sandbox-</Prefix>
  <MaxKeys>20</MaxKeys>
  <StartAfter>sandbox-000</StartAfter>
  <IsTruncated>false</IsTruncated>
  <BucketSpaces>
    <BucketSpace>
      <Name>sandbox-001</Name>
      <Location>oss-cn-hangzhou</Location>
      <CreationDate>2024-01-01T00:00:00.000Z</CreationDate>
      <StorageClass>Standard</StorageClass>
    </BucketSpace>
  </BucketSpaces>
</ListBucketSpacesResult>`

	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &ListBucketSpacesResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotNil(t, result.Owner)
	assert.Equal(t, "1234567890123456", *result.Owner.ID)
	assert.Equal(t, "sandbox-", *result.Prefix)
	assert.Equal(t, 20, *result.MaxKeys)
	assert.Equal(t, "sandbox-000", *result.StartAfter)
	assert.False(t, *result.IsTruncated)
	assert.Equal(t, 1, len(result.BucketSpaces))
	assert.Equal(t, "sandbox-001", *result.BucketSpaces[0].Name)
	assert.Equal(t, "oss-cn-hangzhou", *result.BucketSpaces[0].Location)
	assert.Equal(t, "Standard", *result.BucketSpaces[0].StorageClass)
}
