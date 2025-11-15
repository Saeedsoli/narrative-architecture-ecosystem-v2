// apps/platform/components/content/markdown-renderer.tsx

'use client';

import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm'; // برای پشتیبانی از جداول و ...
import rehypeRaw from 'rehype-raw';   // برای رندر کردن HTML درون Markdown (با احتیاط)
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';

interface MarkdownRendererProps {
  content: string;
}

export function MarkdownRenderer({ content }: MarkdownRendererProps) {
  return (
    <ReactMarkdown
      className="prose prose-lg dark:prose-invert max-w-none"
      remarkPlugins={[remarkGfm]}
      rehypePlugins={[rehypeRaw]}
      components={{
        // بازنویسی کامپوننت کد برای افزودن Syntax Highlighting
        code({ node, inline, className, children, ...props }) {
          const match = /language-(\w+)/.exec(className || '');
          return !inline && match ? (
            <SyntaxHighlighter
              style={vscDarkPlus}
              language={match[1]}
              PreTag="div"
              {...props}
            >
              {String(children).replace(/\n$/, '')}
            </SyntaxHighlighter>
          ) : (
            <code className={className} {...props}>
              {children}
            </code>
          );
        },
        // بازنویسی کامپوننت جدول برای استایل‌دهی بهتر
        table({ children }) {
          return (
            <div className="overflow-x-auto">
              <table className="min-w-full">{children}</table>
            </div>
          );
        },
      }}
    >
      {content}
    </ReactMarkdown>
  );
}