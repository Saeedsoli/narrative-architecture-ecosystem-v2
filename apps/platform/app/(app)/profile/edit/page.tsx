// apps/platform/app/(app)/profile/edit/page.tsx

'use client';

import { useAuth } from '@/lib/hooks/use-auth';
import { useForm } from 'react-hook-form';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateUserProfile } from '@/lib/api/users';
import { Button } from '@/packages/ui/src/button';
import { Input } from '@/packages/ui/src/input';

export default function EditProfilePage() {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const { register, handleSubmit } = useForm({
    defaultValues: {
      fullName: user?.fullName,
      username: user?.username,
    },
  });

  const mutation = useMutation({
    mutationFn: updateUserProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth-user'] });
      alert('پروفایل شما با موفقیت به‌روز شد.');
    },
  });

  return (
    <div className="max-w-2xl mx-auto py-12 px-4">
      <h1 className="text-3xl font-bold mb-8">ویرایش پروفایل</h1>
      <form onSubmit={handleSubmit((data) => mutation.mutate(data))} className="space-y-6">
        <div>
          <label htmlFor="fullName">نام کامل</label>
          <Input id="fullName" {...register('fullName')} className="mt-2" />
        </div>
        <div>
          <label htmlFor="username">نام کاربری</label>
          <Input id="username" {...register('username')} className="mt-2" />
        </div>
        <Button type="submit" disabled={mutation.isLoading}>
          {mutation.isLoading ? 'در حال ذخیره...' : 'ذخیره تغییرات'}
        </Button>
      </form>
    </div>
  );
}