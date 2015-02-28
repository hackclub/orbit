package router

import "github.com/gorilla/mux"

func API() *mux.Router {
	m := mux.NewRouter()
	return m
}
