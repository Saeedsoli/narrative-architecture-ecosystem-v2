// apps/platform/app/(app)/community/forum/new/page.tsx

'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';

export default function NewTopicPage() {
  const router = useRouter();
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');

  const createTopicMutation = useMutation({
    mutationFn: (data: { title: string; body: string; locale: 'fa' }) => apiClient.post('/community/topics', data),
    onSuccess: (data) => {
      // به صفحه تاپیک جدید هدایت شو
      router.push(`/community/forum/${data.id}`);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createTopicMutation.mutate({ title, body, locale: 'fa' });
  };

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-8">ایجاد تاپیک جدید</h1>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="title">عنوان تاپیک</label>
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
          <label htmlFor="body">متن پست اول</label>
          <textarea
            id="body"
            value={body}
            onChange={(e) => setBody(e.target.value)}
            required
            rows={10}
            className="w-full mt-2"
          />
        </div>
        <button
          type="submit"
          disabled={createTopicMutation.isLoading}
          className="px-6 py-2 bg-blue-600 text-white rounded-md"
        >
          {createTopicMutation.isLoading ? 'در حال ایجاد...' : 'ایجاد تاپیک'}
        </button>
      </form>
    </div>
  );
}