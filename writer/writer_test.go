package writer

import (
	"encoding/json"
	_ "fmt"
	"github.com/sfomuseum/go-slack"
	"io"
	"net/http"
	"testing"
)

func TestSlackWriter(t *testing.T) {

	wh_uri := "http://localhost:9876"
	wh_channel := "test"

	wr, err := NewSlackWriter(wh_uri, wh_channel)

	if err != nil {
		t.Fatalf("Failed to create new Slack writer, %v", err)
	}

	mw := io.MultiWriter(wr)

	message_handler := func(rsp http.ResponseWriter, req *http.Request) {

		var m *slack.Message

		dec := json.NewDecoder(req.Body)
		err := dec.Decode(&m)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if m.Text != "hello world" {
			http.Error(rsp, "Invalid message", http.StatusBadRequest)
			return
		}

		return
	}

	s := &http.Server{
		Addr:    "localhost:9876",
		Handler: http.HandlerFunc(message_handler),
	}

	go func() {

		err := s.ListenAndServe()

		if err != nil {
			t.Fatalf("Failed to start server, %v", err)
		}

	}()

	_, err = mw.Write([]byte("hello world"))

	if err != nil {
		t.Fatalf("Failed to write data to Slack writer, %v", err)
	}

}
