package openai

import (
	arceus "arceus/api"
	"arceus/internal/provider"
	cfg "arceus/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type openai struct {
	availableModels []string
	endpoint        string
	apiKey          string
}

func New(cfg *cfg.Config) provider.Provider {
	if cfg.GetOpenai() == nil || !cfg.GetOpenai().Enable {
		return &provider.Noop{}
	}

	return &openai{
		availableModels: cfg.Openai.GetModels(),
		endpoint:        cfg.Openai.GetEndpoint(),
		apiKey:          cfg.Openai.GetApiKey(),
	}
}

func (o *openai) GenerateText(model string, messages []*arceus.Message) (*arceus.Message, *arceus.Usage, error) {
	result, usage, err := o.callApiGenerateText(model, messages)
	if err != nil {
		return nil, nil, err
	}

	return &arceus.Message{
		Content: result,
		Role:    arceus.Role_ROLE_BOT,
	}, usage, nil
}

func (o *openai) GetAvailableModels() ([]string, error) {
	return o.availableModels, nil
}

func (o *openai) callApiGenerateText(model string, messages []*arceus.Message) (string, *arceus.Usage, error) {
	url := fmt.Sprintf("%v/%v", o.endpoint, "v1/chat/completions")

	mistralMss := []Message{}
	for _, message := range messages {
		tmp := Message{
			Content: message.Content,
		}
		switch message.Role {
		case arceus.Role_ROLE_USER:
			tmp.Role = "user"
		case arceus.Role_ROLE_BOT:
			tmp.Role = "assistant"
		default:
			return "", nil, fmt.Errorf("invalid role %v", message.Role)
		}

		mistralMss = append(mistralMss, tmp)
	}

	payload := RequestPayload{
		Model:    model,
		Messages: mistralMss,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", nil, fmt.Errorf("error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("HTTP error! Status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("error reading response body: %v", err)
	}

	response, err := ParseChatResponse(body)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing response: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", nil, fmt.Errorf("empty response")
	}

	return response.Choices[0].Message.Content, &arceus.Usage{
		PromptTokens:     int32(response.Usage.PromptTokens),
		TotalTokens:      int32(response.Usage.TotalTokens),
		CompletionTokens: int32(response.Usage.CompletionTokens),
	}, nil
}
