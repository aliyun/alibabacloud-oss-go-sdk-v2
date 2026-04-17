package oss

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func setupGetObjectToFileV2MockServer(data []byte, crcHeader string, halfBodyOnce *bool, mu *sync.Mutex) *httptest.Server {
	length := len(data)
	gmtTime := getNowGMT()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}

		var httpRange *HTTPRange
		if r.Header.Get("Range") != "" {
			httpRange, _ = ParseRange(r.Header.Get("Range"))
		}

		offset := int64(0)
		statusCode := 200
		sendLen := int64(length)
		if httpRange != nil {
			offset = httpRange.Offset
			sendLen = int64(length) - httpRange.Offset
			if httpRange.Count > 0 {
				sendLen = minInt64(httpRange.Count, sendLen)
			}
			cr := httpContentRange{
				Offset: httpRange.Offset,
				Count:  sendLen,
				Total:  int64(length),
			}
			w.Header().Set("Content-Range", ToString(cr.FormatHTTPContentRange()))
			statusCode = 206
		}

		w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(sendLen))
		w.Header().Set(HTTPHeaderLastModified, gmtTime)
		w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
		w.Header().Set(HTTPHeaderContentType, "text/plain")
		if crcHeader != "" {
			w.Header().Set(HeaderOssCRC64, crcHeader)
		}

		w.WriteHeader(statusCode)

		sendData := data[int(offset):int(offset+sendLen)]
		if mu != nil && halfBodyOnce != nil {
			mu.Lock()
			doHalf := *halfBodyOnce
			if doHalf {
				*halfBodyOnce = false
			}
			mu.Unlock()
			if doHalf {
				sendData = sendData[:len(sendData)/2]
			}
		}
		w.Write(sendData)
	}))
}

func newMockClientForV2(serverURL string) *Client {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(serverURL).
		WithReadWriteTimeout(300 * time.Second)
	return NewClient(cfg)
}

func TestMockGetObjectToFileV2_Basic(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	crcVal := fmt.Sprint(hash.Sum64())

	server := setupGetObjectToFileV2MockServer(data, crcVal, nil, nil)
	defer server.Close()

	client := newMockClientForV2(server.URL)

	localFile := randStr(8) + "-v2-basic"
	defer os.Remove(localFile)

	result, err := client.GetObjectToFileV2(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		localFile,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	downloaded, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, data, downloaded)
}

func TestMockGetObjectToFileV2_Progress(t *testing.T) {
	length := 1024*1024 + 567
	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	crcVal := fmt.Sprint(hash.Sum64())

	server := setupGetObjectToFileV2MockServer(data, crcVal, nil, nil)
	defer server.Close()

	client := newMockClientForV2(server.URL)

	localFile := randStr(8) + "-v2-progress"
	defer os.Remove(localFile)

	var lastTransferred, lastTotal int64
	result, err := client.GetObjectToFileV2(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			ProgressFn: func(increment, transferred, total int64) {
				lastTransferred = transferred
				lastTotal = total
			},
		},
		localFile,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(length), lastTransferred)
	assert.Equal(t, int64(length), lastTotal)
}

func TestMockGetObjectToFileV2_CRC64Check(t *testing.T) {
	length := 512*1024 + 789
	data := []byte(randStr(length))

	// Wrong CRC header
	server := setupGetObjectToFileV2MockServer(data, "12345", nil, nil)
	defer server.Close()

	client := newMockClientForV2(server.URL)

	localFile := randStr(8) + "-v2-crc"
	defer os.Remove(localFile)

	// CRC check enabled by default — should fail
	_, err := client.GetObjectToFileV2(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		localFile,
		nil,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "crc is inconsistent")

	// Disable CRC check
	client.options.FeatureFlags = client.options.FeatureFlags & ^FeatureEnableCRC64CheckDownload
	result, err := client.GetObjectToFileV2(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		localFile,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	downloaded, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, data, downloaded)
}

func TestMockGetObjectToFileV2_CRC64SkipRange(t *testing.T) {
	length := 1024 * 100
	data := []byte(randStr(length))

	// Wrong CRC header — but ranged request should skip CRC check
	server := setupGetObjectToFileV2MockServer(data, "12345", nil, nil)
	defer server.Close()

	client := newMockClientForV2(server.URL)

	localFile := randStr(8) + "-v2-range"
	defer os.Remove(localFile)

	rangeStr := "bytes=0-1023"
	result, err := client.GetObjectToFileV2(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			Range:  Ptr(rangeStr),
		},
		localFile,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	downloaded, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, data[:1024], downloaded)
}

func TestMockGetObjectToFileV2_WriteBufferSize(t *testing.T) {
	length := 2*1024*1024 + 333
	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	crcVal := fmt.Sprint(hash.Sum64())

	server := setupGetObjectToFileV2MockServer(data, crcVal, nil, nil)
	defer server.Close()

	client := newMockClientForV2(server.URL)

	localFile := randStr(8) + "-v2-wbuf"
	defer os.Remove(localFile)

	bufSize := 64 * 1024
	result, err := client.GetObjectToFileV2(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		localFile,
		&bufSize,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	downloaded, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, data, downloaded)
}

func TestMockGetObjectToFileV2_Resume(t *testing.T) {
	length := 1024*1024 + 999
	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	crcVal := fmt.Sprint(hash.Sum64())

	halfBody := true
	var mu sync.Mutex

	server := setupGetObjectToFileV2MockServer(data, crcVal, &halfBody, &mu)
	defer server.Close()

	client := newMockClientForV2(server.URL)

	localFile := randStr(8) + "-v2-resume"
	defer os.Remove(localFile)

	result, err := client.GetObjectToFileV2(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		localFile,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	downloaded, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, data, downloaded)

	// Verify CRC of downloaded file
	fileHash := NewCRC64(0)
	fileHash.Write(downloaded)
	assert.Equal(t, hash.Sum64(), fileHash.Sum64())
}

func TestMockGetObjectToFileV2_NilRequest(t *testing.T) {
	_, err := (&Client{}).GetObjectToFileV2(context.TODO(), nil, "test.txt", nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request")
}

func TestMockGetObjectToFileV2_ProgressWithResume(t *testing.T) {
	length := 1024*512 + 111
	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	crcVal := fmt.Sprint(hash.Sum64())

	halfBody := true
	var mu sync.Mutex

	server := setupGetObjectToFileV2MockServer(data, crcVal, &halfBody, &mu)
	defer server.Close()

	client := newMockClientForV2(server.URL)

	localFile := randStr(8) + "-v2-prog-resume"
	defer os.Remove(localFile)

	var lastTransferred int64
	var progressCalls int64
	result, err := client.GetObjectToFileV2(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			ProgressFn: func(increment, transferred, total int64) {
				lastTransferred = transferred
				progressCalls++
			},
		},
		localFile,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(length), lastTransferred)
	assert.True(t, progressCalls > 0)

	downloaded, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(data, downloaded))
}

func TestMockGetObjectToFileV2_CRC64WithResume(t *testing.T) {
	length := 512*1024 + 555
	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	crcVal := fmt.Sprint(hash.Sum64())

	halfBody := true
	var mu sync.Mutex

	server := setupGetObjectToFileV2MockServer(data, crcVal, &halfBody, &mu)
	defer server.Close()

	client := newMockClientForV2(server.URL)

	localFile := randStr(8) + "-v2-crc-resume"
	defer os.Remove(localFile)

	result, err := client.GetObjectToFileV2(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		localFile,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Verify file content
	downloaded, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, data, downloaded)

	// Verify CRC matches
	fileHash := NewCRC64(0)
	fileHash.Write(downloaded)
	assert.Equal(t, hash.Sum64(), fileHash.Sum64())
}
