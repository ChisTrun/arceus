package gemini

import (
	arceus "arceus/api"
	"arceus/internal/provider"
	cfg "arceus/pkg/config"
	"context"
	"fmt"

	"google.golang.org/genai"
)

type gemini struct {
	availableModels []string
	client          *genai.Client
}

func New(cfg *cfg.Config) provider.Provider {

	if cfg.GetGemini() == nil || !cfg.GetGemini().Enable {
		return &provider.Noop{}
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  cfg.GetGemini().ApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return &provider.Noop{}
	}

	return &gemini{
		availableModels: cfg.GetMistral().Models,
		client:          client,
	}
}

func (g *gemini) GenerateText(model string, messages []*arceus.Message) (*arceus.Message, error) {

	result, err := g.callApiGenerateText(model, messages)
	if err != nil {
		return nil, err
	}

	return &arceus.Message{
		Content: result,
		Role:    arceus.Role_ROLE_BOT,
	}, nil
}

func (g *gemini) GetAvailableModels() ([]string, error) {
	return g.availableModels, nil
}

func (g *gemini) callApiGenerateText(model string, messages []*arceus.Message) (string, error) {

	if len(messages) == 0 {
		return "", fmt.Errorf("no messages provided")
	}

	ctx := context.Background()

	history := []*genai.Content{}

	for i := 0; i < len(messages)-1; i++ {
		switch messages[i].Role {
		case arceus.Role_ROLE_USER:
			history = append(history, genai.NewContentFromText(messages[i].Content, genai.RoleUser))
		case arceus.Role_ROLE_BOT:
			history = append(history, genai.NewContentFromText(messages[i].Content, genai.RoleModel))
		default:
			return "", fmt.Errorf("invalid role %v", messages[i].Role)
		}
	}

	chat, err := g.client.Chats.Create(ctx, model, nil, history)
	if err != nil {
		return "", err
	}

	res, _ := chat.SendMessage(ctx, genai.Part{Text: messages[len(messages)-1].Content})
	if err != nil {
		return "", err
	}

	if len(res.Candidates) > 0 {
		return res.Candidates[0].Content.Parts[0].Text, nil
	}
	return "", fmt.Errorf("no response from model")
}
