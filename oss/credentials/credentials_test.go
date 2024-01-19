package credentials

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

type AtomicBool struct {
	value atomic.Value
}

func NewAtomicBool(initialValue bool) *AtomicBool {
	b := &AtomicBool{}
	b.value.Store(initialValue)
	return b
}

func (b *AtomicBool) Load() bool {
	return b.value.Load().(bool)
}

func (b *AtomicBool) Store(newValue bool) {
	b.value.Store(newValue)
}

func TestCredentials(t *testing.T) {
	cred := &Credentials{}
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.False(t, cred.HasKeys())

	expires := time.Now().Add(10 * time.Second)
	cred = &Credentials{
		AccessKeyID:     "ak",
		AccessKeySecret: "sk",
		Expires:         &expires,
	}
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())

	expires = time.Now().Add(-10 * time.Second)
	cred = &Credentials{
		AccessKeyID:     "ak",
		AccessKeySecret: "sk",
		Expires:         &expires,
	}
	assert.NotNil(t, cred)
	assert.True(t, cred.Expired())
	assert.True(t, cred.HasKeys())
}

func TestStaticCredentialsProvider(t *testing.T) {
	provider := NewStaticCredentialsProvider("ak", "sk")
	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())
	assert.Equal(t, "ak", cred.AccessKeyID)
	assert.Equal(t, "sk", cred.AccessKeySecret)
	assert.Equal(t, "", cred.SecurityToken)

	provider = NewStaticCredentialsProvider("ak1", "sk1", "token1")
	cred, err = provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())
	assert.Equal(t, "ak1", cred.AccessKeyID)
	assert.Equal(t, "sk1", cred.AccessKeySecret)
	assert.Equal(t, "token1", cred.SecurityToken)
}

func TestEnvironmentVariableCredentialsProvider(t *testing.T) {
	provider := NewEnvironmentVariableCredentialsProvider()
	assert.NotNil(t, provider)

	oriak := os.Getenv("OSS_ACCESS_KEY_ID")
	orisk := os.Getenv("OSS_ACCESS_KEY_SECRET")
	oritk := os.Getenv("OSS_SESSION_TOKEN")

	defer func() {
		if oriak == "" {
			os.Unsetenv("OSS_ACCESS_KEY_ID")
		} else {
			os.Setenv("OSS_ACCESS_KEY_ID", oriak)
		}
		if orisk == "" {
			os.Unsetenv("OSS_ACCESS_KEY_SECRET")
		} else {
			os.Setenv("OSS_ACCESS_KEY_SECRET", orisk)
		}
		if oritk == "" {
			os.Unsetenv("OSS_SESSION_TOKEN")
		} else {
			os.Setenv("OSS_SESSION_TOKEN", oritk)
		}
	}()

	os.Setenv("OSS_ACCESS_KEY_ID", "myak")
	os.Setenv("OSS_ACCESS_KEY_SECRET", "mysk")
	provider = NewEnvironmentVariableCredentialsProvider()
	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())
	assert.Equal(t, "myak", cred.AccessKeyID)
	assert.Equal(t, "mysk", cred.AccessKeySecret)
	assert.Equal(t, "", cred.SecurityToken)

	err = os.Setenv("OSS_SESSION_TOKEN", "mytoken")

	provider = NewEnvironmentVariableCredentialsProvider()
	cred, err = provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())
	assert.Equal(t, "myak", cred.AccessKeyID)
	assert.Equal(t, "mysk", cred.AccessKeySecret)
	assert.Equal(t, "mytoken", cred.SecurityToken)
}

func TestAnonymousCredentialsProvider(t *testing.T) {
	provider := NewAnonymousCredentialsProvider()
	assert.NotNil(t, provider)

	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.HasKeys())
	assert.False(t, cred.Expired())
}

type stubCredentialsFetcher struct {
	delay time.Duration
	token string
	count int64
}

func (s *stubCredentialsFetcher) Fetch(ctx context.Context) (Credentials, error) {
	var expires *time.Time
	if s.delay > 0 {
		now := time.Now()
		new := now.Add(s.delay)
		expires = &new
	}

	atomic.AddInt64(&s.count, 1)

	return Credentials{
		AccessKeyID:     "ak",
		AccessKeySecret: "sk",
		SecurityToken:   s.token,
		Expires:         expires,
	}, nil
}

func TestCredentialsFetcherProvider(t *testing.T) {
	provider := NewCredentialsFetcherProvider(nil)
	assert.NotNil(t, provider)
	fetcherProvider, ok := provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	assert.Equal(t, defaultExpiredFactor, fetcherProvider.expiredFactor)
	assert.Equal(t, defaultRefreshDuration, fetcherProvider.refreshDuration)
	assert.Nil(t, fetcherProvider.fetcher)

	_, err := provider.GetCredentials(context.TODO())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "fetcher is null.")

	provider = NewCredentialsFetcherProvider(nil, func(o *CredentialsFetcherOptions) {
		o.ExpiredFactor = 0.7
		o.RefreshDuration = 1 * time.Second
	})
	assert.NotNil(t, provider)
	fetcherProvider, ok = provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	assert.Equal(t, 0.7, fetcherProvider.expiredFactor)
	assert.Equal(t, 1*time.Second, fetcherProvider.refreshDuration)

	provider = NewCredentialsFetcherProvider(&stubCredentialsFetcher{})
	assert.NotNil(t, provider)
	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred.AccessKeyID)
	assert.Equal(t, "sk", cred.AccessKeySecret)
	assert.False(t, cred.Expired())

	// with Expired
	provider = NewCredentialsFetcherProvider(&stubCredentialsFetcher{
		token: "token",
		delay: 2 * time.Second,
	})
	assert.NotNil(t, provider)
	fetcherProvider, ok = provider.(*CredentialsFetcherProvider)
	assert.NotNil(t, fetcherProvider)

	// 1st
	cred1, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred1.AccessKeyID)
	assert.Equal(t, "sk", cred1.AccessKeySecret)
	assert.Equal(t, "token", cred1.SecurityToken)
	assert.NotNil(t, cred1.Expires)
	assert.False(t, cred1.Expired())

	// 2st
	cred2, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred2.AccessKeyID)
	assert.Equal(t, "sk", cred2.AccessKeySecret)
	assert.Equal(t, "token", cred2.SecurityToken)
	assert.Equal(t, cred1.Expires, cred2.Expires)

	time.Sleep(3 * time.Second)
	assert.True(t, cred1.Expired())

	// 3st
	cred3, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred3.AccessKeyID)
	assert.Equal(t, "sk", cred3.AccessKeySecret)
	assert.Equal(t, "token", cred3.SecurityToken)
	assert.False(t, cred3.Expired())

	assert.True(t, cred3.Expires.After(*cred1.Expires))

}

func TestCredentialsFetcherProvider_Soon(t *testing.T) {
	// with Expired
	provider := NewCredentialsFetcherProvider(
		&stubCredentialsFetcher{
			token: "token",
			delay: 10 * time.Second,
		},
		func(o *CredentialsFetcherOptions) {
			o.ExpiredFactor = 0.4
			o.RefreshDuration = 1 * time.Second
		},
	)
	assert.NotNil(t, provider)
	fetcherProvider, ok := provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)

	// 1st
	cred1, err := provider.GetCredentials(context.TODO())
	cred1_1, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred1.AccessKeyID)
	assert.Equal(t, "sk", cred1.AccessKeySecret)
	assert.Equal(t, "token", cred1.SecurityToken)
	assert.NotNil(t, cred1.Expires)
	assert.False(t, cred1.Expired())
	assert.EqualValues(t, cred1, cred1_1)

	// 2st
	time.Sleep(6 * time.Second)
	assert.False(t, cred1.Expired())
	cred2, err := provider.GetCredentials(context.TODO())
	cred3, _ := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred2.AccessKeyID)
	assert.Equal(t, "sk", cred2.AccessKeySecret)
	assert.Equal(t, "token", cred2.SecurityToken)
	assert.True(t, cred2.Expires.After(*cred1.Expires))
	assert.EqualValues(t, cred2, cred3)
}

func TestCredentialsFetcherProvider_MultiJobs(t *testing.T) {
	// with Expired
	fetcher := &stubCredentialsFetcher{
		token: "token",
		delay: 10 * time.Second,
	}

	provider := NewCredentialsFetcherProvider(
		fetcher,
		func(o *CredentialsFetcherOptions) {
			o.ExpiredFactor = 0.4
			o.RefreshDuration = 1 * time.Second
		},
	)
	assert.NotNil(t, provider)
	fetcherProvider, ok := provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)

	var run = NewAtomicBool(false)
	run.Store(true)
	testFn := func() {
		count := int64(0)
		for run.Load() {
			cred, err := provider.GetCredentials(context.TODO())
			assert.Nil(t, err)
			assert.Equal(t, "ak", cred.AccessKeyID)
			assert.Equal(t, "sk", cred.AccessKeySecret)
			assert.Equal(t, "token", cred.SecurityToken)
			assert.NotNil(t, cred.Expires)
			assert.False(t, cred.Expired())
			count++
		}
		assert.Greater(t, count, int64(5000))
	}

	for i := 0; i < 20; i++ {
		go testFn()
	}

	time.Sleep(15 * time.Second)
	run.Store(false)
	//assert.Less(t, atomic.LoadInt64(&fetcher.count), int64(6)*2)
}

type stubCredentialsFetcher2 struct {
	delay        time.Duration
	token        string
	returnErr    bool
	returnTimout AtomicBool
}

func (s *stubCredentialsFetcher2) Fetch(ctx context.Context) (Credentials, error) {
	var expires *time.Time
	if s.delay > 0 {
		now := time.Now()
		new := now.Add(s.delay)
		expires = &new
	}
	if s.returnTimout.Load() {
		time.Sleep(10 * time.Second)
		return Credentials{}, fmt.Errorf("returnTimout")
	} else if s.returnErr {
		return Credentials{}, fmt.Errorf("returnErr")
	} else {
		return Credentials{
			AccessKeyID:     "ak",
			AccessKeySecret: "sk",
			SecurityToken:   s.token,
			Expires:         expires,
		}, nil
	}
}

func TestCredentialsFetcherProvider_Error(t *testing.T) {
	fetcher := &stubCredentialsFetcher2{
		token:     "token",
		delay:     10 * time.Second,
		returnErr: true,
	}
	fetcher.returnTimout.Store(false)
	provider := NewCredentialsFetcherProvider(
		fetcher,
		func(o *CredentialsFetcherOptions) {
			o.ExpiredFactor = 0.4
			o.RefreshDuration = 1 * time.Second
		},
	)
	assert.NotNil(t, provider)
	fetcherProvider, ok := provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)

	// Get Fail
	_, err := provider.GetCredentials(context.TODO())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "returnErr")

	// Get OK
	fetcher.returnErr = false
	cred1, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred1.AccessKeyID)
	assert.Equal(t, "sk", cred1.AccessKeySecret)
	assert.Equal(t, "token", cred1.SecurityToken)
	assert.NotNil(t, cred1.Expires)
	assert.False(t, cred1.Expired())

	// 2st Fail
	fetcher.returnErr = true
	time.Sleep(6 * time.Second)
	assert.False(t, cred1.Expired())
	cred2, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred2.AccessKeyID)
	assert.Equal(t, "sk", cred2.AccessKeySecret)
	assert.Equal(t, "token", cred2.SecurityToken)
	assert.Equal(t, *cred1.Expires, *cred2.Expires)

	// Fetch Timeout
	fetcher = &stubCredentialsFetcher2{
		token: "token",
		delay: 6 * time.Second,
	}

	provider = NewCredentialsFetcherProvider(
		fetcher,
		func(o *CredentialsFetcherOptions) {
			o.ExpiredFactor = 0.4
			o.RefreshDuration = 1 * time.Second
		},
	)
	assert.NotNil(t, provider)
	fetcherProvider, ok = provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	ctxt1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel1()
	fetcher.returnTimout.Store(true)
	_, err = provider.GetCredentials(ctxt1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "FetchCredentialsCanceled")

	fetcher.returnTimout.Store(false)
	cred3, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred3.AccessKeyID)
	assert.Equal(t, "sk", cred3.AccessKeySecret)
	assert.Equal(t, "token", cred3.SecurityToken)
	assert.NotNil(t, cred3.Expires)
	assert.False(t, cred3.Expired())

	time.Sleep(4 * time.Second)
	ctxt2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()
	fetcher.returnTimout.Store(true)
	cred4, err := provider.GetCredentials(ctxt2)
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred4.AccessKeyID)
	assert.Equal(t, "sk", cred4.AccessKeySecret)
	assert.Equal(t, "token", cred4.SecurityToken)
	assert.NotNil(t, cred4.Expires)
	assert.Equal(t, *cred3.Expires, *cred4.Expires)
}

func createFileFromByte(t *testing.T, fileName string, content []byte) {
	fout, err := os.Create(fileName)
	assert.Nil(t, err)
	defer fout.Close()
	_, err = fout.Write(content)
	assert.Nil(t, err)
}

func TestProcessCredentialsProvider(t *testing.T) {
	//default
	p := NewProcessCredentialsProvider("")
	processProvider, _ := p.(*ProcessCredentialsProvider)
	assert.Equal(t, 15*time.Second, processProvider.timeout)
	_, err := p.GetCredentials(context.TODO())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "command must not be empty")

	// set timeout
	p = NewProcessCredentialsProvider("", func(pcpo *ProcessCredentialsProviderOptions) {
		pcpo.Timeout = 5 * time.Minute
	})
	processProvider, _ = p.(*ProcessCredentialsProvider)
	assert.Equal(t, 5*time.Minute, processProvider.timeout)

	//run cmd
	localFile := fmt.Sprintf("cred-file-%v-", time.Now().UnixMicro()) + ".tmp"
	defer func() {
		os.Remove(localFile)
	}()
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = fmt.Sprintf("type %s", localFile)
	} else {
		cmd = fmt.Sprintf("cat %s", localFile)
	}

	// all fileds
	data := `
	{
		"AccessKeyId" : "ak",
		"AccessKeySecret" : "sk",
		"Expiration" : "2023-12-29T07:45:02Z",
		"SecurityToken" : "token"
	}`
	createFileFromByte(t, localFile, []byte(data))
	p = NewProcessCredentialsProvider(cmd)
	cred, err := p.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred.AccessKeyID)
	assert.Equal(t, "sk", cred.AccessKeySecret)
	assert.Equal(t, "token", cred.SecurityToken)
	assert.NotNil(t, cred.Expires)

	// only ak, sk
	data = `
	{
		"AccessKeyId" : "ak",
		"AccessKeySecret" : "sk"
	}`
	createFileFromByte(t, localFile, []byte(data))
	p = NewProcessCredentialsProvider(cmd)
	cred, err = p.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred.AccessKeyID)
	assert.Equal(t, "sk", cred.AccessKeySecret)
	assert.Equal(t, "", cred.SecurityToken)
	assert.Nil(t, cred.Expires)

	// only ak or sk, gets error
	data = `
	{
		"AccessKeyId" : "ak"
	}`
	createFileFromByte(t, localFile, []byte(data))
	p = NewProcessCredentialsProvider(cmd)
	cred, err = p.GetCredentials(context.TODO())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing AccessKeyId or AccessKeySecret in process output")

	data = `
	{
		"AccessKeySecret" : "sk"
	}`
	createFileFromByte(t, localFile, []byte(data))
	p = NewProcessCredentialsProvider(cmd)
	cred, err = p.GetCredentials(context.TODO())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing AccessKeyId or AccessKeySecret in process output")

	// invalid json
	data = `
	{
		"AccessKeyId" : "ak",
		"AccessKeySecret" : "sk"
	`
	createFileFromByte(t, localFile, []byte(data))
	p = NewProcessCredentialsProvider(cmd)
	cred, err = p.GetCredentials(context.TODO())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unexpected end of JSON input")

	// invalid command
	data = `
	{
		"AccessKeyId" : "ak",
		"AccessKeySecret" : "sk"
	}`
	createFileFromByte(t, localFile, []byte(data))
	p = NewProcessCredentialsProvider("invalid cmd")
	cred, err = p.GetCredentials(context.TODO())
	assert.Contains(t, err.Error(), "error in credential_process")
}

func TestMixedCredentialsProvider(t *testing.T) {
	// ProcessCredentialsProvider + CredentialsFetcherProvider
	//run cmd
	localFile := fmt.Sprintf("cred-file-%v-", time.Now().UnixMicro()) + ".tmp"
	defer func() {
		os.Remove(localFile)
	}()
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = fmt.Sprintf("type %s", localFile)
	} else {
		cmd = fmt.Sprintf("cat %s", localFile)
	}

	data := `
	{
		"AccessKeyId" : "ak",
		"AccessKeySecret" : "sk",
		"Expiration" : "2023-12-29T07:45:02Z",
		"SecurityToken" : "token"
	}`

	createFileFromByte(t, localFile, []byte(data))
	provider := NewCredentialsFetcherProvider(CredentialsFetcherFunc(func(ctx context.Context) (Credentials, error) {
		return NewProcessCredentialsProvider(cmd).GetCredentials(ctx)
	}))
	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred.AccessKeyID)
	assert.Equal(t, "sk", cred.AccessKeySecret)
	assert.Equal(t, "token", cred.SecurityToken)
	assert.NotNil(t, cred.Expires)
}

func testSetupEcsRoleMockServer(t *testing.T, nowTime time.Time, isExpiration bool) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isExpiration {
			nowTime = time.Now()
		}
		time.Sleep(4 * time.Second)
		switch r.Method {
		case "GET":
			if r.URL.Path == "/latest/meta-data/ram/security-credentials/" {
				responseData := []byte(`EcsRamRoleTest`)
				w.Write(responseData)
			} else if r.URL.Path == "/latest/meta-data/ram/security-credentials/EcsRamRoleTest" {
				update := nowTime.UTC().Format("2006-01-02T15:04:05Z")
				expiration := nowTime.Add(+5 * time.Second).UTC().Format("2006-01-02T15:04:05Z")
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"AccessKeyId": "accessKeyId","AccessKeySecret": "accessKeySecret","SecurityToken": "securityToken","Expiration": "`+expiration+`","LastUpdated" : "`+update+`","Code" : "Success"}`)
			}
		}
	}))
	return server
}

func TestNewEcsRoleCredentialsProviderWithoutRefresh(t *testing.T) {
	defaultProvider := NewEcsRoleCredentialsProviderWithoutRefresh()
	assert.NotNil(t, defaultProvider)
	ecsRoleProvider, ok := defaultProvider.(*ecsRoleCredentialsProvider)
	assert.True(t, ok)
	assert.NotNil(t, ecsRoleProvider)
	assert.Equal(t, "", ecsRoleProvider.ramRole)
	assert.Equal(t, 3, ecsRoleProvider.retries)
	assert.Equal(t, 10*time.Second, ecsRoleProvider.timeout)

	defaultProvider = NewEcsRoleCredentialsProviderWithoutRefresh(EcsRamRole("EcsRamRoleTest"))
	assert.NotNil(t, defaultProvider)
	ecsRoleProvider, ok = defaultProvider.(*ecsRoleCredentialsProvider)
	assert.True(t, ok)
	assert.NotNil(t, ecsRoleProvider)
	assert.Equal(t, "EcsRamRoleTest", ecsRoleProvider.ramRole)
	assert.Equal(t, 3, ecsRoleProvider.retries)
	assert.Equal(t, 10*time.Second, ecsRoleProvider.timeout)

	defaultProvider = NewEcsRoleCredentialsProviderWithoutRefresh(func(o *EcsRoleCredentialsProviderOptions) {
		o.Timeout = 15 * time.Second
		o.Retries = 10
		o.RamRole = "ramosstest"
	})
	assert.NotNil(t, defaultProvider)
	ecsRoleProvider, ok = defaultProvider.(*ecsRoleCredentialsProvider)
	assert.True(t, ok)
	assert.NotNil(t, ecsRoleProvider)
	assert.Equal(t, "ramosstest", ecsRoleProvider.ramRole)
	assert.Equal(t, 10, ecsRoleProvider.retries)
	assert.Equal(t, 15*time.Second, ecsRoleProvider.timeout)

	nowTime := time.Now()
	server := testSetupEcsRoleMockServer(t, nowTime, false)
	defer server.Close()
	provider := newStubEcsRoleCredentialsProviderWithoutRefresh(stubEcsRamRole("EcsRamRoleTest"), func(o *stubEcsRoleCredentialsProviderOptions) {
		o.timeout = 15 * time.Second
		o.retries = 3
		o.ramCredUrl = server.URL + "/latest/meta-data/ram/security-credentials/"
	})
	ctx := context.Background()
	creds, err := provider.GetCredentials(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "accessKeyId", creds.AccessKeyID)
	assert.Equal(t, "accessKeySecret", creds.AccessKeySecret)
	assert.Equal(t, "securityToken", creds.SecurityToken)
	assert.Equal(t, nowTime.Add(+5*time.Second).UTC().Format("2006-01-02T15:04:05Z"), (*creds.Expires).Format("2006-01-02T15:04:05Z"))

	provider = newStubEcsRoleCredentialsProviderWithoutRefresh(func(o *stubEcsRoleCredentialsProviderOptions) {
		o.timeout = 15 * time.Second
		o.retries = 3
		o.ramCredUrl = server.URL + "/latest/meta-data/ram/security-credentials/"
	})
	creds, err = provider.GetCredentials(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "accessKeyId", creds.AccessKeyID)
	assert.Equal(t, "accessKeySecret", creds.AccessKeySecret)
	assert.Equal(t, "securityToken", creds.SecurityToken)
	assert.Equal(t, nowTime.Add(+5*time.Second).UTC().Format("2006-01-02T15:04:05Z"), (*creds.Expires).Format("2006-01-02T15:04:05Z"))

	ctxt1, cancel1 := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel1()
	creds, err = provider.GetCredentials(ctxt1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")

	provider = newStubEcsRoleCredentialsProviderWithoutRefresh(stubEcsRamRole("EcsRamRoleTest"), func(o *stubEcsRoleCredentialsProviderOptions) {
		o.timeout = 1 * time.Second
		o.retries = 2
		o.ramCredUrl = server.URL + "/latest/meta-data/ram/security-credentials/"
	})
	creds, err = provider.GetCredentials(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestNewEcsRoleCredentialsProvider(t *testing.T) {
	ecsRoleProvider := NewEcsRoleCredentialsProvider()
	assert.NotNil(t, ecsRoleProvider)
	fetcherProvider, ok := ecsRoleProvider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	assert.Equal(t, defaultExpiredFactor, fetcherProvider.expiredFactor)
	assert.Equal(t, defaultRefreshDuration, fetcherProvider.refreshDuration)
	assert.NotNil(t, fetcherProvider.fetcher)

	ecsRoleProvider = NewEcsRoleCredentialsProvider(EcsRamRole("EcsRamRoleTest"))
	assert.NotNil(t, ecsRoleProvider)
	fetcherProvider, ok = ecsRoleProvider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	assert.Equal(t, defaultExpiredFactor, fetcherProvider.expiredFactor)
	assert.Equal(t, defaultRefreshDuration, fetcherProvider.refreshDuration)
	assert.NotNil(t, fetcherProvider.fetcher)

	nowTime := time.Now()
	server := testSetupEcsRoleMockServer(t, nowTime, false)
	defer server.Close()
	provider := newStubEcsRoleCredentialsProvider(stubEcsRamRole("EcsRamRoleTest"), func(o *stubEcsRoleCredentialsProviderOptions) {
		o.ramCredUrl = server.URL + "/latest/meta-data/ram/security-credentials/"
	})
	ctx := context.Background()
	creds, err := provider.GetCredentials(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "accessKeyId", creds.AccessKeyID)
	assert.Equal(t, "accessKeySecret", creds.AccessKeySecret)
	assert.Equal(t, "securityToken", creds.SecurityToken)
	assert.Equal(t, nowTime.Add(+5*time.Second).UTC().Format("2006-01-02T15:04:05Z"), (*creds.Expires).Format("2006-01-02T15:04:05Z"))

	provider = newStubEcsRoleCredentialsProvider(func(o *stubEcsRoleCredentialsProviderOptions) {
		o.ramCredUrl = server.URL + "/latest/meta-data/ram/security-credentials/"
	})
	creds, err = provider.GetCredentials(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "accessKeyId", creds.AccessKeyID)
	assert.Equal(t, "accessKeySecret", creds.AccessKeySecret)
	assert.Equal(t, "securityToken", creds.SecurityToken)
	assert.Equal(t, nowTime.Add(+5*time.Second).UTC().Format("2006-01-02T15:04:05Z"), (*creds.Expires).Format("2006-01-02T15:04:05Z"))

	ctxt1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1()
	creds, err = provider.GetCredentials(ctxt1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FetchCredentialsCanceled")

	provider = newStubEcsRoleCredentialsProvider(stubEcsRamRole("EcsRamRoleTest"), func(o *stubEcsRoleCredentialsProviderOptions) {
		o.ramCredUrl = server.URL + "/latest/meta-data/ram/security-credentials/"
		o.timeout = 3 * time.Second
		o.retries = 3
	})
	creds, err = provider.GetCredentials(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")

	nowTime = time.Now()
	server1 := testSetupEcsRoleMockServer(t, nowTime, true)
	defer server1.Close()
	provider = newStubEcsRoleCredentialsProvider(stubEcsRamRole("EcsRamRoleTest"), func(o *stubEcsRoleCredentialsProviderOptions) {
		o.ramCredUrl = server1.URL + "/latest/meta-data/ram/security-credentials/"
		o.timeout = 15 * time.Second
		o.retries = 3
	})
	assert.NotNil(t, provider)
	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "accessKeyId", cred.AccessKeyID)
	assert.Equal(t, "accessKeySecret", cred.AccessKeySecret)
	assert.Equal(t, "securityToken", cred.SecurityToken)
	assert.NotNil(t, cred.Expires)
	assert.False(t, cred.Expired())

	// 1st
	cred1, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "accessKeyId", cred1.AccessKeyID)
	assert.Equal(t, "accessKeySecret", cred1.AccessKeySecret)
	assert.Equal(t, "securityToken", cred1.SecurityToken)
	assert.NotNil(t, cred1.Expires)
	assert.False(t, cred1.Expired())

	// 2st
	cred2, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "accessKeyId", cred2.AccessKeyID)
	assert.Equal(t, "accessKeySecret", cred2.AccessKeySecret)
	assert.Equal(t, "securityToken", cred2.SecurityToken)
	assert.Equal(t, cred1.Expires, cred2.Expires)
	time.Sleep(4 * time.Second)
	assert.True(t, cred.Expired())
	// 3st
	cred3, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "accessKeyId", cred3.AccessKeyID)
	assert.Equal(t, "accessKeySecret", cred3.AccessKeySecret)
	assert.Equal(t, "securityToken", cred3.SecurityToken)
	assert.False(t, cred3.Expired())
	assert.True(t, cred3.Expires.After(*cred.Expires))

	provider = newStubEcsRoleCredentialsProvider(stubEcsRamRole("EcsRamRoleTest"), func(o *stubEcsRoleCredentialsProviderOptions) {
		o.ramCredUrl = server1.URL + "/latest/meta-data/ram/security-credentials/"
		o.timeout = 2 * time.Second
		o.retries = 2
	})
	assert.NotNil(t, provider)
	cred, err = provider.GetCredentials(context.TODO())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

type stubEcsRoleCredentialsProviderOptions ecsRoleCredentialsProvider

func stubEcsRamRole(ramRole string) func(*stubEcsRoleCredentialsProviderOptions) {
	return func(options *stubEcsRoleCredentialsProviderOptions) {
		options.ramRole = ramRole
	}
}

func newStubEcsRoleCredentialsProviderWithoutRefresh(optFns ...func(*stubEcsRoleCredentialsProviderOptions)) CredentialsProvider {
	options := stubEcsRoleCredentialsProviderOptions{
		ramRole: "",
		timeout: time.Second * 10,
		retries: 3,
	}
	for _, fn := range optFns {
		fn(&options)
	}
	return &ecsRoleCredentialsProvider{
		ramCredUrl: options.ramCredUrl,
		ramRole:    options.ramRole,
		timeout:    options.timeout,
		retries:    options.retries,
	}
}

func newStubEcsRoleCredentialsProvider(optFns ...func(*stubEcsRoleCredentialsProviderOptions)) CredentialsProvider {
	p := newStubEcsRoleCredentialsProviderWithoutRefresh(optFns...)
	provider := NewCredentialsFetcherProvider(CredentialsFetcherFunc(func(ctx context.Context) (Credentials, error) {
		return p.GetCredentials(ctx)
	}))
	return provider
}
