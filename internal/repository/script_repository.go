package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/kafkaStreamFlex/internal/models"
	"go.uber.org/zap"
)

type scriptRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewScriptRepository(db *pgxpool.Pool) ScriptRepository {
	return &scriptRepositoryImpl{db: db}
}

func (r *scriptRepositoryImpl) Save(ctx context.Context, logger *zap.Logger, script models.Script) error {
	query := `
		INSERT INTO tbl_script (id, script_key, script_name, version, is_deleted, created_at, create_by)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6)
		ON CONFLICT (id) DO UPDATE SET
			script_key = EXCLUDED.script_key,
			script_name = EXCLUDED.script_name,
			version = EXCLUDED.version,
			is_deleted = EXCLUDED.is_deleted,
			created_at = EXCLUDED.created_at,
			create_by = EXCLUDED.create_by;
	`
	_, err := r.db.Exec(ctx, query, script.ID, script.ScriptKey, script.ScriptName, script.Version, script.IsDeleted, script.CreateBy)
	if err != nil {
		logger.Error("Error saving script", zap.Error(err))
		return err
	}
	logger.Info("Script saved successfully", zap.Int("id", script.ID))
	return nil
}

func (r *scriptRepositoryImpl) FindById(ctx context.Context, logger *zap.Logger, id int) (*models.Script, error) {
	query := `SELECT id, script_key, script_name, version, is_deleted, created_at, create_by FROM tbl_script WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var script models.Script
	err := row.Scan(&script.ID, &script.ScriptKey, &script.ScriptName, &script.Version, &script.IsDeleted, &script.CreatedAt, &script.CreateBy)
	if err != nil {
		logger.Error("Error fetching script", zap.Error(err))
		return nil, err
	}
	logger.Info("Script retrieved successfully", zap.Int("id", script.ID))
	return &script, nil
}

func (r *scriptRepositoryImpl) FindAll(ctx context.Context, logger *zap.Logger) ([]models.Script, error) {
	query := `SELECT id, script_key, script_name, version, is_deleted, created_at, create_by FROM tbl_script`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Error("Error fetching all scripts", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var scripts []models.Script
	for rows.Next() {
		var script models.Script
		if err := rows.Scan(&script.ID, &script.ScriptKey, &script.ScriptName, &script.Version, &script.IsDeleted, &script.CreatedAt, &script.CreateBy); err != nil {
			logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		scripts = append(scripts, script)
	}
	logger.Info("All scripts retrieved successfully", zap.Int("count", len(scripts)))
	return scripts, nil
}

func (r *scriptRepositoryImpl) Delete(ctx context.Context, logger *zap.Logger, id int) error {
	query := `UPDATE tbl_script SET is_deleted = 'Y' WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.Error("Error deleting script", zap.Error(err))
		return err
	}
	logger.Info("Script deleted successfully", zap.Int("id", id))
	return nil
}
