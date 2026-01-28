package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketOverwriteConfig(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketOverwriteConfigRequest
	var input *OperationInput
	var err error

	request = &PutBucketOverwriteConfigRequest{}
	input = &OperationInput{
		OpName: "PutBucketOverwriteConfig",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"overwriteConfig": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"overwriteConfig"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutBucketOverwriteConfigRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucketOverwriteConfig",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"overwriteConfig": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"overwriteConfig"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, OverwriteConfiguration.")

	request = &PutBucketOverwriteConfigRequest{
		Bucket: Ptr("oss-demo"),
		OverwriteConfiguration: &OverwriteConfiguration{
			Rules: []OverwriteRule{
				{
					ID:     Ptr("1"),
					Action: Ptr("forbid"),
				},
				{
					ID:     Ptr("2"),
					Action: Ptr("forbid"),
					Prefix: Ptr("pre"),
					Suffix: Ptr(".txt"),
					Principals: &OverwritePrincipals{
						[]string{"1234567890"},
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketOverwriteConfig",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"overwriteConfig": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"overwriteConfig"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "PUT")
	assert.Equal(t, input.Parameters["overwriteConfig"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<OverwriteConfiguration><Rule><ID>1</ID><Action>forbid</Action></Rule><Rule><ID>2</ID><Action>forbid</Action><Prefix>pre</Prefix><Suffix>.txt</Suffix><Principals><Principal>1234567890</Principal></Principals></Rule></OverwriteConfiguration>")
}

func TestUnmarshalOutput_PutBucketOverwriteConfig(t *testing.T) {
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
	result := &PutBucketOverwriteConfigResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketOverwriteConfigResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
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
	result = &PutBucketOverwriteConfigResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketOverwriteConfig(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketOverwriteConfigRequest
	var input *OperationInput
	var err error

	request = &GetBucketOverwriteConfigRequest{}
	input = &OperationInput{
		OpName: "GetBucketOverwriteConfig",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"overwriteConfig": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"overwriteConfig"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketOverwriteConfigRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketOverwriteConfig",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"overwriteConfig": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"overwriteConfig"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "GET")
	assert.Equal(t, input.Parameters["overwriteConfig"], "")
}

func TestUnmarshalOutput_GetBucketOverwriteConfig(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<OverwriteConfiguration>
  <Rule>
    <ID>1</ID>
    <Action>forbid</Action>
    <Prefix></Prefix>
    <Suffix></Suffix>
    <Principals />
  </Rule>
  <Rule>
    <ID>2</ID>
    <Action>forbid</Action>
    <Prefix>pre</Prefix>
    <Suffix>.txt</Suffix>
    <Principals>
      <Principal>1234567890</Principal>
    </Principals>
  </Rule>
</OverwriteConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketOverwriteConfigResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, len(result.OverwriteConfiguration.Rules), 2)

	assert.Equal(t, *result.OverwriteConfiguration.Rules[0].ID, "1")
	assert.Equal(t, *result.OverwriteConfiguration.Rules[0].Action, "forbid")
	assert.Empty(t, result.OverwriteConfiguration.Rules[0].Prefix)
	assert.Empty(t, result.OverwriteConfiguration.Rules[0].Suffix)
	assert.Equal(t, *result.OverwriteConfiguration.Rules[1].ID, "2")
	assert.Equal(t, *result.OverwriteConfiguration.Rules[1].Action, "forbid")
	assert.Equal(t, *result.OverwriteConfiguration.Rules[1].Prefix, "pre")
	assert.Equal(t, *result.OverwriteConfiguration.Rules[1].Suffix, ".txt")
	assert.Equal(t, result.OverwriteConfiguration.Rules[1].Principals.Principals[0], "1234567890")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",

		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketOverwriteConfigResult{}
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
	result = &GetBucketOverwriteConfigResult{}
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
	result = &GetBucketOverwriteConfigResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteBucketOverwriteConfig(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteBucketOverwriteConfigRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketOverwriteConfigRequest{}
	input = &OperationInput{
		OpName: "DeleteBucketOverwriteConfig",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"overwriteConfig": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"overwriteConfig"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteBucketOverwriteConfigRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketOverwriteConfig",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"overwriteConfig": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"overwriteConfig"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "DELETE")
	assert.Equal(t, input.Parameters["overwriteConfig"], "")
}

func TestUnmarshalOutput_DeleteBucketOverwriteConfig(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteBucketOverwriteConfigResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketOverwriteConfigResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
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
	result = &DeleteBucketOverwriteConfigResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
