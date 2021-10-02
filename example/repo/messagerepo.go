package repo

import (
	"context"

	"github.com/halimath/depot"
	"github.com/halimath/depot/example/models"
)

func NewMessageRepo(factory *depot.Factory) *MessageRepo {
	return &MessageRepo{
		factory: factory,
	}
}

func (r *MessageRepo) FindByText(ctx context.Context, text string) ([]*models.Message, error) {
	return r.find(ctx, depot.Where("text", text))
}
