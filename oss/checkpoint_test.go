package oss

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDownloadCheckpoint(t *testing.T) {
	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}
	destFilePath := randStr(8) + "-no-surfix"
	cpDir := "."
	header := http.Header{
		"Etag":           {"\"D41D8CD98F00B204E9800998ECF8****\""},
		"Content-Length": {"344606"},
		"Last-Modified":  {"Fri, 24 Feb 2012 06:07:48 GMT"},
	}
	partSize := DefaultDownloadPartSize

	cp := newDownloadCheckpoint(request, destFilePath, cpDir, header, partSize)
	assert.NotNil(t, cp)
	assert.Equal(t, "\"D41D8CD98F00B204E9800998ECF8****\"", cp.Info.Data.ObjectMeta.ETag)
	assert.Equal(t, "Fri, 24 Feb 2012 06:07:48 GMT", cp.Info.Data.ObjectMeta.LastModified)
	assert.Equal(t, int64(344606), cp.Info.Data.ObjectMeta.Size)

	assert.Equal(t, "oss://bucket/key", cp.Info.Data.ObjectInfo.Name)
	assert.Equal(t, "", cp.Info.Data.ObjectInfo.VersionId)
	assert.Equal(t, "", cp.Info.Data.ObjectInfo.Range)

	assert.Equal(t, CheckpointMagic, cp.Info.Magic)
	assert.Equal(t, "", cp.Info.MD5)

	assert.Equal(t, destFilePath, cp.Info.Data.FilePath)
	assert.Equal(t, partSize, cp.Info.Data.PartSize)

	//has version id
	request = &GetObjectRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("key"),
		VersionId: Ptr("id"),
	}
	cp_vid := newDownloadCheckpoint(request, destFilePath, cpDir, header, partSize)
	assert.NotNil(t, cp_vid)
	assert.Equal(t, "oss://bucket/key", cp_vid.Info.Data.ObjectInfo.Name)
	assert.Equal(t, "id", cp_vid.Info.Data.ObjectInfo.VersionId)
	assert.Equal(t, "", cp_vid.Info.Data.ObjectInfo.Range)

	//has range
	request = &GetObjectRequest{
		Bucket:    Ptr("bucket"),
		Key:       Ptr("key"),
		VersionId: Ptr("id"),
		Range:     Ptr("bytes=1-10"),
	}
	cp_range := newDownloadCheckpoint(request, destFilePath, cpDir, header, partSize)
	assert.NotNil(t, cp_range)
	assert.Equal(t, "oss://bucket/key", cp_range.Info.Data.ObjectInfo.Name)
	assert.Equal(t, "id", cp_range.Info.Data.ObjectInfo.VersionId)
	assert.Equal(t, "bytes=1-10", cp_range.Info.Data.ObjectInfo.Range)

	assert.NotEqual(t, cp.CpFilePath, cp_vid.CpFilePath)
	assert.NotEqual(t, cp.CpFilePath, cp_range.CpFilePath)
	assert.NotEqual(t, cp_vid.CpFilePath, cp_range.CpFilePath)

	// with other destFilePath
	destFilePath1 := destFilePath + "-123"
	cp_range_dest := newDownloadCheckpoint(request, destFilePath1, cpDir, header, partSize)
	assert.NotNil(t, cp_range_dest)
	assert.NotEqual(t, cp_range.CpFilePath, cp_range_dest.CpFilePath)
	assert.Equal(t, destFilePath1, cp_range_dest.Info.Data.FilePath)

	//check dump
	cp.dump()
	assert.True(t, FileExists(cp.CpFilePath))
	dcp := downloadCheckpoint{}
	content, err := os.ReadFile(cp.CpFilePath)
	assert.Nil(t, err)
	err = json.Unmarshal(content, &dcp.Info)
	assert.Nil(t, err)

	assert.Equal(t, "\"D41D8CD98F00B204E9800998ECF8****\"", dcp.Info.Data.ObjectMeta.ETag)
	assert.Equal(t, "Fri, 24 Feb 2012 06:07:48 GMT", dcp.Info.Data.ObjectMeta.LastModified)
	assert.Equal(t, int64(344606), dcp.Info.Data.ObjectMeta.Size)

	assert.Equal(t, "oss://bucket/key", dcp.Info.Data.ObjectInfo.Name)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.VersionId)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.Range)

	assert.Equal(t, CheckpointMagic, dcp.Info.Magic)
	assert.Equal(t, 32, len(dcp.Info.MD5))

	assert.Equal(t, destFilePath, dcp.Info.Data.FilePath)
	assert.Equal(t, partSize, dcp.Info.Data.PartSize)

	//check load
	err = cp.load()
	assert.Nil(t, err)
	assert.True(t, cp.Loaded)

	//check valid
	assert.True(t, cp.valid())

	//check complete
	assert.True(t, FileExists(cp.CpFilePath))
	err = cp.remove()
	assert.Nil(t, err)
	assert.True(t, !FileExists(cp.CpFilePath))

	//load not match
	cp = newDownloadCheckpoint(request, destFilePath, cpDir, header, partSize)
	assert.False(t, cp.Loaded)
	notMatch := `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"2f132b5bf65640868a47cb52c57492c8","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":5242880,"DownloadInfo":{"Offset":5242880,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(notMatch), FilePermMode)
	assert.Nil(t, err)
	assert.True(t, FileExists(cp.CpFilePath))
	err = cp.load()
	assert.Nil(t, err)
	assert.False(t, cp.Loaded)
	assert.False(t, FileExists(cp.CpFilePath))
}

func TestDownloadCheckpointInvalidCpPath(t *testing.T) {
	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}
	destFilePath := randStr(8) + "-no-surfix"
	cpDir := "./invliad-dir/"
	header := http.Header{
		"Etag":           {"\"D41D8CD98F00B204E9800998ECF8****\""},
		"Content-Length": {"344606"},
		"Last-Modified":  {"Fri, 24 Feb 2012 06:07:48 GMT"},
	}
	partSize := DefaultDownloadPartSize

	cp := newDownloadCheckpoint(request, destFilePath, cpDir, header, partSize)
	assert.NotNil(t, cp)
	assert.Equal(t, destFilePath, cp.Info.Data.FilePath)
	assert.Equal(t, "invliad-dir", cp.CpDirPath)
	assert.Contains(t, cp.CpFilePath, "invliad-dir")

	//dump fail
	err := cp.dump()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "The system cannot find the path specified") || strings.Contains(err.Error(), "no such file or directory"))

	//load fail
	err = cp.load()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invaid checkpoint dir")
}

func TestDownloadCheckpointValid(t *testing.T) {
	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}
	destFilePath := "gthnjXGQ-no-surfix"
	cpDir := "."
	header := http.Header{
		"Etag":           {"\"D41D8CD98F00B204E9800998ECF8****\""},
		"Content-Length": {"344606"},
		"Last-Modified":  {"Fri, 24 Feb 2012 06:07:48 GMT"},
	}
	partSize := int64(5 * 1024 * 1024) //DefaultDownloadPartSize
	cp := newDownloadCheckpoint(request, destFilePath, cpDir, header, partSize)

	os.Remove(destFilePath)
	assert.Equal(t, int64(0), cp.Info.Data.DownloadInfo.Offset)
	cpdata := `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"3f132b5bf65640868a47cb52c57492c8","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":5242880,"DownloadInfo":{"Offset":5242880,"CRC64":0}}}`
	err := os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)

	assert.True(t, cp.valid())
	assert.Equal(t, int64(5242880), cp.Info.Data.DownloadInfo.Offset)

	// md5 fail
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"4f132b5bf65640868a47cb52c57492c8","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":5242880,"DownloadInfo":{"Offset":5242880,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// Magic fail
	cpdata = `{"Magic":"82611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"3f132b5bf65640868a47cb52c57492c8","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":5242880,"DownloadInfo":{"Offset":5242880,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// invalid cp format
	cpdata = `"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"3f132b5bf65640868a47cb52c57492c8","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":5242880,"DownloadInfo":{"Offset":5242880,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// ObjectInfo not equal
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"cc79917c1a6cf33dc7328db5eaec13ce","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"123","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":5242880,"DownloadInfo":{"Offset":5242880,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// ObjectMeta not equal
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"ec2a5aa662eab0e40f6b2fada05374ae","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":3446061,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":5242880,"DownloadInfo":{"Offset":5242880,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// FilePath not equal
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"a60d2de5bea76d94990b7db8bb00a930","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix-1","PartSize":5242880,"DownloadInfo":{"Offset":5242880,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// PartSize not equal
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"1ea14d099d250953eb4dac39758c4148","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":52428800,"DownloadInfo":{"Offset":5242880,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// Offset invalid
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"4bf7d238bc61d53fa37f2682c0dc5aaa","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":5242880,"DownloadInfo":{"Offset":-1,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// Offset %
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"177798af2463db846520c825693bee16","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":5242880,"DownloadInfo":{"Offset":1,"CRC64":0}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// check sum equal
	cp.Info.Data.PartSize = 6
	data := "hello world!"
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"37ab56d53e402a21285972ecbbfdaaf9","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":6,"DownloadInfo":{"Offset":12,"CRC64":9548687815775124833}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	err = os.WriteFile(cp.Info.Data.FilePath, []byte(data), FilePermMode)
	assert.Nil(t, err)
	cp.VerifyData = true
	assert.True(t, cp.valid())
	assert.Equal(t, int64(12), cp.Info.Data.DownloadInfo.Offset)
	assert.Equal(t, uint64(9548687815775124833), cp.Info.Data.DownloadInfo.CRC64)

	// check sum not equal
	cp.Info.Data.DownloadInfo.Offset = 0
	cp.Info.Data.DownloadInfo.CRC64 = 0
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"cba892ede9f6fcab277293e95a6abd4f","Data":{"ObjectInfo":{"Name":"oss://bucket/key","VersionId":"","Range":""},"ObjectMeta":{"Size":344606,"LastModified":"Fri, 24 Feb 2012 06:07:48 GMT","ETag":"\"D41D8CD98F00B204E9800998ECF8****\""},"FilePath":"gthnjXGQ-no-surfix","PartSize":6,"DownloadInfo":{"Offset":12,"CRC64":9548687815775124834}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	err = os.WriteFile(cp.Info.Data.FilePath, []byte(data), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())
	assert.Equal(t, int64(0), cp.Info.Data.DownloadInfo.Offset)

	os.Remove(destFilePath)
	os.Remove(cp.CpFilePath)
}

func TestUploadCheckpoint(t *testing.T) {
	request := &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}
	destFilePath := randStr(8) + "-no-surfix"
	cpDir := "."
	partSize := DefaultDownloadPartSize

	info := &fileInfo{
		modTime: time.Now(),
		size:    int64(100),
	}

	cp := newUploadCheckpoint(request, destFilePath, cpDir, info, partSize)
	assert.NotNil(t, cp)
	assert.Equal(t, info.ModTime().String(), cp.Info.Data.FileMeta.LastModified)
	assert.Equal(t, info.Size(), cp.Info.Data.FileMeta.Size)

	assert.Equal(t, "oss://bucket/key", cp.Info.Data.ObjectInfo.Name)

	assert.Equal(t, CheckpointMagic, cp.Info.Magic)
	assert.Equal(t, "", cp.Info.MD5)

	assert.Equal(t, destFilePath, cp.Info.Data.FilePath)
	assert.Equal(t, partSize, cp.Info.Data.PartSize)

	//check dump
	cp.Info.Data.UploadInfo.UploadId = "upload-id"
	cp.dump()
	assert.True(t, FileExists(cp.CpFilePath))
	ucp := uploadCheckpoint{}
	content, err := os.ReadFile(cp.CpFilePath)
	assert.Nil(t, err)
	err = json.Unmarshal(content, &ucp.Info)
	assert.Nil(t, err)

	assert.Equal(t, info.ModTime().String(), ucp.Info.Data.FileMeta.LastModified)
	assert.Equal(t, info.Size(), ucp.Info.Data.FileMeta.Size)

	assert.Equal(t, "oss://bucket/key", ucp.Info.Data.ObjectInfo.Name)

	assert.Equal(t, CheckpointMagic, ucp.Info.Magic)
	assert.Equal(t, cp.Info.MD5, ucp.Info.MD5)
	assert.NotEmpty(t, ucp.Info.MD5)

	assert.Equal(t, destFilePath, ucp.Info.Data.FilePath)
	assert.Equal(t, partSize, ucp.Info.Data.PartSize)
	assert.Equal(t, "upload-id", ucp.Info.Data.UploadInfo.UploadId)

	//check load
	err = cp.load()
	assert.Nil(t, err)
	assert.True(t, cp.Loaded)

	//check valid
	assert.True(t, cp.valid())

	//check complete
	assert.True(t, FileExists(cp.CpFilePath))
	err = cp.remove()
	assert.Nil(t, err)
	assert.True(t, !FileExists(cp.CpFilePath))

	//load not match
	cp = newUploadCheckpoint(request, destFilePath, cpDir, info, partSize)
	assert.False(t, cp.Loaded)
	notMatch := `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"5ff2e8fbddc007157488c1087105f6d2","Data":{"FilePath":"vhetHfkY-no-surfix","FileMeta":{"Size":100,"LastModified":"2024-01-08 16:46:27.7178907 +0800 CST m=+0.014509001"},"ObjectInfo":{"Name":"oss://bucket/key"},"PartSize":5242880,"UploadInfo":{"UploadId":""}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(notMatch), FilePermMode)
	assert.Nil(t, err)
	assert.True(t, FileExists(cp.CpFilePath))
	err = cp.load()
	assert.Nil(t, err)
	assert.False(t, cp.Loaded)
	assert.False(t, FileExists(cp.CpFilePath))
}

func TestUploadCheckpointInvalidCpPath(t *testing.T) {
	request := &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}
	destFilePath := randStr(8) + "-no-surfix"
	cpDir := "./invliad-dir/"
	partSize := DefaultDownloadPartSize

	info := &fileInfo{
		modTime: time.Now(),
		size:    int64(100),
	}

	cp := newUploadCheckpoint(request, destFilePath, cpDir, info, partSize)
	assert.NotNil(t, cp)
	assert.Equal(t, destFilePath, cp.Info.Data.FilePath)
	assert.Equal(t, "invliad-dir", cp.CpDirPath)
	assert.Contains(t, cp.CpFilePath, "invliad-dir")

	//dump fail
	err := cp.dump()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "The system cannot find the path specified") || strings.Contains(err.Error(), "no such file or directory"))

	//load fail
	err = cp.load()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invaid checkpoint dir")
}

func TestUploadCheckpointValid(t *testing.T) {
	request := &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}
	destFilePath := "athnjXGQ-no-surfix"
	cpDir := "."

	modTime, err := http.ParseTime("Fri, 24 Feb 2012 06:07:48 GMT")
	assert.Nil(t, err)

	info := &fileInfo{
		modTime: modTime,
		size:    int64(100),
	}
	partSize := int64(5 * 1024 * 1024) //DefaultUploadPartSize

	cp := newUploadCheckpoint(request, destFilePath, cpDir, info, partSize)

	os.Remove(destFilePath)
	assert.Equal(t, "", cp.Info.Data.UploadInfo.UploadId)
	cpdata := `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"0c15f37bb935bec463cb3b6c362f4a21","Data":{"FilePath":"athnjXGQ-no-surfix","FileMeta":{"Size":100,"LastModified":"2012-02-24 06:07:48 +0000 UTC"},"ObjectInfo":{"Name":"oss://bucket/key"},"PartSize":5242880,"UploadInfo":{"UploadId":"upload-id"}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)

	assert.True(t, cp.valid())
	assert.Equal(t, "upload-id", cp.Info.Data.UploadInfo.UploadId)

	// md5 fail
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"1c15f37bb935bec463cb3b6c362f4a21","Data":{"FilePath":"athnjXGQ-no-surfix","FileMeta":{"Size":100,"LastModified":"2012-02-24 06:07:48 +0000 UTC"},"ObjectInfo":{"Name":"oss://bucket/key"},"PartSize":5242880,"UploadInfo":{"UploadId":"upload-id"}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// Magic fail
	cpdata = `{"Magic":"82611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"0c15f37bb935bec463cb3b6c362f4a21","Data":{"FilePath":"athnjXGQ-no-surfix","FileMeta":{"Size":100,"LastModified":"2012-02-24 06:07:48 +0000 UTC"},"ObjectInfo":{"Name":"oss://bucket/key"},"PartSize":5242880,"UploadInfo":{"UploadId":"upload-id"}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// invalid cp format
	cpdata = `"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"0c15f37bb935bec463cb3b6c362f4a21","Data":{"FilePath":"athnjXGQ-no-surfix","FileMeta":{"Size":100,"LastModified":"2012-02-24 06:07:48 +0000 UTC"},"ObjectInfo":{"Name":"oss://bucket/key"},"PartSize":5242880,"UploadInfo":{"UploadId":"upload-id"}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// FilePath not equal
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"f0920664695ea3aa8303d9de885303bd","Data":{"FilePath":"1athnjXGQ-no-surfix","FileMeta":{"Size":100,"LastModified":"2012-02-24 06:07:48 +0000 UTC"},"ObjectInfo":{"Name":"oss://bucket/key"},"PartSize":5242880,"UploadInfo":{"UploadId":"upload-id"}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// FileMeta not equal
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"79b8681208845ed7b432cb1e53790aa4","Data":{"FilePath":"athnjXGQ-no-surfix","FileMeta":{"Size":101,"LastModified":"2012-02-24 06:07:48 +0000 UTC"},"ObjectInfo":{"Name":"oss://bucket/key"},"PartSize":5242880,"UploadInfo":{"UploadId":"upload-id"}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// ObjectInfo not equal
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"3943551baf1c421f3c1eceaa890ec451","Data":{"FilePath":"athnjXGQ-no-surfix","FileMeta":{"Size":100,"LastModified":"2012-02-24 06:07:48 +0000 UTC"},"ObjectInfo":{"Name":"oss://bucket/key1"},"PartSize":5242880,"UploadInfo":{"UploadId":"upload-id"}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// PartSize not equal
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"5546cfe23f4bdde0c24c214b3dc83777","Data":{"FilePath":"athnjXGQ-no-surfix","FileMeta":{"Size":100,"LastModified":"2012-02-24 06:07:48 +0000 UTC"},"ObjectInfo":{"Name":"oss://bucket/key"},"PartSize":5242881,"UploadInfo":{"UploadId":"upload-id"}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	// uploadId invalid
	cpdata = `{"Magic":"92611BED-89E2-46B6-89E5-72F273D4B0A3","MD5":"fa7a68973cf4b51e86386db0d3d84fb1","Data":{"FilePath":"athnjXGQ-no-surfix","FileMeta":{"Size":100,"LastModified":"2012-02-24 06:07:48 +0000 UTC"},"ObjectInfo":{"Name":"oss://bucket/key"},"PartSize":5242880,"UploadInfo":{"UploadId":""}}}`
	err = os.WriteFile(cp.CpFilePath, []byte(cpdata), FilePermMode)
	assert.Nil(t, err)
	assert.False(t, cp.valid())

	os.Remove(destFilePath)
	os.Remove(cp.CpFilePath)
}
