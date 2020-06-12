package server

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

type api struct {
	router http.Handler
}

type Server interface {
	Router() http.Handler
}

func (a *api) checkHealth(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("/bin/grpc_health_probe", "-addr=:"+os.Getenv("GRPC_PORT"), "-connect-timeout 250ms", "-rpc-timeout 100ms")
	stdout, err := cmd.Output()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, string(stdout), http.StatusBadGateway)
		return
	}
	json.NewEncoder(w).Encode("Healthy")
	w.WriteHeader(http.StatusOK)
}

func New() Server {
	a := &api{}

	r := mux.NewRouter()
	r.HandleFunc("/health", a.checkHealth).Methods(http.MethodGet)

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}
