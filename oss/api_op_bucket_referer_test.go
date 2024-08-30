package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketReferer(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketRefererRequest
	var input *OperationInput
	var err error

	request = &PutBucketRefererRequest{}
	input = &OperationInput{
		OpName: "PutBucketReferer",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"referer": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"referer"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutBucketRefererRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucketReferer",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"referer": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"referer"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, RefererConfiguration.")

	request = &PutBucketRefererRequest{
		Bucket: Ptr("oss-demo"),
		RefererConfiguration: &RefererConfiguration{
			AllowEmptyReferer:        Ptr(false),
			AllowTruncateQueryString: Ptr(true),
			TruncatePath:             Ptr(true),
			RefererList: &RefererList{
				[]string{
					"http://www.aliyun.com", "https://www.aliyun.com", "http://www.*.com", "https://www.?.aliyuncs.com",
				},
			},
			RefererBlacklist: &RefererBlacklist{
				[]string{
					"http://www.refuse.com", "https://*.hack.com", "http://ban.*.com", "https://www.?.deny.com",
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketReferer",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"referer": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"referer"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<RefererConfiguration><AllowEmptyReferer>false</AllowEmptyReferer><AllowTruncateQueryString>true</AllowTruncateQueryString><TruncatePath>true</TruncatePath><RefererList><Referer>http://www.aliyun.com</Referer><Referer>https://www.aliyun.com</Referer><Referer>http://www.*.com</Referer><Referer>https://www.?.aliyuncs.com</Referer></RefererList><RefererBlacklist><Referer>http://www.refuse.com</Referer><Referer>https://*.hack.com</Referer><Referer>http://ban.*.com</Referer><Referer>https://www.?.deny.com</Referer></RefererBlacklist></RefererConfiguration>")

	request = &PutBucketRefererRequest{
		Bucket: Ptr("oss-demo"),
		RefererConfiguration: &RefererConfiguration{
			AllowEmptyReferer:        Ptr(false),
			AllowTruncateQueryString: Ptr(true),
			TruncatePath:             Ptr(true),
			RefererList: &RefererList{
				[]string{
					"http://www.aliyun.com", "https://www.aliyun.com", "http://www.*.com", "https://www.?.aliyuncs.com",
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketReferer",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"referer": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"referer"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<RefererConfiguration><AllowEmptyReferer>false</AllowEmptyReferer><AllowTruncateQueryString>true</AllowTruncateQueryString><TruncatePath>true</TruncatePath><RefererList><Referer>http://www.aliyun.com</Referer><Referer>https://www.aliyun.com</Referer><Referer>http://www.*.com</Referer><Referer>https://www.?.aliyuncs.com</Referer></RefererList></RefererConfiguration>")

	request = &PutBucketRefererRequest{
		Bucket: Ptr("oss-demo"),
		RefererConfiguration: &RefererConfiguration{
			AllowEmptyReferer: Ptr(false),
			RefererList: &RefererList{
				Referers: []string{""},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketReferer",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"referer": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"referer"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<RefererConfiguration><AllowEmptyReferer>false</AllowEmptyReferer><RefererList><Referer></Referer></RefererList></RefererConfiguration>")
}

func TestUnmarshalOutput_PutBucketReferer(t *testing.T) {
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
	result := &PutBucketRefererResult{}
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
	result = &PutBucketRefererResult{}
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
	result = &PutBucketRefererResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	result = &PutBucketRefererResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketReferer(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketRefererRequest
	var input *OperationInput
	var err error

	request = &GetBucketRefererRequest{}
	input = &OperationInput{
		OpName: "GetBucketReferer",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"referer": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"referer"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketRefererRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketReferer",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"referer": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"referer"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketReferer(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	putBody := `<RefererConfiguration>
  <AllowEmptyReferer>true</AllowEmptyReferer>
</RefererConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result := &GetBucketRefererResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.True(t, *result.RefererConfiguration.AllowEmptyReferer)

	putBody = `<RefererConfiguration>
  <AllowEmptyReferer>false</AllowEmptyReferer>
  <AllowTruncateQueryString>true</AllowTruncateQueryString>
  <TruncatePath>true</TruncatePath>
  <RefererList>
    <Referer>http://www.aliyun.com</Referer>
    <Referer>https://www.aliyun.com</Referer>
    <Referer>http://www.*.com</Referer>
    <Referer>https://www.?.aliyuncs.com</Referer>
  </RefererList>
</RefererConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketRefererResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.False(t, *result.RefererConfiguration.AllowEmptyReferer)
	assert.True(t, *result.RefererConfiguration.AllowTruncateQueryString)
	assert.True(t, *result.RefererConfiguration.TruncatePath)
	assert.Equal(t, len(result.RefererConfiguration.RefererList.Referers), 4)

	assert.Equal(t, result.RefererConfiguration.RefererList.Referers[3], "https://www.?.aliyuncs.com")

	putBody = `<RefererConfiguration>
  <AllowEmptyReferer>false</AllowEmptyReferer>
  <AllowTruncateQueryString>true</AllowTruncateQueryString>
  <TruncatePath>true</TruncatePath>
  <RefererList>
    <Referer>http://www.aliyun.com</Referer>
    <Referer>https://www.aliyun.com</Referer>
    <Referer>http://www.*.com</Referer>
    <Referer>https://www.?.aliyuncs.com</Referer>
  </RefererList>
  <RefererBlacklist>
    <Referer>http://www.refuse.com</Referer>
    <Referer>https://*.hack.com</Referer>
    <Referer>http://ban.*.com</Referer>
    <Referer>https://www.?.deny.com</Referer>
  </RefererBlacklist>
</RefererConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketRefererResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.False(t, *result.RefererConfiguration.AllowEmptyReferer)
	assert.True(t, *result.RefererConfiguration.AllowTruncateQueryString)
	assert.True(t, *result.RefererConfiguration.TruncatePath)
	assert.Equal(t, len(result.RefererConfiguration.RefererList.Referers), 4)
	assert.Equal(t, len(result.RefererConfiguration.RefererBlacklist.Referers), 4)
	assert.Equal(t, result.RefererConfiguration.RefererList.Referers[3], "https://www.?.aliyuncs.com")
	assert.Equal(t, result.RefererConfiguration.RefererBlacklist.Referers[2], "http://ban.*.com")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>66C2FF09FDF07830343C72EC</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
		<EC>0015-00000101</EC>
	</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketRefererResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketRefererResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
