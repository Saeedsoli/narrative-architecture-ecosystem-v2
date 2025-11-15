// apps/backend/internal/infrastructure/database/mongodb/content_repository.go

package mongodb

import (
	"context"
	"errors"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ContentRepository یک رابط عمومی برای مدیریت محتوا در کالکشن‌های مختلف است.
type ContentRepository struct {
	db *mongo.Database
}

// NewContentRepository یک نمونه جدید از ContentRepository ایجاد می‌کند.
func NewContentRepository(db *mongo.Database) *ContentRepository {
	return &ContentRepository{db: db}
}

// Moderate یک سند را در کالکشن مشخص شده، به‌صورت نرم حذف می‌کند.
func (r *ContentRepository) Moderate(ctx context.Context, targetType, targetID, action string) error {
	var collection *mongo.Collection

	switch targetType {
	case "article_comment":
		collection = r.db.Collection("article_comments")
	case "forum_topic":
		collection = r.db.Collection("forum_topics")
	case "forum_post":
		collection = r.db.Collection("forum_posts")
	default:
		return errors.New("unsupported target type for moderation")
	}

	if action == "rejected" {
		filter := bson.M{"_id": targetID}
		update := bson.M{"$set": bson.M{"deletedAt": time.Now()}}
		
		res, err := collection.UpdateOne(ctx, filter, update)
		if err != nil { return err }
		if res.MatchedCount == 0 { return errors.New("content to moderate not found") }
	}

	return nil
}