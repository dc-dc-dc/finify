package finify

import "context"

type Command interface {
	// Give a short description of the capability
	Resource() string
	// Returns the pg command for the pre prompt
	PGCommand() PGCommand

	// Handle the command when the LLM requests it
	Call(ctx context.Context, agent *Agent, args map[string]interface{}) (string, error)
}

type PGCommand struct {
	Label string            `json:"label"`
	Name  string            `json:"name"`
	Args  map[string]string `json:"args,omitempty"`
}
