// apps/platform/app/(content)/podcasts/[slug]/page.tsx

import { getArticleBySlug } from '@/lib/api/articles';
import { PodcastPlayer } from '@/components/content/podcast-player';
import { MarkdownRenderer } from '@/components/content/markdown-renderer';

export default async function PodcastPage({ params }: { params: { slug: string } }) {
  const podcast = await getArticleBySlug(params.slug);

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-4xl font-bold mb-4">{podcast.title}</h1>
      <p className="text-lg text-gray-500 mb-8">{podcast.excerpt}</p>
      
      <PodcastPlayer
        src={podcast.mediaUrl}
        title={podcast.title}
        author={podcast.author.name}
        coverImage={podcast.coverImage.url}
      />
      
      <div className="mt-12">
        <h2 className="text-2xl font-bold mb-4">توضیحات</h2>
        <MarkdownRenderer content={podcast.content} />
      </div>
    </div>
  );
}