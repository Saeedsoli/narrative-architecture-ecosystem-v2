// apps/platform/app/(app)/dashboard/page.tsx

'use client';

import { useAuth } from '@/lib/hooks/use-auth';
import Link from 'next/link';

export default function DashboardPage() {
  const { user, loading, logout } = useAuth();

  if (loading) {
    return <div>در حال بارگذاری اطلاعات...</div>;
  }

  if (!user) {
    // این حالت معمولاً با یک Middleware مدیریت می‌شود، اما به‌عنوان یک لایه حفاظتی اضافه شده
    return <div>لطفاً برای دسترسی به این صفحه وارد شوید.</div>;
  }

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">داشبورد</h1>
        <button
          onClick={logout}
          className="px-4 py-2 text-sm text-white bg-red-600 rounded-md"
        >
          خروج
        </button>
      </div>

      <div className="p-6 bg-gray-100 dark:bg-gray-800 rounded-lg">
        <h2 className="text-2xl font-semibold mb-2">سلام، {user.fullName}!</h2>
        <p className="text-gray-600 dark:text-gray-400">به پلتفرم معماری روایت خوش آمدید.</p>
      </div>

      <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
        <Link href="/articles/new" className="block p-6 border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800">
          <h3 className="font-semibold text-lg">نوشتن مقاله جدید</h3>
          <p className="text-sm text-gray-500 mt-2">اولین پیش‌نویس خود را ایجاد کنید.</p>
        </Link>
        <Link href="/exercises" className="block p-6 border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800">
          <h3 className="font-semibold text-lg">مشاهده تمرینات</h3>
          <p className="text-sm text-gray-500 mt-2">مهارت‌های خود را به چالش بکشید.</p>
        </Link>
        <Link href="/community/forum" className="block p-6 border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800">
          <h3 className="font-semibold text-lg">ورود به انجمن</h3>
          <p className="text-sm text-gray-500 mt-2">با سایر نویسندگان گفتگو کنید.</p>
        </Link>
        <Link href="/profile/edit" className="block p-6 border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800">
          <h3 className="font-semibold text-lg">ویرایش پروفایل</h3>
          <p className="text-sm text-gray-500 mt-2">اطلاعات کاربری خود را به‌روز کنید.</p>
        </Link>
      </div>
    </div>
  );
}