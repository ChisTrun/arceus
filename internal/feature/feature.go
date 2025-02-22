package feature

import (
	"arceus/internal/feature/llm"
	"arceus/internal/provider"
	"arceus/internal/repository"
)

type Feature struct {
	Llm llm.Llm
}

func New(repo *repository.Repository, aiProviders []provider.Provider) *Feature {
	llm := llm.New(aiProviders, repo)
	return &Feature{
		Llm: llm,
	}
}
