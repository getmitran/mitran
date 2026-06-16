import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        accent: '#06b6d4',
        surface: '#0f1629',
        base: '#0a0e1a',
        card: '#141b2d',
        'card-hover': '#1a2340',
      },
    },
  },
  plugins: [],
} satisfies Config
