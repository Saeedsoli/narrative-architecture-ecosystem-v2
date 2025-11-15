// apps/platform/components/content/video-player.tsx

'use client';

import ReactPlayer from 'react-player/lazy';
import { LoadingSpinner } from '@/components/shared/loading-spinner';

interface VideoPlayerProps {
  url: string;
}

export function VideoPlayer({ url }: VideoPlayerProps) {
  return (
    <div className="relative aspect-video w-full rounded-lg overflow-hidden">
      <ReactPlayer
        url={url}
        width="100%"
        height="100%"
        controls={true}
        playing={false}
        light={true} // نمایش یک تصویر پیش‌نمایش برای بهبود عملکرد
        fallback={<LoadingSpinner />}
        className="absolute top-0 left-0"
      />
    </div>
  );
}