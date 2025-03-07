package handlers

import (
	"io"
	"log"
	"net/http"
)

func (h *Handlers) AddURL(w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		log.Printf("Error reading body: %s", err)

		return
	}

	log.Printf("Adding new URL: %s", url)

	h.writer.AddURL(string(url))

	w.WriteHeader(http.StatusAccepted)
}
