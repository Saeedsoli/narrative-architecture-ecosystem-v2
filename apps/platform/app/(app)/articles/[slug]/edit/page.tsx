// apps/platform/app/(app)/articles/[slug]/edit/page.tsx

'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useParams, useRouter } from 'next/navigation';
import { getArticleBySlug, updateArticle } from '@/lib/api/articles';
import { ArticleForm } from '@/components/content/article-form';
import type { Article } from '@narrative-arch/types';

export default function EditArticlePage() {
  const params = useParams();
  const slug = params.slug as string;
  const router = useRouter();
  const queryClient = useQueryClient();

  const { data: article, isLoading } = useQuery<Article>({
    queryKey: ['article', slug],
    queryFn: () => getArticleBySlug(slug),
    enabled: !!slug,
  });

  const updateMutation = useMutation({
    mutationFn: (updatedData: Partial<Article>) => updateArticle(article!.id, updatedData),
    onSuccess: (updatedArticle) => {
      queryClient.invalidateQueries({ queryKey: ['article', slug] });
      queryClient.invalidateQueries({ queryKey: ['articles'] });
      router.push(`/articles/${updatedArticle.slug}`);
    },
    onError: (error: any) => {
      alert(`خطا در آپدیت مقاله: ${error.message}`);
    }
  });

  if (isLoading) return <div>در حال بارگذاری...</div>;
  if (!article) return <div>مقاله یافت نشد.</div>;

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-8">ویرایش مقاله</h1>
      <ArticleForm
        initialData={article}
        onSubmit={(data) => updateMutation.mutate(data)}
        isSubmitting={updateMutation.isLoading}
      />
    </div>
  );
}