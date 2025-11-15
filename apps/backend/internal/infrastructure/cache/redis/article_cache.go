// apps/backend/internal/infrastructure/cache/redis/article_cache.go

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"narrative-architecture/apps/backend/internal/domain/article"
)

// ArticleCache پیاده‌سازی کش برای مقالات با Redis است.
type ArticleCache struct {
	client *redis.Client
}

// NewArticleCache یک نمونه جدید از ArticleCache ایجاد می‌کند.
func NewArticleCache(client *redis.Client) *ArticleCache {
	return &ArticleCache{client: client}
}

func (c *ArticleCache) getKey(slug string) string {
	return fmt.Sprintf("cache:article:%s", slug)
}

// GetBySlug یک مقاله را از کش بر اساس اسلاگ آن می‌خواند.
func (c *ArticleCache) GetBySlug(ctx context.Context, slug string) (*article.Article, error) {
	key := c.getKey(slug)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var art article.Article
	err = json.Unmarshal([]byte(val), &art)
	if err != nil {
		return nil, err
	}
	return &art, nil
}

// Set یک مقاله را در کش ذخیره می‌کند.
func (c *ArticleCache) Set(ctx context.Context, art *article.Article) error {
	key := c.getKey(art.Slug)
	val, err := json.Marshal(art)
	if err != nil {
		return err
	}
	// کش کردن برای 1 ساعت
	return c.client.Set(ctx, key, val, 1*time.Hour).Err()
}

// Delete یک مقاله را از کش حذف می‌کند (برای زمان آپدیت یا حذف).
func (c *ArticleCache) Delete(ctx context.Context, art *article.Article) error {
	key := c.getKey(art.Slug)
	return c.client.Del(ctx, key).Err()
}