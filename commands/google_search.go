package commands

import (
	"context"

	"github.com/dc-dc-dc/finify"
	googlesearch "github.com/rocketlaunchr/google-search"
)

type GoogleSearchCommand struct{}

var _ finify.Command = (*GoogleSearchCommand)(nil)

func init() {
	// googlesearch.RateLimit.SetLimit(10)
}

func (command *GoogleSearchCommand) Resource() string {
	return "query the internet for information"
}

func (command *GoogleSearchCommand) PGCommand() finify.PGCommand {
	return finify.PGCommand{
		Label: "what you want to search for",
		Name:  "google_search",
		Args: map[string]string{
			"query": "the query you want to search for",
		},
	}
}

func (command *GoogleSearchCommand) Call(ctx context.Context, agent *finify.Agent, args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", finify.ErrCommandInvalidArgs
	}
	res, err := googlesearch.Search(ctx, query)
	if err != nil {
		return "", err
	}
	resString := "google results for " + query + "\n"

	for _, s := range res {
		resString += "Title: " + s.Title + "\nDescription:" + s.Description + "\nUrl:" + s.URL + "\n"
	}

	return resString, nil
}
