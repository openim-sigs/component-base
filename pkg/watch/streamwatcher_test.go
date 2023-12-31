// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package watch_test

import (
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/openim-sigs/component-base/pkg/runtime"
	. "github.com/openim-sigs/component-base/pkg/watch"
)

type fakeDecoder struct {
	items chan Event
	err   error
}

func (f fakeDecoder) Decode() (action EventType, object runtime.Object, err error) {
	if f.err != nil {
		return "", nil, f.err
	}
	item, open := <-f.items
	if !open {
		return action, nil, io.EOF
	}
	return item.Type, item.Object, nil
}

func (f fakeDecoder) Close() {
	if f.items != nil {
		close(f.items)
	}
}

type fakeReporter struct {
	err error
}

func (f *fakeReporter) AsObject(err error) runtime.Object {
	f.err = err
	return runtime.Unstructured(nil)
}

func TestStreamWatcher(t *testing.T) {
	table := []Event{
		{Type: Added, Object: testType("foo")},
	}

	fd := fakeDecoder{items: make(chan Event, 5)}
	sw := NewStreamWatcher(fd, nil)

	for _, item := range table {
		fd.items <- item
		got, open := <-sw.ResultChan()
		if !open {
			t.Errorf("unexpected early close")
		}
		if e, a := item, got; !reflect.DeepEqual(e, a) {
			t.Errorf("expected %v, got %v", e, a)
		}
	}

	sw.Stop()
	_, open := <-sw.ResultChan()
	if open {
		t.Errorf("Unexpected failure to close")
	}
}

func TestStreamWatcherError(t *testing.T) {
	fd := fakeDecoder{err: fmt.Errorf("test error")}
	fr := &fakeReporter{}
	sw := NewStreamWatcher(fd, fr)
	evt, ok := <-sw.ResultChan()
	if !ok {
		t.Fatalf("unexpected close")
	}
	if evt.Type != Error || evt.Object != runtime.Unstructured(nil) {
		t.Fatalf("unexpected object: %#v", evt)
	}
	_, ok = <-sw.ResultChan()
	if ok {
		t.Fatalf("unexpected open channel")
	}

	sw.Stop()
	_, ok = <-sw.ResultChan()
	if ok {
		t.Fatalf("unexpected open channel")
	}
}

func TestStreamWatcherRace(t *testing.T) {
	fd := fakeDecoder{err: fmt.Errorf("test error")}
	fr := &fakeReporter{}
	sw := NewStreamWatcher(fd, fr)
	time.Sleep(10 * time.Millisecond)
	sw.Stop()
	time.Sleep(10 * time.Millisecond)
	_, ok := <-sw.ResultChan()
	if ok {
		t.Fatalf("unexpected pending send")
	}
}
