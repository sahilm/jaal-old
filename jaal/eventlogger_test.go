package jaal

import (
	"bytes"
	"testing"

	"encoding/json"

	"github.com/sahilm/jaal/jaal"
)

func TestEventLogger(t *testing.T) {
	t.Run("it logs events", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{})
		el := jaal.NewSystemLogger(bytes.NewBuffer([]byte{}), "")
		l := jaal.NewEventLogger(buf, el, "")
		l.Log(&jaal.Event{
			Type:   "test",
			Source: "127.0.0.1",
		})

		got := &jaal.Event{}
		err := json.Unmarshal(buf.Bytes(), got)
		if err != nil {
			t.Error(err)
		}
		if got.Source != "127.0.0.1" {
			t.Errorf("got :%v, want: %v", got, "127.0.0.1")
		}

		//test.AssertEqualString(t, got.SourceHostName, "localhost")
		//test.AssertEqualString(t, got.Type, "test")
	})
}
