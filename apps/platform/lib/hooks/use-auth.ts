// apps/platform/app/(auth)/login/page.tsx
// نسخه کامل و نهایی - بدون نیاز به تغییر

'use client';

import { useState } from 'react';
import { useAuth } from '@/lib/hooks/use-auth';
import { useRouter } from 'next/navigation';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(email, password);
      router.push('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.error || 'ایمیل یا رمز عبور نامعتبر است.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="w-full max-w-md p-8 space-y-6 bg-white dark:bg-gray-800 shadow-lg rounded-lg">
        <h1 className="text-3xl font-bold text-center text-gray-900 dark:text-white">
          ورود به معماری روایت
        </h1>
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && <p className="text-red-500 text-center text-sm">{error}</p>}
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              ایمیل
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="example@email.com"
              required
              className="w-full px-4 py-2 mt-2 border rounded-md dark:bg-gray-700 dark:border-gray-600 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              رمز عبور
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
              className="w-full px-4 py-2 mt-2 border rounded-md dark:bg-gray-700 dark:border-gray-600 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full px-4 py-2 text-white font-semibold bg-blue-600 rounded-md hover:bg-blue-700 disabled:bg-blue-300 transition-colors"
          >
            {loading ? 'در حال ورود...' : 'ورود'}
          </button>
        </form>
      </div>
    </div>
  );
}