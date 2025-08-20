package oss

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

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

		},
		&OperationInput{
			OpName: "PutVectorBucket",
			Method: "PUT",
			Bucket: Ptr("bucket"),
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
