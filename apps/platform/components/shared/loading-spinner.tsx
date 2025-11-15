// apps/platform/components/shared/loading-spinner.tsx

export function LoadingSpinner() {
  return (
    <div className="flex justify-center items-center py-12">
      <div
        className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"
        role="status"
        aria-label="در حال بارگذاری..."
      ></div>
    </div>
  );
}