package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
)

func startServer() {
	pageIdGenerator := createSequence()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if _, err := writer.Write([]byte(pad(strconv.Itoa(<-pageIdGenerator), 2))); err != nil {
			log.Fatal(err)
		}
	})

	// Ensure socket is created synchronously so server is ready to accept connections before first request is made.
	socket, err := net.Listen(`tcp4`, `localhost:http`)
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{Handler: mux}
	go func() {
		if err := server.Serve(socket); err != nil {
			log.Fatal(err)
		}
	}()
}
