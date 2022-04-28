package main

func createSequence() <-chan int {
	ch := make(chan int)

	go func() {
		for i := 1; ; i++ {
			ch <- i
		}
	}()

	return ch
}
