/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{html,js}",
    "./views/*.html",
  ],
  theme: {
    container: {
      center: true,
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

