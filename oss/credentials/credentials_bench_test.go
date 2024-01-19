package credentials

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkCredentialsFetcher_GetCredentials(b *testing.B) {
	fetcher := CredentialsFetcherFunc(func(ctx context.Context) (Credentials, error) {
		return Credentials{
			AccessKeyID:     "ak",
			AccessKeySecret: "sk",
		}, nil
	})

	cases := []int{1, 10, 100, 500, 1000, 10000}
	for _, c := range cases {
		b.Run(strconv.Itoa(c), func(b *testing.B) {
			p := NewCredentialsFetcherProvider(fetcher)
			var wg sync.WaitGroup
			wg.Add(c)
			for i := 0; i < c; i++ {
				go func() {
					for j := 0; j < b.N; j++ {
						v, err := p.GetCredentials(context.Background())
						if err != nil {
							b.Errorf("expect no error %v, %v", v, err)
						}
					}
					wg.Done()
				}()
			}
			b.ResetTimer()

			wg.Wait()
		})
	}
}

func BenchmarkCredentialsFetcher_GetCredentials_Expires(b *testing.B) {
	count := int64(0)
	fetcher := CredentialsFetcherFunc(func(ctx context.Context) (Credentials, error) {
		time.Sleep(time.Millisecond)
		count++
		return Credentials{
			AccessKeyID:     "ak",
			AccessKeySecret: "sk",
			Expires:         ptr(time.Now().Add(10 * time.Second)),
		}, nil
	})

	expRates := []int{10000, 1000, 100}
	cases := []int{1, 10, 100, 500, 1000, 10000}
	for _, expRate := range expRates {
		for _, c := range cases {
			b.Run(fmt.Sprintf("%d-%d", expRate, c), func(b *testing.B) {
				p := NewCredentialsFetcherProvider(fetcher, func(o *CredentialsFetcherOptions) {
					o.ExpiredFactor = 0.6
					o.RefreshDuration = 1 * time.Second
				})
				fetcherProvider, _ := p.(*CredentialsFetcherProvider)
				var wg sync.WaitGroup
				wg.Add(c)
				for i := 0; i < c; i++ {
					go func(id int) {
						for j := 0; j < b.N; j++ {
							v, err := p.GetCredentials(context.Background())
							if err != nil {
								b.Errorf("expect no error %v, %v", v, err)
							}
							// periodically expire creds to cause rwlock
							if id == 0 && j%expRate == 0 {
								fetcherProvider.credentials.Store((*Credentials)(nil))
							}
						}
						wg.Done()
					}(i)
				}
				b.ResetTimer()

				wg.Wait()
			})
		}
	}
}
