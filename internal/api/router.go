package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/kafkaStreamFlex/config"
	"github.com/natthphong/kafkaStreamFlex/internal/connection/s3"
	"github.com/natthphong/kafkaStreamFlex/internal/repository"
)

func RegisterRoutes(router fiber.Router,
	dbPool *pgxpool.Pool,
	jwtSecret string,
	cfg config.Config,
	awsClient *session.Session,

) {
	scriptRouter := router.Group("/script")
	s3UploadFunc := s3.NewS3Upload(awsClient, cfg.AwsS3Config.BucketName)
	//m := middleware.JWTMiddlewareWithObjects(jwtSecret, userObjectPermission)
	scriptRouter.Post("/upload", NewUploadScript(
		cfg.ScriptConfig,
		repository.NewScriptRepository(dbPool),
		s3UploadFunc,
	))

	//TODO  UPDATE INACTIVE FOR USER REMOVE ACT
}
