package main

import (
	"github.com/gorilla/mux"

	tools "./tools"
)

type Server struct {
	Config *tools.EtcdConfig
	Mux    *mux.Router
}
