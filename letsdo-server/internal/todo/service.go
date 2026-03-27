package todo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type service struct {
	repo   Repository
	logger *slog.Logger
}

type Repository interface {
	Create(ctx context.Context, todo *Todo) error
}

func NewService(repo Repository, logger *slog.Logger) *service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) CreateTodo(ctx context.Context, req CreateTodoRequest) (*Todo, error) {
	todo := &Todo{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		s.logger.Error("failed to create todo", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil
}
