package repository

import (
	"context"
	"github.com/natthphong/kafkaStreamFlex/internal/models"
	"go.uber.org/zap"
)

type TopicRepository interface {
	Save(ctx context.Context, logger *zap.Logger, topic models.Topic) error
	FindById(ctx context.Context, logger *zap.Logger, topicName string) (*models.Topic, error)
	FindAll(ctx context.Context, logger *zap.Logger) ([]models.Topic, error)
	Delete(ctx context.Context, logger *zap.Logger, topicName string) error
}

type ScriptRepository interface {
	Save(ctx context.Context, logger *zap.Logger, script models.Script) error
	FindById(ctx context.Context, logger *zap.Logger, id int) (*models.Script, error)
	FindAll(ctx context.Context, logger *zap.Logger) ([]models.Script, error)
	Delete(ctx context.Context, logger *zap.Logger, id int) error
}

type ConnectionPoolRepository interface {
	Save(ctx context.Context, logger *zap.Logger, connection models.ConnectionPool) error
	FindById(ctx context.Context, logger *zap.Logger, id int) (*models.ConnectionPool, error)
	FindAll(ctx context.Context, logger *zap.Logger) ([]models.ConnectionPool, error)
	Delete(ctx context.Context, logger *zap.Logger, id int) error
}

type UserRepository interface {
	Save(ctx context.Context, logger *zap.Logger, user models.User) error
	FindById(ctx context.Context, logger *zap.Logger, username string) (*models.User, error)
	FindAll(ctx context.Context, logger *zap.Logger) ([]models.User, error)
	Delete(ctx context.Context, logger *zap.Logger, username string) error
}

type RoleRepository interface {
	Save(ctx context.Context, logger *zap.Logger, role models.Role) error
	FindById(ctx context.Context, logger *zap.Logger, roleName string) (*models.Role, error)
	FindAll(ctx context.Context, logger *zap.Logger) ([]models.Role, error)
	Delete(ctx context.Context, logger *zap.Logger, roleName string) error
}
