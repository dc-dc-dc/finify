package finify

import "context"

type Command interface {
	GetPGCommand() PGCommand
	Call(ctx context.Context, agent *Agent, args map[string]interface{}) (string, error)
}

type PGCommand struct {
	Label string            `json:"label"`
	Name  string            `json:"name"`
	Args  map[string]string `json:"args,omitempty"`
}
