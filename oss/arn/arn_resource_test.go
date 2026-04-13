package arn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArnResourceBuilder_WithAllFields(t *testing.T) {
	resourceType := "bucket"
	qualifier := "obj"
	resource, err := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Qualifier(&qualifier).
		Build()

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.Equal(t, "obj", *resource.Qualifier())
}

func TestArnResourceBuilder_WithOnlyResource(t *testing.T) {
	resource, err := ArnResourceBuilderNew().
		Resource("my-bucket").
		Build()

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.Nil(t, resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.Nil(t, resource.Qualifier())
}

func TestArnResourceBuilder_WithResourceAndQualifier(t *testing.T) {
	qualifier := "obj"
	resource, err := ArnResourceBuilderNew().
		Resource("my-bucket").
		Qualifier(&qualifier).
		Build()

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.Nil(t, resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.Equal(t, "obj", *resource.Qualifier())
}

func TestArnResourceBuilder_WithoutResource_ThrowsError(t *testing.T) {
	resource, err := ArnResourceBuilderNew().
		Build()

	assert.Nil(t, resource)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "resource must not be blank")
}

func TestArnResourceBuilder_WithBlankResource_ThrowsError(t *testing.T) {
	resource, err := ArnResourceBuilderNew().
		Resource("   ").
		Build()

	assert.Nil(t, resource)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "resource must not be blank")
}

func TestArnResourceFromString_SimpleResource(t *testing.T) {
	resource, err := ArnResourceFromString("my-bucket")

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.Nil(t, resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.Nil(t, resource.Qualifier())
}

func TestArnResourceFromString_WithResourceTypeAndResource(t *testing.T) {
	resource, err := ArnResourceFromString("bucket:my-bucket")

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.Nil(t, resource.Qualifier())
}

func TestArnResourceFromString_WithResourceTypeResourceAndQualifier(t *testing.T) {
	resource, err := ArnResourceFromString("bucket:my-bucket:obj")

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.NotNil(t, resource.Qualifier())
	assert.Equal(t, "obj", *resource.Qualifier())
}

func TestArnResourceFromString_WithSlashSplitter(t *testing.T) {
	resource, err := ArnResourceFromString("bucket/my-bucket/obj")

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.NotNil(t, resource.Qualifier())
	assert.Equal(t, "obj", *resource.Qualifier())
}

func TestArnResourceFromString_WithSlashSplitterEmptyQualifier(t *testing.T) {
	resource, err := ArnResourceFromString("bucket/my-bucket/")

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.NotNil(t, resource.Qualifier())
	assert.Equal(t, "", *resource.Qualifier())
}

func TestArnResourceFromString_ComplexResourceName(t *testing.T) {
	resource, err := ArnResourceFromString("bucket:my-bucket-name-123")

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "my-bucket-name-123", resource.Resource())
	assert.Nil(t, resource.Qualifier())
}

func TestArnResourceFromString_ComplexQualifier(t *testing.T) {
	resource, err := ArnResourceFromString("bucket:my-bucket:obj/file.txt")

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.NotNil(t, resource.Qualifier())
	assert.Equal(t, "obj/file.txt", *resource.Qualifier())
}

func TestArnResourceFromString_TableArnResource(t *testing.T) {
	resource, err := ArnResourceFromString("bucket/test-bucket-9326/table/ad3fca49-9de8-4e5f-8d7c-e15c2588c2ad")

	assert.NotNil(t, resource)
	assert.Nil(t, err)
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "test-bucket-9326", resource.Resource())
	assert.NotNil(t, resource.Qualifier())
	assert.Equal(t, "table/ad3fca49-9de8-4e5f-8d7c-e15c2588c2ad", *resource.Qualifier())
}

func TestArnResource_String_WithAllFields(t *testing.T) {
	resourceType := "bucket"
	qualifier := "obj"
	resource, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Qualifier(&qualifier).
		Build()

	assert.Equal(t, "bucket:my-bucket:obj", resource.String())
}

func TestArnResource_String_WithNullFields(t *testing.T) {
	resource, _ := ArnResourceBuilderNew().
		Resource("my-bucket").
		Build()

	assert.Equal(t, ":my-bucket:", resource.String())
}

func TestArnResource_Equals_SameObject(t *testing.T) {
	resource, _ := ArnResourceBuilderNew().
		Resource("my-bucket").
		Build()

	assert.True(t, resource.Equals(resource))
}

func TestArnResource_Equals_EqualObjects(t *testing.T) {
	resource1, _ := ArnResourceBuilderNew().
		Resource("my-bucket").
		Build()

	resource2, _ := ArnResourceBuilderNew().
		Resource("my-bucket").
		Build()

	assert.True(t, resource1.Equals(resource2))
	assert.Equal(t, resource1.HashCode(), resource2.HashCode())
}

func TestArnResource_Equals_DifferentResourceType(t *testing.T) {
	bucket := "bucket"
	object := "object"
	resource1, _ := ArnResourceBuilderNew().
		ResourceType(&bucket).
		Resource("my-bucket").
		Build()

	resource2, _ := ArnResourceBuilderNew().
		ResourceType(&object).
		Resource("my-bucket").
		Build()

	assert.False(t, resource1.Equals(resource2))
}

func TestArnResource_Equals_DifferentResource(t *testing.T) {
	resourceType := "bucket"
	resource1, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Build()

	resource2, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("other-bucket").
		Build()

	assert.False(t, resource1.Equals(resource2))
}

func TestArnResource_Equals_DifferentQualifier(t *testing.T) {
	resourceType := "bucket"
	objQualifier := "obj"
	fileQualifier := "file"

	resource1, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Qualifier(&objQualifier).
		Build()

	resource2, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Qualifier(&fileQualifier).
		Build()

	assert.False(t, resource1.Equals(resource2))
}

func TestArnResource_Equals_NullQualifier(t *testing.T) {
	resourceType := "bucket"
	resource1, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Build()

	resource2, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Build()

	assert.True(t, resource1.Equals(resource2))
}

func TestArnResource_Equals_NilObject(t *testing.T) {
	resource, _ := ArnResourceBuilderNew().
		Resource("my-bucket").
		Build()

	assert.False(t, resource.Equals(nil))
}

func TestArnResource_HashCode_ConsistentWithEquals(t *testing.T) {
	resourceType := "bucket"
	qualifier := "obj"

	resource1, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Qualifier(&qualifier).
		Build()

	resource2, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Qualifier(&qualifier).
		Build()

	assert.True(t, resource1.Equals(resource2))
	assert.Equal(t, resource1.HashCode(), resource2.HashCode())
}

func TestArnResource_ToBuilder(t *testing.T) {
	resourceType := "bucket"
	qualifier := "obj"

	original, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Qualifier(&qualifier).
		Build()

	copy, _ := original.ToBuilder().Build()

	assert.True(t, original.Equals(copy))
}

func TestArnResource_ToBuilder_WithModification(t *testing.T) {
	resourceType := "bucket"
	newResourceType := "object"

	original, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Build()

	modified, _ := original.ToBuilder().
		ResourceType(&newResourceType).
		Build()

	assert.False(t, original.Equals(modified))
	assert.NotNil(t, modified.ResourceType())
	assert.Equal(t, "object", *modified.ResourceType())
}

func TestArnResource_EnsureBucketResource_Valid(t *testing.T) {
	resourceType := "bucket"
	resource, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Build()

	err := resource.EnsureBucketResource()
	assert.Nil(t, err)
}

func TestArnResource_EnsureBucketResource_NilResourceType(t *testing.T) {
	resource, _ := ArnResourceBuilderNew().
		Resource("my-bucket").
		Build()

	err := resource.EnsureBucketResource()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "is not bucket resource")
}

func TestArnResource_EnsureBucketResource_WrongResourceType(t *testing.T) {
	resourceType := "object"
	resource, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Build()

	err := resource.EnsureBucketResource()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "is not bucket resource")
}

func TestArnResource_EnsureBucketResource_WithQualifier(t *testing.T) {
	resourceType := "bucket"
	qualifier := "obj"
	resource, _ := ArnResourceBuilderNew().
		ResourceType(&resourceType).
		Resource("my-bucket").
		Qualifier(&qualifier).
		Build()

	err := resource.EnsureBucketResource()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "is not bucket resource")
}
