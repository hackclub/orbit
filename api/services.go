package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hackedu/orbit"
	"github.com/hackedu/orbit/docker"
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

func serveServices(w http.ResponseWriter, r *http.Request) error {
	projectID, err := strconv.Atoi(mux.Vars(r)["ProjectID"])
	if err != nil {
		return err
	}

	service, err := store.Services.List(projectID)
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

	var err error
	service.ContainerID, service.HostPort, err =
		docker.RunContainer(service.ProjectID, service.Type, service.PortExposed)
	if err != nil {
		return err
	}

	if err := store.Services.Create(&service); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return writeJSON(w, service)
}
