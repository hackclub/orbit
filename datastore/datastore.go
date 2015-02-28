package datastore

import (
	"github.com/jmoiron/modl"
	"github.com/zachlatta/orbit"
)

type Datastore struct {
	Projects orbit.ProjectsService

	dbh modl.SqlExecutor
}

func NewDatastore(dbh modl.SqlExecutor) *Datastore {
	if dbh == nil {
		dbh = DBH
	}

	d := &Datastore{dbh: dbh}
	d.Projects = &projectsStore{d}
	return d
}
