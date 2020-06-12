package main

import (
	"log"
	"net/http"

	"github.com/11juan94/grpc_healtcheck_api/pkg/server"
)

func main() {
	s := server.New()
	log.Print("Server starting...")
	log.Fatal(http.ListenAndServe(":9192", s.Router()))
	log.Print("Server stoped!")
}
