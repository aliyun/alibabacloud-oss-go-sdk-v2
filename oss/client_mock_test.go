package oss

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

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

var testInvokeOperationAnonymousCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Input          *OperationInput
	CheckOutputFn  func(t *testing.T, o *OperationOutput)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "5374A2880232A65C2300****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<ListAllMyBucketsResult>
			<Owner>
				<ID>512**</ID>
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
			</Buckets>				
			</ListAllMyBucketsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		&OperationInput{
			OpName: "ListBuckets",
			Method: "GET",
		},
		func(t *testing.T, o *OperationOutput) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5374A2880232A65C2300****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 15 May 2014 11:18:32 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "654605AA6172673135811AB3",
			"Date":             "Sat, 04 Nov 2023 08:49:46 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<AccessControlPolicy>
				<Owner>
					<ID>12345</ID>
					<DisplayName>12345Name</DisplayName>
				</Owner>
				<AccessControlList>
					<Grant>private</Grant>
				</AccessControlList>
			</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&OperationInput{
			OpName: "GetBucketAcl",
			Bucket: Ptr("bucket"),
			Method: "GET",
			Parameters: map[string]string{
				"acl": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "654605AA6172673135811AB3", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Sat, 04 Nov 2023 08:49:46 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "654605AA6172673135811AB3",
			"Date":             "Sat, 04 Nov 2023 08:49:46 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<AccessControlPolicy>
				<Owner>
					<ID>12345</ID>
					<DisplayName>12345Name</DisplayName>
				</Owner>
				<AccessControlList>
					<Grant>private</Grant>
				</AccessControlList>
			</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/key?acl", r.URL.String())
		},
		&OperationInput{
			OpName: "GetObjectAcl",
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			Method: "GET",
			Parameters: map[string]string{
				"acl": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<InitiateMultipartUploadResult>
				<Bucket>oss-example</Bucket>
				<Key>key+ 123.data</Key>
				<UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
			</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/key%2B%20123/test.data?uploads", r.URL.String())
			assert.Equal(t, "POST", r.Method)
		},
		&OperationInput{
			OpName: "InitiateMultipartUpload",
			Bucket: Ptr("bucket"),
			Key:    Ptr("key+ 123/test.data"),
			Method: "POST",
			Parameters: map[string]string{
				"uploads": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "text/txt",
		},
		[]byte(
			`hello world`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket//subfolder/example.txt?versionId=CAEQNhiBgMDJgZCA0BY%2B123", r.URL.String())
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "Etag1234", r.Header.Get("If-Match"))
		},
		&OperationInput{
			OpName: "GetObject",
			Bucket: Ptr("bucket"),
			Key:    Ptr("/subfolder/example.txt"),
			Method: "GET",
			Headers: map[string]string{
				"If-Match": "Etag1234",
			},
			Parameters: map[string]string{
				"versionId": "CAEQNhiBgMDJgZCA0BY+123",
			},
		},
		func(t *testing.T, o *OperationOutput) {
		},
	},
}

func TestInvokeOperation_Anonymous(t *testing.T) {
	for _, c := range testInvokeOperationAnonymousCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output)

		var fns []func(*Options)
		fns = append(fns, func(c *Options) { c.OpReadWriteTimeout = Ptr(1 * time.Second) })
		output, err = client.InvokeOperation(context.TODO(), c.Input, fns...)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output)
	}
}

var testInvokeOperationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Input          *OperationInput
	CheckOutputFn  func(t *testing.T, o *OperationOutput, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>SignatureDoesNotMatch</Code>
				<Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/test-key.txt", r.URL.String())
		},
		&OperationInput{
			OpName: "PutObject",
			Method: "PUT",
			Bucket: Ptr("bucket"),
			Key:    Ptr("test-key.txt"),
		},
		func(t *testing.T, o *OperationOutput, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "SignatureDoesNotMatch", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
			assert.Contains(t, serr.RequestTarget, "/bucket/test-key.txt")
		},
	},
}

func TestInvokeOperation_Error(t *testing.T) {
	for _, c := range testInvokeOperationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&PutBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *PutBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket", r.URL.String())
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "<CreateBucketConfiguration><StorageClass>Archive</StorageClass><DataRedundancyType>LRS</DataRedundancyType></CreateBucketConfiguration>", string(requestBody))
		},
		&PutBucketRequest{
			Bucket:          Ptr("bucket"),
			Acl:             BucketACLPrivate,
			ResourceGroupId: Ptr("rg-aek27tc********"),
			CreateBucketConfiguration: &CreateBucketConfiguration{
				StorageClass:       StorageClassArchive,
				DataRedundancyType: DataRedundancyLRS,
			},
		},
		func(t *testing.T, o *PutBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucket_Success(t *testing.T) {
	for _, c := range testMockPutBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>SignatureDoesNotMatch</Code>
				<Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&PutBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *PutBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>BucketAlreadyExists</Code>
				<Message>The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.</Message>
				<RequestId>6548A043CA31D****</RequestId>
				<EC>0015-00000104</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&PutBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *PutBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

func TestMockPutBucket_Error(t *testing.T) {
	for _, c := range testMockPutBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListBucketsRequest
	CheckOutputFn  func(t *testing.T, o *ListBucketsResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</ListAllMyBucketsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		nil,
		func(t *testing.T, o *ListBucketsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Owner.DisplayName, "51264")
			assert.Equal(t, *o.Owner.ID, "51264")
			assert.Equal(t, len(o.Buckets), 2)
			assert.Equal(t, *o.Buckets[0].CreationDate, time.Date(2014, time.February, 17, 18, 12, 43, 0, time.UTC))
			assert.Equal(t, *o.Buckets[0].ExtranetEndpoint, "oss-cn-shanghai.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].IntranetEndpoint, "oss-cn-shanghai-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].Name, "app-base-oss")
			assert.Equal(t, *o.Buckets[0].Region, "cn-shanghai")
			assert.Equal(t, *o.Buckets[0].StorageClass, "Standard")

			assert.Equal(t, *o.Buckets[1].CreationDate, time.Date(2014, time.February, 25, 11, 21, 04, 0, time.UTC))
			assert.Equal(t, *o.Buckets[1].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.Buckets[1].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets[1].Name, "mybucket")
			assert.Equal(t, *o.Buckets[1].Region, "cn-hangzhou")
			assert.Equal(t, *o.Buckets[1].StorageClass, "IA")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</ListAllMyBucketsResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/?marker&max-keys=10&prefix=%2F", strUrl)
		},
		&ListBucketsRequest{
			Marker:          Ptr(""),
			MaxKeys:         10,
			Prefix:          Ptr("/"),
			ResourceGroupId: Ptr("rg-aek27tc********"),
		},
		func(t *testing.T, o *ListBucketsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Owner.DisplayName, "ut_test_put_bucket")
			assert.Equal(t, *o.Owner.ID, "ut_test_put_bucket")
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
			assert.Equal(t, *o.Buckets[0].StorageClass, "Standard")
		},
	},
}

func TestMockListBuckets_Success(t *testing.T) {
	for _, c := range testMockListBucketsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListBuckets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListBucketsRequest
	CheckOutputFn  func(t *testing.T, o *ListBucketsResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>InvalidAccessKeyId</Code>
				<Message>The OSS Access Key Id you provided does not exist in our records.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		&ListBucketsRequest{},
		func(t *testing.T, o *ListBucketsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		&ListBucketsRequest{},
		func(t *testing.T, o *ListBucketsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListBuckets fail")
		},
	},
}

func TestMockListBuckets_Error(t *testing.T) {
	for _, c := range testMockListBucketsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListBuckets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketResult, err error)
}{
	{
		204,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&DeleteBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteBucket_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&DeleteBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>BucketNotEmpty</Code>
  <Message>The bucket has objects. Please delete them first.</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000301</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&DeleteBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

func TestMockDeleteBucket_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectsResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?encoding-type=url", r.URL.String())
		},
		&ListObjectsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Empty(t, o.Prefix)
			assert.Equal(t, *o.Name, "examplebucket")
			assert.Empty(t, o.Marker)
			assert.Empty(t, o.Delimiter)
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, len(o.Contents), 4)
			assert.Equal(t, *o.Contents[0].Key, "fun/movie/001.avi")
			assert.Equal(t, *o.Contents[1].LastModified, time.Date(2012, time.February, 24, 8, 43, 27, 0, time.UTC))
			assert.Equal(t, *o.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[3].Type, "Normal")
			assert.Equal(t, o.Contents[0].Size, int64(344606))
			assert.Equal(t, *o.Contents[1].StorageClass, "Standard")
			assert.Equal(t, *o.Contents[2].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[3].Owner.DisplayName, "user-example")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
  <Prefix>fun</Prefix>
  <Marker>test1.txt</Marker>
  <MaxKeys>3</MaxKeys>
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
        <Key>exampleobject3.txt</Key>
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
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Name, "examplebucket")
			assert.Equal(t, *o.Prefix, "fun")
			assert.Equal(t, *o.Marker, "test1.txt")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, o.IsTruncated, true)
			assert.Equal(t, o.MaxKeys, int32(3))
			assert.Equal(t, len(o.Contents), 3)
			assert.Equal(t, *o.Contents[0].Key, "exampleobject1.txt")
			assert.Equal(t, *o.Contents[1].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
			assert.Equal(t, *o.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, o.Contents[1].Size, int64(344606))
			assert.Equal(t, *o.Contents[2].StorageClass, "Standard")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")
			assert.Empty(t, o.Contents[0].RestoreInfo)
			assert.Equal(t, *o.Contents[1].RestoreInfo, "ongoing-request=\"true\"")
			assert.Equal(t, *o.Contents[2].RestoreInfo, "ongoing-request=\"false\", expiry-date=\"Thu, 24 Sep 2020 12:40:33 GMT\"")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
  <Prefix>fun</Prefix>
  <Marker>test1.txt</Marker>
  <MaxKeys>3</MaxKeys>
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
        <Key>exampleobject3.txt</Key>
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
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Name, "examplebucket")
			assert.Equal(t, *o.Prefix, "fun")
			assert.Equal(t, *o.Marker, "test1.txt")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, o.IsTruncated, true)
			assert.Equal(t, o.MaxKeys, int32(3))
			assert.Equal(t, len(o.Contents), 3)
			assert.Equal(t, *o.Contents[0].Key, "exampleobject1.txt")
			assert.Equal(t, *o.Contents[1].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
			assert.Equal(t, *o.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, o.Contents[1].Size, int64(344606))
			assert.Equal(t, *o.Contents[2].StorageClass, "Standard")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")
			assert.Empty(t, o.Contents[0].RestoreInfo)
			assert.Equal(t, *o.Contents[1].RestoreInfo, "ongoing-request=\"true\"")
			assert.Equal(t, *o.Contents[2].RestoreInfo, "ongoing-request=\"false\", expiry-date=\"Thu, 24 Sep 2020 12:40:33 GMT\"")
		},
	},
}

func TestMockListObjects_Success(t *testing.T) {
	for _, c := range testMockListObjectsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?encoding-type=url", r.URL.String())
		},
		&ListObjectsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListObjects fail")
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>Access denied for requester pay bucket</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000703</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "AccessDenied", serr.Code)
			assert.Equal(t, "Access denied for requester pay bucket", serr.Message)
			assert.Equal(t, "0003-00000703", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockListObjects_Error(t *testing.T) {
	for _, c := range testMockListObjectsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectsV2SuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectsRequestV2
	CheckOutputFn  func(t *testing.T, o *ListObjectsResultV2, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
  <Name>examplebucket</Name>
  <Prefix></Prefix>
  <MaxKeys>3</MaxKeys>
  <Delimiter></Delimiter>
  <IsTruncated>false</IsTruncated>
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
        <Key>exampleobject3.txt</Key>
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
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&list-type=2", strUrl)
		},
		&ListObjectsRequestV2{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Empty(t, o.Prefix)
			assert.Equal(t, *o.Name, "examplebucket")
			assert.Empty(t, o.Delimiter)
			assert.Equal(t, o.MaxKeys, int32(3))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, len(o.Contents), 3)
			assert.Equal(t, *o.Contents[0].Key, "exampleobject1.txt")
			assert.Equal(t, *o.Contents[0].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
			assert.Equal(t, *o.Contents[0].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, o.Contents[0].Size, int64(344606))
			assert.Equal(t, *o.Contents[0].StorageClass, "ColdArchive")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")

			assert.Equal(t, *o.Contents[1].RestoreInfo, "ongoing-request=\"true\"")
			assert.Equal(t, *o.Contents[2].RestoreInfo, "ongoing-request=\"false\", expiry-date=\"Thu, 24 Sep 2020 12:40:33 GMT\"")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix>a/</Prefix>
    <MaxKeys>3</MaxKeys>
	<StartAfter>b</StartAfter>
    <Delimiter>/</Delimiter>
    <EncodingType>url</EncodingType>
    <IsTruncated>false</IsTruncated>
  	<Contents>
        <Key>a/b</Key>
        <LastModified>2020-05-18T05:45:47.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
		<Type>Normal</Type>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
		<Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
	</Contents>
  	<Contents>
        <Key>a/b/c</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
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
        <Key>a/b/d</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
        <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>Standard</StorageClass>
        <Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
  </Contents>
	<CommonPrefixes>
        <Prefix>a/b/</Prefix>
    </CommonPrefixes>
    <KeyCount>3</KeyCount>
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
		},
		&ListObjectsRequestV2{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Name, "examplebucket")
			assert.Equal(t, *o.Prefix, "a/")
			assert.Equal(t, *o.StartAfter, "b")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxKeys, int32(3))
			assert.Equal(t, o.KeyCount, 3)
			assert.Equal(t, len(o.Contents), 3)
			assert.Equal(t, *o.Contents[0].Key, "a/b")
			assert.Equal(t, *o.Contents[1].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
			assert.Equal(t, *o.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, o.Contents[1].Size, int64(344606))
			assert.Equal(t, *o.Contents[2].StorageClass, "Standard")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")
			assert.Nil(t, o.Contents[0].RestoreInfo)
			assert.Equal(t, *o.CommonPrefixes[0].Prefix, "a/b/")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix>a/</Prefix>
    <MaxKeys>3</MaxKeys>
	<StartAfter>b</StartAfter>
    <Delimiter>/</Delimiter>
    <EncodingType>url</EncodingType>
    <IsTruncated>false</IsTruncated>
  	<Contents>
        <Key>a/b</Key>
        <LastModified>2020-05-18T05:45:47.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
		<Type>Normal</Type>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
		<Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
	</Contents>
  	<Contents>
        <Key>a/b/c</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
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
        <Key>a/b/d</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
        <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>Standard</StorageClass>
        <Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
  </Contents>
	<CommonPrefixes>
        <Prefix>a/b/</Prefix>
    </CommonPrefixes>
    <KeyCount>3</KeyCount>
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListObjectsRequestV2{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Name, "examplebucket")
			assert.Equal(t, *o.Prefix, "a/")
			assert.Equal(t, *o.StartAfter, "b")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxKeys, int32(3))
			assert.Equal(t, o.KeyCount, 3)
			assert.Equal(t, len(o.Contents), 3)
			assert.Equal(t, *o.Contents[0].Key, "a/b")
			assert.Equal(t, *o.Contents[1].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
			assert.Equal(t, *o.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, o.Contents[1].Size, int64(344606))
			assert.Equal(t, *o.Contents[2].StorageClass, "Standard")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")
			assert.Nil(t, o.Contents[0].RestoreInfo)
			assert.Equal(t, *o.CommonPrefixes[0].Prefix, "a/b/")
		},
	},
}

func TestMockListObjectsV2_Success(t *testing.T) {
	for _, c := range testMockListObjectsV2SuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectsV2(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectsV2ErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectsRequestV2
	CheckOutputFn  func(t *testing.T, o *ListObjectsResultV2, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&list-type=2", strUrl)
		},
		&ListObjectsRequestV2{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
		},
		&ListObjectsRequestV2{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
		},
		&ListObjectsRequestV2{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListObjectsV2 fail")
		},
	},
}

func TestMockListObjectsV2_Error(t *testing.T) {
	for _, c := range testMockListObjectsV2ErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectsV2(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketInfoSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketInfoRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketInfoResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</BucketInfo>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?bucketInfo", r.URL.String())
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.BucketInfo.Name, "oss-example")
			assert.Equal(t, *o.BucketInfo.AccessMonitor, "Enabled")
			assert.Equal(t, *o.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.Location, "oss-cn-hangzhou")
			assert.Equal(t, *o.BucketInfo.StorageClass, "Standard")
			assert.Equal(t, *o.BucketInfo.TransferAcceleration, "Disabled")
			assert.Equal(t, *o.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
			assert.Equal(t, *o.BucketInfo.CrossRegionReplication, "Disabled")
			assert.Equal(t, *o.BucketInfo.ResourceGroupId, "rg-aek27tc********")
			assert.Equal(t, *o.BucketInfo.Owner.ID, "27183473914****")
			assert.Equal(t, *o.BucketInfo.Owner.DisplayName, "username")
			assert.Equal(t, *o.BucketInfo.ACL, "private")
			assert.Equal(t, *o.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
			assert.Equal(t, *o.BucketInfo.BucketPolicy.LogPrefix, "log/")
			assert.Empty(t, *o.BucketInfo.SseRule.KMSMasterKeyID)
			assert.Equal(t, *o.BucketInfo.SseRule.SSEAlgorithm, "KMS")
			assert.Equal(t, *o.BucketInfo.SseRule.KMSDataEncryption, "SM4")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</BucketInfo>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?bucketInfo", r.URL.String())
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.BucketInfo.Name, "oss-example")
			assert.Equal(t, *o.BucketInfo.AccessMonitor, "Enabled")
			assert.Equal(t, *o.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.Location, "oss-cn-hangzhou")
			assert.Equal(t, *o.BucketInfo.StorageClass, "Standard")
			assert.Equal(t, *o.BucketInfo.TransferAcceleration, "Disabled")
			assert.Equal(t, *o.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
			assert.Equal(t, *o.BucketInfo.CrossRegionReplication, "Disabled")
			assert.Equal(t, *o.BucketInfo.ResourceGroupId, "rg-aek27tc********")
			assert.Equal(t, *o.BucketInfo.Owner.ID, "27183473914****")
			assert.Equal(t, *o.BucketInfo.Owner.DisplayName, "username")
			assert.Equal(t, *o.BucketInfo.ACL, "private")
			assert.Equal(t, *o.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
			assert.Equal(t, *o.BucketInfo.BucketPolicy.LogPrefix, "log/")

			assert.Empty(t, o.BucketInfo.SseRule.KMSMasterKeyID)
			assert.Nil(t, o.BucketInfo.SseRule.SSEAlgorithm)
			assert.Nil(t, o.BucketInfo.SseRule.KMSDataEncryption)
		},
	},
}

func TestMockGetBucketInfo_Success(t *testing.T) {
	for _, c := range testMockGetBucketInfoSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketInfo(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketInfoErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketInfoRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketInfoResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?bucketInfo", r.URL.String())
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?bucketInfo", strUrl)
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?bucketInfo", strUrl)
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketInfo fail")
		},
	},
}

func TestMockGetBucketInfo_Error(t *testing.T) {
	for _, c := range testMockGetBucketInfoErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketInfo(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketLocationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLocationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLocationResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<LocationConstraint>oss-cn-hangzhou</LocationConstraint>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?location", r.URL.String())
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.LocationConstraint, "oss-cn-hangzhou")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<LocationConstraint>oss-cn-chengdu</LocationConstraint>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?location", r.URL.String())
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.LocationConstraint, "oss-cn-chengdu")
		},
	},
}

func TestMockGetBucketLocation_Success(t *testing.T) {
	for _, c := range testMockGetBucketLocationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketLocationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLocationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLocationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?location", r.URL.String())
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?location", strUrl)
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?location", strUrl)
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketLocation fail")
		},
	},
}

func TestMockGetBucketLocation_Error(t *testing.T) {
	for _, c := range testMockGetBucketLocationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketStatSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketStatRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketStatResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</BucketStat>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?stat", r.URL.String())
		},
		&GetBucketStatRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketStatResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, int64(1600), o.Storage)
			assert.Equal(t, int64(230), o.ObjectCount)
			assert.Equal(t, int64(40), o.MultipartUploadCount)
			assert.Equal(t, int64(4), o.LiveChannelCount)
			assert.Equal(t, int64(1643341269), o.LastModifiedTime)
			assert.Equal(t, int64(430), o.StandardStorage)
			assert.Equal(t, int64(66), o.StandardObjectCount)
			assert.Equal(t, int64(2359296), o.InfrequentAccessStorage)
			assert.Equal(t, int64(360), o.InfrequentAccessRealStorage)
			assert.Equal(t, int64(54), o.InfrequentAccessObjectCount)
			assert.Equal(t, int64(2949120), o.ArchiveStorage)
			assert.Equal(t, int64(450), o.ArchiveRealStorage)
			assert.Equal(t, int64(74), o.ArchiveObjectCount)
			assert.Equal(t, int64(2359296), o.ColdArchiveStorage)
			assert.Equal(t, int64(360), o.ColdArchiveRealStorage)
			assert.Equal(t, int64(36), o.ColdArchiveObjectCount)
		},
	},
}

func TestMockGetBucketStat_Success(t *testing.T) {
	for _, c := range testMockGetBucketStatSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketStat(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketStatErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketStatRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketStatResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?stat", r.URL.String())
		},
		&GetBucketStatRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketStatResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?stat", strUrl)
		},
		&GetBucketStatRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketStatResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?stat", strUrl)
		},
		&GetBucketStatRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketStatResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketStat fail")
		},
	},
}

func TestMockGetBucketStat_Error(t *testing.T) {
	for _, c := range testMockGetBucketStatErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketStat(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketAclRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
			assert.Equal(t, string(BucketACLPublicRead), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPublicRead,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
			assert.Equal(t, string(BucketACLPrivate), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPrivate,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
			assert.Equal(t, string(BucketACLPublicReadWrite), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPublicReadWrite,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketAcl_Success(t *testing.T) {
	for _, c := range testMockPutBucketAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketAclRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketAclResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
			assert.Equal(t, string(BucketACLPrivate), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPrivate,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?acl", strUrl)
			assert.Equal(t, string(BucketACLPrivate), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPrivate,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?acl", strUrl)
			assert.Equal(t, string(BucketACLPrivate), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPrivate,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutBucketAcl fail")
		},
	},
}

func TestMockPutBucketAcl_Error(t *testing.T) {
	for _, c := range testMockPutBucketAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketAclRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012****</ID>
        <DisplayName>user_example</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read", *o.ACL)
			assert.Equal(t, "0022012****", *o.Owner.ID)
			assert.Equal(t, "user_example", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read-write</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read-write", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>private</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "private", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
		},
	},
}

func TestMockGetBucketAcl_Success(t *testing.T) {
	for _, c := range testMockGetBucketAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketAclRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketAclResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?acl", strUrl)
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketAcl fail")
		},
	},
}

func TestMockGetBucketAcl_Error(t *testing.T) {
	for _, c := range testMockGetBucketAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Nil(t, o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",

			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.NotNil(t, r.Header.Get("x-oss-callback"))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",

			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
		},
		[]byte(`{"filename":"object","size":"6","mimeType":""}`),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.NotNil(t, r.Header.Get("x-oss-callback"))
		},
		&PutObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
			Body:     strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			jsonData, err := json.Marshal(o.CallbackResult)
			assert.Nil(t, err)
			assert.NotEmpty(t, string(jsonData))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&PutObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			TrafficLimit: int64(100 * 1024 * 8),
			Body:         strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object.txt", r.URL.String())
			assert.Equal(t, "text/plain", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object.txt"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object.txt", r.URL.String())
			assert.Equal(t, "my-content-type", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket:      Ptr("bucket"),
			Key:         Ptr("object.txt"),
			Body:        strings.NewReader("hi oss"),
			ContentType: Ptr("my-content-type"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object.txt", r.URL.String())
			assert.Equal(t, "my-content-type", r.Header.Get(HTTPHeaderContentType))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&PutObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object.txt"),
			Body:         strings.NewReader("hi oss"),
			ContentType:  Ptr("my-content-type"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
		},
	},
}

func TestMockPutObject_Success(t *testing.T) {
	for _, c := range testMockPutObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectDisableDetectMimeTypeCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockPutObject_DisableDetectMimeType(t *testing.T) {
	for _, c := range testMockPutObjectDisableDetectMimeTypeCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureAutoDetectMimeType
			})
		assert.NotNil(t, c)

		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectWithProgressCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockPutObject_Progress(t *testing.T) {
	for _, c := range testMockPutObjectWithProgressCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		n := int64(0)
		c.Request.ProgressFn = func(increment, transferred, total int64) {
			n = transferred
			//fmt.Printf("got transferred:%v, total:%v\n", transferred, total)
		}
		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, int64(len("hi oss")), n)
	}
}

var testMockPutObjectWithCrcDisableCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "6707180448768400016",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "6707180448768400016")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockPutObject_DisableCRC64(t *testing.T) {
	//Disable
	for _, c := range testMockPutObjectWithCrcDisableCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureEnableCRC64CheckUpload
			})
		assert.NotNil(t, c)
		n := int64(0)
		c.Request.ProgressFn = func(increment, transferred, total int64) {
			n = transferred
		}
		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, int64(len("hi oss")), n)

		//Enable Got Fail
		client = NewClient(cfg)
		assert.NotNil(t, c)
		n = int64(0)
		c.Request.ProgressFn = func(increment, transferred, total int64) {
			n = transferred
		}
		c.Request.Body = strings.NewReader("hi oss")
		_, err = client.PutObject(context.TODO(), c.Request)
		assert.NotNil(t, err)
		assert.Equal(t, int64(len("hi oss")), n)
		assert.Contains(t, err.Error(), "crc is inconsistent, client 8707180448768400016, server 6707180448768400016")
	}
}

var testMockPutObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		203,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "5C3D9175B6FC201293AD****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>CallbackFailed</Code>
  <Message>Error status : 301.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0007-00000203</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0007-00000203</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Body:     strings.NewReader("hi oss"),
			Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(203), serr.StatusCode)
			assert.Equal(t, "CallbackFailed", serr.Code)
			assert.Equal(t, "Error status : 301.", serr.Message)
			assert.Equal(t, "0007-00000203", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
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
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Body:     strings.NewReader("hi oss"),
			Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutObject fail")
		},
	},
}

func TestMockPutObject_Error(t *testing.T) {
	for _, c := range testMockPutObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "316181249502703****",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(`hi oss,this is a demo!`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "316181249502703****")
			content, err := io.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(content), "hi oss,this is a demo!")
			assert.Nil(t, o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"5B3C1A2E05E1B002CC607C****\"",
			"Content-Length":                      "344606",
			"Last-Modified":                       "Fri, 24 Feb 2012 06:07:48 GMT",
			"x-oss-object-type":                   "Normal",
			"Accept-Ranges":                       "bytes",
			"Content-disposition":                 "attachment; filename=testing.txt",
			"Cache-control":                       "no-cache",
			"X-Oss-Storage-Class":                 "Standard",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-tagging-count":                 "2",
			"Content-MD5":                         "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-hash-crc64ecma":                "870718044876840****",
			"x-oss-meta-name":                     "demo",
			"x-oss-meta-email":                    "demo@aliyun.com",
		},
		[]byte(`hi oss`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
			assert.Equal(t, *o.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
			assert.Equal(t, *o.ContentType, "text")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.StorageClass, "Standard")
			content, err := io.ReadAll(o.Body)
			assert.Equal(t, string(content), "hi oss")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, o.TaggingCount, int32(2))
			assert.Equal(t, o.Metadata["name"], "demo")
			assert.Equal(t, o.Metadata["email"], "demo@aliyun.com")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"5B3C1A2E05E1B002CC607C****\"",
			"Content-Length":                      "344606",
			"Last-Modified":                       "Fri, 24 Feb 2012 06:07:48 GMT",
			"x-oss-object-type":                   "Normal",
			"Accept-Ranges":                       "bytes",
			"Content-disposition":                 "attachment; filename=testing.txt",
			"Cache-control":                       "no-cache",
			"X-Oss-Storage-Class":                 "Standard",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-tagging-count":                 "2",
			"Content-MD5":                         "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-hash-crc64ecma":                "870718044876840****",
			"x-oss-meta-name":                     "demo",
			"x-oss-meta-email":                    "demo@aliyun.com",
		},
		[]byte(`hi oss`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&GetObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			TrafficLimit: int64(100 * 1024 * 8),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
			assert.Equal(t, *o.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
			assert.Equal(t, *o.ContentType, "text")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.StorageClass, "Standard")
			content, err := io.ReadAll(o.Body)
			assert.Equal(t, string(content), "hi oss")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, o.TaggingCount, int32(2))
			assert.Equal(t, o.Metadata["name"], "demo")
			assert.Equal(t, o.Metadata["email"], "demo@aliyun.com")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":         "image/jpeg",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                 "\"5B3C1A2E05E1B002CC607C****\"",
			"Content-Length":       "344606",
			"Last-Modified":        "Fri, 24 Feb 2012 06:07:48 GMT",
			"x-oss-object-type":    "Normal",
			"X-Oss-Storage-Class":  "Standard",
			"x-oss-hash-crc64ecma": "870718044876840****",
		},
		[]byte(`hi oss`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?x-oss-process=image%2Fresize%2Cm_fixed%2Cw_100%2Ch_100", r.URL.String())
		},
		&GetObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr("image/resize,m_fixed,w_100,h_100"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
			assert.Equal(t, *o.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
			assert.Equal(t, *o.ContentType, "image/jpeg")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.StorageClass, "Standard")
			content, err := io.ReadAll(o.Body)
			assert.Equal(t, string(content), "hi oss")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":         "image/jpeg",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                 "\"5B3C1A2E05E1B002CC607C****\"",
			"Content-Length":       "344606",
			"Last-Modified":        "Fri, 24 Feb 2012 06:07:48 GMT",
			"x-oss-object-type":    "Normal",
			"X-Oss-Storage-Class":  "Standard",
			"x-oss-hash-crc64ecma": "870718044876840****",
		},
		[]byte(`hi oss`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&GetObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
			assert.Equal(t, *o.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
			assert.Equal(t, *o.ContentType, "image/jpeg")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.StorageClass, "Standard")
			content, err := io.ReadAll(o.Body)
			assert.Equal(t, string(content), "hi oss")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
		},
	},
}

func TestMockGetObject_Success(t *testing.T) {
	for _, c := range testMockGetObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockGetObject_Error(t *testing.T) {
	for _, c := range testMockGetObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCopyObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CopyObjectRequest
	CheckOutputFn  func(t *testing.T, o *CopyObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":                 "application/xml",
			"x-oss-request-id":             "534B371674E88A4D8906****",
			"Date":                         "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                         "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-hash-crc64ecma":         "870718044876840****",
			"x-oss-copy-source-version-id": "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****",
			"x-oss-version-id":             "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object?versionId=CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			SourceKey:       Ptr("copy-object"),
			SourceVersionId: Ptr("CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":         "text",
			"ETag":                 "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-hash-crc64ecma": "870718044876841****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			SourceKey: Ptr("copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-hash-crc64ecma":                "870718044876841****",
			"x-oss-copy-source-version-id":        "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object?versionId=CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			SourceKey:       Ptr("copy-object"),
			SourceVersionId: Ptr("CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-hash-crc64ecma":                "870718044876841****",
			"x-oss-copy-source-version-id":        "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object?versionId=CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****", r.Header.Get("x-oss-copy-source"))
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&CopyObjectRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			SourceKey:       Ptr("copy-object"),
			TrafficLimit:    int64(100 * 1024 * 8),
			SourceVersionId: Ptr("CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-hash-crc64ecma":                "870718044876841****",
			"x-oss-copy-source-version-id":        "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object?versionId=CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****", r.Header.Get("x-oss-copy-source"))
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&CopyObjectRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			SourceKey:       Ptr("copy-object"),
			TrafficLimit:    int64(100 * 1024 * 8),
			SourceVersionId: Ptr("CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****"),
			RequestPayer:    Ptr("requester"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
}

func TestMockCopyObject_Success(t *testing.T) {
	for _, c := range testMockCopyObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CopyObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCopyObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CopyObjectRequest
	CheckOutputFn  func(t *testing.T, o *CopyObjectResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			SourceKey: Ptr("copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			SourceKey: Ptr("copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			SourceKey: Ptr("copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute CopyObject fail")
		},
	},
}

func TestMockCopyObject_Error(t *testing.T) {
	for _, c := range testMockCopyObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CopyObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "1717",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			Body:     strings.NewReader("hi oss,append object"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Nil(t, o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":           "6551DBCF4311A7303980****",
			"Date":                       "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":           "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position": "0",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(100)),
			Body:     strings.NewReader("hi oss,append object,this is a demo"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(0))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":                    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position":          "1717",
			"x-oss-hash-crc64ecma":                "1474161709526656****",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "12f8711f-90df-4e0d-903d-ab972b0f****")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
		},
		&AppendObjectRequest{
			Bucket:                   Ptr("bucket"),
			Key:                      Ptr("object"),
			Position:                 Ptr(int64(100)),
			Body:                     strings.NewReader("hi oss,append object,this is a demo"),
			ServerSideEncryption:     Ptr("KMS"),
			ServerSideDataEncryption: Ptr("SM4"),
			SSEKMSKeyId:              Ptr("12f8711f-90df-4e0d-903d-ab972b0f****"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":                    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position":          "1717",
			"x-oss-hash-crc64ecma":                "1474161709526656****",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "12f8711f-90df-4e0d-903d-ab972b0f****")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&AppendObjectRequest{
			Bucket:                   Ptr("bucket"),
			Key:                      Ptr("object"),
			Position:                 Ptr(int64(100)),
			Body:                     strings.NewReader("hi oss,append object,this is a demo"),
			ServerSideEncryption:     Ptr("KMS"),
			ServerSideDataEncryption: Ptr("SM4"),
			SSEKMSKeyId:              Ptr("12f8711f-90df-4e0d-903d-ab972b0f****"),
			TrafficLimit:             int64(100 * 1024 * 8),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":                    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position":          "1717",
			"x-oss-hash-crc64ecma":                "1474161709526656****",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "12f8711f-90df-4e0d-903d-ab972b0f****")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&AppendObjectRequest{
			Bucket:                   Ptr("bucket"),
			Key:                      Ptr("object"),
			Position:                 Ptr(int64(100)),
			Body:                     strings.NewReader("hi oss,append object,this is a demo"),
			ServerSideEncryption:     Ptr("KMS"),
			ServerSideDataEncryption: Ptr("SM4"),
			SSEKMSKeyId:              Ptr("12f8711f-90df-4e0d-903d-ab972b0f****"),
			TrafficLimit:             int64(100 * 1024 * 8),
			RequestPayer:             Ptr("requester"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
		},
	},
}

func TestMockAppendObject_Success(t *testing.T) {
	for _, c := range testMockAppendObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectDisableDetectMimeTypeCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "1717",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			Body:     strings.NewReader("hi oss,append object"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockAppendObject_DisableDetectMimeType(t *testing.T) {
	for _, c := range testMockAppendObjectDisableDetectMimeTypeCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureAutoDetectMimeType
			})
		assert.NotNil(t, c)

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectWithProgressCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "1717",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			Body:     strings.NewReader("hi oss,append object"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockAppendObject_Progress(t *testing.T) {
	for _, c := range testMockAppendObjectWithProgressCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		n := int64(0)
		c.Request.ProgressFn = func(increment, transferred, total int64) {
			n = transferred
		}
		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, int64(len("hi oss,append object")), n)
	}
}

var testMockAppendObjectCRC64Cases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "20",
			"x-oss-hash-crc64ecma":       "2313496259928504459",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:        Ptr("bucket"),
			Key:           Ptr("object"),
			Position:      Ptr(int64(0)),
			Body:          strings.NewReader("hi oss,append object"),
			InitHashCRC64: Ptr("0"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(20))
			assert.Equal(t, *o.HashCRC64, "2313496259928504459")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":           "6551DBCF4311A7303980****",
			"Date":                       "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-next-append-position": "35",
			"x-oss-hash-crc64ecma":       "8586970469916596321",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader(",this is a demo"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=20", strUrl)
		},
		&AppendObjectRequest{
			Bucket:        Ptr("bucket"),
			Key:           Ptr("object"),
			Position:      Ptr(int64(20)),
			Body:          strings.NewReader(",this is a demo"),
			InitHashCRC64: Ptr("2313496259928504459"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.HashCRC64, "8586970469916596321")
			assert.Equal(t, o.NextPosition, int64(35))
		},
	},
}

func TestMockAppendObject_CRC64(t *testing.T) {
	for _, c := range testMockAppendObjectCRC64Cases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectDisableCRC64Cases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "20",
			"x-oss-hash-crc64ecma":       "4313496259928504459",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:        Ptr("bucket"),
			Key:           Ptr("object"),
			Position:      Ptr(int64(0)),
			Body:          strings.NewReader("hi oss,append object"),
			InitHashCRC64: Ptr("0"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(20))
			assert.Equal(t, *o.HashCRC64, "4313496259928504459")
		},
	},
}

func TestMockAppendObject_DisableCRC64(t *testing.T) {
	for _, c := range testMockAppendObjectDisableCRC64Cases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		//Enable, meets error
		client := NewClient(cfg)
		assert.NotNil(t, c)

		_, err := client.AppendObject(context.TODO(), c.Request)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "crc is inconsistent, client 2313496259928504459, server 4313496259928504459")

		// Disable, no error
		client = NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureEnableCRC64CheckUpload
			})
		assert.NotNil(t, c)
		c.Request.Body = strings.NewReader("hi oss,append object")
		output, err := client.AppendObject(context.TODO(), c.Request)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output, err)

		// don't set initCRC, no error
		client = NewClient(cfg)
		assert.NotNil(t, c)
		c.Request.InitHashCRC64 = nil
		c.Request.Body = strings.NewReader("hi oss,append object")
		output, err = client.AppendObject(context.TODO(), c.Request)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			strUrl := sortQuery(r)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(100)),
			Body:     strings.NewReader("hi oss,append object"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		409,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>PositionNotEqualToLength</Code>
  <Message>Position is not equal to file length</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0026-00000016</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			Body:     strings.NewReader("hi oss,append object,this is a demo"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "PositionNotEqualToLength", serr.Code)
			assert.Equal(t, "Position is not equal to file length", serr.Message)
			assert.Equal(t, "0026-00000016", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockAppendObject_Error(t *testing.T) {
	for _, c := range testMockAppendObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteObjectRequest
	CheckOutputFn  func(t *testing.T, o *DeleteObjectResult, err error)
}{
	{
		204,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&DeleteObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Nil(t, o.VersionId)
			assert.False(t, o.DeleteMarker)
		},
	},
	{
		204,
		map[string]string{
			"x-oss-request-id":    "6551DBCF4311A7303980****",
			"Date":                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-delete-marker": "true",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&DeleteObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.True(t, o.DeleteMarker)
		},
	},
	{
		204,
		map[string]string{
			"x-oss-request-id":    "6551DBCF4311A7303980****",
			"Date":                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-delete-marker": "true",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&DeleteObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.True(t, o.DeleteMarker)
		},
	},
}

func TestMockDeleteObject_Success(t *testing.T) {
	for _, c := range testMockDeleteObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteObjectRequest
	CheckOutputFn  func(t *testing.T, o *DeleteObjectResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&DeleteObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

func TestMockDeleteObject_Error(t *testing.T) {
	for _, c := range testMockDeleteObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteMultipleObjectsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteMultipleObjectsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteMultipleObjectsResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>true</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
			Quiet:   true,
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Nil(t, o.DeletedObjects)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		}, []byte(`<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), ("<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>"))
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.DeletedObjects[0].Key, "key1.txt")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
			assert.Equal(t, *o.DeletedObjects[1].Key, "key2.txt")
			assert.Equal(t, o.DeletedObjects[1].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
			assert.Nil(t, o.DeletedObjects[1].VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			Objects:      []DeleteObject{{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}, {Key: Ptr("key2.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****")}},
			EncodingType: Ptr("url"),
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 2)
			assert.Equal(t, *o.DeletedObjects[0].Key, "key1.txt")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
			assert.Equal(t, *o.DeletedObjects[1].Key, "key2.txt")
			assert.Equal(t, o.DeletedObjects[1].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
			assert.Nil(t, o.DeletedObjects[1].VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>go-sdk-v1%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>go-sdk-v1&#x01;&#x02;&#x03;&#x04;&#x05;&#x06;&#x07;&#x08;&#x9;&#xA;&#x0B;&#x0C;&#xD;&#x0E;&#x0F;&#x10;&#x11;&#x12;&#x13;&#x14;&#x15;&#x16;&#x17;&#x18;&#x19;&#x1A;&#x1B;&#x1C;&#x1D;&#x1E;&#x1F;</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			Objects:      []DeleteObject{{Key: Ptr("go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}},
			EncodingType: Ptr("url"),
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 1)
			assert.Equal(t, *o.DeletedObjects[0].Key, "go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>go-sdk-v1%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>go-sdk-v1&#x01;&#x02;&#x03;&#x04;&#x05;&#x06;&#x07;&#x08;&#x9;&#xA;&#x0B;&#x0C;&#xD;&#x0E;&#x0F;&#x10;&#x11;&#x12;&#x13;&#x14;&#x15;&#x16;&#x17;&#x18;&#x19;&#x1A;&#x1B;&#x1C;&#x1D;&#x1E;&#x1F;</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object></Delete>")
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			Objects:      []DeleteObject{{Key: Ptr("go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}},
			EncodingType: Ptr("url"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 1)
			assert.Equal(t, *o.DeletedObjects[0].Key, "go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
		},
	},
}

func TestMockDeleteMultipleObjects_Success(t *testing.T) {
	for _, c := range testMockDeleteMultipleObjectsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteMultipleObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteMultipleObjectsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteMultipleObjectsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteMultipleObjectsResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}, {Key: Ptr("key2.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MalformedXML</Code>
  <Message>The XML you provided was not well-formed or did not validate against our published schema.</Message>
  <RequestId>6555AC764311A73931E0****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <ErrorDetail>the root node is not named Delete.</ErrorDetail>
  <EC>0016-00000608</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000608</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "MalformedXML", serr.Code)
			assert.Equal(t, "The XML you provided was not well-formed or did not validate against our published schema.", serr.Message)
			assert.Equal(t, "0016-00000608", serr.EC)
			assert.Equal(t, "6555AC764311A73931E0****", serr.RequestID)
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute DeleteMultipleObjects fail")
		},
	},
}

func TestMockDeleteMultipleObjects_Error(t *testing.T) {
	for _, c := range testMockDeleteMultipleObjectsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteMultipleObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockHeadObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *HeadObjectRequest
	CheckOutputFn  func(t *testing.T, o *HeadObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":    "6555A936CA31DC333143****",
			"Date":                "Thu, 16 Nov 2023 05:31:34 GMT",
			"x-oss-object-type":   "Normal",
			"x-oss-storage-class": "Archive",
			"Last-Modified":       "Fri, 24 Feb 2018 09:41:56 GMT",
			"Content-Length":      "344606",
			"Content-Type":        "image/jpg",
			"ETag":                "\"fba9dede5f27731c9771645a3986****\"",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"fba9dede5f27731c9771645a3986****\"")
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ContentType, "image/jpg")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":    "5CAC3B40B7AEADE01700****",
			"Date":                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":        "text/xml",
			"x-oss-object-type":   "Normal",
			"x-oss-storage-class": "Archive",
			"Last-Modified":       "Fri, 24 Feb 2023 09:41:56 GMT",
			"Content-Length":      "481827",
			"ETag":                "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-version-Id":    "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "text/xml")
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":    "534B371674E88A4D8906****",
			"Date":                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":        "image/jpg",
			"x-oss-object-type":   "Normal",
			"x-oss-restore":       "ongoing-request=\"true\"",
			"x-oss-storage-class": "Archive",
			"Last-Modified":       "Fri, 24 Feb 2023 09:41:59 GMT",
			"Content-Length":      "481827",
			"ETag":                "\"A082B659EF78733A5A042FA253B1****\"",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "image/jpg")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.Restore, "ongoing-request=\"true\"")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":                    "534B371674E88A4D8906****",
			"Date":                                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":                        "image/jpg",
			"x-oss-object-type":                   "Normal",
			"x-oss-restore":                       "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"",
			"x-oss-storage-class":                 "Archive",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "9468da86-3509-4f8d-a61e-6eab1eac****",
			"Content-Length":                      "481827",
			"ETag":                                "\"A082B659EF78733A5A042FA253B1****\"",
			"Last-Modified":                       "Fri, 24 Feb 2023 09:41:59 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "image/jpg")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.Restore, "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.SSEKMSKeyId, "9468da86-3509-4f8d-a61e-6eab1eac****")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":                    "534B371674E88A4D8906****",
			"Date":                                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":                        "image/jpg",
			"x-oss-object-type":                   "Normal",
			"x-oss-restore":                       "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"",
			"x-oss-storage-class":                 "Archive",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "9468da86-3509-4f8d-a61e-6eab1eac****",
			"Content-Length":                      "481827",
			"ETag":                                "\"A082B659EF78733A5A042FA253B1****\"",
			"Last-Modified":                       "Fri, 24 Feb 2023 09:41:59 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&HeadObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "image/jpg")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.Restore, "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.SSEKMSKeyId, "9468da86-3509-4f8d-a61e-6eab1eac****")
		},
	},
}

func TestMockHeadObject_Success(t *testing.T) {
	for _, c := range testMockHeadObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.HeadObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockHeadObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *HeadObjectRequest
	CheckOutputFn  func(t *testing.T, o *HeadObjectResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556E3AED11E553933CCDEDF",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPk5vU3VjaEtleTwvQ29kZT4KICA8TWVzc2FnZT5UaGUgc3BlY2lmaWVkIGtleSBkb2VzIG5vdCBleGlzdC48L01lc3NhZ2U+CiAgPFJlcXVlc3RJZD42NTU2RTNBRUQxMUU1NTM5MzNDQ0RFREY8L1JlcXVlc3RJZD4KICA8SG9zdElkPmRlbW8td2Fsa2VyLTY5NjEub3NzLWNuLWhhbmd6aG91LmFsaXl1bmNzLmNvbTwvSG9zdElkPgogIDxLZXk+d2Fsa2VyMmFzZGFzZGFzZC50eHQ8L0tleT4KICA8RUM+MDAyNi0wMDAwMDAwMTwvRUM+CiAgPFJlY29tbWVuZERvYz5odHRwczovL2FwaS5hbGl5dW4uY29tL3Ryb3VibGVzaG9vdD9xPTAwMjYtMDAwMDAwMDE8L1JlY29tbWVuZERvYz4KPC9FcnJvcj4K",
			"x-oss-ec":         "0026-00000001",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "6556E3AED11E553933CCDEDF", serr.RequestID)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "0026-00000001", serr.EC)
		},
	},
	{
		304,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(304), serr.StatusCode)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556FF5BD11E5536368607E8",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPkludmFsaWRUYXJnZXRUeXBlPC9Db2RlPgogIDxNZXNzYWdlPlRoZSBzeW1ib2xpYydzIHRhcmdldCBmaWxlIHR5cGUgaXMgaW52YWxpZDwvTWVzc2FnZT4KICA8UmVxdWVzdElkPjY1NTZGRjVCRDExRTU1MzYzNjg2MDdFODwvUmVxdWVzdElkPgogIDxIb3N0SWQ+ZGVtby13YWxrZXItNjk2MS5vc3MtY24taGFuZ3pob3UuYWxpeXVuY3MuY29tPC9Ib3N0SWQ+CiAgPEVDPjAwMjYtMDAwMDAwMTE8L0VDPgogIDxSZWNvbW1lbmREb2M+aHR0cHM6Ly9hcGkuYWxpeXVuLmNvbS90cm91Ymxlc2hvb3Q/cT0wMDI2LTAwMDAwMDExPC9SZWNvbW1lbmREb2M+CjwvRXJyb3I+Cg==",
			"x-oss-ec":         "0026-00000011",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidTargetType", serr.Code)
			assert.Equal(t, "6556FF5BD11E5536368607E8", serr.RequestID)
			assert.Equal(t, "The symbolic's target file type is invalid", serr.Message)
			assert.Equal(t, "0026-00000011", serr.EC)
		},
	},
}

func TestMockHeadObject_Error(t *testing.T) {
	for _, c := range testMockHeadObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.HeadObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectMetaSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectMetaRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectMetaResult, err error)
}{
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "6555A936CA31DC333143****",
			"Date":             "Thu, 16 Nov 2023 05:31:34 GMT",
			"Last-Modified":    "Fri, 24 Feb 2018 09:41:56 GMT",
			"Content-Length":   "344606",
			"ETag":             "\"fba9dede5f27731c9771645a3986****\"",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"fba9dede5f27731c9771645a3986****\"")
			assert.Equal(t, *o.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(344606))
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "5CAC3B40B7AEADE01700****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
			"Last-Modified":    "Fri, 24 Feb 2023 09:41:56 GMT",
			"Content-Length":   "481827",
			"ETag":             "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-version-Id": "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?objectMeta&versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", strUrl)
		},
		&GetObjectMetaRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":       "534B371674E88A4D8906****",
			"Date":                   "Tue, 04 Dec 2018 15:56:38 GMT",
			"Last-Modified":          "Fri, 24 Feb 2023 09:41:59 GMT",
			"Content-Length":         "481827",
			"ETag":                   "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-last-access-time": "Thu, 14 Oct 2021 11:49:05 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.LastAccessTime, time.Date(2021, time.October, 14, 11, 49, 05, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":       "534B371674E88A4D8906****",
			"Date":                   "Tue, 04 Dec 2018 15:56:38 GMT",
			"Last-Modified":          "Fri, 24 Feb 2023 09:41:59 GMT",
			"Content-Length":         "481827",
			"ETag":                   "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-last-access-time": "Thu, 14 Oct 2021 11:49:05 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&GetObjectMetaRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.LastAccessTime, time.Date(2021, time.October, 14, 11, 49, 05, 0, time.UTC))
		},
	},
}

func TestMockGetObjectMeta_Success(t *testing.T) {
	for _, c := range testMockGetObjectMetaSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectMetaErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectMetaRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectMetaResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556E3AED11E553933CCDEDF",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPk5vU3VjaEtleTwvQ29kZT4KICA8TWVzc2FnZT5UaGUgc3BlY2lmaWVkIGtleSBkb2VzIG5vdCBleGlzdC48L01lc3NhZ2U+CiAgPFJlcXVlc3RJZD42NTU2RTNBRUQxMUU1NTM5MzNDQ0RFREY8L1JlcXVlc3RJZD4KICA8SG9zdElkPmRlbW8td2Fsa2VyLTY5NjEub3NzLWNuLWhhbmd6aG91LmFsaXl1bmNzLmNvbTwvSG9zdElkPgogIDxLZXk+d2Fsa2VyMmFzZGFzZGFzZC50eHQ8L0tleT4KICA8RUM+MDAyNi0wMDAwMDAwMTwvRUM+CiAgPFJlY29tbWVuZERvYz5odHRwczovL2FwaS5hbGl5dW4uY29tL3Ryb3VibGVzaG9vdD9xPTAwMjYtMDAwMDAwMDE8L1JlY29tbWVuZERvYz4KPC9FcnJvcj4K",
			"x-oss-ec":         "0026-00000001",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "6556E3AED11E553933CCDEDF", serr.RequestID)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "0026-00000001", serr.EC)
		},
	},
	{
		304,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(304), serr.StatusCode)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556FF5BD11E5536368607E8",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPkludmFsaWRUYXJnZXRUeXBlPC9Db2RlPgogIDxNZXNzYWdlPlRoZSBzeW1ib2xpYydzIHRhcmdldCBmaWxlIHR5cGUgaXMgaW52YWxpZDwvTWVzc2FnZT4KICA8UmVxdWVzdElkPjY1NTZGRjVCRDExRTU1MzYzNjg2MDdFODwvUmVxdWVzdElkPgogIDxIb3N0SWQ+ZGVtby13YWxrZXItNjk2MS5vc3MtY24taGFuZ3pob3UuYWxpeXVuY3MuY29tPC9Ib3N0SWQ+CiAgPEVDPjAwMjYtMDAwMDAwMTE8L0VDPgogIDxSZWNvbW1lbmREb2M+aHR0cHM6Ly9hcGkuYWxpeXVuLmNvbS90cm91Ymxlc2hvb3Q/cT0wMDI2LTAwMDAwMDExPC9SZWNvbW1lbmREb2M+CjwvRXJyb3I+Cg==",
			"x-oss-ec":         "0026-00000011",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidTargetType", serr.Code)
			assert.Equal(t, "6556FF5BD11E5536368607E8", serr.RequestID)
			assert.Equal(t, "The symbolic's target file type is invalid", serr.Message)
			assert.Equal(t, "0026-00000011", serr.EC)
		},
	},
}

func TestMockGetObjectMeta_Error(t *testing.T) {
	for _, c := range testMockGetObjectMetaErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockRestoreObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RestoreObjectRequest
	CheckOutputFn  func(t *testing.T, o *RestoreObjectResult, err error)
}{
	{
		202,
		map[string]string{
			"X-Oss-Request-Id": "6555A936CA31DC333143****",
			"Date":             "Thu, 16 Nov 2023 05:31:34 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 202, o.StatusCode)
			assert.Equal(t, "202 Accepted", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "5CAC3B40B7AEADE01700****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
			"x-oss-version-Id": "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?restore&versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days></RestoreRequest>")
		},
		&RestoreObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
			},
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "534B371674E88A4D8906****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())

			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days><JobParameters><Tier>Standard</Tier></JobParameters></RestoreRequest>")
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
				Tier: Ptr("Standard"),
			},
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "534B371674E88A4D8906****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())

			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days><JobParameters><Tier>Standard</Tier></JobParameters></RestoreRequest>")
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
				Tier: Ptr("Standard"),
			},
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockRestoreObject_Success(t *testing.T) {
	for _, c := range testMockRestoreObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RestoreObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockRestoreObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RestoreObjectRequest
	CheckOutputFn  func(t *testing.T, o *RestoreObjectResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6557176CD11E5535303C****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchKey</Code>
		<Message>The specified key does not exist.</Message>
		<RequestId>6557176CD11E5535303C****</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<Key>walker-not-.txt</Key>
		<EC>0026-00000001</EC>
		<RecommendDoc>https://api.aliyun.com/troubleshoot?q=0026-00000001</RecommendDoc>
		</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "6557176CD11E5535303C****", serr.RequestID)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "0026-00000001", serr.EC)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>OperationNotSupported</Code>
  <Message>The operation is not supported for this resource</Message>
  <RequestId>6555AC764311A73931E0****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <Detail>RestoreObject operation does not support this object storage class</Detail>
  <EC>0016-00000702</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "OperationNotSupported", serr.Code)
			assert.Equal(t, "6555AC764311A73931E0****", serr.RequestID)
			assert.Equal(t, "The operation is not supported for this resource", serr.Message)
			assert.Equal(t, "0016-00000702", serr.EC)
		},
	},
	{
		409,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556FF5BD11E55363686****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0026-00000011",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>RestoreAlreadyInProgress</Code>
  <Message>The restore operation is in progress.</Message>
  <RequestId>6556FF5BD11E55363686****</RequestId>
  <HostId>10.101.XX.XX</HostId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "RestoreAlreadyInProgress", serr.Code)
			assert.Equal(t, "6556FF5BD11E55363686****", serr.RequestID)
			assert.Equal(t, "The restore operation is in progress.", serr.Message)
		},
	},
}

func TestMockRestoreObject_Error(t *testing.T) {
	for _, c := range testMockRestoreObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RestoreObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
			assert.Equal(t, string(ObjectACLPublicRead), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPublicRead,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
			assert.Equal(t, string(ObjectACLPrivate), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPrivate,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0%2A%2A%2A%2A", strUrl)
			assert.Equal(t, string(ObjectACLPublicReadWrite), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			Acl:       ObjectACLPublicReadWrite,
			VersionId: Ptr("CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****"),
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get(HTTPHeaderContentType))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get(HeaderOssRequestID))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get(HTTPHeaderDate))

			assert.Equal(t, "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****", *o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0%2A%2A%2A%2A", strUrl)
			assert.Equal(t, string(ObjectACLPublicReadWrite), r.Header.Get(HeaderOssObjectACL))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&PutObjectAclRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			Acl:          ObjectACLPublicReadWrite,
			VersionId:    Ptr("CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get(HTTPHeaderContentType))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get(HeaderOssRequestID))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get(HTTPHeaderDate))

			assert.Equal(t, "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****", *o.VersionId)
		},
	},
}

func TestMockPutObjectAcl_Success(t *testing.T) {
	for _, c := range testMockPutObjectAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectAclResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
			assert.Equal(t, string(ObjectACLPrivate), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPrivate,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl", strUrl)
			assert.Equal(t, string(ObjectACLPrivate), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPrivate,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

func TestMockPutObjectAcl_Error(t *testing.T) {
	for _, c := range testMockPutObjectAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012****</ID>
        <DisplayName>user_example</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read", *o.ACL)
			assert.Equal(t, "0022012****", *o.Owner.ID)
			assert.Equal(t, "user_example", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read-write</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read-write", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>private</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk%2A%2A%2A%2A", strUrl)
		},
		&GetObjectAclRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "private", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
			assert.Equal(t, "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****", *o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>private</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk%2A%2A%2A%2A", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&GetObjectAclRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			VersionId:    Ptr("CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "private", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
			assert.Equal(t, "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****", *o.VersionId)
		},
	},
}

func TestMockGetObjectAcl_Success(t *testing.T) {
	for _, c := range testMockGetObjectAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectAclResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", strUrl)
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockGetObjectAcl_Error(t *testing.T) {
	for _, c := range testMockGetObjectAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockInitiateMultipartUploadSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *InitiateMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *InitiateMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
    <Bucket>oss-example</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<InitiateMultipartUploadResult>
		<Bucket>oss-example</Bucket>
		<Key>multipart.data</Key>
		<UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
		</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object.txt?encoding-type=url&uploads", strUrl)
			assert.Equal(t, r.Header.Get("Cache-Control"), "no-cache")
			assert.Equal(t, r.Header.Get("Content-Disposition"), "attachment")
			assert.Equal(t, r.Header.Get("x-oss-meta-name"), "walker")
			assert.Equal(t, r.Header.Get("x-oss-meta-email"), "demo@aliyun.com")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "9468da86-3509-4f8d-a61e-6eab1eac****")
			assert.Equal(t, r.Header.Get("x-oss-storage-class"), string(StorageClassStandard))
			assert.Equal(t, r.Header.Get("x-oss-forbid-overwrite"), "false")
			assert.Equal(t, r.Header.Get("Content-Encoding"), "utf-8")
			assert.Equal(t, r.Header.Get("Content-MD5"), "1B2M2Y8AsgTpgAmY7PhCfg==")
			assert.Equal(t, r.Header.Get("Expires"), "2022-10-12T00:00:00.000Z")
			assert.Equal(t, r.Header.Get("x-oss-tagging"), "TagA=B&TagC=D")
			assert.Equal(t, "text/plain", r.Header.Get(HTTPHeaderContentType))
		},
		&InitiateMultipartUploadRequest{
			Bucket:                   Ptr("bucket"),
			Key:                      Ptr("object.txt"),
			CacheControl:             Ptr("no-cache"),
			ContentDisposition:       Ptr("attachment"),
			ContentEncoding:          Ptr("utf-8"),
			Expires:                  Ptr("2022-10-12T00:00:00.000Z"),
			ForbidOverwrite:          Ptr("false"),
			ServerSideEncryption:     Ptr("KMS"),
			ServerSideDataEncryption: Ptr("SM4"),
			SSEKMSKeyId:              Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
			StorageClass:             StorageClassStandard,
			Metadata: map[string]string{
				"name":  "walker",
				"email": "demo@aliyun.com",
			},
			Tagging: Ptr("TagA=B&TagC=D"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<InitiateMultipartUploadResult>
		<Bucket>oss-example</Bucket>
		<Key>multipart.data</Key>
		<UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
		</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object.txt?encoding-type=url&uploads", strUrl)
			assert.Equal(t, r.Header.Get("Cache-Control"), "no-cache")
			assert.Equal(t, r.Header.Get("Content-Disposition"), "attachment")
			assert.Equal(t, r.Header.Get("x-oss-meta-name"), "walker")
			assert.Equal(t, r.Header.Get("x-oss-meta-email"), "demo@aliyun.com")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "9468da86-3509-4f8d-a61e-6eab1eac****")
			assert.Equal(t, r.Header.Get("x-oss-storage-class"), string(StorageClassStandard))
			assert.Equal(t, r.Header.Get("x-oss-forbid-overwrite"), "false")
			assert.Equal(t, r.Header.Get("Content-Encoding"), "utf-8")
			assert.Equal(t, r.Header.Get("Content-MD5"), "1B2M2Y8AsgTpgAmY7PhCfg==")
			assert.Equal(t, r.Header.Get("Expires"), "2022-10-12T00:00:00.000Z")
			assert.Equal(t, r.Header.Get("x-oss-tagging"), "TagA=B&TagC=D")
			assert.Equal(t, "text/plain", r.Header.Get(HTTPHeaderContentType))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&InitiateMultipartUploadRequest{
			Bucket:                   Ptr("bucket"),
			Key:                      Ptr("object.txt"),
			CacheControl:             Ptr("no-cache"),
			ContentDisposition:       Ptr("attachment"),
			ContentEncoding:          Ptr("utf-8"),
			Expires:                  Ptr("2022-10-12T00:00:00.000Z"),
			ForbidOverwrite:          Ptr("false"),
			ServerSideEncryption:     Ptr("KMS"),
			ServerSideDataEncryption: Ptr("SM4"),
			SSEKMSKeyId:              Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
			StorageClass:             StorageClassStandard,
			Metadata: map[string]string{
				"name":  "walker",
				"email": "demo@aliyun.com",
			},
			Tagging:      Ptr("TagA=B&TagC=D"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
}

func TestMockInitiateMultipartUpload_Success(t *testing.T) {
	for _, c := range testMockInitiateMultipartUploadSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InitiateMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockInitiateMultipartUploadDisableDetectMimeTypeCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *InitiateMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *InitiateMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
    <Bucket>oss-example</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
			assert.Equal(t, "", r.Header.Get(HTTPHeaderContentType))
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
}

func TestMockInitiateMultipartUpload_DisableDetectMimeType(t *testing.T) {
	for _, c := range testMockInitiateMultipartUploadDisableDetectMimeTypeCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureAutoDetectMimeType
			})
		assert.NotNil(t, c)

		output, err := client.InitiateMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockInitiateMultipartUploadErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *InitiateMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *InitiateMultipartUploadResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute InitiateMultipartUpload fail")
		},
	},
}

func TestMockInitiateMultipartUpload_Error(t *testing.T) {
	for _, c := range testMockInitiateMultipartUploadErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InitiateMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F586****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Ph****",
			"x-oss-hash-crc64ecma": "6571598172666981661",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
			assert.Equal(t, "bce8f3d48247c5d555bb5697bf277b35", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
			ContentMD5: Ptr("bce8f3d48247c5d555bb5697bf277b35"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F586****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
			assert.Equal(t, *o.HashCRC64, "6571598172666981661")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F587****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Pp****",
			"x-oss-hash-crc64ecma": "2060813895736234537",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 2")
			assert.Equal(t, "f811b746eb3e256f97cb3a190d528353", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(2),
			Body:       strings.NewReader("upload part 2"),
			ContentMD5: Ptr("f811b746eb3e256f97cb3a190d528353"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F587****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Pp****")
			assert.Equal(t, *o.HashCRC64, "2060813895736234537")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F587****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Pp****",
			"x-oss-hash-crc64ecma": "2060813895736234537",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 2")
			assert.Equal(t, "f811b746eb3e256f97cb3a190d528353", r.Header.Get("Content-MD5"))
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&UploadPartRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(2),
			Body:         strings.NewReader("upload part 2"),
			ContentMD5:   Ptr("f811b746eb3e256f97cb3a190d528353"),
			TrafficLimit: int64(100 * 1024 * 8),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F587****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Pp****")
			assert.Equal(t, *o.HashCRC64, "2060813895736234537")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F586****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Ph****",
			"x-oss-hash-crc64ecma": "6571598172666981661",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
			assert.Equal(t, "bce8f3d48247c5d555bb5697bf277b35", r.Header.Get("Content-MD5"))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&UploadPartRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			Body:         strings.NewReader("upload part 1"),
			ContentMD5:   Ptr("bce8f3d48247c5d555bb5697bf277b35"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F586****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
			assert.Equal(t, *o.HashCRC64, "6571598172666981661")
		},
	},
}

func TestMockUploadPart_Success(t *testing.T) {
	for _, c := range testMockUploadPartSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartWithProgressCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F586****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Ph****",
			"x-oss-hash-crc64ecma": "6571598172666981661",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
			assert.Equal(t, "bce8f3d48247c5d555bb5697bf277b35", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
			ContentMD5: Ptr("bce8f3d48247c5d555bb5697bf277b35"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F586****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
			assert.Equal(t, *o.HashCRC64, "6571598172666981661")
		},
	},
}

func TestMockUploadPart_Progress(t *testing.T) {
	for _, c := range testMockUploadPartWithProgressCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		n := int64(0)
		c.Request.ProgressFn = func(increment, transferred, total int64) {
			n = transferred
			//fmt.Printf("got transferred:%v, total:%v\n", transferred, total)
		}
		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, int64(len("upload part 1")), n)

	}
}

var testMockUploadPartDisableCRC64Cases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F586****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Ph****",
			"x-oss-hash-crc64ecma": "8571598172666981661",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
			assert.Equal(t, "bce8f3d48247c5d555bb5697bf277b35", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
			ContentMD5: Ptr("bce8f3d48247c5d555bb5697bf277b35"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F586****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
			assert.Equal(t, *o.HashCRC64, "8571598172666981661")
		},
	},
}

func TestMockUploadPart_DisableCRC64(t *testing.T) {
	for _, c := range testMockUploadPartDisableCRC64Cases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		//Disable
		client := NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureEnableCRC64CheckUpload
			})
		assert.NotNil(t, c)
		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)

		client = NewClient(cfg)
		assert.NotNil(t, c)
		c.Request.Body = strings.NewReader("upload part 1")
		_, err = client.UploadPart(context.TODO(), c.Request)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "crc is inconsistent, client 6571598172666981661, server 8571598172666981661")
	}
}

var testMockUploadPartErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockUploadPart_Error(t *testing.T) {
	for _, c := range testMockUploadPartErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartCopySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartCopyRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartCopyResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object"))
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":             "6551DBCF4311A7303980****",
			"Date":                         "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-copy-source-version-id": "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object")+"?versionId=CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY")

			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfMatch), "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfNoneMatch), "\"D41D8CD98F00B204E9800998ECF9****\"")
			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfModifiedSince), "Fri, 13 Nov 2023 14:47:53 GMT")
			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfUnmodifiedSince), "Fri, 13 Nov 2015 14:47:53 GMT")
		},
		&UploadPartCopyRequest{
			Bucket:            Ptr("bucket"),
			Key:               Ptr("object"),
			UploadId:          Ptr("0004B9895DBBB6EC9"),
			SourceKey:         Ptr("oss-src-object"),
			SourceBucket:      Ptr("oss-src-bucket"),
			PartNumber:        int32(2),
			IfMatch:           Ptr("\"D41D8CD98F00B204E9800998ECF8****\""),
			IfNoneMatch:       Ptr("\"D41D8CD98F00B204E9800998ECF9****\""),
			IfModifiedSince:   Ptr("Fri, 13 Nov 2023 14:47:53 GMT"),
			IfUnmodifiedSince: Ptr("Fri, 13 Nov 2015 14:47:53 GMT"),
			SourceVersionId:   Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
			assert.Equal(t, *o.VersionId, "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":             "6551DBCF4311A7303980****",
			"Date":                         "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-copy-source-version-id": "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object"))
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
			PartNumber:   int32(2),
			TrafficLimit: int64(100 * 1024 * 8),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
			assert.Equal(t, *o.VersionId, "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object"))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
		},
	},
}

func TestMockUploadPartCopy_Success(t *testing.T) {
	for _, c := range testMockUploadPartCopySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPartCopy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartCopyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartCopyRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartCopyResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute UploadPartCopy fail")
		},
	},
}

func TestMockUploadPartCopy_Error(t *testing.T) {
	for _, c := range testMockUploadPartCopyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPartCopy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCompleteMultipartUploadSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CompleteMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *CompleteMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult>
  <EncodingType>url</EncodingType>
  <Location>http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data</Location>
  <Bucket>oss-example</Bucket>
  <Key>demo%2Fmultipart.data</Key>
  <ETag>"097DE458AD02B5F89F9D0530231876****"</ETag>
</CompleteMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), `<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>&#34;8EFDA8BE206636A695359836FE0A****&#34;</ETag></Part><Part><PartNumber>2</PartNumber><ETag>&#34;8C315065167132444177411FDA14****&#34;</ETag></Part><Part><PartNumber>3</PartNumber><ETag>&#34;3349DC700140D7F86A0784842780****&#34;</ETag></Part></CompleteMultipartUpload>`)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
			CompleteMultipartUpload: &CompleteMultipartUpload{
				Parts: []UploadPart{
					{PartNumber: int32(3), ETag: Ptr("\"3349DC700140D7F86A0784842780****\"")},
					{PartNumber: int32(1), ETag: Ptr("\"8EFDA8BE206636A695359836FE0A****\"")},
					{PartNumber: int32(2), ETag: Ptr("\"8C315065167132444177411FDA14****\"")},
				},
			},
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"097DE458AD02B5F89F9D0530231876****\"")
			assert.Equal(t, *o.Location, "http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data")
			assert.Equal(t, *o.EncodingType, "url")
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "demo/multipart.data")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":     "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****",
			"Content-Type":         "application/json",
			"x-oss-hash-crc64ecma": "1206617243528768****",
		},
		[]byte(`{"filename":"oss-obj.txt","size":"100","mimeType":"","x":"a","b":"b"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, "false", r.Header.Get(HeaderOssForbidOverWrite))
			assert.Equal(t, "yes", r.Header.Get("x-oss-complete-all"))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`)), r.Header.Get(HeaderOssCallback))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`)), r.Header.Get(HeaderOssCallbackVar))
		},
		&CompleteMultipartUploadRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			UploadId:        Ptr("0004B9895DBBB6EC9"),
			ForbidOverwrite: Ptr("false"),
			CompleteAll:     Ptr("yes"),
			Callback:        Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`))),
			CallbackVar:     Ptr(base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`))),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get(HTTPHeaderContentType), "application/json")
			jsonData, _ := json.Marshal(o.CallbackResult)
			assert.Nil(t, err)
			assert.NotEmpty(t, string(jsonData))
			assert.Equal(t, *o.VersionId, "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":     "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****",
			"Content-Type":         "application/json",
			"x-oss-hash-crc64ecma": "1206617243528768****",
		},
		[]byte(`{"filename":"oss-obj.txt","size":"100","mimeType":"","x":"a","b":"b"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, "false", r.Header.Get(HeaderOssForbidOverWrite))
			assert.Equal(t, "yes", r.Header.Get("x-oss-complete-all"))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`)), r.Header.Get(HeaderOssCallback))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`)), r.Header.Get(HeaderOssCallbackVar))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&CompleteMultipartUploadRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			UploadId:        Ptr("0004B9895DBBB6EC9"),
			ForbidOverwrite: Ptr("false"),
			CompleteAll:     Ptr("yes"),
			Callback:        Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`))),
			CallbackVar:     Ptr(base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`))),
			RequestPayer:    Ptr("requester"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get(HTTPHeaderContentType), "application/json")
			jsonData, _ := json.Marshal(o.CallbackResult)
			assert.Nil(t, err)
			assert.NotEmpty(t, string(jsonData))
			assert.Equal(t, *o.VersionId, "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****")
		},
	},
}

func TestMockCompleteMultipartUpload_Success(t *testing.T) {
	for _, c := range testMockCompleteMultipartUploadSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CompleteMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCompleteMultipartUploadErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CompleteMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *CompleteMultipartUploadResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "655D94CCD11E55313348****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MalformedXML</Code>
  <Message>The XML you provided was not well-formed or did not validate against our published schema.</Message>
  <RequestId>655D94CCD11E55313348****</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0042-00000205</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0042-00000205</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "MalformedXML", serr.Code)
			assert.Equal(t, "The XML you provided was not well-formed or did not validate against our published schema.", serr.Message)
			assert.Equal(t, "655D94CCD11E55313348****", serr.RequestID)
			assert.Equal(t, "0042-00000205", serr.EC)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "655D9598CA31DC313626****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>Should not speficy both complete all header and http body.</Message>
  <RequestId>655D9598CA31DC313626****</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <ArgumentName>x-oss-complete-all</ArgumentName>
  <ArgumentValue>yes</ArgumentValue>
  <EC>0042-00000216</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0042-00000216</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:      Ptr("bucket"),
			Key:         Ptr("object"),
			UploadId:    Ptr("0004B9895DBBB6EC9"),
			CompleteAll: Ptr("yes"),
			CompleteMultipartUpload: &CompleteMultipartUpload{
				Parts: []UploadPart{
					{PartNumber: int32(3), ETag: Ptr("\"3349DC700140D7F86A0784842780****\"")},
					{PartNumber: int32(1), ETag: Ptr("\"8EFDA8BE206636A695359836FE0A****\"")},
					{PartNumber: int32(2), ETag: Ptr("\"8C315065167132444177411FDA14****\"")},
				},
			},
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "Should not speficy both complete all header and http body.", serr.Message)
			assert.Equal(t, "655D9598CA31DC313626****", serr.RequestID)
			assert.Equal(t, "0042-00000216", serr.EC)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute CompleteMultipartUpload fail")
		},
	},
}

func TestMockCompleteMultipartUpload_Error(t *testing.T) {
	for _, c := range testMockCompleteMultipartUploadErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CompleteMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAbortMultipartUploadSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AbortMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *AbortMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
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
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&AbortMultipartUploadRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6E"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockAbortMultipartUpload_Success(t *testing.T) {
	for _, c := range testMockAbortMultipartUploadSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AbortMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAbortMultipartUploadErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AbortMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *AbortMultipartUploadResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockAbortMultipartUpload_Error(t *testing.T) {
	for _, c := range testMockAbortMultipartUploadErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AbortMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListMultipartUploadsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListMultipartUploadsRequest
	CheckOutputFn  func(t *testing.T, o *ListMultipartUploadsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
    <Bucket>oss-example</Bucket>
    <KeyMarker></KeyMarker>
    <UploadIdMarker></UploadIdMarker>
    <NextKeyMarker>oss.avi</NextKeyMarker>
    <NextUploadIdMarker>0004B99B8E707874FC2D692FA5D77D3F</NextUploadIdMarker>
    <Delimiter></Delimiter>
    <Prefix></Prefix>
    <MaxUploads>1000</MaxUploads>
    <IsTruncated>false</IsTruncated>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF518A1FE585B0C9360DC4C8</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>oss.avi</Key>
        <UploadId>0004B99B8E707874FC2D692FA5D7****</UploadId>
        <Initiated>2012-02-23T06:14:27.000Z</Initiated>
    </Upload>
</ListMultipartUploadsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.UploadIdMarker, "")
			assert.Equal(t, *o.NextKeyMarker, "oss.avi")
			assert.Equal(t, *o.NextUploadIdMarker, "0004B99B8E707874FC2D692FA5D77D3F")
			assert.Equal(t, *o.Delimiter, "")
			assert.Equal(t, *o.Prefix, "")
			assert.Equal(t, o.MaxUploads, int32(1000))
			assert.Equal(t, o.IsTruncated, false)
			assert.Len(t, o.Uploads, 3)
			assert.Equal(t, *o.Uploads[0].Key, "multipart.data")
			assert.Equal(t, *o.Uploads[0].UploadId, "0004B999EF518A1FE585B0C9360DC4C8")
			assert.Equal(t, *o.Uploads[0].Initiated, time.Date(2012, time.February, 23, 4, 18, 23, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
  <EncodingType>url</EncodingType>
  <Bucket>oss-example</Bucket>
  <KeyMarker></KeyMarker>
  <UploadIdMarker></UploadIdMarker>
  <NextKeyMarker>oss.avi</NextKeyMarker>
  <NextUploadIdMarker>89F0105AA66942638E35300618DF****</NextUploadIdMarker>
  <Delimiter>/</Delimiter>
  <Prefix>pre</Prefix>
  <MaxUploads>1000</MaxUploads>
  <IsTruncated>false</IsTruncated>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>0214A87687F040F1BA4D83AB17C9****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:45:57.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>3AE2ED7A60E04AFE9A5287055D37****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:03:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>47E0E90F5DCB4AD5B3C4CD886CB0****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:02:11.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>A89E0E28E2E948A1BFF6FD5CDAFF****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T06:57:03.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>B18E1DCDB6964F5CB197F5F6B26A****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:02.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>5209986C3A96486EA16B9C52C160****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:34:47.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>63B652FA2C1342DCB3CCCC86D748****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:28:46.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>6F67B34BCA3C481F887D73508A07****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:32:12.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>89F0105AA66942638E35300618D****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:37:53.000Z</Initiated>
  </Upload>
</ListMultipartUploadsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&key-marker&max-uploads=10&prefix=pre&upload-id-marker&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket:         Ptr("bucket"),
			Delimiter:      Ptr("/"),
			Prefix:         Ptr("pre"),
			EncodingType:   Ptr("url"),
			KeyMarker:      Ptr(""),
			MaxUploads:     int32(10),
			UploadIdMarker: Ptr(""),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.UploadIdMarker, "")
			assert.Equal(t, *o.NextKeyMarker, "oss.avi")
			assert.Equal(t, *o.NextUploadIdMarker, "89F0105AA66942638E35300618DF****")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, *o.Prefix, "pre")
			assert.Equal(t, o.MaxUploads, int32(1000))
			assert.Equal(t, o.IsTruncated, false)
			assert.Len(t, o.Uploads, 10)
			assert.Equal(t, *o.Uploads[0].Key, "demo/gp-\f\n\v")
			assert.Equal(t, *o.Uploads[0].UploadId, "0214A87687F040F1BA4D83AB17C9****")
			assert.Equal(t, *o.Uploads[0].Initiated, time.Date(2023, time.November, 22, 5, 45, 57, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
  <EncodingType>url</EncodingType>
  <Bucket>oss-example</Bucket>
  <KeyMarker></KeyMarker>
  <UploadIdMarker></UploadIdMarker>
  <NextKeyMarker>oss.avi</NextKeyMarker>
  <NextUploadIdMarker>89F0105AA66942638E35300618DF****</NextUploadIdMarker>
  <Delimiter>/</Delimiter>
  <Prefix>pre</Prefix>
  <MaxUploads>1000</MaxUploads>
  <IsTruncated>false</IsTruncated>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>0214A87687F040F1BA4D83AB17C9****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:45:57.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>3AE2ED7A60E04AFE9A5287055D37****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:03:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>47E0E90F5DCB4AD5B3C4CD886CB0****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:02:11.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>A89E0E28E2E948A1BFF6FD5CDAFF****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T06:57:03.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>B18E1DCDB6964F5CB197F5F6B26A****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:02.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>5209986C3A96486EA16B9C52C160****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:34:47.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>63B652FA2C1342DCB3CCCC86D748****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:28:46.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>6F67B34BCA3C481F887D73508A07****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:32:12.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>89F0105AA66942638E35300618D****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:37:53.000Z</Initiated>
  </Upload>
</ListMultipartUploadsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&key-marker&max-uploads=10&prefix=pre&upload-id-marker&uploads", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListMultipartUploadsRequest{
			Bucket:         Ptr("bucket"),
			Delimiter:      Ptr("/"),
			Prefix:         Ptr("pre"),
			EncodingType:   Ptr("url"),
			KeyMarker:      Ptr(""),
			MaxUploads:     int32(10),
			UploadIdMarker: Ptr(""),
			RequestPayer:   Ptr("requester"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.UploadIdMarker, "")
			assert.Equal(t, *o.NextKeyMarker, "oss.avi")
			assert.Equal(t, *o.NextUploadIdMarker, "89F0105AA66942638E35300618DF****")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, *o.Prefix, "pre")
			assert.Equal(t, o.MaxUploads, int32(1000))
			assert.Equal(t, o.IsTruncated, false)
			assert.Len(t, o.Uploads, 10)
			assert.Equal(t, *o.Uploads[0].Key, "demo/gp-\f\n\v")
			assert.Equal(t, *o.Uploads[0].UploadId, "0214A87687F040F1BA4D83AB17C9****")
			assert.Equal(t, *o.Uploads[0].Initiated, time.Date(2023, time.November, 22, 5, 45, 57, 0, time.UTC))
		},
	},
}

func TestMockListMultipartUploads_Success(t *testing.T) {
	for _, c := range testMockListMultipartUploadsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListMultipartUploads(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListMultipartUploadsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListMultipartUploadsRequest
	CheckOutputFn  func(t *testing.T, o *ListMultipartUploadsResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockListMultipartUploads_Error(t *testing.T) {
	for _, c := range testMockListMultipartUploadsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListMultipartUploads(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListPartsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListPartsRequest
	CheckOutputFn  func(t *testing.T, o *ListPartsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
    <Bucket>bucket</Bucket>
    <Key>object</Key>
    <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
    <NextPartNumberMarker>5</NextPartNumberMarker>
    <MaxParts>1000</MaxParts>
    <IsTruncated>false</IsTruncated>
    <Part>
        <PartNumber>1</PartNumber>
        <LastModified>2012-02-23T07:01:34.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>2</PartNumber>
        <LastModified>2012-02-23T07:01:12.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>5</PartNumber>
        <LastModified>2012-02-23T07:02:03.000Z</LastModified>
        <ETag>"7265F4D211B56873A381D321F586****"</ETag>
        <Size>1024</Size>
    </Part>
</ListPartsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Bucket, "bucket")
			assert.Equal(t, *o.Key, "object")
			assert.Empty(t, o.PartNumberMarker)
			assert.Equal(t, o.NextPartNumberMarker, int32(5))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxParts, int32(1000))
			assert.Len(t, o.Parts, 3)
			assert.Equal(t, o.Parts[0].PartNumber, int32(1))
			assert.Equal(t, *o.Parts[0].ETag, "\"3349DC700140D7F86A0784842780****\"")
			assert.Equal(t, *o.Parts[0].LastModified, time.Date(2012, time.February, 23, 7, 1, 34, 0, time.UTC))
			assert.Equal(t, o.Parts[0].Size, int64(6291456))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
  <EncodingType>url</EncodingType>
  <Bucket>bucket</Bucket>
  <Key>demo%2Fgp-%0C%0A%0B</Key>
  <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
  <StorageClass>Standard</StorageClass>
  <PartNumberMarker>0</PartNumberMarker>
  <NextPartNumberMarker>1</NextPartNumberMarker>
  <MaxParts>1000</MaxParts>
  <IsTruncated>false</IsTruncated>
  <Part>
    <PartNumber>1</PartNumber>
    <LastModified>2023-11-22T05:42:34.000Z</LastModified>
    <ETag>"CF3F46D505093571E916FCDD4967****"</ETag>
    <HashCrc64ecma>12066172435287683848</HashCrc64ecma>
    <Size>96316</Size>
  </Part>
</ListPartsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/demo/gp-%0C%0A%0B?encoding-type=url&uploadId=D4E111D4EA834F3ABCE4877B2779%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("demo/gp-\f\n\v"),
			UploadId: Ptr("D4E111D4EA834F3ABCE4877B2779****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Bucket, "bucket")
			key, _ := url.QueryUnescape("demo%2Fgp-%0C%0A%0B")
			assert.Equal(t, *o.Key, key)
			assert.Empty(t, o.PartNumberMarker)
			assert.Equal(t, o.NextPartNumberMarker, int32(1))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxParts, int32(1000))
			assert.Len(t, o.Parts, 1)
			assert.Equal(t, o.Parts[0].PartNumber, int32(1))
			assert.Equal(t, *o.Parts[0].ETag, "\"CF3F46D505093571E916FCDD4967****\"")
			assert.Equal(t, *o.Parts[0].LastModified, time.Date(2023, time.November, 22, 5, 42, 34, 0, time.UTC))
			assert.Equal(t, o.Parts[0].Size, int64(96316))
			assert.Equal(t, *o.Parts[0].HashCRC64, "12066172435287683848")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
    <Bucket>bucket</Bucket>
    <Key>object</Key>
    <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
    <NextPartNumberMarker>5</NextPartNumberMarker>
    <MaxParts>1000</MaxParts>
    <IsTruncated>false</IsTruncated>
    <Part>
        <PartNumber>1</PartNumber>
        <LastModified>2012-02-23T07:01:34.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>2</PartNumber>
        <LastModified>2012-02-23T07:01:12.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>5</PartNumber>
        <LastModified>2012-02-23T07:02:03.000Z</LastModified>
        <ETag>"7265F4D211B56873A381D321F586****"</ETag>
        <Size>1024</Size>
    </Part>
</ListPartsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListPartsRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B999EF5A239BB9138C6227D6****"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Bucket, "bucket")
			assert.Equal(t, *o.Key, "object")
			assert.Empty(t, o.PartNumberMarker)
			assert.Equal(t, o.NextPartNumberMarker, int32(5))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxParts, int32(1000))
			assert.Len(t, o.Parts, 3)
			assert.Equal(t, o.Parts[0].PartNumber, int32(1))
			assert.Equal(t, *o.Parts[0].ETag, "\"3349DC700140D7F86A0784842780****\"")
			assert.Equal(t, *o.Parts[0].LastModified, time.Date(2012, time.February, 23, 7, 1, 34, 0, time.UTC))
			assert.Equal(t, o.Parts[0].Size, int64(6291456))
		},
	},
}

func TestMockListParts_Success(t *testing.T) {
	for _, c := range testMockListPartsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListParts(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListPartsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListPartsRequest
	CheckOutputFn  func(t *testing.T, o *ListPartsResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListParts fail")
		},
	},
}

func TestMockListParts_Error(t *testing.T) {
	for _, c := range testMockListPartsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListParts(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketVersioningSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketVersioningResult, err error)
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
			assert.Equal(t, "/bucket?versioning", r.URL.String())
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Suspended</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionSuspended,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
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
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Enabled</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionEnabled,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketVersioning_Success(t *testing.T) {
	for _, c := range testMockPutBucketVersioningSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketVersioningErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketVersioningResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Enabled</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionEnabled,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Enabled</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionEnabled,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutBucketVersioning fail")
		},
	},
}

func TestMockPutBucketVersioning_Error(t *testing.T) {
	for _, c := range testMockPutBucketVersioningErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketVersioningSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketVersioningResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
</VersioningConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Nil(t, o.VersionStatus)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
<Status>Enabled</Status>
</VersioningConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionStatus, "Enabled")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
<Status>Suspended</Status>
</VersioningConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionStatus, "Suspended")
		},
	},
}

func TestMockGetBucketVersioning_Success(t *testing.T) {
	for _, c := range testMockGetBucketVersioningSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketVersioningErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketVersioningResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketVersioning fail")
		},
	},
}

func TestMockGetBucketVersioning_Error(t *testing.T) {
	for _, c := range testMockGetBucketVersioningErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectVersionsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectVersionsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectVersionsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2F</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>20</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>false</IsTruncated>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&key-marker&max-keys=20&prefix=demo%2F&version-id-marker&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket:          Ptr("bucket"),
			Delimiter:       Ptr("/"),
			Prefix:          Ptr("demo/"),
			KeyMarker:       Ptr(""),
			VersionIdMarker: Ptr(""),
			MaxKeys:         int32(20),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "demo-bucket")
			prefix, _ := url.QueryUnescape(*o.Prefix)
			assert.Equal(t, *o.Prefix, prefix)
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(20))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 0)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "examplebucket-1250000000")
			assert.Equal(t, *o.Prefix, "")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(1000))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 3)
			assert.Equal(t, *o.ObjectVersions[0].Key, "example-object-1.jpg")
			assert.Empty(t, *o.ObjectVersions[1].VersionId)
			assert.True(t, o.ObjectVersions[2].IsLatest)
			assert.NotEmpty(t, *o.ObjectVersions[0].LastModified)
			assert.Equal(t, *o.ObjectVersions[1].ETag, "5B3C1A2E053D763E1B002CC607C5A0FE1****")
			assert.Equal(t, o.ObjectVersions[2].Size, int64(20))
			assert.Equal(t, *o.ObjectVersions[2].Owner.ID, "1250000000")
			assert.Equal(t, *o.ObjectVersions[2].Owner.DisplayName, "1250000000")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&key-marker&max-keys=5&prefix=demo%2Fgp-&version-id-marker&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket:          Ptr("bucket"),
			KeyMarker:       Ptr(""),
			VersionIdMarker: Ptr(""),
			Delimiter:       Ptr("/"),
			MaxKeys:         int32(5),
			Prefix:          Ptr("demo/gp-"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "demo-bucket")
			prefix, _ := url.QueryUnescape(*o.Prefix)
			assert.Equal(t, *o.Prefix, prefix)
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(5))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 3)
			key, _ := url.QueryUnescape(*o.ObjectVersions[0].Key)
			assert.Equal(t, *o.ObjectVersions[0].Key, key)
			assert.Equal(t, *o.ObjectVersions[1].VersionId, "CAEQHxiBgMDYseHV3xgiIDg2Mzk0Zjg3MjQ0MTRhM2FiMzgxOGY1NjdmN2Rk****")
			assert.False(t, o.ObjectVersions[2].IsLatest)
			assert.NotEmpty(t, *o.ObjectVersions[0].LastModified)
			assert.Equal(t, *o.ObjectVersions[1].ETag, "\"29B94424BC241D80B0AF488A4E4B86AF-1\"")
			assert.Equal(t, o.ObjectVersions[2].Size, int64(96316))
			assert.Equal(t, *o.ObjectVersions[2].Owner.ID, "150692521021****")
			assert.Equal(t, *o.ObjectVersions[2].Owner.DisplayName, "150692521021****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2F</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>20</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>false</IsTruncated>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&key-marker&max-keys=20&prefix=demo%2F&version-id-marker&versions", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListObjectVersionsRequest{
			Bucket:          Ptr("bucket"),
			Delimiter:       Ptr("/"),
			Prefix:          Ptr("demo/"),
			KeyMarker:       Ptr(""),
			VersionIdMarker: Ptr(""),
			MaxKeys:         int32(20),
			RequestPayer:    Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "demo-bucket")
			prefix, _ := url.QueryUnescape(*o.Prefix)
			assert.Equal(t, *o.Prefix, prefix)
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(20))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 0)
		},
	},
}

func TestMockListObjectVersions_Success(t *testing.T) {
	for _, c := range testMockListObjectVersionsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectVersions(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectVersionsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectVersionsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectVersionsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListObjectVersions fail")
		},
	},
}

func TestMockListObjectVersions_Error(t *testing.T) {
	for _, c := range testMockListObjectVersionsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectVersions(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutSymlinkSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutSymlinkRequest
	CheckOutputFn  func(t *testing.T, o *PutSymlinkResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&PutSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Target: Ptr("src-object"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
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
			"Content-Type":     "application/xml",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)

		},
		&PutSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Target: Ptr("src-object"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-symlink-target"), "target-object")
			assert.Equal(t, r.Header.Get("x-oss-forbid-overwrite"), "true")
			assert.Equal(t, r.Header.Get("x-oss-object-acl"), string(ObjectACLPrivate))
			assert.Equal(t, r.Header.Get("x-oss-storage-class"), string(StorageClassStandard))
			assert.Equal(t, r.Header.Get("x-oss-meta-name"), "demo")
			assert.Equal(t, r.Header.Get("x-oss-meta-email"), "demo@aliyun.com")
		},
		&PutSymlinkRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			Target:          Ptr("target-object"),
			ForbidOverwrite: Ptr("true"),
			Acl:             ObjectACLPrivate,
			StorageClass:    StorageClassStandard,
			Metadata: map[string]string{
				"name":  "demo",
				"email": "demo@aliyun.com",
			},
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
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
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&PutSymlinkRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			Target:       Ptr("src-object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutSymlink_Success(t *testing.T) {
	for _, c := range testMockPutSymlinkSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutSymlink(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutSymlinkErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutSymlinkRequest
	CheckOutputFn  func(t *testing.T, o *PutSymlinkResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&PutSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Target: Ptr("target-object"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&PutSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Target: Ptr("target-object"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

func TestMockPutSymlink_Error(t *testing.T) {
	for _, c := range testMockPutSymlinkErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutSymlink(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetSymlinkSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetSymlinkRequest
	CheckOutputFn  func(t *testing.T, o *GetSymlinkResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-symlink-target": "example.jpg",
			"ETag":                 "A797938C31D59EDD08D86188F6D5****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&GetSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Target, "example.jpg")
			assert.Equal(t, *o.ETag, "A797938C31D59EDD08D86188F6D5****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":         "application/xml",
			"x-oss-version-id":     "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
			"x-oss-symlink-target": "example.jpg",
			"ETag":                 "A797938C31D59EDD08D86188F6D5****",
			"x-oss-meta-name":      "demo",
			"x-oss-meta-email":     "demo@aliyun.com",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink&versionId=CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh%2A%2A%2A%2A", strUrl)
		},
		&GetSymlinkRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
			assert.Equal(t, *o.Target, "example.jpg")
			assert.Equal(t, *o.ETag, "A797938C31D59EDD08D86188F6D5****")
			assert.Equal(t, o.Metadata["name"], "demo")
			assert.Equal(t, o.Metadata["email"], "demo@aliyun.com")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-symlink-target": "example.jpg",
			"ETag":                 "A797938C31D59EDD08D86188F6D5****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&GetSymlinkRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Target, "example.jpg")
			assert.Equal(t, *o.ETag, "A797938C31D59EDD08D86188F6D5****")
		},
	},
}

func TestMockGetSymlink_Success(t *testing.T) {
	for _, c := range testMockGetSymlinkSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetSymlink(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetSymlinkErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetSymlinkRequest
	CheckOutputFn  func(t *testing.T, o *GetSymlinkResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&GetSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&GetSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

func TestMockGetSymlink_Error(t *testing.T) {
	for _, c := range testMockGetSymlinkErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetSymlink(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectTaggingSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectTaggingRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectTaggingResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `<Tagging><TagSet><Tag><Key>k1</Key><Value>v1</Value></Tag></TagSet></Tagging>`)
		},
		&PutObjectTaggingRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Tagging: &Tagging{
				TagSet{
					Tags: []Tag{
						{
							Key:   Ptr("k1"),
							Value: Ptr("v1"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutObjectTaggingResult, err error) {
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
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging&versionId=CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh%2A%2A%2A%2A", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `<Tagging><TagSet><Tag><Key>k1</Key><Value>v1</Value></Tag><Tag><Key>k2</Key><Value>v2</Value></Tag></TagSet></Tagging>`)
		},
		&PutObjectTaggingRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
			Tagging: &Tagging{
				TagSet{
					Tags: []Tag{
						{
							Key:   Ptr("k1"),
							Value: Ptr("v1"),
						},
						{
							Key:   Ptr("k2"),
							Value: Ptr("v2"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutObjectTaggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging&versionId=CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh%2A%2A%2A%2A", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `<Tagging><TagSet><Tag><Key>k1</Key><Value>v1</Value></Tag><Tag><Key>k2</Key><Value>v2</Value></Tag></TagSet></Tagging>`)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&PutObjectTaggingRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
			Tagging: &Tagging{
				TagSet{
					Tags: []Tag{
						{
							Key:   Ptr("k1"),
							Value: Ptr("v1"),
						},
						{
							Key:   Ptr("k2"),
							Value: Ptr("v2"),
						},
					},
				},
			},
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *PutObjectTaggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
		},
	},
}

func TestMockPutObjectTagging_Success(t *testing.T) {
	for _, c := range testMockPutObjectTaggingSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObjectTagging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectTaggingErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectTaggingRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectTaggingResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `<Tagging><TagSet><Tag><Key>k1</Key><Value>v1</Value></Tag><Tag><Key>k2</Key><Value>v2</Value></Tag></TagSet></Tagging>`)
		},
		&PutObjectTaggingRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Tagging: &Tagging{
				TagSet{
					Tags: []Tag{
						{
							Key:   Ptr("k1"),
							Value: Ptr("v1"),
						},
						{
							Key:   Ptr("k2"),
							Value: Ptr("v2"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutObjectTaggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `<Tagging><TagSet><Tag><Key>k1</Key><Value>v1</Value></Tag><Tag><Key>k2</Key><Value>v2</Value></Tag></TagSet></Tagging>`)
		},
		&PutObjectTaggingRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Tagging: &Tagging{
				TagSet{
					Tags: []Tag{
						{
							Key:   Ptr("k1"),
							Value: Ptr("v1"),
						},
						{
							Key:   Ptr("k2"),
							Value: Ptr("v2"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutObjectTaggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

func TestMockPutObjectTagging_Error(t *testing.T) {
	for _, c := range testMockPutObjectTaggingErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObjectTagging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectTaggingSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectTaggingRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectTaggingResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Tagging>
  <TagSet>
    <Tag>
      <Key>age</Key>
      <Value>18</Value>
    </Tag>
  </TagSet>
</Tagging>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging", strUrl)
		},
		&GetObjectTaggingRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectTaggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.Tags, 1)
			assert.Equal(t, *o.Tags[0].Key, "age")
			assert.Equal(t, *o.Tags[0].Value, "18")

		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Tagging>
  <TagSet>
    <Tag>
      <Key>a</Key>
      <Value>1</Value>
    </Tag>
    <Tag>
      <Key>b</Key>
      <Value>2</Value>
    </Tag>
  </TagSet>
</Tagging>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging&versionId=CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh%2A%2A%2A%2A", strUrl)
		},
		&GetObjectTaggingRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
		},
		func(t *testing.T, o *GetObjectTaggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
			assert.Len(t, o.Tags, 2)
			assert.Equal(t, *o.Tags[0].Key, "a")
			assert.Equal(t, *o.Tags[0].Value, "1")
			assert.Equal(t, *o.Tags[1].Key, "b")
			assert.Equal(t, *o.Tags[1].Value, "2")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Tagging>
  <TagSet>
    <Tag>
      <Key>age</Key>
      <Value>18</Value>
    </Tag>
  </TagSet>
</Tagging>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&GetObjectTaggingRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *GetObjectTaggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.Tags, 1)
			assert.Equal(t, *o.Tags[0].Key, "age")
			assert.Equal(t, *o.Tags[0].Value, "18")

		},
	},
}

func TestMockGetObjectTagging_Success(t *testing.T) {
	for _, c := range testMockGetObjectTaggingSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectTagging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectTaggingErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectTaggingRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectTaggingResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging", strUrl)
		},
		&GetObjectTaggingRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectTaggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging", strUrl)
		},
		&GetObjectTaggingRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectTaggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

func TestMockGetObjectTagging_Error(t *testing.T) {
	for _, c := range testMockGetObjectTaggingErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectTagging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteObjectTaggingSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteObjectTaggingRequest
	CheckOutputFn  func(t *testing.T, o *DeleteObjectTaggingResult, err error)
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
			assert.Equal(t, "/bucket/object?tagging", strUrl)
		},
		&DeleteObjectTaggingRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectTaggingResult, err error) {
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
			assert.Equal(t, "/bucket/object?tagging&versionId=CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh%2A%2A%2A%2A", strUrl)
		},
		&DeleteObjectTaggingRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
		},
		func(t *testing.T, o *DeleteObjectTaggingResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
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
			assert.Equal(t, "/bucket/object?tagging&versionId=CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh%2A%2A%2A%2A", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&DeleteObjectTaggingRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			VersionId:    Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *DeleteObjectTaggingResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
		},
	},
}

func TestMockDeleteObjectTagging_Success(t *testing.T) {
	for _, c := range testMockDeleteObjectTaggingSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteObjectTagging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteObjectTaggingErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteObjectTaggingRequest
	CheckOutputFn  func(t *testing.T, o *DeleteObjectTaggingResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging", strUrl)
		},
		&DeleteObjectTaggingRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectTaggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?tagging", strUrl)
		},
		&DeleteObjectTaggingRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectTaggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

func TestMockDeleteObjectTagging_Error(t *testing.T) {
	for _, c := range testMockDeleteObjectTaggingErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteObjectTagging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateSelectObjectMetaSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateSelectObjectMetaRequest
	CheckOutputFn  func(t *testing.T, o *CreateSelectObjectMetaResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7CC12DFE23958A855DED537EF256****\"",
			"Last-Modified":        "Thu, 30 Nov 2023 08:45:29 GMT",
			"x-oss-object-type":    "Normal",
			"x-oss-hash-crc64ecma": "1178629034631640****",
			"x-oss-storage-class":  "Standard",
			"x-oss-version-id":     "CAEQHxiBgMDLj8_94BgiIGRmNmJjNjg0ZjVlYTRiZjdhYzczYjU1NGRiMjY3****",
			"Content-MD5":          "fMEt/iOVioVd7VN+8lYs****",
		},
		[]byte(hexStrToByte("01800006000000250000000000000000000000000000000000000000000000c8000000010000000000000130000000182e46a93f70")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fmeta", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<CsvMetaRequest></CsvMetaRequest>", string(data))
		},
		&CreateSelectObjectMetaRequest{
			Bucket:      Ptr("bucket"),
			Key:         Ptr("object"),
			MetaRequest: &CsvMetaRequest{},
		},
		func(t *testing.T, o *CreateSelectObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.TotalScanned, int64(0))
			assert.Equal(t, o.MetaStatus, 200)
			assert.Equal(t, o.SplitsCount, int32(1))
			assert.Equal(t, o.RowsCount, int64(304))
			assert.Equal(t, o.ColumnsCount, int32(24))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("0180000600000025000000000000000000012efc0000000000012efc000000c8000000010000000000000130000000182ebf898940")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fmeta", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<CsvMetaRequest><InputSerialization><CSV><RecordDelimiter>Cg==</RecordDelimiter><FieldDelimiter>LA==</FieldDelimiter><QuoteCharacter>Ig==</QuoteCharacter></CSV></InputSerialization><OverwriteIfExists>true</OverwriteIfExists></CsvMetaRequest>", string(data))
		},
		&CreateSelectObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			MetaRequest: &CsvMetaRequest{
				OverwriteIfExists: Ptr(true),
				InputSerialization: &InputSerialization{
					CSV: &InputSerializationCSV{
						RecordDelimiter: Ptr("\n"),
						FieldDelimiter:  Ptr(","),
						QuoteCharacter:  Ptr("\""),
					},
				},
			},
		},
		func(t *testing.T, o *CreateSelectObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.TotalScanned, int64(77564))
			assert.Equal(t, o.MetaStatus, 200)
			assert.Equal(t, o.SplitsCount, int32(1))
			assert.Equal(t, o.RowsCount, int64(304))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("01800007000000210000000000000000000000000000000000000000000000c80000000100000000000000642eb96af4fa")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=json%2Fmeta", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<JsonMetaRequest><InputSerialization><JSON><Type>LINES</Type></JSON></InputSerialization></JsonMetaRequest>", string(data))
		},
		&CreateSelectObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			MetaRequest: &JsonMetaRequest{
				InputSerialization: &InputSerialization{
					JSON: &InputSerializationJSON{
						JSONType: Ptr("LINES"),
					},
				},
			},
		},
		func(t *testing.T, o *CreateSelectObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.TotalScanned, int64(0))
			assert.Equal(t, o.MetaStatus, 200)
			assert.Equal(t, o.SplitsCount, int32(1))
			assert.Equal(t, o.RowsCount, int64(100))
			assert.Equal(t, o.ColumnsCount, int32(0))
		},
	},
}

func TestMockCreateSelectObjectMeta_Success(t *testing.T) {
	for _, c := range testMockCreateSelectObjectMetaSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateSelectObjectMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateSelectObjectMetaErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateSelectObjectMetaRequest
	CheckOutputFn  func(t *testing.T, o *CreateSelectObjectMetaResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidJsonData</Code>
  <Message>Invalid value.The current offset is:0.The incorrect json data:'Year,StateA'.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0016-00000801</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000801</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=json%2Fmeta", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<JsonMetaRequest><InputSerialization><JSON><Type>LINES</Type></JSON></InputSerialization></JsonMetaRequest>", string(data))
		},
		&CreateSelectObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			MetaRequest: &JsonMetaRequest{
				InputSerialization: &InputSerialization{
					JSON: &InputSerializationJSON{
						JSONType: Ptr("LINES"),
					},
				},
			},
		},
		func(t *testing.T, o *CreateSelectObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidJsonData", serr.Code)
			assert.Equal(t, "Invalid value.The current offset is:0.The incorrect json data:'Year,StateA'.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "65699DB6E6F906F45A83****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchKey</Code>
  <Message>The specified key does not exist.</Message>
  <RequestId>65699DB6E6F906F45A83****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <Key>object</Key>
  <EC>0026-00000001</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0026-00000001</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=json%2Fmeta", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<JsonMetaRequest><InputSerialization><JSON><Type>LINES</Type></JSON></InputSerialization></JsonMetaRequest>", string(data))
		},
		&CreateSelectObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			MetaRequest: &JsonMetaRequest{
				InputSerialization: &InputSerialization{
					JSON: &InputSerializationJSON{
						JSONType: Ptr("LINES"),
					},
				},
			},
		},
		func(t *testing.T, o *CreateSelectObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "65699DB6E6F906F45A83****", serr.RequestID)
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=json%2Fmeta", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<JsonMetaRequest><InputSerialization><JSON><Type>LINES</Type></JSON></InputSerialization></JsonMetaRequest>", string(data))
		},
		&CreateSelectObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			MetaRequest: &JsonMetaRequest{
				InputSerialization: &InputSerialization{
					JSON: &InputSerializationJSON{
						JSONType: Ptr("LINES"),
					},
				},
			},
		},
		func(t *testing.T, o *CreateSelectObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute CreateSelectObjectMeta fail")
		},
	},
}

func TestMockCreateSelectObjectMeta_Error(t *testing.T) {
	for _, c := range testMockCreateSelectObjectMetaErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateSelectObjectMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSelectObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SelectObjectRequest
	CheckOutputFn  func(t *testing.T, o *SelectObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7CC12DFE23958A855DED537EF256****\"",
			"Last-Modified":        "Thu, 30 Nov 2023 08:45:29 GMT",
			"x-oss-object-type":    "Normal",
			"x-oss-hash-crc64ecma": "1178629034631640****",
			"x-oss-storage-class":  "Standard",
			"x-oss-version-id":     "CAEQHxiBgMDLj8_94BgiIGRmNmJjNjg0ZjVlYTRiZjdhYzczYjU1NGRiMjY3****",
			"Content-MD5":          "fMEt/iOVioVd7VN+8lYs****",
		},
		[]byte(hexStrToByte("0180000100000044000000000000000000012dfb323031352c55532c2c4869676820426c6f6f642050726573737572650d0a323031352c55532c2c4869676820426c6f6f642050726573737572650d0a000000000180000500000014000000000000000000012efc0000000000012efc000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IFllYXIsIFN0YXRlQWJiciwgQ2l0eU5hbWUsIFNob3J0X1F1ZXN0aW9uX1RleHQgZnJvbSBvc3NvYmplY3Qgd2hlcmUgTWVhc3VyZSBsaWtlICclYmxvb2QgcHJlc3N1cmUlWWVhcnMn</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select Year, StateAbbr, CityName, Short_Question_Text from ossobject where Measure like '%blood pressure%Years'"),
				InputSerializationSelect: InputSerializationSelect{
					CsvBodyInput: &CSVSelectInput{
						FileHeaderInfo: Ptr("Use"),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, err := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "2015,US,,High Blood Press")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("0180000100012e03000000000000000000012dfb596561722c5374617465416262722c5374617465446573632c436974794e616d652c47656f677261706869634c6576656c2c44617461536f757263652c43617465676f72792c556e6971756549442c4d6561737572652c446174615f56616c75655f556e69742c4461746156616c75655479706549442c446174615f56616c75655f547970652c446174615f56616c75652c4c6f775f436f6e666964656e63655f4c696d69742c486967685f436f6e666964656e63655f4c696d69742c446174615f56616c75655f466f6f746e6f74655f53796d626f6c2c446174615f56616c75655f466f6f746e6f74652c506f70756c6174696f6e436f756e742c47656f4c6f636174696f6e2c43617465676f727949442c4d65617375726549642c43697479464950532c5472616374464950532c53686f72745f5175657374696f6e5f546578740d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31352e342c31352e312c31352e372c2c2c3330383734353533382c2c50524556454e542c414343455353322c2c2c4865616c746820496e737572616e63650d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e382c31342e352c31352e302c2c2c3330383734353533382c2c50524556454e542c414343455353322c2c2c4865616c746820496e737572616e63650d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c32322e352c32322e332c32322e372c2c2c3330383734353533382c2c484c54484f55542c4152544852495449532c2c2c4172746872697469730d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32342e372c32342e352c32342e392c2c2c3330383734353533382c2c484c54484f55542c4152544852495449532c2c2c4172746872697469730d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31372e322c31362e392c31372e342c2c2c3330383734353533382c2c554e484245482c42494e47452c2c2c42696e6765204472696e6b696e670d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31362e332c31362e312c31362e352c2c2c3330383734353533382c2c554e484245482c42494e47452c2c2c42696e6765204472696e6b696e670d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4869676820626c6f6f6420707265737375726520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c32392e342c32392e322c32392e372c2c2c3330383734353533382c2c484c54484f55542c4250484947482c2c2c4869676820426c6f6f642050726573737572650d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4869676820626c6f6f6420707265737375726520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33312e392c33312e362c33322e322c2c2c3330383734353533382c2c484c54484f55542c4250484947482c2c2c4869676820426c6f6f642050726573737572650d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c54616b696e67206d65646963696e6520666f72206869676820626c6f6f6420707265737375726520636f6e74726f6c20616d6f6e67206164756c74732061676564203e3d31382059656172732077697468206869676820626c6f6f642070726573737572652c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c35372e372c35372e312c35382e342c2c2c3330383734353533382c2c50524556454e542c42504d45442c2c2c54616b696e67204250204d656469636174696f6e0d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c54616b696e67206d65646963696e6520666f72206869676820626c6f6f6420707265737375726520636f6e74726f6c20616d6f6e67206164756c74732061676564203e3d31382059656172732077697468206869676820626c6f6f642070726573737572652c252c4372645072762c43727564652070726576616c656e63652c37372e322c37362e382c37372e372c2c2c3330383734353533382c2c50524556454e542c42504d45442c2c2c54616b696e67204250204d656469636174696f6e0d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c43616e63657220286578636c7564696e6720736b696e2063616e6365722920616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c362e302c352e392c362e312c2c2c3330383734353533382c2c484c54484f55542c43414e4345522c2c2c43616e636572202865786365707420736b696e290d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c43616e63657220286578636c7564696e6720736b696e2063616e6365722920616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c362e362c362e352c362e382c2c2c3330383734353533382c2c484c54484f55542c43414e4345522c2c2c43616e636572202865786365707420736b696e290d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c43757272656e7420617374686d6120616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c382e372c382e362c382e392c2c2c3330383734353533382c2c484c54484f55542c43415354484d412c2c2c43757272656e7420417374686d610d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c43757272656e7420617374686d6120616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e382c382e362c392e302c2c2c3330383734353533382c2c484c54484f55542c43415354484d412c2c2c43757272656e7420417374686d610d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c436f726f6e617279206865617274206469736561736520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c352e362c352e352c352e382c2c2c3330383734353533382c2c484c54484f55542c4348442c2c2c436f726f6e61727920486561727420446973656173650d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c436f726f6e617279206865617274206469736561736520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c362e332c362e322c362e352c2c2c3330383734353533382c2c484c54484f55542c4348442c2c2c436f726f6e61727920486561727420446973656173650d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c56697369747320746f20646f63746f7220666f7220726f7574696e6520636865636b75702077697468696e207468652070617374205965617220616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c36382e362c36382e332c36382e392c2c2c3330383734353533382c2c50524556454e542c434845434b55502c2c2c416e6e75616c20436865636b75700d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c56697369747320746f20646f63746f7220666f7220726f7574696e6520636865636b75702077697468696e207468652070617374205965617220616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c37302e302c36392e372c37302e332c2c2c3330383734353533382c2c50524556454e542c434845434b55502c2c2c416e6e75616c20436865636b75700d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c43686f6c65737465726f6c2073637265656e696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c37352e322c37342e392c37352e352c2c2c3330383734353533382c2c50524556454e542c43484f4c53435245454e2c2c2c43686f6c65737465726f6c2053637265656e696e670d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c43686f6c65737465726f6c2073637265656e696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c37372e302c37362e372c37372e332c2c2c3330383734353533382c2c50524556454e542c43484f4c53435245454e2c2c2c43686f6c65737465726f6c2053637265656e696e670d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c22466563616c206f6363756c7420626c6f6f6420746573742c207369676d6f69646f73636f70792c206f7220636f6c6f6e6f73636f707920616d6f6e67206164756c74732061676564203530e280933735205965617273222c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c36342e302c36332e352c36342e352c2c2c3330383734353533382c2c50524556454e542c434f4c4f4e5f53435245454e2c2c2c436f6c6f72656374616c2043616e6365722053637265656e696e670d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c22466563616c206f6363756c7420626c6f6f6420746573742c207369676d6f69646f73636f70792c206f7220636f6c6f6e6f73636f707920616d6f6e67206164756c74732061676564203530e280933735205965617273222c252c4372645072762c43727564652070726576616c656e63652c36332e372c36332e332c36342e312c2c2c3330383734353533382c2c50524556454e542c434f4c4f4e5f53435245454e2c2c2c436f6c6f72656374616c2043616e6365722053637265656e696e670d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4368726f6e6963206f627374727563746976652070756c6d6f6e617279206469736561736520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c352e372c352e362c352e392c2c2c3330383734353533382c2c484c54484f55542c434f50442c2c2c434f50440d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4368726f6e6963206f627374727563746976652070756c6d6f6e617279206469736561736520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c362e332c362e322c362e342c2c2c3330383734353533382c2c484c54484f55542c434f50442c2c2c434f50440d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c506879736963616c206865616c7468206e6f7420676f6f6420666f72203e3d3134206461797320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31312e352c31312e332c31312e372c2c2c3330383734353533382c2c484c54484f55542c50484c54482c2c2c506879736963616c204865616c74680d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c224f6c646572206164756c74206d656e2061676564203e3d36352059656172732077686f2061726520757020746f2064617465206f6e206120636f726520736574206f6620636c696e6963616c2070726576656e746976652073657276696365733a20466c752073686f74207061737420596561722c205050562073686f7420657665722c20436f6c6f72656374616c2063616e6365722073637265656e696e67222c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c33322e392c33322e312c33332e362c2c2c3330383734353533382c2c50524556454e542c434f52454d2c2c2c436f72652070726576656e7469766520736572766963657320666f72206f6c646572206d656e0d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c224f6c646572206164756c74206d656e2061676564203e3d36352059656172732077686f2061726520757020746f2064617465206f6e206120636f726520736574206f6620636c696e6963616c2070726576656e746976652073657276696365733a20466c752073686f74207061737420596561722c205050562073686f7420657665722c20436f6c6f72656374616c2063616e6365722073637265656e696e67222c252c4372645072762c43727564652070726576616c656e63652c33322e332c33312e352c33332e302c2c2c3330383734353533382c2c50524556454e542c434f52454d2c2c2c436f72652070726576656e7469766520736572766963657320666f72206f6c646572206d656e0d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c224f6c646572206164756c7420776f6d656e2061676564203e3d36352059656172732077686f2061726520757020746f2064617465206f6e206120636f726520736574206f6620636c696e6963616c2070726576656e746976652073657276696365733a20466c752073686f74207061737420596561722c205050562073686f7420657665722c20436f6c6f72656374616c2063616e6365722073637265656e696e672c20616e64204d616d6d6f6772616d20706173742032205965617273222c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c33302e372c33302e322c33312e342c2c2c3330383734353533382c2c50524556454e542c434f5245572c2c2c436f72652070726576656e7469766520736572766963657320666f72206f6c64657220776f6d656e0d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c224f6c646572206164756c7420776f6d656e2061676564203e3d36352059656172732077686f2061726520757020746f2064617465206f6e206120636f726520736574206f6620636c696e6963616c2070726576656e746976652073657276696365733a20466c752073686f74207061737420596561722c205050562073686f7420657665722c20436f6c6f72656374616c2063616e6365722073637265656e696e672c20616e64204d616d6d6f6772616d20706173742032205965617273222c252c4372645072762c43727564652070726576616c656e63652c33302e372c33302e312c33312e332c2c2c3330383734353533382c2c50524556454e542c434f5245572c2c2c436f72652070726576656e7469766520736572766963657320666f72206f6c64657220776f6d656e0d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c43757272656e7420736d6f6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31372e312c31362e382c31372e332c2c2c3330383734353533382c2c554e484245482c43534d4f4b494e472c2c2c43757272656e7420536d6f6b696e670d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c43757272656e7420736d6f6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31362e382c31362e362c31372e302c2c2c3330383734353533382c2c554e484245482c43534d4f4b494e472c2c2c43757272656e7420536d6f6b696e670d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c56697369747320746f2064656e74697374206f722064656e74616c20636c696e696320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c36342e312c36332e382c36342e342c2c2c3330383734353533382c2c50524556454e542c44454e54414c2c2c2c44656e74616c2056697369740d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c56697369747320746f2064656e74697374206f722064656e74616c20636c696e696320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c36342e342c36342e312c36342e372c2c2c3330383734353533382c2c50524556454e542c44454e54414c2c2c2c44656e74616c2056697369740d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c446961676e6f73656420646961626574657320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c392e332c392e322c392e352c2c2c3330383734353533382c2c484c54484f55542c44494142455445532c2c2c44696162657465730d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c446961676e6f73656420646961626574657320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e342c31302e332c31302e362c2c2c3330383734353533382c2c484c54484f55542c44494142455445532c2c2c44696162657465730d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c486967682063686f6c65737465726f6c20616d6f6e67206164756c74732061676564203e3d31382059656172732077686f2068617665206265656e2073637265656e656420696e20746865207061737420352059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c33312e312c33302e382c33312e342c2c2c3330383734353533382c2c484c54484f55542c4849474843484f4c2c2c2c486967682043686f6c65737465726f6c0d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c486967682063686f6c65737465726f6c20616d6f6e67206164756c74732061676564203e3d31382059656172732077686f2068617665206265656e2073637265656e656420696e20746865207061737420352059656172732c252c4372645072762c43727564652070726576616c656e63652c33372e312c33362e382c33372e342c2c2c3330383734353533382c2c484c54484f55542c4849474843484f4c2c2c2c486967682043686f6c65737465726f6c0d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4368726f6e6963206b69646e6579206469736561736520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c322e352c322e342c322e362c2c2c3330383734353533382c2c484c54484f55542c4b49444e45592c2c2c4368726f6e6963204b69646e657920446973656173650d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4368726f6e6963206b69646e6579206469736561736520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c322e372c322e362c322e382c2c2c3330383734353533382c2c484c54484f55542c4b49444e45592c2c2c4368726f6e6963204b69646e657920446973656173650d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c4e6f206c6569737572652d74696d6520706879736963616c20616374697669747920616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c32352e352c32352e322c32352e382c2c2c3330383734353533382c2c554e484245482c4c50412c2c2c506879736963616c20496e61637469766974790d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c4e6f206c6569737572652d74696d6520706879736963616c20616374697669747920616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32352e392c32352e362c32362e312c2c2c3330383734353533382c2c554e484245482c4c50412c2c2c506879736963616c20496e61637469766974790d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c4d616d6d6f6772617068792075736520616d6f6e6720776f6d656e2061676564203530e2809337342059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c37352e352c37352e312c37352e392c2c2c3330383734353533382c2c50524556454e542c4d414d4d4f5553452c2c2c4d616d6d6f6772617068790d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c4d616d6d6f6772617068792075736520616d6f6e6720776f6d656e2061676564203530e2809337342059656172732c252c4372645072762c43727564652070726576616c656e63652c37352e382c37352e342c37362e322c2c2c3330383734353533382c2c50524556454e542c4d414d4d4f5553452c2c2c4d616d6d6f6772617068790d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4d656e74616c206865616c7468206e6f7420676f6f6420666f72203e3d3134206461797320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31312e362c31312e342c31312e382c2c2c3330383734353533382c2c484c54484f55542c4d484c54482c2c2c4d656e74616c204865616c74680d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4d656e74616c206865616c7468206e6f7420676f6f6420666f72203e3d3134206461797320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31312e342c31312e332c31312e362c2c2c3330383734353533382c2c484c54484f55542c4d484c54482c2c2c4d656e74616c204865616c74680d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c4f62657369747920616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c32382e372c32382e342c32392e302c2c2c3330383734353533382c2c554e484245482c4f4245534954592c2c2c4f6265736974790d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c4f62657369747920616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32382e382c32382e362c32392e312c2c2c3330383734353533382c2c554e484245482c4f4245534954592c2c2c4f6265736974790d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c506170616e69636f6c616f7520736d6561722075736520616d6f6e67206164756c7420776f6d656e2061676564203231e2809336352059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c38312e312c38302e362c38312e362c2c2c3330383734353533382c2c50524556454e542c504150544553542c2c2c50617020536d65617220546573740d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c506170616e69636f6c616f7520736d6561722075736520616d6f6e67206164756c7420776f6d656e2061676564203231e2809336352059656172732c252c4372645072762c43727564652070726576616c656e63652c38312e382c38312e332c38322e322c2c2c3330383734353533382c2c50524556454e542c504150544553542c2c2c50617020536d65617220546573740d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c506879736963616c206865616c7468206e6f7420676f6f6420666f72203e3d3134206461797320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31322e302c31312e382c31322e322c2c2c3330383734353533382c2c484c54484f55542c50484c54482c2c2c506879736963616c204865616c74680d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c536c656570696e67206c657373207468616e203720686f75727320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c33352e312c33342e382c33352e352c2c2c3330383734353533382c2c554e484245482c534c4545502c2c2c536c656570203c203720686f7572730d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c536c656570696e67206c657373207468616e203720686f75727320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33342e382c33342e352c33352e312c2c2c3330383734353533382c2c554e484245482c534c4545502c2c2c536c656570203c203720686f7572730d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c5374726f6b6520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c322e382c322e372c322e382c2c2c3330383734353533382c2c484c54484f55542c5354524f4b452c2c2c5374726f6b650d0a323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c5374726f6b6520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c332e302c332e302c332e312c2c2c3330383734353533382c2c484c54484f55542c5354524f4b452c2c2c5374726f6b650d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c416c6c207465657468206c6f737420616d6f6e67206164756c74732061676564203e3d36352059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31352e342c31352e302c31352e382c2c2c3330383734353533382c2c484c54484f55542c54454554484c4f53542c2c2c5465657468204c6f73730d0a323031342c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c416c6c207465657468206c6f737420616d6f6e67206164756c74732061676564203e3d36352059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e392c31342e362c31352e332c2c2c3330383734353533382c2c484c54484f55542c54454554484c4f53542c2c2c5465657468204c6f73730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c436974792c42524653532c50726576656e74696f6e2c303130373030302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31392e382c31392e352c32302e322c2c2c3231323233372c222833332e353237353636333737332c202d38362e3739383831373436373829222c50524556454e542c414343455353322c303130373030302c2c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c436974792c42524653532c50726576656e74696f6e2c303130373030302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31392e362c31392e322c32302e302c2c2c3231323233372c222833332e353237353636333737332c202d38362e3739383831373436373829222c50524556454e542c414343455353322c303130373030302c2c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030303130302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32332e392c32312e322c32372e322c2c2c333034322c222833332e353739343332383332362c202d38362e3732323833323339323629222c50524556454e542c414343455353322c303130373030302c30313037333030303130302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030303330302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32382e382c32352e342c33322e342c2c2c323733352c222833332e353432383230383638362c202d38362e37353234333339373829222c50524556454e542c414343455353322c303130373030302c30313037333030303330302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030303430302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32362e312c32322e362c32392e392c2c2c333333382c222833332e353633323434393633332c202d38362e3736343034373430363429222c50524556454e542c414343455353322c303130373030302c30313037333030303430302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030303530302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32382e312c32342e362c33322e302c2c2c323836342c222833332e353434323430343539342c202d38362e3737343931333037313929222c50524556454e542c414343455353322c303130373030302c30313037333030303530302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030303730302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c33312e382c32372e302c33362e372c2c2c323537372c222833332e353532353430363133392c202d38362e3830313638393337303629222c50524556454e542c414343455353322c303130373030302c30313037333030303730302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030303830302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32322e342c31392e312c32362e312c2c2c333835392c222833332e3534393639373738392c202d38362e3833333039343437343429222c50524556454e542c414343455353322c303130373030302c30313037333030303830302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030313130302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31362e382c31332e372c32302e352c2c2c353335342c222833332e353432393134333332352c202d38362e3837353637383238353229222c50524556454e542c414343455353322c303130373030302c30313037333030313130302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030313230302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32342e362c32322e322c32372e322c2c2c323837362c222833332e353237383736373730362c202d38362e3836303431363136383629222c50524556454e542c414343455353322c303130373030302c30313037333030313230302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030313430302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32322e302c31382e342c32352e372c2c2c323138312c222833332e353236313439373235382c202d38362e38333531343636303629222c50524556454e542c414343455353322c303130373030302c30313037333030313430302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030313530302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32362e332c32332e312c32392e342c2c2c333138392c222833332e353239383732373334322c202d38362e3831393731393136383529222c50524556454e542c414343455353322c303130373030302c30313037333030313530302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030313630302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32362e382c32322e392c33302e382c2c2c333339302c222833332e353337323939333432332c202d38362e3830333635393034383229222c50524556454e542c414343455353322c303130373030302c30313037333030313630302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030313930322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32362e392c32342e312c32392e382c2c2c313839342c222833332e353533323035303939372c202d38362e3734323938303136303329222c50524556454e542c414343455353322c303130373030302c30313037333030313930322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030323030302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32342e342c32302e382c32382e322c2c2c333838352c222833332e353534313537343130362c202d38362e3731363732323939313529222c50524556454e542c414343455353322c303130373030302c30313037333030323030302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030323130302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32312e352c31382e312c32352e302c2c2c333138362c222833332e353635303031353934322c202d38362e3731303130323437363629222c50524556454e542c414343455353322c303130373030302c30313037333030323130302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030323230302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32312e372c31382e372c32352e302c2c2c323633302c222833332e353532313330313230352c202d38362e3732373637353935303829222c50524556454e542c414343455353322c303130373030302c30313037333030323230302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030323330332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32372e302c32332e372c33302e372c2c2c323933362c222833332e353338333135333230372c202d38362e3732373034343534323829222c50524556454e542c414343455353322c303130373030302c30313037333030323330332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030323330352c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c392e392c382e322c31322e312c2c2c323935322c222833332e353333333431353937362c202d38362e3734373935363630383429222c50524556454e542c414343455353322c303130373030302c30313037333030323330352c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030323330362c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c392e352c382e312c31312e322c2c2c333235372c222833332e353231333837333536342c202d38362e3734393030333132383929222c50524556454e542c414343455353322c303130373030302c30313037333030323330362c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030323430302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32352e312c32322e342c32372e382c2c2c333632392c222833332e353236303734383330392c202d38362e3738333033313534383829222c50524556454e542c414343455353322c303130373030302c30313037333030323430302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030323730302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32302e372c31362e382c32342e392c2c2c333939322c222833332e353137363030383431392c202d38362e3831303638383734353229222c50524556454e542c414343455353322c303130373030302c30313037333030323730302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030323930302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32352e372c32312e352c32392e392c2c2c323036342c222833332e353133323439383836342c202d38362e383330303437343929222c50524556454e542c414343455353322c303130373030302c30313037333030323930302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333030312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31382e342c31352e332c32322e342c2c2c333737392c222833332e353132353135383039342c202d38362e3835373731363439343629222c50524556454e542c414343455353322c303130373030302c30313037333030333030312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333030322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32362e332c32312e372c33312e332c2c2c323230332c222833332e3531323235383130392c202d38362e3834343134333939303729222c50524556454e542c414343455353322c303130373030302c30313037333030333030322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333130302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32332e322c32302e322c32362e372c2c2c333633372c222833332e353035393635353735362c202d38362e3837343535303630383629222c50524556454e542c414343455353322c303130373030302c30313037333030333130302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333230302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32382e372c32352e302c33322e382c2c2c3933312c222833332e353039343031383530322c202d38362e3838353930383139363129222c50524556454e542c414343455353322c303130373030302c30313037333030333230302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333330302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32322e322c31392e312c32352e372c2c2c3934372c222833332e353137313236313130382c202d38362e3839313338313937343929222c50524556454e542c414343455353322c303130373030302c30313037333030333330302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333430302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32362e382c32332e302c33302e392c2c2c323437372c222833332e353035323232393233342c202d38362e3930313438343436353629222c50524556454e542c414343455353322c303130373030302c30313037333030333430302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333530302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32322e332c31392e312c32352e372c2c2c323738302c222833332e353036353731343031312c202d38362e3931393539313030363329222c50524556454e542c414343455353322c303130373030302c30313037333030333530302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333630302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31372e302c31342e352c31392e362c2c2c343638332c222833332e34383437363339372c202d38362e3839383133393239343729222c50524556454e542c414343455353322c303130373030302c30313037333030333630302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333730302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32312e362c31392e322c32332e392c2c2c353036332c222833332e343936393031383538392c202d38362e3839303737323934323629222c50524556454e542c414343455353322c303130373030302c30313037333030333730302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333830322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31392e302c31362e372c32312e362c2c2c353430392c222833332e343738353730373739342c202d38362e38393030303039303729222c50524556454e542c414343455353322c303130373030302c30313037333030333830322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333830332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32302e392c31372e372c32342e322c2c2c343139392c222833332e3438353934353231342c202d38362e38363936393231383629222c50524556454e542c414343455353322c303130373030302c30313037333030333830332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030333930302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32392e392c32362e392c33332e302c2c2c313738332c222833332e343938393935393332372c202d38362e3836343736303030333829222c50524556454e542c414343455353322c303130373030302c30313037333030333930302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030343030302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32352e382c32312e342c33302e332c2c2c333737322c222833332e343935333234363031352c202d38362e3835313632333230373329222c50524556454e542c414343455353322c303130373030302c30313037333030343030302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030343230302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32332e342c32302e382c32362e312c2c2c323334312c222833332e353030373433393336312c202d38362e3832373037323033373929222c50524556454e542c414343455353322c303130373030302c30313037333030343230302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030343530302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32322e372c31392e332c32372e342c2c2c353030332c222833332e353034313835373535362c202d38362e3830333337393833343629222c50524556454e542c414343455353322c303130373030302c30313037333030343530302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030343730312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e362c382e362c31332e342c2c2c333438302c222833332e353037353234323134382c202d38362e3738333636373538333829222c50524556454e542c414343455353322c303130373030302c30313037333030343730312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030343730322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c372e312c362e312c382e372c2c2c323934342c222833332e353131393930323636312c202d38362e3736393435353039383929222c50524556454e542c414343455353322c303130373030302c30313037333030343730322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030343830302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c382e372c372e332c31302e332c2c2c313836312c222833332e343938393036343030382c202d38362e373832363939313429222c50524556454e542c414343455353322c303130373030302c30313037333030343830302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030343930312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31312e302c392e312c31332e322c2c2c313136372c222833332e343937313539353634352c202d38362e3739313734343036363829222c50524556454e542c414343455353322c303130373030302c30313037333030343930312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030343930322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31352e322c31322e352c31382e352c2c2c333134362c222833332e343933353832343034332c202d38362e3830303932393436303329222c50524556454e542c414343455353322c303130373030302c30313037333030343930322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353030302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31372e392c31352e342c32302e372c2c2c333438322c222833332e343836363638393739352c202d38362e3831373332363238333129222c50524556454e542c414343455353322c303130373030302c30313037333030353030302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353130312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c33302e352c32362e352c33342e342c2c2c313530372c222833332e343934353930393030382c202d38362e38333437363339333629222c50524556454e542c414343455353322c303130373030302c30313037333030353130312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353130332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31372e352c31342e342c32302e372c2c2c323538372c222833332e3438353731343838352c202d38362e3833323738313734363729222c50524556454e542c414343455353322c303130373030302c30313037333030353130332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353130342c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32342e392c32312e332c32382e392c2c2c323838312c222833332e343734393934313831362c202d38362e3833333534323137343729222c50524556454e542c414343455353322c303130373030302c30313037333030353130342c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353230302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32302e322c31362e342c32342e322c2c2c333734302c222833332e343830363730383737352c202d38362e3835303836373135313429222c50524556454e542c414343455353322c303130373030302c30313037333030353230302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353330322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31322e382c31302e392c31352e312c2c2c333436332c222833332e353736363435363434392c202d38362e3639363535393033313629222c50524556454e542c414343455353322c303130373030302c30313037333030353330322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353530302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32362e352c32322e392c33302e372c2c2c313832342c222833332e353637303239333839382c202d38362e3830303535363732313329222c50524556454e542c414343455353322c303130373030302c30313037333030353530302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353630302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31322e332c31302e302c31352e302c2c2c343336372c222833332e353230303938313233342c202d38362e3732373230363331393829222c50524556454e542c414343455353322c303130373030302c30313037333030353630302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353730312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31382e392c31352e372c32322e332c2c2c323337322c222833332e343632393534333332392c202d38362e3838393834303636323829222c50524556454e542c414343455353322c303130373030302c30313037333030353730312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353730322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31392e382c31362e342c32332e332c2c2c333431332c222833332e343639383639313338352c202d38362e38373432393036303229222c50524556454e542c414343455353322c303130373030302c30313037333030353730322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353830302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31392e332c31362e342c32322e332c2c2c343231362c222833332e3437393036323332352c202d38362e3831323830363330383929222c50524556454e542c414343455353322c303130373030302c30313037333030353830302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353930332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31352e392c31332e392c31382e302c2c2c343933332c222833332e3539373136323330392c202d38362e3637363637333633353129222c50524556454e542c414343455353322c303130373030302c30313037333030353930332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353930352c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32312e342c31382e322c32342e372c2c2c353033392c222833332e3630333938383435362c202d38362e3730303831323334313829222c50524556454e542c414343455353322c303130373030302c30313037333030353930352c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353930372c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31322e352c31302e332c31342e392c2c2c313937352c222833332e363134323836313530312c202d38362e3636393139393634313729222c50524556454e542c414343455353322c303130373030302c30313037333030353930372c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353930382c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32302e392c31382e312c32332e372c2c2c313632312c222833332e363138323932343433322c202d38362e3638303134383330383429222c50524556454e542c414343455353322c303130373030302c30313037333030353930382c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353930392c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e342c31312e382c31372e352c2c2c323532342c222833332e3631313831313737332c202d38362e3732313432323135313429222c50524556454e542c414343455353322c303130373030302c30313037333030353930392c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333030353931302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e372c31322e362c31372e312c2c2c343631322c222833332e363239393031373439392c202d38362e3731393433313132323929222c50524556454e542c414343455353322c303130373030302c30313037333030353931302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031303530302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32312e382c31382e322c32352e392c2c2c3131342c222833332e343336333738363830362c202d38362e3931323839323330373229222c50524556454e542c414343455353322c303130373030302c30313037333031303530302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031303730312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31382e352c31342e342c32332e362c2c2c37342c222833332e3437333838363135352c202d38362e3831343634383737363229222c50524556454e542c414343455353322c303130373030302c30313037333031303730312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031303730362c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31332e352c31312e302c31362e332c2c2c313532382c222833332e343434333730393434322c202d38362e3834303533353236343529222c50524556454e542c414343455353322c303130373030302c30313037333031303730362c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031303830312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c352e322c342e342c362e332c2c2c3136382c222833332e3531343039373835332c202d38362e3734363639373133363229222c50524556454e542c414343455353322c303130373030302c30313037333031303830312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031303830322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c352e302c332e392c362e352c2c2c3137322c222833332e343838353439333437372c202d38362e37383038343330323429222c50524556454e542c414343455353322c303130373030302c30313037333031303830322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031303830332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31312e332c392e302c31342e302c2c2c3531342c222833332e353232393039333839322c202d38362e3731303236313836343229222c50524556454e542c414343455353322c303130373030302c30313037333031303830332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031303830352c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c372e332c352e352c392e392c2c2c38362c222833332e343935323739323437322c202d38362e3639383731383439373429222c50524556454e542c414343455353322c303130373030302c30313037333031303830352c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313130342c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31332e362c31312e342c31362e312c2c2c313638382c222833332e363135393433363433332c202d38362e3635353738393235303729222c50524556454e542c414343455353322c303130373030302c30313037333031313130342c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313130372c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c34322c222833332e353830343834353234392c202d38362e3633303131313039363129222c50524556454e542c414343455353322c303130373030302c30313037333031313130372c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313130382c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c392c222833332e363035303734323539362c202d38362e3633313637323933383629222c50524556454e542c414343455353322c303130373030302c30313037333031313130382c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313230372c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31372e372c31352e312c32302e362c2c2c3831352c222833332e363731383834383730362c202d38362e3637373235313034363529222c50524556454e542c414343455353322c303130373030302c30313037333031313230372c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313230392c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32302e322c31372e322c32342e322c2c2c313036322c222833332e363535373138393939322c202d38362e3730353036393833343929222c50524556454e542c414343455353322c303130373030302c30313037333031313230392c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313231302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32312e352c31372e302c32362e372c2c2c313338352c222833332e363634313839333735352c202d38362e3639353631373036383629222c50524556454e542c414343455353322c303130373030302c30313037333031313231302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313830332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31382e312c31352e342c32312e332c2c2c3932382c222833332e363235323537353137332c202d38362e3639393836313034303929222c50524556454e542c414343455353322c303130373030302c30313037333031313830332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313830342c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31382e332c31352e322c32312e392c2c2c313135372c222833332e363437343931363137352c202d38362e3730343239373434323429222c50524556454e542c414343455353322c303130373030302c30313037333031313830342c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313930312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c362c222833332e363335353431343634362c202d38362e37333639343636393129222c50524556454e542c414343455353322c303130373030302c30313037333031313930312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031313930342c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e342c31322e312c31362e382c2c2c313931352c222833332e3539333134303138352c202d38362e3733353739333035343129222c50524556454e542c414343455353322c303130373030302c30313037333031313930342c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323030312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e352c31322e352c31362e372c2c2c3330342c222833332e353931393438363131322c202d38362e38363433383430363829222c50524556454e542c414343455353322c303130373030302c30313037333031323030312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323030322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c34342c222833332e353835393233343139372c202d38362e3833353730303731383829222c50524556454e542c414343455353322c303130373030302c30313037333031323030322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323230302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c32332c222833332e353936373434313930342c202d38372e3038373933393638353729222c50524556454e542c414343455353322c303130373030302c30313037333031323230302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323330322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31322e372c31302e362c31352e322c2c2c3134342c222833332e353534323831363335322c202d38372e3035343436393134313629222c50524556454e542c414343455353322c303130373030302c30313037333031323330322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323330352c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31322e352c31302e322c31352e312c2c2c3430332c222833332e343639353335383036342c202d38362e393638333137333929222c50524556454e542c414343455353322c303130373030302c30313037333031323330352c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323430312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31312e342c382e392c31342e322c2c2c313036362c222833332e353537313935313034382c202d38362e3837373739333530343929222c50524556454e542c414343455353322c303130373030302c30313037333031323430312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323430322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31362e372c31342e362c31392e302c2c2c3431382c222833332e3534393938343137322c202d38362e3839393435343333323729222c50524556454e542c414343455353322c303130373030302c30313037333031323430322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323530302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32312e302c31372e362c32342e342c2c2c3431302c222833332e3532393136303438362c202d38362e3933343634373634343529222c50524556454e542c414343455353322c303130373030302c30313037333031323530302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323630322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31372e352c31352e372c31392e362c2c2c3337312c222833332e3537303136343637342c202d38362e36363634333035383229222c50524556454e542c414343455353322c303130373030302c30313037333031323630322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323730312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c34342c222833332e353438343037383037312c202d38362e3633323337373334353529222c50524556454e542c414343455353322c303130373030302c30313037333031323730312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323730332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31312e332c382e362c31342e392c2c2c3439382c222833332e343638313138303934332c202d38362e3636373138383832313329222c50524556454e542c414343455353322c303130373030302c30313037333031323730332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323730342c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c382e392c362e382c31312e352c2c2c3131332c222833332e353033343139353930382c202d38362e3631383039383334303329222c50524556454e542c414343455353322c303130373030302c30313037333031323730342c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323830332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c382e362c362e372c31312e322c2c2c313236312c222833332e343433393432353836352c202d38362e3732313239333639333829222c50524556454e542c414343455353322c303130373030302c30313037333031323830332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031323931302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c392c222833332e343334353830353034322c202d38362e3732363332393230353929222c50524556454e542c414343455353322c303130373030302c30313037333031323931302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031333030322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32312e332c31372e312c32352e362c2c2c313531342c222833332e34363630343138312c202d38362e3835363732383737393729222c50524556454e542c414343455353322c303130373030302c30313037333031333030322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031333130302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32312e382c31382e382c32342e392c2c2c343432342c222833332e34343838303231342c202d38362e3838373834303135373929222c50524556454e542c414343455353322c303130373030302c30313037333031333130302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031333330302c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c32332e352c31392e352c32372e352c2c2c313738322c222833332e343339363434333536392c202d38362e3932343837363836363529222c50524556454e542c414343455353322c303130373030302c30313037333031333330302c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031333930312c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31382e302c31342e382c32312e372c2c2c3935322c222833332e343732393338343930362c202d38362e3935343733333736343829222c50524556454e542c414343455353322c303130373030302c30313037333031333930312c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031343330322c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e302c382e332c31322e312c2c2c323737382c222833332e343234343635383832392c202d38362e3838343134373432313729222c50524556454e542c414343455353322c303130373030302c30313037333031343330322c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313037333031343431332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c372e362c352e392c392e382c2c2c3339372c222833332e343232363539333131372c202d38362e3835303836323037353129222c50524556454e542c414343455353322c303130373030302c30313037333031343431332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313131373033303231332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c392e332c372e372c31312e342c2c2c3634342c222833332e343339353937353139332c202d38362e3637333539353933353929222c50524556454e542c414343455353322c303130373030302c30313131373033303231332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313131373033303231372c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c31362c222833332e343535363939353736332c202d38362e3635323032303836333929222c50524556454e542c414343455353322c303130373030302c30313131373033303231372c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c50726576656e74696f6e2c303130373030302d30313131373033303330332c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c392e392c382e312c31322e302c2c2c3936382c222833332e343235383636313233392c202d38362e37313338313933353629222c50524556454e542c414343455353322c303130373030302c30313131373033303330332c4865616c746820496e737572616e63650d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c436974792c42524653532c4865616c7468204f7574636f6d65732c303130373030302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c33312e302c33302e382c33312e312c2c2c3231323233372c222833332e353237353636333737332c202d38362e3739383831373436373829222c484c54484f55542c4152544852495449532c303130373030302c2c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c436974792c42524653532c4865616c7468204f7574636f6d65732c303130373030302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33302e392c33302e382c33312e312c2c2c3231323233372c222833332e353237353636333737332c202d38362e3739383831373436373829222c484c54484f55542c4152544852495449532c303130373030302c2c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030303130302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33322e352c33312e352c33332e362c2c2c333034322c222833332e353739343332383332362c202d38362e3732323833323339323629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030303130302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030303330302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33312e332c33302e302c33322e342c2c2c323733352c222833332e353432383230383638362c202d38362e37353234333339373829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030303330302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030303430302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33342e362c33332e322c33352e392c2c2c333333382c222833332e353633323434393633332c202d38362e3736343034373430363429222c484c54484f55542c4152544852495449532c303130373030302c30313037333030303430302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030303530302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33372e382c33362e332c33392e322c2c2c323836342c222833332e353434323430343539342c202d38362e3737343931333037313929222c484c54484f55542c4152544852495449532c303130373030302c30313037333030303530302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030303730302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33382e352c33372e312c33392e392c2c2c323537372c222833332e353532353430363133392c202d38362e3830313638393337303629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030303730302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030303830302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33382e302c33362e352c33392e332c2c2c333835392c222833332e3534393639373738392c202d38362e3833333039343437343429222c484c54484f55542c4152544852495449532c303130373030302c30313037333030303830302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030313130302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33342e302c33322e342c33352e352c2c2c353335342c222833332e353432393134333332352c202d38362e3837353637383238353229222c484c54484f55542c4152544852495449532c303130373030302c30313037333030313130302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030313230302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33362e352c33352e342c33372e362c2c2c323837362c222833332e353237383736373730362c202d38362e3836303431363136383629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030313230302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030313430302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33372e312c33352e362c33382e362c2c2c323138312c222833332e353236313439373235382c202d38362e38333531343636303629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030313430302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030313530302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33352e322c33342e302c33362e342c2c2c333138392c222833332e353239383732373334322c202d38362e3831393731393136383529222c484c54484f55542c4152544852495449532c303130373030302c30313037333030313530302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030313630302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33392e392c33382e342c34312e332c2c2c333339302c222833332e353337323939333432332c202d38362e3830333635393034383229222c484c54484f55542c4152544852495449532c303130373030302c30313037333030313630302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030313930322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33352e312c33342e302c33362e312c2c2c313839342c222833332e353533323035303939372c202d38362e3734323938303136303329222c484c54484f55542c4152544852495449532c303130373030302c30313037333030313930322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030323030302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33362e332c33342e392c33372e372c2c2c333838352c222833332e353534313537343130362c202d38362e3731363732323939313529222c484c54484f55542c4152544852495449532c303130373030302c30313037333030323030302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030323130302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33332e312c33312e382c33342e332c2c2c333138362c222833332e353635303031353934322c202d38362e3731303130323437363629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030323130302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030323230302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33322e342c33312e322c33332e362c2c2c323633302c222833332e353532313330313230352c202d38362e3732373637353935303829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030323230302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030323330332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33332e342c33322e322c33342e362c2c2c323933362c222833332e353338333135333230372c202d38362e3732373034343534323829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030323330332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030323330352c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32322e332c32312e332c32332e332c2c2c323935322c222833332e353333333431353937362c202d38362e3734373935363630383429222c484c54484f55542c4152544852495449532c303130373030302c30313037333030323330352c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030323330362c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32362e392c32352e392c32372e392c2c2c333235372c222833332e353231333837333536342c202d38362e3734393030333132383929222c484c54484f55542c4152544852495449532c303130373030302c30313037333030323330362c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030323430302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33312e302c33302e302c33312e392c2c2c333632392c222833332e353236303734383330392c202d38362e3738333033313534383829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030323430302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030323730302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32372e322c32352e392c32382e352c2c2c333939322c222833332e353137363030383431392c202d38362e3831303638383734353229222c484c54484f55542c4152544852495449532c303130373030302c30313037333030323730302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030323930302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33382e372c33372e302c34302e332c2c2c323036342c222833332e353133323439383836342c202d38362e383330303437343929222c484c54484f55542c4152544852495449532c303130373030302c30313037333030323930302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333030312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32342e332c32332e362c32352e302c2c2c333737392c222833332e353132353135383039342c202d38362e3835373731363439343629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333030312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333030322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c34302e392c33392e312c34322e372c2c2c323230332c222833332e3531323235383130392c202d38362e3834343134333939303729222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333030322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333130302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33342e312c33322e372c33352e342c2c2c333633372c222833332e353035393635353735362c202d38362e3837343535303630383629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333130302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333230302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33392e302c33372e362c34302e332c2c2c3933312c222833332e353039343031383530322c202d38362e3838353930383139363129222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333230302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333330302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33382e362c33372e332c33392e392c2c2c3934372c222833332e353137313236313130382c202d38362e3839313338313937343929222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333330302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333430302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33362e332c33342e392c33372e372c2c2c323437372c222833332e353035323232393233342c202d38362e3930313438343436353629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333430302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333530302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33322e382c33312e352c33342e302c2c2c323738302c222833332e353036353731343031312c202d38362e3931393539313030363329222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333530302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333630302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33332e342c33322e312c33342e372c2c2c343638332c222833332e34383437363339372c202d38362e3839383133393239343729222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333630302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333730302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33312e322c33302e312c33322e312c2c2c353036332c222833332e343936393031383538392c202d38362e3839303737323934323629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333730302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333830322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33322e312c33312e302c33332e322c2c2c353430392c222833332e343738353730373739342c202d38362e38393030303039303729222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333830322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333830332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33362e302c33342e352c33372e332c2c2c343139392c222833332e3438353934353231342c202d38362e38363936393231383629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333830332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030333930302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33322e392c33312e382c33332e392c2c2c313738332c222833332e343938393935393332372c202d38362e3836343736303030333829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030333930302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030343030302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33372e382c33362e342c33392e332c2c2c333737322c222833332e343935333234363031352c202d38362e3835313632333230373329222c484c54484f55542c4152544852495449532c303130373030302c30313037333030343030302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030343230302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33362e362c33352e352c33372e372c2c2c323334312c222833332e353030373433393336312c202d38362e3832373037323033373929222c484c54484f55542c4152544852495449532c303130373030302c30313037333030343230302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030343530302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31352e312c31342e342c31352e382c2c2c353030332c222833332e353034313835373535362c202d38362e3830333337393833343629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030343530302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030343730312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32332e372c32322e352c32342e392c2c2c333438302c222833332e353037353234323134382c202d38362e3738333636373538333829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030343730312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030343730322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32372e302c32352e392c32382e312c2c2c323934342c222833332e353131393930323636312c202d38362e3736393435353039383929222c484c54484f55542c4152544852495449532c303130373030302c30313037333030343730322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030343830302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33302e312c32392e302c33312e322c2c2c313836312c222833332e343938393036343030382c202d38362e373832363939313429222c484c54484f55542c4152544852495449532c303130373030302c30313037333030343830302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030343930312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32322e312c32312e322c32332e302c2c2c313136372c222833332e343937313539353634352c202d38362e3739313734343036363829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030343930312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030343930322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31382e342c31372e352c31392e342c2c2c333134362c222833332e343933353832343034332c202d38362e3830303932393436303329222c484c54484f55542c4152544852495449532c303130373030302c30313037333030343930322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353030302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31392e362c31382e372c32302e342c2c2c333438322c222833332e343836363638393739352c202d38362e3831373332363238333129222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353030302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353130312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33362e362c33352e332c33372e382c2c2c313530372c222833332e343934353930393030382c202d38362e38333437363339333629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353130312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353130332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c34302e312c33382e362c34312e372c2c2c323538372c222833332e3438353731343838352c202d38362e3833323738313734363729222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353130332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353130342c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31382e392c31382e302c31392e382c2c2c323838312c222833332e343734393934313831362c202d38362e3833333534323137343729222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353130342c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353230302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33372e392c33362e312c33392e372c2c2c333734302c222833332e343830363730383737352c202d38362e3835303836373135313429222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353230302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353330322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33332e322c33312e382c33342e352c2c2c333436332c222833332e353736363435363434392c202d38362e3639363535393033313629222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353330322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353530302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33362e372c33352e342c33372e392c2c2c313832342c222833332e353637303239333839382c202d38362e3830303535363732313329222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353530302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353630302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33312e382c33302e332c33332e322c2c2c343336372c222833332e353230303938313233342c202d38362e3732373230363331393829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353630302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353730312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33352e382c33342e312c33372e332c2c2c323337322c222833332e343632393534333332392c202d38362e3838393834303636323829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353730312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353730322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33372e382c33362e302c33392e332c2c2c333431332c222833332e343639383639313338352c202d38362e38373432393036303229222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353730322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353830302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31382e342c31372e372c31392e312c2c2c343231362c222833332e3437393036323332352c202d38362e3831323830363330383929222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353830302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353930332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33322e322c33312e312c33332e332c2c2c343933332c222833332e3539373136323330392c202d38362e3637363637333633353129222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353930332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353930352c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33342e322c33322e392c33352e342c2c2c353033392c222833332e3630333938383435362c202d38362e3730303831323334313829222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353930352c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353930372c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33302e312c32382e392c33312e332c2c2c313937352c222833332e363134323836313530312c202d38362e3636393139393634313729222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353930372c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353930382c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33322e352c33312e352c33332e352c2c2c313632312c222833332e363138323932343433322c202d38362e3638303134383330383429222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353930382c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353930392c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32382e302c32362e352c32392e342c2c2c323532342c222833332e3631313831313737332c202d38362e3732313432323135313429222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353930392c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333030353931302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32382e302c32362e382c32392e312c2c2c343631322c222833332e363239393031373439392c202d38362e3731393433313132323929222c484c54484f55542c4152544852495449532c303130373030302c30313037333030353931302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031303530302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c34312e332c33392e342c34332e312c2c2c3131342c222833332e343336333738363830362c202d38362e3931323839323330373229222c484c54484f55542c4152544852495449532c303130373030302c30313037333031303530302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031303730312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31352e312c31342e302c31362e332c2c2c37342c222833332e3437333838363135352c202d38362e3831343634383737363229222c484c54484f55542c4152544852495449532c303130373030302c30313037333031303730312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031303730362c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31352e332c31342e362c31362e302c2c2c313532382c222833332e343434333730393434322c202d38362e3834303533353236343529222c484c54484f55542c4152544852495449532c303130373030302c30313037333031303730362c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031303830312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32342e392c32332e392c32362e322c2c2c3136382c222833332e3531343039373835332c202d38362e3734363639373133363229222c484c54484f55542c4152544852495449532c303130373030302c30313037333031303830312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031303830322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33332e322c33312e372c33352e302c2c2c3137322c222833332e343838353439333437372c202d38362e37383038343330323429222c484c54484f55542c4152544852495449532c303130373030302c30313037333031303830322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031303830332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32332e382c32322e362c32352e302c2c2c3531342c222833332e353232393039333839322c202d38362e3731303236313836343229222c484c54484f55542c4152544852495449532c303130373030302c30313037333031303830332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031303830352c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33322e302c33302e312c33342e302c2c2c38362c222833332e343935323739323437322c202d38362e3639383731383439373429222c484c54484f55542c4152544852495449532c303130373030302c30313037333031303830352c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313130342c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32392e382c32382e352c33312e312c2c2c313638382c222833332e363135393433363433332c202d38362e3635353738393235303729222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313130342c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313130372c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c34322c222833332e353830343834353234392c202d38362e3633303131313039363129222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313130372c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313130382c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c392c222833332e363035303734323539362c202d38362e3633313637323933383629222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313130382c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313230372c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32352e332c32342e332c32362e322c2c2c3831352c222833332e363731383834383730362c202d38362e3637373235313034363529222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313230372c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313230392c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32372e322c32362e312c32382e332c2c2c313036322c222833332e363535373138393939322c202d38362e3730353036393833343929222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313230392c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313231302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32332e382c32322e332c32352e332c2c2c313338352c222833332e363634313839333735352c202d38362e3639353631373036383629222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313231302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313830332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32382e322c32372e302c32392e342c2c2c3932382c222833332e363235323537353137332c202d38362e3639393836313034303929222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313830332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313830342c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32352e382c32342e362c32372e302c2c2c313135372c222833332e363437343931363137352c202d38362e3730343239373434323429222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313830342c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313930312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c362c222833332e363335353431343634362c202d38362e37333639343636393129222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313930312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031313930342c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33342e382c33332e332c33362e332c2c2c313931352c222833332e3539333134303138352c202d38362e3733353739333035343129222c484c54484f55542c4152544852495449532c303130373030302c30313037333031313930342c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323030312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c34342e382c34322e392c34362e362c2c2c3330342c222833332e353931393438363131322c202d38362e38363433383430363829222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323030312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323030322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c34342c222833332e353835393233343139372c202d38362e3833353730303731383829222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323030322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323230302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c32332c222833332e353936373434313930342c202d38372e3038373933393638353729222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323230302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323330322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32362e332c32352e322c32372e352c2c2c3134342c222833332e353534323831363335322c202d38372e3035343436393134313629222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323330322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323330352c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33302e352c32382e362c33322e322c2c2c3430332c222833332e343639353335383036342c202d38362e393638333137333929222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323330352c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323430312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32372e342c32352e362c32392e342c2c2c313036362c222833332e353537313935313034382c202d38362e3837373739333530343929222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323430312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323430322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33312e392c33302e382c33332e302c2c2c3431382c222833332e3534393938343137322c202d38362e3839393435343333323729222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323430322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323530302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33372e392c33362e352c33392e342c2c2c3431302c222833332e3532393136303438362c202d38362e3933343634373634343529222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323530302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323630322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33352e342c33342e352c33362e332c2c2c3337312c222833332e3537303136343637342c202d38362e36363634333035383229222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323630322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323730312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c34342c222833332e353438343037383037312c202d38362e3633323337373334353529222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323730312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323730332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e382c31332e382c31352e392c2c2c3439382c222833332e343638313138303934332c202d38362e3636373138383832313329222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323730332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323730342c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32372e302c32352e322c32392e302c2c2c3131332c222833332e353033343139353930382c202d38362e3631383039383334303329222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323730342c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323830332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31312e392c31312e312c31322e372c2c2c313236312c222833332e343433393432353836352c202d38362e3732313239333639333829222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323830332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031323931302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c392c222833332e343334353830353034322c202d38362e3732363332393230353929222c484c54484f55542c4152544852495449532c303130373030302c30313037333031323931302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031333030322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c34312e302c33392e312c34322e392c2c2c313531342c222833332e34363630343138312c202d38362e3835363732383737393729222c484c54484f55542c4152544852495449532c303130373030302c30313037333031333030322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031333130302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33382e372c33372e352c33392e392c2c2c343432342c222833332e34343838303231342c202d38362e3838373834303135373929222c484c54484f55542c4152544852495449532c303130373030302c30313037333031333130302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031333330302c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33392e322c33372e362c34302e372c2c2c313738322c222833332e343339363434333536392c202d38362e3932343837363836363529222c484c54484f55542c4152544852495449532c303130373030302c30313037333031333330302c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031333930312c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33352e302c33332e342c33362e362c2c2c3935322c222833332e343732393338343930362c202d38362e3935343733333736343829222c484c54484f55542c4152544852495449532c303130373030302c30313037333031333930312c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031343330322c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31392e302c31382e312c31392e382c2c2c323737382c222833332e343234343635383832392c202d38362e3838343134373432313729222c484c54484f55542c4152544852495449532c303130373030302c30313037333031343330322c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313037333031343431332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31372e372c31362e392c31382e362c2c2c3339372c222833332e343232363539333131372c202d38362e3835303836323037353129222c484c54484f55542c4152544852495449532c303130373030302c30313037333031343431332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313131373033303231332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32302e332c31392e312c32312e352c2c2c3634342c222833332e343339353937353139332c202d38362e3637333539353933353929222c484c54484f55542c4152544852495449532c303130373030302c30313131373033303231332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313131373033303231372c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c2c2c2c2a2c457374696d61746573207375707072657373656420666f7220706f70756c6174696f6e206c657373207468616e2035302c31362c222833332e343535363939353736332c202d38362e3635323032303836333929222c484c54484f55542c4152544852495449532c303130373030302c30313131373033303231372c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c4865616c7468204f7574636f6d65732c303130373030302d30313131373033303330332c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31322e382c31322e312c31332e362c2c2c3936382c222833332e343235383636313233392c202d38362e37313338313933353629222c484c54484f55542c4152544852495449532c303130373030302c30313131373033303330332c4172746872697469730d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c436974792c42524653532c556e6865616c746879204265686176696f72732c303130373030302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31312e322c31312e312c31312e332c2c2c3231323233372c222833332e353237353636333737332c202d38362e3739383831373436373829222c554e484245482c42494e47452c303130373030302c2c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c436974792c42524653532c556e6865616c746879204265686176696f72732c303130373030302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31312e332c31312e332c31312e342c2c2c3231323233372c222833332e353237353636333737332c202d38362e3739383831373436373829222c554e484245482c42494e47452c303130373030302c2c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030303130302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e312c392e372c31302e352c2c2c333034322c222833332e353739343332383332362c202d38362e3732323833323339323629222c554e484245482c42494e47452c303130373030302c30313037333030303130302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030303330302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e382c31302e332c31312e322c2c2c323733352c222833332e353432383230383638362c202d38362e37353234333339373829222c554e484245482c42494e47452c303130373030302c30313037333030303330302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030303430302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e352c392e302c31302e302c2c2c333333382c222833332e353633323434393633332c202d38362e3736343034373430363429222c554e484245482c42494e47452c303130373030302c30313037333030303430302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030303530302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e362c382e312c392e302c2c2c323836342c222833332e353434323430343539342c202d38362e3737343931333037313929222c554e484245482c42494e47452c303130373030302c30313037333030303530302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030303730302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c372e342c362e392c372e392c2c2c323537372c222833332e353532353430363133392c202d38362e3830313638393337303629222c554e484245482c42494e47452c303130373030302c30313037333030303730302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030303830302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e392c382e352c392e332c2c2c333835392c222833332e3534393639373738392c202d38362e3833333039343437343429222c554e484245482c42494e47452c303130373030302c30313037333030303830302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030313130302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e362c392e322c31302e312c2c2c353335342c222833332e353432393134333332352c202d38362e3837353637383238353229222c554e484245482c42494e47452c303130373030302c30313037333030313130302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030313230302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e342c392e312c392e372c2c2c323837362c222833332e353237383736373730362c202d38362e3836303431363136383629222c554e484245482c42494e47452c303130373030302c30313037333030313230302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030313430302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e392c382e352c392e332c2c2c323138312c222833332e353236313439373235382c202d38362e38333531343636303629222c554e484245482c42494e47452c303130373030302c30313037333030313430302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030313530302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e342c392e312c392e382c2c2c333138392c222833332e353239383732373334322c202d38362e3831393731393136383529222c554e484245482c42494e47452c303130373030302c30313037333030313530302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030313630302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e322c372e372c382e362c2c2c333339302c222833332e353337323939333432332c202d38362e3830333635393034383229222c554e484245482c42494e47452c303130373030302c30313037333030313630302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030313930322c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e332c382e392c392e382c2c2c313839342c222833332e353533323035303939372c202d38362e3734323938303136303329222c554e484245482c42494e47452c303130373030302c30313037333030313930322c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323030302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e312c382e372c392e362c2c2c333838352c222833332e353534313537343130362c202d38362e3731363732323939313529222c554e484245482c42494e47452c303130373030302c30313037333030323030302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323130302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e302c392e362c31302e362c2c2c333138362c222833332e353635303031353934322c202d38362e3731303130323437363629222c554e484245482c42494e47452c303130373030302c30313037333030323130302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323230302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e312c392e372c31302e362c2c2c323633302c222833332e353532313330313230352c202d38362e3732373637353935303829222c554e484245482c42494e47452c303130373030302c30313037333030323230302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323330332c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e392c382e352c392e342c2c2c323933362c222833332e353338333135333230372c202d38362e3732373034343534323829222c554e484245482c42494e47452c303130373030302c30313037333030323330332c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323330352c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31362e332c31352e382c31362e372c2c2c323935322c222833332e353333333431353937362c202d38362e3734373935363630383429222c554e484245482c42494e47452c303130373030302c30313037333030323330352c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323330362c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31352e362c31352e322c31362e312c2c2c333235372c222833332e353231333837333536342c202d38362e3734393030333132383929222c554e484245482c42494e47452c303130373030302c30313037333030323330362c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323430302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31312e312c31302e382c31312e352c2c2c333632392c222833332e353236303734383330392c202d38362e3738333033313534383829222c554e484245482c42494e47452c303130373030302c30313037333030323430302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323730302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31332e362c31332e302c31342e322c2c2c333939322c222833332e353137363030383431392c202d38362e3831303638383734353229222c554e484245482c42494e47452c303130373030302c30313037333030323730302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323930302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e332c372e382c382e382c2c2c323036342c222833332e353133323439383836342c202d38362e383330303437343929222c554e484245482c42494e47452c303130373030302c30313037333030323930302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333030312c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e332c31332e332c31352e332c2c2c333737392c222833332e353132353135383039342c202d38362e3835373731363439343629222c554e484245482c42494e47452c303130373030302c30313037333030333030312c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333030322c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c372e342c372e302c372e382c2c2c323230332c222833332e3531323235383130392c202d38362e3834343134333939303729222c554e484245482c42494e47452c303130373030302c30313037333030333030322c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333130302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e352c392e312c31302e302c2c2c333633372c222833332e353035393635353735362c202d38362e3837343535303630383629222c554e484245482c42494e47452c303130373030302c30313037333030333130302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333230302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c372e362c372e332c382e302c2c2c3933312c222833332e353039343031383530322c202d38362e3838353930383139363129222c554e484245482c42494e47452c303130373030302c30313037333030333230302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333330302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e342c382e302c382e392c2c2c3934372c222833332e353137313236313130382c202d38362e3839313338313937343929222c554e484245482c42494e47452c303130373030302c30313037333030333330302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030323330352c4f62657369747920616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33312e372c33302e342c33332e302c2c2c323935322c222833332e353333333431353937362c202d38362e3734373935363630383429222c554e484245482c4f4245534954592c303130373030302c30313037333030323330352c4f6265736974790d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333430302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e382c382e332c392e332c2c2c323437372c222833332e353035323232393233342c202d38362e3930313438343436353629222c554e484245482c42494e47452c303130373030302c30313037333030333430302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333530302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e302c392e352c31302e352c2c2c323738302c222833332e353036353731343031312c202d38362e3931393539313030363329222c554e484245482c42494e47452c303130373030302c30313037333030333530302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333630302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e392c392e342c31302e332c2c2c343638332c222833332e34383437363339372c202d38362e3839383133393239343729222c554e484245482c42494e47452c303130373030302c30313037333030333630302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333730302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e332c392e392c31302e372c2c2c353036332c222833332e343936393031383538392c202d38362e3839303737323934323629222c554e484245482c42494e47452c303130373030302c30313037333030333730302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333830322c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31302e332c392e392c31302e372c2c2c353430392c222833332e343738353730373739342c202d38362e38393030303039303729222c554e484245482c42494e47452c303130373030302c30313037333030333830322c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333830332c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e322c382e372c392e372c2c2c343139392c222833332e3438353934353231342c202d38362e38363936393231383629222c554e484245482c42494e47452c303130373030302c30313037333030333830332c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030333930302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e352c392e312c392e382c2c2c313738332c222833332e343938393935393332372c202d38362e3836343736303030333829222c554e484245482c42494e47452c303130373030302c30313037333030333930302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030343030302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c382e352c382e302c392e302c2c2c333737322c222833332e343935333234363031352c202d38362e3835313632333230373329222c554e484245482c42494e47452c303130373030302c30313037333030343030302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030343230302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c392e312c382e382c392e342c2c2c323334312c222833332e353030373433393336312c202d38362e3832373037323033373929222c554e484245482c42494e47452c303130373030302c30313037333030343230302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030343530302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31332e382c31322e392c31342e382c2c2c353030332c222833332e353034313835373535362c202d38362e3830333337393833343629222c554e484245482c42494e47452c303130373030302c30313037333030343530302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030343730312c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31362e342c31352e372c31372e322c2c2c333438302c222833332e353037353234323134382c202d38362e3738333636373538333829222c554e484245482c42494e47452c303130373030302c30313037333030343730312c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030343730322c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31352e332c31342e392c31352e372c2c2c323934342c222833332e353131393930323636312c202d38362e3736393435353039383929222c554e484245482c42494e47452c303130373030302c30313037333030343730322c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030343830302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e312c31332e372c31342e352c2c2c313836312c222833332e343938393036343030382c202d38362e373832363939313429222c554e484245482c42494e47452c303130373030302c30313037333030343830302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030343930312c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31372e302c31362e342c31372e372c2c2c313136372c222833332e343937313539353634352c202d38362e3739313734343036363829222c554e484245482c42494e47452c303130373030302c30313037333030343930312c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030343930322c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31362e332c31352e362c31372e312c2c2c333134362c222833332e343933353832343034332c202d38362e3830303932393436303329222c554e484245482c42494e47452c303130373030302c30313037333030343930322c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030353030302c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31362e312c31352e362c31362e372c2c2c333438322c222833332e343836363638393739352c202d38362e3831373332363238333129222c554e484245482c42494e47452c303130373030302c30313037333030353030302c42696e6765204472696e6b696e670d0a323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030353130312c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c372e392c372e352c382e332c2c2c313530372c222833332e343934353930393030382c202d38362e38333437363339333629222c554e484245482c42494e47452c303130373030302c30313037333030353130312c42696e6765204472696e6b696e670d0a000000000180000100000109000000000000000000012efc323031352c414c2c416c6162616d612c4269726d696e6768616d2c43656e7375732054726163742c42524653532c556e6865616c746879204265686176696f72732c303130373030302d30313037333030353130332c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c372e372c372e342c382e302c2c2c323538372c222833332e3438353731343838352c202d38362e3833323738313734363729222c554e484245482c42494e47452c303130373030302c30313037333030353130332c42696e6765204472696e6b696e670d0a000000000180000500000014000000000000000000012efc0000000000012efc000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0ICogZnJvbSBvc3NvYmplY3Qn</Expression><InputSerialization><CSV><FileHeaderInfo>None</FileHeaderInfo><RecordDelimiter>Cg==</RecordDelimiter><FieldDelimiter>LA==</FieldDelimiter><QuoteCharacter>Ig==</QuoteCharacter><CommentCharacter>Iw==</CommentCharacter><Range></Range></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select * from ossobject'"),
				InputSerializationSelect: InputSerializationSelect{
					CsvBodyInput: &CSVSelectInput{
						FileHeaderInfo:   Ptr("None"),
						CommentCharacter: Ptr("#"),
						RecordDelimiter:  Ptr("\n"),
						FieldDelimiter:   Ptr(","),
						QuoteCharacter:   Ptr("\""),
						Range:            Ptr(""),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, err := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "Year,StateAbbr,StateDesc,")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("018000010000003e0000000000000000000002e1323031352c55532c2c4865616c746820496e737572616e63650d0a323031352c55532c2c4865616c746820496e737572616e63650d0a0000000001800005000000140000000000000000000002e100000000000002e1000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IFllYXIsU3RhdGVBYmJyLCBDaXR5TmFtZSwgU2hvcnRfUXVlc3Rpb25fVGV4dCBmcm9tIG9zc29iamVjdA==</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo><Range>line-range=0-2</Range></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select Year,StateAbbr, CityName, Short_Question_Text from ossobject"),
				InputSerializationSelect: InputSerializationSelect{
					CsvBodyInput: &CSVSelectInput{
						FileHeaderInfo: Ptr("Use"),
						Range:          Ptr("0-2"),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, err := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "2015,US,,Health Insurance")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("0180000100000017000000000000000000012efc323031352c323031352c323031350a000000000180000500000014000000000000000000012efc0000000000012efc000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IGF2ZyhjYXN0KHllYXIgYXMgaW50KSksIG1heChjYXN0KHllYXIgYXMgaW50KSksIG1pbihjYXN0KHllYXIgYXMgaW50KSkgZnJvbSBvc3NvYmplY3Qgd2hlcmUgeWVhciA9IDIwMTU=</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select avg(cast(year as int)), max(cast(year as int)), min(cast(year as int)) from ossobject where year = 2015"),
				InputSerializationSelect: InputSerializationSelect{
					CsvBodyInput: &CSVSelectInput{
						FileHeaderInfo: Ptr("Use"),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, err := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:14]), "2015,2015,2015")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("018000010000005e000000000000000000012efc32332e363737373030333438343332303635323536303432313738383131303936353633353731372c38312e382c363739352e353030303030303030303032373238343834313035333138373834373133373435310a000000000180000500000014000000000000000000012efc0000000000012efc000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IGF2ZyhjYXN0KGRhdGFfdmFsdWUgYXMgZG91YmxlKSksIG1heChjYXN0KGRhdGFfdmFsdWUgYXMgZG91YmxlKSksIHN1bShjYXN0KGRhdGFfdmFsdWUgYXMgZG91YmxlKSkgZnJvbSBvc3NvYmplY3Q=</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select avg(cast(data_value as double)), max(cast(data_value as double)), sum(cast(data_value as double)) from ossobject"),
				InputSerializationSelect: InputSerializationSelect{
					CsvBodyInput: &CSVSelectInput{
						FileHeaderInfo: Ptr("Use"),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, err := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:14]), "23.67770034843")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("01800001000000110000000000000000000000086162637c6465660d0a0000000001800005000000140000000000000000000000080000000000000008000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xLCBfMiBmcm9tIG9zc29iamVjdA==</Expression><InputSerialization></InputSerialization><OutputSerialization><CSV><RecordDelimiter>&#xD;&#xA;</RecordDelimiter><FieldDelimiter>|</FieldDelimiter></CSV></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select _1, _2 from ossobject"),
				OutputSerializationSelect: OutputSerializationSelect{
					CsvBodyOutput: &CSVSelectOutput{
						RecordDelimiter: Ptr("\r\n"),
						FieldDelimiter:  Ptr("|"),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, err := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte), "abc|def\r\n")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("0180000100000987000000000000000000012dfb596561722c5374617465416262720a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031342c55530a323031342c55530a323031352c55530a323031352c55530a323031352c55530a323031342c55530a323031342c55530a323031342c55530a323031342c55530a323031352c55530a323031352c55530a323031342c55530a323031342c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031342c55530a323031342c55530a323031352c55530a323031352c55530a323031352c55530a323031352c55530a323031342c55530a323031342c55530a323031352c55530a323031342c55530a323031342c55530a323031352c55530a323031352c55530a323031342c55530a323031342c55530a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a323031352c414c0a000000000180000100000010000000000000000000012efc323031352c414c0a000000000180000500000014000000000000000000012efc0000000000012efc000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xLCBfMiBmcm9tIG9zc29iamVjdA==</Expression><InputSerialization></InputSerialization><OutputSerialization></OutputSerialization><Options><SkipPartialDataRecord>true</SkipPartialDataRecord></Options></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select _1, _2 from ossobject"),
				SelectOptions: &SelectOptions{
					SkipPartialDataRecord: Ptr(true),
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "Year,StateAbbr\n2015,US\n20")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":        "534B371674E88A4D8906****",
			"Date":                    "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id":        "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
			"x-oss-select-output-raw": "true",
		},
		[]byte(hexStrToByte("596561720a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031340a323031340a323031350a323031350a323031350a323031340a323031340a323031340a323031340a323031350a323031350a323031340a323031340a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031340a323031340a323031350a323031350a323031350a323031350a323031340a323031340a323031350a323031340a323031340a323031350a323031350a323031340a323031340a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xIGZyb20gb3Nzb2JqZWN0</Expression><InputSerialization></InputSerialization><OutputSerialization><OutputRawData>true</OutputRawData></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select _1 from ossobject"),
				OutputSerializationSelect: OutputSerializationSelect{
					OutputRawData: Ptr(true),
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "Year\n2015\n2015\n2015\n2015\n")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":        "534B371674E88A4D8906****",
			"Date":                    "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id":        "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
			"x-oss-select-output-raw": "true",
		},
		[]byte(hexStrToByte("596561722c5374617465416262722c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031342c55532c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a323031352c414c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c0a")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xLF8yIGZyb20gb3Nzb2JqZWN0</Expression><InputSerialization></InputSerialization><OutputSerialization><OutputRawData>true</OutputRawData><KeepAllColumns>true</KeepAllColumns></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select _1,_2 from ossobject"),
				OutputSerializationSelect: OutputSerializationSelect{
					OutputRawData:  Ptr(true),
					KeepAllColumns: Ptr(true),
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "Year,StateAbbr,,,,,,,,,,,")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("01800001000005f3000000000000000000012dfb796561720a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031340a323031340a323031350a323031350a323031350a323031340a323031340a323031340a323031340a323031350a323031350a323031340a323031340a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031340a323031340a323031350a323031350a323031350a323031350a323031340a323031340a323031350a323031340a323031340a323031350a323031350a323031340a323031340a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a323031350a00000000018000010000000d000000000000000000012efc323031350a000000000180000500000014000000000000000000012efc0000000000012efc000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xIGZyb20gb3Nzb2JqZWN0</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization><OutputHeader>true</OutputHeader></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select _1 from ossobject"),
				OutputSerializationSelect: OutputSerializationSelect{
					OutputHeader: Ptr(true),
				},
				InputSerializationSelect: InputSerializationSelect{
					CsvBodyInput: &CSVSelectInput{
						FileHeaderInfo: Ptr("Use"),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "year\n2015\n2015\n2015\n2015\n")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("01800001000000cd0000000000000000000009377374617465616262722c63617465676f72790a55532c50726576656e74696f6e0a55532c50726576656e74696f6e0a55532c4865616c7468204f7574636f6d65730a55532c4865616c7468204f7574636f6d65730a55532c556e6865616c746879204265686176696f72730a55532c556e6865616c746879204265686176696f72730a55532c4865616c7468204f7574636f6d65730a55532c4865616c7468204f7574636f6d65730a55532c50726576656e74696f6e0a55532c50726576656e74696f6e0a1f53a7e401800005000000140000000000000000000009370000000000000937000000c87a38f066")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IFN0YXRlQWJicixDYXRlZ29yeSBmcm9tIG9zc29iamVjdCBsaW1pdCAxMA==</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization><EnablePayloadCrc>true</EnablePayloadCrc><OutputHeader>true</OutputHeader></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select StateAbbr,Category from ossobject limit 10"),
				OutputSerializationSelect: OutputSerializationSelect{
					OutputHeader:     Ptr(true),
					EnablePayloadCrc: Ptr(true),
				},
				InputSerializationSelect: InputSerializationSelect{
					CsvBodyInput: &CSVSelectInput{
						FileHeaderInfo: Ptr("Use"),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "stateabbr,category\nUS,Pre")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("01800001000008040000000000000000000007f2323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31352e342c31352e312c31352e372c2c2c3330383734353533382c2c50524556454e542c414343455353322c2c2c4865616c746820496e737572616e63650d0000323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c43757272656e74206c61636b206f66206865616c746820696e737572616e636520616d6f6e67206164756c74732061676564203138e2809336342059656172732c252c4372645072762c43727564652070726576616c656e63652c31342e382c31342e352c31352e302c2c2c3330383734353533382c2c50524556454e542c414343455353322c2c2c4865616c746820496e737572616e63650d0000323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c32322e352c32322e332c32322e372c2c2c3330383734353533382c2c484c54484f55542c4152544852495449532c2c2c4172746872697469730d0000323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c41727468726974697320616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c32342e372c32342e352c32342e392c2c2c3330383734353533382c2c484c54484f55542c4152544852495449532c2c2c4172746872697469730d0000323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c31372e322c31362e392c31372e342c2c2c3330383734353533382c2c554e484245482c42494e47452c2c2c42696e6765204472696e6b696e670d0000323031352c55532c556e69746564205374617465732c2c55532c42524653532c556e6865616c746879204265686176696f72732c35392c42696e6765206472696e6b696e6720616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c31362e332c31362e312c31362e352c2c2c3330383734353533382c2c554e484245482c42494e47452c2c2c42696e6765204472696e6b696e670d0000323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4869676820626c6f6f6420707265737375726520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c32392e342c32392e322c32392e372c2c2c3330383734353533382c2c484c54484f55542c4250484947482c2c2c4869676820426c6f6f642050726573737572650d0000323031352c55532c556e69746564205374617465732c2c55532c42524653532c4865616c7468204f7574636f6d65732c35392c4869676820626c6f6f6420707265737375726520616d6f6e67206164756c74732061676564203e3d31382059656172732c252c4372645072762c43727564652070726576616c656e63652c33312e392c33312e362c33322e322c2c2c3330383734353533382c2c484c54484f55542c4250484947482c2c2c4869676820426c6f6f642050726573737572650d0000323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c54616b696e67206d65646963696e6520666f72206869676820626c6f6f6420707265737375726520636f6e74726f6c20616d6f6e67206164756c74732061676564203e3d31382059656172732077697468206869676820626c6f6f642070726573737572652c252c41676541646a5072762c4167652d61646a75737465642070726576616c656e63652c35372e372c35372e312c35382e342c2c2c3330383734353533382c2c50524556454e542c42504d45442c2c2c54616b696e67204250204d656469636174696f6e0d0000323031352c55532c556e69746564205374617465732c2c55532c42524653532c50726576656e74696f6e2c35392c54616b696e67206d65646963696e6520666f72206869676820626c6f6f6420707265737375726520636f6e74726f6c20616d6f6e67206164756c74732061676564203e3d31382059656172732077697468206869676820626c6f6f642070726573737572652c252c4372645072762c43727564652070726576616c656e63652c37372e322c37362e382c37372e372c2c2c3330383734353533382c2c50524556454e542c42504d45442c2c2c54616b696e67204250204d656469636174696f6e0d0000a8bfed2b01800005000000140000000000000000000007f200000000000007f2000000c8aa94a492")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0ICogZnJvbSBvc3NvYmplY3QgbGltaXQgMTA=</Expression><InputSerialization><CSV><FileHeaderInfo>IGNORE</FileHeaderInfo><RecordDelimiter>Cg==</RecordDelimiter><FieldDelimiter>LA==</FieldDelimiter><QuoteCharacter>Ig==</QuoteCharacter><CommentCharacter>Iw==</CommentCharacter><Range>split-range=0-12</Range><AllowQuotedRecordDelimiter>true</AllowQuotedRecordDelimiter></CSV><CompressionType>NONE</CompressionType></InputSerialization><OutputSerialization><CSV><RecordDelimiter>&#xA;</RecordDelimiter><FieldDelimiter></FieldDelimiter></CSV><OutputRawData>false</OutputRawData><KeepAllColumns>false</KeepAllColumns><EnablePayloadCrc>true</EnablePayloadCrc><OutputHeader>false</OutputHeader></OutputSerialization><Options><SkipPartialDataRecord>false</SkipPartialDataRecord><MaxSkippedRecordsAllowed>2</MaxSkippedRecordsAllowed></Options></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select * from ossobject limit 10"),
				InputSerializationSelect: InputSerializationSelect{
					CompressionType: Ptr("NONE"),
					CsvBodyInput: &CSVSelectInput{
						FileHeaderInfo:             Ptr("IGNORE"),
						RecordDelimiter:            Ptr("\n"),
						FieldDelimiter:             Ptr(","),
						QuoteCharacter:             Ptr("\""),
						CommentCharacter:           Ptr("#"),
						SplitRange:                 Ptr("0-12"),
						AllowQuotedRecordDelimiter: Ptr(true),
					},
				},
				OutputSerializationSelect: OutputSerializationSelect{
					CsvBodyOutput: &CSVSelectOutput{
						RecordDelimiter: Ptr("\n"),
						FieldDelimiter:  Ptr(""),
					},
					KeepAllColumns:   Ptr(false),
					OutputRawData:    Ptr(false),
					EnablePayloadCrc: Ptr(true),
					OutputHeader:     Ptr(false),
				},
				SelectOptions: &SelectOptions{
					SkipPartialDataRecord:    Ptr(false),
					MaxSkippedRecordsAllowed: Ptr(2),
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "2015,US,United States,,US")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("018000010000096f000000000000000000031a3e7b22636175637573223a6e756c6c2c22636f6e67726573735f6e756d62657273223a5b3131332c3131342c3131355d2c2263757272656e74223a747275652c226465736372697074696f6e223a224a756e696f722053656e61746f7220666f7220576973636f6e73696e222c226469737472696374223a6e756c6c2c22656e6464617465223a22323031392d30312d3033222c226578747261223a7b2261646472657373223a2237303920486172742053656e617465204f6666696365204275696c64696e672057617368696e67746f6e204443203230353130222c22636f6e746163745f666f726d223a2268747470733a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f765c2f666565646261636b222c22666178223a223230322d3232352d36393432222c226f6666696365223a2237303920486172742053656e617465204f6666696365204275696c64696e67222c227273735f75726c223a22687474703a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f765c2f7273735c2f66656564735c2f3f747970653d616c6c227d2c226c6561646572736869705f7469746c65223a6e756c6c2c227061727479223a2244656d6f63726174222c22706572736f6e223a7b2262696f67756964656964223a2242303031323330222c226269727468646179223a22313936322d30322d3131222c22637370616e6964223a35373838342c2266697273746e616d65223a2254616d6d79222c2267656e646572223a2266656d616c65222c2267656e6465725f6c6162656c223a2246656d616c65222c226c6173746e616d65223a2242616c6477696e222c226c696e6b223a2268747470733a5c2f5c2f7777772e676f76747261636b2e75735c2f636f6e67726573735c2f6d656d626572735c2f74616d6d795f62616c6477696e5c2f343030303133222c226d6964646c656e616d65223a22222c226e616d65223a2253656e2e2054616d6d792042616c6477696e205b442d57495d222c226e616d656d6f64223a22222c226e69636b6e616d65223a22222c226f736964223a224e3030303034333637222c227076736964223a2233343730222c22736f72746e616d65223a2242616c6477696e2c2054616d6d79202853656e2e29205b442d57495d222c22747769747465726964223a2253656e61746f7242616c6477696e222c22796f75747562656964223a22776974616d6d7962616c6477696e227d2c2270686f6e65223a223230322d3232342d35363533222c22726f6c655f74797065223a2273656e61746f72222c22726f6c655f747970655f6c6162656c223a2253656e61746f72222c2273656e61746f725f636c617373223a22636c61737331222c2273656e61746f725f636c6173735f6c6162656c223a22436c6173732031222c2273656e61746f725f72616e6b223a226a756e696f72222c2273656e61746f725f72616e6b5f6c6162656c223a224a756e696f72222c22737461727464617465223a22323031332d30312d3033222c227374617465223a225749222c227469746c65223a2253656e2e222c227469746c655f6c6f6e67223a2253656e61746f72222c2277656273697465223a2268747470733a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f76227d2c7b22636175637573223a6e756c6c2c22636f6e67726573735f6e756d62657273223a5b3131332c3131342c3131355d2c2263757272656e74223a747275652c226465736372697074696f6e223a2253656e696f722053656e61746f7220666f72204f68696f222c226469737472696374223a6e756c6c2c22656e6464617465223a22323031392d30312d3033222c226578747261223a7b2261646472657373223a2237313320486172742053656e617465204f6666696365204275696c64696e672057617368696e67746f6e204443203230353130222c22636f6e746163745f666f726d223a22687474703a5c2f5c2f7777772e62726f776e2e73656e6174652e676f765c2f636f6e746163745c2f222c22666178223a223230322d3232382d36333231222c226f6666696365223a2237313320486172742053656e617465204f6666696365204275696c64696e67222c227273735f75726c223a22687474703a5c2f5c2f7777772e62726f776e2e73656e6174652e676f765c2f7273735c2f66656564735c2f3f747970653d616c6c26616d703b227d2c226c6561646572736869705f7469746c65223a6e756c6c2c227061727479223a2244656d6f63726174222c22706572736f6e223a7b2262696f67756964656964223a2242303030393434222c226269727468646179223a22313935322d31312d3039222c22637370616e6964223a353035312c2266697273746e616d65223a2253686572726f64222c2267656e646572223a226d616c65222c2267656e6465725f6c6162656c223a224d616c65222c226c6173746e616d65223a2242726f776e222c226c696e6b223a2268747470733a5c2f5c2f7777772e676f76747261636b2e75735c2f636f6e67726573735c2f6d656d626572735c2f73686572726f645f62726f776e5c2f343030303530222c226d6964646c656e616d65223a22222c226e616d65223a2253656e2e2053686572726f642042726f776e205b442d4f485d222c226e616d656d6f64223a22222c226e69636b6e616d65223a22222c226f736964223a224e3030303033353335222c227076736964223a223237303138222c22736f72746e616d65223a2242726f776e2c2053686572726f64202853656e2e29205b442d4f485d222c22747769747465726964223a2253656e53686572726f6442726f776e222c22796f75747562656964223a2253686572726f6442726f776e4f68696f227d2c2270686f6e65223a223230322d3232342d32333135222c22726f6c655f74797065223a2273656e61746f72222c22726f6c655f747970655f6c6162656c223a2253656e61746f72222c2273656e61746f725f636c617373223a22636c61737331222c2273656e61746f725f636c6173735f6c6162656c223a22436c6173732031222c2273656e61746f725f72616e6b223a2273656e696f72222c2273656e61746f725f72616e6b5f6c6162656c223a2253656e696f72222c22737461727464617465223a22323031332d30312d3033222c227374617465223a224f48222c227469746c65223a2253656e2e222c227469746c655f6c6f6e67223a2253656e61746f72222c2277656273697465223a2268747470733a5c2f5c2f7777772e62726f776e2e73656e6174652e676f76227d2c000000000180000500000014000000000000000000031a3e0000000000031a3e000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=json%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0ICogZnJvbSBvc3NvYmplY3Qub2JqZWN0c1sqXSB3aGVyZSBwYXJ0eSA9ICdEZW1vY3JhdCcgbGltaXQgMTA=</Expression><InputSerialization><JSON><Type>DOCUMENT</Type></JSON></InputSerialization><OutputSerialization><JSON><RecordDelimiter>LA==</RecordDelimiter></JSON></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select * from ossobject.objects[*] where party = 'Democrat' limit 10"),
				InputSerializationSelect: InputSerializationSelect{
					JsonBodyInput: &JSONSelectInput{
						JSONType: Ptr("DOCUMENT"),
					},
				},
				OutputSerializationSelect: OutputSerializationSelect{
					JsonBodyOutput: &JSONSelectOutput{
						RecordDelimiter: Ptr(","),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "{\"caucus\":null,\"congress_")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("018000010000096f0000000000000000000278f77b22636175637573223a6e756c6c2c22636f6e67726573735f6e756d62657273223a5b3131332c3131342c3131355d2c2263757272656e74223a747275652c226465736372697074696f6e223a224a756e696f722053656e61746f7220666f7220576973636f6e73696e222c226469737472696374223a6e756c6c2c22656e6464617465223a22323031392d30312d3033222c226578747261223a7b2261646472657373223a2237303920486172742053656e617465204f6666696365204275696c64696e672057617368696e67746f6e204443203230353130222c22636f6e746163745f666f726d223a2268747470733a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f765c2f666565646261636b222c22666178223a223230322d3232352d36393432222c226f6666696365223a2237303920486172742053656e617465204f6666696365204275696c64696e67222c227273735f75726c223a22687474703a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f765c2f7273735c2f66656564735c2f3f747970653d616c6c227d2c226c6561646572736869705f7469746c65223a6e756c6c2c227061727479223a2244656d6f63726174222c22706572736f6e223a7b2262696f67756964656964223a2242303031323330222c226269727468646179223a22313936322d30322d3131222c22637370616e6964223a35373838342c2266697273746e616d65223a2254616d6d79222c2267656e646572223a2266656d616c65222c2267656e6465725f6c6162656c223a2246656d616c65222c226c6173746e616d65223a2242616c6477696e222c226c696e6b223a2268747470733a5c2f5c2f7777772e676f76747261636b2e75735c2f636f6e67726573735c2f6d656d626572735c2f74616d6d795f62616c6477696e5c2f343030303133222c226d6964646c656e616d65223a22222c226e616d65223a2253656e2e2054616d6d792042616c6477696e205b442d57495d222c226e616d656d6f64223a22222c226e69636b6e616d65223a22222c226f736964223a224e3030303034333637222c227076736964223a2233343730222c22736f72746e616d65223a2242616c6477696e2c2054616d6d79202853656e2e29205b442d57495d222c22747769747465726964223a2253656e61746f7242616c6477696e222c22796f75747562656964223a22776974616d6d7962616c6477696e227d2c2270686f6e65223a223230322d3232342d35363533222c22726f6c655f74797065223a2273656e61746f72222c22726f6c655f747970655f6c6162656c223a2253656e61746f72222c2273656e61746f725f636c617373223a22636c61737331222c2273656e61746f725f636c6173735f6c6162656c223a22436c6173732031222c2273656e61746f725f72616e6b223a226a756e696f72222c2273656e61746f725f72616e6b5f6c6162656c223a224a756e696f72222c22737461727464617465223a22323031332d30312d3033222c227374617465223a225749222c227469746c65223a2253656e2e222c227469746c655f6c6f6e67223a2253656e61746f72222c2277656273697465223a2268747470733a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f76227d2c7b22636175637573223a6e756c6c2c22636f6e67726573735f6e756d62657273223a5b3131332c3131342c3131355d2c2263757272656e74223a747275652c226465736372697074696f6e223a2253656e696f722053656e61746f7220666f72204f68696f222c226469737472696374223a6e756c6c2c22656e6464617465223a22323031392d30312d3033222c226578747261223a7b2261646472657373223a2237313320486172742053656e617465204f6666696365204275696c64696e672057617368696e67746f6e204443203230353130222c22636f6e746163745f666f726d223a22687474703a5c2f5c2f7777772e62726f776e2e73656e6174652e676f765c2f636f6e746163745c2f222c22666178223a223230322d3232382d36333231222c226f6666696365223a2237313320486172742053656e617465204f6666696365204275696c64696e67222c227273735f75726c223a22687474703a5c2f5c2f7777772e62726f776e2e73656e6174652e676f765c2f7273735c2f66656564735c2f3f747970653d616c6c26616d703b227d2c226c6561646572736869705f7469746c65223a6e756c6c2c227061727479223a2244656d6f63726174222c22706572736f6e223a7b2262696f67756964656964223a2242303030393434222c226269727468646179223a22313935322d31312d3039222c22637370616e6964223a353035312c2266697273746e616d65223a2253686572726f64222c2267656e646572223a226d616c65222c2267656e6465725f6c6162656c223a224d616c65222c226c6173746e616d65223a2242726f776e222c226c696e6b223a2268747470733a5c2f5c2f7777772e676f76747261636b2e75735c2f636f6e67726573735c2f6d656d626572735c2f73686572726f645f62726f776e5c2f343030303530222c226d6964646c656e616d65223a22222c226e616d65223a2253656e2e2053686572726f642042726f776e205b442d4f485d222c226e616d656d6f64223a22222c226e69636b6e616d65223a22222c226f736964223a224e3030303033353335222c227076736964223a223237303138222c22736f72746e616d65223a2242726f776e2c2053686572726f64202853656e2e29205b442d4f485d222c22747769747465726964223a2253656e53686572726f6442726f776e222c22796f75747562656964223a2253686572726f6442726f776e4f68696f227d2c2270686f6e65223a223230322d3232342d32333135222c22726f6c655f74797065223a2273656e61746f72222c22726f6c655f747970655f6c6162656c223a2253656e61746f72222c2273656e61746f725f636c617373223a22636c61737331222c2273656e61746f725f636c6173735f6c6162656c223a22436c6173732031222c2273656e61746f725f72616e6b223a2273656e696f72222c2273656e61746f725f72616e6b5f6c6162656c223a2253656e696f72222c22737461727464617465223a22323031332d30312d3033222c227374617465223a224f48222c227469746c65223a2253656e2e222c227469746c655f6c6f6e67223a2253656e61746f72222c2277656273697465223a2268747470733a5c2f5c2f7777772e62726f776e2e73656e6174652e676f76227d2c0000000001800005000000140000000000000000000278f700000000000278f7000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=json%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0ICogZnJvbSBvc3NvYmplY3Qgd2hlcmUgcGFydHkgPSAnRGVtb2NyYXQnIGxpbWl0IDI=</Expression><InputSerialization><JSON><Type>LINES</Type></JSON></InputSerialization><OutputSerialization><JSON><RecordDelimiter>LA==</RecordDelimiter></JSON></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select * from ossobject where party = 'Democrat' limit 2"),
				InputSerializationSelect: InputSerializationSelect{
					JsonBodyInput: &JSONSelectInput{
						JSONType: Ptr("LINES"),
					},
				},
				OutputSerializationSelect: OutputSerializationSelect{
					JsonBodyOutput: &JSONSelectOutput{
						RecordDelimiter: Ptr(","),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "{\"caucus\":null,\"congress_")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(hexStrToByte("018000010000096f0000000000000000000278f77b22636175637573223a6e756c6c2c22636f6e67726573735f6e756d62657273223a5b3131332c3131342c3131355d2c2263757272656e74223a747275652c226465736372697074696f6e223a224a756e696f722053656e61746f7220666f7220576973636f6e73696e222c226469737472696374223a6e756c6c2c22656e6464617465223a22323031392d30312d3033222c226578747261223a7b2261646472657373223a2237303920486172742053656e617465204f6666696365204275696c64696e672057617368696e67746f6e204443203230353130222c22636f6e746163745f666f726d223a2268747470733a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f765c2f666565646261636b222c22666178223a223230322d3232352d36393432222c226f6666696365223a2237303920486172742053656e617465204f6666696365204275696c64696e67222c227273735f75726c223a22687474703a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f765c2f7273735c2f66656564735c2f3f747970653d616c6c227d2c226c6561646572736869705f7469746c65223a6e756c6c2c227061727479223a2244656d6f63726174222c22706572736f6e223a7b2262696f67756964656964223a2242303031323330222c226269727468646179223a22313936322d30322d3131222c22637370616e6964223a35373838342c2266697273746e616d65223a2254616d6d79222c2267656e646572223a2266656d616c65222c2267656e6465725f6c6162656c223a2246656d616c65222c226c6173746e616d65223a2242616c6477696e222c226c696e6b223a2268747470733a5c2f5c2f7777772e676f76747261636b2e75735c2f636f6e67726573735c2f6d656d626572735c2f74616d6d795f62616c6477696e5c2f343030303133222c226d6964646c656e616d65223a22222c226e616d65223a2253656e2e2054616d6d792042616c6477696e205b442d57495d222c226e616d656d6f64223a22222c226e69636b6e616d65223a22222c226f736964223a224e3030303034333637222c227076736964223a2233343730222c22736f72746e616d65223a2242616c6477696e2c2054616d6d79202853656e2e29205b442d57495d222c22747769747465726964223a2253656e61746f7242616c6477696e222c22796f75747562656964223a22776974616d6d7962616c6477696e227d2c2270686f6e65223a223230322d3232342d35363533222c22726f6c655f74797065223a2273656e61746f72222c22726f6c655f747970655f6c6162656c223a2253656e61746f72222c2273656e61746f725f636c617373223a22636c61737331222c2273656e61746f725f636c6173735f6c6162656c223a22436c6173732031222c2273656e61746f725f72616e6b223a226a756e696f72222c2273656e61746f725f72616e6b5f6c6162656c223a224a756e696f72222c22737461727464617465223a22323031332d30312d3033222c227374617465223a225749222c227469746c65223a2253656e2e222c227469746c655f6c6f6e67223a2253656e61746f72222c2277656273697465223a2268747470733a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f76227d2c7b22636175637573223a6e756c6c2c22636f6e67726573735f6e756d62657273223a5b3131332c3131342c3131355d2c2263757272656e74223a747275652c226465736372697074696f6e223a2253656e696f722053656e61746f7220666f72204f68696f222c226469737472696374223a6e756c6c2c22656e6464617465223a22323031392d30312d3033222c226578747261223a7b2261646472657373223a2237313320486172742053656e617465204f6666696365204275696c64696e672057617368696e67746f6e204443203230353130222c22636f6e746163745f666f726d223a22687474703a5c2f5c2f7777772e62726f776e2e73656e6174652e676f765c2f636f6e746163745c2f222c22666178223a223230322d3232382d36333231222c226f6666696365223a2237313320486172742053656e617465204f6666696365204275696c64696e67222c227273735f75726c223a22687474703a5c2f5c2f7777772e62726f776e2e73656e6174652e676f765c2f7273735c2f66656564735c2f3f747970653d616c6c26616d703b227d2c226c6561646572736869705f7469746c65223a6e756c6c2c227061727479223a2244656d6f63726174222c22706572736f6e223a7b2262696f67756964656964223a2242303030393434222c226269727468646179223a22313935322d31312d3039222c22637370616e6964223a353035312c2266697273746e616d65223a2253686572726f64222c2267656e646572223a226d616c65222c2267656e6465725f6c6162656c223a224d616c65222c226c6173746e616d65223a2242726f776e222c226c696e6b223a2268747470733a5c2f5c2f7777772e676f76747261636b2e75735c2f636f6e67726573735c2f6d656d626572735c2f73686572726f645f62726f776e5c2f343030303530222c226d6964646c656e616d65223a22222c226e616d65223a2253656e2e2053686572726f642042726f776e205b442d4f485d222c226e616d656d6f64223a22222c226e69636b6e616d65223a22222c226f736964223a224e3030303033353335222c227076736964223a223237303138222c22736f72746e616d65223a2242726f776e2c2053686572726f64202853656e2e29205b442d4f485d222c22747769747465726964223a2253656e53686572726f6442726f776e222c22796f75747562656964223a2253686572726f6442726f776e4f68696f227d2c2270686f6e65223a223230322d3232342d32333135222c22726f6c655f74797065223a2273656e61746f72222c22726f6c655f747970655f6c6162656c223a2253656e61746f72222c2273656e61746f725f636c617373223a22636c61737331222c2273656e61746f725f636c6173735f6c6162656c223a22436c6173732031222c2273656e61746f725f72616e6b223a2273656e696f72222c2273656e61746f725f72616e6b5f6c6162656c223a2253656e696f72222c22737461727464617465223a22323031332d30312d3033222c227374617465223a224f48222c227469746c65223a2253656e2e222c227469746c655f6c6f6e67223a2253656e61746f72222c2277656273697465223a2268747470733a5c2f5c2f7777772e62726f776e2e73656e6174652e676f76227d2c0000000001800005000000140000000000000000000278f700000000000278f7000000c800000000")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=json%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0ICogZnJvbSBvc3NvYmplY3Qgd2hlcmUgcGFydHkgPSAnRGVtb2NyYXQnIGxpbWl0IDI=</Expression><InputSerialization><JSON><Type>LINES</Type><Range>line-range=0-10</Range></JSON></InputSerialization><OutputSerialization><JSON><RecordDelimiter>LA==</RecordDelimiter></JSON></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select * from ossobject where party = 'Democrat' limit 2"),
				InputSerializationSelect: InputSerializationSelect{
					JsonBodyInput: &JSONSelectInput{
						JSONType: Ptr("LINES"),
						Range:    Ptr("0-10"),
					},
				},
				OutputSerializationSelect: OutputSerializationSelect{
					JsonBodyOutput: &JSONSelectOutput{
						RecordDelimiter: Ptr(","),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "{\"caucus\":null,\"congress_")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":        "534B371674E88A4D8906****",
			"Date":                    "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id":        "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
			"x-oss-select-output-raw": "true",
		},
		[]byte(hexStrToByte("7b2266697273746e616d65223a2254616d6d79222c226c6173746e616d65223a2242616c6477696e222c226578747261223a7b2261646472657373223a2237303920486172742053656e617465204f6666696365204275696c64696e672057617368696e67746f6e204443203230353130222c22636f6e746163745f666f726d223a2268747470733a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f765c2f666565646261636b222c22666178223a223230322d3232352d36393432222c226f6666696365223a2237303920486172742053656e617465204f6666696365204275696c64696e67222c227273735f75726c223a22687474703a5c2f5c2f7777772e62616c6477696e2e73656e6174652e676f765c2f7273735c2f66656564735c2f3f747970653d616c6c227d7d2c7b2266697273746e616d65223a2253686572726f64222c226c6173746e616d65223a2242726f776e222c226578747261223a7b2261646472657373223a2237313320486172742053656e617465204f6666696365204275696c64696e672057617368696e67746f6e204443203230353130222c22636f6e746163745f666f726d223a22687474703a5c2f5c2f7777772e62726f776e2e73656e6174652e676f765c2f636f6e746163745c2f222c22666178223a223230322d3232382d36333231222c226f6666696365223a2237313320486172742053656e617465204f6666696365204275696c64696e67222c227273735f75726c223a22687474703a5c2f5c2f7777772e62726f776e2e73656e6174652e676f765c2f7273735c2f66656564735c2f3f747970653d616c6c26616d703b227d7d2c")),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=json%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IHBlcnNvbi5maXJzdG5hbWUgYXMgYWFhIGFzIGZpcnN0bmFtZSwgcGVyc29uLmxhc3RuYW1lLCBleHRyYSBmcm9tIG9zc29iamVjdCBsaW1pdCAy</Expression><InputSerialization><JSON><Type>LINES</Type><Range>split-range=0-12</Range><ParseJsonNumberAsString>true</ParseJsonNumberAsString></JSON><CompressionType>NONE</CompressionType></InputSerialization><OutputSerialization><JSON><RecordDelimiter>LA==</RecordDelimiter></JSON><OutputRawData>true</OutputRawData><EnablePayloadCrc>false</EnablePayloadCrc></OutputSerialization><Options><SkipPartialDataRecord>false</SkipPartialDataRecord><MaxSkippedRecordsAllowed>2</MaxSkippedRecordsAllowed></Options></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select person.firstname as aaa as firstname, person.lastname, extra from ossobject limit 2"),
				InputSerializationSelect: InputSerializationSelect{
					CompressionType: Ptr("NONE"),
					JsonBodyInput: &JSONSelectInput{
						JSONType:                Ptr("LINES"),
						ParseJSONNumberAsString: Ptr(true),
						SplitRange:              Ptr("0-12"),
					},
				},
				OutputSerializationSelect: OutputSerializationSelect{
					JsonBodyOutput: &JSONSelectOutput{
						RecordDelimiter: Ptr(","),
					},
					OutputRawData:    Ptr(true),
					EnablePayloadCrc: Ptr(false),
				},
				SelectOptions: &SelectOptions{
					SkipPartialDataRecord:    Ptr(false),
					MaxSkippedRecordsAllowed: Ptr(2),
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			dataByte, _ := io.ReadAll(o.Body)
			assert.Equal(t, string(dataByte[:25]), "{\"firstname\":\"Tammy\",\"las")
		},
	},
}

func TestMockSelectObject_Success(t *testing.T) {
	for _, c := range testMockSelectObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SelectObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSelectObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SelectObjectRequest
	CheckOutputFn  func(t *testing.T, o *SelectObjectResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidOutputFieldDelimiter</Code>
  <Message>Invalid FieldDelimiter parameter:|</Message>
  <RequestId>6569ADEDF1BF4B6AE588****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0016-00000828</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000828</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xLCBfMiBmcm9tIG9zc29iamVjdA==</Expression><InputSerialization></InputSerialization><OutputSerialization><CSV><RecordDelimiter>&#xD;&#xA;</RecordDelimiter><FieldDelimiter>,</FieldDelimiter></CSV></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			SelectRequest: &SelectRequest{
				Expression: Ptr("select _1, _2 from ossobject"),
				OutputSerializationSelect: OutputSerializationSelect{
					CsvBodyOutput: &CSVSelectOutput{
						RecordDelimiter: Ptr("\r\n"),
						FieldDelimiter:  Ptr(","),
					},
				},
			},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidOutputFieldDelimiter", serr.Code)
			assert.Equal(t, "Invalid FieldDelimiter parameter:|", serr.Message)
			assert.Equal(t, "6569ADEDF1BF4B6AE588****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "65699DB6E6F906F45A83****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchKey</Code>
  <Message>The specified key does not exist.</Message>
  <RequestId>65699DB6E6F906F45A83****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <Key>object</Key>
  <EC>0026-00000001</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0026-00000001</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process=csv%2Fselect", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<SelectRequest><InputSerialization></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))
		},
		&SelectObjectRequest{
			Bucket:        Ptr("bucket"),
			Key:           Ptr("object"),
			SelectRequest: &SelectRequest{},
		},
		func(t *testing.T, o *SelectObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "65699DB6E6F906F45A83****", serr.RequestID)
		},
	},
}

func TestMockSelectObject_Error(t *testing.T) {
	for _, c := range testMockSelectObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SelectObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockProcessObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ProcessObjectRequest
	CheckOutputFn  func(t *testing.T, o *ProcessObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
    "bucket": "",
    "fileSize": 3267,
    "object": "dest.jpg",
    "status": "OK"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_ZGVzdC5qcGc=", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("dest.jpg")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Bucket, "")
			assert.Equal(t, o.FileSize, 3267)
			assert.Equal(t, o.Object, "dest.jpg")
			assert.Equal(t, o.ProcessStatus, "OK")
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
    "bucket": "dest-bucket",
    "fileSize": 3267,
    "object": "dest.jpg",
    "status": "OK"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_ZGVzdC5qcGc=,b_ZGVzdC1idWNrZXQ=", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v,b_%v", base64.URLEncoding.EncodeToString([]byte("dest.jpg")), base64.URLEncoding.EncodeToString([]byte("dest-bucket")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Bucket, "dest-bucket")
			assert.Equal(t, o.FileSize, 3267)
			assert.Equal(t, o.Object, "dest.jpg")
			assert.Equal(t, o.ProcessStatus, "OK")
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
    "bucket": "",
    "fileSize": 3267,
    "object": "dest.jpg",
    "status": "OK"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_ZGVzdC5qcGc=", string(data))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			Process:      Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("dest.jpg")))),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Bucket, "")
			assert.Equal(t, o.FileSize, 3267)
			assert.Equal(t, o.Object, "dest.jpg")
			assert.Equal(t, o.ProcessStatus, "OK")
		},
	},
}

func TestMockProcessObject_Success(t *testing.T) {
	for _, c := range testMockProcessObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ProcessObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockProcessObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ProcessObjectRequest
	CheckOutputFn  func(t *testing.T, o *ProcessObjectResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>SignatureDoesNotMatch</Code>
				<Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_a2V5LWRlc3QuanBn", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("key-dest.jpg")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "SignatureDoesNotMatch", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
		},
	},
	{
		400,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>operation not support post: test</Message>
  <RequestId>65467C42E001B4333337****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_a2V5LWRlc3QuanBn", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("key-dest.jpg")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "operation not support post: test")
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_a2V5LWRlc3QuanBn", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("key-dest.jpg")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ProcessObject fail")
		},
	},
}

func TestMockProcessObject_Error(t *testing.T) {
	for _, c := range testMockProcessObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ProcessObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAsyncProcessObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AsyncProcessObjectRequest
	CheckOutputFn  func(t *testing.T, o *AsyncProcessObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{"EventId":"181-1kZUlN60OH4fWOcOjZEnGnG****","RequestId":"1D99637F-F59E-5B41-9200-C4892F52****","TaskId":"MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.EventId, "181-1kZUlN60OH4fWOcOjZEnGnG****")
			assert.Equal(t, o.RequestId, "1D99637F-F59E-5B41-9200-C4892F52****")
			assert.Equal(t, o.TaskId, "MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{"EventId":"181-1kZUlN60OH4fWOcOjZEnGnG****","RequestId":"1D99637F-F59E-5B41-9200-C4892F52****","TaskId":"MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.EventId, "181-1kZUlN60OH4fWOcOjZEnGnG****")
			assert.Equal(t, o.RequestId, "1D99637F-F59E-5B41-9200-C4892F52****")
			assert.Equal(t, o.TaskId, "MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****")
		},
	},
}

func TestMockAsyncProcessObject_Success(t *testing.T) {
	for _, c := range testMockAsyncProcessObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AsyncProcessObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAsyncProcessObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AsyncProcessObjectRequest
	CheckOutputFn  func(t *testing.T, o *AsyncProcessObjectResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>SignatureDoesNotMatch</Code>
				<Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "SignatureDoesNotMatch", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
		},
	},
	{
		400,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidRequest</Code>
  <Message>no x-oss-async-process parameter found</Message>
  <RequestId>65467C42E001B4333337****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidRequest", serr.Code)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "no x-oss-async-process parameter found")
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute AsyncProcessObject fail")
		},
	},
}

func TestMockAsyncProcessObject_Error(t *testing.T) {
	for _, c := range testMockAsyncProcessObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AsyncProcessObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketRequestPaymentSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRequestPaymentRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketRequestPaymentResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?requestPayment", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RequestPaymentConfiguration><Payer>Requester</Payer></RequestPaymentConfiguration>")
		},
		&PutBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
			PaymentConfiguration: &RequestPaymentConfiguration{
				Payer: Requester,
			},
		},
		func(t *testing.T, o *PutBucketRequestPaymentResult, err error) {
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
			assert.Equal(t, "/bucket?requestPayment", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RequestPaymentConfiguration><Payer>BucketOwner</Payer></RequestPaymentConfiguration>")
		},
		&PutBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
			PaymentConfiguration: &RequestPaymentConfiguration{
				Payer: BucketOwner,
			},
		},
		func(t *testing.T, o *PutBucketRequestPaymentResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketRequestPayment_Success(t *testing.T) {
	for _, c := range testMockPutBucketRequestPaymentSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketRequestPayment(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketRequestPaymentErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRequestPaymentRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketRequestPaymentResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?requestPayment", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RequestPaymentConfiguration><Payer>BucketOwner</Payer></RequestPaymentConfiguration>")
		},
		&PutBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
			PaymentConfiguration: &RequestPaymentConfiguration{
				Payer: BucketOwner,
			},
		},
		func(t *testing.T, o *PutBucketRequestPaymentResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?requestPayment", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RequestPaymentConfiguration><Payer>BucketOwner</Payer></RequestPaymentConfiguration>")
		},
		&PutBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
			PaymentConfiguration: &RequestPaymentConfiguration{
				Payer: BucketOwner,
			},
		},
		func(t *testing.T, o *PutBucketRequestPaymentResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?requestPayment", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RequestPaymentConfiguration><Payer>BucketOwner</Payer></RequestPaymentConfiguration>")
		},
		&PutBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
			PaymentConfiguration: &RequestPaymentConfiguration{
				Payer: BucketOwner,
			},
		},
		func(t *testing.T, o *PutBucketRequestPaymentResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutBucketRequestPayment fail")
		},
	},
}

func TestMockPutBucketRequestPayment_Error(t *testing.T) {
	for _, c := range testMockPutBucketRequestPaymentErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketRequestPayment(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketRequestPaymentSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketRequestPaymentRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketRequestPaymentResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RequestPaymentConfiguration>
  <Payer>Requester</Payer>
</RequestPaymentConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?requestPayment", r.URL.String())
		},
		&GetBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRequestPaymentResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "Requester", *o.Payer)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RequestPaymentConfiguration>
  <Payer>BucketOwner</Payer>
</RequestPaymentConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?requestPayment", r.URL.String())
		},
		&GetBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRequestPaymentResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "BucketOwner", *o.Payer)
		},
	},
}

func TestMockGetBucketRequestPayment_Success(t *testing.T) {
	for _, c := range testMockGetBucketRequestPaymentSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketRequestPayment(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketRequestPaymentErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketRequestPaymentRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketRequestPaymentResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?requestPayment", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRequestPaymentResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?requestPayment", strUrl)
		},
		&GetBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRequestPaymentResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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
			assert.Equal(t, "/bucket?requestPayment", strUrl)
		},
		&GetBucketRequestPaymentRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRequestPaymentResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketRequestPayment fail")
		},
	},
}

func TestMockGetBucketRequestPayment_Error(t *testing.T) {
	for _, c := range testMockGetBucketRequestPaymentErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketRequestPayment(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
