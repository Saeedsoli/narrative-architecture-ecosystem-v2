// apps/platform/app/layout.tsx

import { Metadata, Viewport } from 'next';
import { Vazirmatn, Noto_Serif, Fira_Code } from 'next/font/google';
import { Providers } from './providers';
import { CookieConsent } from '@/components/shared/cookie-consent';
import '@/styles/globals.css';

// Fonts
const vazirmatn = Vazirmatn({
  subsets: ['arabic'],
  variable: '--font-vazirmatn',
  display: 'swap',
});

const notoSerif = Noto_Serif({
  subsets: ['latin'],
  weight: ['400', '700'],
  variable: '--font-noto-serif',
  display: 'swap',
});

const firaCode = Fira_Code({
  subsets: ['latin'],
  variable: '--font-fira-code',
  display: 'swap',
});

// Metadata
export const metadata: Metadata = {
  title: {
    default: 'معماری روایت - راهنمای جامع طراحی و ساخت داستان',
    template: '%s | معماری روایت',
  },
  description: 'پلتفرم آموزش و تمرین نویسندگی با تمرکز بر معماری و ساختار روایت',
  keywords: [
    'نویسندگی',
    'داستان‌نویسی',
    'فیلمنامه',
    'پیرنگ',
    'شخصیت‌پردازی',
    'معماری روایت',
  ],
  authors: [{ name: 'معماری روایت' }],
  creator: 'معماری روایت',
  publisher: 'معماری روایت',
  
  // Open Graph
  openGraph: {
    type: 'website',
    locale: 'fa_IR',
    alternateLocale: 'en_US',
    url: 'https://narrative-arch.com',
    siteName: 'معماری روایت',
    title: 'معماری روایت - راهنمای جامع طراحی و ساخت داستان',
    description: 'پلتفرم آموزش و تمرین نویسندگی',
    images: [
      {
        url: 'https://cdn.narrative-arch.com/og-image.jpg',
        width: 1200,
        height: 630,
        alt: 'معماری روایت',
      },
    ],
  },
  
  // Twitter
  twitter: {
    card: 'summary_large_image',
    title: 'معماری روایت',
    description: 'پلتفرم آموزش و تمرین نویسندگی',
    images: ['https://cdn.narrative-arch.com/twitter-image.jpg'],
  },
  
  // Icons
  icons: {
    icon: '/favicon.ico',
    apple: '/apple-touch-icon.png',
  },
  
  // Manifest
  manifest: '/manifest.json',
  
  // Other
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-video-preview': -1,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
};

export const viewport: Viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 5,
  themeColor: [
    { media: '(prefers-color-scheme: light)', color: '#ffffff' },
    { media: '(prefers-color-scheme: dark)', color: '#1a1a1a' },
  ],
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html
      lang="fa"
      dir="rtl"
      className={`${vazirmatn.variable} ${notoSerif.variable} ${firaCode.variable}`}
      suppressHydrationWarning
    >
      <body className="font-sans antialiased bg-white dark:bg-gray-900 text-ink dark:text-gray-100">
        <Providers>
          {children}
          <CookieConsent />
        </Providers>
      </body>
    </html>
  );
}