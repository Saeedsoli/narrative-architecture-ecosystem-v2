// apps/backend/internal/infrastructure/cache/cache_repository.go

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"narrative-architecture/apps/backend/internal/domain/article"
)

// --- Article Cache Implementation ---

// ArticleCache پیاده‌سازی رابط کش برای مقالات با Redis است.
type ArticleCache struct {
	client *redis.Client
}

// NewArticleCache یک نمونه جدید از ArticleCache ایجاد می‌کند.
func NewArticleCache(client *redis.Client) *ArticleCache {
	return &ArticleCache{client: client}
}

func (c *ArticleCache) getArticleKey(slug string) string {
	return fmt.Sprintf("cache:article:%s", slug)
}

// GetBySlug یک مقاله را از کش بر اساس اسلاگ آن می‌خواند.
func (c *ArticleCache) GetBySlug(ctx context.Context, slug string) (*article.Article, error) {
	key := c.getArticleKey(slug)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var art article.Article
	if err := json.Unmarshal([]byte(val), &art); err != nil {
		return nil, err
	}
	return &art, nil
}

// Set یک مقاله را در کش ذخیره می‌کند.
func (c *ArticleCache) Set(ctx context.Context, art *article.Article) error {
	key := c.getArticleKey(art.Slug)
	val, err := json.Marshal(art)
	if err != nil {
		return err
	}
	// کش کردن برای 1 ساعت
	return c.client.Set(ctx, key, val, 1*time.Hour).Err()
}

// Delete یک مقاله را از کش حذف می‌کند.
func (c *ArticleCache) Delete(ctx context.Context, art *article.Article) error {
	key := c.getArticleKey(art.Slug)
	return c.client.Del(ctx, key).Err()
}


// --- Session Cache Implementation (مثال) ---

type SessionCache struct {
	client *redis.Client
}

func NewSessionCache(client *redis.Client) *SessionCache {
	return &SessionCache{client: client}
}

func (c *SessionCache) SetSession(ctx context.Context, sessionID string, data map[string]interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, val, ttl).Err()
}

func (c *SessionCache) GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return data, nil
}