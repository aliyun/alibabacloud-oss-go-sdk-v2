package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketWebsite(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketWebsiteRequest
	var input *OperationInput
	var err error

	request = &PutBucketWebsiteRequest{}
	input = &OperationInput{
		OpName: "PutBucketWebsite",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"website": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"website"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutBucketWebsiteRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucketWebsite",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"website": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"website"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, WebsiteConfiguration.")

	request = &PutBucketWebsiteRequest{
		Bucket: Ptr("oss-demo"),
		WebsiteConfiguration: &WebsiteConfiguration{
			IndexDocument: &IndexDocument{
				Suffix:        Ptr("index.html"),
				SupportSubDir: Ptr(true),
				Type:          Ptr(int64(0)),
			},
			ErrorDocument: &ErrorDocument{
				Key:        Ptr("error.html"),
				HttpStatus: Ptr(int64(404)),
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketWebsite",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"website": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"website"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<WebsiteConfiguration><IndexDocument><Suffix>index.html</Suffix><SupportSubDir>true</SupportSubDir><Type>0</Type></IndexDocument><ErrorDocument><Key>error.html</Key><HttpStatus>404</HttpStatus></ErrorDocument></WebsiteConfiguration>")

	request = &PutBucketWebsiteRequest{
		Bucket: Ptr("oss-demo"),
		WebsiteConfiguration: &WebsiteConfiguration{
			IndexDocument: &IndexDocument{
				Suffix:        Ptr("index.html"),
				SupportSubDir: Ptr(true),
				Type:          Ptr(int64(0)),
			},
			ErrorDocument: &ErrorDocument{
				Key:        Ptr("error.html"),
				HttpStatus: Ptr(int64(404)),
			},
			RoutingRules: &RoutingRules{
				[]RoutingRule{
					{
						RuleNumber: Ptr(int64(1)),
						Condition: &RoutingRuleCondition{
							KeyPrefixEquals:             Ptr("abc/"),
							HttpErrorCodeReturnedEquals: Ptr(int64(404)),
						},
						Redirect: &RoutingRuleRedirect{
							RedirectType:          Ptr("Mirror"),
							PassQueryString:       Ptr(true),
							MirrorURL:             Ptr("http://example.com/"),
							MirrorPassQueryString: Ptr(true),
							MirrorFollowRedirect:  Ptr(true),
							MirrorCheckMd5:        Ptr(false),
							MirrorHeaders: &MirrorHeaders{
								PassAll: Ptr(true),
								Passes:   []string{"myheader-key1", "myheader-key2"},
								Removes: []string{"myheader-key3", "myheader-key4"},
								Sets: []MirrorHeadersSet{
									{
										Key:   Ptr("myheader-key5"),
										Value: Ptr("myheader-value5"),
									},
								},
							},
						},
					},
					{
						RuleNumber: Ptr(int64(2)),
						Condition: &RoutingRuleCondition{
							KeyPrefixEquals:             Ptr("abc/"),
							HttpErrorCodeReturnedEquals: Ptr(int64(404)),
							IncludeHeaders: []RoutingRuleIncludeHeader{
								{
									Key:    Ptr("host"),
									Equals: Ptr("test.oss-cn-beijing-internal.aliyuncs.com"),
								},
							},
						},
						Redirect: &RoutingRuleRedirect{
							RedirectType:     Ptr("AliCDN"),
							PassQueryString:  Ptr(false),
							HostName:         Ptr("example.com"),
							ReplaceKeyWith:   Ptr("prefix/${key}.suffix"),
							HttpRedirectCode: Ptr(int64(301)),
							Protocol:         Ptr("http"),
						},
					},
					{
						RuleNumber: Ptr(int64(3)),
						Condition: &RoutingRuleCondition{
							HttpErrorCodeReturnedEquals: Ptr(int64(404)),
						},
						Redirect: &RoutingRuleRedirect{
							RedirectType:        Ptr("External"),
							PassQueryString:     Ptr(false),
							HostName:            Ptr("example.com"),
							ReplaceKeyWith:      Ptr("prefix/${key}"),
							HttpRedirectCode:    Ptr(int64(302)),
							Protocol:            Ptr("http"),
							EnableReplacePrefix: Ptr(false),
						},
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketWebsite",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"website": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"website"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<WebsiteConfiguration><IndexDocument><Suffix>index.html</Suffix><SupportSubDir>true</SupportSubDir><Type>0</Type></IndexDocument><ErrorDocument><Key>error.html</Key><HttpStatus>404</HttpStatus></ErrorDocument><RoutingRules><RoutingRule><RuleNumber>1</RuleNumber><Condition><KeyPrefixEquals>abc/</KeyPrefixEquals><HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals></Condition><Redirect><MirrorURL>http://example.com/</MirrorURL><MirrorFollowRedirect>true</MirrorFollowRedirect><PassQueryString>true</PassQueryString><MirrorHeaders><PassAll>true</PassAll><Pass>myheader-key1</Pass><Pass>myheader-key2</Pass><Remove>myheader-key3</Remove><Remove>myheader-key4</Remove><Set><Key>myheader-key5</Key><Value>myheader-value5</Value></Set></MirrorHeaders><RedirectType>Mirror</RedirectType><MirrorCheckMd5>false</MirrorCheckMd5><MirrorPassQueryString>true</MirrorPassQueryString></Redirect></RoutingRule><RoutingRule><RuleNumber>2</RuleNumber><Condition><KeyPrefixEquals>abc/</KeyPrefixEquals><HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals><IncludeHeader><Key>host</Key><Equals>test.oss-cn-beijing-internal.aliyuncs.com</Equals></IncludeHeader></Condition><Redirect><ReplaceKeyWith>prefix/${key}.suffix</ReplaceKeyWith><HostName>example.com</HostName><PassQueryString>false</PassQueryString><RedirectType>AliCDN</RedirectType><Protocol>http</Protocol><HttpRedirectCode>301</HttpRedirectCode></Redirect></RoutingRule><RoutingRule><RuleNumber>3</RuleNumber><Condition><HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals></Condition><Redirect><EnableReplacePrefix>false</EnableReplacePrefix><ReplaceKeyWith>prefix/${key}</ReplaceKeyWith><HostName>example.com</HostName><PassQueryString>false</PassQueryString><RedirectType>External</RedirectType><Protocol>http</Protocol><HttpRedirectCode>302</HttpRedirectCode></Redirect></RoutingRule></RoutingRules></WebsiteConfiguration>")

	request = &PutBucketWebsiteRequest{
		Bucket: Ptr("oss-demo"),
		WebsiteConfiguration: &WebsiteConfiguration{
			IndexDocument: &IndexDocument{
				Suffix:        Ptr("index.html"),
				SupportSubDir: Ptr(true),
				Type:          Ptr(int64(0)),
			},
			ErrorDocument: &ErrorDocument{
				Key:        Ptr("error.html"),
				HttpStatus: Ptr(int64(404)),
			},
			RoutingRules: &RoutingRules{
				[]RoutingRule{
					{
						RuleNumber: Ptr(int64(1)),
						Condition: &RoutingRuleCondition{
							KeyPrefixEquals:             Ptr("abc/"),
							HttpErrorCodeReturnedEquals: Ptr(int64(404)),
							KeySuffixEquals:             Ptr(".txt"),
						},
						LuaConfig: &RoutingRuleLuaConfig{
							Script: Ptr("test.lua"),
						},
						Redirect: &RoutingRuleRedirect{
							MirrorPassOriginalSlashes:      Ptr(false),
							RedirectType:                   Ptr("Mirror"),
							PassQueryString:                Ptr(true),
							MirrorURL:                      Ptr("http://example.com/"),
							MirrorPassQueryString:          Ptr(true),
							MirrorSNI:                      Ptr(true),
							ReplaceKeyPrefixWith:           Ptr("def/"),
							MirrorFollowRedirect:           Ptr(true),
							HostName:                       Ptr("example.com"),
							MirrorCheckMd5:                 Ptr(true),
							EnableReplacePrefix:            Ptr(true),
							HttpRedirectCode:               Ptr(int64(301)),
							MirrorURLSlave:                 Ptr("http://example.com/"),
							MirrorSaveOssMeta:              Ptr(true),
							MirrorProxyPass:                Ptr(false),
							MirrorAllowGetImageInfo:        Ptr(true),
							MirrorAllowVideoSnapshot:       Ptr(false),
							MirrorIsExpressTunnel:          Ptr(true),
							MirrorDstRegion:                Ptr("cn-hangzhou"),
							MirrorDstVpcId:                 Ptr("vpc-test-id"),
							MirrorDstSlaveVpcId:            Ptr("vpc-test-id"),
							MirrorUserLastModified:         Ptr(false),
							MirrorSwitchAllErrors:          Ptr(true),
							MirrorTunnelId:                 Ptr("test-tunnel-id"),
							MirrorUsingRole:                Ptr(false),
							MirrorRole:                     Ptr("aliyun-test-role"),
							MirrorAllowHeadObject:          Ptr(true),
							TransparentMirrorResponseCodes: Ptr("400"),
							MirrorAsyncStatus:              Ptr(int64(303)),
							MirrorTaggings: &MirrorTaggings{
								Taggings: []MirrorTagging{
									{
										Key:   Ptr("k"),
										Value: Ptr("v"),
									},
								},
							},
							MirrorReturnHeaders: &MirrorReturnHeaders{
								ReturnHeaders: []ReturnHeader{
									{
										Key:   Ptr("k"),
										Value: Ptr("v"),
									},
								},
							},
							MirrorAuth: &MirrorAuth{
								AuthType:        Ptr("S3V4"),
								Region:          Ptr("ap-southeast-1"),
								AccessKeyId:     Ptr("TESTAK"),
								AccessKeySecret: Ptr("TESTSK"),
							},
							MirrorMultiAlternates: &MirrorMultiAlternates{
								MirrorMultiAlternates: []MirrorMultiAlternate{
									{
										MirrorMultiAlternateNumber:    Ptr(int64(32)),
										MirrorMultiAlternateURL:       Ptr("https://test-multi-alter.example.com"),
										MirrorMultiAlternateVpcId:     Ptr("vpc-test-id"),
										MirrorMultiAlternateDstRegion: Ptr("ap-southeast-1"),
									},
								},
							},
							MirrorHeaders: &MirrorHeaders{
								PassAll: Ptr(true),
								Passes:   []string{"myheader-key1", "myheader-key2"},
								Removes: []string{"myheader-key3", "myheader-key4"},
								Sets: []MirrorHeadersSet{
									{
										Key:   Ptr("myheader-key5"),
										Value: Ptr("myheader-value5"),
									},
								},
							},
						},
					},
					{
						RuleNumber: Ptr(int64(2)),
						Condition: &RoutingRuleCondition{
							KeyPrefixEquals:             Ptr("abc/"),
							HttpErrorCodeReturnedEquals: Ptr(int64(404)),
							KeySuffixEquals:             Ptr(".txt"),
						},
						LuaConfig: &RoutingRuleLuaConfig{
							Script: Ptr("test.lua"),
						},
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketWebsite",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"website": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"website"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<WebsiteConfiguration><IndexDocument><Suffix>index.html</Suffix><SupportSubDir>true</SupportSubDir><Type>0</Type></IndexDocument><ErrorDocument><Key>error.html</Key><HttpStatus>404</HttpStatus></ErrorDocument><RoutingRules><RoutingRule><RuleNumber>1</RuleNumber><Condition><KeyPrefixEquals>abc/</KeyPrefixEquals><KeySuffixEquals>.txt</KeySuffixEquals><HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals></Condition><Redirect><MirrorURL>http://example.com/</MirrorURL><MirrorFollowRedirect>true</MirrorFollowRedirect><EnableReplacePrefix>true</EnableReplacePrefix><HostName>example.com</HostName><PassQueryString>true</PassQueryString><MirrorHeaders><PassAll>true</PassAll><Pass>myheader-key1</Pass><Pass>myheader-key2</Pass><Remove>myheader-key3</Remove><Remove>myheader-key4</Remove><Set><Key>myheader-key5</Key><Value>myheader-value5</Value></Set></MirrorHeaders><ReplaceKeyPrefixWith>def/</ReplaceKeyPrefixWith><RedirectType>Mirror</RedirectType><MirrorSNI>true</MirrorSNI><MirrorCheckMd5>true</MirrorCheckMd5><HttpRedirectCode>301</HttpRedirectCode><MirrorPassOriginalSlashes>false</MirrorPassOriginalSlashes><MirrorPassQueryString>true</MirrorPassQueryString><MirrorAsyncStatus>303</MirrorAsyncStatus><MirrorAuth><AuthType>S3V4</AuthType><Region>ap-southeast-1</Region><AccessKeyId>TESTAK</AccessKeyId><AccessKeySecret>TESTSK</AccessKeySecret></MirrorAuth><MirrorAllowVideoSnapshot>false</MirrorAllowVideoSnapshot><MirrorURLSlave>http://example.com/</MirrorURLSlave><MirrorDstVpcId>vpc-test-id</MirrorDstVpcId><MirrorUserLastModified>false</MirrorUserLastModified><MirrorUsingRole>false</MirrorUsingRole><MirrorIsExpressTunnel>true</MirrorIsExpressTunnel><MirrorProxyPass>false</MirrorProxyPass><MirrorTaggings><Taggings><Key>k</Key><Value>v</Value></Taggings></MirrorTaggings><MirrorDstSlaveVpcId>vpc-test-id</MirrorDstSlaveVpcId><MirrorDstRegion>cn-hangzhou</MirrorDstRegion><MirrorSwitchAllErrors>true</MirrorSwitchAllErrors><MirrorTunnelId>test-tunnel-id</MirrorTunnelId><MirrorRole>aliyun-test-role</MirrorRole><MirrorAllowGetImageInfo>true</MirrorAllowGetImageInfo><MirrorSaveOssMeta>true</MirrorSaveOssMeta><MirrorAllowHeadObject>true</MirrorAllowHeadObject><MirrorMultiAlternates><MirrorMultiAlternate><MirrorMultiAlternateDstRegion>ap-southeast-1</MirrorMultiAlternateDstRegion><MirrorMultiAlternateNumber>32</MirrorMultiAlternateNumber><MirrorMultiAlternateURL>https://test-multi-alter.example.com</MirrorMultiAlternateURL><MirrorMultiAlternateVpcId>vpc-test-id</MirrorMultiAlternateVpcId></MirrorMultiAlternate></MirrorMultiAlternates><TransparentMirrorResponseCodes>400</TransparentMirrorResponseCodes><MirrorReturnHeaders><ReturnHeader><Key>k</Key><Value>v</Value></ReturnHeader></MirrorReturnHeaders></Redirect><LuaConfig><Script>test.lua</Script></LuaConfig></RoutingRule><RoutingRule><RuleNumber>2</RuleNumber><Condition><KeyPrefixEquals>abc/</KeyPrefixEquals><KeySuffixEquals>.txt</KeySuffixEquals><HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals></Condition><LuaConfig><Script>test.lua</Script></LuaConfig></RoutingRule></RoutingRules></WebsiteConfiguration>")
}

func TestUnmarshalOutput_PutBucketWebsite(t *testing.T) {
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
	result := &PutBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
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
	result = &PutBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	result = &PutBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketWebsite(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketWebsiteRequest
	var input *OperationInput
	var err error

	request = &GetBucketWebsiteRequest{}
	input = &OperationInput{
		OpName: "GetBucketWebsite",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"website": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"website"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketWebsiteRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketWebsite",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"website": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"website"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketWebsite(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
			<WebsiteConfiguration>
				<IndexDocument>
					<Suffix>index.html</Suffix>
				</IndexDocument>
				<ErrorDocument>
				   <Key>error.html</Key>
				   <HttpStatus>404</HttpStatus>
				</ErrorDocument>
			</WebsiteConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.WebsiteConfiguration.IndexDocument.Suffix, "index.html")
	assert.Nil(t, result.WebsiteConfiguration.IndexDocument.SupportSubDir)
	assert.Nil(t, result.WebsiteConfiguration.IndexDocument.Type)
	assert.Equal(t, *result.WebsiteConfiguration.ErrorDocument.Key, "error.html")
	assert.Equal(t, *result.WebsiteConfiguration.ErrorDocument.HttpStatus, int64(404))

	body = `<WebsiteConfiguration>
		  <IndexDocument>
			<Suffix>index.html</Suffix>
			<SupportSubDir>true</SupportSubDir>
			<Type>0</Type>
		  </IndexDocument>
		  <ErrorDocument>
			<Key>error.html</Key>
			<HttpStatus>404</HttpStatus>
		  </ErrorDocument>
		  <RoutingRules>
			<RoutingRule>
			  <RuleNumber>1</RuleNumber>
			  <Condition>
				<KeyPrefixEquals>abc/</KeyPrefixEquals>
				<HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals>
			  </Condition>
			  <Redirect>
				<RedirectType>Mirror</RedirectType>
				<PassQueryString>true</PassQueryString>
				<MirrorURL>http://example.com/</MirrorURL>   
				<MirrorPassQueryString>true</MirrorPassQueryString>
				<MirrorFollowRedirect>true</MirrorFollowRedirect>
				<MirrorCheckMd5>false</MirrorCheckMd5>
				<MirrorHeaders>
				  <PassAll>true</PassAll>
				  <Pass>myheader-key1</Pass>
				  <Pass>myheader-key2</Pass>
				  <Remove>myheader-key3</Remove>
				  <Remove>myheader-key4</Remove>
				  <Set>
					<Key>myheader-key5</Key>
					<Value>myheader-value5</Value>
				  </Set>
				</MirrorHeaders>
			  </Redirect>
			</RoutingRule>
			<RoutingRule>
			  <RuleNumber>2</RuleNumber>
			  <Condition>
				<KeyPrefixEquals>abc/</KeyPrefixEquals>
				<HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals>
				<IncludeHeader>
				  <Key>host</Key>
				  <Equals>test.oss-cn-beijing-internal.aliyuncs.com</Equals>
				</IncludeHeader>
			  </Condition>
			  <Redirect>
				<RedirectType>AliCDN</RedirectType>
				<Protocol>http</Protocol>
				<HostName>example.com</HostName>
				<PassQueryString>false</PassQueryString>
				<ReplaceKeyWith>prefix/${key}.suffix</ReplaceKeyWith>
				<HttpRedirectCode>301</HttpRedirectCode>
			  </Redirect>
			</RoutingRule>
			<RoutingRule>
			  <Condition>
				<HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals>
			  </Condition>
			  <RuleNumber>3</RuleNumber>
			  <Redirect>
				<ReplaceKeyWith>prefix/${key}</ReplaceKeyWith>
				<HttpRedirectCode>302</HttpRedirectCode>
				<EnableReplacePrefix>false</EnableReplacePrefix>
				<PassQueryString>false</PassQueryString>
				<Protocol>http</Protocol>
				<HostName>example.com</HostName>
				<RedirectType>External</RedirectType>
			  </Redirect>
			</RoutingRule>
		  </RoutingRules>
		</WebsiteConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &GetBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.WebsiteConfiguration.IndexDocument.Suffix, "index.html")
	assert.True(t, *result.WebsiteConfiguration.IndexDocument.SupportSubDir)
	assert.Equal(t, *result.WebsiteConfiguration.IndexDocument.Type, int64(0))
	assert.Equal(t, *result.WebsiteConfiguration.ErrorDocument.Key, "error.html")
	assert.Equal(t, *result.WebsiteConfiguration.ErrorDocument.HttpStatus, int64(404))
	assert.Equal(t, len(result.WebsiteConfiguration.RoutingRules.RoutingRules), 3)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].RuleNumber, int64(1))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Condition.KeyPrefixEquals, "abc/")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Condition.HttpErrorCodeReturnedEquals, int64(404))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.RedirectType, "Mirror")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.PassQueryString, true)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorURL, "http://example.com/")
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorPassQueryString)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorFollowRedirect)
	assert.False(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorCheckMd5)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.PassAll)
	assert.Equal(t, result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Passes[0], "myheader-key1")
	assert.Equal(t, result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Passes[1], "myheader-key2")
	assert.Equal(t, result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Removes[0], "myheader-key3")
	assert.Equal(t, result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Removes[1], "myheader-key4")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Sets[0].Key, "myheader-key5")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Sets[0].Value, "myheader-value5")

	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].RuleNumber, int64(2))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Condition.KeyPrefixEquals, "abc/")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Condition.HttpErrorCodeReturnedEquals, int64(404))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Condition.IncludeHeaders[0].Key, "host")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Condition.IncludeHeaders[0].Equals, "test.oss-cn-beijing-internal.aliyuncs.com")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.RedirectType, "AliCDN")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.Protocol, "http")
	assert.False(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.PassQueryString)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.ReplaceKeyWith, "prefix/${key}.suffix")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.HttpRedirectCode, int64(301))

	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[2].RuleNumber, int64(3))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[2].Redirect.RedirectType, "External")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[2].Redirect.PassQueryString, false)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[2].Redirect.ReplaceKeyWith, "prefix/${key}")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[2].Redirect.HttpRedirectCode, int64(302))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[2].Redirect.EnableReplacePrefix, false)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[2].Redirect.Protocol, "http")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[2].Redirect.HostName, "example.com")

	body = `<WebsiteConfiguration><IndexDocument><Suffix>index.html</Suffix><SupportSubDir>true</SupportSubDir><Type>0</Type></IndexDocument><ErrorDocument><Key>error.html</Key><HttpStatus>404</HttpStatus></ErrorDocument><RoutingRules><RoutingRule><Redirect><MirrorPassOriginalSlashes>false</MirrorPassOriginalSlashes><RedirectType>Mirror</RedirectType><MirrorURL>http://example.com/</MirrorURL><MirrorPassQueryString>true</MirrorPassQueryString><MirrorCheckMd5>true</MirrorCheckMd5><MirrorSNI>true</MirrorSNI><ReplaceKeyPrefixWith>def/</ReplaceKeyPrefixWith><MirrorFollowRedirect>true</MirrorFollowRedirect><HostName>example.com</HostName><MirrorHeaders><Pass>myheader-key1</Pass><Pass>myheader-key2</Pass><Set><Key>myheader-key5</Key><Value>myheader-value5</Value></Set><PassAll>true</PassAll></MirrorHeaders><PassQueryString>true</PassQueryString><EnableReplacePrefix>true</EnableReplacePrefix><HttpRedirectCode>301</HttpRedirectCode><MirrorURLSlave>http://example.com/</MirrorURLSlave><MirrorSaveOssMeta>true</MirrorSaveOssMeta><MirrorProxyPass>false</MirrorProxyPass><MirrorAllowGetImageInfo>true</MirrorAllowGetImageInfo><MirrorAllowVideoSnapshot>false</MirrorAllowVideoSnapshot><MirrorIsExpressTunnel>true</MirrorIsExpressTunnel><MirrorDstRegion>cn-hangzhou</MirrorDstRegion><MirrorDstVpcId>vpc-test-id</MirrorDstVpcId><MirrorDstSlaveVpcId>vpc-test-id</MirrorDstSlaveVpcId><MirrorUserLastModified>false</MirrorUserLastModified><MirrorSwitchAllErrors>true</MirrorSwitchAllErrors><MirrorTunnelId>test-tunnel-id</MirrorTunnelId><MirrorUsingRole>false</MirrorUsingRole><MirrorRole>aliyun-test-role</MirrorRole><MirrorAllowHeadObject>true</MirrorAllowHeadObject><TransparentMirrorResponseCodes>400</TransparentMirrorResponseCodes><MirrorAsyncStatus>303</MirrorAsyncStatus><MirrorTaggings><Taggings><Value>v</Value><Key>k</Key></Taggings></MirrorTaggings><MirrorReturnHeaders><ReturnHeader><Key>k</Key><Value>v</Value></ReturnHeader></MirrorReturnHeaders><MirrorAuth><AuthType>S3V4</AuthType><Region>ap-southeast-1</Region><AccessKeyId>TESTAK</AccessKeyId><AccessKeySecret>TESTSK</AccessKeySecret></MirrorAuth><MirrorMultiAlternates><MirrorMultiAlternate><MirrorMultiAlternateNumber>32</MirrorMultiAlternateNumber><MirrorMultiAlternateURL>https://test-multi-alter.example.com</MirrorMultiAlternateURL><MirrorMultiAlternateVpcId>vpc-test-id</MirrorMultiAlternateVpcId><MirrorMultiAlternateDstRegion>ap-southeast-1</MirrorMultiAlternateDstRegion></MirrorMultiAlternate></MirrorMultiAlternates></Redirect><RuleNumber>1</RuleNumber><Condition><KeyPrefixEquals>abc/</KeyPrefixEquals><KeySuffixEquals>.txt</KeySuffixEquals><HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals></Condition><LuaConfig><Script>test.lua</Script></LuaConfig></RoutingRule><RoutingRule><Redirect><RedirectType>AliCDN</RedirectType><MirrorURL>http://example.com/</MirrorURL><MirrorPassQueryString>true</MirrorPassQueryString><MirrorCheckMd5>true</MirrorCheckMd5><MirrorSNI>true</MirrorSNI><Protocol>http</Protocol><MirrorFollowRedirect>true</MirrorFollowRedirect><MirrorHeaders><Pass>myheader-key1</Pass><Pass>myheader-key2</Pass><Set><Key>myheader-key5</Key><Value>myheader-value5</Value></Set><PassAll>true</PassAll></MirrorHeaders><PassQueryString>true</PassQueryString><ReplaceKeyWith>abc</ReplaceKeyWith></Redirect><RuleNumber>2</RuleNumber><Condition><KeyPrefixEquals>abc/</KeyPrefixEquals><KeySuffixEquals>.txt</KeySuffixEquals><HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals></Condition><LuaConfig><Script>test.lua</Script></LuaConfig></RoutingRule></RoutingRules></WebsiteConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &GetBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	assert.Equal(t, *result.WebsiteConfiguration.IndexDocument.Suffix, "index.html")
	assert.True(t, *result.WebsiteConfiguration.IndexDocument.SupportSubDir)
	assert.Equal(t, *result.WebsiteConfiguration.IndexDocument.Type, int64(0))
	assert.Equal(t, *result.WebsiteConfiguration.ErrorDocument.Key, "error.html")
	assert.Equal(t, *result.WebsiteConfiguration.ErrorDocument.HttpStatus, int64(404))
	assert.Equal(t, len(result.WebsiteConfiguration.RoutingRules.RoutingRules), 2)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].RuleNumber, int64(1))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Condition.KeyPrefixEquals, "abc/")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Condition.HttpErrorCodeReturnedEquals, int64(404))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Condition.KeySuffixEquals, ".txt")

	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.RedirectType, "Mirror")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorPassOriginalSlashes, false)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorURL, "http://example.com/")
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorPassQueryString)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorFollowRedirect)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.PassQueryString, true)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorCheckMd5)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorSNI)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.ReplaceKeyPrefixWith, "def/")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.HostName, "example.com")
	assert.Equal(t, result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Passes[0], "myheader-key1")
	assert.Equal(t, result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Passes[1], "myheader-key2")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Sets[0].Key, "myheader-key5")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.Sets[0].Value, "myheader-value5")
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorHeaders.PassAll)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.EnableReplacePrefix)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.HttpRedirectCode, int64(301))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorURLSlave, "http://example.com/")
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorSaveOssMeta)
	assert.False(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorProxyPass)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorAllowGetImageInfo)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorIsExpressTunnel)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorDstRegion, "cn-hangzhou")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorDstVpcId, "vpc-test-id")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorDstSlaveVpcId, "vpc-test-id")
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorSwitchAllErrors)
	assert.False(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorUsingRole)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorRole, "aliyun-test-role")
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorAllowHeadObject)
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.TransparentMirrorResponseCodes, "400")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorAsyncStatus, int64(303))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorTaggings.Taggings[0].Key, "k")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorTaggings.Taggings[0].Value, "v")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorReturnHeaders.ReturnHeaders[0].Key, "k")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorReturnHeaders.ReturnHeaders[0].Value, "v")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorAuth.AuthType, "S3V4")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorAuth.Region, "ap-southeast-1")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorAuth.AccessKeyId, "TESTAK")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorAuth.AccessKeySecret, "TESTSK")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorMultiAlternates.MirrorMultiAlternates[0].MirrorMultiAlternateURL, "https://test-multi-alter.example.com")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].Redirect.MirrorMultiAlternates.MirrorMultiAlternates[0].MirrorMultiAlternateDstRegion, "ap-southeast-1")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[0].LuaConfig.Script, "test.lua")

	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].RuleNumber, int64(2))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Condition.KeyPrefixEquals, "abc/")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Condition.KeySuffixEquals, ".txt")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Condition.HttpErrorCodeReturnedEquals, int64(404))
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.RedirectType, "AliCDN")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.Protocol, "http")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].LuaConfig.Script, "test.lua")
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.PassQueryString)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.MirrorSNI)
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.MirrorCheckMd5)
	assert.Equal(t, result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.MirrorHeaders.Passes[0], "myheader-key1")
	assert.Equal(t, result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.MirrorHeaders.Passes[1], "myheader-key2")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.MirrorHeaders.Sets[0].Key, "myheader-key5")
	assert.Equal(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.MirrorHeaders.Sets[0].Value, "myheader-value5")
	assert.True(t, *result.WebsiteConfiguration.RoutingRules.RoutingRules[1].Redirect.MirrorHeaders.PassAll)

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	result = &GetBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteBucketWebsite(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteBucketWebsiteRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketWebsiteRequest{}
	input = &OperationInput{
		OpName: "DeleteBucketWebsite",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"website": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"website"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteBucketWebsiteRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketWebsite",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"website": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"website"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteBucketWebsite(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	body := `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>66C2FF09FDF07830343C72EC</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
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
	result = &DeleteBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	result = &DeleteBucketWebsiteResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
