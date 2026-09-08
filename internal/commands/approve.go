package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/go-github/v91/github"
)

// Roles that are allowed to call this command
var allowedRoles = map[string]bool{
	"admin":    true,
	"maintain": true,
}

// Command to approve a PR
type ApproveCommand struct{}

func (command *ApproveCommand) Handle(ctx context.Context, client *github.Client, event *github.IssueCommentEvent, args []string) error {

	prNumber := *event.Issue.Number
	prNumberAttr := slog.Int("number", prNumber)
	slog.Debug("Requested approval for pull request", prNumberAttr)

	// Check that caller has enough permission
	permissionLevel, _, err := client.Repositories.GetPermissionLevel(
		ctx,
		*event.Repo.Owner.Login,
		*event.Repo.Name,
		*event.Sender.Login,
	)
	if err != nil {
		return fmt.Errorf("failed to retrieve permission level: %w", err)
	}

	if !allowedRoles[*permissionLevel.RoleName] {
		_, _, err := client.Issues.CreateComment(
			ctx,
			*event.Repo.Owner.Login,
			*event.Repo.Name,
			prNumber,
			github.IssueCommentRequest{
				Body: fmt.Sprintf(
					"@%s Must be maintainer or admin to call this command",
					*event.Sender.Login,
				),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to post comment: %w", err)
		}
		return nil
	}

	slog.Debug("Approving pull request", prNumberAttr)

	// Submit approval
	eventValue := "APPROVE"
	_, _, err = client.PullRequests.CreateReview(
		ctx,
		*event.Repo.Owner.Login,
		*event.Repo.Name,
		prNumber,
		&github.PullRequestReviewRequest{
			Event: &eventValue,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create PR review: %w", err)
	}

	slog.Debug("Pull request approved", prNumberAttr)

	return nil

}

func (command *ApproveCommand) Identifier() string {
	return "approve"
}

func (command *ApproveCommand) Scope() CommandScope {
	return PULL_REQUEST_ONLY
}
