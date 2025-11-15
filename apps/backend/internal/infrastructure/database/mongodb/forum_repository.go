// apps/backend/internal/infrastructure/database/mongodb/forum_repository.go

package mongodb

import (
	"context"
	"errors"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"narrative-architecture/apps/backend/internal/domain/community"
)

// ForumRepository پیاده‌سازی رابط برای عملیات مربوط به انجمن است.
type ForumRepository struct {
	topics *mongo.Collection
	posts  *mongo.Collection
	votes  *mongo.Collection
}

// NewForumRepository یک نمونه جدید ایجاد می‌کند.
func NewForumRepository(db *mongo.Database) *ForumRepository {
	return &ForumRepository{
		topics: db.Collection("forum_topics"),
		posts:  db.Collection("forum_posts"),
		votes:  db.Collection("comment_votes"),
	}
}

// ... (پیاده‌سازی کامل متدهای CreateTopic, FindTopicByID, ListTopics, CreatePost, FindPostsByTopicID, AddVote که قبلاً در community_repository.go ارائه شد)