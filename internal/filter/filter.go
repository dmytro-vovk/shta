package filter

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Filter struct {
	next     URLAdder
	incoming chan string
	limit    chan struct{}
}

type URLAdder interface {
	AddURL(string)
}

func New(next URLAdder) *Filter {
	f := Filter{
		incoming: make(chan string, 1000),
		limit:    make(chan struct{}, 3),
		next:     next,
	}

	go f.fetcher()

	return &f
}

func (f *Filter) AddURL(url string) { f.incoming <- url }

func (f *Filter) fetcher() {
	for url := range f.incoming {
		go f.withLimiter(func() {
			f.fetch(url)
		})
	}
}

func (f *Filter) withLimiter(fn func()) {
	f.limit <- struct{}{}

	fn()

	<-f.limit
}

func (f *Filter) fetch(url string) {
	r, err := doRequest(http.MethodHead, url)
	if err != nil {
		log.Printf("HEAD request to %s failed: %s", url, err)

		r, err = doRequest(http.MethodGet, url)
	}

	if err != nil {
		log.Printf("Error fetching URL %s: %s", url, err)

		return
	}

	_ = r.Body.Close()

	if r.StatusCode/100 == 2 {
		log.Printf("URL %s is alive, adding", url)

		f.next.AddURL(url)
	} else {
		log.Printf("URL %s returned code %d, skipping", url, r.StatusCode)
	}

	f.next.AddURL(url)
}

func doRequest(method, url string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, method, url, nil)

	return http.DefaultClient.Do(req)
}
