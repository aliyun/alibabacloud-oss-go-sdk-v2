package arn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArnBuilder_WithAllFields(t *testing.T) {
	region := "cn-hangzhou"
	accountId := "123456789012"
	arn, err := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	assert.NotNil(t, arn)
	assert.Nil(t, err)
	assert.Equal(t, "oss", arn.Service())
	assert.NotNil(t, arn.Region())
	assert.Equal(t, "cn-hangzhou", *arn.Region())
	assert.NotNil(t, arn.AccountId())
	assert.Equal(t, "123456789012", *arn.AccountId())
	assert.Equal(t, "bucket:my-bucket", arn.ResourceAsString())
	assert.NotNil(t, arn.Resource())
}

func TestArnBuilder_WithOnlyRequiredFields(t *testing.T) {
	arn, err := ArnBuilderNew().
		Service("oss").
		Resource("my-bucket").
		Build()

	assert.NotNil(t, arn)
	assert.Nil(t, err)
	assert.Equal(t, "oss", arn.Service())
	assert.Nil(t, arn.Region())
	assert.Nil(t, arn.AccountId())
	assert.Equal(t, "my-bucket", arn.ResourceAsString())
}

func TestArnBuilder_WithoutService_ThrowsError(t *testing.T) {
	arn, err := ArnBuilderNew().
		Resource("my-bucket").
		Build()

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "service must not be blank")
}

func TestArnBuilder_WithBlankService_ThrowsError(t *testing.T) {
	arn, err := ArnBuilderNew().
		Service("").
		Resource("my-bucket").
		Build()

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "service must not be blank")
}

func TestArnBuilder_WithoutResource_ThrowsError(t *testing.T) {
	arn, err := ArnBuilderNew().
		Service("oss").
		Build()

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "resource must not be blank")
}

func TestArnBuilder_WithBlankResource_ThrowsError(t *testing.T) {
	arn, err := ArnBuilderNew().
		Service("oss").
		Resource("").
		Build()

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "resource must not be blank")
}

func TestParseArn_ValidArn_WithAllComponents(t *testing.T) {
	arnString := "acs:oss:cn-hangzhou:123456789012:bucket:my-bucket"
	arn, err := ParseArn(arnString)

	assert.NotNil(t, arn)
	assert.Nil(t, err)
	assert.Equal(t, "oss", arn.Service())
	assert.NotNil(t, arn.Region())
	assert.Equal(t, "cn-hangzhou", *arn.Region())
	assert.NotNil(t, arn.AccountId())
	assert.Equal(t, "123456789012", *arn.AccountId())
	assert.Equal(t, "bucket:my-bucket", arn.ResourceAsString())
}

func TestParseArn_ValidArn_WithEmptyRegionAndAccount(t *testing.T) {
	arnString := "acs:oss:::bucket:my-bucket"
	arn, err := ParseArn(arnString)

	assert.NotNil(t, arn)
	assert.Nil(t, err)
	assert.Equal(t, "oss", arn.Service())
	assert.Nil(t, arn.Region())
	assert.Nil(t, arn.AccountId())
	assert.Equal(t, "bucket:my-bucket", arn.ResourceAsString())
}

func TestParseArn_ValidArn_WithEmptyRegion(t *testing.T) {
	arnString := "acs:oss::123456789012:bucket:my-bucket"
	arn, err := ParseArn(arnString)

	assert.NotNil(t, arn)
	assert.Nil(t, err)
	assert.Equal(t, "oss", arn.Service())
	assert.Nil(t, arn.Region())
	assert.NotNil(t, arn.AccountId())
	assert.Equal(t, "123456789012", *arn.AccountId())
}

func TestParseArn_ValidArn_WithEmptyAccount(t *testing.T) {
	arnString := "acs:oss:cn-hangzhou::bucket:my-bucket"
	arn, err := ParseArn(arnString)

	assert.NotNil(t, arn)
	assert.Nil(t, err)
	assert.Equal(t, "oss", arn.Service())
	assert.NotNil(t, arn.Region())
	assert.Equal(t, "cn-hangzhou", *arn.Region())
	assert.Nil(t, arn.AccountId())
}

func TestParseArn_ComplexResource(t *testing.T) {
	arnString := "acs:oss:cn-hangzhou:123456789012:bucket:my-bucket:obj/file.txt"
	arn, err := ParseArn(arnString)

	assert.NotNil(t, arn)
	assert.Nil(t, err)
	assert.Equal(t, "oss", arn.Service())
	resource := arn.Resource()
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.NotNil(t, resource.Qualifier())
	assert.Equal(t, "obj/file.txt", *resource.Qualifier())
}

func TestParseArn_TableArn(t *testing.T) {
	arnString := "acs:osstables:cn-beijing:123456:bucket/test-bucket-9326/table/ad3fca49-9de8-4e5f-8d7c-e15c2588c2ad"
	arn, err := ParseArn(arnString)

	assert.NotNil(t, arn)
	assert.Nil(t, err)
	assert.Equal(t, "osstables", arn.Service())
	assert.NotNil(t, arn.Region())
	assert.Equal(t, "cn-beijing", *arn.Region())
	assert.NotNil(t, arn.AccountId())
	assert.Equal(t, "123456", *arn.AccountId())

	resource := arn.Resource()
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "test-bucket-9326", resource.Resource())
	assert.NotNil(t, resource.Qualifier())
	assert.Equal(t, "table/ad3fca49-9de8-4e5f-8d7c-e15c2588c2ad", *resource.Qualifier())
}

func TestParseArn_InvalidArn_NoAcsPrefix(t *testing.T) {
	arn, err := ParseArn("ats:oss:cn-hangzhou:123456789012:bucket:my-bucket")

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed ARN - doesn't start with 'acs:'")
}

func TestParseArn_InvalidArn_NoService(t *testing.T) {
	arn, err := ParseArn("acs::cn-hangzhou:123456789012:bucket:my-bucket")

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "service must not be blank")
}

func TestParseArn_InvalidArn_NoRegion(t *testing.T) {
	arn, err := ParseArn("acs:oss")

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed ARN - no service specified")
}

func TestParseArn_InvalidArn_NoAccount(t *testing.T) {
	arn, err := ParseArn("acs:oss:cn-hangzhou")

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed ARN - no region specified")
}

func TestParseArn_InvalidArn_NoResource(t *testing.T) {
	arn, err := ParseArn("acs:oss:cn-hangzhou:123456789012:")

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed ARN - no resource specified")
}

func TestParseArn_EmptyString(t *testing.T) {
	arn, err := ParseArn("")

	assert.Nil(t, arn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ARN parsing failed")
}

func TestTryParseArn_ValidArn(t *testing.T) {
	arnString := "acs:oss:cn-hangzhou:123456789012:bucket:my-bucket"
	arn, err := TryParseArn(arnString)

	assert.NotNil(t, arn)
	assert.Nil(t, err)
	assert.Equal(t, "oss", arn.Service())
}

func TestTryParseArn_InvalidArn_NoAcsPrefix(t *testing.T) {
	arn, err := TryParseArn("ats:oss:cn-hangzhou:123456789012:bucket:my-bucket")

	assert.Nil(t, arn)
	assert.Nil(t, err)
}

func TestTryParseArn_InvalidArn_NoRegion(t *testing.T) {
	arn, err := TryParseArn("acs:oss")

	assert.Nil(t, arn)
	assert.Nil(t, err)
}

func TestTryParseArn_InvalidArn_NoAccount(t *testing.T) {
	arn, err := TryParseArn("acs:oss:cn-hangzhou")

	assert.Nil(t, arn)
	assert.Nil(t, err)
}

func TestTryParseArn_InvalidArn_NoResource(t *testing.T) {
	arn, err := TryParseArn("acs:oss:cn-hangzhou:123456789012:")

	assert.Nil(t, arn)
	assert.Nil(t, err)
}

func TestArn_String_WithAllFields(t *testing.T) {
	region := "cn-hangzhou"
	accountId := "123456789012"
	arn, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	assert.Equal(t, "acs:oss:cn-hangzhou:123456789012:bucket:my-bucket", arn.String())
}

func TestArn_String_WithEmptyRegionAndAccount(t *testing.T) {
	arn, _ := ArnBuilderNew().
		Service("oss").
		Resource("my-bucket").
		Build()

	assert.Equal(t, "acs:oss:::my-bucket", arn.String())
}

func TestArn_String_WithEmptyRegion(t *testing.T) {
	accountId := "123456789012"
	arn, _ := ArnBuilderNew().
		Service("oss").
		AccountId(&accountId).
		Resource("my-bucket").
		Build()

	assert.Equal(t, "acs:oss::123456789012:my-bucket", arn.String())
}

func TestArn_String_WithEmptyAccount(t *testing.T) {
	region := "cn-hangzhou"
	arn, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		Resource("my-bucket").
		Build()

	assert.Equal(t, "acs:oss:cn-hangzhou::my-bucket", arn.String())
}

func TestArn_Equals_SameObject(t *testing.T) {
	arn, _ := ArnBuilderNew().
		Service("oss").
		Resource("my-bucket").
		Build()

	assert.True(t, arn.Equals(arn))
}

func TestArn_Equals_EqualObjects(t *testing.T) {
	region := "cn-hangzhou"
	accountId := "123456789012"

	arn1, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	arn2, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	assert.True(t, arn1.Equals(arn2))
	assert.Equal(t, arn1.HashCode(), arn2.HashCode())
}

func TestArn_Equals_DifferentService(t *testing.T) {
	region := "cn-hangzhou"
	accountId := "123456789012"

	arn1, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	arn2, _ := ArnBuilderNew().
		Service("ots").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	assert.False(t, arn1.Equals(arn2))
}

func TestArn_Equals_DifferentRegion(t *testing.T) {
	hangzhou := "cn-hangzhou"
	shanghai := "cn-shanghai"
	accountId := "123456789012"

	arn1, _ := ArnBuilderNew().
		Service("oss").
		Region(&hangzhou).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	arn2, _ := ArnBuilderNew().
		Service("oss").
		Region(&shanghai).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	assert.False(t, arn1.Equals(arn2))
}

func TestArn_Equals_DifferentAccountId(t *testing.T) {
	region := "cn-hangzhou"
	account1 := "123456789012"
	account2 := "987654321098"

	arn1, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&account1).
		Resource("bucket:my-bucket").
		Build()

	arn2, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&account2).
		Resource("bucket:my-bucket").
		Build()

	assert.False(t, arn1.Equals(arn2))
}

func TestArn_Equals_DifferentResource(t *testing.T) {
	region := "cn-hangzhou"
	accountId := "123456789012"

	arn1, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	arn2, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:other-bucket").
		Build()

	assert.False(t, arn1.Equals(arn2))
}

func TestArn_Equals_NilObject(t *testing.T) {
	arn, _ := ArnBuilderNew().
		Service("oss").
		Resource("my-bucket").
		Build()

	assert.False(t, arn.Equals(nil))
}

func TestArn_HashCode_ConsistentWithEquals(t *testing.T) {
	region := "cn-hangzhou"
	accountId := "123456789012"

	arn1, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	arn2, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	assert.True(t, arn1.Equals(arn2))
	assert.Equal(t, arn1.HashCode(), arn2.HashCode())
}

func TestArn_ToBuilder(t *testing.T) {
	region := "cn-hangzhou"
	accountId := "123456789012"

	original, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	copy, _ := original.ToBuilder().Build()

	assert.True(t, original.Equals(copy))
}

func TestArn_ToBuilder_WithModification(t *testing.T) {
	region := "cn-hangzhou"
	newRegion := "cn-shanghai"
	accountId := "123456789012"

	original, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket").
		Build()

	modified, _ := original.ToBuilder().
		Region(&newRegion).
		Build()

	assert.False(t, original.Equals(modified))
	assert.NotNil(t, modified.Region())
	assert.Equal(t, "cn-shanghai", *modified.Region())
}

func TestArn_Resource_AsArnResource(t *testing.T) {
	region := "cn-hangzhou"
	accountId := "123456789012"

	arn, _ := ArnBuilderNew().
		Service("oss").
		Region(&region).
		AccountId(&accountId).
		Resource("bucket:my-bucket:obj/file.txt").
		Build()

	resource := arn.Resource()
	assert.NotNil(t, resource.ResourceType())
	assert.Equal(t, "bucket", *resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.NotNil(t, resource.Qualifier())
	assert.Equal(t, "obj/file.txt", *resource.Qualifier())
}

func TestArn_Resource_SimpleResource(t *testing.T) {
	arn, _ := ArnBuilderNew().
		Service("oss").
		Resource("my-bucket").
		Build()

	resource := arn.Resource()
	assert.Nil(t, resource.ResourceType())
	assert.Equal(t, "my-bucket", resource.Resource())
	assert.Nil(t, resource.Qualifier())
}
