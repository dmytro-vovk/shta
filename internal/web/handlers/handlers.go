package handlers

import (
	"github.com/dmytro-vovk/shta/internal/types"
)

type Handlers struct {
	reader urlReader
	writer urlWriter
}

type (
	urlReader interface {
		GetURLs(string, string) (*types.URLList, error)
	}

	urlWriter interface {
		AddURL(string)
	}
)

func New(reader urlReader, writer urlWriter) *Handlers {
	return &Handlers{
		reader: reader,
		writer: writer,
	}
}
