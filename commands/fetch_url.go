package commands

import (
	"context"

	"github.com/dc-dc-dc/finify"
)

type FetchUrlCommand struct{}

var _ finify.Command = (*FetchUrlCommand)(nil)

func (command *FetchUrlCommand) Resource() string {
	return "load up a url and return the contents"
}

func (command *FetchUrlCommand) PGCommand() finify.PGCommand {
	return finify.PGCommand{
		Label: "load the contents of a url",
		Name:  "fetch_url",
		Args: map[string]string{
			"url": "the url to fetch",
		},
	}
}

func (command *FetchUrlCommand) Call(ctx context.Context, agent *finify.Agent, args map[string]interface{}) (string, error) {
	return "", nil
}
