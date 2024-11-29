package oss

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateAccessPointForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CreateAccessPointForObjectProcessRequest
	var input *OperationInput
	var err error

	request = &CreateAccessPointForObjectProcessRequest{}
	input = &OperationInput{
		OpName: "CreateAccessPointForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &CreateAccessPointForObjectProcessRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "CreateAccessPointForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, AccessPointForObjectProcessName.")

	request = &CreateAccessPointForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("fc-ap-01"),
	}
	input = &OperationInput{
		OpName: "CreateAccessPointForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, CreateAccessPointForObjectProcessConfiguration.")

	request = &CreateAccessPointForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		CreateAccessPointForObjectProcessConfiguration: &CreateAccessPointForObjectProcessConfiguration{
			AccessPointName: Ptr("ap-01"),
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: []string{"GetObject-Range"},
				TransformationConfigurations: []TransformationConfiguration{
					{
						Actions: &AccessPointActions{
							[]string{"GetObject"},
						},
						ContentTransformation: &ContentTransformation{
							FunctionArn:           Ptr("acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01"),
							FunctionAssumeRoleArn: Ptr("acs:ram::111933544165****:role/aliyunfcdefaultrole"),
						},
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "CreateAccessPointForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<CreateAccessPointForObjectProcessConfiguration><AccessPointName>ap-01</AccessPointName><ObjectProcessConfiguration><AllowedFeatures><AllowedFeature>GetObject-Range</AllowedFeature></AllowedFeatures><TransformationConfigurations><TransformationConfiguration><Actions><Action>GetObject</Action></Actions><ContentTransformation><FunctionCompute><FunctionAssumeRoleArn>acs:ram::111933544165****:role/aliyunfcdefaultrole</FunctionAssumeRoleArn><FunctionArn>acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01</FunctionArn></FunctionCompute></ContentTransformation></TransformationConfiguration></TransformationConfigurations></ObjectProcessConfiguration></CreateAccessPointForObjectProcessConfiguration>")
}

func TestUnmarshalOutput_CreateAccessPointForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<CreateAccessPointForObjectProcessResult>
  <AccessPointForObjectProcessArn>acs:oss:cn-qingdao:119335441657143:accesspointforobjectprocess/fc-ap-01</AccessPointForObjectProcessArn>
  <AccessPointForObjectProcessAlias>fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias</AccessPointForObjectProcessAlias>
</CreateAccessPointForObjectProcessResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &CreateAccessPointForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.AccessPointForObjectProcessArn, "acs:oss:cn-qingdao:119335441657143:accesspointforobjectprocess/fc-ap-01")
	assert.Equal(t, *result.AccessPointForObjectProcessAlias, "fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateAccessPointForObjectProcessResult{}
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
	result = &CreateAccessPointForObjectProcessResult{}
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
	result = &CreateAccessPointForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetAccessPointForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetAccessPointForObjectProcessRequest
	var input *OperationInput
	var err error

	request = &GetAccessPointForObjectProcessRequest{}
	input = &OperationInput{
		OpName: "GetAccessPointForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetAccessPointForObjectProcessRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetAccessPointForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, AccessPointForObjectProcessName.")

	request = &GetAccessPointForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("fc-ap-01"),
	}
	input = &OperationInput{
		OpName: "GetAccessPointForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-access-point-for-object-process-name"], "fc-ap-01")
}

func TestUnmarshalOutput_GetAccessPointForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<GetAccessPointForObjectProcessResult>
  <AccessPointNameForObjectProcess>fc-ap-01</AccessPointNameForObjectProcess>
  <AccessPointForObjectProcessAlias>fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias</AccessPointForObjectProcessAlias>
  <AccessPointName>ap-01</AccessPointName>
  <AccountId>111933544165****</AccountId>
  <AccessPointForObjectProcessArn>acs:oss:cn-qingdao:11933544165****:accesspointforobjectprocess/fc-ap-01</AccessPointForObjectProcessArn>
  <CreationDate>1626769503</CreationDate>
  <Status>enable</Status>
  <Endpoints>
    <PublicEndpoint>fc-ap-01-111933544165****.oss-cn-qingdao.oss-object-process.aliyuncs.com</PublicEndpoint>
    <InternalEndpoint>fc-ap-01-111933544165****.oss-cn-qingdao-internal.oss-object-process.aliyuncs.com</InternalEndpoint>
  </Endpoints>
  <PublicAccessBlockConfiguration>
    <BlockPublicAccess>true</BlockPublicAccess>
  </PublicAccessBlockConfiguration>
</GetAccessPointForObjectProcessResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetAccessPointForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.AccessPointNameForObjectProcess, "fc-ap-01")
	assert.Equal(t, *result.AccessPointForObjectProcessAlias, "fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias")
	assert.Equal(t, *result.AccessPointName, "ap-01")
	assert.Equal(t, *result.AccountId, "111933544165****")
	assert.Equal(t, *result.AccessPointForObjectProcessArn, "acs:oss:cn-qingdao:11933544165****:accesspointforobjectprocess/fc-ap-01")
	assert.Equal(t, *result.CreationDate, "1626769503")
	assert.Equal(t, *result.AccessPointForObjectProcessStatus, "enable")
	assert.Equal(t, *result.Endpoints.PublicEndpoint, "fc-ap-01-111933544165****.oss-cn-qingdao.oss-object-process.aliyuncs.com")
	assert.Equal(t, *result.Endpoints.InternalEndpoint, "fc-ap-01-111933544165****.oss-cn-qingdao-internal.oss-object-process.aliyuncs.com")
	assert.True(t, *result.PublicAccessBlockConfiguration.BlockPublicAccess)
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetAccessPointForObjectProcessResult{}
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
	result = &GetAccessPointForObjectProcessResult{}
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
	result = &GetAccessPointForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteAccessPointForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteAccessPointForObjectProcessRequest
	var input *OperationInput
	var err error

	request = &DeleteAccessPointForObjectProcessRequest{}
	input = &OperationInput{
		OpName: "DeleteAccessPointForObjectProcess",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteAccessPointForObjectProcessRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteAccessPointForObjectProcess",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, AccessPointForObjectProcessName.")

	request = &DeleteAccessPointForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("fc-ap-01"),
	}
	input = &OperationInput{
		OpName: "DeleteAccessPointForObjectProcess",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-access-point-for-object-process-name"], "fc-ap-01")
}

func TestUnmarshalOutput_DeleteAccessPointForObjectProcess(t *testing.T) {
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
	result := &DeleteAccessPointForObjectProcessResult{}
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
	result = &DeleteAccessPointForObjectProcessResult{}
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
	result = &DeleteAccessPointForObjectProcessResult{}
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
	result = &DeleteAccessPointForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListAccessPointsForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListAccessPointsForObjectProcessRequest
	var input *OperationInput
	var err error

	request = &ListAccessPointsForObjectProcessRequest{}
	input = &OperationInput{
		OpName: "ListAccessPointsForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Nil(t, input.Bucket)

	request = &ListAccessPointsForObjectProcessRequest{
		MaxKeys:           int64(10),
		ContinuationToken: Ptr("token"),
	}
	input = &OperationInput{
		OpName: "ListAccessPointsForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointForObjectProcess": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["max-keys"], "10")
	assert.Equal(t, input.Parameters["continuation-token"], "token")
}

func TestUnmarshalOutput_ListAccessPointsForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListAccessPointsForObjectProcessResult>
   <IsTruncated>true</IsTruncated>
   <NextContinuationToken>abc</NextContinuationToken>
   <AccountId>111933544165****</AccountId>
   <AccessPointsForObjectProcess>
      <AccessPointForObjectProcess>
          <AccessPointNameForObjectProcess>fc-ap-01</AccessPointNameForObjectProcess>
          <AccessPointForObjectProcessAlias>fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias</AccessPointForObjectProcessAlias>
          <AccessPointName>fc-01</AccessPointName>
          <Status>enable</Status>
      </AccessPointForObjectProcess>
   </AccessPointsForObjectProcess>
</ListAccessPointsForObjectProcessResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &ListAccessPointsForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.NextContinuationToken, "abc")
	assert.True(t, *result.IsTruncated)
	assert.Equal(t, *result.AccountId, "111933544165****")
	assert.Equal(t, len(result.AccessPointsForObjectProcess.AccessPointForObjectProcesss), 1)
	assert.Equal(t, *result.AccessPointsForObjectProcess.AccessPointForObjectProcesss[0].AccessPointNameForObjectProcess, "fc-ap-01")
	assert.Equal(t, *result.AccessPointsForObjectProcess.AccessPointForObjectProcesss[0].AccessPointNameForObjectProcess, "fc-ap-01")
	assert.Equal(t, *result.AccessPointsForObjectProcess.AccessPointForObjectProcesss[0].AccessPointForObjectProcessAlias, "fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias")
	assert.Equal(t, *result.AccessPointsForObjectProcess.AccessPointForObjectProcesss[0].AccessPointName, "fc-01")
	assert.Equal(t, *result.AccessPointsForObjectProcess.AccessPointForObjectProcesss[0].Status, "enable")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListAccessPointsForObjectProcessResult{}
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
	result = &ListAccessPointsForObjectProcessResult{}
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
	result = &ListAccessPointsForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_PutAccessPointPolicyForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutAccessPointPolicyForObjectProcessRequest
	var input *OperationInput
	var err error

	request = &PutAccessPointPolicyForObjectProcessRequest{}
	input = &OperationInput{
		OpName: "PutAccessPointPolicyForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutAccessPointPolicyForObjectProcessRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutAccessPointPolicyForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, AccessPointForObjectProcessName.")

	request = &PutAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("ap-01"),
	}
	input = &OperationInput{
		OpName: "PutAccessPointPolicyForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Body.")

	body := `{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["23721626365841xxxx"],"Resource":["acs:oss:cn-qingdao:111933544165xxxx:accesspointforobjectprocess/fc-ap-01/object/*"]}]}`
	request = &PutAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		Body:                            strings.NewReader(body),
	}
	input = &OperationInput{
		OpName: "PutAccessPointPolicyForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	rBody, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(rBody), "{\"Version\":\"1\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"oss:GetObject\"],\"Principal\":[\"23721626365841xxxx\"],\"Resource\":[\"acs:oss:cn-qingdao:111933544165xxxx:accesspointforobjectprocess/fc-ap-01/object/*\"]}]}")
}

func TestUnmarshalOutput_PutAccessPointPolicyForObjectProcess(t *testing.T) {
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
	result := &PutAccessPointPolicyForObjectProcessResult{}
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
	result = &PutAccessPointPolicyForObjectProcessResult{}
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
	result = &PutAccessPointPolicyForObjectProcessResult{}
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
	result = &PutAccessPointPolicyForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetAccessPointPolicyForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetAccessPointPolicyForObjectProcessRequest
	var input *OperationInput
	var err error

	request = &GetAccessPointPolicyForObjectProcessRequest{}
	input = &OperationInput{
		OpName: "GetAccessPointPolicyForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetAccessPointPolicyForObjectProcessRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetAccessPointPolicyForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, AccessPointForObjectProcessName.")

	request = &GetAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("ap-01"),
	}
	input = &OperationInput{
		OpName: "GetAccessPointPolicyForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-access-point-for-object-process-name"], "ap-01")
}

func TestUnmarshalOutput_GetAccessPointPolicyForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `{
   "Version":"1",
   "Statement":[
   {
     "Action":[
       "oss:PutObject",
       "oss:GetObject"
    ],
    "Effect":"Deny",
    "Principal":["27737962156157xxxx"],
    "Resource":[
       "acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/$ap-01",
       "acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/$ap-01/object/*"
     ]
   }
  ]
 }`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	oBody, err := ioutil.ReadAll(output.Body)
	assert.Nil(t, err)
	result := &GetAccessPointPolicyForObjectProcessResult{
		Body: string(oBody),
	}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Body, "{\n   \"Version\":\"1\",\n   \"Statement\":[\n   {\n     \"Action\":[\n       \"oss:PutObject\",\n       \"oss:GetObject\"\n    ],\n    \"Effect\":\"Deny\",\n    \"Principal\":[\"27737962156157xxxx\"],\n    \"Resource\":[\n       \"acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/$ap-01\",\n       \"acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/$ap-01/object/*\"\n     ]\n   }\n  ]\n }")
	body = `<?xml version="1.0" encoding="UTF-8"?>
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
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	oBody, err = ioutil.ReadAll(output.Body)
	assert.Nil(t, err)
	result = &GetAccessPointPolicyForObjectProcessResult{
		Body: string(oBody),
	}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
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
	oBody, err = ioutil.ReadAll(output.Body)
	assert.Nil(t, err)
	result = &GetAccessPointPolicyForObjectProcessResult{
		Body: string(oBody),
	}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteAccessPointPolicyForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteAccessPointPolicyForObjectProcessRequest
	var input *OperationInput
	var err error

	request = &DeleteAccessPointPolicyForObjectProcessRequest{}
	input = &OperationInput{
		OpName: "DeleteAccessPointPolicyForObjectProcess",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteAccessPointPolicyForObjectProcessRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteAccessPointPolicyForObjectProcess",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, AccessPointForObjectProcessName.")

	request = &DeleteAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("ap-01"),
	}
	input = &OperationInput{
		OpName: "DeleteAccessPointPolicyForObjectProcess",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointPolicyForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointPolicyForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-access-point-for-object-process-name"], "ap-01")
}

func TestUnmarshalOutput_DeleteAccessPointPolicyForObjectProcess(t *testing.T) {
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
	result := &DeleteAccessPointPolicyForObjectProcessResult{}
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
	result = &DeleteAccessPointPolicyForObjectProcessResult{}
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
	result = &DeleteAccessPointPolicyForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_PutAccessPointConfigForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutAccessPointConfigForObjectProcessRequest
	var input *OperationInput
	var err error

	request = &PutAccessPointConfigForObjectProcessRequest{}
	input = &OperationInput{
		OpName: "PutAccessPointConfigForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointConfigForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointConfigForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutAccessPointConfigForObjectProcessRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutAccessPointConfigForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointConfigForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointConfigForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, AccessPointForObjectProcessName.")

	request = &PutAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("ap-01"),
	}
	input = &OperationInput{
		OpName: "PutAccessPointConfigForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointConfigForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointConfigForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, PutAccessPointConfigForObjectProcessConfiguration.")

	request = &PutAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("ap-01"),
		PutAccessPointConfigForObjectProcessConfiguration: &PutAccessPointConfigForObjectProcessConfiguration{
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: []string{"GetObject-Range"},
				TransformationConfigurations: []TransformationConfiguration{
					{
						Actions: &AccessPointActions{
							[]string{"GetObject"},
						},
						ContentTransformation: &ContentTransformation{
							FunctionArn:           Ptr("acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01"),
							FunctionAssumeRoleArn: Ptr("acs:ram::111933544165****:role/aliyunfcdefaultrole"),
						},
					},
				},
			},
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				Ptr(true),
			},
		},
	}
	input = &OperationInput{
		OpName: "PutAccessPointConfigForObjectProcess",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointConfigForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointConfigForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	rBody, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(rBody), "<PutAccessPointConfigForObjectProcessConfiguration><PublicAccessBlockConfiguration><BlockPublicAccess>true</BlockPublicAccess></PublicAccessBlockConfiguration><ObjectProcessConfiguration><AllowedFeatures><AllowedFeature>GetObject-Range</AllowedFeature></AllowedFeatures><TransformationConfigurations><TransformationConfiguration><Actions><Action>GetObject</Action></Actions><ContentTransformation><FunctionCompute><FunctionAssumeRoleArn>acs:ram::111933544165****:role/aliyunfcdefaultrole</FunctionAssumeRoleArn><FunctionArn>acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01</FunctionArn></FunctionCompute></ContentTransformation></TransformationConfiguration></TransformationConfigurations></ObjectProcessConfiguration></PutAccessPointConfigForObjectProcessConfiguration>")
}

func TestUnmarshalOutput_PutAccessPointConfigForObjectProcess(t *testing.T) {
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
	result := &PutAccessPointConfigForObjectProcessResult{}
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
	result = &PutAccessPointConfigForObjectProcessResult{}
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
	result = &PutAccessPointConfigForObjectProcessResult{}
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
	result = &PutAccessPointConfigForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetAccessPointConfigForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetAccessPointConfigForObjectProcessRequest
	var input *OperationInput
	var err error

	request = &GetAccessPointConfigForObjectProcessRequest{}
	input = &OperationInput{
		OpName: "GetAccessPointConfigForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointConfigForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointConfigForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetAccessPointConfigForObjectProcessRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetAccessPointConfigForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointConfigForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointConfigForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, AccessPointForObjectProcessName.")

	request = &GetAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr("oss-demo"),
		AccessPointForObjectProcessName: Ptr("ap-01"),
	}
	input = &OperationInput{
		OpName: "GetAccessPointConfigForObjectProcess",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"accessPointConfigForObjectProcess": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"accessPointConfigForObjectProcess"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-access-point-for-object-process-name"], "ap-01")
}

func TestUnmarshalOutput_GetAccessPointConfigForObjectProcess(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<GetAccessPointConfigForObjectProcessResult>
  <ObjectProcessConfiguration>
    <AllowedFeatures/>
    <TransformationConfigurations>
      <TransformationConfiguration>
        <Actions>
          <Action>getobject</Action>
        </Actions>
        <ContentTransformation>
          <FunctionCompute>
            <FunctionAssumeRoleArn>acs:ram::111933544165****:role/aliyunfcdefaultrole</FunctionAssumeRoleArn>
            <FunctionArn>acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01</FunctionArn>
          </FunctionCompute>
        </ContentTransformation>
      </TransformationConfiguration>
    </TransformationConfigurations>
  </ObjectProcessConfiguration>
  <PublicAccessBlockConfiguration>
    <BlockPublicAccess>true</BlockPublicAccess>
  </PublicAccessBlockConfiguration> 
</GetAccessPointConfigForObjectProcessResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result := &GetAccessPointConfigForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.True(t, *result.PublicAccessBlockConfiguration.BlockPublicAccess)
	assert.Equal(t, result.ObjectProcessConfiguration.TransformationConfigurations[0].Actions.Actions[0], "getobject")
	assert.Equal(t, *result.ObjectProcessConfiguration.TransformationConfigurations[0].ContentTransformation.FunctionAssumeRoleArn, "acs:ram::111933544165****:role/aliyunfcdefaultrole")
	assert.Equal(t, *result.ObjectProcessConfiguration.TransformationConfigurations[0].ContentTransformation.FunctionArn, "acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01")
	body = `<?xml version="1.0" encoding="UTF-8"?>
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
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetAccessPointConfigForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
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
	assert.Nil(t, err)
	result = &GetAccessPointConfigForObjectProcessResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_WriteGetObjectResponse(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *WriteGetObjectResponseRequest
	var input *OperationInput
	var err error

	request = &WriteGetObjectResponseRequest{}
	input = &OperationInput{
		OpName: "WriteGetObjectResponse",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"x-oss-write-get-object-response": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"x-oss-write-get-object-response"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, RequestRoute.")

	request = &WriteGetObjectResponseRequest{
		RequestRoute: Ptr("fc-ap-01-176022554508***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com"),
	}
	input = &OperationInput{
		OpName: "WriteGetObjectResponse",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"x-oss-write-get-object-response": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"x-oss-write-get-object-response"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, RequestToken.")

	request = &WriteGetObjectResponseRequest{
		RequestRoute: Ptr("fc-ap-01-176022554508***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com"),
		RequestToken: Ptr("token"),
	}
	input = &OperationInput{
		OpName: "WriteGetObjectResponse",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"x-oss-write-get-object-response": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"x-oss-write-get-object-response"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, FwdStatus.")

	request = &WriteGetObjectResponseRequest{
		RequestRoute: Ptr("fc-ap-01-***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com"),
		RequestToken: Ptr("token"),
		FwdStatus:    Ptr("200"),
	}
	input = &OperationInput{
		OpName: "WriteGetObjectResponse",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"x-oss-write-get-object-response": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"x-oss-write-get-object-response"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-request-route"], "fc-ap-01-***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com")
	assert.Equal(t, input.Headers["x-oss-request-token"], "token")
	assert.Equal(t, input.Headers["x-oss-fwd-status"], "200")
}

func TestUnmarshalOutput_WriteGetObjectResponse(t *testing.T) {
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
	assert.Nil(t, err)
	result := &WriteGetObjectResponseResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	body := `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>534B371674E88A4D8906****</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
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
	result = &WriteGetObjectResponseResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
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
	assert.Nil(t, err)
	result = &WriteGetObjectResponseResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
