package periodic

import (
	"fmt"
	"net/http"
	"time"
)

type Verifier struct {
	interval time.Duration
	source   LinkSourcer
	fetcher  Fetcher
}

type (
	LinkSourcer interface {
		Top() []string
	}
	Fetcher interface {
		Fetch(method, url string) (*http.Response, error)
	}
)

func New(source LinkSourcer, fetcher Fetcher, interval time.Duration) *Verifier {
	v := Verifier{
		interval: interval,
		source:   source,
		fetcher:  fetcher,
	}

	go v.processor()

	return &v
}

func (v *Verifier) processor() {
	for range time.Tick(v.interval) {
		for _, url := range v.source.Top() {
			go func(link string) {
				start := time.Now()

				resp, err := v.fetcher.Fetch(http.MethodGet, link)
				if err != nil {
					fmt.Printf("Failed to GET: %s ", link)
				} else {
					elapsed := time.Since(start).Truncate(100 * time.Millisecond)

					_ = resp.Body.Close()

					fmt.Printf("Got %d from %s in %s", resp.StatusCode, link, elapsed)
				}
			}(url)
		}
	}
}
