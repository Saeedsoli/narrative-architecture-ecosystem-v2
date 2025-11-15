// apps/platform/components/community/post-editor.tsx

'use client';

import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createPost } from '@/lib/api/community';

interface PostEditorProps {
  topicId: string;
}

export function PostEditor({ topicId }: PostEditorProps) {
  const [body, setBody] = useState('');
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (newPostBody: string) => createPost(topicId, newPostBody),
    onSuccess: () => {
      setBody('');
      queryClient.invalidateQueries({ queryKey: ['forum-posts', topicId] });
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!body.trim()) return;
    mutation.mutate(body);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <textarea
        value={body}
        onChange={(e) => setBody(e.target.value)}
        rows={5}
        placeholder="پاسخ خود را بنویسید..."
        className="w-full p-2 border rounded-md"
        required
      />
      <button
        type="submit"
        disabled={mutation.isLoading}
        className="px-6 py-2 bg-blue-600 text-white rounded-md"
      >
        {mutation.isLoading ? 'در حال ارسال...' : 'ارسال پاسخ'}
      </button>
    </form>
  );
}