package s3

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
	"image"
	"image/jpeg"
	"os"
	"strings"
	"time"
)

func CreateSessionAws(url *string, accessKeyID, secretAccessKey, token, region string) (*session.Session, error) {
	config := &aws.Config{
		Endpoint:    url,
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
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

		filenameGen := fmt.Sprintf("%d_%s", time.Now().Unix(), strings.ReplaceAll(filename, "/", "_"))
		file, err := os.Create(filenameGen)
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
		logger.Info(fmt.Sprintf("File: %s, Size: %d", filenameGen, numBytes))

		return &filenameGen, nil

	}
}

type S3UploadBase64Func func(logger *zap.Logger, base64File, objectKey string) error

func S3UploadBase64(sess *session.Session, bucket string) S3UploadBase64Func {
	return func(logger *zap.Logger, base64Data, objectKey string) error {

		contentType := base64Data[strings.IndexByte(base64Data, ':')+1 : strings.IndexByte(base64Data, ',')]
		b64data := base64Data[strings.IndexByte(base64Data, ',')+1:]

		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))
		image, _, err := image.Decode(reader)
		if err != nil {
			return err
		}
		image.Bounds()

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, image, nil)
		if err != nil {
			return err
		}
		body := buf.Bytes()

		_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
			Bucket:          aws.String(bucket),
			Key:             aws.String(objectKey),
			Body:            bytes.NewReader(body),
			ContentEncoding: aws.String("base64"),
			ContentType:     aws.String(contentType),
		})

		if err != nil {
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
