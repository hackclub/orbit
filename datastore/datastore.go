package datastore

import (
	"github.com/jmoiron/modl"
	"github.com/hackedu/orbit"
)

type Datastore struct {
	Projects orbit.ProjectsService
	Services orbit.ServicesService

	dbh modl.SqlExecutor
}

func NewDatastore(dbh modl.SqlExecutor) *Datastore {
	if dbh == nil {
		dbh = DBH
	}

	d := &Datastore{dbh: dbh}
	d.Projects = &projectsStore{d}
	d.Services = &servicesStore{d}
	return d
}
