package sdk

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/connection/cache"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/connection/kafka"
)

type StreamFlexClient struct {
	Ctx           context.Context
	DB            *pgxpool.Pool
	RedisClient   *cache.RedisClient
	HTTPClient    *http.Client
	S3Client      *session.Session
	KafkaProducer *kafka.SendMessageSyncFunc // or an interface type
	Payload       []byte                     // Raw message payload from Kafka
}

func NewStreamFlexClient(
	ctx context.Context,
	db *pgxpool.Pool,
	redisClient *cache.RedisClient,
	httpClient *http.Client,
	s3Client *session.Session,
	kafkaProducer *kafka.SendMessageSyncFunc,
	payload []byte,
) *StreamFlexClient {
	return &StreamFlexClient{
		Ctx:           ctx,
		DB:            db,
		RedisClient:   redisClient,
		HTTPClient:    httpClient,
		S3Client:      s3Client,
		KafkaProducer: kafkaProducer,
		Payload:       payload,
	}
}
