// apps/platform/components/shared/persian-date.tsx

'use client';

interface PersianDateProps {
  date: string | Date;
  className?: string;
}

export function PersianDate({ date, className }: PersianDateProps) {
  try {
    const formattedDate = new Date(date).toLocaleDateString('fa-IR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
    return <span className={className}>{formattedDate}</span>;
  } catch (error) {
    return <span className={className}>تاریخ نامعتبر</span>;
  }
}