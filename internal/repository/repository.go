package repository

import (
	"arceus/internal/repository/conversation"
	"arceus/pkg/ent"
)

type Repository struct {
	Ent          *ent.Client
	Conversation conversation.Conversation
}

func New(ent *ent.Client) *Repository {
	conversation := conversation.New(ent)
	return &Repository{
		Ent:          ent,
		Conversation: conversation,
	}
}
