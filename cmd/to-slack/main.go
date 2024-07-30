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
	"github.com/sfomuseum/runtimevar"
)

func main() {

	channel := flag.String("channel", "", "A valid Slack channel to post to.")
	runtimevar_uri := flag.String("webhook-uri", "", "A valid gocloud.dev/runtimevar URI containing a Slack webhook URL to post to.")
	stdin := flag.Bool("stdin", false, "Read input from STDIN")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Post a message to a Slack channel.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] message\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	ctx := context.Background()

	// Create webhook thingy

	webhook_uri, err := runtimevar.StringVar(ctx, *runtimevar_uri)

	if err != nil {
		log.Fatalf("Failed to derive webhook uri from runtimevar, %v", err)
	}

	wh, err := slack.NewWebhook(ctx, webhook_uri)

	if err != nil {
		log.Fatalf("Failed to create new webhook, %v", err)
	}

	// Create message

	foo := false
	
	raw := flag.Args()

	if *stdin {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {

			t := scanner.Text()

			if !foo {						
				raw = append(raw, t)
				continue
			}

			m := &slack.Message{
				Channel: *channel,
				Text:    t,
			}
			
			// Post message
			
			err = wh.Post(ctx, m)
			
			if err != nil {
				log.Fatalf("Failed to post message, %v", err)
			}
			
		}

		err := scanner.Err()

		if err != nil {
			log.Fatalf("Failed to read from STDIN, %v", err)
		}

	}

	if foo {
		return
	}
	
	text := strings.Join(raw, " ")

	m := &slack.Message{
		Channel: *channel,
		Text:    text,
	}

	// Post message

	err = wh.Post(ctx, m)

	if err != nil {
		log.Fatalf("Failed to post message, %v", err)
	}
}
