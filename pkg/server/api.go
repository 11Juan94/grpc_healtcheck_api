package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
)

type api struct {
	router http.Handler
}

type Server interface {
	Router() http.Handler
}

func (a *api) checkHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cmd := exec.Command("/bin/grpc_health_probe", "-addr="+os.Getenv("GRPC_HOST")+":"+os.Getenv("GRPC_PORT"), "-connect-timeout", "250ms", "-rpc-timeout", "100ms")
	if err := cmd.Start(); err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	timer := time.AfterFunc(2*time.Second, func() {
		cmd.Process.Kill()
	})
	err := cmd.Wait()
	timer.Stop()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(err.Error())
	w.WriteHeader(http.StatusOK)
}

func New() Server {
	a := &api{}

	r := mux.NewRouter()
	r.HandleFunc("/health", a.checkHealth).Methods(http.MethodGet)
	muxWithMiddlewares := http.TimeoutHandler(r, time.Second*3, "Timeout from server.")
	a.router = muxWithMiddlewares
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}
