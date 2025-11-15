package community

import "time"

type Topic struct {
	ID        string
	Locale    string
	Title     string
	Body      string
	Tags      []string
	Author    struct {
		ID       string
		Username string
	}
	Status    string // "open", "locked"
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type Post struct {
	ID        string
	TopicID   string
	ParentID  *string
	User      struct {
		ID       string
		Username string
		Avatar   string
	}
	Body         string
	LikesCount   int
	DislikesCount int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
}