// tailwind.config.js

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@heroui/theme/dist/components/(button|card|divider|input|modal|skeleton|spinner|tabs|toast|ripple|form).js",
  ],
  theme: {
    extend: {},
  },
  darkMode: "class",
  plugins: [],
};
