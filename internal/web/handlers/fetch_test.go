package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmytro-vovk/shta/internal/types"
	"github.com/dmytro-vovk/shta/internal/web/handlers"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	trw := &testReaderWriter{}
	h := handlers.New(trw, trw)

	t.Run("invalid sort", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/urls?sort=INVALID", http.NoBody)
		rec := httptest.NewRecorder()

		h.LatestURLs(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("invalid order", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/urls?order=INVALID", http.NoBody) // nolint:noctx
		rec := httptest.NewRecorder()

		h.LatestURLs(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	testCases := []struct {
		name        string
		url         string
		expectSort  string
		expectOrder string
	}{
		{
			name:        "no params",
			url:         "/v1/urls",
			expectSort:  "frequency",
			expectOrder: "asc",
		},
		{
			name:        "sort=frequency order=desc",
			url:         "/v1/urls?sort=frequency&order=desc",
			expectSort:  "frequency",
			expectOrder: "desc",
		},
		{
			name:        "sort=frequency order=asc",
			url:         "/v1/urls?sort=frequency&order=asc",
			expectSort:  "frequency",
			expectOrder: "asc",
		},
		{
			name:        "sort=time order=desc",
			url:         "/v1/urls?sort=time&order=desc",
			expectSort:  "time",
			expectOrder: "desc",
		},
		{
			name:        "sort=time order=asc",
			url:         "/v1/urls?sort=time&order=asc",
			expectSort:  "time",
			expectOrder: "asc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.url, http.NoBody) // nolint:noctx

			h.LatestURLs(rec, req)

			assert.Equal(t, tc.expectSort, trw.sort)
			assert.Equal(t, tc.expectOrder, trw.order)
		})
	}

	t.Run("invalid order", func(t *testing.T) {
		trw.err = errors.New("something bad happened")
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/v1/urls", http.NoBody) // nolint:noctx

		h.LatestURLs(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	})
}

type testReaderWriter struct {
	urls  []string
	sort  string
	order string
	err   error
}

func (t *testReaderWriter) GetURLs(sort, order string) (*types.URLList, error) {
	if t.err != nil {
		return nil, t.err
	}

	t.sort, t.order = sort, order

	return &types.URLList{
		URLs: []types.URL{
			{
				URL:   "http://example1.com",
				Count: 1,
			},
			{
				URL:   "http://example1.com",
				Count: 1,
			},
		},
		Sort: struct {
			By    string `json:"by"`
			Order string `json:"order"`
		}{
			By:    "",
			Order: "",
		},
	}, nil
}

func (t *testReaderWriter) AddURL(url string) {
	t.urls = append(t.urls, url)
}
