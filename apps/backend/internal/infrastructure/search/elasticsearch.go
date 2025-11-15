// apps/backend/internal/infrastructure/search/elasticsearch.go

package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v8"
	"narrative-architecture/apps/backend/internal/domain/article"
)

type ElasticsearchService struct {
	client *elasticsearch.Client
}

func NewElasticsearchService(client *elasticsearch.Client) *ElasticsearchService {
	return &ElasticsearchService{client: client}
}

// SearchArticles مقالات را بر اساس کوئری جستجو می‌کند.
func (s *ElasticsearchService) SearchArticles(ctx context.Context, query, locale string, page, pageSize int) ([]*article.Article, int64, error) {
	indexAlias := "articles_fa"
	if locale == "en" {
		indexAlias = "articles_en"
	}

	esQuery := map[string]interface{}{
		"from": (page - 1) * pageSize,
		"size": pageSize,
		"query": {
			"multi_match": {
				"query":  query,
				"fields": []string{"title^3", "excerpt^2", "content"},
				"fuzziness": "AUTO",
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, 0, err
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(indexAlias),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, 0, fmt.Errorf("elasticsearch error: %s", string(body))
	}

	var r struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, err
	}

	var articles []*article.Article
	for _, hit := range r.Hits.Hits {
		var art article.Article
		if err := json.Unmarshal(hit.Source, &art); err != nil {
			continue // یا لاگ بگیرید
		}
		articles = append(articles, &art)
	}

	return articles, r.Hits.Total.Value, nil
}