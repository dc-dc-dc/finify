package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dc-dc-dc/finify"
)

const (
	newsApiEndpoint = "https://newsapi.org/v2/everything"
)

type QueryNewsCommand struct {
	newsApiKey string
}

type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"source"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Content     string `json:"content"`
	} `json:"articles"`
}

func NewQueryNewsCommand(newsApiKey string) finify.Command {
	return &QueryNewsCommand{
		newsApiKey: newsApiKey,
	}
}

func (qnCommand *QueryNewsCommand) GetPGCommand() finify.PGCommand {
	return finify.PGCommand{
		Label: "Query for news articles",
		Name:  "query_news",
		Args: map[string]string{
			"query": "query",
			"limit": "limit",
		},
	}
}

func (qnCommand *QueryNewsCommand) Call(ctx context.Context, agent *finify.Agent, args map[string]interface{}) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, newsApiEndpoint, nil)
	query, ok := args["query"].(string)
	if !ok {
		return "", err
	}
	urlQ := req.URL.Query()
	urlQ.Add("apiKey", qnCommand.newsApiKey)
	urlQ.Add("q", query)
	urlQ.Add("pageSize", "10")
	req.URL.RawQuery = urlQ.Encode()
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	defer res.Body.Close()
	var articleData *NewsAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&articleData); err != nil {
		return "", err
	}

	// format as Article 1: <title> <description> <content> \n
	var finalStr string

	for i, s := range articleData.Articles {
		finalStr += fmt.Sprintf("Article %d: %s %s %s\n", i+1, s.Title, s.Description, s.Content)
	}

	return finalStr, nil
}
