package storage_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dmytro-vovk/shta/internal/storage"
	"github.com/dmytro-vovk/shta/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	td := &testDependency{}
	s := storage.New(td, td)

	testCases := []struct {
		name       string
		sortBy     string
		sortOrder  string
		expectErr  error
		expectURLs *types.URLList
	}{
		{
			name:      "time/asc",
			sortBy:    types.SortByTime,
			sortOrder: types.OrderAsc,
			expectErr: nil,
			expectURLs: &types.URLList{
				URLs: []types.URL{
					{
						URL:   "http://example1.com",
						Count: 1,
						Seen:  time.Date(2025, 3, 9, 21, 47, 0, 0, time.UTC),
					},
					{
						URL:   "http://example2.com",
						Count: 2,
						Seen:  time.Date(2025, 3, 9, 20, 47, 0, 0, time.UTC),
					},
					{
						URL:   "http://example3.com",
						Count: 3,
						Seen:  time.Date(2025, 3, 9, 19, 47, 0, 0, time.UTC),
					},
				},
				Sort: struct {
					By    string `json:"by"`
					Order string `json:"order"`
				}{
					By:    types.SortByTime,
					Order: types.OrderAsc,
				},
			},
		},
		{
			name:      "time/desc",
			sortBy:    types.SortByTime,
			sortOrder: types.OrderDesc,
			expectErr: nil,
			expectURLs: &types.URLList{
				URLs: []types.URL{
					{
						URL:   "http://example3.com",
						Count: 3,
						Seen:  time.Date(2025, 3, 9, 19, 47, 0, 0, time.UTC),
					},
					{
						URL:   "http://example2.com",
						Count: 2,
						Seen:  time.Date(2025, 3, 9, 20, 47, 0, 0, time.UTC),
					},
					{
						URL:   "http://example1.com",
						Count: 1,
						Seen:  time.Date(2025, 3, 9, 21, 47, 0, 0, time.UTC),
					},
				},
				Sort: struct {
					By    string `json:"by"`
					Order string `json:"order"`
				}{
					By:    types.SortByTime,
					Order: types.OrderDesc,
				},
			},
		},
		{
			name:      "frequency/asc",
			sortBy:    types.SortByFrequency,
			sortOrder: types.OrderAsc,
			expectErr: nil,
			expectURLs: &types.URLList{
				URLs: []types.URL{
					{
						URL:   "http://example1.com",
						Count: 1,
						Seen:  time.Date(2025, 3, 9, 21, 47, 0, 0, time.UTC),
					},
					{
						URL:   "http://example2.com",
						Count: 2,
						Seen:  time.Date(2025, 3, 9, 20, 47, 0, 0, time.UTC),
					},
					{
						URL:   "http://example3.com",
						Count: 3,
						Seen:  time.Date(2025, 3, 9, 19, 47, 0, 0, time.UTC),
					},
				},
				Sort: struct {
					By    string `json:"by"`
					Order string `json:"order"`
				}{
					By:    types.SortByFrequency,
					Order: types.OrderAsc,
				},
			},
		},
		{
			name:      "frequency/desc",
			sortBy:    types.SortByFrequency,
			sortOrder: types.OrderDesc,
			expectErr: nil,
			expectURLs: &types.URLList{
				URLs: []types.URL{
					{
						URL:   "http://example3.com",
						Count: 3,
						Seen:  time.Date(2025, 3, 9, 19, 47, 0, 0, time.UTC),
					},
					{
						URL:   "http://example2.com",
						Count: 2,
						Seen:  time.Date(2025, 3, 9, 20, 47, 0, 0, time.UTC),
					},
					{
						URL:   "http://example1.com",
						Count: 1,
						Seen:  time.Date(2025, 3, 9, 21, 47, 0, 0, time.UTC),
					},
				},
				Sort: struct {
					By    string `json:"by"`
					Order string `json:"order"`
				}{
					By:    types.SortByFrequency,
					Order: types.OrderDesc,
				},
			},
		},
		{
			name:       "error",
			sortBy:     types.SortByFrequency,
			sortOrder:  types.OrderDesc,
			expectErr:  errors.New("some error"),
			expectURLs: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			td.err = tc.expectErr

			urls, err := s.GetURLs(tc.sortBy, tc.sortOrder)
			if tc.expectErr != nil {
				require.EqualError(t, err, "fetch urls: "+tc.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectURLs, urls)
			}
		})
	}

	t.Run("upsert", func(t *testing.T) {
		td.err = nil

		assert.Zero(t, td.added)

		s.AddURL("http://example1.com")

		assert.Equal(t, 1, td.added)

		td.err = errors.New("some error")

		s.AddURL("http://example2.com")

		assert.Equal(t, 1, td.added)
	})
}

type testDependency struct {
	err   error
	added int
}

func (t *testDependency) Add(string) {
	t.added++
}

func (t *testDependency) Get(url string) int {
	return map[string]int{
		"http://example1.com": 1,
		"http://example2.com": 2,
		"http://example3.com": 3,
	}[url]
}

func (t *testDependency) UpsertURL(string) error {
	return t.err
}

func (t *testDependency) FetchURLs(int) ([]*types.URLRecord, error) {
	if t.err != nil {
		return nil, t.err
	}

	return []*types.URLRecord{
		{
			URL:      "http://example1.com",
			LastSeen: time.Date(2025, 3, 9, 21, 47, 0, 0, time.UTC),
		},
		{
			URL:      "http://example2.com",
			LastSeen: time.Date(2025, 3, 9, 20, 47, 0, 0, time.UTC),
		},
		{
			URL:      "http://example3.com",
			LastSeen: time.Date(2025, 3, 9, 19, 47, 0, 0, time.UTC),
		},
	}, nil
}
