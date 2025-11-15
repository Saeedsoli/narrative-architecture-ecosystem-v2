// packages/types/src/community.ts

interface Author {
  id: string;
  username: string;
  avatar?: string;
}

export interface Topic {
  id: string;
  title: string;
  author: Author;
  // ... سایر فیلدها
}

export interface Post {
  id: string;
  body: string;
  user: Author;
  likes_count: number;
  dislikes_count: number;
  createdAt: string;
}