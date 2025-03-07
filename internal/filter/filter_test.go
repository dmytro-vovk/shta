package filter_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dmytro-vovk/shta/internal/filter"
)

func TestFilter(t *testing.T) {
	ta := &testAdder{t: t}
	tf := &testFetcher{}
	f := filter.New(ta, tf)

	tf.returnResponse = &http.Response{
		Status:           "",
		StatusCode:       http.StatusOK,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             io.NopCloser(strings.NewReader("")),
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}

	f.AddURL("http://example.com/path")
	// TODO complete this
	time.Sleep(10 * time.Millisecond)
}

type testAdder struct {
	t *testing.T
}

func (ta testAdder) AddURL(url string) {
	ta.t.Logf("URL added: %s", url)
}

type testFetcher struct {
	returnResponse *http.Response
	err            error
}

func (t testFetcher) Fetch(string, string) (*http.Response, error) {
	return t.returnResponse, t.err
}
