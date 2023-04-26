package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dc-dc-dc/finify"
	"github.com/dc-dc-dc/finify/commands"
)

var (
	SaveOutput = flag.Bool("save", false, "save output to file")
)

func init() {
	flag.Parse()
}

func main() {
	// load openAI GPT token
	// generate a prompt with commands and context
	// send to openAI,
	// it should respond with json, will have to fix the response.
	// continue execution based on the response
	// repeat until "finished" commnand is sent

	// summarize earnings reports and any other fillings for companys on the stock market
	// get a sentiment analysis for a company
	// generate charts in python
	fmt.Println(len(os.Args))

	OpenAIApiKey := os.Getenv("OPENAI_API_KEY")
	NewsAPIKey := os.Getenv("NEWS_API_KEY")
	if OpenAIApiKey == "" {
		panic("OPENAI_API_KEY not set")
	}

	if NewsAPIKey == "" {
		panic("NEWS_API_KEY not set")
	}
	return
	cm := finify.NewCommandManager()
	cm.AddCommands(commands.DefaultCommands())
	cm.AddCommand(commands.NewQueryNewsCommand(NewsAPIKey))
	prompt := `
		figure out how the public feels about meta announcement that it will now focus more on AI versus VR
	`
	ctx := context.Background()
	agent := finify.NewAgent(
		"finify",
		finify.GenerateSystemPrompt(prompt, cm.GetCommands()),
		finify.DefaultTriggerPrompt,
		OpenAIApiKey,
		cm.HandleCommand,
		commandLineHandler,
	)
	if err := agent.Start(ctx); err != nil {
		panic(err)
	}
	// fmt.Printf("%+v \n", res)
}

func commandLineHandler(ctx context.Context, agent *finify.Agent, res *finify.DefaultFormatResponse) (bool, error) {
	PrintFormattedResponse(agent.Name, agent.GetCount(), res)
	fmt.Println("Continue [y/N]?")
	var userResponse string
	fmt.Scanf("%s", &userResponse)
	if strings.ToLower(userResponse) != "y" {
		return false, nil
	}
	return true, nil
}

func PrintFormattedResponse(name string, count int, res *finify.DefaultFormatResponse) {
	fmt.Printf("Agent %s, count %d\n\n", name, count)
	fmt.Printf("\tThought - %s\n\n", res.Thoughts.Text)
	fmt.Printf("\tReasoning - %s\n\n", res.Thoughts.Reasoning)
	if res.Command.Name != "" {
		fmt.Printf("\tCommand: %s, args: %+v\n", res.Command.Name, res.Command.Args)
	}
}
