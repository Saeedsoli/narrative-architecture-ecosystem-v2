// apps/platform/components/content/podcast-player.tsx

'use client';

import { useState, useRef, useEffect } from 'react';

interface PodcastPlayerProps {
  src: string;
  title: string;
  author: string;
  coverImage: string;
}

export function PodcastPlayer({ src, title, author, coverImage }: PodcastPlayerProps) {
  const audioRef = useRef<HTMLAudioElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [progress, setProgress] = useState(0);
  const [duration, setDuration] = useState(0);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const updateProgress = () => {
      setProgress((audio.currentTime / audio.duration) * 100);
    };
    const setAudioDuration = () => {
      if (audio.duration && isFinite(audio.duration)) {
        setDuration(audio.duration);
      }
    };

    audio.addEventListener('timeupdate', updateProgress);
    audio.addEventListener('loadedmetadata', setAudioDuration);
    audio.addEventListener('ended', () => setIsPlaying(false));

    return () => {
      audio.removeEventListener('timeupdate', updateProgress);
      audio.removeEventListener('loadedmetadata', setAudioDuration);
      audio.removeEventListener('ended', () => setIsPlaying(false));
    };
  }, []);

  const togglePlayPause = () => {
    if (isPlaying) {
      audioRef.current?.pause();
    } else {
      audioRef.current?.play();
    }
    setIsPlaying(!isPlaying);
  };

  const formatTime = (time: number) => {
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60);
    return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
  };

  return (
    <div className="flex items-center p-4 bg-gray-100 dark:bg-gray-800 rounded-lg shadow-md">
      <img src={coverImage} alt={title} className="w-24 h-24 rounded-lg object-cover" />
      <div className="flex-1 mx-4">
        <h3 className="font-bold text-lg">{title}</h3>
        <p className="text-sm text-gray-500">{author}</p>
        <div className="mt-2">
          <audio ref={audioRef} src={src} preload="metadata" />
          <div className="w-full bg-gray-300 dark:bg-gray-600 rounded-full h-2">
            <div
              className="bg-blue-600 h-2 rounded-full"
              style={{ width: `${progress}%` }}
            ></div>
          </div>
          <div className="flex justify-between text-xs mt-1 text-gray-500">
            <span>{formatTime(audioRef.current?.currentTime || 0)}</span>
            <span>{formatTime(duration)}</span>
          </div>
        </div>
      </div>
      <button onClick={togglePlayPause} className="p-3 bg-blue-600 text-white rounded-full">
        {isPlaying ? (
          <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 20 20"><path d="M5 4h3v12H5V4zm7 0h3v12h-3V4z"></path></svg>
        ) : (
          <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 20 20"><path d="M4.555 15.168V4.832a.5.5 0 01.74-.447l8.361 5.168a.5.5 0 010 .894l-8.361 5.168a.5.5 0 01-.74-.447z"></path></svg>
        )}
      </button>
    </div>
  );
}