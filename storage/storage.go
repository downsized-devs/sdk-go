package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/logger"
)

type Interface interface {
	Upload(ctx context.Context, key string, filename, filemimetype string, data []byte) (url string, err error)
	Download(ctx context.Context, url string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	GetPresignedUrl(ctx context.Context, key string) (string, error)
	GetPresignedUrlWithDuration(ctx context.Context, key string, presignedDuration time.Duration) (string, error)
	CreateUrlByKey(key string) string
}

type Config struct {
	AWSS3 AWSS3Config
}

type AWSS3Config struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	PresignDuration time.Duration
}

type storage struct {
	s3     *s3.S3
	config Config
	log    logger.Interface
}

func Init(cfg Config, log logger.Interface) Interface {
	if cfg.AWSS3.AccessKeyID == "" || cfg.AWSS3.SecretAccessKey == "" {
		log.Fatal(context.Background(), "storage credentials not found")
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWSS3.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AWSS3.AccessKeyID, cfg.AWSS3.SecretAccessKey, ""),
	}))
	s3 := s3.New(sess)

	return &storage{
		s3:     s3,
		config: cfg,
		log:    log,
	}
}

func (s *storage) Upload(ctx context.Context, key string, filename, filemimetype string, data []byte) (string, error) {
	obj, err := s.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(s.config.AWSS3.BucketName),
		Key:                aws.String(key),
		Body:               bytes.NewReader(data),
		ContentDisposition: aws.String(fmt.Sprintf(`attachment; filename="%s"`, filename)),
		ContentType:        aws.String(filemimetype),
	})
	if err != nil {
		return "", errors.NewWithCode(codes.CodeStorageS3Upload, "failed uploading file :%v, with err: %v", obj.String(), err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.config.AWSS3.BucketName, s.config.AWSS3.Region, key), nil
}

func (s *storage) Download(ctx context.Context, key string) ([]byte, error) {
	s3ObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.config.AWSS3.BucketName),
		Key:    aws.String(key),
	}
	obj, err := s.s3.GetObjectWithContext(ctx, s3ObjectInput)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeStorageS3Download, "failed to download file, with err: %v", err)
	}

	size := int(*obj.ContentLength)
	buffer := make([]byte, size)
	defer obj.Body.Close()
	var bbuffer bytes.Buffer
	for {
		byteSize, err := obj.Body.Read(buffer)
		if byteSize > 0 {
			bbuffer.Write(buffer[:byteSize])
		} else if err == io.EOF || err != nil {
			break
		}
	}

	return bbuffer.Bytes(), nil
}

func (s *storage) Delete(ctx context.Context, key string) error {
	obj, err := s.s3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.AWSS3.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return errors.NewWithCode(codes.CodeStorageS3Delete, "failed to delete obj: %v, with err: %v", obj.String(), err)
	}

	return nil
}

func (s *storage) GetPresignedUrl(ctx context.Context, key string) (string, error) {
	s3ObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.config.AWSS3.BucketName),
		Key:    aws.String(key),
	}

	req, _ := s.s3.GetObjectRequest(s3ObjectInput)

	urlStr, err := req.Presign(s.config.AWSS3.PresignDuration)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeStorageS3Download, "failed to get presigned url key: %v, with err: %v", key, err)
	}

	return urlStr, nil
}

func (s *storage) GetPresignedUrlWithDuration(ctx context.Context, key string, presignedDuration time.Duration) (string, error) {
	s3ObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.config.AWSS3.BucketName),
		Key:    aws.String(key),
	}

	req, _ := s.s3.GetObjectRequest(s3ObjectInput)

	urlStr, err := req.Presign(presignedDuration)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeStorageS3Download, "failed to get presigned url key: %v, with err: %v", key, err)
	}

	return urlStr, nil
}

func (s *storage) CreateUrlByKey(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.config.AWSS3.BucketName, s.config.AWSS3.Region, key)
}
