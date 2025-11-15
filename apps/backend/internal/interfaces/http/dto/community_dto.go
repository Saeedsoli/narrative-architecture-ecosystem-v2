// apps/backend/internal/interfaces/http/dto/community_dto.go

package dto

import (
	"time"
	"narrative-architecture/apps/backend/internal/domain/community"
)

// --- Topic DTO ---

// AuthorInfo ساختار اطلاعات نویسنده برای نمایش در API است.
type AuthorInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

// TopicResponse ساختار پاسخ API برای یک تاپیک است.
type TopicResponse struct {
	ID        string     `json:"id"`
	Locale    string     `json:"locale"`
	Title     string     `json:"title"`
	Body      string     `json:"body,omitempty"` // بدنه معمولاً در لیست تاپیک‌ها لازم نیست
	Tags      []string   `json:"tags"`
	Author    AuthorInfo `json:"author"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// ToTopicResponse موجودیت Topic دامنه را به DTO پاسخ تبدیل می‌کند.
// پارامتر withBody برای کنترل ارسال بدنه کامل تاپیک است.
func ToTopicResponse(topic *community.Topic, withBody bool) *TopicResponse {
	if topic == nil {
		return nil
	}

	resp := &TopicResponse{
		ID:     topic.ID,
		Locale: topic.Locale,
		Title:  topic.Title,
		Tags:   topic.Tags,
		Author: AuthorInfo{
			ID:       topic.Author.ID,
			Username: topic.Author.Username,
		},
		Status:    topic.Status,
		CreatedAt: topic.CreatedAt,
		UpdatedAt: topic.UpdatedAt,
	}

	if withBody {
		resp.Body = topic.Body
	}

	return resp
}

// --- Post DTO ---

// PostResponse ساختار پاسخ API برای یک پست در انجمن است.
type PostResponse struct {
	ID            string     `json:"id"`
	TopicID       string     `json:"topicId"`
	ParentID      *string    `json:"parentId,omitempty"`
	User          AuthorInfo `json:"user"`
	Body          string     `json:"body"`
	LikesCount    int        `json:"likesCount"`
	DislikesCount int        `json:"dislikesCount"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

// ToPostResponse موجودیت Post دامنه را به DTO پاسخ تبدیل می‌کند.
func ToPostResponse(post *community.Post) *PostResponse {
	if post == nil {
		return nil
	}

	return &PostResponse{
		ID:       post.ID,
		TopicID:  post.TopicID,
		ParentID: post.ParentID,
		User: AuthorInfo{
			ID:       post.User.ID,
			Username: post.User.Username,
			Avatar:   post.User.Avatar,
		},
		Body:          post.Body,
		LikesCount:    post.LikesCount,
		DislikesCount: post.DislikesCount,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}
}