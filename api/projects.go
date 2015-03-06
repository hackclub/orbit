package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hackedu/orbit"
	"github.com/hackedu/orbit/docker"
	"github.com/hackedu/orbit/git"
)

func serveProject(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		return err
	}

	project, err := store.Projects.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, project)
}

func serveCreateProject(w http.ResponseWriter, r *http.Request) error {
	var project orbit.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		return err
	}

	// TODO create logic

	if err := store.Projects.Create(&project); err != nil {
		return err
	}

	projectPath, err := git.InitializeProject(&project)
	if err != nil {
		return err
	}
	project.GitPath = projectPath

	if err := store.Projects.Update(&project); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return writeJSON(w, project)
}

func serveUpdateProject(w http.ResponseWriter, r *http.Request) error {
	var project orbit.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		return err
	}

	if err := store.Projects.Update(&project); err != nil {
		return err
	}
	return writeJSON(w, project)
}

func serveProjectCommand(w http.ResponseWriter, r *http.Request) error {
	var projectCmd orbit.ProjectCmd
	if err := json.NewDecoder(r.Body).Decode(&projectCmd); err != nil {
		return err
	}

	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		return err
	}

	services, err := store.Services.List(id)
	if err != nil {
		return err
	}

	var service *orbit.Service
	for _, s := range services {
		if s.Type == projectCmd.ContainerType {
			service = s
			break
		}
	}
	if service == nil {
		return fmt.Errorf("service of type %s not found\n", projectCmd.ContainerType)
	}

	cmd := docker.CommandInContainer(service.ContainerID, projectCmd.Command...)
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := docker.CommandInContainer(service.ContainerID,
		"git", "config", "--global", "user.email", "orbit@hackedu.us",
	).Run(); err != nil {
		return fmt.Errorf("error committing files: %q", err.Error())
	}

	if err := docker.CommandInContainer(service.ContainerID,
		"git", "config", "--global", "user.name", "Orbit",
	).Run(); err != nil {
		return fmt.Errorf("error committing files: %q", err.Error())
	}

	if err := docker.CommandInContainer(service.ContainerID,
		"git", "add", "-A", ":/",
	).Run(); err != nil {
		return fmt.Errorf("error adding files: %q", err.Error())
	}

	if err := docker.CommandInContainer(service.ContainerID,
		"git", "commit", "-m", "", "--allow-empty-message", "--allow-empty",
	).Run(); err != nil {
		return fmt.Errorf("error committing files: %q", err.Error())
	}

	if err := docker.CommandInContainer(service.ContainerID,
		"git", "push",
	).Run(); err != nil {
		return fmt.Errorf("error pushing files: %q", err.Error())
	}

	return nil
}
