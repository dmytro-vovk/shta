package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dmytro-vovk/shta/internal/types"
)

func (h *Handlers) LatestURLs(w http.ResponseWriter, r *http.Request) {
	var sortBy string

	switch sort := r.URL.Query().Get("sort"); sort {
	case "frequency":
		sortBy = types.SortByFrequency
	case "time":
		sortBy = types.SortByTime
	case "":
		sortBy = types.SortByFrequency
	default:
		http.Error(w, "Invalid sort parameter: "+sort, http.StatusBadRequest)
		log.Printf("Error: invalid sort parameter: %s", sort)

		return
	}

	var sortOrder string

	switch order := r.URL.Query().Get("order"); order {
	case "asc":
		sortOrder = types.OrderAsc
	case "desc":
		sortOrder = types.OrderDesc
	case "":
		sortOrder = types.OrderAsc
	default:
		http.Error(w, "Invalid order parameter: "+order, http.StatusBadRequest)
		log.Printf("Error: invalid order parameter: %s", order)

		return
	}

	data, err := h.reader.GetURLs(sortBy, sortOrder)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Error getting URLs: %s", err)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil { // coverage-ignore
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		log.Printf("Error encoding URLs: %s", err)
	}
}
