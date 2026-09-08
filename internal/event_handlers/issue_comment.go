package event_handlers

import (
	"context"
	"log"

	"github.com/google/go-github/v91/github"
)

func HandleIssueComment(ctx context.Context, deliveryID string, eventName string, event *github.IssueCommentEvent) error {

	log.Printf("%s made a comment!", *event.Sender.Login)
	return nil

}
