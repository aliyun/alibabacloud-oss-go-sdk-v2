package dataprocess

import (
	"context"
	"encoding/xml"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// Image represents image information
type Image struct {
	XMLName             xml.Name             `xml:"Image"`
	ImageWidth          *int64               `xml:"ImageWidth,omitempty"`
	ImageHeight         *int64               `xml:"ImageHeight,omitempty"`
	Exif                *string              `xml:"EXIF,omitempty"`
	ImageScore          *ImageScore          `xml:"ImageScore,omitempty"`
	CroppingSuggestions []CroppingSuggestion `xml:"CroppingSuggestions>CroppingSuggestion,omitempty"`
	OCRContents         []OCRContents        `xml:"OCRContents>OCRContents,omitempty"`
}

// VideoStream represents video stream information
type VideoStream struct {
	XMLName            xml.Name `xml:"VideoStream"`
	Index              *int64   `xml:"Index,omitempty"`
	Language           *string  `xml:"Language,omitempty"`
	CodecName          *string  `xml:"CodecName,omitempty"`
	CodecLongName      *string  `xml:"CodecLongName,omitempty"`
	Profile            *string  `xml:"Profile,omitempty"`
	CodecTimeBase      *string  `xml:"CodecTimeBase,omitempty"`
	CodecTagString     *string  `xml:"CodecTagString,omitempty"`
	CodecTag           *string  `xml:"CodecTag,omitempty"`
	Width              *int64   `xml:"Width,omitempty"`
	Height             *int64   `xml:"Height,omitempty"`
	HasBFrames         *int64   `xml:"HasBFrames,omitempty"`
	SampleAspectRatio  *string  `xml:"SampleAspectRatio,omitempty"`
	DisplayAspectRatio *string  `xml:"DisplayAspectRatio,omitempty"`
	PixelFormat        *string  `xml:"PixelFormat,omitempty"`
	Level              *int64   `xml:"Level,omitempty"`
	FrameRate          *string  `xml:"FrameRate,omitempty"`
	AverageFrameRate   *string  `xml:"AverageFrameRate,omitempty"`
	TimeBase           *string  `xml:"TimeBase,omitempty"`
	StartTime          *float64 `xml:"StartTime,omitempty"`
	Duration           *float64 `xml:"Duration,omitempty"`
	Bitrate            *int64   `xml:"Bitrate,omitempty"`
	FrameCount         *int64   `xml:"FrameCount,omitempty"`
	Rotate             *string  `xml:"Rotate,omitempty"`
	BitDepth           *int64   `xml:"BitDepth,omitempty"`
	ColorSpace         *string  `xml:"ColorSpace,omitempty"`
	ColorRange         *string  `xml:"ColorRange,omitempty"`
	ColorTransfer      *string  `xml:"ColorTransfer,omitempty"`
	ColorPrimaries     *string  `xml:"ColorPrimaries,omitempty"`
}

// AudioStream represents audio stream information
type AudioStream struct {
	XMLName        xml.Name `xml:"AudioStream"`
	Index          *int64   `xml:"Index,omitempty"`
	Language       *string  `xml:"Language,omitempty"`
	CodecName      *string  `xml:"CodecName,omitempty"`
	CodecLongName  *string  `xml:"CodecLongName,omitempty"`
	CodecTimeBase  *string  `xml:"CodecTimeBase,omitempty"`
	CodecTagString *string  `xml:"CodecTagString,omitempty"`
	CodecTag       *string  `xml:"CodecTag,omitempty"`
	TimeBase       *string  `xml:"TimeBase,omitempty"`
	StartTime      *float64 `xml:"StartTime,omitempty"`
	Duration       *float64 `xml:"Duration,omitempty"`
	Bitrate        *int64   `xml:"Bitrate,omitempty"`
	FrameCount     *int64   `xml:"FrameCount,omitempty"`
	Lyric          *string  `xml:"Lyric,omitempty"`
	SampleFormat   *string  `xml:"SampleFormat,omitempty"`
	SampleRate     *int64   `xml:"SampleRate,omitempty"`
	Channels       *int64   `xml:"Channels,omitempty"`
	ChannelLayout  *string  `xml:"ChannelLayout,omitempty"`
}

// SubtitleStream represents subtitle stream information
type SubtitleStream struct {
	XMLName        xml.Name `xml:"SubtitleStream"`
	Index          *int64   `xml:"Index,omitempty"`
	Language       *string  `xml:"Language,omitempty"`
	CodecName      *string  `xml:"CodecName,omitempty"`
	CodecLongName  *string  `xml:"CodecLongName,omitempty"`
	CodecTagString *string  `xml:"CodecTagString,omitempty"`
	CodecTag       *string  `xml:"CodecTag,omitempty"`
	StartTime      *float64 `xml:"StartTime,omitempty"`
	Duration       *float64 `xml:"Duration,omitempty"`
	Bitrate        *int64   `xml:"Bitrate,omitempty"`
	Content        *string  `xml:"Content,omitempty"`
	Width          *int64   `xml:"Width,omitempty"`
	Height         *int64   `xml:"Height,omitempty"`
}

// Label represents a label with confidence score
type Label struct {
	XMLName         xml.Name `xml:"Label"`
	Language        *string  `xml:"Language,omitempty"`
	LabelName       *string  `xml:"LabelName,omitempty"`
	LabelLevel      *int64   `xml:"LabelLevel,omitempty"`
	LabelConfidence *float64 `xml:"LabelConfidence,omitempty"`
	ParentLabelName *string  `xml:"ParentLabelName,omitempty"`
	CentricScore    *float64 `xml:"CentricScore,omitempty"`
	LabelAlias      *string  `xml:"LabelAlias,omitempty"`
	Clips           []Clip   `xml:"Clips>Clip,omitempty"`
}

// Clip represents a time range clip
type Clip struct {
	XMLName   xml.Name `xml:"Clip"`
	TimeRange []int64  `xml:"TimeRange"`
}

// OCRContents represents OCR recognition results
type OCRContents struct {
	XMLName    xml.Name  `xml:"OCRContents"`
	Language   *string   `xml:"Language,omitempty"`
	Contents   *string   `xml:"Contents,omitempty"`
	Confidence *float64  `xml:"Confidence,omitempty"`
	Boundary   *Boundary `xml:"Boundary,omitempty"`
}

// ImageScore represents image quality score
type ImageScore struct {
	XMLName             xml.Name `xml:"ImageScore"`
	OverallQualityScore *float64 `xml:"OverallQualityScore,omitempty"`
}

// Boundary represents a rectangular boundary
type Boundary struct {
	XMLName xml.Name     `xml:"Boundary"`
	Width   *int64       `xml:"Width,omitempty"`
	Height  *int64       `xml:"Height,omitempty"`
	Left    *int64       `xml:"Left,omitempty"`
	Top     *int64       `xml:"Top,omitempty"`
	Polygon []PointInt64 `xml:"Polygon>PointInt64,omitempty"`
}

// PointInt64 represents a 2D point with int64 coordinates
type PointInt64 struct {
	XMLName xml.Name `xml:"PointInt64"`
	X       *int64   `xml:"X,omitempty"`
	Y       *int64   `xml:"Y,omitempty"`
}

// Figure represents a figure/shape in an image
type Figure struct {
	XMLName                 xml.Name  `xml:"Figure"`
	FigureId                *string   `xml:"FigureId,omitempty"`
	FigureConfidence        *float64  `xml:"FigureConfidence,omitempty"`
	FigureClusterId         *string   `xml:"FigureClusterId,omitempty"`
	FigureClusterConfidence *float64  `xml:"FigureClusterConfidence,omitempty"`
	FigureType              *string   `xml:"FigureType,omitempty"`
	Age                     *int64    `xml:"Age,omitempty"`
	AgeSD                   *float64  `xml:"AgeSD,omitempty"`
	Gender                  *string   `xml:"Gender,omitempty"`
	GenderConfidence        *float64  `xml:"GenderConfidence,omitempty"`
	Emotion                 *string   `xml:"Emotion,omitempty"`
	EmotionConfidence       *float64  `xml:"EmotionConfidence,omitempty"`
	FaceQuality             *float64  `xml:"FaceQuality,omitempty"`
	Boundary                *Boundary `xml:"Boundary,omitempty"`
	Mouth                   *string   `xml:"Mouth,omitempty"`
	MouthConfidence         *float64  `xml:"MouthConfidence,omitempty"`
	Beard                   *string   `xml:"Beard,omitempty"`
	BeardConfidence         *float64  `xml:"BeardConfidence,omitempty"`
	Hat                     *string   `xml:"Hat,omitempty"`
	HatConfidence           *float64  `xml:"HatConfidence,omitempty"`
	Mask                    *string   `xml:"Mask,omitempty"`
	MaskConfidence          *float64  `xml:"MaskConfidence,omitempty"`
	Glasses                 *string   `xml:"Glasses,omitempty"`
	GlassesConfidence       *float64  `xml:"GlassesConfidence,omitempty"`
	Sharpness               *float64  `xml:"Sharpness,omitempty"`
	Attractive              *float64  `xml:"Attractive,omitempty"`
	HeadPose                *HeadPose `xml:"HeadPose,omitempty"`
}

// ElementContent represents the content of an element
type ElementContent struct {
	XMLName xml.Name `xml:"ElementContent"`
	Type    *string  `xml:"Type,omitempty"`
	URI     *string  `xml:"URI,omitempty"`
	Value   *string  `xml:"Value,omitempty"`
}

// Element represents a detected element in media
type Element struct {
	XMLName            xml.Name          `xml:"Element"`
	ElementContents    []ElementContent  `xml:"ElementContents>ElementContent,omitempty"`
	ObjectId           *string           `xml:"ObjectId,omitempty"`
	ElementType        *string           `xml:"ElementType,omitempty"`
	SemanticSimilarity *float64          `xml:"SemanticSimilarity,omitempty"`
	ElementRelations   []ElementRelation `xml:"ElementRelations>ElementRelation,omitempty"`
}

// SceneElement represents an element in a scene
type SceneElement struct {
	XMLName          xml.Name `xml:"SceneElement"`
	TimeRange        []int64  `xml:"TimeRange"`
	FrameTimes       []int64  `xml:"FrameTimes>FrameTime,omitempty"`
	VideoStreamIndex *int64   `xml:"VideoStreamIndex,omitempty"`
	Labels           []Label  `xml:"Labels>Label,omitempty"`
}

// Address represents a geographic address
type Address struct {
	XMLName      xml.Name `xml:"Address"`
	Country      *string  `xml:"Country,omitempty"`
	Province     *string  `xml:"Province,omitempty"`
	City         *string  `xml:"City,omitempty"`
	District     *string  `xml:"District,omitempty"`
	Town         *string  `xml:"Town,omitempty"`
	Street       *string  `xml:"Street,omitempty"`
	StreetNumber *string  `xml:"StreetNumber,omitempty"`
	PostalCode   *string  `xml:"PostalCode,omitempty"`
}

// HeadPose represents head pose information
type HeadPose struct {
	XMLName xml.Name `xml:"HeadPose"`
	Pitch   *float64 `xml:"Pitch,omitempty"`
	Yaw     *float64 `xml:"Yaw,omitempty"`
	Roll    *float64 `xml:"Roll,omitempty"`
}

// ImageInsight represents image insight information
type ImageInsight struct {
	XMLName     xml.Name `xml:"Image"`
	Caption     *string  `xml:"Caption,omitempty"`
	Description *string  `xml:"Description,omitempty"`
}

// VideoInsight represents video insight information
type VideoInsight struct {
	XMLName     xml.Name `xml:"Video"`
	Caption     *string  `xml:"Caption,omitempty"`
	Description *string  `xml:"Description,omitempty"`
}

// Insights represents comprehensive media insights
type Insights struct {
	XMLName xml.Name      `xml:"Insights"`
	Video   *VideoInsight `xml:"Video,omitempty"`
	Image   *ImageInsight `xml:"Image,omitempty"`
}

// AggregationGroup represents a group of aggregation results
type AggregationGroup struct {
	XMLName xml.Name `xml:"Group"`
	Value   *string  `xml:"Value,omitempty"`
	Count   *int64   `xml:"Count,omitempty"`
}

// Aggregation represents an aggregation result
type Aggregation struct {
	XMLName           xml.Name           `xml:"Aggregation"`
	Operation         *string            `xml:"Operation,omitempty"`
	Field             *string            `xml:"Field,omitempty"`
	Value             *string            `xml:"Value,omitempty"`
	AggregationGroups []AggregationGroup `xml:"Groups>Group,omitempty"`
}

// CroppingSuggestion represents an image cropping suggestion
type CroppingSuggestion struct {
	XMLName     xml.Name  `xml:"CroppingSuggestion"`
	AspectRatio *string   `xml:"AspectRatio,omitempty"`
	Confidence  *float64  `xml:"Confidence,omitempty"`
	Boundary    *Boundary `xml:"Boundary,omitempty"`
}

// ElementRelation represents a relation between elements
type ElementRelation struct {
	XMLName  xml.Name `xml:"ElementRelation"`
	ObjectId *string  `xml:"ObjectId,omitempty"`
	Type     *string  `xml:"RelationType,omitempty"`
}

// SimpleQueryRequest defines the request for simple query operation
type SimpleQueryRequest struct {
	Bucket           *string `input:"host,bucket,required"`
	DatasetName      *string `input:"query,datasetName,required"`
	NextToken        *string `input:"query,nextToken"`
	MaxResults       *int32  `input:"query,maxResults"`
	Query            *string `input:"query,query"`
	Sort             *string `input:"query,sort"`
	Order            *string `input:"query,order"`
	Aggregations     *string `input:"query,aggregations"`
	WithFields       *string `input:"query,withFields"`
	WithoutTotalHits *bool   `input:"query,withoutTotalHits"`
	oss.RequestCommon
}

// SimpleQueryResult defines the result for SimpleQuery operation
type SimpleQueryResult struct {
	XMLName      xml.Name      `xml:"SimpleQueryResponse"`
	Files        []File        `xml:"Files>File,omitempty"`
	NextToken    *string       `xml:"NextToken,omitempty"`
	MaxResults   *int32        `xml:"MaxResults,omitempty"`
	TotalHits    *int64        `xml:"TotalHits,omitempty"`
	Aggregations []Aggregation `xml:"Aggregations>Aggregation,omitempty"`
	oss.ResultCommon
}

// SimpleQuery queries files in a dataset using structured query language.
func (c *Client) SimpleQuery(ctx context.Context, request *SimpleQueryRequest, optFns ...func(*oss.Options)) (*SimpleQueryResult, error) {
	var err error
	if request == nil {
		request = &SimpleQueryRequest{}
	}

	input := &oss.OperationInput{
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

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &SimpleQueryResult{}

	if err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	}); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

// SemanticQueryRequest defines the request for semantic query operation
type SemanticQueryRequest struct {
	Bucket          *string `input:"host,bucket,required"`
	DatasetName     *string `input:"query,datasetName,required"`
	NextToken       *string `input:"query,nextToken"`
	MaxResults      *int32  `input:"query,maxResults"`
	Query           *string `input:"query,query"`
	SimpleQuery     *string `input:"query,simpleQuery"`
	WithFields      *string `input:"query,withFields"`
	MediaTypes      *string `input:"query,mediaTypes"`
	SourceUri       *string `input:"query,sourceURI"`
	SmartClusterIds *string `input:"query,smartClusterIds"`
	oss.RequestCommon
}

// File represents a file result in query
type File struct {
	XMLName                               xml.Name             `xml:"File"`
	OwnerId                               *string              `xml:"OwnerId,omitempty"`
	DatasetName                           *string              `xml:"DatasetName,omitempty"`
	ObjectType                            *string              `xml:"ObjectType,omitempty"`
	ObjectId                              *string              `xml:"ObjectId,omitempty"`
	UpdateTime                            *string              `xml:"UpdateTime,omitempty"`
	CreateTime                            *string              `xml:"CreateTime,omitempty"`
	URI                                   *string              `xml:"URI,omitempty"`
	OSSURI                                *string              `xml:"OSSURI,omitempty"`
	Filename                              *string              `xml:"Filename,omitempty"`
	MediaType                             *string              `xml:"MediaType,omitempty"`
	ContentType                           *string              `xml:"ContentType,omitempty"`
	Size                                  *int64               `xml:"Size,omitempty"`
	FileHash                              *string              `xml:"FileHash,omitempty"`
	FileModifiedTime                      *string              `xml:"FileModifiedTime,omitempty"`
	FileCreateTime                        *string              `xml:"FileCreateTime,omitempty"`
	FileAccessTime                        *string              `xml:"FileAccessTime,omitempty"`
	ProduceTime                           *string              `xml:"ProduceTime,omitempty"`
	LatLong                               *string              `xml:"LatLong,omitempty"`
	Timezone                              *string              `xml:"Timezone,omitempty"`
	Addresses                             []Address            `xml:"Addresses>Address"`
	TravelClusterId                       *string              `xml:"TravelClusterId,omitempty"`
	Orientation                           *int64               `xml:"Orientation,omitempty"`
	Figures                               []Figure             `xml:"Figures>Figure,omitempty"`
	FigureCount                           *int64               `xml:"FigureCount,omitempty"`
	Labels                                []Label              `xml:"Labels>Label,omitempty"`
	Title                                 *string              `xml:"Title,omitempty"`
	ImageWidth                            *int64               `xml:"ImageWidth,omitempty"`
	ImageHeight                           *int64               `xml:"ImageHeight,omitempty"`
	EXIF                                  *string              `xml:"EXIF,omitempty"`
	ImageScore                            *ImageScore          `xml:"ImageScore,omitempty"`
	CroppingSuggestions                   []CroppingSuggestion `xml:"CroppingSuggestions>CroppingSuggestion,omitempty"`
	OCRContents                           []OCRContents        `xml:"OCRContents>OCRContents,omitempty"`
	VideoWidth                            *int64               `xml:"VideoWidth,omitempty"`
	VideoHeight                           *int64               `xml:"VideoHeight,omitempty"`
	VideoStreams                          []VideoStream        `xml:"VideoStreams>VideoStream,omitempty"`
	Subtitles                             []SubtitleStream     `xml:"Subtitles>Subtitle,omitempty"`
	AudioStreams                          []AudioStream        `xml:"AudioStreams>AudioStream,omitempty"`
	Artist                                *string              `xml:"Artist,omitempty"`
	AlbumArtist                           *string              `xml:"AlbumArtist,omitempty"`
	AudioCovers                           []Image              `xml:"AudioCovers>AudioCover,omitempty"`
	Composer                              *string              `xml:"Composer,omitempty"`
	Performer                             *string              `xml:"Performer,omitempty"`
	Language                              *string              `xml:"Language,omitempty"`
	Album                                 *string              `xml:"Album,omitempty"`
	PageCount                             *int64               `xml:"PageCount,omitempty"`
	ETag                                  *string              `xml:"ETag,omitempty"`
	CacheControl                          *string              `xml:"CacheControl,omitempty"`
	ContentDisposition                    *string              `xml:"ContentDisposition,omitempty"`
	ContentEncoding                       *string              `xml:"ContentEncoding,omitempty"`
	ContentLanguage                       *string              `xml:"ContentLanguage,omitempty"`
	AccessControlAllowOrigin              *string              `xml:"AccessControlAllowOrigin,omitempty"`
	AccessControlRequestMethod            *string              `xml:"AccessControlRequestMethod,omitempty"`
	ServerSideEncryptionCustomerAlgorithm *string              `xml:"ServerSideEncryptionCustomerAlgorithm,omitempty"`
	ServerSideEncryption                  *string              `xml:"ServerSideEncryption,omitempty"`
	ServerSideDataEncryption              *string              `xml:"ServerSideDataEncryption,omitempty"`
	ServerSideEncryptionKeyId             *string              `xml:"ServerSideEncryptionKeyId,omitempty"`
	OSSStorageClass                       *string              `xml:"OSSStorageClass,omitempty"`
	OSSCRC64                              *string              `xml:"OSSCRC64,omitempty"`
	ObjectACL                             *string              `xml:"ObjectACL,omitempty"`
	ContentMd5                            *string              `xml:"ContentMd5,omitempty"`
	SequenceNumber                        *int64               `xml:"SequenceNumber,omitempty"`
	SemanticSimilarity                    *float64             `xml:"SemanticSimilarity,omitempty"`
	OSSUserMeta                           MapEntry             `xml:"OSSUserMeta,omitempty"`
	OSSTaggingCount                       *int64               `xml:"OSSTaggingCount,omitempty"`
	OSSTagging                            MapEntry             `xml:"OSSTagging,omitempty"`
	OSSExpiration                         *string              `xml:"OSSExpiration,omitempty"`
	OSSVersionId                          *string              `xml:"OSSVersionId,omitempty"`
	OSSDeleteMarker                       *string              `xml:"OSSDeleteMarker,omitempty"`
	OSSObjectType                         *string              `xml:"OSSObjectType,omitempty"`
	CustomId                              *string              `xml:"CustomId,omitempty"`
	CustomLabels                          MapEntry             `xml:"CustomLabels,omitempty"`
	StreamCount                           *int64               `xml:"StreamCount,omitempty"`
	ProgramCount                          *int64               `xml:"ProgramCount,omitempty"`
	FormatName                            *string              `xml:"FormatName,omitempty"`
	FormatLongName                        *string              `xml:"FormatLongName,omitempty"`
	StartTime                             *float64             `xml:"StartTime,omitempty"`
	Bitrate                               *int64               `xml:"Bitrate,omitempty"`
	Duration                              *float64             `xml:"Duration,omitempty"`
	SemanticTypes                         []string             `xml:"SemanticTypes>SemanticType,omitempty"`
	Elements                              []Element            `xml:"Elements>Element,omitempty"`
	SceneElements                         []SceneElement       `xml:"SceneElements>SceneElement,omitempty"`
	OCRTexts                              *string              `xml:"OCRTexts,omitempty"`
	Reason                                *string              `xml:"Reason,omitempty"`
	ObjectStatus                          *string              `xml:"ObjectStatus,omitempty"`
	Insights                              *Insights            `xml:"Insights,omitempty"`
}

// SemanticQueryResult defines the result for SemanticQuery operation
type SemanticQueryResult struct {
	XMLName    xml.Name `xml:"SemanticQueryResponse"`
	Files      []File   `xml:"Files>File,omitempty"`
	NextToken  *string  `xml:"NextToken,omitempty"`
	MaxResults *int32   `xml:"MaxResults,omitempty"`
	TotalCount *int64   `xml:"TotalCount,omitempty"`
	oss.ResultCommon
}

// SemanticQuery queries files in a dataset using natural language.
func (c *Client) SemanticQuery(ctx context.Context, request *SemanticQueryRequest, optFns ...func(*oss.Options)) (*SemanticQueryResult, error) {
	var err error
	if request == nil {
		request = &SemanticQueryRequest{}
	}

	input := &oss.OperationInput{
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

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &SemanticQueryResult{}

	if err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	}); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}
