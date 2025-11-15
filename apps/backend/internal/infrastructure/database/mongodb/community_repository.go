// apps/backend/internal/infrastructure/database/mongodb/community_repository.go

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

type CommunityRepository struct {
	topics *mongo.Collection
	posts  *mongo.Collection
	votes  *mongo.Collection
}

func NewCommunityRepository(db *mongo.Database) *CommunityRepository {
	return &CommunityRepository{
		topics: db.Collection("forum_topics"),
		posts:  db.Collection("forum_posts"),
		votes:  db.Collection("comment_votes"),
	}
}

func (r *CommunityRepository) CreateTopic(ctx context.Context, topic *community.Topic) error {
	_, err := r.topics.InsertOne(ctx, topic)
	return err
}

func (r *CommunityRepository) FindTopicByID(ctx context.Context, id string) (*community.Topic, error) {
	var topic community.Topic
	filter := bson.M{"_id": id, "deletedAt": nil}
	err := r.topics.FindOne(ctx, filter).Decode(&topic)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("topic not found")
		}
		return nil, err
	}
	return &topic, nil
}

func (r *CommunityRepository) ListTopics(ctx context.Context, locale string, page, pageSize int) ([]*community.Topic, int64, error) {
	filter := bson.M{"locale": locale, "deletedAt": nil}
	
	total, err := r.topics.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "updatedAt", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.topics.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var topics []*community.Topic
	if err := cursor.All(ctx, &topics); err != nil {
		return nil, 0, err
	}
	
	return topics, total, nil
}

func (r *CommunityRepository) CreatePost(ctx context.Context, post *community.Post) error {
	_, err := r.posts.InsertOne(ctx, post)
	if err != nil {
		return err
	}

	// Denormalization: Update topic's updatedAt to bring it to the top
	_, err = r.topics.UpdateOne(
		ctx,
		bson.M{"_id": post.TopicID},
		bson.M{"$set": bson.M{"updatedAt": post.CreatedAt}},
	)
	return err
}

func (r *CommunityRepository) FindPostsByTopicID(ctx context.Context, topicID string, page, pageSize int) ([]*community.Post, int64, error) {
	filter := bson.M{"topic_id": topicID, "deletedAt": nil}

	total, err := r.posts.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.posts.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []*community.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, 0, err
	}
	
	return posts, total, nil
}

func (r *CommunityRepository) AddVote(ctx context.Context, userID, targetID, targetType string, value int) error {
	filter := bson.M{"target_id": targetID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"value":      value,
			"updatedAt":  time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":       ulid.New().String(),
			"createdAt": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	
	res, err := r.votes.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	// اگر رأی جدید بود، شمارنده را در پست/کامنت آپدیت کن
	if res.UpsertedCount > 0 {
		var updateField string
		if value == 1 {
			updateField = "likes_count"
		} else {
			updateField = "dislikes_count"
		}
		
		var targetCollection *mongo.Collection
		if targetType == "forum_post" {
			targetCollection = r.posts
		} else {
			// targetCollection = r.article_comments (نیاز به تزریق دارد)
			return nil
		}
		
		_, err = targetCollection.UpdateOne(
			ctx,
			bson.M{"_id": targetID},
			bson.M{"$inc": bson.M{updateField: 1}},
		)
	}

	return err
}