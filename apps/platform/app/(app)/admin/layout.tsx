// apps/platform/app/(app)/admin/layout.tsx

'use client';

import { useAuth } from '@/lib/hooks/use-auth';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && (!user || !user.roles.includes('admin'))) {
      router.replace('/dashboard');
    }
  }, [user, loading, router]);

  if (loading || !user || !user.roles.includes('admin')) {
    return <div>درحال بررسی دسترسی...</div>;
  }

  return (
    <div className="flex">
      <aside className="w-64 bg-gray-800 text-white p-4">
        <h2 className="font-bold text-xl mb-8">پنل ادمین</h2>
        <nav className="space-y-2">
          <Link href="/admin/users" className="block p-2 rounded-md hover:bg-gray-700">مدیریت کاربران</Link>
          <Link href="/admin/moderation" className="block p-2 rounded-md hover:bg-gray-700">صف بررسی محتوا</Link>
          <Link href="/admin/analytics" className="block p-2 rounded-md hover:bg-gray-700">آمار و تحلیل</Link>
        </nav>
      </aside>
      <main className="flex-1 p-8">
        {children}
      </main>
    </div>
  );
}