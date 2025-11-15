// apps/platform/components/exercises/essay-question.tsx

'use client';

import { useState } from 'react';

interface EssayQuestionProps {
  question: string;
  placeholder?: string;
  onSubmit: (answer: string) => void;
  isSubmitting: boolean;
}

export function EssayQuestion({ question, placeholder, onSubmit, isSubmitting }: EssayQuestionProps) {
  const [answer, setAnswer] = useState('');
  const wordCount = answer.trim().split(/\s+/).filter(Boolean).length;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!answer.trim()) return;
    onSubmit(answer);
  };

  return (
    <div className="p-6 border rounded-lg">
      <p className="font-semibold mb-4">{question}</p>
      <form onSubmit={handleSubmit}>
        <textarea
          value={answer}
          onChange={(e) => setAnswer(e.target.value)}
          rows={15}
          placeholder={placeholder || "پاسخ خود را اینجا بنویسید..."}
          className="w-full p-4 border rounded-md"
          required
        />
        <div className="flex justify-between items-center mt-4">
          <span className="text-sm text-gray-500">تعداد کلمات: {wordCount}</span>
          <button
            type="submit"
            disabled={isSubmitting || wordCount === 0}
            className="px-6 py-2 bg-blue-600 text-white rounded-md disabled:opacity-50"
          >
            {isSubmitting ? 'در حال ارسال...' : 'ارسال پاسخ'}
          </button>
        </div>
      </form>
    </div>
  );
}