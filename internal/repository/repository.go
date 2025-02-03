package repository

import (
	"context"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/models"
	"go.uber.org/zap"
)

type TopicRepository interface {
	Save(ctx context.Context, logger *zap.Logger, topic models.Topic) error
	FindById(ctx context.Context, logger *zap.Logger, topicName string) (*models.Topic, error)
	FindAll(ctx context.Context, logger *zap.Logger) ([]models.Topic, error)
	Delete(ctx context.Context, logger *zap.Logger, topicName string) error
}
