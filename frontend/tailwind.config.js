/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js}'],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0fdf4',
          100: '#dcfce7',
          500: '#142c4b',
          600: '#112b4e',
          700: '#112b4e',
          800: '#112b4e',
          900: '#112b4e',
        }
      }
    }
  },
  plugins: []
}
