// Copyright 2014 Peter Waller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package barrier

import (
	"sync"
)

type Barrier struct {
	channel chan struct{}
	o       sync.Once
	init    sync.Once
}

func (b *Barrier) Init() {
	b.init.Do(func() {
		b.channel = make(chan struct{})
	})
}

func (b *Barrier) Fall() {
	b.Init()
	b.o.Do(func() { close(b.channel) })
}

func (b *Barrier) Barrier() <-chan struct{} {
	b.Init()
	return b.channel
}
