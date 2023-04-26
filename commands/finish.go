package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dc-dc-dc/finify"
)

type FinishedCommand struct{}

var _ finify.Command = &FinishedCommand{}

func (fc *FinishedCommand) Resource() string {
	return ""
}

func (fc *FinishedCommand) PGCommand() finify.PGCommand {
	return finify.PGCommand{
		Label: "when you completed your original goal",
		Name:  "goal_completed",
		Args: map[string]string{
			"summary": "summary of what you did",
		},
	}
}

func (fc *FinishedCommand) Call(ctx context.Context, agent *finify.Agent, args map[string]interface{}) (string, error) {
	// dump all the messages to file
	res, err := os.Create(fmt.Sprintf("./messages/%s.json", time.Now().Format("2006-01-02-15-04-05")))
	if err != nil {
		return "", err
	}
	defer res.Close()

	if json.NewEncoder(res).Encode(agent.GetHistory()); err != nil {
		return "", err
	}

	return "", nil
}
