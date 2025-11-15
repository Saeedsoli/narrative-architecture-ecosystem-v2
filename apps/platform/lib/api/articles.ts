// apps/platform/lib/api/articles.ts

import { apiClient } from './client';
import type { Article } from '@narrative-arch/types';

// ساختار پاسخ برای لیست مقالات
interface ListArticlesResponse {
  data: Article[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// دریافت یک مقاله بر اساس اسلاگ
export const getArticleBySlug = async (slug: string): Promise<Article> => {
  const { data } = await apiClient.get<Article>(`/articles/${slug}`);
  return data;
};

// دریافت لیست مقالات با فیلتر و صفحه‌بندی
export const listArticles = async (params: { page?: number; pageSize?: number; tags?: string[] }): Promise<ListArticlesResponse> => {
  const { data } = await apiClient.get<ListArticlesResponse>('/articles', { params });
  return data;
};

// ایجاد یک مقاله جدید
export const createArticle = async (articleData: { title: string; content: string; locale: 'fa' | 'en' }): Promise<Article> => {
  const { data } = await apiClient.post<Article>('/articles', articleData);
  return data;
};

// آپدیت یک مقاله
export const updateArticle = async (id: string, articleData: Partial<Article>): Promise<Article> => {
  const { data } = await apiClient.put<Article>(`/articles/${id}`, articleData);
  return data;
};

// حذف یک مقاله
export const deleteArticle = async (id: string): Promise<void> => {
  await apiClient.delete(`/articles/${id}`);
};

// بوکمارک کردن یک مقاله
export const addBookmark = async (articleId: string): Promise<void> => {
  await apiClient.post(`/articles/${articleId}/bookmark`);
};

// حذف بوکمارک
export const removeBookmark = async (articleId: string): Promise<void> => {
  await apiClient.delete(`/articles/${articleId}/bookmark`);
};

// دریافت لیست بوکمارک‌های کاربر
export const listBookmarks = async (): Promise<Article[]> => {
  const { data } = await apiClient.get<Article[]>('/bookmarks');
  return data;
};