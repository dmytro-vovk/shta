package fetcher_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmytro-vovk/shta/internal/fetcher"
	"github.com/stretchr/testify/assert"
)

func TestFetcher(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(ts.Close)

	testCases := []struct {
		name               string
		url                string
		method             string
		expectResponseCode int
		expectErr          error
	}{
		{
			name:               "Valid URL",
			url:                ts.URL,
			method:             http.MethodGet,
			expectErr:          nil,
			expectResponseCode: http.StatusOK,
		},
		{
			name:               "Valid URL",
			url:                ts.URL + "/not-found",
			method:             http.MethodGet,
			expectErr:          nil,
			expectResponseCode: http.StatusNotFound,
		},
		{
			name:      "Timeout",
			url:       ts.URL + "/timeout",
			method:    http.MethodGet,
			expectErr: fmt.Errorf("Get %q: context deadline exceeded", ts.URL+"/timeout"),
		},
		{
			name:      "Bad method",
			url:       ts.URL + "/timeout",
			method:    " ",
			expectErr: errors.New(`create request: net/http: invalid method " "`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := fetcher.New(2)

			resp, err := f.Fetch(tc.method, tc.url)

			if resp != nil {
				_ = resp.Body.Close()
			}

			if tc.expectErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResponseCode, resp.StatusCode)
			} else {
				assert.EqualError(t, err, tc.expectErr.Error())
			}
		})
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/":
		w.WriteHeader(http.StatusOK)
	case "/not-found":
		w.WriteHeader(http.StatusNotFound)
	case "/timeout":
		time.Sleep(4 * time.Second)
		w.WriteHeader(http.StatusTeapot)
	}
}
