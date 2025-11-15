// apps/platform/components/content/article-card.tsx

import Link from 'next/link';
import Image from 'next/image';
import type { Article } from '@narrative-arch/types';
import { PersianDate } from '@/components/shared/persian-date';

interface ArticleCardProps {
  article: Article;
}

export function ArticleCard({ article }: ArticleCardProps) {
  return (
    <Link href={`/articles/${article.slug}`} className="group block">
      <div className="overflow-hidden rounded-lg">
        <Image
          src={article.coverImage.url}
          alt={article.coverImage.alt || article.title}
          width={400}
          height={225}
          className="w-full h-48 object-cover transition-transform duration-300 group-hover:scale-105"
        />
      </div>
      <div className="mt-4">
        <p className="text-sm text-blue-500 font-semibold">{article.metadata.category}</p>
        <h3 className="mt-2 text-xl font-bold group-hover:text-blue-600 transition-colors">
          {article.title}
        </h3>
        <p className="mt-2 text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
          {article.excerpt}
        </p>
        <div className="mt-4 flex items-center justify-between text-xs text-gray-500">
          <span>{article.author.name}</span>
          <PersianDate date={article.publishedAt} />
        </div>
      </div>
    </Link>
  );
}