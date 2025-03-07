package fetcher

import (
	"context"
	"net/http"
	"time"
)

type Fetcher struct {
	limit chan struct{}
}

func New(concurrencyLimit int) *Fetcher {
	return &Fetcher{
		limit: make(chan struct{}, concurrencyLimit),
	}
}

// Fetch requests given URL using given HTTP method
func (f *Fetcher) Fetch(method, url string) (*http.Response, error) {
	f.limit <- struct{}{}
	defer func() { <-f.limit }()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, method, url, nil)

	return http.DefaultClient.Do(req)
}
