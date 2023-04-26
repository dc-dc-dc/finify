package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dc-dc-dc/finify"
	"github.com/dc-dc-dc/finify/commands"
	"github.com/fatih/color"
)

var (
	SaveOutput = flag.Bool("save", false, "save output to file")
	Prompt     = flag.String("prompt", "figure out how the public feels about meta announcement that it will now focus more on AI versus VR", "prompt to use")
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

	OpenAIApiKey := os.Getenv("OPENAI_API_KEY")
	NewsAPIKey := os.Getenv("NEWS_API_KEY")
	if OpenAIApiKey == "" {
		panic("OPENAI_API_KEY not set")
	}

	if NewsAPIKey == "" {
		panic("NEWS_API_KEY not set")
	}
	cm := finify.NewCommandManager()
	cm.AddCommands(commands.DefaultCommands())
	cm.AddCommand(commands.NewQueryNewsCommand(NewsAPIKey))
	ctx := context.Background()
	agent := finify.NewAgent(
		"finify",
		finify.GenerateSystemPrompt(*Prompt, cm.GetCommands(), cm.GetResources()),
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
	fmt.Println()
	if res.Command.Name != "" {
		fmt.Println("Continue [y/N]?")
	} else {
		fmt.Println("Execute command ? [y/N]")
	}

	var userResponse string
	fmt.Scanf("%s", &userResponse)
	if strings.ToLower(userResponse) != "y" {
		return false, nil
	}
	return true, nil
}

func PrintFormattedResponse(name string, count int, res *finify.DefaultFormatResponse) {
	color.Set(color.Underline).Add(color.BgGreen)
	fmt.Printf("\nAgent %s, count %d\n\n", name, count)
	color.Set(color.Reset)

	color.Set(color.Underline).Add(color.BgMagenta)
	fmt.Printf("\tThought")
	color.Set(color.Reset)

	fmt.Printf(" - %s\n\n", res.Thoughts.Text)

	color.Set(color.Underline).Add(color.BgMagenta)
	fmt.Printf("\tReasoning ")
	color.Set(color.Reset)

	fmt.Printf("- %s\n\n", res.Thoughts.Reasoning)

	if res.Command.Name != "" {
		color.Set(color.Underline).Add(color.BgCyan)
		fmt.Printf("\tCommand: %s, args: %+v\n", res.Command.Name, res.Command.Args)
		color.Set(color.Reset)
	}
}
