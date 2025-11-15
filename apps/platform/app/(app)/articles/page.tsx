// apps/platform/app/(app)/articles/page.tsx

'use client';

import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'next/navigation';
import { listArticles } from '@/lib/api/articles';
import { ArticleCard } from '@/components/content/article-card';
import { Pagination } from '@/components/shared/pagination';
import { LoadingSpinner } from '@/components/shared/loading-spinner';

export default function ArticlesPage() {
  const searchParams = useSearchParams();
  const page = parseInt(searchParams.get('page') || '1', 10);

  const { data, isLoading, error } = useQuery({
    queryKey: ['articles', page],
    queryFn: () => listArticles({ page }),
    keepPreviousData: true,
  });

  if (isLoading) return <LoadingSpinner />;
  if (error) return <div>خطا در دریافت مقالات.</div>;

  return (
    <div className="max-w-7xl mx-auto py-12 px-4">
      <h1 className="text-4xl font-bold mb-8">مقالات</h1>
      
      {data && data.data.length > 0 ? (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {data.data.map((article) => (
              <ArticleCard key={article.id} article={article} />
            ))}
          </div>
          <div className="mt-12">
            <Pagination
              currentPage={data.page}
              totalPages={data.totalPages}
              baseUrl="/articles"
            />
          </div>
        </>
      ) : (
        <p>هیچ مقاله‌ای برای نمایش وجود ندارد.</p>
      )}
    </div>
  );
}