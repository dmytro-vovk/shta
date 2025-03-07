package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handlers) LatestURLs(w http.ResponseWriter, r *http.Request) {
	var (
		sortBy  string
		sortDir string
	)

	sort := r.URL.Query().Get("sort")

	const (
		frequency = "frequency"
		time      = "time"
		asc       = "asc"
		desc      = "desc"
	)

	switch sort {
	case "frequency,asc":
		sortBy, sortDir = frequency, asc
	case "frequency,desc":
		sortBy, sortDir = frequency, desc
	case "time,asc":
		sortBy, sortDir = time, asc
	case "time,desc":
		sortBy, sortDir = time, desc
	case "":
		sortBy, sortDir = frequency, asc
	default:
		http.Error(w, "Invalid sort parameter: "+sort, http.StatusBadRequest)

		log.Printf("Error: invalid sort parameter: %s", sort)

		return
	}

	data, err := h.reader.GetURLs(sortBy, sortDir)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		log.Printf("Error getting URLs: %s", err)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		log.Printf("Error encoding URLs: %s", err)
	}
}
