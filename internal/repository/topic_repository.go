package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/models"
	"go.uber.org/zap"
)

type topicRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewTopicRepository(db *pgxpool.Pool) TopicRepository {
	return &topicRepositoryImpl{db: db}
}
func (r *topicRepositoryImpl) Save(ctx context.Context, logger *zap.Logger, topic models.Topic) error {
	query := `
		INSERT INTO tbl_topic (topic, pod_run, script_version, script_name, in_active, created_at, create_by, update_at, update_by)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6, NOW(), $7)
		ON CONFLICT (topic) DO UPDATE SET
			pod_run = EXCLUDED.pod_run,
			script_version = EXCLUDED.script_version,
			script_name = EXCLUDED.script_name,
			in_active = EXCLUDED.in_active,
			update_at = NOW(),
			update_by = EXCLUDED.update_by;
	`
	_, err := r.db.Exec(ctx, query, topic.Topic, topic.PodRun, topic.ScriptVersion, topic.ScriptName, topic.InActive, topic.CreateBy, topic.UpdateBy)
	if err != nil {
		logger.Error("Error saving topic", zap.Error(err))
		return err
	}
	logger.Info("Topic saved successfully", zap.String("topic", topic.Topic))
	return nil
}
func (r *topicRepositoryImpl) FindById(ctx context.Context, logger *zap.Logger, topicName string) (*models.Topic, error) {
	query := `SELECT topic, pod_run, script_version, script_name, in_active, created_at, create_by, update_at, update_by FROM tbl_topic WHERE topic = $1`
	row := r.db.QueryRow(ctx, query, topicName)

	var topic models.Topic
	err := row.Scan(&topic.Topic, &topic.PodRun, &topic.ScriptVersion, &topic.ScriptName, &topic.InActive, &topic.CreatedAt, &topic.CreateBy, &topic.UpdateAt, &topic.UpdateBy)
	if err != nil {
		logger.Error("Error fetching topic", zap.Error(err))
		return nil, err
	}
	logger.Info("Topic retrieved successfully", zap.String("topic", topic.Topic))
	return &topic, nil
}

func (r *topicRepositoryImpl) FindAll(ctx context.Context, logger *zap.Logger) ([]models.Topic, error) {
	query := `SELECT topic, pod_run, script_version, script_name, in_active, created_at, create_by, update_at, update_by FROM tbl_topic`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Error("Error fetching all topics", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var topic models.Topic
		if err := rows.Scan(&topic.Topic, &topic.PodRun, &topic.ScriptVersion, &topic.ScriptName, &topic.InActive, &topic.CreatedAt, &topic.CreateBy, &topic.UpdateAt, &topic.UpdateBy); err != nil {
			logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		topics = append(topics, topic)
	}
	logger.Info("All topics retrieved successfully", zap.Int("count", len(topics)))
	return topics, nil
}

func (r *topicRepositoryImpl) Delete(ctx context.Context, logger *zap.Logger, topicName string) error {
	query := `UPDATE tbl_topic SET in_active = 'Y', update_at = NOW() WHERE topic = $1`
	_, err := r.db.Exec(ctx, query, topicName)
	if err != nil {
		logger.Error("Error deleting topic", zap.Error(err))
		return err
	}
	logger.Info("Topic deleted successfully", zap.String("topic", topicName))
	return nil
}
