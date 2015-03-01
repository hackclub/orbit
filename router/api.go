package router

import "github.com/gorilla/mux"

func API() *mux.Router {
	m := mux.NewRouter()
	m.Path("/projects").Methods("POST").Name(CreateProject)
	m.Path("/projects/{ProjectID:[0-9]+}/services").Methods("GET").Name(Services)
	m.Path("/projects/{ID:.+}").Methods("GET").Name(Project)
	m.Path("/projects/{ID:.+}").Methods("PUT").Name(UpdateProject)

	m.Path("/services").Methods("POST").Name(CreateService)
	m.Path("/services/{ID:.+}").Methods("GET").Name(Service)
	return m
}
