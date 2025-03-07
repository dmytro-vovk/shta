package storage

import (
	"fmt"
	"log"
	"slices"
	"sort"

	"github.com/dmytro-vovk/shta/internal/db"
	"github.com/dmytro-vovk/shta/internal/types"
)

type Storage struct {
	db     *db.DB
	counts Counter
}

type Counter interface {
	Add(string)
	Get(string) int
}

func New(db *db.DB, counter Counter) *Storage {
	return &Storage{
		db:     db,
		counts: counter,
	}
}

func (s *Storage) AddURL(url string) {
	if err := s.db.UpsertURL(url); err != nil {
		log.Printf("Error upsertting URL %q: %s", url, err)
	} else {
		s.counts.Add(url)
	}
}

func (s *Storage) GetURLs(sortBy, sortOrder string) (*types.URLList, error) {
	log.Printf("Getting URLs sorted by '%s' in %s order", sortBy, sortOrder)

	urls, err := s.db.FetchURLs(50)
	if err != nil {
		return nil, fmt.Errorf("fetch urls: %w", err)
	}

	sorted := make([]types.URL, len(urls))

	for i := range urls {
		sorted[i] = types.URL{
			URL:   urls[i].URL,
			Count: s.counts.Get(urls[i].URL),
			Seen:  urls[i].LastSeen,
		}
	}

	var (
		timeSort = func(i, j int) bool {
			return sorted[i].Seen.After(sorted[j].Seen)
		}
		frequencySort = func(i, j int) bool {
			return sorted[i].Count > sorted[j].Count
		}
	)

	switch {
	case sortBy == "time":
		sort.Slice(sorted, timeSort)
	case sortBy == "frequency":
		sort.Slice(sorted, frequencySort)
	}

	if sortOrder == "asc" {
		slices.Reverse(sorted)
	}

	result := types.URLList{
		URLs: sorted,
	}

	result.Sort.By = sortBy
	result.Sort.Order = sortOrder

	return &result, nil
}
