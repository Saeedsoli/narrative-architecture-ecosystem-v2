// apps/platform/tailwind.config.js

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    '../../packages/ui/src/**/*.{js,ts,jsx,tsx}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0f9ff',
          100: '#e0f2fe',
          200: '#bae6fd',
          300: '#7dd3fc',
          400: '#38bdf8',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
          800: '#075985',
          900: '#0c4a6e',
          950: '#082f49',
        },
        manuscript: {
          light: '#faf8f3',
          DEFAULT: '#f5e6d3',
          dark: '#e8d4b8',
        },
        ink: {
          light: '#4a4a4a',
          DEFAULT: '#1a1a1a',
          dark: '#000000',
        }
      },
      fontFamily: {
        sans: ['var(--font-vazirmatn)', 'system-ui', 'sans-serif'],
        serif: ['var(--font-noto-serif)', 'Georgia', 'serif'],
        mono: ['var(--font-fira-code)', 'Courier New', 'monospace'],
      },
      typography: (theme) => ({
        DEFAULT: {
          css: {
            direction: 'rtl',
            textAlign: 'right',
            maxWidth: '75ch',
            color: theme('colors.ink.DEFAULT'),
            lineHeight: '2',
            
            // Links
            a: {
              color: theme('colors.primary.600'),
              textDecoration: 'underline',
              fontWeight: '500',
              '&:hover': {
                color: theme('colors.primary.700'),
              },
            },
            
            // Headings
            'h1, h2, h3, h4': {
              fontFamily: theme('fontFamily.serif'),
              fontWeight: '700',
              color: theme('colors.ink.DEFAULT'),
            },
            
            h1: {
              fontSize: '2.25rem',
              marginTop: '0',
              marginBottom: '1.5rem',
            },
            
            h2: {
              fontSize: '1.875rem',
              marginTop: '2.5rem',
              marginBottom: '1.25rem',
            },
            
            // Blockquotes
            blockquote: {
              fontStyle: 'normal',
              color: theme('colors.ink.light'),
              borderRightWidth: '0.25rem',
              borderRightColor: theme('colors.primary.500'),
              paddingRight: '1.5rem',
              marginRight: '0',
              quotes: '"\\201C""\\201D""\\2018""\\2019"',
            },
            
            // Code blocks
            'code::before': {
              content: '""',
            },
            'code::after': {
              content: '""',
            },
            code: {
              fontFamily: theme('fontFamily.mono'),
              backgroundColor: theme('colors.gray.100'),
              padding: '0.25rem 0.5rem',
              borderRadius: '0.25rem',
              fontSize: '0.875em',
            },
            
            pre: {
              direction: 'ltr',
              textAlign: 'left',
              backgroundColor: theme('colors.gray.900'),
              color: theme('colors.gray.100'),
              borderRadius: '0.5rem',
              padding: '1.25rem 1.5rem',
              overflow: 'auto',
            },
            
            'pre code': {
              backgroundColor: 'transparent',
              padding: '0',
            },
          },
        },
        
        // Dark mode
        invert: {
          css: {
            color: theme('colors.gray.100'),
            a: {
              color: theme('colors.primary.400'),
            },
            'h1, h2, h3, h4': {
              color: theme('colors.gray.100'),
            },
            blockquote: {
              color: theme('colors.gray.400'),
              borderRightColor: theme('colors.primary.400'),
            },
            code: {
              backgroundColor: theme('colors.gray.800'),
              color: theme('colors.gray.100'),
            },
          },
        },
      }),
      
      // Animation
      animation: {
        'fade-in': 'fadeIn 0.3s ease-in-out',
        'slide-up': 'slideUp 0.4s ease-out',
        'slide-down': 'slideDown 0.4s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        slideDown: {
          '0%': { transform: 'translateY(-10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
      
      // Spacing for Persian typography
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
    require('@tailwindcss/aspect-ratio'),
  ],
};