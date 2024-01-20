package oss

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalOutput_encodetype(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	body := `<?xml version="1.0" encoding="UTF-8"?>
			<ListBucketResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
			<Name>oss-example</Name>
			<Prefix>hello%20world%21</Prefix>
			<Marker>hello%20</Marker>
			<MaxKeys>100</MaxKeys>
			<Delimiter>hello%20%21world</Delimiter>
			<IsTruncated>false</IsTruncated>
			<EncodingType>url</EncodingType>
			<Contents>
				<Key>fun%2Fmovie%2F001.avi</Key>
				<LastModified>2012-02-24T08:43:07.000Z</LastModified>
				<ETag>&quot;5B3C1A2E053D763E1B002CC607C5A0FE&quot;</ETag>
				<Type>Normal</Type>
				<Size>344606</Size>
				<StorageClass>Standard</StorageClass>
				<Owner>
					<ID>00220120222</ID>
					<DisplayName>user-example</DisplayName>
				</Owner>
			</Contents>
			<Contents>
				<Key>fun%2Fmovie%2F007.avi</Key>
				<LastModified>2012-02-24T08:43:27.000Z</LastModified>
				<ETag>&quot;5B3C1A2E053D763E1B002CC607C5A0FE&quot;</ETag>
				<Type>Normal</Type>
				<Size>344606</Size>
				<StorageClass>Standard</StorageClass>
				<Owner>
					<ID>00220120222</ID>
					<DisplayName>user-example</DisplayName>
				</Owner>
			</Contents>
			<Contents>
				<Key>oss.jpg</Key>
				<LastModified>2012-02-24T06:07:48.000Z</LastModified>
				<ETag>&quot;5B3C1A2E053D763E1B002CC607C5A0FE&quot;</ETag>
				<Type>Normal</Type>
				<Size>344606</Size>
				<StorageClass>Standard</StorageClass>
				<Owner>
					<ID>00220120222</ID>
					<DisplayName>user-example</DisplayName>
				</Owner>
			</Contents>
		</ListBucketResult>`

	// unsupport content-type
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	result := &ListObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, "hello ", *result.Marker)
	assert.Equal(t, "hello world!", *result.Prefix)
	assert.Equal(t, "hello !world", *result.Delimiter)
	assert.Equal(t, "hello !world", *result.Delimiter)
	assert.Equal(t, "url", *result.EncodingType)
	assert.Equal(t, "oss-example", *result.Name)
	assert.Equal(t, false, result.IsTruncated)
	assert.Nil(t, result.NextMarker)
	assert.Len(t, result.Contents, 3)
	assert.Equal(t, "fun/movie/001.avi", *result.Contents[0].Key)
	assert.Equal(t, "\"5B3C1A2E053D763E1B002CC607C5A0FE\"", *result.Contents[0].ETag)
	assert.Equal(t, "fun/movie/007.avi", *result.Contents[1].Key)
	assert.Equal(t, "oss.jpg", *result.Contents[2].Key)
	assert.Len(t, result.CommonPrefixes, 0)
}

func TestUnmarshalOutput_encodetype1(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	body := `<?xml version="1.0" encoding="UTF-8"?>
		<ListBucketResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
			<Contents>
				<LastModified>2012-02-24T08:43:07.000Z</LastModified>
			</Contents>
		</ListBucketResult>`

	// unsupport content-type
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	result := &ListObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
}

func TestMarshalInput_PutBucket(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketRequest
	var input *OperationInput
	var err error

	request = &PutBucketRequest{}
	input = &OperationInput{
		OpName: "PutBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutBucketRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_PutBucket(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutBucketResult{}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 409,
		Status:     "BucketAlreadyExist",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &PutBucketResult{}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "BucketAlreadyExist")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteBucket(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteBucketRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketRequest{}
	input = &OperationInput{
		OpName: "PutBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteBucketRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteBucket(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	output = &OperationOutput{
		StatusCode: 209,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5C3D9778CC1C2AEDF85B****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &DeleteBucketResult{}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 209)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C3D9778CC1C2AEDF85B****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListObjects(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListObjectsRequest
	var input *OperationInput
	var err error

	request = &ListObjectsRequest{}
	input = &OperationInput{
		OpName: "ListObjects",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ListObjectsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "ListObjects",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &ListObjectsRequest{
		Bucket:       Ptr("oss-demo"),
		Delimiter:    Ptr("/"),
		Marker:       Ptr(""),
		MaxKeys:      int32(10),
		Prefix:       Ptr(""),
		EncodingType: Ptr("URL"),
	}
	input = &OperationInput{
		OpName: "ListObjects",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_ListObjects(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
<Prefix></Prefix>
<Marker></Marker>
<MaxKeys>100</MaxKeys>
<Delimiter></Delimiter>
<IsTruncated>false</IsTruncated>
<Contents>
      <Key>fun/movie/001.avi</Key>
      <LastModified>2012-02-24T08:43:07.000Z</LastModified>
      <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
      <Type>Normal</Type>
      <Size>344606</Size>
      <StorageClass>Standard</StorageClass>
      <Owner>
          <ID>0022012****</ID>
          <DisplayName>user-example</DisplayName>
      </Owner>
</Contents>
<Contents>
      <Key>fun/movie/007.avi</Key>
      <LastModified>2012-02-24T08:43:27.000Z</LastModified>
      <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
      <Type>Normal</Type>
      <Size>344606</Size>
      <StorageClass>Standard</StorageClass>
      <Owner>
          <ID>0022012****</ID>
          <DisplayName>user-example</DisplayName>
      </Owner>
</Contents>
<Contents>
      <Key>fun/test.jpg</Key>
      <LastModified>2012-02-24T08:42:32.000Z</LastModified>
      <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
      <Type>Normal</Type>
      <Size>344606</Size>
      <StorageClass>Standard</StorageClass>
      <Owner>
          <ID>0022012****</ID>
          <DisplayName>user-example</DisplayName>
      </Owner>
</Contents>
<Contents>
      <Key>oss.jpg</Key>
      <LastModified>2012-02-24T06:07:48.000Z</LastModified>
      <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
      <Type>Normal</Type>
      <Size>344606</Size>
      <StorageClass>Standard</StorageClass>
      <Owner>
          <ID>0022012****</ID>
          <DisplayName>user-example</DisplayName>
      </Owner>
</Contents>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &ListObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Empty(t, result.Prefix)
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Empty(t, result.Marker)
	assert.Empty(t, result.Delimiter)
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, len(result.Contents), 4)
	assert.Equal(t, *result.Contents[1].LastModified, time.Date(2012, time.February, 24, 8, 43, 27, 0, time.UTC))
	assert.Equal(t, *result.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
	assert.Equal(t, *result.Contents[3].Type, "Normal")
	assert.Equal(t, result.Contents[0].Size, int64(344606))
	assert.Equal(t, *result.Contents[1].StorageClass, "Standard")
	assert.Equal(t, *result.Contents[2].Owner.ID, "0022012****")
	assert.Equal(t, *result.Contents[3].Owner.DisplayName, "user-example")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
  <Prefix>fun</Prefix>
  <Marker>test1.txt</Marker>
  <MaxKeys>100</MaxKeys>
  <Delimiter>/</Delimiter>
  <IsTruncated>true</IsTruncated>
  <Contents>
        <Key>exampleobject1.txt</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
        <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>ColdArchive</StorageClass>
        <Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
  </Contents>
  <Contents>
        <Key>exampleobject2.txt</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
        <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>Standard</StorageClass>
        <RestoreInfo>ongoing-request="true"</RestoreInfo>
        <Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
  </Contents>
  <Contents>
        <Key>go-sdk-v1%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
        <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>Standard</StorageClass>
        <RestoreInfo>ongoing-request="false", expiry-date="Thu, 24 Sep 2020 12:40:33 GMT"</RestoreInfo>
        <Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
  </Contents>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Equal(t, *result.Prefix, "fun")
	assert.Equal(t, *result.Marker, "test1.txt")
	assert.Equal(t, *result.Delimiter, "/")
	assert.Equal(t, result.IsTruncated, true)
	assert.Equal(t, len(result.Contents), 3)
	assert.Equal(t, *result.Contents[0].Key, "exampleobject1.txt")
	assert.Equal(t, *result.Contents[1].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
	assert.Equal(t, *result.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
	assert.Equal(t, *result.Contents[0].Type, "Normal")
	assert.Equal(t, result.Contents[1].Size, int64(344606))
	assert.Equal(t, *result.Contents[2].StorageClass, "Standard")
	assert.Equal(t, *result.Contents[0].Owner.ID, "0022012****")
	assert.Equal(t, *result.Contents[0].Owner.DisplayName, "user-example")
	assert.Empty(t, result.Contents[0].RestoreInfo)
	assert.Equal(t, *result.Contents[1].RestoreInfo, "ongoing-request=\"true\"")
	assert.Equal(t, *result.Contents[2].RestoreInfo, "ongoing-request=\"false\", expiry-date=\"Thu, 24 Sep 2020 12:40:33 GMT\"")
	output = &OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
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
	resultErr := &ListObjectsResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListObjectsV2(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListObjectsRequestV2
	var input *OperationInput
	var err error

	request = &ListObjectsRequestV2{}
	input = &OperationInput{
		OpName: "ListObjectsV2",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ListObjectsRequestV2{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "ListObjects",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &ListObjectsRequestV2{
		Bucket:       Ptr("oss-demo"),
		Delimiter:    Ptr("/"),
		StartAfter:   Ptr(""),
		MaxKeys:      int32(10),
		Prefix:       Ptr(""),
		EncodingType: Ptr("URL"),
		FetchOwner:   true,
	}
	input = &OperationInput{
		OpName: "ListObjectsV2",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_ListObjectsV2(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix></Prefix>
    <MaxKeys>100</MaxKeys>
    <EncodingType>url</EncodingType>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>a</Key>
        <LastModified>2020-05-18T05:45:43.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>Standard</StorageClass>
    </Contents>
    <Contents>
        <Key>a/b</Key>
        <LastModified>2020-05-18T05:45:47.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>Standard</StorageClass>
    </Contents>
    <Contents>
        <Key>b</Key>
        <LastModified>2020-05-18T05:45:50.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <Contents>
        <Key>b/c</Key>
        <LastModified>2020-05-18T05:45:54.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <Contents>
        <Key>bc</Key>
        <LastModified>2020-05-18T05:45:59.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>Standard</StorageClass>
    </Contents>
    <Contents>
        <Key>c</Key>
        <LastModified>2020-05-18T05:45:57.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>Standard</StorageClass>
    </Contents>
    <KeyCount>6</KeyCount>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &ListObjectsResultV2{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Empty(t, result.Prefix)
	assert.Empty(t, result.Delimiter)
	assert.Empty(t, result.StartAfter)
	assert.Empty(t, result.ContinuationToken)
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Equal(t, *result.EncodingType, "url")
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, result.KeyCount, 6)
	assert.Equal(t, len(result.Contents), 6)
	assert.Equal(t, *result.Contents[0].Key, "a")
	assert.Equal(t, *result.Contents[1].LastModified, time.Date(2020, time.May, 18, 5, 45, 47, 0, time.UTC))
	assert.Equal(t, *result.Contents[2].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[3].Size, int64(25))
	assert.Equal(t, *result.Contents[0].StorageClass, "Standard")
	assert.Equal(t, *result.Contents[1].StorageClass, "Standard")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix>a</Prefix>
    <MaxKeys>100</MaxKeys>
    <EncodingType>url</EncodingType>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>a</Key>
        <LastModified>2020-05-18T05:45:43.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <Contents>
        <Key>a/b</Key>
        <LastModified>2020-05-18T05:45:47.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <KeyCount>2</KeyCount>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectsResultV2{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Equal(t, *result.Prefix, "a")
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Equal(t, len(result.Contents), 2)
	assert.Equal(t, *result.EncodingType, "url")
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, result.KeyCount, 2)
	assert.Equal(t, *result.Contents[0].Key, "a")
	assert.Equal(t, *result.Contents[0].LastModified, time.Date(2020, time.May, 18, 5, 45, 43, 0, time.UTC))
	assert.Equal(t, *result.Contents[0].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[0].Size, int64(25))
	assert.Equal(t, *result.Contents[0].StorageClass, "STANDARD")

	assert.Equal(t, *result.Contents[1].Key, "a/b")
	assert.Equal(t, *result.Contents[1].LastModified, time.Date(2020, time.May, 18, 5, 45, 47, 0, time.UTC))
	assert.Equal(t, *result.Contents[1].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[1].Size, int64(25))
	assert.Equal(t, *result.Contents[1].StorageClass, "STANDARD")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix>a/</Prefix>
    <MaxKeys>100</MaxKeys>
    <Delimiter>/</Delimiter>
    <EncodingType>url</EncodingType>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>a/b</Key>
        <LastModified>2020-05-18T05:45:47.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <CommonPrefixes>
        <Prefix>a/b/</Prefix>
    </CommonPrefixes>
    <KeyCount>2</KeyCount>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectsResultV2{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Equal(t, *result.Prefix, "a/")
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Equal(t, len(result.Contents), 1)
	assert.Equal(t, *result.EncodingType, "url")
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, result.KeyCount, 2)
	assert.Equal(t, *result.Contents[0].Key, "a/b")
	assert.Equal(t, *result.Contents[0].LastModified, time.Date(2020, time.May, 18, 5, 45, 47, 0, time.UTC))
	assert.Equal(t, *result.Contents[0].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[0].Size, int64(25))
	assert.Equal(t, *result.Contents[0].StorageClass, "STANDARD")

	assert.Equal(t, *result.CommonPrefixes[0].Prefix, "a/b/")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix></Prefix>
    <StartAfter>b</StartAfter>
    <MaxKeys>3</MaxKeys>
    <EncodingType>url</EncodingType>
    <IsTruncated>true</IsTruncated>
    <NextContinuationToken>CgJiYw--</NextContinuationToken>
    <Contents>
        <Key>b%2Fc</Key>
        <LastModified>2020-05-18T05:45:54.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1686240967192623</ID>
            <DisplayName>1686240967192623</DisplayName>
        </Owner>
    </Contents>
    <Contents>
        <Key>ba</Key>
        <LastModified>2020-05-18T11:17:58.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1686240967192623</ID>
            <DisplayName>1686240967192623</DisplayName>
        </Owner>
    </Contents>
    <Contents>
        <Key>bc</Key>
        <LastModified>2020-05-18T05:45:59.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1686240967192623</ID>
            <DisplayName>1686240967192623</DisplayName>
        </Owner>
    </Contents>
    <KeyCount>3</KeyCount>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectsResultV2{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Empty(t, result.Prefix)
	assert.Equal(t, *result.StartAfter, "b")
	assert.Equal(t, result.MaxKeys, int32(3))
	assert.Equal(t, len(result.Contents), 3)
	assert.Equal(t, *result.EncodingType, "url")
	assert.Equal(t, result.IsTruncated, true)
	assert.Equal(t, *result.NextContinuationToken, "CgJiYw--")
	assert.Equal(t, result.KeyCount, 3)
	assert.Equal(t, *result.Contents[0].Key, "b/c")
	assert.Equal(t, *result.Contents[0].LastModified, time.Date(2020, time.May, 18, 5, 45, 54, 0, time.UTC))
	assert.Equal(t, *result.Contents[0].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[0].Size, int64(25))
	assert.Equal(t, *result.Contents[0].StorageClass, "STANDARD")
	assert.Equal(t, *result.Contents[0].Owner.DisplayName, "1686240967192623")
	assert.Equal(t, *result.Contents[0].Owner.ID, "1686240967192623")

	body = `<?xml version="1.0" encoding="UTF-8"?>
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
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
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
	resultErr := &ListObjectsResultV2{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketInfoRequest
	var input *OperationInput
	var err error

	request = &GetBucketInfoRequest{}
	input = &OperationInput{
		OpName: "GetBucketInfo",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"bucketInfo": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketInfoRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketInfo",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"bucketInfo": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<BucketInfo>
  <Bucket>
    <AccessMonitor>Enabled</AccessMonitor>
    <CreationDate>2013-07-31T10:56:21.000Z</CreationDate>
    <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
    <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
    <Location>oss-cn-hangzhou</Location>
    <StorageClass>Standard</StorageClass>
    <TransferAcceleration>Disabled</TransferAcceleration>
    <CrossRegionReplication>Disabled</CrossRegionReplication>
    <Name>oss-example</Name>
    <ResourceGroupId>rg-aek27tc********</ResourceGroupId>
    <Owner>
      <DisplayName>username</DisplayName>
      <ID>27183473914****</ID>
    </Owner>
    <AccessControlList>
      <Grant>private</Grant>
    </AccessControlList>  
    <BucketPolicy>
      <LogBucket>examplebucket</LogBucket>
      <LogPrefix>log/</LogPrefix>
    </BucketPolicy>
  </Bucket>
</BucketInfo>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketInfo.Name, "oss-example")
	assert.Equal(t, *result.BucketInfo.AccessMonitor, "Enabled")
	assert.Equal(t, *result.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.Location, "oss-cn-hangzhou")
	assert.Equal(t, *result.BucketInfo.StorageClass, "Standard")
	assert.Equal(t, *result.BucketInfo.TransferAcceleration, "Disabled")
	assert.Equal(t, *result.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
	assert.Equal(t, *result.BucketInfo.CrossRegionReplication, "Disabled")
	assert.Equal(t, *result.BucketInfo.ResourceGroupId, "rg-aek27tc********")
	assert.Equal(t, *result.BucketInfo.Owner.ID, "27183473914****")
	assert.Equal(t, *result.BucketInfo.Owner.DisplayName, "username")
	assert.Equal(t, *result.BucketInfo.ACL, "private")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogPrefix, "log/")

	assert.Empty(t, result.BucketInfo.SseRule.KMSMasterKeyID)
	assert.Nil(t, result.BucketInfo.SseRule.SSEAlgorithm)
	assert.Nil(t, result.BucketInfo.SseRule.KMSDataEncryption)

	body = `<?xml version="1.0" encoding="UTF-8"?>
<BucketInfo>
  <Bucket>
    <AccessMonitor>Enabled</AccessMonitor>
    <CreationDate>2013-07-31T10:56:21.000Z</CreationDate>
    <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
    <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
    <Location>oss-cn-hangzhou</Location>
    <StorageClass>Standard</StorageClass>
    <TransferAcceleration>Disabled</TransferAcceleration>
    <CrossRegionReplication>Disabled</CrossRegionReplication>
    <Name>oss-example</Name>
    <ResourceGroupId>rg-aek27tc********</ResourceGroupId>
    <Owner>
      <DisplayName>username</DisplayName>
      <ID>27183473914****</ID>
    </Owner>
    <AccessControlList>
      <Grant>private</Grant>
    </AccessControlList>  
	<ServerSideEncryptionRule>
		<SSEAlgorithm>None</SSEAlgorithm>
	</ServerSideEncryptionRule>
    <BucketPolicy>
      <LogBucket>examplebucket</LogBucket>
      <LogPrefix>log/</LogPrefix>
    </BucketPolicy>
  </Bucket>
</BucketInfo>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	if result.BucketInfo.SseRule.KMSMasterKeyID != nil && *result.BucketInfo.SseRule.KMSMasterKeyID == "None" {
		*result.BucketInfo.SseRule.KMSMasterKeyID = ""
	}
	if result.BucketInfo.SseRule.SSEAlgorithm != nil && *result.BucketInfo.SseRule.SSEAlgorithm == "None" {
		*result.BucketInfo.SseRule.SSEAlgorithm = ""
	}
	if result.BucketInfo.SseRule.KMSDataEncryption != nil && *result.BucketInfo.SseRule.KMSDataEncryption == "None" {
		*result.BucketInfo.SseRule.KMSDataEncryption = ""
	}
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketInfo.Name, "oss-example")
	assert.Equal(t, *result.BucketInfo.AccessMonitor, "Enabled")
	assert.Equal(t, *result.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.Location, "oss-cn-hangzhou")
	assert.Equal(t, *result.BucketInfo.StorageClass, "Standard")
	assert.Equal(t, *result.BucketInfo.TransferAcceleration, "Disabled")
	assert.Equal(t, *result.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
	assert.Equal(t, *result.BucketInfo.CrossRegionReplication, "Disabled")
	assert.Equal(t, *result.BucketInfo.ResourceGroupId, "rg-aek27tc********")
	assert.Equal(t, *result.BucketInfo.Owner.ID, "27183473914****")
	assert.Equal(t, *result.BucketInfo.Owner.DisplayName, "username")
	assert.Equal(t, *result.BucketInfo.ACL, "private")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogPrefix, "log/")
	assert.Empty(t, result.BucketInfo.SseRule.KMSMasterKeyID)
	assert.Equal(t, *result.BucketInfo.SseRule.SSEAlgorithm, "")
	assert.Nil(t, result.BucketInfo.SseRule.KMSDataEncryption)

	body = `<?xml version="1.0" encoding="UTF-8"?>
<BucketInfo>
  <Bucket>
    <AccessMonitor>Enabled</AccessMonitor>
    <CreationDate>2013-07-31T10:56:21.000Z</CreationDate>
    <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
    <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
    <Location>oss-cn-hangzhou</Location>
    <StorageClass>Standard</StorageClass>
    <TransferAcceleration>Disabled</TransferAcceleration>
    <CrossRegionReplication>Disabled</CrossRegionReplication>
    <Name>oss-example</Name>
    <ResourceGroupId>rg-aek27tc********</ResourceGroupId>
    <Owner>
      <DisplayName>username</DisplayName>
      <ID>27183473914****</ID>
    </Owner>
    <AccessControlList>
      <Grant>private</Grant>
    </AccessControlList>  
	<ServerSideEncryptionRule>
		<SSEAlgorithm>KMS</SSEAlgorithm>
		<KMSMasterKeyID></KMSMasterKeyID>
		<KMSDataEncryption>SM4</KMSDataEncryption>
	</ServerSideEncryptionRule>
    <BucketPolicy>
      <LogBucket>examplebucket</LogBucket>
      <LogPrefix>log/</LogPrefix>
    </BucketPolicy>
  </Bucket>
</BucketInfo>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketInfo.Name, "oss-example")
	assert.Equal(t, *result.BucketInfo.AccessMonitor, "Enabled")
	assert.Equal(t, *result.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.Location, "oss-cn-hangzhou")
	assert.Equal(t, *result.BucketInfo.StorageClass, "Standard")
	assert.Equal(t, *result.BucketInfo.TransferAcceleration, "Disabled")
	assert.Equal(t, *result.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
	assert.Equal(t, *result.BucketInfo.CrossRegionReplication, "Disabled")
	assert.Equal(t, *result.BucketInfo.ResourceGroupId, "rg-aek27tc********")
	assert.Equal(t, *result.BucketInfo.Owner.ID, "27183473914****")
	assert.Equal(t, *result.BucketInfo.Owner.DisplayName, "username")
	assert.Equal(t, *result.BucketInfo.ACL, "private")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogPrefix, "log/")
	assert.Empty(t, *result.BucketInfo.SseRule.KMSMasterKeyID)
	assert.Equal(t, *result.BucketInfo.SseRule.SSEAlgorithm, "KMS")
	assert.Equal(t, *result.BucketInfo.SseRule.KMSDataEncryption, "SM4")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",

		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
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
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	resultErr := &GetBucketInfoResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketLocation(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketLocationRequest
	var input *OperationInput
	var err error

	request = &GetBucketLocationRequest{}
	input = &OperationInput{
		OpName: "GetBucketLocation",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"location": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketLocationRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketLocation",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"location": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketLocation(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<LocationConstraint>oss-cn-hangzhou</LocationConstraint>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketLocationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.LocationConstraint, "oss-cn-hangzhou")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
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
  <RequestId>534B371674E88A4D8906****</RequestId>
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
	resultErr := &GetBucketLocationResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")

}

func TestMarshalInput_GetBucketStat(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketStatRequest
	var input *OperationInput
	var err error

	request = &GetBucketStatRequest{}
	input = &OperationInput{
		OpName: "GetBucketStat",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"stat": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketStatRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketStat",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"stat": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketStat(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<BucketStat>
  <Storage>1600</Storage>
  <ObjectCount>230</ObjectCount>
  <MultipartUploadCount>40</MultipartUploadCount>
  <LiveChannelCount>4</LiveChannelCount>
  <LastModifiedTime>1643341269</LastModifiedTime>
  <StandardStorage>430</StandardStorage>
  <StandardObjectCount>66</StandardObjectCount>
  <InfrequentAccessStorage>2359296</InfrequentAccessStorage>
  <InfrequentAccessRealStorage>360</InfrequentAccessRealStorage>
  <InfrequentAccessObjectCount>54</InfrequentAccessObjectCount>
  <ArchiveStorage>2949120</ArchiveStorage>
  <ArchiveRealStorage>450</ArchiveRealStorage>
  <ArchiveObjectCount>74</ArchiveObjectCount>
  <ColdArchiveStorage>2359296</ColdArchiveStorage>
  <ColdArchiveRealStorage>360</ColdArchiveRealStorage>
  <ColdArchiveObjectCount>36</ColdArchiveObjectCount>
</BucketStat>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketStatResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, int64(1600), result.Storage)
	assert.Equal(t, int64(230), result.ObjectCount)
	assert.Equal(t, int64(40), result.MultipartUploadCount)
	assert.Equal(t, int64(4), result.LiveChannelCount)
	assert.Equal(t, int64(1643341269), result.LastModifiedTime)
	assert.Equal(t, int64(430), result.StandardStorage)
	assert.Equal(t, int64(66), result.StandardObjectCount)
	assert.Equal(t, int64(2359296), result.InfrequentAccessStorage)
	assert.Equal(t, int64(360), result.InfrequentAccessRealStorage)
	assert.Equal(t, int64(54), result.InfrequentAccessObjectCount)
	assert.Equal(t, int64(2949120), result.ArchiveStorage)
	assert.Equal(t, int64(450), result.ArchiveRealStorage)
	assert.Equal(t, int64(74), result.ArchiveObjectCount)
	assert.Equal(t, int64(2359296), result.ColdArchiveStorage)
	assert.Equal(t, int64(360), result.ColdArchiveRealStorage)
	assert.Equal(t, int64(36), result.ColdArchiveObjectCount)

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
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
  <RequestId>534B371674E88A4D8906****</RequestId>
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
	resultErr := &GetBucketStatResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")

}

func TestMarshalInput_PutBucketAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketAclRequest
	var input *OperationInput
	var err error

	request = &PutBucketAclRequest{}
	input = &OperationInput{
		OpName: "PutBucketAcl",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"acl": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutBucketAclRequest{
		Bucket: Ptr("oss-demo"),
		Acl:    BucketACLPrivate,
	}
	input = &OperationInput{
		OpName: "PutBucketAcl",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"acl": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_PutBucketAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &PutBucketAclResult{}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
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
	resultErr := &PutBucketAclResult{}
	err = c.unmarshalOutput(resultErr, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketAclRequest
	var input *OperationInput
	var err error

	request = &GetBucketAclRequest{}
	input = &OperationInput{
		OpName: "GetBucketAcl",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"acl": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketAclRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketAcl",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"acl": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012****</ID>
        <DisplayName>user_example</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read</Grant>
    </AccessControlList>
</AccessControlPolicy>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketAclResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	assert.Equal(t, *result.ACL, "public-read")
	assert.Equal(t, *result.Owner.ID, "0022012****")
	assert.Equal(t, *result.Owner.DisplayName, "user_example")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
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
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	resultErr := &PutBucketAclResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_PutBucketVersioning(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketVersioningRequest
	var input *OperationInput
	var err error

	request = &PutBucketVersioningRequest{}
	input = &OperationInput{
		OpName: "PutBucketVersioning",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"versioning": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutBucketVersioningRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucketVersioning",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"versioning": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, VersioningConfiguration.")

	request = &PutBucketVersioningRequest{
		Bucket: Ptr("oss-demo"),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	input = &OperationInput{
		OpName: "PutBucketVersioning",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"versioning": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<VersioningConfiguration><Status>Enabled</Status></VersioningConfiguration>")

	request = &PutBucketVersioningRequest{
		Bucket: Ptr("oss-demo"),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionSuspended,
		},
	}
	input = &OperationInput{
		OpName: "PutBucketVersioning",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"versioning": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<VersioningConfiguration><Status>Suspended</Status></VersioningConfiguration>")
}

func TestUnmarshalOutput_PutBucketVersioning(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &PutBucketAclResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
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
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
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
	resultErr := &PutBucketAclResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketVersioning(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketVersioningRequest
	var input *OperationInput
	var err error

	request = &GetBucketVersioningRequest{}
	input = &OperationInput{
		OpName: "GetBucketVersioning",
		Method: "GET",
		Parameters: map[string]string{
			"versioning": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketVersioningRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketVersioning",
		Method: "GET",
		Parameters: map[string]string{
			"versioning": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketVersioning(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
</VersioningConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketVersioningResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Nil(t, result.VersionStatus)

	body = `<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
<Status>Enabled</Status>
</VersioningConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketVersioningResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.VersionStatus, "Enabled")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
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
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
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
	resultErr := &GetBucketVersioningResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListObjectVersions(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListObjectVersionsRequest
	var input *OperationInput
	var err error

	request = &ListObjectVersionsRequest{}
	input = &OperationInput{
		OpName: "ListObjectVersions",
		Method: "GET",
		Parameters: map[string]string{
			"versions ":     "",
			"encoding-type": "url",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ListObjectVersionsRequest{
		Bucket: Ptr("oss-demo"),
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &ListObjectVersionsRequest{
		Bucket:          Ptr("oss-demo"),
		KeyMarker:       Ptr(""),
		VersionIdMarker: Ptr(""),
		Delimiter:       Ptr("/"),
		MaxKeys:         int32(100),
		Prefix:          Ptr("abc"),
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["key-marker"], "")
	assert.Equal(t, input.Parameters["version-id-marker"], "")
	assert.Equal(t, input.Parameters["delimiter"], "/")
	assert.Equal(t, input.Parameters["max-keys"], "100")
	assert.Equal(t, input.Parameters["prefix"], "abc")
}

func TestUnmarshalOutput_ListObjectVersions(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
    <Name>examplebucket-1250000000</Name>
    <Prefix/>
    <KeyMarker/>
    <VersionIdMarker/>
    <MaxKeys>1000</MaxKeys>
    <IsTruncated>false</IsTruncated>
    <Version>
        <Key>example-object-1.jpg</Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-05T12:03:10.000Z</LastModified>
        <ETag>5B3C1A2E053D763E1B669CC607C5A0FE1****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
    <Version>
        <Key>example-object-2.jpg</Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-09T12:03:09.000Z</LastModified>
        <ETag>5B3C1A2E053D763E1B002CC607C5A0FE1****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
    <Version>
        <Key>example-object-3.jpg</Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-10T12:03:08.000Z</LastModified>
        <ETag>4B3F1A2E053D763E1B002CC607C5AGTRF****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
</ListVersionsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &ListObjectVersionsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket-1250000000")
	assert.Equal(t, *result.Prefix, "")
	assert.Equal(t, *result.KeyMarker, "")
	assert.Equal(t, *result.VersionIdMarker, "")
	assert.Equal(t, result.MaxKeys, int32(1000))
	assert.False(t, result.IsTruncated)
	assert.Len(t, result.ObjectVersions, 3)
	assert.Equal(t, *result.ObjectVersions[0].Key, "example-object-1.jpg")
	assert.Empty(t, *result.ObjectVersions[1].VersionId)
	assert.True(t, result.ObjectVersions[2].IsLatest)
	assert.NotEmpty(t, *result.ObjectVersions[0].LastModified)
	assert.Equal(t, *result.ObjectVersions[1].ETag, "5B3C1A2E053D763E1B002CC607C5A0FE1****")
	assert.Equal(t, result.ObjectVersions[2].Size, int64(20))
	assert.Equal(t, *result.ObjectVersions[2].Owner.ID, "1250000000")
	assert.Equal(t, *result.ObjectVersions[2].Owner.DisplayName, "1250000000")
	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2Fgp-</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>5</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>false</IsTruncated>
  <Version>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <VersionId>CAEQHxiBgIDAj.jV3xgiIGFjMDI5ZTRmNGNiODQ0NjE4MDFhODM0Y2UxNTI3****</VersionId>
    <IsLatest>true</IsLatest>
    <LastModified>2023-11-22T05:15:05.000Z</LastModified>
    <ETag>"29B94424BC241D80B0AF488A4E4B86AF-1"</ETag>
    <Type>Multipart</Type>
    <Size>96316</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <Version>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <VersionId>CAEQHxiBgMDYseHV3xgiIDg2Mzk0Zjg3MjQ0MTRhM2FiMzgxOGY1NjdmN2Rk****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-11-22T05:11:25.000Z</LastModified>
    <ETag>"29B94424BC241D80B0AF488A4E4B86AF-1"</ETag>
    <Type>Multipart</Type>
    <Size>96316</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <Version>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <VersionId>CAEQHxiBgICCuNrV3xgiIDI2YzMyYTBhM2U1ZTQwNjI4OWQ4OTllZGJiNGIz****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-11-22T05:07:37.000Z</LastModified>
    <ETag>"29B94424BC241D80B0AF488A4E4B86AF-1"</ETag>
    <Type>Multipart</Type>
    <Size>96316</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
</ListVersionsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectVersionsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "demo-bucket")
	prefix, _ := url.QueryUnescape(*result.Prefix)
	assert.Equal(t, *result.Prefix, prefix)
	assert.Equal(t, *result.KeyMarker, "")
	assert.Equal(t, *result.VersionIdMarker, "")
	assert.Equal(t, result.MaxKeys, int32(5))
	assert.False(t, result.IsTruncated)
	assert.Len(t, result.ObjectVersions, 3)
	key, _ := url.QueryUnescape(*result.ObjectVersions[0].Key)
	assert.Equal(t, *result.ObjectVersions[0].Key, key)
	assert.Equal(t, *result.ObjectVersions[1].VersionId, "CAEQHxiBgMDYseHV3xgiIDg2Mzk0Zjg3MjQ0MTRhM2FiMzgxOGY1NjdmN2Rk****")
	assert.False(t, result.ObjectVersions[2].IsLatest)
	assert.NotEmpty(t, *result.ObjectVersions[0].LastModified)
	assert.Equal(t, *result.ObjectVersions[1].ETag, "\"29B94424BC241D80B0AF488A4E4B86AF-1\"")
	assert.Equal(t, result.ObjectVersions[2].Size, int64(96316))
	assert.Equal(t, *result.ObjectVersions[2].Owner.ID, "150692521021****")
	assert.Equal(t, *result.ObjectVersions[2].Owner.DisplayName, "150692521021****")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2F</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>20</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>true</IsTruncated>
  <NextKeyMarker>demo%2FREADME-CN.md</NextKeyMarker>
  <NextVersionIdMarker>CAEQEhiBgICDzK6NnBgiIGRlZWJhYmNlMGUxZDQ4YTZhNTU2MzM4Mzk5NDBl****</NextVersionIdMarker>
  <DeleteMarker>
    <Key>demo%2F</Key>
    <VersionId>CAEQFxiBgIDh3b_tuRgiIGRjMjExMjVmMzcwMTQ2Njc4NjhhNTA0MzEzMDkx****</VersionId>
    <IsLatest>true</IsLatest>
    <LastModified>2023-04-01T05:52:31.000Z</LastModified>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </DeleteMarker>
  <Version>
    <Key>demo%2F</Key>
    <VersionId>CAEQFxiBgICI173TtRgiIDFlMmYyMzFjNmJmMDQ0NTBiNmIyYThkZjA1YjA5****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-03-06T03:02:28.000Z</LastModified>
    <ETag>"D41D8CD98F00B204E9800998ECF8427E"</ETag>
    <Type>Normal</Type>
    <Size>0</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <DeleteMarker>
    <Key>demo%2F</Key>
    <VersionId>CAEQFxiBgICHsJuXtRgiIDMzNzUxNWIwYzEwODRlYTg5MTgxMDhmYTUzNDQz****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-03-03T04:49:26.000Z</LastModified>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </DeleteMarker>
  <Version>
    <Key>demo%2F</Key>
    <VersionId>CAEQFxiBgIC.oPmDtRgiIDNhNGZjMDQxMTYwYTRkYzE4ZDk4YTQ2NmYxYjA1****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-03-02T06:22:36.000Z</LastModified>
    <ETag>"D41D8CD98F00B204E9800998ECF8****"</ETag>
    <Type>Normal</Type>
    <Size>0</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <DeleteMarker>
    <Key>demo%2F</Key>
    <VersionId>CAEQFxiBgMCH__iDtRgiIDk4ZDFjZGY3NTk5ZjQ2NjViMzhjZjA2ODUwNjU0****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-03-02T06:22:27.000Z</LastModified>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </DeleteMarker>
  <Version>
    <Key>demo%2F</Key>
    <VersionId>CAEQFxiBgMD3p_iDtRgiIDljYjdlNzA3ZjE3ZTQ4NzI4ODE1ZWQ1ZWFlYjZl****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-03-02T06:22:05.000Z</LastModified>
    <ETag>"D41D8CD98F00B204E9800998ECF8427E"</ETag>
    <Type>Normal</Type>
    <Size>0</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <DeleteMarker>
    <Key>demo%2F</Key>
    <VersionId>CAEQFxiBgMC_6feDtRgiIGQ4YTIyOTNjZDY4ZjQ1NGY5NGE5YTNlOTBlODlm****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-03-02T06:21:49.000Z</LastModified>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </DeleteMarker>
  <Version>
    <Key>demo%2F</Key>
    <VersionId>CAEQFxiBgIDPjPSDtRgiIDcwMWU0ZDg2Y2NlNzRhZTM5NDM5ZmMxYjMwZGUw****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-03-02T06:19:47.000Z</LastModified>
    <ETag>"D41D8CD98F00B204E9800998ECF8****"</ETag>
    <Type>Normal</Type>
    <Size>0</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <DeleteMarker>
    <Key>demo%2F.gitignore</Key>
    <VersionId>CAEQFBiBgIDd.86GohgiIDMyMmVlZGNhOTI4OTQ3M2M5MGJiYTVmNTBjYjhl****</VersionId>
    <IsLatest>true</IsLatest>
    <LastModified>2022-11-04T08:00:06.000Z</LastModified>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </DeleteMarker>
  <Version>
    <Key>demo%2F.gitignore</Key>
    <VersionId>CAEQEhiBgIDkyq6NnBgiIDMyMGNhN2JjODllMjQwNWFhZThmZGRkZDRmYzlh****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2022-09-28T09:04:39.000Z</LastModified>
    <ETag>"C173E921A7464E5147B26B5F3DF9****"</ETag>
    <Type>Normal</Type>
    <Size>166</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <DeleteMarker>
    <Key>demo%2F.travis.yml</Key>
    <VersionId>CAEQFBiBgMDv.86GohgiIDc5ZmM5MTkxNjJmZDQ1OWU4Njk4MGI5ODI4M2Yw****</VersionId>
    <IsLatest>true</IsLatest>
    <LastModified>2022-11-04T08:00:06.000Z</LastModified>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </DeleteMarker>
  <Version>
    <Key>demo%2F.travis.yml</Key>
    <VersionId>CAEQEhiBgIDOy66NnBgiIDE5MmVkNzRmOGUxNzRmM2I4NTEyMzBhOGZhMWQw****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2022-09-28T09:04:39.000Z</LastModified>
    <ETag>"1D66AB946CCD6C2E4D7FD65D8D80****"</ETag>
    <Type>Normal</Type>
    <Size>4046</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <DeleteMarker>
    <Key>demo%2FCHANGELOG.md</Key>
    <VersionId>CAEQFBiBgIDy.86GohgiIGE0NTU2ZTFlZWQ4ZTQwZmZiMjc4ZmJhZmQ2YzZj****</VersionId>
    <IsLatest>true</IsLatest>
    <LastModified>2022-11-04T08:00:06.000Z</LastModified>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </DeleteMarker>
  <Version>
    <Key>demo%2FCHANGELOG.md</Key>
    <VersionId>CAEQEhiBgID3y66NnBgiIDk2YmIwYmMxZWYzOTQ4Y2JhZjViMzMzZjg5ZjFm****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2022-09-28T09:04:39.000Z</LastModified>
    <ETag>"1CB587ACD1BB5A0442CAD8A972E0****"</ETag>
    <Type>Normal</Type>
    <Size>6745</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <DeleteMarker>
    <Key>demo%2FLICENSE</Key>
    <VersionId>CAEQFBiBgMD0.86GohgiIGZmMmFlM2UwNjdlMzRiMGFhYjk4MjM1ZGUyZDY0****</VersionId>
    <IsLatest>true</IsLatest>
    <LastModified>2022-11-04T08:00:06.000Z</LastModified>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </DeleteMarker>
  <Version>
    <Key>demo%2FLICENSE</Key>
    <VersionId>CAEQEhiBgICIzK6NnBgiIDMxYjM3OTdmN2E0ODRjZjhhOWVhYTE5MTg3NmQw****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2022-09-28T09:04:39.000Z</LastModified>
    <ETag>"877D6894CBE5711A315681C24ED0****"</ETag>
    <Type>Normal</Type>
    <Size>1094</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <DeleteMarker>
    <Key>demo%2FREADME-CN.md</Key>
    <VersionId>CAEQFBiCgID3.86GohgiIDc4ZTE0NTNhZTc5MDQxYzBhYTU5MjY1ZDFjNGJm****</VersionId>
    <IsLatest>true</IsLatest>
    <LastModified>2022-11-04T08:00:06.000Z</LastModified>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </DeleteMarker>
  <Version>
    <Key>demo%2FREADME-CN.md</Key>
    <VersionId>CAEQEhiBgICDzK6NnBgiIGRlZWJhYmNlMGUxZDQ4YTZhNTU2MzM4Mzk5NDBl****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2022-09-28T09:04:39.000Z</LastModified>
    <ETag>"E317049B40462DE37C422CE4FC1B****"</ETag>
    <Type>Normal</Type>
    <Size>2943</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <CommonPrefixes>
    <Prefix>demo%2F.git%2F</Prefix>
  </CommonPrefixes>
  <CommonPrefixes>
    <Prefix>demo%2F.idea%2F</Prefix>
  </CommonPrefixes>
</ListVersionsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectVersionsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "demo-bucket")
	prefix, _ = url.QueryUnescape(*result.Prefix)
	assert.Equal(t, *result.Prefix, prefix)
	assert.Equal(t, *result.KeyMarker, "")
	assert.Equal(t, *result.VersionIdMarker, "")
	assert.Equal(t, result.MaxKeys, int32(20))
	assert.True(t, result.IsTruncated)
	assert.Len(t, result.ObjectVersions, 9)
	assert.Len(t, result.ObjectDeleteMarkers, 9)
	key, _ = url.QueryUnescape(*result.ObjectVersions[0].Key)
	assert.Equal(t, *result.ObjectVersions[0].Key, key)
	assert.False(t, result.ObjectVersions[2].IsLatest)
	assert.NotEmpty(t, *result.ObjectVersions[0].LastModified)
	assert.Equal(t, *result.ObjectVersions[1].ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, result.ObjectVersions[2].Size, int64(0))
	assert.Equal(t, *result.ObjectVersions[2].Owner.ID, "150692521021****")
	assert.Equal(t, *result.ObjectVersions[2].Owner.DisplayName, "150692521021****")
	assert.Len(t, result.CommonPrefixes, 2)
	compPrefix1, _ := url.QueryUnescape(*result.CommonPrefixes[0].Prefix)
	compPrefix2, _ := url.QueryUnescape(*result.CommonPrefixes[1].Prefix)
	assert.Equal(t, *result.CommonPrefixes[0].Prefix, compPrefix1)
	assert.Equal(t, *result.CommonPrefixes[1].Prefix, compPrefix2)
	key, _ = url.QueryUnescape(*result.ObjectDeleteMarkers[0].Key)
	assert.Equal(t, *result.ObjectDeleteMarkers[0].Key, key)
	assert.Equal(t, *result.ObjectDeleteMarkers[0].VersionId, "CAEQFxiBgIDh3b_tuRgiIGRjMjExMjVmMzcwMTQ2Njc4NjhhNTA0MzEzMDkx****")
	assert.True(t, result.ObjectDeleteMarkers[0].IsLatest)
	assert.NotEmpty(t, result.ObjectDeleteMarkers[0].LastModified)
	assert.Equal(t, *result.ObjectDeleteMarkers[0].Owner.ID, "150692521021****")
	assert.Equal(t, *result.ObjectDeleteMarkers[0].Owner.DisplayName, "150692521021****")
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
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
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
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
	resultErr := &ListObjectVersionsResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}
