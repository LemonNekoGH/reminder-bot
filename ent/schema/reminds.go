package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Reminds holds the schema definition for the Reminds entity.
type Reminds struct {
	ent.Schema
}

// Fields of the Reminds.
func (Reminds) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.Int64("chat_id"),             // 属于哪个聊天
		field.Int64("owner"),               // 属于哪个用户，其它用户无权修改此提醒
		field.String("period").Optional(),  // 间隔
		field.String("content").Optional(), // 内容
		field.String("name"),               // 名称
	}
}

// Edges of the Reminds.
func (Reminds) Edges() []ent.Edge {
	return nil
}
