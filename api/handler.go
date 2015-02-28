package api

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"sourcegraph.com/sourcegraph/thesrc/datastore"
	"sourcegraph.com/sourcegraph/thesrc/router"
)

var (
	store         = datastore.NewDatastore(nil)
	schemaDecoder = schema.NewDecoder()
)

func Handler() *mux.Router {
	m := router.API()
	return m
}
