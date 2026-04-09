package registry

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	var reg Registry
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if reg.Protocol == "" {
		reg.Protocol = "http"
	}
	if reg.Status == "" {
		reg.Status = "active"
	}
	if reg.PolicyIDs == nil {
		reg.PolicyIDs = []string{}
	}
	if err := insert(r.Context(), &reg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reg)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var reg Registry
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := update(r.Context(), id, &reg); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reg)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := deleteByID(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	reg, err := getByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reg)
}

func GetByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		var payload struct {
			Name string `json:"name"`
		}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil && !errors.Is(err, io.EOF) {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		name = payload.Name
	}
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	reg, err := getByName(r.Context(), name)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reg)
}

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	regs, err := getAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if regs == nil {
		regs = []*Registry{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(regs)
}
