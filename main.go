package main

import (
	"jwt-demo/transport"
	"log"
	"net/http"
)

const (
	SERVER_ADDR = ":8080"
)

func main() {
	http.HandleFunc("/login", transport.Login)
	http.HandleFunc("/validatetest", transport.TestValidate)
	server := new(http.Server)
	server.Addr = SERVER_ADDR
	log.Println("Starting server at", SERVER_ADDR)
	server.ListenAndServe()
}
