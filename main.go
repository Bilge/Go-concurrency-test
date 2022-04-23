package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var i = 0

type HttpBody struct {
	body string
	error
}

func main() {
	startServer()

	ch := make(chan HttpBody)
	for i := 0; i < 10; i++ {
		go downloadPage(ch)
	}

	for response := range ch {
		if response.error != nil {
			fmt.Print(`Error!:`, response)
		}

		fmt.Println(response.body)
	}
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHttpRequest)

	// Ensure socket is created synchronously.
	socket, _ := net.Listen(`tcp4`, `localhost:http`)
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

func downloadPage(ch chan<- HttpBody) {
	resp, err := http.Get("http://localhost")

	if err != nil {
		ch <- HttpBody{``, err}

		return
	}

	bytes, err := io.ReadAll(resp.Body)

	ch <- HttpBody{string(bytes), err}
}
