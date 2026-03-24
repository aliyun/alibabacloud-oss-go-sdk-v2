package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketObjectWormConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketObjectWormConfigurationRequest
	var input *OperationInput
	var err error

	request = &PutBucketObjectWormConfigurationRequest{}
	input = &OperationInput{
		OpName: "PutBucketObjectWormConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"objectWorm": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"objectWorm"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutBucketObjectWormConfigurationRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucketObjectWormConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"objectWorm": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"objectWorm"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, ObjectWormConfiguration.")

	request = &PutBucketObjectWormConfigurationRequest{
		Bucket: Ptr("oss-demo"),
		ObjectWormConfiguration: &ObjectWormConfiguration{
			ObjectWormEnabled: Ptr("Enabled"),
			Rule: &ObjectWormRule{
				DefaultRetention: &ObjectWormDefaultRetention{
					Mode: Ptr("COMPLIANCE"),
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketObjectWormConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"objectWorm": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"objectWorm"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<ObjectWormConfiguration><ObjectWormEnabled>Enabled</ObjectWormEnabled><Rule><DefaultRetention><Mode>COMPLIANCE</Mode></DefaultRetention></Rule></ObjectWormConfiguration>")

	request = &PutBucketObjectWormConfigurationRequest{
		Bucket: Ptr("oss-demo"),
		ObjectWormConfiguration: &ObjectWormConfiguration{
			ObjectWormEnabled: Ptr("Enabled"),
			Rule: &ObjectWormRule{
				DefaultRetention: &ObjectWormDefaultRetention{
					Mode: Ptr("COMPLIANCE"),
					Days: Ptr(int32(1)),
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketObjectWormConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"objectWorm": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"objectWorm"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<ObjectWormConfiguration><ObjectWormEnabled>Enabled</ObjectWormEnabled><Rule><DefaultRetention><Mode>COMPLIANCE</Mode><Days>1</Days></DefaultRetention></Rule></ObjectWormConfiguration>")

	request = &PutBucketObjectWormConfigurationRequest{
		Bucket: Ptr("oss-demo"),
		ObjectWormConfiguration: &ObjectWormConfiguration{
			ObjectWormEnabled: Ptr("Enabled"),
			Rule: &ObjectWormRule{
				DefaultRetention: &ObjectWormDefaultRetention{
					Mode:  Ptr("GOVERNANCE"),
					Years: Ptr(int32(1)),
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketObjectWormConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"objectWorm": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"objectWorm"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<ObjectWormConfiguration><ObjectWormEnabled>Enabled</ObjectWormEnabled><Rule><DefaultRetention><Mode>GOVERNANCE</Mode><Years>1</Years></DefaultRetention></Rule></ObjectWormConfiguration>")
}

func TestUnmarshalOutput_PutBucketObjectWormConfiguration(t *testing.T) {
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
	result := &PutBucketObjectWormConfigurationResult{}
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
	result = &PutBucketObjectWormConfigurationResult{}
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
	result = &PutBucketObjectWormConfigurationResult{}
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
	result = &PutBucketObjectWormConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketObjectWormConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketObjectWormConfigurationRequest
	var input *OperationInput
	var err error

	request = &GetBucketObjectWormConfigurationRequest{}
	input = &OperationInput{
		OpName: "GetBucketObjectWormConfiguration",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"objectWorm": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"objectWorm"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketObjectWormConfigurationRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketObjectWormConfiguration",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"objectWorm": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"objectWorm"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketObjectWormConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ObjectWormConfiguration>
  <ObjectWormEnabled>Enabled</ObjectWormEnabled>
  <Rule>
    <DefaultRetention>
      <Mode>COMPLIANCE</Mode>
      <Days>1</Days>
    </DefaultRetention>
  </Rule>
</ObjectWormConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetBucketObjectWormConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.ObjectWormConfiguration.ObjectWormEnabled, "Enabled")
	assert.Equal(t, *result.ObjectWormConfiguration.Rule.DefaultRetention.Mode, "COMPLIANCE")
	assert.Equal(t, *result.ObjectWormConfiguration.Rule.DefaultRetention.Days, int32(1))

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ObjectWormConfiguration>
  <ObjectWormEnabled>Enabled</ObjectWormEnabled>
  <Rule>
    <DefaultRetention>
      <Mode>GOVERNANCE</Mode>
      <Years>2</Years>
    </DefaultRetention>
  </Rule>
</ObjectWormConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &GetBucketObjectWormConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.ObjectWormConfiguration.ObjectWormEnabled, "Enabled")
	assert.Equal(t, *result.ObjectWormConfiguration.Rule.DefaultRetention.Mode, "GOVERNANCE")
	assert.Equal(t, *result.ObjectWormConfiguration.Rule.DefaultRetention.Years, int32(2))

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketObjectWormConfigurationResult{}
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
	result = &GetBucketObjectWormConfigurationResult{}
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
	result = &GetBucketObjectWormConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
