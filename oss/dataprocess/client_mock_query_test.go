package dataprocess

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockSimpleQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SimpleQueryRequest
	CheckOutputFn  func(t *testing.T, o *SimpleQueryResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
  <NextToken>MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****</NextToken>
  <TotalHits>150</TotalHits>
  <Files>
    <File>
      <Filename>docs/report.pdf</Filename>
      <Size>5242880</Size>
      <URI>oss://examplebucket/docs/report.pdf</URI>
      <OSSURI>oss://examplebucket/docs/report.pdf</OSSURI>
      <MediaType>document</MediaType>
      <ContentType>application/pdf</ContentType>
      <FileModifiedTime>2025-12-01T10:30:00Z</FileModifiedTime>
      <PageCount>20</PageCount>
    </File>
  </Files>
  <Aggregations>
    <Aggregation>
      <Field>MediaType</Field>
      <Operation>group</Operation>
      <Groups>
        <AggregationGroup>
          <Value>document</Value>
          <Count>80</Count>
        </AggregationGroup>
        <AggregationGroup>
          <Value>image</Value>
          <Count>70</Count>
        </AggregationGroup>
      </Groups>
    </Aggregation>
  </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=simpleQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SimpleQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SimpleQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.NextToken, "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****")
			assert.Equal(t, *o.TotalHits, int64(150))
			assert.Equal(t, len(o.Files), 1)
			assert.Equal(t, *o.Files[0].Filename, "docs/report.pdf")
			assert.Equal(t, *o.Files[0].Size, int64(5242880))
			assert.Equal(t, *o.Files[0].URI, "oss://examplebucket/docs/report.pdf")
			assert.Equal(t, *o.Files[0].OSSURI, "oss://examplebucket/docs/report.pdf")
			assert.Equal(t, *o.Files[0].MediaType, "document")
			assert.Equal(t, *o.Files[0].ContentType, "application/pdf")
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2025-12-01T10:30:00Z")
			assert.Equal(t, *o.Files[0].PageCount, int64(20))
			assert.Equal(t, len(o.Aggregations), 1)
			assert.Equal(t, *o.Aggregations[0].Field, "MediaType")
			assert.Equal(t, *o.Aggregations[0].Operation, "group")
			assert.Equal(t, len(o.Aggregations[0].AggregationGroups), 2)
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Value, "document")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Count, int64(80))
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Value, "image")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Count, int64(70))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
  <NextToken>MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****</NextToken>
  <TotalHits>258</TotalHits>
  <Files>
    <File>
      <Filename>photos/sunset.jpg</Filename>
      <Size>2048000</Size>
      <URI>oss://examplebucket/photos/sunset.jpg</URI>
      <OSSURI>oss://examplebucket/photos/sunset.jpg</OSSURI>
      <MediaType>image</MediaType>
      <ContentType>image/jpeg</ContentType>
      <FileModifiedTime>2025-12-01T10:30:00Z</FileModifiedTime>
      <ImageWidth>3840</ImageWidth>
      <ImageHeight>2160</ImageHeight>
      <Orientation>1</Orientation>
    </File>
    <File>
      <Filename>photos/mountain.png</Filename>
      <Size>5120000</Size>
      <URI>oss://examplebucket/photos/mountain.png</URI>
      <OSSURI>oss://examplebucket/photos/mountain.png</OSSURI>
      <MediaType>image</MediaType>
      <ContentType>image/png</ContentType>
      <FileModifiedTime>2025-11-20T14:00:00Z</FileModifiedTime>
      <ImageWidth>1920</ImageWidth>
      <ImageHeight>1080</ImageHeight>
      <Orientation>1</Orientation>
    </File>
  </Files>
  <Aggregations>
    <Aggregation>
      <Field>MediaType</Field>
      <Operation>group</Operation>
      <Groups>
        <AggregationGroup>
          <Value>image</Value>
          <Count>200</Count>
        </AggregationGroup>
        <AggregationGroup>
          <Value>video</Value>
          <Count>58</Count>
        </AggregationGroup>
      </Groups>
    </Aggregation>
  </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=simpleQuery&aggregations=Size&datasetName=your_dataset&maxResults=99&metaQuery&nextToken=MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw%2A%2A%2A%2A&order=acs&query=%7B%22Field%22%3A+%22Size%22%2C%22Value%22%3A+%221%22%2C%22Operation%22%3A+%22gt%22%7D&sort=Size&withFields=%5B%22Filename%22%2C%22Size%22%5D&withoutTotalHits=true", strUrl)
		},
		&SimpleQueryRequest{
			Bucket:           oss.Ptr("bucket"),
			DatasetName:      oss.Ptr("your_dataset"),
			NextToken:        oss.Ptr("MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****"),
			MaxResults:       oss.Ptr(int32(99)),
			Query:            oss.Ptr("{\"Field\": \"Size\",\"Value\": \"1\",\"Operation\": \"gt\"}"),
			Sort:             oss.Ptr("Size"),
			Order:            oss.Ptr("acs"),
			Aggregations:     oss.Ptr("Size"),
			WithFields:       oss.Ptr(`["Filename","Size"]`),
			WithoutTotalHits: oss.Ptr(true),
		},
		func(t *testing.T, o *SimpleQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.NextToken, "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****")
			assert.Equal(t, *o.TotalHits, int64(258))
			assert.Equal(t, len(o.Files), 2)
			assert.Equal(t, *o.Files[0].Filename, "photos/sunset.jpg")
			assert.Equal(t, *o.Files[0].Size, int64(2048000))
			assert.Equal(t, *o.Files[0].URI, "oss://examplebucket/photos/sunset.jpg")
			assert.Equal(t, *o.Files[0].OSSURI, "oss://examplebucket/photos/sunset.jpg")
			assert.Equal(t, *o.Files[0].MediaType, "image")
			assert.Equal(t, *o.Files[0].ContentType, "image/jpeg")
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2025-12-01T10:30:00Z")
			assert.Equal(t, *o.Files[0].ImageWidth, int64(3840))
			assert.Equal(t, *o.Files[0].ImageHeight, int64(2160))
			assert.Equal(t, *o.Files[0].Orientation, int64(1))

			assert.Equal(t, *o.Files[1].Filename, "photos/mountain.png")
			assert.Equal(t, *o.Files[1].Size, int64(5120000))
			assert.Equal(t, *o.Files[1].URI, "oss://examplebucket/photos/mountain.png")
			assert.Equal(t, *o.Files[1].OSSURI, "oss://examplebucket/photos/mountain.png")
			assert.Equal(t, *o.Files[1].MediaType, "image")
			assert.Equal(t, *o.Files[1].ContentType, "image/png")
			assert.Equal(t, *o.Files[1].FileModifiedTime, "2025-11-20T14:00:00Z")
			assert.Equal(t, *o.Files[1].ImageWidth, int64(1920))
			assert.Equal(t, *o.Files[1].ImageHeight, int64(1080))
			assert.Equal(t, *o.Files[1].Orientation, int64(1))
			assert.Equal(t, len(o.Aggregations), 1)
			assert.Equal(t, *o.Aggregations[0].Field, "MediaType")
			assert.Equal(t, *o.Aggregations[0].Operation, "group")
			assert.Equal(t, len(o.Aggregations[0].AggregationGroups), 2)
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Value, "image")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Count, int64(200))
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Value, "video")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Count, int64(58))
		},
	},
}

func TestMockSimpleQuery_Success(t *testing.T) {
	for _, c := range testMockSimpleQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SimpleQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSimpleQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SimpleQueryRequest
	CheckOutputFn  func(t *testing.T, o *SimpleQueryResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=simpleQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SimpleQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SimpleQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=simpleQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SimpleQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SimpleQueryResult, err error) {
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

func TestMockSimpleQuery_Error(t *testing.T) {
	for _, c := range testMockSimpleQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SimpleQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSemanticQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SemanticQueryRequest
	CheckOutputFn  func(t *testing.T, o *SemanticQueryResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
<Files>
      <File>
          <Addresses/>
          <AudioCovers/>
          <AudioStreams>
              <AudioStream>
                  <Bitrate>128000</Bitrate>
                  <ChannelLayout>stereo</ChannelLayout>
                  <Channels>2</Channels>
                  <CodecLongName>AAC (Advanced Audio Coding)</CodecLongName>
                  <CodecName>aac</CodecName>
                  <CodecTag>0x6134706d</CodecTag>
                  <CodecTagString>mp4a</CodecTagString>
                  <Duration>16.021769</Duration>
                  <FrameCount>690</FrameCount>
                  <Index>1</Index>
                  <SampleFormat>fltp</SampleFormat>
                  <SampleRate>44100</SampleRate>
                  <TimeBase>1/44100</TimeBase>
              </AudioStream>
          </AudioStreams>
          <Bitrate>1656706</Bitrate>
          <ContentMd5>5oJccWuBoqVXS8zrzckPlg==</ContentMd5>
          <ContentType>video/mp4</ContentType>
          <CreateTime>2026-04-21T20:28:17.018858947+08:00</CreateTime>
          <CroppingSuggestions/>
          <DatasetName>test-dataset-sem-vid-1776774492</DatasetName>
          <Duration>16.034</Duration>
          <ETag>\"E6825C716B81A2A5574BCCEBCDC90F96\"</ETag>
          <Elements/>
          <Figures/>
          <FileHash>E6825C716B81A2A5574BCCEBCDC90F96</FileHash>
          <FileModifiedTime>2026-04-21T20:28:13+08:00</FileModifiedTime>
          <Filename>test-temp/sem-vid-1776774492774503000.mp4</Filename>
          <FormatLongName>QuickTime / MOV</FormatLongName>
          <FormatName>mov,mp4,m4a,3gp,3g2,mj2</FormatName>
          <Insights>
              <Video>
                  <Caption>蓝衣男走向餐桌</Caption>
                  <Description>这是一段室内高角度监控录像，场景为一个客厅。</Description>
              </Video>
          </Insights>
          <Labels/>
          <MediaType>video</MediaType>
          <OCRContents/>
          <OSSCRC64>2327801188977127298</OSSCRC64>
          <OSSObjectType>Normal</OSSObjectType>
          <OSSStorageClass>Standard</OSSStorageClass>
          <OSSTagging>
			  <Tagging>
				  <Key>routing-dataset</Key>
				  <Value>test-dataset-sem-vid-1776774492</Value>
			  </Tagging>
          </OSSTagging>
          <OSSTaggingCount>1</OSSTaggingCount>
          <ObjectACL>default</ObjectACL>
          <SequenceNumber>2</SequenceNumber>
          <SemanticSimilarity>0.5583347777557373</SemanticSimilarity>
          <Size>3320455</Size>
          <SmartClusters/>
          <StreamCount>2</StreamCount>
          <Subtitles/>
          <URI>oss://oss-metaquery-dataset-test/test-temp/sem-vid-1776774492774503000.mp4</URI>
          <UpdateTime>2026-04-21T20:28:27.359034257+08:00</UpdateTime>
          <VideoHeight>1080</VideoHeight>
          <VideoStreams>
              <VideoStream>
                  <AverageFrameRate>21645000/721493</AverageFrameRate>
                  <BitDepth>8</BitDepth>
                  <Bitrate>1521221</Bitrate>
                  <CodecLongName>H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10</CodecLongName>
                  <CodecName>h264</CodecName>
                  <CodecTag>0x31637661</CodecTag>
                  <CodecTagString>avc1</CodecTagString>
                  <ColorPrimaries>bt709</ColorPrimaries>
                  <ColorRange>tv</ColorRange>
                  <ColorSpace>bt709</ColorSpace>
                  <ColorTransfer>bt709</ColorTransfer>
                  <DisplayAspectRatio>16:9</DisplayAspectRatio>
                  <Duration>16.033178</Duration>
                  <FrameCount>481</FrameCount>
                  <FrameRate>90000/2999</FrameRate>
                  <Height>1080</Height>
                  <Level>31</Level>
                  <PixelFormat>yuv420p</PixelFormat>
                  <Profile>High</Profile>
                  <SampleAspectRatio>1:1</SampleAspectRatio>
                  <TimeBase>1/90000</TimeBase>
                  <Width>1920</Width>
              </VideoStream>
          </VideoStreams>
          <VideoWidth>1920</VideoWidth>
      </File>
  </Files>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=semanticQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SemanticQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SemanticQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, len(o.Files), 1)
			assert.Equal(t, *o.Files[0].Bitrate, int64(1656706))
			assert.Equal(t, *o.Files[0].ContentMd5, "5oJccWuBoqVXS8zrzckPlg==")
			assert.Equal(t, *o.Files[0].ContentType, "video/mp4")
			assert.Equal(t, *o.Files[0].CreateTime, "2026-04-21T20:28:17.018858947+08:00")
			assert.Equal(t, *o.Files[0].DatasetName, "test-dataset-sem-vid-1776774492")
			assert.Equal(t, *o.Files[0].Duration, float64(16.034))
			assert.Equal(t, *o.Files[0].ETag, "\\\"E6825C716B81A2A5574BCCEBCDC90F96\\\"")
			assert.Equal(t, *o.Files[0].FileHash, "E6825C716B81A2A5574BCCEBCDC90F96")
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2026-04-21T20:28:13+08:00")
			assert.Equal(t, *o.Files[0].Filename, "test-temp/sem-vid-1776774492774503000.mp4")
			assert.Equal(t, *o.Files[0].FormatLongName, "QuickTime / MOV")
			assert.Equal(t, *o.Files[0].FormatName, "mov,mp4,m4a,3gp,3g2,mj2")
			assert.Equal(t, *o.Files[0].MediaType, "video")
			assert.Equal(t, *o.Files[0].Size, int64(3320455))
			assert.Equal(t, *o.Files[0].VideoWidth, int64(1920))
			assert.Equal(t, *o.Files[0].VideoHeight, int64(1080))
			assert.Equal(t, *o.Files[0].StreamCount, int64(2))
			assert.Equal(t, *o.Files[0].OSSObjectType, "Normal")
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, *o.Files[0].OSSTaggingCount, int64(1))
			assert.Equal(t, *o.Files[0].OSSTagging[0].Key, "routing-dataset")
			assert.Equal(t, *o.Files[0].OSSTagging[0].Value, "test-dataset-sem-vid-1776774492")
			assert.Equal(t, *o.Files[0].ObjectACL, "default")
			assert.Equal(t, *o.Files[0].SequenceNumber, int64(2))
			assert.Equal(t, *o.Files[0].SemanticSimilarity, float64(0.5583347777557373))
			assert.Equal(t, *o.Files[0].Size, int64(3320455))
			assert.Equal(t, *o.Files[0].URI, "oss://oss-metaquery-dataset-test/test-temp/sem-vid-1776774492774503000.mp4")
			assert.Equal(t, *o.Files[0].UpdateTime, "2026-04-21T20:28:27.359034257+08:00")

			assert.Equal(t, len(o.Files[0].AudioStreams), 1)
			assert.Equal(t, *o.Files[0].AudioStreams[0].Bitrate, int64(128000))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Channels, int64(2))
			assert.Equal(t, *o.Files[0].AudioStreams[0].ChannelLayout, "stereo")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecLongName, "AAC (Advanced Audio Coding)")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecName, "aac")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecTag, "0x6134706d")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecTagString, "mp4a")
			assert.Equal(t, *o.Files[0].AudioStreams[0].Duration, float64(16.021769))
			assert.Equal(t, *o.Files[0].AudioStreams[0].FrameCount, int64(690))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Index, int64(1))
			assert.Equal(t, *o.Files[0].AudioStreams[0].SampleFormat, "fltp")
			assert.Equal(t, *o.Files[0].AudioStreams[0].SampleRate, int64(44100))
			assert.Equal(t, *o.Files[0].AudioStreams[0].TimeBase, "1/44100")
			assert.Equal(t, *o.Files[0].Insights.Video.Description, "这是一段室内高角度监控录像，场景为一个客厅。")
			assert.Equal(t, *o.Files[0].Insights.Video.Caption, "蓝衣男走向餐桌")
			assert.Equal(t, len(o.Files[0].VideoStreams), 1)
			assert.Equal(t, *o.Files[0].VideoStreams[0].AverageFrameRate, "21645000/721493")
			assert.Equal(t, *o.Files[0].VideoStreams[0].BitDepth, int64(8))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Bitrate, int64(1521221))
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecLongName, "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecName, "h264")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecTag, "0x31637661")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecTagString, "avc1")
			assert.Equal(t, *o.Files[0].VideoStreams[0].ColorPrimaries, "bt709")
			assert.Equal(t, *o.Files[0].VideoStreams[0].ColorRange, "tv")
			assert.Equal(t, *o.Files[0].VideoStreams[0].ColorTransfer, "bt709")
			assert.Equal(t, *o.Files[0].VideoStreams[0].ColorSpace, "bt709")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Duration, float64(16.033178))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameCount, int64(481))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameRate, "90000/2999")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Height, int64(1080))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Width, int64(1920))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Level, int64(31))
			assert.Equal(t, *o.Files[0].VideoStreams[0].PixelFormat, "yuv420p")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Profile, "High")
			assert.Equal(t, *o.Files[0].VideoStreams[0].PixelFormat, "yuv420p")
			assert.Equal(t, *o.Files[0].VideoStreams[0].SampleAspectRatio, "1:1")
			assert.Equal(t, *o.Files[0].VideoStreams[0].TimeBase, "1/90000")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
    <Files>
        <File>
            <Addresses/>
            <AudioCovers/>
            <AudioStreams>
                <AudioStream>
                    <Bitrate>14983</Bitrate>
                    <ChannelLayout>mono</ChannelLayout>
                    <Channels>1</Channels>
                    <CodecLongName>AAC (Advanced Audio Coding)</CodecLongName>
                    <CodecName>aac</CodecName>
                    <CodecTag>0x6134706d</CodecTag>
                    <CodecTagString>mp4a</CodecTagString>
                    <Duration>7.936</Duration>
                    <FrameCount>62</FrameCount>
                    <Index>1</Index>
                    <SampleFormat>fltp</SampleFormat>
                    <SampleRate>8000</SampleRate>
                    <TimeBase>1/8000</TimeBase>
                </AudioStream>
            </AudioStreams>
            <Bitrate>196284</Bitrate>
            <ContentMd5>5/ZLrWYXpuQfDfxEf4+lyA==</ContentMd5>
            <ContentType>video/mp4</ContentType>
            <CreateTime>2026-04-21T10:51:38.264045621+08:00</CreateTime>
            <CroppingSuggestions/>
            <DatasetName>dataset-aianalysis-walk</DatasetName>
            <Duration>8</Duration>
            <ETag>\"E7F64BAD6617A6E41F0DFC447F8FA5C8\"</ETag>
            <Elements/>
            <Figures/>
            <FileHash>E7F64BAD6617A6E41F0DFC447F8FA5C8</FileHash>
            <FileModifiedTime>2026-04-21T10:51:25+08:00</FileModifiedTime>
            <Filename>mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4</Filename>
            <FormatLongName>QuickTime / MOV</FormatLongName>
            <FormatName>mov,mp4,m4a,3gp,3g2,mj2</FormatName>
            <Labels>
                <Label>
                    <LabelConfidence>1</LabelConfidence>
                    <LabelName>有人走过</LabelName>
                    <ParentLabelName>自定义标签</ParentLabelName>
                    <Clips>
                        <Clip>
                            <TimeRange>200</TimeRange>
                            <TimeRange>5533</TimeRange>
                        </Clip>
                    </Clips>
                </Label>
            </Labels>
            <MediaType>video</MediaType>
            <OCRContents/>
            <OSSCRC64>16628192875747293357</OSSCRC64>
            <OSSObjectType>Normal</OSSObjectType>
            <OSSStorageClass>Standard</OSSStorageClass>
            <OSSTagging>
				<Tagging>
					<Key>routing-dataset</Key>
					<Value>photos-2026</Value>
				</Tagging>
				<Tagging>
					<Key>env</Key>
					<Value>production</Value>
				</Tagging>
            </OSSTagging>
            <OSSTaggingCount>2</OSSTaggingCount>
            <OSSUserMeta>
				<UserMeta>
					<Key>category</Key>
					<Value>photo</Value>
				</UserMeta>
				<UserMeta>
					<Key>album</Key>
					<Value>vacation</Value>
				</UserMeta>
            </OSSUserMeta>
 			<CustomLabels>
				<Item>
					<Key>label-ley</Key>
					<Value>label-val</Value>
				</Item>
			</CustomLabels>
            <ObjectACL>default</ObjectACL>
            <ProduceTime>2026-04-21T10:46:10+08:00</ProduceTime>
            <SceneElements>
                <SceneElement>
                    <FrameTimes>6000</FrameTimes>
                    <TimeRange>4133</TimeRange>
                    <TimeRange>8533</TimeRange>
                    <VideoStreamIndex>0</VideoStreamIndex>
                    <Labels/>
                </SceneElement>
            </SceneElements>
            <SemanticSimilarity>0.2536</SemanticSimilarity>
            <SequenceNumber>5</SequenceNumber>
            <Size>196284</Size>
            <SmartClusters/>
            <StreamCount>2</StreamCount>
            <Subtitles/>
            <URI>oss://paas-smart-cloud-test/mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4</URI>
            <UpdateTime>2026-04-21T10:52:39.412605575+08:00</UpdateTime>
            <VideoHeight>360</VideoHeight>
            <VideoStreams>
                <VideoStream>
                    <AverageFrameRate>15/1</AverageFrameRate>
                    <BitDepth>8</BitDepth>
                    <Bitrate>178202</Bitrate>
                    <CodecLongName>H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10</CodecLongName>
                    <CodecName>h264</CodecName>
                    <CodecTag>0x31637661</CodecTag>
                    <CodecTagString>avc1</CodecTagString>
                    <Duration>8</Duration>
                    <FrameCount>120</FrameCount>
                    <FrameRate>500/33</FrameRate>
                    <Height>360</Height>
                    <Level>22</Level>
                    <PixelFormat>yuv420p</PixelFormat>
                    <Profile>Main</Profile>
                    <TimeBase>1/1000</TimeBase>
                    <Width>640</Width>
                </VideoStream>
            </VideoStreams>
            <VideoWidth>640</VideoWidth>
        </File>
    </Files>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=semanticQuery&datasetName=your_dataset&maxResults=10&mediaTypes=%5B%22video%22%2C%22image%22%5D&metaQuery&query=%7B%22Field%22%3A+%22Size%22%2C%22Value%22%3A+%221%22%2C%22Operation%22%3A+%22gt%22%7D&sourceURI=oss%3A%2F%2Fbucket%2Fprefix&withFields=%5B%22Filename%22%2C%22Size%22%5D", strUrl)
		},
		&SemanticQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			MaxResults:  oss.Ptr(int32(10)),
			Query:       oss.Ptr("{\"Field\": \"Size\",\"Value\": \"1\",\"Operation\": \"gt\"}"),
			WithFields:  oss.Ptr(`["Filename","Size"]`),
			MediaTypes:  oss.Ptr(`["video","image"]`),
			SourceUri:   oss.Ptr("oss://bucket/prefix"),
		},
		func(t *testing.T, o *SemanticQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Files), 1)
			assert.Equal(t, *o.Files[0].Bitrate, int64(196284))
			assert.Equal(t, *o.Files[0].ContentMd5, "5/ZLrWYXpuQfDfxEf4+lyA==")
			assert.Equal(t, *o.Files[0].ContentType, "video/mp4")
			assert.Equal(t, *o.Files[0].CreateTime, "2026-04-21T10:51:38.264045621+08:00")
			assert.Equal(t, *o.Files[0].DatasetName, "dataset-aianalysis-walk")
			assert.Equal(t, *o.Files[0].Duration, float64(8))
			assert.Equal(t, *o.Files[0].ETag, "\\\"E7F64BAD6617A6E41F0DFC447F8FA5C8\\\"")
			assert.Equal(t, *o.Files[0].FileHash, "E7F64BAD6617A6E41F0DFC447F8FA5C8")
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2026-04-21T10:51:25+08:00")
			assert.Equal(t, *o.Files[0].Filename, "mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4")
			assert.Equal(t, *o.Files[0].FormatLongName, "QuickTime / MOV")
			assert.Equal(t, *o.Files[0].FormatName, "mov,mp4,m4a,3gp,3g2,mj2")
			assert.Equal(t, *o.Files[0].MediaType, "video")
			assert.Equal(t, *o.Files[0].OSSCRC64, "16628192875747293357")
			assert.Equal(t, *o.Files[0].Size, int64(196284))
			assert.Equal(t, *o.Files[0].VideoWidth, int64(640))
			assert.Equal(t, *o.Files[0].VideoHeight, int64(360))
			assert.Equal(t, *o.Files[0].StreamCount, int64(2))
			assert.Equal(t, *o.Files[0].OSSObjectType, "Normal")
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, len(o.Files[0].OSSTagging), 2)
			assert.Equal(t, *o.Files[0].OSSTagging[0].Key, "routing-dataset")
			assert.Equal(t, *o.Files[0].OSSTagging[1].Key, "env")
			assert.Equal(t, *o.Files[0].OSSTagging[0].Value, "photos-2026")
			assert.Equal(t, *o.Files[0].OSSTagging[1].Value, "production")
			assert.Equal(t, len(o.Files[0].OSSUserMeta), 2)
			assert.Equal(t, *o.Files[0].OSSUserMeta[0].Key, "category")
			assert.Equal(t, *o.Files[0].OSSUserMeta[0].Value, "photo")
			assert.Equal(t, *o.Files[0].OSSUserMeta[1].Key, "album")
			assert.Equal(t, *o.Files[0].OSSUserMeta[1].Value, "vacation")
			assert.Equal(t, *o.Files[0].OSSTaggingCount, int64(2))
			assert.Equal(t, *o.Files[0].ObjectACL, "default")
			assert.Equal(t, len(o.Files[0].OSSUserMeta), 2)
			assert.Equal(t, *o.Files[0].OSSUserMeta[0].Key, "category")
			assert.Equal(t, *o.Files[0].OSSUserMeta[0].Value, "photo")
			assert.Equal(t, *o.Files[0].OSSUserMeta[1].Key, "album")
			assert.Equal(t, *o.Files[0].OSSUserMeta[1].Value, "vacation")
			assert.Equal(t, len(o.Files[0].CustomLabels), 1)
			assert.Equal(t, *o.Files[0].CustomLabels[0].Key, "label-ley")
			assert.Equal(t, *o.Files[0].CustomLabels[0].Value, "label-val")
			assert.Equal(t, *o.Files[0].ProduceTime, "2026-04-21T10:46:10+08:00")
			assert.Equal(t, *o.Files[0].SequenceNumber, int64(5))
			assert.Equal(t, *o.Files[0].SemanticSimilarity, float64(0.2536))
			assert.Equal(t, *o.Files[0].Size, int64(196284))
			assert.Equal(t, *o.Files[0].StreamCount, int64(2))
			assert.Equal(t, *o.Files[0].URI, "oss://paas-smart-cloud-test/mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4")
			assert.Equal(t, *o.Files[0].UpdateTime, "2026-04-21T10:52:39.412605575+08:00")

			assert.Equal(t, len(o.Files[0].Labels), 1)
			assert.Equal(t, *o.Files[0].Labels[0].LabelConfidence, float64(1))
			assert.Equal(t, *o.Files[0].Labels[0].LabelName, "有人走过")
			assert.Equal(t, *o.Files[0].Labels[0].ParentLabelName, "自定义标签")
			assert.Equal(t, len(o.Files[0].Labels[0].Clips), 1)
			assert.Equal(t, o.Files[0].Labels[0].Clips[0].TimeRange[0], int64(200))
			assert.Equal(t, o.Files[0].Labels[0].Clips[0].TimeRange[1], int64(5533))

			assert.Equal(t, len(o.Files[0].AudioStreams), 1)
			assert.Equal(t, *o.Files[0].AudioStreams[0].Bitrate, int64(14983))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Channels, int64(1))
			assert.Equal(t, *o.Files[0].AudioStreams[0].ChannelLayout, "mono")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecLongName, "AAC (Advanced Audio Coding)")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecName, "aac")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecTag, "0x6134706d")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecTagString, "mp4a")
			assert.Equal(t, *o.Files[0].AudioStreams[0].Duration, float64(7.936))
			assert.Equal(t, *o.Files[0].AudioStreams[0].FrameCount, int64(62))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Index, int64(1))
			assert.Equal(t, *o.Files[0].AudioStreams[0].SampleFormat, "fltp")
			assert.Equal(t, *o.Files[0].AudioStreams[0].SampleRate, int64(8000))
			assert.Equal(t, *o.Files[0].AudioStreams[0].TimeBase, "1/8000")

			assert.Equal(t, len(o.Files[0].VideoStreams), 1)
			assert.Equal(t, *o.Files[0].VideoStreams[0].AverageFrameRate, "15/1")
			assert.Equal(t, *o.Files[0].VideoStreams[0].BitDepth, int64(8))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Bitrate, int64(178202))
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecLongName, "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecName, "h264")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecTag, "0x31637661")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecTagString, "avc1")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Duration, float64(8))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameCount, int64(120))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameRate, "500/33")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Height, int64(360))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Width, int64(640))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Level, int64(22))
			assert.Equal(t, *o.Files[0].VideoStreams[0].TimeBase, "1/1000")
		},
	},
}

func TestMockSemanticQuery_Success(t *testing.T) {
	for _, c := range testMockSemanticQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SemanticQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSemanticQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SemanticQueryRequest
	CheckOutputFn  func(t *testing.T, o *SemanticQueryResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=semanticQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SemanticQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SemanticQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=semanticQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SemanticQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SemanticQueryResult, err error) {
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

func TestMockSemanticQuery_Error(t *testing.T) {
	for _, c := range testMockSemanticQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SemanticQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
