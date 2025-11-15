// apps/platform/app/(content)/articles/[slug]/page.tsx

import { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { ArticleViewer } from '@/components/content/article-viewer';
import { getArticle, getAllArticleSlugs } from '@/lib/api/articles';

interface PageProps {
  params: { slug: string };
}

// Generate static params
export async function generateStaticParams() {
  const slugs = await getAllArticleSlugs();
  return slugs.map((slug) => ({ slug }));
}

// Generate metadata
export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  try {
    const article = await getArticle(params.slug);
    
    return {
      title: article.title.fa,
      description: article.excerpt.fa,
      openGraph: {
        title: article.title.fa,
        description: article.excerpt.fa,
        type: 'article',
        publishedTime: article.publishedAt.toISOString(),
        authors: [article.author.name],
        images: [
          {
            url: article.coverImage.url,
            alt: article.coverImage.alt,
          },
        ],
      },
    };
  } catch {
    return {
      title: 'مقاله یافت نشد',
    };
  }
}

// Page component
export default async function ArticlePage({ params }: PageProps) {
  let article;
  
  try {
    article = await getArticle(params.slug);
  } catch {
    notFound();
  }
  
  return <ArticleViewer article={article} />;
}

// Revalidate every hour
export const revalidate = 3600;