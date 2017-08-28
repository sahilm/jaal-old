package web_test

import (
	"testing"

	"time"

	"net/http"

	"github.com/sahilm/jaal"
	"github.com/sahilm/jaal/web"
)

func TestServer_Listen(t *testing.T) {
	t.Run("it calls event handler when a request is made", func(t *testing.T) {
		s := &web.Server{
			Address: ":8080",
		}
		ch := make(chan bool, 1)
		errchan := make(chan error, 1)

		go s.Listen(func(event *jaal.Event) {
			ch <- true
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
		case <-ch:
		}

	})
}
