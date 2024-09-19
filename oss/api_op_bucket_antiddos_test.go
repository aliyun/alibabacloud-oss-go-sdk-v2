package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_InitUserAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *InitUserAntiDDosInfoRequest
	var input *OperationInput
	var err error

	request = &InitUserAntiDDosInfoRequest{}
	input = &OperationInput{
		OpName: "InitUserAntiDDosInfo",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_InitUserAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":        {"534B371674E88A4D8906****"},
			"x-oss-defender-instance": {"cbcac8d2-4f75-4d6d-9f2e-c3447f73****"},
		},
	}
	result := &InitUserAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeader, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.DefenderInstance, "cbcac8d2-4f75-4d6d-9f2e-c3447f73****")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &InitUserAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeader, unmarshalBodyXmlMix)
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
	result = &InitUserAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeader, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_UpdateUserAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UpdateUserAntiDDosInfoRequest
	var input *OperationInput
	var err error

	request = &UpdateUserAntiDDosInfoRequest{}
	input = &OperationInput{
		OpName: "UpdateUserAntiDDosInfo",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DefenderInstance")

	request = &UpdateUserAntiDDosInfoRequest{
		DefenderInstance: Ptr("cbcac8d2-4f75-4d6d-9f2e-c3447f73****"),
	}
	input = &OperationInput{
		OpName: "UpdateUserAntiDDosInfo",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DefenderStatus")

	request = &UpdateUserAntiDDosInfoRequest{
		DefenderInstance: Ptr("cbcac8d2-4f75-4d6d-9f2e-c3447f73****"),
		DefenderStatus:   Ptr("HaltDefending"),
	}
	input = &OperationInput{
		OpName: "UpdateUserAntiDDosInfo",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_UpdateUserAntiDDosInfo(t *testing.T) {
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
	result := &UpdateUserAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &UpdateUserAntiDDosInfoResult{}
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
	result = &UpdateUserAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetUserAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetUserAntiDDosInfoRequest
	var input *OperationInput
	var err error

	request = &GetUserAntiDDosInfoRequest{}
	input = &OperationInput{
		OpName: "GetUserAntiDDosInfo",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetUserAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<AntiDDOSListConfiguration>    
    <AntiDDOSConfiguration>        
        <InstanceId>cbcac8d2-4f75-4d6d-9f2e-c3447f73****</InstanceId>
        <Owner>114893010724****</Owner> 
        <Ctime>12345667</Ctime>
        <Mtime>12345667</Mtime>
        <ActiveTime>12345680</ActiveTime>
        <Status>Init</Status>
    </AntiDDOSConfiguration>
 </AntiDDOSListConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetUserAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.AntiDDOSConfigurations), 1)
	assert.Equal(t, *result.AntiDDOSConfigurations[0].InstanceId, "cbcac8d2-4f75-4d6d-9f2e-c3447f73****")
	assert.Equal(t, *result.AntiDDOSConfigurations[0].Owner, "114893010724****")
	assert.Equal(t, *result.AntiDDOSConfigurations[0].Ctime, int64(12345667))
	assert.Equal(t, *result.AntiDDOSConfigurations[0].Mtime, int64(12345667))
	assert.Equal(t, *result.AntiDDOSConfigurations[0].ActiveTime, int64(12345680))
	assert.Equal(t, *result.AntiDDOSConfigurations[0].Status, "Init")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetUserAntiDDosInfoResult{}
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
	result = &GetUserAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_InitBucketAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *InitBucketAntiDDosInfoRequest
	var input *OperationInput
	var err error

	request = &InitBucketAntiDDosInfoRequest{}
	input = &OperationInput{
		OpName: "InitBucketAntiDDosInfo",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &InitBucketAntiDDosInfoRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "InitBucketAntiDDosInfo",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DefenderInstance.")

	request = &InitBucketAntiDDosInfoRequest{
		Bucket:           Ptr("bucket"),
		DefenderInstance: Ptr("cbcac8d2-4f75-4d6d-9f2e-c3447f73****"),
	}
	input = &OperationInput{
		OpName: "InitBucketAntiDDosInfo",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DefenderType.")

	request = &InitBucketAntiDDosInfoRequest{
		Bucket:           Ptr("bucket"),
		DefenderInstance: Ptr("cbcac8d2-4f75-4d6d-9f2e-c3447f73****"),
		DefenderType:     Ptr("AntiDDosPremimum"),
	}
	input = &OperationInput{
		OpName: "InitBucketAntiDDosInfo",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_InitBucketAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":        {"534B371674E88A4D8906****"},
			"x-oss-defender-instance": {"cbcac8d2-4f75-4d6d-9f2e-c3447f73****"},
		},
	}
	result := &InitBucketAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeader, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.DefenderInstance, "cbcac8d2-4f75-4d6d-9f2e-c3447f73****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &InitBucketAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeader, unmarshalBodyXmlMix)
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
	result = &InitBucketAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeader, unmarshalBodyXmlMix)
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
	result = &InitBucketAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeader, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_UpdateBucketAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UpdateBucketAntiDDosInfoRequest
	var input *OperationInput
	var err error

	request = &UpdateBucketAntiDDosInfoRequest{}
	input = &OperationInput{
		OpName: "UpdateBucketAntiDDosInfo",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &UpdateBucketAntiDDosInfoRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "UpdateBucketAntiDDosInfo",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DefenderInstance.")

	request = &UpdateBucketAntiDDosInfoRequest{
		Bucket:           Ptr("bucket"),
		DefenderInstance: Ptr("cbcac8d2-4f75-4d6d-9f2e-c3447f73****"),
	}
	input = &OperationInput{
		OpName: "UpdateBucketAntiDDosInfo",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DefenderStatus.")

	request = &UpdateBucketAntiDDosInfoRequest{
		Bucket:           Ptr("bucket"),
		DefenderInstance: Ptr("cbcac8d2-4f75-4d6d-9f2e-c3447f73****"),
		DefenderStatus:   Ptr("Init"),
	}
	input = &OperationInput{
		OpName: "UpdateBucketAntiDDosInfo",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"antiDDos": "",
		},
		Bucket: request.Bucket,
	}

	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_UpdateBucketAntiDDosInfo(t *testing.T) {
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
	result := &UpdateBucketAntiDDosInfoResult{}
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
	result = &UpdateBucketAntiDDosInfoResult{}
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
	result = &UpdateBucketAntiDDosInfoResult{}
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
	result = &UpdateBucketAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListBucketAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListBucketAntiDDosInfoRequest
	var input *OperationInput
	var err error

	request = &ListBucketAntiDDosInfoRequest{}
	input = &OperationInput{
		OpName: "ListBucketAntiDDosInfo",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"bucketAntiDDos": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_ListBucketAntiDDosInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<AntiDDOSListConfiguration> 
  <Marker>nextMarker</Marker>
  <IsTruncated>true</IsTruncated>
  <AntiDDOSConfiguration>      
    <InstanceId>cbcac8d2-4f75-4d6d-9f2e-c3447f73****</InstanceId>  
    <Owner>114893010724****</Owner>  
    <Bucket>examplebucket</Bucket>  
    <Ctime>1626769503</Ctime>  
    <Mtime>1626769840</Mtime>  
    <ActiveTime>1626769845</ActiveTime>  
    <Status>Defending</Status>  
    <Type>AntiDDosPremimum</Type>  
    <Cnames> 
      <Domain>abc1.example.cn</Domain>  
      <Domain>abc2.example.cn</Domain> 
    </Cnames> 
  </AntiDDOSConfiguration>  
  <AntiDDOSConfiguration>      
    <InstanceId>cbcae8u6-4f75-4d6d-9f2e-c3446g89****</InstanceId>  
    <Owner>1148930107246818</Owner>  
    <Bucket>test-antiddos2</Bucket>  
    <Ctime>1626769993</Ctime>  
    <Mtime>1626769993</Mtime>  
    <ActiveTime>0</ActiveTime>  
    <Status>Init</Status>  
    <Type>AntiDDosPremimum</Type>  
    <Cnames> 
      <Domain>abc3.example.cn</Domain>  
      <Domain>abc4.example.cn</Domain> 
    </Cnames> 
  </AntiDDOSConfiguration> 
</AntiDDOSListConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &ListBucketAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.True(t, *result.AntiDDOSListConfiguration.IsTruncated)
	assert.Equal(t, *result.AntiDDOSListConfiguration.Marker, "nextMarker")
	assert.Equal(t, len(result.AntiDDOSListConfiguration.AntiDDOSConfigurations), 2)
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].InstanceId, "cbcac8d2-4f75-4d6d-9f2e-c3447f73****")
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].Owner, "114893010724****")
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].Bucket, "examplebucket")
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].Ctime, int64(1626769503))
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].Mtime, int64(1626769840))
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].ActiveTime, int64(1626769845))
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].Status, "Defending")
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].Type, "AntiDDosPremimum")
	assert.Equal(t, len(result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].Domains), 2)
	assert.Equal(t, result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].Domains[0], "abc1.example.cn")
	assert.Equal(t, result.AntiDDOSListConfiguration.AntiDDOSConfigurations[0].Domains[1], "abc2.example.cn")
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].InstanceId, "cbcae8u6-4f75-4d6d-9f2e-c3446g89****")
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].Owner, "1148930107246818")
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].Bucket, "test-antiddos2")
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].Ctime, int64(1626769993))
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].Mtime, int64(1626769993))
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].ActiveTime, int64(0))
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].Status, "Init")
	assert.Equal(t, *result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].Type, "AntiDDosPremimum")
	assert.Equal(t, len(result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].Domains), 2)
	assert.Equal(t, result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].Domains[0], "abc3.example.cn")
	assert.Equal(t, result.AntiDDOSListConfiguration.AntiDDOSConfigurations[1].Domains[1], "abc4.example.cn")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListBucketAntiDDosInfoResult{}
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
	result = &ListBucketAntiDDosInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
