// apps/platform/lib/api/client.ts

import axios, { AxiosError, AxiosRequestConfig } from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const apiClient = axios.create({
  baseURL: API_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // مهم: برای ارسال کوکی‌ها (Refresh Token)
});

// Interceptor برای افزودن Access Token به هدر درخواست‌ها
apiClient.interceptors.request.use(
  (config) => {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('accessToken');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Interceptor برای مدیریت خطای 401 و تازه‌سازی توکن
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      try {
        const { data } = await apiClient.post('/auth/refresh');
        const { accessToken } = data;
        localStorage.setItem('accessToken', accessToken);
        apiClient.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
        originalRequest.headers!.Authorization = `Bearer ${accessToken}`;
        return apiClient(originalRequest);
      } catch (refreshError) {
        // اگر refresh هم شکست خورد، کاربر را logout کن
        if (typeof window !== 'undefined') {
          localStorage.removeItem('accessToken');
          window.location.href = '/auth/login';
        }
        return Promise.reject(refreshError);
      }
    }
    return Promise.reject(error);
  }
);