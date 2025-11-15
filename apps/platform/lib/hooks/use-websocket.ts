// apps/platform/lib/hooks/use-websocket.ts

'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { io, Socket } from 'socket.io-client';
import { useAuth } from './use-auth';

interface UseWebSocketOptions {
  url?: string;
  reconnection?: boolean;
  reconnectionAttempts?: number;
  reconnectionDelay?: number;
}

interface WebSocketState {
  connected: boolean;
  error: string | null;
}

export function useWebSocket(options: UseWebSocketOptions = {}) {
  const { user } = useAuth();
  const [state, setState] = useState<WebSocketState>({
    connected: false,
    error: null,
  });
  
  const socketRef = useRef<Socket | null>(null);
  const listenersRef = useRef<Map<string, Function>>(new Map());

  useEffect(() => {
    if (!user) return;

    const wsUrl = options.url || process.env.NEXT_PUBLIC_WS_URL || 'http://localhost:8080';
    
    // Create socket connection
    const socket = io(wsUrl, {
      auth: {
        token: localStorage.getItem('accessToken'),
      },
      reconnection: options.reconnection ?? true,
      reconnectionAttempts: options.reconnectionAttempts ?? 5,
      reconnectionDelay: options.reconnectionDelay ?? 1000,
      transports: ['websocket', 'polling'],
    });

    // Connection events
    socket.on('connect', () => {
      console.log('WebSocket connected');
      setState({ connected: true, error: null });
    });

    socket.on('disconnect', (reason) => {
      console.log('WebSocket disconnected:', reason);
      setState({ connected: false, error: reason });
    });

    socket.on('connect_error', (error) => {
      console.error('WebSocket connection error:', error);
      setState({ connected: false, error: error.message });
    });

    // Re-attach all listeners
    listenersRef.current.forEach((handler, event) => {
      socket.on(event, handler as any);
    });

    socketRef.current = socket;

    return () => {
      socket.disconnect();
      socketRef.current = null;
    };
  }, [user, options.url, options.reconnection, options.reconnectionAttempts, options.reconnectionDelay]);

  const on = useCallback((event: string, handler: Function) => {
    listenersRef.current.set(event, handler);
    socketRef.current?.on(event, handler as any);
  }, []);

  const off = useCallback((event: string) => {
    listenersRef.current.delete(event);
    socketRef.current?.off(event);
  }, []);

  const emit = useCallback((event: string, data?: any) => {
    socketRef.current?.emit(event, data);
  }, []);

  return {
    ...state,
    on,
    off,
    emit,
    socket: socketRef.current,
  };
}