// apps/platform/lib/api/community.ts

import { apiClient } from './client';
import type { Topic, Post } from '@narrative-arch/types'; // فرض بر وجود این تایپ‌ها

interface ListResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export const listTopics = async (page: number = 1): Promise<ListResponse<Topic>> => {
  const { data } = await apiClient.get('/community/topics', { params: { page } });
  return data;
};

export const getTopic = async (id: string): Promise<Topic> => {
  const { data } = await apiClient.get(`/community/topics/${id}`);
  return data;
};

export const createTopic = async (topicData: { title: string; body: string; locale: 'fa'; tags?: string[] }): Promise<Topic> => {
  const { data } = await apiClient.post('/community/topics', topicData);
  return data;
};

export const listPosts = async (topicId: string, page: number = 1): Promise<ListResponse<Post>> => {
  const { data } = await apiClient.get(`/community/topics/${topicId}/posts`, { params: { page } });
  return data;
};

export const createPost = async (topicId: string, body: string, parentId?: string): Promise<Post> => {
  const { data } = await apiClient.post(`/community/topics/${topicId}/posts`, { body, parentId });
  return data;
};