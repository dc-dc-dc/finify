package finify

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type CommandHandler func(ctx context.Context, agent *Agent, name string, args map[string]interface{}) (string, error)
type ResponseHandler func(ctx context.Context, agent *Agent, response *DefaultFormatResponse) (bool, error)
type Agent struct {
	Name            string
	SystemPrompt    string
	TriggerPrompt   string
	apiKey          string
	memory          []string
	history         []OpenAIChatMessage
	commandHandler  CommandHandler
	responseHandler ResponseHandler
	count           int
}

const (
	DefaultTriggerPrompt = "Determine the next command to use and respond with the format specified above"
)

func NewAgent(name, systemPrompt, triggerPrompt, apiKey string, commandHandler CommandHandler, responseHandler ResponseHandler) *Agent {
	return &Agent{
		Name:            name,
		SystemPrompt:    systemPrompt,
		TriggerPrompt:   triggerPrompt,
		apiKey:          apiKey,
		commandHandler:  commandHandler,
		responseHandler: responseHandler,
		memory:          []string{},
		history:         []OpenAIChatMessage{},
		count:           0,
	}
}

func (a *Agent) GenerateSystemPromptMessage() []OpenAIChatMessage {
	t := time.Now()
	// TODO: Add relevant memory
	res := []OpenAIChatMessage{
		{Role: OPENAI_ROLE_SYSTEM, Content: a.SystemPrompt},
		{Role: OPENAI_ROLE_SYSTEM, Content: "the current time and date is " + t.Format("2006-01-02 15:04:05")},
	}

	// add history
	for i := len(a.history) - 1; i >= 0; i-- {
		res = append(res, a.history[i])
	}

	return res
}

func (agent *Agent) Start(ctx context.Context) error {
	for {
		res, err := agent.Next(ctx)
		if err != nil {
			return err
		}
		// Print the prompt and reasoning and continue if user allows
		canExec, err := agent.responseHandler(ctx, agent, res)
		if err != nil {
			return err
		}
		// check if command needs to be done
		commandRes := "Failed to execute the command"
		if !canExec {
			commandRes = "The user did not allow the command to be executed"
		}

		if canExec && res.Command.Name != "" {
			commandRes, err = agent.commandHandler(ctx, agent, res.Command.Name, res.Command.Args)
			if err != nil {
				return err
			}
			if res.Command.Name == "finished" {
				break
			}
		}
		agent.history = append(agent.history, OpenAIChatMessage{Role: OPENAI_ROLE_SYSTEM, Content: commandRes})
		agent.count += 1
	}
	return nil
}

func (agent *Agent) GetCount() int {
	return agent.count
}

func (agent *Agent) AddToMemory(reply, result string) {
	agent.memory = append(agent.memory, "Assistant Reply: "+reply+"\nResult: "+result+"\n")
}

func (agent *Agent) Next(ctx context.Context) (*DefaultFormatResponse, error) {
	messages := agent.GenerateSystemPromptMessage()

	res, err := OpenAIChatCall(ctx, agent.apiKey, OPENAI_GPT_3_5_TURBO_0301, messages)
	if err != nil {
		return nil, err
	}
	aiResponse := strings.TrimSpace(res.Choices[0].Message.Content)
	agent.history = append(agent.history, OpenAIChatMessage{Role: OPENAI_ROLE_ASSISTANT, Content: aiResponse})
	var found = false
	for i, c := range aiResponse {
		// Find the opening bracket
		if c == '{' {
			aiResponse = aiResponse[i:]
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("could not find opening bracket in response. response:\n %s", aiResponse)
	}
	// Try to parse the json
	// aiResponse, _ = strconv.Unquote(aiResponse)
	var formattedResponse *DefaultFormatResponse
	if err := json.Unmarshal([]byte(aiResponse), &formattedResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: err - %w, response - \n %s", err, aiResponse)
	}
	return formattedResponse, nil
}

func (agent *Agent) GetHistory() []OpenAIChatMessage {
	return agent.history
}
