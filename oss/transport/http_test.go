package transport

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

var copyTestConfig = Config{
	ConnectTimeout:        ptr(1 * time.Second),
	ReadWriteTimeout:      ptr(2 * time.Second),
	IdleConnectionTimeout: ptr(3 * time.Second),
	KeepAliveTimeout:      ptr(5 * time.Second),
}

func TestCopy(t *testing.T) {
	want := copyTestConfig
	got := copyTestConfig.copy()
	assert.Equal(t, want, *got)

	got.ConnectTimeout = ptr(22 * time.Second)
	assert.NotEqual(t, want, *got)

	var zeroValueConfig = Config{}
	want = copyTestConfig
	got = copyTestConfig.copy(&zeroValueConfig)
	assert.Equal(t, want, *got)

	want = copyTestConfig
	got = copyTestConfig.copy(nil)
	assert.Equal(t, want, *got)

	var valueConfig = &Config{
		ConnectTimeout: ptr(10 * time.Second),
	}
	cfg := copyTestConfig.copy(valueConfig)
	assert.Equal(t, 10*time.Second, *cfg.ConnectTimeout)
	assert.Equal(t, 2*time.Second, *cfg.ReadWriteTimeout)
	assert.Equal(t, 3*time.Second, *cfg.IdleConnectionTimeout)
	assert.Equal(t, 5*time.Second, *cfg.KeepAliveTimeout)
	assert.Nil(t, cfg.EnabledRedirect)

	var config1 = &Config{
		ReadWriteTimeout: ptr(30 * time.Second),
	}
	var config2 = &Config{
		ReadWriteTimeout: ptr(330 * time.Second),
		EnabledRedirect:  ptr(true),
	}
	cfg = copyTestConfig.copy(config1, config2)
	assert.Equal(t, 1*time.Second, *cfg.ConnectTimeout)
	assert.Equal(t, 330*time.Second, *cfg.ReadWriteTimeout)
	assert.Equal(t, 3*time.Second, *cfg.IdleConnectionTimeout)
	assert.Equal(t, 5*time.Second, *cfg.KeepAliveTimeout)
	assert.Equal(t, true, *cfg.EnabledRedirect)

	cfg = copyTestConfig.copy()
	cfg.mergeIn(&DefaultConfig)
	assert.Equal(t, 5*time.Second, *cfg.ConnectTimeout)
	assert.Equal(t, 10*time.Second, *cfg.ReadWriteTimeout)
	assert.Equal(t, 50*time.Second, *cfg.IdleConnectionTimeout)
	assert.Equal(t, 30*time.Second, *cfg.KeepAliveTimeout)
}

func TestMerge(t *testing.T) {
	var zeroValueConfig = &Config{}
	want := copyTestConfig
	got := copyTestConfig.copy()
	got.mergeIn(zeroValueConfig)
	assert.Equal(t, want, *got)

	var valueConfig = &Config{
		ConnectTimeout: ptr(10 * time.Second),
	}
	cfg := copyTestConfig.copy()
	cfg.mergeIn(valueConfig)
	assert.Equal(t, 10*time.Second, *cfg.ConnectTimeout)
	assert.Equal(t, 2*time.Second, *cfg.ReadWriteTimeout)
	assert.Equal(t, 3*time.Second, *cfg.IdleConnectionTimeout)
	assert.Equal(t, 5*time.Second, *cfg.KeepAliveTimeout)
	assert.Nil(t, cfg.EnabledRedirect)

	var config1 = &Config{
		ReadWriteTimeout: ptr(30 * time.Second),
	}
	var config2 = &Config{
		ReadWriteTimeout: ptr(34 * time.Second),
		EnabledRedirect:  ptr(true),
	}
	cfg = copyTestConfig.copy()
	cfg.mergeIn(config1, config2)
	assert.Equal(t, 1*time.Second, *cfg.ConnectTimeout)
	assert.Equal(t, 34*time.Second, *cfg.ReadWriteTimeout)
	assert.Equal(t, 3*time.Second, *cfg.IdleConnectionTimeout)
	assert.Equal(t, 5*time.Second, *cfg.KeepAliveTimeout)
	assert.Equal(t, true, *cfg.EnabledRedirect)

	cfg = copyTestConfig.copy()
	cfg.mergeIn(&DefaultConfig)
	assert.Equal(t, 5*time.Second, *cfg.ConnectTimeout)
	assert.Equal(t, 10*time.Second, *cfg.ReadWriteTimeout)
	assert.Equal(t, 50*time.Second, *cfg.IdleConnectionTimeout)
	assert.Equal(t, 30*time.Second, *cfg.KeepAliveTimeout)

	//read and write
	cfg = copyTestConfig.copy()
	cfg.mergeIn(&Config{
		PostRead: []func(n int, err error){
			func(n int, err error) {},
			func(n int, err error) {},
		},
	})
	cfg.mergeIn(&DefaultConfig)
	assert.Len(t, cfg.PostRead, 2)
	assert.Nil(t, cfg.PostWrite)

	cfg = copyTestConfig.copy()
	cfg.mergeIn(&Config{
		PostWrite: []func(n int, err error){
			func(n int, err error) {},
		},
	})
	cfg.mergeIn(&DefaultConfig)
	assert.Len(t, cfg.PostWrite, 1)
	assert.Nil(t, cfg.PostRead)
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig
	assert.NotNil(t, cfg.ConnectTimeout)
	assert.NotNil(t, cfg.ReadWriteTimeout)
	assert.NotNil(t, cfg.IdleConnectionTimeout)
	assert.NotNil(t, cfg.KeepAliveTimeout)
	assert.Nil(t, cfg.EnabledRedirect)

	cfg = Config{}
	assert.Nil(t, cfg.ConnectTimeout)
	assert.Nil(t, cfg.ReadWriteTimeout)
	assert.Nil(t, cfg.IdleConnectionTimeout)
	assert.Nil(t, cfg.KeepAliveTimeout)
	assert.Nil(t, cfg.EnabledRedirect)
}

func TestDialer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(""))
	}))
	defer server.Close()
	assert.NotNil(t, server)
	address, _ := url.Parse(server.URL)
	assert.NotNil(t, address)

	cfg := DefaultConfig.copy()
	d := newDialer(cfg)
	assert.NotNil(t, d)
	assert.Equal(t, DefaultConnectTimeout, d.Dialer.Timeout)
	assert.Equal(t, DefaultKeepAliveTimeout, d.Dialer.KeepAlive)

	dCtxFn := d.DialContext
	assert.NotNil(t, dCtxFn)
	c, err := dCtxFn(context.TODO(), "tcp", address.Host)
	assert.Nil(t, err)
	assert.NotNil(t, c)
	conn, ok := c.(*timeoutConn)
	assert.True(t, ok)
	assert.Equal(t, DefaultReadWriteTimeout, conn.timeout)

	readwriteTimeout := 200 * time.Second
	ctx := context.WithValue(context.TODO(), "OpReadWriteTimeout", &readwriteTimeout)
	c, err = dCtxFn(ctx, "tcp", address.Host)
	assert.Nil(t, err)
	assert.NotNil(t, c)
	conn, ok = c.(*timeoutConn)
	assert.True(t, ok)
	assert.Equal(t, readwriteTimeout, conn.timeout)

	readwriteTimeout = 200 * time.Second
	ctx = context.WithValue(context.TODO(), "ReadWriteTimeout", &readwriteTimeout)
	c, err = dCtxFn(ctx, "tcp", address.Host)
	assert.Nil(t, err)
	assert.NotNil(t, c)
	conn, ok = c.(*timeoutConn)
	assert.True(t, ok)
	assert.Equal(t, DefaultReadWriteTimeout, conn.timeout)

	cfg = DefaultConfig.copy()
	d = newDialer(cfg)
	assert.NotNil(t, d)
	assert.Equal(t, DefaultConnectTimeout, d.Dialer.Timeout)
	assert.Equal(t, DefaultKeepAliveTimeout, d.Dialer.KeepAlive)
	assert.Nil(t, d.postWrite)
	assert.Nil(t, d.postWrite)

	cfg = DefaultConfig.copy()
	cfg.PostRead = []func(n int, err error){
		func(n int, err error) {},
	}
	d = newDialer(cfg)
	assert.NotNil(t, d)
	assert.Equal(t, DefaultConnectTimeout, d.Dialer.Timeout)
	assert.Equal(t, DefaultKeepAliveTimeout, d.Dialer.KeepAlive)
	assert.Len(t, d.postRead, 1)
	assert.Nil(t, d.postWrite)

	cfg = DefaultConfig.copy()
	cfg.PostWrite = []func(n int, err error){
		func(n int, err error) {},
	}
	d = newDialer(cfg)
	assert.NotNil(t, d)
	assert.Equal(t, DefaultConnectTimeout, d.Dialer.Timeout)
	assert.Equal(t, DefaultKeepAliveTimeout, d.Dialer.KeepAlive)
	assert.Len(t, d.postWrite, 1)
	assert.Nil(t, d.postRead)
}

func TestTransport(t *testing.T) {
	cfg := DefaultConfig.copy()
	tr := newTransportCustom(cfg)
	assert.NotNil(t, tr)
	transport, ok := tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.Equal(t, DefaultConnectTimeout, transport.TLSHandshakeTimeout)
	assert.Equal(t, DefaultIdleConnectionTimeout, transport.IdleConnTimeout)
	assert.Equal(t, DefaultMaxConnections, transport.MaxConnsPerHost)
	assert.Equal(t, DefaultExpectContinueTimeout, transport.ExpectContinueTimeout)
	assert.Equal(t, DefaultTLSMinVersion, transport.TLSClientConfig.MinVersion)
	assert.Equal(t, false, transport.TLSClientConfig.InsecureSkipVerify)

	tr = newTransportCustom(cfg,
		InsecureSkipVerify(true),
		MaxConnections(5),
		ExpectContinueTimeout(300*time.Second),
		TLSMinVersion(tls.VersionTLS13))
	assert.NotNil(t, tr)
	transport, ok = tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.Equal(t, DefaultConnectTimeout, transport.TLSHandshakeTimeout)
	assert.Equal(t, DefaultIdleConnectionTimeout, transport.IdleConnTimeout)
	assert.Equal(t, 5, transport.MaxConnsPerHost)
	assert.Equal(t, 300*time.Second, transport.ExpectContinueTimeout)
	assert.Equal(t, uint16(tls.VersionTLS13), transport.TLSClientConfig.MinVersion)
	assert.Equal(t, true, transport.TLSClientConfig.InsecureSkipVerify)
}

func TestTransportWithProxy(t *testing.T) {
	cfg := DefaultConfig.copy()
	tr := newTransportCustom(cfg)
	assert.NotNil(t, tr)
	transport, ok := tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.Nil(t, transport.Proxy)

	proxyUrl, _ := url.Parse("http://127.0.0.1")
	tr = newTransportCustom(cfg,
		HttpProxy(proxyUrl))
	assert.NotNil(t, tr)
	transport, ok = tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.NotNil(t, transport.Proxy)

	tr = newTransportCustom(cfg,
		HttpProxy(proxyUrl),
		ProxyFromEnvironment())
	assert.NotNil(t, tr)
	transport, ok = tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.NotNil(t, transport.Proxy)
}

func TestHttpClient(t *testing.T) {
	client := &http.Client{}
	assert.NotNil(t, client)
	assert.Nil(t, client.CheckRedirect)

	//default
	client = NewHttpClient(nil)
	assert.NotNil(t, client)
	tr := client.Transport
	assert.NotNil(t, tr)
	transport, ok := tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.Equal(t, DefaultConnectTimeout, transport.TLSHandshakeTimeout)
	assert.Equal(t, DefaultIdleConnectionTimeout, transport.IdleConnTimeout)
	assert.Equal(t, DefaultMaxConnections, transport.MaxConnsPerHost)
	assert.Equal(t, DefaultExpectContinueTimeout, transport.ExpectContinueTimeout)
	assert.Equal(t, DefaultTLSMinVersion, transport.TLSClientConfig.MinVersion)
	assert.Equal(t, false, transport.TLSClientConfig.InsecureSkipVerify)
	assert.NotNil(t, client.CheckRedirect)

	//has value
	cfg := DefaultConfig.copy()
	client = NewHttpClient(cfg)
	assert.NotNil(t, client)
	tr = client.Transport
	assert.NotNil(t, tr)
	transport, ok = tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.Equal(t, DefaultConnectTimeout, transport.TLSHandshakeTimeout)
	assert.Equal(t, DefaultIdleConnectionTimeout, transport.IdleConnTimeout)
	assert.Equal(t, DefaultMaxConnections, transport.MaxConnsPerHost)
	assert.Equal(t, DefaultExpectContinueTimeout, transport.ExpectContinueTimeout)
	assert.Equal(t, DefaultTLSMinVersion, transport.TLSClientConfig.MinVersion)
	assert.Equal(t, false, transport.TLSClientConfig.InsecureSkipVerify)
	assert.NotNil(t, client.CheckRedirect)

	//over write
	cfg = DefaultConfig.copy(
		&Config{
			EnabledRedirect: ptr(true),
		},
		&Config{
			ConnectTimeout: ptr(100 * time.Second),
		})
	client = NewHttpClient(cfg,
		InsecureSkipVerify(true),
		MaxConnections(5),
		ExpectContinueTimeout(300*time.Second),
		TLSMinVersion(tls.VersionTLS13))
	assert.NotNil(t, client)
	tr = client.Transport
	assert.NotNil(t, tr)
	transport, ok = tr.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.Equal(t, 100*time.Second, transport.TLSHandshakeTimeout)
	assert.Equal(t, DefaultIdleConnectionTimeout, transport.IdleConnTimeout)
	assert.Equal(t, 5, transport.MaxConnsPerHost)
	assert.Equal(t, 300*time.Second, transport.ExpectContinueTimeout)
	assert.Equal(t, uint16(tls.VersionTLS13), transport.TLSClientConfig.MinVersion)
	assert.Equal(t, true, transport.TLSClientConfig.InsecureSkipVerify)
	assert.Nil(t, client.CheckRedirect)
}
