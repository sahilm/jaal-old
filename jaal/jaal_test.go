package jaal_test

import (
	"errors"
	"testing"

	"bytes"

	"time"

	"github.com/sahilm/jaal/jaal"
)

func TestListen(t *testing.T) {
	t.Run("it logs all events", func(t *testing.T) {
		tl := &testEventLog{}
		listener := newTestListener()
		errLogger := jaal.NewSystemLogger(bytes.NewBuffer([]byte{}), "")
		go jaal.Listen([]jaal.Listener{listener}, tl, errLogger)
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
	LoggedEvents []*jaal.Event
}

func (tl *testEventLog) Log(event *jaal.Event) {
	tl.LoggedEvents = append(tl.LoggedEvents, event)
}

type testListener struct {
	listenDone chan bool
}

type data struct {
	foo string
	bar uint
}

var event *jaal.Event

func (t *testListener) Listen(eventHandler func(*jaal.Event), errHandler func(error)) {
	event = &jaal.Event{
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
		//err := errors.New("I'm wrapped")
		//fe := &jaal.FatalError{
		//	Err: err,
		//}

		//got := fe.Error()
		//want := err.Error()
		//test.AssertEqualString(t, got, want)
	})
}
