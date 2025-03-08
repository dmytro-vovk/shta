package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLValidator(t *testing.T) {
	for _, tc := range []struct {
		name  string
		url   string
		valid bool
	}{
		{
			name:  "Empty URL",
			url:   "",
			valid: false,
		},
		{
			name:  "Path only",
			url:   "/some/path",
			valid: false,
		},
		{
			name:  "Non-HTTP URL",
			url:   "ftp://example.com/some/path",
			valid: false,
		},
		{
			name:  "Empty hostname",
			url:   "http:///path",
			valid: false,
		},
		{
			name:  "HTTP",
			url:   "http://example.com",
			valid: true,
		},
		{
			name:  "HTTPS",
			url:   "HTTPS://example.com",
			valid: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.valid, validURL(tc.url))
		})
	}
}
