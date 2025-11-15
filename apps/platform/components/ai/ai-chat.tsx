// apps/platform/components/ai/ai-chat.tsx

'use client';

import { useState, useRef, useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { MarkdownRenderer } from '@/components/content/markdown-renderer';
import { Button } from '@/packages/ui/src/button';
import { Input } from '@/packages/ui/src/input';

interface Message {
  role: 'user' | 'assistant';
  content: string;
}

const analyzeWithAI = async (text: string): Promise<string> => {
  const { data } = await apiClient.post('/ai/analyze', { text });
  return data.analysis;
};

export function AIChat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const mutation = useMutation({
    mutationFn: analyzeWithAI,
    onSuccess: (analysis) => {
      setMessages((prev) => [...prev, { role: 'assistant', content: analysis }]);
    },
    onError: (error: any) => {
      setMessages((prev) => [...prev, { role: 'assistant', content: `خطا: ${error.message}` }]);
    },
  });

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, mutation.isLoading]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim() || mutation.isLoading) return;
    
    const userMessage: Message = { role: 'user', content: input };
    setMessages((prev) => [...prev, userMessage]);
    mutation.mutate(input);
    setInput('');
  };

  return (
    <div className="flex flex-col h-[70vh] border rounded-lg bg-card">
      <div className="flex-1 p-4 overflow-y-auto space-y-6">
        {messages.map((msg, index) => (
          <div key={index} className={`flex gap-3 ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
            {msg.role === 'assistant' && <span className="text-xl">🤖</span>}
            <div
              className={`max-w-xl p-3 rounded-lg ${
                msg.role === 'user'
                  ? 'bg-primary text-primary-foreground'
                  : 'bg-secondary'
              }`}
            >
              <MarkdownRenderer content={msg.content} />
            </div>
          </div>
        ))}
        {mutation.isLoading && (
          <div className="flex gap-3 justify-start">
            <span className="text-xl">🤖</span>
            <div className="max-w-xl p-3 rounded-lg bg-secondary animate-pulse">
              در حال فکر کردن...
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>
      <form onSubmit={handleSubmit} className="p-4 border-t flex gap-2">
        <Input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="متن خود را برای تحلیل وارد کنید..."
          className="flex-1"
          disabled={mutation.isLoading}
        />
        <Button type="submit" disabled={mutation.isLoading}>
          ارسال
        </Button>
      </form>
    </div>
  );
}