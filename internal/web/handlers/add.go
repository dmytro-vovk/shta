package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func (h *Handlers) AddURL(w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		log.Printf("Error reading body: %s", err)

		return
	}

	if !validURL(string(url)) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	log.Printf("Adding new URL: %s", url)

	h.writer.AddURL(string(url))

	w.WriteHeader(http.StatusAccepted)
}

func validURL(link string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}

	return (u.Scheme == "http" || u.Scheme == "https") && u.Hostname() != ""
}
