package oss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

type httpContentRange struct {
	Offset int64
	Count  int64
	Total  int64
}

func (r httpContentRange) FormatHTTPContentRange() *string {
	if r.Offset == 0 && r.Count == 0 {
		return nil // No specified range
	}
	endOffset := "" // if count == CountToEnd (0)
	if r.Count > 0 {
		endOffset = strconv.FormatInt((r.Offset+r.Count)-1, 10)
	}
	dataRange := fmt.Sprintf("bytes %v-%s/%s", r.Offset, endOffset, strconv.FormatInt(r.Total, 10))
	return &dataRange
}

func TestMockOpenFile_DirectRead(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

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
					sendLen = httpRange.Count
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

			//status code
			w.WriteHeader(statusCode)

			//body
			sendData := data[int(offset):int(offset+sendLen)]
			//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)
			w.Write(sendData)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, false, f.enablePrefetch)
	assert.Equal(t, DefaultPrefetchChunkSize, f.chunkSize)
	assert.Equal(t, DefaultPrefetchNum, f.prefetchNum)
	assert.Equal(t, DefaultPrefetchThreshold, f.prefetchThreshold)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())
	assert.Equal(t, gmtTime, stat.ModTime().Format(http.TimeFormat))
	assert.Equal(t, os.FileMode(0644), stat.Mode())
	assert.Equal(t, false, stat.IsDir())
	assert.Equal(t, "oss://bucket/key", stat.Name())
	h, ok := stat.Sys().(http.Header)
	assert.True(t, ok)
	assert.NotNil(t, h)
	assert.Equal(t, "fba9dede5f27731c9771645a3986****", h.Get(HTTPHeaderETag))
	assert.Equal(t, "text/plain", h.Get(HTTPHeaderContentType))

	//seek ok
	begin, err := f.Seek(0, io.SeekStart)
	end, err := f.Seek(0, io.SeekEnd)
	assert.Equal(t, stat.Size(), end-begin)

	//seek invalid
	begin, err = f.Seek(0, 4)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid whence")
	assert.Equal(t, int64(0), begin)

	begin, err = f.Seek(-1, io.SeekStart)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "negative position")
	assert.Equal(t, int64(0), begin)

	begin, err = f.Seek(100, io.SeekEnd)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "offset is unavailable")
	assert.Equal(t, int64(0), begin)

	//read all
	f.Seek(0, io.SeekStart)
	hash := NewCRC64(0)
	written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, stat.Size(), written)
	//seek readN
	for i := 0; i < 64; i++ {
		offset := rand.Int63n(int64(length / 4))
		n := rand.Int63n(int64(length/2)) + 1
		//fmt.Printf("seek readN check offset%d, len:%d\n", offset, n)
		begin, err = f.Seek(offset, io.SeekStart)
		assert.Nil(t, err)
		assert.Equal(t, offset, begin)

		hash := NewCRC64(0)
		written, err := io.CopyN(io.MultiWriter(io.Discard, hash), f, n)
		assert.Nil(t, err)
		assert.Equal(t, n, written)

		hash1 := NewCRC64(0)
		hash1.Write(data[offset : offset+n])

		assert.Equal(t, hash1.Sum64(), hash.Sum64())
	}

	//seek read from offset to end
	for i := 0; i < 64; i++ {
		offset := rand.Int63n(int64(length / 5))
		begin, err = f.Seek(offset, io.SeekStart)
		n := int64(length) - offset
		//fmt.Printf("seek readAll check offset%d, len:%d\n", offset, n)
		assert.Nil(t, err)
		assert.Equal(t, offset, begin)

		hash := NewCRC64(0)
		written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
		assert.Nil(t, err)
		assert.Equal(t, n, written)

		hash1 := NewCRC64(0)
		hash1.Write(data[offset:])
		assert.Equal(t, hash1.Sum64(), hash.Sum64())
	}

	//seek to end
	begin, err = f.Seek(0, io.SeekEnd)
	assert.Nil(t, err)
	written, err = io.Copy(io.Discard, f)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), written)

	//
	begin, err = f.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	io.CopyN(io.Discard, f, 2)
	time.Sleep(2 * time.Second)

	err = f.Close()
	assert.Nil(t, err)

	//call Close many times
	err = f.Close()
	assert.Nil(t, err)

	_, err = f.Seek(0, io.SeekEnd)
	assert.Equal(t, err, os.ErrClosed)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrClosed)

	bytedata := make([]byte, 5)
	_, err = f.Read(bytedata)
	assert.Equal(t, err, os.ErrClosed)

	f = nil
	err = f.Close()
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.Seek(0, io.SeekEnd)
	assert.Equal(t, err, os.ErrInvalid)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.Read(bytedata)
	assert.Equal(t, err, os.ErrInvalid)
}

func TestMockOpenFile_PrefetchRead(t *testing.T) {
	length := 11*1024*1024 + 13435
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

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
					sendLen = httpRange.Count
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

			//status code
			w.WriteHeader(statusCode)

			//body
			sendData := data[int(offset):int(offset+sendLen)]
			//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)
			w.Write(sendData)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.EnablePrefetch = true
		oo.ChunkSize = 2 * 1024 * 1024
		oo.PrefetchNum = 3
		oo.PrefetchThreshold = int64(0)
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, true, f.enablePrefetch)
	assert.Equal(t, int64(2*1024*1024), f.chunkSize)
	assert.Equal(t, 3, f.prefetchNum)
	assert.Equal(t, int64(0), f.prefetchThreshold)
	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())
	assert.Equal(t, gmtTime, stat.ModTime().Format(http.TimeFormat))
	assert.Equal(t, os.FileMode(0644), stat.Mode())
	assert.Equal(t, false, stat.IsDir())
	assert.Equal(t, "oss://bucket/key", stat.Name())
	h, ok := stat.Sys().(http.Header)
	assert.True(t, ok)
	assert.NotNil(t, h)
	assert.Equal(t, "fba9dede5f27731c9771645a3986****", h.Get(HTTPHeaderETag))
	assert.Equal(t, "text/plain", h.Get(HTTPHeaderContentType))

	//seek ok
	begin, err := f.Seek(0, io.SeekStart)
	end, err := f.Seek(0, io.SeekEnd)
	assert.Equal(t, stat.Size(), end-begin)

	//seek invalid
	begin, err = f.Seek(0, 4)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid whence")
	assert.Equal(t, int64(0), begin)

	begin, err = f.Seek(-1, io.SeekStart)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "negative position")
	assert.Equal(t, int64(0), begin)

	begin, err = f.Seek(100, io.SeekEnd)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "offset is unavailable")
	assert.Equal(t, int64(0), begin)

	//read all
	f.Seek(0, io.SeekStart)
	hash := NewCRC64(0)
	written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, stat.Size(), written)

	//seek readN
	for i := 0; i < 64; i++ {
		offset := rand.Int63n(int64(length / 5))
		n := rand.Int63n(int64(length/4)) + 3*1024*1024
		begin, err = f.Seek(offset, io.SeekStart)
		//fmt.Printf("seek read check offset%d, len:%d\n", offset, n)
		f.numOOORead = 0
		assert.Nil(t, err)
		assert.Equal(t, offset, begin)

		hash := NewCRC64(0)
		written, err := io.CopyN(io.MultiWriter(io.Discard, hash), f, n)
		assert.Nil(t, err)
		assert.Equal(t, n, written)

		hash1 := NewCRC64(0)
		hash1.Write(data[offset : offset+n])

		assert.Equal(t, hash1.Sum64(), hash.Sum64())
	}

	//seek read from offset to end
	for i := 0; i < 64; i++ {
		offset := rand.Int63n(int64(length / 5))
		begin, err = f.Seek(offset, io.SeekStart)
		n := int64(length) - offset
		//fmt.Printf("seek readAll check offset%d, len:%d\n", offset, n)
		f.numOOORead = 0
		assert.Nil(t, err)
		assert.Equal(t, offset, begin)

		hash := NewCRC64(0)
		written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
		assert.Nil(t, err)
		assert.Equal(t, n, written)

		hash1 := NewCRC64(0)
		hash1.Write(data[offset:])
		assert.Equal(t, hash1.Sum64(), hash.Sum64())
	}

	f.numOOORead = 0
	begin, err = f.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	io.CopyN(io.Discard, f, 2)
	time.Sleep(2 * time.Second)

	err = f.Close()
	assert.Nil(t, err)

	//call Close many times
	err = f.Close()
	assert.Nil(t, err)

	_, err = f.Seek(0, io.SeekEnd)
	assert.Equal(t, err, os.ErrClosed)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrClosed)

	bytedata := make([]byte, 5)
	_, err = f.Read(bytedata)
	assert.Equal(t, err, os.ErrClosed)

	f = nil
	err = f.Close()
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.Seek(0, io.SeekEnd)
	assert.Equal(t, err, os.ErrInvalid)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.Read(bytedata)
	assert.Equal(t, err, os.ErrInvalid)
}

func TestMockOpenFile_OutOfOrderReadThreshold(t *testing.T) {
	length := 11*1024*1024 + 13435
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

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
					sendLen = httpRange.Count
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

			//status code
			w.WriteHeader(statusCode)

			//body
			sendData := data[int(offset):int(offset+sendLen)]
			//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)
			w.Write(sendData)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.EnablePrefetch = true
		oo.ChunkSize = 2 * 1024 * 1024
		oo.PrefetchNum = 3
		oo.PrefetchThreshold = int64(0)
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, true, f.enablePrefetch)
	assert.Equal(t, int64(2*1024*1024), f.chunkSize)
	assert.Equal(t, 3, f.prefetchNum)
	assert.Equal(t, int64(0), f.prefetchThreshold)
	assert.Equal(t, int64(3), f.oooReadThreshold)

	f, err = NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.EnablePrefetch = true
		oo.ChunkSize = 2 * 1024 * 1024
		oo.PrefetchNum = 3
		oo.PrefetchThreshold = int64(0)
		oo.OutOfOrderReadThreshold = int64(1)
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, true, f.enablePrefetch)
	assert.Equal(t, int64(2*1024*1024), f.chunkSize)
	assert.Equal(t, 3, f.prefetchNum)
	assert.Equal(t, int64(0), f.prefetchThreshold)
	assert.Equal(t, int64(1), f.oooReadThreshold)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())
	assert.Equal(t, gmtTime, stat.ModTime().Format(http.TimeFormat))
	assert.Equal(t, os.FileMode(0644), stat.Mode())
	assert.Equal(t, false, stat.IsDir())
	assert.Equal(t, "oss://bucket/key", stat.Name())
	h, ok := stat.Sys().(http.Header)
	assert.True(t, ok)
	assert.NotNil(t, h)
	assert.Equal(t, "fba9dede5f27731c9771645a3986****", h.Get(HTTPHeaderETag))
	assert.Equal(t, "text/plain", h.Get(HTTPHeaderContentType))

	//seek ok
	begin, err := f.Seek(0, io.SeekStart)
	end, err := f.Seek(0, io.SeekEnd)
	assert.Equal(t, stat.Size(), end-begin)

	//seek invalid
	begin, err = f.Seek(0, 4)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid whence")
	assert.Equal(t, int64(0), begin)

	begin, err = f.Seek(-1, io.SeekStart)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "negative position")
	assert.Equal(t, int64(0), begin)

	begin, err = f.Seek(100, io.SeekEnd)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "offset is unavailable")
	assert.Equal(t, int64(0), begin)

	//read all
	f.Seek(0, io.SeekStart)
	hash := NewCRC64(0)
	written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, stat.Size(), written)

	//seek readN
	for i := 0; i < 64; i++ {
		offset := rand.Int63n(int64(length / 5))
		n := rand.Int63n(int64(length/4)) + 3*1024*1024
		begin, err = f.Seek(offset, io.SeekStart)
		//fmt.Printf("seek read check offset%d, len:%d\n", offset, n)
		f.numOOORead = 0
		assert.Nil(t, err)
		assert.Equal(t, offset, begin)

		hash := NewCRC64(0)
		written, err := io.CopyN(io.MultiWriter(io.Discard, hash), f, n)
		assert.Nil(t, err)
		assert.Equal(t, n, written)

		hash1 := NewCRC64(0)
		hash1.Write(data[offset : offset+n])

		assert.Equal(t, hash1.Sum64(), hash.Sum64())
	}

	//seek read from offset to end
	for i := 0; i < 64; i++ {
		offset := rand.Int63n(int64(length / 5))
		begin, err = f.Seek(offset, io.SeekStart)
		n := int64(length) - offset
		//fmt.Printf("seek readAll check offset%d, len:%d\n", offset, n)
		f.numOOORead = 0
		assert.Nil(t, err)
		assert.Equal(t, offset, begin)

		hash := NewCRC64(0)
		written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
		assert.Nil(t, err)
		assert.Equal(t, n, written)

		hash1 := NewCRC64(0)
		hash1.Write(data[offset:])
		assert.Equal(t, hash1.Sum64(), hash.Sum64())
	}

	f.numOOORead = 0
	begin, err = f.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	io.CopyN(io.Discard, f, 2)
	time.Sleep(2 * time.Second)

	err = f.Close()
	assert.Nil(t, err)
}

func TestMockOpenFile_PrefetchRead_AsyncReader(t *testing.T) {
	length := 11*1024*1024 + 13435
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")
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
					sendLen = httpRange.Count
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

			//status code
			w.WriteHeader(statusCode)
			//body
			sendData := data[int(offset):int(offset+sendLen)]
			//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)
			w.Write(sendData)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.EnablePrefetch = true
		oo.ChunkSize = 2 * 1024 * 1024
		oo.PrefetchNum = 3
		oo.PrefetchThreshold = int64(0)
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, true, f.enablePrefetch)
	assert.Equal(t, int64(2*1024*1024), f.chunkSize)
	assert.Equal(t, 3, f.prefetchNum)
	assert.Equal(t, int64(0), f.prefetchThreshold)
	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())
	assert.Equal(t, gmtTime, stat.ModTime().Format(http.TimeFormat))
	assert.Equal(t, os.FileMode(0644), stat.Mode())
	assert.Equal(t, false, stat.IsDir())
	assert.Equal(t, "oss://bucket/key", stat.Name())
	h, ok := stat.Sys().(http.Header)
	assert.True(t, ok)
	assert.NotNil(t, h)
	assert.Equal(t, "fba9dede5f27731c9771645a3986****", h.Get(HTTPHeaderETag))
	assert.Equal(t, "text/plain", h.Get(HTTPHeaderContentType))

	bytedata := make([]byte, 2*1024*1024)
	for {
		n, err := f.Read(bytedata)
		assert.NotNil(t, f.asyncReaders)
		if err != nil || n == 0 {
			break
		}
	}
}

func TestMockOpenFile_MixRead(t *testing.T) {
	length := 11*1024*1024 + 13435
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()

	rangeReqCount := int32(0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "GET":
			atomic.AddInt32(&rangeReqCount, 1)
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
					sendLen = httpRange.Count
				}
				cr := httpContentRange{
					Offset: httpRange.Offset,
					Count:  sendLen,
					Total:  int64(length),
				}
				w.Header().Set("Content-Range", ToString(cr.FormatHTTPContentRange()))
				statusCode = 206
			}

			if atomic.LoadInt32(&rangeReqCount) > 3 && httpRange != nil && httpRange.Count > 0 {
				time.Sleep(2 * time.Second)
			}

			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(sendLen))
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(statusCode)

			//body
			sendData := data[int(offset):int(offset+sendLen)]
			//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)
			w.Write(sendData)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(1 * time.Second).
		WithRetryMaxAttempts(1)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.EnablePrefetch = true
		oo.ChunkSize = 2 * 1024 * 1024
		oo.PrefetchNum = 3
		oo.PrefetchThreshold = 2*1024*1024 + 1234
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, true, f.enablePrefetch)
	assert.Equal(t, int64(2*1024*1024), f.chunkSize)
	assert.Equal(t, 3, f.prefetchNum)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())
	assert.Equal(t, gmtTime, stat.ModTime().Format(http.TimeFormat))
	assert.Equal(t, os.FileMode(0644), stat.Mode())
	assert.Equal(t, false, stat.IsDir())
	assert.Equal(t, "oss://bucket/key", stat.Name())
	h, ok := stat.Sys().(http.Header)
	assert.True(t, ok)
	assert.NotNil(t, h)
	assert.Equal(t, "fba9dede5f27731c9771645a3986****", h.Get(HTTPHeaderETag))
	assert.Equal(t, "text/plain", h.Get(HTTPHeaderContentType))

	//seek ok
	begin, err := f.Seek(0, io.SeekStart)
	end, err := f.Seek(0, io.SeekEnd)
	assert.Equal(t, stat.Size(), end-begin)
	curr, err := f.Seek(0, io.SeekCurrent)
	assert.Equal(t, stat.Size(), curr)

	//read all
	f.Seek(0, io.SeekStart)
	hash := NewCRC64(0)

	written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, stat.Size(), written)

}

func TestMockOpenFile_PrefetchRead_Resume(t *testing.T) {
	length := 11*1024*1024 + 13435
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	halfBodyErr := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

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
					sendLen = httpRange.Count
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

			//status code
			w.WriteHeader(statusCode)

			//body
			sendData := data[int(offset):int(offset+sendLen)]

			if halfBodyErr {
				sendData = sendData[0 : len(sendData)/2]
				halfBodyErr = false
			}

			//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)

			w.Write(sendData)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.EnablePrefetch = true
		oo.ChunkSize = 2 * 1024 * 1024
		oo.PrefetchNum = 3
		oo.PrefetchThreshold = int64(0)
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, true, f.enablePrefetch)
	assert.Equal(t, int64(2*1024*1024), f.chunkSize)
	assert.Equal(t, 3, f.prefetchNum)
	assert.Equal(t, int64(0), f.prefetchThreshold)

	//read all
	halfBodyErr = true
	f.Seek(0, io.SeekStart)
	hash := NewCRC64(0)
	written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, int64(length), written)
}

func TestMockOpenFile_DirectRead_Fail(t *testing.T) {
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	setTimeout := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "GET":
			if setTimeout {
				time.Sleep(2 * time.Second)
			}

			// header
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(data)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())

	//seek ok
	_, err = f.Seek(5, io.SeekStart)
	assert.Nil(t, err)

	//read all
	_, err = io.Copy(io.Discard, f)
	assert.Contains(t, err.Error(), "Range get fail, expect offset")

	err = f.Close()
	assert.Nil(t, err)

	// with invalid offset
	f, err = NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.Offset = int64(length * 2)
	})
	assert.NotNil(t, err)
	assert.Nil(t, f)
	assert.Contains(t, err.Error(), "offset is unavailable, offset")

	//timeout
	cfg = LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(1 * time.Second)

	client = NewClient(cfg)
	f, err = NewReadOnlyFile(context.TODO(), client, "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)
	setTimeout = true
	_, err = io.Copy(io.Discard, f)
	assert.Contains(t, err.Error(), "i/o timeout")
}

func TestMockOpenFile_Constructor(t *testing.T) {
	gmtTime := getNowGMT()
	length := 0
	url := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			url = r.URL.String()
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprintf("%d", length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	length = -1
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key")
	assert.NotNil(t, err)
	assert.Nil(t, f)

	length = 0
	f, err = NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.EnablePrefetch = true
		oo.ChunkSize = -1
		oo.PrefetchNum = -1
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, DefaultPrefetchChunkSize, f.chunkSize)
	assert.Equal(t, DefaultPrefetchNum, f.prefetchNum)

	//version Id
	length = 0
	f, err = NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.VersionId = Ptr("123")
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Contains(t, url, "versionId=123")

	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.NotNil(t, stat)
	assert.Equal(t, "oss://bucket/key?versionId=123", stat.Name())
}

func TestMockOpenFile_DirectRead_FileChange(t *testing.T) {
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "GET":
			// header
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderETag, "0ba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(data)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())

	//seek ok
	_, err = f.Seek(0, io.SeekStart)
	assert.Nil(t, err)

	//read all
	_, err = io.Copy(io.Discard, f)
	assert.Contains(t, err.Error(), "Source file is changed")

	err = f.Close()
	assert.Nil(t, err)
}

func TestMockOpenFile_PrefetchRead_Fail(t *testing.T) {
	length := 5*1024*1024 + 13435
	data := []byte(randStr(length))
	gmtTime := getNowGMT()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "GET":
			// header
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(data)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.EnablePrefetch = true
		oo.ChunkSize = 1 * 1024 * 1024
		oo.PrefetchNum = 3
		oo.PrefetchThreshold = int64(0)
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())

	//seek ok
	_, err = f.Seek(5, io.SeekStart)
	assert.Nil(t, err)

	//read all
	_, err = io.Copy(io.Discard, f)
	assert.Contains(t, err.Error(), "Range get fail, expect offset")

	err = f.Close()
	assert.Nil(t, err)
}

func TestMockOpenFile_PrefetchRead_FileChange(t *testing.T) {
	length := 5*1024*1024 + 13435
	data := []byte(randStr(length))
	gmtTime := getNowGMT()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "GET":
			// header
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderETag, "0ba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(data)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(oo *OpenOptions) {
		oo.EnablePrefetch = true
		oo.ChunkSize = 1*1024*1024 + 1
		oo.PrefetchNum = 2
		oo.PrefetchThreshold = int64(0)
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, 2, f.prefetchNum)
	assert.Equal(t, int64(2*1024*1024), f.chunkSize)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())

	//seek ok
	_, err = f.Seek(0, io.SeekStart)
	assert.Nil(t, err)

	//read all
	_, err = io.Copy(io.Discard, f)
	assert.Contains(t, err.Error(), "Source file is changed")

	err = f.Close()
	assert.Nil(t, err)
}

func TestMockAppendFile_Exist(t *testing.T) {
	data := []byte("start:")
	gmtTime := getNowGMT()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(data)))
			w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
			w.Header().Set(HTTPHeaderContentType, "text/plain")
			w.Header().Set(HeaderOssObjectType, "Appendable")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "POST":
			in, err := io.ReadAll(r.Body)
			assert.Nil(t, err)

			var buffer bytes.Buffer
			buffer.Write(data)
			buffer.Write(in)
			data = buffer.Bytes()

			// header
			w.Header().Set(HTTPHeaderContentLength, "0")
			w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
			w.Header().Set(HTTPHeaderContentType, "text/plain")
			w.Header().Set(HeaderOssNextAppendPosition, fmt.Sprintf("%d", len(data)))

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := client.AppendFile(context.TODO(), "bucket", "key", func(*AppendOptions) {})
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(data)), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	n, err := f.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	n, err = f.Write([]byte(" world"))
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	pattern := "start:hello world"
	assert.Equal(t, pattern, string(data))

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(pattern)), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(pattern)), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	length := 1238
	str := randStr(length)
	written, err := f.WriteFrom(io.NopCloser(bytes.NewReader([]byte(str))))
	assert.Nil(t, err)
	assert.Equal(t, int64(length), written)

	assert.Equal(t, pattern+str, string(data))

	err = f.Close()
	assert.Nil(t, err)

	//call Close many times
	err = f.Close()
	assert.Nil(t, err)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrClosed)

	_, err = f.Write([]byte("world"))
	assert.Equal(t, err, os.ErrClosed)

	_, err = f.WriteFrom(io.NopCloser(bytes.NewReader([]byte("world"))))
	assert.Equal(t, err, os.ErrClosed)

	f = nil
	err = f.Close()
	assert.Equal(t, err, os.ErrInvalid)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.Write([]byte("world"))
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.WriteFrom(io.NopCloser(bytes.NewReader([]byte("world"))))
	assert.Equal(t, err, os.ErrInvalid)
}

func TestMockAppendFile_NoExist(t *testing.T) {
	data := make([]byte, 0)
	gmtTime := getNowGMT()
	first := true

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			if first {
				// header

				//status code
				w.WriteHeader(404)

			} else {
				// header
				w.Header().Set(HTTPHeaderLastModified, gmtTime)
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(data)))
				w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
				w.Header().Set(HTTPHeaderContentType, "text/plain")
				w.Header().Set(HeaderOssObjectType, "Appendable")

				//status code
				w.WriteHeader(200)
			}
			first = false
			//body
			w.Write(nil)

		case "POST":
			in, err := io.ReadAll(r.Body)
			assert.Nil(t, err)

			var buffer bytes.Buffer
			buffer.Write(data)
			buffer.Write(in)
			data = buffer.Bytes()

			// header
			w.Header().Set(HTTPHeaderContentLength, "0")
			w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
			w.Header().Set(HTTPHeaderContentType, "text/plain")
			w.Header().Set(HeaderOssNextAppendPosition, fmt.Sprintf("%d", len(data)))

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := client.AppendFile(context.TODO(), "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	n, err := f.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	n, err = f.Write([]byte(" world"))
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	assert.Equal(t, "hello world", string(data))

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(11), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(11), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	length := 1238
	str := randStr(length)
	written, err := f.WriteFrom(io.NopCloser(bytes.NewReader([]byte(str))))
	assert.Nil(t, err)
	assert.Equal(t, int64(length), written)

	assert.Equal(t, "hello world"+str, string(data))

	err = f.Close()
	assert.Nil(t, err)
}

func TestMockAppendFile_NoAppenable(t *testing.T) {
	data := make([]byte, 0)
	gmtTime := getNowGMT()
	count := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":

			if count == 0 {
				// header
				w.Header().Set(HTTPHeaderLastModified, gmtTime)
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(data)))
				w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
				w.Header().Set(HTTPHeaderContentType, "text/plain")
				w.Header().Set(HeaderOssObjectType, "Normal")

				//status code
				w.WriteHeader(200)
				w.Write(nil)
			} else {
				s := `<?xml version="1.0" encoding="UTF-8"?>
				<Error>
					<Code>SignatureDoesNotMatch</Code>
					<Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
					<RequestId>65467C42E001B4333337****</RequestId>
					<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
					<EC>0002-00000040</EC>
				</Error>`
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(s)))
				w.Header().Set(HTTPHeaderContentType, "application/xml")

				w.WriteHeader(403)
				w.Write([]byte(s))
			}
			count++
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := client.AppendFile(context.TODO(), "bucket", "key")
	assert.NotNil(t, err)
	assert.Nil(t, f)
	assert.Contains(t, err.Error(), "Not a appendable file")

	f, err = client.AppendFile(context.TODO(), "bucket", "key")
	assert.NotNil(t, err)
	assert.Nil(t, f)
	assert.Contains(t, err.Error(), "Http Status Code: 403")
}

func TestMockAppendFile_PositionNotEqualToLength(t *testing.T) {
	var dataMu sync.Mutex
	data := []byte("start:")
	gmtTime := getNowGMT()

	count := int32(0)

	getDataLen := func() int {
		dataMu.Lock()
		defer dataMu.Unlock()
		return len(data)
	}

	getDataValue := func() []byte {
		dataMu.Lock()
		defer dataMu.Unlock()
		return data
	}

	setDataValue := func(v []byte) {
		dataMu.Lock()
		defer dataMu.Unlock()
		data = v
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(data)))
			w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
			w.Header().Set(HTTPHeaderContentType, "text/plain")
			w.Header().Set(HeaderOssObjectType, "Appendable")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)

			//fmt.Printf("Get Head, data len%d\n", len(data))
		case "POST":
			query, _ := url.ParseQuery(r.URL.RawQuery)
			position, _ := strconv.ParseInt(query.Get("position"), 10, 64)
			atomic.AddInt32(&count, 1)
			if position != int64(getDataLen()) {
				//fmt.Printf("position != int64(len(data), position:%d, len:%d\n", position, len(data))

				s := `<?xml version="1.0" encoding="UTF-8"?>
				<Error>
					<Code>PositionNotEqualToLength</Code>
					<Message>error</Message>
					<RequestId>65467C42E001B4333337****</RequestId>
				</Error>`
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(s)))
				w.Header().Set(HTTPHeaderContentType, "application/xml")

				w.WriteHeader(409)
				w.Write([]byte(s))

			} else {
				in, err := io.ReadAll(r.Body)
				assert.Nil(t, err)

				var buffer bytes.Buffer
				buffer.Write(getDataValue())
				buffer.Write(in)
				setDataValue(buffer.Bytes())

				//fmt.Printf("position = int64(len(data), position:%d, new len:%d\n", position, len(data))

				if atomic.LoadInt32(&count) <= 1 {
					time.Sleep(5 * time.Second)
				} else {
					// header
					w.Header().Set(HTTPHeaderContentLength, "0")
					w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
					w.Header().Set(HTTPHeaderContentType, "text/plain")
					w.Header().Set(HeaderOssNextAppendPosition, fmt.Sprintf("%d", len(data)))

					//status code
					w.WriteHeader(200)

					//body
					w.Write(nil)
				}
			}
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(1 * time.Second)

	client := NewClient(cfg)
	f, err := client.AppendFile(context.TODO(), "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(data)), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	n, err := f.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	n, err = f.Write([]byte(" world"))
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	pattern := "start:hello world"
	assert.Equal(t, pattern, string(data))

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(pattern)), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	length := 1238
	str := randStr(length)
	written, err := f.WriteFrom(io.NopCloser(bytes.NewReader([]byte(str))))
	assert.Nil(t, err)
	assert.Equal(t, int64(length), written)

	assert.Equal(t, pattern+str, string(data))

	err = f.Close()
	assert.Nil(t, err)
	f = nil

	// next time
	pattern = "second:"
	data = []byte(pattern)
	atomic.StoreInt32(&count, 0)
	f, err = client.AppendFile(context.TODO(), "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(7), stat.Size())

	length = 1234
	str = randStr(length)
	pattern += str
	written, err = f.WriteFrom(bytes.NewReader([]byte(str)))
	assert.Nil(t, err)
	assert.Equal(t, int64(length), written)
	assert.Equal(t, pattern, string(data))

	n, err = f.Write([]byte("hello"))
	pattern += "hello"
	assert.Nil(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, pattern, string(data))

	//path error
	pattern = "path error:"
	data = []byte(pattern)
	n, err = f.Write([]byte("hello"))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "PositionNotEqualToLength")

	written, err = f.WriteFrom(bytes.NewReader([]byte(str)))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "PositionNotEqualToLength")
}

func TestMockAppendFile_CRC(t *testing.T) {
	data := make([]byte, 0)
	gmtTime := getNowGMT()
	first := true

	crc64Err := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "HEAD":
			if first {
				// header

				//status code
				w.WriteHeader(404)

			} else {
				hashall := NewCRC64(0)
				hashall.Write(data)
				crc64ecma := fmt.Sprint(hashall.Sum64())

				// header
				w.Header().Set(HTTPHeaderLastModified, gmtTime)
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(data)))
				w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
				w.Header().Set(HTTPHeaderContentType, "text/plain")
				w.Header().Set(HeaderOssObjectType, "Appendable")
				w.Header().Set(HeaderOssCRC64, crc64ecma)

				//status code
				w.WriteHeader(200)
			}
			first = false
			//body
			w.Write(nil)

		case "POST":
			in, err := io.ReadAll(r.Body)
			assert.Nil(t, err)

			var buffer bytes.Buffer
			buffer.Write(data)
			buffer.Write(in)
			data = buffer.Bytes()

			hashall := NewCRC64(0)
			hashall.Write(data)
			crc64ecma := fmt.Sprint(hashall.Sum64())

			if crc64Err {
				crc64ecma = "6707180448768400016"
			}

			// header
			w.Header().Set(HTTPHeaderContentLength, "0")
			w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
			w.Header().Set(HTTPHeaderContentType, "text/plain")
			w.Header().Set(HeaderOssNextAppendPosition, fmt.Sprintf("%d", len(data)))
			w.Header().Set(HeaderOssCRC64, crc64ecma)

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	//No Object Exist
	f, err := client.AppendFile(context.TODO(), "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	n, err := f.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	n, err = f.Write([]byte(" world"))
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	assert.Equal(t, "hello world", string(data))

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(11), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(11), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	length := 1238
	str := randStr(length)
	written, err := f.WriteFrom(io.NopCloser(bytes.NewReader([]byte(str))))
	assert.Nil(t, err)
	assert.Equal(t, int64(length), written)

	assert.Equal(t, "hello world"+str, string(data))

	err = f.Close()
	assert.Nil(t, err)

	// Ojbect Exist
	data = []byte("object exist")
	f, err = client.AppendFile(context.TODO(), "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err = f.Stat()
	assert.Equal(t, int64(len("object exist")), stat.Size())

	n, err = f.Write([]byte(" 123456"))
	assert.Nil(t, err)
	assert.Equal(t, 7, n)
	f.Close()

	// CRC not equal
	crc64Err = true
	data = []byte("object exist")
	f, err = client.AppendFile(context.TODO(), "bucket", "key")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err = f.Stat()
	assert.Equal(t, int64(len("object exist")), stat.Size())

	n, err = f.Write([]byte(" 123456"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, n)
	f.Close()
}

func TestMockOpenFile_Payer(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-oss-request-payer") == "" {
			errData := []byte(
				`<?xml version="1.0" encoding="UTF-8"?>
				<Error>
				  <Code>AccessDenied</Code>
				  <Message>Access denied for requester pay bucket</Message>
				  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
				  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
				  <BucketName>test</BucketName>
				  <EC>0003-00000703</EC>
				</Error>`)
			w.Header().Set(HTTPHeaderContentType, "application/xml")
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
			w.WriteHeader(403)
			w.Write(errData)
			return
		}
		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

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
					sendLen = httpRange.Count
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

			//status code
			w.WriteHeader(statusCode)

			//body
			sendData := data[int(offset):int(offset+sendLen)]
			//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)
			w.Write(sendData)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := NewReadOnlyFile(context.TODO(), client, "bucket", "key")
	assert.NotNil(t, err)

	f, err = NewReadOnlyFile(context.TODO(), client, "bucket", "key", func(op *OpenOptions) {
		op.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, false, f.enablePrefetch)
	assert.Equal(t, DefaultPrefetchChunkSize, f.chunkSize)
	assert.Equal(t, DefaultPrefetchNum, f.prefetchNum)
	assert.Equal(t, DefaultPrefetchThreshold, f.prefetchThreshold)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(length), stat.Size())
	assert.Equal(t, gmtTime, stat.ModTime().Format(http.TimeFormat))
	assert.Equal(t, os.FileMode(0644), stat.Mode())
	assert.Equal(t, false, stat.IsDir())
	assert.Equal(t, "oss://bucket/key", stat.Name())
	h, ok := stat.Sys().(http.Header)
	assert.True(t, ok)
	assert.NotNil(t, h)
	assert.Equal(t, "fba9dede5f27731c9771645a3986****", h.Get(HTTPHeaderETag))
	assert.Equal(t, "text/plain", h.Get(HTTPHeaderContentType))

	//seek ok
	begin, err := f.Seek(0, io.SeekStart)
	end, err := f.Seek(0, io.SeekEnd)
	assert.Equal(t, stat.Size(), end-begin)

	//seek invalid
	begin, err = f.Seek(0, 4)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid whence")
	assert.Equal(t, int64(0), begin)

	begin, err = f.Seek(-1, io.SeekStart)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "negative position")
	assert.Equal(t, int64(0), begin)

	begin, err = f.Seek(100, io.SeekEnd)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "offset is unavailable")
	assert.Equal(t, int64(0), begin)

	//read all
	f.Seek(0, io.SeekStart)
	hash := NewCRC64(0)
	written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, stat.Size(), written)
	//seek readN
	for i := 0; i < 64; i++ {
		offset := rand.Int63n(int64(length / 4))
		n := rand.Int63n(int64(length/2)) + 1
		//fmt.Printf("seek readN check offset%d, len:%d\n", offset, n)
		begin, err = f.Seek(offset, io.SeekStart)
		assert.Nil(t, err)
		assert.Equal(t, offset, begin)

		hash := NewCRC64(0)
		written, err := io.CopyN(io.MultiWriter(io.Discard, hash), f, n)
		assert.Nil(t, err)
		assert.Equal(t, n, written)

		hash1 := NewCRC64(0)
		hash1.Write(data[offset : offset+n])

		assert.Equal(t, hash1.Sum64(), hash.Sum64())
	}

	//seek read from offset to end
	for i := 0; i < 64; i++ {
		offset := rand.Int63n(int64(length / 5))
		begin, err = f.Seek(offset, io.SeekStart)
		n := int64(length) - offset
		//fmt.Printf("seek readAll check offset%d, len:%d\n", offset, n)
		assert.Nil(t, err)
		assert.Equal(t, offset, begin)

		hash := NewCRC64(0)
		written, err := io.Copy(io.MultiWriter(io.Discard, hash), f)
		assert.Nil(t, err)
		assert.Equal(t, n, written)

		hash1 := NewCRC64(0)
		hash1.Write(data[offset:])
		assert.Equal(t, hash1.Sum64(), hash.Sum64())
	}

	//seek to end
	begin, err = f.Seek(0, io.SeekEnd)
	assert.Nil(t, err)
	written, err = io.Copy(io.Discard, f)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), written)

	//
	begin, err = f.Seek(0, io.SeekStart)
	assert.Nil(t, err)
	io.CopyN(io.Discard, f, 2)
	time.Sleep(2 * time.Second)

	err = f.Close()
	assert.Nil(t, err)

	//call Close many times
	err = f.Close()
	assert.Nil(t, err)

	_, err = f.Seek(0, io.SeekEnd)
	assert.Equal(t, err, os.ErrClosed)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrClosed)

	bytedata := make([]byte, 5)
	_, err = f.Read(bytedata)
	assert.Equal(t, err, os.ErrClosed)

	f = nil
	err = f.Close()
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.Seek(0, io.SeekEnd)
	assert.Equal(t, err, os.ErrInvalid)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.Read(bytedata)
	assert.Equal(t, err, os.ErrInvalid)
}

func TestMockAppendFile_Payer(t *testing.T) {
	data := []byte("start:")
	gmtTime := getNowGMT()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-oss-request-payer") == "" {
			errData := []byte(
				`<?xml version="1.0" encoding="UTF-8"?>
				<Error>
				  <Code>AccessDenied</Code>
				  <Message>Access denied for requester pay bucket</Message>
				  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
				  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
				  <BucketName>test</BucketName>
				  <EC>0003-00000703</EC>
				</Error>`)
			w.Header().Set(HTTPHeaderContentType, "application/xml")
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
			w.WriteHeader(403)
			w.Write(errData)
			return
		}
		switch r.Method {
		case "HEAD":
			// header
			w.Header().Set(HTTPHeaderLastModified, gmtTime)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(data)))
			w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
			w.Header().Set(HTTPHeaderContentType, "text/plain")
			w.Header().Set(HeaderOssObjectType, "Appendable")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "POST":
			in, err := io.ReadAll(r.Body)
			assert.Nil(t, err)

			var buffer bytes.Buffer
			buffer.Write(data)
			buffer.Write(in)
			data = buffer.Bytes()

			// header
			w.Header().Set(HTTPHeaderContentLength, "0")
			w.Header().Set(HTTPHeaderETag, fmt.Sprintf("etag-%d", len(data)))
			w.Header().Set(HTTPHeaderContentType, "text/plain")
			w.Header().Set(HeaderOssNextAppendPosition, fmt.Sprintf("%d", len(data)))

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		}
	}))
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	f, err := client.AppendFile(context.TODO(), "bucket", "key")
	assert.NotNil(t, err)

	f, err = client.AppendFile(context.TODO(), "bucket", "key", func(ap *AppendOptions) {
		ap.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)

	//stat
	stat, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(data)), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	n, err := f.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	n, err = f.Write([]byte(" world"))
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	pattern := "start:hello world"
	assert.Equal(t, pattern, string(data))

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(pattern)), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	stat, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(pattern)), stat.Size())
	assert.Equal(t, "oss://bucket/key", stat.Name())

	length := 1238
	str := randStr(length)
	written, err := f.WriteFrom(io.NopCloser(bytes.NewReader([]byte(str))))
	assert.Nil(t, err)
	assert.Equal(t, int64(length), written)

	assert.Equal(t, pattern+str, string(data))

	err = f.Close()
	assert.Nil(t, err)

	//call Close many times
	err = f.Close()
	assert.Nil(t, err)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrClosed)

	_, err = f.Write([]byte("world"))
	assert.Equal(t, err, os.ErrClosed)

	_, err = f.WriteFrom(io.NopCloser(bytes.NewReader([]byte("world"))))
	assert.Equal(t, err, os.ErrClosed)

	f = nil
	err = f.Close()
	assert.Equal(t, err, os.ErrInvalid)

	stat, err = f.Stat()
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.Write([]byte("world"))
	assert.Equal(t, err, os.ErrInvalid)

	_, err = f.WriteFrom(io.NopCloser(bytes.NewReader([]byte("world"))))
	assert.Equal(t, err, os.ErrInvalid)

}
