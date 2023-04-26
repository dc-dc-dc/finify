package finify

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkoukk/tiktoken-go"
)

type CommandHandler func(ctx context.Context, agent *Agent, name string, args map[string]interface{}) (string, error)
type ResponseHandler func(ctx context.Context, agent *Agent, response *DefaultFormatResponse) (bool, error)
type Agent struct {
	Name            string
	SystemPrompt    string
	TriggerPrompt   string
	apiKey          string
	memory          []string
	model           string
	history         []OpenAIChatMessage
	commandHandler  CommandHandler
	responseHandler ResponseHandler
	count           int
	maxRetries      int
	maxTokens       int
	encoder         *tiktoken.Tiktoken
}

const (
	DefaultTriggerPrompt = "Determine the next command to use and respond with the format specified above"
)

func NewAgent(name, systemPrompt, triggerPrompt, apiKey, model string, reservedResponseTokens int, commandHandler CommandHandler, responseHandler ResponseHandler) *Agent {
	maxTokens, ok := modelSet[model]
	if !ok {
		panic(fmt.Errorf("invalid model: %s", model))
	}
	maxTokens = maxTokens - reservedResponseTokens
	encoder, err := tiktoken.EncodingForModel(model)
	if err != nil {
		panic(err)
	}
	return &Agent{
		Name:            name,
		SystemPrompt:    systemPrompt,
		TriggerPrompt:   triggerPrompt,
		apiKey:          apiKey,
		commandHandler:  commandHandler,
		responseHandler: responseHandler,
		memory:          []string{},
		history:         []OpenAIChatMessage{},
		model:           model,
		count:           0,
		maxRetries:      3,
		maxTokens:       maxTokens,
		encoder:         encoder,
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
	index := len(a.history) - 1
	currentTokenCount := a.getTokenCount(res)

	for index >= 0 && currentTokenCount < a.maxTokens {
		fmt.Printf("current token count: %d\n", currentTokenCount)
		nextTokenCount := a.getTokenCount([]OpenAIChatMessage{a.history[index]})
		if nextTokenCount+currentTokenCount > a.maxTokens {
			break
		}
		res = append(res, a.history[index])
		currentTokenCount += nextTokenCount
		index -= 1
	}

	return res
}

func (a *Agent) getTokenCount(messages []OpenAIChatMessage) int {
	var t string
	for _, s := range messages {
		t += s.Content
	}
	return len(a.encoder.Encode(t, nil, nil))
}

func (agent *Agent) Start(ctx context.Context) (err error) {
	for {
		var res *DefaultFormatResponse
		var tries int
		for {
			res, err = agent.Next(ctx)
			if err == nil {
				break
			} else {
				if tries >= agent.maxRetries {
					return err
				}
				tries += 1
			}
			fmt.Printf("got an error, waiting 10 seconds and trying again attempt: %d err: %s\n", tries, err.Error())
			time.Sleep(10 * time.Second)
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
			if res.Command.Name == "goal_completed" {
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

	res, err := OpenAIChatCall(ctx, agent.apiKey, agent.model, messages)
	if err != nil {
		return nil, err
	}
	aiResponse := strings.TrimSpace(res.Choices[0].Message.Content)
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
	agent.history = append(agent.history, OpenAIChatMessage{Role: OPENAI_ROLE_ASSISTANT, Content: aiResponse})
	return formattedResponse, nil
}

func (agent *Agent) GetHistory() []OpenAIChatMessage {
	return agent.history
}
