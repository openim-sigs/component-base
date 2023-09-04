package watch_test

import (
	"reflect"
	"testing"

	. "github.com/openim-sigs/component-base/pkg/watch"
)

func TestFilter(t *testing.T) {
	table := []Event{
		{Type: Added, Object: testType("foo")},
		{Type: Added, Object: testType("bar")},
		{Type: Added, Object: testType("baz")},
		{Type: Added, Object: testType("qux")},
		{Type: Added, Object: testType("zoo")},
	}

	source := NewFake()
	filtered := Filter(source, func(e Event) (Event, bool) {
		return e, e.Object.(testType)[0] != 'b'
	})

	go func() {
		for _, item := range table {
			source.Action(item.Type, item.Object)
		}
		source.Stop()
	}()

	var got []string
	for {
		event, ok := <-filtered.ResultChan()
		if !ok {
			break
		}
		got = append(got, string(event.Object.(testType)))
	}

	if e, a := []string{"foo", "qux", "zoo"}, got; !reflect.DeepEqual(e, a) {
		t.Errorf("got %v, wanted %v", e, a)
	}
}

func TestFilterStop(t *testing.T) {
	source := NewFake()
	filtered := Filter(source, func(e Event) (Event, bool) {
		return e, e.Object.(testType)[0] != 'b'
	})

	go func() {
		source.Add(testType("foo"))
		filtered.Stop()
	}()

	var got []string
	for {
		event, ok := <-filtered.ResultChan()
		if !ok {
			break
		}
		got = append(got, string(event.Object.(testType)))
	}

	if e, a := []string{"foo"}, got; !reflect.DeepEqual(e, a) {
		t.Errorf("got %v, wanted %v", e, a)
	}
}

func TestRecorder(t *testing.T) {
	events := []Event{
		{Type: Added, Object: testType("foo")},
		{Type: Added, Object: testType("bar")},
		{Type: Added, Object: testType("baz")},
		{Type: Added, Object: testType("qux")},
		{Type: Added, Object: testType("zoo")},
	}

	source := NewFake()
	go func() {
		for _, item := range events {
			source.Action(item.Type, item.Object)
		}
		source.Stop()
	}()

	recorder := NewRecorder(source)
	for {
		_, ok := <-recorder.Interface.ResultChan()
		if !ok {
			break
		}
	}
	recordedEvents := recorder.Events()
	if !reflect.DeepEqual(recordedEvents, events) {
		t.Errorf("got %v, expected %v", recordedEvents, events)
	}
}
