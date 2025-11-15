// apps/platform/components/exercises/progress-tracker.tsx

'use client';

interface ProgressTrackerProps {
  label: string;
  percent: number; // A number between 0 and 100
}

export function ProgressTracker({ label, percent }: ProgressTrackerProps) {
  const normalizedPercent = Math.min(100, Math.max(0, percent));

  return (
    <div className="p-4 border rounded-lg">
      <div className="flex justify-between items-center mb-2">
        <span className="font-semibold">{label}</span>
        <span className="text-sm font-mono text-blue-600">{normalizedPercent.toFixed(0)}%</span>
      </div>
      <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-4">
        <div
          className="bg-blue-600 h-4 rounded-full transition-all duration-500"
          style={{ width: `${normalizedPercent}%` }}
        ></div>
      </div>
    </div>
  );
}