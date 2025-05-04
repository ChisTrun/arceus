package llm

import (
	arceus "arceus/api"
	"arceus/internal/provider"
	"arceus/internal/repository"
	"arceus/internal/utils/tx"
	"arceus/pkg/ent"
	"arceus/pkg/logger/pkg/logging"
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Llm interface {
	GenerateText(ctx context.Context, req *arceus.GenerateTextRequest) (*arceus.GenerateTextResponse, error)
}

type llm struct {
	modelMap map[string]provider.Provider
	repo     *repository.Repository
}

func New(aiProviders []provider.Provider, repo *repository.Repository) Llm {
	modelMap := make(map[string]provider.Provider)
	for _, aiProvider := range aiProviders {
		models, err := aiProvider.GetAvailableModels()
		if err != nil {
			continue
		}
		for _, model := range models {
			modelMap[model] = aiProvider
		}
	}

	return &llm{
		modelMap: modelMap,
		repo:     repo,
	}
}

func (l *llm) GenerateText(ctx context.Context, req *arceus.GenerateTextRequest) (*arceus.GenerateTextResponse, error) {
	aiProvider, ok := l.modelMap[req.Model]
	if !ok {
		return nil, fmt.Errorf("model %s not found", req.Model)
	}

	var conversation *ent.Conversation = nil
	var err error = nil
	if req.ConversationId != nil {
		if conversation, err = l.repo.Conversation.Get(ctx, *req.ConversationId); err != nil {
			return nil, err
		}
	} else {
		if conversation, err = l.repo.Conversation.Create(ctx, arceus.Conversation{
			Messages: []*arceus.Message{&arceus.Message{
				Content: req.Content,
				Role:    arceus.Role_ROLE_USER,
			}},
		}); err != nil {
			return nil, err
		}
	}

	conversation.Context.Messages = append(conversation.Context.Messages, &arceus.Message{
		Content: req.Content,
		Role:    arceus.Role_ROLE_USER,
	})

	mss, usage, err := aiProvider.GenerateText(req.Model, conversation.Context.Messages)
	if err != nil {
		return nil, err
	}

	conversation.Context.Messages = append(conversation.Context.Messages, mss)
	go func(conversationContext arceus.Conversation) {
		if txErr := tx.WithTransaction(context.Background(), l.repo.Ent, func(ctx context.Context, tx tx.Tx) error {
			return l.repo.Conversation.Update(ctx, tx, conversation.ID, conversationContext)
		}); txErr != nil {
			logging.Logger(context.Background()).Error("failed to update conversation", zap.Error(txErr))
		}
	}(conversation.Context)

	return &arceus.GenerateTextResponse{
		Content:        mss.Content,
		ConversationId: conversation.ID,
		CreatedAt:      timestamppb.Now(),
		Usage:          usage,
	}, nil
}
