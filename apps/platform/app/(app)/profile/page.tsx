// apps/platform/app/(app)/profile/page.tsx

'use client';

import { useAuth } from '@/lib/hooks/use-auth';
import Link from 'next/link';
import { LoadingSpinner } from '@/components/shared/loading-spinner';
import { Button } from '@/packages/ui/src/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/packages/ui/src/card';

export default function ProfilePage() {
  const { user, loading } = useAuth();

  if (loading) return <LoadingSpinner />;
  if (!user) return <div>لطفا وارد شوید.</div>;

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">پروفایل شما</h1>
        <Button asChild variant="outline">
          <Link href="/profile/edit">ویرایش پروفایل</Link>
        </Button>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{user.fullName}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <p><strong>نام کاربری:</strong> {user.username}</p>
          <p><strong>ایمیل:</strong> {user.email}</p>
          <p><strong>نقش‌ها:</strong> {user.roles.join(', ')}</p>
        </CardContent>
      </Card>
    </div>
  );
}