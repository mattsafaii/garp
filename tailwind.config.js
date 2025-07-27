/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./site/**/*.{html,js,md}",
    "./site/docs/**/*.{html,md}",
    "./site/docs/markdown/**/*.md",
  ],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['system-ui', '-apple-system', 'sans-serif'],
        'mono': ['SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'monospace'],
      },
      colors: {
        'garp': {
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
        }
      },
      typography: {
        DEFAULT: {
          css: {
            maxWidth: 'none',
            color: '#374151',
            a: {
              color: '#2563eb',
              '&:hover': {
                color: '#1d4ed8',
              },
            },
            'h1, h2, h3, h4': {
              color: '#111827',
            },
            code: {
              color: '#374151',
              backgroundColor: '#f3f4f6',
              padding: '0.125rem 0.25rem',
              borderRadius: '0.25rem',
              fontWeight: '400',
            },
            'code::before': {
              content: '""',
            },
            'code::after': {
              content: '""',
            },
            pre: {
              backgroundColor: '#111827',
              color: '#f9fafb',
            },
            'pre code': {
              backgroundColor: 'transparent',
              color: '#f9fafb',
              padding: '0',
            },
            blockquote: {
              borderLeftColor: '#93c5fd',
              backgroundColor: '#eff6ff',
              padding: '1rem',
              borderRadius: '0.375rem',
            },
          },
        },
      },
    },
  },
  plugins: [
    // Add typography plugin if available
    // require('@tailwindcss/typography'),
  ],
}