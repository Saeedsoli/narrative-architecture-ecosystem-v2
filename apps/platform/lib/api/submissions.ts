// apps/platform/lib/api/submissions.ts

import { apiClient } from './client';
import type { Submission, Exercise } from '@narrative-arch/types';

export const getExercise = async (id: string): Promise<Exercise> => {
  const { data } = await apiClient.get(`/exercises/${id}`);
  return data;
};

export const getSubmissions = async (exerciseId: string): Promise<Submission[]> => {
  const { data } = await apiClient.get<Submission[]>('/submissions', { params: { exerciseId } });
  return data;
};

export const createSubmission = async (data: { exerciseId: string; answer: { text: string } }): Promise<{ submissionId: string }> => {
  const { data: response } = await apiClient.post('/submissions', data);
  return response;
};

export const requestAnalysis = async (submissionId: string): Promise<{ message: string }> => {
  const { data } = await apiClient.post(`/submissions/${submissionId}/analyze`);
  return data;
};