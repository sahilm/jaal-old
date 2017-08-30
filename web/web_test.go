package web_test

import (
	"testing"

	"time"

	"net/http"

	"encoding/json"

	"github.com/sahilm/jaal"
	"github.com/sahilm/jaal/web"
)

func TestServer_Listen(t *testing.T) {
	t.Run("it returns event with request metadata", func(t *testing.T) {
		s := &web.Server{
			Address: ":8080",
		}
		ch := make(chan *jaal.Event, 1)
		errchan := make(chan error, 1)

		go s.Listen(func(event *jaal.Event) {
			ch <- event
		}, func(e error) {
			errchan <- e
		})

		_, err := http.Get("http://127.0.0.1:8080/")
		if err != nil {
			t.Errorf("got error: %v, want no error", err)
		}

		timeout := time.After(100 * time.Millisecond)
		select {
		case <-timeout:
			t.Errorf("timed out")
		case e := <-errchan:
			t.Errorf("got error: %v, want no error", e)
		case event := <-ch:
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
		t.Errorf("got: %v, want: %v", got, want)
	}
}
