package oss

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

type uploaderMockTracker struct {
	partNum        int
	saveDate       [][]byte
	checkTime      []time.Time
	timeout        []time.Duration
	uploadPartCnt  int32
	putObjectCnt   int32
	contentType    string
	uploadPartErr  []bool
	InitiateMPErr  bool
	CompleteMPErr  bool
	AbortMPErr     bool
	putObjectErr   bool
	ListPartsErr   bool
	crcPartInvalid []bool
	CompleteMPData []byte
}

func testSetupUploaderMockServer(t *testing.T, tracker *uploaderMockTracker) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf("r.URL :%s\n", r.URL.String())
		errData := []byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>InvalidAccessKeyId</Code>
				<Message>The OSS Access Key Id you provided does not exist in our records.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>ak</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`)

		query := r.URL.Query()
		switch r.Method {
		case "POST":
			//url := r.URL.String()
			//strings.Contains(url, "/bucket/key?uploads")
			if query.Get("uploads") == "" && query.Get("uploadId") == "" {
				// InitiateMultipartUpload
				if tracker.InitiateMPErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				sendData := []byte(`
				<InitiateMultipartUploadResult>
					<Bucket>bucket</Bucket>
					<Key>key</Key>
					<UploadId>uploadId-1234</UploadId>
				</InitiateMultipartUploadResult>`)

				tracker.contentType = r.Header.Get(HTTPHeaderContentType)

				w.Header().Set(HTTPHeaderContentType, "application/xml")
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(sendData)))
				w.WriteHeader(200)
				w.Write(sendData)
			} else if query.Get("uploadId") == "uploadId-1234" {
				// strings.Contains(url, "/bucket/key?uploadId=uploadId-1234")
				// CompleteMultipartUpload
				if tracker.CompleteMPErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				sendData := []byte(`
				<CompleteMultipartUploadResult>
					<EncodingType>url</EncodingType>
					<Location>bucket/key</Location>
					<Bucket>bucket</Bucket>
					<Key>key</Key>
					<ETag>etag</ETag>
			  	</CompleteMultipartUploadResult>`)

				tracker.CompleteMPData, _ = io.ReadAll(r.Body)

				hash := NewCRC64(0)
				mr := NewMultiBytesReader(tracker.saveDate)
				io.Copy(io.MultiWriter(hash), mr)
				w.Header().Set(HTTPHeaderContentType, "application/xml")
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(sendData)))
				crc64ecma := fmt.Sprint(hash.Sum64())
				w.Header().Set(HeaderOssCRC64, crc64ecma)
				w.WriteHeader(200)
				w.Write(sendData)
			} else {
				assert.Fail(t, "not support")
			}
		case "PUT":
			in, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			hash := NewCRC64(0)
			hash.Write(in)
			crc64ecma := fmt.Sprint(hash.Sum64())

			md5hash := md5.New()
			md5hash.Write(in)
			etag := fmt.Sprintf("\"%s\"", strings.ToUpper(hex.EncodeToString(md5hash.Sum(nil))))

			if query.Get("uploadId") == "uploadId-1234" {
				// UploadPart
				//in, err := io.ReadAll(r.Body)
				//assert.Nil(t, err)
				num, err := strconv.Atoi(query.Get("partNumber"))
				assert.Nil(t, err)
				assert.LessOrEqual(t, num, tracker.partNum)
				assert.Nil(t, err)
				assert.Equal(t, "uploadId-1234", query.Get("uploadId"))

				//hash := NewCRC64(0)
				//hash.Write(in)
				//crc64ecma := fmt.Sprint(hash.Sum64())
				if tracker.timeout[num-1] > 0 {
					time.Sleep(tracker.timeout[num-1])
				} else {
					time.Sleep(10 * time.Millisecond)
				}

				if tracker.uploadPartErr[num-1] {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				tracker.saveDate[num-1] = in

				// header
				if tracker.crcPartInvalid != nil && tracker.crcPartInvalid[num-1] {
					w.Header().Set(HeaderOssCRC64, "12345")
				} else {
					w.Header().Set(HeaderOssCRC64, crc64ecma)
				}
				w.Header().Set(HTTPHeaderETag, etag)

				//status code
				w.WriteHeader(200)

				//body
				w.Write(nil)
				tracker.checkTime[num-1] = time.Now()
				//fmt.Printf("UploadPart done, num :%d, %v\n", num, tracker.checkTime[num-1])
				atomic.AddInt32(&tracker.uploadPartCnt, 1)
			} else if query.Get("uploadId") == "" {
				tracker.contentType = r.Header.Get(HTTPHeaderContentType)

				if tracker.putObjectErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				//PutObject
				w.Header().Set(HeaderOssCRC64, crc64ecma)
				w.Header().Set(HTTPHeaderETag, etag)

				//status code
				w.WriteHeader(200)

				//body
				w.Write(nil)
				tracker.saveDate[0] = in
				tracker.checkTime[0] = time.Now()
				atomic.AddInt32(&tracker.putObjectCnt, 1)
			} else {
				assert.Fail(t, "not support")
			}
		case "DELETE":
			if query.Get("uploadId") == "uploadId-1234" {
				// AbortMultipartUpload
				if tracker.AbortMPErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				w.WriteHeader(204)
				w.Write(nil)
			} else {
				assert.Fail(t, "not support")
			}
		case "GET":
			if query.Get("uploadId") == "uploadId-1234" {
				// ListParts
				if tracker.ListPartsErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				var buf strings.Builder
				buf.WriteString("<ListPartsResult>")
				buf.WriteString("  <Bucket>bucket</Bucket>")
				buf.WriteString("  <Key>key</Key>")
				buf.WriteString("  <UploadId>uploadId-1234</UploadId>")
				buf.WriteString("  <IsTruncated>false</IsTruncated>")
				for i, d := range tracker.saveDate {
					if d != nil {
						buf.WriteString("  <Part>")
						buf.WriteString(fmt.Sprintf("    <PartNumber>%v</PartNumber>", i+1))
						buf.WriteString("    <LastModified>2012-02-23T07:01:34.000Z</LastModified>")
						buf.WriteString("    <ETag>etag</ETag>")
						buf.WriteString(fmt.Sprintf("    <Size>%v</Size>", len(d)))
						hash := NewCRC64(0)
						hash.Write(d)
						buf.WriteString(fmt.Sprintf("    <HashCrc64ecma>%v</HashCrc64ecma>", fmt.Sprint(hash.Sum64())))
						buf.WriteString("  </Part>")
					}
				}
				buf.WriteString("</ListPartsResult>")

				data := buf.String()
				w.Header().Set(HTTPHeaderContentType, "application/xml")
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(data)))
				w.WriteHeader(200)
				w.Write([]byte(data))
			}
		}
	}))
	return server
}

func TestMockUploadSinglePart(t *testing.T) {
	partSize := DefaultUploadPartSize
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := NewUploader(client)

	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	result, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)
	assert.Equal(t, int32(1), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.uploadPartCnt))
}

func TestMockUploadSequential(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	result, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	index := 3
	ctime := tracker.checkTime[index]
	for i, t := range tracker.checkTime {
		if t.After(ctime) {
			index = i
			ctime = t
		}
	}
	assert.Equal(t, partsNum-1, index)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
}

func TestMockUploadSequentialWithTeeReader(t *testing.T) {
	partSize := int64(10 * 1024 * 1024)
	length := 20 * 1024 * 1024
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)
	pReader := io.TeeReader(bytes.NewReader(data), md5.New())
	result, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		pReader)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)
}

func TestMockUploadParallel(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond

	result, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		bytes.NewReader(data))
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	index := 3
	ctime := tracker.checkTime[index]
	for i, t := range tracker.checkTime {
		if t.After(ctime) {
			index = i
			ctime = t
		}
	}
	assert.Equal(t, 0, index)
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
}

func TestMockUploadParallelWithTeeReader(t *testing.T) {
	partSize := int64(10 * 1024 * 1024)
	length := 40 * 1024 * 1024
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond

	pReader := io.TeeReader(bytes.NewReader(data), md5.New())
	result, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		pReader)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)
}

func TestMockUploadArgmentCheck(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)
	u := NewUploader(client)
	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	// upload stream
	_, err := u.UploadFrom(context.TODO(), nil, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request")

	_, err = u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: nil,
		Key:    Ptr("key"),
	}, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request.Bucket")

	_, err = u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    nil,
	}, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request.Key")

	// upload file
	_, err = u.UploadFile(context.TODO(), nil, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request")

	_, err = u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: nil,
		Key:    Ptr("key"),
	}, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request.Bucket")

	_, err = u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    nil,
	}, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request.Key")

	//Invalid filePath
	_, err = u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, "#@!Ainvalud-path")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "File not exists,")

	// nil body
	_, err = u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the body is null")
}

type noSeeker struct {
	io.Reader
}

func newNoSeeker(r io.Reader) noSeeker {
	return noSeeker{r}
}

type fakeSeeker struct {
	r io.Reader
	n int64
	i int64
}

func newFakeSeeker(r io.Reader, n int64) fakeSeeker {
	return fakeSeeker{r: r, n: n, i: 0}
}

func (r fakeSeeker) Read(p []byte) (n int, err error) {
	return r.Read(p)
}

func (r fakeSeeker) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = r.n + offset
	default:
		return 0, errors.New("MultiSliceReader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("MultiSliceReader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}

func createFile(t *testing.T, fileName, content string) {
	fout, err := os.Create(fileName)
	assert.Nil(t, err)
	defer fout.Close()
	_, err = fout.WriteString(content)
	assert.Nil(t, err)
}

func createFileFromByte(t *testing.T, fileName string, content []byte) {
	fout, err := os.Create(fileName)
	assert.Nil(t, err)
	defer fout.Close()
	_, err = fout.Write(content)
	assert.Nil(t, err)
}

func TestUpload_newDelegate(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)
	u := NewUploader(client)
	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	// nil body
	d, err := u.newDelegate(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})

	assert.Nil(t, err)
	assert.Equal(t, DefaultUploadParallel, d.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, d.options.PartSize)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(0), d.totalSize)
	assert.Equal(t, "", d.filePath)

	assert.Nil(t, d.partPool)
	assert.Nil(t, d.body)
	assert.NotNil(t, d.client)
	assert.NotNil(t, d.context)

	assert.NotNil(t, d.request)

	// empty body
	d, err = u.newDelegate(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})

	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(0), d.totalSize)
	assert.Nil(t, d.body)

	// non-empty body
	d, err = u.newDelegate(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})
	assert.Nil(t, err)
	d.body = bytes.NewReader([]byte("123"))
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(3), d.totalSize)
	assert.NotNil(t, d.body)

	// non-empty without seek body
	d, err = u.newDelegate(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})
	assert.Nil(t, err)
	d.body = newNoSeeker(bytes.NewReader([]byte("123")))
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(-1), d.totalSize)
	assert.NotNil(t, d.body)

	//file path check
	var localFile = randStr(8) + ".txt"
	createFile(t, localFile, "12345")
	d, err = u.newDelegate(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})
	assert.Nil(t, err)
	f, err := os.Open(localFile)
	assert.Nil(t, err)
	d.body = f
	err = d.applySource()
	f.Close()
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(5), d.totalSize)
	assert.NotNil(t, d.body)
	os.Remove(localFile)

	// options
	// non-empty body
	d, err = u.newDelegate(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		}, func(uo *UploaderOptions) {
			uo.ParallelNum = 10
			uo.PartSize = 10
		})
	assert.Nil(t, err)
	d.body = bytes.NewReader([]byte("123"))
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, 10, d.options.ParallelNum)
	assert.Equal(t, int64(10), d.options.PartSize)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	// non-empty body
	d, err = u.newDelegate(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		}, func(uo *UploaderOptions) {
			uo.ParallelNum = 0
			uo.PartSize = 0
		})
	assert.Nil(t, err)
	d.body = bytes.NewReader([]byte("123"))
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, DefaultUploadParallel, d.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	//adjust partSize
	maxSize := DefaultUploadPartSize * int64(MaxUploadParts*4)
	d, err = u.newDelegate(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})
	assert.Nil(t, err)
	d.body = newFakeSeeker(bytes.NewReader([]byte("123")), maxSize)
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, maxSize, d.totalSize)
	assert.Equal(t, DefaultUploadPartSize*5, d.options.PartSize)
}

func TestMockUploadSinglePartFromFile(t *testing.T) {
	partSize := DefaultUploadPartSize
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := randStr(8) + ".txt"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := NewUploader(client)

	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	result, err := u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)
	assert.Equal(t, int32(1), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.uploadPartCnt))
	assert.Equal(t, "text/plain", tracker.contentType)
}

func TestMockUploadSequentialFromFile(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := randStr(8) + ".tif"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	result, err := u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	index := 3
	ctime := tracker.checkTime[index]
	for i, t := range tracker.checkTime {
		if t.After(ctime) {
			index = i
			ctime = t
		}
	}
	assert.Equal(t, partsNum-1, index)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
	assert.Equal(t, "image/tiff", tracker.contentType)
}

func TestMockUploadParallelFromFile(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := randStr(8) + "-no-surfix"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond

	result, err := u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	index := 3
	ctime := tracker.checkTime[index]
	for i, t := range tracker.checkTime {
		if t.After(ctime) {
			index = i
			ctime = t
		}
	}
	assert.Equal(t, 0, index)
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
}

func TestMockUploadWithEmptyBody(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond

	result, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(nil))
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.UploadId)
	assert.Equal(t, "0", *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(all))
	assert.Equal(t, int32(1), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
}

func TestMockUploadSinglePartFail(t *testing.T) {
	partSize := DefaultUploadPartSize
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
		putObjectErr:  true,
	}

	data := []byte(randStr(length))

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := NewUploader(client)

	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	_, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
}

func TestMockUploadSequentialInitiateMultipartUploadFail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
		InitiateMPErr: true,
	}

	data := []byte(randStr(length))

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	_, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
}

func TestMockUploadSequentialUploadPartFail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}
	tracker.uploadPartErr[1] = true

	data := []byte(randStr(length))

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	_, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "uploadId-1234", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
}

func TestMockUploadSequentialCompleteMultipartUploadFail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}
	tracker.CompleteMPErr = true

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	_, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "uploadId-1234", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)
}

func TestMockUploadParallelUploadPartFail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}
	tracker.uploadPartErr[2] = true

	data := []byte(randStr(length))

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 2
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 2, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	_, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "uploadId-1234", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)

	assert.NotNil(t, tracker.saveDate[0])
	assert.NotNil(t, tracker.saveDate[1])
	assert.Nil(t, tracker.saveDate[2])
	assert.Nil(t, tracker.saveDate[5])
}

func TestMockUploaderUploadFileEnableCheckpointNotUseCp(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := randStr(8) + "-no-surfix"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
			uo.CheckpointDir = "."
			uo.EnableCheckpoint = true
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond

	result, err := u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	index := 3
	ctime := tracker.checkTime[index]
	for i, t := range tracker.checkTime {
		if t.After(ctime) {
			index = i
			ctime = t
		}
	}
	assert.Equal(t, 0, index)
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
}

func TestMockUploaderUploadFileEnableCheckpointUseCp(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := "upload-file-with-cp-no-surfix"
	absPath, _ := filepath.Abs(localFile)
	hashmd5 := md5.New()
	hashmd5.Write([]byte(absPath))
	srcHash := hex.EncodeToString(hashmd5.Sum(nil))
	cpFile := srcHash + "-d36fc07f5d963b319b1b48e20a9b8ae9.ucp"

	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 5
			uo.PartSize = partSize
			uo.CheckpointDir = "."
			uo.EnableCheckpoint = true
		},
	)
	assert.Equal(t, 5, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	// Case 1, fail in part number 4
	tracker.saveDate = make([][]byte, partsNum)
	tracker.checkTime = make([]time.Time, partsNum)
	tracker.timeout = make([]time.Duration, partsNum)
	tracker.uploadPartErr = make([]bool, partsNum)
	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond
	tracker.uploadPartErr[3] = true
	os.Remove(cpFile)

	result, err := u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "uploadId-1234", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)

	assert.NotNil(t, tracker.saveDate[0])
	assert.NotNil(t, tracker.saveDate[1])
	assert.NotNil(t, tracker.saveDate[2])
	assert.Nil(t, tracker.saveDate[3])
	assert.NotNil(t, tracker.saveDate[4])

	assert.FileExists(t, cpFile)

	//retry
	time.Sleep(2 * time.Second)
	retryTime := time.Now()
	tracker.uploadPartErr[3] = false
	atomic.StoreInt32(&tracker.uploadPartCnt, 0)

	result, err = u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.True(t, tracker.checkTime[0].Before(retryTime))
	assert.True(t, tracker.checkTime[1].Before(retryTime))
	assert.True(t, tracker.checkTime[2].Before(retryTime))
	assert.True(t, tracker.checkTime[3].After(retryTime))
	assert.True(t, tracker.checkTime[4].After(retryTime))
	assert.True(t, tracker.checkTime[5].After(retryTime))

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(3), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
	assert.Equal(t, strings.Count(string(tracker.CompleteMPData), "<PartNumber>"), 6)
	assert.Equal(t, strings.Count(string(tracker.CompleteMPData), "<PartNumber>1</PartNumber>"), 1)

	assert.NoFileExists(t, cpFile)

	// Case 2, fail in part number 1
	tracker.saveDate = make([][]byte, partsNum)
	tracker.checkTime = make([]time.Time, partsNum)
	tracker.timeout = make([]time.Duration, partsNum)
	tracker.uploadPartErr = make([]bool, partsNum)
	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond
	tracker.uploadPartErr[0] = true
	os.Remove(cpFile)

	result, err = u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.NotNil(t, err)
	assert.Nil(t, tracker.saveDate[0])
	assert.NotNil(t, tracker.saveDate[1])
	assert.NotNil(t, tracker.saveDate[2])
	assert.NotNil(t, tracker.saveDate[3])
	assert.NotNil(t, tracker.saveDate[4])

	assert.FileExists(t, cpFile)

	//retry
	time.Sleep(2 * time.Second)
	retryTime = time.Now()
	tracker.uploadPartErr[0] = false
	atomic.StoreInt32(&tracker.uploadPartCnt, 0)

	result, err = u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.True(t, tracker.checkTime[0].After(retryTime))
	assert.True(t, tracker.checkTime[1].After(retryTime))
	assert.True(t, tracker.checkTime[2].After(retryTime))
	assert.True(t, tracker.checkTime[3].After(retryTime))
	assert.True(t, tracker.checkTime[4].After(retryTime))
	assert.True(t, tracker.checkTime[5].After(retryTime))

	mr = NewMultiBytesReader(tracker.saveDate)
	all, err = io.ReadAll(mr)
	assert.Nil(t, err)

	hashall = NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma = fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(6), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
	assert.NoFileExists(t, cpFile)

	// Case 3, list Parts Fail
	tracker.saveDate = make([][]byte, partsNum)
	tracker.checkTime = make([]time.Time, partsNum)
	tracker.timeout = make([]time.Duration, partsNum)
	tracker.uploadPartErr = make([]bool, partsNum)
	tracker.uploadPartErr[3] = true
	os.Remove(cpFile)

	result, err = u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.NotNil(t, err)
	assert.NotNil(t, tracker.saveDate[0])
	assert.NotNil(t, tracker.saveDate[1])
	assert.NotNil(t, tracker.saveDate[2])
	assert.Nil(t, tracker.saveDate[3])
	assert.NotNil(t, tracker.saveDate[4])

	assert.FileExists(t, cpFile)

	//retry
	time.Sleep(2 * time.Second)
	retryTime = time.Now()
	tracker.uploadPartErr[3] = false
	tracker.ListPartsErr = true
	atomic.StoreInt32(&tracker.uploadPartCnt, 0)

	result, err = u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.True(t, tracker.checkTime[0].After(retryTime))
	assert.True(t, tracker.checkTime[1].After(retryTime))
	assert.True(t, tracker.checkTime[2].After(retryTime))
	assert.True(t, tracker.checkTime[3].After(retryTime))
	assert.True(t, tracker.checkTime[4].After(retryTime))
	assert.True(t, tracker.checkTime[5].After(retryTime))

	mr = NewMultiBytesReader(tracker.saveDate)
	all, err = io.ReadAll(mr)
	assert.Nil(t, err)

	hashall = NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma = fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(6), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
	assert.NoFileExists(t, cpFile)
}

func TestMockUploadParallelFromStreamWithoutSeeker(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := randStr(8) + "-no-surfix"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond

	file, err := os.Open(localFile)
	assert.Nil(t, err)
	defer file.Close()

	result, err := u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, io.LimitReader(file, int64(length)))
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
}

func TestMockUploadCRC64Fail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:        partsNum,
		saveDate:       make([][]byte, partsNum),
		checkTime:      make([]time.Time, partsNum),
		timeout:        make([]time.Duration, partsNum),
		uploadPartErr:  make([]bool, partsNum),
		crcPartInvalid: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)
	tracker.crcPartInvalid[2] = true
	_, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "crc is inconsistent")

	//disable crc check
	client = NewClient(cfg,
		func(o *Options) {
			o.FeatureFlags = o.FeatureFlags & ^FeatureEnableCRC64CheckUpload
		})

	u = NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)
	tracker.crcPartInvalid[2] = true
	tracker.saveDate = make([][]byte, partsNum)
	result, err := u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)
}

func TestMockUploadWithPayer(t *testing.T) {
	partSize := DefaultUploadPartSize
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)

	localFile := randStr(8) + ".txt"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := NewUploader(client)

	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	result, err := u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket:       Ptr("bucket"),
		Key:          Ptr("key"),
		RequestPayer: Ptr("requester"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.UploadId)

	_, err = u.UploadFrom(
		context.TODO(),
		&PutObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("key"),
			RequestPayer: Ptr("requester"),
		},
		bytes.NewReader(data))
	assert.Nil(t, err)
}

func TestMockUploadSinglePartFromFileWithProgress(t *testing.T) {
	partSize := DefaultUploadPartSize
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := randStr(8) + ".txt"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := NewUploader(client)

	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	n := int64(0)
	result, err := u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		ProgressFn: func(increment, transferred, total int64) {
			n = transferred
			fmt.Printf("increment:%#v, transferred:%#v, total:%#v\n", increment, transferred, total)
		},
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)
	assert.Equal(t, n, int64(length))
}

func TestMockUploaderUploadFileEnableCheckpointUseCpProgress(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := "upload-file-with-cp-no-surfix"
	absPath, _ := filepath.Abs(localFile)
	hashmd5 := md5.New()
	hashmd5.Write([]byte(absPath))
	srcHash := hex.EncodeToString(hashmd5.Sum(nil))
	cpFile := srcHash + "-d36fc07f5d963b319b1b48e20a9b8ae9.ucp"

	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 5
			uo.PartSize = partSize
			uo.CheckpointDir = "."
			uo.EnableCheckpoint = true
		},
	)
	assert.Equal(t, 5, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	// Case 1, fail in part number 4
	tracker.saveDate = make([][]byte, partsNum)
	tracker.checkTime = make([]time.Time, partsNum)
	tracker.timeout = make([]time.Duration, partsNum)
	tracker.uploadPartErr = make([]bool, partsNum)
	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond
	tracker.uploadPartErr[3] = true
	os.Remove(cpFile)

	inc := int64(0)
	result, err := u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		ProgressFn: func(increment, transferred, total int64) {
			inc += increment
			//fmt.Printf("increment:%#v, transferred:%#v, total:%#v\n", increment, transferred, total)
		},
	}, localFile)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "uploadId-1234", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)

	assert.NotNil(t, tracker.saveDate[0])
	assert.NotNil(t, tracker.saveDate[1])
	assert.NotNil(t, tracker.saveDate[2])
	assert.Nil(t, tracker.saveDate[3])
	assert.NotNil(t, tracker.saveDate[4])

	assert.FileExists(t, cpFile)

	assert.Less(t, inc, int64(length))

	//retry
	time.Sleep(2 * time.Second)
	retryTime := time.Now()
	tracker.uploadPartErr[3] = false
	atomic.StoreInt32(&tracker.uploadPartCnt, 0)

	inc = 0
	result, err = u.UploadFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		ProgressFn: func(increment, transferred, total int64) {
			inc += increment
			//fmt.Printf("increment:%#v, transferred:%#v, total:%#v\n", increment, transferred, total)
		},
	}, localFile)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.True(t, tracker.checkTime[0].Before(retryTime))
	assert.True(t, tracker.checkTime[1].Before(retryTime))
	assert.True(t, tracker.checkTime[2].Before(retryTime))
	assert.True(t, tracker.checkTime[3].After(retryTime))
	assert.True(t, tracker.checkTime[4].After(retryTime))
	assert.True(t, tracker.checkTime[5].After(retryTime))

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(3), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
	assert.Equal(t, strings.Count(string(tracker.CompleteMPData), "<PartNumber>"), 6)
	assert.Equal(t, strings.Count(string(tracker.CompleteMPData), "<PartNumber>1</PartNumber>"), 1)

	assert.NoFileExists(t, cpFile)
	assert.Equal(t, int64(length), inc)
}

func TestMockUploadWithMixedError(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}
	tracker.uploadPartErr[1] = true

	data := []byte(randStr(length))

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 2 * time.Second
	ctx, cancfg := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancfg()

	pReader := io.TeeReader(bytes.NewReader(data), io.Discard)

	_, err := u.UploadFrom(
		ctx,
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		pReader)
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "uploadId-1234", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)
	assert.Contains(t, uerr.Error(), "context deadline exceeded")
}
