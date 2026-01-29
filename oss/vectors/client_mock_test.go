package vectors

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func sortQuery(r *http.Request) string {
	u := r.URL
	var buf strings.Builder
	keys := make([]string, 0, len(u.Query()))
	for k := range u.Query() {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := u.Query()[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			if len(v) > 0 {
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(v))
			}
		}
	}
	u.RawQuery = buf.String()
	return u.String()
}

func testSetupMockServer(t *testing.T, statusCode int, headers map[string]string, body []byte,
	chkfunc func(t *testing.T, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check request
		chkfunc(t, r)

		// header s
		for k, v := range headers {
			w.Header().Set(k, v)
		}

		// status code
		w.WriteHeader(statusCode)

		// body
		w.Write(body)
	}))
}

var testVectorsInvokeOperationAnonymousCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Input          *oss.OperationInput
	CheckOutputFn  func(t *testing.T, o *oss.OperationOutput)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "5374A2880232A65C2300****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"encryptionConfiguration\": {\"KMSMasterKeyID\": \"string\",\"SSEAlgorithm\": \"string\"},\"vectorBucketName\": \"string\"}", string(requestBody))

		},
		&oss.OperationInput{
			OpName: "PutVectorBucket",
			Method: "PUT",
			Bucket: oss.Ptr("bucket"),
			Body:   strings.NewReader(`{"encryptionConfiguration": {"KMSMasterKeyID": "string","SSEAlgorithm": "string"},"vectorBucketName": "string"}`),
		},
		func(t *testing.T, o *oss.OperationOutput) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5374A2880232A65C2300****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 15 May 2014 11:18:32 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "5374A2880232A65C2300****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
					  "BucketInfo": {
						"Bucket": {
						  "CreationDate": "2013-07-31T10:56:21.000Z",
						  "ExtranetEndpoint": "oss-cn-hangzhou.aliyuncs.com",
						  "IntranetEndpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
						  "Location": "oss-cn-hangzhou",
						  "Name": "oss-example",
						  "ResourceGroupId": "rg-aek27tc********",
						}
					  }
					}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?bucketInfo", r.URL.String())
			assert.Equal(t, "GET", r.Method)

		},
		&oss.OperationInput{
			OpName: "GetVectorBucket",
			Bucket: oss.Ptr("bucket"),
			Method: "GET",
			Parameters: map[string]string{
				"bucketInfo": "",
			},
		},
		func(t *testing.T, o *oss.OperationOutput) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5374A2880232A65C2300****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 15 May 2014 11:18:32 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			content, err := io.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(content), "{\n\t\t\t\t\t  \"BucketInfo\": {\n\t\t\t\t\t\t\"Bucket\": {\n\t\t\t\t\t\t  \"CreationDate\": \"2013-07-31T10:56:21.000Z\",\n\t\t\t\t\t\t  \"ExtranetEndpoint\": \"oss-cn-hangzhou.aliyuncs.com\",\n\t\t\t\t\t\t  \"IntranetEndpoint\": \"oss-cn-hangzhou-internal.aliyuncs.com\",\n\t\t\t\t\t\t  \"Location\": \"oss-cn-hangzhou\",\n\t\t\t\t\t\t  \"Name\": \"oss-example\",\n\t\t\t\t\t\t  \"ResourceGroupId\": \"rg-aek27tc********\",\n\t\t\t\t\t\t}\n\t\t\t\t\t  }\n\t\t\t\t\t}")
		},
	},
	{
		204,
		map[string]string{
			"x-oss-request-id": "5374A2880232A65C2300****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?deleteVectorIndex", r.URL.String())
			assert.Equal(t, "DELETE", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"indexName\": \"string\"}", string(requestBody))
		},
		&oss.OperationInput{
			OpName: "DeleteVectorIndex",
			Bucket: oss.Ptr("bucket"),
			Method: "DELETE",
			Parameters: map[string]string{
				"deleteVectorIndex": "",
			},
			Body: strings.NewReader(`{"indexName": "string"}`),
		},
		func(t *testing.T, o *oss.OperationOutput) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "5374A2880232A65C2300****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 15 May 2014 11:18:32 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
		},
	},
}

func TestVectorsInvokeOperation_Anonymous(t *testing.T) {
	for _, c := range testVectorsInvokeOperationAnonymousCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output)
	}
}

var testVectorsInvokeOperationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Input          *oss.OperationInput
	CheckOutputFn  func(t *testing.T, o *oss.OperationOutput, err error)
}{
	{
		400,
		map[string]string{
			"x-oss-request-id": "57ABD896CCB80C366955****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(
			`{
  "Error": {
    "Code": "MissingArgument",
    "Message": "Missing Some Required Arguments.",
    "RequestId": "57ABD896CCB80C366955****",
    "HostId": "oss-example.oss-cn-hangzhou.aliyuncs.com",
    "EC": "0016-00000502",
    "RecommendDoc": "https://api.aliyun.com/troubleshoot?q=0016-00000502"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		&oss.OperationInput{
			OpName: "PutVectorBucket",
			Method: "PUT",
		},
		func(t *testing.T, o *oss.OperationOutput, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "MissingArgument", serr.Code)
			assert.Equal(t, "0016-00000502", serr.EC)
			assert.Equal(t, "57ABD896CCB80C366955****", serr.RequestID)
			assert.Contains(t, serr.Message, "Missing Some Required Arguments.")
		},
	},
}

func TestVectorsInvokeOperation_Error(t *testing.T) {
	for _, c := range testVectorsInvokeOperationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutVectorBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutVectorBucketRequest
	CheckOutputFn  func(t *testing.T, o *PutVectorBucketResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/", r.URL.String())
		},
		&PutVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *PutVectorBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-resource-group-id"), "rg-aek27tc****")
			assert.Equal(t, r.Header.Get("x-oss-bucket-tagging"), "k1=v1&k2=v2")
		},
		&PutVectorBucketRequest{
			Bucket:          oss.Ptr("bucket"),
			ResourceGroupId: oss.Ptr("rg-aek27tc****"),
			Tagging:         oss.Ptr("k1=v1&k2=v2"),
		},
		func(t *testing.T, o *PutVectorBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutVectorBucket_Success(t *testing.T) {
	for _, c := range testMockPutVectorBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutVectorBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutVectorBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutVectorBucketRequest
	CheckOutputFn  func(t *testing.T, o *PutVectorBucketResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(
			`{
			  "Error": {
				"Code": "SignatureDoesNotMatch",
				"Message": "The request signature we calculated does not match the signature you provided. Check your key and signing method.",
				"RequestId": "65467C42E001B4333337****",
				"SignatureProvided": "RizTbeKC/QlwxINq8xEdUPowc84=",
				"EC": "0002-00000040"
			  }
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/", r.URL.String())
		},
		&PutVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *PutVectorBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "SignatureDoesNotMatch", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
			assert.Contains(t, serr.RequestTarget, "/bucket")
		},
	},
	{
		409,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(
			`{
			  "Error": {
				"Code": "BucketAlreadyExists",
				"Message": "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.",
				"RequestId": "6548A043CA31D****",
				"EC": "0015-00000104"
			  }
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/", r.URL.String())
		},
		&PutVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *PutVectorBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "BucketAlreadyExists", serr.Code)
			assert.Equal(t, "0015-00000104", serr.EC)
			assert.Equal(t, "6548A043CA31D****", serr.RequestID)
			assert.Contains(t, serr.Message, "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again")
			assert.Contains(t, serr.RequestTarget, "/bucket")
		},
	},
}

func TestMockPutVectorBucket_Error(t *testing.T) {
	for _, c := range testMockPutVectorBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutVectorBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetVectorBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetVectorBucketRequest
	CheckOutputFn  func(t *testing.T, o *GetVectorBucketResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
					  "BucketInfo": {
						  "CreationDate": "2013-07-31T10:56:21.000Z",
						  "ExtranetEndpoint": "oss-cn-hangzhou.aliyuncs.com",
						  "IntranetEndpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
						  "Location": "oss-cn-hangzhou",
						  "Name": "oss-example",
						  "ResourceGroupId": "rg-aek27tc********"
					  }
					}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?bucketInfo", r.URL.String())
		},
		&GetVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetVectorBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.BucketInfo.Name, "oss-example")
			assert.Equal(t, *o.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.Location, "oss-cn-hangzhou")
			assert.Equal(t, *o.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
			assert.Equal(t, *o.BucketInfo.ResourceGroupId, "rg-aek27tc********")
		},
	},
}

func TestMockGetVectorBucket_Success(t *testing.T) {
	for _, c := range testMockGetVectorBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetVectorBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetVectorBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetVectorBucketRequest
	CheckOutputFn  func(t *testing.T, o *GetVectorBucketResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?bucketInfo", r.URL.String())
		},
		&GetVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetVectorBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?bucketInfo", strUrl)
		},
		&GetVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetVectorBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?bucketInfo", strUrl)
		},
		&GetVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetVectorBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetVectorBucket fail")
		},
	},
}

func TestMockGetVectorBucket_Error(t *testing.T) {
	for _, c := range testMockGetVectorBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetVectorBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListVectorBucketsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListVectorBucketsRequest
	CheckOutputFn  func(t *testing.T, o *ListVectorBucketsResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "ListAllMyBucketsResult": {
      "Buckets": [
        {
          "CreationDate": "2014-02-17T18:12:43.000Z",
          "ExtranetEndpoint": "oss-cn-shanghai.aliyuncs.com",
          "IntranetEndpoint": "oss-cn-shanghai-internal.aliyuncs.com",
          "Location": "oss-cn-shanghai",
          "Name": "app-base-oss",
          "Region": "cn-shanghai",
          "ResourceGroupId": "rg-aek27ta********"
        },
        {
          "CreationDate": "2014-02-25T11:21:04.000Z",
          "ExtranetEndpoint": "oss-cn-hangzhou.aliyuncs.com",
          "IntranetEndpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
          "Location": "oss-cn-hangzhou",
          "Name": "mybucket",
          "Region": "cn-hangzhou",
          "ResourceGroupId": "rg-aek27tc********"
        }
      ]
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		nil,
		func(t *testing.T, o *ListVectorBucketsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Buckets), 2)
			assert.Equal(t, *o.Buckets[0].CreationDate, time.Date(2014, time.February, 17, 18, 12, 43, 0, time.UTC))
			assert.Equal(t, *o.Buckets[0].ExtranetEndpoint, "oss-cn-shanghai.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].IntranetEndpoint, "oss-cn-shanghai-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].Name, "app-base-oss")
			assert.Equal(t, *o.Buckets[0].Region, "cn-shanghai")
			assert.Equal(t, *o.Buckets[0].Location, "oss-cn-shanghai")
			assert.Equal(t, *o.Buckets[0].ResourceGroupId, "rg-aek27ta********")

			assert.Equal(t, *o.Buckets[1].CreationDate, time.Date(2014, time.February, 25, 11, 21, 04, 0, time.UTC))
			assert.Equal(t, *o.Buckets[1].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.Buckets[1].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets[1].Name, "mybucket")
			assert.Equal(t, *o.Buckets[1].Region, "cn-hangzhou")
			assert.Equal(t, *o.Buckets[1].Location, "oss-cn-hangzhou")
			assert.Equal(t, *o.Buckets[1].ResourceGroupId, "rg-aek27tc********")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "ListAllMyBucketsResult": {
    "Prefix": "my",
    "Marker": "mybucket",
    "MaxKeys": 10,
    "IsTruncated": true,
    "NextMarker": "mybucket10",
      "Buckets": [{
        "CreationDate": "2014-05-14T11:18:32.000Z",
        "ExtranetEndpoint": "oss-cn-hangzhou.aliyuncs.com",
        "IntranetEndpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
        "Location": "oss-cn-hangzhou",
        "Name": "mybucket01",
        "Region": "cn-hangzhou",
        "ResourceGroupId": "rg-aek27tc********"
      }]
  }
}`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/?marker&max-keys=10&prefix=%2F", strUrl)
			assert.Equal(t, "rg-aek27tc********", r.Header.Get("x-oss-resource-group-id"))
		},
		&ListVectorBucketsRequest{
			Marker:          oss.Ptr(""),
			MaxKeys:         10,
			Prefix:          oss.Ptr("/"),
			ResourceGroupId: oss.Ptr("rg-aek27tc********"),
		},
		func(t *testing.T, o *ListVectorBucketsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Prefix, "my")
			assert.Equal(t, *o.Marker, "mybucket")
			assert.Equal(t, o.MaxKeys, int32(10))
			assert.Equal(t, o.IsTruncated, true)
			assert.Equal(t, *o.NextMarker, "mybucket10")

			assert.Equal(t, len(o.Buckets), 1)
			assert.Equal(t, *o.Buckets[0].CreationDate, time.Date(2014, time.May, 14, 11, 18, 32, 0, time.UTC))
			assert.Equal(t, *o.Buckets[0].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].Name, "mybucket01")
			assert.Equal(t, *o.Buckets[0].Region, "cn-hangzhou")
			assert.Equal(t, *o.Buckets[0].Location, "oss-cn-hangzhou")
			assert.Equal(t, *o.Buckets[0].ResourceGroupId, "rg-aek27tc********")
		},
	},
}

func TestMockListVectorBuckets_Success(t *testing.T) {
	for _, c := range testMockListVectorBucketsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListVectorBuckets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListVectorBucketsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListVectorBucketsRequest
	CheckOutputFn  func(t *testing.T, o *ListVectorBucketsResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(
			`{
  "Error": {
    "Code": "InvalidAccessKeyId",
    "Message": "The OSS Access Key Id you provided does not exist in our records.",
    "RequestId": "65467C42E001B4333337****",
    "SignatureProvided": "RizTbeKC/QlwxINq8xEdUPowc84=",
    "EC": "0002-00000040"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		&ListVectorBucketsRequest{},
		func(t *testing.T, o *ListVectorBucketsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "InvalidAccessKeyId", serr.Code)
			assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		&ListVectorBucketsRequest{},
		func(t *testing.T, o *ListVectorBucketsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListVectorBuckets fail")
		},
	},
}

func TestMockListVectorBuckets_Error(t *testing.T) {
	for _, c := range testMockListVectorBucketsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListVectorBuckets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteVectorBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteVectorBucketRequest
	CheckOutputFn  func(t *testing.T, o *DeleteVectorBucketResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/", r.URL.String())
		},
		&DeleteVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteVectorBucketResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteVectorBucket_Success(t *testing.T) {
	for _, c := range testMockDeleteVectorBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteVectorBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteVectorBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteVectorBucketRequest
	CheckOutputFn  func(t *testing.T, o *DeleteVectorBucketResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/", r.URL.String())
		},
		&DeleteVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteVectorBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		409,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "BucketNotEmpty",
    "Message": "The bucket has objects. Please delete them first.",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000301"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/", r.URL.String())
		},
		&DeleteVectorBucketRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteVectorBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "BucketNotEmpty", serr.Code)
			assert.Equal(t, "The bucket has objects. Please delete them first.", serr.Message)
			assert.Equal(t, "0015-00000301", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteVectorBucket_Error(t *testing.T) {
	for _, c := range testMockDeleteVectorBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteVectorBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketPolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketPolicyResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?policy", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"ossvector:PutVectors\",\"ossvector:GetVectors\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:ossvector:*:1234567890:*\"]}]}")
		},
		&PutBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
			Body:   strings.NewReader(`{"Version":"1","Statement":[{"Action":["ossvector:PutVectors","ossvector:GetVectors"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:ossvector:*:1234567890:*"]}]}`),
		},
		func(t *testing.T, o *PutBucketPolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketPolicy_Success(t *testing.T) {
	for _, c := range testMockPutBucketPolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketPolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketPolicyResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?policy", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"ossvector:PutVectors\",\"ossvector:GetVectors\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:ossvector:*:1234567890:*\"]}]}")
		},
		&PutBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
			Body:   strings.NewReader(`{"Version":"1","Statement":[{"Action":["ossvector:PutVectors","ossvector:GetVectors"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:ossvector:*:1234567890:*"]}]}`),
		},
		func(t *testing.T, o *PutBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"ossvector:PutVectors\",\"ossvector:GetVectors\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:ossvector:*:1234567890:*\"]}]}")
		},
		&PutBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
			Body:   strings.NewReader(`{"Version":"1","Statement":[{"Action":["ossvector:PutVectors","ossvector:GetVectors"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:ossvector:*:1234567890:*"]}]}`),
		},
		func(t *testing.T, o *PutBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutBucketPolicy_Error(t *testing.T) {
	for _, c := range testMockPutBucketPolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketPolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketPolicyResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"Version":"1","Statement":[{"Action":["ossvector:PutVectors","ossvector:GetVectors"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:ossvector:*:1234567890:*"]}]}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?policy", r.URL.String())
		},
		&GetBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketPolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Body, "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"ossvector:PutVectors\",\"ossvector:GetVectors\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:ossvector:*:1234567890:*\"]}]}")
		},
	},
}

func TestMockGetBucketPolicy_Success(t *testing.T) {
	for _, c := range testMockGetBucketPolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketPolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketPolicyResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?policy", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy", strUrl)
		},
		&GetBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetBucketPolicy_Error(t *testing.T) {
	for _, c := range testMockGetBucketPolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketPolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketPolicyResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy", strUrl)
		},
		&DeleteBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketPolicyResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteBucketPolicy_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketPolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketPolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketPolicyResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy", strUrl)
		},
		&DeleteBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy", strUrl)
		},
		&DeleteBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteBucketPolicy_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketPolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutVectorIndexSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutVectorIndexRequest
	CheckOutputFn  func(t *testing.T, o *PutVectorIndexResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?putVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\",\"metadata\":{\"nonFilterableMetadataKeys\":[\"foo\",\"bar\"]}}")
		},
		&PutVectorIndexRequest{
			Bucket:         oss.Ptr("bucket"),
			DataType:       oss.Ptr("string"),
			Dimension:      oss.Ptr(128),
			DistanceMetric: oss.Ptr("cosine"),
			IndexName:      oss.Ptr("exampleIndex"),
			Metadata: map[string]any{
				"nonFilterableMetadataKeys": []string{"foo", "bar"},
			},
		},
		func(t *testing.T, o *PutVectorIndexResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutVectorIndex_Success(t *testing.T) {
	for _, c := range testMockPutVectorIndexSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutVectorIndex(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutVectorIndexErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutVectorIndexRequest
	CheckOutputFn  func(t *testing.T, o *PutVectorIndexResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?putVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\",\"metadata\":{\"nonFilterableMetadataKeys\":[\"foo\",\"bar\"]}}")
		},
		&PutVectorIndexRequest{
			Bucket:         oss.Ptr("bucket"),
			DataType:       oss.Ptr("string"),
			Dimension:      oss.Ptr(128),
			DistanceMetric: oss.Ptr("cosine"),
			IndexName:      oss.Ptr("exampleIndex"),
			Metadata: map[string]any{
				"nonFilterableMetadataKeys": []string{"foo", "bar"},
			},
		},
		func(t *testing.T, o *PutVectorIndexResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?putVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\",\"metadata\":{\"nonFilterableMetadataKeys\":[\"foo\",\"bar\"]}}")
		},
		&PutVectorIndexRequest{
			Bucket:         oss.Ptr("bucket"),
			DataType:       oss.Ptr("string"),
			Dimension:      oss.Ptr(128),
			DistanceMetric: oss.Ptr("cosine"),
			IndexName:      oss.Ptr("exampleIndex"),
			Metadata: map[string]any{
				"nonFilterableMetadataKeys": []string{"foo", "bar"},
			},
		},
		func(t *testing.T, o *PutVectorIndexResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutVectorIndex_Error(t *testing.T) {
	for _, c := range testMockPutVectorIndexErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutVectorIndex(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetVectorIndexSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetVectorIndexRequest
	CheckOutputFn  func(t *testing.T, o *GetVectorIndexResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
   "index": { 
      "createTime": "2025-08-02T10:49:17.289372919Z",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "cosine",
      "indexName": "exampleIndex",
      "metadata": { 
         "nonFilterableMetadataKeys": ["foo", "bar"]
      },
      "status": "running",
      "vectorBucketName": "bucket"
   }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?getVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"exampleIndex\"}")
		},
		&GetVectorIndexRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("exampleIndex"),
		},
		func(t *testing.T, o *GetVectorIndexResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Index.CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.UTC))
			assert.Equal(t, *o.Index.DataType, "string")
			assert.Equal(t, *o.Index.Dimension, 128)
			assert.Equal(t, *o.Index.DistanceMetric, "cosine")
			assert.Equal(t, *o.Index.IndexName, "exampleIndex")
			//assert.Equal(t, len(o.Index.Metadata.NonFilterableMetadataKeys), 2)
			//assert.Equal(t, o.Index.Metadata.NonFilterableMetadataKeys[0], "foo")
			//assert.Equal(t, o.Index.Metadata.NonFilterableMetadataKeys[1], "bar")
			assert.Equal(t, *o.Index.Status, "running")
			assert.Equal(t, *o.Index.VectorBucketName, "bucket")
		},
	},
}

func TestMockGetVectorIndex_Success(t *testing.T) {
	for _, c := range testMockGetVectorIndexSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetVectorIndex(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetVectorIndexErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetVectorIndexRequest
	CheckOutputFn  func(t *testing.T, o *GetVectorIndexResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?getVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"demoIndex\"}")
		},
		&GetVectorIndexRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("demoIndex"),
		},
		func(t *testing.T, o *GetVectorIndexResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?getVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"demoIndex\"}")
		},
		&GetVectorIndexRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("demoIndex"),
		},
		func(t *testing.T, o *GetVectorIndexResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetVectorIndex_Error(t *testing.T) {
	for _, c := range testMockGetVectorIndexErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetVectorIndex(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListVectorIndexesSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListVectorIndexesRequest
	CheckOutputFn  func(t *testing.T, o *ListVectorIndexesResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
  "indexes": [
    { 
      "createTime": "2025-08-02T10:49:17.289372919Z",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "demo1",
      "metadata": { 
        "nonFilterableMetadataKeys": ["foo", "bar"]
      },
      "status": "running",
      "vectorBucketName": "bucket"
    },
    { 
      "createTime": "2025-08-20T10:49:17.289372919Z",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "demo2",
      "metadata": { 
        "nonFilterableMetadataKeys": ["foo2", "bar2"]
      },
      "status": "deleting",
      "vectorBucketName": "bucket"
    }
  ],
  "nextToken": "123"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?listVectorIndexes", r.URL.String())
		},
		&ListVectorIndexesRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListVectorIndexesResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Indexes), 2)
			assert.Equal(t, *o.Indexes[0].CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.UTC))
			assert.Equal(t, *o.Indexes[0].DataType, "string")
			assert.Equal(t, *o.Indexes[0].Dimension, 128)
			assert.Equal(t, *o.Indexes[0].DistanceMetric, "string")
			assert.Equal(t, *o.Indexes[0].IndexName, "demo1")
			assert.Len(t, o.Indexes[0].Metadata["nonFilterableMetadataKeys"], 2)
			if metadataValue, ok := o.Indexes[0].Metadata["nonFilterableMetadataKeys"]; ok {
				if keys, ok := metadataValue.([]any); ok {
					assert.Equal(t, keys[0].(string), "foo")
					assert.Equal(t, keys[1].(string), "bar")
				}
			}
			assert.Equal(t, *o.Indexes[0].Status, "running")
			assert.Equal(t, *o.Indexes[0].VectorBucketName, "bucket")

			assert.Equal(t, *o.Indexes[1].CreateTime, time.Date(2025, time.August, 20, 10, 49, 17, 289372919, time.UTC))
			assert.Equal(t, *o.Indexes[1].DataType, "string")
			assert.Equal(t, *o.Indexes[1].Dimension, 128)
			assert.Equal(t, *o.Indexes[1].DistanceMetric, "string")
			assert.Equal(t, *o.Indexes[1].IndexName, "demo2")
			assert.Len(t, o.Indexes[1].Metadata["nonFilterableMetadataKeys"], 2)
			if metadataValue, ok := o.Indexes[1].Metadata["nonFilterableMetadataKeys"]; ok {
				if keys, ok := metadataValue.([]any); ok {
					assert.Equal(t, keys[0].(string), "foo2")
					assert.Equal(t, keys[1].(string), "bar2")
				}
			}
			assert.Equal(t, *o.Indexes[1].VectorBucketName, "bucket")
			assert.Equal(t, *o.Indexes[1].Status, "deleting")
		},
	},
}

func TestMockListVectorIndexes_Success(t *testing.T) {
	for _, c := range testMockListVectorIndexesSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListVectorIndexes(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListVectorIndexesErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListVectorIndexesRequest
	CheckOutputFn  func(t *testing.T, o *ListVectorIndexesResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
 "Error": {
   "Code": "NoSuchBucket",
   "Message": "The specified bucket does not exist.",
   "RequestId": "5C3D9175B6FC201293AD****",
   "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
   "BucketName": "test",
   "EC": "0015-00000101"
 }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?listVectorIndexes", r.URL.String())
		},
		&ListVectorIndexesRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListVectorIndexesResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
 "Error": {
   "Code": "UserDisable",
   "Message": "UserDisable",
   "RequestId": "5C3D8D2A0ACA54D87B43****",
   "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
   "BucketName": "test",
   "EC": "0003-00000801"
 }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?listVectorIndexes", r.URL.String())
		},
		&ListVectorIndexesRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListVectorIndexesResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockListVectorIndexes_Error(t *testing.T) {
	for _, c := range testMockListVectorIndexesErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListVectorIndexes(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteVectorIndexSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteVectorIndexRequest
	CheckOutputFn  func(t *testing.T, o *DeleteVectorIndexResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?deleteVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"exampleIndex\"}")
		},
		&DeleteVectorIndexRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("exampleIndex"),
		},
		func(t *testing.T, o *DeleteVectorIndexResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteVectorIndex_Success(t *testing.T) {
	for _, c := range testMockDeleteVectorIndexSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteVectorIndex(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteVectorIndexErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteVectorIndexRequest
	CheckOutputFn  func(t *testing.T, o *DeleteVectorIndexResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?deleteVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"demoIndex\"}")
		},
		&DeleteVectorIndexRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("demoIndex"),
		},
		func(t *testing.T, o *DeleteVectorIndexResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?deleteVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"demoIndex\"}")
		},
		&DeleteVectorIndexRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("demoIndex"),
		},
		func(t *testing.T, o *DeleteVectorIndexResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteVectorIndex_Error(t *testing.T) {
	for _, c := range testMockDeleteVectorIndexErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteVectorIndex(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutVectorsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutVectorsRequest
	CheckOutputFn  func(t *testing.T, o *PutVectorsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?putVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `{"indexName":"exampleIndex","vectors":[{"data":{"float32":[1.2,2.5,3]},"key":"vector1","metadata":{"Key1":"value2","Key2":["1","2","3"]}}]}`)
		},
		&PutVectorsRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("exampleIndex"),
			Vectors: []map[string]any{
				{
					"key": "vector1",
					"data": map[string]any{
						"float32": []float32{1.2, 2.5, 3},
					},
					"metadata": map[string]any{
						"Key1": "value2",
						"Key2": []string{"1", "2", "3"},
					},
				},
			},
		},
		func(t *testing.T, o *PutVectorsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutVectors_Success(t *testing.T) {
	for _, c := range testMockPutVectorsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutVectorsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutVectorsRequest
	CheckOutputFn  func(t *testing.T, o *PutVectorsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?putVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `{"indexName":"exampleIndex","vectors":[{"data":{"float32":[1.2,2.5,3]},"key":"vector1","metadata":{"Key2":"value2","Key3":["1","2","3"]}}]}`)
		},
		&PutVectorsRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("exampleIndex"),
			Vectors: []map[string]any{
				{
					"key": "vector1",
					"data": map[string]any{
						"float32": []float32{1.2, 2.5, 3},
					},
					"metadata": map[string]any{
						"Key2": "value2",
						"Key3": []string{"1", "2", "3"},
					},
				},
			},
		},
		func(t *testing.T, o *PutVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?putVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `{"indexName":"exampleIndex","vectors":[{"data":{"float32":[1.2,2.5,3]},"key":"vector1","metadata":{"Key2":"value2","Key3":["1","2","3"]}}]}`)
		},
		&PutVectorsRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("exampleIndex"),
			Vectors: []map[string]any{
				{
					"key": "vector1",
					"data": map[string]any{
						"float32": []float32{1.2, 2.5, 3},
					},
					"metadata": map[string]any{
						"Key2": "value2",
						"Key3": []string{"1", "2", "3"},
					},
				},
			},
		},
		func(t *testing.T, o *PutVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutVectors_Error(t *testing.T) {
	for _, c := range testMockPutVectorsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetVectorsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetVectorsRequest
	CheckOutputFn  func(t *testing.T, o *GetVectorsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
   "indexName": "index",
   "vectors": [ 
      { 
         "data": {
            "float32":[32]
         },
         "key": "key",
         "metadata": {
             "Key1": "value1",
             "Key2": "value2"
         }
      }
   ]
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?getVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"index\",\"keys\":[\"key1\",\"key2\",\"key3\"],\"returnData\":true,\"returnMetadata\":false}")
		},
		&GetVectorsRequest{
			Bucket:         oss.Ptr("bucket"),
			IndexName:      oss.Ptr("index"),
			Keys:           []string{"key1", "key2", "key3"},
			ReturnData:     oss.Ptr(true),
			ReturnMetadata: oss.Ptr(false),
		},
		func(t *testing.T, o *GetVectorsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Vectors), 1)
			for _, vector := range o.Vectors {
				if keyVal, exists := vector["key"]; exists {
					keyStr, ok := keyVal.(string)
					assert.True(t, ok)
					assert.Equal(t, keyStr, "key")
				}

				// 访问 data 字段
				if dataVal, exists := vector["data"]; exists {
					dataMap, ok := dataVal.(map[string]any)
					assert.True(t, ok)
					if float32Val, exists := dataMap["float32"]; exists {
						float32Data, ok := float32Val.([]any)
						assert.True(t, ok)
						assert.Equal(t, float32Data[0], float64(32))
					}
				}

				if metadataVal, exists := vector["metadata"]; exists {
					metadataMap, ok := metadataVal.(map[string]any)
					assert.True(t, ok)
					if key1Val, exists := metadataMap["Key1"]; exists {
						key1Data, ok := key1Val.(string)
						assert.True(t, ok)
						assert.Equal(t, key1Data, "value1")
					}
					if key2Val, exists := metadataMap["Key2"]; exists {
						key2Data, ok := key2Val.(string)
						assert.True(t, ok)
						assert.Equal(t, key2Data, "value2")
					}
				}
			}
		},
	},
}

func TestMockGetVectors_Success(t *testing.T) {
	for _, c := range testMockGetVectorsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetVectorsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetVectorsRequest
	CheckOutputFn  func(t *testing.T, o *GetVectorsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?getVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"index\",\"keys\":[\"key1\",\"key2\",\"key3\"],\"returnData\":true,\"returnMetadata\":false}")
		},
		&GetVectorsRequest{
			Bucket:         oss.Ptr("bucket"),
			IndexName:      oss.Ptr("index"),
			Keys:           []string{"key1", "key2", "key3"},
			ReturnData:     oss.Ptr(true),
			ReturnMetadata: oss.Ptr(false),
		},
		func(t *testing.T, o *GetVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?getVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"index\",\"keys\":[\"key1\",\"key2\",\"key3\"],\"returnData\":true,\"returnMetadata\":false}")
		},
		&GetVectorsRequest{
			Bucket:         oss.Ptr("bucket"),
			IndexName:      oss.Ptr("index"),
			Keys:           []string{"key1", "key2", "key3"},
			ReturnData:     oss.Ptr(true),
			ReturnMetadata: oss.Ptr(false),
		},
		func(t *testing.T, o *GetVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetVectors_Error(t *testing.T) {
	for _, c := range testMockGetVectorsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListVectorsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListVectorsRequest
	CheckOutputFn  func(t *testing.T, o *ListVectorsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
    "nextToken": "123",
   "vectors": [ 
      { 
         "data": {
            "float32":[32]
         },
         "key": "key",
         "metadata": {
             "Key1": "value1",
             "Key2": "value2"
         }
      }
   ]
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?listVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"index\",\"maxResults\":100,\"nextToken\":\"123\",\"returnData\":false,\"returnMetadata\":true,\"segmentCount\":10,\"segmentIndex\":3}")
		},
		&ListVectorsRequest{
			Bucket:         oss.Ptr("bucket"),
			IndexName:      oss.Ptr("index"),
			MaxResults:     100,
			NextToken:      oss.Ptr("123"),
			ReturnMetadata: oss.Ptr(true),
			ReturnData:     oss.Ptr(false),
			SegmentCount:   oss.Ptr(int(10)),
			SegmentIndex:   oss.Ptr(3),
		},
		func(t *testing.T, o *ListVectorsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Vectors), 1)
			assert.Equal(t, *o.NextToken, "123")
			assert.NotEmpty(t, o.Vectors)
		},
	},
}

func TestMockListVectors_Success(t *testing.T) {
	for _, c := range testMockListVectorsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListVectorsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListVectorsRequest
	CheckOutputFn  func(t *testing.T, o *ListVectorsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?listVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"index\",\"maxResults\":100,\"nextToken\":\"123\",\"returnData\":false,\"returnMetadata\":true,\"segmentCount\":10,\"segmentIndex\":3}")
		},
		&ListVectorsRequest{
			Bucket:         oss.Ptr("bucket"),
			IndexName:      oss.Ptr("index"),
			MaxResults:     100,
			NextToken:      oss.Ptr("123"),
			ReturnMetadata: oss.Ptr(true),
			ReturnData:     oss.Ptr(false),
			SegmentCount:   oss.Ptr(int(10)),
			SegmentIndex:   oss.Ptr(3),
		},
		func(t *testing.T, o *ListVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?listVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"index\",\"maxResults\":100,\"nextToken\":\"123\",\"returnData\":false,\"returnMetadata\":true,\"segmentCount\":10,\"segmentIndex\":3}")
		},
		&ListVectorsRequest{
			Bucket:         oss.Ptr("bucket"),
			IndexName:      oss.Ptr("index"),
			MaxResults:     100,
			NextToken:      oss.Ptr("123"),
			ReturnMetadata: oss.Ptr(true),
			ReturnData:     oss.Ptr(false),
			SegmentCount:   oss.Ptr(int(10)),
			SegmentIndex:   oss.Ptr(3),
		},
		func(t *testing.T, o *ListVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockListVectors_Error(t *testing.T) {
	for _, c := range testMockListVectorsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteVectorsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteVectorsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteVectorsResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?deleteVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"index\",\"keys\":[\"key1\",\"key2\"]}")
		},
		&DeleteVectorsRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("index"),
			Keys: []string{
				"key1", "key2",
			},
		},
		func(t *testing.T, o *DeleteVectorsResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteVectors_Success(t *testing.T) {
	for _, c := range testMockDeleteVectorsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteVectorsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteVectorsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteVectorsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?deleteVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"index\",\"keys\":[\"key1\",\"key2\"]}")
		},
		&DeleteVectorsRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("index"),
			Keys: []string{
				"key1", "key2",
			},
		},
		func(t *testing.T, o *DeleteVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?deleteVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"index\",\"keys\":[\"key1\",\"key2\"]}")
		},
		&DeleteVectorsRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("index"),
			Keys: []string{
				"key1", "key2",
			},
		},
		func(t *testing.T, o *DeleteVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteVectors_Error(t *testing.T) {
	for _, c := range testMockDeleteVectorsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockQueryVectorsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *QueryVectorsRequest
	CheckOutputFn  func(t *testing.T, o *QueryVectorsResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?queryVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"filter\":{\"$and\":[{\"type\":{\"$in\":[\"comedy\",\"documentary\"]}}]},\"indexName\":\"index\",\"queryVector\":{\"float32\":[32]},\"returnDistance\":true,\"returnMetadata\":true,\"topK\":10}")
		},
		&QueryVectorsRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("index"),
			Filter: map[string]any{
				"$and": []map[string]any{
					{
						"type": map[string]any{
							"$in": []string{"comedy", "documentary"},
						},
					},
				},
			},
			QueryVector: map[string]any{
				"float32": []float32{float32(32)},
			},
			ReturnMetadata: oss.Ptr(true),
			ReturnDistance: oss.Ptr(true),
			TopK:           oss.Ptr(10),
		},
		func(t *testing.T, o *QueryVectorsResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockQueryVectors_Success(t *testing.T) {
	for _, c := range testMockQueryVectorsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.QueryVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockQueryVectorsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *QueryVectorsRequest
	CheckOutputFn  func(t *testing.T, o *QueryVectorsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?queryVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"filter\":{\"$and\":[{\"type\":{\"$in\":[\"comedy\",\"documentary\"]}}]},\"indexName\":\"index\",\"queryVector\":{\"float32\":[32]},\"returnDistance\":true,\"returnMetadata\":true,\"topK\":10}")
		},
		&QueryVectorsRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("index"),
			Filter: map[string]any{
				"$and": []map[string]any{
					{
						"type": map[string]any{
							"$in": []string{"comedy", "documentary"},
						},
					},
				},
			},
			QueryVector: map[string]any{
				"float32": []float32{float32(32)},
			},
			ReturnMetadata: oss.Ptr(true),
			ReturnDistance: oss.Ptr(true),
			TopK:           oss.Ptr(10),
		},
		func(t *testing.T, o *QueryVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?queryVectors", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"filter\":{\"$and\":[{\"type\":{\"$in\":[\"comedy\",\"documentary\"]}}]},\"indexName\":\"index\",\"queryVector\":{\"float32\":[32]},\"returnDistance\":true,\"returnMetadata\":true,\"topK\":10}")
		},
		&QueryVectorsRequest{
			Bucket:    oss.Ptr("bucket"),
			IndexName: oss.Ptr("index"),
			Filter: map[string]any{
				"$and": []map[string]any{
					{
						"type": map[string]any{
							"$in": []string{"comedy", "documentary"},
						},
					},
				},
			},
			QueryVector: map[string]any{
				"float32": []float32{float32(32)},
			},
			ReturnMetadata: oss.Ptr(true),
			ReturnDistance: oss.Ptr(true),
			TopK:           oss.Ptr(10),
		},
		func(t *testing.T, o *QueryVectorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockQueryVectors_Error(t *testing.T) {
	for _, c := range testMockQueryVectorsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.QueryVectors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketLoggingSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketLoggingResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?logging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"BucketLoggingStatus\":{\"LoggingEnabled\":{\"TargetBucket\":\"TargetBucket\",\"TargetPrefix\":\"TargetPrefix\",\"LoggingRole\":\"AliyunOSSLoggingDefaultRole\"}}}")
		},
		&PutBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				LoggingEnabled: &LoggingEnabled{
					TargetBucket: oss.Ptr("TargetBucket"),
					TargetPrefix: oss.Ptr("TargetPrefix"),
					LoggingRole:  oss.Ptr("AliyunOSSLoggingDefaultRole"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?logging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"BucketLoggingStatus\":{\"LoggingEnabled\":{\"TargetBucket\":\"TargetBucket\"}}}")
		},
		&PutBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				LoggingEnabled: &LoggingEnabled{
					TargetBucket: oss.Ptr("TargetBucket"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketLogging_Success(t *testing.T) {
	for _, c := range testMockPutBucketLoggingSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketLoggingErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketLoggingResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?logging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"BucketLoggingStatus\":{\"LoggingEnabled\":{\"TargetBucket\":\"TargetBucket\",\"TargetPrefix\":\"TargetPrefix\"}}}")
		},
		&PutBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				LoggingEnabled: &LoggingEnabled{
					TargetBucket: oss.Ptr("TargetBucket"),
					TargetPrefix: oss.Ptr("TargetPrefix"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"BucketLoggingStatus\":{\"LoggingEnabled\":{\"TargetBucket\":\"TargetBucket\",\"TargetPrefix\":\"TargetPrefix\"}}}")
		},
		&PutBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				LoggingEnabled: &LoggingEnabled{
					TargetBucket: oss.Ptr("TargetBucket"),
					TargetPrefix: oss.Ptr("TargetPrefix"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutBucketLogging_Error(t *testing.T) {
	for _, c := range testMockPutBucketLoggingErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketLoggingSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLoggingResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "BucketLoggingStatus": {
    "LoggingEnabled": {
      "TargetBucket": "bucket-log",
      "TargetPrefix": "prefix-access_log"
    }
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?logging", r.URL.String())
		},
		&GetBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.BucketLoggingStatus.LoggingEnabled.TargetBucket, "bucket-log")
			assert.Equal(t, *o.BucketLoggingStatus.LoggingEnabled.TargetPrefix, "prefix-access_log")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "BucketLoggingStatus": {
    "LoggingEnabled": {
      "TargetBucket": "bucket-log"
    }
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?logging", r.URL.String())
		},
		&GetBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.BucketLoggingStatus.LoggingEnabled.TargetBucket, "bucket-log")
			assert.Nil(t, o.BucketLoggingStatus.LoggingEnabled.TargetPrefix)
		},
	},
}

func TestMockGetBucketLogging_Success(t *testing.T) {
	for _, c := range testMockGetBucketLoggingSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketLoggingErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLoggingResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?logging", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&GetBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&GetBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketLogging fail")
		},
	},
}

func TestMockGetBucketLogging_Error(t *testing.T) {
	for _, c := range testMockGetBucketLoggingErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketLoggingSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketLoggingResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&DeleteBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLoggingResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&DeleteBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLoggingResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteBucketLogging_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketLoggingSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketLoggingErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketLoggingResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&DeleteBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLoggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "UserDisable",
    "Message": "UserDisable",
    "RequestId": "5C3D8D2A0ACA54D87B43****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0003-00000801"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&DeleteBucketLoggingRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLoggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteBucketLogging_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketLoggingErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
