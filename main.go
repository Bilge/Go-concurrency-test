package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func main() {
	startServer()

	for response := range downloadPages(10) {
		fmt.Println(response)
	}
}

func startServer() {
	pageIdGenerator := createSequence()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(pad(strconv.Itoa(<-pageIdGenerator), 2)))
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

func createSequence() <-chan int {
	ch := make(chan int)

	go func() {
		for i := 1; ; i++ {
			ch <- i
		}
	}()

	return ch
}

func pad(str string, n int) string {
	if len(str) < n {
		return strings.Repeat(`0`, n-len(str)) + str
	}

	return str
}

func downloadPages(n int) <-chan string {
	ch := make(chan string)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer waitGroup.Done()

			if err := downloadPage(ch); err != nil {
				fmt.Println(`Error:`, err)
			}
		}()
	}

	go func() {
		waitGroup.Wait()
		close(ch)
	}()

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
