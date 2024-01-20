package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_ListBuckets(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListBucketsRequest
	var input *OperationInput
	var err error

	request = &ListBucketsRequest{}
	input = &OperationInput{
		OpName: "ListBuckets",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)

	request = &ListBucketsRequest{
		Marker:          Ptr(""),
		MaxKeys:         10,
		Prefix:          Ptr("/"),
		ResourceGroupId: Ptr("rg-aek27tc********"),
	}
	input = &OperationInput{
		OpName: "ListBuckets",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_ListBuckets(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult>
  <Owner>
    <ID>51264</ID>
    <DisplayName>51264</DisplayName>
  </Owner>
  <Buckets>
    <Bucket>
      <CreationDate>2014-02-17T18:12:43.000Z</CreationDate>
      <ExtranetEndpoint>oss-cn-shanghai.aliyuncs.com</ExtranetEndpoint>
      <IntranetEndpoint>oss-cn-shanghai-internal.aliyuncs.com</IntranetEndpoint>
      <Location>oss-cn-shanghai</Location>
      <Name>app-base-oss</Name>
      <Region>cn-shanghai</Region>
      <StorageClass>Standard</StorageClass>
    </Bucket>
    <Bucket>
      <CreationDate>2014-02-25T11:21:04.000Z</CreationDate>
      <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
      <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
      <Location>oss-cn-hangzhou</Location>
      <Name>mybucket</Name>
      <Region>cn-hangzhou</Region>
      <StorageClass>IA</StorageClass>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>`

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"5374A2880232A65C2300****"},
		},
	}
	result := &ListBucketsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5374A2880232A65C2300****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Owner.DisplayName, "51264")
	assert.Equal(t, *result.Owner.ID, "51264")
	assert.Equal(t, len(result.Buckets), 2)
	assert.Equal(t, *result.Buckets[0].CreationDate, time.Date(2014, time.February, 17, 18, 12, 43, 0, time.UTC))
	assert.Equal(t, *result.Buckets[0].ExtranetEndpoint, "oss-cn-shanghai.aliyuncs.com")
	assert.Equal(t, *result.Buckets[0].IntranetEndpoint, "oss-cn-shanghai-internal.aliyuncs.com")
	assert.Equal(t, *result.Buckets[0].Name, "app-base-oss")
	assert.Equal(t, *result.Buckets[0].Region, "cn-shanghai")
	assert.Equal(t, *result.Buckets[0].StorageClass, "Standard")
	assert.Equal(t, *result.Buckets[1].CreationDate, time.Date(2014, time.February, 25, 11, 21, 04, 0, time.UTC))
	assert.Equal(t, *result.Buckets[1].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.Buckets[1].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.Buckets[1].Name, "mybucket")
	assert.Equal(t, *result.Buckets[1].Region, "cn-hangzhou")
	assert.Equal(t, *result.Buckets[1].StorageClass, "IA")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult>
  <Prefix>my</Prefix>
  <Marker>mybucket</Marker>
  <MaxKeys>10</MaxKeys>
  <IsTruncated>true</IsTruncated>
  <NextMarker>mybucket10</NextMarker>
  <Owner>
    <ID>ut_test_put_bucket</ID>
    <DisplayName>ut_test_put_bucket</DisplayName>
  </Owner>
  <Buckets>
    <Bucket>
      <CreationDate>2014-05-14T11:18:32.000Z</CreationDate>
      <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
      <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
      <Location>oss-cn-hangzhou</Location>
      <Name>mybucket01</Name>
      <Region>cn-hangzhou</Region>
      <StorageClass>Standard</StorageClass>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>`

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"5374A2880232A65C2300****"},
		},
	}
	result = &ListBucketsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5374A2880232A65C2300****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Owner.DisplayName, "ut_test_put_bucket")
	assert.Equal(t, *result.Owner.ID, "ut_test_put_bucket")
	assert.Equal(t, *result.Prefix, "my")
	assert.Equal(t, *result.Marker, "mybucket")
	assert.Equal(t, result.MaxKeys, int32(10))
	assert.Equal(t, result.IsTruncated, true)
	assert.Equal(t, *result.NextMarker, "mybucket10")
	assert.Equal(t, len(result.Buckets), 1)

	assert.Equal(t, *result.Buckets[0].CreationDate, time.Date(2014, time.May, 14, 11, 18, 32, 0, time.UTC))
	assert.Equal(t, *result.Buckets[0].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.Buckets[0].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.Buckets[0].Name, "mybucket01")
	assert.Equal(t, *result.Buckets[0].Region, "cn-hangzhou")
	assert.Equal(t, *result.Buckets[0].StorageClass, "Standard")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
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
	resultErr := &ListBucketsResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}
