package commands

import (
	"context"
	"encoding/json"
	"os"

	"github.com/dc-dc-dc/finify"
)

type FinishedCommand struct{}

var _ finify.Command = &FinishedCommand{}

func (fc *FinishedCommand) GetPGCommand() finify.PGCommand {
	return finify.PGCommand{
		Label: "when you accomplished your goal",
		Name:  "finished",
		Args:  map[string]string{},
	}
}

func (fc *FinishedCommand) Call(ctx context.Context, agent *finify.Agent, args map[string]interface{}) (string, error) {
	// dump all the messages to file
	res, err := os.Create("./messages.json")
	if err != nil {
		return "", err
	}
	defer res.Close()

	if json.NewEncoder(res).Encode(agent.GetHistory()); err != nil {
		return "", err
	}

	return "", nil
}
