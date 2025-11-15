// apps/platform/app/(app)/community/forum/[topicId]/page.tsx

'use client';

import { useQuery } from '@tanstack/react-query';
import { useParams } from 'next/navigation';
import { getTopic, listPosts } from '@/lib/api/community';
import { PostCard } from '@/components/community/post-card';
import { PostEditor } from '@/components/community/post-editor';
import { LoadingSpinner } from '@/components/shared/loading-spinner';

export default function TopicPage() {
  const params = useParams();
  const topicId = params.topicId as string;

  const { data: topic, isLoading: isLoadingTopic } = useQuery({
    queryKey: ['forum-topic', topicId],
    queryFn: () => getTopic(topicId),
    enabled: !!topicId,
  });

  const { data: posts, isLoading: isLoadingPosts } = useQuery({
    queryKey: ['forum-posts', topicId],
    queryFn: () => listPosts(topicId),
    enabled: !!topicId,
  });

  if (isLoadingTopic || isLoadingPosts) return <LoadingSpinner />;
  if (!topic) return <div>تاپیک یافت نشد.</div>;

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-2">{topic.title}</h1>
      <p className="text-sm text-gray-500 mb-8">ایجاد شده توسط {topic.author.username}</p>
      
      <div className="space-y-6">
        {posts?.data.map((post) => (
          <PostCard key={post.id} post={post} />
        ))}
      </div>
      
      <div className="mt-12 border-t pt-8">
        <h3 className="text-xl font-semibold mb-4">پاسخ خود را ثبت کنید</h3>
        <PostEditor topicId={topicId} />
      </div>
    </div>
  );
}