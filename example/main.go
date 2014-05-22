package main

import (
	"sync"

	"github.com/pwaller/barrier"
)

func main() {
	var w sync.WaitGroup
	defer w.Wait() // Main should wait for its goroutines

	var b barrier.Barrier

	w.Add(1)
	go func() {
		defer w.Done()
		defer b.Fall()
		println("GO!")
		<-b.Barrier() // Many goroutines can wait on the barrier
	}()

	w.Add(1)
	go func() {
		defer w.Done()
		defer b.Fall()
		println("GO!")
		// When this goroutine happens to return,
		// all barrier waits can be passed.
		return
	}()

}
