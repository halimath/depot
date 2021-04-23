package repo

import (
	"context"

	"github.com/halimath/depot"
	"github.com/halimath/depot/example/models"
)

func (r *MessageRepo) FindByText(ctx context.Context, text string) ([]*models.Message, error) {
	return r.find(ctx, depot.Where("text", text))
}
