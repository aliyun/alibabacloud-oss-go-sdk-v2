package dataprocess

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_SimpleQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *SimpleQueryRequest
	var input *oss.OperationInput
	var err error

	request = &SimpleQueryRequest{}
	input = &oss.OperationInput{
		OpName: "SimpleQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "simpleQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &SimpleQueryRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "SimpleQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "simpleQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &SimpleQueryRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
	}
	input = &oss.OperationInput{
		OpName: "SimpleQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "simpleQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")

	request = &SimpleQueryRequest{
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
	}
	input = &oss.OperationInput{
		OpName: "SimpleQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "simpleQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["nextToken"], "MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****")
	assert.Equal(t, input.Parameters["maxResults"], "99")
	assert.Equal(t, input.Parameters["query"], "{\"Field\": \"Size\",\"Value\": \"1\",\"Operation\": \"gt\"}")
	assert.Equal(t, input.Parameters["sort"], "Size")
	assert.Equal(t, input.Parameters["order"], "acs")
	assert.Equal(t, input.Parameters["aggregations"], "Size")
	assert.Equal(t, input.Parameters["withFields"], "[\"Filename\",\"Size\"]")
	assert.Equal(t, input.Parameters["withoutTotalHits"], "true")
}

func TestUnmarshalOutput_SimpleQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
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
</MetaQuery>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &SimpleQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.NextToken, "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****")
	assert.Equal(t, *result.TotalHits, int64(150))
	assert.Equal(t, len(result.Files), 1)
	assert.Equal(t, *result.Files[0].Filename, "docs/report.pdf")
	assert.Equal(t, *result.Files[0].Size, int64(5242880))
	assert.Equal(t, *result.Files[0].URI, "oss://examplebucket/docs/report.pdf")
	assert.Equal(t, *result.Files[0].OSSURI, "oss://examplebucket/docs/report.pdf")
	assert.Equal(t, *result.Files[0].MediaType, "document")
	assert.Equal(t, *result.Files[0].ContentType, "application/pdf")
	assert.Equal(t, *result.Files[0].FileModifiedTime, "2025-12-01T10:30:00Z")
	assert.Equal(t, *result.Files[0].PageCount, int64(20))
	assert.Equal(t, len(result.Aggregations), 1)
	assert.Equal(t, *result.Aggregations[0].Field, "MediaType")
	assert.Equal(t, *result.Aggregations[0].Operation, "group")
	assert.Equal(t, len(result.Aggregations[0].AggregationGroups), 2)
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[0].Value, "document")
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[0].Count, int64(80))
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[1].Value, "image")
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[1].Count, int64(70))

	body = `<?xml version="1.0" encoding="UTF-8"?>
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
</MetaQuery>
`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &SimpleQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.NextToken, "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****")
	assert.Equal(t, *result.TotalHits, int64(258))
	assert.Equal(t, len(result.Files), 2)
	assert.Equal(t, *result.Files[0].Filename, "photos/sunset.jpg")
	assert.Equal(t, *result.Files[0].Size, int64(2048000))
	assert.Equal(t, *result.Files[0].URI, "oss://examplebucket/photos/sunset.jpg")
	assert.Equal(t, *result.Files[0].OSSURI, "oss://examplebucket/photos/sunset.jpg")
	assert.Equal(t, *result.Files[0].MediaType, "image")
	assert.Equal(t, *result.Files[0].ContentType, "image/jpeg")
	assert.Equal(t, *result.Files[0].FileModifiedTime, "2025-12-01T10:30:00Z")
	assert.Equal(t, *result.Files[0].ImageWidth, int64(3840))
	assert.Equal(t, *result.Files[0].ImageHeight, int64(2160))
	assert.Equal(t, *result.Files[0].Orientation, int64(1))

	assert.Equal(t, *result.Files[1].Filename, "photos/mountain.png")
	assert.Equal(t, *result.Files[1].Size, int64(5120000))
	assert.Equal(t, *result.Files[1].URI, "oss://examplebucket/photos/mountain.png")
	assert.Equal(t, *result.Files[1].OSSURI, "oss://examplebucket/photos/mountain.png")
	assert.Equal(t, *result.Files[1].MediaType, "image")
	assert.Equal(t, *result.Files[1].ContentType, "image/png")
	assert.Equal(t, *result.Files[1].FileModifiedTime, "2025-11-20T14:00:00Z")
	assert.Equal(t, *result.Files[1].ImageWidth, int64(1920))
	assert.Equal(t, *result.Files[1].ImageHeight, int64(1080))
	assert.Equal(t, *result.Files[1].Orientation, int64(1))
	assert.Equal(t, len(result.Aggregations), 1)
	assert.Equal(t, *result.Aggregations[0].Field, "MediaType")
	assert.Equal(t, *result.Aggregations[0].Operation, "group")
	assert.Equal(t, len(result.Aggregations[0].AggregationGroups), 2)
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[0].Value, "image")
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[0].Count, int64(200))
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[1].Value, "video")
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[1].Count, int64(58))

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &SimpleQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_SemanticQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *SemanticQueryRequest
	var input *oss.OperationInput
	var err error

	request = &SemanticQueryRequest{}
	input = &oss.OperationInput{
		OpName: "SemanticQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "semanticQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &SemanticQueryRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "SemanticQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "semanticQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &SemanticQueryRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
	}
	input = &oss.OperationInput{
		OpName: "SemanticQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "semanticQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")

	request = &SemanticQueryRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		MaxResults:  oss.Ptr(int32(10)),
		Query:       oss.Ptr("{\"Field\": \"Size\",\"Value\": \"1\",\"Operation\": \"gt\"}"),
		WithFields:  oss.Ptr(`["Filename","Size"]`),
		MediaTypes:  oss.Ptr(`["video","image"]`),
		SourceUri:   oss.Ptr("oss://bucket/prefix"),
	}
	input = &oss.OperationInput{
		OpName: "SemanticQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "semanticQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["maxResults"], "10")
	assert.Equal(t, input.Parameters["query"], "{\"Field\": \"Size\",\"Value\": \"1\",\"Operation\": \"gt\"}")
	assert.Equal(t, input.Parameters["mediaTypes"], "[\"video\",\"image\"]")
	assert.Equal(t, input.Parameters["withFields"], "[\"Filename\",\"Size\"]")
	assert.Equal(t, input.Parameters["sourceURI"], "oss://bucket/prefix")
}

func TestUnmarshalOutput_SemanticQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<?xml version=\"1.0\" encoding=\"UTF-8\"?>
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
</MetaQuery>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &SemanticQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Files), 1)
	assert.Equal(t, *result.Files[0].Bitrate, int64(1656706))
	assert.Equal(t, *result.Files[0].ContentMd5, "5oJccWuBoqVXS8zrzckPlg==")
	assert.Equal(t, *result.Files[0].ContentType, "video/mp4")
	assert.Equal(t, *result.Files[0].CreateTime, "2026-04-21T20:28:17.018858947+08:00")
	assert.Equal(t, *result.Files[0].DatasetName, "test-dataset-sem-vid-1776774492")
	assert.Equal(t, *result.Files[0].Duration, float64(16.034))
	assert.Equal(t, *result.Files[0].ETag, "\\\"E6825C716B81A2A5574BCCEBCDC90F96\\\"")
	assert.Equal(t, *result.Files[0].FileHash, "E6825C716B81A2A5574BCCEBCDC90F96")
	assert.Equal(t, *result.Files[0].FileModifiedTime, "2026-04-21T20:28:13+08:00")
	assert.Equal(t, *result.Files[0].Filename, "test-temp/sem-vid-1776774492774503000.mp4")
	assert.Equal(t, *result.Files[0].FormatLongName, "QuickTime / MOV")
	assert.Equal(t, *result.Files[0].FormatName, "mov,mp4,m4a,3gp,3g2,mj2")
	assert.Equal(t, *result.Files[0].MediaType, "video")
	assert.Equal(t, *result.Files[0].Size, int64(3320455))
	assert.Equal(t, *result.Files[0].VideoWidth, int64(1920))
	assert.Equal(t, *result.Files[0].VideoHeight, int64(1080))
	assert.Equal(t, *result.Files[0].StreamCount, int64(2))
	assert.Equal(t, *result.Files[0].OSSObjectType, "Normal")
	assert.Equal(t, *result.Files[0].OSSStorageClass, "Standard")
	assert.Equal(t, *result.Files[0].OSSTagging[0].Key, "routing-dataset")
	assert.Equal(t, *result.Files[0].OSSTagging[0].Value, "test-dataset-sem-vid-1776774492")
	assert.Equal(t, *result.Files[0].OSSTaggingCount, int64(1))
	assert.Equal(t, *result.Files[0].ObjectACL, "default")
	assert.Equal(t, *result.Files[0].SequenceNumber, int64(2))
	assert.Equal(t, *result.Files[0].SemanticSimilarity, float64(0.5583347777557373))
	assert.Equal(t, *result.Files[0].Size, int64(3320455))
	assert.Equal(t, *result.Files[0].URI, "oss://oss-metaquery-dataset-test/test-temp/sem-vid-1776774492774503000.mp4")
	assert.Equal(t, *result.Files[0].UpdateTime, "2026-04-21T20:28:27.359034257+08:00")

	assert.Equal(t, len(result.Files[0].AudioStreams), 1)
	assert.Equal(t, *result.Files[0].AudioStreams[0].Bitrate, int64(128000))
	assert.Equal(t, *result.Files[0].AudioStreams[0].Channels, int64(2))
	assert.Equal(t, *result.Files[0].AudioStreams[0].ChannelLayout, "stereo")
	assert.Equal(t, *result.Files[0].AudioStreams[0].CodecLongName, "AAC (Advanced Audio Coding)")
	assert.Equal(t, *result.Files[0].AudioStreams[0].CodecName, "aac")
	assert.Equal(t, *result.Files[0].AudioStreams[0].CodecTag, "0x6134706d")
	assert.Equal(t, *result.Files[0].AudioStreams[0].CodecTagString, "mp4a")
	assert.Equal(t, *result.Files[0].AudioStreams[0].Duration, float64(16.021769))
	assert.Equal(t, *result.Files[0].AudioStreams[0].FrameCount, int64(690))
	assert.Equal(t, *result.Files[0].AudioStreams[0].Index, int64(1))
	assert.Equal(t, *result.Files[0].AudioStreams[0].SampleFormat, "fltp")
	assert.Equal(t, *result.Files[0].AudioStreams[0].SampleRate, int64(44100))
	assert.Equal(t, *result.Files[0].AudioStreams[0].TimeBase, "1/44100")

	assert.Equal(t, *result.Files[0].Insights.Video.Description, "这是一段室内高角度监控录像，场景为一个客厅。")
	assert.Equal(t, *result.Files[0].Insights.Video.Caption, "蓝衣男走向餐桌")

	assert.Equal(t, len(result.Files[0].VideoStreams), 1)
	assert.Equal(t, *result.Files[0].VideoStreams[0].AverageFrameRate, "21645000/721493")
	assert.Equal(t, *result.Files[0].VideoStreams[0].BitDepth, int64(8))
	assert.Equal(t, *result.Files[0].VideoStreams[0].Bitrate, int64(1521221))
	assert.Equal(t, *result.Files[0].VideoStreams[0].CodecLongName, "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10")
	assert.Equal(t, *result.Files[0].VideoStreams[0].CodecName, "h264")
	assert.Equal(t, *result.Files[0].VideoStreams[0].CodecTag, "0x31637661")
	assert.Equal(t, *result.Files[0].VideoStreams[0].CodecTagString, "avc1")
	assert.Equal(t, *result.Files[0].VideoStreams[0].ColorPrimaries, "bt709")
	assert.Equal(t, *result.Files[0].VideoStreams[0].ColorRange, "tv")
	assert.Equal(t, *result.Files[0].VideoStreams[0].ColorTransfer, "bt709")
	assert.Equal(t, *result.Files[0].VideoStreams[0].ColorSpace, "bt709")
	assert.Equal(t, *result.Files[0].VideoStreams[0].Duration, float64(16.033178))
	assert.Equal(t, *result.Files[0].VideoStreams[0].FrameCount, int64(481))
	assert.Equal(t, *result.Files[0].VideoStreams[0].FrameRate, "90000/2999")
	assert.Equal(t, *result.Files[0].VideoStreams[0].Height, int64(1080))
	assert.Equal(t, *result.Files[0].VideoStreams[0].Width, int64(1920))
	assert.Equal(t, *result.Files[0].VideoStreams[0].Level, int64(31))
	assert.Equal(t, *result.Files[0].VideoStreams[0].PixelFormat, "yuv420p")
	assert.Equal(t, *result.Files[0].VideoStreams[0].Profile, "High")
	assert.Equal(t, *result.Files[0].VideoStreams[0].PixelFormat, "yuv420p")
	assert.Equal(t, *result.Files[0].VideoStreams[0].SampleAspectRatio, "1:1")
	assert.Equal(t, *result.Files[0].VideoStreams[0].TimeBase, "1/90000")

	body = `<?xml version="1.0" encoding="UTF-8"?>
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
</MetaQuery>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &SemanticQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Files), 1)
	assert.Equal(t, *result.Files[0].Bitrate, int64(196284))
	assert.Equal(t, *result.Files[0].ContentMd5, "5/ZLrWYXpuQfDfxEf4+lyA==")
	assert.Equal(t, *result.Files[0].ContentType, "video/mp4")
	assert.Equal(t, *result.Files[0].CreateTime, "2026-04-21T10:51:38.264045621+08:00")
	assert.Equal(t, *result.Files[0].DatasetName, "dataset-aianalysis-walk")
	assert.Equal(t, *result.Files[0].Duration, float64(8))
	assert.Equal(t, *result.Files[0].ETag, "\\\"E7F64BAD6617A6E41F0DFC447F8FA5C8\\\"")
	assert.Equal(t, *result.Files[0].FileHash, "E7F64BAD6617A6E41F0DFC447F8FA5C8")
	assert.Equal(t, *result.Files[0].FileModifiedTime, "2026-04-21T10:51:25+08:00")
	assert.Equal(t, *result.Files[0].Filename, "mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4")
	assert.Equal(t, *result.Files[0].FormatLongName, "QuickTime / MOV")
	assert.Equal(t, *result.Files[0].FormatName, "mov,mp4,m4a,3gp,3g2,mj2")
	assert.Equal(t, *result.Files[0].MediaType, "video")
	assert.Equal(t, *result.Files[0].OSSCRC64, "16628192875747293357")
	assert.Equal(t, *result.Files[0].Size, int64(196284))
	assert.Equal(t, *result.Files[0].VideoWidth, int64(640))
	assert.Equal(t, *result.Files[0].VideoHeight, int64(360))
	assert.Equal(t, *result.Files[0].StreamCount, int64(2))
	assert.Equal(t, *result.Files[0].OSSObjectType, "Normal")
	assert.Equal(t, *result.Files[0].OSSStorageClass, "Standard")
	assert.Equal(t, len(result.Files[0].OSSTagging), 2)
	assert.Equal(t, *result.Files[0].OSSTagging[0].Key, "routing-dataset")
	assert.Equal(t, *result.Files[0].OSSTagging[1].Key, "env")
	assert.Equal(t, *result.Files[0].OSSTagging[0].Value, "photos-2026")
	assert.Equal(t, *result.Files[0].OSSTagging[1].Value, "production")
	assert.Equal(t, *result.Files[0].ObjectACL, "default")
	assert.Equal(t, len(result.Files[0].OSSUserMeta), 2)
	assert.Equal(t, *result.Files[0].OSSUserMeta[0].Key, "category")
	assert.Equal(t, *result.Files[0].OSSUserMeta[0].Value, "photo")
	assert.Equal(t, *result.Files[0].OSSUserMeta[1].Key, "album")
	assert.Equal(t, *result.Files[0].OSSUserMeta[1].Value, "vacation")
	assert.Equal(t, len(result.Files[0].CustomLabels), 1)
	assert.Equal(t, *result.Files[0].CustomLabels[0].Key, "label-ley")
	assert.Equal(t, *result.Files[0].CustomLabels[0].Value, "label-val")

	assert.Equal(t, *result.Files[0].ProduceTime, "2026-04-21T10:46:10+08:00")
	assert.Equal(t, *result.Files[0].SequenceNumber, int64(5))
	assert.Equal(t, *result.Files[0].SemanticSimilarity, float64(0.2536))
	assert.Equal(t, *result.Files[0].Size, int64(196284))
	assert.Equal(t, *result.Files[0].StreamCount, int64(2))
	assert.Equal(t, *result.Files[0].URI, "oss://paas-smart-cloud-test/mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4")
	assert.Equal(t, *result.Files[0].UpdateTime, "2026-04-21T10:52:39.412605575+08:00")

	assert.Equal(t, len(result.Files[0].Labels), 1)
	assert.Equal(t, *result.Files[0].Labels[0].LabelConfidence, float64(1))
	assert.Equal(t, *result.Files[0].Labels[0].LabelName, "有人走过")
	assert.Equal(t, *result.Files[0].Labels[0].ParentLabelName, "自定义标签")
	assert.Equal(t, len(result.Files[0].Labels[0].Clips), 1)
	assert.Equal(t, result.Files[0].Labels[0].Clips[0].TimeRange[0], int64(200))
	assert.Equal(t, result.Files[0].Labels[0].Clips[0].TimeRange[1], int64(5533))

	assert.Equal(t, len(result.Files[0].AudioStreams), 1)
	assert.Equal(t, *result.Files[0].AudioStreams[0].Bitrate, int64(14983))
	assert.Equal(t, *result.Files[0].AudioStreams[0].Channels, int64(1))
	assert.Equal(t, *result.Files[0].AudioStreams[0].ChannelLayout, "mono")
	assert.Equal(t, *result.Files[0].AudioStreams[0].CodecLongName, "AAC (Advanced Audio Coding)")
	assert.Equal(t, *result.Files[0].AudioStreams[0].CodecName, "aac")
	assert.Equal(t, *result.Files[0].AudioStreams[0].CodecTag, "0x6134706d")
	assert.Equal(t, *result.Files[0].AudioStreams[0].CodecTagString, "mp4a")
	assert.Equal(t, *result.Files[0].AudioStreams[0].Duration, float64(7.936))
	assert.Equal(t, *result.Files[0].AudioStreams[0].FrameCount, int64(62))
	assert.Equal(t, *result.Files[0].AudioStreams[0].Index, int64(1))
	assert.Equal(t, *result.Files[0].AudioStreams[0].SampleFormat, "fltp")
	assert.Equal(t, *result.Files[0].AudioStreams[0].SampleRate, int64(8000))
	assert.Equal(t, *result.Files[0].AudioStreams[0].TimeBase, "1/8000")

	assert.Equal(t, len(result.Files[0].VideoStreams), 1)
	assert.Equal(t, *result.Files[0].VideoStreams[0].AverageFrameRate, "15/1")
	assert.Equal(t, *result.Files[0].VideoStreams[0].BitDepth, int64(8))
	assert.Equal(t, *result.Files[0].VideoStreams[0].Bitrate, int64(178202))
	assert.Equal(t, *result.Files[0].VideoStreams[0].CodecLongName, "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10")
	assert.Equal(t, *result.Files[0].VideoStreams[0].CodecName, "h264")
	assert.Equal(t, *result.Files[0].VideoStreams[0].CodecTag, "0x31637661")
	assert.Equal(t, *result.Files[0].VideoStreams[0].CodecTagString, "avc1")
	assert.Equal(t, *result.Files[0].VideoStreams[0].Duration, float64(8))
	assert.Equal(t, *result.Files[0].VideoStreams[0].FrameCount, int64(120))
	assert.Equal(t, *result.Files[0].VideoStreams[0].FrameRate, "500/33")
	assert.Equal(t, *result.Files[0].VideoStreams[0].Height, int64(360))
	assert.Equal(t, *result.Files[0].VideoStreams[0].Width, int64(640))
	assert.Equal(t, *result.Files[0].VideoStreams[0].Level, int64(22))
	assert.Equal(t, *result.Files[0].VideoStreams[0].TimeBase, "1/1000")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &SemanticQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_OpenMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *OpenMetaQueryRequest
	var input *oss.OperationInput
	var err error

	request = &OpenMetaQueryRequest{}
	input = &oss.OperationInput{
		OpName: "OpenMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "openMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &OpenMetaQueryRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "OpenMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "openMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Mode.")

	request = &OpenMetaQueryRequest{
		Bucket: oss.Ptr("bucket"),
		Mode:   oss.Ptr("basic"),
		Role:   oss.Ptr("my-role"),
		MetaQuery: &OpenMetaQuery{
			Filters: &Filters{
				Filter: []string{
					"Size > 0",
					"ContentType = 'image/jpeg'",
				},
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "OpenMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "openMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["mode"], "basic")
	assert.Equal(t, input.Parameters["role"], "my-role")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<MetaQuery><Filters><Filter>Size &gt; 0</Filter><Filter>ContentType = &#39;image/jpeg&#39;</Filter></Filters></MetaQuery>")

	request = &OpenMetaQueryRequest{
		Bucket: oss.Ptr("bucket"),
		Role:   oss.Ptr("AliyunMetaQueryDefaultRole"),
		Mode:   oss.Ptr("semantic"),
		MetaQuery: &OpenMetaQuery{
			WorkflowParameters: &WorkflowParameters{
				WorkflowParameter: []WorkflowParameter{
					{
						Name:  oss.Ptr("ImageInsightEnable"),
						Value: oss.Ptr("True"),
					},
					{
						Name:  oss.Ptr("VideoInsightEnable"),
						Value: oss.Ptr("True"),
					},
				},
			},
			NotificationAttributes: &NotificationAttributes{
				Notifications: &Notifications{
					Notification: []Notification{
						{
							MNS: oss.Ptr("imm-index-notification"),
						},
					},
				},
				WithFields: &WithFields{
					[]string{
						"Insights",
						"Labels",
					},
				},
			},
			IndexOptions: &IndexOptions{
				IgnoreObjectDelete: oss.Ptr(true),
			},
			RouteRule: &RouteRule{
				Type:              oss.Ptr("OSSTag"),
				AutoCreateDataset: oss.Ptr(true),
				OSSTagKey:         oss.Ptr("routing-dataset"),
			},
			DatasetConfig: &DatasetConfig{
				Insights: &InsightsConfig{
					Language: oss.Ptr("en"),
				},
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "OpenMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "openMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["mode"], "semantic")
	assert.Equal(t, input.Parameters["role"], "AliyunMetaQueryDefaultRole")
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<MetaQuery><WorkflowParameters><WorkflowParameter><Name>ImageInsightEnable</Name><Value>True</Value></WorkflowParameter><WorkflowParameter><Name>VideoInsightEnable</Name><Value>True</Value></WorkflowParameter></WorkflowParameters><NotificationAttributes><Notifications><Notification><MNS>imm-index-notification</MNS></Notification></Notifications><WithFields><WithField>Insights</WithField><WithField>Labels</WithField></WithFields></NotificationAttributes><DatasetConfig><Insights><Language>en</Language></Insights></DatasetConfig><IndexOptions><IgnoreObjectDelete>true</IgnoreObjectDelete></IndexOptions><RouteRule><Type>OSSTag</Type><AutoCreateDataset>true</AutoCreateDataset><OSSTagKey>routing-dataset</OSSTagKey></RouteRule></MetaQuery>")
}

func TestUnmarshalOutput_OpenMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &OpenMetaQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &OpenMetaQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &OpenMetaQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetMetaQueryStatus(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetMetaQueryStatusRequest
	var input *oss.OperationInput
	var err error

	request = &GetMetaQueryStatusRequest{}
	input = &oss.OperationInput{
		OpName: "GetMetaQueryStatus",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "getMetaQueryStatus",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetMetaQueryStatusRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "GetMetaQueryStatus",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "getMetaQueryStatus",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, "bucket", *input.Bucket)
	assert.Equal(t, "getMetaQueryStatus", input.Parameters["action"])
}

func TestUnmarshalOutput_GetMetaQueryStatus(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<MetaQueryStatus>
  <State>Running</State>
  <Phase>IncrementalScanning</Phase>
  <CreateTime>2026-05-20T08:00:00.000+08:00</CreateTime>
  <UpdateTime>2026-05-20T08:30:00.000+08:00</UpdateTime>
  <MetaQueryMode>semantic</MetaQueryMode>
  <WorkflowParameters>
    <WorkflowParameter>
      <Name>ImageInsightEnable</Name>
      <Value>True</Value>
    </WorkflowParameter>
  </WorkflowParameters>
  <Filters>
    <Filter>Size > 1024,FileModifiedTime > 2025-06-03T09:20:47.999Z</Filter>
    <Filter>Filename prefix (YWEvYmIv)</Filter>
  </Filters>
  <IndexOptions>
    <IgnoreObjectDelete>True</IgnoreObjectDelete>
  </IndexOptions>
  <RouteRule>
    <Type>OSSTag</Type>
    <AutoCreateDataset>True</AutoCreateDataset>
    <OSSTagKey>routing-dataset</OSSTagKey>
  </RouteRule>
  <NotificationAttributes>
    <Notifications>
      <Notification>
        <MNS>imm-index-notification</MNS>
      </Notification>
    </Notifications>
    <WithFields>
      <WithField>Insights</WithField>
    </WithFields>
  </NotificationAttributes>
  <DatasetConfig>
    <Insights>
      <Language>en</Language>
    </Insights>
  </DatasetConfig>
</MetaQueryStatus>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetMetaQueryStatusResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.State, "Running")
	assert.Equal(t, *result.Phase, "IncrementalScanning")
	assert.Equal(t, *result.CreateTime, "2026-05-20T08:00:00.000+08:00")
	assert.Equal(t, *result.UpdateTime, "2026-05-20T08:30:00.000+08:00")
	assert.Equal(t, *result.MetaQueryMode, "semantic")
	assert.Equal(t, *result.WorkflowParameters.WorkflowParameter[0].Name, "ImageInsightEnable")
	assert.Equal(t, *result.WorkflowParameters.WorkflowParameter[0].Value, "True")
	assert.Equal(t, result.Filters.Filter[0], "Size > 1024,FileModifiedTime > 2025-06-03T09:20:47.999Z")
	assert.Equal(t, result.Filters.Filter[1], "Filename prefix (YWEvYmIv)")
	assert.Equal(t, *result.IndexOptions.IgnoreObjectDelete, true)
	assert.Equal(t, *result.RouteRule.Type, "OSSTag")
	assert.Equal(t, *result.RouteRule.AutoCreateDataset, true)
	assert.Equal(t, *result.RouteRule.OSSTagKey, "routing-dataset")
	assert.Equal(t, *result.NotificationAttributes.Notifications.Notification[0].MNS, "imm-index-notification")
	assert.Equal(t, result.NotificationAttributes.WithFields.WithField[0], "Insights")
	assert.Equal(t, *result.DatasetConfig.Insights.Language, "en")

	body = `<MetaQueryStatus>
    <State>Running</State>
    <Phase>IncrementalScanning</Phase>
    <CreateTime>2026-06-05T10:49:57.109295646+08:00</CreateTime>
    <UpdateTime>2026-06-05T10:49:59.411578836+08:00</UpdateTime>
    <MetaQueryMode>basic</MetaQueryMode>
    <RouteRule>
        <Type>default</Type>
        <AutoCreateDataset>True</AutoCreateDataset>
    </RouteRule>
    <WorkflowParameters>
    </WorkflowParameters>
    <Filters>
    </Filters>
    <DatasetConfig>
        <ReverseImage>
            <Video>
                <Enable>false</Enable>
            </Video>
            <Image>
                <Enable>false</Enable>
            </Image>
        </ReverseImage>
        <Insights>
            <Language>zh-Hans</Language>
            <Image>
                <Caption>
                    <Enable>false</Enable>
                    <Prompt>
                    </Prompt>
                </Caption>
            </Image>
            <Video>
                <Caption>
                    <Enable>false</Enable>
                    <Prompt>
                    </Prompt>
                    <PersonReference>
                        <Enable>false</Enable>
                    </PersonReference>
                </Caption>
                <Label>
                    <System>
                        <Enable>false</Enable>
                    </System>
                    <UserDefined>
                        <Enable>false</Enable>
                        <Labels>
                        </Labels>
                    </UserDefined>
                    <Highlight>
                        <Enable>false</Enable>
                        <Labels>
                        </Labels>
                    </Highlight>
                </Label>
                <MultiStream>
                    <Enable>false</Enable>
                </MultiStream>
            </Video>
        </Insights>
        <SmartCluster>
            <Figure>
                <AutoGenerate>false</AutoGenerate>
                <AutoClustering>false</AutoClustering>
                <MinEntityCount>3</MinEntityCount>
                <EnabledFeatures>face</EnabledFeatures>
            </Figure>
        </SmartCluster>
    </DatasetConfig>
</MetaQueryStatus>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetMetaQueryStatusResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.State, "Running")
	assert.Equal(t, *result.Phase, "IncrementalScanning")
	assert.Equal(t, *result.CreateTime, "2026-06-05T10:49:57.109295646+08:00")
	assert.Equal(t, *result.UpdateTime, "2026-06-05T10:49:59.411578836+08:00")
	assert.Equal(t, *result.MetaQueryMode, "basic")
	assert.Equal(t, *result.RouteRule.Type, "default")
	assert.Equal(t, *result.RouteRule.AutoCreateDataset, true)
	assert.Equal(t, *result.DatasetConfig.ReverseImage.Video.Enable, false)
	assert.Equal(t, *result.DatasetConfig.ReverseImage.Image.Enable, false)
	assert.Equal(t, *result.DatasetConfig.Insights.Language, "zh-Hans")
	assert.Equal(t, *result.DatasetConfig.Insights.Image.Caption.Enable, false)
	assert.Equal(t, *result.DatasetConfig.Insights.Video.Caption.Enable, false)
	assert.Equal(t, *result.DatasetConfig.Insights.Video.Caption.PersonReference.Enable, false)
	assert.Equal(t, *result.DatasetConfig.Insights.Video.Label.System.Enable, false)
	assert.Equal(t, *result.DatasetConfig.Insights.Video.Label.UserDefined.Enable, false)
	assert.Equal(t, *result.DatasetConfig.Insights.Video.Label.Highlight.Enable, false)
	assert.Equal(t, *result.DatasetConfig.Insights.Video.MultiStream.Enable, false)
	assert.Equal(t, *result.DatasetConfig.SmartCluster.Figure.AutoGenerate, false)
	assert.Equal(t, *result.DatasetConfig.SmartCluster.Figure.AutoClustering, false)
	assert.Equal(t, *result.DatasetConfig.SmartCluster.Figure.MinEntityCount, int64(3))
	assert.Equal(t, result.DatasetConfig.SmartCluster.Figure.EnabledFeatures[0], "face")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetMetaQueryStatusResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DoMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DoMetaQueryRequest
	var input *oss.OperationInput
	var err error

	request = &DoMetaQueryRequest{}
	input = &oss.OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "doMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DoMetaQueryRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "doMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Mode.")

	request = &DoMetaQueryRequest{
		Bucket: oss.Ptr("bucket"),
		Mode:   oss.Ptr("basic"),
	}
	input = &oss.OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "doMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, MetaQuery.")

	request = &DoMetaQueryRequest{
		Bucket: oss.Ptr("bucket"),
		Mode:   oss.Ptr("basic"),
		MetaQuery: &DoMetaQuery{
			Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
			Sort:  oss.Ptr("Size"),
			Order: oss.Ptr(MetaQueryOrderDesc),
			Aggregations: &MetaQueryAggregations{
				[]Aggregation{
					{
						Field:     oss.Ptr("Size"),
						Operation: oss.Ptr("sum"),
					},
				},
			},
			MaxResults: oss.Ptr(int64(100)),
		},
	}
	input = &oss.OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "doMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, "bucket", *input.Bucket)
	assert.Equal(t, "doMetaQuery", input.Parameters["action"])
	assert.Equal(t, "basic", input.Parameters["mode"])
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")

	request = &DoMetaQueryRequest{
		Bucket: oss.Ptr("bucket"),
		Mode:   oss.Ptr("semantic"),
		MetaQuery: &DoMetaQuery{
			Query: oss.Ptr(`The cat in the living room`),
			MediaTypes: &MetaQueryMediaTypes{
				[]string{"image"},
			},
			SimpleQuery: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"102400"}`),
			MaxResults:  oss.Ptr(int64(20)),
		},
	}
	input = &oss.OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "doMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, "bucket", *input.Bucket)
	assert.Equal(t, "doMetaQuery", input.Parameters["action"])
	assert.Equal(t, "semantic", input.Parameters["mode"])
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<MetaQuery><MaxResults>20</MaxResults><Query>The cat in the living room</Query><MediaTypes><MediaType>image</MediaType></MediaTypes><SimpleQuery>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;102400&#34;}</SimpleQuery></MetaQuery>")

	request = &DoMetaQueryRequest{
		Bucket: oss.Ptr("bucket"),
		Mode:   oss.Ptr("semantic"),
		MetaQuery: &DoMetaQuery{
			Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
			MediaTypes: &MetaQueryMediaTypes{
				[]string{"image"},
			},
			SimpleQuery: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"102400"}`),
			MaxResults:  oss.Ptr(int64(10)),
			NextToken:   oss.Ptr("MTIzNDU2"),
			WithFields: &WithFields{
				[]string{
					"URI", "Size", "FileModifiedTime",
				},
			},
			WithoutTotalHits: oss.Ptr("False"),
			Aggregations: &MetaQueryAggregations{
				Aggregations: []Aggregation{
					{
						Operation: oss.Ptr("sum"),
						Field:     oss.Ptr("Field"),
					},
				},
			},
			SmartClusterIds: &SmartClusterIds{
				SmartClusterId: []string{
					"cluster-001",
				},
			},
		},
	}
	input = &oss.OperationInput{
		OpName: "DoMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "doMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, "bucket", *input.Bucket)
	assert.Equal(t, "doMetaQuery", input.Parameters["action"])
	assert.Equal(t, "semantic", input.Parameters["mode"])
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<MetaQuery><MaxResults>10</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Aggregations><Aggregation><Operation>sum</Operation><Field>Field</Field><Groups></Groups></Aggregation></Aggregations><NextToken>MTIzNDU2</NextToken><MediaTypes><MediaType>image</MediaType></MediaTypes><SimpleQuery>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;102400&#34;}</SimpleQuery><WithoutTotalHits>False</WithoutTotalHits><SmartClusterIds><SmartClusterId>cluster-001</SmartClusterId></SmartClusterIds><WithFields><WithField>URI</WithField><WithField>Size</WithField><WithField>FileModifiedTime</WithField></WithFields></MetaQuery>")
}

func TestUnmarshalOutput_DoMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
    <NextToken>next-page-token-abc</NextToken>
    <TotalHits>123</TotalHits>
    <Files>
        <File>
            <Filename>photos/sunset.jpg</Filename>
            <Size>2097152</Size>
            <FileModifiedTime>2026-05-19T15:30:00.000+08:00</FileModifiedTime>
            <ContentType>image/jpeg</ContentType>
            <ObjectACL>default</ObjectACL>
            <OSSStorageClass>Standard</OSSStorageClass>
        </File>
        <File>
            <Filename>photos/mountain.png</Filename>
            <Size>5242880</Size>
        </File>
    </Files>
    <Aggregations>
        <Aggregation>
            <Field>Size</Field>
            <Operation>sum</Operation>
            <Value>12345678</Value>
        </Aggregation>
    </Aggregations>
</MetaQuery>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &DoMetaQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.NextToken, "next-page-token-abc")
	assert.Equal(t, *result.TotalHits, int64(123))
	assert.Equal(t, *result.Files[0].Filename, "photos/sunset.jpg")
	assert.Equal(t, *result.Files[0].Size, int64(2097152))
	assert.Equal(t, *result.Files[0].FileModifiedTime, "2026-05-19T15:30:00.000+08:00")
	assert.Equal(t, *result.Files[0].ContentType, "image/jpeg")
	assert.Equal(t, *result.Files[0].ObjectACL, "default")
	assert.Equal(t, *result.Files[0].OSSStorageClass, "Standard")
	assert.Equal(t, *result.Files[1].Filename, "photos/mountain.png")
	assert.Equal(t, *result.Files[1].Size, int64(5242880))
	assert.Equal(t, *result.Aggregations[0].Field, "Size")
	assert.Equal(t, *result.Aggregations[0].Operation, "sum")
	assert.Equal(t, *result.Aggregations[0].Value, "12345678")

	body = `<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
    <TotalHits>2</TotalHits>
    <Files>
        <File>
            <Filename>photos/cat-in-living-room.jpg</Filename>
            <Size>3145728</Size>
            <OSSStorageClass>Standard</OSSStorageClass>
            <Labels>
                <Label>
                    <LabelName>cat</LabelName>
                    <LabelConfidence>0.98</LabelConfidence>
                </Label>
            </Labels>
        </File>
        <File>
            <Filename>photos/kitten-sofa.jpg</Filename>
            <Size>2621440</Size>
        </File>
    </Files>
</MetaQuery>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DoMetaQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.TotalHits, int64(2))
	assert.Equal(t, *result.Files[0].Filename, "photos/cat-in-living-room.jpg")
	assert.Equal(t, *result.Files[0].Size, int64(3145728))
	assert.Equal(t, *result.Files[0].OSSStorageClass, "Standard")
	assert.Equal(t, *result.Files[0].Labels[0].LabelName, "cat")
	assert.Equal(t, *result.Files[0].Labels[0].LabelConfidence, float64(0.98))
	assert.Equal(t, *result.Files[1].Filename, "photos/kitten-sofa.jpg")
	assert.Equal(t, *result.Files[1].Size, int64(2621440))

	body = `<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
    <TotalHits>50</TotalHits>
    <Aggregations>
        <Aggregation>
            <Field>StorageClass</Field>
            <Operation>group_by</Operation>
            <Groups>
                <AggregationGroup>
                    <Value>Standard</Value>
                    <Count>30</Count>
                </AggregationGroup>
                <AggregationGroup>
                    <Value>IA</Value>
                    <Count>20</Count>
                </AggregationGroup>
            </Groups>
        </Aggregation>
    </Aggregations>
</MetaQuery>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DoMetaQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.TotalHits, int64(50))
	assert.Equal(t, len(result.Aggregations), 1)
	assert.Equal(t, *result.Aggregations[0].Field, "StorageClass")
	assert.Equal(t, *result.Aggregations[0].Operation, "group_by")
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[0].Value, "Standard")
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[0].Count, int64(30))
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[1].Value, "IA")
	assert.Equal(t, *result.Aggregations[0].AggregationGroups[1].Count, int64(20))

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DoMetaQueryResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_CloseMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CloseMetaQueryRequest
	var input *oss.OperationInput
	var err error

	request = &CloseMetaQueryRequest{}
	input = &oss.OperationInput{
		OpName: "CloseMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "closeMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &CloseMetaQueryRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
		OpName: "CloseMetaQuery",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "closeMetaQuery",
		},
		Bucket: request.Bucket,
	}
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, "bucket", *input.Bucket)
	assert.Equal(t, "closeMetaQuery", input.Parameters["action"])
}

func TestUnmarshalOutput_CloseMetaQuery(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &CloseMetaQueryResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CloseMetaQueryResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
