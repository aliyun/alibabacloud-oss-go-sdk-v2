package oss

import (
	"testing"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutBucketLifecycleSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketLifecycleRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketLifecycleResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>r0</ID><Prefix>prefix0</Prefix><Expiration><Days>40</Days><ExpiredObjectDeleteMarker>false</ExpiredObjectDeleteMarker></Expiration></Rule><Rule><Status>Enabled</Status><Filter><ObjectSizeGreaterThan>500</ObjectSizeGreaterThan><ObjectSizeLessThan>64500</ObjectSizeLessThan></Filter><ID>r1</ID><Prefix>prefix1</Prefix><Expiration><Days>40</Days><ExpiredObjectDeleteMarker>false</ExpiredObjectDeleteMarker></Expiration></Rule><Rule><Status>Enabled</Status><Filter><ObjectSizeGreaterThan>500</ObjectSizeGreaterThan><ObjectSizeLessThan>64500</ObjectSizeLessThan></Filter><ID>r3</ID><Prefix>prefix3</Prefix><Expiration><Days>40</Days><ExpiredObjectDeleteMarker>false</ExpiredObjectDeleteMarker></Expiration><Transition><Days>30</Days><StorageClass>IA</StorageClass><IsAccessTime>false</IsAccessTime></Transition></Rule><Rule><Status>Enabled</Status><AbortMultipartUpload><CreatedBeforeDate>2015-11-11T00:00:00.000Z</CreatedBeforeDate></AbortMultipartUpload><NoncurrentVersionTransition><IsAccessTime>true</IsAccessTime><ReturnToStdWhenVisit>true</ReturnToStdWhenVisit><NoncurrentDays>10</NoncurrentDays><StorageClass>IA</StorageClass></NoncurrentVersionTransition><ID>r4</ID><Prefix>prefix4</Prefix><Expiration><ExpiredObjectDeleteMarker>true</ExpiredObjectDeleteMarker></Expiration></Rule><Rule><Status>Enabled</Status><Prefix>pre_</Prefix><Expiration><ExpiredObjectDeleteMarker>true</ExpiredObjectDeleteMarker></Expiration></Rule></LifecycleConfiguration>")
		},
		&PutBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
			LifecycleConfiguration: &LifecycleConfiguration{
				Rules: []LifecycleRule{
					{
						ID:     Ptr("r0"),
						Prefix: Ptr("prefix0"),
						Status: Ptr("Enabled"),
						Expiration: &LifecycleRuleExpiration{
							Days:                      Ptr(int32(40)),
							ExpiredObjectDeleteMarker: Ptr(false),
						},
					},
					{
						ID:     Ptr("r1"),
						Prefix: Ptr("prefix1"),
						Status: Ptr("Enabled"),
						Expiration: &LifecycleRuleExpiration{
							Days:                      Ptr(int32(40)),
							ExpiredObjectDeleteMarker: Ptr(false),
						},
						Filter: &LifecycleRuleFilter{
							ObjectSizeGreaterThan: Ptr(int64(500)),
							ObjectSizeLessThan:    Ptr(int64(64500)),
						},
					},
					{
						ID:     Ptr("r3"),
						Prefix: Ptr("prefix3"),
						Status: Ptr("Enabled"),
						Expiration: &LifecycleRuleExpiration{
							Days:                      Ptr(int32(40)),
							ExpiredObjectDeleteMarker: Ptr(false),
						},
						Transitions: []LifecycleRuleTransition{
							{
								Days:         Ptr(int32(30)),
								StorageClass: StorageClassIA,
								IsAccessTime: Ptr(false),
							},
						},
						Filter: &LifecycleRuleFilter{
							ObjectSizeGreaterThan: Ptr(int64(500)),
							ObjectSizeLessThan:    Ptr(int64(64500)),
						},
					},
					{
						ID:     Ptr("r4"),
						Prefix: Ptr("prefix4"),
						Status: Ptr("Enabled"),
						Expiration: &LifecycleRuleExpiration{
							ExpiredObjectDeleteMarker: Ptr(true),
						},
						AbortMultipartUpload: &LifecycleRuleAbortMultipartUpload{
							CreatedBeforeDate: Ptr("2015-11-11T00:00:00.000Z"),
						},
						NoncurrentVersionTransitions: []NoncurrentVersionTransition{
							{
								NoncurrentDays:       Ptr(int32(10)),
								StorageClass:         StorageClassIA,
								IsAccessTime:         Ptr(true),
								ReturnToStdWhenVisit: Ptr(true),
							},
						},
					},
					{
						Prefix: Ptr("pre_"),
						Status: Ptr("Enabled"),
						Expiration: &LifecycleRuleExpiration{
							ExpiredObjectDeleteMarker: Ptr(true),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketLifecycleResult, err error) {
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
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule</ID><Prefix>log/</Prefix><Transition><Days>30</Days><StorageClass>IA</StorageClass></Transition></Rule></LifecycleConfiguration>")
		},
		&PutBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
			LifecycleConfiguration: &LifecycleConfiguration{
				Rules: []LifecycleRule{
					{
						Status: Ptr("Enabled"),
						ID:     Ptr("rule"),
						Prefix: Ptr("log/"),
						Transitions: []LifecycleRuleTransition{
							{
								Days:         Ptr(int32(30)),
								StorageClass: StorageClassIA,
							},
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketLifecycleResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketLifecycle_Success(t *testing.T) {
	for _, c := range testMockPutBucketLifecycleSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketLifecycle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketLifecycleErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketLifecycleRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketLifecycleResult, err error)
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
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule</ID><Prefix>log/</Prefix><Transition><Days>30</Days><StorageClass>IA</StorageClass></Transition></Rule></LifecycleConfiguration>")
		},
		&PutBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
			LifecycleConfiguration: &LifecycleConfiguration{
				Rules: []LifecycleRule{
					{
						Status: Ptr("Enabled"),
						ID:     Ptr("rule"),
						Prefix: Ptr("log/"),
						Transitions: []LifecycleRuleTransition{
							{
								Days:         Ptr(int32(30)),
								StorageClass: StorageClassIA,
							},
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketLifecycleResult, err error) {
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
			assert.Equal(t, "/bucket/?lifecycle", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule</ID><Prefix>log/</Prefix><Transition><Days>30</Days><StorageClass>IA</StorageClass></Transition></Rule></LifecycleConfiguration>")
		},
		&PutBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
			LifecycleConfiguration: &LifecycleConfiguration{
				Rules: []LifecycleRule{
					{
						Status: Ptr("Enabled"),
						ID:     Ptr("rule"),
						Prefix: Ptr("log/"),
						Transitions: []LifecycleRuleTransition{
							{
								Days:         Ptr(int32(30)),
								StorageClass: StorageClassIA,
							},
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketLifecycleResult, err error) {
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
			assert.Equal(t, "/bucket/?lifecycle", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule</ID><Prefix>log/</Prefix><Transition><Days>30</Days><StorageClass>IA</StorageClass></Transition></Rule></LifecycleConfiguration>")
		},
		&PutBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
			LifecycleConfiguration: &LifecycleConfiguration{
				Rules: []LifecycleRule{
					{
						Status: Ptr("Enabled"),
						ID:     Ptr("rule"),
						Prefix: Ptr("log/"),
						Transitions: []LifecycleRuleTransition{
							{
								Days:         Ptr(int32(30)),
								StorageClass: StorageClassIA,
							},
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketLifecycleResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutBucketLifecycle fail")
		},
	},
}

func TestMockPutBucketLifecycle_Error(t *testing.T) {
	for _, c := range testMockPutBucketLifecycleErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketLifecycle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

const lifecycleRespWithUnknown = `<?xml version="1.0" encoding="UTF-8"?>
<LifecycleConfiguration>
  <Rule>
    <ID>r0</ID>
    <Prefix>a/</Prefix>
    <Status>Enabled</Status>
    <Expiration>
      <Days>1</Days>
    </Expiration>
    <NewExtension>v1</NewExtension>
  </Rule>
  <Rule>
    <ID>r1</ID>
    <Prefix>b/</Prefix>
    <Status>Enabled</Status>
    <Filter>
      <Not>
        <Prefix>skip/</Prefix>
        <NewFilterExt>x</NewFilterExt>
      </Not>
    </Filter>
  </Rule>
</LifecycleConfiguration>`

const lifecycleRespClean = `<?xml version="1.0" encoding="UTF-8"?>
<LifecycleConfiguration>
  <Rule>
    <ID>r0</ID>
    <Prefix>a/</Prefix>
    <Status>Enabled</Status>
    <Expiration>
      <Days>1</Days>
    </Expiration>
  </Rule>
</LifecycleConfiguration>`

var testMockGetBucketLifecycleSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLifecycleRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLifecycleResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<LifecycleConfiguration>
  <Rule>
    <ID>delete after one day</ID>
    <Prefix>logs1/</Prefix>
    <Status>Enabled</Status>
    <Expiration>
      <Days>1</Days>
    </Expiration>
  </Rule>
  <Rule>
    <ID>mtime transition1</ID>
    <Prefix>logs2/</Prefix>
    <Status>Enabled</Status>
    <Transition>
      <Days>30</Days>
      <StorageClass>IA</StorageClass>
    </Transition>
  </Rule>
  <Rule>
    <ID>mtime transition2</ID>
    <Prefix>logs3/</Prefix>
    <Status>Enabled</Status>
    <Transition>
      <Days>30</Days>
      <StorageClass>IA</StorageClass>
      <IsAccessTime>false</IsAccessTime>
    </Transition>
  </Rule>
</LifecycleConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
		},
		&GetBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			config := *o.LifecycleConfiguration
			assert.Equal(t, 3, len(config.Rules))
			assert.Equal(t, "delete after one day", *config.Rules[0].ID)
			assert.Equal(t, "logs1/", *config.Rules[0].Prefix)
			assert.Equal(t, "Enabled", *config.Rules[0].Status)
			assert.Equal(t, int32(1), *config.Rules[0].Expiration.Days)

			assert.Equal(t, "mtime transition1", *config.Rules[1].ID)
			assert.Equal(t, "logs2/", *config.Rules[1].Prefix)
			assert.Equal(t, "Enabled", *config.Rules[1].Status)
			assert.Equal(t, 1, len(config.Rules[1].Transitions))
			assert.Equal(t, int32(30), *config.Rules[1].Transitions[0].Days)
			assert.Equal(t, StorageClassIA, config.Rules[1].Transitions[0].StorageClass)

			assert.Equal(t, "mtime transition2", *config.Rules[2].ID)
			assert.Equal(t, "logs3/", *config.Rules[2].Prefix)
			assert.Equal(t, "Enabled", *config.Rules[2].Status)
			assert.Equal(t, 1, len(config.Rules[2].Transitions))
			assert.Equal(t, int32(30), *config.Rules[2].Transitions[0].Days)
			assert.Equal(t, StorageClassIA, config.Rules[2].Transitions[0].StorageClass)
			assert.False(t, *config.Rules[2].Transitions[0].IsAccessTime)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<LifecycleConfiguration>
  <Rule>
    <ID>atime transition1</ID>
    <Prefix>logs1/</Prefix>
    <Status>Enabled</Status>
    <Transition>
      <Days>30</Days>
      <StorageClass>IA</StorageClass>
      <IsAccessTime>true</IsAccessTime>
      <ReturnToStdWhenVisit>false</ReturnToStdWhenVisit>
    </Transition>
    <AtimeBase>1631698332</AtimeBase>
  </Rule>
  <Rule>
    <ID>atime transition2</ID>
    <Prefix>logs2/</Prefix>
    <Status>Enabled</Status>
    <NoncurrentVersionTransition>
      <NoncurrentDays>10</NoncurrentDays>
      <StorageClass>IA</StorageClass>
      <IsAccessTime>true</IsAccessTime>
      <ReturnToStdWhenVisit>false</ReturnToStdWhenVisit>
    </NoncurrentVersionTransition>
    <AtimeBase>1631698332</AtimeBase>
  </Rule>
</LifecycleConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
		},
		&GetBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			config := o.LifecycleConfiguration
			assert.Equal(t, 2, len(config.Rules))
			assert.Equal(t, "atime transition1", *config.Rules[0].ID)
			assert.Equal(t, "logs1/", *config.Rules[0].Prefix)
			assert.Equal(t, "Enabled", *config.Rules[0].Status)
			assert.Equal(t, 1, len(config.Rules[0].Transitions))
			assert.Equal(t, int32(30), *config.Rules[0].Transitions[0].Days)
			assert.Equal(t, StorageClassIA, config.Rules[0].Transitions[0].StorageClass)
			assert.False(t, *config.Rules[0].Transitions[0].ReturnToStdWhenVisit)
			assert.True(t, *config.Rules[0].Transitions[0].IsAccessTime)
			assert.Equal(t, int64(1631698332), *config.Rules[0].AtimeBase)

			assert.Equal(t, "atime transition2", *config.Rules[1].ID)
			assert.Equal(t, "logs2/", *config.Rules[1].Prefix)
			assert.Equal(t, "Enabled", *config.Rules[1].Status)
			assert.Equal(t, int32(10), *config.Rules[1].NoncurrentVersionTransitions[0].NoncurrentDays)
			assert.Equal(t, StorageClassIA, config.Rules[1].NoncurrentVersionTransitions[0].StorageClass)
			assert.True(t, *config.Rules[1].NoncurrentVersionTransitions[0].IsAccessTime)
			assert.False(t, *config.Rules[1].NoncurrentVersionTransitions[0].ReturnToStdWhenVisit)
			assert.Equal(t, int64(1631698332), *config.Rules[1].AtimeBase)
		},
	},
	// UnknownFieldsCheck=Fail (default via nil), clean body → no err, Body set, UnknownFields nil
	{
		200,
		map[string]string{
			"x-oss-request-id": "mock-request-id",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(lifecycleRespClean),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
		},
		&GetBucketLifecycleRequest{Bucket: Ptr("bucket")},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, o)
			assert.NotNil(t, o.LifecycleConfiguration)
			assert.Equal(t, 1, len(o.LifecycleConfiguration.Rules))
			assert.NotEmpty(t, o.Body)
			assert.Nil(t, o.UnknownFields)
		},
	},
	// UnknownFieldsCheck=Fail (default), body with unknown fields → *UnknownFieldsError, result preserved
	{
		200,
		map[string]string{
			"x-oss-request-id": "mock-request-id",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(lifecycleRespWithUnknown),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
		},
		&GetBucketLifecycleRequest{Bucket: Ptr("bucket")},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.Error(t, err)
			var ufe *UnknownFieldsError
			assert.True(t, errors.As(err, &ufe), "expected *UnknownFieldsError, got %T: %v", err, err)
			assert.NotNil(t, ufe)
			assert.ElementsMatch(t, []string{
				"LifecycleConfiguration/Rule/NewExtension",
				"LifecycleConfiguration/Rule/Filter/Not/NewFilterExt",
			}, ufe.MissingPaths)
			assert.NotNil(t, o, "result must be preserved when err is *UnknownFieldsError")
			assert.NotNil(t, o.LifecycleConfiguration)
			assert.Equal(t, 2, len(o.LifecycleConfiguration.Rules))
			assert.NotEmpty(t, o.Body)
			assert.ElementsMatch(t, ufe.MissingPaths, o.UnknownFields)
		},
	},
	// UnknownFieldsCheck=Report → no err, UnknownFields populated
	{
		200,
		map[string]string{
			"x-oss-request-id": "mock-request-id",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(lifecycleRespWithUnknown),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
		},
		&GetBucketLifecycleRequest{
			Bucket:             Ptr("bucket"),
			UnknownFieldsCheck: Ptr(UnknownFieldsCheckReport),
		},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, o)
			assert.NotEmpty(t, o.Body)
			assert.ElementsMatch(t, []string{
				"LifecycleConfiguration/Rule/NewExtension",
				"LifecycleConfiguration/Rule/Filter/Not/NewFilterExt",
			}, o.UnknownFields)
			assert.Equal(t, 2, len(o.LifecycleConfiguration.Rules))
		},
	},
	// UnknownFieldsCheck=Skip → no err, Body nil, UnknownFields nil, config parsed
	{
		200,
		map[string]string{
			"x-oss-request-id": "mock-request-id",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(lifecycleRespWithUnknown),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
		},
		&GetBucketLifecycleRequest{
			Bucket:             Ptr("bucket"),
			UnknownFieldsCheck: Ptr(UnknownFieldsCheckSkip),
		},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, o)
			assert.Nil(t, o.Body)
			assert.Nil(t, o.UnknownFields)
			assert.Equal(t, 2, len(o.LifecycleConfiguration.Rules))
		},
	},
	// UnknownFieldsCheck=Fail, empty LifecycleConfiguration → no err, zero rules
	{
		200,
		map[string]string{
			"x-oss-request-id": "mock-request-id",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<LifecycleConfiguration></LifecycleConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
		},
		&GetBucketLifecycleRequest{Bucket: Ptr("bucket")},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, o)
			assert.NotNil(t, o.LifecycleConfiguration)
			assert.Equal(t, 0, len(o.LifecycleConfiguration.Rules))
			assert.NotEmpty(t, o.Body)
			assert.Nil(t, o.UnknownFields)
		},
	},
	// UnknownFieldsCheck=Fail, self-closing empty element → no err, zero rules
	{
		200,
		map[string]string{
			"x-oss-request-id": "mock-request-id",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?><LifecycleConfiguration/>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
		},
		&GetBucketLifecycleRequest{Bucket: Ptr("bucket")},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, o)
			assert.NotNil(t, o.LifecycleConfiguration)
			assert.Equal(t, 0, len(o.LifecycleConfiguration.Rules))
			assert.NotEmpty(t, o.Body)
			assert.Nil(t, o.UnknownFields)
		},
	},
	// UnknownFieldsCheck=Fail explicit (same effect as nil default)
	{
		200,
		map[string]string{
			"x-oss-request-id": "mock-request-id",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(lifecycleRespWithUnknown),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
		},
		&GetBucketLifecycleRequest{
			Bucket:             Ptr("bucket"),
			UnknownFieldsCheck: Ptr(UnknownFieldsCheckFail),
		},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.Error(t, err)
			var ufe *UnknownFieldsError
			assert.True(t, errors.As(err, &ufe))
			assert.NotNil(t, o)
			assert.NotEmpty(t, o.Body)
			assert.Equal(t, 2, len(ufe.MissingPaths))
		},
	},
}

func TestMockGetBucketLifecycle_Success(t *testing.T) {
	for _, c := range testMockGetBucketLifecycleSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLifecycle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketLifecycleErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLifecycleRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLifecycleResult, err error)
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
			assert.Equal(t, "/bucket/?lifecycle", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
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
			assert.Equal(t, "/bucket/?lifecycle", strUrl)
		},
		&GetBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
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
			assert.Equal(t, "/bucket/?lifecycle", strUrl)
		},
		&GetBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLifecycleResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketLifecycle fail")
		},
	},
}

func TestMockGetBucketLifecycle_Error(t *testing.T) {
	for _, c := range testMockGetBucketLifecycleErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLifecycle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketLifecycleSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketLifecycleRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketLifecycleResult, err error)
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
			assert.Equal(t, "/bucket/?lifecycle", strUrl)
		},
		&DeleteBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLifecycleResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteBucketLifecycle_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketLifecycleSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketLifecycle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketLifecycleErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketLifecycleRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketLifecycleResult, err error)
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
			assert.Equal(t, "/bucket/?lifecycle", strUrl)
		},
		&DeleteBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLifecycleResult, err error) {
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
			assert.Equal(t, "/bucket/?lifecycle", strUrl)
		},
		&DeleteBucketLifecycleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLifecycleResult, err error) {
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

func TestMockDeleteBucketLifecycle_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketLifecycleErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketLifecycle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


