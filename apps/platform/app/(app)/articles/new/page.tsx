// apps/platform/app/(app)/articles/new/page.tsx

'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import type { Article } from '@narrative-arch/types';

export default function NewArticlePage() {
  const router = useRouter();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [locale, setLocale] = useState('fa');

  const createArticleMutation = useMutation<Article, Error, { title: string; content: string; locale: string }>({
    mutationFn: (newArticle) => apiClient.post('/articles', newArticle),
    onSuccess: (data) => {
      // پس از ایجاد موفق، به صفحه ویرایش مقاله هدایت شو
      router.push(`/articles/${data.slug}/edit`);
    },
    onError: (error) => {
      alert(`Failed to create article: ${error.message}`);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createArticleMutation.mutate({ title, content, locale });
  };

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-8">نوشتن مقاله جدید</h1>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="locale">زبان</label>
          <select
            id="locale"
            value={locale}
            onChange={(e) => setLocale(e.target.value)}
            className="w-full mt-2"
          >
            <option value="fa">فارسی</option>
            <option value="en">English</option>
          </select>
        </div>
        <div>
          <label htmlFor="title">عنوان</label>
          <input
            id="title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            className="w-full mt-2"
          />
        </div>
        <div>
          <label htmlFor="content">محتوا (Markdown)</label>
          <textarea
            id="content"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            required
            rows={15}
            className="w-full mt-2"
          />
        </div>
        <button
          type="submit"
          disabled={createArticleMutation.isLoading}
          className="px-6 py-2 bg-blue-600 text-white rounded-md"
        >
          {createArticleMutation.isLoading ? 'در حال ایجاد...' : 'ایجاد پیش‌نویس'}
        </button>
      </form>
    </div>
  );
}