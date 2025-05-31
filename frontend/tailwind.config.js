/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#3b82f6',
          light: '#93c5fd',
          dark: '#1d4ed8',
        },
        secondary: {
          DEFAULT: '#8b5cf6',
          light: '#a78bfa',
          dark: '#6d28d9',
        },
        accent: {
          DEFAULT: '#f59e0b',
          light: '#fcd34d',
          dark: '#d97706',
        },
        success: '#10b981',
        warning: '#f59e0b',
        error: '#ef4444',
        surface: '#ffffff',
        background: '#f9fafb',
        text: {
          primary: '#111827',
          secondary: '#6b7280',
        },
        border: '#e5e7eb',
      },
      animation: {
        'spin-slow': 'spin 3s linear infinite',
      },
      backdropBlur: {
        xs: '2px',
      },
    },
  },
  plugins: [],
}