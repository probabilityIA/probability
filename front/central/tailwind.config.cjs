/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    './app/**/*.{js,ts,jsx,tsx}',
    './src/**/*.{js,ts,jsx,tsx}',
    './components/**/*.{js,ts,jsx,tsx}',
    './pages/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        navy: {
          50: '#f0f4ff',
          100: '#e6eeff',
          200: '#cfe0ff',
          300: '#a6c1ff',
          400: '#7ea1ff',
          500: '#5f7df0',
          600: '#3f5be0',
          700: '#1e3a8a',
          800: '#172554',
          900: '#0b1530',
        },
        brand: {
          50: '#f5f3ff',
          100: '#ede9fe',
          200: '#ddd6fe',
          300: '#c4b5fd',
          400: '#a78bfa',
          500: '#7c3aed',
          600: '#6d28d9',
          700: '#5b21b6',
          800: '#4c1d95',
          900: '#3b1464',
        },
        lightPrimary: '#f6f8ff',
      },
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', 'Segoe UI', 'Roboto', 'Helvetica', 'Arial'],
      },
      borderRadius: {
        xl: '1rem',
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
};
