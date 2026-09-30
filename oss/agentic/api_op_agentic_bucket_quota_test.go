package agentic

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutAgenticBucketStorageQuota(t *testing.T) {
	c := AgenticBucketClient{}
	request := &PutAgenticBucketStorageQuotaRequest{}
	input := &oss.OperationInput{
		OpName:     "PutAgenticBucketStorageQuota",
		Method:     "PUT",
		Parameters: map[string]string{"agenticBucket": "", "quota": ""},
		Bucket:     request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutAgenticBucketStorageQuotaRequest{
		Bucket: oss.Ptr("my-agentic"),
		QuotaConfiguration: &oss.QuotaConfiguration{
			StorageQuota: oss.Ptr(int64(10737418240)),
			Mode:         oss.StorageQuotaModeStrict,
		},
	}
	input = &oss.OperationInput{
		OpName: "PutAgenticBucketStorageQuota",
		Method: "PUT",
		Parameters: map[string]string{
			"agenticBucket": "",
			"quota":         "",
		},
		Bucket: request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NoError(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, "<QuotaConfiguration><StorageQuota>10737418240</StorageQuota><Mode>Strict</Mode></QuotaConfiguration>", string(body))
	assert.Equal(t, "", input.Parameters["agenticBucket"])
	assert.Equal(t, "", input.Parameters["quota"])
}

func TestUnmarshalOutput_PutAgenticBucketStorageQuota(t *testing.T) {
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
	result := &PutAgenticBucketStorageQuotaResult{}
	err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, "534B371674E88A4D8906****", result.Headers.Get("X-Oss-Request-Id"))

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &PutAgenticBucketStorageQuotaResult{}
	err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, 403, result.StatusCode)
	assert.Equal(t, "AccessDenied", result.Status)
	assert.Equal(t, "534B371674E88A4D8906****", result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, "application/xml", result.Headers.Get("Content-Type"))
}

func TestMarshalInput_GetAgenticBucketStorageQuota(t *testing.T) {
	c := AgenticBucketClient{}

	request := &GetAgenticBucketStorageQuotaRequest{}
	input := &oss.OperationInput{
		OpName:     "GetAgenticBucketStorageQuota",
		Method:     "GET",
		Parameters: map[string]string{"agenticBucket": "", "quota": ""},
		Bucket:     request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetAgenticBucketStorageQuotaRequest{Bucket: oss.Ptr("my-agentic")}
	input = &oss.OperationInput{
		OpName:     "GetAgenticBucketStorageQuota",
		Method:     "GET",
		Parameters: map[string]string{"agenticBucket": "", "quota": ""},
		Bucket:     request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NoError(t, err)
	assert.Equal(t, "", input.Parameters["agenticBucket"])
	assert.Equal(t, "", input.Parameters["quota"])
}

func TestUnmarshalOutput_GetAgenticBucketStorageQuota(t *testing.T) {
	c := AgenticBucketClient{}
	body := `<QuotaConfiguration><Mode>Warning</Mode><StorageQuota>10737418240</StorageQuota><CurrentUsage>1049600</CurrentUsage></QuotaConfiguration>`
	output := &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers:    http.Header{"X-Oss-Request-Id": {"request-id"}},
	}
	result := &GetAgenticBucketStorageQuotaResult{}
	err := c.clientImpl.UnmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.NoError(t, err)
	assert.Equal(t, oss.StorageQuotaModeWarning, result.QuotaConfiguration.Mode)
	assert.Equal(t, int64(10737418240), *result.QuotaConfiguration.StorageQuota)
	assert.Equal(t, int64(1049600), *result.QuotaConfiguration.CurrentUsage)
	assert.Equal(t, "request-id", result.Headers.Get("X-Oss-Request-Id"))
}

func TestMarshalInput_DeleteAgenticBucketStorageQuota(t *testing.T) {
	c := AgenticBucketClient{}
	request := &DeleteAgenticBucketStorageQuotaRequest{}
	input := &oss.OperationInput{
		OpName:     "DeleteAgenticBucketStorageQuota",
		Method:     "DELETE",
		Parameters: map[string]string{"agenticBucket": "", "quota": ""},
		Bucket:     request.Bucket,
	}
	err := c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteAgenticBucketStorageQuotaRequest{Bucket: oss.Ptr("my-agentic")}
	input = &oss.OperationInput{
		OpName:     "DeleteAgenticBucketStorageQuota",
		Method:     "DELETE",
		Parameters: map[string]string{"agenticBucket": "", "quota": ""},
		Bucket:     request.Bucket,
	}
	err = c.clientImpl.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NoError(t, err)
	assert.Equal(t, "", input.Parameters["agenticBucket"])
	assert.Equal(t, "", input.Parameters["quota"])
}

func TestUnmarshalOutput_DeleteAgenticBucketStorageQuota(t *testing.T) {
	c := AgenticBucketClient{}

	result := &DeleteAgenticBucketStorageQuotaResult{}
	output := &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers:    http.Header{"X-Oss-Request-Id": {"request-id"}},
	}
	err := c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.NoError(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "No Content", result.Status)
	assert.Equal(t, "request-id", result.Headers.Get("X-Oss-Request-Id"))

	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers:    http.Header{"X-Oss-Request-Id": {"request-id"}},
	}
	err = c.clientImpl.UnmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.NoError(t, err)
	assert.Equal(t, 403, result.StatusCode)
	assert.Equal(t, "AccessDenied", result.Status)
	assert.Equal(t, "request-id", result.Headers.Get("X-Oss-Request-Id"))
}
