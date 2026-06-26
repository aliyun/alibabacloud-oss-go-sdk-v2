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
	ClipURI   *string  `xml:"ClipURI"`
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
	XMLName xml.Name `xml:"AggregationGroup"`
	Value   *string  `xml:"Value,omitempty"`
	Count   *int64   `xml:"Count,omitempty"`
}

// Aggregation represents an aggregation result
type Aggregation struct {
	XMLName           xml.Name           `xml:"Aggregation"`
	Operation         *string            `xml:"Operation,omitempty"`
	Field             *string            `xml:"Field,omitempty"`
	Value             *string            `xml:"Value,omitempty"`
	AggregationGroups []AggregationGroup `xml:"Groups>AggregationGroup,omitempty"`
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
	Bucket      *string `input:"host,bucket,required"`
	DatasetName *string `input:"query,datasetName,required"`
	MaxResults  *int32  `input:"query,maxResults"`
	Query       *string `input:"query,query"`
	SimpleQuery *string `input:"query,simpleQuery"`
	WithFields  *string `input:"query,withFields"`
	MediaTypes  *string `input:"query,mediaTypes"`
	SourceUri   *string `input:"query,sourceURI"`
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
	OSSUserMeta                           []UserMeta           `xml:"OSSUserMeta>UserMeta,omitempty"`
	OSSTaggingCount                       *int64               `xml:"OSSTaggingCount,omitempty"`
	OSSTagging                            []Tagging            `xml:"OSSTagging>Tagging,omitempty"`
	OSSExpiration                         *string              `xml:"OSSExpiration,omitempty"`
	OSSVersionId                          *string              `xml:"OSSVersionId,omitempty"`
	OSSDeleteMarker                       *string              `xml:"OSSDeleteMarker,omitempty"`
	OSSObjectType                         *string              `xml:"OSSObjectType,omitempty"`
	CustomId                              *string              `xml:"CustomId,omitempty"`
	CustomLabels                          []CustomLabel        `xml:"CustomLabels>Item,omitempty"`
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

type Tagging struct {
	Key   *string `xml:"Key,omitempty"`
	Value *string `xml:"Value,omitempty"`
}

type UserMeta struct {
	Key   *string `xml:"Key,omitempty"`
	Value *string `xml:"Value,omitempty"`
}

type CustomLabel struct {
	Key   *string `xml:"Key,omitempty"`
	Value *string `xml:"Value,omitempty"`
}

// SemanticQueryResult defines the result for SemanticQuery operation
type SemanticQueryResult struct {
	Files []File `xml:"Files>File,omitempty"`

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

type OpenMetaQueryRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	Mode *string `input:"query,mode,required"`

	Role *string `input:"query,role"`

	MetaQuery *OpenMetaQuery `input:"body,MetaQuery,xml"`

	oss.RequestCommon
}

type OpenMetaQuery struct {
	WorkflowParameters *WorkflowParameters `xml:"WorkflowParameters,omitempty"`

	Filters *Filters `xml:"Filters,omitempty"`

	NotificationAttributes *NotificationAttributes `xml:"NotificationAttributes,omitempty"`

	DatasetConfig *DatasetConfig `xml:"DatasetConfig,omitempty"`

	IndexOptions *IndexOptions `xml:"IndexOptions,omitempty"`

	RouteRule *RouteRule `xml:"RouteRule,omitempty"`
}

type Filters struct {
	Filter []string `xml:"Filter,omitempty"`
}

type NotificationAttributes struct {
	Notifications *Notifications `xml:"Notifications,omitempty"`
	WithFields    *WithFields    `xml:"WithFields,omitempty"`
}

type Notifications struct {
	Notification []Notification `xml:"Notification,omitempty"`
}

type WithFields struct {
	WithField []string `xml:"WithField,omitempty"`
}

type Notification struct {
	MNS *string `xml:"MNS,omitempty"`
}

type IndexOptions struct {
	IgnoreObjectDelete *bool         `xml:"IgnoreObjectDelete,omitempty"`
	IgnoreEvents       *IgnoreEvents `xml:"IgnoreEvents,omitempty"`
}

type IgnoreEvents struct {
	IgnoreEvent []string `xml:"IgnoreEvent,omitempty"`
}

type RouteRule struct {
	Type              *string `xml:"Type,omitempty"`
	AutoCreateDataset *bool   `xml:"AutoCreateDataset,omitempty"`
	OSSTagKey         *string `xml:"OSSTagKey,omitempty"`
}

type OpenMetaQueryResult struct {
	oss.ResultCommon
}

// OpenMetaQuery Enables metadata management for a bucket. After you enable the metadata management feature for a bucket, Object Storage Service (OSS) creates a metadata index library for the bucket and creates metadata indexes for all objects in the bucket. After the metadata index library is created, OSS continues to perform quasi-real-time scans on incremental objects in the bucket and creates metadata indexes for the incremental objects.
func (c *Client) OpenMetaQuery(ctx context.Context, request *OpenMetaQueryRequest, optFns ...func(*oss.Options)) (*OpenMetaQueryResult, error) {
	var err error
	if request == nil {
		request = &OpenMetaQueryRequest{}
	}
	input := &oss.OperationInput{
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
	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &OpenMetaQueryResult{}

	if err = c.client.UnmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetMetaQueryStatusRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

type GetMetaQueryStatusResult struct {
	State                  *string                 `xml:"State,omitempty"`
	Phase                  *string                 `xml:"Phase,omitempty"`
	CreateTime             *string                 `xml:"CreateTime,omitempty"`
	UpdateTime             *string                 `xml:"UpdateTime,omitempty"`
	MetaQueryMode          *string                 `xml:"MetaQueryMode,omitempty"`
	WorkflowParameters     *WorkflowParameters     `xml:"WorkflowParameters,omitempty"`
	Filters                *Filters                `xml:"Filters,omitempty"`
	IndexOptions           *IndexOptions           `xml:"IndexOptions,omitempty"`
	RouteRule              *RouteRule              `xml:"RouteRule,omitempty"`
	NotificationAttributes *NotificationAttributes `xml:"NotificationAttributes,omitempty"`
	DatasetConfig          *DatasetConfig          `xml:"DatasetConfig,omitempty"`

	oss.ResultCommon
}

// GetMetaQueryStatus Queries the information about the metadata index library of a bucket.
func (c *Client) GetMetaQueryStatus(ctx context.Context, request *GetMetaQueryStatusRequest, optFns ...func(*oss.Options)) (*GetMetaQueryStatusResult, error) {
	var err error
	if request == nil {
		request = &GetMetaQueryStatusRequest{}
	}
	input := &oss.OperationInput{
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
	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetMetaQueryStatusResult{}

	if err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	}); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type DoMetaQueryRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	Mode *string `input:"query,mode,required"`

	// The request body schema.
	MetaQuery *DoMetaQuery `input:"body,MetaQuery,xml,required"`

	oss.RequestCommon
}

type DoMetaQuery struct {
	// The maximum number of objects to return. Valid values: 0 to 100. If this parameter is not set or is set to 0, up to 100 objects are returned.
	MaxResults *int64 `xml:"MaxResults"`

	// The query conditions. A query condition includes the following elements:*   Operation: the operator. Valid values: eq (equal to), gt (greater than), gte (greater than or equal to), lt (less than), lte (less than or equal to), match (fuzzy query), prefix (prefix query), and (AND), or (OR), and not (NOT).*   Field: the field name.*   Value: the field value.*   SubQueries: the subquery conditions. Options that are included in this element are the same as those of simple query. You need to set subquery conditions only when Operation is set to and, or, or not.
	Query *string `xml:"Query"`

	// The field based on which the results are sorted.
	Sort *string `xml:"Sort"`

	// The sort order.
	Order *MetaQueryOrderType `xml:"Order"`

	// The container that stores the information about aggregate operations.
	Aggregations *MetaQueryAggregations `xml:"Aggregations"`

	// The pagination token used to obtain information in the next request. The object information is returned in alphabetical order starting from the value of NextToken.
	NextToken *string `xml:"NextToken"`

	// The container that stores the type of multimedia.
	MediaTypes *MetaQueryMediaTypes `xml:"MediaTypes"`

	//The query conditions
	SimpleQuery *string `xml:"SimpleQuery"`

	WithoutTotalHits *string `xml:"WithoutTotalHits"`

	SourceURI *string `xml:"SourceURI"`

	SmartClusterIds *SmartClusterIds `xml:"SmartClusterIds"`

	WithFields *WithFields `xml:"WithFields,omitempty"`
}

type SmartClusterIds struct {
	SmartClusterId []string `xml:"SmartClusterId"`
}

type MetaQueryMediaTypes struct {
	// The type of multimedia that you want to query. Valid values: image, video, audio, document
	MediaTypes []string `xml:"MediaType"`
}

type MetaQueryAggregations struct {
	// The container that stores the information about a single aggregate operation.
	Aggregations []Aggregation `xml:"Aggregation"`
}

type MetaQueryGroups struct {
	// The grouped aggregations.
	Groups []MetaQueryGroup `xml:"Group"`
}

type MetaQueryGroup struct {
	// The value for the grouped aggregation.
	Value *string `xml:"Value"`

	// The number of results in the grouped aggregation.
	Count *int64 `xml:"Count"`
}

type DoMetaQueryResult struct {
	// The token that is used for the next query when the total number of objects exceeds the value of MaxResults.The value of NextToken is used to return the unreturned results in the next query.This parameter has a value only when not all objects are returned.
	NextToken *string `xml:"NextToken"`

	TotalHits *int64 `xml:"TotalHits"`

	// The list of file information.
	Files []File `xml:"Files>File"`

	// The list of file information.
	Aggregations []Aggregation `xml:"Aggregations>Aggregation"`

	oss.ResultCommon
}

// DoMetaQuery Queries the objects in a bucket that meet the specified conditions by using the data indexing feature. The information about the objects is listed based on the specified fields and sorting methods.
func (c *Client) DoMetaQuery(ctx context.Context, request *DoMetaQueryRequest, optFns ...func(*oss.Options)) (*DoMetaQueryResult, error) {
	var err error
	if request == nil {
		request = &DoMetaQueryRequest{}
	}
	input := &oss.OperationInput{
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
	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DoMetaQueryResult{}

	if err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	}); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}

type CloseMetaQueryRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	oss.RequestCommon
}

type CloseMetaQueryResult struct {
	oss.ResultCommon
}

// CloseMetaQuery Disables the metadata management feature for an Object Storage Service (OSS) bucket. After the metadata management feature is disabled for a bucket, OSS automatically deletes the metadata index library of the bucket and you cannot perform metadata indexing.
func (c *Client) CloseMetaQuery(ctx context.Context, request *CloseMetaQueryRequest, optFns ...func(*oss.Options)) (*CloseMetaQueryResult, error) {
	var err error
	if request == nil {
		request = &CloseMetaQueryRequest{}
	}
	input := &oss.OperationInput{
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
	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &CloseMetaQueryResult{}

	if err = c.client.UnmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
