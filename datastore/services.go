package datastore

import (
	"database/sql"

	"github.com/hackedu/orbit"
)

func init() {
	DB.AddTableWithName(orbit.Service{}, "service").SetKeys(true, "ID")
	createSQL = append(createSQL,
		`ALTER TABLE service
		ADD CONSTRAINT service_productid_fkey
		FOREIGN KEY (projectid)
		REFERENCES project ON DELETE CASCADE;`,
		`CREATE UNIQUE INDEX service_containerid ON service(containerid);`)
}

type servicesStore struct{ *Datastore }

func (s *servicesStore) Get(id int) (*orbit.Service, error) {
	var service orbit.Service
	if err := s.dbh.SelectOne(&service, `SELECT * FROM service WHERE id=$1`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, orbit.ErrServiceNotFound
		}
		return nil, err
	}
	return &service, nil
}

func (s *servicesStore) List(projectID int) ([]*orbit.Service, error) {
	var services []*orbit.Service
	if err := s.dbh.Select(&services, `SELECT * FROM service WHERE projectid=$1`, projectID); err != nil {
		if err == sql.ErrNoRows {
			return nil, orbit.ErrServiceNotFound
		}
		return nil, err
	}
	return services, nil
}

func (s *servicesStore) Create(service *orbit.Service) error {
	if err := s.dbh.Insert(service); err != nil {
		return err
	}
	return nil
}

func (s *servicesStore) Update(service *orbit.Service) error {
	if _, err := s.dbh.Update(service); err != nil {
		return err
	}
	return nil
}
