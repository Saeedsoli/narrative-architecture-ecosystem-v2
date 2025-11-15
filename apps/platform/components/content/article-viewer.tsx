// apps/platform/components/content/article-viewer.tsx

'use client';

import { useState } from 'react';
import Image from 'next/image';
import { MarkdownRenderer } from '@/components/shared/markdown-renderer';
import { PersianDate } from '@/components/shared/persian-date';
import { ShareButtons } from '@/components/shared/share-buttons';
import { RelatedArticles } from './related-articles';
import { CommentsSection } from './comments-section';
import type { Article } from '@narrative-arch/types';

interface ArticleViewerProps {
  article: Article;
}

export function ArticleViewer({ article }: ArticleViewerProps) {
  const [isBookmarked, setIsBookmarked] = useState(false);

  return (
    <article className="max-w-4xl mx-auto px-4 py-12">
      {/* Header */}
      <header className="mb-12">
        {/* Cover Image */}
        {article.coverImage && (
          <div className="relative aspect-video mb-8 rounded-xl overflow-hidden">
            <Image
              src={article.coverImage.url}
              alt={article.coverImage.alt}
              fill
              className="object-cover"
              priority
              sizes="(max-width: 768px) 100vw, (max-width: 1200px) 80vw, 1200px"
            />
          </div>
        )}
        
        {/* Title */}
        <h1 className="text-4xl md:text-5xl font-serif font-bold text-ink dark:text-white mb-6">
          {article.title.fa}
        </h1>
        
        {/* Metadata */}
        <div className="flex items-center gap-6 text-sm text-gray-600 dark:text-gray-400 mb-6">
          <div className="flex items-center gap-2">
            <Image
              src={article.author.avatar}
              alt={article.author.name}
              width={40}
              height={40}
              className="rounded-full"
            />
            <span>{article.author.name}</span>
          </div>
          
          <PersianDate date={article.publishedAt} />
          
          <span>{article.readingTime} دقیقه مطالعه</span>
        </div>
        
        {/* Tags */}
        <div className="flex flex-wrap gap-2 mb-6">
          {article.tags.map((tag) => (
            <span
              key={tag}
              className="px-3 py-1 bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 rounded-full text-sm"
            >
              #{tag}
            </span>
          ))}
        </div>
        
        {/* Actions */}
        <div className="flex items-center gap-4">
          <ShareButtons
            url={`https://narrative-arch.com/articles/${article.slug}`}
            title={article.title.fa}
          />
          
          <button
            onClick={() => setIsBookmarked(!isBookmarked)}
            className="flex items-center gap-2 px-4 py-2 rounded-lg border hover:bg-gray-50 dark:hover:bg-gray-800 transition"
          >
            <svg
              className={`w-5 h-5 ${isBookmarked ? 'fill-primary-500' : 'fill-none'}`}
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"
              />
            </svg>
            <span>{isBookmarked ? 'ذخیره شد' : 'ذخیره'}</span>
          </button>
        </div>
      </header>
      
      {/* Content */}
      <div className="prose prose-lg dark:prose-invert max-w-none mb-12">
        <MarkdownRenderer content={article.content.fa} />
      </div>
      
      {/* Related Articles */}
      <RelatedArticles articleId={article.id} />
      
      {/* Comments */}
      <CommentsSection articleId={article.id} />
    </article>
  );
}