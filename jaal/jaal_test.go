package jaal

import (
	"errors"
	"testing"

	"bytes"

	"time"
)

func TestListen(t *testing.T) {
	t.Run("it logs all events", func(t *testing.T) {
		tl := &testEventLog{}
		listener := newTestListener()
		errLogger := NewSystemLogger(bytes.NewBuffer([]byte{}))
		go Listen([]Listener{listener}, tl, errLogger)
		timeout := time.After(100 * time.Millisecond)
		select {
		case <-timeout:
			t.Fatal("timed out")
		case <-listener.listenDone:
			got := tl.LoggedEvents[0]
			want := event
			if got != want {
				t.Errorf("got: %v, want: %v", got, want)
			}
		}
	})
}

type testEventLog struct {
	LoggedEvents []*Event
}

func (tl *testEventLog) Log(event *Event) {
	tl.LoggedEvents = append(tl.LoggedEvents, event)
}

type testListener struct {
	listenDone chan bool
}

type data struct {
	foo string
	bar uint
}

var event *Event

func (t *testListener) Listen(eventHandler func(*Event), errHandler func(interface{})) {
	event = &Event{
		Data: &data{
			foo: "something",
			bar: 9,
		},
		Source:        "test",
		CorrelationID: "1234",
		UnixTime:      1503956846,
	}

	eventHandler(event)
	errHandler(errors.New("an error"))
	t.listenDone <- true
}

func newTestListener() *testListener {
	return &testListener{make(chan bool)}
}

func TestFatalError(t *testing.T) {
	t.Run("it wraps its underlying error", func(t *testing.T) {
		err := errors.New("I'm wrapped")
		fe := &FatalError{
			Err: err,
		}

		got := fe.Error()
		want := err.Error()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
