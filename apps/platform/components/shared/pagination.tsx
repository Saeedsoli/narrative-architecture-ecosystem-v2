// apps/platform/components/shared/pagination.tsx

'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  baseUrl: string;
}

export function Pagination({ currentPage, totalPages, baseUrl }: PaginationProps) {
  const pathname = usePathname();
  if (totalPages <= 1) return null;

  const pages = Array.from({ length: totalPages }, (_, i) => i + 1);

  return (
    <nav className="flex justify-center items-center space-x-2 space-x-reverse">
      {currentPage > 1 && (
        <Link href={`${baseUrl}?page=${currentPage - 1}`} className="px-4 py-2 border rounded-md">
          قبلی
        </Link>
      )}

      {pages.map((page) => (
        <Link
          key={page}
          href={`${baseUrl}?page=${page}`}
          className={`px-4 py-2 border rounded-md ${
            page === currentPage
              ? 'bg-blue-600 text-white border-blue-600'
              : 'hover:bg-gray-100 dark:hover:bg-gray-800'
          }`}
        >
          {page}
        </Link>
      ))}

      {currentPage < totalPages && (
        <Link href={`${baseUrl}?page=${currentPage + 1}`} className="px-4 py-2 border rounded-md">
          بعدی
        </Link>
      )}
    </nav>
  );
}