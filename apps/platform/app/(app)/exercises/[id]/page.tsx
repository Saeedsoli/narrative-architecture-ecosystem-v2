// apps/platform/app/(app)/exercises/[id]/page.tsx

'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useParams } from 'next/navigation';
import { getExercise, getSubmissions, createSubmission, requestAnalysis } from '@/lib/api/submissions';
import { LoadingSpinner } from '@/components/shared/loading-spinner';
import type { Submission, Exercise } from '@narrative-arch/types';

export default function ExercisePage() {
  const params = useParams();
  const exerciseId = params.id as string;
  const queryClient = useQueryClient();
  const [answer, setAnswer] = useState('');

  const { data: exercise, isLoading: isLoadingExercise } = useQuery<Exercise>({
    queryKey: ['exercise', exerciseId],
    queryFn: () => getExercise(exerciseId),
    enabled: !!exerciseId,
  });

  const { data: submissions, isLoading: isLoadingSubmissions } = useQuery<Submission[]>({
    queryKey: ['submissions', exerciseId],
    queryFn: () => getSubmissions(exerciseId),
    enabled: !!exerciseId,
  });

  const submissionMutation = useMutation({
    mutationFn: createSubmission,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['submissions', exerciseId] });
      setAnswer('');
      alert('پاسخ شما با موفقیت ثبت شد.');
    },
  });

  const analysisMutation = useMutation({
    mutationFn: requestAnalysis,
    onSuccess: () => {
      alert('درخواست تحلیل ارسال شد. نتیجه به‌زودی نمایش داده می‌شود.');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!answer.trim()) return;
    submissionMutation.mutate({ exerciseId, answer: { text: answer } });
  };
  
  if (isLoadingExercise) return <LoadingSpinner />;
  if (!exercise) return <div>تمرین یافت نشد.</div>;

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-4">{exercise.title}</h1>
      <p className="text-gray-600 mb-8">{exercise.content.description}</p>

      <form onSubmit={handleSubmit} className="space-y-4 mb-12">
        <textarea
          value={answer}
          onChange={(e) => setAnswer(e.target.value)}
          rows={10}
          placeholder="پاسخ خود را اینجا بنویسید..."
          className="w-full p-4 border rounded-md"
          required
        />
        <button
          type="submit"
          disabled={submissionMutation.isLoading}
          className="px-6 py-2 bg-blue-600 text-white rounded-md"
        >
          {submissionMutation.isLoading ? 'در حال ارسال...' : 'ارسال پاسخ'}
        </button>
      </form>

      <div className="space-y-8">
        <h2 className="text-2xl font-bold">ارسال‌های قبلی شما</h2>
        {isLoadingSubmissions ? <LoadingSpinner /> : (
          submissions && submissions.length > 0 ? (
            submissions.map((sub) => (
              <div key={sub.id} className="p-6 border rounded-lg">
                <p className="text-gray-500 text-sm mb-2">ارسال شده در: {new Date(sub.submittedAt).toLocaleString('fa-IR')}</p>
                <p className="whitespace-pre-wrap mb-4">{sub.answer.text}</p>
                
                <div className="p-4 bg-gray-50 rounded-md">
                  <h4 className="font-semibold mb-2">تحلیل هوش مصنوعی:</h4>
                  {sub.aiSummary ? (
                    <p className="text-gray-700">{sub.aiSummary}</p>
                  ) : (
                    <div className="flex items-center gap-4">
                      <p className="text-gray-500">در انتظار تحلیل...</p>
                      <button 
                        onClick={() => analysisMutation.mutate(sub.id)}
                        disabled={analysisMutation.isLoading}
                        className="text-sm text-blue-600"
                      >
                        درخواست تحلیل
                      </button>
                    </div>
                  )}
                </div>
              </div>
            ))
          ) : (
            <p>شما هنوز پاسخی برای این تمرین ثبت نکرده‌اید.</p>
          )
        )}
      </div>
    </div>
  );
}