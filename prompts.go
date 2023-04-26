package finify

func GenerateSystemPrompt(goal string, commands []PGCommand) string {
	pg := NewPromptGenerator()
	pg.AddConstraint(
		"4000 word limit for short term memory, save important information to memory",
		"if you are unsure or want to recall past events, think about similar events ",
		"no user assitance",
		"only use the commands listed in double quotes eg \"command name\"",
		"use subprocess for commands that take a long time",
	)
	// TODO: need a better way of adding these in
	pg.AddCommand(
	// PGCommand{
	// 	Label: "Memorize a piece of information",
	// 	Name:  "memorize",
	// 	Args: map[string]string{
	// 		"data": "things to memorize",
	// 	},
	// },
	// PGCommand{
	// 	Label: "Recall a piece of information",
	// 	Name:  "recall",
	// 	Args: map[string]string{
	// 		"data": "things to recall",
	// 	},
	// },
	)
	pg.AddCommand(commands...)
	pg.AddResource(
		"Query news articles to gain information",
	)

	return `Your decisions must always be made independently without seeking user assitance. \n Use simple strategies with no legal consequences. Your name is Alex and you are a financial analyst assistant that will help generate reports and charts. \n GOALS: \n` + goal + "\n" + pg.String()
}
