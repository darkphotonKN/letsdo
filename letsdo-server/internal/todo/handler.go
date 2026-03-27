package todo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/darkphotonKN/letsdo/internal/utils/errorutils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
	logger  *slog.Logger
}

type Service interface {
	CreateTodo(ctx context.Context, req CreateTodoRequest) (*Todo, error)
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateTodo(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", slog.String("error", err.Error()))
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	todo, err := h.service.CreateTodo(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, errorutils.ErrDuplicateResource):
			c.JSON(409, gin.H{"error": "Todo already exists"})
		case errors.Is(err, errorutils.ErrInvalidInput):
			c.JSON(400, gin.H{"error": err.Error()})
		case errors.Is(err, errorutils.ErrConstraintViolation):
			c.JSON(400, gin.H{"error": "Invalid data provided"})
		default:
			h.logger.Error("failed to create todo", slog.String("error", err.Error()))
			c.JSON(500, gin.H{"error": "Failed to create todo"})
		}
		return
	}

	c.JSON(201, todo)
}
