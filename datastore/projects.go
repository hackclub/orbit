package datastore

import (
	"github.com/zachlatta/orbit"
)

func init() {
	DB.AddTableWithName(orbit.Project{}, "project").SetKeys(true, "ID")
}

type projectsStore struct{ *Datastore }

func (s *projectsStore) Get(id int) (*orbit.Project, error) {
	var project *orbit.Project
	if err := s.dbh.SelectOne(&project, `SELECT * FROM project WHERE id=$1`, id); err != nil {
		return nil, err
	}
	if project == nil {
		return nil, orbit.ErrProjectNotFound
	}
	return project, nil
}

func (s *projectsStore) Create(project *orbit.Project) error {
	if err := s.dbh.Insert(project); err != nil {
		return err
	}
	return nil
}
