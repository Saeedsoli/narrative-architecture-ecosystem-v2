// apps/backend/internal/infrastructure/cache/session_store.go

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// SessionData ساختار داده‌هایی است که در یک Session ذخیره می‌شود.
type SessionData struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	// می‌توان فیلدهای دیگری مانند IP یا User-Agent را نیز اضافه کرد.
}

// SessionStore رابطی برای مدیریت Sessionها است.
type SessionStore interface {
	// Set یک Session جدید را با TTL مشخص ذخیره می‌کند.
	Set(ctx context.Context, sessionID string, data *SessionData, ttl time.Duration) error
	
	// Get یک Session را بر اساس شناسه آن بازیابی می‌کند.
	Get(ctx context.Context, sessionID string) (*SessionData, error)
	
	// Delete یک Session را حذف می‌کند (برای Logout).
	Delete(ctx context.Context, sessionID string) error
}

// RedisSessionStore پیاده‌سازی SessionStore با استفاده از Redis است.
type RedisSessionStore struct {
	client *redis.Client
}

// NewRedisSessionStore یک نمونه جدید از RedisSessionStore ایجاد می‌کند.
func NewRedisSessionStore(client *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{client: client}
}

func (s *RedisSessionStore) getSessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

// Set یک Session جدید را در Redis ذخیره می‌کند.
func (s *RedisSessionStore) Set(ctx context.Context, sessionID string, data *SessionData, ttl time.Duration) error {
	key := s.getSessionKey(sessionID)
	val, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}
	return s.client.Set(ctx, key, val, ttl).Err()
}

// Get یک Session را از Redis بازیابی می‌کند.
func (s *RedisSessionStore) Get(ctx context.Context, sessionID string) (*SessionData, error) {
	key := s.getSessionKey(sessionID)
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err // Redis `ErrNil` را برمی‌گرداند اگر کلید پیدا نشود.
	}

	var data SessionData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}
	return &data, nil
}

// Delete یک Session را از Redis حذف می‌کند.
func (s *RedisSessionStore) Delete(ctx context.Context, sessionID string) error {
	key := s.getSessionKey(sessionID)
	return s.client.Del(ctx, key).Err()
}