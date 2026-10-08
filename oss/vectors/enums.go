package vectors

// FieldType The field type for vector FieldSchema
type FieldType string

// Enum values for FieldType
const (
	FieldTypeVector   FieldType = "vector"
	FieldTypeDouble   FieldType = "double"
	FieldTypeLong     FieldType = "long"
	FieldTypeBool     FieldType = "bool"
	FieldTypeString   FieldType = "string"
	FieldTypeIp       FieldType = "ip"
	FieldTypeGeoPoint FieldType = "geoPoint"
)

// VectorDataType The vector data type for vector data
type VectorDataType string

// VectorDataTypeFloat32 Enum values for VectorDataType
const (
	VectorDataTypeFloat32 VectorDataType = "float32"
)

// DistanceMetricType The distance metric type for DistanceMetric
type DistanceMetricType string

// Enum values for DistanceMetricType
const (
	DistanceMetricTypeEuclidean    DistanceMetricType = "euclidean"
	DistanceMetricTypeCosine       DistanceMetricType = "cosine"
	DistanceMetricTypeInnerProduct DistanceMetricType = "ip"
)

// SortOrderType The sort order of a sort field
type SortOrderType string

// Enum values for SortOrderType
const (
	SortOrderTypeAsc  SortOrderType = "asc"
	SortOrderTypeDesc SortOrderType = "desc"
)

// NormalizerType The score normalizer of a weight compound retriever component.
type NormalizerType string

// Enum values for NormalizerType
const (
	NormalizerTypeNone   NormalizerType = "none"
	NormalizerTypeMinMax NormalizerType = "minMax"
	NormalizerTypeL2     NormalizerType = "l2"
)

// IndexModeType The mode of a vector index.
type IndexModeType string

// Enum values for IndexModeType
const (
	IndexModeTypeStandard IndexModeType = "standard"
	IndexModeTypeFusion   IndexModeType = "fusion"
)

// AnalyzerType The analyzer type used for the full text search of a string field.
type AnalyzerType string

// Enum values for AnalyzerType
const (
	AnalyzerTypeStandard AnalyzerType = "standard"
	AnalyzerTypeSplit    AnalyzerType = "split"
)
