package datastore

import (
	"github.com/jmoiron/modl"
)

type Datastore struct {
	dbh modl.SqlExecutor
}

func NewDatastore(dbh modl.SqlExecutor) *Datastore {
	if dbh == nil {
		dbh = DBH
	}

	d := &Datastore{dbh: dbh}
	return d
}
