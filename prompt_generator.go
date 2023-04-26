package finify

const (
	PG_DEFAULT_RESPONSE_FORMAT = `{"thoughts": {"text": "thought", "reasoning": "reasoning", "plan": "- short bulleted -\n list that conveys long term plan", "criticism": "constructive self-criticism"}, "command": {"name": "command name", "args": { "arg name": "value" }}}`
)

type DefaultFormatResponse struct {
	Thoughts struct {
		Text      string `json:"text"`
		Reasoning string `json:"reasoning"`
		Plan      string `json:"plan"`
		Criticism string `json:"criticism"`
	} `json:"thoughts"`
	Command struct {
		Name string                 `json:"name"`
		Args map[string]interface{} `json:"args,omitempty"`
	} `json:"command"`
}

type PromptGenerator struct {
	Constraints    []string    `json:"constraints"`
	Resources      []string    `json:"resources"`
	Commands       []PGCommand `json:"commands"`
	ResponseFormat string      `json:"response_format"`
}

func NewPromptGenerator() *PromptGenerator {
	return &PromptGenerator{
		Constraints:    []string{},
		Resources:      []string{},
		Commands:       []PGCommand{},
		ResponseFormat: PG_DEFAULT_RESPONSE_FORMAT,
	}
}

func (pg *PromptGenerator) AddConstraint(constraint ...string) {
	pg.Constraints = append(pg.Constraints, constraint...)
}

func (pg *PromptGenerator) AddResource(resource ...string) {
	pg.Resources = append(pg.Resources, resource...)
}

func (pg *PromptGenerator) AddCommand(command ...PGCommand) {
	pg.Commands = append(pg.Commands, command...)
}

func (pg *PromptGenerator) String() string {
	var res string

	if len(pg.Constraints) > 0 {
		res += "Constraints:\n"
		for _, constraint := range pg.Constraints {
			res += constraint + "\n"
		}
	}

	if len(pg.Commands) > 0 {
		res += "Commands:\n"
		for _, command := range pg.Commands {
			var args string
			for k, v := range command.Args {
				args += "\"" + k + "\": \"" + v + "\", "
			}
			res += command.Label + ": " + command.Name + ", args: " + args + "\n"
		}
	}

	if len(pg.Resources) > 0 {
		res += "Resources:\n"
		for _, resource := range pg.Resources {
			res += resource + "\n"
		}
	}

	res += "\nRespond in valid JSON format as described below\nResponse Format: \n" + pg.ResponseFormat + "\n"
	return res
}
