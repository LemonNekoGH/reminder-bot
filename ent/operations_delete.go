// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lemonnekogh/reminderbot/ent/operations"
	"github.com/lemonnekogh/reminderbot/ent/predicate"
)

// OperationsDelete is the builder for deleting a Operations entity.
type OperationsDelete struct {
	config
	hooks    []Hook
	mutation *OperationsMutation
}

// Where appends a list predicates to the OperationsDelete builder.
func (od *OperationsDelete) Where(ps ...predicate.Operations) *OperationsDelete {
	od.mutation.Where(ps...)
	return od
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (od *OperationsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, od.sqlExec, od.mutation, od.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (od *OperationsDelete) ExecX(ctx context.Context) int {
	n, err := od.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (od *OperationsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(operations.Table, sqlgraph.NewFieldSpec(operations.FieldID, field.TypeUUID))
	if ps := od.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, od.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	od.mutation.done = true
	return affected, err
}

// OperationsDeleteOne is the builder for deleting a single Operations entity.
type OperationsDeleteOne struct {
	od *OperationsDelete
}

// Where appends a list predicates to the OperationsDelete builder.
func (odo *OperationsDeleteOne) Where(ps ...predicate.Operations) *OperationsDeleteOne {
	odo.od.mutation.Where(ps...)
	return odo
}

// Exec executes the deletion query.
func (odo *OperationsDeleteOne) Exec(ctx context.Context) error {
	n, err := odo.od.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{operations.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (odo *OperationsDeleteOne) ExecX(ctx context.Context) {
	if err := odo.Exec(ctx); err != nil {
		panic(err)
	}
}
