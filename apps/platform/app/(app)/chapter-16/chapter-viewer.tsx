'use client';

import { Watermark } from './watermark';
import styles from './viewer.module.css';

interface ChapterViewerProps {
  content: {
    sections: Array<{
      id: string;
      title: string;
      content: string;
      // ... سایر فیلدهای تعاملی
    }>;
  };
  watermark: string;
}

export function ChapterViewer({ content, watermark }: ChapterViewerProps) {
  // جلوگیری از select و right-click در سمت کلاینت (یک لایه حفاظتی ساده)
  const handleContextMenu = (e: React.MouseEvent) => e.preventDefault();

  return (
    <div
      className={styles.viewer}
      onContextMenu={handleContextMenu}
    >
      <Watermark text={watermark} />
      
      <div className={styles.contentWrapper}>
        <h1>فصل ۱۶: هنر شکستن قوانین</h1>
        
        {content.sections.map((section) => (
          <section key={section.id} className={styles.section}>
            <h2>{section.title}</h2>
            {/* 
              استفاده از dangerouslySetInnerHTML با یک کتابخانه sanitize
              اگر محتوا HTML است. اگر Markdown است، از renderer استفاده کنید.
            */}
            <div dangerouslySetInnerHTML={{ __html: section.content }} />
          </section>
        ))}
      </div>
    </div>
  );
}