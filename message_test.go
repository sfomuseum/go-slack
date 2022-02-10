package slack

import (
	"encoding/json"
	"testing"
)

func TestSlackMessage(t *testing.T) {

	m := &Message{
		Channel: "test",
		Text:    "hello world",
	}

	_, err := json.Marshal(m)

	if err != nil {
		t.Fatalf("Failed to encode message, %v", err)
	}
}
