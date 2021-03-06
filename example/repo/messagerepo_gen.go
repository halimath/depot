// This file has been generated by github.com/halimath/depot.
// Any changes will be overwritten when re-generating.

package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/halimath/depot"
	"github.com/halimath/depot/example/models"
)

var (
	messageRepoCols  = depot.Cols("id", "text", "order_index", "len", "attachment", "created", "updated")
	messageRepoTable = depot.Table("messages")
)

type MessageRepo struct {
	db *depot.DB
}

func (r *MessageRepo) Begin(ctx context.Context) (context.Context, error) {
	_, ctx, err := r.db.BeginTx(ctx)
	return ctx, err
}

func (r *MessageRepo) Commit(ctx context.Context) error {
	tx := depot.MustGetTx(ctx)
	return tx.Commit()
}

func (r *MessageRepo) Rollback(ctx context.Context) error {
	tx := depot.MustGetTx(ctx)
	return tx.Rollback()
}

func (r *MessageRepo) fromValues(vals depot.Values) (*models.Message, error) {
	var ok bool

	var id string

	id, ok = vals.GetString("id")

	if !ok {
		return nil, fmt.Errorf("failed to get id for models.Message: invalid value: %#v", vals["id"])
	}

	var text string

	text, ok = vals.GetString("text")

	if !ok {
		return nil, fmt.Errorf("failed to get text for models.Message: invalid value: %#v", vals["text"])
	}

	var orderindex int

	orderindex, ok = vals.GetInt("order_index")

	if !ok {
		return nil, fmt.Errorf("failed to get order_index for models.Message: invalid value: %#v", vals["order_index"])
	}

	var length float32

	length, ok = vals.GetFloat32("len")

	if !ok {
		return nil, fmt.Errorf("failed to get len for models.Message: invalid value: %#v", vals["len"])
	}

	var attachment []byte

	attachment, ok = vals.GetBytes("attachment")

	if !ok {
		return nil, fmt.Errorf("failed to get attachment for models.Message: invalid value: %#v", vals["attachment"])
	}

	var created time.Time

	created, ok = vals.GetTime("created")

	if !ok {
		return nil, fmt.Errorf("failed to get created for models.Message: invalid value: %#v", vals["created"])
	}

	var updated *time.Time

	if !vals.IsNull("updated") {
		if u, k := vals.GetTime("updated"); k {
			updated = &u
		} else {
			ok = false
		}
	}

	if !ok {
		return nil, fmt.Errorf("failed to get updated for models.Message: invalid value: %#v", vals["updated"])
	}

	return &models.Message{
		ID:         id,
		Text:       text,
		OrderIndex: orderindex,
		Length:     length,
		Attachment: attachment,
		Created:    created,
		Updated:    updated,
	}, nil
}

func (r *MessageRepo) find(ctx context.Context, clauses ...depot.SelectClause) ([]*models.Message, error) {
	tx := depot.MustGetTx(ctx)
	vals, err := tx.QueryMany(messageRepoCols, messageRepoTable, clauses...)
	if err != nil {
		err = fmt.Errorf("failed to load models.Message: %w", err)
		tx.Error(err)
		return nil, err
	}

	res := make([]*models.Message, 0, len(vals))
	for _, v := range vals {
		entity, err := r.fromValues(v)
		if err != nil {
			return nil, err
		}
		res = append(res, entity)
	}
	return res, nil
}

func (r *MessageRepo) count(ctx context.Context, clauses ...depot.WhereClause) (int, error) {
	tx := depot.MustGetTx(ctx)
	count, err := tx.QueryCount(messageRepoTable, clauses...)
	if err != nil {
		err = fmt.Errorf("failed to count models.Message: %w", err)
		tx.Error(err)
		return 0, err
	}

	return count, err
}

func (r *MessageRepo) LoadByID(ctx context.Context, ID string) (*models.Message, error) {
	tx := depot.MustGetTx(ctx)
	vals, err := tx.QueryOne(messageRepoCols, messageRepoTable, depot.Where(depot.Eq("id", ID)))
	if err != nil {
		err = fmt.Errorf("failed to load models.Message by ID: %w", err)
		if !errors.Is(err, depot.ErrNoResult) {
			tx.Error(err)
		}
		return nil, err
	}
	return r.fromValues(vals)
}

func (r *MessageRepo) toValues(entity *models.Message) depot.Values {
	return depot.Values{
		"id":          entity.ID,
		"text":        entity.Text,
		"order_index": entity.OrderIndex,
		"len":         entity.Length,
		"attachment":  entity.Attachment,
		"created":     entity.Created,
		"updated":     entity.Updated,
	}
}

func (r *MessageRepo) Insert(ctx context.Context, entity *models.Message) error {
	tx := depot.MustGetTx(ctx)
	err := tx.InsertOne(messageRepoTable, r.toValues(entity))
	if err != nil {
		err = fmt.Errorf("failed to insert models.Message: %w", err)
	}
	return err
}

func (r *MessageRepo) delete(ctx context.Context, clauses ...depot.WhereClause) error {
	tx := depot.MustGetTx(ctx)
	err := tx.DeleteMany(messageRepoTable, clauses...)
	if err != nil {
		err = fmt.Errorf("failed to delete models.Message: %w", err)
	}
	return err
}

func (r *MessageRepo) Update(ctx context.Context, entity *models.Message) error {
	tx := depot.MustGetTx(ctx)
	err := tx.UpdateMany(messageRepoTable, r.toValues(entity), depot.Where(depot.Eq("id", entity.ID)))
	if err != nil {
		err = fmt.Errorf("failed to update models.Message: %w", err)
	}
	return err
}

func (r *MessageRepo) DeleteByID(ctx context.Context, ID string) error {
	return r.delete(ctx, depot.Where(depot.Eq("id", ID)))
}

func (r *MessageRepo) Delete(ctx context.Context, entity *models.Message) error {
	return r.delete(ctx, depot.Where(depot.Eq("id", entity.ID)))
}
