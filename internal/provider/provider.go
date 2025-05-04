package provider

import (
	arceus "arceus/api"
)

type Provider interface {
	GetAvailableModels() ([]string, error)
	GenerateText(model string, messages []*arceus.Message) (*arceus.Message, *arceus.Usage, error)
}
