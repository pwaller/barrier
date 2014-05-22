// Copyright 2014 Peter Waller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A barrier primitive which can be used to signal a permanent state change,
// for example to signal that shutdown should occur.
//
// Example:
//
//	package main
//
//	import (
//		"sync"
//
//		"github.com/pwaller/barrier"
//	)
//
//	func main() {
//		var w sync.WaitGroup
//		defer w.Wait() // Main should wait for its goroutines
//
//		var b barrier.Barrier
//
//		w.Add(1)
//		go func() {
//			defer w.Done()
//			defer b.Fall()
//			println("GO!")
//			<-b.Barrier() // Many goroutines can wait on the barrier
//		}()
//
//		w.Add(1)
//		go func() {
//			defer w.Done()
//			defer b.Fall()
//			println("GO!")
//			// When this goroutine happens to return,
//			// all barrier waits can be passed.
//			return
//		}()
//
//	}
//
//
package barrier

import (
	"sync"
)

// The zero of Barrier is a ready-to-use value
type Barrier struct {
	channel  chan struct{}
	o        sync.Once
	initOnce sync.Once
}

func (b *Barrier) init() {
	b.initOnce.Do(func() {
		b.channel = make(chan struct{})
	})
}

// `b.Fall()` can be called any number of times and causes the channel returned
// by `b.Barrier()` to become closed (permanently available for immediate reading)
func (b *Barrier) Fall() {
	b.init()
	b.o.Do(func() { close(b.channel) })
}

// When `b.Fall()` is called, the channel returned by Barrier() is closed
// (and becomes always readable)
func (b *Barrier) Barrier() <-chan struct{} {
	b.init()
	return b.channel
}
