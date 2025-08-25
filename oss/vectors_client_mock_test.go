package oss

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testVectorsInvokeOperationAnonymousCases = []struct {
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
		&OperationInput{
			OpName: "PutVectorBucket",
			Method: "PUT",
			Bucket: Ptr("bucket"),
			Body:   strings.NewReader(`{"encryptionConfiguration": {"KMSMasterKeyID": "string","SSEAlgorithm": "string"},"vectorBucketName": "string"}`),
		},
		func(t *testing.T, o *OperationOutput) {
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
		&OperationInput{
			OpName: "GetVectorBucket",
			Bucket: Ptr("bucket"),
			Method: "GET",
			Parameters: map[string]string{
				"bucketInfo": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
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
			assert.Equal(t, "/bucket/?DeleteVectorIndex", r.URL.String())
			assert.Equal(t, "DELETE", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "{\"indexName\": \"string\"}", string(requestBody))
		},
		&OperationInput{
			OpName: "DeleteVectorIndex",
			Bucket: Ptr("bucket"),
			Method: "DELETE",
			Parameters: map[string]string{
				"DeleteVectorIndex": "",
			},
			Body: strings.NewReader(`{"indexName": "string"}`),
		},
		func(t *testing.T, o *OperationOutput) {
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

		cfg := LoadDefaultConfig().
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
	Input          *OperationInput
	CheckOutputFn  func(t *testing.T, o *OperationOutput, err error)
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
		&OperationInput{
			OpName: "PutVectorBucket",
			Method: "PUT",
		},
		func(t *testing.T, o *OperationOutput, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
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

		cfg := LoadDefaultConfig().
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
			Bucket: Ptr("bucket"),
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
			Bucket:          Ptr("bucket"),
			ResourceGroupId: Ptr("rg-aek27tc****"),
			Tagging:         Ptr("k1=v1&k2=v2"),
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

		cfg := LoadDefaultConfig().
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *PutVectorBucketResult, err error) {
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *PutVectorBucketResult, err error) {
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

func TestMockPutVectorBucket_Error(t *testing.T) {
	for _, c := range testMockPutVectorBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
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
						"Bucket": {
						  "CreationDate": "2013-07-31T10:56:21.000Z",
						  "ExtranetEndpoint": "oss-cn-hangzhou.aliyuncs.com",
						  "IntranetEndpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
						  "Location": "oss-cn-hangzhou",
						  "Name": "oss-example",
						  "ResourceGroupId": "rg-aek27tc********"
						}
					  }
					}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?bucketInfo", r.URL.String())
		},
		&GetVectorBucketRequest{
			Bucket: Ptr("bucket"),
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

		cfg := LoadDefaultConfig().
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetVectorBucketResult, err error) {
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetVectorBucketResult, err error) {
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
			assert.Equal(t, "/bucket/?bucketInfo", strUrl)
		},
		&GetVectorBucketRequest{
			Bucket: Ptr("bucket"),
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

		cfg := LoadDefaultConfig().
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
    "Buckets": {
      "Bucket": [
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
			assert.Equal(t, len(o.Buckets.Bucket), 2)
			assert.Equal(t, *o.Buckets.Bucket[0].CreationDate, time.Date(2014, time.February, 17, 18, 12, 43, 0, time.UTC))
			assert.Equal(t, *o.Buckets.Bucket[0].ExtranetEndpoint, "oss-cn-shanghai.aliyuncs.com")
			assert.Equal(t, *o.Buckets.Bucket[0].IntranetEndpoint, "oss-cn-shanghai-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets.Bucket[0].Name, "app-base-oss")
			assert.Equal(t, *o.Buckets.Bucket[0].Region, "cn-shanghai")
			assert.Equal(t, *o.Buckets.Bucket[0].Location, "oss-cn-shanghai")
			assert.Equal(t, *o.Buckets.Bucket[0].ResourceGroupId, "rg-aek27ta********")

			assert.Equal(t, *o.Buckets.Bucket[1].CreationDate, time.Date(2014, time.February, 25, 11, 21, 04, 0, time.UTC))
			assert.Equal(t, *o.Buckets.Bucket[1].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.Buckets.Bucket[1].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets.Bucket[1].Name, "mybucket")
			assert.Equal(t, *o.Buckets.Bucket[1].Region, "cn-hangzhou")
			assert.Equal(t, *o.Buckets.Bucket[1].Location, "oss-cn-hangzhou")
			assert.Equal(t, *o.Buckets.Bucket[1].ResourceGroupId, "rg-aek27tc********")
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
    "Buckets": {
      "Bucket": [{
        "CreationDate": "2014-05-14T11:18:32.000Z",
        "ExtranetEndpoint": "oss-cn-hangzhou.aliyuncs.com",
        "IntranetEndpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
        "Location": "oss-cn-hangzhou",
        "Name": "mybucket01",
        "Region": "cn-hangzhou",
        "ResourceGroupId": "rg-aek27tc********"
      }]
    }
  }
}`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/?marker&max-keys=10&prefix=%2F", strUrl)
			assert.Equal(t, "rg-aek27tc********", r.Header.Get("x-oss-resource-group-id"))
		},
		&ListVectorBucketsRequest{
			Marker:          Ptr(""),
			MaxKeys:         10,
			Prefix:          Ptr("/"),
			ResourceGroupId: Ptr("rg-aek27tc********"),
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

			assert.Equal(t, len(o.Buckets.Bucket), 1)
			assert.Equal(t, *o.Buckets.Bucket[0].CreationDate, time.Date(2014, time.May, 14, 11, 18, 32, 0, time.UTC))
			assert.Equal(t, *o.Buckets.Bucket[0].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.Buckets.Bucket[0].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets.Bucket[0].Name, "mybucket01")
			assert.Equal(t, *o.Buckets.Bucket[0].Region, "cn-hangzhou")
			assert.Equal(t, *o.Buckets.Bucket[0].Location, "oss-cn-hangzhou")
			assert.Equal(t, *o.Buckets.Bucket[0].ResourceGroupId, "rg-aek27tc********")
		},
	},
}

func TestMockListVectorBuckets_Success(t *testing.T) {
	for _, c := range testMockListVectorBucketsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
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

		cfg := LoadDefaultConfig().
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
			Bucket: Ptr("bucket"),
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

		cfg := LoadDefaultConfig().
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteVectorBucketResult, err error) {
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteVectorBucketResult, err error) {
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

func TestMockDeleteVectorBucket_Error(t *testing.T) {
	for _, c := range testMockDeleteVectorBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteVectorBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketPolicyForVectorSuccessCases = []struct {
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
			Bucket: Ptr("bucket"),
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

func TestMockPutBucketPolicyForVector_Success(t *testing.T) {
	for _, c := range testMockPutBucketPolicyForVectorSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketPolicyForVectorErrorCases = []struct {
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
			Bucket: Ptr("bucket"),
			Body:   strings.NewReader(`{"Version":"1","Statement":[{"Action":["ossvector:PutVectors","ossvector:GetVectors"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:ossvector:*:1234567890:*"]}]}`),
		},
		func(t *testing.T, o *PutBucketPolicyResult, err error) {
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
			Bucket: Ptr("bucket"),
			Body:   strings.NewReader(`{"Version":"1","Statement":[{"Action":["ossvector:PutVectors","ossvector:GetVectors"],"Effect":"Deny","Principal":["1234567890"],"Resource":["acs:ossvector:*:1234567890:*"]}]}`),
		},
		func(t *testing.T, o *PutBucketPolicyResult, err error) {
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

func TestMockPutBucketPolicyForVector_Error(t *testing.T) {
	for _, c := range testMockPutBucketPolicyForVectorErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketPolicyForVectorSuccessCases = []struct {
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
			Bucket: Ptr("bucket"),
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

func TestMockGetBucketPolicyForVector_Success(t *testing.T) {
	for _, c := range testMockGetBucketPolicyForVectorSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketPolicyForVectorErrorCases = []struct {
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketPolicyResult, err error) {
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketPolicyResult, err error) {
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

func TestMockGetBucketPolicyForVector_Error(t *testing.T) {
	for _, c := range testMockGetBucketPolicyForVectorErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketPolicyForVectorSuccessCases = []struct {
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketPolicyResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteBucketPolicyForVector_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketPolicyForVectorSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketPolicyForVectorErrorCases = []struct {
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketPolicyResult, err error) {
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
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketPolicyResult, err error) {
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

func TestMockDeleteBucketPolicyForVector_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketPolicyForVectorErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketResourceGroupForVectorForVectorSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketResourceGroupRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketResourceGroupResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?resourceGroup", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"BucketResourceGroupConfiguration\":{\"ResourceGroupId\":\"rg-aekz****\"}}")
		},
		&PutBucketResourceGroupRequest{
			Bucket: Ptr("bucket"),
			BucketResourceGroupConfiguration: &BucketResourceGroupConfiguration{
				Ptr("rg-aekz****"),
			},
		},
		func(t *testing.T, o *PutBucketResourceGroupResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketResourceGroupForVector_Success(t *testing.T) {
	for _, c := range testMockPutBucketResourceGroupForVectorForVectorSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketResourceGroup(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketResourceGroupForVectorErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketResourceGroupRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketResourceGroupResult, err error)
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
			assert.Equal(t, "/bucket/?resourceGroup", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"BucketResourceGroupConfiguration\":{\"ResourceGroupId\":\"rg-aekz****\"}}")
		},
		&PutBucketResourceGroupRequest{
			Bucket: Ptr("bucket"),
			BucketResourceGroupConfiguration: &BucketResourceGroupConfiguration{
				Ptr("rg-aekz****"),
			},
		},
		func(t *testing.T, o *PutBucketResourceGroupResult, err error) {
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
			assert.Equal(t, "/bucket/?resourceGroup", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"BucketResourceGroupConfiguration\":{\"ResourceGroupId\":\"rg-aekz****\"}}")
		},
		&PutBucketResourceGroupRequest{
			Bucket: Ptr("bucket"),
			BucketResourceGroupConfiguration: &BucketResourceGroupConfiguration{
				Ptr("rg-aekz****"),
			},
		},
		func(t *testing.T, o *PutBucketResourceGroupResult, err error) {
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

func TestMockPutBucketResourceGroupForVector_Error(t *testing.T) {
	for _, c := range testMockPutBucketResourceGroupForVectorErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketResourceGroup(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketResourceGroupForVectorForVectorSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketResourceGroupRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketResourceGroupResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "BucketResourceGroupConfiguration": {
    "ResourceGroupId": "rg-aekz****"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?resourceGroup", r.URL.String())
		},
		&GetBucketResourceGroupRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketResourceGroupResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "rg-aekz****", *o.BucketResourceGroupConfiguration.ResourceGroupId)
		},
	},
}

func TestMockGetBucketResourceGroupForVector_Success(t *testing.T) {
	for _, c := range testMockGetBucketResourceGroupForVectorForVectorSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketResourceGroup(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketResourceGroupForVectorErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketResourceGroupRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketResourceGroupResult, err error)
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
			assert.Equal(t, "/bucket/?resourceGroup", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketResourceGroupRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketResourceGroupResult, err error) {
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
			assert.Equal(t, "/bucket/?resourceGroup", strUrl)
		},
		&GetBucketResourceGroupRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketResourceGroupResult, err error) {
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
			assert.Equal(t, "/bucket/?resourceGroup", strUrl)
		},
		&GetBucketResourceGroupRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketResourceGroupResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketResourceGroup fail")
		},
	},
}

func TestMockGetBucketResourceGroupForVector_Error(t *testing.T) {
	for _, c := range testMockGetBucketResourceGroupForVectorErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketResourceGroup(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketTagsForVectorSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketTagsRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketTagsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?tagging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"Tagging\":{\"TagSet\":{\"Tag\":[{\"Key\":\"key1\",\"Value\":\"value1\"},{\"Key\":\"key2\",\"Value\":\"value2\"}]}}}")
		},
		&PutBucketTagsRequest{
			Bucket: Ptr("bucket"),
			Tagging: &Tagging{
				&TagSet{
					[]Tag{
						{
							Ptr("key1"),
							Ptr("value1"),
						},
						{
							Ptr("key2"),
							Ptr("value2"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketTagsResult, err error) {
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
			assert.Equal(t, "/bucket/?tagging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"Tagging\":{\"TagSet\":{\"Tag\":[{\"Key\":\"key1\",\"Value\":\"value1\"}]}}}")
		},
		&PutBucketTagsRequest{
			Bucket: Ptr("bucket"),
			Tagging: &Tagging{
				&TagSet{
					[]Tag{
						{
							Ptr("key1"),
							Ptr("value1"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketTagsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketTagsForVector_Success(t *testing.T) {
	for _, c := range testMockPutBucketTagsForVectorSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketTags(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketTagsForVectorErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketTagsRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketTagsResult, err error)
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
			assert.Equal(t, "/bucket/?tagging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"Tagging\":{\"TagSet\":{\"Tag\":[{\"Key\":\"key1\",\"Value\":\"value1\"},{\"Key\":\"key2\",\"Value\":\"value2\"}]}}}")
		},
		&PutBucketTagsRequest{
			Bucket: Ptr("bucket"),
			Tagging: &Tagging{
				&TagSet{
					[]Tag{
						{
							Ptr("key1"),
							Ptr("value1"),
						},
						{
							Ptr("key2"),
							Ptr("value2"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketTagsResult, err error) {
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
			assert.Equal(t, "/bucket/?tagging", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"Tagging\":{\"TagSet\":{\"Tag\":[{\"Key\":\"key1\",\"Value\":\"value1\"},{\"Key\":\"key2\",\"Value\":\"value2\"}]}}}")
		},
		&PutBucketTagsRequest{
			Bucket: Ptr("bucket"),
			Tagging: &Tagging{
				&TagSet{
					[]Tag{
						{
							Ptr("key1"),
							Ptr("value1"),
						},
						{
							Ptr("key2"),
							Ptr("value2"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketTagsResult, err error) {
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

func TestMockPutBucketTagsForVector_Error(t *testing.T) {
	for _, c := range testMockPutBucketTagsForVectorErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketTags(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketTagsForVectorSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketTagsRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketTagsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
  "Tagging": {
    "TagSet": {
      "Tag": [
        {
          "Key": "testa",
          "Value": "value1"
        },
        {
          "Key": "testb",
          "Value": "value2"
        }
      ]
    }
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?tagging", r.URL.String())
		},
		&GetBucketTagsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketTagsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/json", o.Headers.Get("Content-Type"))
			assert.Equal(t, len(o.Tagging.TagSet.Tags), 2)
			assert.Equal(t, *o.Tagging.TagSet.Tags[0].Key, "testa")
			assert.Equal(t, *o.Tagging.TagSet.Tags[1].Value, "value2")
		},
	},
}

func TestMockGetBucketTagsForVector_Success(t *testing.T) {
	for _, c := range testMockGetBucketTagsForVectorSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketTags(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketTagsForVectorErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketTagsRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketTagsResult, err error)
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
			assert.Equal(t, "/bucket/?tagging", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketTagsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketTagsResult, err error) {
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
			assert.Equal(t, "/bucket/?tagging", strUrl)
		},
		&GetBucketTagsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketTagsResult, err error) {
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
			assert.Equal(t, "/bucket/?tagging", strUrl)
		},
		&GetBucketTagsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketTagsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketTags fail")
		},
	},
}

func TestMockGetBucketTagsForVector_Error(t *testing.T) {
	for _, c := range testMockGetBucketTagsForVectorErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketTags(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketTagsForVectorSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketTagsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketTagsResult, err error)
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
			assert.Equal(t, "/bucket/?tagging", strUrl)
		},
		&DeleteBucketTagsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketTagsResult, err error) {
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
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?tagging=k1%2Ck2", strUrl)
		},
		&DeleteBucketTagsRequest{
			Bucket:  Ptr("bucket"),
			Tagging: Ptr("k1,k2"),
		},
		func(t *testing.T, o *DeleteBucketTagsResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteBucketTagsForVector_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketTagsForVectorSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketTags(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketTagsForVectorErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketTagsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketTagsResult, err error)
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
			assert.Equal(t, "/bucket/?tagging", strUrl)
		},
		&DeleteBucketTagsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketTagsResult, err error) {
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
			assert.Equal(t, "/bucket/?tagging", strUrl)
		},
		&DeleteBucketTagsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketTagsResult, err error) {
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

func TestMockDeleteBucketTagsForVector_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketTagsForVectorErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketTags(context.TODO(), c.Request)
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
			assert.Equal(t, "/bucket/?PutVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\",\"metadata\":{\"nonFilterableMetadataKeys\":[\"foo\",\"bar\"]}}")
		},
		&PutVectorIndexRequest{
			Bucket:         Ptr("bucket"),
			DataType:       Ptr("string"),
			Dimension:      Ptr(128),
			DistanceMetric: Ptr("cosine"),
			IndexName:      Ptr("exampleIndex"),
			Metadata: map[string]interface{}{
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

		cfg := LoadDefaultConfig().
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
			assert.Equal(t, "/bucket/?PutVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\",\"metadata\":{\"nonFilterableMetadataKeys\":[\"foo\",\"bar\"]}}")
		},
		&PutVectorIndexRequest{
			Bucket:         Ptr("bucket"),
			DataType:       Ptr("string"),
			Dimension:      Ptr(128),
			DistanceMetric: Ptr("cosine"),
			IndexName:      Ptr("exampleIndex"),
			Metadata: map[string]interface{}{
				"nonFilterableMetadataKeys": []string{"foo", "bar"},
			},
		},
		func(t *testing.T, o *PutVectorIndexResult, err error) {
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
			assert.Equal(t, "/bucket/?PutVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\",\"metadata\":{\"nonFilterableMetadataKeys\":[\"foo\",\"bar\"]}}")
		},
		&PutVectorIndexRequest{
			Bucket:         Ptr("bucket"),
			DataType:       Ptr("string"),
			Dimension:      Ptr(128),
			DistanceMetric: Ptr("cosine"),
			IndexName:      Ptr("exampleIndex"),
			Metadata: map[string]interface{}{
				"nonFilterableMetadataKeys": []string{"foo", "bar"},
			},
		},
		func(t *testing.T, o *PutVectorIndexResult, err error) {
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

func TestMockPutVectorIndex_Error(t *testing.T) {
	for _, c := range testMockPutVectorIndexErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
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
      "createTime": "2025-08-02T10:49:17.289372919+08:00",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "cosine",
      "indexName": "exampleIndex",
      "metadata": { 
         "nonFilterableMetadataKeys": ["foo", "bar"]
      },
      "status": "running"
   },
   "vectorBucketName": "bucket"
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/?GetVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"exampleIndex\"}")
		},
		&GetVectorIndexRequest{
			Bucket:    Ptr("bucket"),
			IndexName: Ptr("exampleIndex"),
		},
		func(t *testing.T, o *GetVectorIndexResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Index.CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.Local))
			assert.Equal(t, *o.Index.DataType, "string")
			assert.Equal(t, *o.Index.Dimension, 128)
			assert.Equal(t, *o.Index.DistanceMetric, "cosine")
			assert.Equal(t, *o.Index.IndexName, "exampleIndex")
			//assert.Equal(t, len(o.Index.Metadata.NonFilterableMetadataKeys), 2)
			//assert.Equal(t, o.Index.Metadata.NonFilterableMetadataKeys[0], "foo")
			//assert.Equal(t, o.Index.Metadata.NonFilterableMetadataKeys[1], "bar")
			assert.Equal(t, *o.Index.Status, "running")
			assert.Equal(t, *o.VectorBucketName, "bucket")
		},
	},
}

func TestMockGetVectorIndex_Success(t *testing.T) {
	for _, c := range testMockGetVectorIndexSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
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
			assert.Equal(t, "/bucket/?GetVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"demoIndex\"}")
		},
		&GetVectorIndexRequest{
			Bucket:    Ptr("bucket"),
			IndexName: Ptr("demoIndex"),
		},
		func(t *testing.T, o *GetVectorIndexResult, err error) {
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
			assert.Equal(t, "/bucket/?GetVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"demoIndex\"}")
		},
		&GetVectorIndexRequest{
			Bucket:    Ptr("bucket"),
			IndexName: Ptr("demoIndex"),
		},
		func(t *testing.T, o *GetVectorIndexResult, err error) {
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

func TestMockGetVectorIndex_Error(t *testing.T) {
	for _, c := range testMockGetVectorIndexErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
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
      "createTime": "2025-08-02T10:49:17.289372919+08:00",
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
      "createTime": "2025-08-20T10:49:17.289372919+08:00",
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
			assert.Equal(t, "/bucket/?ListVectorIndexes", r.URL.String())
		},
		&ListVectorIndexesRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListVectorIndexesResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Indexes), 2)
			assert.Equal(t, *o.Indexes[0].CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.Local))
			assert.Equal(t, *o.Indexes[0].DataType, "string")
			assert.Equal(t, *o.Indexes[0].Dimension, 128)
			assert.Equal(t, *o.Indexes[0].DistanceMetric, "string")
			assert.Equal(t, *o.Indexes[0].IndexName, "demo1")
			assert.Len(t, o.Indexes[0].Metadata["nonFilterableMetadataKeys"], 2)
			if metadataValue, ok := o.Indexes[0].Metadata["nonFilterableMetadataKeys"]; ok {
				if keys, ok := metadataValue.([]interface{}); ok {
					assert.Equal(t, keys[0].(string), "foo")
					assert.Equal(t, keys[1].(string), "bar")
				}
			}
			assert.Equal(t, *o.Indexes[0].Status, "running")
			assert.Equal(t, *o.Indexes[0].VectorBucketName, "bucket")

			assert.Equal(t, *o.Indexes[1].CreateTime, time.Date(2025, time.August, 20, 10, 49, 17, 289372919, time.Local))
			assert.Equal(t, *o.Indexes[1].DataType, "string")
			assert.Equal(t, *o.Indexes[1].Dimension, 128)
			assert.Equal(t, *o.Indexes[1].DistanceMetric, "string")
			assert.Equal(t, *o.Indexes[1].IndexName, "demo2")
			assert.Len(t, o.Indexes[1].Metadata["nonFilterableMetadataKeys"], 2)
			if metadataValue, ok := o.Indexes[1].Metadata["nonFilterableMetadataKeys"]; ok {
				if keys, ok := metadataValue.([]interface{}); ok {
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

		cfg := LoadDefaultConfig().
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
			assert.Equal(t, "/bucket/?ListVectorIndexes", r.URL.String())
		},
		&ListVectorIndexesRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListVectorIndexesResult, err error) {
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
			assert.Equal(t, "/bucket/?ListVectorIndexes", r.URL.String())
		},
		&ListVectorIndexesRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListVectorIndexesResult, err error) {
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

func TestMockListVectorIndexes_Error(t *testing.T) {
	for _, c := range testMockListVectorIndexesErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
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
			assert.Equal(t, "/bucket/?DeleteVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"exampleIndex\"}")
		},
		&DeleteVectorIndexRequest{
			Bucket:    Ptr("bucket"),
			IndexName: Ptr("exampleIndex"),
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

		cfg := LoadDefaultConfig().
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
			assert.Equal(t, "/bucket/?DeleteVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"demoIndex\"}")
		},
		&DeleteVectorIndexRequest{
			Bucket:    Ptr("bucket"),
			IndexName: Ptr("demoIndex"),
		},
		func(t *testing.T, o *DeleteVectorIndexResult, err error) {
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
			assert.Equal(t, "/bucket/?DeleteVectorIndex", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "{\"indexName\":\"demoIndex\"}")
		},
		&DeleteVectorIndexRequest{
			Bucket:    Ptr("bucket"),
			IndexName: Ptr("demoIndex"),
		},
		func(t *testing.T, o *DeleteVectorIndexResult, err error) {
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

func TestMockDeleteVectorIndex_Error(t *testing.T) {
	for _, c := range testMockDeleteVectorIndexErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewVectorsClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteVectorIndex(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
