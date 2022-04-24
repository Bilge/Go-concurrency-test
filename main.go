package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var i = 0

func main() {
	startServer()

	ch := downloadPages(10)
	for response := range ch {
		fmt.Println(response)
	}
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHttpRequest)

	// Ensure socket is created synchronously.
	socket, err := net.Listen(`tcp4`, `localhost:http`)
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{Handler: mux}
	go server.Serve(socket)
}

func handleHttpRequest(writer http.ResponseWriter, request *http.Request) {
	i++
	writer.Write([]byte(pad(strconv.Itoa(i), 2)))
}

func pad(str string, n int) string {
	if len(str) < n {
		return strings.Repeat(`0`, n-len(str)) + str
	}

	return str
}

func downloadPages(n int) <-chan string {
	ch := make(chan string)

	for i := 0; i < n; i++ {
		go func() {
			if err := downloadPage(ch); err != nil {
				fmt.Println(`Error:`, err)
			}
		}()
	}

	return ch
}

func downloadPage(ch chan<- string) error {
	resp, err := http.Get("http://localhost")
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err == nil {
		ch <- string(bytes)
	}

	return err
}
