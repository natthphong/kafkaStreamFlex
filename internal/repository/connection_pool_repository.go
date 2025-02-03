package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/models"
	"go.uber.org/zap"
)

type connectionPoolRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewConnectionPoolRepository(db *pgxpool.Pool) ConnectionPoolRepository {
	return &connectionPoolRepositoryImpl{db: db}
}

func (r *connectionPoolRepositoryImpl) Save(ctx context.Context, logger *zap.Logger, connection models.ConnectionPool) error {
	query := `
		INSERT INTO tbl_connection_pool (id, name, type, end_point, username, password, key, is_deleted, created_at, create_by, update_at, update_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), $9, NOW(), $10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			type = EXCLUDED.type,
			end_point = EXCLUDED.end_point,
			username = EXCLUDED.username,
			password = EXCLUDED.password,
			key = EXCLUDED.key,
			is_deleted = EXCLUDED.is_deleted,
			update_at = NOW(),
			update_by = EXCLUDED.update_by;
	`
	_, err := r.db.Exec(ctx, query, connection.ID, connection.Name, connection.Type, connection.EndPoint, connection.Username, connection.Password, connection.Key, connection.IsDeleted, connection.CreateBy, connection.UpdateBy)
	if err != nil {
		logger.Error("Error saving connection pool", zap.Error(err))
		return err
	}
	logger.Info("Connection pool saved successfully", zap.Int("id", connection.ID))
	return nil
}

func (r *connectionPoolRepositoryImpl) FindById(ctx context.Context, logger *zap.Logger, id int) (*models.ConnectionPool, error) {
	query := `SELECT id, name, type, end_point, username, password, key, is_deleted, created_at, create_by, update_at, update_by FROM tbl_connection_pool WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var connection models.ConnectionPool
	err := row.Scan(&connection.ID, &connection.Name, &connection.Type, &connection.EndPoint, &connection.Username, &connection.Password, &connection.Key, &connection.IsDeleted, &connection.CreatedAt, &connection.CreateBy, &connection.UpdateAt, &connection.UpdateBy)
	if err != nil {
		logger.Error("Error fetching connection pool", zap.Error(err))
		return nil, err
	}
	logger.Info("Connection pool retrieved successfully", zap.Int("id", connection.ID))
	return &connection, nil
}

func (r *connectionPoolRepositoryImpl) FindAll(ctx context.Context, logger *zap.Logger) ([]models.ConnectionPool, error) {
	query := `SELECT id, name, type, end_point, username, password, key, is_deleted, created_at, create_by, update_at, update_by FROM tbl_connection_pool`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Error("Error fetching all connection pools", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var connections []models.ConnectionPool
	for rows.Next() {
		var connection models.ConnectionPool
		if err := rows.Scan(&connection.ID, &connection.Name, &connection.Type, &connection.EndPoint, &connection.Username, &connection.Password, &connection.Key, &connection.IsDeleted, &connection.CreatedAt, &connection.CreateBy, &connection.UpdateAt, &connection.UpdateBy); err != nil {
			logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		connections = append(connections, connection)
	}
	logger.Info("All connection pools retrieved successfully", zap.Int("count", len(connections)))
	return connections, nil
}

func (r *connectionPoolRepositoryImpl) Delete(ctx context.Context, logger *zap.Logger, id int) error {
	query := `UPDATE tbl_connection_pool SET is_deleted = 'Y', update_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.Error("Error deleting connection pool", zap.Error(err))
		return err
	}
	logger.Info("Connection pool deleted successfully", zap.Int("id", id))
	return nil
}
