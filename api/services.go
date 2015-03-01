package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zachlatta/orbit"
)

func serveService(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		return err
	}

	service, err := store.Services.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, service)
}

func serveCreateService(w http.ResponseWriter, r *http.Request) error {
	var service orbit.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		return err
	}

	if err := store.Services.Create(&service); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return writeJSON(w, service)
}
