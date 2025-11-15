// packages/types/src/user.ts

export interface User {
  id: string;
  email: string;
  username: string;
  fullName: string;
  avatarUrl?: string;
  roles: string[];
  status: 'active' | 'suspended' | 'deleted';
}