package slack

import (
	"context"
	"encoding/json"
	_ "fmt"
	"net/http"
	"testing"
)

func TestWebhook(t *testing.T) {

	ctx := context.Background()

	m := &Message{
		Channel: "test",
		Text:    "hello world",
	}

	wh, err := NewWebhook(ctx, "http://localhost:9876")

	if err != nil {
		t.Fatalf("Failed to create new webhook, %v", err)
	}

	message_handler := func(rsp http.ResponseWriter, req *http.Request) {

		var m *Message

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

	err = wh.Post(ctx, m)

	if err != nil {
		t.Fatalf("Failed to post message, %v", err)
	}
}
