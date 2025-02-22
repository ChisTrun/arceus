package openai

import (
	arceus "arceus/api"
	"arceus/internal/provider"
	cfg "arceus/pkg/config"
)

type openai struct {
}

func New(cfg *cfg.Config) provider.Provider {
	if cfg.GetMistral() == nil || !cfg.GetMistral().Enable {
		return &Noop{}
	}

	return &openai{}
}

func (o *openai) GenerateText(model string, messages []*arceus.Message) (*arceus.Message, error) {
	return nil, nil
}

func (o *openai) GetAvailableModels() ([]string, error) {
	return nil, nil
}
