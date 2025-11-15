// apps/platform/app/(app)/admin/users/page.tsx

'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import type { User } from '@narrative-arch/types';

const fetchUsers = async (page: number = 1): Promise<{ data: User[], total: number }> => {
  const { data } = await apiClient.get('/admin/users', { params: { page } });
  return data;
};

const updateUserStatus = async ({ userId, status }: { userId: string, status: string }) => {
  await apiClient.put(`/admin/users/${userId}/status`, { status });
};

export default function AdminUsersPage() {
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ['admin-users'],
    queryFn: () => fetchUsers(),
  });

  const mutation = useMutation({
    mutationFn: updateUserStatus,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] });
    },
  });

  if (isLoading) return <div>در حال بارگذاری کاربران...</div>;

  return (
    <div className="max-w-7xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-8">مدیریت کاربران</h1>
      <table className="w-full text-right">
        <thead>
          <tr>
            <th>نام کامل</th>
            <th>ایمیل</th>
            <th>نام کاربری</th>
            <th>وضعیت</th>
            <th>عملیات</th>
          </tr>
        </thead>
        <tbody>
          {data?.data.map((user) => (
            <tr key={user.id}>
              <td>{user.fullName}</td>
              <td>{user.email}</td>
              <td>{user.username}</td>
              <td>{user.status}</td>
              <td>
                <select
                  defaultValue={user.status}
                  onChange={(e) => mutation.mutate({ userId: user.id, status: e.target.value })}
                >
                  <option value="active">فعال</option>
                  <option value="suspended">معلق</option>
                  <option value="deleted">حذف شده</option>
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}