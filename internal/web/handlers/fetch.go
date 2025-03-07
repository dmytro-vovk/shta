package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dmytro-vovk/shta/internal/types"
)

func (h *Handlers) LatestURLs(w http.ResponseWriter, r *http.Request) {
	var (
		sortBy  string
		sortDir string
	)

	switch sort := r.URL.Query().Get("sort"); sort {
	case "frequency,asc":
		sortBy, sortDir = types.SortByFrequency, types.OrderAsc
	case "frequency,desc":
		sortBy, sortDir = types.SortByFrequency, types.OrderDesc
	case "time,asc":
		sortBy, sortDir = types.SortByTime, types.OrderAsc
	case "time,desc":
		sortBy, sortDir = types.SortByTime, types.OrderDesc
	case "":
		sortBy, sortDir = types.SortByFrequency, types.OrderAsc
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
