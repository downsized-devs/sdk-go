// Package storage wraps the AWS S3 v1 SDK with a small Interface
// covering upload, download, delete, presigned URL generation, and
// static URL construction.
package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/logger"
)

// Interface is the public surface of the storage package. Mockable for tests.
type Interface interface {
	Upload(ctx context.Context, key string, filename, filemimetype string, data []byte) (url string, err error)
	Download(ctx context.Context, url string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	GetPresignedUrl(ctx context.Context, key string) (string, error)
	GetPresignedUrlWithDuration(ctx context.Context, key string, presignedDuration time.Duration) (string, error)
	CreateUrlByKey(key string) string
}

// Config wraps the AWS-specific options consumed by Init.
type Config struct {
	AWSS3 AWSS3Config
}

// AWSS3Config carries the bucket, region, IAM credentials, and the
// default presign duration applied by GetPresignedUrl.
type AWSS3Config struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	PresignDuration time.Duration
}

// s3API is the internal seam over *s3.S3 so tests can swap in a fake.
// *s3.S3 satisfies it automatically.
type s3API interface {
	PutObjectWithContext(ctx aws.Context, in *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error)
	GetObjectWithContext(ctx aws.Context, in *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error)
	DeleteObjectWithContext(ctx aws.Context, in *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error)
	GetObjectRequest(in *s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput)
}

type storage struct {
	s3     s3API
	config Config
	log    logger.Interface
}

// Init constructs an S3-backed storage client. It calls log.Fatal when
// the configured credentials are empty so misconfigured services fail
// fast at startup rather than at the first request.
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
	defer obj.Body.Close()

	data, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeStorageS3Download, "failed to read object body for key %s: %v", key, err)
	}

	if obj.ContentLength != nil && int64(len(data)) != *obj.ContentLength {
		return nil, errors.NewWithCode(codes.CodeStorageS3Download, "short read for key %s: got %d bytes, expected %d", key, len(data), *obj.ContentLength)
	}

	return data, nil
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
