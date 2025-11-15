// apps/platform/lib/api/auth.ts

import { get, post } from './client';
import type { User } from '@narrative-arch/types';

interface LoginRequest {
  email: string;
  password: string;
}

interface LoginResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}

interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  fullName: string;
}

interface RegisterResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}

interface RefreshTokenResponse {
  accessToken: string;
  refreshToken?: string;
}

export async function login(data: LoginRequest): Promise<LoginResponse> {
  return post<LoginResponse>('/api/v1/auth/login', data);
}

export async function register(data: RegisterRequest): Promise<RegisterResponse> {
  return post<RegisterResponse>('/api/v1/auth/register', data);
}

export async function logout(): Promise<void> {
  return post<void>('/api/v1/auth/logout');
}

export async function getCurrentUser(): Promise<User> {
  return get<User>('/api/v1/auth/me');
}

export async function refreshToken(): Promise<RefreshTokenResponse> {
  return post<RefreshTokenResponse>('/api/v1/auth/refresh', {}, {
    withCredentials: true,
  });
}

export async function forgotPassword(email: string): Promise<void> {
  return post<void>('/api/v1/auth/forgot-password', { email });
}

export async function resetPassword(token: string, password: string): Promise<void> {
  return post<void>('/api/v1/auth/reset-password', { token, password });
}