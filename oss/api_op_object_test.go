package oss

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutObjectRequest
	var input *OperationInput
	var err error

	request = &PutObjectRequest{}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutObjectRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.Body)
	assert.Equal(t, input.Headers["x-oss-object-acl"], "")

	request = &PutObjectRequest{
		Bucket:               Ptr("oss-bucket"),
		Key:                  Ptr("oss-key"),
		CacheControl:         Ptr("no-cache"),
		ContentDisposition:   Ptr("attachment"),
		ContentEncoding:      Ptr("utf-8"),
		ContentMD5:           Ptr("eB5eJF1ptWaXm4bijSPyxw=="),
		ContentLength:        Ptr(int64(100)),
		Expires:              Ptr("2022-10-12T00:00:00.000Z"),
		ForbidOverwrite:      Ptr("true"),
		ServerSideEncryption: Ptr("AES256"),
		Acl:                  ObjectACLPrivate,
		StorageClass:         StorageClassStandard,
		Metadata: map[string]string{
			"location": "demo",
			"user":     "walker",
		},
		Tagging: Ptr("TagA=A&TagB=B"),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.Body)
	assert.Equal(t, input.Headers["Cache-Control"], "no-cache")
	assert.Equal(t, input.Headers["Content-Disposition"], "attachment")
	assert.Equal(t, input.Headers["x-oss-meta-user"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-location"], "demo")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "AES256")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "true")
	assert.Equal(t, input.Headers["Content-Encoding"], "utf-8")
	assert.Equal(t, input.Headers["Content-Length"], "100")
	assert.Equal(t, input.Headers["Content-MD5"], "eB5eJF1ptWaXm4bijSPyxw==")
	assert.Equal(t, input.Headers["Expires"], "2022-10-12T00:00:00.000Z")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=A&TagB=B")
	assert.Nil(t, input.Parameters)
	assert.Nil(t, input.OpMetadata.values)

	body := randLowStr(1000)
	request = &PutObjectRequest{
		Bucket:                    Ptr("oss-bucket"),
		Key:                       Ptr("oss-key"),
		CacheControl:              Ptr("no-cache"),
		ContentDisposition:        Ptr("attachment"),
		ContentEncoding:           Ptr("utf-8"),
		ContentMD5:                Ptr("eB5eJF1ptWaXm4bijSPyxw=="),
		ContentLength:             Ptr(int64(100)),
		Expires:                   Ptr("2022-10-12T00:00:00.000Z"),
		ForbidOverwrite:           Ptr("false"),
		ServerSideEncryption:      Ptr("KMS"),
		ServerSideDataEncryption:  Ptr("SM4"),
		ServerSideEncryptionKeyId: Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
		Acl:                       ObjectACLPrivate,
		StorageClass:              StorageClassStandard,
		Metadata: map[string]string{
			"name":  "walker",
			"email": "demo@aliyun.com",
		},
		Tagging: Ptr("TagA=B&TagC=D"),
		Body:    strings.NewReader(body),
	}

	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Body, strings.NewReader(body))
	assert.Equal(t, input.Headers["Cache-Control"], "no-cache")
	assert.Equal(t, input.Headers["Content-Disposition"], "attachment")
	assert.Equal(t, input.Headers["x-oss-meta-name"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-email"], "demo@aliyun.com")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "KMS")
	assert.Equal(t, input.Headers["x-oss-server-side-data-encryption"], "SM4")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption-key-id"], "9468da86-3509-4f8d-a61e-6eab1eac****")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "false")
	assert.Equal(t, input.Headers["Content-Encoding"], "utf-8")
	assert.Equal(t, input.Headers["Content-Length"], "100")
	assert.Equal(t, input.Headers["Content-MD5"], "eB5eJF1ptWaXm4bijSPyxw==")
	assert.Equal(t, input.Headers["Expires"], "2022-10-12T00:00:00.000Z")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=B&TagC=D")
	assert.Nil(t, input.Parameters)
	assert.Nil(t, input.OpMetadata.values)

	callbackMap := map[string]string{}
	callbackMap["callbackUrl"] = "www.aliyuncs.com"
	callbackMap["callbackBody"] = "filename=${object}&size=${size}&mimeType=${mimeType}"
	callbackMap["callbackBodyType"] = "application/x-www-form-urlencoded"
	callbackBuffer := bytes.NewBuffer([]byte{})
	callbackEncoder := json.NewEncoder(callbackBuffer)
	callbackEncoder.SetEscapeHTML(false)
	err = callbackEncoder.Encode(callbackMap)
	assert.Nil(t, err)

	callbackVal := base64.StdEncoding.EncodeToString(callbackBuffer.Bytes())
	callbackVar := base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`))
	request = &PutObjectRequest{
		Bucket:      Ptr("oss-bucket"),
		Key:         Ptr("oss-key"),
		Body:        strings.NewReader(body),
		Callback:    Ptr(callbackVal),
		CallbackVar: Ptr(callbackVar),
	}

	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Body, strings.NewReader(body))
	assert.Equal(t, input.Headers["x-oss-callback"], callbackVal)
	assert.Equal(t, input.Headers["x-oss-callback-var"], callbackVar)
	assert.Nil(t, input.Parameters)
	assert.Nil(t, input.OpMetadata.values)

	request = &PutObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		Body:         strings.NewReader(body),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Body, strings.NewReader(body))
	assert.Equal(t, input.Headers["x-oss-traffic-limit"], strconv.FormatInt(100*1024*8, 10))
	assert.Nil(t, input.Parameters)
	assert.Nil(t, input.OpMetadata.values)

	request = &PutObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		Body:         strings.NewReader(body),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Body, strings.NewReader(body))
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
	assert.Nil(t, input.Parameters)
	assert.Nil(t, input.OpMetadata.values)
}

func TestUnmarshalOutput_PutObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                 {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
		},
	}
	result := &PutObjectResult{}
	var unmarshalFns []func(result any, output *OperationOutput) error
	unmarshalFns = append(unmarshalFns, unmarshalHeader)
	unmarshalFns = append(unmarshalFns, discardBody)
	err = c.unmarshalOutput(result, output, unmarshalFns...)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, *result.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Nil(t, result.VersionId)

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                 {"\"A797938C31D59EDD08D86188F6D5****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
			"x-oss-version-id":     {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
		},
	}
	result = &PutObjectResult{}
	err = c.unmarshalOutput(result, output, unmarshalFns...)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"A797938C31D59EDD08D86188F6D5****\"")
	assert.Equal(t, *result.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, *result.VersionId, "CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                 {"\"A797938C31D59EDD08D86188F6D5****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
			"x-oss-version-id":     {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
		},
		Body: io.NopCloser(strings.NewReader(`{"filename":"object.txt","size":"100","mimeType":""}`)),
	}
	result = &PutObjectResult{}
	unmarshalFns = []func(result any, output *OperationOutput) error{}
	unmarshalFns = append(unmarshalFns, unmarshalHeader)
	unmarshalFns = append(unmarshalFns, unmarshalCallbackBody)
	err = c.unmarshalOutput(result, output, unmarshalFns...)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"A797938C31D59EDD08D86188F6D5****\"")
	assert.Equal(t, *result.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, *result.VersionId, "CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****")
	jsonData, err := json.Marshal(result.CallbackResult)
	assert.Nil(t, err)
	assert.NotEmpty(t, string(jsonData))

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalFns...)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalFns...)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 203,
		Status:     "Non-Authoritative Information",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalFns...)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 203)
	assert.Equal(t, result.Status, "Non-Authoritative Information")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

}

func TestMarshalInput_GetObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetObjectRequest
	var input *OperationInput
	var err error

	request = &GetObjectRequest{}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetObjectRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")

	request = &GetObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Headers["x-oss-traffic-limit"], strconv.FormatInt(100*1024*8, 10))

	request = &GetObjectRequest{
		Bucket:                     Ptr("oss-bucket"),
		Key:                        Ptr("oss-key"),
		IfMatch:                    Ptr("\"D41D8CD98F00B204E9800998ECF8****\""),
		IfNoneMatch:                Ptr("\"D41D8CD98F00B204E9800998ECF9****\""),
		IfModifiedSince:            Ptr("Fri, 13 Nov 2023 14:47:53 GMT"),
		IfUnmodifiedSince:          Ptr("Fri, 13 Nov 2015 14:47:53 GMT"),
		Range:                      Ptr("bytes 0~9/44"),
		ResponseCacheControl:       Ptr("gzip"),
		ResponseContentDisposition: Ptr("attachment; filename=testing.txt"),
		ResponseContentEncoding:    Ptr("utf-8"),
		ResponseContentLanguage:    Ptr("中文"),
		ResponseContentType:        Ptr("text"),
		ResponseExpires:            Ptr("Fri, 24 Feb 2012 17:00:00 GMT"),
		VersionId:                  Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")

	assert.Equal(t, input.Headers["If-Match"], "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, input.Headers["If-None-Match"], "\"D41D8CD98F00B204E9800998ECF9****\"")
	assert.Equal(t, input.Headers["If-Modified-Since"], "Fri, 13 Nov 2023 14:47:53 GMT")
	assert.Equal(t, input.Headers["If-Unmodified-Since"], "Fri, 13 Nov 2015 14:47:53 GMT")
	assert.Equal(t, input.Headers["Range"], "bytes 0~9/44")
	assert.Equal(t, input.Parameters["response-cache-control"], "gzip")
	assert.Equal(t, input.Parameters["response-content-disposition"], "attachment; filename=testing.txt")
	assert.Equal(t, input.Parameters["response-content-encoding"], "utf-8")
	assert.Equal(t, input.Parameters["response-content-language"], "中文")
	assert.Equal(t, input.Parameters["response-expires"], "Fri, 24 Feb 2012 17:00:00 GMT")
	assert.Equal(t, input.Parameters["versionId"], "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****")
	assert.Nil(t, input.OpMetadata.values)

	request = &GetObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_GetObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := randLowStr(344606)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":  {"3a8f-2e2d-7965-3ff9-51c875b*****"},
			"Content-Type":      {"image/jpg"},
			"Date":              {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":              {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"Content-Length":    {"344606"},
			"Last-Modified":     {"Fri, 24 Feb 2012 06:07:48 GMT"},
			"x-oss-object-type": {"Normal"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &GetObjectResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "3a8f-2e2d-7965-3ff9-51c875b*****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, *result.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
	assert.Equal(t, *result.ContentType, "image/jpg")
	assert.Equal(t, result.ContentLength, int64(344606))
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, result.Body, io.NopCloser(bytes.NewReader([]byte(body))))

	body = randLowStr(34460)
	output = &OperationOutput{
		StatusCode: 206,
		Status:     "Partial Content",
		Headers: http.Header{
			"X-Oss-Request-Id":  {"28f6-15ea-8224-234e-c0ce407****"},
			"Content-Type":      {"image/jpg"},
			"Date":              {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":              {"\"5B3C1A2E05E1B002CC607C****\""},
			"Content-Length":    {"801"},
			"Last-Modified":     {"Fri, 24 Feb 2012 06:07:48 GMT"},
			"x-oss-object-type": {"Normal"},
			"Accept-Ranges":     {"bytes"},
			"Content-Range":     {"bytes 100-900/34460"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &GetObjectResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 206)
	assert.Equal(t, result.Status, "Partial Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
	assert.Equal(t, *result.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
	assert.Equal(t, *result.ContentType, "image/jpg")
	assert.Equal(t, result.ContentLength, int64(801))
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, result.Body, io.NopCloser(bytes.NewReader([]byte(body))))
	assert.Equal(t, *result.ContentRange, "bytes 100-900/34460")

	body = randLowStr(344606)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":                    {"28f6-15ea-8224-234e-c0ce407****"},
			"Content-Type":                        {"text"},
			"Date":                                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                                {"\"5B3C1A2E05E1B002CC607C****\""},
			"Content-Length":                      {"344606"},
			"Last-Modified":                       {"Fri, 24 Feb 2012 06:07:48 GMT"},
			"x-oss-object-type":                   {"Normal"},
			"Accept-Ranges":                       {"bytes"},
			"Content-disposition":                 {"attachment; filename=testing.txt"},
			"Cache-control":                       {"no-cache"},
			"X-Oss-Storage-Class":                 {"Standard"},
			"x-oss-server-side-encryption":        {"KMS"},
			"x-oss-server-side-data-encryption":   {"SM4"},
			"x-oss-server-side-encryption-key-id": {"12f8711f-90df-4e0d-903d-ab972b0f****"},
			"x-oss-tagging-count":                 {"2"},
			"Content-MD5":                         {"si4Nw3Cn9wZ/rPX3XX+j****"},
			"x-oss-hash-crc64ecma":                {"870718044876840****"},
			"x-oss-meta-name":                     {"demo"},
			"x-oss-meta-email":                    {"demo@aliyun.com"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &GetObjectResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, *result.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
	assert.Equal(t, *result.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
	assert.Equal(t, *result.ContentType, "text")
	assert.Equal(t, result.ContentLength, int64(344606))
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, *result.StorageClass, "Standard")
	assert.Equal(t, result.Body, io.NopCloser(bytes.NewReader([]byte(body))))
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
	assert.Equal(t, result.TaggingCount, int32(2))
	assert.Equal(t, result.Metadata["name"], "demo")
	assert.Equal(t, result.Metadata["email"], "demo@aliyun.com")
	assert.Equal(t, *result.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
	assert.Equal(t, *result.HashCRC64, "870718044876840****")
	body = randLowStr(344606)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":  {"28f6-15ea-8224-234e-c0ce407****"},
			"Content-Type":      {"text"},
			"Date":              {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":              {"\"5B3C1A2E05E1B002CC607C****\""},
			"Content-Length":    {"344606"},
			"Last-Modified":     {"Fri, 24 Feb 2012 06:07:48 GMT"},
			"x-oss-object-type": {"Normal"},
			"x-oss-restore":     {"ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\""},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &GetObjectResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.Headers.Get("Cache-control"), "")
	assert.Equal(t, *result.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
	assert.Equal(t, *result.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
	assert.Equal(t, *result.ContentType, "text")
	assert.Equal(t, result.ContentLength, int64(344606))
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, result.Body, io.NopCloser(bytes.NewReader([]byte(body))))
	assert.Equal(t, *result.Restore, "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_CopyObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CopyObjectRequest
	var input *OperationInput
	var err error

	request = &CopyObjectRequest{}
	source := encodeSourceObject(request)
	input = &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CopyObjectRequest{
		Bucket: Ptr("oss-bucket"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CopyObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CopyObjectRequest{
		Bucket:    Ptr("oss-bucket"),
		Key:       Ptr("oss-copy-key"),
		SourceKey: Ptr("oss-src-key"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-copy-key")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-bucket/oss-src-key")

	request = &CopyObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-copy-key"),
		SourceKey:    Ptr("oss-key"),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-copy-key")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-bucket/oss-key")
	assert.Equal(t, input.Headers["x-oss-traffic-limit"], strconv.FormatInt(100*1024*8, 10))

	request = &CopyObjectRequest{
		Bucket:            Ptr("oss-bucket"),
		Key:               Ptr("oss-copy-key"),
		SourceKey:         Ptr("oss-dir/oss-obj"),
		IfMatch:           Ptr("\"D41D8CD98F00B204E9800998ECF8****\""),
		IfNoneMatch:       Ptr("\"D41D8CD98F00B204E9800998ECF9****\""),
		IfModifiedSince:   Ptr("Fri, 13 Nov 2023 14:47:53 GMT"),
		IfUnmodifiedSince: Ptr("Fri, 13 Nov 2015 14:47:53 GMT"),
		SourceVersionId:   Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-copy-key")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-match"], "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-none-match"], "\"D41D8CD98F00B204E9800998ECF9****\"")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-modified-since"], "Fri, 13 Nov 2023 14:47:53 GMT")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-unmodified-since"], "Fri, 13 Nov 2015 14:47:53 GMT")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-bucket/oss-dir/oss-obj"+"?versionId=CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****")
	assert.Nil(t, input.OpMetadata.values)

	request = &CopyObjectRequest{
		Bucket:                    Ptr("oss-copy-bucket"),
		Key:                       Ptr("oss-copy-key"),
		SourceKey:                 Ptr("oss-key"),
		IfMatch:                   Ptr("\"D41D8CD98F00B204E9800998ECF8****\""),
		IfNoneMatch:               Ptr("\"D41D8CD98F00B204E9800998ECF9****\""),
		IfModifiedSince:           Ptr("Fri, 13 Nov 2023 14:47:53 GMT"),
		IfUnmodifiedSince:         Ptr("Fri, 13 Nov 2015 14:47:53 GMT"),
		SourceVersionId:           Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****"),
		ForbidOverwrite:           Ptr("false"),
		ServerSideEncryption:      Ptr("KMS"),
		ServerSideDataEncryption:  Ptr("SM4"),
		ServerSideEncryptionKeyId: Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
		MetadataDirective:         Ptr("REPLACE"),
		TaggingDirective:          Ptr("Replace"),
		Acl:                       ObjectACLPrivate,
		StorageClass:              StorageClassStandard,
		Metadata: map[string]string{
			"name":  "walker",
			"email": "demo@aliyun.com",
		},
		Tagging: Ptr("TagA=B&TagC=D"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-copy-bucket")
	assert.Equal(t, *input.Key, "oss-copy-key")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-match"], "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-none-match"], "\"D41D8CD98F00B204E9800998ECF9****\"")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-modified-since"], "Fri, 13 Nov 2023 14:47:53 GMT")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-unmodified-since"], "Fri, 13 Nov 2015 14:47:53 GMT")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-copy-bucket/oss-key?versionId=CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****")
	assert.Equal(t, input.Headers["x-oss-meta-name"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-email"], "demo@aliyun.com")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "KMS")
	assert.Equal(t, input.Headers["x-oss-server-side-data-encryption"], "SM4")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption-key-id"], "9468da86-3509-4f8d-a61e-6eab1eac****")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "false")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=B&TagC=D")
	assert.Equal(t, input.Headers["x-oss-tagging-directive"], "Replace")
	assert.Equal(t, input.Headers["x-oss-metadata-directive"], "REPLACE")
	assert.Nil(t, input.OpMetadata.values)

	request = &CopyObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-copy-key"),
		SourceKey:    Ptr("oss-key"),
		RequestPayer: Ptr("requester"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-copy-key")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-bucket/oss-key")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_CopyObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<CopyObjectResult>
  <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
  <LastModified>2018-02-24T09:41:56.000Z</LastModified>
</CopyObjectResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"3a8f-2e2d-7965-3ff9-51c875b*****"},
			"Content-Type":         {"image/jpg"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                 {"\"F2064A169EE92E9775EE5324D0B1****\""},
			"Content-Length":       {"344606"},
			"x-oss-hash-crc64ecma": {"1275300285919610****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &CopyObjectResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "3a8f-2e2d-7965-3ff9-51c875b*****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
	assert.Equal(t, *result.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, *result.HashCRC64, "1275300285919610****")

	body = `<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":                    {"28f6-15ea-8224-234e-c0ce407****"},
			"Content-Type":                        {"text"},
			"Date":                                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                                {"\"F2064A169EE92E9775EE5324D0B1****\""},
			"Content-Length":                      {"344606"},
			"x-oss-server-side-encryption":        {"KMS"},
			"x-oss-server-side-data-encryption":   {"SM4"},
			"x-oss-server-side-encryption-key-id": {"12f8711f-90df-4e0d-903d-ab972b0f****"},
			"x-oss-hash-crc64ecma":                {"870718044876840****"},
			"x-oss-copy-source-version-id":        {"CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****"},
			"x-oss-version-id":                    {"CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &CopyObjectResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
	assert.Equal(t, *result.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
	assert.Equal(t, *result.HashCRC64, "870718044876840****")
	assert.Equal(t, *result.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
	assert.Equal(t, *result.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
	<Error>
	 <Code>AccessDenied</Code>
	 <Message>AccessDenied</Message>
	 <RequestId>568D5566F2D0F89F5C0E****</RequestId>
	 <HostId>test.oss.aliyuncs.com</HostId>
	</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_AppendObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *AppendObjectRequest
	var input *OperationInput
	var err error

	request = &AppendObjectRequest{}
	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AppendObjectRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AppendObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	p := int64(0)
	request = &AppendObjectRequest{
		Bucket:               Ptr("oss-bucket"),
		Key:                  Ptr("oss-key"),
		Position:             Ptr(p),
		CacheControl:         Ptr("no-cache"),
		ContentDisposition:   Ptr("attachment"),
		ContentEncoding:      Ptr("gzip"),
		ContentMD5:           Ptr("eB5eJF1ptWaXm4bijSPyxw=="),
		ContentLength:        Ptr(int64(100)),
		Expires:              Ptr("2022-10-12T00:00:00.000Z"),
		ForbidOverwrite:      Ptr("true"),
		ServerSideEncryption: Ptr("AES256"),
		Acl:                  ObjectACLPrivate,
		StorageClass:         StorageClassStandard,
		Metadata: map[string]string{
			"location": "demo",
			"user":     "walker",
		},
		Tagging: Ptr("TagA=A&TagB=B"),
	}
	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.Body)
	assert.Equal(t, input.Headers["Cache-Control"], "no-cache")
	assert.Equal(t, input.Headers["Content-Disposition"], "attachment")
	assert.Equal(t, input.Headers["x-oss-meta-user"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-location"], "demo")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "AES256")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "true")
	assert.Equal(t, input.Headers["Content-Encoding"], "gzip")
	assert.Equal(t, input.Headers["Content-Length"], "100")
	assert.Equal(t, input.Headers["Content-MD5"], "eB5eJF1ptWaXm4bijSPyxw==")
	assert.Equal(t, input.Headers["Expires"], "2022-10-12T00:00:00.000Z")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=A&TagB=B")
	assert.Empty(t, input.Parameters["append"])
	assert.Equal(t, input.Parameters["position"], strconv.FormatInt(p, 10))
	assert.Nil(t, input.OpMetadata.values)

	body := randLowStr(1000)
	request = &AppendObjectRequest{
		Bucket:                    Ptr("oss-bucket"),
		Key:                       Ptr("oss-key"),
		Position:                  Ptr(int64(0)),
		CacheControl:              Ptr("no-cache"),
		ContentDisposition:        Ptr("attachment"),
		ContentEncoding:           Ptr("utf-8"),
		ContentMD5:                Ptr("eB5eJF1ptWaXm4bijSPyxw=="),
		ContentLength:             Ptr(int64(100)),
		Expires:                   Ptr("2022-10-12T00:00:00.000Z"),
		ForbidOverwrite:           Ptr("false"),
		ServerSideEncryption:      Ptr("KMS"),
		ServerSideDataEncryption:  Ptr("SM4"),
		ServerSideEncryptionKeyId: Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
		Acl:                       ObjectACLPrivate,
		StorageClass:              StorageClassStandard,
		Metadata: map[string]string{
			"name":  "walker",
			"email": "demo@aliyun.com",
		},
		Tagging: Ptr("TagA=B&TagC=D"),
		Body:    strings.NewReader(body),
	}

	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Body, strings.NewReader(body))
	assert.Equal(t, input.Headers["Cache-Control"], "no-cache")
	assert.Equal(t, input.Headers["Content-Disposition"], "attachment")
	assert.Equal(t, input.Headers["x-oss-meta-name"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-email"], "demo@aliyun.com")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "KMS")
	assert.Equal(t, input.Headers["x-oss-server-side-data-encryption"], "SM4")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption-key-id"], "9468da86-3509-4f8d-a61e-6eab1eac****")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "false")
	assert.Equal(t, input.Headers["Content-Encoding"], "utf-8")
	assert.Equal(t, input.Headers["Content-Length"], "100")
	assert.Equal(t, input.Headers["Content-MD5"], "eB5eJF1ptWaXm4bijSPyxw==")
	assert.Equal(t, input.Headers["Expires"], "2022-10-12T00:00:00.000Z")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=B&TagC=D")
	assert.Empty(t, input.Parameters["append"])
	assert.Equal(t, input.Parameters["position"], strconv.FormatInt(p, 10))
	assert.Nil(t, input.OpMetadata.values)

	request = &AppendObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		Position:     Ptr(int64(0)),
		Body:         strings.NewReader(body),
		TrafficLimit: int64(100 * 1024 * 8),
	}

	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Body, strings.NewReader(body))
	assert.Empty(t, input.Parameters["append"])
	assert.Equal(t, input.Parameters["position"], strconv.FormatInt(p, 10))
	assert.Nil(t, input.OpMetadata.values)
	assert.Equal(t, input.Headers["x-oss-traffic-limit"], strconv.FormatInt(100*1024*8, 10))

	request = &AppendObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		Position:     Ptr(int64(0)),
		Body:         strings.NewReader(body),
		RequestPayer: Ptr("requester"),
	}

	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Body, strings.NewReader(body))
	assert.Empty(t, input.Parameters["append"])
	assert.Equal(t, input.Parameters["position"], strconv.FormatInt(p, 10))
	assert.Nil(t, input.OpMetadata.values)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_AppendObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":           {"5C06A3B67B8B5A3DA422****"},
			"Date":                       {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                       {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"x-oss-hash-crc64ecma":       {"316181249502703****"},
			"x-oss-next-append-position": {"0"},
		},
	}
	result := &AppendObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, result.NextPosition, int64(0))
	assert.Nil(t, result.VersionId)

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":           {"5C06A3B67B8B5A3DA422****"},
			"Date":                       {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                       {"\"A797938C31D59EDD08D86188F6D5****\""},
			"x-oss-hash-crc64ecma":       {"316181249502703****"},
			"x-oss-version-id":           {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
			"x-oss-next-append-position": {"1717"},
		},
	}
	result = &AppendObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, *result.VersionId, "CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****")
	assert.Equal(t, result.NextPosition, int64(1717))

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":           {"5C06A3B67B8B5A3DA422****"},
			"Date":                       {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                       {"\"A797938C31D59EDD08D86188F6D5****\""},
			"x-oss-hash-crc64ecma":       {"316181249502703****"},
			"x-oss-version-id":           {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
			"x-oss-next-append-position": {"1717"},
		},
	}
	result = &AppendObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, *result.VersionId, "CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****")
	assert.Equal(t, result.NextPosition, int64(1717))
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":                    {"5C06A3B67B8B5A3DA422****"},
			"Date":                                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                                {"\"A797938C31D59EDD08D86188F6D5****\""},
			"x-oss-hash-crc64ecma":                {"316181249502703****"},
			"x-oss-version-id":                    {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
			"x-oss-next-append-position":          {"1717"},
			"x-oss-server-side-encryption":        {"KMS"},
			"x-oss-server-side-data-encryption":   {"SM4"},
			"x-oss-server-side-encryption-key-id": {"12f8711f-90df-4e0d-903d-ab972b0f****"},
		},
	}
	result = &AppendObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, *result.VersionId, "CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****")
	assert.Equal(t, result.NextPosition, int64(1717))
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 409,
		Status:     "ObjectNotAppendable",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "ObjectNotAppendable")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteObjectRequest
	var input *OperationInput
	var err error

	request = &DeleteObjectRequest{}
	input = &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteObjectRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Nil(t, input.Parameters)
	request = &DeleteObjectRequest{
		Bucket:    Ptr("oss-bucket"),
		Key:       Ptr("oss-key"),
		VersionId: Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY****"),
	}
	input = &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Parameters["versionId"], "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY****")
	assert.Nil(t, input.OpMetadata.values)

	request = &DeleteObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
	assert.Nil(t, input.OpMetadata.values)
}

func TestUnmarshalOutput_DeleteObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"3a8f-2e2d-7965-3ff9-51c875b*****"},
			"Content-Type":     {"image/jpg"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
		},
	}
	result := &DeleteObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "3a8f-2e2d-7965-3ff9-51c875b*****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id":    {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"x-oss-version-id":    {"CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****"},
			"x-oss-delete-marker": {"true"},
		},
	}
	result = &DeleteObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
	assert.True(t, result.DeleteMarker)

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteMultipleObjects(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteMultipleObjectsRequest
	var input *OperationInput
	var err error

	request = &DeleteMultipleObjectsRequest{}
	input = &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"delete": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteMultipleObjectsRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"delete": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteMultipleObjectsRequest{
		Bucket:  Ptr("oss-bucket"),
		Objects: []DeleteObject{},
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteMultipleObjectsRequest{
		Bucket:  Ptr("oss-bucket"),
		Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
	}
	input = &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"delete": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, input.Body, strings.NewReader("<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>"))
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["delete"])
	request = &DeleteMultipleObjectsRequest{
		Bucket:       Ptr("oss-bucket"),
		Objects:      []DeleteObject{{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}, {Key: Ptr("key2.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****")}},
		EncodingType: Ptr("url"),
		Quiet:        true,
	}
	input = &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"delete": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, input.Body, strings.NewReader("<Delete><Quiet>true</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****</VersionId></Object></Delete>"))
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["delete"])
	assert.Equal(t, input.Parameters["encoding-type"], "url")

	request = &DeleteMultipleObjectsRequest{
		Bucket:       Ptr("oss-bucket"),
		Objects:      []DeleteObject{{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}, {Key: Ptr("key2.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****")}},
		EncodingType: Ptr("url"),
		Quiet:        true,
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"delete": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, input.Body, strings.NewReader("<Delete><Quiet>true</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****</VersionId></Object></Delete>"))
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["delete"])
	assert.Equal(t, input.Parameters["encoding-type"], "url")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_DeleteMultipleObjects(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"6555A936CA31DC333143****"},
			"Date":             {"Thu, 16 Nov 2023 05:31:34 GMT"},
		},
	}
	result := &DeleteMultipleObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555A936CA31DC333143****")
	assert.Equal(t, result.Headers.Get("Date"), "Thu, 16 Nov 2023 05:31:34 GMT")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &DeleteMultipleObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Len(t, result.DeletedObjects, 2)
	assert.Equal(t, *result.DeletedObjects[0].Key, "key1.txt")
	assert.Equal(t, result.DeletedObjects[0].DeleteMarker, true)
	assert.Equal(t, *result.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
	assert.Nil(t, result.DeletedObjects[0].VersionId)
	assert.Equal(t, *result.DeletedObjects[1].Key, "key2.txt")
	assert.Equal(t, result.DeletedObjects[1].DeleteMarker, true)
	assert.Equal(t, *result.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
	assert.Nil(t, result.DeletedObjects[1].VersionId)

	body = `<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <VersionId>CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmYx****</VersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <VersionId>CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmY1****</VersionId>
  </Deleted>
</DeleteResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &DeleteMultipleObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Len(t, result.DeletedObjects, 2)
	assert.Equal(t, *result.DeletedObjects[0].Key, "key1.txt")
	assert.False(t, result.DeletedObjects[0].DeleteMarker)
	assert.Nil(t, result.DeletedObjects[0].DeleteMarkerVersionId)
	assert.Equal(t, *result.DeletedObjects[0].VersionId, "CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmYx****")
	assert.Equal(t, *result.DeletedObjects[1].Key, "key2.txt")
	assert.False(t, result.DeletedObjects[1].DeleteMarker)
	assert.Nil(t, result.DeletedObjects[1].DeleteMarkerVersionId)
	assert.Equal(t, *result.DeletedObjects[1].VersionId, "CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmY1****")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>go-sdk-v1%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F</Key>
    <VersionId>CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmYx****</VersionId>
  </Deleted>
</DeleteResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &DeleteMultipleObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Len(t, result.DeletedObjects, 1)
	assert.Equal(t, *result.DeletedObjects[0].Key, "go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f")
	assert.False(t, result.DeletedObjects[0].DeleteMarker)
	assert.Nil(t, result.DeletedObjects[0].DeleteMarkerVersionId)
	assert.Equal(t, *result.DeletedObjects[0].VersionId, "CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmYx****")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "MalformedXML",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "MalformedXML")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MalformedXML</Code>
  <Message>The XML you provided was not well-formed or did not validate against our published schema.</Message>
  <RequestId>6555AC764311A73931E0****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <ErrorDetail>the root node is not named Delete.</ErrorDetail>
  <EC>0016-00000608</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000608</RecommendDoc>
</Error>`
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "MalformedXML",
		Headers: http.Header{
			"X-Oss-Request-Id": {"6555AC764311A73931E0****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "MalformedXML")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555AC764311A73931E0****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_HeadObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *HeadObjectRequest
	var input *OperationInput
	var err error

	request = &HeadObjectRequest{}
	input = &OperationInput{
		OpName: "HeadObject",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &HeadObjectRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "HeadObject",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &HeadObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "HeadObject",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)

	request = &HeadObjectRequest{
		Bucket:    Ptr("oss-bucket"),
		Key:       Ptr("oss-key"),
		VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
	}
	input = &OperationInput{
		OpName: "HeadObject",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Equal(t, input.Parameters["versionId"], "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")

	request = &HeadObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "HeadObject",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_HeadObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":    {"6555A936CA31DC333143****"},
			"Date":                {"Thu, 16 Nov 2023 05:31:34 GMT"},
			"x-oss-object-type":   {"Normal"},
			"x-oss-storage-class": {"Archive"},
			"Last-Modified":       {"Fri, 24 Feb 2018 09:41:56 GMT"},
			"Content-Length":      {"344606"},
			"Content-Type":        {"image/jpg"},
			"ETag":                {"\"fba9dede5f27731c9771645a3986****\""},
		},
	}
	result := &HeadObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555A936CA31DC333143****")
	assert.Equal(t, result.Headers.Get("Date"), "Thu, 16 Nov 2023 05:31:34 GMT")
	assert.Equal(t, *result.ETag, "\"fba9dede5f27731c9771645a3986****\"")
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, *result.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, *result.StorageClass, "Archive")
	assert.Equal(t, result.ContentLength, int64(344606))
	assert.Equal(t, *result.ContentType, "image/jpg")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":      {"6555A936CA31DC333143****"},
			"Date":                  {"Thu, 16 Nov 2023 05:31:34 GMT"},
			"x-oss-object-type":     {"Normal"},
			"x-oss-storage-class":   {"ColdArchive"},
			"Last-Modified":         {"Fri, 24 Feb 2018 09:41:56 GMT"},
			"Content-Length":        {"344606"},
			"Content-Type":          {"image/jpg"},
			"ETag":                  {"\"fba9dede5f27731c9771645a3986****\""},
			"x-oss-transition-time": {"Thu, 31 Oct 2024 00:24:17 GMT"},
			"x-oss-restore":         {"ongoing-request=\"false\", expiry-date=\"Fri, 08 Nov 2024 08:15:52 GMT\""},
		},
	}
	result = &HeadObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555A936CA31DC333143****")
	assert.Equal(t, result.Headers.Get("Date"), "Thu, 16 Nov 2023 05:31:34 GMT")
	assert.Equal(t, *result.ETag, "\"fba9dede5f27731c9771645a3986****\"")
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, *result.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, *result.StorageClass, "ColdArchive")
	assert.Equal(t, result.ContentLength, int64(344606))
	assert.Equal(t, *result.ContentType, "image/jpg")
	assert.Equal(t, *result.TransitionTime, time.Date(2024, time.October, 31, 00, 24, 17, 0, time.UTC))
	assert.Equal(t, *result.Restore, "ongoing-request=\"false\", expiry-date=\"Fri, 08 Nov 2024 08:15:52 GMT\"")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":    {"5CAC3B40B7AEADE01700****"},
			"Date":                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":        {"text/xml"},
			"x-oss-object-type":   {"Normal"},
			"x-oss-storage-class": {"Archive"},
			"Last-Modified":       {"Fri, 24 Feb 2023 09:41:56 GMT"},
			"Content-Length":      {"481827"},
			"ETag":                {"\"A082B659EF78733A5A042FA253B1****\""},
			"x-oss-version-Id":    {"CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"},
		},
	}
	result = &HeadObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5CAC3B40B7AEADE01700****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, *result.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, *result.StorageClass, "Archive")
	assert.Equal(t, result.ContentLength, int64(481827))
	assert.Equal(t, *result.ContentType, "text/xml")
	assert.Equal(t, *result.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
	assert.Equal(t, *result.ETag, "\"A082B659EF78733A5A042FA253B1****\"")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":    {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":        {"image/jpg"},
			"x-oss-object-type":   {"Normal"},
			"x-oss-restore":       {"ongoing-request=\"true\""},
			"x-oss-storage-class": {"Archive"},
			"Last-Modified":       {"Fri, 24 Feb 2023 09:41:59 GMT"},
			"Content-Length":      {"481827"},
			"ETag":                {"\"A082B659EF78733A5A042FA253B1****\""},
		},
	}
	result = &HeadObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, *result.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
	assert.Equal(t, *result.StorageClass, "Archive")
	assert.Equal(t, result.ContentLength, int64(481827))
	assert.Equal(t, *result.ContentType, "image/jpg")
	assert.Equal(t, *result.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
	assert.Equal(t, *result.Restore, "ongoing-request=\"true\"")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":                    {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":                                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":                        {"image/jpg"},
			"x-oss-object-type":                   {"Normal"},
			"x-oss-restore":                       {"ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\""},
			"x-oss-storage-class":                 {"Archive"},
			"x-oss-server-side-encryption":        {"KMS"},
			"x-oss-server-side-data-encryption":   {"SM4"},
			"x-oss-server-side-encryption-key-id": {"9468da86-3509-4f8d-a61e-6eab1eac****"},
			"Content-Length":                      {"481827"},
			"ETag":                                {"\"A082B659EF78733A5A042FA253B1****\""},
			"Last-Modified":                       {"Fri, 24 Feb 2023 09:41:59 GMT"},
		},
	}
	result = &HeadObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, *result.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
	assert.Equal(t, *result.StorageClass, "Archive")
	assert.Equal(t, result.ContentLength, int64(481827))
	assert.Equal(t, *result.ContentType, "image/jpg")
	assert.Equal(t, *result.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
	assert.Equal(t, *result.Restore, "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"")
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.Equal(t, *result.ServerSideEncryptionKeyId, "9468da86-3509-4f8d-a61e-6eab1eac****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchKey",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchKey")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidTargetType",
		Headers: http.Header{
			"X-Oss-Request-Id": {"6555AC764311A73931E0****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidTargetType")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555AC764311A73931E0****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetObjectMeta(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetObjectMetaRequest
	var input *OperationInput
	var err error

	request = &GetObjectMetaRequest{}
	input = &OperationInput{
		OpName: "GetObjectMeta",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"objectMeta": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetObjectMetaRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "GetObjectMeta",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"objectMeta": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetObjectMetaRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "GetObjectMeta",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"objectMeta": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["objectMeta"])
	request = &GetObjectMetaRequest{
		Bucket:    Ptr("oss-bucket"),
		Key:       Ptr("oss-key"),
		VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
	}
	input = &OperationInput{
		OpName: "GetObjectMeta",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"objectMeta": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["objectMeta"])
	assert.Equal(t, input.Parameters["versionId"], "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")

	request = &GetObjectMetaRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "GetObjectMeta",
		Method: "HEAD",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"objectMeta": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["objectMeta"])
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_GetObjectMeta(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"6555A936CA31DC333143****"},
			"Date":             {"Thu, 16 Nov 2023 05:31:34 GMT"},
			"Last-Modified":    {"Fri, 24 Feb 2018 09:41:56 GMT"},
			"Content-Length":   {"344606"},
			"ETag":             {"\"fba9dede5f27731c9771645a3986****\""},
		},
	}
	result := &GetObjectMetaResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555A936CA31DC333143****")
	assert.Equal(t, result.Headers.Get("Date"), "Thu, 16 Nov 2023 05:31:34 GMT")
	assert.Equal(t, *result.ETag, "\"fba9dede5f27731c9771645a3986****\"")
	assert.Equal(t, *result.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, result.ContentLength, int64(344606))

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":      {"6555A936CA31DC333143****"},
			"Date":                  {"Thu, 16 Nov 2023 05:31:34 GMT"},
			"x-oss-object-type":     {"Normal"},
			"Last-Modified":         {"Fri, 24 Feb 2018 09:41:56 GMT"},
			"Content-Length":        {"344606"},
			"ETag":                  {"\"fba9dede5f27731c9771645a3986****\""},
			"x-oss-transition-time": {"Thu, 31 Oct 2024 00:24:17 GMT"},
		},
	}
	result = &GetObjectMetaResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555A936CA31DC333143****")
	assert.Equal(t, result.Headers.Get("Date"), "Thu, 16 Nov 2023 05:31:34 GMT")
	assert.Equal(t, *result.ETag, "\"fba9dede5f27731c9771645a3986****\"")
	assert.Equal(t, *result.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, result.ContentLength, int64(344606))
	assert.Equal(t, *result.TransitionTime, time.Date(2024, time.October, 31, 00, 24, 17, 0, time.UTC))

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5CAC3B40B7AEADE01700****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Last-Modified":    {"Fri, 24 Feb 2023 09:41:56 GMT"},
			"Content-Length":   {"481827"},
			"ETag":             {"\"A082B659EF78733A5A042FA253B1****\""},
			"x-oss-version-Id": {"CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"},
		},
	}
	result = &GetObjectMetaResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5CAC3B40B7AEADE01700****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, *result.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, result.ContentLength, int64(481827))
	assert.Equal(t, *result.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
	assert.Equal(t, *result.ETag, "\"A082B659EF78733A5A042FA253B1****\"")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":       {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":                   {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"x-oss-last-access-time": {"Thu, 14 Oct 2021 11:49:05 GMT"},
			"Last-Modified":          {"Fri, 24 Feb 2020 09:41:59 GMT"},
			"Content-Length":         {"481827"},
			"ETag":                   {"\"A082B659EF78733A5A042FA253B1****\""},
		},
	}
	result = &GetObjectMetaResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, *result.LastModified, time.Date(2020, time.February, 24, 9, 41, 59, 0, time.UTC))
	assert.Equal(t, *result.LastAccessTime, time.Date(2021, time.October, 14, 11, 49, 05, 0, time.UTC))
	assert.Equal(t, result.ContentLength, int64(481827))
	assert.Equal(t, *result.ETag, "\"A082B659EF78733A5A042FA253B1****\"")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchKey",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchKey")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidTargetType",
		Headers: http.Header{
			"X-Oss-Request-Id": {"6555AC764311A73931E0****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidTargetType")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555AC764311A73931E0****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_RestoreObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *RestoreObjectRequest
	var input *OperationInput
	var err error

	request = &RestoreObjectRequest{}
	input = &OperationInput{
		OpName: "RestoreObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"restore": "",
		},
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &RestoreObjectRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "RestoreObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"restore": "",
		},
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &RestoreObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "RestoreObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"restore": "",
		},
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["restore"])
	request = &RestoreObjectRequest{
		Bucket:    Ptr("oss-bucket"),
		Key:       Ptr("oss-key"),
		VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
	}
	input = &OperationInput{
		OpName: "RestoreObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"restore": "",
		},
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["restore"])
	assert.Equal(t, input.Parameters["versionId"], "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")

	request = &RestoreObjectRequest{
		Bucket:    Ptr("oss-bucket"),
		Key:       Ptr("oss-key"),
		VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
		RestoreRequest: &RestoreRequest{
			Days: int32(2),
			Tier: Ptr("Standard"),
		},
	}
	input = &OperationInput{
		OpName: "RestoreObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"restore": "",
		},
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["restore"])
	assert.Equal(t, input.Parameters["versionId"], "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
	data, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days><JobParameters><Tier>Standard</Tier></JobParameters></RestoreRequest>")

	request = &RestoreObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		RequestPayer: Ptr("requester"),
		RestoreRequest: &RestoreRequest{
			Days: int32(2),
			Tier: Ptr("Standard"),
		},
	}
	input = &OperationInput{
		OpName: "RestoreObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"restore": "",
		},
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["restore"])
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days><JobParameters><Tier>Standard</Tier></JobParameters></RestoreRequest>")
}

func TestUnmarshalOutput_RestoreObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 202,
		Status:     "Accepted",
		Headers: http.Header{
			"X-Oss-Request-Id": {"6555A936CA31DC333143****"},
			"Date":             {"Thu, 16 Nov 2023 05:31:34 GMT"},
		},
	}
	result := &RestoreObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 202)
	assert.Equal(t, result.Status, "Accepted")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555A936CA31DC333143****")
	assert.Equal(t, result.Headers.Get("Date"), "Thu, 16 Nov 2023 05:31:34 GMT")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":              {"5CAC3B40B7AEADE01700****"},
			"Date":                          {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"x-oss-object-restore-priority": {"Standard"},
		},
	}
	result = &RestoreObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5CAC3B40B7AEADE01700****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, *result.RestorePriority, "Standard")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"x-oss-version-id": {"CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"},
		},
	}
	result = &RestoreObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, *result.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchKey",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchKey")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "OperationNotSupported",
		Headers: http.Header{
			"X-Oss-Request-Id": {"6555AC764311A73931E0****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "OperationNotSupported")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555AC764311A73931E0****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_PutObjectAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutObjectAclRequest
	var input *OperationInput
	var err error

	request = &PutObjectAclRequest{}
	input = &OperationInput{
		OpName: "PutObjectAcl",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutObjectAclRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "PutObjectAcl",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutObjectAclRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
		Acl:    ObjectACLPrivate,
	}
	input = &OperationInput{
		OpName: "PutObjectAcl",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &PutObjectAclRequest{
		Bucket:    Ptr("oss-demo"),
		Key:       Ptr("oss-object"),
		Acl:       ObjectACLPrivate,
		VersionId: Ptr("CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****"),
	}
	input = &OperationInput{
		OpName: "PutObjectAcl",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &PutObjectAclRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		Acl:          ObjectACLPrivate,
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "PutObjectAcl",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_PutObjectAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &PutObjectAclResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
			"X-Oss-Version-Id": {"CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.VersionId, "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****")
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetObjectAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetObjectAclRequest
	var input *OperationInput
	var err error

	request = &GetObjectAclRequest{}
	input = &OperationInput{
		OpName: "GetObjectAcl",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetObjectAclRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetObjectAcl",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetObjectAclRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "GetObjectAcl",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &GetObjectAclRequest{
		Bucket:    Ptr("oss-demo"),
		Key:       Ptr("oss-object"),
		VersionId: Ptr("CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****"),
	}
	input = &OperationInput{
		OpName: "GetObjectAcl",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	request = &GetObjectAclRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "GetObjectAcl",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_GetObjectAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012****</ID>
        <DisplayName>0022012****</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read</Grant>
    </AccessControlList>
</AccessControlPolicy>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &GetObjectAclResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
	assert.Equal(t, *result.ACL, "public-read")
	assert.Equal(t, *result.Owner.DisplayName, "0022012****")
	assert.Equal(t, *result.Owner.ID, "0022012****")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<AccessControlPolicy>
    <Owner>
        <ID>1234513715092****</ID>
        <DisplayName>1234513715092****</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>private</Grant>
    </AccessControlList>
</AccessControlPolicy>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
			"X-Oss-Version-Id": {"CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
	assert.Equal(t, *result.ACL, "private")
	assert.Equal(t, *result.Owner.DisplayName, "1234513715092****")
	assert.Equal(t, *result.Owner.ID, "1234513715092****")
	assert.Equal(t, *result.VersionId, "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
}

func TestMarshalInput_InitiateMultipartUpload(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *InitiateMultipartUploadRequest
	var input *OperationInput
	var err error

	request = &InitiateMultipartUploadRequest{}
	input = &OperationInput{
		OpName: "InitiateMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"uploads": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &InitiateMultipartUploadRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "InitiateMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"uploads": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &InitiateMultipartUploadRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "InitiateMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"uploads": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	request = &InitiateMultipartUploadRequest{
		Bucket:                    Ptr("oss-bucket"),
		Key:                       Ptr("oss-key"),
		CacheControl:              Ptr("no-cache"),
		ContentDisposition:        Ptr("attachment"),
		ContentEncoding:           Ptr("utf-8"),
		Expires:                   Ptr("2022-10-12T00:00:00.000Z"),
		ForbidOverwrite:           Ptr("false"),
		ServerSideEncryption:      Ptr("KMS"),
		ServerSideDataEncryption:  Ptr("SM4"),
		ServerSideEncryptionKeyId: Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
		StorageClass:              StorageClassStandard,
		Metadata: map[string]string{
			"name":  "walker",
			"email": "demo@aliyun.com",
		},
		Tagging: Ptr("TagA=B&TagC=D"),
	}

	input = &OperationInput{
		OpName: "InitiateMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"uploads": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Headers["Cache-Control"], "no-cache")
	assert.Equal(t, input.Headers["Content-Disposition"], "attachment")
	assert.Equal(t, input.Headers["x-oss-meta-name"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-email"], "demo@aliyun.com")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "KMS")
	assert.Equal(t, input.Headers["x-oss-server-side-data-encryption"], "SM4")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption-key-id"], "9468da86-3509-4f8d-a61e-6eab1eac****")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "false")
	assert.Equal(t, input.Headers["Content-Encoding"], "utf-8")
	assert.Equal(t, input.Headers["Content-MD5"], "1B2M2Y8AsgTpgAmY7PhCfg==")
	assert.Equal(t, input.Headers["Expires"], "2022-10-12T00:00:00.000Z")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=B&TagC=D")
	assert.Empty(t, input.Parameters["uploads"])
	assert.Nil(t, input.OpMetadata.values)

	request = &InitiateMultipartUploadRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "InitiateMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"uploads": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_InitiateMultipartUpload(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
    <Bucket>oss-example</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
</InitiateMultipartUploadResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &InitiateMultipartUploadResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	assert.Equal(t, *result.Bucket, "oss-example")
	assert.Equal(t, *result.Key, "multipart.data")
	assert.Equal(t, *result.UploadId, "0004B9894A22E5B1888A1E29F823****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
}

func TestMarshalInput_UploadPart(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UploadPartRequest
	var input *OperationInput
	var err error

	request = &UploadPartRequest{}
	input = &OperationInput{
		OpName: "UploadPart",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &UploadPartRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "UploadPart",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &UploadPartRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "UploadPart",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &UploadPartRequest{
		Bucket:     Ptr("oss-demo"),
		Key:        Ptr("oss-object"),
		PartNumber: int32(1),
	}
	input = &OperationInput{
		OpName: "UploadPart",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &UploadPartRequest{
		Bucket:     Ptr("oss-demo"),
		Key:        Ptr("oss-object"),
		PartNumber: int32(1),
		UploadId:   Ptr("0004B9895DBBB6EC9****"),
	}
	input = &OperationInput{
		OpName: "UploadPart",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["partNumber"], "1")
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Nil(t, input.OpMetadata.values)

	request = &UploadPartRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		PartNumber:   int32(1),
		UploadId:     Ptr("0004B9895DBBB6EC9****"),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["partNumber"], "1")
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Equal(t, input.Headers["x-oss-traffic-limit"], strconv.FormatInt(100*1024*8, 10))
	assert.Nil(t, input.OpMetadata.values)

	request = &UploadPartRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		PartNumber:   int32(1),
		UploadId:     Ptr("0004B9895DBBB6EC9****"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "UploadPart",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["partNumber"], "1")
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Nil(t, input.OpMetadata.values)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_UploadPart(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"534B371674E88A4D8906****"},
			"ETag":                 {"\"7265F4D211B56873A381D321F586****\""},
			"Date":                 {"Wed, 22 Feb 2012 08:32:21 GMT"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7Ph****"},
			"x-oss-hash-crc64ecma": {"316181249502703*****"},
		},
	}
	result := &UploadPartResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, *result.ETag, "\"7265F4D211B56873A381D321F586****\"")
	assert.Equal(t, *result.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
	assert.Equal(t, *result.HashCRC64, "316181249502703*****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
}

func TestMarshalInput_UploadPartCopy(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UploadPartCopyRequest
	var input *OperationInput
	var err error

	request = &UploadPartCopyRequest{}
	source := encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &UploadPartCopyRequest{
		Bucket: Ptr("oss-demo"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &UploadPartCopyRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &UploadPartCopyRequest{
		Bucket:     Ptr("oss-demo"),
		Key:        Ptr("oss-object"),
		PartNumber: int32(1),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &UploadPartCopyRequest{
		Bucket:     Ptr("oss-demo"),
		Key:        Ptr("oss-object"),
		PartNumber: int32(1),
		UploadId:   Ptr("0004B9895DBBB6EC9****"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &UploadPartCopyRequest{
		Bucket:       Ptr("oss-dest-bucket"),
		Key:          Ptr("oss-dest-object"),
		PartNumber:   int32(1),
		UploadId:     Ptr("0004B9895DBBB6EC9****"),
		SourceKey:    Ptr("oss-src-dir/oss-src-obj+123"),
		SourceBucket: Ptr("oss-src-bucket"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["partNumber"], "1")
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-src-bucket/oss-src-dir/oss-src-obj%2B123")
	assert.Nil(t, input.OpMetadata.values)

	request = &UploadPartCopyRequest{
		Bucket:       Ptr("oss-dest-bucket"),
		Key:          Ptr("oss-dest-object"),
		PartNumber:   int32(1),
		UploadId:     Ptr("0004B9895DBBB6EC9****"),
		SourceKey:    Ptr("oss-src-dir/oss-src-obj"),
		SourceBucket: Ptr("oss-src-bucket"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["partNumber"], "1")
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-src-bucket/oss-src-dir/oss-src-obj")
	assert.Nil(t, input.OpMetadata.values)

	request = &UploadPartCopyRequest{
		Bucket:          Ptr("oss-dest-bucket"),
		Key:             Ptr("oss-dest-object"),
		PartNumber:      int32(1),
		UploadId:        Ptr("0004B9895DBBB6EC9****"),
		SourceKey:       Ptr("oss-src-dir/oss-src-obj"),
		SourceBucket:    Ptr("oss-src-bucket"),
		SourceVersionId: Ptr("CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["partNumber"], "1")
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-src-bucket/oss-src-dir/oss-src-obj"+"?versionId=CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****")
	assert.Nil(t, input.OpMetadata.values)

	request = &UploadPartCopyRequest{
		Bucket:       Ptr("oss-dest-bucket"),
		Key:          Ptr("oss-dest-object"),
		PartNumber:   int32(1),
		UploadId:     Ptr("0004B9895DBBB6EC9****"),
		SourceKey:    Ptr("oss-src-dir/oss-src-obj"),
		SourceBucket: Ptr("oss-src-bucket"),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["partNumber"], "1")
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-src-bucket/oss-src-dir/oss-src-obj")
	assert.Equal(t, input.Headers["x-oss-traffic-limit"], strconv.FormatInt(100*1024*8, 10))
	assert.Nil(t, input.OpMetadata.values)

	request = &UploadPartCopyRequest{
		Bucket:       Ptr("oss-dest-bucket"),
		Key:          Ptr("oss-dest-object"),
		PartNumber:   int32(1),
		UploadId:     Ptr("0004B9895DBBB6EC9****"),
		SourceKey:    Ptr("oss-src-dir/oss-src-obj+123"),
		SourceBucket: Ptr("oss-src-bucket"),
		RequestPayer: Ptr("requester"),
	}
	source = encodeSourceObject(request)
	input = &OperationInput{
		OpName: "UploadPartCopy",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			"x-oss-copy-source": source,
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["partNumber"], "1")
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-src-bucket/oss-src-dir/oss-src-obj%2B123")
	assert.Nil(t, input.OpMetadata.values)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_UploadPartCopy(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Date":             {"Wed, 22 Feb 2012 08:32:21 GMT"},
		},
	}
	result := &UploadPartCopyResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderDate), "Wed, 22 Feb 2012 08:32:21 GMT")
	assert.Equal(t, *result.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
	assert.Equal(t, *result.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))

	body = `<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id":             {"534B371674E88A4D8906****"},
			"Date":                         {"Wed, 22 Feb 2012 08:32:21 GMT"},
			"x-oss-copy-source-version-id": {"CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderDate), "Wed, 22 Feb 2012 08:32:21 GMT")
	assert.Equal(t, *result.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
	assert.Equal(t, *result.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
	assert.Equal(t, *result.VersionId, "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
}

func TestMarshalInput_CompleteMultipartUpload(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CompleteMultipartUploadRequest
	var input *OperationInput
	var err error

	request = &CompleteMultipartUploadRequest{}
	input = &OperationInput{
		OpName: "CompleteMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CompleteMultipartUploadRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "CompleteMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CompleteMultipartUploadRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		UploadId:     Ptr("0004B9895DBBB6EC9****"),
		EncodingType: Ptr("url"),
		CompleteMultipartUpload: &CompleteMultipartUpload{
			Parts: []UploadPart{
				{PartNumber: int32(3), ETag: Ptr("\"3349DC700140D7F86A0784842780****\"")},
				{PartNumber: int32(1), ETag: Ptr("\"8EFDA8BE206636A695359836FE0A****\"")},
				{PartNumber: int32(2), ETag: Ptr("\"8C315065167132444177411FDA14****\"")},
			},
		},
	}
	if request.CompleteMultipartUpload != nil && len(request.CompleteMultipartUpload.Parts) > 0 {
		sort.Sort(UploadParts(request.CompleteMultipartUpload.Parts))
	}
	input = &OperationInput{
		OpName: "CompleteMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Equal(t, input.Parameters["encoding-type"], "url")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), `<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>&#34;8EFDA8BE206636A695359836FE0A****&#34;</ETag></Part><Part><PartNumber>2</PartNumber><ETag>&#34;8C315065167132444177411FDA14****&#34;</ETag></Part><Part><PartNumber>3</PartNumber><ETag>&#34;3349DC700140D7F86A0784842780****&#34;</ETag></Part></CompleteMultipartUpload>`)
	callbackVal := base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`))
	callbackVar := base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`))
	request = &CompleteMultipartUploadRequest{
		Bucket:          Ptr("oss-dest-bucket"),
		Key:             Ptr("oss-dest-object"),
		UploadId:        Ptr("0004B9895DBBB6EC9****"),
		ForbidOverwrite: Ptr("false"),
		CompleteAll:     Ptr("yes"),
		Callback:        Ptr(callbackVal),
		CallbackVar:     Ptr(callbackVar),
	}
	input = &OperationInput{
		OpName: "CompleteMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "false")
	assert.Equal(t, input.Headers["x-oss-callback"], callbackVal)
	assert.Equal(t, input.Headers["x-oss-callback-var"], callbackVar)
	assert.Nil(t, input.OpMetadata.values)

	request = &CompleteMultipartUploadRequest{
		Bucket:          Ptr("oss-dest-bucket"),
		Key:             Ptr("oss-dest-object"),
		UploadId:        Ptr("0004B9895DBBB6EC9****"),
		ForbidOverwrite: Ptr("false"),
		CompleteAll:     Ptr("yes"),
		RequestPayer:    Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "CompleteMultipartUpload",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["uploadId"], "0004B9895DBBB6EC9****")
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "false")
	assert.Nil(t, input.OpMetadata.values)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_CompleteMultipartUpload(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult>
  <EncodingType>url</EncodingType>
  <Location>http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data</Location>
  <Bucket>oss-example</Bucket>
  <Key>demo%2Fmultipart.data</Key>
  <ETag>"097DE458AD02B5F89F9D0530231876****"</ETag>
</CompleteMultipartUploadResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Date":             {"Wed, 22 Feb 2012 08:32:21 GMT"},
		},
	}
	result := &CompleteMultipartUploadResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderDate), "Wed, 22 Feb 2012 08:32:21 GMT")
	assert.Equal(t, *result.ETag, "\"097DE458AD02B5F89F9D0530231876****\"")
	assert.Equal(t, *result.Location, "http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data")
	assert.Equal(t, *result.EncodingType, "url")
	assert.Equal(t, *result.Bucket, "oss-example")
	assert.Equal(t, *result.Key, "demo/multipart.data")

	body = `{"filename":"oss-obj.txt","size":"100","mimeType":"","x":"a","b":"b"}`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id":     {"534B371674E88A4D8906****"},
			"Date":                 {"Wed, 22 Feb 2012 08:32:21 GMT"},
			"x-oss-version-id":     {"CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****"},
			"Content-Type":         {"application/json"},
			"x-oss-hash-crc64ecma": {"1206617243528768****"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalCallbackBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderDate), "Wed, 22 Feb 2012 08:32:21 GMT")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/json")

	jsonData, _ := json.Marshal(result.CallbackResult)
	assert.Nil(t, err)
	assert.NotEmpty(t, string(jsonData))
	assert.Equal(t, *result.VersionId, "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
}

func TestMarshalInput_AbortMultipartUpload(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *AbortMultipartUploadRequest
	var input *OperationInput
	var err error

	request = &AbortMultipartUploadRequest{}
	input = &OperationInput{
		OpName: "AbortMultipartUpload",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AbortMultipartUploadRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "AbortMultipartUpload",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AbortMultipartUploadRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "AbortMultipartUpload",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	uploadId := "0004B9895DBBB6E****"
	request = &AbortMultipartUploadRequest{
		Bucket:   Ptr("oss-demo"),
		Key:      Ptr("oss-object"),
		UploadId: Ptr(uploadId),
	}
	input = &OperationInput{
		OpName: "AbortMultipartUpload",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["uploadId"], uploadId)
	assert.Nil(t, input.OpMetadata.values)

	request = &AbortMultipartUploadRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		UploadId:     Ptr(uploadId),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "AbortMultipartUpload",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["uploadId"], uploadId)
	assert.Nil(t, input.OpMetadata.values)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_AbortMultipartUpload(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Date":             {"Wed, 22 Feb 2012 08:32:21 GMT"},
		},
	}
	result := &InitiateMultipartUploadResult{}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderDate), "Wed, 22 Feb 2012 08:32:21 GMT")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
}

func TestMarshalInput_ListMultipartUploads(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListMultipartUploadsRequest
	var input *OperationInput
	var err error

	request = &ListMultipartUploadsRequest{}
	input = &OperationInput{
		OpName: "ListMultipartUploads",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
			"uploads":       "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ListMultipartUploadsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "ListMultipartUploads",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
			"uploads":       "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &ListMultipartUploadsRequest{
		Bucket:         Ptr("oss-demo"),
		Delimiter:      Ptr("/"),
		Prefix:         Ptr("prefix"),
		EncodingType:   Ptr("url"),
		KeyMarker:      Ptr("89F0105AA66942638E35300618DF5EE7"),
		MaxUploads:     int32(10),
		UploadIdMarker: Ptr("89F0105AA66942638E35300618DF5EE7"),
	}
	input = &OperationInput{
		OpName: "ListMultipartUploads",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
			"uploads":       "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Empty(t, input.Parameters["uploads"])
	assert.Equal(t, input.Parameters["delimiter"], "/")
	assert.Equal(t, input.Parameters["prefix"], "prefix")
	assert.Equal(t, input.Parameters["key-marker"], "89F0105AA66942638E35300618DF5EE7")
	assert.Equal(t, input.Parameters["max-uploads"], "10")
	assert.Equal(t, input.Parameters["upload-id-marker"], "89F0105AA66942638E35300618DF5EE7")
	assert.Nil(t, input.OpMetadata.values)

	request = &ListMultipartUploadsRequest{
		Bucket:       Ptr("oss-demo"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "ListMultipartUploads",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
			"uploads":       "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_ListMultipartUploads(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
    <Bucket>oss-example</Bucket>
    <KeyMarker></KeyMarker>
    <UploadIdMarker></UploadIdMarker>
    <NextKeyMarker>oss.avi</NextKeyMarker>
    <NextUploadIdMarker>0004B99B8E707874FC2D692FA5D77D3F</NextUploadIdMarker>
    <Delimiter></Delimiter>
    <Prefix></Prefix>
    <MaxUploads>1000</MaxUploads>
    <IsTruncated>false</IsTruncated>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF518A1FE585B0C9360DC4C8</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>oss.avi</Key>
        <UploadId>0004B99B8E707874FC2D692FA5D7****</UploadId>
        <Initiated>2012-02-23T06:14:27.000Z</Initiated>
    </Upload>
</ListMultipartUploadsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Date":             {"Wed, 22 Feb 2012 08:32:21 GMT"},
		},
	}
	result := &ListMultipartUploadsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderDate), "Wed, 22 Feb 2012 08:32:21 GMT")
	assert.Equal(t, *result.Bucket, "oss-example")
	assert.Equal(t, *result.KeyMarker, "")
	assert.Equal(t, *result.UploadIdMarker, "")
	assert.Equal(t, *result.NextKeyMarker, "oss.avi")
	assert.Equal(t, *result.NextUploadIdMarker, "0004B99B8E707874FC2D692FA5D77D3F")
	assert.Equal(t, *result.Delimiter, "")
	assert.Equal(t, *result.Prefix, "")
	assert.Equal(t, result.MaxUploads, int32(1000))
	assert.Equal(t, result.IsTruncated, false)
	assert.Len(t, result.Uploads, 3)
	assert.Equal(t, *result.Uploads[0].Key, "multipart.data")
	assert.Equal(t, *result.Uploads[0].UploadId, "0004B999EF518A1FE585B0C9360DC4C8")
	assert.Equal(t, *result.Uploads[0].Initiated, time.Date(2012, time.February, 23, 4, 18, 23, 0, time.UTC))

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
  <EncodingType>url</EncodingType>
  <Bucket>oss-example</Bucket>
  <KeyMarker></KeyMarker>
  <UploadIdMarker></UploadIdMarker>
  <NextKeyMarker>oss.avi</NextKeyMarker>
  <NextUploadIdMarker>89F0105AA66942638E35300618DF****</NextUploadIdMarker>
  <Delimiter></Delimiter>
  <Prefix></Prefix>
  <MaxUploads>1000</MaxUploads>
  <IsTruncated>false</IsTruncated>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>0214A87687F040F1BA4D83AB17C9****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:45:57.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>3AE2ED7A60E04AFE9A5287055D37****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:03:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>47E0E90F5DCB4AD5B3C4CD886CB0****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:02:11.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>A89E0E28E2E948A1BFF6FD5CDAFF****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T06:57:03.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>B18E1DCDB6964F5CB197F5F6B26A****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:02.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>5209986C3A96486EA16B9C52C160****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:34:47.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>63B652FA2C1342DCB3CCCC86D748****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:28:46.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>6F67B34BCA3C481F887D73508A07****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:32:12.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>89F0105AA66942638E35300618D****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:37:53.000Z</Initiated>
  </Upload>
</ListMultipartUploadsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Date":             {"Wed, 22 Feb 2012 08:32:21 GMT"},
		},
	}
	result = &ListMultipartUploadsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderDate), "Wed, 22 Feb 2012 08:32:21 GMT")
	assert.Equal(t, *result.Bucket, "oss-example")
	assert.Equal(t, *result.KeyMarker, "")
	assert.Equal(t, *result.UploadIdMarker, "")
	assert.Equal(t, *result.NextKeyMarker, "oss.avi")
	assert.Equal(t, *result.NextUploadIdMarker, "89F0105AA66942638E35300618DF****")
	assert.Equal(t, *result.Delimiter, "")
	assert.Equal(t, *result.Prefix, "")
	assert.Equal(t, result.MaxUploads, int32(1000))
	assert.Equal(t, result.IsTruncated, false)
	assert.Len(t, result.Uploads, 10)
	assert.Equal(t, *result.Uploads[0].Key, "demo/gp-\f\n\v")
	assert.Equal(t, *result.Uploads[0].UploadId, "0214A87687F040F1BA4D83AB17C9****")
	assert.Equal(t, *result.Uploads[0].Initiated, time.Date(2023, time.November, 22, 5, 45, 57, 0, time.UTC))

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
}

func TestMarshalInput_ListParts(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListPartsRequest
	var input *OperationInput
	var err error

	request = &ListPartsRequest{}
	input = &OperationInput{
		OpName: "ListParts",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ListPartsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "ListParts",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ListPartsRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "ListParts",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ListPartsRequest{
		Bucket:   Ptr("oss-demo"),
		Key:      Ptr("oss-object"),
		UploadId: Ptr("89F0105AA66942638E35300618DF5EE7"),
	}
	input = &OperationInput{
		OpName: "ListParts",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["uploadId"], "89F0105AA66942638E35300618DF5EE7")
	assert.Nil(t, input.OpMetadata.values)

	request = &ListPartsRequest{
		Bucket:   Ptr("oss-demo"),
		Key:      Ptr("oss-object"),
		UploadId: Ptr("89F0105AA66942638E35300618DF5EE7"),
		MaxParts: int32(10),
	}
	input = &OperationInput{
		OpName: "ListParts",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["max-parts"], "10")
	assert.Equal(t, input.Parameters["uploadId"], "89F0105AA66942638E35300618DF5EE7")
	assert.Empty(t, input.Parameters["part-number-marker"])
	assert.Nil(t, input.OpMetadata.values)

	request = &ListPartsRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		UploadId:     Ptr("89F0105AA66942638E35300618DF5EE7"),
		MaxParts:     int32(10),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "ListParts",
		Method: "GET",
		Bucket: request.Bucket,
		Parameters: map[string]string{
			"encoding-type": "url",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["max-parts"], "10")
	assert.Equal(t, input.Parameters["uploadId"], "89F0105AA66942638E35300618DF5EE7")
	assert.Empty(t, input.Parameters["part-number-marker"])
	assert.Nil(t, input.OpMetadata.values)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_ListParts(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
    <Bucket>multipart_upload</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
    <NextPartNumberMarker>5</NextPartNumberMarker>
    <MaxParts>1000</MaxParts>
    <IsTruncated>false</IsTruncated>
    <Part>
        <PartNumber>1</PartNumber>
        <LastModified>2012-02-23T07:01:34.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>2</PartNumber>
        <LastModified>2012-02-23T07:01:12.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>5</PartNumber>
        <LastModified>2012-02-23T07:02:03.000Z</LastModified>
        <ETag>"7265F4D211B56873A381D321F586****"</ETag>
        <Size>1024</Size>
    </Part>
</ListPartsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Date":             {"Wed, 22 Feb 2012 08:32:21 GMT"},
		},
	}
	result := &ListPartsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderDate), "Wed, 22 Feb 2012 08:32:21 GMT")
	assert.Equal(t, *result.Bucket, "multipart_upload")
	assert.Equal(t, *result.Key, "multipart.data")
	assert.Empty(t, result.PartNumberMarker)
	assert.Equal(t, result.NextPartNumberMarker, int32(5))
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, result.MaxParts, int32(1000))
	assert.Len(t, result.Parts, 3)
	assert.Equal(t, result.Parts[0].PartNumber, int32(1))
	assert.Equal(t, *result.Parts[0].ETag, "\"3349DC700140D7F86A0784842780****\"")
	assert.Equal(t, *result.Parts[0].LastModified, time.Date(2012, time.February, 23, 7, 1, 34, 0, time.UTC))
	assert.Equal(t, result.Parts[0].Size, int64(6291456))

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
  <EncodingType>url</EncodingType>
  <Bucket>oss-bucket</Bucket>
  <Key>demo%2Fgp-%0C%0A%0B</Key>
  <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
  <StorageClass>Standard</StorageClass>
  <PartNumberMarker>0</PartNumberMarker>
  <NextPartNumberMarker>1</NextPartNumberMarker>
  <MaxParts>1000</MaxParts>
  <IsTruncated>false</IsTruncated>
  <Part>
    <PartNumber>1</PartNumber>
    <LastModified>2023-11-22T05:42:34.000Z</LastModified>
    <ETag>"CF3F46D505093571E916FCDD4967****"</ETag>
    <HashCrc64ecma>12066172435287683848</HashCrc64ecma>
    <Size>96316</Size>
  </Part>
</ListPartsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Date":             {"Wed, 22 Feb 2012 08:32:21 GMT"},
		},
	}
	result = &ListPartsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderDate), "Wed, 22 Feb 2012 08:32:21 GMT")
	assert.Equal(t, *result.Bucket, "oss-bucket")
	key, _ := url.QueryUnescape("demo%2Fgp-%0C%0A%0B")
	assert.Equal(t, *result.Key, key)
	assert.Empty(t, result.PartNumberMarker)
	assert.Equal(t, result.NextPartNumberMarker, int32(1))
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, result.MaxParts, int32(1000))
	assert.Len(t, result.Parts, 1)
	assert.Equal(t, result.Parts[0].PartNumber, int32(1))
	assert.Equal(t, *result.Parts[0].ETag, "\"CF3F46D505093571E916FCDD4967****\"")
	assert.Equal(t, *result.Parts[0].LastModified, time.Date(2023, time.November, 22, 5, 42, 34, 0, time.UTC))
	assert.Equal(t, result.Parts[0].Size, int64(96316))
	assert.Equal(t, *result.Parts[0].HashCRC64, "12066172435287683848")
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get(HeaderOssRequestID), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get(HTTPHeaderContentType), "application/xml")
}

func TestMarshalInput_PutSymlink(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutSymlinkRequest
	var input *OperationInput
	var err error

	request = &PutSymlinkRequest{}
	input = &OperationInput{
		OpName: "PutSymlink",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutSymlinkRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutSymlink",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key")

	request = &PutSymlinkRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "PutSymlink",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Target")

	request = &PutSymlinkRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
		Target: Ptr("oss-target-object"),
	}
	input = &OperationInput{
		OpName: "PutSymlink",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-symlink-target"], "oss-target-object")

	request = &PutSymlinkRequest{
		Bucket:          Ptr("oss-demo"),
		Key:             Ptr("oss-object"),
		Target:          Ptr("oss-target-object"),
		ForbidOverwrite: Ptr("true"),
		Acl:             ObjectACLPrivate,
		StorageClass:    StorageClassStandard,
		Metadata: map[string]string{
			"name":  "demo",
			"email": "demo@aliyun.com",
		},
	}
	input = &OperationInput{
		OpName: "PutSymlink",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-symlink-target"], "oss-target-object")
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "true")
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-meta-name"], "demo")
	assert.Equal(t, input.Headers["x-oss-meta-email"], "demo@aliyun.com")
	assert.Nil(t, input.OpMetadata.values)

	request = &PutSymlinkRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		Target:       Ptr("oss-target-object"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "PutSymlink",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-symlink-target"], "oss-target-object")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_PutSymlink(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutSymlinkResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"x-oss-version-id": {"CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"},
		},
	}
	result = &PutSymlinkResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutSymlinkResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetSymlink(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetSymlinkRequest
	var input *OperationInput
	var err error

	request = &GetSymlinkRequest{}
	input = &OperationInput{
		OpName: "GetSymlink",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &GetSymlinkRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetSymlink",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key")

	request = &GetSymlinkRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "GetSymlink",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &GetSymlinkRequest{
		Bucket:    Ptr("oss-demo"),
		Key:       Ptr("oss-object"),
		VersionId: Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
	}
	input = &OperationInput{
		OpName: "GetSymlink",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["versionId"], "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")

	request = &GetSymlinkRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "GetSymlink",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"symlink": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_GetSymlink(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"534B371674E88A4D8906****"},
			"x-oss-symlink-target": {"example.jpg"},
			"ETag":                 {"A797938C31D59EDD08D86188F6D5****"},
		},
	}
	result := &GetSymlinkResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Target, "example.jpg")
	assert.Equal(t, *result.ETag, "A797938C31D59EDD08D86188F6D5****")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"534B371674E88A4D8906****"},
			"x-oss-version-id":     {"CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"},
			"x-oss-symlink-target": {"example.jpg"},
			"ETag":                 {"A797938C31D59EDD08D86188F6D5****"},
		},
	}
	result = &GetSymlinkResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
	assert.Equal(t, *result.Target, "example.jpg")
	assert.Equal(t, *result.ETag, "A797938C31D59EDD08D86188F6D5****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetSymlinkResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_PutObjectTagging(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutObjectTaggingRequest
	var input *OperationInput
	var err error

	request = &PutObjectTaggingRequest{}
	input = &OperationInput{
		OpName: "PutObjectTagging",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutObjectTaggingRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutObjectTagging",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key")

	request = &PutObjectTaggingRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "PutObjectTagging",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Tagging")

	request = &PutObjectTaggingRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutObjectTagging",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	data, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(data), `<Tagging><TagSet><Tag><Key>k1</Key><Value>v1</Value></Tag></TagSet></Tagging>`)

	request = &PutObjectTaggingRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
				},
			},
		},
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "PutObjectTagging",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(data), `<Tagging><TagSet><Tag><Key>k1</Key><Value>v1</Value></Tag></TagSet></Tagging>`)
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_PutObjectTagging(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Date":             {"Mon, 18 Mar 2019 08:25:17 GMT"},
		},
	}
	result := &PutObjectTaggingResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Date"), "Mon, 18 Mar 2019 08:25:17 GMT")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"x-oss-version-id": {"CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"},
		},
	}
	result = &PutObjectTaggingResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutObjectTaggingResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetObjectTagging(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetObjectTaggingRequest
	var input *OperationInput
	var err error

	request = &GetObjectTaggingRequest{}
	input = &OperationInput{
		OpName: "GetObjectTagging",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &GetObjectTaggingRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetObjectTagging",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key")

	request = &GetObjectTaggingRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "GetObjectTagging",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &GetObjectTaggingRequest{
		Bucket:    Ptr("oss-demo"),
		Key:       Ptr("oss-object"),
		VersionId: Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
	}
	input = &OperationInput{
		OpName: "GetObjectTagging",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["versionId"], "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")

	request = &GetObjectTaggingRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		VersionId:    Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "GetObjectTagging",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["versionId"], "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_GetObjectTagging(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<Tagging>
  <TagSet>
    <Tag>
      <Key>a</Key>
      <Value>1</Value>
    </Tag>
    <Tag>
      <Key>b</Key>
      <Value>2</Value>
    </Tag>
  </TagSet>
</Tagging>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
			"Date":             {"Mon, 18 Mar 2019 08:25:17 GMT"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &GetObjectTaggingResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Date"), "Mon, 18 Mar 2019 08:25:17 GMT")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	assert.Len(t, result.Tags, 2)
	assert.Equal(t, *result.Tags[0].Key, "a")
	assert.Equal(t, *result.Tags[0].Value, "1")
	assert.Equal(t, *result.Tags[1].Key, "b")
	assert.Equal(t, *result.Tags[1].Value, "2")
	body = `<?xml version="1.0" encoding="UTF-8"?>
<Tagging>
  <TagSet>
    <Tag>
      <Key>age</Key>
      <Value>18</Value>
    </Tag>
  </TagSet>
</Tagging>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
			"x-oss-version-id": {"CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &GetObjectTaggingResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Len(t, result.Tags, 1)
	assert.Equal(t, *result.Tags[0].Key, "age")
	assert.Equal(t, *result.Tags[0].Value, "18")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteObjectTagging(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteObjectTaggingRequest
	var input *OperationInput
	var err error

	request = &DeleteObjectTaggingRequest{}
	input = &OperationInput{
		OpName: "DeleteObjectTagging",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &DeleteObjectTaggingRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteObjectTagging",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key")

	request = &DeleteObjectTaggingRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-object"),
	}
	input = &OperationInput{
		OpName: "DeleteObjectTagging",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &DeleteObjectTaggingRequest{
		Bucket:    Ptr("oss-demo"),
		Key:       Ptr("oss-object"),
		VersionId: Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
	}
	input = &OperationInput{
		OpName: "DeleteObjectTagging",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["versionId"], "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")

	request = &DeleteObjectTaggingRequest{
		Bucket:       Ptr("oss-demo"),
		Key:          Ptr("oss-object"),
		VersionId:    Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "DeleteObjectTagging",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"tagging": "",
		},
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["versionId"], "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_DeleteObjectTagging(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Date":             {"Mon, 18 Mar 2019 08:25:17 GMT"},
		},
	}
	result := &DeleteObjectTaggingResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Date"), "Mon, 18 Mar 2019 08:25:17 GMT")
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
			"x-oss-version-id": {"CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"},
		},
	}
	result = &DeleteObjectTaggingResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_PutObject_ContentType(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var err error
	var input *OperationInput

	//PutObjectRequest
	var putRequest *PutObjectRequest
	putRequest = &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}

	// No auto detect content-type
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: putRequest.Bucket,
		Key:    putRequest.Key,
	}
	err = c.marshalInput(putRequest, input)
	assert.Nil(t, err)
	assert.Empty(t, input.Headers[HTTPHeaderContentType])

	// auto detect content-type, not match
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: putRequest.Bucket,
		Key:    putRequest.Key,
	}
	err = c.marshalInput(putRequest, input, updateContentType)
	assert.Nil(t, err)
	assert.Equal(t, "application/octet-stream", input.Headers[HTTPHeaderContentType])

	// auto detect content-type, match
	putRequest = &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key.txt"),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: putRequest.Bucket,
		Key:    putRequest.Key,
	}
	err = c.marshalInput(putRequest, input, updateContentType)
	assert.Nil(t, err)
	assert.Equal(t, "text/plain", input.Headers[HTTPHeaderContentType])

	//auto detect content-type + set by user
	putRequest = &PutObjectRequest{
		Bucket:      Ptr("bucket"),
		Key:         Ptr("key.txt"),
		ContentType: Ptr("set-by-user"),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: putRequest.Bucket,
		Key:    putRequest.Key,
	}
	err = c.marshalInput(putRequest, input, updateContentType)
	assert.Nil(t, err)
	assert.Equal(t, "set-by-user", input.Headers[HTTPHeaderContentType])
}

func TestMarshalInput_AppendObject_ContentType(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var err error
	var input *OperationInput
	var request *AppendObjectRequest

	//AppendObjectRequest
	request = &AppendObjectRequest{
		Bucket:   Ptr("bucket"),
		Key:      Ptr("key"),
		Position: Ptr(int64(0)),
	}

	// No auto detect content-type
	input = &OperationInput{
		OpName: "AppendObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Empty(t, input.Headers[HTTPHeaderContentType])

	// auto detect content-type, not match
	input = &OperationInput{
		OpName: "AppendObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentType)
	assert.Nil(t, err)
	assert.Equal(t, "application/octet-stream", input.Headers[HTTPHeaderContentType])

	// auto detect content-type, match
	request = &AppendObjectRequest{
		Bucket:   Ptr("bucket"),
		Key:      Ptr("key.txt"),
		Position: Ptr(int64(0)),
	}
	input = &OperationInput{
		OpName: "AppendObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentType)
	assert.Nil(t, err)
	assert.Equal(t, "text/plain", input.Headers[HTTPHeaderContentType])

	//auto detect content-type + set by user
	request = &AppendObjectRequest{
		Bucket:      Ptr("bucket"),
		Key:         Ptr("key.txt"),
		ContentType: Ptr("set-by-user"),
		Position:    Ptr(int64(0)),
	}
	input = &OperationInput{
		OpName: "AppendObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentType)
	assert.Nil(t, err)
	assert.Equal(t, "set-by-user", input.Headers[HTTPHeaderContentType])
}

func TestMarshalInput_InitiateMultipartUploadRequest_ContentType(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var err error
	var input *OperationInput
	var request *InitiateMultipartUploadRequest

	//InitiateMultipartUploadRequest
	request = &InitiateMultipartUploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}

	// No auto detect content-type
	input = &OperationInput{
		OpName: "InitiateMultipartUpload",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Empty(t, input.Headers[HTTPHeaderContentType])

	// auto detect content-type, not match
	input = &OperationInput{
		OpName: "InitiateMultipartUpload",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentType)
	assert.Nil(t, err)
	assert.Equal(t, "application/octet-stream", input.Headers[HTTPHeaderContentType])

	// auto detect content-type, match
	request = &InitiateMultipartUploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key.txt"),
	}
	input = &OperationInput{
		OpName: "InitiateMultipartUpload",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentType)
	assert.Nil(t, err)
	assert.Equal(t, "text/plain", input.Headers[HTTPHeaderContentType])

	//auto detect content-type + set by user
	request = &InitiateMultipartUploadRequest{
		Bucket:      Ptr("bucket"),
		Key:         Ptr("key.txt"),
		ContentType: Ptr("set-by-user"),
	}
	input = &OperationInput{
		OpName: "InitiateMultipartUpload",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentType)
	assert.Nil(t, err)
	assert.Equal(t, "set-by-user", input.Headers[HTTPHeaderContentType])
}

func TestMarshalInput_ProcessObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ProcessObjectRequest
	var input *OperationInput
	var err error

	request = &ProcessObjectRequest{}
	input = &OperationInput{
		OpName: "ProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ProcessObjectRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "ProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ProcessObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "ProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	destObjName := "dest.jpg"
	process := fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte(destObjName)))
	request = &ProcessObjectRequest{
		Bucket:  Ptr("oss-bucket"),
		Key:     Ptr("oss-key"),
		Process: Ptr(process),
	}
	input = &OperationInput{
		OpName: "ProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.NotEmpty(t, input.Body)
	assert.Empty(t, input.Parameters["x-oss-process"])
	data, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(data), "x-oss-process=image/resize,w_100|sys/saveas,o_ZGVzdC5qcGc=")

	request = &ProcessObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		Process:      Ptr(process),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "ProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.NotEmpty(t, input.Body)
	assert.Empty(t, input.Parameters["x-oss-process"])
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(data), "x-oss-process=image/resize,w_100|sys/saveas,o_ZGVzdC5qcGc=")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_ProcessObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5C06A3B67B8B5A3DA422****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":     {"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(`{
    "bucket": "",
    "fileSize": 3267,
    "object": "dest.jpg",
    "status": "OK"}`)),
	}
	result := &ProcessObjectResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyDefault, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.Bucket, "")
	assert.Equal(t, result.FileSize, 3267)
	assert.Equal(t, result.Object, "dest.jpg")
	assert.Equal(t, result.ProcessStatus, "OK")

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5C06A3B67B8B5A3DA422****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":     {"application/json"},
			"x-oss-version-id": {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
		},
		Body: io.NopCloser(strings.NewReader(`{
    "bucket": "desct-bucket",
    "fileSize": 3267,
    "object": "dest.jpg",
    "status": "OK"}`)),
	}
	result = &ProcessObjectResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyDefault, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, result.Bucket, "desct-bucket")
	assert.Equal(t, result.FileSize, 3267)
	assert.Equal(t, result.Object, "dest.jpg")
	assert.Equal(t, result.ProcessStatus, "OK")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyDefault, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyDefault, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 203,
		Status:     "Non-Authoritative Information",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyDefault, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 203)
	assert.Equal(t, result.Status, "Non-Authoritative Information")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_AsyncProcessObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *AsyncProcessObjectRequest
	var input *OperationInput
	var err error

	request = &AsyncProcessObjectRequest{}
	input = &OperationInput{
		OpName: "AsyncProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-async-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AsyncProcessObjectRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "AsyncProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-async-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AsyncProcessObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "AsyncProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-async-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	process := fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))
	request = &AsyncProcessObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		AsyncProcess: Ptr(process),
	}
	input = &OperationInput{
		OpName: "AsyncProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-async-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.NotEmpty(t, input.Body)
	assert.Empty(t, input.Parameters["x-oss-async-process"])
	data, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(data), "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ")

	request = &AsyncProcessObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		AsyncProcess: Ptr(process),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "AsyncProcessObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Parameters: map[string]string{
			"x-oss-async-process": "",
		},
	}
	err = c.marshalInput(request, input, addProcess, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.NotEmpty(t, input.Body)
	assert.Empty(t, input.Parameters["x-oss-async-process"])
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(data), "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")

}

func TestUnmarshalOutput_AsyncProcessObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5C06A3B67B8B5A3DA422****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":     {"application/json;charset=utf-8"},
		},
		Body: io.NopCloser(strings.NewReader(`{"EventId":"181-1kZUlN60OH4fWOcOjZEnGnG****","RequestId":"1D99637F-F59E-5B41-9200-C4892F52****","TaskId":"MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****"}`)),
	}
	result := &AsyncProcessObjectResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyDefault, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.EventId, "181-1kZUlN60OH4fWOcOjZEnGnG****")
	assert.Equal(t, result.RequestId, "1D99637F-F59E-5B41-9200-C4892F52****")
	assert.Equal(t, result.TaskId, "MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyDefault, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyDefault, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 203,
		Status:     "Non-Authoritative Information",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyDefault, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 203)
	assert.Equal(t, result.Status, "Non-Authoritative Information")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetObjectRequest_Process(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetObjectRequest
	var input *OperationInput
	var err error

	request = &GetObjectRequest{
		Bucket:  Ptr("oss-bucket"),
		Key:     Ptr("oss-key"),
		Process: Ptr("image/resize,m_fixed,w_100,h_100"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Parameters["x-oss-process"], "image/resize,m_fixed,w_100,h_100")
}

func TestMarshalInput_CleanRestoredObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CleanRestoredObjectRequest
	var input *OperationInput
	var err error

	request = &CleanRestoredObjectRequest{}
	input = &OperationInput{
		OpName: "CleanRestoredObject",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"cleanRestoredObject": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &CleanRestoredObjectRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "CleanRestoredObject",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"cleanRestoredObject": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Key")

	request = &CleanRestoredObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "CleanRestoredObject",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"cleanRestoredObject": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")

	request = &CleanRestoredObjectRequest{
		Bucket:       Ptr("oss-bucket"),
		Key:          Ptr("oss-key"),
		VersionId:    Ptr("version-id"),
		RequestPayer: Ptr("requester"),
	}
	input = &OperationInput{
		OpName: "CleanRestoredObject",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"cleanRestoredObject": "",
		},
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Parameters["versionId"], "version-id")
	assert.Equal(t, input.Headers["x-oss-request-payer"], "requester")
}

func TestUnmarshalOutput_CleanRestoredObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5C06A3B67B8B5A3DA422****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
		},
	}
	result := &CleanRestoredObjectResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>ArchiveRestoreNotFinished</Code>
  <Message>The archive file's restore is not finished.</Message>
  <RequestId>672C880CDF727138392C****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0016-00000719</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000719</RecommendDoc>
</Error>`)),
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
