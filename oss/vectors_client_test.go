package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_ForVectorsClient(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *stubRequest
	var err error

	// nil request
	input = &OperationInput{}
	request = nil

	err = c.client.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 0, len(input.Parameters))

	// emtpy request
	input = &OperationInput{}
	request = &stubRequest{}

	err = c.client.marshalInput(request, input)
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

	err = c.client.marshalInput(request, input)
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

	err = c.client.marshalInput(request, input)
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

	err = c.client.marshalInput(request, input)
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

	err = c.client.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Parameters))
	assert.Equal(t, 3, len(input.Headers))
	assert.Equal(t, "str2", input.Headers["x-oss-str-field"])
	assert.Equal(t, "223", input.Headers["x-oss-int32-field"])
	assert.Equal(t, "false", input.Headers["x-oss-bool-field"])
}

type jsonbodyRequest struct {
	StrHostPrtField    *string         `input:"host,bucket,required"`
	StrQueryPrtField   *string         `input:"query,str-field"`
	StrHeaderPrtField  *string         `input:"header,x-oss-str-field"`
	StructBodyPrtField *jsonBodyConfig `input:"body,BodyConfiguration,json"`
}

type jsonBodyConfig struct {
	StrField1 *string `json:"StrField1"`
	StrField2 string  `json:"StrField2"`
}

func TestMarshalInput_JsonBody_ForVectorsClient(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *jsonbodyRequest
	var err error

	input = &OperationInput{
		Method: "GET",
	}
	request = &jsonbodyRequest{
		StrHostPrtField:   Ptr("bucket"),
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		StructBodyPrtField: &jsonBodyConfig{
			StrField1: Ptr("StrField1"),
			StrField2: "StrField2",
		},
	}

	err = c.client.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(input.Parameters))
	assert.Equal(t, "query", input.Parameters["str-field"])
	assert.Equal(t, 1, len(input.Headers))
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.NotNil(t, input.Body)

	body, err := io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"BodyConfiguration\":{\"StrField1\":\"StrField1\",\"StrField2\":\"StrField2\"}}", string(body))
}

func TestMarshalInput_body_ForVectorsClient(t *testing.T) {
	c := VectorsClient{}
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

	err = c.client.marshalInput(request, input)
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
	err = c.client.marshalInput(request1, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "type not support, StringBodyField")
}

type commonStubVectorRequest struct {
	StrHostPrtField    *string         `input:"host,bucket"`
	StrQueryPrtField   *string         `input:"query,str-field"`
	StrHeaderPrtField  *string         `input:"header,x-oss-str-field"`
	StructBodyPrtField *jsonBodyConfig `input:"body,BodyConfiguration,json"`
	RequestCommon
}

func TestMarshalInput_CommonFields_ForVectorsClient(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *commonStubVectorRequest
	var err error

	//default
	request = &commonStubVectorRequest{}
	input = &OperationInput{
		Method: "GET",
	}
	err = c.client.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Nil(t, input.Body)
	assert.Nil(t, input.Headers)
	assert.Nil(t, input.Parameters)

	//set by request
	request = &commonStubVectorRequest{}
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
	err = c.client.marshalInput(request, input)
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
	request = &commonStubVectorRequest{}
	request.Headers = map[string]string{
		"key": "value2",
	}
	request.Parameters = map[string]string{
		"p": "value3",
	}
	err = c.client.marshalInput(request, input)
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
	request = &commonStubVectorRequest{
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		StructBodyPrtField: &jsonBodyConfig{
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
	err = c.client.marshalInput(request, input)
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
	assert.Equal(t, "{\"BodyConfiguration\":{\"StrField1\":\"StrField1\",\"StrField2\":\"StrField2\"}}", string(data))

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
	request = &commonStubVectorRequest{
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
	err = c.client.marshalInput(request, input)
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

func TestMarshalInput_UserMeta_ForVectorsClient(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *usermetaRequest
	var err error

	input = &OperationInput{
		Method: "GET",
	}
	request = &usermetaRequest{}
	err = c.client.marshalInput(request, input)
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
	err = c.client.marshalInput(request, input)
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

func TestMarshalInput_required_ForVectorsClient(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *requiredRequest
	var err error

	input = &OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{}
	err = c.client.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "StrHostPtrField")

	input = &OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField: Ptr("host"),
	}
	err = c.client.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "StrPathField")

	input = &OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField: Ptr("host"),
		StrPathField:    "path",
	}
	err = c.client.marshalInput(request, input)
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
	err = c.client.marshalInput(request, input)
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
	err = c.client.marshalInput(request, input)
	assert.Nil(t, err)
}

type jsonBodyResult struct {
	ResultCommon
	StrField1 *string `json:"StrField1"`
	StrField2 *string `json:"StrField2"`
}

func TestUnmarshalOutput_ForVectorsClient(t *testing.T) {
	c := VectorsClient{}
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
	err = c.client.unmarshalOutput(result, output)
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
	err = c.client.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, "-1", result.Headers.Get("Expires"))
	assert.Equal(t, "0", result.Headers.Get("Content-Length"))
	assert.Equal(t, "gzip", result.Headers.Get("Content-Encoding"))

	// extract body
	body := "{\"BodyConfiguration\":{\"StrField1\":\"StrField1\",\"StrField2\":\"StrField2\"}}"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/json"},
		},
	}
	jsonresult := &jsonBodyResult{}
	err = c.client.unmarshalOutput(jsonresult, output, unmarshalBodyDefaultV2)

	assert.Nil(t, err)
	assert.Equal(t, 200, jsonresult.StatusCode)
	assert.Equal(t, "OK", jsonresult.Status)
	assert.Equal(t, "StrField1", *jsonresult.StrField1)
	assert.Equal(t, "StrField2", *jsonresult.StrField2)
}

func TestUnmarshalOutput_header_ForVectorsClient(t *testing.T) {
	c := VectorsClient{}
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
	err = c.client.unmarshalOutput(result, output, unmarshalHeader)
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
	err = c.client.unmarshalOutput(result, output, unmarshalHeader)
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
	err = c.client.unmarshalOutput(result, output, unmarshalHeaderLite)
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

func TestUnmarshalOutput_error_ForVectorsClient(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	// unsupport content-type
	body := "{\"BodyConfiguration\":{\"StrField1\":\"StrField1\",\"StrField2\":\"StrField2\"}}"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/text"},
		},
	}
	jsonresult := &jsonBodyResult{}
	err = c.client.unmarshalOutput(jsonresult, output, unmarshalBodyDefaultV2)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupport contentType:application/text")

	// xml decode fail
	body = "StrField1>StrField1</StrField1><StrField2>StrField2<"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/json"},
		},
	}
	result := &stubResult{}
	err = c.client.unmarshalOutput(result, output, unmarshalBodyDefaultV2)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "deserialization failed, invalid json format: no opening brace found")
}

func TestResolveEndpoint_ForVectorsClient(t *testing.T) {
	cfg := NewConfig()
	assert.Nil(t, cfg.Region)
	assert.Nil(t, cfg.Endpoint)

	//No Endpont & No Region
	resolveVectorsEndpoint(cfg)
	assert.Nil(t, cfg.Endpoint)

	//Endpont + ssl
	cfg = NewConfig()
	cfg.Endpoint = Ptr("test-endpoint")
	resolveVectorsEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "https://test-endpoint", ToString(cfg.Endpoint))

	cfg = NewConfig()
	cfg.Endpoint = Ptr("test-endpoint")
	cfg.DisableSSL = Ptr(true)
	resolveVectorsEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "http://test-endpoint", ToString(cfg.Endpoint))

	cfg = NewConfig()
	cfg.Endpoint = Ptr("test-endpoint")
	cfg.DisableSSL = Ptr(false)
	resolveVectorsEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "https://test-endpoint", ToString(cfg.Endpoint))

	cfg = NewConfig()
	cfg.Endpoint = Ptr("http://test-endpoint")
	resolveVectorsEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "http://test-endpoint", ToString(cfg.Endpoint))

	//Region + ssl
	cfg = NewConfig()
	cfg.Region = Ptr("test-region")
	resolveVectorsEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "https://oss-test-region.oss-vectors.aliyuncs.com", ToString(cfg.Endpoint))

	cfg = NewConfig()
	cfg.Region = Ptr("test-region")
	cfg.DisableSSL = Ptr(true)
	resolveVectorsEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "http://oss-test-region.oss-vectors.aliyuncs.com", ToString(cfg.Endpoint))

	cfg = NewConfig()
	cfg.Region = Ptr("test-region")
	cfg.DisableSSL = Ptr(false)
	resolveVectorsEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "https://oss-test-region.oss-vectors.aliyuncs.com", ToString(cfg.Endpoint))

	cfg = NewConfig()
	cfg.Region = Ptr("test-region")
	cfg.UseInternalEndpoint = Ptr(true)
	resolveVectorsEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "https://oss-test-region-internal.oss-vectors.aliyuncs.com", ToString(cfg.Endpoint))

	cfg = NewConfig()
	cfg.Region = Ptr("test-region")
	cfg.UseInternalEndpoint = Ptr(false)
	cfg.UseDualStackEndpoint = Ptr(false)
	cfg.UseAccelerateEndpoint = Ptr(false)
	resolveVectorsEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "https://oss-test-region.oss-vectors.aliyuncs.com", ToString(cfg.Endpoint))
}

func TestResolveVectorsUserAgent_ForVectorsClient(t *testing.T) {
	cfg := NewConfig()
	resolveVectorsUserAgent(cfg)
	assert.Equal(t, ToString(cfg.UserAgent), VectorsUserAgentPrefix)

	cfg = NewConfig()
	c := NewVectorsClient(cfg)
	assert.Equal(t, defaultUserAgent+"/"+VectorsUserAgentPrefix, c.client.inner.UserAgent)

	cfg = NewConfig()
	cfg.UserAgent = Ptr("my-user-agent")
	resolveVectorsUserAgent(cfg)
	assert.Equal(t, ToString(cfg.UserAgent), VectorsUserAgentPrefix+"/"+"my-user-agent")

	cfg = NewConfig()
	cfg.UserAgent = Ptr("my-user-agent")
	c = NewVectorsClient(cfg)
	assert.Equal(t, defaultUserAgent+"/"+VectorsUserAgentPrefix+"/"+"my-user-agent", c.client.inner.UserAgent)

}

func TestSinger_ForVectorsClient(t *testing.T) {
	cfg := NewConfig()
	c := NewVectorsClient(cfg)
	assert.NotNil(t, c.client.options.Signer)

	v4, ok := c.client.options.Signer.(*signer.SignerVectorsV4)
	assert.NotNil(t, v4)
	assert.True(t, ok)

	cfg = NewConfig()
	cfg.WithSignatureVersion(SignatureVersionV4)
	c = NewVectorsClient(cfg)
	assert.NotNil(t, c.client.options.Signer)
	v4, ok = c.client.options.Signer.(*signer.SignerVectorsV4)
	assert.NotNil(t, v4)
	assert.True(t, ok)
}
