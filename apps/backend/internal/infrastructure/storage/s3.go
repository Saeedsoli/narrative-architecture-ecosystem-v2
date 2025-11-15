// apps/backend/internal/infrastructure/storage/s3.go

package storage

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Config شامل تنظیمات اتصال به S3 است.
type S3Config struct {
	Region string
}

// NewS3Client یک کلاینت جدید S3 ایجاد و برمی‌گرداند.
func NewS3Client(cfg S3Config) *s3.Client {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.Region))
	if err != nil {
		log.Fatalf("FATAL: unable to load AWS SDK config, %v", err)
	}
	return s3.NewFromConfig(awsCfg)
}