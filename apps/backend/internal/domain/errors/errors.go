// apps/backend/internal/domain/errors/errors.go

package errors

import "errors"

// خطاهای استاندارد دامنه که در کل اپلیکیشن استفاده می‌شوند.
var (
	ErrNotFound         = errors.New("resource not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidInput     = errors.New("invalid input data")
	ErrConflict         = errors.New("resource conflict or already exists")
	ErrUnauthorized     = errors.New("unauthorized")
)