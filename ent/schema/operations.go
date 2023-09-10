package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type OperationType int

func (o OperationType) String() string {
	switch o {
	case EnumOperationTypeCreate:
		return "创建"
	case EnumOperationTypeSetContent:
		return "设置内容"
	case EnumOperationTypeSetName:
		return "设置名称"
	case EnumOperationTypeSetCron:
		return "设置定时"
	default:
		return "未知"
	}
}

const (
	EnumOperationTypeCreate OperationType = iota
	EnumOperationTypeSetName
	EnumOperationTypeSetContent
	EnumOperationTypeSetCron
)

// Operations holds the schema definition for the Operations entity.
type Operations struct {
	ent.Schema
}

// Fields of the Operations.
func (Operations) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.Int("type"),                              // 操作类型
		field.Bool("completed").Default(false),         // 是否已完成操作
		field.Bool("success").Default(false),           // 是否成功操作
		field.Int64("operator"),                        // 操作者
		field.Int("message_id"),                        // 需要回复的消息 ID
		field.UUID("remind_id", uuid.New()).Optional(), // 被操作的提醒项 ID
	}
}

// Edges of the Operations.
func (Operations) Edges() []ent.Edge {
	return nil
}
