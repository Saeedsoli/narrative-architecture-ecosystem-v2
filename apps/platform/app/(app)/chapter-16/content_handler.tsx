import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { ChapterViewer } from '@/components/chapter16/chapter-viewer';
import { getApiServer } from '@/lib/api/server'; // یک کلاینت API برای سمت سرور

async function getChapter16Content() {
  try {
    const api = getApiServer(cookies());
    const data = await api.get('/content/chapter-16');
    return data;
  } catch (error) {
    console.error("Failed to fetch Chapter 16:", error);
    return null;
  }
}

export default async function Chapter16Page() {
  const data = await getChapter16Content();

  // اگر کاربر دسترسی نداشته باشد، سرور خطا برمی‌گرداند و ما او را redirect می‌کنیم
  if (!data) {
    redirect('/shop?reason=no-access-to-chapter-16');
  }

  // داده‌های رمزگشایی شده فقط در سمت سرور وجود دارند و به کامپوننت ارسال می‌شوند
  return (
    <div className="chapter-16-container">
      <ChapterViewer content={data.content} watermark={data.watermark} />
    </div>
  );
}