package dataprocess

import (
	"context"
	"errors"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutDataPipelineConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutDataPipelineConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *PutDataPipelineConfigurationResult, err error)
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
			assert.Equal(t, "/?action=putDataPipelineConfiguration&dataPipeline&dataPipelineName=data-pipeline&role=AliyunOSSDataPipelineRole", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<DataPipelineConfiguration><DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription><Sources><InputBucket>bucket</InputBucket><InputDataScope>All</InputDataScope><FilterConfiguration><PrefixSet>prefix1</PrefixSet><PrefixSet>prefix2/prefix3</PrefixSet><ObjectMediaTypes>text</ObjectMediaTypes><ObjectMediaTypes>image</ObjectMediaTypes><ObjectMediaTypes>video</ObjectMediaTypes></FilterConfiguration></Sources><DataPipelineEmbeddingConfiguration><EmbeddingProvider>bailian</EmbeddingProvider><ApiKey>your_api_key</ApiKey><Model>qwen2.5-vl-embedding</Model><FPS>1</FPS></DataPipelineEmbeddingConfiguration><Destination><VectorBucketName>my-vector-bucket</VectorBucketName><VectorKeyPrefix>prefix</VectorKeyPrefix><VectorIndexNames>my-index</VectorIndexNames><ObjectTagToMetadata>key1</ObjectTagToMetadata><ObjectTagToMetadata>key2</ObjectTagToMetadata><UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata></Destination><DataPipelineError><ErrorMode>ignoreAndRecord</ErrorMode><ErrorBucket>my-error-bucket</ErrorBucket><ErrorPrefix>error-output/</ErrorPrefix></DataPipelineError></DataPipelineConfiguration>")
		},
		&PutDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("data-pipeline"),
			Role:             oss.Ptr("AliyunOSSDataPipelineRole"),
			DataPipelineConfiguration: &DataPipelineConfiguration{
				DataPipelineDescription: oss.Ptr("使用百炼多模态模型为业务数据向量化"),
				Sources: []DataPipelineSource{
					{
						InputBucket:    oss.Ptr("bucket"),
						InputDataScope: oss.Ptr("All"),
						FilterConfiguration: &DataPipelineSourceFilterConfiguration{
							PrefixSet:        []string{"prefix1", "prefix2/prefix3"},
							ObjectMediaTypes: []string{"text", "image", "video"},
						},
					},
				},
				DataPipelineEmbeddingConfiguration: &DataPipelineEmbeddingConfiguration{
					ApiKey:            oss.Ptr("your_api_key"),
					EmbeddingProvider: oss.Ptr("bailian"),
					FPS:               oss.Ptr(float64(1)),
					Model:             oss.Ptr("qwen2.5-vl-embedding"),
				},
				Destination: &DataPipelineDestination{
					VectorBucketName:    oss.Ptr("my-vector-bucket"),
					VectorIndexNames:    []string{"my-index"},
					VectorKeyPrefix:     oss.Ptr("prefix"),
					ObjectTagToMetadata: []string{"key1", "key2"},
					UsermetaToMetadata:  []string{"x-oss-meta-key1"},
				},
				DataPipelineError: &DataPipelineError{
					ErrorMode:   oss.Ptr("ignoreAndRecord"),
					ErrorBucket: oss.Ptr("my-error-bucket"),
					ErrorPrefix: oss.Ptr("error-output/"),
				},
			},
		},
		func(t *testing.T, o *PutDataPipelineConfigurationResult, err error) {
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
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=putDataPipelineConfiguration&dataPipeline&dataPipelineName=data-pipeline-v2&role=AliyunOSSDataPipelineRole", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<DataPipelineConfiguration><Sources><InputBucket>source-bucket</InputBucket><InputDataScope>All</InputDataScope><IgnoreDelete>false</IgnoreDelete></Sources><Destination><ObjectTagToMetadata>category</ObjectTagToMetadata><UsermetaToMetadata>x-oss-meta-source</UsermetaToMetadata><ImageEmbedding><Bucket>vector-bucket</Bucket><IndexName>image</IndexName><Prefix>v2</Prefix></ImageEmbedding><ImageTextEmbedding><Bucket>vector-bucket</Bucket><IndexName>image-text</IndexName><Prefix>v2</Prefix></ImageTextEmbedding><VideoFrameEmbedding><Bucket>vector-bucket</Bucket><IndexName>video-frame</IndexName><Prefix>v2</Prefix></VideoFrameEmbedding><VideoTextEmbedding><Bucket>vector-bucket</Bucket><IndexName>video-text</IndexName><Prefix>v2</Prefix></VideoTextEmbedding><DocumentChunkEmbedding><Bucket>vector-bucket</Bucket><IndexName>document</IndexName><Prefix>v2</Prefix></DocumentChunkEmbedding></Destination><ModelTier>standard</ModelTier><DataPipelineDataProcessConfiguration><SearchMode>balanced</SearchMode><Insights><Image><Caption><Prompt>Describe the image.</Prompt></Caption></Image><Video><Caption><Prompt>Describe each video scene.</Prompt></Caption><FrameEmbedding><Snapshot><Mode>interval</Mode><Interval>1</Interval></Snapshot></FrameEmbedding></Video></Insights></DataPipelineDataProcessConfiguration></DataPipelineConfiguration>")
		},
		&PutDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("data-pipeline-v2"),
			Role:             oss.Ptr("AliyunOSSDataPipelineRole"),
			DataPipelineConfiguration: &DataPipelineConfiguration{
				Sources: []DataPipelineSource{{
					InputBucket:    oss.Ptr("source-bucket"),
					InputDataScope: oss.Ptr("All"),
					IgnoreDelete:   oss.Ptr(false),
				}},
				ModelTier: oss.Ptr("standard"),
				DataPipelineDataProcessConfiguration: &DataPipelineDataProcessConfiguration{
					SearchMode: oss.Ptr("balanced"),
					Insights: &DataPipelineInsights{
						Image: &DataPipelineInsightsImage{Caption: &DataPipelineInsightsCaption{Prompt: oss.Ptr("Describe the image.")}},
						Video: &DataPipelineInsightsVideo{
							Caption: &DataPipelineInsightsCaption{Prompt: oss.Ptr("Describe each video scene.")},
							FrameEmbedding: &DataPipelineInsightsFrameEmbedding{Snapshot: &DataPipelineInsightsSnapshot{
								Mode:     oss.Ptr("interval"),
								Interval: oss.Ptr(1.0),
							}},
						},
					},
				},
				Destination: &DataPipelineDestination{
					ImageEmbedding:         &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("image"), Prefix: oss.Ptr("v2")},
					ImageTextEmbedding:     &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("image-text"), Prefix: oss.Ptr("v2")},
					VideoFrameEmbedding:    &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("video-frame"), Prefix: oss.Ptr("v2")},
					VideoTextEmbedding:     &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("video-text"), Prefix: oss.Ptr("v2")},
					DocumentChunkEmbedding: &VectorDestination{Bucket: oss.Ptr("vector-bucket"), IndexName: oss.Ptr("document"), Prefix: oss.Ptr("v2")},
					ObjectTagToMetadata:    []string{"category"},
					UsermetaToMetadata:     []string{"x-oss-meta-source"},
				},
			},
		},
		func(t *testing.T, o *PutDataPipelineConfigurationResult, err error) {
			assert.NoError(t, err)
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
}

func TestMockPutDataPipelineConfiguration_Success(t *testing.T) {
	for _, c := range testMockPutDataPipelineConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutDataPipelineConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutDataPipelineConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutDataPipelineConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *PutDataPipelineConfigurationResult, err error)
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
			assert.Equal(t, "/?action=putDataPipelineConfiguration&dataPipeline&dataPipelineName=data-pipeline&role=AliyunOSSDataPipelineRole", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<DataPipelineConfiguration><DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription><Sources><InputBucket>bucket</InputBucket><InputDataScope>All</InputDataScope><FilterConfiguration><PrefixSet>prefix1</PrefixSet><PrefixSet>prefix2/prefix3</PrefixSet><ObjectMediaTypes>text</ObjectMediaTypes><ObjectMediaTypes>image</ObjectMediaTypes><ObjectMediaTypes>video</ObjectMediaTypes></FilterConfiguration></Sources><DataPipelineEmbeddingConfiguration><EmbeddingProvider>bailian</EmbeddingProvider><ApiKey>your_api_key</ApiKey><Model>qwen2.5-vl-embedding</Model><FPS>1</FPS></DataPipelineEmbeddingConfiguration><Destination><VectorBucketName>my-vector-bucket</VectorBucketName><VectorKeyPrefix>prefix</VectorKeyPrefix><VectorIndexNames>my-index</VectorIndexNames><ObjectTagToMetadata>key1</ObjectTagToMetadata><ObjectTagToMetadata>key2</ObjectTagToMetadata><UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata></Destination><DataPipelineError><ErrorMode>ignoreAndRecord</ErrorMode><ErrorBucket>my-error-bucket</ErrorBucket><ErrorPrefix>error-output/</ErrorPrefix></DataPipelineError></DataPipelineConfiguration>")
		},
		&PutDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("data-pipeline"),
			Role:             oss.Ptr("AliyunOSSDataPipelineRole"),
			DataPipelineConfiguration: &DataPipelineConfiguration{
				DataPipelineDescription: oss.Ptr("使用百炼多模态模型为业务数据向量化"),
				Sources: []DataPipelineSource{
					{
						InputBucket:    oss.Ptr("bucket"),
						InputDataScope: oss.Ptr("All"),
						FilterConfiguration: &DataPipelineSourceFilterConfiguration{
							PrefixSet:        []string{"prefix1", "prefix2/prefix3"},
							ObjectMediaTypes: []string{"text", "image", "video"},
						},
					},
				},
				DataPipelineEmbeddingConfiguration: &DataPipelineEmbeddingConfiguration{
					ApiKey:            oss.Ptr("your_api_key"),
					EmbeddingProvider: oss.Ptr("bailian"),
					FPS:               oss.Ptr(float64(1)),
					Model:             oss.Ptr("qwen2.5-vl-embedding"),
				},
				Destination: &DataPipelineDestination{
					VectorBucketName:    oss.Ptr("my-vector-bucket"),
					VectorIndexNames:    []string{"my-index"},
					VectorKeyPrefix:     oss.Ptr("prefix"),
					ObjectTagToMetadata: []string{"key1", "key2"},
					UsermetaToMetadata:  []string{"x-oss-meta-key1"},
				},
				DataPipelineError: &DataPipelineError{
					ErrorMode:   oss.Ptr("ignoreAndRecord"),
					ErrorBucket: oss.Ptr("my-error-bucket"),
					ErrorPrefix: oss.Ptr("error-output/"),
				},
			},
		},
		func(t *testing.T, o *PutDataPipelineConfigurationResult, err error) {
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
			assert.Equal(t, "/?action=putDataPipelineConfiguration&dataPipeline&dataPipelineName=data-pipeline&role=AliyunOSSDataPipelineRole", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<DataPipelineConfiguration><DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription><Sources><InputBucket>bucket</InputBucket><InputDataScope>All</InputDataScope><FilterConfiguration><PrefixSet>prefix1</PrefixSet><PrefixSet>prefix2/prefix3</PrefixSet><ObjectMediaTypes>text</ObjectMediaTypes><ObjectMediaTypes>image</ObjectMediaTypes><ObjectMediaTypes>video</ObjectMediaTypes></FilterConfiguration></Sources><DataPipelineEmbeddingConfiguration><EmbeddingProvider>bailian</EmbeddingProvider><ApiKey>your_api_key</ApiKey><Model>qwen2.5-vl-embedding</Model><FPS>1</FPS></DataPipelineEmbeddingConfiguration><Destination><VectorBucketName>my-vector-bucket</VectorBucketName><VectorKeyPrefix>prefix</VectorKeyPrefix><VectorIndexNames>my-index</VectorIndexNames><ObjectTagToMetadata>key1</ObjectTagToMetadata><ObjectTagToMetadata>key2</ObjectTagToMetadata><UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata></Destination><DataPipelineError><ErrorMode>ignoreAndRecord</ErrorMode><ErrorBucket>my-error-bucket</ErrorBucket><ErrorPrefix>error-output/</ErrorPrefix></DataPipelineError></DataPipelineConfiguration>")
		},
		&PutDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("data-pipeline"),
			Role:             oss.Ptr("AliyunOSSDataPipelineRole"),
			DataPipelineConfiguration: &DataPipelineConfiguration{
				DataPipelineDescription: oss.Ptr("使用百炼多模态模型为业务数据向量化"),
				Sources: []DataPipelineSource{
					{
						InputBucket:    oss.Ptr("bucket"),
						InputDataScope: oss.Ptr("All"),
						FilterConfiguration: &DataPipelineSourceFilterConfiguration{
							PrefixSet:        []string{"prefix1", "prefix2/prefix3"},
							ObjectMediaTypes: []string{"text", "image", "video"},
						},
					},
				},
				DataPipelineEmbeddingConfiguration: &DataPipelineEmbeddingConfiguration{
					ApiKey:            oss.Ptr("your_api_key"),
					EmbeddingProvider: oss.Ptr("bailian"),
					FPS:               oss.Ptr(float64(1)),
					Model:             oss.Ptr("qwen2.5-vl-embedding"),
				},
				Destination: &DataPipelineDestination{
					VectorBucketName:    oss.Ptr("my-vector-bucket"),
					VectorIndexNames:    []string{"my-index"},
					VectorKeyPrefix:     oss.Ptr("prefix"),
					ObjectTagToMetadata: []string{"key1", "key2"},
					UsermetaToMetadata:  []string{"x-oss-meta-key1"},
				},
				DataPipelineError: &DataPipelineError{
					ErrorMode:   oss.Ptr("ignoreAndRecord"),
					ErrorBucket: oss.Ptr("my-error-bucket"),
					ErrorPrefix: oss.Ptr("error-output/"),
				},
			},
		},
		func(t *testing.T, o *PutDataPipelineConfigurationResult, err error) {
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

func TestMockPutDataPipelineConfiguration_Error(t *testing.T) {
	for _, c := range testMockPutDataPipelineConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutDataPipelineConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetDataPipelineConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetDataPipelineConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *GetDataPipelineConfigurationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8" ?>
<DataPipelineConfiguration>
    <DataPipelineName>my-data-pipeline</DataPipelineName>
    <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>
    <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>
    <Status>Running</Status>
    <Phase>IncrementalScanning</Phase>
    <Sources>
        <InputBucket>my-bucket</InputBucket>
        <InputDataScope>All</InputDataScope>
        <IgnoreDelete>true</IgnoreDelete>
        <FilterConfiguration>
            <PrefixSet>prefix1/</PrefixSet>
            <PrefixSet>prefix2/prefix3/</PrefixSet>
            <ObjectMediaTypes>text</ObjectMediaTypes>
            <ObjectMediaTypes>image</ObjectMediaTypes>
            <ObjectMediaTypes>video</ObjectMediaTypes>
        </FilterConfiguration>
    </Sources>
    <DataPipelineEmbeddingConfiguration>
        <EmbeddingProvider>bailian</EmbeddingProvider>
        <ApiKey>sk-12345678901234556</ApiKey>
        <Model>qwen2.5-vl-embedding</Model>
        <FPS>1</FPS>
    </DataPipelineEmbeddingConfiguration>
    <Destination>
        <VectorBucketName>my-vector-bucket</VectorBucketName>
        <VectorIndexNames>my-index</VectorIndexNames>
        <VectorKeyPrefix></VectorKeyPrefix>
        <ObjectTagToMetadata>key1</ObjectTagToMetadata>
        <ObjectTagToMetadata>key2</ObjectTagToMetadata>
        <UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata>
    </Destination>
    <DataPipelineError>
        <ErrorMode>ignoreAndRecord</ErrorMode>
        <ErrorBucket>my-error-bucket</ErrorBucket>
        <ErrorPrefix>error-output/</ErrorPrefix>
    </DataPipelineError>
    <CreateTime>2021-06-29T14:50:13.011643661+08:00</CreateTime>
</DataPipelineConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=getDataPipelineConfiguration&dataPipeline&dataPipelineName=my-data-pipeline", urlStr)
		},
		&GetDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *GetDataPipelineConfigurationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineDescription, "使用百炼多模态模型为业务数据向量化")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineName, "my-data-pipeline")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineRole, "my-data-pipeline-role")
			assert.Equal(t, *o.DataPipelineConfiguration.Status, "Running")
			assert.Equal(t, *o.DataPipelineConfiguration.Phase, "IncrementalScanning")
			assert.Equal(t, *o.DataPipelineConfiguration.Sources[0].InputBucket, "my-bucket")
			assert.Equal(t, *o.DataPipelineConfiguration.Sources[0].InputDataScope, "All")
			assert.Equal(t, *o.DataPipelineConfiguration.Sources[0].IgnoreDelete, true)
			assert.Equal(t, o.DataPipelineConfiguration.Sources[0].FilterConfiguration.PrefixSet[0], "prefix1/")
			assert.Equal(t, o.DataPipelineConfiguration.Sources[0].FilterConfiguration.PrefixSet[1], "prefix2/prefix3/")
			assert.Equal(t, o.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[0], "text")
			assert.Equal(t, o.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[1], "image")
			assert.Equal(t, o.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[2], "video")
			assert.Equal(t, o.DataPipelineConfiguration.Sources[0].FilterConfiguration.PrefixSet[1], "prefix2/prefix3/")
			assert.Equal(t, o.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[0], "text")
			assert.Equal(t, o.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[1], "image")
			assert.Equal(t, o.DataPipelineConfiguration.Sources[0].FilterConfiguration.ObjectMediaTypes[2], "video")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineEmbeddingConfiguration.EmbeddingProvider, "bailian")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineEmbeddingConfiguration.ApiKey, "sk-12345678901234556")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineEmbeddingConfiguration.Model, "qwen2.5-vl-embedding")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineEmbeddingConfiguration.FPS, float64(1))
			assert.Equal(t, *o.DataPipelineConfiguration.Destination.VectorBucketName, "my-vector-bucket")
			assert.Equal(t, o.DataPipelineConfiguration.Destination.VectorIndexNames[0], "my-index")
			assert.Equal(t, *o.DataPipelineConfiguration.Destination.VectorKeyPrefix, "")
			assert.Equal(t, o.DataPipelineConfiguration.Destination.ObjectTagToMetadata[0], "key1")
			assert.Equal(t, o.DataPipelineConfiguration.Destination.ObjectTagToMetadata[1], "key2")
			assert.Equal(t, o.DataPipelineConfiguration.Destination.UsermetaToMetadata[0], "x-oss-meta-key1")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineError.ErrorBucket, "my-error-bucket")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineError.ErrorPrefix, "error-output/")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineError.ErrorMode, "ignoreAndRecord")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineError.ErrorPrefix, "error-output/")
			assert.Equal(t, *o.DataPipelineConfiguration.DataPipelineError.ErrorMode, "ignoreAndRecord")
			assert.Equal(t, *o.DataPipelineConfiguration.CreateTime, "2021-06-29T14:50:13.011643661+08:00")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<DataPipelineConfiguration>
    <DataPipelineName>data-pipeline-v2</DataPipelineName>
    <ModelTier>standard</ModelTier>
    <DataPipelineDataProcessConfiguration>
        <SearchMode>balanced</SearchMode>
        <Insights>
            <Image><Caption><Prompt>Describe the image.</Prompt></Caption></Image>
            <Video><Caption><Prompt>Describe each video scene.</Prompt></Caption><FrameEmbedding><Snapshot><Mode>dhash</Mode><Number>10</Number></Snapshot></FrameEmbedding></Video>
        </Insights>
    </DataPipelineDataProcessConfiguration>
    <Destination>
        <ImageEmbedding><Bucket>vector-bucket</Bucket><IndexName>image</IndexName><Prefix>v2</Prefix></ImageEmbedding>
        <ImageTextEmbedding><Bucket>vector-bucket</Bucket><IndexName>image-text</IndexName><Prefix>v2</Prefix></ImageTextEmbedding>
        <VideoFrameEmbedding><Bucket>vector-bucket</Bucket><IndexName>video-frame</IndexName><Prefix>v2</Prefix></VideoFrameEmbedding>
        <VideoTextEmbedding><Bucket>vector-bucket</Bucket><IndexName>video-text</IndexName><Prefix>v2</Prefix></VideoTextEmbedding>
        <DocumentChunkEmbedding><Bucket>vector-bucket</Bucket><IndexName>document</IndexName><Prefix>v2</Prefix></DocumentChunkEmbedding>
        <ObjectTagToMetadata>category</ObjectTagToMetadata>
        <UsermetaToMetadata>x-oss-meta-source</UsermetaToMetadata>
    </Destination>
</DataPipelineConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=getDataPipelineConfiguration&dataPipeline&dataPipelineName=data-pipeline-v2", urlStr)
		},
		&GetDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("data-pipeline-v2"),
		},
		func(t *testing.T, o *GetDataPipelineConfigurationResult, err error) {
			assert.NoError(t, err)
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "standard", oss.ToString(o.DataPipelineConfiguration.ModelTier))
			assert.Equal(t, "balanced", oss.ToString(o.DataPipelineConfiguration.DataPipelineDataProcessConfiguration.SearchMode))
			assert.Equal(t, "Describe the image.", oss.ToString(o.DataPipelineConfiguration.DataPipelineDataProcessConfiguration.Insights.Image.Caption.Prompt))
			assert.Equal(t, "dhash", oss.ToString(o.DataPipelineConfiguration.DataPipelineDataProcessConfiguration.Insights.Video.FrameEmbedding.Snapshot.Mode))
			assert.Equal(t, int64(10), oss.ToInt64(o.DataPipelineConfiguration.DataPipelineDataProcessConfiguration.Insights.Video.FrameEmbedding.Snapshot.Number))
			assert.Equal(t, "image", oss.ToString(o.DataPipelineConfiguration.Destination.ImageEmbedding.IndexName))
			assert.Equal(t, "image-text", oss.ToString(o.DataPipelineConfiguration.Destination.ImageTextEmbedding.IndexName))
			assert.Equal(t, "video-frame", oss.ToString(o.DataPipelineConfiguration.Destination.VideoFrameEmbedding.IndexName))
			assert.Equal(t, "video-text", oss.ToString(o.DataPipelineConfiguration.Destination.VideoTextEmbedding.IndexName))
			assert.Equal(t, "document", oss.ToString(o.DataPipelineConfiguration.Destination.DocumentChunkEmbedding.IndexName))
			assert.Equal(t, []string{"category"}, o.DataPipelineConfiguration.Destination.ObjectTagToMetadata)
			assert.Equal(t, []string{"x-oss-meta-source"}, o.DataPipelineConfiguration.Destination.UsermetaToMetadata)
		},
	},
}

func TestMockGetDataPipelineConfiguration_Success(t *testing.T) {
	for _, c := range testMockGetDataPipelineConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetDataPipelineConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetDataPipelineConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetDataPipelineConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *GetDataPipelineConfigurationResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=getDataPipelineConfiguration&dataPipeline&dataPipelineName=my-data-pipeline", urlStr)
			assert.Equal(t, "POST", r.Method)
		},
		&GetDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *GetDataPipelineConfigurationResult, err error) {
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
			assert.Equal(t, "/?action=getDataPipelineConfiguration&dataPipeline&dataPipelineName=my-data-pipeline", strUrl)
		},
		&GetDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *GetDataPipelineConfigurationResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?action=getDataPipelineConfiguration&dataPipeline&dataPipelineName=my-data-pipeline", strUrl)
		},
		&GetDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *GetDataPipelineConfigurationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetDataPipelineConfiguration fail")
		},
	},
}

func TestMockGetDataPipelineConfiguration_Error(t *testing.T) {
	for _, c := range testMockGetDataPipelineConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetDataPipelineConfiguration(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteDataPipelineConfigurationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteDataPipelineConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *DeleteDataPipelineConfigurationResult, err error)
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
			assert.Equal(t, "/?action=deleteDataPipelineConfiguration&dataPipeline&dataPipelineName=my-data-pipeline", urlStr)
		},
		&DeleteDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *DeleteDataPipelineConfigurationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
}

func TestMockDeleteDataPipelineConfiguration_Success(t *testing.T) {
	for _, c := range testMockDeleteDataPipelineConfigurationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteDataPipelineConfiguration(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteDataPipelineConfigurationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteDataPipelineConfigurationRequest
	CheckOutputFn  func(t *testing.T, o *DeleteDataPipelineConfigurationResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=deleteDataPipelineConfiguration&dataPipeline&dataPipelineName=my-data-pipeline", urlStr)
			assert.Equal(t, "POST", r.Method)
		},
		&DeleteDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *DeleteDataPipelineConfigurationResult, err error) {
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
			assert.Equal(t, "/?action=deleteDataPipelineConfiguration&dataPipeline&dataPipelineName=my-data-pipeline", strUrl)
		},
		&DeleteDataPipelineConfigurationRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *DeleteDataPipelineConfigurationResult, err error) {
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

func TestMockDeleteDataPipelineConfiguration_Error(t *testing.T) {
	for _, c := range testMockDeleteDataPipelineConfigurationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteDataPipelineConfiguration(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockListDataPipelineConfigurationsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListDataPipelineConfigurationsRequest
	CheckOutputFn  func(t *testing.T, o *ListDataPipelineConfigurationsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8" ?>
<ListDataPipelineConfigurationsResult>
  <DataPipelineConfigurations>
    <DataPipelineConfiguration>
      <DataPipelineName>my-data-pipeline</DataPipelineName>
      <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>
      <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>
      <Status>Running</Status>
      <Phase>IncrementalScanning</Phase>
      <Sources>
          <InputBucket>my-bucket</InputBucket>
          <InputDataScope>All</InputDataScope>
          <IgnoreDelete>true</IgnoreDelete>
          <FilterConfiguration>
              <PrefixSet>prefix1/</PrefixSet>
              <PrefixSet>prefix2/prefix3/</PrefixSet>
              <ObjectMediaTypes>text</ObjectMediaTypes>
              <ObjectMediaTypes>image</ObjectMediaTypes>
              <ObjectMediaTypes>video</ObjectMediaTypes>
          </FilterConfiguration>
      </Sources>
      <DataPipelineEmbeddingConfiguration>
          <EmbeddingProvider>bailian</EmbeddingProvider>
          <ApiKey>sk-12345678901234556</ApiKey>
          <Model>qwen2.5-vl-embedding</Model>
          <FPS>1</FPS>
      </DataPipelineEmbeddingConfiguration>
      <Destination>
          <VectorBucketName>my-vector-bucket</VectorBucketName>
          <VectorIndexNames>my-index</VectorIndexNames>
          <VectorKeyPrefix></VectorKeyPrefix>
          <ObjectTagToMetadata>key1</ObjectTagToMetadata>
          <ObjectTagToMetadata>key2</ObjectTagToMetadata>
          <UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata>
      </Destination>
      <DataPipelineError>
          <ErrorMode>ignoreAndRecord</ErrorMode>
          <ErrorBucket>my-error-bucket</ErrorBucket>
          <ErrorPrefix>error-output/</ErrorPrefix>
      </DataPipelineError>
      <CreateTime>2021-06-29T14:50:13.011643661+08:00</CreateTime>
    </DataPipelineConfiguration>
  </DataPipelineConfigurations>
  <NextToken>xxx</NextToken>
</ListDataPipelineConfigurationsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=listDataPipelineConfigurations&dataPipeline&maxResults=100&nextToken=next-token&prefix=prefix", urlStr)
		},
		&ListDataPipelineConfigurationsRequest{
			MaxResults: oss.Ptr(int64(100)),
			NextToken:  oss.Ptr("next-token"),
			Prefix:     oss.Ptr("prefix"),
		},
		func(t *testing.T, o *ListDataPipelineConfigurationsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineDescription, "使用百炼多模态模型为业务数据向量化")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineName, "my-data-pipeline")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineRole, "my-data-pipeline-role")
			assert.Equal(t, *o.DataPipelineConfigurations[0].Status, "Running")
			assert.Equal(t, *o.DataPipelineConfigurations[0].Phase, "IncrementalScanning")
			assert.Equal(t, *o.DataPipelineConfigurations[0].Sources[0].InputBucket, "my-bucket")
			assert.Equal(t, *o.DataPipelineConfigurations[0].Sources[0].InputDataScope, "All")
			assert.Equal(t, *o.DataPipelineConfigurations[0].Sources[0].IgnoreDelete, true)
			assert.Equal(t, o.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.PrefixSet[0], "prefix1/")
			assert.Equal(t, o.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.PrefixSet[1], "prefix2/prefix3/")
			assert.Equal(t, o.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[0], "text")
			assert.Equal(t, o.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[1], "image")
			assert.Equal(t, o.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[2], "video")
			assert.Equal(t, o.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.PrefixSet[1], "prefix2/prefix3/")
			assert.Equal(t, o.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[0], "text")
			assert.Equal(t, o.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[1], "image")
			assert.Equal(t, o.DataPipelineConfigurations[0].Sources[0].FilterConfiguration.ObjectMediaTypes[2], "video")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineEmbeddingConfiguration.EmbeddingProvider, "bailian")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineEmbeddingConfiguration.ApiKey, "sk-12345678901234556")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineEmbeddingConfiguration.Model, "qwen2.5-vl-embedding")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineEmbeddingConfiguration.FPS, float64(1))
			assert.Equal(t, *o.DataPipelineConfigurations[0].Destination.VectorBucketName, "my-vector-bucket")
			assert.Equal(t, o.DataPipelineConfigurations[0].Destination.VectorIndexNames[0], "my-index")
			assert.Equal(t, *o.DataPipelineConfigurations[0].Destination.VectorKeyPrefix, "")
			assert.Equal(t, o.DataPipelineConfigurations[0].Destination.ObjectTagToMetadata[0], "key1")
			assert.Equal(t, o.DataPipelineConfigurations[0].Destination.ObjectTagToMetadata[1], "key2")
			assert.Equal(t, o.DataPipelineConfigurations[0].Destination.UsermetaToMetadata[0], "x-oss-meta-key1")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineError.ErrorBucket, "my-error-bucket")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineError.ErrorPrefix, "error-output/")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineError.ErrorMode, "ignoreAndRecord")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineError.ErrorPrefix, "error-output/")
			assert.Equal(t, *o.DataPipelineConfigurations[0].DataPipelineError.ErrorMode, "ignoreAndRecord")
			assert.Equal(t, *o.DataPipelineConfigurations[0].CreateTime, "2021-06-29T14:50:13.011643661+08:00")
			assert.Equal(t, *o.NextToken, "xxx")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<ListDataPipelineConfigurationsResult>
  <DataPipelineConfigurations>
    <DataPipelineConfiguration>
      <DataPipelineName>data-pipeline-v2</DataPipelineName>
      <ModelTier>standard</ModelTier>
      <DataPipelineDataProcessConfiguration><SearchMode>fast</SearchMode></DataPipelineDataProcessConfiguration>
      <Destination><ImageEmbedding><Bucket>vector-bucket</Bucket><IndexName>image</IndexName></ImageEmbedding></Destination>
    </DataPipelineConfiguration>
  </DataPipelineConfigurations>
  <NextToken>next-token</NextToken>
</ListDataPipelineConfigurationsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=listDataPipelineConfigurations&dataPipeline&inputBucket=source-bucket&maxResults=0&nextToken=previous-token&prefix=data-pipeline", urlStr)
		},
		&ListDataPipelineConfigurationsRequest{
			InputBucket: oss.Ptr("source-bucket"),
			MaxResults:  oss.Ptr(int64(0)),
			NextToken:   oss.Ptr("previous-token"),
			Prefix:      oss.Ptr("data-pipeline"),
		},
		func(t *testing.T, o *ListDataPipelineConfigurationsResult, err error) {
			assert.NoError(t, err)
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "next-token", oss.ToString(o.NextToken))
			assert.Len(t, o.DataPipelineConfigurations, 1)
			assert.Equal(t, "standard", oss.ToString(o.DataPipelineConfigurations[0].ModelTier))
			assert.Equal(t, "fast", oss.ToString(o.DataPipelineConfigurations[0].DataPipelineDataProcessConfiguration.SearchMode))
			assert.Equal(t, "image", oss.ToString(o.DataPipelineConfigurations[0].Destination.ImageEmbedding.IndexName))
		},
	},
}

func TestMockListDataPipelineConfigurations_Success(t *testing.T) {
	for _, c := range testMockListDataPipelineConfigurationsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListDataPipelineConfigurations(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListDataPipelineConfigurationsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListDataPipelineConfigurationsRequest
	CheckOutputFn  func(t *testing.T, o *ListDataPipelineConfigurationsResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=listDataPipelineConfigurations&dataPipeline", urlStr)
			assert.Equal(t, "POST", r.Method)
		},
		&ListDataPipelineConfigurationsRequest{},
		func(t *testing.T, o *ListDataPipelineConfigurationsResult, err error) {
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
			assert.Equal(t, "/?action=listDataPipelineConfigurations&dataPipeline", strUrl)
		},
		&ListDataPipelineConfigurationsRequest{},
		func(t *testing.T, o *ListDataPipelineConfigurationsResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?action=listDataPipelineConfigurations&dataPipeline", strUrl)
		},
		&ListDataPipelineConfigurationsRequest{},
		func(t *testing.T, o *ListDataPipelineConfigurationsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListDataPipelineConfigurations fail")
		},
	},
}

func TestMockListDataPipelineConfigurations_Error(t *testing.T) {
	for _, c := range testMockListDataPipelineConfigurationsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListDataPipelineConfigurations(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockPauseDataPipelineSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PauseDataPipelineRequest
	CheckOutputFn  func(t *testing.T, o *PauseDataPipelineResult, err error)
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
			assert.Equal(t, "/bucket/?action=pauseDataPipeline&dataPipeline&dataPipelineName=my-data-pipeline", urlStr)
		},
		&PauseDataPipelineRequest{
			Bucket:           oss.Ptr("bucket"),
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *PauseDataPipelineResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
}

func TestMockPauseDataPipeline_Success(t *testing.T) {
	for _, c := range testMockPauseDataPipelineSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PauseDataPipeline(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPauseDataPipelineErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PauseDataPipelineRequest
	CheckOutputFn  func(t *testing.T, o *PauseDataPipelineResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?action=pauseDataPipeline&dataPipeline&dataPipelineName=my-data-pipeline", urlStr)
			assert.Equal(t, "POST", r.Method)
		},
		&PauseDataPipelineRequest{
			Bucket:           oss.Ptr("bucket"),
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *PauseDataPipelineResult, err error) {
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
			assert.Equal(t, "/bucket/?action=pauseDataPipeline&dataPipeline&dataPipelineName=my-data-pipeline", strUrl)
		},
		&PauseDataPipelineRequest{
			Bucket:           oss.Ptr("bucket"),
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *PauseDataPipelineResult, err error) {
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

func TestMockPauseDataPipeline_Error(t *testing.T) {
	for _, c := range testMockPauseDataPipelineErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PauseDataPipeline(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockRestartDataPipelineSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RestartDataPipelineRequest
	CheckOutputFn  func(t *testing.T, o *RestartDataPipelineResult, err error)
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
			assert.Equal(t, "/?action=restartDataPipeline&dataPipeline&dataPipelineName=my-data-pipeline", urlStr)
		},
		&RestartDataPipelineRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *RestartDataPipelineResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
}

func TestMockRestartDataPipeline_Success(t *testing.T) {
	for _, c := range testMockRestartDataPipelineSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RestartDataPipeline(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockRestartDataPipelineErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RestartDataPipelineRequest
	CheckOutputFn  func(t *testing.T, o *RestartDataPipelineResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=restartDataPipeline&dataPipeline&dataPipelineName=my-data-pipeline", urlStr)
			assert.Equal(t, "POST", r.Method)
		},
		&RestartDataPipelineRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *RestartDataPipelineResult, err error) {
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
			assert.Equal(t, "/?action=restartDataPipeline&dataPipeline&dataPipelineName=my-data-pipeline", strUrl)
		},
		&RestartDataPipelineRequest{
			DataPipelineName: oss.Ptr("my-data-pipeline"),
		},
		func(t *testing.T, o *RestartDataPipelineResult, err error) {
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

func TestMockRestartDataPipeline_Error(t *testing.T) {
	for _, c := range testMockRestartDataPipelineErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RestartDataPipeline(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}
