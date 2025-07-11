/** @type {import('tailwindcss').Config} */
module.exports = {
  // Thêm dòng này vào để kích hoạt dark mode bằng class
  darkMode: 'class',

  content: ["./public/index.html", "./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {},
  },
  plugins: [],
};