# 🍳 کتابچه راهنمای توسعه: پیاده‌سازی یک ویژگی جدید

این مستند به شما نشان می‌دهد که چگونه یک ویژگی جدید (Use Case) را از ابتدا تا انتها با استفاده از معماری و الگوهای موجود در پروژه پیاده‌سازی کنید.

**سناریوی نمونه:** افزودن قابلیت **"بوکمارک کردن مقالات"**.

---

### **مرحله 1: طراحی دیتابیس (لایه PostgreSQL)**

ما به یک جدول جدید برای ذخیره بوکمارک‌ها نیاز داریم.

#### **1.1. ایجاد فایل مایگریشن**

در ترمینال، از ریشه پروژه، دستور زیر را اجرا کنید:

```bash
pnpm run db:migrate:create -- create_bookmarks_table
```
این دستور دو فایل جدید در `apps/backend/internal/infrastructure/database/migrations/` ایجاد می‌کند.

#### **1.2. نوشتن کد مایگریشن**

**`..._create_bookmarks_table.up.sql`:**
```sql
CREATE TABLE IF NOT EXISTS article_bookmarks (
  id         ulid PRIMARY KEY,
  user_id    ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  article_id ulid NOT NULL, -- ارجاع به ID مقاله در MongoDB
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, article_id)
);

CREATE INDEX IF NOT EXISTS ix_article_bookmarks_user ON article_bookmarks (user_id, created_at DESC);
```

**`..._create_bookmarks_table.down.sql`:**
```sql
DROP INDEX IF EXISTS ix_article_bookmarks_user;
DROP TABLE IF EXISTS article_bookmarks;
```

#### **1.3. اجرای مایگریشن**

```bash
pnpm run db:migrate:up
```

---

### **مرحله 2: لایه Domain**

موجودیت و رابط `Repository` را تعریف می‌کنیم.

**`apps/backend/internal/domain/bookmark/bookmark.go`:**
```go
package bookmark

import "time"

type Bookmark struct {
    ID        string
    UserID    string
    ArticleID string
    CreatedAt time.Time
}
```

**`apps/backend/internal/domain/bookmark/repository.go`:**
```go
package bookmark

import "context"

type Repository interface {
    Create(ctx context.Context, b *Bookmark) error
    Delete(ctx context.Context, userID, articleID string) error
    FindByUser(ctx context.Context, userID string) ([]*Bookmark, error)
}
```

---

### **مرحله 3: لایه Infrastructure**

پیاده‌سازی واقعی `Repository`.

**`apps/backend/internal/infrastructure/database/postgres/bookmark_repository.go`:**
```go
package postgres

import (
	"context"
	"database/sql"
	// ... سایر import ها
)

type BookmarkRepository struct {
	db *sql.DB
}

func NewBookmarkRepository(db *sql.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

func (r *BookmarkRepository) Create(ctx context.Context, b *bookmark.Bookmark) error {
	query := `INSERT INTO article_bookmarks (id, user_id, article_id) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, b.ID, b.UserID, b.ArticleID)
	return err
}

// ... پیاده‌سازی سایر متدها (Delete, FindByUser)
```

---

### **مرحله 4: لایه Application (Use Cases)**

منطق تجاری را در Use Caseها کپسوله می‌کنیم.

**`apps/backend/internal/application/bookmark/add_bookmark.go`:**
```go
package bookmark

import (
	"context"
	"github.com/oklog/ulid/v2"
	"narrative-architecture/apps/backend/internal/domain/bookmark"
)

type AddBookmarkRequest struct {
	UserID    string
	ArticleID string
}

type AddBookmarkUseCase struct {
	repo bookmark.Repository
}

func NewAddBookmarkUseCase(repo bookmark.Repository) *AddBookmarkUseCase {
	return &AddBookmarkUseCase{repo: repo}
}

func (uc *AddBookmarkUseCase) Execute(ctx context.Context, req AddBookmarkRequest) error {
	newBookmark := &bookmark.Bookmark{
		ID:        ulid.New().String(),
		UserID:    req.UserID,
		ArticleID: req.ArticleID,
	}
	return uc.repo.Create(ctx, newBookmark)
}
```

---

### **مرحله 5: لایه Presentation (Handler & Route)**

یک API Endpoint برای این ویژگی ایجاد می‌کنیم.

**`apps/backend/internal/interfaces/http/handlers/bookmark_handler.go`:**
```go
package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	app_bookmark "narrative-architecture/apps/backend/internal/application/bookmark"
)

type BookmarkHandler struct {
	addBookmarkUC *app_bookmark.AddBookmarkUseCase
}

func NewBookmarkHandler(addUC *app_bookmark.AddBookmarkUseCase) *BookmarkHandler {
	return &BookmarkHandler{addBookmarkUC: addUC}
}

func (h *BookmarkHandler) AddBookmark(c *gin.Context) {
	var req struct {
		ArticleID string `json:"articleId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetString("userID") // از AuthMiddleware

	err := h.addBookmarkUC.Execute(c.Request.Context(), app_bookmark.AddBookmarkRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add bookmark"})
		return
	}

	c.Status(http.StatusCreated)
}
```

**`apps/backend/cmd/api/main.go` (بخش اتصال):**
```go
// ... در بخش Dependency Injection
bookmarkRepo := postgres.NewBookmarkRepository(db)
addBookmarkUC := app_bookmark.NewAddBookmarkUseCase(bookmarkRepo)
bookmarkHandler := handlers.NewBookmarkHandler(addBookmarkUC)

// ... در بخش Router Setup (داخل گروه protected)
bookmarks := protected.Group("/bookmarks")
{
	bookmarks.POST("", bookmarkHandler.AddBookmark)
	// bookmarks.DELETE("/:articleId", bookmarkHandler.RemoveBookmark)
	// bookmarks.GET("", bookmarkHandler.ListBookmarks)
}
```

---

### **مرحله 6: Frontend**

در نهایت، در Frontend یک دکمه برای فراخوانی این API اضافه می‌کنیم.

**در کامپوننت `ArticleViewer`:**
```tsx
// ...
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';

// ...
const addBookmarkMutation = useMutation({
  mutationFn: (articleId: string) => apiClient.post('/bookmarks', { articleId }),
  onSuccess: () => {
    // نمایش پیام موفقیت با Toast
    console.log('Bookmarked!');
  },
});

// ...
<button onClick={() => addBookmarkMutation.mutate(article.id)}>
  بوکمارک کردن
</button>
```

با دنبال کردن این 6 مرحله، شما می‌توانید هر ویژگی جدیدی را به‌صورت استاندارد و هماهنگ با معماری کلی پروژه پیاده‌سازی کنید.