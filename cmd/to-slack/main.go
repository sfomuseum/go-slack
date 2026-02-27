// post is a command-line tool for posting messages to a Slack channel. For example:
//
//	$> ./bin/post -channel test -webhook-uri 'constant://?val=https://hooks.slack.com/services/.../.../...' testing
//
// Or:
//
//	$> echo "wub wub wub" | ./bin/post -channel test -webhook-uri 'constant://?val=https://hooks.slack.com/services/.../.../...' -stdin
//
// Or:
//
//	$> echo "wub wub wub" | ./bin/post -channel test -webhook 'awsparamstore://{SECRET_NAME}?region={REGION}&credentials={CREDENTIALS}' -stdin
package main

import (
	_ "gocloud.dev/runtimevar/awsparamstore"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
)

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sfomuseum/go-slack"
	"github.com/aaronland/gocloud/runtimevar"
)

func main() {

	var channel string
	var runtimevar_uri string
	var stdin bool
	var message_per_line bool
	var prefix string

	flag.StringVar(&channel, "channel", "", "A valid Slack channel to post to.")
	flag.StringVar(&runtimevar_uri, "webhook-uri", "", "A valid gocloud.dev/runtimevar URI containing a Slack webhook URL to post to.")
	flag.StringVar(&prefix, "prefix", "", "Optional prefix to prepend each message with.")
	flag.BoolVar(&stdin, "stdin", false, "Read input from STDIN")
	flag.BoolVar(&message_per_line, "message-per-line", false, "Send a message for each line when reading input from STDIN.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Post a message to a Slack channel.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] message\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	ctx := context.Background()

	// Create webhook thingy

	webhook_uri, err := runtimevar.StringVar(ctx, runtimevar_uri)

	if err != nil {
		log.Fatalf("Failed to derive webhook uri from runtimevar, %v", err)
	}

	wh, err := slack.NewWebhook(ctx, webhook_uri)

	if err != nil {
		log.Fatalf("Failed to create new webhook, %v", err)
	}

	// Create message

	raw := flag.Args()

	if stdin {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {

			t := scanner.Text()

			if !message_per_line {
				raw = append(raw, t)
				continue
			}

			if prefix != "" {
				t = fmt.Sprintf("%s %s", prefix, t)
			}

			m := &slack.Message{
				Channel: channel,
				Text:    t,
			}

			// Post message

			err := wh.Post(ctx, m)

			if err != nil {
				log.Fatalf("Failed to post message, %v", err)
			}

		}

		err := scanner.Err()

		if err != nil {
			log.Fatalf("Failed to read from STDIN, %v", err)
		}

	}

	text := strings.Join(raw, " ")

	if text == "" {
		return
	}

	if prefix != "" {
		text = fmt.Sprintf("%s %s", prefix, text)
	}

	m := &slack.Message{
		Channel: channel,
		Text:    text,
	}

	err = wh.Post(ctx, m)

	if err != nil {
		log.Fatalf("Failed to post message, %v", err)
	}
}
