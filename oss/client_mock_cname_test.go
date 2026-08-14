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

var testMockPutCnameSuccessCasesLegacy = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutCnameRequest
	CheckOutputFn  func(t *testing.T, o *PutCnameResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=add", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&PutCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *PutCnameResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=add", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain><CertificateConfiguration><CertId>493****-cn-hangzhou</CertId><Certificate>-----BEGIN CERTIFICATE----- MIIDhDCCAmwCCQCFs8ixARsyrDANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC **** -----END CERTIFICATE-----</Certificate><PrivateKey>-----BEGIN CERTIFICATE----- MIIDhDCCAmwCCQCFs8ixARsyrDANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC **** -----END CERTIFICATE-----</PrivateKey><PreviousCertId>493****-cn-hangzhou</PreviousCertId><Force>true</Force></CertificateConfiguration></Cname></BucketCnameConfiguration>")
		},
		&PutCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
				CertificateConfiguration: &CertificateConfiguration{
					CertId:         Ptr("493****-cn-hangzhou"),
					Certificate:    Ptr("-----BEGIN CERTIFICATE----- MIIDhDCCAmwCCQCFs8ixARsyrDANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC **** -----END CERTIFICATE-----"),
					PrivateKey:     Ptr("-----BEGIN CERTIFICATE----- MIIDhDCCAmwCCQCFs8ixARsyrDANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC **** -----END CERTIFICATE-----"),
					PreviousCertId: Ptr("493****-cn-hangzhou"),
					Force:          Ptr(true),
				},
			},
		},
		func(t *testing.T, o *PutCnameResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=add", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain><CertificateConfiguration><DeleteCertificate>true</DeleteCertificate></CertificateConfiguration></Cname></BucketCnameConfiguration>")
		},
		&PutCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
				CertificateConfiguration: &CertificateConfiguration{
					DeleteCertificate: Ptr(true),
				},
			},
		},
		func(t *testing.T, o *PutCnameResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutCnameLegacy_Success(t *testing.T) {
	for _, c := range testMockPutCnameSuccessCasesLegacy {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutCname(context.TODO(), c.Request)
		assert.Nil(t, c.Request.BucketCnameConfiguration.Cname)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutCnameSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutCnameRequest
	CheckOutputFn  func(t *testing.T, o *PutCnameResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=add", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain><IsWildCard>true</IsWildCard></Cname></BucketCnameConfiguration>")
		},
		&PutCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Cname: &Cname{
					Domain:     Ptr("example.com"),
					IsWildCard: Ptr(true),
				},
			},
		},
		func(t *testing.T, o *PutCnameResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=add", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain><CertificateConfiguration><CertId>493****-cn-hangzhou</CertId><Certificate>-----BEGIN CERTIFICATE----- MIIDhDCCAmwCCQCFs8ixARsyrDANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC **** -----END CERTIFICATE-----</Certificate><PrivateKey>-----BEGIN CERTIFICATE----- MIIDhDCCAmwCCQCFs8ixARsyrDANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC **** -----END CERTIFICATE-----</PrivateKey><PreviousCertId>493****-cn-hangzhou</PreviousCertId><Force>true</Force></CertificateConfiguration></Cname></BucketCnameConfiguration>")
		},
		&PutCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Cname: &Cname{
					Domain: Ptr("example.com"),
					CertificateConfiguration: &CertificateConfiguration{
						CertId:         Ptr("493****-cn-hangzhou"),
						Certificate:    Ptr("-----BEGIN CERTIFICATE----- MIIDhDCCAmwCCQCFs8ixARsyrDANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC **** -----END CERTIFICATE-----"),
						PrivateKey:     Ptr("-----BEGIN CERTIFICATE----- MIIDhDCCAmwCCQCFs8ixARsyrDANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMC **** -----END CERTIFICATE-----"),
						PreviousCertId: Ptr("493****-cn-hangzhou"),
						Force:          Ptr(true),
					},
				},
			},
		},
		func(t *testing.T, o *PutCnameResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=add", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain><CertificateConfiguration><DeleteCertificate>true</DeleteCertificate></CertificateConfiguration></Cname></BucketCnameConfiguration>")
		},
		&PutCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Cname: &Cname{
					Domain: Ptr("example.com"),
					CertificateConfiguration: &CertificateConfiguration{
						DeleteCertificate: Ptr(true),
					},
				},
			},
		},
		func(t *testing.T, o *PutCnameResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutCname_Success(t *testing.T) {
	for _, c := range testMockPutCnameSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutCname(context.TODO(), c.Request)
		assert.NotNil(t, c.Request.BucketCnameConfiguration.Cname)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutCnameErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutCnameRequest
	CheckOutputFn  func(t *testing.T, o *PutCnameResult, err error)
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=add", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&PutCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *PutCnameResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=add", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&PutCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *PutCnameResult, err error) {
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

func TestMockPutCname_Error(t *testing.T) {
	for _, c := range testMockPutCnameErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutCname(context.TODO(), c.Request)
		assert.Nil(t, c.Request.BucketCnameConfiguration.Cname)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateCnameTokenLegacySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateCnameTokenRequest
	CheckOutputFn  func(t *testing.T, o *CreateCnameTokenResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<CnameToken>
  <Bucket>bucket</Bucket>
  <Cname>example.com</Cname>;
  <Token>be1d49d863dea9ffeff3df7d6455****</Token>
  <ExpireTime>Wed, 23 Feb 2022 21:16:37 GMT</ExpireTime>
</CnameToken>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=token", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&CreateCnameTokenRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *CreateCnameTokenResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.CnameToken.Bucket, "bucket")
			assert.Equal(t, *o.CnameToken.Cname, "example.com")
			assert.Equal(t, *o.CnameToken.Token, "be1d49d863dea9ffeff3df7d6455****")
			assert.Equal(t, *o.CnameToken.ExpireTime, "Wed, 23 Feb 2022 21:16:37 GMT")
		},
	},
}

func TestMockCreateCnameTokenLegacy_Success(t *testing.T) {
	for _, c := range testMockCreateCnameTokenLegacySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CreateCnameToken(context.TODO(), c.Request)
		assert.Nil(t, c.Request.BucketCnameConfiguration.Cname)
		c.CheckOutputFn(t, output, err)
	}
}

var createRequest = func() *CreateCnameTokenRequest {
	req := &CreateCnameTokenRequest{
		Bucket: Ptr("bucket"),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Cname: &Cname{
				Domain: Ptr("example.com"),
			},
		},
	}
	req.AddParameter("wildcard", "false")
	return req
}()

var testMockCreateCnameTokenSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateCnameTokenRequest
	CheckOutputFn  func(t *testing.T, o *CreateCnameTokenResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<CnameToken>
  <Bucket>bucket</Bucket>
  <Cname>example.com</Cname>;
  <Token>be1d49d863dea9ffeff3df7d6455****</Token>
  <ExpireTime>Wed, 23 Feb 2022 21:16:37 GMT</ExpireTime>
</CnameToken>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=token&wildcard=false", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		createRequest,
		func(t *testing.T, o *CreateCnameTokenResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.CnameToken.Bucket, "bucket")
			assert.Equal(t, *o.CnameToken.Cname, "example.com")
			assert.Equal(t, *o.CnameToken.Token, "be1d49d863dea9ffeff3df7d6455****")
			assert.Equal(t, *o.CnameToken.ExpireTime, "Wed, 23 Feb 2022 21:16:37 GMT")
		},
	},
}

func TestMockCreateCnameToken_Success(t *testing.T) {
	for _, c := range testMockCreateCnameTokenSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CreateCnameToken(context.TODO(), c.Request)
		assert.NotNil(t, c.Request.BucketCnameConfiguration.Cname)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateCnameTokenErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateCnameTokenRequest
	CheckOutputFn  func(t *testing.T, o *CreateCnameTokenResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<Error>
  <Code>TooManyCnameToken</Code>
  <Message>You have attempted to create more cname token than allowed.</Message>
  <RequestId>6215FD21DA0E27393F004E9E</RequestId>
  <HostId>127.0.0.1</HostId>
  <Bucket>bucket</Bucket>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=token", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&CreateCnameTokenRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *CreateCnameTokenResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "TooManyCnameToken", serr.Code)
			assert.Equal(t, "You have attempted to create more cname token than allowed.", serr.Message)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=token", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&CreateCnameTokenRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *CreateCnameTokenResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=token", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&CreateCnameTokenRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *CreateCnameTokenResult, err error) {
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=token", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&CreateCnameTokenRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *CreateCnameTokenResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute CreateCnameToken fail")
		},
	},
}

func TestMockCreateCnameToken_Error(t *testing.T) {
	for _, c := range testMockCreateCnameTokenErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CreateCnameToken(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetCnameTokenSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetCnameTokenRequest
	CheckOutputFn  func(t *testing.T, o *GetCnameTokenResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CnameToken>
  <Bucket>bucket</Bucket>
  <Cname>example.com</Cname>
  <Token>be1d49d863dea9ffeff3df7d6455****</Token>
  <ExpireTime>Wed, 23 Feb 2022 21:39:42 GMT</ExpireTime>
</CnameToken>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname=example.com&comp=token&wildcard=false", urlStr)
		},
		&GetCnameTokenRequest{
			Bucket:   Ptr("bucket"),
			Cname:    Ptr("example.com"),
			Wildcard: Ptr(false),
		},
		func(t *testing.T, o *GetCnameTokenResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.CnameToken.Bucket, "bucket")
			assert.Equal(t, *o.CnameToken.Cname, "example.com")
			assert.Equal(t, *o.CnameToken.Token, "be1d49d863dea9ffeff3df7d6455****")
			assert.Equal(t, *o.CnameToken.ExpireTime, "Wed, 23 Feb 2022 21:39:42 GMT")
		},
	},
}

func TestMockGetCnameToken_Success(t *testing.T) {
	for _, c := range testMockGetCnameTokenSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetCnameToken(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetCnameTokenErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetCnameTokenRequest
	CheckOutputFn  func(t *testing.T, o *GetCnameTokenResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname=example.com&comp=token", urlStr)
		},
		&GetCnameTokenRequest{
			Bucket: Ptr("bucket"),
			Cname:  Ptr("example.com"),
		},
		func(t *testing.T, o *GetCnameTokenResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname=example.com&comp=token", urlStr)
		},
		&GetCnameTokenRequest{
			Bucket: Ptr("bucket"),
			Cname:  Ptr("example.com"),
		},
		func(t *testing.T, o *GetCnameTokenResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname=example.com&comp=token", urlStr)
		},
		&GetCnameTokenRequest{
			Bucket: Ptr("bucket"),
			Cname:  Ptr("example.com"),
		},
		func(t *testing.T, o *GetCnameTokenResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetCnameToken fail")
		},
	},
}

func TestMockGetCnameToken_Error(t *testing.T) {
	for _, c := range testMockGetCnameTokenErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetCnameToken(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListCnameSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListCnameRequest
	CheckOutputFn  func(t *testing.T, o *ListCnameResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<ListCnameResult>
  <Bucket>bucket</Bucket>
  <Owner>owner</Owner>
  <Cname>
    <Domain>example.com</Domain>
    <LastModified>2021-09-15T02:35:07.000Z</LastModified>
    <Status>Enabled</Status>
    <Certificate>
      <Type>CAS</Type>
      <CertId>493****-cn-hangzhou</CertId>
      <Status>Enabled</Status>
      <CreationDate>Wed, 15 Sep 2021 02:35:06 GMT</CreationDate>
      <Fingerprint>DE:01:CF:EC:7C:A7:98:CB:D8:6E:FB:1D:97:EB:A9:64:1D:4E:**:**</Fingerprint>
      <ValidStartDate>Wed, 12 Apr 2023 10:14:51 GMT</ValidStartDate>
      <ValidEndDate>Mon, 4 May 2048 10:14:51 GMT</ValidEndDate>
    </Certificate>
  </Cname>
  <Cname>
    <Domain>example.org</Domain>
    <LastModified>2021-09-15T02:34:58.000Z</LastModified>
    <Status>Enabled</Status>
  </Cname>
  <Cname>
    <Domain>example.edu</Domain>
    <LastModified>2021-09-15T02:50:34.000Z</LastModified>
    <Status>Enabled</Status>
  </Cname>
  <Cname>
	<Domain>example.net</Domain>
	<LastModified>2021-09-15T02:50:34.000Z</LastModified>
	<Status>Enabled</Status>
	<IsWildCard>true</IsWildCard>
  </Cname>
</ListCnameResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname", urlStr)
		},
		&ListCnameRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListCnameResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Bucket, "bucket")
			assert.Equal(t, *o.Owner, "owner")
			assert.Equal(t, len(o.Cnames), 4)
			assert.Equal(t, *o.Cnames[0].Domain, "example.com")
			assert.Equal(t, *o.Cnames[0].LastModified, "2021-09-15T02:35:07.000Z")
			assert.Equal(t, *o.Cnames[0].Status, "Enabled")
			assert.Equal(t, *o.Cnames[0].Certificate.Type, "CAS")
			assert.Equal(t, *o.Cnames[0].Certificate.CertId, "493****-cn-hangzhou")
			assert.Equal(t, *o.Cnames[0].Certificate.Status, "Enabled")
			assert.Equal(t, *o.Cnames[0].Certificate.CreationDate, "Wed, 15 Sep 2021 02:35:06 GMT")
			assert.Equal(t, *o.Cnames[0].Certificate.Fingerprint, "DE:01:CF:EC:7C:A7:98:CB:D8:6E:FB:1D:97:EB:A9:64:1D:4E:**:**")
			assert.Equal(t, *o.Cnames[0].Certificate.ValidStartDate, "Wed, 12 Apr 2023 10:14:51 GMT")
			assert.Equal(t, *o.Cnames[0].Certificate.ValidEndDate, "Mon, 4 May 2048 10:14:51 GMT")
			assert.Equal(t, *o.Cnames[1].Domain, "example.org")
			assert.Equal(t, *o.Cnames[1].LastModified, "2021-09-15T02:34:58.000Z")
			assert.Equal(t, *o.Cnames[1].Status, "Enabled")
			assert.Equal(t, *o.Cnames[2].Domain, "example.edu")
			assert.Equal(t, *o.Cnames[2].LastModified, "2021-09-15T02:50:34.000Z")
			assert.Equal(t, *o.Cnames[2].Status, "Enabled")
			assert.Equal(t, *o.Cnames[3].Domain, "example.net")
			assert.Equal(t, *o.Cnames[3].LastModified, "2021-09-15T02:50:34.000Z")
			assert.Equal(t, *o.Cnames[3].Status, "Enabled")
			assert.Equal(t, *o.Cnames[3].IsWildCard, true)
		},
	},
}

func TestMockListCname_Success(t *testing.T) {
	for _, c := range testMockListCnameSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListCname(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListCnameErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListCnameRequest
	CheckOutputFn  func(t *testing.T, o *ListCnameResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname", urlStr)
		},
		&ListCnameRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListCnameResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname", urlStr)
		},
		&ListCnameRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListCnameResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname", urlStr)
		},
		&ListCnameRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListCnameResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListCname fail")
		},
	},
}

func TestMockListCname_Error(t *testing.T) {
	for _, c := range testMockListCnameErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListCname(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteCnameLegacySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteCnameRequest
	CheckOutputFn  func(t *testing.T, o *DeleteCnameResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=delete", urlStr)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&DeleteCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *DeleteCnameResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteCnameLegacy_Success(t *testing.T) {
	for _, c := range testMockDeleteCnameLegacySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteCname(context.TODO(), c.Request)
		assert.Nil(t, c.Request.BucketCnameConfiguration.Cname)
		c.CheckOutputFn(t, output, err)
	}
}

var deleteRequest = func() *DeleteCnameRequest {
	req := &DeleteCnameRequest{
		Bucket: Ptr("bucket"),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Cname: &Cname{
				Domain: Ptr("example.com"),
			},
		},
	}
	req.AddParameter("wildcard", "false")
	return req
}()

var testMockDeleteCnameSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteCnameRequest
	CheckOutputFn  func(t *testing.T, o *DeleteCnameResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=delete&wildcard=false", urlStr)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		deleteRequest,
		func(t *testing.T, o *DeleteCnameResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteCname_Success(t *testing.T) {
	for _, c := range testMockDeleteCnameSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteCname(context.TODO(), c.Request)
		assert.NotNil(t, c.Request.BucketCnameConfiguration.Cname)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteCnameErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteCnameRequest
	CheckOutputFn  func(t *testing.T, o *DeleteCnameResult, err error)
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=delete", urlStr)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&DeleteCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *DeleteCnameResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cname&comp=delete", urlStr)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<BucketCnameConfiguration><Cname><Domain>example.com</Domain></Cname></BucketCnameConfiguration>")
		},
		&DeleteCnameRequest{
			Bucket: Ptr("bucket"),
			BucketCnameConfiguration: &BucketCnameConfiguration{
				Domain: Ptr("example.com"),
			},
		},
		func(t *testing.T, o *DeleteCnameResult, err error) {
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

func TestMockDeleteCname_Error(t *testing.T) {
	for _, c := range testMockDeleteCnameErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteCname(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
