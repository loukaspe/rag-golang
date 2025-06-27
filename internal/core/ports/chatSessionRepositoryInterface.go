package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/loukaspe/rag-golang/internal/core/domain"
)

type ChatSessionRepositoryInterface interface {
	GetChatSession(context.Context, uuid.UUID) (*domain.ChatSession, error)
	GetUserChatSessions(context.Context, uuid.UUID) ([]*domain.ChatSession, error)
	CreateChatSession(context.Context, *domain.ChatSession) (uuid.UUID, error)
	UpdateChatSessionTitle(context.Context, uuid.UUID, string) error
}
