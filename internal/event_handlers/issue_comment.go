package event_handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/go-github/v91/github"
	"github.com/tmarback/github-helper-app/internal/commands"
)

// Handles issue comment events
func (handler *EventHandler) HandleIssueComment(ctx context.Context, deliveryID string, eventName string, event *github.IssueCommentEvent) error {

	// Check for command prefix
	message, isCommand := strings.CutPrefix(*event.Comment.Body, "/")
	if !isCommand {
		return nil
	}

	// Split the command into its parts
	parts := strings.Fields(message)
	if len(parts) == 0 {
		slog.Debug("Blank command")
		return nil
	}
	identifier := parts[0]
	args := parts[1:]

	// Identify the command to run
	command := handler.commandRegistry.GetCommand(identifier)
	if command == nil {
		slog.Debug("Command not found", slog.String("identifier", identifier))
		return nil
	}

	// Check if scope is appropriate
	if event.Issue.IsPullRequest() {
		if command.Scope() == commands.ISSUE_ONLY {
			slog.Debug("Issue-only command invoked in PR", slog.String("identifier", identifier))
			return nil
		}
	} else {
		if command.Scope() == commands.PULL_REQUEST_ONLY {
			slog.Debug("PR-only command invoked in issue", slog.String("identifier", identifier))
			return nil
		}
	}

	// Get Github client
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
		return fmt.Errorf("error while acking command: %w", err)
	}

	// Run the command
	if err := command.Handle(ctx, client, event, args); err != nil {
		return fmt.Errorf("error while handling command: %w", err)
	}

	return nil

}
