// Package waitgroup test cases reference golang sync.WaitGroup https://golang.org/src/sync/waitgroup_test.go.
package waitgroup

import (
	"testing"
)

func TestWaitGroup(t *testing.T) {
	wg1 := &SafeWaitGroup{}
	wg2 := &SafeWaitGroup{}
	n := 16
	wg1.Add(n)
	wg2.Add(n)
	exited := make(chan bool, n)
	for i := 0; i != n; i++ {
		go func(i int) {
			wg1.Done()
			wg2.Wait()
			exited <- true
		}(i)
	}
	wg1.Wait()
	for i := 0; i != n; i++ {
		select {
		case <-exited:
			t.Fatal("SafeWaitGroup released group too soon")
		default:
		}
		wg2.Done()
	}
	for i := 0; i != n; i++ {
		<-exited // Will block if barrier fails to unlock someone.
	}
}

func TestWaitGroupAddFail(t *testing.T) {
	wg := &SafeWaitGroup{}
	wg.Add(1)
	wg.Done()
	wg.Wait()
	if err := wg.Add(1); err == nil {
		t.Errorf("Should return error when add positive after Wait")
	}
}
