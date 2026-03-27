package todo

import (
	"context"

	"github.com/darkphotonKN/letsdo/internal/utils/errorutils"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, todo *Todo) error {
	query := `
		INSERT INTO todos (id, name, description, created_at, updated_at)
		VALUES (:id, :name, :description, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, todo, todo)
	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}
