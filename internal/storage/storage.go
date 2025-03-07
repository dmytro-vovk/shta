package storage

import (
	"fmt"
	"log"

	"github.com/dmytro-vovk/shta/internal/db"
)

type Storage struct {
	db *db.DB
}

func New(db *db.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) AddURL(url string) {
	if err := s.db.UpsertURL(url); err != nil {
		log.Printf("Error upsertting URL %q: %s", url, err)
	}
}

type URLList struct {
	URLs []db.URLRecord `json:"urls"`
	Sort struct {
		Field     string `json:"field"`
		Direction string `json:"direction"`
	} `json:"sort"`
}

func (s *Storage) GetURLs(sortBy, sortDir string) (*URLList, error) {
	var (
		sort  string
		order string
	)

	switch sortBy {
	case "frequency":
		sort = "seen"
	case "time":
		sort = "created_at"
	default:
		return nil, fmt.Errorf("unknown sort by: %q", sortBy)
	}

	switch sortDir {
	case "asc":
		order = "ASC"
	case "desc":
		order = "DESC"
	default:
		return nil, fmt.Errorf("unknown sort direction: %q", sortDir)
	}

	log.Printf("Getting URLs sorted by '%s' in %s order", sort, order)

	urls, err := s.db.FetchURLs(sort, order, 50)
	if err != nil {
		return nil, err
	}

	if urls == nil {
		urls = []db.URLRecord{}
	}

	result := URLList{
		URLs: urls,
	}

	result.Sort.Field = sortBy
	result.Sort.Direction = sortDir

	return &result, nil
}
