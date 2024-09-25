package oss

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

type downloaderMockTracker struct {
	lastModified string
	data         []byte

	maxRangeCount int64
	getPartCnt    int32

	etagChangeOffset int64
	failPartNum      int32

	partSize     int32
	gotMinOffset int64

	rStart            int64
	headReqeustCRCErr bool

	halfBodyErr bool

	mu sync.Mutex
}

func testSetupDownloaderMockServer(_ *testing.T, tracker *downloaderMockTracker) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		length := len(tracker.data)
		data := tracker.data
		errData := []byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>InvalidAccessKeyId</Code>
				<Message>The OSS Access Key Id you provided does not exist in our records.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>ak</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`)

		switch r.Method {
		case "HEAD":
			tracker.gotMinOffset = int64(length)
			hash := NewCRC64(0)
			hash.Write(data)
			// header
			w.Header().Set(HTTPHeaderLastModified, tracker.lastModified)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			if tracker.headReqeustCRCErr {
				w.Header().Set(HeaderOssCRC64, "12345")
			} else {
				w.Header().Set(HeaderOssCRC64, fmt.Sprint(hash.Sum64()))
			}

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "GET":
			// header
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
					tracker.mu.Lock()
					tracker.maxRangeCount = maxInt64(httpRange.Count, tracker.maxRangeCount)
					tracker.mu.Unlock()
				}
				cr := httpContentRange{
					Offset: httpRange.Offset,
					Count:  sendLen,
					Total:  int64(length),
				}
				w.Header().Set("Content-Range", ToString(cr.FormatHTTPContentRange()))
				statusCode = 206
			}

			tracker.mu.Lock()
			tracker.gotMinOffset = minInt64(tracker.gotMinOffset, offset)
			tracker.mu.Unlock()

			if tracker.failPartNum > 0 && (int64(tracker.partSize*tracker.failPartNum)+tracker.rStart) == offset {
				w.Header().Set(HTTPHeaderContentType, "application/xml")
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
				w.WriteHeader(403)
				w.Write(errData)
			} else {
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(sendLen))
				w.Header().Set(HTTPHeaderLastModified, tracker.lastModified)
				if tracker.etagChangeOffset > 0 && offset > 0 && offset > tracker.etagChangeOffset {
					w.Header().Set(HTTPHeaderETag, "2ba9dede5f27731c9771645a3986****")
				} else {
					w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
				}
				w.Header().Set(HTTPHeaderContentType, "text/plain")

				//status code
				w.WriteHeader(statusCode)

				//body
				sendData := data[int(offset):int(offset+sendLen)]
				if tracker.halfBodyErr {
					sendData = sendData[0 : len(sendData)/2]
					tracker.halfBodyErr = false
				}
				//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)
				w.Write(sendData)
			}
		}
	}))
	return server
}

func TestMockDownloaderSingleRead(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client, func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile)
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)
	defer func() {
		rfile.Close()
		os.Remove(localFile)
	}()
	io.Copy(hash, rfile)
	assert.Equal(t, datasum, hash.Sum64())
}

func TestMockDownloaderLoopSingleRead(t *testing.T) {
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	for i := 1; i <= 20; i++ {
		if FileExists(localFile) {
			assert.Nil(t, os.Remove(localFile))
		}
		tracker.maxRangeCount = 0
		result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
			func(do *DownloaderOptions) {
				do.ParallelNum = 1
				do.PartSize = int64(i)
			})
		assert.Nil(t, err)
		assert.Equal(t, int64(length), result.Written)
		hash := NewCRC64(0)
		rfile, err := os.Open(localFile)
		assert.Nil(t, err)
		io.Copy(hash, rfile)
		rfile.Close()
		os.Remove(localFile)
		assert.Equal(t, datasum, hash.Sum64())
		assert.Equal(t, int64(i), tracker.maxRangeCount)
	}
}

func TestMockDownloaderLoopSingleReadWithRange(t *testing.T) {
	length := 63
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	for rs := 0; rs < 7; rs++ {
		for rcount := 1; rcount < length; rcount++ {
			for i := 1; i <= 3; i++ {
				//fmt.Printf("rs:%v, rcount:%v, i:%v\n", rs, rcount, i)
				if FileExists(localFile) {
					assert.Nil(t, os.Remove(localFile))
				}
				tracker.maxRangeCount = 0
				httpRange := HTTPRange{Offset: int64(rs), Count: int64(rcount)}
				result, err := d.DownloadFile(context.TODO(),
					&GetObjectRequest{
						Bucket: Ptr("bucket"),
						Key:    Ptr("key"),
						Range:  httpRange.FormatHTTPRange()},
					localFile,
					func(do *DownloaderOptions) {
						do.ParallelNum = 1
						do.PartSize = int64(i)
					})
				assert.Nil(t, err)
				expectLen := minInt64(int64(length-rs), int64(rcount))
				assert.Equal(t, expectLen, result.Written)
				hash := NewCRC64(0)
				rfile, err := os.Open(localFile)
				assert.Nil(t, err)
				io.Copy(hash, rfile)
				rfile.Close()
				//ldata, err := os.ReadFile(localFile)
				//assert.Nil(t, err)
				os.Remove(localFile)
				hdata := NewCRC64(0)
				pat := data[rs:int(minInt64(int64(rs+rcount), int64(length)))]
				hdata.Write(pat)
				//assert.EqualValues(t, ldata, pat)
				assert.Equal(t, hdata.Sum64(), hash.Sum64())
				assert.Equal(t, minInt64(int64(i), expectLen), tracker.maxRangeCount)
			}
		}
	}
}

func TestMockDownloaderParalleRead(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client, func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 3, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile)
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)
	defer func() {
		rfile.Close()
		os.Remove(localFile)
	}()

	io.Copy(hash, rfile)
	assert.Equal(t, datasum, hash.Sum64())
}

func TestMockDownloaderLoopParalleRead(t *testing.T) {
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	for i := 1; i <= 20; i++ {
		if FileExists(localFile) {
			assert.Nil(t, os.Remove(localFile))
		}
		tracker.maxRangeCount = 0
		result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
			func(do *DownloaderOptions) {
				do.ParallelNum = 4
				do.PartSize = int64(i)
			})
		assert.Nil(t, err)
		assert.Equal(t, int64(length), result.Written)
		hash := NewCRC64(0)
		rfile, err := os.Open(localFile)
		assert.Nil(t, err)
		io.Copy(hash, rfile)
		assert.Nil(t, rfile.Close())
		assert.Nil(t, os.Remove(localFile))
		assert.Equal(t, datasum, hash.Sum64())
		assert.Equal(t, int64(i), tracker.maxRangeCount)
	}
}

func TestMockDownloaderLoopParalleReadWithRange(t *testing.T) {
	length := 63
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	for rs := 0; rs < 7; rs++ {
		for rcount := 1; rcount < length; rcount++ {
			for i := 1; i <= 3; i++ {
				//fmt.Printf("rs:%v, rcount:%v, i:%v\n", rs, rcount, i)
				if FileExists(localFile) {
					assert.Nil(t, os.Remove(localFile))
				}
				tracker.maxRangeCount = 0
				httpRange := HTTPRange{Offset: int64(rs), Count: int64(rcount)}
				result, err := d.DownloadFile(context.TODO(),
					&GetObjectRequest{
						Bucket: Ptr("bucket"),
						Key:    Ptr("key"),
						Range:  httpRange.FormatHTTPRange()},
					localFile,
					func(do *DownloaderOptions) {
						do.ParallelNum = 3
						do.PartSize = int64(i)
					})
				assert.Nil(t, err)
				expectLen := minInt64(int64(length-rs), int64(rcount))
				assert.Equal(t, expectLen, result.Written)
				hash := NewCRC64(0)
				rfile, err := os.Open(localFile)
				assert.Nil(t, err)
				io.Copy(hash, rfile)
				rfile.Close()
				//ldata, err := os.ReadFile(localFile)
				//assert.Nil(t, err)
				os.Remove(localFile)
				hdata := NewCRC64(0)
				pat := data[rs:int(minInt64(int64(rs+rcount), int64(length)))]
				hdata.Write(pat)
				//assert.EqualValues(t, ldata, pat)
				assert.Equal(t, hdata.Sum64(), hash.Sum64())
				assert.Equal(t, minInt64(int64(i), expectLen), tracker.maxRangeCount)
			}
		}
	}
}

func TestDownloaderConstruct(t *testing.T) {
	c := &Client{}
	d := NewDownloader(c)
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.True(t, d.options.UseTempFile)
	assert.False(t, d.options.EnableCheckpoint)
	assert.Equal(t, "", d.options.CheckpointDir)

	d = NewDownloader(c, func(do *DownloaderOptions) {
		do.CheckpointDir = "dir"
		do.EnableCheckpoint = true
		do.ParallelNum = 1
		do.PartSize = 2
		do.UseTempFile = false
	})
	assert.Equal(t, 1, d.options.ParallelNum)
	assert.Equal(t, int64(2), d.options.PartSize)
	assert.False(t, d.options.UseTempFile)
	assert.True(t, d.options.EnableCheckpoint)
	assert.Equal(t, "dir", d.options.CheckpointDir)
}

func TestDownloaderDelegateConstruct(t *testing.T) {
	c := &Client{}
	d := NewDownloader(c)

	_, err := d.newDelegate(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")

	_, err = d.newDelegate(context.TODO(), &GetObjectRequest{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request.Bucket")

	_, err = d.newDelegate(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket")})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request.Key")

	delegate, err := d.newDelegate(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")})
	assert.Nil(t, err)
	assert.NotNil(t, delegate)
	assert.Equal(t, DefaultDownloadParallel, delegate.options.ParallelNum)
	assert.Equal(t, DefaultDownloadPartSize, delegate.options.PartSize)
	assert.True(t, delegate.options.UseTempFile)
	assert.False(t, delegate.options.EnableCheckpoint)
	assert.Empty(t, delegate.options.CheckpointDir)

	delegate, err = d.newDelegate(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")},
		func(do *DownloaderOptions) {
			do.ParallelNum = 5
			do.PartSize = 1
		})
	assert.Nil(t, err)
	assert.NotNil(t, delegate)
	assert.Equal(t, 5, delegate.options.ParallelNum)
	assert.Equal(t, int64(1), delegate.options.PartSize)

	delegate, err = d.newDelegate(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")},
		func(do *DownloaderOptions) {
			do.ParallelNum = 0
			do.PartSize = 0
		})
	assert.Nil(t, err)
	assert.NotNil(t, delegate)
	assert.Equal(t, DefaultDownloadParallel, delegate.options.ParallelNum)
	assert.Equal(t, DefaultDownloadPartSize, delegate.options.PartSize)

	delegate, err = d.newDelegate(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")},
		func(do *DownloaderOptions) {
			do.ParallelNum = -1
			do.PartSize = -1
			do.CheckpointDir = "dir"
			do.EnableCheckpoint = true
			do.UseTempFile = false
		})
	assert.Nil(t, err)
	assert.NotNil(t, delegate)
	assert.Equal(t, DefaultDownloadParallel, delegate.options.ParallelNum)
	assert.Equal(t, DefaultDownloadPartSize, delegate.options.PartSize)
	assert.False(t, delegate.options.UseTempFile)
	assert.True(t, delegate.options.EnableCheckpoint)
	assert.Equal(t, "dir", delegate.options.CheckpointDir)
}

func TestDownloaderDownloadFileArgument(t *testing.T) {
	c := &Client{}
	d := NewDownloader(c)

	_, err := d.DownloadFile(context.TODO(), nil, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")

	_, err = d.DownloadFile(context.TODO(), &GetObjectRequest{}, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request.Bucket")

	_, err = d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket")}, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request.Key")

	_, err = d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "operation error HeadObject")
}

func TestMockDownloaderDownloadFileWithoutTempFile(t *testing.T) {
	length := 3*int(DefaultDownloadPartSize) + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client, func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.UseTempFile = false
			do.PartSize = 1024
			do.ParallelNum = 2
		})
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	io.Copy(hash, rfile)
	defer func() {
		rfile.Close()
		os.Remove(localFile)
	}()
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, int64(1024), tracker.maxRangeCount)
}

func TestMockDownloaderDownloadFileInvalidPartSizeAndParallelNum(t *testing.T) {
	length := int(DefaultDownloadPartSize*2) + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client, func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = 0
			do.ParallelNum = 0
		})
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	io.Copy(hash, rfile)
	rfile.Close()
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, DefaultDownloadPartSize, tracker.maxRangeCount)
}

func TestMockDownloaderDownloadFileFileSizeLessPartSize(t *testing.T) {
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client, func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = 0
			do.ParallelNum = 0
		})
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	io.Copy(hash, rfile)
	rfile.Close()
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, int64(length), tracker.maxRangeCount)
}

func TestMockDownloaderDownloadFileFileChange(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()

	tracker := &downloaderMockTracker{
		lastModified:     gmtTime,
		data:             data,
		etagChangeOffset: 700,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	_, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
		})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Source file is changed")
	assert.False(t, FileExists(localFile))
	assert.False(t, FileExists(localFile+TempFileSuffix))
}

func TestMockDownloaderDownloadFileEnableCheckpointNormal(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()

	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	_, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})
	assert.Nil(t, err)
}

func TestMockDownloaderDownloadFileEnableCheckpoint2(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
		failPartNum:  6,
		partSize:     int32(partSize),
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := "check-point-to-check-no-surfix"
	localFileTmep := "check-point-to-check-no-surfix" + TempFileSuffix
	absPath, _ := filepath.Abs(localFileTmep)
	hashmd5 := md5.New()
	hashmd5.Reset()
	hashmd5.Write([]byte(absPath))
	destHash := hex.EncodeToString(hashmd5.Sum(nil))
	cpFileName := "ddaf063c8f69766ecc8e4a93b6402e3e-" + destHash + ".dcp"
	defer func() {
		os.Remove(localFile)
	}()
	os.Remove(localFile)
	os.Remove(localFileTmep)
	os.Remove(cpFileName)
	_, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})

	assert.NotNil(t, err)
	assert.True(t, FileExists(localFileTmep))
	assert.True(t, FileExists(cpFileName))

	//load CheckPointFile
	content, err := os.ReadFile(cpFileName)
	assert.Nil(t, err)
	dcp := downloadCheckpoint{}
	err = json.Unmarshal(content, &dcp.Info)
	assert.Nil(t, err)

	assert.Equal(t, "fba9dede5f27731c9771645a3986****", dcp.Info.Data.ObjectMeta.ETag)
	assert.Equal(t, gmtTime, dcp.Info.Data.ObjectMeta.LastModified)
	assert.Equal(t, int64(length), dcp.Info.Data.ObjectMeta.Size)

	assert.Equal(t, "oss://bucket/key", dcp.Info.Data.ObjectInfo.Name)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.VersionId)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.Range)

	abslocalFileTmep, _ := filepath.Abs(localFileTmep)
	assert.Equal(t, abslocalFileTmep, dcp.Info.Data.FilePath)
	assert.Equal(t, int64(partSize), dcp.Info.Data.PartSize)

	assert.Equal(t, int64(tracker.failPartNum*tracker.partSize), dcp.Info.Data.DownloadInfo.Offset)
	h := NewCRC64(0)
	h.Write(data[0:int(dcp.Info.Data.DownloadInfo.Offset)])
	assert.Equal(t, h.Sum64(), dcp.Info.Data.DownloadInfo.CRC64)

	// resume from checkpoint
	tracker.failPartNum = 0
	result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
			do.VerifyData = true
		})

	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	io.Copy(hash, rfile)
	rfile.Close()
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, dcp.Info.Data.DownloadInfo.Offset, tracker.gotMinOffset)
}

func TestMockDownloaderDownloadFileEnableCheckpointWithRange(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	rs := 5
	rcount := 832
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
		failPartNum:  6,
		partSize:     int32(partSize),
		rStart:       int64(rs),
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := "check-point-to-check-no-surfix"
	localFileTmep := "check-point-to-check-no-surfix" + TempFileSuffix
	absPath, _ := filepath.Abs(localFileTmep)
	hashmd5 := md5.New()
	hashmd5.Reset()
	hashmd5.Write([]byte(absPath))
	destHash := hex.EncodeToString(hashmd5.Sum(nil))
	cpFileName := "0fbbf3bb7c80debbecb37dca52a646eb-" + destHash + ".dcp"
	defer func() {
		os.Remove(localFile)
	}()
	os.Remove(localFile)
	os.Remove(localFileTmep)
	os.Remove(cpFileName)
	httpRange := HTTPRange{Offset: int64(rs), Count: int64(rcount)}
	_, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			Range:  httpRange.FormatHTTPRange()},
		localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})

	assert.NotNil(t, err)
	assert.True(t, FileExists(localFileTmep))
	assert.True(t, FileExists(cpFileName))

	//load CheckPointFile
	content, err := os.ReadFile(cpFileName)
	assert.Nil(t, err)
	dcp := downloadCheckpoint{}
	err = json.Unmarshal(content, &dcp.Info)
	assert.Nil(t, err)

	assert.Equal(t, "fba9dede5f27731c9771645a3986****", dcp.Info.Data.ObjectMeta.ETag)
	assert.Equal(t, gmtTime, dcp.Info.Data.ObjectMeta.LastModified)
	assert.Equal(t, int64(length), dcp.Info.Data.ObjectMeta.Size)

	assert.Equal(t, "oss://bucket/key", dcp.Info.Data.ObjectInfo.Name)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.VersionId)
	assert.Equal(t, ToString(httpRange.FormatHTTPRange()), dcp.Info.Data.ObjectInfo.Range)

	abslocalFileTmep, _ := filepath.Abs(localFileTmep)
	assert.Equal(t, abslocalFileTmep, dcp.Info.Data.FilePath)
	assert.Equal(t, int64(partSize), dcp.Info.Data.PartSize)

	assert.Equal(t, int64(tracker.failPartNum*tracker.partSize)+int64(rs), dcp.Info.Data.DownloadInfo.Offset)
	assert.Equal(t, uint64(0), dcp.Info.Data.DownloadInfo.CRC64)

	// resume from checkpoint
	tracker.failPartNum = 0
	result, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			Range:  httpRange.FormatHTTPRange()},
		localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})

	assert.Nil(t, err)
	expectLen := minInt64(int64(length-rs), int64(rcount))
	assert.Equal(t, expectLen, result.Written)
	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)
	io.Copy(hash, rfile)
	rfile.Close()
	os.Remove(localFile)
	hdata := NewCRC64(0)
	pat := data[rs:int(minInt64(int64(rs+rcount), int64(length)))]
	hdata.Write(pat)
	assert.Equal(t, hdata.Sum64(), hash.Sum64())
	assert.Equal(t, dcp.Info.Data.DownloadInfo.Offset, tracker.gotMinOffset)
}

func TestMockDownloaderDownloadWithError(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client, func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	// filePath is invalid
	_, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid field, filePath")

	localFile := "./no-exist-folder/file-no-surfix"
	_, err = d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "The system cannot find the path specified") || strings.Contains(err.Error(), "no such file or directory"))

	// Range is invalid
	localFile = randStr(8) + "-no-surfix"
	_, err = d.DownloadFile(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Range:  Ptr("invalid range")},
		localFile)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid field, request.Range")
}

func TestDownloadedChunksSort(t *testing.T) {
	chunks := downloadedChunks{}
	chunks = append(chunks, downloadedChunk{start: 0, size: 10})
	chunks = append(chunks, downloadedChunk{start: 10, size: 5})
	chunks = append(chunks, downloadedChunk{start: 15, size: 30})
	chunks = append(chunks, downloadedChunk{start: 45, size: 1})
	sort.Sort(chunks)

	assert.Equal(t, 4, len(chunks))
	assert.Equal(t, int64(0), chunks[0].start)
	assert.Equal(t, int64(10), chunks[1].start)
	assert.Equal(t, int64(15), chunks[2].start)
	assert.Equal(t, int64(45), chunks[3].start)

	chunks = downloadedChunks{}
	chunks = append(chunks, downloadedChunk{start: 10, size: 5})
	chunks = append(chunks, downloadedChunk{start: 0, size: 10})
	chunks = append(chunks, downloadedChunk{start: 45, size: 1})
	chunks = append(chunks, downloadedChunk{start: 15, size: 30})

	assert.Equal(t, 4, len(chunks))
	assert.Equal(t, int64(10), chunks[0].start)
	assert.Equal(t, int64(0), chunks[1].start)
	assert.Equal(t, int64(45), chunks[2].start)
	assert.Equal(t, int64(15), chunks[3].start)

	sort.Sort(chunks)

	assert.Equal(t, int64(0), chunks[0].start)
	assert.Equal(t, int64(10), chunks[1].start)
	assert.Equal(t, int64(15), chunks[2].start)
	assert.Equal(t, int64(45), chunks[3].start)
}

func TestMockDownloaderCRCCheck(t *testing.T) {
	length := 5*100*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified:      gmtTime,
		data:              data,
		headReqeustCRCErr: true,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client, func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 100 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := randStr(8) + "-no-surfix"

	result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cause: crc is inconsistent")
	assert.NotEmpty(t, datasum)
	os.Remove(localFile)
	assert.NoFileExists(t, localFile)

	// Disable CRC
	d.featureFlags = d.featureFlags & ^FeatureEnableCRC64CheckDownload
	result, err = d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)

	defer func() {
		rfile.Close()
		os.Remove(localFile)
	}()

	io.Copy(hash, rfile)
	assert.Equal(t, datasum, hash.Sum64())
}

func TestMockDownloaderCRCCheckWithResume(t *testing.T) {
	length := 5*100*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
		halfBodyErr:  true,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client, func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 100 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := randStr(8) + "-no-surfix"

	result, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)

	defer func() {
		rfile.Close()
		os.Remove(localFile)
	}()

	io.Copy(hash, rfile)
	assert.Equal(t, datasum, hash.Sum64())
}

func TestMockDownloaderDownloadFileProcess(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()

	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	_, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key"), Process: Ptr("image/resize,m_fixed,w_100,h_100/rotate,90")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})
	assert.Nil(t, err)
}

func TestMockDownloaderDownloadFilePayer(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()

	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	_, err := d.DownloadFile(context.TODO(), &GetObjectRequest{Bucket: Ptr("bucket"), Key: Ptr("key"), RequestPayer: Ptr("requester")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})
	assert.Nil(t, err)
}

func TestMockDownloaderWithProgress(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)
	client := NewClient(cfg)
	var n int64
	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)
	// filePath is invalid
	_, err := d.DownloadFile(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			ProgressFn: func(increment, transferred, total int64) {
				n = transferred
			},
		}, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid field, filePath")
	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()
	_, err = d.DownloadFile(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			ProgressFn: func(increment, transferred, total int64) {
				n = transferred
				fmt.Printf("increment:%#v, transferred:%#v, total:%#v\n", increment, transferred, total)
			},
		}, localFile)
	assert.Nil(t, err)
	assert.Equal(t, n, int64(length))
	n = int64(0)
	d = client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 3 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(3*1024*1024), d.options.PartSize)
	assert.Equal(t, 3, d.options.ParallelNum)
	// filePath is invalid
	_, err = d.DownloadFile(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			ProgressFn: func(increment, transferred, total int64) {
				n = transferred
			},
		}, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid field, filePath")
	_, err = d.DownloadFile(
		context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			ProgressFn: func(increment, transferred, total int64) {
				n = transferred
			},
		}, localFile)
	assert.Nil(t, err)
	assert.Equal(t, n, int64(length))
}

func TestMockDownloaderDownloadFileEnableCheckpointProgress(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
		failPartNum:  6,
		partSize:     int32(partSize),
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := NewDownloader(client)
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := "check-point-to-check-no-surfix"
	localFileTmep := "check-point-to-check-no-surfix" + TempFileSuffix
	absPath, _ := filepath.Abs(localFileTmep)
	hashmd5 := md5.New()
	hashmd5.Reset()
	hashmd5.Write([]byte(absPath))
	destHash := hex.EncodeToString(hashmd5.Sum(nil))
	cpFileName := "ddaf063c8f69766ecc8e4a93b6402e3e-" + destHash + ".dcp"
	defer func() {
		os.Remove(localFile)
	}()
	os.Remove(localFile)
	os.Remove(localFileTmep)
	os.Remove(cpFileName)
	inc := int64(0)
	_, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			ProgressFn: func(increment, transferred, total int64) {
				inc += increment
			},
		},
		localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})

	assert.NotNil(t, err)
	assert.True(t, FileExists(localFileTmep))
	assert.True(t, FileExists(cpFileName))
	assert.Less(t, inc, int64(length))

	//load CheckPointFile
	content, err := os.ReadFile(cpFileName)
	assert.Nil(t, err)
	dcp := downloadCheckpoint{}
	err = json.Unmarshal(content, &dcp.Info)
	assert.Nil(t, err)

	assert.Equal(t, "fba9dede5f27731c9771645a3986****", dcp.Info.Data.ObjectMeta.ETag)
	assert.Equal(t, gmtTime, dcp.Info.Data.ObjectMeta.LastModified)
	assert.Equal(t, int64(length), dcp.Info.Data.ObjectMeta.Size)

	assert.Equal(t, "oss://bucket/key", dcp.Info.Data.ObjectInfo.Name)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.VersionId)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.Range)

	abslocalFileTmep, _ := filepath.Abs(localFileTmep)
	assert.Equal(t, abslocalFileTmep, dcp.Info.Data.FilePath)
	assert.Equal(t, int64(partSize), dcp.Info.Data.PartSize)

	assert.Equal(t, int64(tracker.failPartNum*tracker.partSize), dcp.Info.Data.DownloadInfo.Offset)
	h := NewCRC64(0)
	h.Write(data[0:int(dcp.Info.Data.DownloadInfo.Offset)])
	assert.Equal(t, h.Sum64(), dcp.Info.Data.DownloadInfo.CRC64)

	// resume from checkpoint
	tracker.failPartNum = 0
	inc = int64(0)
	result, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			ProgressFn: func(increment, transferred, total int64) {
				inc += increment
			},
		},
		localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
			do.VerifyData = true
		})

	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	io.Copy(hash, rfile)
	rfile.Close()
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, dcp.Info.Data.DownloadInfo.Offset, tracker.gotMinOffset)
	assert.Equal(t, int64(length), inc)
}
