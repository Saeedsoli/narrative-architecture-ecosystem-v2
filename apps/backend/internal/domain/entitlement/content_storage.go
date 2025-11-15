package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ContentStorage struct {
	s3Client *s3.Client
	bucket   string
}

func NewContentStorage(client *s3.Client, bucket string) *ContentStorage {
	return &ContentStorage{s3Client: client, bucket: bucket}
}

// GetEncryptedChapter16 محتوای رمزنگاری شده فصل 16 را از S3 می‌خواند.
func (s *ContentStorage) GetEncryptedChapter16(ctx context.Context) ([]byte, error) {
	key := "book-content/chapter16.encrypted" // مسیر فایل در S3

	output, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	return io.ReadAll(output.Body)
}