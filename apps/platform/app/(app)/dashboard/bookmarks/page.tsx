// apps/platform/app/(app)/dashboard/bookmarks/page.tsx

'use client';

import { useQuery } from '@tanstack/react-query';
import { listBookmarks } from '@/lib/api/articles';
import { ArticleCard } from '@/components/content/article-card';
import { LoadingSpinner } from '@/components/shared/loading-spinner';
import type { Article } from '@narrative-arch/types';

export default function BookmarksPage() {
  const { data: bookmarks, isLoading, error } = useQuery<Article[]>({
    queryKey: ['bookmarks'],
    queryFn: listBookmarks,
  });

  if (isLoading) return <LoadingSpinner />;
  if (error) return <div>خطا در دریافت بوکمارک‌ها.</div>;

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-8">مقالات ذخیره شده</h1>
      {bookmarks && bookmarks.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {bookmarks.map((article) => (
            <ArticleCard key={article.id} article={article} />
          ))}
        </div>
      ) : (
        <p>شما هنوز هیچ مقاله‌ای را بوکمارک نکرده‌اید.</p>
      )}
    </div>
  );
}