// Code generated by ent, DO NOT EDIT.

package operations

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the operations type in the database.
	Label = "operations"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldCompleted holds the string denoting the completed field in the database.
	FieldCompleted = "completed"
	// FieldSuccess holds the string denoting the success field in the database.
	FieldSuccess = "success"
	// FieldOperator holds the string denoting the operator field in the database.
	FieldOperator = "operator"
	// FieldMessageID holds the string denoting the message_id field in the database.
	FieldMessageID = "message_id"
	// FieldRemindID holds the string denoting the remind_id field in the database.
	FieldRemindID = "remind_id"
	// Table holds the table name of the operations in the database.
	Table = "operations"
)

// Columns holds all SQL columns for operations fields.
var Columns = []string{
	FieldID,
	FieldType,
	FieldCompleted,
	FieldSuccess,
	FieldOperator,
	FieldMessageID,
	FieldRemindID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCompleted holds the default value on creation for the "completed" field.
	DefaultCompleted bool
	// DefaultSuccess holds the default value on creation for the "success" field.
	DefaultSuccess bool
)

// OrderOption defines the ordering options for the Operations queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByCompleted orders the results by the completed field.
func ByCompleted(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCompleted, opts...).ToFunc()
}

// BySuccess orders the results by the success field.
func BySuccess(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSuccess, opts...).ToFunc()
}

// ByOperator orders the results by the operator field.
func ByOperator(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOperator, opts...).ToFunc()
}

// ByMessageID orders the results by the message_id field.
func ByMessageID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMessageID, opts...).ToFunc()
}

// ByRemindID orders the results by the remind_id field.
func ByRemindID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRemindID, opts...).ToFunc()
}
