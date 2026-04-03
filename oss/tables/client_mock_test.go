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
			Bucket: oss.Ptr("bucket"),
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
			assert.Equal(t, "acs:osstables:cn-beijing:1234567890:bucket/bucket", *o.BucketArn)
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
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *CreateTableBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "acs:osstables:cn-beijing:1234567890:bucket/bucket", *o.BucketArn)
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
			Bucket: oss.Ptr("bucket"),
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
			Bucket: oss.Ptr("bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
		},
		func(t *testing.T, o *GetTableBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, o.Headers.Get("Content-Type"), "application/json")
			assert.Equal(t, *o.BucketArn, "acs:osstables:cn-beijing:12345657890:bucket/demo-bucket")
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			assert.Equal(t, "/?buckets", r.URL.String())
		},
		nil,
		func(t *testing.T, o *ListTableBucketsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ContinuationToken, "token-123")
			assert.Equal(t, len(o.Buckets), 2)
			assert.Equal(t, *o.Buckets[0].CreatedAt, "2026-04-02T05:27:31.000000+00:00")
			assert.Equal(t, *o.Buckets[0].BucketArn, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket")
			assert.Equal(t, *o.Buckets[0].Name, "demo-bucket")
			assert.Equal(t, *o.Buckets[0].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f5")
			assert.Equal(t, *o.Buckets[0].OwnerAccountId, "1234567890")
			assert.Equal(t, *o.Buckets[0].Type, "customer")

			assert.Equal(t, *o.Buckets[1].CreatedAt, "2026-04-02T05:27:32.000000+00:00")
			assert.Equal(t, *o.Buckets[1].BucketArn, "acs:osstables:cn-beijing:1234567890:bucket/demo-bucket-1")
			assert.Equal(t, *o.Buckets[1].Name, "demo-bucket-1")
			assert.Equal(t, *o.Buckets[1].TableBucketId, "340c6672-0a1f-4426-aff9-1a8e2ac7b0f3")
			assert.Equal(t, *o.Buckets[1].OwnerAccountId, "1234567890")
			assert.Equal(t, *o.Buckets[1].Type, "customer")
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
			assert.Equal(t, "/?buckets", r.URL.String())
		},
		&ListTableBucketsRequest{},
		func(t *testing.T, o *ListTableBucketsResult, err error) {
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
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/?buckets", r.URL.String())
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			BucketArn: oss.Ptr("acs:osstables:cn-beijing:1234567890:bucket/demo-bucket"),
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
			assert.Equal(t, "/bucket/?policy", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"osstable:GetTable\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/table\"]}]}")
		},
		&PutTableBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
			Body:   strings.NewReader(`{"Version":"1","Statement":[{"Action":["osstable:GetTable"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/table"]}]}`),
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
			assert.Equal(t, string(data), "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"osstable:GetTable\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/table\"]}]}")
		},
		&PutTableBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
			Body:   strings.NewReader(`{"Version":"1","Statement":[{"Action":["osstable:GetTable"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/table"]}]}`),
		},
		func(t *testing.T, o *PutTableBucketPolicyResult, err error) {
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
			assert.Equal(t, string(data), "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"osstable:GetTable\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/table\"]}]}")
		},
		&PutTableBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
			Body:   strings.NewReader(`{"Version":"1","Statement":[{"Action":["osstable:GetTable"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/table"]}]}`),
		},
		func(t *testing.T, o *PutTableBucketPolicyResult, err error) {
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
		},
		[]byte(`{"Version":"1","Statement":[{"Action":["osstable:GetTable"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:osstable:cn-hangzhou:1234567890:bucket/table"]}]}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?policy", r.URL.String())
		},
		&GetTableBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetTableBucketPolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Body, "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"osstable:GetTable\"],\"Effect\":\"Deny\",\"Principal\":[\"1234567890\"],\"Resource\":[\"acs:osstable:cn-hangzhou:1234567890:bucket/table\"]}]}")
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
		&GetTableBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetTableBucketPolicyResult, err error) {
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
		&GetTableBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetTableBucketPolicyResult, err error) {
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
			assert.Equal(t, "/bucket/?policy", strUrl)
		},
		&DeleteTableBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
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
		&DeleteTableBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketPolicyResult, err error) {
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
		&DeleteTableBucketPolicyRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteTableBucketPolicyResult, err error) {
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
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?maintenance", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"icebergUnreferencedFileRemoval\":{\"settings\":{\"unreferencedDays\":4,\"nonCurrentDays\":10},\"status\":\"enable\"}}")
		},
		&PutTableBucketMaintenanceConfigurationRequest{
			Bucket: oss.Ptr("bucket"),
			IcebergUnreferencedFileRemoval: &IcebergUnreferencedFileRemoval{
				Settings: &MaintenanceSettings{
					UnreferencedDays: oss.Ptr(int64(4)),
					NonCurrentDays:   oss.Ptr(int64(10)),
				},
				Status: oss.Ptr("enable"),
			},
		},
		func(t *testing.T, o *PutTableBucketMaintenanceConfigurationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
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
			assert.Equal(t, "/bucket/?maintenance", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"icebergUnreferencedFileRemoval\":{\"settings\":{\"unreferencedDays\":4,\"nonCurrentDays\":10},\"status\":\"enable\"}}")
		},
		&PutTableBucketMaintenanceConfigurationRequest{
			Bucket: oss.Ptr("bucket"),
			IcebergUnreferencedFileRemoval: &IcebergUnreferencedFileRemoval{
				Settings: &MaintenanceSettings{
					UnreferencedDays: oss.Ptr(int64(4)),
					NonCurrentDays:   oss.Ptr(int64(10)),
				},
				Status: oss.Ptr("enable"),
			},
		},
		func(t *testing.T, o *PutTableBucketMaintenanceConfigurationResult, err error) {
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
			assert.Equal(t, "/bucket/?maintenance", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"icebergUnreferencedFileRemoval\":{\"settings\":{\"unreferencedDays\":4,\"nonCurrentDays\":10},\"status\":\"enable\"}}")
		},
		&PutTableBucketMaintenanceConfigurationRequest{
			Bucket: oss.Ptr("bucket"),
			IcebergUnreferencedFileRemoval: &IcebergUnreferencedFileRemoval{
				Settings: &MaintenanceSettings{
					UnreferencedDays: oss.Ptr(int64(4)),
					NonCurrentDays:   oss.Ptr(int64(10)),
				},
				Status: oss.Ptr("enable"),
			},
		},
		func(t *testing.T, o *PutTableBucketMaintenanceConfigurationResult, err error) {
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
   "configuration": { 
      "icebergUnreferencedFileRemoval": {
        "settings": {
          "unreferencedDays":4,
          "nonCurrentDays":10
        },
        "status": "enable"
     }
   },
   "tableBucketARN": "test-arn"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?maintenance", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetTableBucketMaintenanceConfigurationRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetTableBucketMaintenanceConfigurationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Configuration.IcebergUnreferencedFileRemoval.Settings.UnreferencedDays, int64(4))
			assert.Equal(t, *o.Configuration.IcebergUnreferencedFileRemoval.Settings.NonCurrentDays, int64(10))
			assert.Equal(t, *o.Configuration.IcebergUnreferencedFileRemoval.Status, "enable")
			assert.Equal(t, *o.TableBucketARN, "test-arn")
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
			assert.Equal(t, "/bucket/?maintenance", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetTableBucketMaintenanceConfigurationRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetTableBucketMaintenanceConfigurationResult, err error) {
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
			assert.Equal(t, "/bucket/?maintenance", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetTableBucketMaintenanceConfigurationRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetTableBucketMaintenanceConfigurationResult, err error) {
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
   "tableBucketARN": "bucket-arn"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/?namespaces", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"namespace\":[\"space\"]}", string(requestBody))
		},
		&CreateNamespaceRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: []string{"space"},
		},
		func(t *testing.T, o *CreateNamespaceResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "bucket-arn", *o.TableBucketARN)
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
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/?namespaces", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"namespace\":[\"space\"]}", string(requestBody))
		},
		&CreateNamespaceRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: []string{"space"},
		},
		func(t *testing.T, o *CreateNamespaceResult, err error) {
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
			assert.Contains(t, serr.RequestTarget, "/?namespaces")
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
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/?namespaces", r.URL.String())
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"namespace\":[\"space\"]}", string(requestBody))
		},
		&CreateNamespaceRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: []string{"space"},
		},
		func(t *testing.T, o *CreateNamespaceResult, err error) {
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
			assert.Contains(t, serr.RequestTarget, "/?namespaces")
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
   "createdAt": "2013-07-31T10:56:21.000Z",
   "createdBy": "aliyun",
   "namespace": ["123"],
   "namespaceId": "123",
   "ownerAccountId": "123456",
   "tableBucketId": "1"
}`),
		func(t *testing.T, r *http.Request) {
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?namespaces&space", urlStr)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetNamespaceRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *GetNamespaceResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/json")
			assert.Equal(t, *o.CreatedBy, "aliyun")
			assert.Equal(t, *o.CreatedAt, "2013-07-31T10:56:21.000Z")
			assert.Equal(t, *o.OwnerAccountId, "123456")
			assert.Equal(t, o.Namespace[0], "123")
			assert.Equal(t, *o.NamespaceId, "123")
			assert.Equal(t, *o.TableBucketId, "1")
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?namespaces&space", urlStr)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetNamespaceRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *GetNamespaceResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?namespaces&space", urlStr)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
		},
		&GetNamespaceRequest{
			Bucket:    oss.Ptr("bucket"),
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
  "continuationToken": "token-123",
  "Namespaces": [{
    "createdAt": "2026-01-31T10:56:21.000Z",
    "createdBy": "aliyun",
    "namespace": ["demo-space"],
    "namespaceId": "123",
    "ownerAccountId": "123456",
    "tableBucketId": "1"
  },
  {
     "createdAt": "2026-02-31T10:56:21.000Z",
    "createdBy": "aliyun",
    "namespace": ["oss-space"],
    "namespaceId": "123457",
    "ownerAccountId": "123456",
    "tableBucketId": "2"
  }]
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?namespaces", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&ListNamespacesRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListNamespacesResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ContinuationToken, "token-123")
			assert.Equal(t, len(o.Namespaces), 2)
			assert.Equal(t, *o.Namespaces[0].CreatedAt, "2026-01-31T10:56:21.000Z")
			assert.Equal(t, *o.Namespaces[0].CreatedBy, "aliyun")
			assert.Equal(t, o.Namespaces[0].Namespace[0], "demo-space")
			assert.Equal(t, *o.Namespaces[0].NamespaceId, "123")
			assert.Equal(t, *o.Namespaces[0].OwnerAccountId, "123456")
			assert.Equal(t, *o.Namespaces[0].TableBucketId, "1")

			assert.Equal(t, *o.Namespaces[1].CreatedAt, "2026-02-31T10:56:21.000Z")
			assert.Equal(t, *o.Namespaces[1].CreatedBy, "aliyun")
			assert.Equal(t, o.Namespaces[1].Namespace[0], "oss-space")
			assert.Equal(t, *o.Namespaces[1].NamespaceId, "123457")
			assert.Equal(t, *o.Namespaces[1].OwnerAccountId, "123456")
			assert.Equal(t, *o.Namespaces[1].TableBucketId, "2")
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
			assert.Equal(t, "/bucket/?namespaces", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&ListNamespacesRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListNamespacesResult, err error) {
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
			assert.Equal(t, "/bucket/?namespaces", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&ListNamespacesRequest{
			Bucket: oss.Ptr("bucket"),
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
			assert.Equal(t, "/bucket/?namespaces&space", urlStr)
			assert.Equal(t, "DELETE", r.Method)
		},
		&DeleteNamespaceRequest{
			Bucket:    oss.Ptr("bucket"),
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
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchNamespace",
    "Message": "The specified namespace does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?namespaces&space", urlStr)
			assert.Equal(t, "DELETE", r.Method)
		},
		&DeleteNamespaceRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *DeleteNamespaceResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchNamespace", serr.Code)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?namespaces&space", urlStr)
			assert.Equal(t, "DELETE", r.Method)
		},
		&DeleteNamespaceRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *DeleteNamespaceResult, err error) {
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
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?space&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"format\":\"iceberg\",\"metadata\":{\"iceberg\":{\"schema\":{\"fields\":[{\"name\":\"id\",\"required\":true,\"type\":\"int\"},{\"name\":\"name\",\"type\":\"string\"}]}}},\"name\":\"table\"}", string(requestBody))
		},
		&CreateTableRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Format:    oss.Ptr("iceberg"),
			Table:     oss.Ptr("table"),
			Metadata: &TableMetadata{
				Iceberg: &MetadataIceberg{
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
		},
		func(t *testing.T, o *CreateTableResult, err error) {
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
			assert.Equal(t, "/bucket/?space&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"encryptionConfiguration\":{\"kmsKeyArn\":\"arn\",\"sseAlgorithm\":\"AES256\"},\"format\":\"iceberg\",\"metadata\":{\"iceberg\":{\"schema\":{\"fields\":[{\"name\":\"id\",\"required\":true,\"type\":\"int\"},{\"name\":\"name\",\"type\":\"string\"}]}}},\"name\":\"table\",\"storageClassConfiguration\":{\"storageClass\":\"Standard\"},\"tags\":{\"k1\":\"v1\",\"k2\":\"v2\"}}", string(requestBody))
		},
		&CreateTableRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Format:    oss.Ptr("iceberg"),
			Table:     oss.Ptr("table"),
			Metadata: &TableMetadata{
				Iceberg: &MetadataIceberg{
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
				KmsKeyArn:    oss.Ptr("arn"),
				SseAlgorithm: oss.Ptr("AES256"),
			},
			Tags: map[string]any{
				"k1": "v1", "k2": "v2",
			},
		},
		func(t *testing.T, o *CreateTableResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
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
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?space&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"format\":\"iceberg\",\"metadata\":{\"iceberg\":{\"schema\":{\"fields\":[{\"name\":\"id\",\"required\":true,\"type\":\"int\"},{\"name\":\"name\",\"type\":\"string\"}]}}},\"name\":\"table\"}", string(requestBody))
		},
		&CreateTableRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Format:    oss.Ptr("iceberg"),
			Table:     oss.Ptr("table"),
			Metadata: &TableMetadata{
				Iceberg: &MetadataIceberg{
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
		},
		func(t *testing.T, o *CreateTableResult, err error) {
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
				"Code": "ConflictException",
				"Message": "The request failed because there is a conflict with a previous write. You can retry the request.",
				"RequestId": "6548A043CA31D****",
				"EC": "0015-00000104"
			  }
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?space&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"format\":\"iceberg\",\"metadata\":{\"iceberg\":{\"schema\":{\"fields\":[{\"name\":\"id\",\"required\":true,\"type\":\"int\"},{\"name\":\"name\",\"type\":\"string\"}]}}},\"name\":\"table\"}", string(requestBody))
		},
		&CreateTableRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Format:    oss.Ptr("iceberg"),
			Table:     oss.Ptr("table"),
			Metadata: &TableMetadata{
				Iceberg: &MetadataIceberg{
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
		},
		func(t *testing.T, o *CreateTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "ConflictException", serr.Code)
			assert.Equal(t, "0015-00000104", serr.EC)
			assert.Equal(t, "6548A043CA31D****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request failed because there is a conflict with a previous write. You can retry the request.")
			assert.Contains(t, serr.RequestTarget, "/bucket")
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
   "createdAt": "2026-02-31T10:56:21.000Z",
   "createdBy": "oss-create",
   "format": "demo-format",
   "metadataLocation": "location",
   "modifiedAt": "2026-03-01T10:56:21.000Z",
   "modifiedBy": "oss-modify",
   "name": "table",
   "namespace": [ "space" ],
   "namespaceId": "space-01",
   "ownerAccountId": "123",
   "tableARN": "acs:osstable:cn-hangzhou:123:bucket/table_bucket/table/table_123",
   "tableBucketId": "table_bucket_123",
   "type": "oss",
   "versionToken": "aaa",
   "warehouseLocation": "bbb"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?get-table&name=table&namespace=space&tableArn=table-arn&tableBucketARN=table-bucket-arn", strUrl)
		},
		&GetTableRequest{
			Bucket:         oss.Ptr("bucket"),
			Table:          oss.Ptr("table"),
			Namespace:      oss.Ptr("space"),
			TableArn:       oss.Ptr("table-arn"),
			TableBucketARN: oss.Ptr("table-bucket-arn"),
		},
		func(t *testing.T, o *GetTableResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/json")
			assert.Equal(t, *o.CreatedAt, "2026-02-31T10:56:21.000Z")
			assert.Equal(t, *o.CreatedBy, "oss-create")
			assert.Equal(t, *o.Format, "demo-format")
			assert.Equal(t, *o.MetadataLocation, "location")
			assert.Equal(t, *o.ModifiedAt, "2026-03-01T10:56:21.000Z")
			assert.Equal(t, *o.ModifiedBy, "oss-modify")
			assert.Equal(t, *o.Name, "table")
			assert.Equal(t, o.Namespace[0], "space")
			assert.Equal(t, *o.NamespaceId, "space-01")
			assert.Equal(t, *o.Type, "oss")
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
			assert.Equal(t, "/bucket/?get-table&name=table&namespace=space&tableArn=table-arn&tableBucketARN=table-bucket-arn", strUrl)
		},
		&GetTableRequest{
			Bucket:         oss.Ptr("bucket"),
			Table:          oss.Ptr("table"),
			Namespace:      oss.Ptr("space"),
			TableArn:       oss.Ptr("table-arn"),
			TableBucketARN: oss.Ptr("table-bucket-arn"),
		},
		func(t *testing.T, o *GetTableResult, err error) {
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
			assert.Equal(t, "/bucket/?get-table&name=table&namespace=space&tableArn=table-arn&tableBucketARN=table-bucket-arn", strUrl)
		},
		&GetTableRequest{
			Bucket:         oss.Ptr("bucket"),
			Table:          oss.Ptr("table"),
			Namespace:      oss.Ptr("space"),
			TableArn:       oss.Ptr("table-arn"),
			TableBucketARN: oss.Ptr("table-bucket-arn"),
		},
		func(t *testing.T, o *GetTableResult, err error) {
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
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?get-table&name=table&namespace=space&tableArn=table-arn&tableBucketARN=table-bucket-arn", strUrl)
		},
		&GetTableRequest{
			Bucket:         oss.Ptr("bucket"),
			Table:          oss.Ptr("table"),
			Namespace:      oss.Ptr("space"),
			TableArn:       oss.Ptr("table-arn"),
			TableBucketARN: oss.Ptr("table-bucket-arn"),
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
      "createdAt": "2026-01-26T03:11:51.527997035Z",
      "modifiedAt": "2026-01-26T03:11:51.527997035Z",
      "name": "example_table",
      "namespace": [
        "my_namespace"
      ],
      "tableARN": "arn:aws:s3tables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/7568a090-50f8-4808-8c8d-930a2c264076",
      "type": "customer"
    },
    {
      "createdAt": "2026-01-26T03:16:46.622650810Z",
      "modifiedAt": "2026-01-26T03:16:46.622650810Z",
      "name": "example_table1",
      "namespace": [
        "my_namespace"
      ],
      "tableARN": "arn:aws:s3tables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/757c17c1-532e-4a45-b5b3-d8783374fc2a",
      "type": "customer"
    }
  ]
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?continuationToken=token&maxTables=1000&namespace=space&prefix=prefix&tables", strUrl)
		},
		&ListTablesRequest{
			Bucket:            oss.Ptr("bucket"),
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
			assert.Equal(t, *o.Tables[0].CreatedAt, "2026-01-26T03:11:51.527997035Z")
			assert.Equal(t, *o.Tables[0].ModifiedAt, "2026-01-26T03:11:51.527997035Z")
			assert.Equal(t, *o.Tables[0].Name, "example_table")
			assert.Equal(t, o.Tables[0].Namespace[0], "my_namespace")
			assert.Equal(t, *o.Tables[0].TableARN, "arn:aws:s3tables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/7568a090-50f8-4808-8c8d-930a2c264076")
			assert.Equal(t, *o.Tables[0].Type, "customer")

			assert.Equal(t, *o.Tables[1].CreatedAt, "2026-01-26T03:16:46.622650810Z")
			assert.Equal(t, *o.Tables[1].ModifiedAt, "2026-01-26T03:16:46.622650810Z")
			assert.Equal(t, *o.Tables[1].Name, "example_table1")
			assert.Equal(t, o.Tables[1].Namespace[0], "my_namespace")
			assert.Equal(t, *o.Tables[1].TableARN, "arn:aws:s3tables:ap-southeast-1:651322719100:bucket/donggu-table-bucket-test/table/757c17c1-532e-4a45-b5b3-d8783374fc2a")
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
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?namespace=space&tables", strUrl)
		},
		&ListTablesRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *ListTablesResult, err error) {
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
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?namespace=space&tables", strUrl)
		},
		&ListTablesRequest{
			Bucket:    oss.Ptr("bucket"),
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
			assert.Equal(t, "/bucket/?name=table&namespace=space&tables", strUrl)
		},
		&DeleteTableRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
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
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?name=table&namespace=space&tables", strUrl)
		},
		&DeleteTableRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTableResult, err error) {
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
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?name=table&namespace=space&tables", strUrl)
		},
		&DeleteTableRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTableResult, err error) {
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
			assert.Equal(t, "/bucket/?rename&space&table&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"namespace\":\"new-space\",\"newName\":\"new-table\",\"versionToken\":\"version-token\"}", string(requestBody))
		},
		&RenameTableRequest{
			Bucket:       oss.Ptr("bucket"),
			Namespace:    oss.Ptr("space"),
			Table:        oss.Ptr("table"),
			NewNamespace: oss.Ptr("new-space"),
			NewTable:     oss.Ptr("new-table"),
			VersionToken: oss.Ptr("version-token"),
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
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?rename&space&table&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"namespace\":\"new-space\",\"newName\":\"new-table\",\"versionToken\":\"version-token\"}", string(requestBody))
		},
		&RenameTableRequest{
			Bucket:       oss.Ptr("bucket"),
			Namespace:    oss.Ptr("space"),
			Table:        oss.Ptr("table"),
			NewNamespace: oss.Ptr("new-space"),
			NewTable:     oss.Ptr("new-table"),
			VersionToken: oss.Ptr("version-token"),
		},
		func(t *testing.T, o *RenameTableResult, err error) {
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
				"Code": "ConflictException",
				"Message": "The request failed because there is a conflict with a previous write. You can retry the request.",
				"RequestId": "6548A043CA31D****",
				"EC": "0015-00000104"
			  }
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?rename&space&table&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"namespace\":\"new-space\",\"newName\":\"new-table\",\"versionToken\":\"version-token\"}", string(requestBody))
		},
		&RenameTableRequest{
			Bucket:       oss.Ptr("bucket"),
			Namespace:    oss.Ptr("space"),
			Table:        oss.Ptr("table"),
			NewNamespace: oss.Ptr("new-space"),
			NewTable:     oss.Ptr("new-table"),
			VersionToken: oss.Ptr("version-token"),
		},
		func(t *testing.T, o *RenameTableResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "ConflictException", serr.Code)
			assert.Equal(t, "0015-00000104", serr.EC)
			assert.Equal(t, "6548A043CA31D****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request failed because there is a conflict with a previous write. You can retry the request.")
			assert.Contains(t, serr.RequestTarget, "/bucket")
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
			assert.Equal(t, "/bucket/?policy&space&table&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/table/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`, string(requestBody))
		},
		&PutTablePolicyRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
			Body:      strings.NewReader(`{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/table/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`),
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
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy&space&table&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/table/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`, string(requestBody))
		},
		&PutTablePolicyRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
			Body:      strings.NewReader(`{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/table/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`),
		},
		func(t *testing.T, o *PutTablePolicyResult, err error) {
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
				"Code": "ConflictException",
				"Message": "The request failed because there is a conflict with a previous write. You can retry the request.",
				"RequestId": "6548A043CA31D****",
				"EC": "0015-00000104"
			  }
			}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy&space&table&tables", strUrl)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/table/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`, string(requestBody))
		},
		&PutTablePolicyRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
			Body:      strings.NewReader(`{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/table/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`),
		},
		func(t *testing.T, o *PutTablePolicyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "ConflictException", serr.Code)
			assert.Equal(t, "0015-00000104", serr.EC)
			assert.Equal(t, "6548A043CA31D****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request failed because there is a conflict with a previous write. You can retry the request.")
			assert.Contains(t, serr.RequestTarget, "/bucket")
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
		[]byte(`{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/xfz-table-bucket/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy&space&table&tables", strUrl)
		},
		&GetTablePolicyRequest{
			Bucket:    oss.Ptr("bucket"),
			Table:     oss.Ptr("table"),
			Namespace: oss.Ptr("space"),
		},
		func(t *testing.T, o *GetTablePolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/json")
			assert.Equal(t, o.Body, `{"resourcePolicy":"\"Version\":\"2012-10-17\",\"Id\":\"DeleteTable\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"OSS\":\"arn:oss:iam::651322719100:user/jiangqi\"},\"Action\":[\"osstables:DeleteTable\",\"osstables:UpdateTableMetadataLocation\",\"osstables:PutTableData\",\"osstables:GetTableMetadataLocation\"],\"Resource\":\"arn:oss:osstables:cn-hangzhou:651322719100:bucket/xfz-table-bucket/table/af5ab6a4-f9a5-4d9b-8e89-eb9c6f1c0c8f\""}`)
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
			assert.Equal(t, "/bucket/?policy&space&table&tables", strUrl)
		},
		&GetTablePolicyRequest{
			Bucket:    oss.Ptr("bucket"),
			Table:     oss.Ptr("table"),
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
			assert.Equal(t, "/bucket/?policy&space&table&tables", strUrl)
		},
		&GetTablePolicyRequest{
			Bucket:    oss.Ptr("bucket"),
			Table:     oss.Ptr("table"),
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
			assert.Equal(t, "/bucket/?policy&space&table&tables", strUrl)
		},
		&DeleteTablePolicyRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
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
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy&space&table&tables", strUrl)
		},
		&DeleteTablePolicyRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTablePolicyResult, err error) {
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
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?policy&space&table&tables", strUrl)
		},
		&DeleteTablePolicyRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *DeleteTablePolicyResult, err error) {
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
      "kmsKeyArn": "test-arn",
      "sseAlgorithm": "AES256"
   }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?encryption&space&table&tables", strUrl)
		},
		&GetTableEncryptionRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableEncryptionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.EncryptionConfiguration.KmsKeyArn, "test-arn")
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
			assert.Equal(t, "/bucket/?encryption&space&table&tables", strUrl)
		},
		&GetTableEncryptionRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableEncryptionResult, err error) {
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
			assert.Equal(t, "/bucket/?encryption&space&table&tables", strUrl)
		},
		&GetTableEncryptionRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableEncryptionResult, err error) {
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
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?encryption&space&table&tables", strUrl)
		},
		&GetTableEncryptionRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
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
   "metadataLocation": "location",
   "warehouseLocation": "bbb",
   "versionToken": "aaa"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?metadata-location&space&table&tables", strUrl)
		},
		&GetTableMetadataLocationRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMetadataLocationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.MetadataLocation, "location")
			assert.Equal(t, *o.VersionToken, "aaa")
			assert.Equal(t, *o.WarehouseLocation, "bbb")
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
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?metadata-location&space&table&tables", strUrl)
		},
		&GetTableMetadataLocationRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMetadataLocationResult, err error) {
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
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?metadata-location&space&table&tables", strUrl)
		},
		&GetTableMetadataLocationRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMetadataLocationResult, err error) {
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
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"metadataLocation\":\"location\",\"versionToken\":\"version-token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?metadata-location&space&table&tables", strUrl)
		},
		&UpdateTableMetadataLocationRequest{
			Bucket:           oss.Ptr("bucket"),
			Namespace:        oss.Ptr("space"),
			Table:            oss.Ptr("table"),
			MetadataLocation: oss.Ptr("location"),
			VersionToken:     oss.Ptr("version-token"),
		},
		func(t *testing.T, o *UpdateTableMetadataLocationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
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
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"metadataLocation\":\"location\",\"versionToken\":\"version-token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?metadata-location&space&table&tables", strUrl)
		},
		&UpdateTableMetadataLocationRequest{
			Bucket:           oss.Ptr("bucket"),
			Namespace:        oss.Ptr("space"),
			Table:            oss.Ptr("table"),
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
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"metadataLocation\":\"location\",\"versionToken\":\"version-token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?metadata-location&space&table&tables", strUrl)
		},
		&UpdateTableMetadataLocationRequest{
			Bucket:           oss.Ptr("bucket"),
			Namespace:        oss.Ptr("space"),
			Table:            oss.Ptr("table"),
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
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
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
		200,
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
			assert.Equal(t, string(body), "{\"icebergUnreferencedFileRemoval\":{\"settings\":{\"unreferencedDays\":4,\"nonCurrentDays\":10},\"status\":\"enable\"}}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance&space&table&tables", strUrl)
		},
		&PutTableMaintenanceConfigurationRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
			IcebergUnreferencedFileRemoval: &IcebergUnreferencedFileRemoval{
				Settings: &MaintenanceSettings{
					UnreferencedDays: oss.Ptr(int64(4)),
					NonCurrentDays:   oss.Ptr(int64(10)),
				},
				Status: oss.Ptr("enable"),
			},
		},
		func(t *testing.T, o *PutTableMaintenanceConfigurationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
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
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"icebergUnreferencedFileRemoval\":{\"settings\":{\"unreferencedDays\":4,\"nonCurrentDays\":10},\"status\":\"enable\"}}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance&space&table&tables", strUrl)
		},
		&PutTableMaintenanceConfigurationRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
			IcebergUnreferencedFileRemoval: &IcebergUnreferencedFileRemoval{
				Settings: &MaintenanceSettings{
					UnreferencedDays: oss.Ptr(int64(4)),
					NonCurrentDays:   oss.Ptr(int64(10)),
				},
				Status: oss.Ptr("enable"),
			},
		},
		func(t *testing.T, o *PutTableMaintenanceConfigurationResult, err error) {
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
			assert.Equal(t, r.Method, "PUT")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"icebergUnreferencedFileRemoval\":{\"settings\":{\"unreferencedDays\":4,\"nonCurrentDays\":10},\"status\":\"enable\"}}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance&space&table&tables", strUrl)
		},
		&PutTableMaintenanceConfigurationRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
			IcebergUnreferencedFileRemoval: &IcebergUnreferencedFileRemoval{
				Settings: &MaintenanceSettings{
					UnreferencedDays: oss.Ptr(int64(4)),
					NonCurrentDays:   oss.Ptr(int64(10)),
				},
				Status: oss.Ptr("enable"),
			},
		},
		func(t *testing.T, o *PutTableMaintenanceConfigurationResult, err error) {
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
   "configuration": { 
      "icebergUnreferencedFileRemoval": {
        "settings": {
          "unreferencedDays":4,
          "nonCurrentDays":10
        },
        "status": "enable"
     }
   },
   "tableARN": "test-arn"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance&space&table&tables", strUrl)
		},
		&GetTableMaintenanceConfigurationRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceConfigurationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Configuration.IcebergUnreferencedFileRemoval.Settings.UnreferencedDays, int64(4))
			assert.Equal(t, *o.Configuration.IcebergUnreferencedFileRemoval.Settings.NonCurrentDays, int64(10))
			assert.Equal(t, *o.Configuration.IcebergUnreferencedFileRemoval.Status, "enable")
			assert.Equal(t, *o.TableARN, "test-arn")
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
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance&space&table&tables", strUrl)
		},
		&GetTableMaintenanceConfigurationRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceConfigurationResult, err error) {
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
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance&space&table&tables", strUrl)
		},
		&GetTableMaintenanceConfigurationRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceConfigurationResult, err error) {
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

var testMockSetTableMaintenanceJobStatusSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SetTableMaintenanceJobStatusRequest
	CheckOutputFn  func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "POST")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			assert.Equal(t, r.Header.Get("x-oss-tables-operation"), "")
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"status\":{\"job\":{\"failureMessage\":\"no message\",\"lastRunTimestamp\":\"2026-02-31T10:56:21.000Z\",\"status\":\"success\"}},\"versionToken\":\"token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance-job-status&space&table&tables", strUrl)
		},
		&SetTableMaintenanceJobStatusRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
			Status: &MaintenanceJobStatus{
				Job: &MaintenanceJob{
					FailureMessage:   oss.Ptr("no message"),
					LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
					Status:           oss.Ptr("success"),
				},
			},
			VersionToken: oss.Ptr("token"),
		},
		func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockSetTableMaintenanceJobStatus_Success(t *testing.T) {
	for _, c := range testMockSetTableMaintenanceJobStatusSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SetTableMaintenanceJobStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSetTableMaintenanceJobStatusErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SetTableMaintenanceJobStatusRequest
	CheckOutputFn  func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error)
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
			assert.Equal(t, r.Method, "POST")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			assert.Equal(t, r.Header.Get("x-oss-tables-operation"), "")
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"status\":{\"job\":{\"failureMessage\":\"no message\",\"lastRunTimestamp\":\"2026-02-31T10:56:21.000Z\",\"status\":\"success\"}},\"versionToken\":\"token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance-job-status&space&table&tables", strUrl)
		},
		&SetTableMaintenanceJobStatusRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
			Status: &MaintenanceJobStatus{
				Job: &MaintenanceJob{
					FailureMessage:   oss.Ptr("no message"),
					LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
					Status:           oss.Ptr("success"),
				},
			},
			VersionToken: oss.Ptr("token"),
		},
		func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error) {
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
			assert.Equal(t, r.Method, "POST")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			assert.Equal(t, r.Header.Get("x-oss-tables-operation"), "")
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"status\":{\"job\":{\"failureMessage\":\"no message\",\"lastRunTimestamp\":\"2026-02-31T10:56:21.000Z\",\"status\":\"success\"}},\"versionToken\":\"token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance-job-status&space&table&tables", strUrl)
		},
		&SetTableMaintenanceJobStatusRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
			Status: &MaintenanceJobStatus{
				Job: &MaintenanceJob{
					FailureMessage:   oss.Ptr("no message"),
					LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
					Status:           oss.Ptr("success"),
				},
			},
			VersionToken: oss.Ptr("token"),
		},
		func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error) {
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

func TestMockSetTableMaintenanceJobStatus_Error(t *testing.T) {
	for _, c := range testMockSetTableMaintenanceJobStatusErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SetTableMaintenanceJobStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSetTableMaintenanceJobStatusByTableArnSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SetTableMaintenanceJobStatusByTableArnRequest
	CheckOutputFn  func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "POST")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			assert.Equal(t, r.Header.Get("x-oss-table-arn"), "acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123")
			assert.Equal(t, r.Header.Get("x-oss-tables-operation"), "")
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"status\":{\"job\":{\"failureMessage\":\"no message\",\"lastRunTimestamp\":\"2026-02-31T10:56:21.000Z\",\"status\":\"success\"}},\"versionToken\":\"token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/?maintenance-job-status&tables", strUrl)
		},
		&SetTableMaintenanceJobStatusByTableArnRequest{
			TableArn: oss.Ptr("acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123"),
			Status: &MaintenanceJobStatus{
				Job: &MaintenanceJob{
					FailureMessage:   oss.Ptr("no message"),
					LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
					Status:           oss.Ptr("success"),
				},
			},
			VersionToken: oss.Ptr("token"),
		},
		func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockSetTableMaintenanceJobStatusByTableArn_Success(t *testing.T) {
	for _, c := range testMockSetTableMaintenanceJobStatusByTableArnSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SetTableMaintenanceJobStatusByTableArn(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSetTableMaintenanceJobStatusByTableArnErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SetTableMaintenanceJobStatusByTableArnRequest
	CheckOutputFn  func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error)
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
			assert.Equal(t, r.Method, "POST")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			assert.Equal(t, r.Header.Get("x-oss-tables-operation"), "")
			assert.Equal(t, r.Header.Get("x-oss-table-arn"), "acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123")
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"status\":{\"job\":{\"failureMessage\":\"no message\",\"lastRunTimestamp\":\"2026-02-31T10:56:21.000Z\",\"status\":\"success\"}},\"versionToken\":\"token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/?maintenance-job-status&tables", strUrl)
		},
		&SetTableMaintenanceJobStatusByTableArnRequest{
			TableArn: oss.Ptr("acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123"),
			Status: &MaintenanceJobStatus{
				Job: &MaintenanceJob{
					FailureMessage:   oss.Ptr("no message"),
					LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
					Status:           oss.Ptr("success"),
				},
			},
			VersionToken: oss.Ptr("token"),
		},
		func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error) {
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
			assert.Equal(t, r.Method, "POST")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			assert.Equal(t, r.Header.Get("x-oss-tables-operation"), "")
			assert.Equal(t, r.Header.Get("x-oss-table-arn"), "acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123")
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(body), "{\"status\":{\"job\":{\"failureMessage\":\"no message\",\"lastRunTimestamp\":\"2026-02-31T10:56:21.000Z\",\"status\":\"success\"}},\"versionToken\":\"token\"}")
			strUrl := sortQuery(r)
			assert.Equal(t, "/?maintenance-job-status&tables", strUrl)
		},
		&SetTableMaintenanceJobStatusByTableArnRequest{
			TableArn: oss.Ptr("acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123"),
			Status: &MaintenanceJobStatus{
				Job: &MaintenanceJob{
					FailureMessage:   oss.Ptr("no message"),
					LastRunTimestamp: oss.Ptr("2026-02-31T10:56:21.000Z"),
					Status:           oss.Ptr("success"),
				},
			},
			VersionToken: oss.Ptr("token"),
		},
		func(t *testing.T, o *SetTableMaintenanceJobStatusResult, err error) {
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

func TestMockSetTableMaintenanceJobStatusByTableArn_Error(t *testing.T) {
	for _, c := range testMockSetTableMaintenanceJobStatusByTableArnErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SetTableMaintenanceJobStatusByTableArn(context.TODO(), c.Request)
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
      "job" : { 
         "failureMessage": "no message",
         "lastRunTimestamp": "2026-02-31T10:56:21.000Z",
         "status": "success"
      }
   },
   "versionToken": "aaa",
   "tableARN": "test-arn"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance-job-status&space&table&tables", strUrl)
		},
		&GetTableMaintenanceJobStatusRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.MaintenanceJobStatus.Job.FailureMessage, "no message")
			assert.Equal(t, *o.MaintenanceJobStatus.Job.Status, "success")
			assert.Equal(t, *o.MaintenanceJobStatus.Job.LastRunTimestamp, "2026-02-31T10:56:21.000Z")
			assert.Equal(t, *o.VersionToken, "aaa")
			assert.Equal(t, *o.TableARN, "test-arn")
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
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance-job-status&space&table&tables", strUrl)
		},
		&GetTableMaintenanceJobStatusRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error) {
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
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?maintenance-job-status&space&table&tables", strUrl)
		},
		&GetTableMaintenanceJobStatusRequest{
			Bucket:    oss.Ptr("bucket"),
			Namespace: oss.Ptr("space"),
			Table:     oss.Ptr("table"),
		},
		func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error) {
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

var testMockGetTableMaintenanceJobStatusByTableArnSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableMaintenanceJobStatusByTableArnRequest
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
      "job" : { 
         "failureMessage": "no message",
         "lastRunTimestamp": "2026-02-31T10:56:21.000Z",
         "status": "success"
      }
   },
   "versionToken": "aaa",
   "tableARN": "test-arn"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			assert.Equal(t, r.Header.Get("x-oss-table-arn"), "acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123")
			strUrl := sortQuery(r)
			assert.Equal(t, "/?maintenance-job-status&tables", strUrl)
		},
		&GetTableMaintenanceJobStatusByTableArnRequest{
			TableArn: oss.Ptr("acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123"),
		},
		func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.MaintenanceJobStatus.Job.FailureMessage, "no message")
			assert.Equal(t, *o.MaintenanceJobStatus.Job.Status, "success")
			assert.Equal(t, *o.MaintenanceJobStatus.Job.LastRunTimestamp, "2026-02-31T10:56:21.000Z")
			assert.Equal(t, *o.VersionToken, "aaa")
			assert.Equal(t, *o.TableARN, "test-arn")
		},
	},
}

func TestMockGetTableMaintenanceJobStatusByTableArn_Success(t *testing.T) {
	for _, c := range testMockGetTableMaintenanceJobStatusByTableArnSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableMaintenanceJobStatusByTableArn(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetTableMaintenanceJobStatusByTableArnErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetTableMaintenanceJobStatusByTableArnRequest
	CheckOutputFn  func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error)
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
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			assert.Equal(t, r.Header.Get("x-oss-table-arn"), "acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123")
			strUrl := sortQuery(r)
			assert.Equal(t, "/?maintenance-job-status&tables", strUrl)
		},
		&GetTableMaintenanceJobStatusByTableArnRequest{
			TableArn: oss.Ptr("acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123"),
		},
		func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error) {
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
			assert.Equal(t, r.Method, "GET")
			assert.Equal(t, r.Header.Get(oss.HTTPHeaderContentType), contentTypeJSON)
			assert.Equal(t, r.Header.Get("x-oss-table-arn"), "acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123")
			strUrl := sortQuery(r)
			assert.Equal(t, "/?maintenance-job-status&tables", strUrl)
		},
		&GetTableMaintenanceJobStatusByTableArnRequest{
			TableArn: oss.Ptr("acs:osstables:cn-hangzhou:123:bucket/oss-demo-bucket/table/table-123"),
		},
		func(t *testing.T, o *GetTableMaintenanceJobStatusResult, err error) {
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

func TestMockGetTableMaintenanceJobStatusByTableArn_Error(t *testing.T) {
	for _, c := range testMockGetTableMaintenanceJobStatusByTableArnErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewTablesClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetTableMaintenanceJobStatusByTableArn(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
