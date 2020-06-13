package server

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	// Create the command with our context
	cmd := exec.CommandContext(ctx, "/bin/grpc_health_probe", "-addr="+os.Getenv("GRPC_HOST")+":"+os.Getenv("GRPC_PORT"), "-connect-timeout", "250ms", "-rpc-timeout", "100ms")
	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		log.Print("Command timed out")
		http.Error(w, "Command timed out", http.StatusInternalServerError)
		return
	}
	// If there's no context error, we know the command completed (or errored).
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(out)
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
