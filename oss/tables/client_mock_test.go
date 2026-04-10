package tables

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"

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

func addIsBucketArn() oss.OperationMetadata {
	om := oss.OperationMetadata{}
	om.Add(oss.OpMetaKeyRequestIsBucketArn, true)
	return om
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

var testTablesInvokeOperationAnonymousCases = []struct {
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
			assert.Equal(t, "/buckets", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, contentTypeJSON, r.Header.Get(oss.HTTPHeaderContentType))
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"\",\"sseAlgorithm\":\"AES256\"},\"name\":\"bucket\"}", string(requestBody))

		},
		&oss.OperationInput{
			OpName: "CreateTableBucket",
			Method: "PUT",
			Headers: map[string]string{
				oss.HTTPHeaderContentType: contentTypeJSON,
			},
			Key:  oss.Ptr("buckets"),
			Body: strings.NewReader(`{"encryptionConfiguration":{"kmsKeyArn":"","sseAlgorithm":"AES256"},"name":"bucket"}`),
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
   "arn": "test-arn",
   "createdAt": "2013-07-31T10:56:21.000Z",
   "name": "oss-bucket",
   "ownerAccountId": "123456",
   "tableBucketId": "123",
   "type": "oss"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, contentTypeJSON, r.Header.Get(oss.HTTPHeaderContentType))
		},
		&oss.OperationInput{
			OpName: "GetTableBucket",
			Method: "GET",
			Headers: map[string]string{
				oss.HTTPHeaderContentType: contentTypeJSON,
			},
			Bucket:     oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Key:        oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"))),
			OpMetadata: addIsBucketArn(),
		},
		func(t *testing.T, o *oss.OperationOutput) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5374A2880232A65C2300****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 15 May 2014 11:18:32 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			content, err := io.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(content), "{\n   \"arn\": \"test-arn\",\n   \"createdAt\": \"2013-07-31T10:56:21.000Z\",\n   \"name\": \"oss-bucket\",\n   \"ownerAccountId\": \"123456\",\n   \"tableBucketId\": \"123\",\n   \"type\": \"oss\"\n}")
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
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "DELETE", r.Method)
		},
		&oss.OperationInput{
			OpName: "DeleteTableBucket",
			Method: "DELETE",
			Headers: map[string]string{
				oss.HTTPHeaderContentType: contentTypeJSON,
			},
			Bucket:     oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Key:        oss.Ptr(fmt.Sprintf("buckets/%s", url.QueryEscape("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"))),
			OpMetadata: addIsBucketArn(),
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

func TestTablesInvokeOperation_Anonymous(t *testing.T) {
	for _, c := range testTablesInvokeOperationAnonymousCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output)
	}
}

var testTablesInvokeOperationErrorCases = []struct {
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
			"x-oss-ec":         "0016-00000502",
		},
		[]byte(
			`{"message": "Missing Some Required Arguments."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, contentTypeJSON, r.Header.Get(oss.HTTPHeaderContentType))
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"name\":\"bucket\"}", string(requestBody))
		},
		&oss.OperationInput{
			OpName: "CreateTableBucket",
			Method: "PUT",
			Headers: map[string]string{
				oss.HTTPHeaderContentType: contentTypeJSON,
			},
			Bucket: oss.Ptr("bucket"),
			Body:   strings.NewReader(`{"name":"bucket"}`),
		},
		func(t *testing.T, o *oss.OperationOutput, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "Bad Request", serr.Code)
			assert.Equal(t, "0016-00000502", serr.EC)
			assert.Equal(t, "57ABD896CCB80C366955****", serr.RequestID)
			assert.Contains(t, serr.Message, "Missing Some Required Arguments.")
		},
	},
	{
		405,
		map[string]string{
			"x-oss-request-id": "57ABD896CCB80C366955****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0016-00000502",
		},
		[]byte(
			`{"message": "The specified method is not allowed against this resource."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, contentTypeJSON, r.Header.Get(oss.HTTPHeaderContentType))
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"name\":\"bucket\"}", string(requestBody))
		},
		&oss.OperationInput{
			OpName: "CreateTableBucket",
			Method: "PUT",
			Headers: map[string]string{
				oss.HTTPHeaderContentType: contentTypeJSON,
			},
			Bucket: oss.Ptr("bucket"),
			Body:   strings.NewReader(`{"name":"bucket"}`),
		},
		func(t *testing.T, o *oss.OperationOutput, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(405), serr.StatusCode)
			assert.Equal(t, "Method Not Allowed", serr.Code)
			assert.Equal(t, "0016-00000502", serr.EC)
			assert.Equal(t, "57ABD896CCB80C366955****", serr.RequestID)
			assert.Contains(t, serr.Message, "The specified method is not allowed against this resource.")
		},
	},
}

func TestTablesInvokeOperation_Error(t *testing.T) {
	for _, c := range testTablesInvokeOperationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateTableBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateTableBucketRequest
	CheckOutputFn  func(t *testing.T, o *CreateTableBucketResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{"arn": "acs:osstables:cn-beijing:1234567890:bucket/bucket"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/buckets", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"arn\",\"sseAlgorithm\":\"AES256\"},\"name\":\"bucket\"}", string(requestBody))
		},
		&CreateTableBucketRequest{
				Name: oss.Ptr("bucket"),
				EncryptionConfiguration: &EncryptionConfiguration{
				KmsKeyArn:    oss.Ptr("arn"),
				SseAlgorithm: oss.Ptr("AES256"),
			},
		},
		func(t *testing.T, o *CreateTableBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "acs:osstables:cn-beijing:1234567890:bucket/bucket", *o.Arn)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{"arn": "acs:osstables:cn-beijing:1234567890:bucket/bucket"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/buckets", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&CreateTableBucketRequest{
				Name: oss.Ptr("bucket"),
			},
		func(t *testing.T, o *CreateTableBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "acs:osstables:cn-beijing:1234567890:bucket/bucket", *o.Arn)
		},
	},
}

func TestMockCreateTableBucket_Success(t *testing.T) {
	for _, c := range testMockCreateTableBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateTableBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateTableBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateTableBucketRequest
	CheckOutputFn  func(t *testing.T, o *CreateTableBucketResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0002-00000040",
		},
		[]byte(
			`{
				"message": "The request signature we calculated does not match the signature you provided. Check your key and signing method."
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/buckets", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"name\":\"bucket\"}", string(requestBody))
		},
		&CreateTableBucketRequest{
				Name: oss.Ptr("bucket"),
			},
		func(t *testing.T, o *CreateTableBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
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
			"x-oss-ec":         "0015-00000104",
		},
		[]byte(
			`{
				"message": "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again."
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/buckets", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"name\":\"bucket\"}", string(requestBody))
		},
		&CreateTableBucketRequest{
				Name: oss.Ptr("bucket"),
			},
		func(t *testing.T, o *CreateTableBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "Conflict", serr.Code)
			assert.Equal(t, "0015-00000104", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again")
			assert.Contains(t, serr.RequestTarget, "/bucket")
		},
	},
}

func TestMockCreateTableBucket_Error(t *testing.T) {
	for _, c := range testMockCreateTableBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateTableBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableBucketRequest
	CheckOutputFn  func(t *testing.T, o *GetTableBucketResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
   "arn": "acs:osstables:cn-beijing:12345657890:bucket/demo-bucket",
   "createdAt": "2026-04-01T09:42:50.000000+00:00",
   "name": "demo-bucket",
   "ownerAccountId": "12345657890",
   "tableBucketId": "50859410-3482-401c-b500-605c22848ef4",
   "type": "oss"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetTableBucketRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, o.Headers.Get("Content-Type"), "application/json")
			assert.Equal(t, *o.Arn, "acs:osstables:cn-beijing:12345657890:bucket/demo-bucket")
			assert.Equal(t, *o.Name, "demo-bucket")
			assert.Equal(t, *o.CreatedAt, "2026-04-01T09:42:50.000000+00:00")
			assert.Equal(t, *o.OwnerAccountId, "12345657890")
			assert.Equal(t, *o.TableBucketId, "50859410-3482-401c-b500-605c22848ef4")
			assert.Equal(t, *o.Type, "oss")
		},
	},
}

func TestMockGetTableBucket_Success(t *testing.T) {
	for _, c := range testMockGetTableBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableBucketRequest
	CheckOutputFn  func(t *testing.T, o *GetTableBucketResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetTableBucketRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetTableBucketRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
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
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetTableBucketRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetTableBucket fail")
		},
	},
}

func TestMockGetTableBucket_Error(t *testing.T) {
	for _, c := range testMockGetTableBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListTableBucketsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListTableBucketsRequest
	CheckOutputFn  func(t *testing.T, o *ListTableBucketsResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "continuationToken": "Cj5hY3M6b3NzdGFibGVzOmNuLWJlaWppbmc6MTc2MDIyNTU0NTA4NDMzMTpidWNrZXQvZGVtby13YWxrZXItMQ--",
  "tableBuckets": [{
    "arn": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket",
    "createdAt": "2026-04-02T05:27:31.000000+00:00",
    "name": "demo-bucket",
    "ownerAccountId": "1234567890",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5",
    "type": "customer"
  },
  {
    "arn": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket-1",
    "createdAt": "2026-04-02T05:27:32.000000+00:00",
    "name": "demo-bucket-1",
    "ownerAccountId": "1234567890",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f3",
    "type": "customer"
  }]
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/buckets", r.URL.String())
		},
		nil,
		func(t *testing.T, o *ListTableBucketsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ContinuationToken, "Cj5hY3M6b3NzdGFibGVzOmNuLWJlaWppbmc6MTc2MDIyNTU0NTA4NDMzMTpidWNrZXQvZGVtby13YWxrZXItMQ--")
			assert.Equal(t, len(o.TableBuckets), 2)
			assert.Equal(t, *o.TableBuckets[0].CreatedAt, "2026-04-02T05:27:31.000000+00:00")
			assert.Equal(t, *o.TableBuckets[0].Arn, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
			assert.Equal(t, *o.TableBuckets[0].Name, "demo-bucket")
			assert.Equal(t, *o.TableBuckets[0].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")
			assert.Equal(t, *o.TableBuckets[0].OwnerAccountId, "1234567890")
			assert.Equal(t, *o.TableBuckets[0].Type, "customer")

			assert.Equal(t, *o.TableBuckets[1].CreatedAt, "2026-04-02T05:27:32.000000+00:00")
			assert.Equal(t, *o.TableBuckets[1].Arn, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket-1")
			assert.Equal(t, *o.TableBuckets[1].Name, "demo-bucket-1")
			assert.Equal(t, *o.TableBuckets[1].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f3")
			assert.Equal(t, *o.TableBuckets[1].OwnerAccountId, "1234567890")
			assert.Equal(t, *o.TableBuckets[1].Type, "customer")
		},
	},
}

func TestMockListTableBuckets_Success(t *testing.T) {
	for _, c := range testMockListTableBucketsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListTableBuckets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListTableBucketsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListTableBucketsRequest
	CheckOutputFn  func(t *testing.T, o *ListTableBucketsResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0002-00000040",
		},
		[]byte(
			`{"message": "The OSS Access Key Id you provided does not exist in our records."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/buckets", r.URL.String())
		},
		&ListTableBucketsRequest{},
		func(t *testing.T, o *ListTableBucketsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
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
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/buckets", r.URL.String())
		},
		&ListTableBucketsRequest{},
		func(t *testing.T, o *ListTableBucketsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListTableBuckets fail")
		},
	},
}

func TestMockListTableBuckets_Error(t *testing.T) {
	for _, c := range testMockListTableBucketsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListTableBuckets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTableBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTableBucketRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTableBucketResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "DELETE", r.Method)
		},
		&DeleteTableBucketRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteTableBucket_Success(t *testing.T) {
	for _, c := range testMockDeleteTableBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTableBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTableBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTableBucketRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTableBucketResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "DELETE", r.Method)
		},
		&DeleteTableBucketRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0015-00000301",
		},
		[]byte(`{"message": "The bucket has objects. Please delete them first."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "DELETE", r.Method)
		},
		&DeleteTableBucketRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "Conflict", serr.Code)
			assert.Equal(t, "The bucket has objects. Please delete them first.", serr.Message)
			assert.Equal(t, "0015-00000301", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteTableBucket_Error(t *testing.T) {
	for _, c := range testMockDeleteTableBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTableBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTableBucketEncryptionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTableBucketEncryptionRequest
	CheckOutputFn  func(t *testing.T, o *PutTableBucketEncryptionResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", r.URL.String())
			assert.Equal(t, contentTypeJSON, r.Header.Get(oss.HTTPHeaderContentType))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"\",\"sseAlgorithm\":\"AES256\"}}")
		},
		&PutTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			EncryptionConfiguration: &EncryptionConfiguration{
				KmsKeyArn:    oss.Ptr(""),
				SseAlgorithm: oss.Ptr("AES256"),
			},
		},
		func(t *testing.T, o *PutTableBucketEncryptionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutTableBucketEncryption_Success(t *testing.T) {
	for _, c := range testMockPutTableBucketEncryptionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutTableBucketEncryption(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTableBucketEncryptionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTableBucketEncryptionRequest
	CheckOutputFn  func(t *testing.T, o *PutTableBucketEncryptionResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", r.URL.String())
			assert.Equal(t, contentTypeJSON, r.Header.Get(oss.HTTPHeaderContentType))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"\",\"sseAlgorithm\":\"AES256\"}}")
		},
		&PutTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			EncryptionConfiguration: &EncryptionConfiguration{
				KmsKeyArn:    oss.Ptr(""),
				SseAlgorithm: oss.Ptr("AES256"),
			},
		},
		func(t *testing.T, o *PutTableBucketEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{
    "message": "UserDisable"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", r.URL.String())
			assert.Equal(t, contentTypeJSON, r.Header.Get(oss.HTTPHeaderContentType))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"\",\"sseAlgorithm\":\"AES256\"}}")
		},
		&PutTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			EncryptionConfiguration: &EncryptionConfiguration{
				KmsKeyArn:    oss.Ptr(""),
				SseAlgorithm: oss.Ptr("AES256"),
			},
		},
		func(t *testing.T, o *PutTableBucketEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutTableBucketEncryption_Error(t *testing.T) {
	for _, c := range testMockPutTableBucketEncryptionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutTableBucketEncryption(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableBucketEncryptionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableBucketEncryptionRequest
	CheckOutputFn  func(t *testing.T, o *GetTableBucketEncryptionResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
   "encryptionConfiguration": { 
      "kmsKeyArn": "",
      "sseAlgorithm": "AES256"
   }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketEncryptionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.EncryptionConfiguration.KmsKeyArn, "")
			assert.Equal(t, *o.EncryptionConfiguration.SseAlgorithm, "AES256")
		},
	},
}

func TestMockGetTableBucketEncryption_Success(t *testing.T) {
	for _, c := range testMockGetTableBucketEncryptionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableBucketEncryption(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableBucketEncryptionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableBucketEncryptionRequest
	CheckOutputFn  func(t *testing.T, o *GetTableBucketEncryptionResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", r.URL.String())
		},
		&GetTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
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
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", r.URL.String())
		},
		&GetTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetTableBucketEncryption fail")
		},
	},
}

func TestMockGetTableBucketEncryption_Error(t *testing.T) {
	for _, c := range testMockGetTableBucketEncryptionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableBucketEncryption(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTableBucketEncryptionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTableBucketEncryptionRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTableBucketEncryptionResult, err error)
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
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", strUrl)
		},
		&DeleteTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketEncryptionResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteTableBucketEncryption_Success(t *testing.T) {
	for _, c := range testMockDeleteTableBucketEncryptionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTableBucketEncryption(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTableBucketEncryptionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTableBucketEncryptionRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTableBucketEncryptionResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", strUrl)
		},
		&DeleteTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/encryption", strUrl)
		},
		&DeleteTableBucketEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteTableBucketEncryption_Error(t *testing.T) {
	for _, c := range testMockDeleteTableBucketEncryptionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTableBucketEncryption(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTableBucketPolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTableBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *PutTableBucketPolicyResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/policy", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"resourcePolicy\":\"{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:GetTable\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"1234567890\\\"],\\\"Resource\\\":[\\\"acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket\\\"]}]}\"}")
		},
		&PutTableBucketPolicyRequest{
			TableBucketARN:      oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			ResourcePolicy: oss.Ptr(`{"Version":"1","Statement":[{"Action":["oss:GetTable"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket"]}]}`),
		},
		func(t *testing.T, o *PutTableBucketPolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutTableBucketPolicy_Success(t *testing.T) {
	for _, c := range testMockPutTableBucketPolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutTableBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTableBucketPolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTableBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *PutTableBucketPolicyResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{
    "Message": "The specified bucket does not exist."
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/policy", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"resourcePolicy\":\"{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:GetTable\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"1234567890\\\"],\\\"Resource\\\":[\\\"acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket\\\"]}]}\"}")
		},
		&PutTableBucketPolicyRequest{
			TableBucketARN:      oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			ResourcePolicy: oss.Ptr(`{"Version":"1","Statement":[{"Action":["oss:GetTable"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket"]}]}`),
		},
		func(t *testing.T, o *PutTableBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/policy", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"resourcePolicy\":\"{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:GetTable\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"1234567890\\\"],\\\"Resource\\\":[\\\"acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket\\\"]}]}\"}")
		},
		&PutTableBucketPolicyRequest{
			TableBucketARN:      oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			ResourcePolicy: oss.Ptr(`{"Version":"1","Statement":[{"Action":["oss:GetTable"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket"]}]}`),
		},
		func(t *testing.T, o *PutTableBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutTableBucketPolicy_Error(t *testing.T) {
	for _, c := range testMockPutTableBucketPolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutTableBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableBucketPolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *GetTableBucketPolicyResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte("{\"resourcePolicy\":\"{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:GetTable\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"1234567890\\\"],\\\"Resource\\\":[\\\"acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket\\\"]}]}\"}"),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/policy", r.URL.String())
		},
		&GetTableBucketPolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketPolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ResourcePolicy, "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket\"]}]}")
		},
	},
}

func TestMockGetTableBucketPolicy_Success(t *testing.T) {
	for _, c := range testMockGetTableBucketPolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableBucketPolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *GetTableBucketPolicyResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/policy", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetTableBucketPolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{
    "message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/policy", strUrl)
		},
		&GetTableBucketPolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetTableBucketPolicy_Error(t *testing.T) {
	for _, c := range testMockGetTableBucketPolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTableBucketPolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTableBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTableBucketPolicyResult, err error)
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
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/policy", strUrl)
		},
		&DeleteTableBucketPolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketPolicyResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteTableBucketPolicy_Success(t *testing.T) {
	for _, c := range testMockDeleteTableBucketPolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTableBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTableBucketPolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTableBucketPolicyRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTableBucketPolicyResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{
    "message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/policy", strUrl)
		},
		&DeleteTableBucketPolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/policy", strUrl)
		},
		&DeleteTableBucketPolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketPolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteTableBucketPolicy_Error(t *testing.T) {
	for _, c := range testMockDeleteTableBucketPolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTableBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTableBucketMaintenanceConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTableBucketMaintenanceConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *PutTableBucketMaintenanceConfigurationResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/maintenance/icebergUnreferencedFileRemoval", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"value\":{\"settings\":{\"icebergUnreferencedFileRemoval\":{\"nonCurrentDays\":10,\"unreferencedDays\":4}},\"status\":\"disabled\"}}")
		},
		&PutTableBucketMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Type:      oss.Ptr("icebergUnreferencedFileRemoval"),
			Value: &MaintenanceValue{
				Settings: &MaintenanceSettings{
					IcebergUnreferencedFileRemoval: &SettingsDetail{
						UnreferencedDays: oss.Ptr(int(4)),
						NonCurrentDays:   oss.Ptr(10),
					},
				},
				Status: oss.Ptr("disabled"),
			},
		},
		func(t *testing.T, o *PutTableBucketMaintenanceConfigurationResult, err error) {
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
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/maintenance/icebergUnreferencedFileRemoval", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"value\":{\"settings\":{\"icebergUnreferencedFileRemoval\":{\"nonCurrentDays\":1,\"unreferencedDays\":2147483647}},\"status\":\"enabled\"}}")
		},
		&PutTableBucketMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Type:      oss.Ptr("icebergUnreferencedFileRemoval"),
			Value: &MaintenanceValue{
				Settings: &MaintenanceSettings{
					IcebergUnreferencedFileRemoval: &SettingsDetail{
						UnreferencedDays: oss.Ptr(2147483647),
						NonCurrentDays:   oss.Ptr(1),
					},
				},
				Status: oss.Ptr("enabled"),
			},
		},
		func(t *testing.T, o *PutTableBucketMaintenanceConfigurationResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutTableBucketMaintenanceConfiguration_Success(t *testing.T) {
	for _, c := range testMockPutTableBucketMaintenanceConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutTableBucketMaintenanceConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTableBucketMaintenanceConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTableBucketMaintenanceConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *PutTableBucketMaintenanceConfigurationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "5C3D9175B6FC201293AD****",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/maintenance/icebergUnreferencedFileRemoval", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"value\":{\"settings\":{\"icebergUnreferencedFileRemoval\":{\"nonCurrentDays\":1,\"unreferencedDays\":2147483647}},\"status\":\"enabled\"}}")
		},
		&PutTableBucketMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Type:      oss.Ptr("icebergUnreferencedFileRemoval"),
			Value: &MaintenanceValue{
				Settings: &MaintenanceSettings{
					IcebergUnreferencedFileRemoval: &SettingsDetail{
						UnreferencedDays: oss.Ptr(2147483647),
						NonCurrentDays:   oss.Ptr(1),
					},
				},
				Status: oss.Ptr("enabled"),
			},
		},
		func(t *testing.T, o *PutTableBucketMaintenanceConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/maintenance/icebergUnreferencedFileRemoval", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"value\":{\"settings\":{\"icebergUnreferencedFileRemoval\":{\"nonCurrentDays\":1,\"unreferencedDays\":2147483647}},\"status\":\"enabled\"}}")
		},
		&PutTableBucketMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Type:      oss.Ptr("icebergUnreferencedFileRemoval"),
			Value: &MaintenanceValue{
				Settings: &MaintenanceSettings{
					IcebergUnreferencedFileRemoval: &SettingsDetail{
						UnreferencedDays: oss.Ptr(2147483647),
						NonCurrentDays:   oss.Ptr(1),
					},
				},
				Status: oss.Ptr("enabled"),
			},
		},
		func(t *testing.T, o *PutTableBucketMaintenanceConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutTableBucketMaintenanceConfiguration_Error(t *testing.T) {
	for _, c := range testMockPutTableBucketMaintenanceConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutTableBucketMaintenanceConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableBucketMaintenanceConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableBucketMaintenanceConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *GetTableBucketMaintenanceConfigurationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
    "configuration": {"icebergUnreferencedFileRemoval": {
            "settings": {"icebergUnreferencedFileRemoval": {
                    "nonCurrentDays": 2147483647,
                    "unreferencedDays": 10}},
            "status": "enabled"}},
    "tableBucketARN": "acs:osstables:cn-beijing:123456:bucket/demo-bucket"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/maintenance", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetTableBucketMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketMaintenanceConfigurationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Configuration.IcebergUnreferencedFileRemoval.Settings.IcebergUnreferencedFileRemoval.UnreferencedDays, 10)
			assert.Equal(t, *o.Configuration.IcebergUnreferencedFileRemoval.Settings.IcebergUnreferencedFileRemoval.NonCurrentDays, 2147483647)
			assert.Equal(t, *o.Configuration.IcebergUnreferencedFileRemoval.Status, "enabled")
			assert.Equal(t, *o.TableBucketARN, "acs:osstables:cn-beijing:123456:bucket/demo-bucket")
		},
	},
}

func TestMockGetTableBucketMaintenanceConfiguration_Success(t *testing.T) {
	for _, c := range testMockGetTableBucketMaintenanceConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableBucketMaintenanceConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableBucketMaintenanceConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableBucketMaintenanceConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *GetTableBucketMaintenanceConfigurationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/maintenance", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetTableBucketMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketMaintenanceConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/buckets/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/maintenance", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetTableBucketMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketMaintenanceConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetTableBucketMaintenanceConfiguration_Error(t *testing.T) {
	for _, c := range testMockGetTableBucketMaintenanceConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableBucketMaintenanceConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateNamespaceSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateNamespaceRequest
	CheckOutputFn  func(t *testing.T, o *CreateNamespaceResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
   "namespace": [ "space" ],
   "tableBucketARN": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"namespace\":[\"space\"]}", string(requestBody))
		},
		&CreateNamespaceRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: []string{"space"},
		},
		func(t *testing.T, o *CreateNamespaceResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket", *o.TableBucketARN)
			assert.Equal(t, "space", o.Namespace[0])
		},
	},
}

func TestMockCreateNamespace_Success(t *testing.T) {
	for _, c := range testMockCreateNamespaceSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateNamespace(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateNamespaceErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateNamespaceRequest
	CheckOutputFn  func(t *testing.T, o *CreateNamespaceResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0002-00000040",
		},
		[]byte(
			`{
				"message": "The request signature we calculated does not match the signature you provided. Check your key and signing method."
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"namespace\":[\"space\"]}", string(requestBody))
		},
		&CreateNamespaceRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: []string{"space"},
		},
		func(t *testing.T, o *CreateNamespaceResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
		},
	},
	{
		409,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0015-00000104",
		},
		[]byte(
			`{
				"message": "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again."
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"namespace\":[\"space\"]}", string(requestBody))
		},
		&CreateNamespaceRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: []string{"space"},
		},
		func(t *testing.T, o *CreateNamespaceResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "Conflict", serr.Code)
			assert.Equal(t, "0015-00000104", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again")
		},
	},
}

func TestMockCreateNamespace_Error(t *testing.T) {
	for _, c := range testMockCreateNamespaceErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateNamespace(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetNamespaceSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetNamespaceRequest
	CheckOutputFn  func(t *testing.T, o *GetNamespaceResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
   "createdAt": "2026-04-03T09:00:44.014637+00:00",
   "createdBy": "1234567890",
   "namespace": ["my_space"],
   "namespaceId": "0a8fcd4d-a22a-42a4-a3f6-d4a88027018f",
   "ownerAccountId": "1234567890",
   "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5"
}`),
		func(t *testing.T, r *http.Request) {
			urlStr := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", urlStr)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetNamespaceRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *GetNamespaceResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/json")
			assert.Equal(t, *o.CreatedBy, "1234567890")
			assert.Equal(t, *o.CreatedAt, "2026-04-03T09:00:44.014637+00:00")
			assert.Equal(t, *o.OwnerAccountId, "1234567890")
			assert.Equal(t, o.Namespace[0], "my_space")
			assert.Equal(t, *o.NamespaceId, "0a8fcd4d-a22a-42a4-a3f6-d4a88027018f")
			assert.Equal(t, *o.TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")
		},
	},
}

func TestMockGetNamespace_Success(t *testing.T) {
	for _, c := range testMockGetNamespaceSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetNamespace(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetNamespaceErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetNamespaceRequest
	CheckOutputFn  func(t *testing.T, o *GetNamespaceResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			urlStr := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", urlStr)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetNamespaceRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *GetNamespaceResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", urlStr)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetNamespaceRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *GetNamespaceResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetNamespace fail")
		},
	},
}

func TestMockGetNamespace_Error(t *testing.T) {
	for _, c := range testMockGetNamespaceErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetNamespace(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListNamespacesSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListNamespacesRequest
	CheckOutputFn  func(t *testing.T, o *ListNamespacesResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "continuationToken": "CgxteV9uYW1lc3BhY2U-",
  "Namespaces": [{
    "createdAt": "2026-04-03T08:54:25.205905+00:00",
    "createdBy": "1760225545089999",
    "namespace": ["my_namespace"],
    "namespaceId": "22af7160-82b5-4d6a-b9fb-4d14c6e01199",
    "ownerAccountId": "1760225545089999",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f4"
  },
  {
     "createdAt": "2026-04-03T08:59:25.205905+00:00",
    "createdBy": "1760225545089999",
    "namespace": ["demo_namespace"],
    "namespaceId": "22af7160-82b5-4d6a-b9fb-4d14c6e01198",
    "ownerAccountId": "1760225545089999",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5"
  }]
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&ListNamespacesRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *ListNamespacesResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ContinuationToken, "CgxteV9uYW1lc3BhY2U-")
			assert.Equal(t, len(o.Namespaces), 2)
			assert.Equal(t, *o.Namespaces[0].CreatedAt, "2026-04-03T08:54:25.205905+00:00")
			assert.Equal(t, *o.Namespaces[0].CreatedBy, "1760225545089999")
			assert.Equal(t, o.Namespaces[0].Namespace[0], "my_namespace")
			assert.Equal(t, *o.Namespaces[0].NamespaceId, "22af7160-82b5-4d6a-b9fb-4d14c6e01199")
			assert.Equal(t, *o.Namespaces[0].OwnerAccountId, "1760225545089999")
			assert.Equal(t, *o.Namespaces[0].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f4")

			assert.Equal(t, *o.Namespaces[1].CreatedAt, "2026-04-03T08:59:25.205905+00:00")
			assert.Equal(t, *o.Namespaces[1].CreatedBy, "1760225545089999")
			assert.Equal(t, o.Namespaces[1].Namespace[0], "demo_namespace")
			assert.Equal(t, *o.Namespaces[1].NamespaceId, "22af7160-82b5-4d6a-b9fb-4d14c6e01198")
			assert.Equal(t, *o.Namespaces[1].OwnerAccountId, "1760225545089999")
			assert.Equal(t, *o.Namespaces[1].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")
		},
	},
}

func TestMockListNamespaces_Success(t *testing.T) {
	for _, c := range testMockListNamespacesSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListNamespaces(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListNamespacesErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListNamespacesRequest
	CheckOutputFn  func(t *testing.T, o *ListNamespacesResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0002-00000040",
		},
		[]byte(
			`{"message": "The OSS Access Key Id you provided does not exist in our records."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&ListNamespacesRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *ListNamespacesResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
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
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&ListNamespacesRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *ListNamespacesResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListNamespaces fail")
		},
	},
}

func TestMockListNamespaces_Error(t *testing.T) {
	for _, c := range testMockListNamespacesErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListNamespaces(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteNamespaceSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteNamespaceRequest
	CheckOutputFn  func(t *testing.T, o *DeleteNamespaceResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			urlStr := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", urlStr)
			assert.Equal(t, "DELETE", r.Method)
		},
		&DeleteNamespaceRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *DeleteNamespaceResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteNamespace_Success(t *testing.T) {
	for _, c := range testMockDeleteNamespaceSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteNamespace(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteNamespaceErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteNamespaceRequest
	CheckOutputFn  func(t *testing.T, o *DeleteNamespaceResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{"message": "The specified namespace does not exist."}`),
		func(t *testing.T, r *http.Request) {
			urlStr := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", urlStr)
			assert.Equal(t, "DELETE", r.Method)
		},
		&DeleteNamespaceRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *DeleteNamespaceResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
			assert.Equal(t, "The specified namespace does not exist.", serr.Message)
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
			"x-oss-ec":         "0015-00000301",
		},
		[]byte(`{"message": "The bucket has objects. Please delete them first."}`),
		func(t *testing.T, r *http.Request) {
			urlStr := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/namespaces/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", urlStr)
			assert.Equal(t, "DELETE", r.Method)
		},
		&DeleteNamespaceRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *DeleteNamespaceResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "Conflict", serr.Code)
			assert.Equal(t, "The bucket has objects. Please delete them first.", serr.Message)
			assert.Equal(t, "0015-00000301", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteNamespace_Error(t *testing.T) {
	for _, c := range testMockDeleteNamespaceErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteNamespace(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateTableSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateTableRequest
	CheckOutputFn  func(t *testing.T, o *CreateTableResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
			"tableARN": "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/16dc6c23-7a64-4f55-af2f-ee243524a5cc",
			"versionToken": "8c651fb37897499092bd95e1bc2816a9"
		}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(requestBody), "{\"format\":\"ICEBERG\",\"name\":\"table\"}")
		},
		&CreateTableRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
			Format:    oss.Ptr("ICEBERG"),
		},
		func(t *testing.T, o *CreateTableResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, oss.ToString(o.TableARN), "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/16dc6c23-7a64-4f55-af2f-ee243524a5cc")
			assert.Equal(t, oss.ToString(o.VersionToken), "8c651fb37897499092bd95e1bc2816a9")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
			"tableARN": "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/16dc6c23-7a64-4f55-af2f-ee243524a5cc",
			"versionToken": "8c651fb37897499092bd95e1bc2816a9"
		}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(requestBody), "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"\",\"sseAlgorithm\":\"AES256\"},\"format\":\"ICEBERG\",\"metadata\":{\"iceberg\":{\"schema\":{\"fields\":[{\"name\":\"id\",\"required\":true,\"type\":\"int\"},{\"name\":\"name\",\"type\":\"string\"}]}}},\"name\":\"table\"}")
		},
		&CreateTableRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
			Format:    oss.Ptr("ICEBERG"),
			Metadata: &TableMetadata{
				Iceberg: &IcebergMetadata{
					Schema: map[string]any{
						"fields": []map[string]any{
							{
								"name": "id", "type": "int", "required": true,
							},
							{
								"name": "name", "type": "string",
							},
						},
					},
				},
			},
			EncryptionConfiguration: &EncryptionConfiguration{
				KmsKeyArn:    oss.Ptr(""),
				SseAlgorithm: oss.Ptr("AES256"),
			},
		},
		func(t *testing.T, o *CreateTableResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, oss.ToString(o.TableARN), "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/16dc6c23-7a64-4f55-af2f-ee243524a5cc")
			assert.Equal(t, oss.ToString(o.VersionToken), "8c651fb37897499092bd95e1bc2816a9")
		},
	},
}

func TestMockCreateTable_Success(t *testing.T) {
	for _, c := range testMockCreateTableSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateTable(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateTableErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateTableRequest
	CheckOutputFn  func(t *testing.T, o *CreateTableResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0002-00000040",
		},
		[]byte(
			`{
				"message": "The request signature we calculated does not match the signature you provided. Check your key and signing method."
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(requestBody), "{\"format\":\"ICEBERG\",\"name\":\"table\"}")
		},
		&CreateTableRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
			Format:    oss.Ptr("ICEBERG"),
		},
		func(t *testing.T, o *CreateTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
			assert.Contains(t, serr.RequestTarget, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space")
		},
	},
	{
		409,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0015-00000104",
		},
		[]byte(
			`{
				"message": "The request failed because there is a conflict with a previous write. You can retry the request."
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(requestBody), "{\"format\":\"ICEBERG\",\"name\":\"table\"}")
		},
		&CreateTableRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
			Format:    oss.Ptr("ICEBERG"),
		},
		func(t *testing.T, o *CreateTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "Conflict", serr.Code)
			assert.Equal(t, "0015-00000104", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request failed because there is a conflict with a previous write. You can retry the request.")
			assert.Contains(t, serr.RequestTarget, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space")
		},
	},
}

func TestMockCreateTable_Error(t *testing.T) {
	for _, c := range testMockCreateTableErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateTable(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableRequest
	CheckOutputFn  func(t *testing.T, o *GetTableResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
    "createdAt": "2026-04-07T05:27:18.397920+00:00",
    "createdBy": "1234567890",
    "format": "ICEBERG",
    "metadataLocation": "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss/metadata/00000-edb683a9-ce46-492a-a495-35e5b2f7a649.metadata.json",
    "modifiedAt": "2026-04-07T05:27:18.397920+00:00",
    "modifiedBy": "1234567890",
    "name": "my_table",
    "namespace": ["my_namespace"],
    "namespaceId": "22af7160-82b5-4d6a-b9fb-4d14c6e01198",
    "ownerAccountId": "1234567890",
    "tableARN": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5",
    "type": "customer",
    "versionToken": "365f934c6e234f35ace5ae48f0a0d871",
    "warehouseLocation": "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/get-table?name=table&namespace=space&tableBucketARN=acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket", strUrl)
		},
		&GetTableRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/json")
			assert.Equal(t, *o.CreatedAt, "2026-04-07T05:27:18.397920+00:00")
			assert.Equal(t, *o.CreatedBy, "1234567890")
			assert.Equal(t, *o.Format, "ICEBERG")
			assert.Equal(t, *o.MetadataLocation, "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss/metadata/00000-edb683a9-ce46-492a-a495-35e5b2f7a649.metadata.json")
			assert.Equal(t, *o.ModifiedAt, "2026-04-07T05:27:18.397920+00:00")
			assert.Equal(t, *o.ModifiedBy, "1234567890")
			assert.Equal(t, *o.Name, "my_table")
			assert.Equal(t, o.Namespace[0], "my_namespace")
			assert.Equal(t, *o.NamespaceId, "22af7160-82b5-4d6a-b9fb-4d14c6e01198")
			assert.Equal(t, *o.OwnerAccountId, "1234567890")
			assert.Equal(t, *o.TableARN, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9")
			assert.Equal(t, *o.TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")
			assert.Equal(t, *o.Type, "customer")
			assert.Equal(t, *o.VersionToken, "365f934c6e234f35ace5ae48f0a0d871")
			assert.Equal(t, *o.WarehouseLocation, "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss")
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
    "createdAt": "2026-04-07T05:27:18.397920+00:00",
    "createdBy": "1234567890",
    "format": "ICEBERG",
    "metadataLocation": "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss/metadata/00000-edb683a9-ce46-492a-a495-35e5b2f7a649.metadata.json",
    "modifiedAt": "2026-04-07T05:27:18.397920+00:00",
    "modifiedBy": "1234567890",
    "name": "my_table",
    "namespace": ["my_namespace"],
    "namespaceId": "22af7160-82b5-4d6a-b9fb-4d14c6e01198",
    "ownerAccountId": "1234567890",
    "tableARN": "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9",
    "tableBucketId": "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5",
    "type": "customer",
    "versionToken": "365f934c6e234f35ace5ae48f0a0d871",
    "warehouseLocation": "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/get-table?tableArn=acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket%2Ftable%2Ff13de3a6-de93-4801-bd7f-a09c124177d9", strUrl)
		},
		&GetTableRequest{
			TableARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9"),
		},
		func(t *testing.T, o *GetTableResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/json")
			assert.Equal(t, *o.CreatedAt, "2026-04-07T05:27:18.397920+00:00")
			assert.Equal(t, *o.CreatedBy, "1234567890")
			assert.Equal(t, *o.Format, "ICEBERG")
			assert.Equal(t, *o.MetadataLocation, "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss/metadata/00000-edb683a9-ce46-492a-a495-35e5b2f7a649.metadata.json")
			assert.Equal(t, *o.ModifiedAt, "2026-04-07T05:27:18.397920+00:00")
			assert.Equal(t, *o.ModifiedBy, "1234567890")
			assert.Equal(t, *o.Name, "my_table")
			assert.Equal(t, o.Namespace[0], "my_namespace")
			assert.Equal(t, *o.NamespaceId, "22af7160-82b5-4d6a-b9fb-4d14c6e01198")
			assert.Equal(t, *o.OwnerAccountId, "1234567890")
			assert.Equal(t, *o.TableARN, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9")
			assert.Equal(t, *o.TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")
			assert.Equal(t, *o.Type, "customer")
			assert.Equal(t, *o.VersionToken, "365f934c6e234f35ace5ae48f0a0d871")
			assert.Equal(t, *o.WarehouseLocation, "oss://f13de3a6-de93-4801-vlz6uao35255n4bbo5q3sujl1fy83su13--table-oss")
		},
	},
}

func TestMockGetTable_Success(t *testing.T) {
	for _, c := range testMockGetTableSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTable(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableRequest
	CheckOutputFn  func(t *testing.T, o *GetTableResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/get-table?tableArn=acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket%2Ftable%2Ff13de3a6-de93-4801-bd7f-a09c124177d9", strUrl)
		},
		&GetTableRequest{
			TableARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9"),
		},
		func(t *testing.T, o *GetTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
			"x-oss-ec":         "0003-00000801",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/get-table?tableArn=acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket%2Ftable%2Ff13de3a6-de93-4801-bd7f-a09c124177d9", strUrl)
		},
		&GetTableRequest{
			TableARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9"),
		},
		func(t *testing.T, o *GetTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
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
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/get-table?tableArn=acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket%2Ftable%2Ff13de3a6-de93-4801-bd7f-a09c124177d9", strUrl)
		},
		&GetTableRequest{
			TableARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/table/f13de3a6-de93-4801-bd7f-a09c124177d9"),
		},
		func(t *testing.T, o *GetTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetTable fail")
		},
	},
}

func TestMockGetTable_Error(t *testing.T) {
	for _, c := range testMockGetTableErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTable(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListTablesSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListTablesRequest
	CheckOutputFn  func(t *testing.T, o *ListTablesResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "continuationToken": "AAMA-EFRSURBSGk2VFNsQXNjVHQ2QU05UU5YN2xkME53VWI3U1B5RTl6WEh1UTRVc",
  "tables": [
    {
      "createdAt": "2026-04-07T02:15:12.186626+00:00",
      "modifiedAt": "2026-04-07T02:15:12.186626+00:00",
      "name": "example_table",
      "namespace": [
        "my_namespace"
      ],
      "tableARN": "acs:osstables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/7568a090-50f8-4808-8c8d-930a2c264076",
      "type": "customer"
    },
    {
      "createdAt": "2026-04-07T02:15:12.186626+00:00",
      "modifiedAt": "2026-04-07T02:15:12.186626+00:00",
      "name": "example_table1",
      "namespace": [
        "my_namespace"
      ],
      "tableARN": "acs:osstables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/757c17c1-532e-4a45-b5b3-d8783374fc2a",
      "type": "customer"
    }
  ]
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket?continuationToken=token&maxTables=1000&namespace=space&prefix=prefix", strUrl)
		},
		&ListTablesRequest{
			TableBucketARN:         oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:         oss.Ptr("space"),
			ContinuationToken: oss.Ptr("token"),
			MaxTables:         int32(1000),
			Prefix:            oss.Ptr("prefix"),
		},
		func(t *testing.T, o *ListTablesResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ContinuationToken, "AAMA-EFRSURBSGk2VFNsQXNjVHQ2QU05UU5YN2xkME53VWI3U1B5RTl6WEh1UTRVc")
			assert.Equal(t, len(o.Tables), 2)
			assert.Equal(t, *o.Tables[0].CreatedAt, "2026-04-07T02:15:12.186626+00:00")
			assert.Equal(t, *o.Tables[0].ModifiedAt, "2026-04-07T02:15:12.186626+00:00")
			assert.Equal(t, *o.Tables[0].Name, "example_table")
			assert.Equal(t, o.Tables[0].Namespace[0], "my_namespace")
			assert.Equal(t, *o.Tables[0].TableARN, "acs:osstables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/7568a090-50f8-4808-8c8d-930a2c264076")
			assert.Equal(t, *o.Tables[0].Type, "customer")

			assert.Equal(t, *o.Tables[1].CreatedAt, "2026-04-07T02:15:12.186626+00:00")
			assert.Equal(t, *o.Tables[1].ModifiedAt, "2026-04-07T02:15:12.186626+00:00")
			assert.Equal(t, *o.Tables[1].Name, "example_table1")
			assert.Equal(t, o.Tables[1].Namespace[0], "my_namespace")
			assert.Equal(t, *o.Tables[1].TableARN, "acs:osstables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/757c17c1-532e-4a45-b5b3-d8783374fc2a")
			assert.Equal(t, *o.Tables[1].Type, "customer")
		},
	},
}

func TestMockListTables_Success(t *testing.T) {
	for _, c := range testMockListTablesSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListTables(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListTablesErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListTablesRequest
	CheckOutputFn  func(t *testing.T, o *ListTablesResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0002-00000040",
		},
		[]byte(
			`{"message": "The OSS Access Key Id you provided does not exist in our records."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket?namespace=space", strUrl)
		},
		&ListTablesRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *ListTablesResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
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
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket?namespace=space", strUrl)
		},
		&ListTablesRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *ListTablesResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListTables fail")
		},
	},
}

func TestMockListTables_Error(t *testing.T) {
	for _, c := range testMockListTablesErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListTables(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTableSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTableRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTableResult, err error)
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
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table", strUrl)
		},
		&DeleteTableRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTableResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteTable_Success(t *testing.T) {
	for _, c := range testMockDeleteTableSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTable(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTableErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTableRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTableResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
    "message": "The specified bucket does not exist."
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table", strUrl)
		},
		&DeleteTableRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
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
    "message": "The bucket has objects. Please delete them first."
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table", strUrl)
		},
		&DeleteTableRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "Conflict", serr.Code)
			assert.Equal(t, "The bucket has objects. Please delete them first.", serr.Message)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteTable_Error(t *testing.T) {
	for _, c := range testMockDeleteTableErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTable(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockRenameTableSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RenameTableRequest
	CheckOutputFn  func(t *testing.T, o *RenameTableResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/rename", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"newName\":\"new-table\",\"newNamespaceName\":\"new-space\",\"versionToken\":\"version-token\"}", string(requestBody))
		},
		&RenameTableRequest{
			TableBucketARN:        oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:        oss.Ptr("space"),
			Name:             oss.Ptr("table"),
			NewNamespace: oss.Ptr("new-space"),
			NewName:          oss.Ptr("new-table"),
			VersionToken:     oss.Ptr("version-token"),
		},
		func(t *testing.T, o *RenameTableResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/rename", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"newName\":\"new-table\"}", string(requestBody))
		},
		&RenameTableRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
			NewName:   oss.Ptr("new-table"),
		},
		func(t *testing.T, o *RenameTableResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockRenameTable_Success(t *testing.T) {
	for _, c := range testMockRenameTableSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RenameTable(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockRenameTableErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RenameTableRequest
	CheckOutputFn  func(t *testing.T, o *RenameTableResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0002-00000040",
		},
		[]byte(
			`{
				"message": "The request signature we calculated does not match the signature you provided. Check your key and signing method."
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/rename", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"newName\":\"new-table\",\"newNamespaceName\":\"new-space\",\"versionToken\":\"version-token\"}", string(requestBody))
		},
		&RenameTableRequest{
			TableBucketARN:        oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:        oss.Ptr("space"),
			Name:             oss.Ptr("table"),
			NewNamespace: oss.Ptr("new-space"),
			NewName:          oss.Ptr("new-table"),
			VersionToken:     oss.Ptr("version-token"),
		},
		func(t *testing.T, o *RenameTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
			assert.Contains(t, serr.RequestTarget, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/rename")
		},
	},
	{
		409,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
			"x-oss-ec":         "0015-00000104",
		},
		[]byte(
			`{
				"Message": "The request failed because there is a conflict with a previous write. You can retry the request."
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/rename", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"newName\":\"new-table\",\"newNamespaceName\":\"new-space\",\"versionToken\":\"version-token\"}", string(requestBody))
		},
		&RenameTableRequest{
			TableBucketARN:        oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:        oss.Ptr("space"),
			Name:             oss.Ptr("table"),
			NewNamespace: oss.Ptr("new-space"),
			NewName:          oss.Ptr("new-table"),
			VersionToken:     oss.Ptr("version-token"),
		},
		func(t *testing.T, o *RenameTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "Conflict", serr.Code)
			assert.Equal(t, "0015-00000104", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request failed because there is a conflict with a previous write. You can retry the request.")
			assert.Contains(t, serr.RequestTarget, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/rename")
		},
	},
}

func TestMockRenameTable_Error(t *testing.T) {
	for _, c := range testMockRenameTableErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RenameTable(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTablePolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTablePolicyRequest
	CheckOutputFn  func(t *testing.T, o *PutTablePolicyResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"resourcePolicy":"{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"}`, string(requestBody))
		},
		&PutTablePolicyRequest{
			TableBucketARN:      oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:      oss.Ptr("space"),
			Name:           oss.Ptr("table"),
			ResourcePolicy: oss.Ptr("{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"),
		},
		func(t *testing.T, o *PutTablePolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutTablePolicy_Success(t *testing.T) {
	for _, c := range testMockPutTablePolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutTablePolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTablePolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTablePolicyRequest
	CheckOutputFn  func(t *testing.T, o *PutTablePolicyResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(
			`{"message": "The request signature we calculated does not match the signature you provided. Check your key and signing method."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"resourcePolicy":"{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"}`, string(requestBody))
		},
		&PutTablePolicyRequest{
			TableBucketARN:      oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:      oss.Ptr("space"),
			Name:           oss.Ptr("table"),
			ResourcePolicy: oss.Ptr("{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"),
		},
		func(t *testing.T, o *PutTablePolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
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
			`{"message": "The request failed because there is a conflict with a previous write. You can retry the request."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"resourcePolicy":"{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"}`, string(requestBody))
		},
		&PutTablePolicyRequest{
			TableBucketARN:      oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:      oss.Ptr("space"),
			Name:           oss.Ptr("table"),
			ResourcePolicy: oss.Ptr("{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"),
		},
		func(t *testing.T, o *PutTablePolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "Conflict", serr.Code)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request failed because there is a conflict with a previous write. You can retry the request.")
		},
	},
}

func TestMockPutTablePolicy_Error(t *testing.T) {
	for _, c := range testMockPutTablePolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutTablePolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTablePolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTablePolicyRequest
	CheckOutputFn  func(t *testing.T, o *GetTablePolicyResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"resourcePolicy":"{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"oss:GetTable\"],\"Effect\":\"Allow\",\"Principal\":[\"9876543210\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*\"]}]}"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy", strUrl)
		},
		&GetTablePolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Name:      oss.Ptr("table"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *GetTablePolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/json")
			assert.Equal(t, *o.ResourcePolicy, `{"Version":"1","Statement":[{"Action":["oss:GetTable"],"Effect":"Allow","Principal":["9876543210"],"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/my-table-bucket/table/*"]}]}`)
		},
	},
}

func TestMockGetTablePolicy_Success(t *testing.T) {
	for _, c := range testMockGetTablePolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTablePolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTablePolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTablePolicyRequest
	CheckOutputFn  func(t *testing.T, o *GetTablePolicyResult, err error)
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
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy", strUrl)
		},
		&GetTablePolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Name:      oss.Ptr("table"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *GetTablePolicyResult, err error) {
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
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy", strUrl)
		},
		&GetTablePolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Name:      oss.Ptr("table"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *GetTablePolicyResult, err error) {
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

func TestMockGetTablePolicy_Error(t *testing.T) {
	for _, c := range testMockGetTablePolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTablePolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTablePolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTablePolicyRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTablePolicyResult, err error)
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
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy", strUrl)
		},
		&DeleteTablePolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTablePolicyResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteTablePolicy_Success(t *testing.T) {
	for _, c := range testMockDeleteTablePolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTablePolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteTablePolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteTablePolicyRequest
	CheckOutputFn  func(t *testing.T, o *DeleteTablePolicyResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy", strUrl)
		},
		&DeleteTablePolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTablePolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
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
		[]byte(`{"message": "The bucket has objects. Please delete them first."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/policy", strUrl)
		},
		&DeleteTablePolicyRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTablePolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "Conflict", serr.Code)
			assert.Equal(t, "The bucket has objects. Please delete them first.", serr.Message)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteTablePolicy_Error(t *testing.T) {
	for _, c := range testMockDeleteTablePolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteTablePolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableEncryptionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableEncryptionRequest
	CheckOutputFn  func(t *testing.T, o *GetTableEncryptionResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
   "encryptionConfiguration": { 
      "kmsKeyArn": "",
      "sseAlgorithm": "AES256"
   }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/encryption", strUrl)
		},
		&GetTableEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableEncryptionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.EncryptionConfiguration.KmsKeyArn, "")
			assert.Equal(t, *o.EncryptionConfiguration.SseAlgorithm, "AES256")
		},
	},
}

func TestMockGetTableEncryption_Success(t *testing.T) {
	for _, c := range testMockGetTableEncryptionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableEncryption(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableEncryptionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableEncryptionRequest
	CheckOutputFn  func(t *testing.T, o *GetTableEncryptionResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/encryption", strUrl)
		},
		&GetTableEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/encryption", strUrl)
		},
		&GetTableEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
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
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/encryption", strUrl)
		},
		&GetTableEncryptionRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableEncryptionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetTableEncryption fail")
		},
	},
}

func TestMockGetTableEncryption_Error(t *testing.T) {
	for _, c := range testMockGetTableEncryptionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableEncryption(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableMetadataLocationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableMetadataLocationRequest
	CheckOutputFn  func(t *testing.T, o *GetTableMetadataLocationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
   "metadataLocation": "oss://data-bucket/metadata/00000-xxx.metadata.json",
   "warehouseLocation": "oss://eb998f10-d20c-4f22-bmz18s1enia50ot33z1i51zzrrb51b9tc--table-oss",
   "versionToken": "f62eb60ebcd1405db129f4ac86569e2d"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/metadata-location", strUrl)
		},
		&GetTableMetadataLocationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMetadataLocationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.MetadataLocation, "oss://data-bucket/metadata/00000-xxx.metadata.json")
			assert.Equal(t, *o.VersionToken, "f62eb60ebcd1405db129f4ac86569e2d")
			assert.Equal(t, *o.WarehouseLocation, "oss://eb998f10-d20c-4f22-bmz18s1enia50ot33z1i51zzrrb51b9tc--table-oss")
		},
	},
}

func TestMockGetTableMetadataLocation_Success(t *testing.T) {
	for _, c := range testMockGetTableMetadataLocationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableMetadataLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableMetadataLocationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableMetadataLocationRequest
	CheckOutputFn  func(t *testing.T, o *GetTableMetadataLocationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"message": "The specified table does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/metadata-location", strUrl)
		},
		&GetTableMetadataLocationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMetadataLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
			assert.Equal(t, "The specified table does not exist.", serr.Message)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/metadata-location", strUrl)
		},
		&GetTableMetadataLocationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMetadataLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetTableMetadataLocation_Error(t *testing.T) {
	for _, c := range testMockGetTableMetadataLocationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableMetadataLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateTableMetadataLocationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateTableMetadataLocationRequest
	CheckOutputFn  func(t *testing.T, o *UpdateTableMetadataLocationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
   "metadataLocation": "location",
   "name": "table",
   "namespace": [ "space" ],
   "tableARN": "acs:osstable:cn-hangzhou:123:bucket/demo-bucket/table/table_123",
   "versionToken": "aaa"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"metadataLocation\":\"location\",\"versionToken\":\"version-token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/metadata-location", strUrl)
		},
		&UpdateTableMetadataLocationRequest{
			TableBucketARN:        oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:        oss.Ptr("space"),
			Name:             oss.Ptr("table"),
			MetadataLocation: oss.Ptr("location"),
			VersionToken:     oss.Ptr("version-token"),
		},
		func(t *testing.T, o *UpdateTableMetadataLocationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.MetadataLocation, "location")
			assert.Equal(t, *o.Name, "table")
			assert.Equal(t, o.Namespace[0], "space")
			assert.Equal(t, *o.TableARN, "acs:osstable:cn-hangzhou:123:bucket/demo-bucket/table/table_123")
			assert.Equal(t, *o.VersionToken, "aaa")
		},
	},
}

func TestMockUpdateTableMetadataLocation_Success(t *testing.T) {
	for _, c := range testMockUpdateTableMetadataLocationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateTableMetadataLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateTableMetadataLocationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateTableMetadataLocationRequest
	CheckOutputFn  func(t *testing.T, o *UpdateTableMetadataLocationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"metadataLocation\":\"location\",\"versionToken\":\"version-token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/metadata-location", strUrl)
		},
		&UpdateTableMetadataLocationRequest{
			TableBucketARN:        oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:        oss.Ptr("space"),
			Name:             oss.Ptr("table"),
			MetadataLocation: oss.Ptr("location"),
			VersionToken:     oss.Ptr("version-token"),
		},
		func(t *testing.T, o *UpdateTableMetadataLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"metadataLocation\":\"location\",\"versionToken\":\"version-token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/metadata-location", strUrl)
		},
		&UpdateTableMetadataLocationRequest{
			TableBucketARN:        oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace:        oss.Ptr("space"),
			Name:             oss.Ptr("table"),
			MetadataLocation: oss.Ptr("location"),
			VersionToken:     oss.Ptr("version-token"),
		},
		func(t *testing.T, o *UpdateTableMetadataLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockUpdateTableMetadataLocation_Error(t *testing.T) {
	for _, c := range testMockUpdateTableMetadataLocationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateTableMetadataLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTableMaintenanceConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTableMaintenanceConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *PutTableMaintenanceConfigurationResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"value\":{\"settings\":{\"icebergSnapshotManagement\":{\"maxSnapshotAgeHours\":350,\"minSnapshotsToKeep\":1}},\"status\":\"enabled\"}}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance/icebergSnapshotManagement", strUrl)
		},
		&PutTableMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
			Type:      oss.Ptr("icebergSnapshotManagement"),
			Value: &TableMaintenanceValue{
				Status: oss.Ptr("enabled"),
				Settings: &TableMaintenanceSettings{
					IcebergSnapshotManagement: &IcebergSnapshotManagementSettingsDetail{
						MaxSnapshotAgeHours: oss.Ptr(int(350)),
						MinSnapshotsToKeep:  oss.Ptr(1),
					},
				},
			},
		},
		func(t *testing.T, o *PutTableMaintenanceConfigurationResult, err error) {
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
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"value\":{\"settings\":{\"icebergCompaction\":{\"strategy\":\"auto\",\"targetFileSizeMB\":400}},\"status\":\"enabled\"}}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance/icebergCompaction", strUrl)
		},
		&PutTableMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
			Type:      oss.Ptr("icebergCompaction"),
			Value: &TableMaintenanceValue{
				Status: oss.Ptr("enabled"),
				Settings: &TableMaintenanceSettings{
					IcebergCompaction: &IcebergCompactionSettingsDetail{
						TargetFileSizeMB: oss.Ptr(400),
						Strategy:         oss.Ptr("auto"),
					},
				},
			},
		},
		func(t *testing.T, o *PutTableMaintenanceConfigurationResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutTableMaintenanceConfiguration_Success(t *testing.T) {
	for _, c := range testMockPutTableMaintenanceConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutTableMaintenanceConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutTableMaintenanceConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutTableMaintenanceConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *PutTableMaintenanceConfigurationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"value\":{\"settings\":{\"icebergCompaction\":{\"strategy\":\"auto\",\"targetFileSizeMB\":400}},\"status\":\"enabled\"}}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance/icebergCompaction", strUrl)
		},
		&PutTableMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
			Type:      oss.Ptr("icebergCompaction"),
			Value: &TableMaintenanceValue{
				Status: oss.Ptr("enabled"),
				Settings: &TableMaintenanceSettings{
					IcebergCompaction: &IcebergCompactionSettingsDetail{
						TargetFileSizeMB: oss.Ptr(400),
						Strategy:         oss.Ptr("auto"),
					},
				},
			},
		},
		func(t *testing.T, o *PutTableMaintenanceConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"value\":{\"settings\":{\"icebergCompaction\":{\"strategy\":\"auto\",\"targetFileSizeMB\":400}},\"status\":\"enabled\"}}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance/icebergCompaction", strUrl)
		},
		&PutTableMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
			Type:      oss.Ptr("icebergCompaction"),
			Value: &TableMaintenanceValue{
				Status: oss.Ptr("enabled"),
				Settings: &TableMaintenanceSettings{
					IcebergCompaction: &IcebergCompactionSettingsDetail{
						TargetFileSizeMB: oss.Ptr(400),
						Strategy:         oss.Ptr("auto"),
					},
				},
			},
		},
		func(t *testing.T, o *PutTableMaintenanceConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutTableMaintenanceConfiguration_Error(t *testing.T) {
	for _, c := range testMockPutTableMaintenanceConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutTableMaintenanceConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableMaintenanceConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableMaintenanceConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *GetTableMaintenanceConfigurationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
   "tableARN": "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/table_id",
   "configuration": {
      "icebergCompaction": {
         "status": "enabled",
         "settings": {
            "icebergCompaction": {
               "targetFileSizeMB": 512,
               "strategy": "binpack"
            }
         }
      },
      "icebergSnapshotManagement": {
         "status": "enabled",
         "settings": {
            "icebergSnapshotManagement": {
               "minSnapshotsToKeep": 1,
               "maxSnapshotAgeHours": 720
            }
         }
      }
   }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance", strUrl)
		},
		&GetTableMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceConfigurationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Configuration.IcebergCompaction.Settings.IcebergCompaction.TargetFileSizeMB, 512)
			assert.Equal(t, *o.Configuration.IcebergCompaction.Settings.IcebergCompaction.Strategy, "binpack")
			assert.Equal(t, *o.Configuration.IcebergCompaction.Status, "enabled")
			assert.Equal(t, *o.Configuration.IcebergSnapshotManagement.Settings.IcebergSnapshotManagement.MaxSnapshotAgeHours, 720)
			assert.Equal(t, *o.Configuration.IcebergSnapshotManagement.Settings.IcebergSnapshotManagement.MinSnapshotsToKeep, 1)
			assert.Equal(t, *o.Configuration.IcebergSnapshotManagement.Status, "enabled")
			assert.Equal(t, *o.TableARN, "acs:osstable:cn-hangzhou:1234567890:bucket/demo-bucket/table/table_id")
		},
	},
}

func TestMockGetTableMaintenanceConfiguration_Success(t *testing.T) {
	for _, c := range testMockGetTableMaintenanceConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableMaintenanceConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableMaintenanceConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableMaintenanceConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *GetTableMaintenanceConfigurationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance", strUrl)
		},
		&GetTableMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance", strUrl)
		},
		&GetTableMaintenanceConfigurationRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetTableMaintenanceConfiguration_Error(t *testing.T) {
	for _, c := range testMockGetTableMaintenanceConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableMaintenanceConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableMaintenanceJobStatusSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableMaintenanceJobStatusRequest
	CheckOutputFn  func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
    "status": {
        "icebergCompaction": {
            "failureMessage": "internal error",
            "lastRunTimestamp": "2026-04-08T08:37:07.988892Z",
            "status": "Failed"},
        "icebergSnapshotManagement": {
            "failureMessage": "internal error",
            "lastRunTimestamp": "2026-04-08T08:36:07.426846Z",
            "status": "Failed"},
        "icebergUnreferencedFileRemoval": {"status": "Disabled"}},
    "tableARN": "acs:osstables:cn-beijing:123456:bucket/demo-bucket/table/eb998f10-d20c-4f22-9a76-ed64e9668f56"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance-job-status", strUrl)
		},
		&GetTableMaintenanceJobStatusRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.JobStatus.IcebergCompaction.FailureMessage, "internal error")
			assert.Equal(t, *o.JobStatus.IcebergCompaction.Status, "Failed")
			assert.Equal(t, *o.JobStatus.IcebergCompaction.LastRunTimestamp, "2026-04-08T08:37:07.988892Z")
			assert.Equal(t, *o.JobStatus.IcebergSnapshotManagement.FailureMessage, "internal error")
			assert.Equal(t, *o.JobStatus.IcebergSnapshotManagement.Status, "Failed")
			assert.Equal(t, *o.JobStatus.IcebergSnapshotManagement.LastRunTimestamp, "2026-04-08T08:36:07.426846Z")
			assert.Equal(t, *o.JobStatus.IcebergUnreferencedFileRemoval.Status, "Disabled")
			assert.Equal(t, *o.TableARN, "acs:osstables:cn-beijing:123456:bucket/demo-bucket/table/eb998f10-d20c-4f22-9a76-ed64e9668f56")
		},
	},
}

func TestMockGetTableMaintenanceJobStatus_Success(t *testing.T) {
	for _, c := range testMockGetTableMaintenanceJobStatusSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableMaintenanceJobStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableMaintenanceJobStatusErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableMaintenanceJobStatusRequest
	CheckOutputFn  func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"message": "The specified bucket does not exist."}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance-job-status", strUrl)
		},
		&GetTableMaintenanceJobStatusRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "Not Found", serr.Code)
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
		[]byte(`{"message": "UserDisable"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/acs:osstables:cn-beijing:1234567890:bucket/demo-bucket/tables/acs%3Aosstables%3Acn-beijing%3A1234567890%3Abucket%2Fdemo-bucket/space/table/maintenance-job-status", strUrl)
		},
		&GetTableMaintenanceJobStatusRequest{
			TableBucketARN: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
			Namespace: oss.Ptr("space"),
			Name:      oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "Forbidden", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetTableMaintenanceJobStatus_Error(t *testing.T) {
	for _, c := range testMockGetTableMaintenanceJobStatusErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableMaintenanceJobStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
