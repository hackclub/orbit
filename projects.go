package orbit

import "errors"

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
	return nil, nil
}

func (s *projectsService) Create(project *Project) error {
	return nil
}

func (s *projectsService) Update(project *Project) error {
	return nil
}
