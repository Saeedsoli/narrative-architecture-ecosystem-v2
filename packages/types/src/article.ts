// packages/types/src/article.ts

export interface ArticleTranslation {
  locale: 'fa' | 'en';
  slug: string;
}

export interface Article {
  id: string;
  locale: 'fa' | 'en';
  slug: string;
  title: string;
  excerpt: string;
  content: string;
  coverImage: {
    url: string;
    alt: string;
  };
  author: {
    id: string;
    name: string;
    avatar: string;
  };
  metadata: {
    tags: string[];
    category: string;
    readTime: number;
    difficulty: 'beginner' | 'intermediate' | 'advanced';
  };
  publishedAt: string;
  updatedAt: string;
  translations: ArticleTranslation[];
}