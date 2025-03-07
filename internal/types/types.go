package types

import "time"

// URL contain all the information about a link
type URL struct {
	URL   string    `json:"url"`   // The URL itself
	Count int       `json:"count"` // Number of submissions
	Seen  time.Time `json:"-"`     // Last submission timestamp
}

// URLList is a list of recent URLs with sorting parameters
type URLList struct {
	URLs []URL `json:"urls"` // List of URLs
	Sort struct {
		By    string `json:"by"`    // How the list is sorted
		Order string `json:"order"` // Sort order
	} `json:"sort"`
}
