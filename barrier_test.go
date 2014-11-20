package barrier

import (
	"sync"
	"testing"
)

func TestFall(t *testing.T) {
	var w sync.WaitGroup
	defer w.Wait() // Main should wait for its goroutines

	var b Barrier

	w.Add(1)
	go func() {
		defer w.Done()
		defer b.Fall()
		<-b.Barrier() // Many goroutines can wait on the barrier
	}()

	w.Add(1)
	go func() {
		defer w.Done()
		defer b.Fall()
		// If any goroutine finishes, they all finish
	}()
}

func TestForwards(t *testing.T) {
	var w sync.WaitGroup
	defer w.Wait() // Main should wait for its goroutines

	var b Barrier
	var b2 Barrier

	// If b falls, b2 should fall
	b.Forward(&b2)

	w.Add(1)
	go func() {
		defer w.Done()
		defer b.Fall()
		<-b2.Barrier()
		<-b.Barrier()
	}()

	w.Add(1)
	go func() {
		defer w.Done()
		defer b.Fall()
	}()
}

func TestForwardFallen(t *testing.T) {
	// When a fallen barrier is forwarded, the forwardee should immediately fall
	var b, f Barrier
	b.Fall()
	b.Forward(&f)
	<-f.Barrier()
}

func BenchmarkBarrier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b Barrier
		go b.Fall()
		<-b.Barrier()
	}
}

func BenchmarkBarrierNoGoroutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b Barrier
		b.Fall()
		<-b.Barrier()
	}
}

func BenchmarkForward(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b, f Barrier
		b.Forward(&f)
		go b.Fall()
		<-f.Barrier()
	}
}

func BenchmarkForward2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b, f, g Barrier
		b.Forward(&f)
		f.Forward(&g)
		go b.Fall()
		<-g.Barrier()
	}
}

func BenchmarkForward3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b, f, g, h Barrier
		b.Forward(&f)
		f.Forward(&g)
		g.Forward(&h)
		go b.Fall()
		<-h.Barrier()
	}
}
