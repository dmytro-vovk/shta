package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"os"
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
func (f *Fetcher) Fetch(method, url string) (resp *http.Response, err error) {
	f.limit <- struct{}{}
	defer func() { <-f.limit }()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	start := time.Now()

	defer func() {
		if resp != nil {
			// fmt.Printf() by default outputs to STDOUT, but fmt.Fprintf(os.Stdout, ...) is more explicit
			_, _ = fmt.Fprintf(os.Stdout, "Got %d from %s in %s\n", resp.StatusCode, url, time.Since(start))
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "Failed %s %s\n", method, url)
		}
	}()

	return http.DefaultClient.Do(req)
}
