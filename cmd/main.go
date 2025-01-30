package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/config"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/api"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/connection/cache"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/connection/db"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/connection/httputil"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/connection/logz"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/connection/s3"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
)

func main() {
	currentTime := time.Now()
	versionDeploy := currentTime.Unix()
	ctx := context.Background()
	app := initFiber()
	config.InitTimeZone()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(errors.New("Unable to initial config."))
	}
	logz.Init(cfg.LogConfig.Level, cfg.Server.Name)
	defer logz.Drop()

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	logger := zap.L()
	logger.Info("version " + strconv.FormatInt(versionDeploy, 10))
	dbPool, err := db.Open(ctx, cfg.DBConfig)
	if err != nil {
		logger.Fatal("server connect to db", zap.Error(err))
	}
	defer dbPool.Close()
	logger.Info("DB CONNECT")

	httpClient := httputil.InitHttpClient(
		cfg.HTTP.TimeOut,
		cfg.HTTP.MaxIdleConn,
		cfg.HTTP.MaxIdleConnPerHost,
		cfg.HTTP.MaxConnPerHost,
	)
	_ = httpClient

	redisClient, err := cache.Initialize(ctx, cfg.RedisConfig)
	if err != nil {
		logger.Fatal("server connect to redis", zap.Error(err))
	}
	redisCMD := redisClient.UniversalClient()
	defer func() {
		err = redisCMD.Close()
		if err != nil {
			logger.Fatal("closing redis connection error", zap.Error(err))
		}
	}()
	logger.Info("Redis Connected")

	//configSftp := sftp.Config{
	//	Username: "",
	//	Password: "",
	//	Server:   "host:port",
	//	Timeout:  time.Second * 30,
	//}
	//
	//client, err := sftp.New(configSftp)
	//if err != nil {
	//	logger.Fatal("server connect to sftp", zap.Error(err))
	//}
	//defer client.Close()
	s3Config := cfg.AwsS3Config
	awsClient, err := s3.CreateSessionAws(&s3Config.DoSpaceEndpoint, s3Config.AccessKey, s3Config.SecretKey, s3Config.DoSpaceRegion)
	if err != nil {
		logger.Fatal("server connect to s3", zap.Error(err))
	}
	_ = awsClient
	logger.Info("S3 Connected")

	group := app.Group(fmt.Sprintf("/%s/api/v1", cfg.Server.Name))

	group.Get("/health", func(c *fiber.Ctx) error {
		return api.Ok(c, versionDeploy)
	})
	logger.Info(fmt.Sprintf("/%s/api/v1", cfg.Server.Name), zap.Any("port", cfg.Server.Port))
	if err = app.Listen(fmt.Sprintf(":%v", cfg.Server.Port)); err != nil {
		logger.Fatal(err.Error())
	}

}

func initFiber() *fiber.App {
	app := fiber.New(
		fiber.Config{
			ReadTimeout:           5 * time.Second,
			WriteTimeout:          5 * time.Second,
			IdleTimeout:           30 * time.Second,
			DisableStartupMessage: true,
			CaseSensitive:         true,
			StrictRouting:         true,
		},
	)
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(SetHeaderID())
	return app
}

func SetHeaderID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		randomTrace := uuid.New().String()
		traceId := c.Get("traceId")
		//refId := c.Get("RequestRef")
		if traceId == "" {
			traceId = randomTrace
		}

		c.Accepts(fiber.MIMEApplicationJSON)
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		c.Request().Header.Set("traceId", traceId)
		return c.Next()
	}
}
