package oss

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutBucketStorageQuotaSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketStorageQuotaRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketStorageQuotaResult, err error)
}{
	{
		http.StatusOK,
		map[string]string{"x-oss-request-id": "534B371674E88A4D8906****", "Date": "Fri, 24 Feb 2017 03:15:40 GMT"},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, http.MethodPut, r.Method)
			assert.Equal(t, "/bucket/?quota", sortQuery(r))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "<QuotaConfiguration><StorageQuota>10737418240</StorageQuota><Mode>Strict</Mode></QuotaConfiguration>", string(data))
		},
		&PutBucketStorageQuotaRequest{Bucket: Ptr("bucket"), QuotaConfiguration: &QuotaConfiguration{StorageQuota: Ptr(int64(10737418240)), Mode: StorageQuotaModeStrict}},
		func(t *testing.T, o *PutBucketStorageQuotaResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketStorageQuota_Success(t *testing.T) {
	for _, c := range testMockPutBucketStorageQuotaSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketStorageQuota(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketStorageQuotaSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketStorageQuotaRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketStorageQuotaResult, err error)
}{
	{
		http.StatusOK,
		map[string]string{"x-oss-request-id": "534B371674E88A4D8906****", "Date": "Fri, 24 Feb 2017 03:15:40 GMT"},
		[]byte(`<QuotaConfiguration><StorageQuota>10737418240</StorageQuota><Mode>Warning</Mode><CurrentUsage>1049600</CurrentUsage></QuotaConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/bucket/?quota", sortQuery(r))
		},
		&GetBucketStorageQuotaRequest{Bucket: Ptr("bucket")},
		func(t *testing.T, o *GetBucketStorageQuotaResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, o.StatusCode)
			assert.Equal(t, int64(10737418240), *o.QuotaConfiguration.StorageQuota)
			assert.Equal(t, StorageQuotaModeWarning, o.QuotaConfiguration.Mode)
			assert.Equal(t, int64(1049600), *o.QuotaConfiguration.CurrentUsage)
		},
	},
}

func TestMockGetBucketStorageQuota_Success(t *testing.T) {
	for _, c := range testMockGetBucketStorageQuotaSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetBucketStorageQuota(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketStorageQuotaSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketStorageQuotaRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketStorageQuotaResult, err error)
}{
	{
		http.StatusNoContent,
		map[string]string{"x-oss-request-id": "534B371674E88A4D8906****", "Date": "Fri, 24 Feb 2017 03:15:40 GMT"},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/bucket/?quota", sortQuery(r))
		},
		&DeleteBucketStorageQuotaRequest{Bucket: Ptr("bucket")},
		func(t *testing.T, o *DeleteBucketStorageQuotaResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, http.StatusNoContent, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
		},
	},
}

func TestMockDeleteBucketStorageQuota_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketStorageQuotaSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteBucketStorageQuota(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketStorageQuotaErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketStorageQuotaRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketStorageQuotaResult, err error)
}{
	{
		http.StatusNotFound,
		map[string]string{"Content-Type": "application/xml", "x-oss-request-id": "5C3D9175B6FC201293AD****"},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?><Error><Code>NoSuchBucket</Code><Message>The specified bucket does not exist.</Message><RequestId>5C3D9175B6FC201293AD****</RequestId><EC>0015-00000101</EC></Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, http.MethodPut, r.Method)
			assert.Equal(t, "/bucket/?quota", sortQuery(r))
		},
		&PutBucketStorageQuotaRequest{Bucket: Ptr("bucket"), QuotaConfiguration: &QuotaConfiguration{StorageQuota: Ptr(int64(10737418240)), Mode: StorageQuotaModeStrict}},
		func(t *testing.T, o *PutBucketStorageQuotaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serviceErr *ServiceError
			assert.True(t, errors.As(err, &serviceErr))
			assert.Equal(t, http.StatusNotFound, serviceErr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serviceErr.Code)
			assert.Equal(t, "0015-00000101", serviceErr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serviceErr.RequestID)
		},
	},
}

func TestMockPutBucketStorageQuota_Error(t *testing.T) {
	for _, c := range testMockPutBucketStorageQuotaErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketStorageQuota(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketStorageQuotaErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketStorageQuotaRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketStorageQuotaResult, err error)
}{
	{
		http.StatusNotFound,
		map[string]string{"Content-Type": "application/xml", "x-oss-request-id": "5C3D9175B6FC201293AD****"},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?><Error><Code>NoSuchBucket</Code><Message>The specified bucket does not exist.</Message><RequestId>5C3D9175B6FC201293AD****</RequestId><EC>0015-00000101</EC></Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/bucket/?quota", sortQuery(r))
		},
		&GetBucketStorageQuotaRequest{Bucket: Ptr("bucket")},
		func(t *testing.T, o *GetBucketStorageQuotaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serviceErr *ServiceError
			assert.True(t, errors.As(err, &serviceErr))
			assert.Equal(t, http.StatusNotFound, serviceErr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serviceErr.Code)
			assert.Equal(t, "0015-00000101", serviceErr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serviceErr.RequestID)
		},
	},
}

func TestMockGetBucketStorageQuota_Error(t *testing.T) {
	for _, c := range testMockGetBucketStorageQuotaErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetBucketStorageQuota(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketStorageQuotaErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketStorageQuotaRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketStorageQuotaResult, err error)
}{
	{
		http.StatusNotFound,
		map[string]string{"Content-Type": "application/xml", "x-oss-request-id": "5C3D9175B6FC201293AD****"},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?><Error><Code>NoSuchBucket</Code><Message>The specified bucket does not exist.</Message><RequestId>5C3D9175B6FC201293AD****</RequestId><EC>0015-00000101</EC></Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/bucket/?quota", sortQuery(r))
		},
		&DeleteBucketStorageQuotaRequest{Bucket: Ptr("bucket")},
		func(t *testing.T, o *DeleteBucketStorageQuotaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serviceErr *ServiceError
			assert.True(t, errors.As(err, &serviceErr))
			assert.Equal(t, http.StatusNotFound, serviceErr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serviceErr.Code)
			assert.Equal(t, "0015-00000101", serviceErr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serviceErr.RequestID)
		},
	},
}

func TestMockDeleteBucketStorageQuota_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketStorageQuotaErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteBucketStorageQuota(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
