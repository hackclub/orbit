package datastore

import (
	"database/sql"

	"github.com/hackedu/orbit"
)

func init() {
	DB.AddTableWithName(orbit.Project{}, "project").SetKeys(true, "ID")
}

type projectsStore struct{ *Datastore }

func (s *projectsStore) Get(id int) (*orbit.Project, error) {
	var project orbit.Project
	if err := s.dbh.SelectOne(&project, `SELECT * FROM project WHERE id=$1`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, orbit.ErrProjectNotFound
		}
		return nil, err
	}
	return &project, nil
}

func (s *projectsStore) Create(project *orbit.Project) error {
	if err := s.dbh.Insert(project); err != nil {
		return err
	}
	return nil
}

func (s *projectsStore) Update(project *orbit.Project) error {
	if _, err := s.dbh.Update(project); err != nil {
		return err
	}
	return nil
}
