# Platform Frontend

این اپلیکیشن، رابط کاربری اصلی پلتفرم است که با Next.js و TypeScript ساخته شده است.

## راه‌اندازی

1.  از ریشه پروژه، دستور `pnpm install` را اجرا کنید.
2.  فایل `.env.local` را بر اساس `.env.example` ایجاد کنید.
3.  برای اجرای سرور توسعه، از ریشه پروژه دستور `pnpm run dev` را اجرا کنید.

## افزودن یک صفحه جدید (مثال: `/bookmarks`)

1.  **API Client:**
    *   در `lib/api/bookmarks.ts`، یک تابع `listBookmarks` برای فراخوانی `GET /api/v1/bookmarks` ایجاد کنید.
2.  **React Query:**
    *   از `useQuery` برای fetch کردن داده‌ها با استفاده از `listBookmarks` استفاده کنید.
3.  **Page Component:**
    *   فایل `app/(app)/bookmarks/page.tsx` را ایجاد کنید.
    *   از هوک `useQuery` برای دریافت داده‌ها استفاده کرده و آن‌ها را نمایش دهید.
4.  **Routing:**
    *   یک `Link` از داشبورد به `/bookmarks` اضافه کنید.