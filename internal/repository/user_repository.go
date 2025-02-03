package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/models"
	"go.uber.org/zap"
)

type userRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Save(ctx context.Context, logger *zap.Logger, user models.User) error {
	query := `
		INSERT INTO tbl_user (username, password, role, is_deleted, created_at, create_by, update_at, update_by)
		VALUES ($1, $2, $3, $4, NOW(), $5, NOW(), $6)
		ON CONFLICT (username) DO UPDATE SET
			password = EXCLUDED.password,
			role = EXCLUDED.role,
			is_deleted = EXCLUDED.is_deleted,
			update_at = NOW(),
			update_by = EXCLUDED.update_by;
	`
	_, err := r.db.Exec(ctx, query, user.Username, user.Password, user.Role, user.IsDeleted, user.CreateBy, user.UpdateBy)
	if err != nil {
		logger.Error("Error saving user", zap.Error(err))
		return err
	}
	logger.Info("User saved successfully", zap.String("username", user.Username))
	return nil
}

func (r *userRepositoryImpl) FindById(ctx context.Context, logger *zap.Logger, username string) (*models.User, error) {
	query := `SELECT username, password, role, is_deleted, created_at, create_by, update_at, update_by FROM tbl_user WHERE username = $1`
	row := r.db.QueryRow(ctx, query, username)

	var user models.User
	err := row.Scan(&user.Username, &user.Password, &user.Role, &user.IsDeleted, &user.CreatedAt, &user.CreateBy, &user.UpdateAt, &user.UpdateBy)
	if err != nil {
		logger.Error("Error fetching user", zap.Error(err))
		return nil, err
	}
	logger.Info("User retrieved successfully", zap.String("username", user.Username))
	return &user, nil
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, logger *zap.Logger) ([]models.User, error) {
	query := `SELECT username, password, role, is_deleted, created_at, create_by, update_at, update_by FROM tbl_user`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Error("Error fetching all users", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Username, &user.Password, &user.Role, &user.IsDeleted, &user.CreatedAt, &user.CreateBy, &user.UpdateAt, &user.UpdateBy); err != nil {
			logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}
	logger.Info("All users retrieved successfully", zap.Int("count", len(users)))
	return users, nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, logger *zap.Logger, username string) error {
	query := `UPDATE tbl_user SET is_deleted = 'Y', update_at = NOW() WHERE username = $1`
	_, err := r.db.Exec(ctx, query, username)
	if err != nil {
		logger.Error("Error deleting user", zap.Error(err))
		return err
	}
	logger.Info("User deleted successfully", zap.String("username", username))
	return nil
}
