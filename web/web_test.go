package web_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/sahilm/jaal/jaal"
	"github.com/sahilm/jaal/web"
)

func TestServer_Listen(t *testing.T) {
	t.Run("it logs events with request metadata", func(t *testing.T) {
		s := &web.Server{
			Address: ":8080",
		}
		eventChan := make(chan *jaal.Event, 1)

		go s.Listen(func(event *jaal.Event) {
			eventChan <- event
		}, func(i interface{}) {
		})

		_, err := http.Get("http://localhost:8080/")
		if err != nil {
			t.Errorf("got error: %v, want no error", err)
		}

		timeout := time.After(100 * time.Millisecond)
		select {
		case <-timeout:
			t.Errorf("timed out")
		case event := <-eventChan:
			validateEvent(event, t)
		}
	})
}
func validateEvent(event *jaal.Event, t *testing.T) {
	b, err := json.Marshal(event.Data)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := `{"URI":"/","Method":"GET","Header":{"Accept-Encoding":["gzip"],"User-Agent":["Go-http-client/1.1"]}}`

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	got = event.Summary
	want = "received GET at / from localhost (::1)"

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	got = event.Source
	want = "::1"

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	got = event.Type
	want = "http"

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
