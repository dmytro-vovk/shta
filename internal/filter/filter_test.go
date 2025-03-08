package filter_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dmytro-vovk/shta/internal/filter"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	testCases := []struct {
		name      string
		url       string
		response  *http.Response
		err       error
		expectURL bool
		expectLog string
	}{
		{
			name: "Happy path",
			url:  "http://example.com/",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			err:       nil,
			expectURL: true,
			expectLog: "URL http://example.com/ is alive",
		},
		{
			name:      "Request error",
			url:       "http://example.com/unreachable",
			response:  nil,
			err:       io.EOF,
			expectURL: false,
			expectLog: "HEAD request to http://example.com/unreachable failed: EOF",
		},
		{
			name: "Bad response",
			url:  "http://example.com/not-found",
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			err:       nil,
			expectURL: false,
			expectLog: "URL http://example.com/not-found returned code 404, skipping",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ta := &testAdder{}
			tf := &testFetcher{
				returnResponse: tc.response,
				err:            tc.err,
			}
			f := filter.New(ta, tf)

			var buf bytes.Buffer

			log.SetOutput(&buf)

			f.AddURL(tc.url)

			time.Sleep(10 * time.Millisecond) // Should wait a bit, adding URL is asynchronous

			log.SetOutput(os.Stderr)

			assert.Contains(t, buf.String(), tc.expectLog)

			if tc.expectURL && assert.Len(t, ta.getURLs(), 1) {
				assert.Equal(t, ta.getURLs()[0], tc.url)
			}
		})
	}
}

type testAdder struct {
	m    sync.Mutex
	urls []string
}

func (ta *testAdder) AddURL(url string) {
	ta.m.Lock()
	ta.urls = append(ta.urls, url)
	ta.m.Unlock()
}

func (ta *testAdder) getURLs() []string {
	ta.m.Lock()
	defer ta.m.Unlock()

	return ta.urls
}

type testFetcher struct {
	returnResponse *http.Response
	err            error
}

func (t testFetcher) Fetch(string, string) (*http.Response, error) {
	return t.returnResponse, t.err
}
