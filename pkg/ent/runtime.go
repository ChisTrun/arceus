// Code generated by ent, DO NOT EDIT.

package ent

import (
	"arceus/pkg/ent/conversation"
	"arceus/schema"
	"time"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	conversationMixin := schema.Conversation{}.Mixin()
	conversationMixinFields0 := conversationMixin[0].Fields()
	_ = conversationMixinFields0
	conversationFields := schema.Conversation{}.Fields()
	_ = conversationFields
	// conversationDescCreatedAt is the schema descriptor for created_at field.
	conversationDescCreatedAt := conversationMixinFields0[1].Descriptor()
	// conversation.DefaultCreatedAt holds the default value on creation for the created_at field.
	conversation.DefaultCreatedAt = conversationDescCreatedAt.Default.(func() time.Time)
	// conversationDescUpdatedAt is the schema descriptor for updated_at field.
	conversationDescUpdatedAt := conversationMixinFields0[2].Descriptor()
	// conversation.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	conversation.DefaultUpdatedAt = conversationDescUpdatedAt.Default.(func() time.Time)
	// conversation.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	conversation.UpdateDefaultUpdatedAt = conversationDescUpdatedAt.UpdateDefault.(func() time.Time)
}
