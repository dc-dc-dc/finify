package commands

import (
	"context"
	"fmt"

	"github.com/dc-dc-dc/finify"
)

type UserPromptCommand struct{}

var _ finify.Command = &UserPromptCommand{}

func (fc *UserPromptCommand) GetPGCommand() finify.PGCommand {
	return finify.PGCommand{
		Label: "prompt the user for input",
		Name:  "user_prompt",
		Args: map[string]string{
			"prompt": "prompt",
		},
	}
}

func (fc *UserPromptCommand) Call(ctx context.Context, agent *finify.Agent, args map[string]interface{}) (string, error) {
	// dump all the messages to file
	prompt, ok := args["prompt"].(string)
	if !ok {
		return "", fmt.Errorf("you did not provide a prompt")
	}

	fmt.Println(agent.Name, ":", prompt)
	var userResponse string
	_, err := fmt.Scanln(&userResponse)
	if err != nil {
		return "", err
	}
	return agent.Name + ":" + prompt + "\nUser Response: " + userResponse, nil
}
