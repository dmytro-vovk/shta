package filter

import (
	"log"
	"net/http"
)

type Filter struct {
	next     URLAdder
	fetcher  Fetcher
	incoming chan string
}

type (
	URLAdder interface {
		AddURL(string)
	}
	Fetcher interface {
		Fetch(method, url string) (*http.Response, error)
	}
)

func New(next URLAdder, fetcher Fetcher) *Filter {
	f := Filter{
		incoming: make(chan string, 1000),
		fetcher:  fetcher,
		next:     next,
	}

	go f.processor()

	return &f
}

// AddURL puts new URL into processing queue
func (f *Filter) AddURL(url string) { f.incoming <- url }

func (f *Filter) processor() {
	for url := range f.incoming {
		f.fetch(url)
	}
}

func (f *Filter) fetch(url string) {
	r, err := f.fetcher.Fetch(http.MethodHead, url)
	if err != nil {
		log.Printf("HEAD request to %s failed: %s", url, err)

		r, err = f.fetcher.Fetch(http.MethodGet, url)
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
}
