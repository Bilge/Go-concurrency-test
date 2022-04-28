package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func main() {
	startServer()

	for response := range downloadPages(10) {
		fmt.Println(response)
	}
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
