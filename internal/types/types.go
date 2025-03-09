package types

import "time"

// URL contain all the information about a link
type URL struct {
	URL   string    `json:"url"`   // The URL itself
	Count int       `json:"count"` // Number of submissions
	Seen  time.Time `json:"-"`     // Last submission timestamp
}

// URLRecord is database URL entity
type URLRecord struct {
	URL      string    `db:"url"`
	LastSeen time.Time `db:"last_seen"`
}

// URLList is a list of recent URLs with sorting parameters
type URLList struct {
	URLs []URL `json:"urls"` // List of URLs
	Sort struct {
		By    string `json:"by"`    // How the list is sorted
		Order string `json:"order"` // Sort order
	} `json:"sort"`
}

const (
	SortByFrequency = "frequency"
	SortByTime      = "time"
	OrderAsc        = "asc"
	OrderDesc       = "desc"
)
