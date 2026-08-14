/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{vue,js}'],
  theme: {
    extend: {
      colors: {
        primary: {
          50:  '#eef3f8',
          100: '#d5e2ed',
          200: '#a8c2da',
          300: '#74a0c3',
          400: '#4a7fad',
          500: '#2a5f96',
          600: '#1a4571',
          700: '#112b4e',
          800: '#0d2140',
          900: '#081630',
        }
      }
    }
  },
  plugins: []
}
