package schema

import (
	arceus "arceus/api"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Conversation struct {
	ent.Schema
}

func (Conversation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Base{},
	}
}

func (Conversation) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Optional(),
		field.JSON("context", arceus.Conversation{}),
	}
}

func (Conversation) Edges() []ent.Edge {
	return nil
}
