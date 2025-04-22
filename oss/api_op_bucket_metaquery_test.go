package oss

import (
	"bytes"
	"html"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_OpenMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *OpenMetaQueryRequest
	var input *OperationInput
	var err error

	request = &OpenMetaQueryRequest{}
	input = &OperationInput{
		OpName: "OpenMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "add",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &OpenMetaQueryRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "OpenMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "add",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &OpenMetaQueryRequest{
		Bucket: Ptr("oss-demo"),
		Mode:   Ptr("basic"),
	}
	input = &OperationInput{
		OpName: "OpenMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "add",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["mode"], "basic")

	request = &OpenMetaQueryRequest{
		Bucket: Ptr("oss-demo"),
		Mode:   Ptr("semantic"),
	}
	input = &OperationInput{
		OpName: "OpenMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "add",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["mode"], "semantic")
}

func TestUnmarshalOutput_OpenMetaQuery(t *testing.T) {
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
	result := &OpenMetaQueryResult{}
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
	result = &OpenMetaQueryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "MetaQueryAlreadyExist",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &OpenMetaQueryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "MetaQueryAlreadyExist")
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
	result = &OpenMetaQueryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetMetaQueryStatus(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetMetaQueryStatusRequest
	var input *OperationInput
	var err error

	request = &GetMetaQueryStatusRequest{}
	input = &OperationInput{
		OpName: "GetMetaQueryStatus",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetMetaQueryStatusRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetMetaQueryStatus",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetMetaQueryStatus(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<MetaQueryStatus>
  <State>Running</State>
  <Phase>FullScanning</Phase>
  <CreateTime>2021-08-02T10:49:17.289372919+08:00</CreateTime>
  <UpdateTime>2021-08-02T10:49:17.289372919+08:00</UpdateTime>
</MetaQueryStatus>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetMetaQueryStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.MetaQueryStatus.State, "Running")
	assert.Equal(t, *result.MetaQueryStatus.Phase, "FullScanning")
	assert.Equal(t, *result.MetaQueryStatus.CreateTime, "2021-08-02T10:49:17.289372919+08:00")
	assert.Equal(t, *result.MetaQueryStatus.UpdateTime, "2021-08-02T10:49:17.289372919+08:00")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",

		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetMetaQueryStatusResult{}
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
	result = &GetMetaQueryStatusResult{}
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
	result = &GetMetaQueryStatusResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DoMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DoMetaQueryRequest
	var input *OperationInput
	var err error

	request = &DoMetaQueryRequest{}
	input = &OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "query",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DoMetaQueryRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "query",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, MetaQuery.")

	request = &DoMetaQueryRequest{
		Bucket: Ptr("bucket"),
		MetaQuery: &MetaQuery{
			NextToken:  Ptr("MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****"),
			MaxResults: Ptr(int64(5)),
			Query:      Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
			Sort:       Ptr("Size"),
			Order:      Ptr(MetaQueryOrderAsc),
			Aggregations: &MetaQueryAggregations{
				[]MetaQueryAggregation{
					{
						Field:     Ptr("Size"),
						Operation: Ptr("sum"),
					},
					{
						Field:     Ptr("Size"),
						Operation: Ptr("max"),
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "query",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, html.UnescapeString(string(body)), "<MetaQuery><MaxResults>5</MaxResults><Query>{\"Field\": \"Size\",\"Value\": \"1048576\",\"Operation\": \"gt\"}</Query><Sort>Size</Sort><Order>asc</Order><Aggregations><Aggregation><Field>Size</Field><Operation>sum</Operation></Aggregation><Aggregation><Field>Size</Field><Operation>max</Operation></Aggregation></Aggregations><NextToken>MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****</NextToken></MetaQuery>")

	request = &DoMetaQueryRequest{
		Bucket: Ptr("bucket"),
		Mode:   Ptr("semantic"),
		MetaQuery: &MetaQuery{
			NextToken:   Ptr("MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****"),
			MaxResults:  Ptr(int64(99)),
			Query:       Ptr(`Overlook the snow-covered forest`),
			MediaType:   Ptr("image"),
			SimpleQuery: Ptr(`{"Operation":"gt", "Field": "Size", "Value": "30"}`),
		},
	}
	input = &OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "query",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["mode"], "semantic")
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, html.UnescapeString(string(body)), `<MetaQuery><MaxResults>99</MaxResults><Query>Overlook the snow-covered forest</Query><NextToken>MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****</NextToken><MediaTypes><MediaType>image</MediaType></MediaTypes><SimpleQuery>{"Operation":"gt", "Field": "Size", "Value": "30"}</SimpleQuery></MetaQuery>`)
}

func TestUnmarshalOutput_DoMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
  <NextToken>MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****</NextToken>
  <Files>
    <File>
      <Filename>exampleobject.txt</Filename>
      <Size>120</Size>
      <FileModifiedTime>2021-06-29T15:04:05.000000000Z07:00</FileModifiedTime>
      <OSSObjectType>Normal</OSSObjectType>
      <OSSStorageClass>Standard</OSSStorageClass>
      <ObjectACL>default</ObjectACL>
      <ETag>"fba9dede5f27731c9771645a3986****"</ETag>
      <OSSCRC64>4858A48BD1466884</OSSCRC64>
      <OSSTaggingCount>2</OSSTaggingCount>
      <OSSTagging>
        <Tagging>
          <Key>owner</Key>
          <Value>John</Value>
        </Tagging>
        <Tagging>
          <Key>type</Key>
          <Value>document</Value>
        </Tagging>
      </OSSTagging>
      <OSSUserMeta>
        <UserMeta>
          <Key>x-oss-meta-location</Key>
          <Value>hangzhou</Value>
        </UserMeta>
      </OSSUserMeta>
    </File>
  </Files>
  <Aggregations>
    <Aggregation>
      <Field>Size</Field>
      <Operation>sum</Operation>
      <Value>4859250309</Value>
    </Aggregation>
    <Aggregation>
      <Field>Size</Field>
      <Operation>max</Operation>
      <Value>2235483240</Value>
    </Aggregation>
  </Aggregations>
</MetaQuery>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &DoMetaQueryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.NextToken, "MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****")
	assert.Equal(t, len(result.Files), 1)
	assert.Equal(t, *result.Files[0].Filename, "exampleobject.txt")
	assert.Equal(t, *result.Files[0].Size, int64(120))
	assert.Equal(t, *result.Files[0].FileModifiedTime, "2021-06-29T15:04:05.000000000Z07:00")
	assert.Equal(t, *result.Files[0].OSSObjectType, "Normal")
	assert.Equal(t, *result.Files[0].OSSStorageClass, "Standard")
	assert.Equal(t, *result.Files[0].ObjectACL, "default")
	assert.Equal(t, *result.Files[0].ETag, "\"fba9dede5f27731c9771645a3986****\"")
	assert.Equal(t, *result.Files[0].OSSTaggingCount, int64(2))
	assert.Equal(t, *result.Files[0].OSSTagging[0].Key, "owner")
	assert.Equal(t, *result.Files[0].OSSTagging[0].Value, "John")
	assert.Equal(t, *result.Files[0].OSSTagging[1].Key, "type")
	assert.Equal(t, *result.Files[0].OSSTagging[1].Value, "document")
	assert.Equal(t, len(result.Files[0].OSSUserMeta), 1)
	assert.Equal(t, *result.Files[0].OSSUserMeta[0].Key, "x-oss-meta-location")
	assert.Equal(t, *result.Files[0].OSSUserMeta[0].Value, "hangzhou")
	assert.Equal(t, len(result.Aggregations), 2)
	assert.Equal(t, *result.Aggregations[0].Field, "Size")
	assert.Equal(t, *result.Aggregations[0].Operation, "sum")
	assert.Equal(t, *result.Aggregations[0].Value, float64(4859250309))
	assert.Equal(t, *result.Aggregations[1].Field, "Size")
	assert.Equal(t, *result.Aggregations[1].Operation, "max")
	assert.Equal(t, *result.Aggregations[1].Value, float64(2235483240))

	body = `<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
  <NextToken></NextToken>
  <Aggregations>
    <Aggregation>
      <Field>Size</Field>
      <Operation>sum</Operation>
      <Value>30930054</Value>
    </Aggregation>
    <Aggregation>
      <Field>Size</Field>
      <Operation>group</Operation>
      <Groups>
        <Group>
          <Value>1536000</Value>
          <Count>1</Count>
        </Group>
        <Group>
          <Value>5472362</Value>
          <Count>1</Count>
        </Group>
        <Group>
          <Value>10354204</Value>
          <Count>1</Count>
        </Group>
        <Group>
          <Value>1890304</Value>
          <Count>3</Count>
        </Group>
        <Group>
          <Value>2632192</Value>
          <Count>3</Count>
        </Group>
      </Groups>
    </Aggregation>
  </Aggregations>
</MetaQuery>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DoMetaQueryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, len(result.Aggregations), 2)
	assert.Equal(t, *result.Aggregations[0].Field, "Size")
	assert.Equal(t, *result.Aggregations[0].Operation, "sum")
	assert.Equal(t, *result.Aggregations[0].Value, float64(30930054))
	assert.Equal(t, *result.Aggregations[1].Field, "Size")
	assert.Equal(t, *result.Aggregations[1].Operation, "group")
	assert.Equal(t, len(result.Aggregations[1].Groups.Groups), 5)
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[0].Value, "1536000")
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[0].Count, int64(1))
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[1].Value, "5472362")
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[1].Count, int64(1))
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[2].Value, "10354204")
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[2].Count, int64(1))
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[3].Value, "1890304")
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[3].Count, int64(3))
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[4].Value, "2632192")
	assert.Equal(t, *result.Aggregations[1].Groups.Groups[4].Count, int64(3))

	body = `<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
  <Files>
    <File>
      <URI>oss://bucket/sample-object.jpg</URI>
      <Filename>sample-object.jpg</Filename>
      <Size>1000</Size>
      <ObjectACL>default</ObjectACL>
      <FileModifiedTime>2021-06-29T14:50:14.011643661+08:00</FileModifiedTime>
      <ServerSideEncryption>AES256</ServerSideEncryption>
      <ServerSideEncryptionCustomerAlgorithm>SM4</ServerSideEncryptionCustomerAlgorithm>
      <ETag>"1D9C280A7C4F67F7EF873E28449****"</ETag>
      <OSSCRC64>559890638950338001</OSSCRC64>
      <ProduceTime>2021-06-29T14:50:15.011643661+08:00</ProduceTime>
      <ContentType>image/jpeg</ContentType>
      <MediaType>image</MediaType>
      <LatLong>30.134390,120.074997</LatLong>
      <Title>test</Title>
      <OSSExpiration>2024-12-01T12:00:00.000Z</OSSExpiration>
      <AccessControlAllowOrigin>https://aliyundoc.com</AccessControlAllowOrigin>
      <AccessControlRequestMethod>PUT</AccessControlRequestMethod>
      <ServerSideDataEncryption>SM4</ServerSideDataEncryption>
      <ServerSideEncryptionKeyId>9468da86-3509-4f8d-a61e-6eab1eac****</ServerSideEncryptionKeyId>
      <CacheControl>no-cache</CacheControl>
      <ContentDisposition>attachment; filename=test.jpg</ContentDisposition>
      <ContentEncoding>UTF-8</ContentEncoding>
      <ContentLanguage>zh-CN</ContentLanguage>
      <ImageHeight>500</ImageHeight>
      <ImageWidth>270</ImageWidth>
      <VideoWidth>1080</VideoWidth>
      <VideoHeight>1920</VideoHeight>
      <VideoStreams>
        <VideoStream>
          <CodecName>h264</CodecName>
          <Language>en</Language>
          <Bitrate>5407765</Bitrate>
          <FrameRate>25/1</FrameRate>
          <StartTime>0</StartTime>
          <Duration>22.88</Duration>
          <FrameCount>572</FrameCount>
          <BitDepth>8</BitDepth>
          <PixelFormat>yuv420p</PixelFormat>
          <ColorSpace>bt709</ColorSpace>
          <Height>720</Height>
          <Width>1280</Width>
        </VideoStream>
        <VideoStream>
          <CodecName>h264</CodecName>
          <Language>en</Language>
          <Bitrate>5407765</Bitrate>
          <FrameRate>25/1</FrameRate>
          <StartTime>0</StartTime>
          <Duration>22.88</Duration>
          <FrameCount>572</FrameCount>
          <BitDepth>8</BitDepth>
          <PixelFormat>yuv420p</PixelFormat>
          <ColorSpace>bt709</ColorSpace>
          <Height>720</Height>
          <Width>1280</Width>
        </VideoStream>
      </VideoStreams>
      <AudioStreams>
        <AudioStream>
          <CodecName>aac</CodecName>
          <Bitrate>1048576</Bitrate>
          <SampleRate>48000</SampleRate>
          <StartTime>0.0235</StartTime>
          <Duration>3.690667</Duration>
          <Channels>2</Channels>
          <Language>en</Language>
        </AudioStream>
      </AudioStreams>
      <Subtitles>
        <Subtitle>
          <CodecName>mov_text</CodecName>
          <Language>en</Language>
          <StartTime>0</StartTime>
          <Duration>71.378</Duration>
        </Subtitle>
        <Subtitle>
          <CodecName>mov_text</CodecName>
          <Language>en</Language>
          <StartTime>72</StartTime>
          <Duration>71.378</Duration>
        </Subtitle>
      </Subtitles>
      <Bitrate>5407765</Bitrate>
      <Artist>Jane</Artist>
      <AlbumArtist>Jenny</AlbumArtist>
      <Composer>Jane</Composer>
      <Performer>Jane</Performer>
      <Album>FirstAlbum</Album>
      <Duration>71.378</Duration>
      <Addresses>
        <Address>
          <AddressLine>中国浙江省杭州市余杭区文一西路969号</AddressLine>
          <City>杭州市</City>
          <Country>中国</Country>
          <District>余杭区</District>
          <Language>zh-Hans</Language>
          <Province>浙江省</Province>
          <Township>文一西路</Township>
        </Address>
        <Address>
          <AddressLine>中国浙江省杭州市余杭区文一西路970号</AddressLine>
          <City>杭州市</City>
          <Country>中国</Country>
          <District>余杭区</District>
          <Language>zh-Hans</Language>
          <Province>浙江省</Province>
          <Township>文一西路</Township>
        </Address>
      </Addresses>
      <OSSObjectType>Normal</OSSObjectType>
      <OSSStorageClass>Standard</OSSStorageClass>
      <OSSTaggingCount>2</OSSTaggingCount>
      <OSSTagging>
        <Tagging>
          <Key>key</Key>
          <Value>val</Value>
        </Tagging>
        <Tagging>
          <Key>key2</Key>
          <Value>val2</Value>
        </Tagging>
      </OSSTagging>
      <OSSUserMeta>
        <UserMeta>
          <Key>key</Key>
          <Value>val</Value>
        </UserMeta>
      </OSSUserMeta>
    </File>
  </Files>
</MetaQuery>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DoMetaQueryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, len(result.Files), 1)
	assert.Equal(t, *result.Files[0].URI, "oss://bucket/sample-object.jpg")
	assert.Equal(t, *result.Files[0].Filename, "sample-object.jpg")
	assert.Equal(t, *result.Files[0].Size, int64(1000))
	assert.Equal(t, *result.Files[0].FileModifiedTime, "2021-06-29T14:50:14.011643661+08:00")
	assert.Equal(t, *result.Files[0].ServerSideEncryption, "AES256")
	assert.Equal(t, *result.Files[0].ServerSideEncryptionCustomerAlgorithm, "SM4")
	assert.Equal(t, *result.Files[0].ETag, "\"1D9C280A7C4F67F7EF873E28449****\"")
	assert.Equal(t, *result.Files[0].OSSCRC64, "559890638950338001")
	assert.Equal(t, *result.Files[0].ProduceTime, "2021-06-29T14:50:15.011643661+08:00")
	assert.Equal(t, *result.Files[0].ContentType, "image/jpeg")
	assert.Equal(t, *result.Files[0].MediaType, "image")
	assert.Equal(t, *result.Files[0].LatLong, "30.134390,120.074997")
	assert.Equal(t, *result.Files[0].Title, "test")
	assert.Equal(t, *result.Files[0].OSSExpiration, "2024-12-01T12:00:00.000Z")
	assert.Equal(t, *result.Files[0].AccessControlAllowOrigin, "https://aliyundoc.com")
	assert.Equal(t, *result.Files[0].AccessControlRequestMethod, "PUT")
	assert.Equal(t, *result.Files[0].ServerSideDataEncryption, "SM4")
	assert.Equal(t, *result.Files[0].ServerSideEncryptionKeyId, "9468da86-3509-4f8d-a61e-6eab1eac****")
	assert.Equal(t, *result.Files[0].CacheControl, "no-cache")
	assert.Equal(t, *result.Files[0].ContentDisposition, "attachment; filename=test.jpg")
	assert.Equal(t, *result.Files[0].ContentEncoding, "UTF-8")
	assert.Equal(t, *result.Files[0].ContentLanguage, "zh-CN")
	assert.Equal(t, *result.Files[0].ImageHeight, int64(500))
	assert.Equal(t, *result.Files[0].ImageWidth, int64(270))
	assert.Equal(t, *result.Files[0].VideoWidth, int64(1080))
	assert.Equal(t, *result.Files[0].VideoHeight, int64(1920))
	assert.Equal(t, len(result.Files[0].VideoStreams), 2)
	assert.Equal(t, *result.Files[0].VideoStreams[0].CodecName, "h264")
	assert.Equal(t, *result.Files[0].VideoStreams[0].Language, "en")
	assert.Equal(t, *result.Files[0].VideoStreams[0].Bitrate, int64(5407765))
	assert.Equal(t, *result.Files[0].VideoStreams[0].FrameRate, "25/1")
	assert.Equal(t, *result.Files[0].VideoStreams[0].StartTime, float64(0))
	assert.Equal(t, *result.Files[0].VideoStreams[0].Duration, float64(22.88))
	assert.Equal(t, *result.Files[0].VideoStreams[0].FrameCount, int64(572))
	assert.Equal(t, *result.Files[0].VideoStreams[0].BitDepth, int64(8))
	assert.Equal(t, *result.Files[0].VideoStreams[0].PixelFormat, "yuv420p")
	assert.Equal(t, *result.Files[0].VideoStreams[0].ColorSpace, "bt709")
	assert.Equal(t, *result.Files[0].VideoStreams[0].Height, int64(720))
	assert.Equal(t, *result.Files[0].VideoStreams[0].Width, int64(1280))

	assert.Equal(t, *result.Files[0].VideoStreams[1].CodecName, "h264")
	assert.Equal(t, *result.Files[0].VideoStreams[1].Language, "en")
	assert.Equal(t, *result.Files[0].VideoStreams[1].Bitrate, int64(5407765))
	assert.Equal(t, *result.Files[0].VideoStreams[1].FrameRate, "25/1")
	assert.Equal(t, *result.Files[0].VideoStreams[1].StartTime, float64(0))
	assert.Equal(t, *result.Files[0].VideoStreams[1].Duration, float64(22.88))
	assert.Equal(t, *result.Files[0].VideoStreams[1].FrameCount, int64(572))
	assert.Equal(t, *result.Files[0].VideoStreams[1].BitDepth, int64(8))
	assert.Equal(t, *result.Files[0].VideoStreams[1].PixelFormat, "yuv420p")
	assert.Equal(t, *result.Files[0].VideoStreams[1].ColorSpace, "bt709")
	assert.Equal(t, *result.Files[0].VideoStreams[1].Height, int64(720))
	assert.Equal(t, *result.Files[0].VideoStreams[1].Width, int64(1280))

	assert.Equal(t, len(result.Files[0].AudioStreams), 1)
	assert.Equal(t, *result.Files[0].AudioStreams[0].CodecName, "aac")
	assert.Equal(t, *result.Files[0].AudioStreams[0].Bitrate, int64(1048576))
	assert.Equal(t, *result.Files[0].AudioStreams[0].SampleRate, int64(48000))
	assert.Equal(t, *result.Files[0].AudioStreams[0].StartTime, float64(0.0235))
	assert.Equal(t, *result.Files[0].AudioStreams[0].Duration, float64(3.690667))
	assert.Equal(t, *result.Files[0].AudioStreams[0].Channels, int64(2))
	assert.Equal(t, *result.Files[0].AudioStreams[0].Language, "en")

	assert.Equal(t, len(result.Files[0].Subtitles), 2)
	assert.Equal(t, *result.Files[0].Subtitles[0].CodecName, "mov_text")
	assert.Equal(t, *result.Files[0].Subtitles[0].Language, "en")
	assert.Equal(t, *result.Files[0].Subtitles[0].StartTime, float64(0))
	assert.Equal(t, *result.Files[0].Subtitles[0].Duration, float64(71.378))
	assert.Equal(t, *result.Files[0].Subtitles[1].CodecName, "mov_text")
	assert.Equal(t, *result.Files[0].Subtitles[1].Language, "en")
	assert.Equal(t, *result.Files[0].Subtitles[1].StartTime, float64(72))
	assert.Equal(t, *result.Files[0].Subtitles[1].Duration, float64(71.378))

	assert.Equal(t, *result.Files[0].Bitrate, int64(5407765))
	assert.Equal(t, *result.Files[0].Artist, "Jane")
	assert.Equal(t, *result.Files[0].AlbumArtist, "Jenny")
	assert.Equal(t, *result.Files[0].Composer, "Jane")
	assert.Equal(t, *result.Files[0].Performer, "Jane")
	assert.Equal(t, *result.Files[0].Album, "FirstAlbum")
	assert.Equal(t, *result.Files[0].Duration, float64(71.378))

	assert.Equal(t, len(result.Files[0].Addresses), 2)
	assert.Equal(t, *result.Files[0].Addresses[0].AddressLine, "中国浙江省杭州市余杭区文一西路969号")
	assert.Equal(t, *result.Files[0].Addresses[0].City, "杭州市")
	assert.Equal(t, *result.Files[0].Addresses[0].Country, "中国")
	assert.Equal(t, *result.Files[0].Addresses[0].District, "余杭区")
	assert.Equal(t, *result.Files[0].Addresses[0].Language, "zh-Hans")
	assert.Equal(t, *result.Files[0].Addresses[0].Province, "浙江省")
	assert.Equal(t, *result.Files[0].Addresses[0].Township, "文一西路")

	assert.Equal(t, *result.Files[0].Addresses[1].AddressLine, "中国浙江省杭州市余杭区文一西路970号")
	assert.Equal(t, *result.Files[0].Addresses[1].City, "杭州市")
	assert.Equal(t, *result.Files[0].Addresses[1].Country, "中国")
	assert.Equal(t, *result.Files[0].Addresses[1].District, "余杭区")
	assert.Equal(t, *result.Files[0].Addresses[1].Language, "zh-Hans")
	assert.Equal(t, *result.Files[0].Addresses[1].Province, "浙江省")
	assert.Equal(t, *result.Files[0].Addresses[1].Township, "文一西路")

	assert.Equal(t, *result.Files[0].OSSObjectType, "Normal")
	assert.Equal(t, *result.Files[0].OSSStorageClass, "Standard")
	assert.Equal(t, *result.Files[0].OSSTaggingCount, int64(2))
	assert.Equal(t, *result.Files[0].OSSTagging[0].Key, "key")
	assert.Equal(t, *result.Files[0].OSSTagging[0].Value, "val")
	assert.Equal(t, *result.Files[0].OSSTagging[1].Key, "key2")
	assert.Equal(t, *result.Files[0].OSSTagging[1].Value, "val2")
	assert.Equal(t, len(result.Files[0].OSSUserMeta), 1)
	assert.Equal(t, *result.Files[0].OSSUserMeta[0].Key, "key")
	assert.Equal(t, *result.Files[0].OSSUserMeta[0].Value, "val")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",

		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DoMetaQueryResult{}
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
	result = &DoMetaQueryResult{}
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
	result = &DoMetaQueryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_CloseMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CloseMetaQueryRequest
	var input *OperationInput
	var err error

	request = &CloseMetaQueryRequest{}
	input = &OperationInput{
		OpName: "CloseMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "delete",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &CloseMetaQueryRequest{
		Bucket: Ptr("bucket"),
	}
	request = &CloseMetaQueryRequest{
		Bucket: Ptr("bucket"),
	}
	input = &OperationInput{
		OpName: "CloseMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp":      "delete",
			"metaQuery": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"metaQuery", "comp"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_CloseMetaQuery(t *testing.T) {
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
	result := &CloseMetaQueryResult{}
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
	result = &CloseMetaQueryResult{}
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
	result = &CloseMetaQueryResult{}
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
	result = &CloseMetaQueryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
