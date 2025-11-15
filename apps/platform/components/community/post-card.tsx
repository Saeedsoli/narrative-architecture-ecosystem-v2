// apps/platform/components/community/post-card.tsx

import type { Post } from '@narrative-arch/types';
import { PersianDate } from '@/components/shared/persian-date';

interface PostCardProps {
  post: Post;
}

export function PostCard({ post }: PostCardProps) {
  return (
    <div className="flex space-x-4 space-x-reverse p-4 border rounded-lg">
      <div className="flex-shrink-0">
        <img src={post.user.avatar} alt={post.user.username} className="w-12 h-12 rounded-full" />
      </div>
      <div className="flex-1">
        <div className="flex justify-between items-center">
          <span className="font-semibold">{post.user.username}</span>
          <PersianDate date={post.createdAt} className="text-xs text-gray-500" />
        </div>
        <p className="mt-2 text-gray-700">{post.body}</p>
        <div className="mt-4 flex items-center space-x-4 space-x-reverse">
          <button className="text-sm text-gray-500">پاسخ</button>
          <button className="text-sm text-gray-500">👍 {post.likes_count}</button>
          <button className="text-sm text-gray-500">👎 {post.dislikes_count}</button>
        </div>
      </div>
    </div>
  );
}