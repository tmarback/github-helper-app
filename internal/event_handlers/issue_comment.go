package event_handlers

import (
	"context"
	"log"
	"log/slog"

	"github.com/google/go-github/v91/github"
)

// Handles issue comment events
func (handler *EventHandler) HandleIssueComment(ctx context.Context, deliveryID string, eventName string, event *github.IssueCommentEvent) error {

	client, err := handler.getGithubClient(ctx, event.Installation)
	if err != nil {
		return err
	}

	// Acknowledge with a reaction
	if _, _, err = client.Reactions.CreateIssueCommentReaction(
		ctx,
		*event.Repo.Owner.Login,
		*event.Repo.Name,
		*event.Comment.ID,
		"+1",
	); err != nil {
		slog.Error("Error during event handling", slog.Any("error", err))
	}

	// TODO
	log.Printf("%s made a comment!", *event.Sender.Login)

	return nil

}
