package commands

import "github.com/dc-dc-dc/finify"

func DefaultCommands() []finify.Command {
	return []finify.Command{
		&FinishedCommand{},
		&UserPromptCommand{},
		&GoogleSearchCommand{},
	}
}
