package orbit

import (
	"errors"
	"strconv"

	"github.com/hackedu/orbit/router"
)

type Service struct {
	ID          int
	ProjectID   int
	ContainerID string
	Type        string

	// The port exposed from the local Docker container to the host.
	PortExposed string

	// The port on the host that PortExposed corresponds to.
	HostPort string
}

type ServicesService interface {
	Get(id int) (*Service, error)
	List(projectID int) ([]*Service, error)
	Create(service *Service) error
}

var (
	ErrServiceNotFound = errors.New("service not found")
)

type servicesService struct {
	client *Client
}

func (s *servicesService) Get(id int) (*Service, error) {
	url, err := s.client.url(router.Service, map[string]string{"ID": strconv.Itoa(id)}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var service *Service
	_, err = s.client.Do(req, &service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *servicesService) List(projectID int) ([]*Service, error) {
	url, err := s.client.url(router.Services, map[string]string{"ProjectID": strconv.Itoa(projectID)}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var services []*Service
	_, err = s.client.Do(req, &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (s *servicesService) Create(service *Service) error {
	url, err := s.client.url(router.CreateService, nil, nil)
	if err != nil {
		return err
	}

	req, err := s.client.NewRequest("POST", url.String(), service)
	if err != nil {
		return err
	}

	if _, err = s.client.Do(req, &service); err != nil {
		return err
	}

	return nil
}
