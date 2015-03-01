package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/zachlatta/orbit/datastore"
	"github.com/zachlatta/orbit/router"
)

var (
	store         = datastore.NewDatastore(nil)
	schemaDecoder = schema.NewDecoder()
)

func Handler() *mux.Router {
	m := router.API()
	m.Get(router.Project).Handler(handler(serveProject))
	m.Get(router.CreateProject).Handler(handler(serveCreateProject))
	m.Get(router.UpdateProject).Handler(handler(serveUpdateProject))

	m.Get(router.Service).Handler(handler(serveService))
	m.Get(router.CreateService).Handler(handler(serveCreateService))
	return m
}

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		log.Println(err)
	}
}
