package vectors

// FieldType The field type for vector SchemaField
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

type SortOrderType string

// Enum values for SortOrderType
const (
	SortOrderTypeAsc  SortOrderType = "asc"
	SortOrderTypeDesc SortOrderType = "desc"
)
