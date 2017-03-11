package util

import "sync"

var wg sync.WaitGroup

// WaitGroupWrapper 异步协程执行体封装
type WaitGroupWrapper struct {
	sync.WaitGroup
}

// Wrap 封装
func (wgw *WaitGroupWrapper) Wrap(cb func()) {
	wg.Add(1)
	go func() {
		cb()
		wg.Done()
	}()
}
