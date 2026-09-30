package oss

import (
	"context"
	"errors"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

var testMockCreateBucketMetadataTableConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateBucketMetadataTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *CreateBucketMetadataTableConfigurationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, contentTypeXML, r.Header.Get(HTTPHeaderContentType))
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataConfiguration", strUrl)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "<MetadataConfiguration><JournalTableConfiguration><RecordExpiration><Expiration>ENABLED</Expiration><Days>7</Days></RecordExpiration></JournalTableConfiguration><InventoryTableConfiguration><ConfigurationState>ENABLED</ConfigurationState><EncryptionConfiguration><SseAlgorithm>AES256</SseAlgorithm></EncryptionConfiguration></InventoryTableConfiguration></MetadataConfiguration>", string(requestBody))
		},
		&CreateBucketMetadataTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
			MetadataConfiguration: &MetadataConfiguration{
				JournalTableConfiguration: &JournalTableConfiguration{
					RecordExpiration: &RecordExpiration{
						Expiration: Ptr("ENABLED"),
						Days:       Ptr(7),
					},
				},
				InventoryTableConfiguration: &InventoryTableConfiguration{
					ConfigurationState: Ptr("ENABLED"),
					EncryptionConfiguration: &MetadataTableEncryptionConfiguration{
						SseAlgorithm: Ptr("AES256"),
					},
				},
			},
		},
		func(t *testing.T, o *CreateBucketMetadataTableConfigurationResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
		},
	},
}

func TestMockCreateBucketMetadataTableConfiguration_Success(t *testing.T) {
	for _, c := range testMockCreateBucketMetadataTableConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateBucketMetadataTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateBucketMetadataTableConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateBucketMetadataTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *CreateBucketMetadataTableConfigurationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataConfiguration", strUrl)
		},
		&CreateBucketMetadataTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
			MetadataConfiguration: &MetadataConfiguration{
				InventoryTableConfiguration: &InventoryTableConfiguration{
					ConfigurationState: Ptr("ENABLED"),
				},
			},
		},
		func(t *testing.T, o *CreateBucketMetadataTableConfigurationResult, err error) {
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
}

func TestMockCreateBucketMetadataTableConfiguration_Error(t *testing.T) {
	for _, c := range testMockCreateBucketMetadataTableConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateBucketMetadataTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketMetadataTableConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketMetadataTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketMetadataTableConfigurationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<GetBucketMetadataConfigurationResult>
  <MetadataConfigurationResult>
    <DestinationResult>
      <TableBucketType>oss</TableBucketType>
      <TableBucketArn>arn:acs:oss:cn-hangzhou:123456789012:bucket/examplebucket</TableBucketArn>
      <TableNamespace>b_examplebucket</TableNamespace>
    </DestinationResult>
    <JournalTableConfigurationResult>
      <TableStatus>ACTIVE</TableStatus>
      <TableName>journal</TableName>
      <TableArn>arn:acs:oss:cn-hangzhou:123456789012:bucket/examplebucket/journal</TableArn>
      <RecordExpiration>
        <Expiration>ENABLED</Expiration>
        <Days>7</Days>
      </RecordExpiration>
      <EncryptionConfiguration>
        <SseAlgorithm>AES256</SseAlgorithm>
      </EncryptionConfiguration>
    </JournalTableConfigurationResult>
    <InventoryTableConfigurationResult>
      <ConfigurationState>ENABLED</ConfigurationState>
      <TableStatus>ACTIVE</TableStatus>
      <TableName>inventory</TableName>
      <TableArn>arn:acs:oss:cn-hangzhou:123456789012:bucket/examplebucket/inventory</TableArn>
      <EncryptionConfiguration>
        <SseAlgorithm>AES256</SseAlgorithm>
      </EncryptionConfiguration>
    </InventoryTableConfigurationResult>
  </MetadataConfigurationResult>
</GetBucketMetadataConfigurationResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, contentTypeXML, r.Header.Get(HTTPHeaderContentType))
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataConfiguration", strUrl)
		},
		&GetBucketMetadataTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
		},
		func(t *testing.T, o *GetBucketMetadataTableConfigurationResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.NotNil(t, o.GetBucketMetadataConfigurationResult)
			assert.NotNil(t, o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult)
			assert.Equal(t, "oss", *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.DestinationResult.TableBucketType)
			assert.Equal(t, "b_examplebucket", *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.DestinationResult.TableNamespace)
			assert.Equal(t, "ACTIVE", *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.TableStatus)
			assert.Equal(t, "journal", *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.TableName)
			assert.Equal(t, "ENABLED", *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.RecordExpiration.Expiration)
			assert.Equal(t, int(7), *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.RecordExpiration.Days)
			assert.Equal(t, "AES256", *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.EncryptionConfiguration.SseAlgorithm)
			assert.Equal(t, "ENABLED", *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.InventoryTableConfigurationResult.ConfigurationState)
			assert.Equal(t, "ACTIVE", *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.InventoryTableConfigurationResult.TableStatus)
			assert.Equal(t, "inventory", *o.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.InventoryTableConfigurationResult.TableName)
		},
	},
}

func TestMockGetBucketMetadataTableConfiguration_Success(t *testing.T) {
	for _, c := range testMockGetBucketMetadataTableConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketMetadataTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketMetadataTableConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketMetadataTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketMetadataTableConfigurationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataConfiguration", strUrl)
		},
		&GetBucketMetadataTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
		},
		func(t *testing.T, o *GetBucketMetadataTableConfigurationResult, err error) {
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
}

func TestMockGetBucketMetadataTableConfiguration_Error(t *testing.T) {
	for _, c := range testMockGetBucketMetadataTableConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketMetadataTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketMetadataTableConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketMetadataTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketMetadataTableConfigurationResult, err error)
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
			assert.Equal(t, contentTypeXML, r.Header.Get(HTTPHeaderContentType))
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataConfiguration", strUrl)
		},
		&DeleteBucketMetadataTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
		},
		func(t *testing.T, o *DeleteBucketMetadataTableConfigurationResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
		},
	},
}

func TestMockDeleteBucketMetadataTableConfiguration_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketMetadataTableConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketMetadataTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketMetadataTableConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketMetadataTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketMetadataTableConfigurationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MetadataConfigurationNotFound</Code>
  <Message>The metadata table configuration does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataConfiguration", strUrl)
		},
		&DeleteBucketMetadataTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
		},
		func(t *testing.T, o *DeleteBucketMetadataTableConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "MetadataConfigurationNotFound", serr.Code)
			assert.Equal(t, "The metadata table configuration does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockDeleteBucketMetadataTableConfiguration_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketMetadataTableConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketMetadataTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateBucketMetadataInventoryTableConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateBucketMetadataInventoryTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *UpdateBucketMetadataInventoryTableConfigurationResult, err error)
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
			assert.Equal(t, contentTypeXML, r.Header.Get(HTTPHeaderContentType))
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataInventoryTable", strUrl)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "<InventoryTableConfiguration><ConfigurationState>DISABLED</ConfigurationState></InventoryTableConfiguration>", string(requestBody))
		},
		&UpdateBucketMetadataInventoryTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
			InventoryTableConfiguration: &InventoryTableConfiguration{
				ConfigurationState: Ptr("DISABLED"),
			},
		},
		func(t *testing.T, o *UpdateBucketMetadataInventoryTableConfigurationResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
		},
	},
}

func TestMockUpdateBucketMetadataInventoryTableConfiguration_Success(t *testing.T) {
	for _, c := range testMockUpdateBucketMetadataInventoryTableConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateBucketMetadataInventoryTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateBucketMetadataInventoryTableConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateBucketMetadataInventoryTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *UpdateBucketMetadataInventoryTableConfigurationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MetadataConfigurationNotFound</Code>
  <Message>The metadata table configuration does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataInventoryTable", strUrl)
		},
		&UpdateBucketMetadataInventoryTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
			InventoryTableConfiguration: &InventoryTableConfiguration{
				ConfigurationState: Ptr("DISABLED"),
			},
		},
		func(t *testing.T, o *UpdateBucketMetadataInventoryTableConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "MetadataConfigurationNotFound", serr.Code)
			assert.Equal(t, "The metadata table configuration does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockUpdateBucketMetadataInventoryTableConfiguration_Error(t *testing.T) {
	for _, c := range testMockUpdateBucketMetadataInventoryTableConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateBucketMetadataInventoryTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateBucketMetadataJournalTableConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateBucketMetadataJournalTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *UpdateBucketMetadataJournalTableConfigurationResult, err error)
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
			assert.Equal(t, contentTypeXML, r.Header.Get(HTTPHeaderContentType))
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataJournalTable", strUrl)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "<JournalTableConfiguration><RecordExpiration><Expiration>DISABLED</Expiration></RecordExpiration></JournalTableConfiguration>", string(requestBody))
		},
		&UpdateBucketMetadataJournalTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
			JournalTableConfiguration: &JournalTableConfiguration{
				RecordExpiration: &RecordExpiration{
					Expiration: Ptr("DISABLED"),
				},
			},
		},
		func(t *testing.T, o *UpdateBucketMetadataJournalTableConfigurationResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
		},
	},
}

func TestMockUpdateBucketMetadataJournalTableConfiguration_Success(t *testing.T) {
	for _, c := range testMockUpdateBucketMetadataJournalTableConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateBucketMetadataJournalTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateBucketMetadataJournalTableConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateBucketMetadataJournalTableConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *UpdateBucketMetadataJournalTableConfigurationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0015-00000101",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MetadataConfigurationNotFound</Code>
  <Message>The metadata table configuration does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/oss-demo/?metadataJournalTable", strUrl)
		},
		&UpdateBucketMetadataJournalTableConfigurationRequest{
			Bucket: Ptr("oss-demo"),
			JournalTableConfiguration: &JournalTableConfiguration{
				RecordExpiration: &RecordExpiration{
					Expiration: Ptr("DISABLED"),
				},
			},
		},
		func(t *testing.T, o *UpdateBucketMetadataJournalTableConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "MetadataConfigurationNotFound", serr.Code)
			assert.Equal(t, "The metadata table configuration does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockUpdateBucketMetadataJournalTableConfiguration_Error(t *testing.T) {
	for _, c := range testMockUpdateBucketMetadataJournalTableConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateBucketMetadataJournalTableConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
