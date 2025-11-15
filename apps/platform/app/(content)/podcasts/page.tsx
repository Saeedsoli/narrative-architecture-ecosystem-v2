// apps/platform/app/(content)/podcasts/page.tsx

'use client';

import { useQuery } from '@tanstack/react-query';
import { listArticles } from '@/lib/api/articles'; // از همان API مقالات استفاده می‌کنیم
import { ArticleCard } from '@/components/content/article-card';
import { LoadingSpinner } from '@/components/shared/loading-spinner';

export default function PodcastsPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['podcasts'],
    queryFn: () => listArticles({ type: 'podcast' }), // فیلتر بر اساس نوع
  });

  if (isLoading) return <LoadingSpinner />;

  return (
    <div className="max-w-7xl mx-auto py-12 px-4">
      <h1 className="text-4xl font-bold mb-8">پادکست‌ها</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {data?.data.map((podcast) => (
          <ArticleCard key={podcast.id} article={podcast} />
        ))}
      </div>
    </div>
  );
}