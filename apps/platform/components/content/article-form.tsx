// apps/platform/components/content/article-form.tsx

'use client';

import { useState, useEffect } from 'react';
import type { Article } from '@narrative-arch/types';

interface ArticleFormProps {
  initialData?: Partial<Article>;
  onSubmit: (data: Partial<Article>) => void;
  isSubmitting: boolean;
}

export function ArticleForm({ initialData, onSubmit, isSubmitting }: ArticleFormProps) {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  
  useEffect(() => {
    if (initialData) {
      setTitle(initialData.title || '');
      setContent(initialData.content || '');
    }
  }, [initialData]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ title, content });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label htmlFor="title">عنوان</label>
        <input
          id="title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          className="w-full mt-2 p-2 border rounded-md"
        />
      </div>
      <div>
        <label htmlFor="content">محتوا (Markdown)</label>
        <textarea
          id="content"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          required
          rows={20}
          className="w-full mt-2 p-2 border rounded-md"
        />
      </div>
      <button
        type="submit"
        disabled={isSubmitting}
        className="px-6 py-2 bg-blue-600 text-white rounded-md"
      >
        {isSubmitting ? 'در حال ذخیره...' : 'ذخیره'}
      </button>
    </form>
  );
}