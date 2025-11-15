// packages/types/src/submission.ts

export interface Exercise {
  id: string;
  title: string;
  content: {
    description: string;
    // ... سایر فیلدهای محتوای تمرین
  };
}

export interface Submission {
  id: string;
  exerciseId: string;
  userId: string;
  answer: {
    text?: string;
    // ... سایر انواع پاسخ
  };
  status: 'pending' | 'graded' | 'reviewed';
  score?: number;
  feedback?: string;
  aiSummary?: string;
  submittedAt: string;
}