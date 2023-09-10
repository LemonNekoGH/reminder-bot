// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lemonnekogh/reminderbot/ent/settings"
)

// SettingsCreate is the builder for creating a Settings entity.
type SettingsCreate struct {
	config
	mutation *SettingsMutation
	hooks    []Hook
}

// SetChatID sets the "chat_id" field.
func (sc *SettingsCreate) SetChatID(i int64) *SettingsCreate {
	sc.mutation.SetChatID(i)
	return sc
}

// SetAllowAllUser sets the "allow_all_user" field.
func (sc *SettingsCreate) SetAllowAllUser(b bool) *SettingsCreate {
	sc.mutation.SetAllowAllUser(b)
	return sc
}

// Mutation returns the SettingsMutation object of the builder.
func (sc *SettingsCreate) Mutation() *SettingsMutation {
	return sc.mutation
}

// Save creates the Settings in the database.
func (sc *SettingsCreate) Save(ctx context.Context) (*Settings, error) {
	return withHooks(ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SettingsCreate) SaveX(ctx context.Context) *Settings {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SettingsCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SettingsCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SettingsCreate) check() error {
	if _, ok := sc.mutation.ChatID(); !ok {
		return &ValidationError{Name: "chat_id", err: errors.New(`ent: missing required field "Settings.chat_id"`)}
	}
	if _, ok := sc.mutation.AllowAllUser(); !ok {
		return &ValidationError{Name: "allow_all_user", err: errors.New(`ent: missing required field "Settings.allow_all_user"`)}
	}
	return nil
}

func (sc *SettingsCreate) sqlSave(ctx context.Context) (*Settings, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *SettingsCreate) createSpec() (*Settings, *sqlgraph.CreateSpec) {
	var (
		_node = &Settings{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(settings.Table, sqlgraph.NewFieldSpec(settings.FieldID, field.TypeInt))
	)
	if value, ok := sc.mutation.ChatID(); ok {
		_spec.SetField(settings.FieldChatID, field.TypeInt64, value)
		_node.ChatID = value
	}
	if value, ok := sc.mutation.AllowAllUser(); ok {
		_spec.SetField(settings.FieldAllowAllUser, field.TypeBool, value)
		_node.AllowAllUser = value
	}
	return _node, _spec
}

// SettingsCreateBulk is the builder for creating many Settings entities in bulk.
type SettingsCreateBulk struct {
	config
	builders []*SettingsCreate
}

// Save creates the Settings entities in the database.
func (scb *SettingsCreateBulk) Save(ctx context.Context) ([]*Settings, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Settings, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SettingsMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SettingsCreateBulk) SaveX(ctx context.Context) []*Settings {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SettingsCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SettingsCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}
