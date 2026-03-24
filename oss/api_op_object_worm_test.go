package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutObjectRetention(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutObjectRetentionRequest
	var input *OperationInput
	var err error

	request = &PutObjectRetentionRequest{}
	input = &OperationInput{
		OpName: "PutObjectRetention",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"retention": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"retention"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutObjectRetentionRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "PutObjectRetention",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"retention": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"retention"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key.")

	request = &PutObjectRetentionRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}
	input = &OperationInput{
		OpName: "PutObjectRetention",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"retention": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"retention"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Retention.")

	request = &PutObjectRetentionRequest{
		Bucket:                    Ptr("bucket"),
		Key:                       Ptr("key"),
		VersionId:                 Ptr("123"),
		BypassGovernanceRetention: Ptr(true),
		Retention: &ObjectWormRetention{
			Mode:            Ptr("GOVERNANC"),
			RetainUntilDate: Ptr("2025-11-10T16:00:00.000Z"),
		},
	}
	input = &OperationInput{
		OpName: "PutObjectRetention",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"retention": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"retention"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["retention"], "")
	assert.Equal(t, input.Parameters["versionId"], "123")
	assert.Equal(t, input.Headers["x-oss-bypass-governance-retention"], "true")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<Retention><Mode>GOVERNANC</Mode><RetainUntilDate>2025-11-10T16:00:00.000Z</RetainUntilDate></Retention>")
}

func TestUnmarshalOutput_PutObjectRetention(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutObjectRetentionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	body := `<?xml version="1.0" encoding="UTF-8"?>
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
	result = &PutObjectRetentionResult{}
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
	result = &PutObjectRetentionResult{}
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
	result = &PutObjectRetentionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetObjectRetention(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetObjectRetentionRequest
	var input *OperationInput
	var err error

	request = &GetObjectRetentionRequest{}
	input = &OperationInput{
		OpName: "GetObjectRetention",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"retention": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"retention"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetObjectRetentionRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "GetObjectRetention",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"retention": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"retention"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key.")

	request = &GetObjectRetentionRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("key"),
		VersionId: Ptr("123"),
	}
	input = &OperationInput{
		OpName: "GetObjectRetention",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"retention": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"retention"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["retention"], "")
	assert.Equal(t, input.Parameters["versionId"], "123")
}

func TestUnmarshalOutput_GetObjectRetention(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<Retention>
  <Mode>COMPLIANCE</Mode>
  <RetainUntilDate>2025-11-10T16:00:00.000Z</RetainUntilDate>
</Retention>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetObjectRetentionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Retention.Mode, "COMPLIANCE")
	assert.Equal(t, *result.Retention.RetainUntilDate, "2025-11-10T16:00:00.000Z")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetObjectRetentionResult{}
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
	result = &GetObjectRetentionResult{}
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
	result = &GetObjectRetentionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_PutObjectLegalHold(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutObjectLegalHoldRequest
	var input *OperationInput
	var err error

	request = &PutObjectLegalHoldRequest{}
	input = &OperationInput{
		OpName: "PutObjectLegalHold",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"legalHold": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"legalHold"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutObjectLegalHoldRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "PutObjectLegalHold",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"legalHold": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"legalHold"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key.")

	request = &PutObjectLegalHoldRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}
	input = &OperationInput{
		OpName: "PutObjectLegalHold",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"legalHold": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"legalHold"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, LegalHold.")

	request = &PutObjectLegalHoldRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("key"),
		VersionId: Ptr("123"),
		LegalHold: &ObjectWormLegalHold{
			Status: Ptr("ON"),
		},
	}
	input = &OperationInput{
		OpName: "PutObjectLegalHold",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"legalHold": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"legalHold"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["legalHold"], "")
	assert.Equal(t, input.Parameters["versionId"], "123")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LegalHold><Status>ON</Status></LegalHold>")
}

func TestUnmarshalOutput_PutObjectLegalHold(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutObjectLegalHoldResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	body := `<?xml version="1.0" encoding="UTF-8"?>
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
	result = &PutObjectLegalHoldResult{}
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
	result = &PutObjectLegalHoldResult{}
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
	result = &PutObjectLegalHoldResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetObjectLegalHold(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetObjectLegalHoldRequest
	var input *OperationInput
	var err error

	request = &GetObjectLegalHoldRequest{}
	input = &OperationInput{
		OpName: "GetObjectLegalHold",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"legalHold": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"legalHold"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetObjectLegalHoldRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "GetObjectLegalHold",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"legalHold": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"legalHold"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key.")

	request = &GetObjectLegalHoldRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("key"),
		VersionId: Ptr("123"),
	}
	input = &OperationInput{
		OpName: "GetObjectLegalHold",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"legalHold": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"legalHold"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["legalHold"], "")
	assert.Equal(t, input.Parameters["versionId"], "123")
}

func TestUnmarshalOutput_GetObjectLegalHold(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<LegalHold>
   <Status>ON</Status>
</LegalHold>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetObjectLegalHoldResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.LegalHold.Status, "ON")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetObjectLegalHoldResult{}
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
	result = &GetObjectLegalHoldResult{}
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
	result = &GetObjectLegalHoldResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
