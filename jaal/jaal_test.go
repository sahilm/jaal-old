package jaal_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sahilm/jaal/jaal"
	"github.com/sahilm/jaal/test"
	logrus "github.com/sirupsen/logrus/hooks/test"
)

func TestListen(t *testing.T) {
	t.Run("it logs all events", func(t *testing.T) {
		logger, hook := logrus.NewNullLogger()
		listener := newTestListener()
		go jaal.Listen([]jaal.Listener{listener}, logger)

		timeout := time.After(100 * time.Millisecond)
		select {
		case <-timeout:
			t.Fatal("timed out")
		case <-listener.listenDone:
		}

		test.AssertEqualInt(t, len(hook.AllEntries()), 1)
	})
}

type testListener struct {
	listenDone chan bool
}

type data struct {
	foo string
	bar uint
}

func (d data) String() string {
	return fmt.Sprintf("foo=%v bar=%v", d.foo, d.bar)
}

func (t *testListener) Listen(eventHandler func(*jaal.Event), errHandler func(error)) {
	data := &data{
		foo: "something",
		bar: 9,
	}

	event := &jaal.Event{
		Data:          data,
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
		fe := &jaal.FatalError{
			Err: err,
		}

		got := fe.Error()
		want := err.Error()
		test.AssertEqualString(t, got, want)
	})
}
