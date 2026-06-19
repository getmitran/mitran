import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  darkMode: ['selector', '[data-theme="dark"]'],
  theme: {
    extend: {
      colors: {
        accent: 'var(--accent)',
        surface: 'var(--bg)',
        base: 'var(--bg)',
        card: 'var(--card)',
        'card-hover': 'var(--bg-hover)',
        border: 'var(--border)',
      },
      backgroundColor: {
        DEFAULT: 'var(--bg)',
      },
      textColor: {
        DEFAULT: 'var(--text)',
      },
      borderColor: {
        DEFAULT: 'var(--border)',
      },
    },
  },
  plugins: [],
} satisfies Config
