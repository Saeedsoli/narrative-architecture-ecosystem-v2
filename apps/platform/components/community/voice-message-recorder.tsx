// apps/platform/components/community/voice-message-recorder.tsx

'use client';

import { useVoiceRecorder } from '@/lib/hooks/use-voice-recorder';
import { formatDuration } from '@/lib/utils/format';

interface VoiceMessageRecorderProps {
  onSend: (audioBlob: Blob) => void;
}

export function VoiceMessageRecorder({ onSend }: VoiceMessageRecorderProps) {
  const {
    isRecording,
    duration,
    audioBlob,
    startRecording,
    stopRecording,
    cancelRecording,
  } = useVoiceRecorder();

  const handleSend = () => {
    if (audioBlob) {
      onSend(audioBlob);
      cancelRecording();
    }
  };

  if (isRecording) {
    return (
      <div className="flex items-center gap-2 px-4 py-2 bg-red-50 border border-red-200 rounded-lg">
        <div className="w-3 h-3 bg-red-500 rounded-full animate-pulse" />
        <span className="font-mono text-sm">{formatDuration(duration)}</span>
        
        <button
          onClick={stopRecording}
          className="px-3 py-1 bg-red-500 text-white rounded-lg text-sm hover:bg-red-600 transition"
        >
          توقف
        </button>
        
        <button
          onClick={cancelRecording}
          className="px-3 py-1 bg-gray-500 text-white rounded-lg text-sm hover:bg-gray-600 transition"
        >
          لغو
        </button>
      </div>
    );
  }

  if (audioBlob) {
    return (
      <div className="flex items-center gap-2 px-4 py-2 bg-green-50 border border-green-200 rounded-lg">
        <audio src={URL.createObjectURL(audioBlob)} controls className="h-10" />
        
        <button
          onClick={handleSend}
          className="px-3 py-1 bg-green-500 text-white rounded-lg text-sm hover:bg-green-600 transition"
        >
          ارسال
        </button>
        
        <button
          onClick={cancelRecording}
          className="px-3 py-1 bg-gray-500 text-white rounded-lg text-sm hover:bg-gray-600 transition"
        >
          حذف
        </button>
      </div>
    );
  }

  return (
    <button
      onClick={startRecording}
      className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition"
      title="ضبط پیام صوتی"
    >
      <svg
        className="w-6 h-6 text-gray-600 dark:text-gray-400"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"
        />
      </svg>
    </button>
  );
}