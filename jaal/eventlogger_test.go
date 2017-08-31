package jaal

import (
	"bytes"
	"testing"

	"encoding/json"
)

func TestEventLogger(t *testing.T) {
	t.Run("it logs events", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{})
		el := NewSystemLogger(bytes.NewBuffer([]byte{}), "")
		l := NewEventLogger(buf, el, "")
		testEvent := &Event{
			Type:   "test",
			Source: "127.0.0.1",
		}

		l.Log(testEvent)

		got := &Event{}
		err := json.Unmarshal(buf.Bytes(), got)

		if err != nil {
			t.Error(err)
		}

		if got.Source != "127.0.0.1" {
			t.Errorf("got :%v, want: %v", got, "127.0.0.1")
		}

		if got.Type != "test" {
			t.Errorf("got :%v, want: %v", got, "test")
		}
	})
}
