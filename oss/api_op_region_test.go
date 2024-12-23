package oss

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_DescribeRegions(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DescribeRegionsRequest
	var input *OperationInput
	var err error

	request = &DescribeRegionsRequest{}
	input = &OperationInput{
		OpName: "DescribeRegions",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"regions": "",
		},
	}

	input.OpMetadata.Set(signer.SubResource, []string{"regions"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["regions"], "")

	request = &DescribeRegionsRequest{
		Regions: Ptr("oss-cn-hangzhou"),
	}
	input = &OperationInput{
		OpName: "DescribeRegions",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"regions": "",
		},
	}

	input.OpMetadata.Set(signer.SubResource, []string{"regions"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["regions"], "oss-cn-hangzhou")
}

func TestUnmarshalOutput_DescribeRegions(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<RegionInfoList>
  <RegionInfo>
     <Region>oss-cn-hangzhou</Region>
     <InternetEndpoint>oss-cn-hangzhou.aliyuncs.com</InternetEndpoint>
     <InternalEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</InternalEndpoint>
     <AccelerateEndpoint>oss-accelerate.aliyuncs.com</AccelerateEndpoint>  
  </RegionInfo>
  <RegionInfo>
     <Region>oss-cn-shanghai</Region>
     <InternetEndpoint>oss-cn-shanghai.aliyuncs.com</InternetEndpoint>
     <InternalEndpoint>oss-cn-shanghai-internal.aliyuncs.com</InternalEndpoint>
     <AccelerateEndpoint>oss-accelerate.aliyuncs.com</AccelerateEndpoint>  
  </RegionInfo>
</RegionInfoList>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
	result := &DescribeRegionsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, len(result.RegionInfoList.RegionInfos), 2)
	assert.Equal(t, *result.RegionInfoList.RegionInfos[0].InternetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.RegionInfoList.RegionInfos[0].InternalEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.RegionInfoList.RegionInfos[0].Region, "oss-cn-hangzhou")
	assert.Equal(t, *result.RegionInfoList.RegionInfos[0].AccelerateEndpoint, "oss-accelerate.aliyuncs.com")

	assert.Equal(t, *result.RegionInfoList.RegionInfos[1].InternetEndpoint, "oss-cn-shanghai.aliyuncs.com")
	assert.Equal(t, *result.RegionInfoList.RegionInfos[1].InternalEndpoint, "oss-cn-shanghai-internal.aliyuncs.com")
	assert.Equal(t, *result.RegionInfoList.RegionInfos[1].Region, "oss-cn-shanghai")
	assert.Equal(t, *result.RegionInfoList.RegionInfos[1].AccelerateEndpoint, "oss-accelerate.aliyuncs.com")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<RegionInfoList>
  <RegionInfo>
     <Region>oss-cn-hangzhou</Region>
     <InternetEndpoint>oss-cn-hangzhou.aliyuncs.com</InternetEndpoint>
     <InternalEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</InternalEndpoint>
     <AccelerateEndpoint>oss-accelerate.aliyuncs.com</AccelerateEndpoint>  
  </RegionInfo>
</RegionInfoList>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
	result = &DescribeRegionsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, len(result.RegionInfoList.RegionInfos), 1)
	assert.Equal(t, *result.RegionInfoList.RegionInfos[0].InternetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.RegionInfoList.RegionInfos[0].InternalEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.RegionInfoList.RegionInfos[0].Region, "oss-cn-hangzhou")
	assert.Equal(t, *result.RegionInfoList.RegionInfos[0].AccelerateEndpoint, "oss-accelerate.aliyuncs.com")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DescribeRegionsResult{}
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
	result = &DescribeRegionsResult{}
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
	result = &DescribeRegionsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
