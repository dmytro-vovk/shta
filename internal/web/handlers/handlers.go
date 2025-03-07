package handlers

import (
	"github.com/dmytro-vovk/shta/internal/types"
)

type Handlers struct {
	reader reader
	writer writer
}

type (
	reader interface {
		GetURLs(string, string) (*types.URLList, error)
	}

	writer interface {
		AddURL(string)
	}
)

func New(reader reader, writer writer) *Handlers {
	return &Handlers{
		reader: reader,
		writer: writer,
	}
}
