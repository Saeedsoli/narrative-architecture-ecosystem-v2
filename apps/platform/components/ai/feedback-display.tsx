// apps/platform/components/ai/feedback-display.tsx

'use client';

interface FeedbackDisplayProps {
  analysis: string; // فرض می‌کنیم تحلیل یک متن Markdown است
}

export function FeedbackDisplay({ analysis }: FeedbackDisplayProps) {
  // در آینده، می‌توانیم پاسخ AI را به‌صورت JSON دریافت کرده و به‌صورت ساختاریافته نمایش دهیم
  return (
    <div className="p-6 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
      <h3 className="text-xl font-bold mb-4 text-blue-800 dark:text-blue-200">تحلیل هوش مصنوعی</h3>
      <div className="prose prose-sm dark:prose-invert max-w-none">
        {analysis}
      </div>
    </div>
  );
}