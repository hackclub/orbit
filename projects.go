package orbit

import (
	"errors"
	"strconv"

	"github.com/hackedu/orbit/router"
)

type ProjectCmd struct {
	ContainerType string
	Command       []string
}

type Project struct {
	ID int

	// GitPath is a relative path to the project's files from the git root.
	GitPath string
}

type ProjectsService interface {
	Get(id int) (*Project, error)
	Create(project *Project) error
	Update(project *Project) error
}

var (
	ErrProjectNotFound = errors.New("project not found")
)

type projectsService struct {
	client *Client
}

func (s *projectsService) Get(id int) (*Project, error) {
	url, err := s.client.url(router.Project, map[string]string{"ID": strconv.Itoa(id)}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var project *Project
	_, err = s.client.Do(req, &project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectsService) Create(project *Project) error {
	url, err := s.client.url(router.CreateProject, nil, nil)
	if err != nil {
		return err
	}

	req, err := s.client.NewRequest("POST", url.String(), project)
	if err != nil {
		return err
	}

	if _, err = s.client.Do(req, &project); err != nil {
		return err
	}

	return nil
}

func (s *projectsService) Update(project *Project) error {
	url, err := s.client.url(router.UpdateProject, nil, nil)
	if err != nil {
		return err
	}

	req, err := s.client.NewRequest("PUT", url.String(), project)
	if err != nil {
		return err
	}

	if _, err := s.client.Do(req, &project); err != nil {
		return err
	}

	return nil
}
