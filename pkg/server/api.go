package server

import (
	"os/exec"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

type api struct {
	router http.Handler
}

type Server interface {
	Router() http.Handler
}

func (a *api) fetchGophers(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("/bin/grpc_health_probe", "-addr=:5000")
	stdout, err := cmd.Output()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(stdout)
        return
	}
	json.NewEncoder(w).Encode("Healthy")
}

func New() Server {
	a := &api{}

	r := mux.NewRouter()
	r.HandleFunc("/health", a.fetchGophers).Methods(http.MethodGet)

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}