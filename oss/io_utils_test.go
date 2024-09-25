package oss

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsyncRangeReader(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	buf := io.NopCloser(bytes.NewBufferString(data))
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: buf,
		}, nil
	}
	ar, err := NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)

	var dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, []byte(data), dst[:n])

	n, err = ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)

	// Test read after error
	n, err = ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)

	err = ar.Close()
	require.NoError(t, err)
	// Test double close
	err = ar.Close()
	require.NoError(t, err)

	// Test Close without reading everything
	buf = io.NopCloser(bytes.NewBuffer(make([]byte, 50000)))
	getFn = func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: buf,
		}, nil
	}
	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)
	err = ar.Close()
	require.NoError(t, err)
}

func TestAsyncRangeReaderErrors(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"

	// test nil reader
	_, err := NewAsyncRangeReader(ctx, nil, nil, "", 4)
	require.Error(t, err)

	// invalid buffer number
	buf := io.NopCloser(bytes.NewBufferString(data))
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: buf,
		}, nil
		//return buf, 0, "", nil
	}
	_, err = NewAsyncRangeReader(ctx, getFn, nil, "", 0)
	require.Error(t, err)
	_, err = NewAsyncRangeReader(ctx, getFn, nil, "", -1)
	require.Error(t, err)
}

type readMaker struct {
	name string
	fn   func(io.Reader) io.Reader
}

var readMakers = []readMaker{
	{"full", func(r io.Reader) io.Reader { return r }},
	{"byte", iotest.OneByteReader},
	{"half", iotest.HalfReader},
	{"data+err", iotest.DataErrReader},
	{"timeout", iotest.TimeoutReader},
}

// Call Read to accumulate the text of a file
func reads(buf io.Reader, m int) string {
	var b [1000]byte
	nb := 0
	for {
		n, err := buf.Read(b[nb : nb+m])
		nb += n
		if err == io.EOF {
			break
		} else if err != nil && err != iotest.ErrTimeout {
			panic("Data: " + err.Error())
		} else if err != nil {
			break
		}
	}
	return string(b[0:nb])
}

type bufReader struct {
	name string
	fn   func(io.Reader) string
}

var bufreaders = []bufReader{
	{"1", func(b io.Reader) string { return reads(b, 1) }},
	{"2", func(b io.Reader) string { return reads(b, 2) }},
	{"3", func(b io.Reader) string { return reads(b, 3) }},
	{"4", func(b io.Reader) string { return reads(b, 4) }},
	{"5", func(b io.Reader) string { return reads(b, 5) }},
	{"7", func(b io.Reader) string { return reads(b, 7) }},
}

const minReadBufferSize = 16

var bufsizes = []int{
	0, minReadBufferSize, 23, 32, 46, 64, 93, 128, 1024, 4096,
}

// Test various  input buffer sizes, number of buffers and read sizes.
func TestAsyncRangeReaderSizes(t *testing.T) {
	ctx := context.Background()
	var texts [31]string
	str := ""
	all := ""
	for i := 0; i < len(texts)-1; i++ {
		texts[i] = str + "\n"
		all += texts[i]
		str += string(rune(i)%26 + 'a')
	}
	texts[len(texts)-1] = all

	for h := 0; h < len(texts); h++ {
		text := texts[h]
		for i := 0; i < len(readMakers); i++ {
			for j := 0; j < len(bufreaders); j++ {
				for k := 0; k < len(bufsizes); k++ {
					for l := 1; l < 10; l++ {
						readmaker := readMakers[i]
						bufreader := bufreaders[j]
						bufsize := bufsizes[k]
						read := readmaker.fn(strings.NewReader(text))
						buf := bufio.NewReaderSize(read, bufsize)
						getFn := func(_ context.Context, httpRange HTTPRange) (output *ReaderRangeGetOutput, err error) {
							contentRange := fmt.Sprintf("bytes %v-%v/*", httpRange.Offset, httpRange.Offset+int64(bufsize))
							return &ReaderRangeGetOutput{
								Body:         io.NopCloser(buf),
								ContentRange: Ptr(contentRange),
							}, nil
						}

						ar, _ := NewAsyncRangeReader(ctx, getFn, nil, "", 1)
						s := bufreader.fn(ar)
						// "timeout" expects the Reader to recover, AsyncRangeReader does not.
						if s != text && readmaker.name != "timeout" {
							t.Errorf("reader=%s fn=%s bufsize=%d want=%q got=%q",
								readmaker.name, bufreader.name, bufsize, text, s)
						}
						err := ar.Close()
						require.NoError(t, err)
					}
				}
			}
		}
	}
}

// Read an infinite number of zeros
type zeroReader struct {
	closed bool
}

func (z *zeroReader) Read(p []byte) (n int, err error) {
	if z.closed {
		return 0, io.EOF
	}
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func (z *zeroReader) Close() error {
	if z.closed {
		panic("double close on zeroReader")
	}
	z.closed = true
	return nil
}

// Test closing and abandoning
func TestAsyncRangeReaderClose(t *testing.T) {
	ctx := context.Background()
	zr := &zeroReader{}
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: zr,
		}, nil
	}
	a, err := NewAsyncRangeReader(ctx, getFn, nil, "", 16)
	require.NoError(t, err)
	var copyN int64
	var copyErr error
	var wg sync.WaitGroup
	started := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		close(started)
		{
			// exercise the Read path
			buf := make([]byte, 64*1024)
			for {
				var n int
				n, copyErr = a.Read(buf)
				copyN += int64(n)
				if copyErr != nil {
					break
				}
			}
		}
	}()
	// Do some copying
	<-started
	time.Sleep(100 * time.Millisecond)
	// abandon the copy
	a.abandon()
	wg.Wait()
	assert.Contains(t, copyErr.Error(), "stream abandoned")
	// t.Logf("Copied %d bytes, err %v", copyN, copyErr)
	assert.True(t, copyN > 0)
}

func TestAsyncRangeReaderEtagCheck(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: io.NopCloser(bytes.NewBufferString(data)),
			ETag: Ptr("etag"),
		}, nil
	}

	// don't set etag
	ar, err := NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)

	var dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, data, string(dst[0:n]))

	// set etag to "etag"
	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, data, string(dst[:n]))

	// set etag to "invalid-etag"
	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "invalid-etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Contains(t, err.Error(), "Source file is changed, expect etag:invalid-etag")
}

func TestAsyncRangeReaderOffsetCheck(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBufferString(data)))),
			ETag: Ptr("etag"),
		}, nil
	}

	// don't set etag
	ar, err := NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)

	var dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 1, n)
	n, err = ar.Read(dst)
	assert.Equal(t, 0, n)
	assert.Contains(t, err.Error(), "Range get fail, expect offset")

	//
	getFn = func(ctx context.Context, range_ HTTPRange) (output *ReaderRangeGetOutput, err error) {
		b := []byte(data)
		if range_.Offset == 0 {
			return &ReaderRangeGetOutput{
				Body: io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBuffer(b[range_.Offset:])))),
				ETag: Ptr("etag"),
			}, nil
		} else {
			contentRange := fmt.Sprintf("bytes %v-%v/*", range_.Offset, int64(len(data))-range_.Offset-1)
			return &ReaderRangeGetOutput{
				Body:         io.NopCloser(bytes.NewBuffer(b[range_.Offset:])),
				ContentRange: Ptr(contentRange),
				ETag:         Ptr("etag"),
			}, nil
		}
	}

	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 1, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 9, n)
	assert.Equal(t, data, string(dst[:10]))
}

func TestAsyncRangeReaderOffsetAndLimiterRange(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(ctx context.Context, range_ HTTPRange) (output *ReaderRangeGetOutput, err error) {
		b := []byte(data)
		contentRange := fmt.Sprintf("bytes %v-%v/%v", range_.Offset, int64(len(data))-range_.Offset-1, len(data))
		return &ReaderRangeGetOutput{
			Body:         io.NopCloser(bytes.NewBuffer(b[range_.Offset:])),
			ContentRange: Ptr(contentRange),
			ETag:         Ptr("etag"),
		}, nil
	}

	// has range count
	httpRange := HTTPRange{
		Offset: 2,
		Count:  3,
	}

	ar, err := NewAsyncRangeReader(ctx, getFn, &httpRange, "etag", 4)
	require.NoError(t, err)

	dst := make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 3, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[2:5], string(dst[:3]))

	// range count is 0 or < 0
	httpRange = HTTPRange{
		Offset: 2,
		Count:  0,
	}

	ar, err = NewAsyncRangeReader(ctx, getFn, &httpRange, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 8, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[2:], string(dst[:8]))

	httpRange = HTTPRange{
		Offset: 3,
		Count:  -1,
	}

	ar, err = NewAsyncRangeReader(ctx, getFn, &httpRange, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 7, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[3:], string(dst[:7]))
}

func TestAsyncRangeReaderOffsetAndResume(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	errCount := 0
	getFn := func(ctx context.Context, range_ HTTPRange) (output *ReaderRangeGetOutput, err error) {
		b := []byte(data)
		var range_count int64
		if range_.Count > 0 {
			range_count = range_.Count
		} else {
			range_count = int64(len(data)) - range_.Offset
		}
		contentRange := fmt.Sprintf("bytes %v-%v/%v", range_.Offset, range_.Offset+range_count-1, len(data))

		body := io.NopCloser(bytes.NewBuffer(b[range_.Offset : range_.Offset+range_count]))
		if errCount > 0 {
			body = io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(body)))
			errCount -= 1
		}
		return &ReaderRangeGetOutput{
			Body:         body,
			ContentRange: Ptr(contentRange),
			ETag:         Ptr("etag"),
		}, nil
	}

	// read fail and resume read
	// has range count, read pattern, 1 byte, 2 bytes
	httpRange := HTTPRange{
		Offset: 2,
		Count:  3,
	}
	errCount = 1
	ar, err := NewAsyncRangeReader(ctx, getFn, &httpRange, "etag", 4)
	require.NoError(t, err)
	got := 0
	dst := make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 1, n)
	got += n
	n, err = ar.Read(dst[got:])
	assert.Equal(t, 2, n)
	got += n
	n, err = ar.Read(dst[got:])
	assert.Equal(t, 0, n)

	assert.Equal(t, data[2:5], string(dst[:3]))

	// range count is 0 or < 0
	httpRange = HTTPRange{
		Offset: 2,
		Count:  0,
	}
	errCount = 2
	ar, err = NewAsyncRangeReader(ctx, getFn, &httpRange, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	got = 0
	n, err = ar.Read(dst)
	assert.Equal(t, 1, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 1, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 6, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 0, n)

	assert.Equal(t, data[2:], string(dst[:8]))

	httpRange = HTTPRange{
		Offset: 3,
		Count:  -1,
	}
	errCount = 2
	ar, err = NewAsyncRangeReader(ctx, getFn, &httpRange, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	got = 0
	n, err = ar.Read(dst)
	assert.Equal(t, 1, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 1, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 5, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 0, n)

	assert.Equal(t, data[3:], string(dst[:7]))

	// not fail
	// has range count
	httpRange = HTTPRange{
		Offset: 2,
		Count:  3,
	}

	ar, err = NewAsyncRangeReader(ctx, getFn, &httpRange, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 3, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[2:5], string(dst[:3]))

	// range count is 0 or < 0
	httpRange = HTTPRange{
		Offset: 2,
		Count:  0,
	}

	ar, err = NewAsyncRangeReader(ctx, getFn, &httpRange, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 8, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[2:], string(dst[:8]))

	httpRange = HTTPRange{
		Offset: 3,
		Count:  -1,
	}

	ar, err = NewAsyncRangeReader(ctx, getFn, &httpRange, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 7, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[3:], string(dst[:7]))
}

func TestAsyncRangeReaderGetError(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{}, errors.New("range get fail")
	}

	// don't set etag
	ar, err := NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)

	var dst = make([]byte, 100)
	_, err = ar.Read(dst)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "range get fail")

	getFn = func(ctx context.Context, range_ HTTPRange) (output *ReaderRangeGetOutput, err error) {
		b := []byte(data)
		if range_.Offset == 0 {
			contentRange := fmt.Sprintf("bytes %v-%v/*", range_.Offset, len(data)-1)
			return &ReaderRangeGetOutput{
				Body:         io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBuffer(b[range_.Offset:])))),
				ContentRange: Ptr(contentRange),
				ETag:         Ptr("etag"),
			}, nil
		} else {
			return &ReaderRangeGetOutput{}, errors.New("range get fail")
		}
	}

	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "etag", 4)
	require.NoError(t, err)
	dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 1, n)
	assert.Equal(t, int64(1), ar.offset)
	n, err = ar.Read(dst[n:])
	require.Error(t, err)
	assert.Contains(t, err.Error(), "range get fail")
}

func TestMultiBytesReader(t *testing.T) {
	datas := []struct {
		data [][]byte
	}{
		{[][]byte{[]byte("0123456789")}},
		{[][]byte{[]byte("01"), []byte("23"), []byte("45"), []byte("67"), []byte("89")}},
		{[][]byte{[]byte("0"), []byte("12"), []byte("345"), []byte("6789")}},
		{[][]byte{nil, []byte("0"), []byte("12"), nil, []byte("345"), []byte("6789"), nil}},
	}

	for _, d := range datas {
		r := NewMultiBytesReader(d.data)
		tests := []struct {
			off     int64
			seek    int
			n       int
			want    string
			wantpos int64
			readerr error
			seekerr string
		}{
			{seek: io.SeekStart, off: 0, n: 20, want: "0123456789"},
			{seek: io.SeekStart, off: 1, n: 1, want: "1"},
			{seek: io.SeekCurrent, off: 1, wantpos: 3, n: 2, want: "34"},
			{seek: io.SeekStart, off: -1, seekerr: "MultiSliceReader.Seek: negative position"},
			{seek: io.SeekStart, off: 1 << 33, wantpos: 1 << 33, readerr: io.EOF},
			{seek: io.SeekCurrent, off: 1, wantpos: 1<<33 + 1, readerr: io.EOF},
			{seek: io.SeekStart, n: 5, want: "01234"},
			{seek: io.SeekCurrent, n: 5, want: "56789"},
			{seek: io.SeekEnd, off: -1, n: 1, wantpos: 9, want: "9"},
		}

		for i, tt := range tests {
			pos, err := r.Seek(tt.off, tt.seek)
			if err == nil && tt.seekerr != "" {
				t.Errorf("%d. want seek error %q", i, tt.seekerr)
				continue
			}
			if err != nil && err.Error() != tt.seekerr {
				t.Errorf("%d. seek error = %q; want %q", i, err.Error(), tt.seekerr)
				continue
			}
			if tt.wantpos != 0 && tt.wantpos != pos {
				t.Errorf("%d. pos = %d, want %d", i, pos, tt.wantpos)
			}
			buf := make([]byte, tt.n)
			n, err := r.Read(buf)
			if err != tt.readerr {
				t.Errorf("%d. read = %v; want %v", i, err, tt.readerr)
				continue
			}
			got := string(buf[:n])
			if got != tt.want {
				t.Errorf("%d. got %q; want %q", i, got, tt.want)
			}
		}
	}
}

func TestMultiBytesReaderAfterBigSeek(t *testing.T) {
	datas := []struct {
		data [][]byte
	}{
		{[][]byte{[]byte("0123456789")}},
		{[][]byte{[]byte("01"), []byte("23"), []byte("45"), []byte("67"), []byte("89")}},
		{[][]byte{[]byte("0"), []byte("12"), []byte("345"), []byte("6789")}},
		{[][]byte{nil, []byte("0"), []byte("12"), nil, []byte("345"), []byte("6789"), nil}},
	}

	for _, d := range datas {
		r := NewMultiBytesReader(d.data)
		if _, err := r.Seek(1<<31+5, io.SeekStart); err != nil {
			t.Fatal(err)
		}
		if n, err := r.Read(make([]byte, 10)); n != 0 || err != io.EOF {
			t.Errorf("Read = %d, %v; want 0, EOF", n, err)
		}
	}
}

func testEmptyMultiBytesReaderConcurrent(_ *testing.T) {
	datas := []struct {
		data [][]byte
	}{
		{[][]byte{[]byte("0123456789")}},
		{[][]byte{[]byte("01"), []byte("23"), []byte("45"), []byte("67"), []byte("89")}},
		{[][]byte{[]byte("0"), []byte("12"), []byte("345"), []byte("6789")}},
		{[][]byte{nil, []byte("0"), []byte("12"), nil, []byte("345"), []byte("6789"), nil}},
	}

	for _, d := range datas {
		r := NewMultiBytesReader(d.data)
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				var buf [1]byte
				r.Read(buf[:])
			}()
			go func() {
				defer wg.Done()
				r.Read(nil)
			}()
		}
		wg.Wait()
	}
}

func TestMultiBytesReaderLen(t *testing.T) {
	datas := []struct {
		data [][]byte
	}{
		{[][]byte{[]byte("hello world")}},
		{[][]byte{[]byte("hello"), []byte(" world")}},
		{[][]byte{[]byte("hello"), nil, []byte(" world")}},
	}

	for _, d := range datas {
		r := NewMultiBytesReader(d.data)
		if got, want := r.Len(), 11; got != want {
			t.Errorf("r.Len(): got %d, want %d", got, want)
		}
		if n, err := r.Read(make([]byte, 10)); err != nil || n != 10 {
			t.Errorf("Read failed: read %d %v", n, err)
		}
		if got, want := r.Len(), 1; got != want {
			t.Errorf("r.Len(): got %d, want %d", got, want)
		}
		if n, err := r.Read(make([]byte, 1)); err != nil || n != 1 {
			t.Errorf("Read failed: read %d %v; want 1, nil", n, err)
		}
		if got, want := r.Len(), 0; got != want {
			t.Errorf("r.Len(): got %d, want %d", got, want)
		}
	}
}
func TestMultiBytesReaderCopyNothing(t *testing.T) {
	type nErr struct {
		n   int64
		err error
	}
	type justReader struct {
		io.Reader
	}
	type justWriter struct {
		io.Writer
	}
	discard := justWriter{io.Discard} // hide ReadFrom

	var with, withOut nErr
	with.n, with.err = io.Copy(discard, NewMultiBytesReader(nil))
	withOut.n, withOut.err = io.Copy(discard, justReader{NewMultiBytesReader(nil)})
	if with != withOut {
		t.Errorf("behavior differs: with = %#v; without: %#v", with, withOut)
	}
}

// tests that Len is affected by reads, but Size is not.
func TestMultiBytesReaderLenSize(t *testing.T) {
	datas := []struct {
		data [][]byte
	}{
		{[][]byte{[]byte("abc")}},
		{[][]byte{[]byte("a"), []byte("bc")}},
		{[][]byte{[]byte("ab"), nil, []byte("c")}},
	}

	for _, d := range datas {
		r := NewMultiBytesReader(d.data)
		io.CopyN(io.Discard, r, 1)
		if r.Len() != 2 {
			t.Errorf("Len = %d; want 2", r.Len())
		}
		if r.Size() != 3 {
			t.Errorf("Size = %d; want 3", r.Size())
		}
	}
}

func TestMultiBytesReaderReset(t *testing.T) {
	r := NewMultiBytesReader([][]byte{[]byte("世界")})
	const want = "abcdef"
	r.Reset([][]byte{[]byte(want)})

	buf, err := io.ReadAll(r)
	if err != nil {
		t.Errorf("ReadAll: unexpected error: %v", err)
	}
	if got := string(buf); got != want {
		t.Errorf("ReadAll: got %q, want %q", got, want)
	}
}

func TestReaderZero(t *testing.T) {
	if l := (&MultiBytesReader{}).Len(); l != 0 {
		t.Errorf("Len: got %d, want 0", l)
	}

	if n, err := (&MultiBytesReader{}).Read(nil); n != 0 || err != io.EOF {
		t.Errorf("Read: got %d, %v; want 0, io.EOF", n, err)
	}

	if offset, err := (&MultiBytesReader{}).Seek(11, io.SeekStart); offset != 11 || err != nil {
		t.Errorf("Seek: got %d, %v; want 11, nil", offset, err)
	}

	if s := (&MultiBytesReader{}).Size(); s != 0 {
		t.Errorf("Size: got %d, want 0", s)
	}
}

func TestRangeReader(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	buf := io.NopCloser(bytes.NewBufferString(data))
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: buf,
		}, nil
	}
	ar, err := NewRangeReader(ctx, getFn, nil, "")
	require.NoError(t, err)

	var dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 10, n)
	assert.Equal(t, []byte(data), dst[:n])

	n, err = ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)

	// Test read after error
	n, err = ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)

	err = ar.Close()
	require.NoError(t, err)
	// Test double close
	err = ar.Close()
	require.NoError(t, err)

	// Test Close without reading everything
	buf = io.NopCloser(bytes.NewBuffer(make([]byte, 50000)))
	getFn = func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: buf,
		}, nil
	}
	ar, err = NewRangeReader(ctx, getFn, nil, "")
	require.NoError(t, err)
	err = ar.Close()
	require.NoError(t, err)
}

func TestRangeReaderErrors(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"

	// test nil reader
	_, err := NewRangeReader(ctx, nil, nil, "")
	require.Error(t, err)

	// invalid buffer number
	buf := io.NopCloser(bytes.NewBufferString(data))
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: buf,
		}, nil
		//return buf, 0, "", nil
	}
	_, err = NewAsyncRangeReader(ctx, getFn, nil, "", 0)
	require.Error(t, err)
	_, err = NewAsyncRangeReader(ctx, getFn, nil, "", -1)
	require.Error(t, err)
}

func TestRangeReaderSizes(t *testing.T) {
	ctx := context.Background()
	var texts [31]string
	str := ""
	all := ""
	for i := 0; i < len(texts)-1; i++ {
		texts[i] = str + "\n"
		all += texts[i]
		str += string(rune(i)%26 + 'a')
	}
	texts[len(texts)-1] = all

	for h := 0; h < len(texts); h++ {
		text := texts[h]
		for i := 0; i < len(readMakers); i++ {
			for j := 0; j < len(bufreaders); j++ {
				for k := 0; k < len(bufsizes); k++ {
					for l := 1; l < 10; l++ {
						readmaker := readMakers[i]
						bufreader := bufreaders[j]
						bufsize := bufsizes[k]
						read := readmaker.fn(strings.NewReader(text))
						buf := bufio.NewReaderSize(read, bufsize)
						getFn := func(_ context.Context, httpRange HTTPRange) (output *ReaderRangeGetOutput, err error) {
							contentRange := fmt.Sprintf("bytes %v-%v/*", httpRange.Offset, httpRange.Offset+int64(bufsize))
							return &ReaderRangeGetOutput{
								Body:         io.NopCloser(buf),
								ContentRange: Ptr(contentRange),
							}, nil
						}

						ar, _ := NewRangeReader(ctx, getFn, nil, "")
						s := bufreader.fn(ar)
						// "timeout" expects the Reader to recover, AsyncRangeReader does not.
						if s != text && readmaker.name != "timeout" {
							t.Errorf("reader=%s fn=%s bufsize=%d want=%q got=%q",
								readmaker.name, bufreader.name, bufsize, text, s)
						}
						err := ar.Close()
						require.NoError(t, err)
					}
				}
			}
		}
	}
}

func TestRangeReaderEtagCheck(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: io.NopCloser(bytes.NewBufferString(data)),
			ETag: Ptr("etag"),
		}, nil
	}

	// don't set etag
	ar, err := NewRangeReader(ctx, getFn, nil, "")
	require.NoError(t, err)

	var dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 10, n)
	assert.Equal(t, data, string(dst[0:n]))

	// set etag to "etag"
	ar, err = NewRangeReader(ctx, getFn, nil, "etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 10, n)
	assert.Equal(t, data, string(dst[:n]))

	// set etag to "invalid-etag"
	ar, err = NewRangeReader(ctx, getFn, nil, "invalid-etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Contains(t, err.Error(), "Source file is changed, expect etag:invalid-etag")
}

func TestRangeReaderOffsetCheck(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{
			Body: io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBufferString(data)))),
			ETag: Ptr("etag"),
		}, nil
	}

	// don't set etag
	ar, err := NewRangeReader(ctx, getFn, nil, "")
	require.NoError(t, err)

	var dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 1, n)
	n, err = ar.Read(dst)
	assert.Equal(t, 0, n)
	assert.Nil(t, err)
	n, err = ar.Read(dst)
	assert.Equal(t, 0, n)
	assert.Contains(t, err.Error(), "Range get fail, expect offset")

	//
	getFn = func(ctx context.Context, range_ HTTPRange) (output *ReaderRangeGetOutput, err error) {
		b := []byte(data)
		if range_.Offset == 0 {
			return &ReaderRangeGetOutput{
				Body: io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBuffer(b[range_.Offset:])))),
				ETag: Ptr("etag"),
			}, nil
		} else {
			contentRange := fmt.Sprintf("bytes %v-%v/*", range_.Offset, int64(len(data))-range_.Offset-1)
			return &ReaderRangeGetOutput{
				Body:         io.NopCloser(bytes.NewBuffer(b[range_.Offset:])),
				ContentRange: Ptr(contentRange),
				ETag:         Ptr("etag"),
			}, nil
		}
	}

	ar, err = NewRangeReader(ctx, getFn, nil, "etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 1, n)
	n, err = io.ReadFull(ar, dst[n:])
	assert.Equal(t, 9, n)
	assert.Equal(t, data, string(dst[:10]))
}

func TestRangeReaderOffsetAndLimiterRange(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(ctx context.Context, range_ HTTPRange) (output *ReaderRangeGetOutput, err error) {
		b := []byte(data)
		contentRange := fmt.Sprintf("bytes %v-%v/%v", range_.Offset, int64(len(data))-range_.Offset-1, len(data))
		return &ReaderRangeGetOutput{
			Body:         io.NopCloser(bytes.NewBuffer(b[range_.Offset:])),
			ContentRange: Ptr(contentRange),
			ETag:         Ptr("etag"),
		}, nil
	}

	// has range count
	httpRange := HTTPRange{
		Offset: 2,
		Count:  3,
	}

	ar, err := NewRangeReader(ctx, getFn, &httpRange, "etag")
	require.NoError(t, err)

	dst := make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 3, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[2:5], string(dst[:3]))

	// range count is 0 or < 0
	httpRange = HTTPRange{
		Offset: 2,
		Count:  0,
	}

	ar, err = NewRangeReader(ctx, getFn, &httpRange, "etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 8, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[2:], string(dst[:8]))

	httpRange = HTTPRange{
		Offset: 3,
		Count:  -1,
	}

	ar, err = NewRangeReader(ctx, getFn, &httpRange, "etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 7, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[3:], string(dst[:7]))
}

func TestRangeReaderOffsetAndResume(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	errCount := 0
	getFn := func(ctx context.Context, range_ HTTPRange) (output *ReaderRangeGetOutput, err error) {
		b := []byte(data)
		var range_count int64
		if range_.Count > 0 {
			range_count = range_.Count
		} else {
			range_count = int64(len(data)) - range_.Offset
		}
		contentRange := fmt.Sprintf("bytes %v-%v/%v", range_.Offset, range_.Offset+range_count-1, len(data))

		body := io.NopCloser(bytes.NewBuffer(b[range_.Offset : range_.Offset+range_count]))
		if errCount > 0 {
			body = io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(body)))
			errCount -= 1
		}
		return &ReaderRangeGetOutput{
			Body:         body,
			ContentRange: Ptr(contentRange),
			ETag:         Ptr("etag"),
		}, nil
	}

	// read fail and resume read
	// has range count, read pattern, 1 byte, 0 byts, 2 bytes
	httpRange := HTTPRange{
		Offset: 2,
		Count:  3,
	}
	errCount = 1
	ar, err := NewRangeReader(ctx, getFn, &httpRange, "etag")
	require.NoError(t, err)
	got := 0
	dst := make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 1, n)
	got += n
	n, err = ar.Read(dst[got:])
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	got += n
	n, err = ar.Read(dst[got:])
	assert.Equal(t, 2, n)
	got += n
	n, err = ar.Read(dst[got:])
	assert.Equal(t, 0, n)

	assert.Equal(t, data[2:5], string(dst[:3]))

	// range count is 0 or < 0
	httpRange = HTTPRange{
		Offset: 2,
		Count:  0,
	}
	errCount = 2
	ar, err = NewRangeReader(ctx, getFn, &httpRange, "etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	got = 0
	n, err = ar.Read(dst)
	assert.Equal(t, 1, n)
	got += n

	n, err = ar.Read(dst[got:])
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 1, n)
	got += n

	n, err = ar.Read(dst[got:])
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 6, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 0, n)

	assert.Equal(t, data[2:], string(dst[:8]))

	httpRange = HTTPRange{
		Offset: 3,
		Count:  -1,
	}
	errCount = 2
	ar, err = NewRangeReader(ctx, getFn, &httpRange, "etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	got = 0
	n, err = ar.Read(dst)
	assert.Equal(t, 1, n)
	got += n

	n, err = ar.Read(dst[got:])
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 1, n)
	got += n

	n, err = ar.Read(dst[got:])
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 5, n)
	got += n

	n, err = ar.Read(dst[got:])
	assert.Equal(t, 0, n)

	assert.Equal(t, data[3:], string(dst[:7]))

	// not fail
	// has range count
	httpRange = HTTPRange{
		Offset: 2,
		Count:  3,
	}

	ar, err = NewRangeReader(ctx, getFn, &httpRange, "etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 3, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[2:5], string(dst[:3]))

	// range count is 0 or < 0
	httpRange = HTTPRange{
		Offset: 2,
		Count:  0,
	}

	ar, err = NewRangeReader(ctx, getFn, &httpRange, "etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 8, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[2:], string(dst[:8]))

	httpRange = HTTPRange{
		Offset: 3,
		Count:  -1,
	}

	ar, err = NewRangeReader(ctx, getFn, &httpRange, "etag")
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 7, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 0, n)
	assert.Equal(t, data[3:], string(dst[:7]))
}

func TestRangeReaderGetError(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(context.Context, HTTPRange) (output *ReaderRangeGetOutput, err error) {
		return &ReaderRangeGetOutput{}, errors.New("range get fail")
	}

	// don't set etag
	ar, err := NewRangeReader(ctx, getFn, nil, "")
	require.NoError(t, err)

	var dst = make([]byte, 100)
	_, err = ar.Read(dst)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "range get fail")

	getFn = func(ctx context.Context, range_ HTTPRange) (output *ReaderRangeGetOutput, err error) {
		b := []byte(data)
		if range_.Offset == 0 {
			contentRange := fmt.Sprintf("bytes %v-%v/*", range_.Offset, len(data)-1)
			return &ReaderRangeGetOutput{
				Body:         io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBuffer(b[range_.Offset:])))),
				ContentRange: Ptr(contentRange),
				ETag:         Ptr("etag"),
			}, nil
		} else {
			return &ReaderRangeGetOutput{}, errors.New("range get fail")
		}
	}

	ar, err = NewRangeReader(ctx, getFn, nil, "etag")
	require.NoError(t, err)
	dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 1, n)
	assert.Equal(t, int64(1), ar.offset)
	n, err = io.ReadFull(ar, dst[n:])
	require.Error(t, err)
	assert.Contains(t, err.Error(), "range get fail")
}

type seekerReaderStub struct {
	r    io.ReadSeeker
	bErr bool
	cErr bool
	eErr bool
}

func (r *seekerReaderStub) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r *seekerReaderStub) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		if r.bErr {
			return 0, errors.New("SeekStart error")
		}
	case io.SeekCurrent:
		if r.cErr {
			return 0, errors.New("SeekCurrent error")
		}
	case io.SeekEnd:
		if r.eErr {
			return 0, errors.New("SeekEnd error")
		}
	}
	return r.r.Seek(offset, whence)
}

func TestGetReaderLen(t *testing.T) {
	data := "hello world"

	// bytes.Buffer
	b := bytes.NewBuffer([]byte(data))
	n := GetReaderLen(b)
	assert.Equal(t, int64(len(data)), n)

	// bytes.Reader
	br := bytes.NewReader([]byte(data))
	n = GetReaderLen(br)
	assert.Equal(t, int64(len(data)), n)

	// strings.Reader
	sr := strings.NewReader(data)
	n = GetReaderLen(sr)
	assert.Equal(t, int64(len(data)), n)

	// os.File
	filePath := randStr(8) + ".txt"
	createFile(t, filePath, data)
	f, err := os.Open(filePath)
	defer os.Remove(filePath)
	defer f.Close()
	assert.Nil(t, err)
	n = GetReaderLen(f)
	assert.Equal(t, int64(len(data)), n)

	f.Seek(0, io.SeekEnd)
	n = GetReaderLen(f)
	assert.Equal(t, int64(0), n)

	f.Seek(2, io.SeekStart)
	n = GetReaderLen(f)
	assert.Equal(t, int64(len(data)-2), n)

	// err
	n = GetReaderLen(nil)
	assert.Equal(t, int64(-1), n)

	//has not Len() , Seek(), N
	b = bytes.NewBuffer([]byte(data))
	bc := io.NopCloser(b)
	n = GetReaderLen(bc)
	assert.Equal(t, int64(-1), n)

	//Seek error
	sef := &seekerReaderStub{
		r: f,
	}
	sef.Seek(0, io.SeekStart)
	n = GetReaderLen(sef)
	assert.Equal(t, int64(len(data)), n)

	sef.bErr = true
	n = GetReaderLen(sef)
	assert.Equal(t, int64(-1), n)

	sef.bErr = false
	sef.cErr = true
	n = GetReaderLen(sef)
	assert.Equal(t, int64(-1), n)

	sef.cErr = false
	sef.eErr = true
	n = GetReaderLen(sef)
	assert.Equal(t, int64(-1), n)
}
