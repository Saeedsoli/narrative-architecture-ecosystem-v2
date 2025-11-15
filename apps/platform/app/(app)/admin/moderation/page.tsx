// apps/platform/app/(app)/admin/moderation/page.tsx

'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { useState } from 'react';

// ... (تعریف نوع داده ModerationItem)

const fetchModerationQueue = async (status: string) => {
  const { data } = await apiClient.get('/admin/moderation', { params: { status } });
  return data;
};

const moderateContent = async ({ itemId, action, reason }: { itemId: string, action: string, reason?: string }) => {
  await apiClient.post(`/admin/moderation/${itemId}`, { action, reason });
};

export default function ModerationPage() {
  const [status, setStatus] = useState('pending');
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['moderation-queue', status],
    queryFn: () => fetchModerationQueue(status),
  });

  const mutation = useMutation({
    mutationFn: moderateContent,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['moderation-queue'] });
    },
  });

  const handleModerate = (itemId: string, action: 'approve' | 'reject') => {
    const reason = action === 'reject' ? prompt('دلیل رد کردن را وارد کنید:') : undefined;
    mutation.mutate({ itemId, action, reason });
  };

  return (
    <div className="max-w-7xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-8">صف بررسی محتوا</h1>
      
      {/* Tabs for status */}
      <div className="flex space-x-4 mb-8">
        <button onClick={() => setStatus('pending')} className={status === 'pending' ? 'font-bold' : ''}>در انتظار</button>
        <button onClick={() => setStatus('approved')} className={status === 'approved' ? 'font-bold' : ''}>تایید شده</button>
        <button onClick={() => setStatus('rejected')} className={status === 'rejected' ? 'font-bold' : ''}>رد شده</button>
      </div>

      {isLoading && <div>در حال بارگذاری...</div>}
      
      {/* Table of items */}
      <table>
        {/* ... (Header) */}
        <tbody>
          {data?.data.map((item) => (
            <tr key={item.id}>
              <td>{item.target_type}</td>
              <td>{item.reason}</td>
              <td>
                <button onClick={() => handleModerate(item.id, 'approve')}>تایید</button>
                <button onClick={() => handleModerate(item.id, 'reject')}>رد</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}