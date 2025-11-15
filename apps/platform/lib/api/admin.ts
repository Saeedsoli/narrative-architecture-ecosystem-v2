// apps/platform/lib/api/admin.ts

import { apiClient } from './client';
import type { User } from '@narrative-arch/types';

// دریافت لیست کاربران
export const listUsers = async (page: number = 1): Promise<{ data: User[] }> => {
  const { data } = await apiClient.get('/admin/users', { params: { page } });
  return data;
};

// تغییر وضعیت یک کاربر
export const updateUserStatus = async ({ userId, status }: { userId: string, status: string }): Promise<void> => {
  await apiClient.put(`/admin/users/${userId}/status`, { status });
};