package commands

import (
	"context"

	"github.com/google/go-github/v91/github"
)

// Scope of a command
type CommandScope int

const (
	// Command is allowed everywhere
	BOTH CommandScope = iota
	// Command is only allowed in issues
	ISSUE_ONLY
	// Command is only allowed in pull requests
	PULL_REQUEST_ONLY
)

// A command that can be invoked by a comment
type CommentCommand interface {
	// Handles the invocation
	Handle(ctx context.Context, client *github.Client, event *github.IssueCommentEvent, args []string) error
	// The identifier of the command
	Identifier() string
	// The scope that the command is allowed in
	Scope() CommandScope
}

// Centralized command registry
type CommandRegistry struct {
	// Supported commands
	commands map[string]CommentCommand
}

// Creates a new registry
func NewRegistry() *CommandRegistry {

	// Initialize command list
	commands := [...]CommentCommand{
		&ApproveCommand{},
	}

	// Convert list into a map
	commandMap := make(map[string]CommentCommand, len(commands))
	for _, command := range commands {
		commandMap[command.Identifier()] = command
	}

	return &CommandRegistry{
		commands: commandMap,
	}
}

// Finds the command with the given identifier, or returns nil if there is none
func (registry *CommandRegistry) GetCommand(identifier string) CommentCommand {
	return registry.commands[identifier]
}
