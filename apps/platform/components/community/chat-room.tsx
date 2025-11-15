// apps/platform/components/community/chat-room.tsx

'use client';

import { useState, useEffect, useRef } from 'react';
import { useAuth } from '@/lib/hooks/use-auth';
import { useWebSocket } from '@/lib/hooks/use-websocket';
import { ChatMessage } from './chat-message';
import { VoiceMessageRecorder } from './voice-message-recorder';
import type { Message } from '@narrative-arch/types';

interface ChatRoomProps {
  roomId: string;
}

export function ChatRoom({ roomId }: ChatRoomProps) {
  const { user } = useAuth();
  const { connected, on, off, emit } = useWebSocket();
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [isTyping, setIsTyping] = useState<string[]>([]);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const typingTimeoutRef = useRef<NodeJS.Timeout>();

  useEffect(() => {
    if (!connected) return;

    // Join room
    emit('join_room', { roomId });

    // Listen for new messages
    on('new_message', (message: Message) => {
      setMessages((prev) => [...prev, message]);
    });

    // Listen for typing indicators
    on('user_typing', ({ userId, username }: { userId: string; username: string }) => {
      setIsTyping((prev) => {
        if (prev.includes(username)) return prev;
        return [...prev, username];
      });

      setTimeout(() => {
        setIsTyping((prev) => prev.filter((u) => u !== username));
      }, 3000);
    });

    // Listen for user stopped typing
    on('user_stop_typing', ({ username }: { username: string }) => {
      setIsTyping((prev) => prev.filter((u) => u !== username));
    });

    return () => {
      emit('leave_room', { roomId });
      off('new_message');
      off('user_typing');
      off('user_stop_typing');
    };
  }, [connected, roomId, on, off, emit]);

  // Auto-scroll to bottom
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!inputValue.trim()) return;

    emit('send_message', {
      roomId,
      text: inputValue,
    });

    setInputValue('');
    
    // Stop typing indicator
    emit('stop_typing', { roomId });
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);

    // Emit typing indicator
    emit('typing', { roomId });

    // Clear previous timeout
    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }

    // Stop typing after 2 seconds of inactivity
    typingTimeoutRef.current = setTimeout(() => {
      emit('stop_typing', { roomId });
    }, 2000);
  };

  const handleVoiceMessageSent = (audioBlob: Blob) => {
    // Upload voice message
    const formData = new FormData();
    formData.append('audio', audioBlob, 'voice-message.webm');
    formData.append('roomId', roomId);

    fetch('/api/voice-messages/upload', {
      method: 'POST',
      body: formData,
    })
      .then((res) => res.json())
      .then((data) => {
        emit('send_voice_message', {
          roomId,
          voiceMessageUrl: data.url,
          duration: data.duration,
        });
      })
      .catch((error) => {
        console.error('Failed to upload voice message:', error);
      });
  };

  if (!user) {
    return (
      <div className="flex items-center justify-center h-full">
        <p>لطفاً ابتدا وارد شوید</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b">
        <h2 className="text-xl font-bold">اتاق گفتگو</h2>
        <div className="flex items-center gap-2">
          <div
            className={`w-3 h-3 rounded-full ${connected ? 'bg-green-500' : 'bg-red-500'}`}
          />
          <span className="text-sm text-gray-600">
            {connected ? 'متصل' : 'قطع شده'}
          </span>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message) => (
          <ChatMessage
            key={message.id}
            message={message}
            isOwn={message.userId === user.id}
          />
        ))}

        {/* Typing indicator */}
        {isTyping.length > 0 && (
          <div className="text-sm text-gray-500">
            {isTyping.join(', ')} در حال نوشتن...
          </div>
        )}

        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <div className="p-4 border-t">
        <form onSubmit={handleSendMessage} className="flex gap-2">
          <input
            type="text"
            value={inputValue}
            onChange={handleInputChange}
            placeholder="پیام خود را بنویسید..."
            className="flex-1 px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
          
          <VoiceMessageRecorder onSend={handleVoiceMessageSent} />
          
          <button
            type="submit"
            disabled={!inputValue.trim()}
            className="px-6 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600 disabled:opacity-50 disabled:cursor-not-allowed transition"
          >
            ارسال
          </button>
        </form>
      </div>
    </div>
  );
}