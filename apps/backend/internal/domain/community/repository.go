// apps/backend/internal/domain/community/repository.go

package community

import "context"

type Repository interface {
	// Topic methods
	CreateTopic(ctx context.Context, topic *Topic) error
	FindTopicByID(ctx context.Context, id string) (*Topic, error)
	ListTopics(ctx context.Context, locale string, page, pageSize int) ([]*Topic, int64, error)
	
	// Post methods
	CreatePost(ctx context.Context, post *Post) error
	FindPostsByTopicID(ctx context.Context, topicID string, page, pageSize int) ([]*Post, int64, error)
	
	// Comment methods (برای article_comments)
	CreateComment(ctx context.Context, comment *Comment) error
	FindCommentByID(ctx context.Context, id string) (*Comment, error)
	
	// Vote methods
	AddVote(ctx context.Context, userID, targetID, targetType string, value int) error
}