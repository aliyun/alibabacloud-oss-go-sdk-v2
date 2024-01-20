package oss

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOffsetAndSizeFromHeaders(t *testing.T) {
	// no header
	header := http.Header{}
	offset, size := parseOffsetAndSizeFromHeaders(header)
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(-1), size)

	// only Content-Length
	header = http.Header{}
	header.Set("Content-Length", "12335")
	offset, size = parseOffsetAndSizeFromHeaders(header)
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(12335), size)

	// Content-Length and Content-Range
	header = http.Header{}
	header.Set("Content-Length", "1000")
	header.Set("Content-Range", "bytes 1-499/1000")
	offset, size = parseOffsetAndSizeFromHeaders(header)
	assert.Equal(t, int64(1), offset)
	assert.Equal(t, int64(1000), size)

	// Content-Length and Content-Range
	header = http.Header{}
	header.Set("Content-Length", "1000")
	header.Set("Content-Range", "bytes 100-/1000")
	offset, size = parseOffsetAndSizeFromHeaders(header)
	assert.Equal(t, int64(100), offset)
	assert.Equal(t, int64(1000), size)

	// invalid Content-Length
	header = http.Header{}
	header.Set("Content-Length", "abcde")
	offset, size = parseOffsetAndSizeFromHeaders(header)
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(-1), size)

	// invalid Content-Range
	header = http.Header{}
	header.Set("Content-Length", "1000")
	header.Set("Content-Range", "byte 100-/1000")
	offset, size = parseOffsetAndSizeFromHeaders(header)
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(-1), size)

	// invalid Content-Range
	header = http.Header{}
	header.Set("Content-Length", "1000")
	header.Set("Content-Range", "bytes abc-/1000")
	offset, size = parseOffsetAndSizeFromHeaders(header)
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(-1), size)

	// invalid Content-Range
	header = http.Header{}
	header.Set("Content-Length", "1000")
	header.Set("Content-Range", "bytes 123-456")
	offset, size = parseOffsetAndSizeFromHeaders(header)
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(-1), size)

	// invalid Content-Range
	header = http.Header{}
	header.Set("Content-Length", "1000")
	header.Set("Content-Range", "bytes 123-456/abc")
	offset, size = parseOffsetAndSizeFromHeaders(header)
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(-1), size)
}

func TestParseContentRange(t *testing.T) {
	from, to, total, err := ParseContentRange("")
	assert.Equal(t, int64(0), from)
	assert.Equal(t, int64(0), to)
	assert.Equal(t, int64(0), total)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid content range", err.Error())

	from, to, total, err = ParseContentRange("invalid")
	assert.Equal(t, int64(0), from)
	assert.Equal(t, int64(0), to)
	assert.Equal(t, int64(0), total)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid content range", err.Error())

	from, to, total, err = ParseContentRange("otherUnit 22-33/42")
	assert.Equal(t, int64(0), from)
	assert.Equal(t, int64(0), to)
	assert.Equal(t, int64(0), total)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid content range", err.Error())

	from, to, total, err = ParseContentRange("bytes */42")
	assert.Equal(t, int64(0), from)
	assert.Equal(t, int64(0), to)
	assert.Equal(t, int64(0), total)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid content range", err.Error())

	from, to, total, err = ParseContentRange("bytes 22-33/42")
	assert.Equal(t, int64(22), from)
	assert.Equal(t, int64(33), to)
	assert.Equal(t, int64(42), total)
	assert.Nil(t, err)

	from, to, total, err = ParseContentRange("bytes 22-33/*")
	assert.Equal(t, int64(22), from)
	assert.Equal(t, int64(33), to)
	assert.Equal(t, int64(-1), total)
	assert.Nil(t, err)
}

type copyRequestStub struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the object.
	Key *string `input:"path,key,required"`

	Acl ObjectACLType `input:"header,x-oss-object-acl"`

	RequestCommon
}

func TestCopyRequest(t *testing.T) {
	requestStub := &copyRequestStub{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Acl:    "acl-abc",
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"header-1": "hvalue-1",
			},
			Parameters: map[string]string{
				"query-1": "qvalue-1",
			},
		},
	}
	request := &PutObjectRequest{}
	copyRequest(request, requestStub)

	assert.Equal(t, "bucket", ToString(request.Bucket))
	assert.Equal(t, "key", ToString(request.Key))
	assert.Equal(t, "acl-abc", ToString((*string)(&request.Acl)))
	assert.Equal(t, "hvalue-1", request.Headers["header-1"])
	assert.Equal(t, "qvalue-1", request.Parameters["query-1"])
}
