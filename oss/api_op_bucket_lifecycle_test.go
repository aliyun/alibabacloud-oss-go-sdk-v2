package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketLifecycle(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketLifecycleRequest
	var input *OperationInput
	var err error

	request = &PutBucketLifecycleRequest{}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, LifecycleConfiguration.")

	// demo 1
	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
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
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule</ID><Prefix>log/</Prefix><Transition><Days>30</Days><StorageClass>IA</StorageClass></Transition></Rule></LifecycleConfiguration>")

	// demo 2
	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					ID:     Ptr("rule"),
					Prefix: Ptr("log/"),
					Status: Ptr("Enabled"),
					Expiration: &LifecycleRuleExpiration{
						Days: Ptr(int32(90)),
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule</ID><Prefix>log/</Prefix><Expiration><Days>90</Days></Expiration></Rule></LifecycleConfiguration>")
	// demo 3
	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					ID:     Ptr("rule"),
					Prefix: Ptr("log/"),
					Status: Ptr("Enabled"),
					Transitions: []LifecycleRuleTransition{
						{
							Days:         Ptr(int32(30)),
							StorageClass: StorageClassIA,
						},
						{
							Days:         Ptr(int32(60)),
							StorageClass: StorageClassArchive,
						},
					},
					Expiration: &LifecycleRuleExpiration{
						Days: Ptr(int32(3600)),
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule</ID><Prefix>log/</Prefix><Expiration><Days>3600</Days></Expiration><Transition><Days>30</Days><StorageClass>IA</StorageClass></Transition><Transition><Days>60</Days><StorageClass>Archive</StorageClass></Transition></Rule></LifecycleConfiguration>")

	// demo 4
	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					ID:     Ptr("rule"),
					Prefix: Ptr(""),
					Status: Ptr("Enabled"),
					Expiration: &LifecycleRuleExpiration{
						ExpiredObjectDeleteMarker: Ptr(true),
					},
					NoncurrentVersionExpiration: &NoncurrentVersionExpiration{
						NoncurrentDays: Ptr(int32(5)),
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule</ID><Prefix></Prefix><Expiration><ExpiredObjectDeleteMarker>true</ExpiredObjectDeleteMarker></Expiration><NoncurrentVersionExpiration><NoncurrentDays>5</NoncurrentDays></NoncurrentVersionExpiration></Rule></LifecycleConfiguration>")

	// demo 5
	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					ID:     Ptr("rule"),
					Prefix: Ptr(""),
					Status: Ptr("Enabled"),
					Filter: &LifecycleRuleFilter{
						Not: &LifecycleRuleNot{
							Prefix: Ptr("log"),
							Tag: &Tag{
								Key:   Ptr("key1"),
								Value: Ptr("value1"),
							},
						},
					},
					Transitions: []LifecycleRuleTransition{
						{
							Days:         Ptr(int32(30)),
							StorageClass: StorageClassArchive,
						},
					},
					Expiration: &LifecycleRuleExpiration{
						Days: Ptr(int32(100)),
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LifecycleConfiguration><Rule><Status>Enabled</Status><Filter><Not><Tag><Key>key1</Key><Value>value1</Value></Tag><Prefix>log</Prefix></Not></Filter><ID>rule</ID><Prefix></Prefix><Expiration><Days>100</Days></Expiration><Transition><Days>30</Days><StorageClass>Archive</StorageClass></Transition></Rule></LifecycleConfiguration>")

	// demo 6
	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					ID:     Ptr("rule"),
					Prefix: Ptr("log/"),
					Status: Ptr("Enabled"),
					Transitions: []LifecycleRuleTransition{
						{
							Days:                 Ptr(int32(30)),
							StorageClass:         StorageClassArchive,
							IsAccessTime:         Ptr(true),
							ReturnToStdWhenVisit: Ptr(true),
						},
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule</ID><Prefix>log/</Prefix><Transition><Days>30</Days><StorageClass>Archive</StorageClass><IsAccessTime>true</IsAccessTime><ReturnToStdWhenVisit>true</ReturnToStdWhenVisit></Transition></Rule></LifecycleConfiguration>")

	// demo 7
	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					ID:     Ptr("rule"),
					Prefix: Ptr(""),
					Status: Ptr("Enabled"),
					AbortMultipartUpload: &LifecycleRuleAbortMultipartUpload{
						Days: Ptr(int32(30)),
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LifecycleConfiguration><Rule><Status>Enabled</Status><AbortMultipartUpload><Days>30</Days></AbortMultipartUpload><ID>rule</ID><Prefix></Prefix></Rule></LifecycleConfiguration>")

	// demo 8
	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					ID:     Ptr("rule1"),
					Prefix: Ptr("dir1"),
					Status: Ptr("Enabled"),
					Expiration: &LifecycleRuleExpiration{
						Days: Ptr(int32(180)),
					},
				},
				{
					ID:     Ptr("rule2"),
					Prefix: Ptr("dir1/dir2/"),
					Status: Ptr("Enabled"),
					Expiration: &LifecycleRuleExpiration{
						Days: Ptr(int32(30)),
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>rule1</ID><Prefix>dir1</Prefix><Expiration><Days>180</Days></Expiration></Rule><Rule><Status>Enabled</Status><ID>rule2</ID><Prefix>dir1/dir2/</Prefix><Expiration><Days>30</Days></Expiration></Rule></LifecycleConfiguration>")

	// demo 9
	expirationDays := &LifecycleRuleExpiration{
		Days:                      Ptr(int32(40)),
		ExpiredObjectDeleteMarker: Ptr(false),
	}
	now := time.Now()
	future := now.AddDate(0, 0, 40)
	utcTime := time.Date(future.Year(), future.Month(), future.Day(), 0, 0, 0, 0, time.UTC)
	expirationDate := &LifecycleRuleExpiration{
		CreatedBeforeDate: Ptr(utcTime.Format("2006-01-02T15:04:05.000Z")),
	}
	expirationDelete := &LifecycleRuleExpiration{
		ExpiredObjectDeleteMarker: Ptr(true),
	}
	request = &PutBucketLifecycleRequest{
		Bucket: Ptr("demo-walker-6961"),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					ID:         Ptr("r0"),
					Prefix:     Ptr("prefix0"),
					Status:     Ptr("Enabled"),
					Expiration: expirationDays,
				},
				{
					ID:         Ptr("r1"),
					Prefix:     Ptr("prefix1"),
					Status:     Ptr("Enabled"),
					Expiration: expirationDate,
					Filter: &LifecycleRuleFilter{
						ObjectSizeGreaterThan: Ptr(int64(500)),
						ObjectSizeLessThan:    Ptr(int64(64500)),
					},
				},
				{
					ID:         Ptr("r3"),
					Prefix:     Ptr("prefix3"),
					Status:     Ptr("Enabled"),
					Expiration: expirationDays,
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
					ID:         Ptr("r4"),
					Prefix:     Ptr("prefix4"),
					Status:     Ptr("Enabled"),
					Expiration: expirationDelete,
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
					Prefix:     Ptr("pre_"),
					Status:     Ptr("Enabled"),
					Expiration: expirationDate,
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketLifecycle",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<LifecycleConfiguration><Rule><Status>Enabled</Status><ID>r0</ID><Prefix>prefix0</Prefix><Expiration><Days>40</Days><ExpiredObjectDeleteMarker>false</ExpiredObjectDeleteMarker></Expiration></Rule><Rule><Status>Enabled</Status><Filter><ObjectSizeGreaterThan>500</ObjectSizeGreaterThan><ObjectSizeLessThan>64500</ObjectSizeLessThan></Filter><ID>r1</ID><Prefix>prefix1</Prefix><Expiration><CreatedBeforeDate>"+utcTime.Format("2006-01-02T15:04:05.000Z")+"</CreatedBeforeDate></Expiration></Rule><Rule><Status>Enabled</Status><Filter><ObjectSizeGreaterThan>500</ObjectSizeGreaterThan><ObjectSizeLessThan>64500</ObjectSizeLessThan></Filter><ID>r3</ID><Prefix>prefix3</Prefix><Expiration><Days>40</Days><ExpiredObjectDeleteMarker>false</ExpiredObjectDeleteMarker></Expiration><Transition><Days>30</Days><StorageClass>IA</StorageClass><IsAccessTime>false</IsAccessTime></Transition></Rule><Rule><Status>Enabled</Status><AbortMultipartUpload><CreatedBeforeDate>2015-11-11T00:00:00.000Z</CreatedBeforeDate></AbortMultipartUpload><NoncurrentVersionTransition><IsAccessTime>true</IsAccessTime><ReturnToStdWhenVisit>true</ReturnToStdWhenVisit><NoncurrentDays>10</NoncurrentDays><StorageClass>IA</StorageClass></NoncurrentVersionTransition><ID>r4</ID><Prefix>prefix4</Prefix><Expiration><ExpiredObjectDeleteMarker>true</ExpiredObjectDeleteMarker></Expiration></Rule><Rule><Status>Enabled</Status><Prefix>pre_</Prefix><Expiration><CreatedBeforeDate>"+utcTime.Format("2006-01-02T15:04:05.000Z")+"</CreatedBeforeDate></Expiration></Rule></LifecycleConfiguration>")
}

func TestUnmarshalOutput_PutBucketLifecycle(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	body := `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>534B371674E88A4D8906****</RequestId>
		<HostId>bucket-not-exist.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket-not-exist</BucketName>
		<EC>0015-00000101</EC>
	</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidRequest</Code>
  <Message>Found mixed CreatedBeforeDate and Days based Expiration and Transition for StorageClass IA actions</Message>
  <RequestId>534B371674E88A4D8906****</RequestId>
  <HostId>oss.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0014-00000080</EC>
</Error>`
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidRequest",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidRequest")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketLifecycle(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketLifecycleRequest
	var input *OperationInput
	var err error

	request = &GetBucketLifecycleRequest{}
	input = &OperationInput{
		OpName: "GetBucketLifecycle",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketLifecycle",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketLifecycle(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
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
</LifecycleConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	config := *result.LifecycleConfiguration
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

	body = `<?xml version="1.0" encoding="UTF-8"?>
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
</LifecycleConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	config = *result.LifecycleConfiguration
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

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <BucketName>oss-example</BucketName>
  <Code>NoSuchLifecycle</Code>
  <Message>No Row found in Lifecycle Table.</Message>
  <RequestId>534B371674E88A4D8906****</RequestId>
  <HostId>BucketName.oss.example.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchLifecycle",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchLifecycle")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteBucketLifecycle(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteBucketLifecycleRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketLifecycleRequest{}
	input = &OperationInput{
		OpName: "DeleteBucketLifecycle",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteBucketLifecycleRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketLifecycle",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"lifecycle": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"lifecycle"})
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteBucketLifecycle(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &DeleteBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>534B371674E88A4D8906****</RequestId>
  <HostId>bucket-not-exist.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>bucket-not-exist</BucketName>
  <EC>0015-00000101</EC>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketLifecycleResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
