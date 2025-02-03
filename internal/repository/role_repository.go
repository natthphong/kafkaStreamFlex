package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/kafkaStreamFlex/internal/models"
	"go.uber.org/zap"
)

type roleRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &roleRepositoryImpl{db: db}
}

func (r *roleRepositoryImpl) Save(ctx context.Context, logger *zap.Logger, role models.Role) error {
	query := `
		INSERT INTO tbl_role (role_name, is_deleted, created_at, create_by, update_at, update_by)
		VALUES ($1, $2, NOW(), $3, NOW(), $4)
		ON CONFLICT (role_name) DO UPDATE SET
			is_deleted = EXCLUDED.is_deleted,
			update_at = NOW(),
			update_by = EXCLUDED.update_by;
	`
	_, err := r.db.Exec(ctx, query, role.RoleName, role.IsDeleted, role.CreateBy, role.UpdateBy)
	if err != nil {
		logger.Error("Error saving role", zap.Error(err))
		return err
	}
	logger.Info("Role saved successfully", zap.String("role_name", role.RoleName))
	return nil
}

func (r *roleRepositoryImpl) FindById(ctx context.Context, logger *zap.Logger, roleName string) (*models.Role, error) {
	query := `SELECT role_name, is_deleted, created_at, create_by, update_at, update_by FROM tbl_role WHERE role_name = $1`
	row := r.db.QueryRow(ctx, query, roleName)

	var role models.Role
	err := row.Scan(&role.RoleName, &role.IsDeleted, &role.CreatedAt, &role.CreateBy, &role.UpdateAt, &role.UpdateBy)
	if err != nil {
		logger.Error("Error fetching role", zap.Error(err))
		return nil, err
	}
	logger.Info("Role retrieved successfully", zap.String("role_name", role.RoleName))
	return &role, nil
}

func (r *roleRepositoryImpl) FindAll(ctx context.Context, logger *zap.Logger) ([]models.Role, error) {
	query := `SELECT role_name, is_deleted, created_at, create_by, update_at, update_by FROM tbl_role`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Error("Error fetching all roles", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.RoleName, &role.IsDeleted, &role.CreatedAt, &role.CreateBy, &role.UpdateAt, &role.UpdateBy); err != nil {
			logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}
		roles = append(roles, role)
	}
	logger.Info("All roles retrieved successfully", zap.Int("count", len(roles)))
	return roles, nil
}

func (r *roleRepositoryImpl) Delete(ctx context.Context, logger *zap.Logger, roleName string) error {
	query := `UPDATE tbl_role SET is_deleted = 'Y', update_at = NOW() WHERE role_name = $1`
	_, err := r.db.Exec(ctx, query, roleName)
	if err != nil {
		logger.Error("Error deleting role", zap.Error(err))
		return err
	}
	logger.Info("Role deleted successfully", zap.String("role_name", roleName))
	return nil
}
