package conversation

import (
	arceus "arceus/api"
	"arceus/internal/utils/tx"
	"arceus/pkg/ent"
	entconversation "arceus/pkg/ent/conversation"
	"context"
)

type Conversation interface {
	Create(ctx context.Context, conversation arceus.Conversation) (*ent.Conversation, error)
	Get(ctx context.Context, id uint64) (*ent.Conversation, error)
	Update(ctx context.Context, tx tx.Tx, id uint64, conversation arceus.Conversation) error
}

type conversation struct {
	ent *ent.Client
}

func New(ent *ent.Client) Conversation {
	return &conversation{
		ent: ent,
	}
}

func (c *conversation) Create(ctx context.Context, conversation arceus.Conversation) (*ent.Conversation, error) {
	return c.ent.Conversation.Create().SetContext(conversation).Save(ctx)
}

func (c *conversation) Get(ctx context.Context, id uint64) (*ent.Conversation, error) {
	return c.ent.Conversation.Query().Where(entconversation.ID(id)).Only(ctx)
}

func (c *conversation) Update(ctx context.Context, tx tx.Tx, id uint64, conversation arceus.Conversation) error {
	return tx.Client().Conversation.UpdateOneID(id).SetContext(conversation).Exec(ctx)
}
