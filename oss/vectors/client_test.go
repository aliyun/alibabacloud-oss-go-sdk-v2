package vectors

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
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

func TestMarshalInput(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *oss.OperationInput
	var request *stubRequest
	var err error

	// nil request
	input = &oss.OperationInput{}
	request = nil

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 0, len(input.Parameters))

	// emtpy request
	input = &oss.OperationInput{}
	request = &stubRequest{}

	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid field, OperationInput.Method")
	//assert.Equal(t, 0, len(input.Headers))
	//assert.Equal(t, 0, len(input.Parameters))

	// query ptr
	input = &oss.OperationInput{
		Method: "GET",
	}

	request = &stubRequest{
		StrPrtField:  oss.Ptr("str1"),
		IntPtrFiled:  oss.Ptr(int32(123)),
		BoolPtrFiled: oss.Ptr(true),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 3, len(input.Parameters))
	assert.Equal(t, "str1", input.Parameters["str-field"])
	assert.Equal(t, "123", input.Parameters["int32-field"])
	assert.Equal(t, "true", input.Parameters["bool-field"])

	// query value
	input = &oss.OperationInput{
		Method: "GET",
	}

	request = &stubRequest{
		StrField:     "str2",
		IntFiled:     int32(223),
		BoolPtrFiled: oss.Ptr(false),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 3, len(input.Parameters))
	assert.Equal(t, "str2", input.Parameters["str-field"])
	assert.Equal(t, "223", input.Parameters["int32-field"])
	assert.Equal(t, "false", input.Parameters["bool-field"])

	// header ptr
	input = &oss.OperationInput{
		Method: "GET",
	}

	request = &stubRequest{
		HStrPrtField:  oss.Ptr("str1"),
		HIntPtrFiled:  oss.Ptr(int32(123)),
		HBoolPtrFiled: oss.Ptr(true),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Parameters))
	assert.Equal(t, 3, len(input.Headers))
	assert.Equal(t, "str1", input.Headers["x-oss-str-field"])
	assert.Equal(t, "123", input.Headers["x-oss-int32-field"])
	assert.Equal(t, "true", input.Headers["x-oss-bool-field"])

	// header value
	input = &oss.OperationInput{
		Method: "GET",
	}

	request = &stubRequest{
		HStrField:     "str2",
		HIntFiled:     int32(223),
		HBoolPtrFiled: oss.Ptr(false),
	}

	err = c.marshalInput(request, input)
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

func TestMarshalInput_JsonBody(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *oss.OperationInput
	var request *jsonbodyRequest
	var err error

	input = &oss.OperationInput{
		Method: "GET",
	}
	request = &jsonbodyRequest{
		StrHostPrtField:   oss.Ptr("bucket"),
		StrQueryPrtField:  oss.Ptr("query"),
		StrHeaderPrtField: oss.Ptr("header"),
		StructBodyPrtField: &jsonBodyConfig{
			StrField1: oss.Ptr("StrField1"),
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
	assert.Equal(t, "{\"BodyConfiguration\":{\"StrField1\":\"StrField1\",\"StrField2\":\"StrField2\"}}", string(body))
}

func TestMarshalInput_body(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *oss.OperationInput
	var request *readerBodyRequest
	var request1 *notSupportBodyTypeRequest
	var err error

	input = &oss.OperationInput{
		Method: "GET",
	}
	request = &readerBodyRequest{
		StrHostPrtField:   oss.Ptr("bucket"),
		StrQueryPrtField:  oss.Ptr("query"),
		StrHeaderPrtField: oss.Ptr("header"),
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
	input = &oss.OperationInput{
		Method: "GET",
	}
	request1 = &notSupportBodyTypeRequest{
		StrHostPrtField:   oss.Ptr("bucket"),
		StrQueryPrtField:  oss.Ptr("query"),
		StrHeaderPrtField: oss.Ptr("header"),
		StringBodyField:   "hello world",
	}
	err = c.marshalInput(request1, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "type not support, StringBodyField")
}

type commonStubVectorRequest struct {
	StrHostPrtField    *string         `input:"host,bucket"`
	StrQueryPrtField   *string         `input:"query,str-field"`
	StrHeaderPrtField  *string         `input:"header,x-oss-str-field"`
	StructBodyPrtField *jsonBodyConfig `input:"body,BodyConfiguration,json"`
	oss.RequestCommon
}

func TestMarshalInput_CommonFields(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *oss.OperationInput
	var request *commonStubVectorRequest
	var err error

	//default
	request = &commonStubVectorRequest{}
	input = &oss.OperationInput{
		Method: "GET",
	}
	err = c.marshalInput(request, input)
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
	input = &oss.OperationInput{
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
	input = &oss.OperationInput{
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
	input = &oss.OperationInput{
		Method: "GET",
	}
	request = &commonStubVectorRequest{
		StrQueryPrtField:  oss.Ptr("query"),
		StrHeaderPrtField: oss.Ptr("header"),
		StructBodyPrtField: &jsonBodyConfig{
			StrField1: oss.Ptr("StrField1"),
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
	assert.Equal(t, "{\"BodyConfiguration\":{\"StrField1\":\"StrField1\",\"StrField2\":\"StrField2\"}}", string(data))

	// merge, replace
	//reuqest filed parametr > request commmn > input
	input = &oss.OperationInput{
		Method: "GET",
		Headers: map[string]string{
			"input-key": "value1",
		},
		Parameters: map[string]string{
			"input-param":  "value2",
			"input-param1": "value2-1",
		}}
	request = &commonStubVectorRequest{
		StrQueryPrtField:  oss.Ptr("query"),
		StrHeaderPrtField: oss.Ptr("header"),
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

func TestMarshalInput_UserMeta(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *oss.OperationInput
	var request *usermetaRequest
	var err error

	input = &oss.OperationInput{
		Method: "GET",
	}
	request = &usermetaRequest{}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Nil(t, input.Headers)

	input = &oss.OperationInput{
		Method: "GET",
		Headers: map[string]string{
			"input-key": "value1",
		},
		Parameters: map[string]string{
			"input-param":  "value2",
			"input-param1": "value2-1",
		}}
	request = &usermetaRequest{
		StrQueryPrtField:  oss.Ptr("query"),
		StrHeaderPrtField: oss.Ptr("header"),
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
	c := VectorsClient{}
	assert.NotNil(t, c)
	var input *oss.OperationInput
	var request *requiredRequest
	var err error

	input = &oss.OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "StrHostPtrField")

	input = &oss.OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField: oss.Ptr("host"),
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "StrPathField")

	input = &oss.OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField: oss.Ptr("host"),
		StrPathField:    "path",
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "BoolQueryPtrField")

	input = &oss.OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField:   oss.Ptr("host"),
		StrPathField:      "path",
		BoolQueryPtrField: oss.Ptr(false),
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "IntHeaderField")

	input = &oss.OperationInput{
		Method: "GET",
	}
	request = &requiredRequest{
		StrHostPtrField:   oss.Ptr("host"),
		StrPathField:      "path",
		BoolQueryPtrField: oss.Ptr(true),
		IntHeaderField:    int(32),
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
}

type jsonBodyResult struct {
	oss.ResultCommon
	StrField1 *string `json:"StrField1"`
	StrField2 *string `json:"StrField2"`
}

type stubResult struct {
	oss.ResultCommon
}

func TestUnmarshalOutput(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var result *stubResult
	var err error

	//empty
	output = &oss.OperationOutput{}
	assert.Nil(t, output.Input)
	assert.Nil(t, output.Body)
	assert.Nil(t, output.Headers)
	assert.Empty(t, output.Status)
	assert.Empty(t, output.StatusCode)

	// with default values
	output = &oss.OperationOutput{
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
	output = &oss.OperationOutput{
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
	body := "{\"BodyConfiguration\":{\"StrField1\":\"StrField1\",\"StrField2\":\"StrField2\"}}"
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       oss.ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/json"},
		},
	}
	jsonresult := &jsonBodyResult{}
	err = c.unmarshalOutput(jsonresult, output, unmarshalBodyLikeXmlJson2)

	assert.Nil(t, err)
	assert.Equal(t, 200, jsonresult.StatusCode)
	assert.Equal(t, "OK", jsonresult.Status)
	assert.Equal(t, "StrField1", *jsonresult.StrField1)
	assert.Equal(t, "StrField2", *jsonresult.StrField2)
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

	oss.ResultCommon
}

func TestUnmarshalOutput_error(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	/*
		// unsupport content-type
		body := "{\"BodyConfiguration\":{\"StrField1\":\"StrField1\",\"StrField2\":\"StrField2\"}}"
		output = &oss.OperationOutput{
			StatusCode: 200,
			Status:     "OK",
			Body:       oss.ReadSeekNopCloser(bytes.NewReader([]byte(body))),
			Headers: http.Header{
				"Content-Type": {"application/text"},
			},
		}
		jsonresult := &jsonBodyResult{}
		err = c.unmarshalOutput(jsonresult, output, unmarshalBodyLikeXmlJson2)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unsupport contentType:application/text")
	*/
	// xml decode fail
	body := "StrField1>StrField1</StrField1><StrField2>StrField2<"
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       oss.ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/json"},
		},
	}
	result := &stubResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyLikeXmlJson2)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "eserialization failed,")
}

func TestResolveEndpoint(t *testing.T) {
	cfg := oss.NewConfig()
	assert.Nil(t, cfg.Region)
	assert.Nil(t, cfg.Endpoint)

	//No Endpont & No Region
	updateEndpoint(cfg)
	assert.Nil(t, cfg.Endpoint)

	//Endpont + ssl
	cfg = oss.NewConfig()
	cfg.Endpoint = oss.Ptr("test-endpoint")
	updateEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "test-endpoint", oss.ToString(cfg.Endpoint))

	cfg = oss.NewConfig()
	cfg.Endpoint = oss.Ptr("test-endpoint")
	cfg.DisableSSL = oss.Ptr(true)
	updateEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "test-endpoint", oss.ToString(cfg.Endpoint))

	cfg = oss.NewConfig()
	cfg.Endpoint = oss.Ptr("test-endpoint")
	cfg.DisableSSL = oss.Ptr(false)
	updateEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "test-endpoint", oss.ToString(cfg.Endpoint))

	cfg = oss.NewConfig()
	cfg.Endpoint = oss.Ptr("http://test-endpoint")
	updateEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "http://test-endpoint", oss.ToString(cfg.Endpoint))

	//Region + ssl
	cfg = oss.NewConfig()
	cfg.Region = oss.Ptr("test-region")
	updateEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "oss-test-region.oss-vectors.aliyuncs.com", oss.ToString(cfg.Endpoint))

	cfg = oss.NewConfig()
	cfg.Region = oss.Ptr("test-region")
	cfg.DisableSSL = oss.Ptr(true)
	updateEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "oss-test-region.oss-vectors.aliyuncs.com", oss.ToString(cfg.Endpoint))

	cfg = oss.NewConfig()
	cfg.Region = oss.Ptr("test-region")
	cfg.DisableSSL = oss.Ptr(false)
	updateEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "oss-test-region.oss-vectors.aliyuncs.com", oss.ToString(cfg.Endpoint))

	cfg = oss.NewConfig()
	cfg.Region = oss.Ptr("test-region")
	cfg.UseInternalEndpoint = oss.Ptr(true)
	updateEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "oss-test-region-internal.oss-vectors.aliyuncs.com", oss.ToString(cfg.Endpoint))

	cfg = oss.NewConfig()
	cfg.Region = oss.Ptr("test-region")
	cfg.UseDualStackEndpoint = oss.Ptr(true)
	cfg.UseAccelerateEndpoint = oss.Ptr(true)
	updateEndpoint(cfg)
	assert.NotNil(t, cfg.Endpoint)
	assert.Equal(t, "oss-test-region.oss-vectors.aliyuncs.com", oss.ToString(cfg.Endpoint))
}

func TestResolveVectorsUserAgent(t *testing.T) {
	cfg := oss.NewConfig()
	updateUserAgent(cfg)
	assert.Equal(t, oss.ToString(cfg.UserAgent), "vectors-client")

	cfg = oss.NewConfig()
	cfg.WithUserAgent("my-user-agent")
	updateUserAgent(cfg)
	assert.Equal(t, oss.ToString(cfg.UserAgent), "vectors-client/my-user-agent")
}
