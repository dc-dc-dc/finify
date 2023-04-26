package finify

import (
	"context"
	"fmt"
)

type CommandManager struct {
	commands map[string]Command
}

func NewCommandManager() *CommandManager {
	cm := &CommandManager{
		commands: make(map[string]Command),
	}
	return cm
}

func (cm *CommandManager) HandleCommand(ctx context.Context, agent *Agent, name string, args map[string]interface{}) (string, error) {
	command, ok := cm.commands[name]
	if !ok {
		return "", fmt.Errorf("command %s not found", name)
	}
	return command.Call(ctx, agent, args)
}

func (cm *CommandManager) AddCommand(command Command) error {
	name := command.GetPGCommand().Name
	if _, ok := cm.commands[name]; ok {
		return fmt.Errorf("command %s already exists", name)
	}
	cm.commands[name] = command
	return nil
}

func (cm *CommandManager) AddCommands(command []Command) error {
	for _, command := range command {
		if err := cm.AddCommand(command); err != nil {
			return err
		}
	}
	return nil
}

func (cm *CommandManager) GetCommands() []PGCommand {
	commands := make([]PGCommand, len(cm.commands))
	var index int
	for _, command := range cm.commands {
		commands[index] = command.GetPGCommand()
		index += 1
	}
	return commands
}
