package event_handlers

import (
	"context"
	"log"

	"github.com/google/go-github/v91/github"
)

// Handles issue comment events
func HandleIssueComment(ctx context.Context, deliveryID string, eventName string, event *github.IssueCommentEvent) error {

	// TODO
	log.Printf("%s made a comment!", *event.Sender.Login)
	return nil

}
