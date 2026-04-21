package dataprocess

import (
	"context"
	"encoding/xml"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// Image represents image information
type Image struct {
	XMLName    xml.Name `xml:"Image"`
	Width      *int64   `xml:"Width,omitempty"`
	Height     *int64   `xml:"Height,omitempty"`
	Format     *string  `xml:"Format,omitempty"`
	ColorSpace *string  `xml:"ColorSpace,omitempty"`
	Exif       *string  `xml:"Exif,omitempty"`
}

// VideoStream represents video stream information
type VideoStream struct {
	XMLName       xml.Name `xml:"VideoStream"`
	Index         *int32   `xml:"Index,omitempty"`
	CodecName     *string  `xml:"CodecName,omitempty"`
	CodecTimeBase *string  `xml:"CodecTimeBase,omitempty"`
	CodecTag      *string  `xml:"CodecTag,omitempty"`
	Width         *int64   `xml:"Width,omitempty"`
	Height        *int64   `xml:"Height,omitempty"`
	FrameRate     *string  `xml:"FrameRate,omitempty"`
	Duration      *string  `xml:"Duration,omitempty"`
	BitRate       *int64   `xml:"BitRate,omitempty"`
}

// AudioStream represents audio stream information
type AudioStream struct {
	XMLName       xml.Name `xml:"AudioStream"`
	Index         *int32   `xml:"Index,omitempty"`
	CodecName     *string  `xml:"CodecName,omitempty"`
	CodecTimeBase *string  `xml:"CodecTimeBase,omitempty"`
	CodecTag      *string  `xml:"CodecTag,omitempty"`
	SampleRate    *string  `xml:"SampleRate,omitempty"`
	Channels      *int32   `xml:"Channels,omitempty"`
	Duration      *string  `xml:"Duration,omitempty"`
	BitRate       *int64   `xml:"BitRate,omitempty"`
}

// SubtitleStream represents subtitle stream information
type SubtitleStream struct {
	XMLName   xml.Name `xml:"SubtitleStream"`
	Index     *int32   `xml:"Index,omitempty"`
	CodecName *string  `xml:"CodecName,omitempty"`
	Language  *string  `xml:"Language,omitempty"`
	Duration  *string  `xml:"Duration,omitempty"`
}

// Label represents a label with confidence score
type Label struct {
	XMLName     xml.Name `xml:"Label"`
	Name        *string  `xml:"Name,omitempty"`
	Score       *float64 `xml:"Score,omitempty"`
	ParentLabel *string  `xml:"ParentLabel,omitempty"`
}

// OCRContents represents OCR recognition results
type OCRContents struct {
	XMLName xml.Name `xml:"OCRContents"`
	Text    *string  `xml:"Text,omitempty"`
}

// ImageScore represents image quality score
type ImageScore struct {
	XMLName     xml.Name `xml:"ImageScore"`
	Clarity     *float64 `xml:"Clarity,omitempty"`
	Composition *float64 `xml:"Composition,omitempty"`
}

// Boundary represents a rectangular boundary
type Boundary struct {
	XMLName xml.Name `xml:"Boundary"`
	Left    *int64   `xml:"Left,omitempty"`
	Top     *int64   `xml:"Top,omitempty"`
	Width   *int64   `xml:"Width,omitempty"`
	Height  *int64   `xml:"Height,omitempty"`
}

// PointInt64 represents a 2D point with int64 coordinates
type PointInt64 struct {
	XMLName xml.Name `xml:"PointInt64"`
	X       *int64   `xml:"X,omitempty"`
	Y       *int64   `xml:"Y,omitempty"`
}

// Figure represents a figure/shape in an image
type Figure struct {
	XMLName  xml.Name     `xml:"Figure"`
	Type     *string      `xml:"Type,omitempty"`
	Points   []PointInt64 `xml:"Points,omitempty"`
	Boundary *Boundary    `xml:"Boundary,omitempty"`
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
	XMLName    xml.Name        `xml:"Element"`
	Type       *string         `xml:"Type,omitempty"`
	SubType    *string         `xml:"SubType,omitempty"`
	Confidence *float64        `xml:"Confidence,omitempty"`
	Boundary   *Boundary       `xml:"Boundary,omitempty"`
	Content    *ElementContent `xml:"ElementContent,omitempty"`
	OCRContent *string         `xml:"OCRContent,omitempty"`
	FaceId     *string         `xml:"FaceId,omitempty"`
}

// SceneElement represents an element in a scene
type SceneElement struct {
	XMLName    xml.Name `xml:"SceneElement"`
	Type       *string  `xml:"Type,omitempty"`
	Confidence *float64 `xml:"Confidence,omitempty"`
	Labels     []Label  `xml:"Labels,omitempty"`
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
	XMLName     xml.Name      `xml:"ImageInsight"`
	Labels      []Label       `xml:"Labels,omitempty"`
	OCRContents *OCRContents  `xml:"OCRContents,omitempty"`
	ImageScore  *ImageScore   `xml:"ImageScore,omitempty"`
	Elements    []Element     `xml:"Elements,omitempty"`
}

// VideoInsight represents video insight information
type VideoInsight struct {
	XMLName       xml.Name       `xml:"VideoInsight"`
	Labels        []Label        `xml:"Labels,omitempty"`
	SceneElements []SceneElement `xml:"SceneElements,omitempty"`
}

// Insights represents comprehensive media insights
type Insights struct {
	XMLName       xml.Name      `xml:"Insights"`
	ImageInsight  *ImageInsight `xml:"ImageInsight,omitempty"`
	VideoInsight  *VideoInsight `xml:"VideoInsight,omitempty"`
	AudioDuration *string       `xml:"AudioDuration,omitempty"`
}

// AggregationInfo represents aggregation information
type AggregationInfo struct {
	XMLName xml.Name `xml:"AggregationInfo"`
	Type    *string  `xml:"Type,omitempty"`
	Field   *string  `xml:"Field,omitempty"`
	Value   *string  `xml:"Value,omitempty"`
	Count   *int64   `xml:"Count,omitempty"`
}

// AggregationGroup represents a group of aggregation results
type AggregationGroup struct {
	XMLName          xml.Name          `xml:"AggregationGroup"`
	AggregationInfos []AggregationInfo `xml:"AggregationInfos,omitempty"`
}

// Aggregation represents an aggregation result
type Aggregation struct {
	XMLName           xml.Name           `xml:"Aggregation"`
	Type              *string            `xml:"Type,omitempty"`
	Field             *string            `xml:"Field,omitempty"`
	Value             *string            `xml:"Value,omitempty"`
	Count             *int64             `xml:"Count,omitempty"`
	AggregationGroups []AggregationGroup `xml:"AggregationGroups,omitempty"`
}

// CroppingSuggestion represents an image cropping suggestion
type CroppingSuggestion struct {
	XMLName  xml.Name  `xml:"CroppingSuggestion"`
	Priority *int32    `xml:"Priority,omitempty"`
	Boundary *Boundary `xml:"Boundary,omitempty"`
	Score    *float64  `xml:"Score,omitempty"`
}

// Clip represents a media clip
type Clip struct {
	XMLName   xml.Name `xml:"Clip"`
	StartTime *string  `xml:"StartTime,omitempty"`
	EndTime   *string  `xml:"EndTime,omitempty"`
	URI       *string  `xml:"URI,omitempty"`
}

// ElementRelation represents a relation between elements
type ElementRelation struct {
	XMLName      xml.Name `xml:"ElementRelation"`
	SubjectId    *string  `xml:"SubjectId,omitempty"`
	ObjectId     *string  `xml:"ObjectId,omitempty"`
	RelationType *string  `xml:"RelationType,omitempty"`
	Confidence   *float64 `xml:"Confidence,omitempty"`
}

// File represents a file in the dataset
type File struct {
	XMLName         xml.Name         `xml:"File"`
	URI             *string          `xml:"URI,omitempty"`
	CustomContent   *string          `xml:"CustomContent,omitempty"`
	FileSize        *int64           `xml:"FileSize,omitempty"`
	CreateTime      *string          `xml:"CreateTime,omitempty"`
	LastModified    *string          `xml:"LastModified,omitempty"`
	Image           *Image           `xml:"Image,omitempty"`
	VideoStreams    []VideoStream    `xml:"VideoStreams,omitempty"`
	AudioStreams    []AudioStream    `xml:"AudioStreams,omitempty"`
	SubtitleStreams []SubtitleStream `xml:"SubtitleStreams,omitempty"`
	Duration        *string          `xml:"Duration,omitempty"`
	Labels          []Label          `xml:"Labels,omitempty"`
	Elements        []Element        `xml:"Elements,omitempty"`
	Insights        *Insights        `xml:"Insights,omitempty"`
}

// SimpleQuery represents a simple query structure
type SimpleQuery struct {
	XMLName xml.Name `xml:"SimpleQuery"`
	Query   *string  `xml:"Query,omitempty"`
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

// SimpleQueryFile represents a file result in simple query
type SimpleQueryFile struct {
	XMLName       xml.Name  `xml:"File"`
	URI           *string   `xml:"URI,omitempty"`
	CustomContent *string   `xml:"CustomContent,omitempty"`
	FileSize      *int64    `xml:"FileSize,omitempty"`
	CreateTime    *string   `xml:"CreateTime,omitempty"`
	LastModified  *string   `xml:"LastModified,omitempty"`
	Image         *Image    `xml:"Image,omitempty"`
	Duration      *string   `xml:"Duration,omitempty"`
	Labels        []Label   `xml:"Labels,omitempty"`
	Elements      []Element `xml:"Elements,omitempty"`
	Insights      *Insights `xml:"Insights,omitempty"`
}

// SimpleQueryResult defines the result for SimpleQuery operation
type SimpleQueryResult struct {
	XMLName      xml.Name          `xml:"SimpleQueryResult"`
	Files        []SimpleQueryFile `xml:"Files>File,omitempty"`
	NextToken    *string           `xml:"NextToken,omitempty"`
	MaxResults   *int32            `xml:"MaxResults,omitempty"`
	TotalCount   *int64            `xml:"TotalCount,omitempty"`
	Aggregations []Aggregation     `xml:"Aggregations>Aggregation,omitempty"`
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
	NextToken   *string `input:"query,nextToken"`
	MaxResults  *int32  `input:"query,maxResults"`
	Query       *string `input:"query,query"`
	WithFields  *string `input:"query,withFields"`
	MediaTypes  *string `input:"query,mediaTypes"`
	SourceUri   *string `input:"query,sourceURI"`
	oss.RequestCommon
}

// SemanticQueryFile represents a file result in semantic query
type SemanticQueryFile struct {
	XMLName       xml.Name  `xml:"File"`
	URI           *string   `xml:"URI,omitempty"`
	CustomContent *string   `xml:"CustomContent,omitempty"`
	FileSize      *int64    `xml:"FileSize,omitempty"`
	CreateTime    *string   `xml:"CreateTime,omitempty"`
	LastModified  *string   `xml:"LastModified,omitempty"`
	Image         *Image    `xml:"Image,omitempty"`
	Duration      *string   `xml:"Duration,omitempty"`
	Labels        []Label   `xml:"Labels,omitempty"`
	Elements      []Element `xml:"Elements,omitempty"`
	Insights      *Insights `xml:"Insights,omitempty"`
	Score         *float64  `xml:"Score,omitempty"`
}

// SemanticQueryResult defines the result for SemanticQuery operation
type SemanticQueryResult struct {
	XMLName    xml.Name            `xml:"SemanticQueryResult"`
	Files      []SemanticQueryFile `xml:"Files>File,omitempty"`
	NextToken  *string             `xml:"NextToken,omitempty"`
	MaxResults *int32              `xml:"MaxResults,omitempty"`
	TotalCount *int64              `xml:"TotalCount,omitempty"`
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
