package finify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	OPENAI_URL = "https://api.openai.com"

	OPENAI_ROLE_SYSTEM    = "system"
	OPENAI_ROLE_USER      = "user"
	OPENAI_ROLE_ASSISTANT = "assistant"

	OPENAI_GPT_3_5_TURBO      = "gpt-3.5-turbo"
	OPENAI_GPT_4              = "gpt-4"
	OPENAI_GPT_4_0314         = "gpt-4-0314"
	OPENAI_GPT_4_32K          = "gpt-4-32k"
	OPENAI_GPT_4_32K_0314     = "gpt-4-32k-0314"
	OPENAI_GPT_3_5_TURBO_0301 = "gpt-3.5-turbo-0301"
)

var (
	modelSet = map[string]interface{}{
		OPENAI_GPT_3_5_TURBO:      nil,
		OPENAI_GPT_4:              nil,
		OPENAI_GPT_4_0314:         nil,
		OPENAI_GPT_4_32K:          nil,
		OPENAI_GPT_4_32K_0314:     nil,
		OPENAI_GPT_3_5_TURBO_0301: nil,
	}
)

type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type OpenAIChatMessageCall struct {
	Model       string              `json:"model"`
	Messages    []OpenAIChatMessage `json:"messages"`
	Temperature int                 `json:"temperature,omitempty"`
	TopP        int                 `json:"top_p,omitempty"`
	Choices     int                 `json:"n,omitempty"`
	Stream      bool                `json:"stream,omitempty"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
	User        string              `json:"user,omitempty"`
}

type OpenAIChatMessageChoice struct {
	Index        int               `json:"index"`
	Message      OpenAIChatMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

type OpenAIChatMessageResponse struct {
	ID      string                    `json:"id"`
	Object  string                    `json:"object"`
	Created int64                     `json:"created"`
	Choices []OpenAIChatMessageChoice `json:"choices"`
}

func OpenAIChatCall(ctx context.Context, apiKey, model string, messages []OpenAIChatMessage) (*OpenAIChatMessageResponse, error) {
	if _, ok := modelSet[model]; !ok {
		return nil, fmt.Errorf("invalid model: %s", model)
	}
	res, err := OpenAICall(ctx, apiKey, OPENAI_URL+"/v1/chat/completions", OpenAIChatMessageCall{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}
	var data *OpenAIChatMessageResponse
	defer res.Close()
	if err := json.NewDecoder(res).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func OpenAICall(ctx context.Context, apiKey, location string, body interface{}) (io.ReadCloser, error) {
	var bodyReader io.Reader = nil
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, location, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp.Body, nil
}
