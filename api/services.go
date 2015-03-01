package api

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

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

	var err error
	service.ContainerID, service.HostPort, err = runContainer(service.Type, service.PortExposed)
	if err != nil {
		return err
	}

	if err := store.Services.Create(&service); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return writeJSON(w, service)
}

func runContainer(image, port string) (containerID, hostPort string, err error) {
	output, err := exec.Command("docker", "run",
		"-d",
		"-p", port,
		image,
		"/bin/sh", "-c", "while true; do sleep 1; done",
	).Output()
	if err != nil {
		return "", "", err
	}
	containerID = strings.TrimSpace(string(output))

	output, err = exec.Command("docker", "inspect", containerID).Output()
	if err != nil {
		return "", "", err
	}

	var inspectOutput []struct {
		NetworkSettings struct {
			Ports map[string][]struct {
				HostPort string
			}
		}
	}
	if err := json.Unmarshal(output, &inspectOutput); err != nil {
		return "", "", err
	}

	hostPort = inspectOutput[0].NetworkSettings.Ports[port+"/tcp"][0].HostPort

	return containerID, hostPort, nil
}
