package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/loukaspe/rag-golang/internal/core/domain"
)

type MessageRepositoryInterface interface {
	CreateMessage(context.Context, *domain.Message) (uuid.UUID, error)
	GetMessage(context.Context, uuid.UUID) (*domain.Message, error)
	UpdateMessageFeedback(context.Context, uuid.UUID, string) error
}
