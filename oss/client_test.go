package oss

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/retry"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/transport"
	"github.com/stretchr/testify/assert"
)

type stubRequest struct {
	StrPrtField   *string `input:"query,str-field"`
	StrField      string  `input:"query,str-field"`
	IntPtrFiled   *int32  `input:"query,int32-field"`
	IntFiled      int32   `input:"query,int32-field"`
	BoolPtrFiled  *bool   `input:"query,bool-field"`
	HStrPrtField  *string `input:"header,x-oss-str-field"`
	HStrField     string  `input:"header,x-oss-str-field"`
	HIntPtrFiled  *int32  `input:"header,x-oss-int32-field"`
	HIntFiled     int32   `input:"header,x-oss-int32-field"`
	HBoolPtrFiled *bool   `input:"header,x-oss-bool-field"`
}

func TestMarshalInput(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *stubRequest
	var err error

	// nil request
	input = &OperationInput{}
	request = nil

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 0, len(input.Parameters))

	// emtpy request
	input = &OperationInput{}
	request = &stubRequest{}

	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "OperationInput.Method")
	//assert.Equal(t, 0, len(input.Headers))
	//assert.Equal(t, 0, len(input.Parameters))

	// query ptr
	input = &OperationInput{
		Method: "GET",
	}

	request = &stubRequest{
		StrPrtField:  Ptr("str1"),
		IntPtrFiled:  Ptr(int32(123)),
		BoolPtrFiled: Ptr(true),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 3, len(input.Parameters))
	assert.Equal(t, "str1", input.Parameters["str-field"])
	assert.Equal(t, "123", input.Parameters["int32-field"])
	assert.Equal(t, "true", input.Parameters["bool-field"])

	// query value
	input = &OperationInput{
		Method: "GET",
	}

	request = &stubRequest{
		StrField:     "str2",
		IntFiled:     int32(223),
		BoolPtrFiled: Ptr(false),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 3, len(input.Parameters))
	assert.Equal(t, "str2", input.Parameters["str-field"])
	assert.Equal(t, "223", input.Parameters["int32-field"])
	assert.Equal(t, "false", input.Parameters["bool-field"])

	// header ptr
	input = &OperationInput{
		Method: "GET",
	}

	request = &stubRequest{
		HStrPrtField:  Ptr("str1"),
		HIntPtrFiled:  Ptr(int32(123)),
		HBoolPtrFiled: Ptr(true),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Parameters))
	assert.Equal(t, 3, len(input.Headers))
	assert.Equal(t, "str1", input.Headers["x-oss-str-field"])
	assert.Equal(t, "123", input.Headers["x-oss-int32-field"])
	assert.Equal(t, "true", input.Headers["x-oss-bool-field"])

	// header value
	input = &OperationInput{
		Method: "GET",
	}

	request = &stubRequest{
		HStrField:     "str2",
		HIntFiled:     int32(223),
		HBoolPtrFiled: Ptr(false),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Parameters))
	assert.Equal(t, 3, len(input.Headers))
	assert.Equal(t, "str2", input.Headers["x-oss-str-field"])
	assert.Equal(t, "223", input.Headers["x-oss-int32-field"])
	assert.Equal(t, "false", input.Headers["x-oss-bool-field"])
}

type xmlbodyRequest struct {
	StrHostPrtField    *string        `input:"host,bucket,required"`
	StrQueryPrtField   *string        `input:"query,str-field"`
	StrHeaderPrtField  *string        `input:"header,x-oss-str-field"`
	StructBodyPrtField *xmlBodyConfig `input:"body,BodyConfiguration,xml"`
}

type xmlBodyConfig struct {
	StrField1 *string `xml:"StrField1"`
	StrField2 string  `xml:"StrField2"`
}

func TestMarshalInput_xmlbody(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *xmlbodyRequest
	var err error

	input = &OperationInput{
		Method: "GET",
	}
	request = &xmlbodyRequest{
		StrHostPrtField:   Ptr("bucket"),
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		StructBodyPrtField: &xmlBodyConfig{
			StrField1: Ptr("StrField1"),
			StrField2: "StrField2",
		},
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(input.Parameters))
	assert.Equal(t, "query", input.Parameters["str-field"])
	assert.Equal(t, 1, len(input.Headers))
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.NotNil(t, input.Body)

	body, err := io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "<BodyConfiguration><StrField1>StrField1</StrField1><StrField2>StrField2</StrField2></BodyConfiguration>", string(body))
}

type readerBodyRequest struct {
	StrHostPrtField   *string   `input:"host,bucket,required"`
	StrQueryPrtField  *string   `input:"query,str-field"`
	StrHeaderPrtField *string   `input:"header,x-oss-str-field"`
	IoReaderBodyField io.Reader `input:"body,nop"`
}

type notSupportBodyTypeRequest struct {
	StrHostPrtField   *string `input:"host,bucket,required"`
	StrQueryPrtField  *string `input:"query,str-field"`
	StrHeaderPrtField *string `input:"header,x-oss-str-field"`
	StringBodyField   string  `input:"body,nop"`
}

func TestMarshalInput_body(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *readerBodyRequest
	var request1 *notSupportBodyTypeRequest
	var err error

	input = &OperationInput{
		Method: "GET",
	}
	request = &readerBodyRequest{
		StrHostPrtField:   Ptr("bucket"),
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		IoReaderBodyField: bytes.NewReader([]byte("hello world")),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(input.Parameters))
	assert.Equal(t, "query", input.Parameters["str-field"])
	assert.Equal(t, 1, len(input.Headers))
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.NotNil(t, input.Body)
	data, err := io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", string(data))

	// not support body format
	input = &OperationInput{
		Method: "GET",
	}
	request1 = &notSupportBodyTypeRequest{
		StrHostPrtField:   Ptr("bucket"),
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		StringBodyField:   "hello world",
	}
	err = c.marshalInput(request1, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "type not support, StringBodyField")
}

type commonStubRequest struct {
	StrHostPrtField    *string        `input:"host,bucket"`
	StrQueryPrtField   *string        `input:"query,str-field"`
	StrHeaderPrtField  *string        `input:"header,x-oss-str-field"`
	StructBodyPrtField *xmlBodyConfig `input:"body,BodyConfiguration,xml"`
	RequestCommon
}

func TestMarshalInput_CommonFields(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *commonStubRequest
	var err error

	//default
	request = &commonStubRequest{}
	input = &OperationInput{
		Method: "GET",
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Nil(t, input.Body)
	assert.Nil(t, input.Headers)
	assert.Nil(t, input.Parameters)

	//set by request
	request = &commonStubRequest{}
	request.Headers = map[string]string{
		"key": "value",
	}
	request.Parameters = map[string]string{
		"p": "",
	}
	request.Payload = bytes.NewReader([]byte("hello"))
	input = &OperationInput{
		Method: "GET",
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 1)
	assert.Equal(t, "value", input.Headers["key"])
	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 1)
	assert.Equal(t, "", input.Parameters["p"])
	assert.NotNil(t, input.Body)
	data, err := io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(data))

	// priority
	// request commmn > input default
	input = &OperationInput{
		Method: "GET",
		Headers: map[string]string{
			"key": "value1",
		},
		Parameters: map[string]string{
			"p": "value1",
		},
	}
	request = &commonStubRequest{}
	request.Headers = map[string]string{
		"key": "value2",
	}
	request.Parameters = map[string]string{
		"p": "value3",
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 1)
	assert.Equal(t, "value2", input.Headers["key"])
	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 1)
	assert.Equal(t, "value3", input.Parameters["p"])
	assert.Nil(t, input.Body)

	// reuqest filed parametr > request commmn
	input = &OperationInput{
		Method: "GET",
	}
	request = &commonStubRequest{
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		StructBodyPrtField: &xmlBodyConfig{
			StrField1: Ptr("StrField1"),
			StrField2: "StrField2",
		},
	}
	request.Headers = map[string]string{
		"x-oss-str-field": "value2",
	}
	request.Parameters = map[string]string{
		"str-field": "value3",
	}
	request.Payload = bytes.NewReader([]byte("hello"))
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 1)
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 1)
	assert.Equal(t, "query", input.Parameters["str-field"])
	assert.NotNil(t, input.Body)
	data, err = io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "<BodyConfiguration><StrField1>StrField1</StrField1><StrField2>StrField2</StrField2></BodyConfiguration>", string(data))

	// merge, replace
	//reuqest filed parametr > request commmn > input
	input = &OperationInput{
		Method: "GET",
		Headers: map[string]string{
			"input-key": "value1",
		},
		Parameters: map[string]string{
			"input-param":  "value2",
			"input-param1": "value2-1",
		}}
	request = &commonStubRequest{
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
	}
	request.Headers = map[string]string{
		"x-oss-str-field":  "value2",
		"x-oss-str-field1": "value2-1",
	}
	request.Parameters = map[string]string{
		"str-field1": "value3",
	}
	request.Payload = bytes.NewReader([]byte("hello"))
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 3)
	assert.Equal(t, "value1", input.Headers["input-key"])
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.Equal(t, "value2-1", input.Headers["x-oss-str-field1"])
	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 4)
	assert.Equal(t, "value2", input.Parameters["input-param"])
	assert.Equal(t, "value2-1", input.Parameters["input-param1"])
	assert.Equal(t, "query", input.Parameters["str-field"])
	assert.Equal(t, "value3", input.Parameters["str-field1"])
	assert.NotNil(t, input.Body)
	data, err = io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(data))
}

type usermetaRequest struct {
	StrQueryPrtField  *string           `input:"query,str-field"`
	StrHeaderPrtField *string           `input:"header,x-oss-str-field"`
	UserMetaField1    map[string]string `input:"header,x-oss-meta-,usermeta"`
	UserMetaField2    map[string]string `input:"header,x-oss-meta1-,usermeta"`
}

func TestMarshalInput_usermeta(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *usermetaRequest
	var err error

	input = &OperationInput{
		Method: "GET",
	}
	request = &usermetaRequest{}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Nil(t, input.Headers)

	input = &OperationInput{
		Method: "GET",
		Headers: map[string]string{
			"input-key": "value1",
		},
		Parameters: map[string]string{
			"input-param":  "value2",
			"input-param1": "value2-1",
		}}
	request = &usermetaRequest{
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		UserMetaField1: map[string]string{
			"user1": "value1",
			"user2": "value2",
		},
		UserMetaField2: map[string]string{
			"user3": "value3",
			"user4": "value4",
		},
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 6)
	assert.Equal(t, "value1", input.Headers["input-key"])
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.Equal(t, "value1", input.Headers["x-oss-meta-user1"])
	assert.Equal(t, "value2", input.Headers["x-oss-meta-user2"])
	assert.Equal(t, "value3", input.Headers["x-oss-meta1-user3"])
	assert.Equal(t, "value4", input.Headers["x-oss-meta1-user4"])

	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 3)
	assert.Equal(t, "value2", input.Parameters["input-param"])
	assert.Equal(t, "value2-1", input.Parameters["input-param1"])
	assert.Equal(t, "query", input.Parameters["str-field"])
}

type requiredRequest struct {
	StrHostPtrField   *string `input:"host,bucket,required"`
	StrPathField      string  `input:"path,key,required"`
	BoolQueryPtrField *bool   `input:"query,bool-ptr-field,required"`
	IntHeaderField    int     `input:"header,x-oss-str-field,required"`
}

func TestMarshalInput_required(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *requiredRequest
	var err error

	input = &OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "StrHostPtrField")

	input = &OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField: Ptr("host"),
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "StrPathField")

	input = &OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField: Ptr("host"),
		StrPathField:    "path",
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "BoolQueryPtrField")

	input = &OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField:   Ptr("host"),
		StrPathField:      "path",
		BoolQueryPtrField: Ptr(false),
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "IntHeaderField")

	input = &OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField:   Ptr("host"),
		StrPathField:      "path",
		BoolQueryPtrField: Ptr(true),
		IntHeaderField:    int(32),
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
}

type stubResult struct {
	ResultCommon
}

type xmlBodyResult struct {
	ResultCommon
	StrField1 *string `xml:"StrField1"`
	StrField2 *string `xml:"StrField2"`
}

func TestUnmarshalOutput(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var result *stubResult
	var err error

	//empty
	output = &OperationOutput{}
	assert.Nil(t, output.Input)
	assert.Nil(t, output.Body)
	assert.Nil(t, output.Headers)
	assert.Empty(t, output.Status)
	assert.Empty(t, output.StatusCode)

	// with default values
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
	}

	result = &stubResult{}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Nil(t, result.Headers)

	// has header
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Expires":          {"-1"},
			"Content-Length":   {"0"},
			"Content-Encoding": {"gzip"},
		},
	}

	result = &stubResult{}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, "-1", result.Headers.Get("Expires"))
	assert.Equal(t, "0", result.Headers.Get("Content-Length"))
	assert.Equal(t, "gzip", result.Headers.Get("Content-Encoding"))

	// extract body
	body := "<BodyConfiguration><StrField1>StrField1</StrField1><StrField2>StrField2</StrField2></BodyConfiguration>"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	xmlresult := &xmlBodyResult{}
	err = c.unmarshalOutput(xmlresult, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, 200, xmlresult.StatusCode)
	assert.Equal(t, "OK", xmlresult.Status)
	assert.Equal(t, "StrField1", *xmlresult.StrField1)
	assert.Equal(t, "StrField2", *xmlresult.StrField2)
}

type bodyConfiguration struct {
	XMLName   xml.Name `xml:"BodyConfiguration"`
	StrField1 *string  `xml:"StrField1"`
	StrField2 *string  `xml:"StrField2"`
}

type xmlBodyResult2 struct {
	ResultCommon
	XmlStruct *bodyConfiguration `output:"body,BodyConfiguration,xml"`
}

func TestUnmarshalOutput2(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	//empty
	output = &OperationOutput{}
	assert.Nil(t, output.Input)
	assert.Nil(t, output.Body)
	assert.Nil(t, output.Headers)
	assert.Empty(t, output.Status)
	assert.Empty(t, output.StatusCode)

	// extract body to inner filed
	body := "<BodyConfiguration><StrField1>StrField1</StrField1><StrField2>StrField2</StrField2></BodyConfiguration>"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	xmlresult := &xmlBodyResult{}
	err = c.unmarshalOutput(xmlresult, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, xmlresult.StatusCode)
	assert.Equal(t, "OK", xmlresult.Status)
	assert.Equal(t, "StrField1", *xmlresult.StrField1)
	assert.Equal(t, "StrField2", *xmlresult.StrField2)

	// extract body to outer filed, without init
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	xmlresult2 := &xmlBodyResult2{}
	assert.Nil(t, xmlresult2.XmlStruct)
	err = c.unmarshalOutput(xmlresult2, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, xmlresult2.StatusCode)
	assert.Equal(t, "OK", xmlresult2.Status)
	assert.NotNil(t, xmlresult2.XmlStruct)
	assert.Equal(t, "StrField1", *xmlresult2.XmlStruct.StrField1)
	assert.Equal(t, "StrField2", *xmlresult2.XmlStruct.StrField2)

	// extract body to outer filed, init
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	xmlresult2 = &xmlBodyResult2{
		XmlStruct: &bodyConfiguration{},
	}
	assert.NotNil(t, xmlresult2.XmlStruct)
	assert.Nil(t, xmlresult2.XmlStruct.StrField1)
	assert.Nil(t, xmlresult2.XmlStruct.StrField2)
	err = c.unmarshalOutput(xmlresult2, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, xmlresult2.StatusCode)
	assert.Equal(t, "OK", xmlresult2.Status)
	assert.NotNil(t, xmlresult2.XmlStruct)
	assert.Equal(t, "StrField1", *xmlresult2.XmlStruct.StrField1)
	assert.Equal(t, "StrField2", *xmlresult2.XmlStruct.StrField2)
}

type headerStubResult struct {
	Int64PtrField  *int64  `output:"header,int64-ptr-field"`
	Int64Field     int64   `output:"header,Int64-Field"`
	StringPtrField *string `output:"header,String-Ptr-Field"`
	StringField    string  `output:"header,String-Field"`
	BoolPtrField   *bool   `output:"header,Bool-Ptr-Field"`
	BoolField      bool    `output:"header,Bool-Field"`
	EmptyFiled     *string `output:"header,Empty-Field"`
	NoTagField     *string
	TimePrtFiled   *time.Time `output:"header,Time-Ptr-Field,time"`
	TimeFiled      time.Time  `output:"header,Time-Field,time"`

	ResultCommon
}

func TestUnmarshalOutput_header(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var result *headerStubResult
	var err error

	// has header
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Int64-Ptr-Field":  {"-1"},
			"Int64-Field":      {"1000"},
			"String-Ptr-Field": {"text"},
			"String-Field":     {"xml"},
			"Bool-Ptr-Field":   {"false"},
			"Bool-Field":       {"true"},
			"Time-Ptr-Field":   {"Fri, 24 Feb 2012 08:43:27 GMT"},
			"Time-Field":       {"Fri, 24 Feb 2013 08:43:27 GMT"},
		},
	}
	result = &headerStubResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Nil(t, result.EmptyFiled)
	assert.Nil(t, result.NoTagField)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, int64(-1), *result.Int64PtrField)
	assert.Equal(t, int64(1000), result.Int64Field)
	assert.Equal(t, "text", *result.StringPtrField)
	assert.Equal(t, "xml", result.StringField)
	assert.Equal(t, false, *result.BoolPtrField)
	assert.Equal(t, true, result.BoolField)

	assert.Equal(t, "2012-02-24T08:43:27Z", (*result.TimePrtFiled).Format(time.RFC3339))
	assert.Equal(t, "2013-02-24T08:43:27Z", result.TimeFiled.Format(time.RFC3339))

	//low case
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"int64-ptr-field":  {"-1"},
			"int64-field":      {"10001"},
			"string-ptr-field": {"text1"},
			"string-field":     {"xml1"},
			"bool-Ptr-Field":   {"false"},
			"BOOL-FIELD":       {"true"},
			"TIME-Ptr-Field":   {"Fri, 24 Feb 2014 08:43:27 GMT"},
			"Time-FIELD":       {"Fri, 24 Feb 2010 08:43:27 GMT"},
		},
	}
	result = &headerStubResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Nil(t, result.EmptyFiled)
	assert.Nil(t, result.NoTagField)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, int64(-1), *result.Int64PtrField)
	assert.Equal(t, int64(10001), result.Int64Field)
	assert.Equal(t, "text1", *result.StringPtrField)
	assert.Equal(t, "xml1", result.StringField)
	assert.Equal(t, false, *result.BoolPtrField)
	assert.Equal(t, true, result.BoolField)

	assert.Equal(t, "2014-02-24T08:43:27Z", (*result.TimePrtFiled).Format(time.RFC3339))
	assert.Equal(t, "2010-02-24T08:43:27Z", result.TimeFiled.Format(time.RFC3339))

	//primitive type
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Int64-Ptr-Field":  {"-1"},
			"Int64-Field":      {"1000"},
			"String-Ptr-Field": {"text"},
			"String-Field":     {"xml"},
			"Bool-Ptr-Field":   {"false"},
			"Bool-Field":       {"true"},
			"Time-Ptr-Field":   {"Fri, 24 Feb 2012 08:43:27 GMT"},
			"Time-Field":       {"Fri, 24 Feb 2013 08:43:27 GMT"},
		},
	}
	result = &headerStubResult{}
	err = c.unmarshalOutput(result, output, unmarshalHeaderLite)
	assert.Nil(t, err)
	assert.Nil(t, result.EmptyFiled)
	assert.Nil(t, result.NoTagField)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, int64(-1), *result.Int64PtrField)
	assert.Equal(t, int64(1000), result.Int64Field)
	assert.Equal(t, "text", *result.StringPtrField)
	assert.Equal(t, "xml", result.StringField)
	assert.Equal(t, false, *result.BoolPtrField)
	assert.Equal(t, true, result.BoolField)

	assert.Nil(t, result.TimePrtFiled)
	assert.Empty(t, result.TimeFiled)
}

func TestUnmarshalOutput_error(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	// unsupport content-type
	body := "<BodyConfiguration><StrField1>StrField1</StrField1><StrField2>StrField2</StrField2></BodyConfiguration>"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/text"},
		},
	}
	xmlresult := &xmlBodyResult{}
	err = c.unmarshalOutput(xmlresult, output, unmarshalBodyDefault)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupport contentType:application/text")

	// xml decode fail
	body = "StrField1>StrField1</StrField1><StrField2>StrField2<"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	result := &stubResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "XML syntax error on line 1")
}

func TestResolveHTTPClient(t *testing.T) {
	cfg := NewConfig()
	opt := &Options{}
	assert.Nil(t, cfg.HttpClient)
	assert.Nil(t, opt.HttpClient)
	resolveHTTPClient(cfg, opt, nil)
	assert.NotNil(t, opt.HttpClient)

	tr := opt.HttpClient.(*http.Client).Transport
	assert.NotNil(t, tr)
	tran, ok := tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, tran)
	assert.Equal(t, transport.DefaultConnectTimeout, tran.TLSHandshakeTimeout)
	assert.Equal(t, transport.DefaultIdleConnectionTimeout, tran.IdleConnTimeout)
	assert.Equal(t, transport.DefaultMaxConnections, tran.MaxConnsPerHost)
	assert.Equal(t, transport.DefaultExpectContinueTimeout, tran.ExpectContinueTimeout)
	assert.Equal(t, transport.DefaultTLSMinVersion, tran.TLSClientConfig.MinVersion)
	assert.Equal(t, false, tran.TLSClientConfig.InsecureSkipVerify)
	assert.NotNil(t, opt.HttpClient.(*http.Client).CheckRedirect)
	assert.Nil(t, tran.Proxy)

	cfg = &Config{
		ConnectTimeout:       Ptr(101 * time.Second),
		ReadWriteTimeout:     Ptr(101 * time.Second),
		InsecureSkipVerify:   Ptr(true),
		EnabledRedirect:      Ptr(true),
		ProxyFromEnvironment: Ptr(true),
	}
	opt = &Options{}
	resolveHTTPClient(cfg, opt, nil)
	assert.NotNil(t, opt.HttpClient)
	tr = opt.HttpClient.(*http.Client).Transport
	assert.NotNil(t, tr)
	tran, ok = tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, tran)
	assert.Equal(t, 101*time.Second, tran.TLSHandshakeTimeout)
	assert.Equal(t, transport.DefaultIdleConnectionTimeout, tran.IdleConnTimeout)
	assert.Equal(t, transport.DefaultMaxConnections, tran.MaxConnsPerHost)
	assert.Equal(t, transport.DefaultExpectContinueTimeout, tran.ExpectContinueTimeout)
	assert.Equal(t, transport.DefaultTLSMinVersion, tran.TLSClientConfig.MinVersion)
	assert.Equal(t, true, tran.TLSClientConfig.InsecureSkipVerify)
	assert.Nil(t, opt.HttpClient.(*http.Client).CheckRedirect)
	assert.NotNil(t, tran.Proxy)
}

func TestHasFeature(t *testing.T) {
	c := &Client{}
	assert.False(t, c.hasFeature(FeatureCorrectClockSkew))

	c.options.FeatureFlags = FeatureFlagsDefault
	assert.True(t, c.hasFeature(FeatureCorrectClockSkew))

	c.options.FeatureFlags = 0xf001
	assert.True(t, c.hasFeature(FeatureCorrectClockSkew))
	assert.True(t, c.hasFeature(0x3))
	assert.False(t, c.hasFeature(0x2))
	assert.True(t, c.hasFeature(0xf000))

	cfg := NewConfig()
	c = NewClient(cfg)
	assert.False(t, c.hasFeature(0))
	assert.True(t, c.hasFeature(FeatureCorrectClockSkew))
	assert.True(t, c.hasFeature(FeatureAutoDetectMimeType))
	assert.True(t, c.hasFeature(FeatureEnableCRC64CheckUpload))
	assert.True(t, c.hasFeature(FeatureEnableCRC64CheckDownload))

	// Disable FeatureEnableCRC64CheckUpload
	cfg = NewConfig()
	cfg.WithDisableUploadCRC64Check(true)
	c = NewClient(cfg)
	assert.False(t, c.hasFeature(0))
	assert.True(t, c.hasFeature(FeatureCorrectClockSkew))
	assert.True(t, c.hasFeature(FeatureAutoDetectMimeType))
	assert.False(t, c.hasFeature(FeatureEnableCRC64CheckUpload))
	assert.True(t, c.hasFeature(FeatureEnableCRC64CheckDownload))

	// Disable FeatureEnableCRC64CheckDownload
	cfg.WithDisableDownloadCRC64Check(true)
	c = NewClient(cfg)
	assert.False(t, c.hasFeature(0))
	assert.True(t, c.hasFeature(FeatureCorrectClockSkew))
	assert.True(t, c.hasFeature(FeatureAutoDetectMimeType))
	assert.False(t, c.hasFeature(FeatureEnableCRC64CheckUpload))
	assert.False(t, c.hasFeature(FeatureEnableCRC64CheckDownload))

	cfg = NewConfig()
	cfg.WithDisableUploadCRC64Check(false)
	cfg.WithDisableDownloadCRC64Check(false)
	c = NewClient(cfg)
	assert.False(t, c.hasFeature(0))
	assert.True(t, c.hasFeature(FeatureCorrectClockSkew))
	assert.True(t, c.hasFeature(FeatureAutoDetectMimeType))
	assert.True(t, c.hasFeature(FeatureEnableCRC64CheckUpload))
	assert.True(t, c.hasFeature(FeatureEnableCRC64CheckDownload))
}

func TestFeatureCorrectClockSkew(t *testing.T) {
	serverTime, _ := http.ParseTime("Sun, 12 Nov 2023 16:56:44 GMT")

	// current time < servertime
	clientTime, _ := http.ParseTime("Sun, 12 Nov 2022 16:56:44 GMT")
	clockOffset := serverTime.Sub(clientTime)
	assert.True(t, clockOffset > 0)
	curr := clientTime.Add(clockOffset)
	assert.Equal(t, serverTime.UTC(), curr.UTC())

	// current time > servertime
	clientTime = time.Now()
	clockOffset = serverTime.Sub(clientTime)
	assert.True(t, clockOffset < 0)
	curr = clientTime.Add(clockOffset)
	assert.Equal(t, serverTime.UTC(), curr.UTC())
}

func TestResolveEndpoint(t *testing.T) {
	cfg := NewConfig()
	opt := &Options{}
	assert.Nil(t, cfg.Region)
	assert.Nil(t, cfg.Endpoint)

	//No Endpont & No Region
	resolveEndpoint(cfg, opt)
	assert.Nil(t, opt.Endpoint)

	//Endpont + ssl
	cfg = NewConfig()
	opt = &Options{}
	cfg.Endpoint = Ptr("test-endpoint")
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "test-endpoint", opt.Endpoint.Host)
	assert.Equal(t, "https", opt.Endpoint.Scheme)

	cfg = NewConfig()
	opt = &Options{}
	cfg.Endpoint = Ptr("test-endpoint")
	cfg.DisableSSL = Ptr(true)
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "test-endpoint", opt.Endpoint.Host)
	assert.Equal(t, "http", opt.Endpoint.Scheme)

	cfg = NewConfig()
	opt = &Options{}
	cfg.Endpoint = Ptr("test-endpoint")
	cfg.DisableSSL = Ptr(false)
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "test-endpoint", opt.Endpoint.Host)
	assert.Equal(t, "https", opt.Endpoint.Scheme)

	cfg = NewConfig()
	opt = &Options{}
	cfg.Endpoint = Ptr("http://test-endpoint")
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "test-endpoint", opt.Endpoint.Host)
	assert.Equal(t, "http", opt.Endpoint.Scheme)

	//Region + ssl
	cfg = NewConfig()
	opt = &Options{}
	cfg.Region = Ptr("test-region")
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "oss-test-region.aliyuncs.com", opt.Endpoint.Host)
	assert.Equal(t, "https", opt.Endpoint.Scheme)

	cfg = NewConfig()
	opt = &Options{}
	cfg.Region = Ptr("test-region")
	cfg.DisableSSL = Ptr(true)
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "oss-test-region.aliyuncs.com", opt.Endpoint.Host)
	assert.Equal(t, "http", opt.Endpoint.Scheme)

	cfg = NewConfig()
	opt = &Options{}
	cfg.Region = Ptr("test-region")
	cfg.DisableSSL = Ptr(false)
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "oss-test-region.aliyuncs.com", opt.Endpoint.Host)
	assert.Equal(t, "https", opt.Endpoint.Scheme)

	cfg = NewConfig()
	opt = &Options{}
	cfg.Region = Ptr("test-region")
	cfg.UseInternalEndpoint = Ptr(true)
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "oss-test-region-internal.aliyuncs.com", opt.Endpoint.Host)
	assert.Equal(t, "https", opt.Endpoint.Scheme)

	cfg = NewConfig()
	opt = &Options{}
	cfg.Region = Ptr("test-region")
	cfg.UseDualStackEndpoint = Ptr(true)
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "test-region.oss.aliyuncs.com", opt.Endpoint.Host)
	assert.Equal(t, "https", opt.Endpoint.Scheme)

	cfg = NewConfig()
	opt = &Options{}
	cfg.Region = Ptr("test-region")
	cfg.UseAccelerateEndpoint = Ptr(true)
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "oss-accelerate.aliyuncs.com", opt.Endpoint.Host)
	assert.Equal(t, "https", opt.Endpoint.Scheme)

	cfg = NewConfig()
	opt = &Options{}
	cfg.Region = Ptr("test-region")
	cfg.UseInternalEndpoint = Ptr(false)
	cfg.UseDualStackEndpoint = Ptr(false)
	cfg.UseAccelerateEndpoint = Ptr(false)
	resolveEndpoint(cfg, opt)
	assert.NotNil(t, opt.Endpoint)
	assert.Equal(t, "oss-test-region.aliyuncs.com", opt.Endpoint.Host)
	assert.Equal(t, "https", opt.Endpoint.Scheme)
}

func TestEndpoint(t *testing.T) {
	//No Endpont & No Region
	cfg := NewConfig()
	c := NewClient(cfg)
	assert.Nil(t, c.options.Endpoint)

	//Endpont + ssl
	cfg = NewConfig()
	cfg.WithEndpoint("test-endpoint")
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "test-endpoint", c.options.Endpoint.Host)
	assert.Equal(t, "https", c.options.Endpoint.Scheme)

	cfg = NewConfig()
	cfg.WithEndpoint("test-endpoint")
	cfg.WithDisableSSL(true)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "test-endpoint", c.options.Endpoint.Host)
	assert.Equal(t, "http", c.options.Endpoint.Scheme)

	cfg = NewConfig()
	cfg.WithEndpoint("test-endpoint")
	cfg.WithDisableSSL(false)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "test-endpoint", c.options.Endpoint.Host)
	assert.Equal(t, "https", c.options.Endpoint.Scheme)

	cfg = NewConfig()
	cfg.WithEndpoint("http://test-endpoint")
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "test-endpoint", c.options.Endpoint.Host)
	assert.Equal(t, "http", c.options.Endpoint.Scheme)

	//Region + ssl
	cfg = NewConfig()
	cfg.WithRegion("test-region")
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "oss-test-region.aliyuncs.com", c.options.Endpoint.Host)
	assert.Equal(t, "https", c.options.Endpoint.Scheme)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithDisableSSL(true)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "oss-test-region.aliyuncs.com", c.options.Endpoint.Host)
	assert.Equal(t, "http", c.options.Endpoint.Scheme)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithUseInternalEndpoint(true)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "oss-test-region-internal.aliyuncs.com", c.options.Endpoint.Host)
	assert.Equal(t, "https", c.options.Endpoint.Scheme)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithUseDualStackEndpoint(true)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "test-region.oss.aliyuncs.com", c.options.Endpoint.Host)
	assert.Equal(t, "https", c.options.Endpoint.Scheme)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithUseAccelerateEndpoint(true)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "oss-accelerate.aliyuncs.com", c.options.Endpoint.Host)
	assert.Equal(t, "https", c.options.Endpoint.Scheme)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithUseInternalEndpoint(false)
	cfg.WithUseDualStackEndpoint(false)
	cfg.WithUseAccelerateEndpoint(false)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "oss-test-region.aliyuncs.com", c.options.Endpoint.Host)
	assert.Equal(t, "https", c.options.Endpoint.Scheme)
}

func TestSinger(t *testing.T) {
	cfg := NewConfig()
	c := NewClient(cfg)
	assert.NotNil(t, c.options.Signer)

	v4, ok := c.options.Signer.(*signer.SignerV4)
	assert.NotNil(t, v4)
	assert.True(t, ok)

	cfg = NewConfig()
	cfg.WithSignatureVersion(SignatureVersionV1)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Signer)
	v1, ok := c.options.Signer.(*signer.SignerV1)
	assert.NotNil(t, v1)
	assert.True(t, ok)

	cfg = NewConfig()
	cfg.WithSignatureVersion(SignatureVersionV4)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.Signer)
	v4, ok = c.options.Signer.(*signer.SignerV4)
	assert.NotNil(t, v4)
	assert.True(t, ok)
}

func TestRetryMaxAttempts(t *testing.T) {
	cfg := NewConfig()
	c := NewClient(cfg)
	assert.Nil(t, c.options.RetryMaxAttempts)

	assert.Equal(t, retry.DefaultMaxAttempts, c.retryMaxAttempts(nil))

	cfg = NewConfig()
	cfg.RetryMaxAttempts = Ptr(5)
	c = NewClient(cfg)
	assert.NotNil(t, c.options.RetryMaxAttempts)
	assert.Equal(t, 5, c.retryMaxAttempts(nil))
}

func TestUserAgent(t *testing.T) {
	cfg := NewConfig()
	c := NewClient(cfg)
	assert.NotEmpty(t, defaultUserAgent)
	assert.Equal(t, defaultUserAgent, c.inner.UserAgent)

	cfg = NewConfig()
	cfg.UserAgent = Ptr("my-user-agent")
	c = NewClient(cfg)
	assert.Equal(t, defaultUserAgent+"/my-user-agent", c.inner.UserAgent)
}

func TestCloudBoxId(t *testing.T) {
	//default product
	cfg := NewConfig()
	c := NewClient(cfg)
	assert.Equal(t, "oss", c.options.Product)

	// default region, endpiont
	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("test-endpoint")
	c = NewClient(cfg)
	assert.Equal(t, "oss", c.options.Product)
	assert.Equal(t, "test-region", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "test-endpoint", c.options.Endpoint.Host)

	// set cloudbox id
	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("test-endpoint")
	cfg.WithCloudBoxId("test-cloudbox-id")
	c = NewClient(cfg)
	assert.Equal(t, "oss-cloudbox", c.options.Product)
	assert.Equal(t, "test-cloudbox-id", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "test-endpoint", c.options.Endpoint.Host)

	//cb-***.{region}.oss-cloudbox-control.aliyuncs.com
	//cb-***.{region}.oss-cloudbox.aliyuncs.com

	// auto detect cloudbox id default
	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("cb-123.test-region.oss-cloudbox-control.aliyuncs.com")
	c = NewClient(cfg)
	assert.Equal(t, "oss", c.options.Product)
	assert.Equal(t, "test-region", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "cb-123.test-region.oss-cloudbox-control.aliyuncs.com", c.options.Endpoint.Host)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("cb-123.test-region.oss-cloudbox.aliyuncs.com")
	c = NewClient(cfg)
	assert.Equal(t, "oss", c.options.Product)
	assert.Equal(t, "test-region", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "cb-123.test-region.oss-cloudbox.aliyuncs.com", c.options.Endpoint.Host)

	// auto detect cloudbox id set false
	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("cb-123.test-region.oss-cloudbox-control.aliyuncs.com")
	cfg.WithEnableAutoDetectCloudBoxId(false)
	c = NewClient(cfg)
	assert.Equal(t, "oss", c.options.Product)
	assert.Equal(t, "test-region", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "cb-123.test-region.oss-cloudbox-control.aliyuncs.com", c.options.Endpoint.Host)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("cb-123.test-region.oss-cloudbox.aliyuncs.com")
	cfg.WithEnableAutoDetectCloudBoxId(false)
	c = NewClient(cfg)
	assert.Equal(t, "oss", c.options.Product)
	assert.Equal(t, "test-region", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "cb-123.test-region.oss-cloudbox.aliyuncs.com", c.options.Endpoint.Host)

	// auto detect cloudbox id set true
	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("cb-123.test-region.oss-cloudbox-control.aliyuncs.com")
	cfg.WithEnableAutoDetectCloudBoxId(true)
	c = NewClient(cfg)
	assert.Equal(t, "oss-cloudbox", c.options.Product)
	assert.Equal(t, "cb-123", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "cb-123.test-region.oss-cloudbox-control.aliyuncs.com", c.options.Endpoint.Host)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("cb-123.test-region.oss-cloudbox.aliyuncs.com")
	cfg.WithEnableAutoDetectCloudBoxId(true)
	c = NewClient(cfg)
	assert.Equal(t, "oss-cloudbox", c.options.Product)
	assert.Equal(t, "cb-123", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "cb-123.test-region.oss-cloudbox.aliyuncs.com", c.options.Endpoint.Host)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("cb-123.test-region.oss-cloudbox.aliyuncs.com/test?123")
	cfg.WithEnableAutoDetectCloudBoxId(true)
	c = NewClient(cfg)
	assert.Equal(t, "oss-cloudbox", c.options.Product)
	assert.Equal(t, "cb-123", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "cb-123.test-region.oss-cloudbox.aliyuncs.com", c.options.Endpoint.Host)

	// auto detect cloudbox id set true + non cloud box endpoint
	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("cb-123.test-region.oss.aliyuncs.com")
	cfg.WithEnableAutoDetectCloudBoxId(true)
	c = NewClient(cfg)
	assert.Equal(t, "oss", c.options.Product)
	assert.Equal(t, "test-region", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "cb-123.test-region.oss.aliyuncs.com", c.options.Endpoint.Host)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("ncb-123.test-region.oss-cloudbox.aliyuncs.com")
	cfg.WithEnableAutoDetectCloudBoxId(true)
	c = NewClient(cfg)
	assert.Equal(t, "oss", c.options.Product)
	assert.Equal(t, "test-region", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "ncb-123.test-region.oss-cloudbox.aliyuncs.com", c.options.Endpoint.Host)

	cfg = NewConfig()
	cfg.WithRegion("test-region")
	cfg.WithEndpoint("cb-123.oss-cloudbox.aliyuncs.com")
	cfg.WithEnableAutoDetectCloudBoxId(true)
	c = NewClient(cfg)
	assert.Equal(t, "oss", c.options.Product)
	assert.Equal(t, "test-region", c.options.Region)
	assert.NotNil(t, c.options.Endpoint)
	assert.Equal(t, "cb-123.oss-cloudbox.aliyuncs.com", c.options.Endpoint.Host)
}
