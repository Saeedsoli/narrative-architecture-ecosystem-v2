// apps/backend/internal/infrastructure/database/mongodb/article_repository.go

package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"narrative-architecture/apps/backend/internal/domain/article"
	"narrative-architecture/apps/backend/lib/utils"
)

// MongoArticle ساختار یک سند مقاله در MongoDB را تعریف می‌کند.
type MongoArticle struct {
	ID             string     `bson:"_id"`
	ContentGroupID string     `bson:"content_group_id"`
	Locale         string     `bson:"locale"`
	Slug           string     `bson:"slug"`
	Title          string     `bson:"title"`
	Excerpt        string     `bson:"excerpt,omitempty"`
	Content        string     `bson:"content"`
	CoverImage     MongoCoverImage `bson:"coverImage,omitempty"`
	Author         MongoAuthor     `bson:"author"`
	Metadata       MongoMetadata   `bson:"metadata,omitempty"`
	Status         string     `bson:"status"`
	PublishedAt    *time.Time `bson:"publishedAt,omitempty"`
	CreatedAt      time.Time  `bson:"createdAt"`
	UpdatedAt      time.Time  `bson:"updatedAt"`
	DeletedAt      *time.Time `bson:"deletedAt,omitempty"`
}

type MongoCoverImage struct {
	URL string `bson:"url"`
	Alt string `bson:"alt"`
}

type MongoAuthor struct {
	ID     string `bson:"id"`
	Name   string `bson:"name"`
	Avatar string `bson:"avatar,omitempty"`
}

type MongoMetadata struct {
	Tags       []string `bson:"tags,omitempty"`
	Category   string   `bson:"category,omitempty"`
	ReadTime   int      `bson:"readTime,omitempty"`
	Difficulty string   `bson:"difficulty,omitempty"`
}

// ArticleRepository پیاده‌سازی رابط Repository با MongoDB است.
type ArticleRepository struct {
	collection *mongo.Collection
}

// NewArticleRepository یک نمونه جدید از ArticleRepository ایجاد می‌کند.
func NewArticleRepository(db *mongo.Database) *ArticleRepository {
	return &ArticleRepository{collection: db.Collection("articles")}
}

// FindByID یک مقاله را بر اساس ID پیدا می‌کند.
func (r *ArticleRepository) FindByID(ctx context.Context, id string) (*article.Article, error) {
	var result MongoArticle
	filter := bson.M{"_id": id, "deletedAt": nil}

	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("article not found")
		}
		return nil, err
	}
	return toDomainArticle(&result), nil
}

// FindBySlug یک مقاله را بر اساس اسلاگ آن پیدا می‌کند.
func (r *ArticleRepository) FindBySlug(ctx context.Context, slug string) (*article.Article, error) {
	var result MongoArticle
	filter := bson.M{"slug": slug, "deletedAt": nil}

	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("article not found with this slug")
		}
		return nil, err
	}

	return toDomainArticle(&result), nil
}

// Find لیستی از مقالات را با فیلتر و صفحه‌بندی برمی‌گرداند.
func (r *ArticleRepository) Find(ctx context.Context, filter article.Filter) ([]*article.Article, int64, error) {
	mongoFilter := bson.M{"deletedAt": nil}

	if len(filter.IDs) > 0 {
		mongoFilter["_id"] = bson.M{"$in": filter.IDs}
	}
	if filter.Locale != "" {
		mongoFilter["locale"] = filter.Locale
	}
	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}
	if filter.AuthorID != "" {
		mongoFilter["author.id"] = filter.AuthorID
	}
	if filter.Category != "" {
		mongoFilter["metadata.category"] = filter.Category
	}
	if len(filter.Tags) > 0 {
		mongoFilter["metadata.tags"] = bson.M{"$in": filter.Tags}
	}
	
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "publishedAt", Value: -1}})
	opts.SetSkip(int64((filter.Page - 1) * filter.PageSize))
	opts.SetLimit(int64(filter.PageSize))

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var articles []*article.Article
	for cursor.Next(ctx) {
		var result MongoArticle
		if err := cursor.Decode(&result); err != nil {
			continue // یا لاگ بگیرید
		}
		articles = append(articles, toDomainArticle(&result))
	}

	return articles, total, nil
}

// FindTranslations لیستی از ترجمه‌های موجود را برمی‌گرداند.
func (r *ArticleRepository) FindTranslations(ctx context.Context, contentGroupID string) ([]article.Translation, error) {
	filter := bson.M{"content_group_id": contentGroupID, "deletedAt": nil}
	projection := options.Find().SetProjection(bson.M{"locale": 1, "slug": 1})

	cursor, err := r.collection.Find(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var translations []article.Translation
	for cursor.Next(ctx) {
		var result struct {
			Locale string `bson:"locale"`
			Slug   string `bson:"slug"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		translations = append(translations, article.Translation{
			Locale: result.Locale,
			Slug:   result.Slug,
		})
	}
	return translations, nil
}

// Save یک مقاله را ذخیره یا آپدیت می‌کند.
func (r *ArticleRepository) Save(ctx context.Context, art *article.Article) error {
	// تولید خودکار اسلاگ بر اساس عنوان
	art.Slug = utils.CreateSlug(art.Title)

	mongoDoc := fromDomainArticle(art)
	mongoDoc.UpdatedAt = time.Now()
	if mongoDoc.CreatedAt.IsZero() {
		mongoDoc.CreatedAt = mongoDoc.UpdatedAt
	}
	
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": art.ID}
	
	// ما از $set استفاده می‌کنیم تا فیلدهایی که در آپدیت ارسال نمی‌شوند، حفظ شوند.
	update := bson.M{"$set": mongoDoc}

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// Delete یک مقاله را به‌صورت نرم حذف می‌کند.
func (r *ArticleRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"deletedAt": time.Now()}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("article not found to delete")
	}
	return nil
}

// toDomainArticle سند MongoDB را به موجودیت Domain تبدیل می‌کند.
func toDomainArticle(m *MongoArticle) *article.Article {
	return &article.Article{
		ID:             m.ID,
		ContentGroupID: m.ContentGroupID,
		Locale:         m.Locale,
		Slug:           m.Slug,
		Title:          m.Title,
		Excerpt:        m.Excerpt,
		Content:        m.Content,
		CoverImage: struct {
			URL string
			Alt string
		}{URL: m.CoverImage.URL, Alt: m.CoverImage.Alt},
		Author: struct {
			ID     string
			Name   string
			Avatar string
		}{ID: m.Author.ID, Name: m.Author.Name, Avatar: m.Author.Avatar},
		Metadata: struct {
			Tags       []string
			Category   string
			ReadTime   int
			Difficulty string
		}{Tags: m.Metadata.Tags, Category: m.Metadata.Category, ReadTime: m.Metadata.ReadTime, Difficulty: m.Metadata.Difficulty},
		Status:      article.ArticleStatus(m.Status),
		PublishedAt: m.PublishedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

// fromDomainArticle موجودیت Domain را به سند MongoDB تبدیل می‌کند.
func fromDomainArticle(art *article.Article) *MongoArticle {
	return &MongoArticle{
		ID:             art.ID,
		ContentGroupID: art.ContentGroupID,
		Locale:         art.Locale,
		Slug:           art.Slug,
		Title:          art.Title,
		Excerpt:        art.Excerpt,
		Content:        art.Content,
		CoverImage:     MongoCoverImage{URL: art.CoverImage.URL, Alt: art.CoverImage.Alt},
		Author:         MongoAuthor{ID: art.Author.ID, Name: art.Author.Name, Avatar: art.Author.Avatar},
		Metadata:       MongoMetadata{Tags: art.Metadata.Tags, Category: art.Metadata.Category, ReadTime: art.Metadata.ReadTime, Difficulty: art.Metadata.Difficulty},
		Status:         string(art.Status),
		PublishedAt:    art.PublishedAt,
		CreatedAt:      art.CreatedAt,
        UpdatedAt:      art.UpdatedAt,
		DeletedAt:      art.DeletedAt,
	}
}