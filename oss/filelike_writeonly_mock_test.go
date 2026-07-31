package oss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

// mpServer is an in-memory multipart-upload mock.
type mpServer struct {
	mu          sync.Mutex
	parts       map[int32][]byte // partNumber -> latest bytes
	partUploads map[int32]int    // partNumber -> number of UploadPart calls
	partSizes   map[int32][]int  // partNumber -> size of each UploadPart call
	completed   []byte           // assembled object after CompleteMultipartUpload
	putObject   []byte           // body of a single PutObject (small-object path)
	initiateHdr http.Header      // headers seen on the last InitiateMultipartUpload
	completeHdr http.Header      // headers seen on the last CompleteMultipartUpload
	putHdr      http.Header      // headers seen on the last PutObject
	uploadId    string
	partDelay   time.Duration // optional per-UploadPart delay (backpressure tests)
}

func newMPServer() *mpServer {
	return &mpServer{
		parts:       make(map[int32][]byte),
		partUploads: make(map[int32]int),
		partSizes:   make(map[int32][]int),
		uploadId:    "upload-id-xyz",
	}
}

func (s *mpServer) crc(b []byte) string {
	h := NewCRC64(0)
	h.Write(b)
	return fmt.Sprint(h.Sum64())
}

func (s *mpServer) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method, path, query := r.Method, r.URL.Path, r.URL.RawQuery
		switch {
		case method == "POST" && strings.Contains(query, "uploads"):
			// InitiateMultipartUpload
			s.mu.Lock()
			s.initiateHdr = r.Header.Clone()
			s.mu.Unlock()
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(200)
			fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult><Bucket>bucket</Bucket><Key>key</Key><UploadId>%s</UploadId></InitiateMultipartUploadResult>`, s.uploadId)

		case method == "PUT" && strings.Contains(query, "partNumber="):
			// UploadPart
			if s.partDelay > 0 {
				time.Sleep(s.partDelay)
			}
			pnStr := r.URL.Query().Get("partNumber")
			pn64, _ := strconv.ParseInt(pnStr, 10, 32)
			body, _ := io.ReadAll(r.Body)
			s.mu.Lock()
			s.parts[int32(pn64)] = body
			s.partUploads[int32(pn64)]++
			s.partSizes[int32(pn64)] = append(s.partSizes[int32(pn64)], len(body))
			s.mu.Unlock()
			w.Header().Set("ETag", fmt.Sprintf("\"etag-%s\"", pnStr))
			w.Header().Set("x-oss-hash-crc64ecma", s.crc(body))
			w.WriteHeader(200)

		case method == "GET" && strings.Contains(query, "uploadId"):
			// ListParts
			s.mu.Lock()
			nums := make([]int, 0, len(s.parts))
			for n := range s.parts {
				nums = append(nums, int(n))
			}
			sort.Ints(nums)
			var b strings.Builder
			b.WriteString(`<?xml version="1.0" encoding="UTF-8"?><ListPartsResult>`)
			fmt.Fprintf(&b, `<Bucket>bucket</Bucket><Key>key</Key><UploadId>%s</UploadId><IsTruncated>false</IsTruncated>`, s.uploadId)
			for _, n := range nums {
				body := s.parts[int32(n)]
				fmt.Fprintf(&b, `<Part><PartNumber>%d</PartNumber><ETag>"etag-%d"</ETag><Size>%d</Size><HashCrc64ecma>%s</HashCrc64ecma></Part>`,
					n, n, len(body), s.crc(body))
			}
			b.WriteString(`</ListPartsResult>`)
			s.mu.Unlock()
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(200)
			fmt.Fprint(w, b.String())

		case method == "POST" && !strings.Contains(query, "uploads"):
			// CompleteMultipartUpload — assemble parts in ascending order
			s.mu.Lock()
			s.completeHdr = r.Header.Clone()
			nums := make([]int, 0, len(s.parts))
			for n := range s.parts {
				nums = append(nums, int(n))
			}
			sort.Ints(nums)
			var full []byte
			for _, n := range nums {
				full = append(full, s.parts[int32(n)]...)
			}
			s.completed = full
			crc := s.crc(full)
			s.mu.Unlock()
			w.Header().Set("Content-Type", "application/xml")
			w.Header().Set("x-oss-hash-crc64ecma", crc)
			w.WriteHeader(200)
			fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult><Bucket>bucket</Bucket><Key>key</Key><ETag>"final-etag"</ETag></CompleteMultipartUploadResult>`)

		case method == "DELETE":
			// AbortMultipartUpload
			w.WriteHeader(204)

		case method == "PUT":
			// PutObject (small-object path)
			body, _ := io.ReadAll(r.Body)
			s.mu.Lock()
			s.putObject = body
			s.putHdr = r.Header.Clone()
			s.mu.Unlock()
			w.Header().Set("ETag", "\"put-etag\"")
			w.Header().Set("x-oss-hash-crc64ecma", s.crc(body))
			w.WriteHeader(200)

		default:
			t.Fatalf("unexpected request %s %s?%s", method, path, query)
		}
	}
}

func newTestClient(t *testing.T, endpoint string, optFns ...func(*Options)) *Client {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(endpoint)
	return NewClient(cfg, optFns...)
}

func TestWriteOnlyFileScaffolding(t *testing.T) {
	client := newTestClient(t, "http://oss-cn-hangzhou.aliyuncs.com")

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 8
		o.ParallelNum = 2
	})
	assert.NoError(t, err)
	assert.NotNil(t, f)

	// no network yet: nothing initiated
	cp := f.StatCheckpoint()
	assert.Equal(t, "", cp.UploadId)
	assert.Equal(t, int64(0), cp.Offset)
	assert.Equal(t, int64(8), cp.PartSize)
	assert.NoError(t, cp.LastError)

	// Seek query returns 0 write cursor
	pos, err := f.Seek(0, io.SeekCurrent)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), pos)

	// illegal seeks
	_, err = f.Seek(1, io.SeekCurrent)
	assert.Error(t, err)
	_, err = f.Seek(0, io.SeekStart)
	assert.Error(t, err)

	// defaults applied when unset
	f2, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key")
	assert.NoError(t, err)
	assert.Equal(t, DefaultUploadPartSize, f2.partSize)
	assert.Equal(t, DefaultUploadParallel, f2.parallelNum)
}

func TestWriteOnlyFileMultiPartComplete(t *testing.T) {
	s := newMPServer()
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 8
		o.ParallelNum = 3
	})
	assert.NoError(t, err)

	data := []byte("0123456789ABCDEFGHIJ") // 20 bytes -> parts 8,8,4
	n, err := f.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, 20, n)

	assert.NoError(t, f.Close())
	assert.Equal(t, data, s.completed)

	// terminal: further calls rejected
	_, err = f.Write([]byte("x"))
	assert.ErrorIs(t, err, os.ErrClosed)
}

func TestWriteOnlyFileContiguousOffset(t *testing.T) {
	s := newMPServer()
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 8
		o.ParallelNum = 4
	})
	assert.NoError(t, err)

	// 3 full parts; poll checkpoint after writing all data
	_, err = f.Write([]byte("AAAAAAAABBBBBBBBCCCCCCCC")) // 24 bytes -> parts 1,2,3
	assert.NoError(t, err)
	assert.NoError(t, f.Close())

	cp := f.StatCheckpoint()
	assert.Equal(t, int64(24), cp.Offset) // all three full parts contiguous
	assert.Equal(t, s.uploadId, cp.UploadId)
	// committed CRC equals crc of first 24 bytes
	h := NewCRC64(0)
	h.Write([]byte("AAAAAAAABBBBBBBBCCCCCCCC"))
	assert.Equal(t, h.Sum64(), cp.CRC64)
}

func TestWriteOnlyFileSmallObjectPut(t *testing.T) {
	s := newMPServer()
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 16
		o.ParallelNum = 2
	})
	assert.NoError(t, err)

	_, err = f.Write([]byte("hello")) // < 1 part, never flushed
	assert.NoError(t, err)
	assert.NoError(t, f.Close())

	assert.Equal(t, []byte("hello"), s.putObject)
	assert.Nil(t, s.completed)                       // no multipart
	assert.Equal(t, "", f.StatCheckpoint().UploadId) // never initiated
}

func TestWriteOnlyFileFlushDoesNotAdvanceOffset(t *testing.T) {
	s := newMPServer()
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 8
		o.ParallelNum = 2
	})
	assert.NoError(t, err)

	// write 5 bytes (< 1 part), then Flush
	_, err = f.Write([]byte("HELLO"))
	assert.NoError(t, err)
	assert.NoError(t, f.Flush())

	// flush uploaded part 1 as a short (5-byte) part; offset NOT advanced
	s.mu.Lock()
	assert.Equal(t, 1, s.partUploads[1])
	assert.Equal(t, []int{5}, s.partSizes[1])
	s.mu.Unlock()
	assert.Equal(t, int64(0), f.StatCheckpoint().Offset)

	// write cursor reflects flushed bytes
	pos, _ := f.Seek(0, io.SeekCurrent)
	assert.Equal(t, int64(5), pos)

	// continue writing to fill part 1 (needs 3 more bytes) plus part 2
	_, err = f.Write([]byte("WORLD!!")) // 7 bytes -> fills part1 (8), 4 into part2
	assert.NoError(t, err)

	assert.NoError(t, f.Close())

	// part 1 re-uploaded as a full 8-byte version (overwrote the flush short part)
	s.mu.Lock()
	assert.Equal(t, 2, s.partUploads[1]) // flush + full seal
	assert.Equal(t, []int{5, 8}, s.partSizes[1])
	assert.Equal(t, []byte("HELLOWORLD!!"), s.completed)
	s.mu.Unlock()

	// after Close, part 1 is full and part 2 is the last part -> offset covers part 1
	assert.Equal(t, int64(8), f.StatCheckpoint().Offset)
}

func TestWriteOnlyFileAbortClose(t *testing.T) {
	var aborted int32
	s := newMPServer()
	base := s.handler(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			atomic.AddInt32(&aborted, 1)
		}
		base(w, r)
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 8
		o.ParallelNum = 2
	})
	assert.NoError(t, err)
	_, err = f.Write([]byte("AAAAAAAABBBB")) // seals part 1, part 2 partial
	assert.NoError(t, err)

	assert.NoError(t, f.AbortClose())
	assert.Equal(t, int32(1), atomic.LoadInt32(&aborted))
	assert.Nil(t, s.completed)

	// idempotent: repeated AbortClose and a later Close are no-ops (no second
	// DELETE, no COMPLETE), while Write still rejects the closed handle.
	assert.NoError(t, f.AbortClose())
	assert.NoError(t, f.Close())
	assert.Equal(t, int32(1), atomic.LoadInt32(&aborted))
	assert.Nil(t, s.completed)
	_, err = f.Write([]byte("x"))
	assert.ErrorIs(t, err, os.ErrClosed)
}

func TestWriteOnlyFileCloseIdempotent(t *testing.T) {
	s := newMPServer()
	server := httptest.NewServer(http.HandlerFunc(s.handler(t)))
	defer server.Close()
	client := newTestClient(t, server.URL)

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 8
		o.ParallelNum = 2
	})
	assert.NoError(t, err)
	_, err = f.Write([]byte("AAAAAAAABBBB")) // seals part 1, part 2 partial
	assert.NoError(t, err)

	assert.NoError(t, f.Close())
	assert.Equal(t, []byte("AAAAAAAABBBB"), s.completed)

	// idempotent: a second Close is a no-op (object unchanged), AbortClose does
	// not abort a committed object, and Write still rejects the closed handle.
	assert.NoError(t, f.Close())
	assert.NoError(t, f.AbortClose())
	assert.Equal(t, []byte("AAAAAAAABBBB"), s.completed)
	_, err = f.Write([]byte("x"))
	assert.ErrorIs(t, err, os.ErrClosed)
}

func TestWriteOnlyFileStickyErrorNoAutoAbort(t *testing.T) {
	var aborted, completed int32
	s := newMPServer()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, query := r.Method, r.URL.RawQuery
		switch {
		case method == "POST" && strings.Contains(query, "uploads"):
			w.WriteHeader(200)
			fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult><UploadId>%s</UploadId></InitiateMultipartUploadResult>`, s.uploadId)
		case method == "PUT" && strings.Contains(query, "partNumber="):
			// every UploadPart fails hard
			w.Header().Set("x-oss-request-id", "req-part-fail")
			w.WriteHeader(403)
			fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
<Error><Code>AccessDenied</Code><Message>denied</Message></Error>`)
		case method == "POST":
			atomic.AddInt32(&completed, 1)
			w.WriteHeader(200)
		case method == "DELETE":
			atomic.AddInt32(&aborted, 1)
			w.WriteHeader(204)
		default:
			w.WriteHeader(200)
		}
	}))
	defer server.Close()
	client := newTestClient(t, server.URL, func(o *Options) { o.RetryMaxAttempts = Ptr(1) })

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 4
		o.ParallelNum = 1
	})
	assert.NoError(t, err)

	// write enough to seal at least one part and trigger the failing UploadPart
	_, _ = f.Write([]byte("AAAABBBBCCCC"))

	// the sticky error surfaces via StatCheckpoint and terminates Write eventually
	// (poll: worker must observe the failure)
	f.partsWg.Wait()
	assert.Error(t, f.StatCheckpoint().LastError)

	// Close returns the sticky error and does NOT auto-abort (parts kept for resume)
	err = f.Close()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AccessDenied")
	assert.Equal(t, int32(0), atomic.LoadInt32(&aborted))
	assert.Equal(t, int32(0), atomic.LoadInt32(&completed))
}

func TestWriteOnlyFileResumeValid(t *testing.T) {
	s := newMPServer()
	// simulate a prior interrupted session: two full 8-byte parts + a trailing
	// 3-byte Flush snapshot at part 3.
	s.parts[1] = []byte("AAAAAAAA")
	s.parts[2] = []byte("BBBBBBBB")
	s.parts[3] = []byte("CCC") // trailing short part (flush snapshot) -> ignored
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	h := NewCRC64(0)
	h.Write([]byte("AAAAAAAABBBBBBBB"))
	cp := WriteCheckpoint{UploadId: s.uploadId, Offset: 16, CRC64: h.Sum64(), PartSize: 8}

	f, err := client.OpenWriteOnlyFile(context.TODO(), "bucket", "key", cp)
	assert.NoError(t, err)
	// seeded from full prefix [1,2]; trailing short part 3 ignored
	got := f.StatCheckpoint()
	assert.Equal(t, int64(16), got.Offset)
	assert.Equal(t, s.uploadId, got.UploadId)

	// replay from offset 16: write the real part-3 payload and finish
	_, err = f.Write([]byte("CCCCCCCCDD")) // fills part 3 (8) then part 4 (2)
	assert.NoError(t, err)
	assert.NoError(t, f.Close())
	assert.Equal(t, []byte("AAAAAAAABBBBBBBBCCCCCCCCDD"), s.completed)
}

func TestWriteOnlyFileResumeInvalidMiddleShort(t *testing.T) {
	s := newMPServer()
	s.parts[1] = []byte("AAAAAAAA")
	s.parts[2] = []byte("BB") // middle short part -> invalid
	s.parts[3] = []byte("CCCCCCCC")
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	cp := WriteCheckpoint{UploadId: s.uploadId, Offset: 16, PartSize: 8}
	f, err := client.OpenWriteOnlyFile(context.TODO(), "bucket", "key", cp)
	// invalid checkpoint (middle short part) is refused, not silently restarted
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrCheckpointInvalid)
}

func TestWriteOnlyFileResumeOffsetMismatch(t *testing.T) {
	s := newMPServer()
	// server holds one full part (durable prefix = 8) ...
	s.parts[1] = []byte("AAAAAAAA")
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	// ... but the checkpoint claims 16 -> mismatch, refused
	cp := WriteCheckpoint{UploadId: s.uploadId, Offset: 16, PartSize: 8}
	f, err := client.OpenWriteOnlyFile(context.TODO(), "bucket", "key", cp)
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrCheckpointInvalid)
}

func TestWriteOnlyFileResumeMissingPartSize(t *testing.T) {
	client := newTestClient(t, "http://oss-cn-hangzhou.aliyuncs.com")

	// a resume token carrying an UploadId but no PartSize is malformed: the
	// durable prefix cannot be validated against a guessed part size.
	cp := WriteCheckpoint{UploadId: "upload-id-xyz", Offset: 16}
	f, err := client.OpenWriteOnlyFile(context.TODO(), "bucket", "key", cp)
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrCheckpointInvalid)
}

func TestWriteOnlyFileResumeNoSuchUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.RawQuery, "uploadId") {
			w.Header().Set("x-oss-request-id", "req-nosuch")
			w.WriteHeader(404)
			fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
<Error><Code>NoSuchUpload</Code><Message>The specified upload does not exist.</Message></Error>`)
			return
		}
		w.WriteHeader(200)
	}))
	defer server.Close()
	client := newTestClient(t, server.URL, func(o *Options) { o.RetryMaxAttempts = Ptr(1) })

	cp := WriteCheckpoint{UploadId: "gone", Offset: 24, PartSize: 8}
	f, err := client.OpenWriteOnlyFile(context.TODO(), "bucket", "key", cp)
	// a vanished upload (NoSuchUpload) is refused, not silently restarted
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrCheckpointInvalid)
}

func TestWriteOnlyFileCreateParameterMultipart(t *testing.T) {
	s := newMPServer()
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 8
		o.ParallelNum = 2
		o.CreateParameter = &PutObjectRequest{
			ContentType:  Ptr("text/plain"),
			StorageClass: StorageClassIA,
			Metadata:     map[string]string{"foo": "bar"},
			Tagging:      Ptr("a=b"),
			Acl:          ObjectACLPrivate, // no multipart equivalent -> must be dropped
		}
	})
	assert.NoError(t, err)

	_, err = f.Write([]byte("0123456789")) // 10 bytes -> seals part 1, triggers Initiate
	assert.NoError(t, err)
	assert.NoError(t, f.Close())

	s.mu.Lock()
	hdr := s.initiateHdr
	completeHdr := s.completeHdr
	s.mu.Unlock()
	assert.NotNil(t, hdr) // multipart path was taken

	// common object attributes mapped onto InitiateMultipartUpload
	assert.Equal(t, "text/plain", hdr.Get("Content-Type"))
	assert.Equal(t, string(StorageClassIA), hdr.Get("x-oss-storage-class"))
	assert.Equal(t, "bar", hdr.Get("x-oss-meta-foo"))
	assert.Equal(t, "a=b", hdr.Get("x-oss-tagging"))
	// Acl has no Initiate equivalent: it must NOT be on Initiate ...
	assert.Equal(t, "", hdr.Get("x-oss-object-acl"))
	// ... but IS applied at CompleteMultipartUpload
	assert.NotNil(t, completeHdr)
	assert.Equal(t, string(ObjectACLPrivate), completeHdr.Get("x-oss-object-acl"))
}

func TestWriteOnlyFileCreateParameterSmallObject(t *testing.T) {
	s := newMPServer()
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
		o.PartSize = 16
		o.ParallelNum = 2
		o.CreateParameter = &PutObjectRequest{
			ContentType: Ptr("application/json"),
			Metadata:    map[string]string{"k": "v"},
			Acl:         ObjectACLPublicRead, // takes effect on the PutObject path
		}
	})
	assert.NoError(t, err)

	_, err = f.Write([]byte("small")) // < 1 part, never flushed -> single PutObject
	assert.NoError(t, err)
	assert.NoError(t, f.Close())

	s.mu.Lock()
	hdr := s.putHdr
	body := s.putObject
	s.mu.Unlock()
	assert.NotNil(t, hdr)
	assert.Equal(t, []byte("small"), body)

	assert.Equal(t, "application/json", hdr.Get("Content-Type"))
	assert.Equal(t, "v", hdr.Get("x-oss-meta-k"))
	// Acl does apply on the small-object PutObject path
	assert.Equal(t, string(ObjectACLPublicRead), hdr.Get("x-oss-object-acl"))
	// never initiated a multipart upload
	assert.Nil(t, s.initiateHdr)
	assert.Equal(t, "", f.StatCheckpoint().UploadId)
}

// countWorkerGoroutines returns the number of goroutines currently executing
// (*writeOnlyCore).worker, isolating WriteOnlyFile's background workers from
// unrelated HTTP keep-alive/transport goroutines.
func countWorkerGoroutines() int {
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	return bytes.Count(buf[:n], []byte("writeOnlyCore).worker"))
}

// A handle dropped without Close/AbortClose must not leak its workers: the
// finalizer closes the channel and reaps them. This works only because workers
// close over the core, not the handle, so the handle can become unreachable.
func TestWriteOnlyFileFinalizerReapsWorkers(t *testing.T) {
	s := newMPServer()
	server := httptest.NewServer(s.handler(t))
	defer server.Close()
	client := newTestClient(t, server.URL)

	func() {
		f, err := client.NewWriteOnlyFile(context.TODO(), "bucket", "key", func(o *WriteOnlyOptions) {
			o.PartSize = 8
			o.ParallelNum = 2
		})
		assert.NoError(t, err)
		_, err = f.Write([]byte("AAAAAAAABBBBBBBB")) // seals 2 parts -> starts 2 workers
		assert.NoError(t, err)
		assert.NoError(t, f.Flush()) // drain dispatched parts; workers idle on range ch
	}()

	assert.Equal(t, 2, countWorkerGoroutines(), "workers should be live before GC")

	var n int
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		runtime.GC()
		time.Sleep(30 * time.Millisecond)
		if n = countWorkerGoroutines(); n == 0 {
			break
		}
	}
	assert.Equal(t, 0, n, "finalizer should have reaped all worker goroutines")
}
