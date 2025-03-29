package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
	"os"
	"time"
)

func CreateSessionAws(url *string, accessKeyID, secretAccessKey, token, region string) (*session.Session, error) {
	config := &aws.Config{
		Endpoint:         url,
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	}
	if accessKeyID == "" || secretAccessKey == "" {
		config.Credentials = nil
	}

	sess, err := session.NewSession(
		config,
	)
	if err != nil {
		return nil, err
	}

	return sess, nil

}

type DownloadFileFromS3Func func(logger *zap.Logger, filename string) (*string, error)

func NewDownloadFileFromS3Func(sess *session.Session, bucket string) DownloadFileFromS3Func {
	return func(logger *zap.Logger, filename string) (*string, error) {

		logger.Info(fmt.Sprintf("Download file %s from S3", filename))

		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		downloader := s3manager.NewDownloader(sess)
		numBytes, err := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(filename),
			})
		if err != nil {
			logger.Error(fmt.Sprintf("download err: %s", err.Error()))
			return nil, err
		}
		logger.Info(fmt.Sprintf("File: %s, Size: %d", filename, numBytes))

		return &filename, nil

	}
}

type S3UploadFunc func(logger *zap.Logger, objectKey string, body []byte) error

func NewS3Upload(sess *session.Session, bucket string) S3UploadFunc {
	return func(logger *zap.Logger, objectKey string, body []byte) error {
		_, err := s3.New(sess).PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
			Body:   bytes.NewReader(body),
		})

		if err != nil {
			logger.Error(fmt.Sprintf("upload err: %s", err.Error()))
			return err
		}
		return nil

	}
}

func GenerateTemporarilyUrl(sess *session.Session, bucket, key string) (string, error) {

	svc := s3.New(sess)

	resp, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	//2:15
	return resp.Presign(time.Minute * 15)
}
