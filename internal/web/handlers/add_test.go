package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmytro-vovk/shta/internal/web/handlers"
	"github.com/stretchr/testify/assert"
)

func TestAddURL(t *testing.T) {
	trw := &testReaderWriter{}
	h := handlers.New(trw, trw)

	t.Run("adding urls", func(t *testing.T) {
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/v1/urls", strings.NewReader("https://example.com")) // nolint:noctx

		h.AddURL(rec, req)

		assert.Equal(t, http.StatusAccepted, rec.Code)

		if assert.Len(t, trw.urls, 1) {
			assert.Equal(t, trw.urls, []string{"https://example.com"})
		}

		req, _ = http.NewRequest(http.MethodPost, "/v1/urls", strings.NewReader("https://another.com")) // nolint:noctx

		h.AddURL(rec, req)

		if assert.Len(t, trw.urls, 2) {
			assert.Equal(t, trw.urls, []string{"https://example.com", "https://another.com"})
		}
	})

	t.Run("error reading body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/v1/urls", &errorReader{}) // nolint:noctx

		h.AddURL(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("test error")
}

func TestURLValidator(t *testing.T) {
	trw := &testReaderWriter{}
	h := handlers.New(trw, trw)

	for _, tc := range []struct {
		name        string
		url         string
		expectValid bool
	}{
		{
			name:        "Empty URL",
			url:         "",
			expectValid: false,
		},
		{
			name:        "Path only",
			url:         "/some/path",
			expectValid: false,
		},
		{
			name:        "Non-HTTP URL",
			url:         "ftp://example.com/some/path",
			expectValid: false,
		},
		{
			name:        "Empty hostname",
			url:         "http:///path",
			expectValid: false,
		},
		{
			name:        "HTTP",
			url:         "http://example.com",
			expectValid: true,
		},
		{
			name:        "HTTPS",
			url:         "HTTPS://example.com",
			expectValid: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader(tc.url)) // nolint:noctx

			h.AddURL(rec, req)

			if tc.expectValid {
				assert.Equal(t, http.StatusAccepted, rec.Result().StatusCode)
			} else {
				assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
			}
		})
	}
}
