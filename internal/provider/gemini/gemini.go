package gemini

import (
	arceus "arceus/api"
	"arceus/internal/provider"
	cfg "arceus/pkg/config"
	"context"
	"fmt"

	"github.com/pkoukk/tiktoken-go"
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
		availableModels: cfg.GetGemini().Models,
		client:          client,
	}
}

func (g *gemini) GenerateText(model string, messages []*arceus.Message) (*arceus.Message, *arceus.Usage, error) {

	result, usage, err := g.callApiGenerateText(model, messages)
	if err != nil {
		return nil, nil, err
	}

	return &arceus.Message{
		Content: result,
		Role:    arceus.Role_ROLE_BOT,
	}, usage, nil
}

func (g *gemini) GetAvailableModels() ([]string, error) {
	return g.availableModels, nil
}

func (g *gemini) callApiGenerateText(model string, messages []*arceus.Message) (string, *arceus.Usage, error) {

	if len(messages) == 0 {
		return "", nil, fmt.Errorf("no messages provided")
	}

	ctx := context.Background()

	history := []*genai.Content{}

	promptTokens := int32(0)

	for i := 0; i < len(messages)-1; i++ {
		switch messages[i].Role {
		case arceus.Role_ROLE_USER:
			history = append(history, genai.NewContentFromText(messages[i].Content, genai.RoleUser))
		case arceus.Role_ROLE_BOT:
			history = append(history, genai.NewContentFromText(messages[i].Content, genai.RoleModel))
		default:
			return "", nil, fmt.Errorf("invalid role %v", messages[i].Role)
		}
		tokens, _ := countTokens(messages[i].Content)
		promptTokens += tokens
	}

	chat, err := g.client.Chats.Create(ctx, model, nil, history)
	if err != nil {
		return "", nil, err
	}

	res, err := chat.SendMessage(ctx, genai.Part{Text: messages[len(messages)-1].Content})
	if err != nil {
		return "", nil, err
	}

	if len(res.Candidates) > 0 {
		mss := res.Candidates[0].Content.Parts[0].Text
		completionTokens, _ := countTokens(mss)

		return mss, &arceus.Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		}, nil
	}
	return "", nil, fmt.Errorf("no response from model")
}

func countTokens(text string) (int32, error) {
	enc, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
	if err != nil {
		return 0, fmt.Errorf("failed to load encoder: %v", err)
	}
	tokens := enc.Encode(text, nil, nil)
	return int32(len(tokens)), nil
}
