package handlers

import (
	"github.com/dmytro-vovk/shta/internal/storage"
)

type Handlers struct {
	reader reader
	writer writer
}

type (
	reader interface {
		GetURLs(string, string) (*storage.URLList, error)
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
