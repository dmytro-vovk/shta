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
			go v.verifyURL(url)
		}
	}
}

func (v *Verifier) verifyURL(url string) {
	resp, err := v.fetcher.Fetch(http.MethodGet, url)
	if err != nil {
		fmt.Printf("Failed to GET: %s ", url)
	} else {
		_ = resp.Body.Close()
	}
}
