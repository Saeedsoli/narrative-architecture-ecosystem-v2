// apps/platform/lib/api/users.ts

import { apiClient } from './client';
import type { User } from '@narrative-arch/types';

interface UpdateProfileData {
  fullName?: string;
  username?: string;
}

export const updateUserProfile = async (data: UpdateProfileData): Promise<User> => {
  const { data: updatedUser } = await apiClient.put('/users/me/profile', data);
  return updatedUser;
};