package storage

import (
	"fmt"
	"log"
	"sort"

	"github.com/dmytro-vovk/shta/internal/types"
)

type Storage struct {
	db     database
	counts Counter
}

type (
	database interface {
		UpsertURL(string) error
		FetchURLs(int) ([]*types.URLRecord, error)
	}
)

type Counter interface {
	Add(string)
	Get(string) int
}

func New(db database, counter Counter) *Storage {
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

	urls, err := s.db.FetchURLs(50) // NOTE would be nice to have '50' configurable
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

	switch sortBy {
	case types.SortByTime:
		if sortOrder == types.OrderAsc {
			sort.Slice(sorted, func(i, j int) bool { return sorted[i].Seen.After(sorted[j].Seen) })
		} else {
			sort.Slice(sorted, func(i, j int) bool { return sorted[i].Seen.Before(sorted[j].Seen) })
		}
	case types.SortByFrequency:
		if sortOrder == types.OrderAsc {
			sort.Slice(sorted, func(i, j int) bool { return sorted[i].Count < sorted[j].Count })
		} else {
			sort.Slice(sorted, func(i, j int) bool { return sorted[i].Count > sorted[j].Count })
		}
	}

	result := types.URLList{
		URLs: sorted,
	}

	result.Sort.By = sortBy
	result.Sort.Order = sortOrder

	return &result, nil
}
