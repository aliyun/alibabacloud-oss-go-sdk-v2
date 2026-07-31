//go:build integration

package oss

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/crypto"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestEncryptionClient(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	objectName := objectNamePrefix + randLowStr(6)

	length := 3*100*1024 + 123
	partSize := int64(200 * 1024)
	partsNum := length/int(partSize) + 1
	data := []byte(randStr(length))
	hashData := NewCRC64(0)
	hashData.Write(data)

	client := getDefaultClient()
	assert.NotNil(t, client)

	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	mc, err := crypto.CreateMasterRsa(map[string]string{"tag": "value"}, rsaPublicKey, rsaPrivateKey)
	assert.Nil(t, err)
	eclient, err := NewEncryptionClient(client, mc)
	assert.Nil(t, err)

	initResult, err := eclient.InitiateMultipartUpload(context.TODO(), &InitiateMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectName),
		CSEPartSize: Ptr(partSize),
		CSEDataSize: Ptr(int64(length)),
	})
	assert.Nil(t, err)
	assert.NotNil(t, initResult)
	assert.NotNil(t, initResult.CSEMultiPartContext)
	assert.NotNil(t, initResult.CSEMultiPartContext.ContentCipher)
	assert.Equal(t, partSize, initResult.CSEMultiPartContext.PartSize)
	assert.Equal(t, int64(length), initResult.CSEMultiPartContext.DataSize)

	var parts UploadParts
	for i := 0; i < partsNum; i++ {
		start := i * int(partSize)
		end := start + int(partSize)
		end = minInt(end, length)
		var contentLength *int64 = nil
		if i%2 == 0 {
			contentLength = Ptr(int64(end - start))
		}
		upResult, err := eclient.UploadPart(context.TODO(), &UploadPartRequest{
			Bucket:              Ptr(bucketName),
			Key:                 Ptr(objectName),
			UploadId:            initResult.UploadId,
			PartNumber:          int32(i + 1),
			CSEMultiPartContext: initResult.CSEMultiPartContext,
			ContentLength:       contentLength,
			Body:                bytes.NewReader(data[start:end]),
		})
		assert.Nil(t, err)
		assert.NotNil(t, upResult)
		parts = append(parts, UploadPart{PartNumber: int32(i + 1), ETag: upResult.ETag})
	}

	lsResult, err := eclient.ListParts(context.TODO(), &ListPartsRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: initResult.UploadId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, lsResult)

	sort.Sort(parts)
	cmResult, err := eclient.CompleteMultipartUpload(context.TODO(), &CompleteMultipartUploadRequest{
		Bucket:                  Ptr(bucketName),
		Key:                     Ptr(objectName),
		UploadId:                initResult.UploadId,
		CompleteMultipartUpload: &CompleteMultipartUpload{Parts: parts},
	})
	assert.Nil(t, err)
	assert.NotNil(t, cmResult)

	// GetObject
	gResult, err := eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err := io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.Len(t, gData, length)
	assert.EqualValues(t, data, gData)

	assert.NotEmpty(t, gResult.Headers.Get(OssClientSideEncryptionKey))
	assert.NotEmpty(t, gResult.Headers.Get(OssClientSideEncryptionStart))
	assert.Equal(t, crypto.AesCtrAlgorithm, gResult.Headers.Get(OssClientSideEncryptionCekAlg))
	assert.Equal(t, crypto.RsaCryptoWrap, gResult.Headers.Get(OssClientSideEncryptionWrapAlg))
	assert.Equal(t, "{\"tag\":\"value\"}", gResult.Headers.Get(OssClientSideEncryptionMatDesc))
	assert.Equal(t, fmt.Sprint(partSize), gResult.Headers.Get(OssClientSideEncryptionPartSize))
	assert.Equal(t, fmt.Sprint(length), gResult.Headers.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, gResult.Headers.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, gResult.Headers.Get(OssClientSideEncryptionUnencryptedContentMD5))

	// HeadObject
	hResult, err := eclient.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.NotEmpty(t, hResult.Headers.Get(OssClientSideEncryptionKey))
	assert.NotEmpty(t, hResult.Headers.Get(OssClientSideEncryptionStart))
	assert.Equal(t, crypto.AesCtrAlgorithm, hResult.Headers.Get(OssClientSideEncryptionCekAlg))
	assert.Equal(t, crypto.RsaCryptoWrap, hResult.Headers.Get(OssClientSideEncryptionWrapAlg))
	assert.Equal(t, "{\"tag\":\"value\"}", hResult.Headers.Get(OssClientSideEncryptionMatDesc))
	assert.Equal(t, fmt.Sprint(partSize), hResult.Headers.Get(OssClientSideEncryptionPartSize))
	assert.Equal(t, fmt.Sprint(length), hResult.Headers.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, hResult.Headers.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, hResult.Headers.Get(OssClientSideEncryptionUnencryptedContentMD5))
	assert.Equal(t, int64(length), hResult.ContentLength)

	// HeadObject
	gmResult, err := eclient.GetObjectMeta(context.TODO(), &GetObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionKey))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionStart))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionCekAlg))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionWrapAlg))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionMatDesc))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionPartSize))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, gmResult.Headers.Get(OssClientSideEncryptionUnencryptedContentMD5))
	assert.Equal(t, int64(length), gmResult.ContentLength)

	// Downloader with not 16 align partSize
	d := eclient.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 123 * 1024
	})
	assert.NotNil(t, d)
	assert.Equal(t, int64(123*1024), d.options.PartSize)
	assert.Equal(t, 3, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"
	dResult, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName)},
		localFile)
	defer os.Remove(localFile)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(gData)), dResult.Written)
	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)
	io.Copy(hash, rfile)
	rfile.Close()
	assert.Equal(t, hash.Sum64(), hashData.Sum64())

	//Use ReadOnlyFile
	f, err := eclient.OpenFile(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	for i := 13; i < 42; i++ {
		for len := 100*1024 + 123; len < 100*1024+123+17; len++ {
			_, err := f.Seek(int64(i), io.SeekStart)
			assert.Nil(t, err)
			gData, err := io.ReadAll(io.LimitReader(f, int64(len)))
			assert.Nil(t, err)
			assert.EqualValues(t, data[i:i+len], gData)
		}
	}
	f.Close()
	time.Sleep(2 * time.Second)

	// Use Uploader
	lastEtag := hResult.Headers.Get(HTTPHeaderETag)
	assert.NotEmpty(t, lastEtag)
	u := eclient.NewUploader()
	assert.NotNil(t, u)
	urResult, err := u.UploadFrom(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
		},
		bytes.NewReader(data),
		func(uo *UploaderOptions) {
			uo.ParallelNum = 2
			uo.PartSize = 100 * 1024
		},
	)
	if !assert.Nil(t, err) {
		fmt.Printf("%s", err.Error())
	}
	assert.NotNil(t, urResult)

	// GetObject again
	gResult, err = eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err = io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.Len(t, gData, length)
	assert.EqualValues(t, data, gData)

	assert.NotEmpty(t, gResult.Headers.Get(OssClientSideEncryptionKey))
	assert.NotEmpty(t, gResult.Headers.Get(OssClientSideEncryptionStart))
	assert.Equal(t, crypto.AesCtrAlgorithm, gResult.Headers.Get(OssClientSideEncryptionCekAlg))
	assert.Equal(t, crypto.RsaCryptoWrap, gResult.Headers.Get(OssClientSideEncryptionWrapAlg))
	assert.Equal(t, "{\"tag\":\"value\"}", gResult.Headers.Get(OssClientSideEncryptionMatDesc))
	assert.Equal(t, fmt.Sprint(100*1024), gResult.Headers.Get(OssClientSideEncryptionPartSize))
	assert.Equal(t, fmt.Sprint(length), gResult.Headers.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, gResult.Headers.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, gResult.Headers.Get(OssClientSideEncryptionUnencryptedContentMD5))

	assert.NotEqual(t, lastEtag, ToString(gResult.ETag))
}

func TestClientExtension(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)
	bucketNameNoExist := bucketName + "-no-exist"
	objectNameNoExist := objectName + "-no-exist"

	client := getDefaultClient()
	assert.NotNil(t, client)

	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	assert.NotNil(t, noPermClient)

	errorClient := getClientWithCredentialsProvider("", "",
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	assert.NotNil(t, errorClient)

	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)

	// IsBucketExist
	exist, err := client.IsBucketExist(context.TODO(), bucketName)
	assert.Nil(t, err)
	assert.True(t, exist)

	exist, err = client.IsBucketExist(context.TODO(), bucketNameNoExist)
	assert.Nil(t, err)
	assert.False(t, exist)

	exist, err = noPermClient.IsBucketExist(context.TODO(), bucketName)
	assert.Nil(t, err)
	assert.True(t, exist)

	exist, err = noPermClient.IsBucketExist(context.TODO(), bucketNameNoExist)
	assert.Nil(t, err)
	assert.False(t, exist)

	exist, err = errorClient.IsBucketExist(context.TODO(), bucketName)
	assert.NotNil(t, err)
	var serr *ServiceError
	assert.False(t, errors.As(err, &serr))

	// IsObjectExist
	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	assert.True(t, exist)

	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectNameNoExist)
	assert.Nil(t, err)
	assert.False(t, exist)

	exist, err = client.IsObjectExist(context.TODO(), bucketNameNoExist, objectName)
	assert.NotNil(t, err)
	assert.False(t, exist)
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "NoSuchBucket", serr.Code)

	exist, err = client.IsObjectExist(context.TODO(), bucketNameNoExist, objectNameNoExist)
	assert.NotNil(t, err)
	assert.False(t, exist)
	assert.NotNil(t, serr)
	assert.Equal(t, "NoSuchBucket", serr.Code)

	exist, err = noPermClient.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.NotNil(t, err)
	assert.False(t, exist)
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)

	exist, err = noPermClient.IsObjectExist(context.TODO(), bucketNameNoExist, objectName)
	assert.NotNil(t, err)
	assert.False(t, exist)
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "NoSuchBucket", serr.Code)

	exist, err = errorClient.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.NotNil(t, err)
	assert.False(t, exist)
	assert.False(t, errors.As(err, &serr))

	//PutObjectFromFile
	objectNameFromFile := objectName + "-from-file"
	var localFile = randStr(8) + ".txt"
	length := 1234
	content := randStr(length)
	hashContent := NewCRC64(0)
	hashContent.Write([]byte(content))
	createFile(t, localFile, content)
	defer func() { os.Remove(localFile) }()

	result, err := client.PutObjectFromFile(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameFromFile),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	gResult, err := client.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameFromFile),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, fmt.Sprint(hashContent.Sum64()), ToString(gResult.HashCRC64))
	_, err = io.ReadAll(gResult.Body)
	gResult.Body.Close()
	dumpErrIfNotNil(err)

	// Use Uploader, set meta and acl
	objectNameBig := objectName + "-big"
	bigLength := 5*100*1024 + 1234
	bigContent := randStr(bigLength)
	bigHash := NewCRC64(0)
	bigHash.Write([]byte(bigContent))
	u := client.NewUploader()
	assert.NotNil(t, u)
	urResult, err := u.UploadFrom(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectNameBig),
			Metadata: map[string]string{
				"author": "test",
				"magic":  "123",
			},
			Acl: ObjectACLPrivate,
		},
		bytes.NewReader([]byte(bigContent)),
		func(uo *UploaderOptions) {
			uo.ParallelNum = 3
			uo.PartSize = 100 * 1024
		},
	)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.NotNil(t, urResult)

	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectNameBig)
	assert.Nil(t, err)
	assert.True(t, exist)

	hResult, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameBig),
	})
	assert.Nil(t, err)
	assert.NotNil(t, hResult)
	assert.Contains(t, hResult.Headers.Get(HTTPHeaderETag), "-6")
	assert.Equal(t, "Multipart", hResult.Headers.Get(HeaderOssObjectType))
	assert.Equal(t, "test", hResult.Headers.Get("x-oss-meta-author"))
	assert.Equal(t, "123", hResult.Headers.Get("x-oss-meta-magic"))

	aclResult, err := client.GetObjectAcl(context.TODO(), &GetObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameBig),
	})
	assert.Nil(t, err)
	assert.NotNil(t, hResult)
	assert.Equal(t, "private", ToString(aclResult.ACL))

	// Downloader with not align partSize
	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 100*1024 + 123
	})
	assert.NotNil(t, d)
	assert.Equal(t, int64(100*1024+123), d.options.PartSize)
	assert.Equal(t, 3, d.options.ParallelNum)
	localFileBig := randStr(8) + "-downloader"
	dResult, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectNameBig)},
		localFileBig)

	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, int64(bigLength), dResult.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFileBig)
	assert.Nil(t, err)
	defer func() {
		rfile.Close()
		os.Remove(localFileBig)
	}()
	io.Copy(hash, rfile)
	assert.Equal(t, bigHash.Sum64(), hash.Sum64())

	// Downloader with not align partSize + write-bufer-size
	d2 := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 100*1024 + 123
		do.WriteBufferSize = 16 * 1024
	})
	assert.NotNil(t, d2)
	assert.Equal(t, int64(100*1024+123), d.options.PartSize)
	assert.Equal(t, 3, d2.options.ParallelNum)
	assert.Equal(t, 16*1024, d2.options.WriteBufferSize)
	localFileBig2 := randStr(8) + "-downloader-2"
	dResult2, err := d2.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectNameBig)},
		localFileBig2)

	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, int64(bigLength), dResult2.Written)

	hash2 := NewCRC64(0)
	rfile2, err := os.Open(localFileBig2)
	assert.Nil(t, err)
	defer func() {
		rfile.Close()
		os.Remove(localFileBig2)
	}()
	io.Copy(hash2, rfile2)
	assert.Equal(t, bigHash.Sum64(), hash2.Sum64())

	//Use ReadOnlyFile
	f, err := client.OpenFile(context.TODO(), bucketName, objectNameBig)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	for i := 13; i < 42; i++ {
		for len := 100*1024 + 123; len < 100*1024+123+17; len++ {
			_, err := f.Seek(int64(i), io.SeekStart)
			assert.Nil(t, err)
			gData, err := io.ReadAll(io.LimitReader(f, int64(len)))
			assert.Nil(t, err)
			assert.EqualValues(t, []byte(bigContent)[i:i+len], gData)
		}
	}
	f.Close()

	// AppenableFile
	objectNameAppend := objectName + "-append"
	dataa1 := []byte("helle world")
	dataa2 := []byte(randStr(12345))
	dataa3 := []byte(randStr(100*1024*5 + 13))
	var localFileData3 = randStr(8) + ".txt"
	createFile(t, localFileData3, string(dataa3))
	defer func() {
		os.Remove(localFileData3)
	}()

	af, err := client.AppendFile(context.TODO(), bucketName, objectNameAppend)
	n, err := af.Write(dataa1)
	assert.Nil(t, err)
	assert.Equal(t, len(dataa1), n)

	hResult, err = client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend),
	})
	assert.Nil(t, err)
	assert.NotNil(t, hResult)
	assert.Equal(t, int64(n), hResult.ContentLength)
	nl, err := af.WriteFrom(bytes.NewReader(dataa2))
	assert.Equal(t, int64(len(dataa2)), nl)

	filedataa3, err := os.Open(localFileData3)
	assert.Nil(t, err)
	nl, err = io.Copy(af, filedataa3)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(dataa3)), nl)
	defer func() {
		filedataa3.Close()
	}()

	af.Close()
	hashA := NewCRC64(0)
	hashA.Write(dataa1)
	hashA.Write(dataa2)
	hashA.Write(dataa3)
	hResult, err = client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend),
	})
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(hashA.Sum64()), ToString(hResult.HashCRC64))

	//GetObjectToFile
	var localFileToFile = randStr(8) + "-to-file"
	defer func() {
		os.Remove(localFileToFile)
	}()
	gResult, err = client.GetObjectToFile(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend)},
		localFileToFile,
	)
	assert.Nil(t, err)
	gResult, err = client.GetObjectToFile(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend),
		ProgressFn: func(increment, transferred, total int64) {
			//fmt.Printf("increment:%v, transferred:%v, total:%v\n", increment, transferred, total)
		}},
		localFileToFile,
	)
	assert.Nil(t, err)
	hash = NewCRC64(0)
	rfiletoFile, err := os.Open(localFileToFile)
	assert.Nil(t, err)
	defer func() {
		rfiletoFile.Close()
	}()
	io.Copy(hash, rfiletoFile)
	assert.Equal(t, hashA.Sum64(), hash.Sum64())
}

func TestGetObjectToFileV2(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	content := randStr(3*1024*1024 + 1234)
	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	})
	assert.Nil(t, err)

	// Compute expected CRC64
	expectedHash := NewCRC64(0)
	expectedHash.Write([]byte(content))
	expectedCRC := expectedHash.Sum64()

	// 1. Basic full download
	localFile := randStr(8) + "-v2-integ"
	defer os.Remove(localFile)

	result, err := client.GetObjectToFileV2(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
		},
		localFile,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, ToString(result.ETag))
	assert.NotEmpty(t, ToString(result.HashCRC64))

	downloaded, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, len(content), len(downloaded))
	assert.Equal(t, content, string(downloaded))

	fileHash := NewCRC64(0)
	fileHash.Write(downloaded)
	assert.Equal(t, expectedCRC, fileHash.Sum64())

	// 2. Full download with progress
	localFileProgress := randStr(8) + "-v2-integ-prog"
	defer os.Remove(localFileProgress)

	var lastTransferred, lastTotal int64
	var progressCalled bool
	result, err = client.GetObjectToFileV2(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
			ProgressFn: func(increment, transferred, total int64) {
				progressCalled = true
				lastTransferred = transferred
				lastTotal = total
			},
		},
		localFileProgress,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.True(t, progressCalled)
	assert.Equal(t, int64(len(content)), lastTransferred)
	assert.Equal(t, int64(len(content)), lastTotal)

	downloaded, err = os.ReadFile(localFileProgress)
	assert.Nil(t, err)
	assert.Equal(t, content, string(downloaded))

	// 3. Full download with writeBufferSize
	localFileBuf := randStr(8) + "-v2-integ-buf"
	defer os.Remove(localFileBuf)

	bufSize := 64 * 1024
	result, err = client.GetObjectToFileV2(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
		},
		localFileBuf,
		&bufSize,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	downloaded, err = os.ReadFile(localFileBuf)
	assert.Nil(t, err)
	assert.Equal(t, content, string(downloaded))

	// 4. Range download — first 1024 bytes
	localFileRange := randStr(8) + "-v2-integ-range"
	defer os.Remove(localFileRange)

	result, err = client.GetObjectToFileV2(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
			Range:  Ptr("bytes=0-1023"),
		},
		localFileRange,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	downloaded, err = os.ReadFile(localFileRange)
	assert.Nil(t, err)
	assert.Equal(t, 1024, len(downloaded))
	assert.Equal(t, content[:1024], string(downloaded))

	// 5. Range download — middle range
	localFileRange2 := randStr(8) + "-v2-integ-range2"
	defer os.Remove(localFileRange2)

	result, err = client.GetObjectToFileV2(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
			Range:  Ptr("bytes=1024-2047"),
		},
		localFileRange2,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	downloaded, err = os.ReadFile(localFileRange2)
	assert.Nil(t, err)
	assert.Equal(t, 1024, len(downloaded))
	assert.Equal(t, content[1024:2048], string(downloaded))

	// 6. Range download with progress
	localFileRangeProg := randStr(8) + "-v2-integ-range-prog"
	defer os.Remove(localFileRangeProg)

	progressCalled = false
	lastTransferred = 0
	result, err = client.GetObjectToFileV2(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
			Range:  Ptr("bytes=0-1023"),
			ProgressFn: func(increment, transferred, total int64) {
				progressCalled = true
				lastTransferred = transferred
			},
		},
		localFileRangeProg,
		nil,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.True(t, progressCalled)
	assert.Equal(t, int64(1024), lastTransferred)

	// 7. Nil request
	_, err = client.GetObjectToFileV2(context.TODO(), nil, "test.txt", nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request")

	// 8. Non-existent bucket
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	localFileErr := randStr(8) + "-v2-integ-err"
	defer os.Remove(localFileErr)

	_, err = client.GetObjectToFileV2(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketNameNotExist),
			Key:    Ptr(objectName),
		},
		localFileErr,
		nil,
	)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, 404, serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)

	// 9. Non-existent key
	_, err = client.GetObjectToFileV2(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectNamePrefix + "not-exist-key"),
		},
		localFileErr,
		nil,
	)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 404, serr.StatusCode)
	assert.Equal(t, "NoSuchKey", serr.Code)
}

func testClientExtensionWithPayer(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)
	client := getDefaultClient()
	assert.NotNil(t, client)

	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	policyInfo := `
	{
		"Version":"1",
		"Statement":[
			{
				"Action":[
					"oss:*"
				],
				"Effect":"Allow",
				"Principal":["` + payerUID_ + `"],
				"Resource":["acs:oss:*:*:` + bucketName + `", "acs:oss:*:*:` + bucketName + `/*"]
			}
		]
	}`
	input := &OperationInput{
		OpName: "PutBucketPolicy",
		Bucket: Ptr(bucketName),
		Method: "PUT",
		Parameters: map[string]string{
			"policy": "",
		},
		Body: strings.NewReader(policyInfo),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	_, err = client.InvokeOperation(context.TODO(), input)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	request := &PutBucketRequestPaymentRequest{
		Bucket: Ptr(bucketName),
		PaymentConfiguration: &RequestPaymentConfiguration{
			Payer: Requester,
		},
	}
	_, err = client.PutBucketRequestPayment(context.TODO(), request)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	creClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider(payerAccessID_, payerAccessKey_))

	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)

	// IsObjectExist
	exist, err := creClient.IsObjectExist(context.TODO(), bucketName, objectName, func(op *IsObjectExistOptions) {
		op.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	assert.True(t, exist)

	//PutObjectFromFile
	objectNameFromFile := objectName + "-from-file"
	var localFile = randStr(8) + ".txt"
	length := 1234
	content := randStr(length)
	hashContent := NewCRC64(0)
	hashContent.Write([]byte(content))
	createFile(t, localFile, content)
	defer func() { os.Remove(localFile) }()

	result, err := creClient.PutObjectFromFile(context.TODO(), &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectNameFromFile),
		RequestPayer: Ptr("requester"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Use Uploader, set meta and acl
	objectNameBig := objectName + "-big"
	bigLength := 5*100*1024 + 1234
	bigContent := randStr(bigLength)
	bigHash := NewCRC64(0)
	bigHash.Write([]byte(bigContent))
	u := creClient.NewUploader()
	assert.NotNil(t, u)
	urResult, err := u.UploadFrom(context.TODO(),
		&PutObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectNameBig),
			Metadata: map[string]string{
				"author": "test",
				"magic":  "123",
			},
			Acl:          ObjectACLPrivate,
			RequestPayer: Ptr("requester"),
		},
		bytes.NewReader([]byte(bigContent)),
		func(uo *UploaderOptions) {
			uo.ParallelNum = 3
			uo.PartSize = 100 * 1024
		},
	)
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.NotNil(t, urResult)

	exist, err = creClient.IsObjectExist(context.TODO(), bucketName, objectNameBig, func(op *IsObjectExistOptions) {
		op.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	assert.True(t, exist)

	// Downloader with not align partSize
	d := creClient.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 100*1024 + 123
	})
	assert.NotNil(t, d)
	assert.Equal(t, int64(100*1024+123), d.options.PartSize)
	assert.Equal(t, 3, d.options.ParallelNum)
	localFileBig := randStr(8) + "-downloader"
	dResult, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket:       Ptr(bucketName),
			Key:          Ptr(objectNameBig),
			RequestPayer: Ptr("requester"),
		},
		localFileBig)

	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, int64(bigLength), dResult.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFileBig)
	assert.Nil(t, err)
	defer func() {
		rfile.Close()
		os.Remove(localFileBig)
	}()
	io.Copy(hash, rfile)
	assert.Equal(t, bigHash.Sum64(), hash.Sum64())

	//Use ReadOnlyFile
	f, err := creClient.OpenFile(context.TODO(), bucketName, objectNameBig, func(op *OpenOptions) {
		op.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	for i := 13; i < 42; i++ {
		for len := 100*1024 + 123; len < 100*1024+123+17; len++ {
			_, err := f.Seek(int64(i), io.SeekStart)
			assert.Nil(t, err)
			gData, err := io.ReadAll(io.LimitReader(f, int64(len)))
			assert.Nil(t, err)
			assert.EqualValues(t, []byte(bigContent)[i:i+len], gData)
		}
	}
	f.Close()

	// AppenableFile
	objectNameAppend := objectName + "-append"
	dataa3 := []byte(randStr(100*1024*5 + 13))
	var localFileData3 = randStr(8) + ".txt"
	createFile(t, localFileData3, string(dataa3))
	defer func() {
		os.Remove(localFileData3)
	}()

	af, err := creClient.AppendFile(context.TODO(), bucketName, objectNameAppend, func(op *AppendOptions) {
		op.RequestPayer = Ptr("requester")
	})
	assert.Nil(t, err)
	_, err = af.Write([]byte(content))
	assert.Nil(t, err)
	_, err = af.WriteFrom(strings.NewReader(content))
	assert.Nil(t, err)
	_, err = af.Stat()
	assert.Nil(t, err)
	//GetObjectToFile
	var localFileToFile = randStr(8) + "-to-file"
	defer func() {
		os.Remove(localFileToFile)
	}()
	_, err = creClient.GetObjectToFile(context.TODO(), &GetObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectNameAppend),
		RequestPayer: Ptr("requester"),
	},
		localFileToFile,
	)
	assert.Nil(t, err)
	_, err = creClient.GetObjectToFile(context.TODO(), &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameAppend),
	},
		localFileToFile,
	)
	assert.NotNil(t, err)
}

func TestClientAppendFile(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6) + ".append"

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	f, err := client.AppendFile(context.TODO(), bucketName, objectName, func(ao *AppendOptions) {
		ao.CreateParameter = &AppendObjectRequest{
			Acl:          ObjectACLPrivate,
			CacheControl: Ptr("no-cache"),
			Metadata: map[string]string{
				"user": "jack",
			},
			Tagging: Ptr("key=value"),
		}
	})
	assert.Nil(t, err)

	n, err := f.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, int(5), n)

	n, err = f.Write([]byte(" world"))
	assert.Nil(t, err)
	assert.Equal(t, int(6), n)

	nn, err := f.WriteFrom(strings.NewReader(" 123"))
	assert.Nil(t, err)
	assert.Equal(t, int64(4), nn)

	info, err := f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(15), info.Size())
	header, ok := info.Sys().(http.Header)
	assert.True(t, ok)

	assert.Equal(t, "no-cache", header.Get("Cache-Control"))
	assert.Equal(t, "jack", header.Get("x-oss-meta-user"))
	assert.Equal(t, "1", header.Get("x-oss-tagging-count"))

	f.Close()

	//Open Again
	f, err = client.AppendFile(context.TODO(), bucketName, objectName, func(ao *AppendOptions) {
		ao.CreateParameter = &AppendObjectRequest{
			Acl:          ObjectACLPrivate,
			CacheControl: Ptr("no-cache"),
			Metadata: map[string]string{
				"user": "jack",
			},
			Tagging: Ptr("key=value"),
		}
	})
	n, err = f.Write([]byte("abc"))
	assert.Nil(t, err)
	info, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(18), info.Size())
	f.Close()

	objectName1 := objectNamePrefix + randLowStr(6) + "-1-.append"
	f, err = client.AppendFile(context.TODO(), bucketName, objectName1, func(ao *AppendOptions) {
		ao.CreateParameter = &AppendObjectRequest{
			Acl:          ObjectACLPrivate,
			CacheControl: Ptr("no-cache"),
			Metadata: map[string]string{
				"user": "jack-1",
			},
			Tagging: Ptr("key1=value1"),
		}
	})
	assert.Nil(t, err)

	nn, err = f.WriteFrom(strings.NewReader("123"))
	assert.Nil(t, err)
	assert.Equal(t, int64(3), nn)

	nn, err = f.WriteFrom(strings.NewReader("-abc-321"))
	assert.Nil(t, err)
	assert.Equal(t, int64(8), nn)
	info, err = f.Stat()
	assert.Nil(t, err)
	assert.Equal(t, int64(11), info.Size())
	header, ok = info.Sys().(http.Header)
	assert.True(t, ok)

	assert.Equal(t, "no-cache", header.Get("Cache-Control"))
	assert.Equal(t, "jack-1", header.Get("x-oss-meta-user"))
	assert.Equal(t, "1", header.Get("x-oss-tagging-count"))
	f.Close()
}

func TestDownloaderTruncate(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	objectLen := 100*1024*4 + 123
	content := randLowStr(objectLen)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	result, err := client.PutObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 100 * 1024
		do.UseTempFile = false
	})
	assert.NotNil(t, d)
	assert.Equal(t, 3, d.options.ParallelNum)
	localFile := randStr(8) + "-downloader"
	defer os.Remove(localFile)

	// file size is 0, file size < get size
	dResult, err := d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
		},
		localFile)

	assert.Nil(t, err)

	assert.Equal(t, objectLen, (int)(dResult.Written))
	gotCotent, err := os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, []byte(content), gotCotent)

	//range get
	// file size > get size
	dResult, err = d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
			Range:  Ptr("bytes=123-123456"),
		},
		localFile)

	gotCotent, err = os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, 123456-123+1, (int)(dResult.Written))
	assert.Equal(t, []byte(content[123:123456+1]), []byte(gotCotent))

	// file size == get size
	content1 := randLowStr(123456 - 123 + 1)
	os.WriteFile(localFile, []byte(content1), 0644)
	gotCotent, err = os.ReadFile(localFile)
	assert.Equal(t, []byte(content1), []byte(gotCotent))
	dResult, err = d.DownloadFile(context.TODO(),
		&GetObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr(objectName),
			Range:  Ptr("bytes=123-123456"),
		},
		localFile)

	gotCotent, err = os.ReadFile(localFile)
	assert.Nil(t, err)
	assert.Equal(t, 123456-123+1, (int)(dResult.Written))
	assert.Equal(t, []byte(content[123:123456+1]), []byte(gotCotent))
}

func TestUploaderWithSequential(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	objectLen := 100*1024*4 + 123
	content := randLowStr(objectLen)
	partSize := int64(100 * 1024)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.PartSize = partSize
		},
	)

	result, err := u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}, strings.NewReader(content))
	assert.Nil(t, err)
	assert.NotNil(t, result)

	request := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	headResult, err := client.HeadObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, headResult.ContentMD5)

	result, err = u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"sequential": "",
			},
		},
	}, strings.NewReader(content))
	assert.Nil(t, err)
	assert.NotNil(t, result)

	request = &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	headResult, err = client.HeadObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.NotEmpty(t, headResult.ContentMD5)
}

func TestUploaderUploadFromTeeReader(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	partSize := int64(100 * 1024)

	u := NewUploader(client,
		func(uo *UploaderOptions) {
			uo.PartSize = partSize
		},
	)

	// empty reader
	content := strings.Repeat("a", 0)
	reader := io.TeeReader(strings.NewReader(content), io.Discard)

	result, err := u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}, reader)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	request := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	headResult, err := client.HeadObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), headResult.ContentLength)

	// 2 * part-size
	content = strings.Repeat("a", 200*1024)
	reader = io.TeeReader(strings.NewReader(content), io.Discard)
	result, err = u.UploadFrom(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}, strings.NewReader(content))
	assert.Nil(t, err)
	assert.NotNil(t, result)

	request = &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	headResult, err = client.HeadObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, int64(200*1024), headResult.ContentLength)
}

func TestCopierWithNoCheckSSEFlags(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	length := 2*100*1024 + 1234
	content := randStr(length)
	hash := NewCRC64(0)
	hash.Write([]byte(content))

	request := &PutObjectRequest{
		Bucket:               Ptr(bucketName),
		Key:                  Ptr(objectName),
		ServerSideEncryption: Ptr("AES256"),
		Body:                 strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	dstObjectName := objectName + "-copy"
	copyRequest := &CopyObjectRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(dstObjectName),
		SourceKey: Ptr(objectName),
	}

	c := client.NewCopier(func(co *CopierOptions) {
		co.ParallelNum = 1
		co.PartSize = 100 * 1024
		co.MultipartCopyThreshold = 100 * 1024
	})
	result, err := c.Copy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	headResult, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(dstObjectName),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Multipart", ToString(headResult.ObjectType))

	// Disable SSE Check From Copy
	result, err = c.Copy(context.TODO(), copyRequest, WithCopierNoCheckSSE(true))
	assert.Nil(t, err)
	assert.NotNil(t, result)
	headResult, err = client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(dstObjectName),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Normal", ToString(headResult.ObjectType))

	// Disable SSE Check From Copier
	dstObjectName1 := objectName + "-copy-client"
	c1 := client.NewCopier(func(co *CopierOptions) {
		co.ParallelNum = 1
		co.PartSize = 100 * 1024
		co.MultipartCopyThreshold = 100 * 1024
	},
		WithCopierNoCheckSSE(true),
	)
	result1, err := c1.Copy(context.TODO(), &CopyObjectRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(dstObjectName1),
		SourceKey: Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result1)

	headResult1, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(dstObjectName1),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Normal", ToString(headResult1.ObjectType))

}

func TestCopierWithNoCheckCrossBucketFlags(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	crossBucketName := bucketName + "-cross"

	//TODO
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	_, err = client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(crossBucketName),
	})
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	length := 2*100*1024 + 1234
	content := randStr(length)
	hash := NewCRC64(0)
	hash.Write([]byte(content))

	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	dstObjectName := objectName + "-copy"
	copyRequest := &CopyObjectRequest{
		Bucket:       Ptr(crossBucketName),
		Key:          Ptr(dstObjectName),
		SourceBucket: Ptr(bucketName),
		SourceKey:    Ptr(objectName),
	}

	c := client.NewCopier(func(co *CopierOptions) {
		co.ParallelNum = 1
		co.PartSize = 100 * 1024
		co.MultipartCopyThreshold = 100 * 1024
	})
	result, err := c.Copy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	headResult, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(crossBucketName),
		Key:    Ptr(dstObjectName),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Multipart", ToString(headResult.ObjectType))

	// Disable Check From Copy
	result, err = c.Copy(context.TODO(), copyRequest, WithCopierNoCheckCrossBucket(true))
	assert.Nil(t, err)
	assert.NotNil(t, result)
	headResult, err = client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(crossBucketName),
		Key:    Ptr(dstObjectName),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Normal", ToString(headResult.ObjectType))

	// Disable Check From Copier
	dstObjectName1 := objectName + "-copy-client"
	c1 := client.NewCopier(func(co *CopierOptions) {
		co.ParallelNum = 1
		co.PartSize = 100 * 1024
		co.MultipartCopyThreshold = 100 * 1024
	},
		WithCopierNoCheckCrossBucket(true),
	)
	result1, err := c1.Copy(context.TODO(), &CopyObjectRequest{
		Bucket:       Ptr(crossBucketName),
		Key:          Ptr(dstObjectName1),
		SourceBucket: Ptr(bucketName),
		SourceKey:    Ptr(objectName),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result1)

	headResult1, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
		Bucket: Ptr(crossBucketName),
		Key:    Ptr(dstObjectName1),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Normal", ToString(headResult1.ObjectType))
}

func TestIntegrationWriteOnlyFileMultiPart(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{Bucket: Ptr(bucketName)})
	assert.Nil(t, err)

	partSize := int64(200 * 1024)
	length := int(partSize)*3 + 12345 // 3 full parts + a short last part
	data := []byte(randStr(length))
	wantCRC := NewCRC64(0)
	wantCRC.Write(data)

	f, err := client.NewWriteOnlyFile(context.TODO(), bucketName, objectName, func(o *WriteOnlyOptions) {
		o.PartSize = partSize
		o.ParallelNum = 3
		// Acl has no Initiate equivalent and is applied at Complete; Metadata
		// flows through Initiate. Verify both survive the multipart path.
		// (public ACLs are blocked by the test account's Block Public Access.)
		o.CreateParameter = &PutObjectRequest{
			Acl:      ObjectACLPrivate,
			Metadata: map[string]string{"author": "test"},
		}
	})
	assert.Nil(t, err)

	// write in 64KB chunks to exercise buffering across Write calls
	for off := 0; off < length; off += 64 * 1024 {
		end := off + 64*1024
		if end > length {
			end = length
		}
		_, err = f.Write(data[off:end])
		assert.Nil(t, err)
	}
	assert.Nil(t, f.Close())

	// download and verify content + CRC
	gr, err := client.GetObject(context.TODO(), &GetObjectRequest{Bucket: Ptr(bucketName), Key: Ptr(objectName)})
	assert.Nil(t, err)
	got, err := io.ReadAll(gr.Body)
	gr.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, length, len(got))
	gotCRC := NewCRC64(0)
	gotCRC.Write(got)
	assert.Equal(t, wantCRC.Sum64(), gotCRC.Sum64())
	assert.Equal(t, "test", gr.Metadata["author"])

	// Acl applied at Complete
	aclResult, err := client.GetObjectAcl(context.TODO(), &GetObjectAclRequest{Bucket: Ptr(bucketName), Key: Ptr(objectName)})
	assert.Nil(t, err)
	assert.Equal(t, string(ObjectACLPrivate), ToString(aclResult.ACL))
}

func TestIntegrationWriteOnlyFileResume(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{Bucket: Ptr(bucketName)})
	assert.Nil(t, err)

	partSize := int64(200 * 1024)
	length := int(partSize) * 4
	data := []byte(randStr(length))
	wantCRC := NewCRC64(0)
	wantCRC.Write(data)

	// session 1: write 2 full parts, then "crash" (no Close)
	f1, err := client.NewWriteOnlyFile(context.TODO(), bucketName, objectName, func(o *WriteOnlyOptions) {
		o.PartSize = partSize
		o.ParallelNum = 1 // serial so committedOffset advances deterministically
	})
	assert.Nil(t, err)
	_, err = f1.Write(data[:partSize*2])
	assert.Nil(t, err)
	// ensure the two parts are durably uploaded before snapshotting
	assert.Nil(t, f1.Flush())
	cp := f1.StatCheckpoint()
	assert.Equal(t, int64(partSize*2), cp.Offset)
	assert.NotEqual(t, "", cp.UploadId)

	// session 2: resume from checkpoint, replay from offset, finish
	f2, err := client.OpenWriteOnlyFile(context.TODO(), bucketName, objectName, cp)
	assert.Nil(t, err)
	assert.Equal(t, int64(partSize*2), f2.StatCheckpoint().Offset)
	_, err = f2.Write(data[cp.Offset:])
	assert.Nil(t, err)
	assert.Nil(t, f2.Close())

	gr, err := client.GetObject(context.TODO(), &GetObjectRequest{Bucket: Ptr(bucketName), Key: Ptr(objectName)})
	assert.Nil(t, err)
	got, err := io.ReadAll(gr.Body)
	gr.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, length, len(got))
	gotCRC := NewCRC64(0)
	gotCRC.Write(got)
	assert.Equal(t, wantCRC.Sum64(), gotCRC.Sum64())
}

func TestIntegrationWriteOnlyFileResumeNoSuchUpload(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{Bucket: Ptr(bucketName)})
	assert.Nil(t, err)

	// non-existent uploadId -> resume is refused with ErrCheckpointInvalid
	cp := WriteCheckpoint{UploadId: "non-existent-upload-id", Offset: 400 * 1024, PartSize: 200 * 1024}
	f, err := client.OpenWriteOnlyFile(context.TODO(), bucketName, objectName, cp)
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrCheckpointInvalid)

	// recommended recovery: discard the stale checkpoint and start fresh
	f, err = client.NewWriteOnlyFile(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	_, err = f.Write([]byte("recovered"))
	assert.Nil(t, err)
	assert.Nil(t, f.Close())

	res, err := client.GetObject(context.TODO(), &GetObjectRequest{Bucket: Ptr(bucketName), Key: Ptr(objectName)})
	assert.Nil(t, err)
	got2, _ := io.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, []byte("recovered"), got2)
}

func TestIntegrationWriteOnlyFileResumeInvalidCheckpoint(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{Bucket: Ptr(bucketName)})
	assert.Nil(t, err)

	partSize := int64(200 * 1024)
	data := []byte(randStr(int(partSize) * 2))

	// establish a real upload with two durable parts, snapshot a valid checkpoint
	f1, err := client.NewWriteOnlyFile(context.TODO(), bucketName, objectName, func(o *WriteOnlyOptions) {
		o.PartSize = partSize
		o.ParallelNum = 1
	})
	assert.Nil(t, err)
	_, err = f1.Write(data)
	assert.Nil(t, err)
	assert.Nil(t, f1.Flush())
	cp := f1.StatCheckpoint()
	assert.Equal(t, int64(partSize*2), cp.Offset)
	assert.NotEqual(t, "", cp.UploadId)

	// invalid input: Offset disagrees with server truth -> refused
	badOffset := cp
	badOffset.Offset = partSize // server actually holds two parts (partSize*2)
	f, err := client.OpenWriteOnlyFile(context.TODO(), bucketName, objectName, badOffset)
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrCheckpointInvalid)

	// invalid input: resume token present but PartSize missing -> refused
	badPartSize := WriteCheckpoint{UploadId: cp.UploadId, Offset: cp.Offset, PartSize: 0}
	f, err = client.OpenWriteOnlyFile(context.TODO(), bucketName, objectName, badPartSize)
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrCheckpointInvalid)

	// the untampered checkpoint still resumes cleanly, proving only the bad ones were rejected
	f2, err := client.OpenWriteOnlyFile(context.TODO(), bucketName, objectName, cp)
	assert.Nil(t, err)
	assert.Equal(t, int64(partSize*2), f2.StatCheckpoint().Offset)
	assert.Nil(t, f2.AbortClose()) // discard the multipart upload
}

func TestIntegrationWriteOnlyFileInvalidCredentials(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6)

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{Bucket: Ptr(bucketName)})
	assert.Nil(t, err)

	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	assert.NotNil(t, noPermClient)

	var serr *ServiceError

	// multipart path: sealing a full part triggers InitiateMultipartUpload, whose
	// auth failure surfaces directly from Write.
	fm, err := noPermClient.NewWriteOnlyFile(context.TODO(), bucketName, objectName, func(o *WriteOnlyOptions) {
		o.PartSize = 8
		o.ParallelNum = 2
	})
	assert.Nil(t, err)
	_, err = fm.Write([]byte("AAAAAAAABBBB")) // seals part 1 -> Initiate -> auth failure
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Nil(t, fm.AbortClose()) // never initiated: no upload to abort

	// small-object path: no full part sealed, so Close falls back to PutObject,
	// which fails auth.
	serr = nil
	fs, err := noPermClient.NewWriteOnlyFile(context.TODO(), bucketName, objectName, func(o *WriteOnlyOptions) {
		o.PartSize = 1024 * 1024
	})
	assert.Nil(t, err)
	_, err = fs.Write([]byte("small"))
	assert.Nil(t, err)
	err = fs.Close()
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
}
