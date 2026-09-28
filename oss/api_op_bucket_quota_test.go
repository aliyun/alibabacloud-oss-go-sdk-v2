package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketStorageQuota(t *testing.T) {
	c := Client{}
	request := &PutBucketStorageQuotaRequest{}
	input := &OperationInput{OpName: "PutBucketStorageQuota", Method: "PUT", Parameters: map[string]string{"quota": ""}, Bucket: request.Bucket}
	err := c.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	quota := int64(10737418240)
	request = &PutBucketStorageQuotaRequest{
		Bucket: Ptr("bucket-space"),
		QuotaConfiguration: &QuotaConfiguration{
			StorageQuota: &quota,
			Mode:         StorageQuotaModeStrict,
		},
	}
	input = &OperationInput{OpName: "PutBucketStorageQuota", Method: "PUT", Parameters: map[string]string{"quota": ""}, Bucket: request.Bucket}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NoError(t, err)
	assert.Equal(t, "", input.Parameters["quota"])
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, "<QuotaConfiguration><StorageQuota>10737418240</StorageQuota><Mode>Strict</Mode></QuotaConfiguration>", string(body))
}

func TestUnmarshalOutput_PutBucketStorageQuota(t *testing.T) {
	c := Client{}
	output := &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutBucketStorageQuotaResult{}
	err := c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, "534B371674E88A4D8906****", result.Headers.Get("X-Oss-Request-Id"))

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &PutBucketStorageQuotaResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 403, result.StatusCode)
	assert.Equal(t, "AccessDenied", result.Status)
	assert.Equal(t, "534B371674E88A4D8906****", result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, "application/xml", result.Headers.Get("Content-Type"))
}

func TestUnmarshalOutput_GetBucketStorageQuota(t *testing.T) {
	c := Client{}
	body := `<QuotaConfiguration><Mode>Warning</Mode><StorageQuota>10737418240</StorageQuota><CurrentUsage>1049600</CurrentUsage></QuotaConfiguration>`
	output := &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers:    http.Header{"X-Oss-Request-Id": {"request-id"}},
	}
	result := &GetBucketStorageQuotaResult{}
	err := c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.NoError(t, err)
	assert.Equal(t, int64(10737418240), *result.QuotaConfiguration.StorageQuota)
	assert.Equal(t, StorageQuotaModeWarning, result.QuotaConfiguration.Mode)
	assert.Equal(t, int64(1049600), *result.QuotaConfiguration.CurrentUsage)
	assert.Equal(t, "request-id", result.Headers.Get("X-Oss-Request-Id"))
}

func TestMarshalInput_GetBucketStorageQuota(t *testing.T) {
	c := Client{}
	request := &GetBucketStorageQuotaRequest{}
	input := &OperationInput{
		OpName:     "GetBucketStorageQuota",
		Method:     "GET",
		Parameters: map[string]string{"quota": ""},
		Bucket:     request.Bucket,
	}
	err := c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketStorageQuotaRequest{Bucket: Ptr("bucket-space")}
	input = &OperationInput{
		OpName:     "GetBucketStorageQuota",
		Method:     "GET",
		Parameters: map[string]string{"quota": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NoError(t, err)
	assert.Equal(t, "", input.Parameters["quota"])
}

func TestMarshalInput_DeleteBucketStorageQuota(t *testing.T) {
	c := Client{}
	request := &DeleteBucketStorageQuotaRequest{}
	input := &OperationInput{
		OpName:     "DeleteBucketStorageQuota",
		Method:     "DELETE",
		Parameters: map[string]string{"quota": ""},
		Bucket:     request.Bucket,
	}
	err := c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteBucketStorageQuotaRequest{Bucket: Ptr("bucket-space")}
	input = &OperationInput{
		OpName:     "DeleteBucketStorageQuota",
		Method:     "DELETE",
		Parameters: map[string]string{"quota": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NoError(t, err)
	assert.Equal(t, "", input.Parameters["quota"])
}

func TestUnmarshalOutput_DeleteBucketStorageQuota(t *testing.T) {
	c := Client{}
	result := &DeleteBucketStorageQuotaResult{}
	output := &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers:    http.Header{"X-Oss-Request-Id": {"request-id"}},
	}
	err := c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.NoError(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "No Content", result.Status)
	assert.Equal(t, "request-id", result.Headers.Get("X-Oss-Request-Id"))

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers:    http.Header{"X-Oss-Request-Id": {"request-id"}},
	}
	result = &DeleteBucketStorageQuotaResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.NoError(t, err)
	assert.Equal(t, 403, result.StatusCode)
	assert.Equal(t, "AccessDenied", result.Status)
	assert.Equal(t, "request-id", result.Headers.Get("X-Oss-Request-Id"))
}
